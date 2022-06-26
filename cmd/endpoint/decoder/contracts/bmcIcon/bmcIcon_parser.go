package bmcIcon

import (
	"errors"
	"math/big"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
)

func parseMessage(log icon.TxnEventLog) (*MessageIconTransfer, error) {
	if len(log.Indexed) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	vstr := log.Indexed[2]
	if strings.HasPrefix(vstr, "0x") {
		vstr = vstr[2:]
	}
	value := new(big.Int)
	value.SetString(vstr, 16)
	ts := &MessageIconTransfer{
		Next: log.Indexed[1],
		Seq:  value,
	}
	return ts, nil
}

type MessageIconTransfer struct {
	Next string
	Seq  *big.Int
}
