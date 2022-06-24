package icon

import (
	"context"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewSubscriptionAPI(l log.Logger, cfg chain.SubscriberConfig, endpoint string) (chain.SubscritionAPI, error) {

	recv, err := NewReceiver(cfg.Src, cfg.Dst, []string{endpoint}, cfg.Opts, l)
	if err != nil {
		panic(err)
	}
	return recv, nil
}

func (r *receiver) Start(ctx context.Context, sinkChan chan<- *chain.SubscribedEvent, errChan chan<- error) error {
	var seq uint64 = 0
	var height uint64 = 1211 // used to be fetched by BMC Status
	err := r.Subscribe(ctx, sinkChan, errChan,
		chain.SubscribeOptions{Height: height, Seq: seq})
	if err != nil {
		return err
	}
	return nil
}
