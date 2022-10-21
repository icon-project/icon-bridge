package near

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/reactivex/rxgo/v2"
	"math/rand"
)

type ReceiverOptions struct {
	SyncConcurrency uint `json:"syncConcurrency"`
}
type Receiver struct {
	clients     []IClient
	source      chain.BTPAddress
	destination chain.BTPAddress
	logger      log.Logger
	options     ReceiverOptions
}

func NewReceiver(source, destination chain.BTPAddress, urls []string, options json.RawMessage, logger log.Logger) (chain.Receiver, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	clients, err := newClients(urls, logger)
	if err != nil {
		return nil, err
	}
	
	r := &Receiver{
		clients,
		source,
		destination,
		logger,
		ReceiverOptions{},
	}

	if err := json.Unmarshal(options, &r.options); err != nil {
		logger.Panicf("fail to unmarshal opt:%#v err:%+v", options, err)
		return nil, err
	}

	return r, nil
}

func newMockReceiver(source, destination chain.BTPAddress, client *Client, urls []string, _ json.RawMessage, logger log.Logger) (*Receiver, error) {
	clients := make([]IClient, 0)
	clients = append(clients, client)
	receiver := &Receiver{
		clients:     clients,
		source:      source,
		destination: destination,
		logger:      logger,
	}

	return receiver, nil
}

func (r *Receiver) receiveBlocks(height uint64, source string, processBlockNotification func(blockNotification *types.BlockNotification)) error {
	return r.client().MonitorBlocks(height, r.source.ContractAddress(), r.options.SyncConcurrency, func(observable rxgo.Observable) error {
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

		if err := r.receiveBlocks(opts.Height, r.source.ContractAddress(), func(blockNotification *types.BlockNotification) {
			r.logger.WithFields(log.Fields{"height": blockNotification.Block().Height()}).Debug("block notification")
			receipts := blockNotification.Receipts()

			for _, receipt := range receipts {
				events := receipt.Events[:0]
				for _, event := range receipt.Events {
					switch {

					case event.Sequence == opts.Seq:
						events = append(events, event)
						opts.Seq++

					case event.Sequence > opts.Seq:
						r.logger.WithFields(log.Fields{
							"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
						}).Error("invalid event seq")

						_errCh <- fmt.Errorf("invalid event seq")
					}

					receipt.Events = events
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

	return errCh, nil
}

func (r *Receiver) client() IClient {
	return r.clients[rand.Intn(len(r.clients))]
}

func (r *Receiver) StopReceivingBlocks() {
	r.client().CloseMonitor()
}
