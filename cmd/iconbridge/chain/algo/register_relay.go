package algo

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["algo"] = NewSender
	relay.Receivers["algo"] = NewReceiver
}
