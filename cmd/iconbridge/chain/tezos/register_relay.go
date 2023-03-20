package tezos 

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["tz"] = NewSender
	relay.Receivers["tz"] = NewReceiver
}
