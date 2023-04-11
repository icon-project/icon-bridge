package tezos

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/common/log"

	// "blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/pkg/errors"
)

const (
	BlockInterval              = 30 * time.Second
	BlockHeightPollInterval    = BlockInterval * 5
	BlockFinalityConfirmations = 2
	MonitorBlockMaxConcurrency = 10 // number of concurrent requests to synchronize older blocks from source chain
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

	verifier, err := r.NewVerifier(ctx, int64(opts.Height)) 
	fmt.Println("returned by the new verifier")
	if err != nil {
		_errCh <- err
		return _errCh, err 
	}

	err = r.syncVerifier(ctx, verifier, int64(opts.Height + 80))

	if err != nil {
		_errCh <- err
		return _errCh, err
	}

	fmt.Println("reached to before monitor block")

	go func() {
		defer close(_errCh)
		
		if err := r.client.MonitorBlock(ctx, int64(opts.Height + 80), verifier, 
		func (v []*chain.Receipt) error {
			fmt.Println(v[0].Events[0].Message)
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
							events = append(events, event)
							opts.Seq++ 


					}
				}
				receipt.Events = events
				vCP = append(vCP, &chain.Receipt{Events: receipt.Events})
			} 
			if len(v) > 0 {
				fmt.Println("reached to sending message")
				fmt.Println(vCP[0].Events[0].Message)
				msgCh <- &chain.Message{Receipts: vCP}
			}
			fmt.Println("returned nill")
			return nil 
		}); err != nil {
			_errCh <- err 
		}

		fmt.Println("Printing from inside the receiver")
	}()
	
	return _errCh, nil
}

// func (r *receiver) getRelayReceipts(v *chain.BlockNotification) []*chain.Receipt {
// 	sc := tezos.MustParseAddress(r.src.ContractAddress())
// 	var receipts []*chain.Receipt
// 	var events []*chain.Event
// 	for _, receipt := range v.Receipts{
// 		events = append(events, &chain.Event{
// 			Next: chain.BTPAddress(r.dst),

// 		})
// 	}
// }

func NewReceiver(src, dst chain.BTPAddress, urls []string, rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error){
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

	if receiver.opts.SyncConcurrency < 1 {
		receiver.opts.SyncConcurrency = MonitorBlockMaxConcurrency
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

func (r *receiver) syncVerifier(ctx context.Context, vr IVerifier, height int64) error {
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
				fmt.Println("should reach in eliminating the null")
				sres = append(sres, v)
			}
		}

		fmt.Printf("The lenght of sres is %d\n", len(sres))

		if len(sres) > 0 {
			sort.SliceStable(sres, func(i, j int) bool {
				return sres[i].Height < sres[j].Height
			})
			for i := range sres {
				fmt.Println("Has to reach in the first time only")
				cursor++
				next := sres[i]
				if prevHeader == nil {
					fmt.Println("Previous header is nil ")
					fmt.Println(next.Header.Level)
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
					return errors.Wrapf(err, "syncVerifier: Verify: %v", err)
				}

				fmt.Println("verified block now updating ")

				err = vr.Update(prevHeader)
				if err != nil {
					return errors.Wrapf(err, "syncVerifier: Update: %v", err)
				}
				prevHeader = next.Header
			}
			// r.log.WithFields(log.Fields{"height": vr.Next(), "target": height}).Debug("syncVerifier: syncing")
		}
	}

	// r.log.WithFields(log.Fields{"height": vr.Next()}).Info("syncVerifier: complete")

	fmt.Println("sync complete")
	return nil 
}
