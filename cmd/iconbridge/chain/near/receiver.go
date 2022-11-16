package near

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sync"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/reactivex/rxgo/v2"
)

type ReceiverOptions struct {
	SyncConcurrency uint `json:"syncConcurrency"`
}

type ReceiverConfig struct {
	source      chain.BTPAddress
	destination chain.BTPAddress
	options     types.ReceiverOptions
}

type Receiver struct {
	clients     []IClient
	source      chain.BTPAddress
	destination chain.BTPAddress
	logger      log.Logger
	verifier    *Verifier
	options     types.ReceiverOptions
}

func receiverFactory(source, destination chain.BTPAddress, urls []string, opt json.RawMessage, logger log.Logger) (chain.Receiver, error) {
	var options types.ReceiverOptions
	clients, err := newClients(urls, logger)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(opt, &options); err != nil {
		logger.Panicf("fail to unmarshal options:%#v err:%+v", opt, err)
		return nil, err
	}

	return NewReceiver(ReceiverConfig{source, destination, options}, logger, clients...)
}

func NewReceiver(config ReceiverConfig, logger log.Logger, clients ...IClient) (*Receiver, error) {
	if len(clients) == 0 {
		return nil, fmt.Errorf("nil clients")
	}
	
	r := &Receiver{
		clients:     clients,
		logger:      logger,
		source:      config.source,
		destination: config.destination,
		options:     config.options,
	}

	return r, nil
}

func (r *Receiver) MapReceipts(height uint64, source string, observable rxgo.Observable) rxgo.Observable {
	return observable.Map(
		r.client().FetchReceipts,
		rxgo.WithPool(r.options.SyncConcurrency),
		rxgo.WithContext(context.WithValue(context.Background(), Source{}, source)),
	).Serialize(
		int(height),
		r.client().SerializeBlocks,
	).Filter(
		r.client().FilterUnknownBlocks,
	)
}

func (r *Receiver) ReceiveBlocks(height uint64, source string, processBlockNotification func(blockNotification *types.BlockNotification)) error {

	return r.client().MonitorBlocks(height, math.MaxInt64, r.options.SyncConcurrency, func(observable rxgo.Observable) error {
		result := r.MapReceipts(height, source, observable).Scan(
			func(_ context.Context, acc interface{}, bn interface{}) (interface{}, error) {
				blockNotification, _ := bn.(*types.BlockNotification)

				if r.verifier != nil {
					if err := r.verifier.ValidateHeader(blockNotification); err != nil {
						return nil, err
					}
				}

				r.logger.WithFields(log.Fields{"height": blockNotification.Block().Height()}).Debug("block notification")

				return blockNotification, nil
			},
		).Observe()

		for item := range result {
			if err := item.E; err != nil {
				return err
			}

			if bn, ok := item.V.(*types.BlockNotification); ok {
				processBlockNotification(bn)
			} else {
				return fmt.Errorf("expected *types.BlockNotification but got: %v", reflect.TypeOf(item.V))
			}
		}

		return nil
	})
}

func (r *Receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	opts.Seq++
	_errCh := make(chan error)

	if r.options.Verifier != nil {
		r.verifier, err = NewVerifier(
			r.options.Verifier.BlockHeight,
			r.options.Verifier.PreviousBlockHash,
			r.options.Verifier.CurrentEpochId,
			r.options.Verifier.NextEpochId,
			r.options.Verifier.CurrentBpsHash,
			r.options.Verifier.NextBpsHash,
			r.options.SyncConcurrency,
			r.client(),
		)
	}

	if err != nil {
		return _errCh, err
	}

	go func() {
		defer close(_errCh)

		if r.verifier != nil {
			wg := new(sync.WaitGroup)
			wg.Add(1)

			r.logger.WithFields(log.Fields{"start": r.options.Verifier.BlockHeight, "target": opts.Height - 1}).Debug("syncing verifier head")
			if err := r.verifier.SyncHeader(wg, opts.Height-1); err != nil {
				_errCh <- err
			}

			wg.Wait()
			r.logger.Debug("syncing complete")
		}

		if err := r.ReceiveBlocks(opts.Height, r.source.ContractAddress(), func(blockNotification *types.BlockNotification) {
			receipts := make([]*chain.Receipt, 0)

			for _, receipt := range blockNotification.Receipts() {
				events := receipt.Events[:0]
				for _, event := range receipt.Events {
					switch {
					case event.Sequence == opts.Seq && event.Next == r.destination:
						events = append(events, event)
						opts.Seq++

					case event.Sequence > opts.Seq && event.Next == r.destination:
						r.logger.WithFields(log.Fields{
							"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
						}).Error("invalid event seq")

						_errCh <- fmt.Errorf("invalid event seq")

						return
					}

					receipt.Events = events
				}

				if len(events) > 0 {
					receipts = append(receipts, receipt)
				}
			}

			if len(receipts) > 0 {
				msgCh <- &chain.Message{
					From:     r.source,
					Receipts: receipts,
				}
			}
		}); err != nil {
			_errCh <- err
		}
	}()

	return _errCh, nil
}

func (r *Receiver) client() IClient {
	return r.clients[rand.Intn(len(r.clients))]
}

func (r *Receiver) StopReceivingBlocks() {
	r.client().CloseMonitor()
}
