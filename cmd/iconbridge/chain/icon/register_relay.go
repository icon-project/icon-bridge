package icon

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["icon"] = NewSender
	relay.Receivers["icon"] = NewReceiver
}
