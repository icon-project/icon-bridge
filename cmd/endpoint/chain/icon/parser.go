package icon

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
)

type parser struct {
	addressToContractName map[string]chain.ContractName
}

func NewParser(nameToAddr map[chain.ContractName]string) (*parser, error) {
	addrToName := map[string]chain.ContractName{}
	for name, addr := range nameToAddr {
		addrToName[addr] = name
	}
	return &parser{addressToContractName: addrToName}, nil
}

func (p *parser) Parse(log *TxnEventLog) (resLog interface{}, eventType chain.EventLogType, err error) {
	cName, ok := p.addressToContractName[string(log.Addr)]
	if !ok {
		err = errors.New("Couldn't find contract matching the log " + string(log.Addr))
		return
	}
	eventName := strings.Split(log.Indexed[0], "(")
	eventType = chain.EventLogType(strings.TrimSpace(eventName[0]))
	if cName == chain.NativeBSHIcon {
		if eventType == chain.TransferStart {
			resLog, err = parseTransferStartNativeCoin(log)
		} else if eventType == chain.TransferReceived {
			resLog, err = parseTransferReceivedNativeCoin(log)
		} else if eventType == chain.TransferEnd {
			resLog, err = parseTransferEndNativeCoin(log)
		} else {
			err = errors.New("No matching signature ")
		}
	} else if cName == chain.TokenBSHIcon {
		if eventType == chain.TransferStart {
			resLog, err = parseTransferStartToken(log)
		} else if eventType == chain.TransferReceived {
			resLog, err = parseTransferReceivedToken(log)
		} else if eventType == chain.TransferEnd {
			resLog, err = parseTransferEndToken(log)
		} else {
			err = errors.New("No matching signature ")
		}
	} else {
		err = errors.New("Contract not amongst processed ones")
	}
	return
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

func parseTransferStartNativeCoin(log *TxnEventLog) (*chain.TransferStartEvent, error) {
	if len(log.Data) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	data := log.Data
	res := &[]chain.AssetTransferDetails{}
	err := rlpDecodeHex(data[len(data)-1], res)
	if err != nil {
		return nil, err
	}
	sn := new(big.Int)
	if strings.HasPrefix(data[1], "0x") {
		data[1] = data[1][2:]
	}
	sn.SetString(data[1], 16)
	ts := &chain.TransferStartEvent{
		From:   log.Indexed[1],
		To:     data[0],
		Sn:     sn,
		Assets: *res,
	}
	return ts, nil
}

func parseTransferReceivedNativeCoin(log *TxnEventLog) (*chain.TransferReceivedEvent, error) {
	if len(log.Data) != 2 || len(log.Indexed) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	data := log.Data
	res := &[]chain.AssetDetails{}
	err := rlpDecodeHex(data[len(data)-1], res)
	if err != nil {
		return nil, err
	}
	sn := new(big.Int)
	if strings.HasPrefix(data[0], "0x") {
		data[0] = data[0][2:]
	}
	sn.SetString(data[0], 16)
	newAssetDetails := make([]chain.AssetTransferDetails, len(*res))
	for i, v := range *res {
		newAssetDetails[i].Name = v.Name
		newAssetDetails[i].Value = v.Value
	}
	ts := &chain.TransferReceivedEvent{
		From:   log.Indexed[1],
		To:     log.Indexed[2],
		Sn:     sn,
		Assets: newAssetDetails,
	}
	return ts, nil
}

func parseTransferEndNativeCoin(log *TxnEventLog) (*chain.TransferEndEvent, error) {
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
	te := &chain.TransferEndEvent{
		From: log.Indexed[1],
		Sn:   sn,
		Code: cd,
	}
	return te, nil
}

func parseTransferStartToken(log *TxnEventLog) (*chain.TransferStartEvent, error) {
	if len(log.Data) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	data := log.Data
	res := &[]chain.AssetTransferDetails{}
	err := rlpDecodeHex(data[len(data)-1], res)
	if err != nil {
		return nil, err
	}
	sn := new(big.Int)
	if strings.HasPrefix(data[1], "0x") {
		data[1] = data[1][2:]
	}
	sn.SetString(data[1], 16)
	ts := &chain.TransferStartEvent{
		From:   log.Indexed[1],
		To:     data[0],
		Sn:     sn,
		Assets: *res,
	}
	return ts, nil
}

func parseTransferReceivedToken(log *TxnEventLog) (*chain.TransferReceivedEvent, error) {
	if len(log.Data) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
	data := log.Data
	res := &[]chain.AssetTransferDetails{}
	err := rlpDecodeHex(data[len(data)-1], res)
	if err != nil {
		return nil, err
	}
	sn := new(big.Int)
	if strings.HasPrefix(data[1], "0x") {
		data[1] = data[1][2:]
	}
	sn.SetString(data[1], 16)
	ts := &chain.TransferReceivedEvent{
		From:   log.Indexed[1],
		To:     data[0],
		Sn:     sn,
		Assets: *res,
	}
	return ts, nil
}

func parseTransferEndToken(log *TxnEventLog) (*chain.TransferEndEvent, error) {
	if len(log.Data) != 3 {
		return nil, errors.New("Unexpected length of log.Data")
	}
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
	te := &chain.TransferEndEvent{
		From: log.Indexed[1],
		Sn:   sn,
		Code: cd,
	}
	return te, nil
}
