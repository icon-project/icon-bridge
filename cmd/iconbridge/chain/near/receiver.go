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

type Receiver struct {
	clients     []*Client
	source      chain.BTPAddress
	destination chain.BTPAddress
	logger      log.Logger
	options     struct{}
}

func NewReceiver(src, dst chain.BTPAddress, urls []string, opt map[string]interface{}, logger log.Logger) (chain.Receiver, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	r := &Receiver{
		clients:     newClients(urls, logger),
		source:      src,
		destination: dst,
		logger:      logger,
	}
	b, err := json.Marshal(opt)
	if err != nil {
		logger.Panicf("fail to marshal opt:%#v err:%+v", opt, err)
	}
	if err = json.Unmarshal(b, &r.options); err != nil {
		logger.Panicf("fail to unmarshal opt:%#v err:%+v", opt, err)
	}

	if err != nil {
		return nil, err
	}
	return r, nil
}

func newMockReceiver(source, destination chain.BTPAddress, client *Client, urls []string, _ map[string]interface{}, logger log.Logger) (*Receiver, error) {
	clients := make([]*Client, 0)
	clients = append(clients, client)
	receiver := &Receiver{
		clients:     clients,
		source:      source,
		destination: destination,
		logger:      logger,
	}

	return receiver, nil
}

func (r *Receiver) receiveBlocks(height uint64, processBlock func(block *types.Block)) error {
	return r.client().MonitorBlocks(height, func(observable rxgo.Observable) error {
		result := observable.Observe()

		for item := range result {
			if err := item.E; err != nil {
				return err
			}

			block, _ := item.V.(types.Block)
			processBlock(&block)
		}
		return nil
	})
}

func (r *Receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	opts.Seq++
	_errCh := make(chan error)

	go func() {
		defer close(_errCh)

		if err := r.receiveBlocks(opts.Height, func(block *types.Block) {
			receipts, err := r.client().GetReceipts(block, r.source.ContractAddress())
			if err != nil {
				_errCh <- err
			}

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

func (r *Receiver) client() *Client {
	return r.clients[rand.Intn(len(r.clients))]
}

func (r *Receiver) StopReceivingBlocks() {
	r.client().CloseMonitor()
}
