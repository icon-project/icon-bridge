package icon

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

//  const (
// 	 EventIndexSignature = 0
// 	 EventIndexNext      = 1
// 	 EventIndexSequence  = 2
//  )

//  const (
// 	 SyncVerifierMaxConcurrency = 300 // 150
// 	 MonitorBlockMaxConcurrency = 300
//  )

type ReceiverCore struct {
	Log      log.Logger
	Cl       *Client
	Opts     ReceiverOptions
	BlockReq BlockRequest
}

func (r *ReceiverCore) newVerifer(opts *VerifierOptions) (*Verifier, error) {
	validators, err := r.Cl.getValidatorsByHash(opts.ValidatorsHash)
	if err != nil {
		return nil, err
	}
	vr := Verifier{
		next:               int64(opts.BlockHeight),
		nextValidatorsHash: opts.ValidatorsHash,
		validators: map[string][]common.Address{
			opts.ValidatorsHash.String(): validators,
		},
	}
	header, err := r.Cl.getBlockHeaderByHeight(int64(vr.next))
	if err != nil {
		return nil, err
	}
	votes, err := r.Cl.GetVotesByHeight(
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

func (r *ReceiverCore) syncVerifier(vr *Verifier, height int64) error {
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
		NextValidators []common.Address
	}

	type req struct {
		height int64
		err    error
		res    *res
		retry  int64
	}

	r.Log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Debug("syncVerifier: start")

	for vr.Next() < height {
		rqch := make(chan *req, MonitorBlockMaxConcurrency)
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
				r.Log.WithFields(log.Fields{
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
					q.res.Header, q.err = r.Cl.getBlockHeaderByHeight(q.height)
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "syncVerifier: getBlockHeader: %v", q.err)
						return
					}
					q.res.Votes, q.err = r.Cl.GetVotesByHeight(
						&BlockHeightParam{Height: NewHexInt(int64(q.height))})
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "syncVerifier: GetVotesByHeight: %v", q.err)
						return
					}
					if len(vr.Validators(q.res.Header.NextValidatorsHash)) == 0 {
						q.res.NextValidators, q.err = r.Cl.getValidatorsByHash(q.res.Header.NextValidatorsHash)
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
			r.Log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Debug("syncVerifier: syncing")
		}
	}

	r.Log.WithFields(log.Fields{"height": vr.Next()}).Debug("syncVerifier: complete")
	return nil
}

func (r *ReceiverCore) ReceiveLoop(ctx context.Context, startHeight, startSeq uint64, callback func(txrs []*TxResult) error) (err error) {

	blockReq := r.BlockReq // copy

	blockReq.Height = NewHexInt(int64(startHeight))

	var vr *Verifier
	if r.Opts.Verifier != nil {
		vr, err = r.newVerifer(r.Opts.Verifier)
		if err != nil {
			return err
		}
	}

	type res struct {
		Height         int64
		Hash           common.HexHash
		Header         *BlockHeader
		Votes          []byte
		NextValidators []common.Address
		Txrs           []*TxResult
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
				err := r.Cl.MonitorBlock(ctx, &blockReq,
					func(conn *websocket.Conn, v *BlockNotification) error {
						if !errors.Is(ctx.Err(), context.Canceled) {
							bnch <- v
						}
						return nil
					},
					func(conn *websocket.Conn) {},
					func(c *websocket.Conn, err error) {})
				if err != nil {
					if websocket.IsUnexpectedCloseError(err) {
						reconnect() // unexpected error
						r.Log.WithFields(log.Fields{"error": err}).Error("reconnect: monitor block error")
					} else if !errors.Is(err, context.Canceled) {
						ech <- err
					}
				}
			}(ctxMonitorBlock, cancelMonitorBlock)

			// sync verifier
			if vr != nil {
				if err := r.syncVerifier(vr, next); err != nil {
					return errors.Wrapf(err, "sync verifier: %v", err)
				}
			}

		case br := <-brch:
			for ; br != nil; next++ {
				if br.Height%100 == 0 {
					r.Log.WithFields(log.Fields{"height": br.Height}).Debug("block notification")
				}
				if vr != nil {
					ok, err := vr.Verify(br.Header, br.Votes)
					if !ok || err != nil {
						if err != nil {
							r.Log.WithFields(log.Fields{"height": br.Height, "error": err}).Debug("receiveLoop: verification error")
						} else if !ok {
							r.Log.WithFields(log.Fields{"height": br.Height, "hash": br.Hash}).Debug("receiveLoop: invalid header")
						}
						reconnect() // reconnect websocket
						r.Log.WithFields(log.Fields{"height": br.Height, "hash": br.Hash}).Error("reconnect: verification failed")
						break
					}
					if err := vr.Update(br.Header, br.NextValidators); err != nil {
						return errors.Wrapf(err, "receiveLoop: update verifier: %v", err)
					}
				}
				if err := callback(br.Txrs); err != nil {
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
						r.Log.WithFields(log.Fields{
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
						r.Log.WithFields(log.Fields{"height": q.height, "error": q.err}).Debug("receiveLoop: req error")
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

							q.res.Header, q.err = r.Cl.getBlockHeaderByHeight(q.height)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "getBlockHeader: %v", q.err)
								return
							}
							// fetch votes, next validators only if verifier exists
							if vr != nil {
								q.res.Votes, q.err = r.Cl.GetVotesByHeight(
									&BlockHeightParam{Height: NewHexInt(int64(q.height))})
								if q.err != nil {
									q.err = errors.Wrapf(q.err, "GetVotesByHeight: %v", q.err)
									return
								}
								if len(vr.Validators(q.res.Header.NextValidatorsHash)) == 0 {
									q.res.NextValidators, q.err = r.Cl.getValidatorsByHash(q.res.Header.NextValidatorsHash)
									if q.err != nil {
										q.err = errors.Wrapf(q.err, "getValidatorsByHash: %v", q.err)
										return
									}
								}
							}

							if len(q.indexes) > 0 && len(q.events) > 0 {
								if len(q.indexes) != len(q.events) {
									q.err = fmt.Errorf("Got unequal values of len(indexes)=%v len(events)=%v", len(q.indexes), len(q.events))
								}
								var hr BlockHeaderResult
								_, err := codec.RLP.UnmarshalFromBytes(q.res.Header.Result, &hr)
								if q.err != nil {
									q.err = errors.Wrapf(q.err, "BlockHeaderResult.UnmarshalFromBytes: %v", err)
									return
								}
								for id := 0; id < len(q.indexes); id++ {
									for i, index := range q.indexes[id] {
										p := &ProofEventsParam{
											Index:     index,
											BlockHash: q.hash,
											Events:    q.events[id][i],
										}
										proofs, err := r.Cl.GetProofForEvents(p)
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

										idx, err := index.Value()
										if err != nil {
											q.err = errors.Wrapf(err, "Index value: %v", index)
										}
										result.EventLogs = result.EventLogs[:0]
										result.TxIndex = NewHexInt(int64(idx))
										result.BlockHeight = NewHexInt(int64(q.height))
										for j := 0; j < len(p.Events); j++ {
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
											result.EventLogs = append(result.EventLogs, el)
										}
										if len(result.EventLogs) > 0 {
											if len(result.EventLogs) == len(p.Events) {
												q.res.Txrs = append(q.res.Txrs, &result)
											} else {
												r.Log.WithFields(log.Fields{
													"height":              q.height,
													"receipt_index":       index,
													"got_num_events":      len(result.EventLogs),
													"expected_num_events": len(p.Events)}).Info("failed to verify all events for the receipt")
												q.err = errors.New("failed to verify all events for the receipt")
												return
											}
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
