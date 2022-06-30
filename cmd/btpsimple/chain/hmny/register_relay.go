//go:build hmny

package hmny

import "github.com/icon-project/icon-bridge/cmd/btpsimple/relay"

func init() {
	relay.Senders["hmny"] = NewSender
	relay.Receivers["hmny"] = NewReceiver
}
