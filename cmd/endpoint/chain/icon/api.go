package icon

import (
	"context"
	"math/big"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature             = "Message(str,int,bytes)"
	MonitorBlockMaxConcurrency = 10
)

type api struct {
	log       log.Logger
	src       chain.BTPAddress
	dst       chain.BTPAddress
	cl        *client
	networkID string
	blockReq  BlockRequest
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	par       *parser
	fd        *finder
	requester *requestAPI
}

func NewApi(l log.Logger, cfg *chain.ChainConfig) (chain.ChainAPI, error) {
	if len(cfg.URL) == 0 {
		return nil, errors.New("List of Urls is empty")
	}
	client, err := newClient(cfg.URL, l)
	if err != nil {
		return nil, err
	}

	dstAddr := cfg.Dst.String()
	ef := &EventFilter{
		Addr:      Address(cfg.Src.ContractAddress()),
		Signature: EventSignature,
		Indexed:   []*string{&dstAddr},
	}
	evtReq := BlockRequest{
		EventFilters: []*EventFilter{ef},
	} // fill height later

	recvr := &api{
		log:       l,
		src:       cfg.Src,
		dst:       cfg.Dst,
		cl:        client,
		blockReq:  evtReq,
		sinkChan:  make(chan *chain.EventLogInfo),
		errChan:   make(chan error),
		fd:        NewFinder(l),
		networkID: cfg.NetworkID,
	}
	recvr.par, err = NewParser(cfg.ConftractAddresses)
	if err != nil {
		return nil, err
	}
	recvr.requester, err = newRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
	return recvr, nil
}

func (r *api) receiveLoop(ctx context.Context, startHeight uint64, callback func(ts []*TxnEventLog) error) (err error) {

	blockReq := r.blockReq // copy

	blockReq.Height = NewHexInt(int64(startHeight))

	type res struct {
		Height  int64
		Hash    common.HexBytes
		TxnLogs []*TxnEventLog
	}
	ech := make(chan error)                                           // error channel
	rech := make(chan struct{}, 1)                                    // reconnect channel
	bnch := make(chan *BlockNotification, MonitorBlockMaxConcurrency) // block notification channel
	brch := make(chan *res, cap(bnch))                                // block result channel

	reconnect := func() {
		select {
		case rech <- struct{}{}:
		default:
		}
		for len(brch) > 0 || len(bnch) > 0 {
			select {
			case <-brch: // clear block result channel
			case <-bnch: // clear block notification channel
			}
		}
	}

	next := int64(startHeight) // next block height to process

	// subscribe to monitor block
	ctxMonitorBlock, cancelMonitorBlock := context.WithCancel(ctx)
	reconnect()

loop:
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-ech:
			return err

		case <-rech:
			cancelMonitorBlock()
			ctxMonitorBlock, cancelMonitorBlock = context.WithCancel(ctx)

			// start new monitor loop
			go func(ctx context.Context, cancel context.CancelFunc) {
				defer cancel()
				blockReq.Height = NewHexInt(next)
				err := r.cl.MonitorBlock(ctx, &blockReq,
					func(conn *websocket.Conn, v *BlockNotification) error {
						if !errors.Is(ctx.Err(), context.Canceled) {
							bnch <- v
						}
						return nil
					},
					func(conn *websocket.Conn) {},
					func(c *websocket.Conn, err error) {})
				if err != nil && !errors.Is(err, context.Canceled) {
					ech <- err
				}
			}(ctxMonitorBlock, cancelMonitorBlock)

		case br := <-brch:
			for ; br != nil; next++ {
				if br.Height%100 == 0 {
					r.log.WithFields(log.Fields{"height": br.Height}).Debug("block notification")
				}
				if err := callback(br.TxnLogs); err != nil {
					return errors.Wrapf(err, "receiveLoop: callback: %v", err)
				}
				if br = nil; len(brch) > 0 {
					br = <-brch
				}
			}
		default:
			select {
			default:
			case bn := <-bnch:

				type req struct {
					height int64
					hash   HexBytes
					retry  int
					err    error
					res    *res
				}

				qch := make(chan *req, cap(bnch))
				for i := int64(0); bn != nil; i++ {
					height, err := bn.Height.Value()
					if err != nil {
						panic(err)
					} else if height != next+i {
						r.log.WithFields(log.Fields{
							"height": log.Fields{"got": height, "expected": next + i},
						}).Error("reconnect: missing block notification")
						reconnect()
						continue loop
					}
					qch <- &req{
						height: height,
						retry:  3,
					} // fill qch with requests
					if bn = nil; len(bnch) > 0 && len(qch) < cap(qch) {
						bn = <-bnch
					}
				}

				brs := make([]*res, 0, len(qch))
				for q := range qch {
					switch {
					case q.err != nil:
						if q.retry > 0 {
							q.retry--
							q.res, q.err = nil, nil
							qch <- q
							continue
						}
						r.log.WithFields(log.Fields{"height": q.height, "error": q.err}).Debug("receiveLoop: req error")
						brs = append(brs, nil)
						if len(brs) == cap(brs) {
							close(qch)
						}

					case q.res != nil:
						brs = append(brs, q.res)
						if len(brs) == cap(brs) {
							close(qch)
						}

					default:
						go func(q *req) {
							defer func() {
								time.Sleep(500 * time.Millisecond)
								qch <- q
							}()
							if q.res == nil {
								q.res = &res{}
							}
							q.res.Height = q.height
							q.res.Hash, q.err = q.hash.Value()
							if q.err != nil {
								q.err = errors.Wrapf(q.err,
									"invalid hash: height=%v, hash=%v, %v", q.height, q.hash, q.err)
								return
							}

							blk, err := r.cl.GetBlockByHeight(&BlockHeightParam{Height: NewHexInt(q.height)})
							if err != nil {
								q.err = errors.Wrapf(err, "GetBlockByHeight %v", q.height)
								return
							}
							q.res.TxnLogs = []*TxnEventLog{}
							for _, txn := range blk.NormalTransactions {
								res, err := r.cl.GetTransactionResult(&TransactionHashParam{Hash: txn.TxHash})
								if err != nil {
									switch re := err.(type) {
									case *jsonrpc.Error:
										switch re.Code {
										case JsonrpcErrorCodePending, JsonrpcErrorCodeExecuting:
											time.Sleep(2 * time.Second)
											res, err = r.cl.GetTransactionResult(&TransactionHashParam{Hash: txn.TxHash})
										}
									}
									if err != nil {
										q.err = err
										return
									}
								}
								if len(res.EventLogs) > 0 {
									for i := 0; i < len(res.EventLogs); i++ {
										q.res.TxnLogs = append(q.res.TxnLogs, &res.EventLogs[i])
									}
								}

							}
						}(q)
					}
				}
				// filter nil
				_brs, brs := brs, brs[:0]
				for _, v := range _brs {
					if v != nil {
						brs = append(brs, v)
					}
				}
				// sort and forward notifications
				if len(brs) > 0 {
					sort.SliceStable(brs, func(i, j int) bool {
						return brs[i].Height < brs[j].Height
					})
					for i, d := range brs {
						if d.Height == int64(next)+int64(i) {
							brch <- d
						}
					}
				}
			}
		}
	}

}

