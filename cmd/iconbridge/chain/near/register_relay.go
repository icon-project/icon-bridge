package near

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["near"] = senderFactory
	relay.Receivers["near"] = receiverFactory
}
