package bsc

import "github.com/icon-project/icon-bridge/cmd/btpsimple/relay"

func init() {
	relay.Senders["bsc"] = NewSender
	relay.Receivers["bsc"] = NewReceiver
}
