package tezos 

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["tezos"] = NewSender
	relay.Receivers["tezos"] = NewReceiver
}
