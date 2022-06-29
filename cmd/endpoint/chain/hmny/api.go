package hmny

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	BlockInterval              = 2 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	defaultReadTimeout         = 15 * time.Second
	monitorBlockMaxConcurrency = 50 // number of concurrent requests to synchronize older blocks from source chain
)

func NewApi(l log.Logger, cfg *chain.ChainConfig) (chain.ChainAPI, error) {
	if len(cfg.URL) == 0 {
		return nil, errors.New("empty urls")
	}
	var err error
	r := &api{
		log:       l,
		src:       cfg.Src,
		dst:       cfg.Dst,
		networkID: cfg.NetworkID,
		fd:        NewFinder(l),
		sinkChan:  make(chan *chain.EventLogInfo),
		errChan:   make(chan error),
	}

	r.cls, err = newClients([]string{cfg.URL}, cfg.Src.ContractAddress(), r.log)
	if err != nil {
		return nil, err
	}
	r.par, err = NewParser(cfg.URL, cfg.ConftractAddresses)
	if err != nil {
		return nil, err
	}
	r.requester, err = newRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type api struct {
	log       log.Logger
	src       chain.BTPAddress
	dst       chain.BTPAddress
	cls       []*client
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	par       *parser
	fd        *finder
	requester *requestAPI
	networkID string
}

func (r *api) client() *client {
	return r.cls[0]
}

// Options for a new block notifications channel

func (r *api) receiveLoop(ctx context.Context, opts *bnOptions, callback func(v *BlockNotification) error) error {

	if opts == nil {
		return errors.New("receiveLoop: invalid options: <nil>")
	}

	// block notification channel
	// (buffered: to avoid deadlock)
	// increase concurrency parameter for faster sync
	bnch := make(chan *BlockNotification, opts.Concurrency)

	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	latestHeight := func() uint64 {
		height, err := r.client().GetBlockNumber()
		if err != nil {
			r.log.WithFields(log.Fields{"error": err}).Error("receiveLoop: failed to GetBlockNumber")
			return 0
		}
		return height
	}

	next, latest := opts.StartHeight, latestHeight()

	// last unverified block notification
	var lbn *BlockNotification

	// start monitor loop
	for {
		select {
		case <-ctx.Done():
			r.log.Warn("Context Cancelled Exiting Hmny Subscription Loop")
			return nil

		case <-heightTicker.C:
			latest++

		case <-heightPoller.C:
			if height := latestHeight(); height > latest {
				latest = height
				if next > latest {
					r.log.Debugf("receiveLoop: skipping; latest=%d, next=%d", latest, next)
				}
			}

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					if err := callback(lbn); err != nil {
						return errors.Wrapf(err, "receiveLoop: callback: %v", err)
					}
				}
				if lbn, bn = bn, nil; len(bnch) > 0 {
					bn = <-bnch
				}
			}
			// remove unprocessed notifications
			for len(bnch) > 0 {
				<-bnch
			}

		default:
			if next >= latest {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			type bnq struct {
				h     uint64
				v     *BlockNotification
				err   error
				retry int
			}

			qch := make(chan *bnq, cap(bnch))
			for i := next; i < latest &&
				len(qch) < cap(qch); i++ {
				qch <- &bnq{i, nil, nil, 3} // fill bch with requests
			}
			bns := make([]*BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil:
					if q.retry > 0 {
						if !strings.HasSuffix(q.err.Error(), "requested block number greater than current block number") {
							q.retry--
							q.v, q.err = nil, nil
							qch <- q
							continue
						}
						if latest >= q.h {
							latest = q.h - 1
						}
					}
					r.log.Errorf("receiveLoop: bnq: h=%d:%v, %v", q.h, q.v.Header.Hash(), q.err)
					bns = append(bns, nil)
					if len(bns) == cap(bns) {
						close(qch)
					}

				case q.v != nil:
					bns = append(bns, q.v)
					if len(bns) == cap(bns) {
						close(qch)
					}
				default:
					go func(q *bnq) {
						defer func() {
							time.Sleep(500 * time.Millisecond)
							qch <- q
						}()
						if q.v == nil {
							q.v = &BlockNotification{}
						}
						q.v.Height = (&big.Int{}).SetUint64(q.h)
						q.v.Header, q.err = r.client().GetHmyV2HeaderByHeight(q.v.Height)
						if q.err != nil {
							q.err = errors.Wrapf(q.err, "GetHmyHeaderByHeight: %v", q.err)
							return
						}
						q.v.Hash = q.v.Header.Hash()
						if q.v.Header.GasUsed > 0 {
							q.v.Receipts, q.err = r.client().GetBlockReceipts(q.v.Hash)
							if q.err == nil {
								receiptsRoot := types.DeriveSha(q.v.Receipts)
								if !bytes.Equal(receiptsRoot.Bytes(), q.v.Header.ReceiptsRoot.Bytes()) {
									q.err = fmt.Errorf(
										"invalid receipts: remote=%v, local=%v",
										q.v.Header.ReceiptsRoot, receiptsRoot)
								}
							}
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "GetBlockReceipts: %v", q.err)
								return
							}
						}
					}(q)
				}
			}
			// filter nil
			_bns_, bns := bns, bns[:0]
			for _, v := range _bns_ {
				if v != nil {
					bns = append(bns, v)
				}
			}
			// sort and forward notifications
			if len(bns) > 0 {
				sort.SliceStable(bns, func(i, j int) bool {
					return bns[i].Height.Uint64() < bns[j].Height.Uint64()
				})
				for i, v := range bns {
					if v.Height.Uint64() == next+uint64(i) {
						bnch <- v
					}
				}
			}
		}
	}
}

