package tezos

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/common/log"

	// "blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/pkg/errors"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/tezos/types"

)

const (
	BlockInterval              = 30 * time.Second
	BlockHeightPollInterval    = BlockInterval * 5
	BlockFinalityConfirmations = 2
	MonitorBlockMaxConcurrency = 300 // number of concurrent requests to synchronize older blocks from source chain
	RPCCallRetry               = 5
)

type receiver struct {
	log log.Logger
	src chain.BTPAddress
	dst chain.BTPAddress
	opts ReceiverOptions
	client *Client
}

func (r *receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	fmt.Println("reached to subscribe")
	src := tezos.MustParseAddress(string(r.src))
	r.client.Contract = contract.NewContract(src, r.client.Cl)

	opts.Seq++

	_errCh := make(chan error)

	verifierOpts := &VerifierOptions{
		BlockHeight: int64(opts.Height),
	}

	r.opts.Verifier = verifierOpts

	verifier, err := r.NewVerifier(ctx, r.opts.Verifier.BlockHeight) 
	fmt.Println("returned by the new verifier")
	if err != nil {
		_errCh <- err
		return _errCh, err 
	}

	if err = r.SyncVerifier(ctx, verifier, r.opts.Verifier.BlockHeight + 1,
	func (v []*chain.Receipt) error {
			fmt.Println("has to reach in this callback ")
			var vCP []*chain.Receipt
			var events []*chain.Event
			for _, receipt := range v{
				for _, event := range receipt.Events {
					switch {
						case event.Sequence == opts.Seq:
							events = append(events, event)
							opts.Seq++ 
						case event.Sequence > opts.Seq:
							return fmt.Errorf("invalid event seq")
						default:
							fmt.Println("default?????")
							events = append(events, event)
							opts.Seq++ 


					}
				}
				receipt.Events = events
				vCP = append(vCP, &chain.Receipt{Events: receipt.Events})
			} 
			if len(v) > 0 {
				msgCh <- &chain.Message{Receipts: vCP}
			}
			fmt.Println("returned nill")
			return nil 
		}); err != nil {
			_errCh <- err 
		}

	fmt.Println("reached to before monitor block")

	go func() {
		defer close(_errCh)
		lastHeight := opts.Height

		bn := &BnOptions{
			StartHeight: int64(opts.Height),
			Concurrnecy: r.opts.SyncConcurrency,
		}
		if err := r.receiveLoop(ctx, bn,
		func (blN *types.BlockNotification) error {
			fmt.Println("has to reach in this callback ")

			if blN.Height.Uint64() != lastHeight + 1{
				return fmt.Errorf(
					"block notification: expected=%d, got %d", lastHeight + 1, blN.Height.Uint64())
			}

			var vCP []*chain.Receipt
			var events []*chain.Event
			v := blN.Receipts
			for _, receipt := range v{
				for _, event := range receipt.Events {
					switch {
						case event.Sequence == opts.Seq:
							events = append(events, event)
							opts.Seq++ 
						case event.Sequence > opts.Seq:
							return fmt.Errorf("invalid event seq")
						default:
							events = append(events, event)
							opts.Seq++ 
					}
				}
				receipt.Events = events
				vCP = append(vCP, &chain.Receipt{Events: receipt.Events})
			} 
			if len(v) > 0 {
				fmt.Println("reached to sending message")
				msgCh <- &chain.Message{Receipts: vCP}
			}
			fmt.Println("returned nill")
			lastHeight++
			return nil 
		}); err != nil {
			fmt.Println(err)
			_errCh <- err 
		}

		fmt.Println("Printing from inside the receiver")
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

func NewReceiver(src, dst chain.BTPAddress, urls []string, rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error){
	var newClient *Client
	var err error

	if len(urls) == 0 {
		return nil, fmt.Errorf("Empty urls")
	}
	verifierOpts := &VerifierOptions{
		BlockHeight: int64(2468690),
	}

	receiverOpts := &ReceiverOptions{
		SyncConcurrency: 50,
		Verifier: verifierOpts,
	}

	receiver := &receiver{
		log: l,
		src: src,
		dst: dst,
		opts: *receiverOpts,
	}

	if receiver.opts.SyncConcurrency < 1 {
		receiver.opts.SyncConcurrency = 1
	} else if receiver.opts.SyncConcurrency > MonitorBlockMaxConcurrency {
		receiver.opts.SyncConcurrency = MonitorBlockMaxConcurrency
	}

	srcAddr := tezos.MustParseAddress(string(src))

	newClient, err = NewClient(urls[0], srcAddr, receiver.log)

	if err != nil {
		return nil, err
	}
	receiver.client = newClient

	return receiver, nil
}

type ReceiverOptions struct {
	SyncConcurrency uint64           `json:"syncConcurrency"`
	Verifier        *VerifierOptions `json:"verifier"`
}

func (r *receiver) NewVerifier(ctx context.Context, previousHeight int64) (vri IVerifier, err error) {
	fmt.Println("reached to verifyer")
	header, err := r.client.GetBlockHeaderByHeight(ctx, r.client.Cl, previousHeight)
	fmt.Println("reached to after block header ")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("returned from here?")
	fittness, err := strconv.ParseInt(string(header.Fitness[1].String()), 16, 64)
	if err != nil {
		return nil, err
	}

	fmt.Println("before chain id")
	chainIdHash, err := r.client.Cl.GetChainId(ctx)
	if err != nil {
		return nil, err
	}

	id := chainIdHash.Uint32()

	if err != nil {
		return nil, err 
	} 

	vr := &Verifier{
		mu: sync.RWMutex{},
		next: header.Level + 1,
		parentHash: header.Hash,
		parentFittness: fittness,
		chainID: id,
	}
	fmt.Println("returned to the original")
	fmt.Println(vr.parentHash)
	return vr, nil
}

func (r *receiver) SyncVerifier(ctx context.Context, vr IVerifier, height int64, callback func([]*chain.Receipt) error) error {
	if height == vr.Next() {
		fmt.Println("returned from here")
		return nil 
	}

	if vr.Next() > height {
		return fmt.Errorf("Invalida target height: Verifier height (%d) > target height (%d)", vr.Next(), height)
	}

	type res struct {
		Height int64
		Header *rpc.BlockHeader
		Block *rpc.Block 
		Votes int64
	}

	type req struct {
		height int64
		err error 
		res *res
		retry int64
	}
	fmt.Println("reached before starting to log")
	// r.log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Info("syncVerifier: start")
	
	fmt.Println("reached in sync verifier")
	var prevHeader *rpc.BlockHeader 

	cursor := vr.Next()

	for cursor <= height {
		fmt.Println("reached inside for")
		fmt.Println(r.opts.SyncConcurrency)
		
		rqch := make(chan *req, r.opts.SyncConcurrency)
		fmt.Println(len(rqch))
		fmt.Println(cap(rqch))
		for i := cursor; len(rqch) < cap(rqch); i++{
			rqch <- &req{height: i, retry: 5}
		}
		sres := make([]*res, 0, len(rqch))
		fmt.Println("reached here after sres")
		for q := range rqch {
			switch {
			case q.err != nil:
				if q.retry > 0 {
					q.retry--
					q.res, q.err = nil, nil 
					rqch <- q
					continue
				}
				// r.log.WithFields(log.Fields{"height": q.height, "error": q.err.Error()}).Debug("syncVerifier: req error")
				sres = append(sres, nil)
				if len(sres) == cap(sres) {
					close(rqch)
				}
			case q.res != nil:
				fmt.Println("should reach here in the second loop ")
				sres = append(sres, q.res)
				fmt.Println(cap(sres))
				if len(sres) == cap(sres){
					fmt.Println("closes channel")
					close(rqch)
				}
			default:
				fmt.Println("has to reach in this default ")
				go func(q *req) {
					defer func() {
						time.Sleep(500 * time.Millisecond)
						rqch <- q
					}()
					if q.res == nil {
						fmt.Println("should reach here in nil portion")
						q.res = &res{}
					}
					q.res.Height = q.height
					q.res.Header, q.err = r.client.GetBlockHeaderByHeight(ctx, r.client.Cl, q.height)
					fmt.Println(q.res.Header)
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "syncVerifier: getBlockHeader: %v", q.err)
						return
					}
					q.res.Block, q.err = r.client.GetBlockByHeight(ctx, r.client.Cl, q.height)
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "syncVerifier: getBlock: %v", q.err)
						return
					}
					fmt.Println(q.res.Block)
				}(q)
			}
		
		}
		_sres, sres := sres, sres[:0]
		for _, v := range _sres {
			if v != nil {
				fmt.Println("should reach in eliminating the null ", v.Height)
				sres = append(sres, v)
			}
		}

		fmt.Printf("The lenght of sres is %d\n", len(sres))

		if len(sres) > 0 {
			sort.SliceStable(sres, func(i, j int) bool {
				return sres[i].Height < sres[j].Height
			})
			for i := range sres {
				cursor++
				next := sres[i]
				if prevHeader == nil {
					prevHeader = next.Header
					continue 
				}
				if vr.Next() >= height {
					fmt.Println("did it just break")
					break
				}

				fmt.Println("has it reached to verification")
				fmt.Println(next.Header.Level)

				err := vr.Verify(ctx, prevHeader, next.Block.Metadata.Baker, r.client.Cl, next.Header.Hash)

				if err != nil {
					cursor = vr.Height() + 1
					prevHeader = nil 
					fmt.Println(cursor)
					fmt.Println("when some verification is failed prompts it to get the data again from that point")
					time.Sleep(15 * time.Second)
					break 
					// return errors.Wrapf(err, "syncVerifier: Verify: %v", err)
				}

				fmt.Println("verified block now updating ")

				err = vr.Update(prevHeader)
				if err != nil {
					return errors.Wrapf(err, "syncVerifier: Update: %v", err)
				}

				prevHeader = next.Header
			}

		}
			// r.log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Debug("syncVerifier: syncing")
	}
	// r.log.WithFields(log.Fields{"height": vr.Next()}).Info("syncVerifier: complete")

	fmt.Println("sync complete")
	return nil 
}

