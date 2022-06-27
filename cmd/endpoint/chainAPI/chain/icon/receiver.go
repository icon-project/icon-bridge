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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/common/trie/ompt"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature      = "Message(str,int,bytes)"
	EventIndexSignature = 0
	EventIndexNext      = 1
	EventIndexSequence  = 2
)
const MAX_RETRY = 3
const RECONNECT_ON_UNEXPECTED_HEIGHT = "Unexpected Block Height. Should Reconnect"
const (
	SyncVerifierMaxConcurrency = 300 // 150
	MonitorBlockMaxConcurrency = 300
)

type receiverOptions struct {
	SyncConcurrency uint64           `json:"syncConcurrency"`
	Verifier        *VerifierOptions `json:"verifier"`
}

func (opts *receiverOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

type eventLogRawFilter struct {
	addr      []byte
	signature []byte
	next      []byte
	seq       uint64
}

type receiver struct {
	log       log.Logger
	src       chain.BTPAddress
	dst       chain.BTPAddress
	cl        *client
	opts      receiverOptions
	blockReq  BlockRequest
	logFilter eventLogRawFilter
	sinkChan  chan *chain.SubscribedEvent
	errChan   chan error
}

func NewReceiver(src, dst chain.BTPAddress, urls []string, opts map[string]interface{}, l log.Logger) (chain.SubscriptionAPI, error) {
	if len(urls) == 0 {
		return nil, errors.New("List of Urls is empty")
	}
	client, err := newClient(urls[0], l)
	if err != nil {
		return nil, err
	}
	var recvOpts receiverOptions
	if err := recvOpts.Unmarshal(opts); err != nil {
		return nil, errors.Wrapf(err, "recvOpts.Unmarshal: %v", err)
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

	efAddr, err := ef.Addr.Value()
	if err != nil {
		return nil, errors.Wrapf(err, "ef.Addr.Value: %v", err)
	}

	recvr := &receiver{
		log:      l,
		src:      src,
		dst:      dst,
		cl:       client,
		opts:     recvOpts,
		blockReq: evtReq,
		logFilter: eventLogRawFilter{
			addr:      efAddr,
			signature: []byte(EventSignature),
			next:      []byte(dstAddr),
		}, // fill seq later
		sinkChan: make(chan *chain.SubscribedEvent),
		errChan:  make(chan error),
	}

	return recvr, nil
}

func (r *receiver) newVerifer(opts *VerifierOptions) (*Verifier, error) {
	validators, err := r.cl.getValidatorsByHash(opts.ValidatorsHash)
	if err != nil {
		return nil, err
	}
	vr := Verifier{
		next:               int64(opts.BlockHeight),
		nextValidatorsHash: opts.ValidatorsHash,
		validators: map[string][]common.HexBytes{
			opts.ValidatorsHash.String(): validators,
		},
	}
	header, err := r.cl.getBlockHeaderByHeight(int64(vr.next))
	if err != nil {
		return nil, err
	}
	votes, err := r.cl.GetVotesByHeight(
		&BlockHeightParam{Height: NewHexInt(vr.next)})
	if err != nil {
		return nil, err
	}
	ok, err := vr.Verify(header, votes)
	if !ok {
		err = errors.New("verification failed")
	}
	if err != nil {
		return nil, err
	}
	return &vr, nil
}

func (r *receiver) syncVerifier(vr *Verifier, height int64) error {
	if height == vr.Next() {
		return nil
	}
	if vr.Next() > height {
		return fmt.Errorf(
			"invalid target height: verifier height (%d) > target height (%d)",
			vr.Next(), height)
	}

	type res struct {
		Height         int64
		Header         *BlockHeader
		Votes          []byte
		NextValidators []common.HexBytes
	}

	type req struct {
		height int64
		err    error
		res    *res
		retry  int64
	}

	r.log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Debug("syncVerifier: start")

	for vr.Next() < height {
		rqch := make(chan *req, SyncVerifierMaxConcurrency)
		for i := vr.Next(); len(rqch) < cap(rqch); i++ {
			rqch <- &req{height: i}
		}
		sres := make([]*res, 0, len(rqch))
		for q := range rqch {
			switch {
			case q.err != nil:
				if q.retry > 0 {
					q.retry--
					q.res, q.err = nil, nil
					rqch <- q
					continue
				}
				r.log.WithFields(log.Fields{
					"height": q.height, "error": q.err.Error()}).Debug("syncVerifier: req error")
				sres = append(sres, nil)
				if len(sres) == cap(sres) {
					close(rqch)
				}
			case q.res != nil:
				sres = append(sres, q.res)
				if len(sres) == cap(sres) {
					close(rqch)
				}
			default:
				go func(q *req) {
					defer func() {
						time.Sleep(500 * time.Millisecond)
						rqch <- q
					}()
					if q.res == nil {
						q.res = &res{}
					}
					q.res.Height = q.height
					q.res.Header, q.err = r.cl.getBlockHeaderByHeight(q.height)
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "syncVerifier: getBlockHeader: %v", q.err)
						return
					}
					q.res.Votes, q.err = r.cl.GetVotesByHeight(
						&BlockHeightParam{Height: NewHexInt(int64(q.height))})
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "syncVerifier: GetVotesByHeight: %v", q.err)
						return
					}
					if len(vr.Validators(q.res.Header.NextValidatorsHash)) == 0 {
						q.res.NextValidators, q.err = r.cl.getValidatorsByHash(q.res.Header.NextValidatorsHash)
						if q.err != nil {
							q.err = errors.Wrapf(q.err, "syncVerifier: getValidatorsByHash: %v", q.err)
							return
						}
					}
				}(q)
			}
		}
		// filter nil
		_sres, sres := sres, sres[:0]
		for _, v := range _sres {
			if v != nil {
				sres = append(sres, v)
			}
		}
		// sort and forward notifications
		if len(sres) > 0 {
			sort.SliceStable(sres, func(i, j int) bool {
				return sres[i].Height < sres[j].Height
			})
			for _, r := range sres {
				if vr.Next() == r.Height {
					err := vr.Update(r.Header, r.NextValidators)
					if err != nil {
						return errors.Wrapf(err, "syncVerifier: Update: %v", err)
					}
				}
			}
			//r.log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Debug("syncVerifier: syncing")
		}
	}

	r.log.WithFields(log.Fields{"height": vr.Next()}).Debug("syncVerifier: complete")
	return nil
}

