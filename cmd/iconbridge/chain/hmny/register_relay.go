//go:build hmny
// +build hmny

package hmny

import "github.com/icon-project/icon-bridge/cmd/iconbridge/relay"

func init() {
	relay.Senders["hmny"] = NewSender
	relay.Receivers["hmny"] = NewReceiver
}