type BnOptions struct {
	StartHeight int64
	Concurrnecy uint64
}

func (r *receiver) receiveLoop(ctx context.Context, opts *BnOptions, callback func(v *types.BlockNotification) error) (err error){
	fmt.Println("reached to receivelopp")
	if opts == nil {
		return errors.New("receiveLoop: invalid options: <nil>")
	}

	var vr IVerifier

	if r.opts.Verifier != nil{
		vr, err = r.NewVerifier(ctx, r.opts.Verifier.BlockHeight)
		if err != nil {
			return err
		}
		err = r.SyncVerifier(ctx, vr, r.opts.Verifier.BlockHeight + 1, func(r []*chain.Receipt) error {return nil})
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
	next, latest := r.opts.Verifier.BlockHeight + 1, latestHeight()

	var lbn *types.BlockNotification

	for {
		select {
		case <- ctx.Done():
			return nil
		case <- heightTicker.C:
			latest++
		case <- heightPoller.C:
			if height := latestHeight(); height > 0 {
				latest = height 
				// r.log.WithFields(log.Fields{"latest": latest, "next": next}).Debug("poll height")
			}
		case bn := <-bnch:
			fmt.Println("has it reached in the block notification channel")
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					if bn.Height.Cmp(lbn.Height) == 0 {
						if bn.Header.Predecessor != lbn.Header.Predecessor {
							// r.log.WithFields(log.Fields{"lbnParentHash": lbn.Header.Predecessor, "bnParentHash": bn.Header.Predecessor}).Error("verification failed on retry ")
							break
						}
					} else {
						if vr != nil {
							fmt.Println("vr is not nil")
							// header := bn.Header
							if err := vr.Verify(ctx, lbn.Header, bn.Proposer, r.client.Cl, bn.Header.Hash); err != nil { // change accordingly 
								// r.log.WithFields(log.Fields{
								// 	"height":     lbn.Height,
								// 	"lbnHash":    lbn.Hash,
								// 	"nextHeight": next,
								// 	"bnHash":     bn.Hash}).Error("verification failed. refetching block ", err)
								fmt.Println("error in verifying ")
								time.Sleep(20 * time.Second)
								next--
								break
							}
							if err := vr.Update(lbn.Header); err != nil {
								return errors.Wrapf(err, "receiveLoop: vr.Update: %v", err)
							}
						}
						if err := callback(lbn); err != nil {
							return errors.Wrapf(err, "receiveLoop: callback: %v", err)
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
				//r.log.WithFields(log.Fields{"lenBnch": len(bnch), "height": t.Height}).Info("remove unprocessed block noitification")
			}
			
		default:
			if next >= latest {
				time.Sleep(10 * time.Second)
				continue
			}

			type bnq struct {
				h     int64
				v     *types.BlockNotification
				err   error
				retry int
			}

			qch := make(chan *bnq, cap(bnch))

			for i:= next; i < latest && len(qch) < cap(qch); i++{
				qch <- &bnq{i, nil, nil, RPCCallRetry}
			}

			if len(qch) == 0 {
				// r.log.Error("Fatal: Zero length of query channel. Avoiding deadlock")
				continue
			}
			bns := make([]*types.BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil :
					if q.retry > 0 {
						q.retry --
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
					fmt.Println("reached in default of receive loop")
					go func(q *bnq) {
						defer func() {
							time.Sleep(500 * time.Millisecond)
							qch <- q
						}()

						if q.v == nil {
							q.v = &types.BlockNotification{}
						}
						q.v.Height = (&big.Int{}).SetInt64(q.h)

						if q.v.Header == nil {
							header, err := r.client.GetBlockHeaderByHeight(ctx, r.client.Cl, q.v.Height.Int64())
							if err != nil {
								q.err = errors.Wrapf(err, "GetHeaderByHeight: %v", err)
								return
							}
						q.v.Header = header // change accordingly  
						q.v.Hash = q.v.Hash	// change accordingly 
						}
					
						if q.v.HasBTPMessage == nil {
							fmt.Println("height: ", q.v.Height.Int64())
							block, err := r.client.GetBlockByHeight(ctx, r.client.Cl, q.v.Height.Int64())

							if err != nil {
								return
							}
							q.v.Proposer = block.Metadata.Proposer
							

							hasBTPMessage, receipt, err := returnTxMetadata2(block, r.client.Contract.Address(), q.v.Height.Int64(), r.client)

							if err != nil {
								q.err = errors.Wrapf(err, "hasBTPMessage: %v", err)
								return
							}
							q.v.HasBTPMessage = &hasBTPMessage

							if receipt != nil {
								q.v.Receipts = receipt
							}
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
					if v.Height.Int64() == next + int64(i) {
						bnch <- v
					}
				}
			}

		}

	}
}

