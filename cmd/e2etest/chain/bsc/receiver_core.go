package bsc

import (
	"context"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type ReceiverCore struct {
	Log      log.Logger
	Opts     ReceiverOptions
	Cls      []*ethclient.Client
	BlockReq ethereum.FilterQuery
}

const (
	BlockInterval              = 2 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	defaultReadTimeout         = 15 * time.Second
	monitorBlockMaxConcurrency = 1000 // number of concurrent requests to synchronize older blocks from source chain
)

type BnOptions struct {
	StartHeight     uint64
	Concurrency     uint64
	VerifierOptions *VerifierOptions
}

type VerifierOptions struct {
}
type ReceiverOptions struct {
	Verifier        *VerifierOptions `json:"verifier"`
	SyncConcurrency uint64           `json:"syncConcurrency"`
}

func (r *ReceiverCore) client() *ethclient.Client {
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

	// block notification channel
	// (buffered: to avoid deadlock)
	// increase concurrency parameter for faster sync
	bnch := make(chan *BlockNotification, opts.Concurrency)

	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	latestHeight := func() uint64 {
		height, err := r.client().BlockNumber(context.TODO())
		if err != nil {
			r.Log.WithFields(log.Fields{"error": err}).Error("receiveLoop: failed to GetBlockNumber")
			return 0
		}
		return height
	}

	next, latest := opts.StartHeight, latestHeight()

	// last unverified block notification
	var lbn *BlockNotification

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
					//r.Log.Debugf("receiveLoop: bnq: h=%d:%v, %v", q.h, q.v.Header.Hash(), q.err)
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
						q.v.Header, q.err = r.client().HeaderByNumber(context.TODO(), q.v.Height)
						if q.err != nil {
							//q.err = errors.Wrapf(q.err, "GetHmyHeaderByHeight: %v", q.err)
							return
						}
						if q.v.Header.GasUsed > 0 {
							ht := big.NewInt(q.v.Height.Int64())
							r.BlockReq.FromBlock = ht
							r.BlockReq.ToBlock = ht
							q.v.Logs, q.err = r.client().FilterLogs(context.TODO(), r.BlockReq)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "FilterLogs: %v", q.err)
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