func (r *api) Subscribe(ctx context.Context, height uint64) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {

	go func() {
		lastHeight := height - 1
		if err := r.receiveLoop(ctx,
			&bnOptions{
				StartHeight: height,
				Concurrency: monitorBlockMaxConcurrency,
			},
			func(v *BlockNotification) error {
				if v.Height.Int64()%100 == 0 {
					r.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")
				}
				if v.Height.Uint64() != lastHeight+1 {
					r.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, v.Height.Uint64())
					return fmt.Errorf(
						"block notification: expected=%d, got=%d",
						lastHeight+1, v.Height.Uint64())
				}
				if len(v.Receipts) > 0 {
					for _, sev := range v.Receipts {
						for _, txnLog := range sev.Logs {
							res, evtType, err := r.par.Parse(txnLog)
							if err != nil {
								//r.log.Error(err)
								continue
							}
							el := &chain.EventLogInfo{ContractAddress: txnLog.Address.String(), EventType: evtType, EventLog: res}
							if r.fd.Match(el) {
								//r.log.Infof("Matched %+v", el)
								r.sinkChan <- el
							}
						}
					}
				}
				lastHeight++
				return nil
			}); err != nil {
			r.log.Errorf("receiveLoop terminated: %+v", err)
			r.errChan <- err
		}
	}()

	return r.sinkChan, r.errChan, nil
}

func (r *api) GetCoinBalance(addr string, coinType chain.TokenType) (*big.Int, error) {
	if coinType == chain.ONEToken {
		return r.requester.getHmnyBalance(addr)
	} else if coinType == chain.ERC20Token {
		return r.requester.getHmnyErc20Balance(addr)
	} else if coinType == chain.ICXToken {
		return r.requester.getHmnyWrappedICX(addr)
	}
	return nil, errors.New("Unsupported Token Type ")
}

func (r *api) Transfer(param *chain.RequestParam) (txnHash string, err error) {
	if param.FromChain != chain.HMNY {
		err = errors.New("Source Chan should be Hmny")
		return
	}
	if param.ToChain == chain.HMNY {
		if param.Token == chain.ONEToken {
			txnHash, _, err = r.requester.transferHmnyOne(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.ERC20Token {
			txnHash, _, err = r.requester.transferErc20(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else if param.ToChain == chain.ICON {
		if param.Token == chain.ONEToken {
			txnHash, _, err = r.requester.transferOneToIcon(param.SenderKey, param.ToAddress, param.Amount)
		} else if param.Token == chain.ERC20Token {
			_, _, txnHash, _, err = r.requester.transferERC20ToIcon(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.ICXToken {
			txnHash, _, err = r.requester.transferWrappedICXFromHmnyToIcon(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else {
		err = errors.New("Unsupport Transaction Parameters ")
	}
	return
}

func (r *api) Approve(ownerKey string, amount big.Int) (txnHash string, err error) {
	txnHash, _, _, err = r.requester.approveHmnyNativeBSHCoreToAccessICX(ownerKey, amount)
	return
}

// func (r *api) WaitForTxnResult(hash string) (txr interface{}, err error) {
// 	txr, err = r.requester.waitForResults(context.TODO(), )
// 	return
// }
func (r *api) WaitForTxnResult(ctx context.Context, hash string) (interface{}, []*chain.EventLogInfo, error) {
	txRes, err := r.requester.waitForResults(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, nil, err
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.Logs {
		decodedLog, eventType, err := r.par.ParseEth(v)
		if err != nil {
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: v.Address.String(), EventType: eventType, EventLog: decodedLog})
	}
	return txRes, plogs, nil
}

func (r *api) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + r.networkID + ".hmny/" + addr
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
	return r.fd.WatchFor(args{id: id, eventType: eventType, seq: seq, contractAddress: contractAddress})
}
