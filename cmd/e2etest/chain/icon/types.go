package icon

import "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"

type eventLogRawFilter struct {
	addr      []byte
	signature []byte
	next      []byte
	seq       uint64
}
type TxnEventLog struct {
	Addr    icon.Address `json:"scoreAddress"`
	Indexed []string     `json:"indexed"`
	Data    []string     `json:"data"`
}
