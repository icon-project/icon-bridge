package bsc

import "github.com/icon-project/icon-bridge/bmr/cmd/iconbridge/relay"

func init() {
	relay.Senders["bsc"] = NewSender
	relay.Receivers["bsc"] = NewReceiver
}
