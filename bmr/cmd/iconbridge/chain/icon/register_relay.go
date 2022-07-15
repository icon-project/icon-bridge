package icon

import "github.com/icon-project/icon-bridge/bmr/cmd/iconbridge/relay"

func init() {
	relay.Senders["icon"] = NewSender
	relay.Receivers["icon"] = NewReceiver
}
