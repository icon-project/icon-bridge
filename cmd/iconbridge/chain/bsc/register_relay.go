package bsc

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["bsc"] = NewSender
	relay.Receivers["bsc"] = NewReceiver
}
