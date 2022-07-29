package bsc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	BlockInterval              = 5 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	monitorBlockMaxConcurrency = 1000 // number of concurrent requests to synchronize older blocks from source chain
)

func NewReceiver(
	src, dst chain.BTPAddress, urls []string,
	opts map[string]interface{}, l log.Logger) (chain.Receiver, error) {
	r := &receiver{
		log: l,
		src: src,
		dst: dst,
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}
	err := r.opts.Unmarshal(opts)
	if err != nil {
		return nil, err
	}
	r.cls, r.bmcs, err = newClients(urls, src.ContractAddress(), r.log)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type ReceiverOptions struct {
	SyncConcurrency uint64 `json:"syncConcurrency"`
}

func (opts *ReceiverOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

type receiver struct {
	log  log.Logger
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts ReceiverOptions
	cls  []*Client
	bmcs []*BMC
}

func (r *receiver) client() *Client {
	randInt := rand.Intn(len(r.cls))
	return r.cls[randInt]
}

func (r *receiver) bmcClient() *BMC {
	randInt := rand.Intn(len(r.cls))
	return r.bmcs[randInt]
}

type BnOptions struct {
	StartHeight uint64
	Concurrency uint64
}

func (r *receiver) receiveLoop(ctx context.Context, opts *BnOptions, callback func(v *BlockNotification) error) error {

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
		r.log.Warnf("receiveLoop: opts.Concurrency (%d): value out of range [%d, %d]: setting to default %d",
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
					r.log.Debugf("receiveLoop: bnq: h=%d:%v, %v", q.h, q.v.Header.Hash(), q.err)
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
						q.v.Header, q.err = r.client().GetHeaderByHeight(q.v.Height)
						if q.err != nil {
							q.err = errors.Wrapf(q.err, "GetHeaderByHeight: %v", q.err)
							return
						}
						q.v.Hash = q.v.Header.Hash()
						if q.v.Header.GasUsed > 0 {
							q.v.Receipts, q.err = r.client().GetBlockReceipts(q.v.Hash)
							if q.err == nil {
								receiptsRoot := ethTypes.DeriveSha(q.v.Receipts, trie.NewStackTrie(nil))
								if !bytes.Equal(receiptsRoot.Bytes(), q.v.Header.ReceiptHash.Bytes()) {
									q.err = fmt.Errorf(
										"invalid receipts: remote=%v, local=%v",
										q.v.Header.ReceiptHash, receiptsRoot)
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

func (r *receiver) Subscribe(
	ctx context.Context, msgCh chan<- *chain.Message,
	opts chain.SubscribeOptions) (errCh <-chan error, err error) {

	opts.Seq++

	_errCh := make(chan error)

	go func() {
		defer close(_errCh)
		lastHeight := opts.Height - 1
		if err := r.receiveLoop(ctx,
			&BnOptions{
				StartHeight: opts.Height,
				Concurrency: r.opts.SyncConcurrency,
			},
			func(v *BlockNotification) error {
				r.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")

				if v.Height.Uint64() != lastHeight+1 {
					r.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, v.Height.Uint64())
					return fmt.Errorf(
						"block notification: expected=%d, got=%d",
						lastHeight+1, v.Height.Uint64())
				}

				receipts := r.getRelayReceipts(v)
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
						}
					}
					receipt.Events = events
				}
				if len(receipts) > 0 {
					msgCh <- &chain.Message{Receipts: receipts}
				}
				lastHeight++
				return nil
			}); err != nil {
			r.log.Errorf("receiveLoop terminated: %v", err)
			_errCh <- err
		}
	}()

	return _errCh, nil
}

func (r *receiver) getRelayReceipts(v *BlockNotification) []*chain.Receipt {
	sc := common.HexToAddress(r.src.ContractAddress())
	var receipts []*chain.Receipt
	var events []*chain.Event
	for i, receipt := range v.Receipts {
		events := events[:0]
		for _, log := range receipt.Logs {
			if !bytes.Equal(log.Address.Bytes(), sc.Bytes()) {
				continue
			}
			msg, err := r.bmcClient().ParseMessage(ethTypes.Log{
				Data: log.Data, Topics: log.Topics,
			})
			if err == nil {
				events = append(events, &chain.Event{
					Next:     chain.BTPAddress(msg.Next),
					Sequence: msg.Seq.Uint64(),
					Message:  msg.Msg,
				})
			}
		}
		if len(events) > 0 {
			rp := &chain.Receipt{}
			rp.Index, rp.Height = uint64(i), v.Height.Uint64()
			rp.Events = append(rp.Events, events...)
			receipts = append(receipts, rp)
		}
	}
	return receipts
}
