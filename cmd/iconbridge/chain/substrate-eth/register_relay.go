package substrate_eth

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["snow"] = NewSender
	relay.Receivers["snow"] = NewReceiver
}
