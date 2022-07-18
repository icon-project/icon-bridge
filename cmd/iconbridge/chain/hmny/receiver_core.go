//go:build hmny
// +build hmny

package hmny

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type ReceiverCore struct {
	Log  log.Logger
	Opts ReceiverOptions
	Cls  []*Client
}

func (r *ReceiverCore) client() *Client {
	randInt := rand.Intn(len(r.Cls))
	return r.Cls[randInt]
}

func (r *ReceiverCore) ReceiveLoop(ctx context.Context, opts *BnOptions, callback func(v *BlockNotification) error) error {

	if opts == nil {
		return errors.New("receiveLoop: invalid options: <nil>")
	}

	if opts.Concurrency < 1 || opts.Concurrency > monitorBlockMaxConcurrency {
		concurrency := opts.Concurrency
		if concurrency < 1 {
			opts.Concurrency = 1
		} else {
			opts.Concurrency = monitorBlockMaxConcurrency
		}
		r.Log.Warnf("receiveLoop: opts.Concurrency (%d): value out of range [%d, %d]: setting to default %d",
			concurrency, 1, monitorBlockMaxConcurrency, opts.Concurrency)
	}

	if opts.VerifierOptions != nil &&
		opts.StartHeight < opts.VerifierOptions.BlockHeight {
		return fmt.Errorf(
			"receiveLoop: start height (%d) < verifier height (%d)",
			opts.StartHeight, opts.VerifierOptions.BlockHeight,
		)
	}
	var vr Verifier
	if opts.VerifierOptions != nil {
		var err error
		vr, err = r.client().newVerifier(opts.VerifierOptions)
		if err != nil {
			return errors.Wrapf(err, "receiveLoop: NewVerifier: %v", err)
		}
		err = r.client().syncVerifier(vr, opts.StartHeight)
		if err != nil {
			return errors.Wrapf(err, "receiveLoop: cl.syncVerifier: %v", err)
		}
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
			r.Log.WithFields(log.Fields{"error": err}).Error("receiveLoop: failed to GetBlockNumber")
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
			return nil

		case <-heightTicker.C:
			latest++

		case <-heightPoller.C:
			if height := latestHeight(); height > latest {
				latest = height
				if next > latest {
					r.Log.Debugf("receiveLoop: skipping; latest=%d, next=%d", latest, next)
				}
			}

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					if vr != nil {
						ok, err := vr.Verify(lbn.Header,
							bn.Header.LastCommitBitmap, bn.Header.LastCommitSignature)
						if err != nil {
							r.Log.Errorf("receiveLoop: signature validation failed: h=%d, %v", lbn.Header.Number, err)
							break
						}
						if !ok {
							r.Log.Errorf("receiveLoop: invalid header: signature validation failed: h=%d", lbn.Header.Number)
							break
						}
						if err := vr.Update(lbn.Header); err != nil {
							return errors.Wrapf(err, "receiveLoop: update verifier: %v", err)
						}
					}
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
					r.Log.Debugf("receiveLoop: bnq: h=%d:%v, %v", q.h, q.v.Header.Hash(), q.err)
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
