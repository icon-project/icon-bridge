package irc2Icon

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
)

func parseTransfer(log icon.TxnEventLog) (*Irc2IconTransfer, error) {
	if len(log.Indexed) != 4 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	vstr := log.Indexed[3]
	if strings.HasPrefix(vstr, "0x") {
		vstr = vstr[2:]
	}
	value := new(big.Int)
	value.SetString(vstr, 16)
	// err := rlpDecodeHex(log.Indexed[3], value)
	// if err != nil {
	// 	return nil, err
	// }
	ts := &Irc2IconTransfer{
		From:  log.Indexed[1],
		To:    log.Indexed[2],
		Value: value,
	}
	return ts, nil
}

type Irc2IconTransfer struct {
	From  string
	To    string
	Value *big.Int
}

func rlpDecodeHex(str string, out interface{}) error {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	input, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	err = rlp.Decode(bytes.NewReader(input), out)
	if err != nil {
		return err
	}
	return nil
}
