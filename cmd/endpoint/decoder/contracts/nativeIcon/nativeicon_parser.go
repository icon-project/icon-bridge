package nativeIcon

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
)

func parseTransferStart(log icon.TxnEventLog) (*NativeIconTransferStart, error) {
	if len(log.Data) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	data := log.Data
	res := &[]AssetTransferDetails{}
	err := rlpDecodeHex(data[len(data)-1], res)
	if err != nil {
		return nil, err
	}
	sn := new(big.Int)
	if strings.HasPrefix(data[1], "0x") {
		data[1] = data[1][2:]
	}
	sn.SetString(data[1], 16)
	ts := &NativeIconTransferStart{
		From:   log.Indexed[1],
		To:     data[0],
		Sn:     sn,
		Assets: *res,
	}
	return ts, nil
}

func parseTransferReceived(log icon.TxnEventLog) (*NativeIconTransferReceived, error) {
	if len(log.Data) != 2 || len(log.Indexed) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	data := log.Data
	res := &[]AssetDetails{}
	err := rlpDecodeHex(data[len(data)-1], res)
	if err != nil {
		return nil, err
	}
	sn := new(big.Int)
	if strings.HasPrefix(data[0], "0x") {
		data[0] = data[0][2:]
	}
	sn.SetString(data[0], 16)
	ts := &NativeIconTransferReceived{
		From:   log.Indexed[1],
		To:     log.Indexed[2],
		Sn:     sn,
		Assets: *res,
	}
	return ts, nil
}

func parseTransferEnd(log icon.TxnEventLog) (*NativeIconTransferEnd, error) {
	data := log.Data
	sn := new(big.Int)
	if strings.HasPrefix(data[0], "0x") {
		data[0] = data[0][2:]
	}
	sn.SetString(data[0], 16)

	cd := new(big.Int)
	if strings.HasPrefix(data[1], "0x") {
		data[1] = data[1][2:]
	}
	cd.SetString(data[1], 16)
	te := &NativeIconTransferEnd{
		From: log.Indexed[1],
		Sn:   sn,
		Code: cd,
	}
	return te, nil
}

type NativeIconTransferStart struct {
	From   string
	To     string
	Sn     *big.Int
	Assets []AssetTransferDetails
}

type NativeIconTransferReceived struct {
	From   string
	To     string
	Sn     *big.Int
	Assets []AssetDetails
}

type AssetTransferDetails struct {
	Name  string
	Value *big.Int
	Fee   *big.Int
}

type AssetDetails struct {
	Name  string
	Value *big.Int
}

type NativeIconTransferEnd struct {
	From string
	Sn   *big.Int
	Code *big.Int
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
