package hmny

import (
	"context"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewSubscriptionAPI(l log.Logger, cfg chain.SubscriberConfig, endpoint string) (chain.SubscriptionAPI, error) {
	rx, err := NewReceiver(cfg.Src, cfg.Dst, []string{endpoint}, cfg.Opts, l)
	if err != nil {
		return nil, err
	}
	return rx, nil
}

func (r *receiver) Start(ctx context.Context) error {
	err := r.Subscribe(ctx,
		chain.SubscribeOptions{
			Seq:    0,
			Height: 69,
		})
	if err != nil {
		return err
	}
	return nil
}
func (r *receiver) OutputChan() <-chan *chain.SubscribedEvent {
	return r.sinkChan
}

func (r *receiver) ErrChan() <-chan error {
	return r.errChan
}
