package near

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"math/rand"
)

type receiver struct {
	clients     []*Client
	source      chain.BTPAddress
	destination chain.BTPAddress
	log         log.Logger
	options     struct{}
}

func NewReceiver(src, dst chain.BTPAddress, urls []string, opt map[string]interface{}, logger log.Logger) (chain.Receiver, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	r := &receiver{
		clients:     newClients(urls, logger),
		source:      src,
		destination: dst,
		log:         logger,
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

func (r *receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	return nil, nil
}

func (r *receiver) client() *Client {
	return r.clients[rand.Intn(len(r.clients))]
}
