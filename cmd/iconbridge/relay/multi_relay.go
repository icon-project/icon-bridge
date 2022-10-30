package relay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

type NewSenderFunc func(
	src, dst chain.BTPAddress, urls []string, w wallet.Wallet,
	opts json.RawMessage, l log.Logger) (chain.Sender, error)

type NewReceiverFunc func(
	src, dst chain.BTPAddress, urls []string,
	opts json.RawMessage, l log.Logger) (chain.Receiver, error)

// type NewSenderFunc func(
// 	src, dst chain.BTPAddress, urls []string, w wallet.Wallet,
// 	opts map[string]interface{}, l log.Logger) (chain.Sender, error)

// type NewReceiverFunc func(
// 	src, dst chain.BTPAddress, urls []string,
// 	opts map[string]interface{}, l log.Logger) (chain.Receiver, error)

var (
	Senders   = map[string]NewSenderFunc{}
	Receivers = map[string]NewReceiverFunc{}
)

func NewMultiRelay(cfg *Config, l log.Logger) (Relay, error) {
	mr := &multiRelay{log: l}

	for _, rc := range cfg.Relays {

		var dst chain.Sender
		var src chain.Receiver

		w, err := rc.Dst.Wallet()
		if err != nil {
			return nil, fmt.Errorf("dst.wallet chain %v err %v", rc.Name, err)
		}
		chainName := rc.Dst.Address.BlockChain()
		srvName := "BMR-"
		if strings.ToUpper(chainName) == "ICON" {
			srvName += strings.ToUpper(rc.Src.Address.BlockChain())
		} else {
			srvName += strings.ToUpper(chainName)
		}
		l := l.WithFields(log.Fields{
			log.FieldKeyModule:  rc.Name,
			log.FieldKeyWallet:  w.Address(),
			log.FieldKeyService: srvName,
		})

		if sender, ok := Senders[chainName]; ok {
			if dst, err = sender(
				rc.Src.Address,
				rc.Dst.Address,
				rc.Dst.Endpoint,
				w,
				rc.Dst.Options,
				l.WithFields(log.Fields{
					log.FieldKeyPrefix: "tx_",
					log.FieldKeyChain:  chainName,
				})); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported blockchain: sender=%s", chainName)
		}

		chainName = rc.Src.Address.BlockChain()
		if receiver, ok := Receivers[chainName]; ok {
			if src, err = receiver(
				rc.Src.Address,
				rc.Dst.Address,
				rc.Src.Endpoint,
				rc.Src.Options,
				l.WithFields(log.Fields{
					log.FieldKeyPrefix: "rx_",
					log.FieldKeyChain:  chainName,
				}),
			); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported blockchain: receiver=%s", chainName)
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
				if err := relay.Start(ctx); err != nil {
					if !errors.Is(err, context.Canceled) {
						mr.log.Errorf("%v", err)
						mr.log.Info("restarting relay in 5s...")
						time.Sleep(5 * time.Second)
						rch <- relay
					}
				}
			}(r)
		}
	}
}
