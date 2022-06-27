package icon

import (
	"context"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewSubscriptionAPI(l log.Logger, cfg chain.SubscriberConfig, endpoint string) (chain.SubscriptionAPI, error) {

	recv, err := NewReceiver(cfg.Src, cfg.Dst, []string{endpoint}, cfg.Opts, l)
	if err != nil {
		panic(err)
	}
	return recv, nil
}

func (r *receiver) Start(ctx context.Context) error {
	var seq uint64 = 0
	var height uint64 = 42000 // used to be fetched by BMC Status
	err := r.Subscribe(ctx, chain.SubscribeOptions{Height: height, Seq: seq})
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
