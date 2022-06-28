/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package icon

import (
	"context"
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
	MonitorBlockMaxConcurrency = 300
)

type eventLogRawFilter struct {
	addr      []byte
	signature []byte
	next      []byte
	seq       uint64
}

type receiver struct {
	log      log.Logger
	src      chain.BTPAddress
	dst      chain.BTPAddress
	cl       *client
	blockReq BlockRequest
	sinkChan chan *chain.SubscribedEvent
	errChan  chan error
	par      *parser
	fd       *finder
}

func NewReceiver(src, dst chain.BTPAddress, urls []string, l log.Logger, addrToContractName map[string]chain.ContractName) (*receiver, error) {
	if len(urls) == 0 {
		return nil, errors.New("List of Urls is empty")
	}
	client, err := newClient(urls[0], l)
	if err != nil {
		return nil, err
	}

	dstAddr := dst.String()
	ef := &EventFilter{
		Addr:      Address(src.ContractAddress()),
		Signature: EventSignature,
		Indexed:   []*string{&dstAddr},
	}
	evtReq := BlockRequest{
		EventFilters: []*EventFilter{ef},
	} // fill height later

	recvr := &receiver{
		log:      l,
		src:      src,
		dst:      dst,
		cl:       client,
		blockReq: evtReq,
		sinkChan: make(chan *chain.SubscribedEvent),
		errChan:  make(chan error),
		fd:       NewFinder(l),
	}
	recvr.par, err = NewParser(addrToContractName)
	if err != nil {
		return nil, err
	}
	return recvr, nil
}

func (r *receiver) receiveLoop(ctx context.Context, startHeight uint64, callback func(ts []*TxnEventLog) error) (err error) {

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
									q.err = err
									return
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

func (r *receiver) Subscribe(ctx context.Context, height uint64) (err error) {

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
				el := eventLogInfo{contractAddress: string(txnLog.Addr), eventType: evtType, eventLog: res}
				if r.fd.Match(el) {
					r.log.Infof("Matched %+v", el)
					r.sinkChan <- &chain.SubscribedEvent{Res: []*TxnEventLog{txnLog}, ChainName: chain.ICON}
				}

			}
			return nil
		})
		if err != nil {
			r.log.Errorf("receiveLoop terminated: %v", err)
			r.errChan <- err
		}
	}()
	return nil
}
