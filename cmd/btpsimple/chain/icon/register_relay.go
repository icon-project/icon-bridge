package icon

import "github.com/icon-project/icon-bridge/cmd/btpsimple/relay"

func init() {
	relay.Senders["icon"] = NewSender
	relay.Receivers["icon"] = NewReceiver
}
