package nativeIcon

import (
	"errors"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type nativeIconContract struct {
	name contracts.ContractName
}

func (ti *nativeIconContract) GetName() contracts.ContractName {
	return ti.name
}

func NewContract(name contracts.ContractName) (contracts.Contract, error) {
	tic := &nativeIconContract{name: name}
	return tic, nil
}

func (ti *nativeIconContract) Decode(li interface{}) (res map[string]interface{}, err error) {
	log, ok := li.(*icon.TxnEventLog)
	if !ok {
		return nil, errors.New("Log of wrong type. Expected icon.TxnLog")
	}
	if len(log.Indexed) == 0 {
		err = errors.New("log.Indexed is of size zero")
		return nil, err
	}
	eventName := strings.Split(log.Indexed[0], "(")
	sign := strings.TrimSpace(eventName[0])
	res = map[string]interface{}{}
	if sign == "TransferStart" {
		res[sign], err = parseTransferStart(log)
		if err != nil {
			return nil, err
		}
	} else if sign == "TransferEnd" {
		res[sign], err = parseTransferEnd(log)
		if err != nil {
			return nil, err
		}
	} else if sign == "TransferReceived" {
		res[sign], err = parseTransferReceived(log)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