func (r *api) Subscribe(ctx context.Context, height uint64) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {

	// _errCh := make(chan error)
	go func() {
		// defer close(_errCh)
		err := r.receiveLoop(ctx, height, func(txnLogs []*TxnEventLog) error {
			for _, txnLog := range txnLogs {
				res, evtType, err := r.par.Parse(txnLog)
				if err != nil {
					//r.log.Error(err)
					continue
				}
				el := &chain.EventLogInfo{ContractAddress: string(txnLog.Addr), EventType: evtType, EventLog: res}

				if r.fd.Match(el) { //el.IDs is updated by match if matched
					//r.log.Infof("Matched %+v", el)
					r.sinkChan <- el
				}

			}
			return nil
		})
		if err != nil {
			r.log.Errorf("receiveLoop terminated: %v", err)
			r.errChan <- err
		}
	}()
	return r.sinkChan, r.errChan, nil
}

func (r *api) Transfer(param *chain.RequestParam) (txnHash string, err error) {
	if param.FromChain != chain.ICON {
		err = errors.New("Source Chan should be Icon")
		return
	}
	if param.ToChain == chain.ICON {
		if param.Token == chain.ICXToken {
			txnHash, _, err = r.requester.transferICX(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.IRC2Token {
			txnHash, _, err = r.requester.transferIrc2(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else if param.ToChain == chain.HMNY {
		if param.Token == chain.ICXToken {
			txnHash, _, err = r.requester.TransferICXToHarmony(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.IRC2Token {
			txnHash, _, err = r.requester.transferIrc2ToHmny(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.ONEToken {
			txnHash, _, err = r.requester.transferWrappedOneFromIconToHmny(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else {
		err = errors.New("Unsupport Transaction Parameters ")
	}
	return
}

func (r *api) GetCoinBalance(addr string, coinType chain.TokenType) (*big.Int, error) {
	if coinType == chain.ICXToken {
		return r.requester.getICXBalance(addr)
	} else if coinType == chain.IRC2Token {
		return r.requester.getIrc2Balance(addr)
	} else if coinType == chain.ONEToken {
		return r.requester.getIconWrappedOne(addr)
	}
	return nil, errors.New("Unsupported Token Type ")
}

func (r *api) WaitForTxnResult(hash string) (interface{}, []*chain.EventLogInfo, error) {
	_, txRes, err := r.cl.waitForResults(context.TODO(), &TransactionHashParam{Hash: HexBytes(hash)})
	if err != nil {
		return nil, nil, err
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.EventLogs {
		decodedLog, eventType, err := r.par.Parse(&v)
		if err != nil {
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: string(v.Addr), EventType: eventType, EventLog: decodedLog})
	}
	return txRes, plogs, nil
}

func (r *api) Approve(ownerKey string, amount big.Int) (txnHash string, err error) {
	txnHash, _, _, err = r.requester.approveIconNativeCoinBSHToAccessHmnyOne(ownerKey, amount)
	return
}

func (r *api) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + r.networkID + ".icon/" + addr
	return &fullAddr
}

func (r *api) GetKeyPairs(num int) ([][2]string, error) {
	var err error
	res := make([][2]string, num)
	for i := 0; i < num; i++ {
		res[i], err = generateKeyPair()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *api) WatchFor(id uint64, eventType chain.EventLogType, seq int64, contractAddress string) error {
	return r.fd.WatchFor(args{eventType: eventType, seq: seq, contractAddress: contractAddress, id: id})
}
