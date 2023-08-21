package tezos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/common/log"

	// "blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/tezos/types"
	"github.com/pkg/errors"
)

const (
	BlockInterval              = 15 * time.Second
	BlockHeightPollInterval    = BlockInterval * 5
	BlockFinalityConfirmations = 2
	MonitorBlockMaxConcurrency = 300 // number of concurrent requests to synchronize older blocks from source chain
	RPCCallRetry               = 5
)

type receiver struct {
	log    log.Logger
	src    chain.BTPAddress
	dst    chain.BTPAddress
	opts   ReceiverOptions
	client *Client
}

var RelaySyncStatusLog bool

func (r *receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	src := tezos.MustParseAddress(r.src.ContractAddress())
	r.client.Contract = contract.NewContract(src, r.client.Cl)

	opts.Seq++

	_errCh := make(chan error)

	go func() {
		defer close(_errCh)
		lastHeight := opts.Height + 1

		bn := &BnOptions{
			StartHeight: int64(opts.Height),
			Concurrnecy: r.opts.SyncConcurrency,
		}
		if err := r.receiveLoop(ctx, bn,
			func(blN *types.BlockNotification) error {
				r.log.WithFields(log.Fields{"height": blN.Height}).Debug("block notification")

				if blN.Height.Uint64() != lastHeight {
					return fmt.Errorf(
						"block notification: expected=%d, got %d", lastHeight, blN.Height.Uint64())
				}

				// var events []*chain.Event
				receipts := blN.Receipts
				for _, receipt := range receipts {
					events := receipt.Events[:0]
					for _, event := range receipt.Events {
						switch {
						case event.Sequence == opts.Seq:
							events = append(events, event)
							opts.Seq++
						case event.Sequence > opts.Seq:
							r.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, blN.Height.Uint64())
							return fmt.Errorf(
								"block notification: expected=%d, got=%d",
								lastHeight+1, blN.Height.Uint64())
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
			_errCh <- err
		}
	}()

	return _errCh, nil
}

// func (r *receiver) getRelayReceipts(v *chain.BlockNotification) []*chain.Receipt {
// 	sc := common.HexToAddress(string(r.src))
// 	var receipts[]*chain.Receipt
// 	var events []*chain.Event

// 	for i, receipt := range v.Receipts {
// 		events := events[:0]

// 		}
// 	}

func NewReceiver(src, dst chain.BTPAddress, urls []string, rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error) {
	var newClient *Client
	var err error

	if len(urls) == 0 {
		return nil, fmt.Errorf("Empty urls")
	}

	receiver := &receiver{
		log: l,
		src: src,
		dst: dst,
	}

	err = json.Unmarshal(rawOpts, &receiver.opts)

	if receiver.opts.SyncConcurrency < 1 {
		receiver.opts.SyncConcurrency = 1
	} else if receiver.opts.SyncConcurrency > MonitorBlockMaxConcurrency {
		receiver.opts.SyncConcurrency = MonitorBlockMaxConcurrency
	}

	srcAddr := tezos.MustParseAddress(src.ContractAddress())
	bmcManagement := tezos.MustParseAddress(receiver.opts.BMCManagment)

	newClient, err = NewClient(urls[0], srcAddr, bmcManagement, receiver.log)

	if err != nil {
		return nil, err
	}
	receiver.client = newClient

	return receiver, nil
}

type ReceiverOptions struct {
	SyncConcurrency uint64           `json:"syncConcurrency"`
	Verifier        *VerifierOptions `json:"verifier"`
	BMCManagment    string           `json:"bmcManagement"`
}

func (r *receiver) NewVerifier(ctx context.Context, previousHeight int64) (vri IVerifier, err error) {
	block, err := r.client.GetBlockByHeight(ctx, r.client.Cl, previousHeight)
	if err != nil {
		return nil, err
	}

	fittness, err := strconv.ParseInt(string(block.Header.Fitness[1].String()), 16, 64)
	if err != nil {
		return nil, err
	}

	chainIdHash, err := r.client.Cl.GetChainId(ctx)
	if err != nil {
		return nil, err
	}

	id := chainIdHash.Uint32()

	if err != nil {
		return nil, err
	}

	vr := &Verifier{
		mu:                  sync.RWMutex{},
		next:                block.Header.Level + 1,
		parentHash:          block.Hash,
		parentFittness:      fittness,
		chainID:             id,
		cl:                  r.client,
		validators:          make(map[tezos.Address]bool),
		validatorsPublicKey: make(map[tezos.Address]tezos.Key),
	}

	vr.updateValidatorsAndCycle(ctx, previousHeight, block.Metadata.LevelInfo.Cycle)
	return vr, nil
}


type BnOptions struct {
	StartHeight int64
	Concurrnecy uint64
}

// merging the syncing and receiving function

func (r *receiver) receiveLoop(ctx context.Context, opts *BnOptions, callback func(v *types.BlockNotification) error) (err error) {
	if opts == nil {
		return errors.New("receiveLoop: invalid options: <nil>")
	}

	RelaySyncStatusLog = false

	var vr IVerifier

	if r.opts.Verifier != nil {
		vr, err = r.NewVerifier(ctx, r.opts.Verifier.BlockHeight)
		if err != nil {
			return err
		}
	}
	bnch := make(chan *types.BlockNotification, r.opts.SyncConcurrency)
	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	latestHeight := func() int64 {
		block, err := r.client.GetLastBlock(ctx, r.client.Cl)
		if err != nil {
			return 0
		}
		return block.GetLevel()
	}
	next, latest := r.opts.Verifier.BlockHeight+1, latestHeight()

	var lbn *types.BlockNotification

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-heightTicker.C:
			latest++
		case <-heightPoller.C:
			if height := latestHeight(); height > 0 {
				latest = height - 5
				r.log.WithFields(log.Fields{"latest": latest, "next": next}).Debug("poll height")
			}
		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					if bn.Height.Cmp(lbn.Height) == 0 {
						if bn.Header.Predecessor != lbn.Header.Predecessor {
							r.log.WithFields(log.Fields{"lbnParentHash": lbn.Header.Predecessor, "bnParentHash": bn.Header.Predecessor}).Error("verification failed on retry ")
							break
						}
					} else {
						if vr != nil {
							if err := vr.Verify(ctx, lbn); err != nil { 
								r.log.WithFields(log.Fields{
									"height":     lbn.Height,
									"lbnHash":    lbn.Hash,
									"nextHeight": next,
									"bnHash":     bn.Hash}).Error("verification failed. refetching block ", err)
								next--
								break
							}
							if err := vr.Update(ctx, lbn); err != nil {
								return errors.Wrapf(err, "receiveLoop: vr.Update: %v", err)
							}
						}
						if lbn.Header.Level > opts.StartHeight {
							if !RelaySyncStatusLog {
								r.log.WithFields(log.Fields{"height": vr.Next()}).Info("syncVerifier: complete")
								RelaySyncStatusLog = true
							}
							if vr.LastVerifiedBn() != nil  && vr.LastVerifiedBn().Header.Level > opts.StartHeight{
								if err := callback(vr.LastVerifiedBn()); err != nil {
									return errors.Wrapf(err, "receiveLoop: callback: %v", err)
								}
							}
						} else {
							r.log.WithFields(log.Fields{"height": vr.Next(), "target": opts.StartHeight}).Debug("syncVerifier: syncing")
						}
					}
				}
				if lbn, bn = bn, nil; len(bnch) > 0 {
					bn = <-bnch
				}
			}
			// remove unprocessed notifications
			for len(bnch) > 0 {
				<-bnch
				// r.log.WithFields(log.Fields{"lenBnch": len(bnch), "height": t.Height}).Info("remove unprocessed block noitification")
			}

		default:
			if next >= latest {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			type bnq struct {
				h     int64
				v     *types.BlockNotification
				err   error
				retry int
			}

			qch := make(chan *bnq, cap(bnch))

			for i := next; i < latest && len(qch) < cap(qch); i++ {
				qch <- &bnq{i, nil, nil, RPCCallRetry}
			}

			if len(qch) == 0 {
				r.log.Error("Fatal: Zero length of query channel. Avoiding deadlock")
				continue
			}
			bns := make([]*types.BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil:
					if q.retry > 0 {
						q.retry--
						q.v, q.err = nil, nil
						qch <- q
						continue
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
							q.v = &types.BlockNotification{}
						}
						q.v.Height = (&big.Int{}).SetInt64(q.h)
						block := &rpc.Block{}

						if q.v.Header == nil {
							block, err := r.client.GetBlockByHeight(ctx, r.client.Cl, q.v.Height.Int64())
							if err != nil {
								q.err = errors.Wrapf(err, "GetHeaderByHeight: %v", err)
								return
							}
							q.v.Header = &block.Header 
							q.v.Hash = block.Hash 
							q.v.Block = block
						}

						if q.v.HasBTPMessage == nil && q.v.Height.Int64() > opts.StartHeight {
							if err != nil {
								return
							}
							q.v.Proposer = block.Metadata.Proposer
							
							hasBTPMessage, receipt, err := filterTransactionOperations(q.v.Block, r.client.Contract.Address(), q.v.Height.Int64(), r.client, r.dst.String())

							if err != nil {
								q.err = errors.Wrapf(err, "hasBTPMessage: %v", err)
								return
							}
							q.v.HasBTPMessage = &hasBTPMessage

							if receipt != nil {
								q.v.Receipts = receipt
							}
						} else {
							return
						}

						if !*q.v.HasBTPMessage {
							return
						}
					}(q)
				}
			}
			// filtering nil
			_bns_, bns := bns, bns[:0]

			for _, v := range _bns_ {
				if v != nil {
					bns = append(bns, v)
				}
			}

			if len(bns) > 0 {
				sort.SliceStable(bns, func(i, j int) bool {
					return bns[i].Height.Int64() < bns[j].Height.Int64()
				})
				for i, v := range bns {
					if v.Height.Int64() == next+int64(i) {
						bnch <- v
					}
				}
			}

		}

	}
}
