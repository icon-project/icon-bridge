package relay

import (
	"context"
	"runtime/debug"

	"github.com/icon-project/btp/cmd/btpsimple/chain"
	"github.com/icon-project/btp/cmd/btpsimple/chain/bsc"
	"github.com/icon-project/btp/cmd/btpsimple/chain/icon"
	"github.com/icon-project/btp/common/log"
	"github.com/icon-project/btp/common/wallet"
)

func NewMultiRelay(cfg *Config, l log.Logger) (Relay, error) {
	mr := &multiRelay{log: l}

	for _, rc := range cfg.Relays {

		var dst chain.Sender
		var src chain.Receiver

		w, err := rc.Dst.Wallet()
		if err != nil {
			return nil, err
		}

		l := l.WithFields(log.Fields{
			log.FieldKeyModule: rc.Name,
			log.FieldKeyWallet: w.Address(),
		})

		blockChain := rc.Dst.Address.BlockChain()
		switch blockChain {
		case "icon":
			if dst, err = icon.NewSender(
				rc.Src.Address,
				rc.Dst.Address,
				rc.Dst.Endpoint,
				w,
				rc.Dst.Options,
				l.WithFields(log.Fields{
					log.FieldKeyPrefix: "tx_",
					log.FieldKeyChain:  blockChain,
				}),
			); err != nil {
				return nil, err
			}
		case "bsc":
			if dst, err = bsc.NewSender(
				rc.Src.Address,
				rc.Dst.Address,
				rc.Dst.Endpoint,
				w.(*wallet.EvmWallet),
				rc.Dst.Options,
				l.WithFields(log.Fields{
					log.FieldKeyPrefix: "tx_",
					log.FieldKeyChain:  blockChain,
				}),
			); err != nil {
				return nil, err
			}
		}

		blockChain = rc.Src.Address.BlockChain()
		switch blockChain {
		case "icon":
			if src, err = icon.NewReceiver(
				rc.Src.Address,
				rc.Dst.Address,
				rc.Src.Endpoint,
				rc.Src.Options,
				l.WithFields(log.Fields{
					log.FieldKeyPrefix: "rx_",
					log.FieldKeyChain:  blockChain,
				}),
			); err != nil {
				return nil, err
			}
		case "bsc":
			if src, err = bsc.NewReceiver(
				rc.Src.Address,
				rc.Dst.Address,
				rc.Src.Endpoint,
				rc.Src.Options,
				l.WithFields(log.Fields{
					log.FieldKeyPrefix: "rx_",
					log.FieldKeyChain:  blockChain,
				}),
			); err != nil {
				return nil, err
			}
		}

		relay, err := NewRelay(rc, src, dst, l.WithFields(log.Fields{log.FieldKeyChain: "relay"}))
		if err != nil {
			return nil, err
		}
		mr.relays = append(mr.relays, relay)

	}

	return mr, nil
}

type multiRelay struct {
	log    log.Logger
	relays []Relay
}

func (mr *multiRelay) Start(ctx context.Context) error {
	rch := make(chan Relay, len(mr.relays))
	for _, relay := range mr.relays {
		rch <- relay
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case r := <-rch:
			go func(relay Relay) {
				defer func() {
					if r := recover(); r != nil {
						debug.PrintStack()
						rch <- relay
					}
				}()
				if relay.Start(ctx) != nil {
					rch <- relay
				}
			}(r)
		}
	}
}