func (r *receiver) receiveLoop(ctx context.Context, startHeight, startSeq uint64, callback func(rs []*chain.Receipt, ts []*TxnEventLog) error) (err error) {

	blockReq, logFilter := r.blockReq, r.logFilter // copy

	blockReq.Height, logFilter.seq = NewHexInt(int64(startHeight)), startSeq

	var vr *Verifier
	if r.opts.Verifier != nil {
		vr, err = r.newVerifer(r.opts.Verifier)
		if err != nil {
			return err
		}
	}

	type res struct {
		Height         int64
		Hash           common.HexBytes
		Header         *BlockHeader
		Votes          []byte
		NextValidators []common.HexBytes
		Receipts       []*chain.Receipt
		TxnLogs        []*TxnEventLog
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

			// sync verifier
			if err := r.syncVerifier(vr, next); err != nil {
				return errors.Wrapf(err, "sync verifier: %v", err)
			}

		case br := <-brch:
			for ; br != nil; next++ {
				if br.Height%100 == 0 {
					r.log.WithFields(log.Fields{"height": br.Height}).Debug("icon block notification")
				}
				if vr != nil {
					ok, err := vr.Verify(br.Header, br.Votes)
					if !ok || err != nil {
						if err != nil {
							r.log.WithFields(log.Fields{"height": br.Height, "error": err}).Debug("receiveLoop: verification error")
						} else if !ok {
							r.log.WithFields(log.Fields{"height": br.Height, "hash": br.Hash}).Debug("receiveLoop: invalid header")
						}
						reconnect() // reconnect websocket
						r.log.WithFields(log.Fields{"height": br.Height, "hash": br.Hash}).Error("reconnect: verification failed")
						break
					}
					if err := vr.Update(br.Header, br.NextValidators); err != nil {
						return errors.Wrapf(err, "receiveLoop: update verifier: %v", err)
					}
				}
				if err := callback(br.Receipts, br.TxnLogs); err != nil {
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
					height  int64
					hash    HexBytes
					indexes [][]HexInt
					events  [][][]HexInt

					retry int

					err error
					res *res
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
						height:  height,
						hash:    bn.Hash,
						indexes: bn.Indexes,
						events:  bn.Events,
						retry:   3,
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

							q.res.Header, q.err = r.cl.getBlockHeaderByHeight(q.height)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "getBlockHeader: %v", q.err)
								return
							}
							// fetch votes, next validators only if verifier exists
							if vr != nil {
								q.res.Votes, q.err = r.cl.GetVotesByHeight(
									&BlockHeightParam{Height: NewHexInt(int64(q.height))})
								if q.err != nil {
									q.err = errors.Wrapf(q.err, "GetVotesByHeight: %v", q.err)
									return
								}
								if len(vr.Validators(q.res.Header.NextValidatorsHash)) == 0 {
									q.res.NextValidators, q.err = r.cl.getValidatorsByHash(q.res.Header.NextValidatorsHash)
									if q.err != nil {
										q.err = errors.Wrapf(q.err, "getValidatorsByHash: %v", q.err)
										return
									}
								}
							}
							if len(q.res.Header.NormalTransactionsHash) > 0 {
								blk, err := r.cl.GetBlockByHeight(&BlockHeightParam{Height: NewHexInt(q.height)})
								if err != nil {
									q.err = errors.Wrapf(err, "GetBlockByHeight %v", q.height)
									return
								}
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
										// &TxnLog{From: txn.From, To: txn.To,
										// 	TxHash: txn.TxHash, BlockHeight: q.height,
										// 	EventLogs: res.EventLogs, Status: res.Status})
									}

								}
							}

							if len(q.indexes) > 0 && len(q.events) > 0 {
								var hr BlockHeaderResult
								_, err := codec.RLP.UnmarshalFromBytes(q.res.Header.Result, &hr)
								if q.err != nil {
									q.err = errors.Wrapf(q.err, "BlockHeaderResult.UnmarshalFromBytes: %v", err)
									return
								}
								for i, index := range q.indexes[0] {
									p := &ProofEventsParam{
										Index:     index,
										BlockHash: q.hash,
										Events:    q.events[0][i],
									}
									proofs, err := r.cl.GetProofForEvents(p)
									if err != nil {
										q.err = errors.Wrapf(err, "GetProofForEvents: %v", err)
										return
									}
									if len(proofs) != 1+len(p.Events) { // num_receipt + num_events
										q.err = errors.Wrapf(q.err,
											"Proof does not include all events: len(proofs)=%d, expected=%d",
											len(proofs), len(p.Events)+1,
										)
										return
									}

									// Processing receipt index
									serializedReceipt, err := mptProve(index, proofs[0], hr.ReceiptHash)
									if err != nil {
										q.err = errors.Wrapf(err, "MPTProve Receipt: %v", err)
										return
									}
									var result TxResult
									_, err = codec.RLP.UnmarshalFromBytes(serializedReceipt, &result)
									if err != nil {
										q.err = errors.Wrapf(err, "Unmarshal Receipt: %v", err)
										return
									}

									idx, _ := index.Value()
									receipt := &chain.Receipt{
										Index:  uint64(idx),
										Height: uint64(q.height),
									}
									for j := 0; j < len(p.Events); j++ {
										// nextEP is pointer to event where sequence has caught up
										serializedEventLog, err := mptProve(
											p.Events[j], proofs[j+1], common.HexBytes(result.EventLogsHash))
										if err != nil {
											q.err = errors.Wrapf(err, "event.MPTProve: %v", err)
											return
										}
										var el EventLog
										_, err = codec.RLP.UnmarshalFromBytes(serializedEventLog, &el)
										if err != nil {
											q.err = errors.Wrapf(err, "event.UnmarshalFromBytes: %v", err)
											return
										}

										if bytes.Equal(el.Addr, logFilter.addr) &&
											bytes.Equal(el.Indexed[EventIndexSignature], logFilter.signature) &&
											bytes.Equal(el.Indexed[EventIndexNext], logFilter.next) {
											var seqGot common.HexInt
											seqGot.SetBytes(el.Indexed[EventIndexSequence])

											receipt.Events = append(receipt.Events,
												&chain.Event{
													Next:     chain.BTPAddress(el.Indexed[EventIndexNext]),
													Sequence: seqGot.Uint64(),
													Message:  el.Data[0],
												})
										} else {
											if !bytes.Equal(el.Addr, logFilter.addr) {
												r.log.WithFields(log.Fields{
													"height":   q.height,
													"got":      common.HexBytes(el.Addr),
													"expected": common.HexBytes(logFilter.addr)}).Error("invalid event: cannot match addr")
											}
											if !bytes.Equal(el.Indexed[EventIndexSignature], logFilter.signature) {
												r.log.WithFields(log.Fields{
													"height":   q.height,
													"got":      common.HexBytes(el.Indexed[EventIndexSignature]),
													"expected": common.HexBytes(logFilter.signature)}).Error("invalid event: cannot match sig")
											}
											if !bytes.Equal(el.Indexed[EventIndexNext], logFilter.next) {
												r.log.WithFields(log.Fields{
													"height":   q.height,
													"got":      common.HexBytes(el.Indexed[EventIndexNext]),
													"expected": common.HexBytes(logFilter.next)}).Error("invalid event: cannot match next")
											}
											q.err = errors.New("invalid event")
											return
										}
									}
									if len(receipt.Events) > 0 {
										if len(receipt.Events) == len(p.Events) {
											q.res.Receipts = append(q.res.Receipts, receipt)
										} else {
											r.log.WithFields(log.Fields{
												"height":              q.height,
												"receipt_index":       index,
												"got_num_events":      len(receipt.Events),
												"expected_num_events": len(p.Events)}).Info("failed to verify all events for the receipt")
											q.err = errors.New("failed to verify all events for the receipt")
											return
										}
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

func (r *receiver) Subscribe(
	ctx context.Context,
	opts chain.SubscribeOptions) (err error) {

	opts.Seq++

	if opts.Height < 1 {
		opts.Height = 1
	}

	// _errCh := make(chan error)
	go func() {
		// defer close(_errCh)
		err := r.receiveLoop(ctx, opts.Height, opts.Seq, func(receipts []*chain.Receipt, txnLogs []*TxnEventLog) error {
			for _, receipt := range receipts {
				events := receipt.Events[:0]
				for _, event := range receipt.Events {
					switch {
					case event.Sequence == opts.Seq:
						events = append(events, event)
						opts.Seq++
					case event.Sequence > opts.Seq:
						r.log.WithFields(log.Fields{
							"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
						}).Error("invalid event seq")
						return fmt.Errorf("invalid event seq")
					default:
						r.log.WithFields(log.Fields{
							"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
						}).Warn("Searching sequence...")

					}
				}
				receipt.Events = events
			}
			if len(txnLogs) > 0 {
				r.sinkChan <- &chain.SubscribedEvent{Res: txnLogs, ChainName: chain.ICON}
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

// func decodeTxnLogToSubscribedEvent(txn *TxnLog) *chain.SubscribedEvent {
// 	newEventLog := make([]chain.EventLog, len(txn.EventLogs))
// 	for li, l := range txn.EventLogs {
// 		newEventLog[li] = chain.EventLog{Addr: string(l.Addr), Indexed: l.Indexed, Data: l.Data}
// 	}
// 	status := int64(1)
// 	if s, err := txn.Status.Value(); err != nil {
// 		status = s
// 	}
// 	return &chain.SubscribedEvent{
// 		TxHash:      string(txn.TxHash),
// 		From:        string(txn.From),
// 		To:          string(txn.To),
// 		EventLogs:   newEventLog,
// 		Status:      status,
// 		BlockHeight: txn.BlockHeight,
// 	}
// }

func mptProve(key HexInt, proofs [][]byte, hash []byte) ([]byte, error) {
	db := db.NewMapDB()
	defer db.Close()
	index, err := key.Value()
	if err != nil {
		return nil, err
	}
	indexKey, err := vlcodec.RLP.MarshalToBytes(index)
	if err != nil {
		return nil, err
	}
	mpt := ompt.NewMPTForBytes(db, hash)
	trie, err1 := mpt.Prove(indexKey, proofs)
	if err1 != nil {
		return nil, err1

	}
	return trie, nil
}
