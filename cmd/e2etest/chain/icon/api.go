package icon

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature             = "Message(str,int,bytes)"
	MonitorBlockMaxConcurrency = 50
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
		return nil, errors.New("List of Urls is empty ")
	}
	client, err := newClient(cfg.URL, l)
	if err != nil {
		return nil, errors.Wrap(err, "newClient ")
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
		fd:        NewFinder(l, cfg.ConftractAddresses),
		networkID: cfg.NetworkID,
	}
	recvr.par, err = NewParser(cfg.ConftractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "NewParser ")
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
			r.log.Warn("Context Cancelled Exiting Icon Subscription Loop")
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
										q.err = errors.Wrapf(err, "GetTransactionResult(%v)", blk.Height)
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

func (r *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	blk, err := r.cl.GetLastBlock()
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetLastBlock ")
	}
	height := uint64(blk.Height)
	r.log.Infof("Subscribe Start Height %v", height)
	// _errCh := make(chan error)
	go func() {
		// defer close(_errCh)
		err := r.receiveLoop(ctx, height, func(txnLogs []*TxnEventLog) error {
			for _, txnLog := range txnLogs {
				res, evtType, err := r.par.Parse(txnLog)
				if err != nil {
					r.log.Trace(errors.Wrap(err, "Parse "))
					err = nil
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

func (r *api) Transfer(coinName, senderKey, recepientAddress string, amount big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp://") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".icon") {
		within = true
		splts := strings.Split(recepientAddress, "/")
		recepientAddress = splts[len(splts)-1]
	}
	if within {
		if coinName == "ICX" {
			txnHash, _, err = r.requester.transferNativeIntraChain(senderKey, recepientAddress, amount)
		} else if coinName == "TICX" {
			txnHash, _, err = r.requester.transferTokenIntraChain(senderKey, recepientAddress, amount)
		} else {
			err = fmt.Errorf("IntraChain transfers are supported for coins ICX and TICX only")
		}
	} else {
		if coinName == "ICX" {
			txnHash, _, err = r.requester.transferNativeCrossChain(senderKey, recepientAddress, amount)
		} else { // ONE, TONE, TICX
			txnHash, _, err = r.requester.transferWrappedCrossChain(coinName, senderKey, recepientAddress, amount)
		}
	}
	return
}

func (r *api) GetCoinBalance(coinName string, addr string) (*big.Int, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".icon") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	if coinName == "ICX" {
		return r.requester.getICXBalance(address)
	}
	return r.requester.getWrappedCoinBalance(coinName, address)
}

func (r *api) WaitForTxnResult(ctx context.Context, hash string) (interface{}, []*chain.EventLogInfo, error) {
	_, txRes, err := r.cl.waitForResults(ctx, &TransactionHashParam{Hash: HexBytes(hash)})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "waitForResults(%v)", hash)
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.EventLogs {
		decodedLog, eventType, err := r.par.Parse(&v)
		if err != nil {
			r.log.Trace(errors.Wrap(err, "waitForResults.Parse "))
			err = nil
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: string(v.Addr), EventType: eventType, EventLog: decodedLog})
	}
	return txRes, plogs, nil
}

func (r *api) Approve(coinName string, ownerKey string, amount big.Int) (txnHash string, err error) {
	if coinName == "ONE" || coinName == "TONE" {
		txnHash, _, err = r.requester.approveCrossNativeCoin(coinName, ownerKey, amount)
	} else if coinName == "TICX" {
		txnHash, _, err = r.requester.approveToken(coinName, ownerKey, amount)
	} else {
		err = errors.Wrapf(err, "CoinName not among accepted Values ONE, ETH. Got %v", coinName)
	}
	return
}

func (r *api) GetChainType() chain.ChainType {
	return chain.ICON
}

func (r *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + r.networkID + ".icon/" + addr
	return fullAddr
}

func (r *api) GetKeyPairs(num int) ([][2]string, error) {
	var err error
	res := make([][2]string, num)
	for i := 0; i < num; i++ {
		res[i], err = generateKeyPair()
		if err != nil {
			return nil, errors.Wrap(err, "generateKeyPair ")
		}
	}
	return res, nil
}

func (r *api) WatchForTransferStart(id uint64, coinName string, seq int64) error {
	return r.fd.watchFor(chain.TransferStart, id, coinName, seq)
}

func (r *api) WatchForTransferReceived(id uint64, coinName string, seq int64) error {
	return r.fd.watchFor(chain.TransferReceived, id, coinName, seq)
}

func (r *api) WatchForTransferEnd(id uint64, coinName string, seq int64) error {
	return r.fd.watchFor(chain.TransferEnd, id, coinName, seq)
}
