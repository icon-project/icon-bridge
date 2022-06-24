package hmny

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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	BlockInterval              = 2 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	defaultReadTimeout         = 15 * time.Second
	monitorBlockMaxConcurrency = 1000 // number of concurrent requests to synchronize older blocks from source chain
)

func NewReceiver(
	src, dst chain.BTPAddress, urls []string,
	opts map[string]interface{}, l log.Logger) (chain.SubscritionAPI, error) {
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
	r.cls, err = newClients(urls, src.ContractAddress(), r.log)
	if err != nil {
		return nil, err
	}
	return r, nil
}

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

type receiver struct {
	log  log.Logger
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts receiverOptions
	cls  []*client
}

func (r *receiver) client() *client {
	return r.cls[rand.Intn(len(r.cls))]
}

// Options for a new block notifications channel
type bnOptions struct {
	StartHeight     uint64
	Concurrency     uint64
	VerifierOptions *VerifierOptions
}

func (r *receiver) receiveLoop(ctx context.Context, opts *bnOptions, callback func(v *BlockNotification) error) error {

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

	if opts.VerifierOptions != nil &&
		opts.StartHeight < opts.VerifierOptions.BlockHeight {
		return fmt.Errorf(
			"receiveLoop: start height (%d) < verifier height (%d)",
			opts.StartHeight, opts.VerifierOptions.BlockHeight,
		)
	}
	vr, err := r.client().newVerifier(opts.VerifierOptions)
	if err != nil {
		return errors.Wrapf(err, "receiveLoop: NewVerifier: %v", err)
	}
	if err = r.client().syncVerifier(vr, opts.StartHeight); err != nil {
		return errors.Wrapf(err, "receiveLoop: cl.syncVerifier: %v", err)
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
					ok, err := vr.Verify(lbn.Header,
						bn.Header.LastCommitBitmap, bn.Header.LastCommitSignature)
					if err != nil {
						r.log.Errorf("receiveLoop: signature validation failed: h=%d, %v", lbn.Header.Number, err)
						break
					}
					if !ok {
						r.log.Errorf("receiveLoop: invalid header: signature validation failed: h=%d", lbn.Header.Number)
						break
					}
					if err := vr.Update(lbn.Header); err != nil {
						return errors.Wrapf(err, "receiveLoop: update verifier: %v", err)
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
			msg, err := r.client().bmc.ParseMessage(ethtypes.Log{
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

func (r *receiver) Subscribe(
	ctx context.Context, sinkChan chan<- *chain.SubscribedEvent, _errCh chan<- error,
	opts chain.SubscribeOptions) (err error) {

	opts.Seq++

	go func() {
		lastHeight := opts.Height - 1
		if err := r.receiveLoop(ctx,
			&bnOptions{
				StartHeight:     opts.Height,
				VerifierOptions: r.opts.Verifier,
				Concurrency:     r.opts.SyncConcurrency,
			},
			func(v *BlockNotification) error {
				//r.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")

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
						default:
							r.log.WithFields(log.Fields{
								"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
							}).Warn("Default: invalid event seq")
						}
					}
					receipt.Events = events
				}

				if len(receipts) > 0 {
					for _, sev := range v.Receipts {
						sinkChan <- &chain.SubscribedEvent{Res: sev, ChainName: chain.HMNY}
					}
				}
				lastHeight++
				return nil
			}); err != nil {
			r.log.Errorf("receiveLoop terminated: %+v", err)
			_errCh <- err
		}
	}()

	return nil
}

func (r *receiver) getFilteredReceipts(v *BlockNotification) []*LogResult {
	const (
		TransferStartSignature         = "0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a" //"TransferStart(address,string,uint256,(string,uint256,uint256)[])" //
		TransferEndSignature           = "0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2" //"TransferEnd(address,uint256,uint256,string)"                      //
		TransferReceivedSignature      = "0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680" //"TransferReceived(string,address,uint256,(string,uint256)[])"      //
		TransferReceivedSignatureToken = "0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c" //"TransferReceived(string,address,uint256,(string,uint256,uint256)[])" //
	)
	signatureMap := map[string]string{
		"0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a": "TransferStart",
		"0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2": "TransferEnd",
		"0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680": "TransferReceived",
		"0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c": "TransferReceived",
	}

	newResults := []*LogResult{}
	for _, receipt := range v.Receipts {
		for _, log := range receipt.Logs {
			var newTopic *common.Hash
			for _, topic := range log.Topics {
				if topic == common.HexToHash(TransferStartSignature) ||
					topic == common.HexToHash(TransferReceivedSignature) ||
					topic == common.HexToHash(TransferReceivedSignatureToken) ||
					topic == common.HexToHash(TransferEndSignature) {
					newTopic = &topic
					break
				}
			}
			if newTopic != nil {

				if res, err := decodeLogData(log.Data, signatureMap[newTopic.String()]); err == nil && res != nil {
					newResults = append(newResults, &LogResult{
						TxHash:   log.TxHash,
						LogIndex: log.Index,
						Address:  log.Address,
						Topic:    signatureMap[newTopic.String()],
						Logs:     res,
					})
				} else if err != nil {
					r.log.Error(err)
				} else if res == nil {
					r.log.Error("Returned nil interface")
				}
			}
		}
	}
	for i, r := range newResults {
		fmt.Println("New ", i, "  ", *r)
	}
	return newResults
}

func decodeLogData(data []byte, topicType string) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty Log Data input to decode")
	}

	abi, err := abi.JSON(strings.NewReader(bshPeripherABI))
	if err != nil {
		return nil, err
	}

	if topicType == "TransferStart" {
		var ev TransferStart
		err = abi.UnpackIntoInterface(&ev, topicType, data)
		if err != nil {
			return nil, err
		}
		return ev, nil
	} else if topicType == "TransferEnd" {
		var ev TransferEnd
		err = abi.UnpackIntoInterface(&ev, topicType, data)
		if err != nil {
			return nil, err
		}
		return ev, nil
	} else if topicType == "TransferReceived" {
		var ev TransferReceived
		err = abi.UnpackIntoInterface(&ev, topicType, data)
		if err != nil {
			return nil, err
		}
		return ev, nil
	}
	return nil, errors.New("Doesn't match any signature")
}
