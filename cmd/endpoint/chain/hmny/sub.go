package hmny

import (
	"context"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewSubscriptionAPI(l log.Logger, cfg chain.SubscriberConfig, endpoint string) (chain.SubscritionAPI, error) {
	rx, err := NewReceiver(cfg.Src, cfg.Dst, []string{endpoint}, cfg.Opts, l)
	if err != nil {
		return nil, err
	}
	return rx, nil
}

func (r *receiver) Start(ctx context.Context, sinkChan chan<- *chain.SubscribedEvent, errChan chan<- error) error {
	err := r.Subscribe(ctx, sinkChan, errChan,
		chain.SubscribeOptions{
			Seq:    138,
			Height: 29076,
		})
	if err != nil {
		return err
	}
	return nil
}
