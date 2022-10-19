package near

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/reactivex/rxgo/v2"
)

type ReceiverConfig struct {
	source      chain.BTPAddress
	destination chain.BTPAddress
	options     json.RawMessage
}

type Receiver struct {
	clients     []IClient
	source      chain.BTPAddress
	destination chain.BTPAddress
	logger      log.Logger
	options     struct {
		SyncConcurrency uint `json:"syncConcurrency"`
	}
}

func receiverFactory(source, destination chain.BTPAddress, urls []string, options json.RawMessage, logger log.Logger) (chain.Receiver, error) {
	clients, err := newClients(urls, logger)
	if err != nil {
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
	}

	if err := json.Unmarshal(config.options, &r.options); err != nil && config.options != nil {
		logger.Panicf("fail to unmarshal opt:%#v err:%+v", config.options, err)
		return nil, err
	}

	return r, nil
}

func (r *Receiver) ReceiveBlocks(height uint64, source string, processBlockNotification func(blockNotification *types.BlockNotification)) error {
	return r.client().MonitorBlocks(height, source, r.options.SyncConcurrency, func(observable rxgo.Observable) error {
		result := observable.Observe()

		for item := range result {
			if err := item.E; err != nil {
				return err
			}

			bn, _ := item.V.(*types.BlockNotification)

			if *bn.Block().Hash() != [32]byte{} {
				processBlockNotification(bn)
			}
		}
		return nil
	}, func() IClient {
		return r.client()
	})
}

func (r *Receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	opts.Seq++
	_errCh := make(chan error)

	go func() {
		defer close(_errCh)

		if err := r.ReceiveBlocks(opts.Height, r.source.ContractAddress(), func(blockNotification *types.BlockNotification) {
			r.logger.WithFields(log.Fields{"height": blockNotification.Block().Height()}).Debug("block notification")
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
