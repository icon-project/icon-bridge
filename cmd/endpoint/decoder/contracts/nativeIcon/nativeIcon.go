package nativeIcon

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type nativeIconContract struct {
}

func (ti *nativeIconContract) GetName() contracts.ContractName {
	return contracts.NativeIcon
}

func NewContract(cAddr common.Address) (contracts.Contract, error) {
	tic := &nativeIconContract{}
	return tic, nil
}

func (ti *nativeIconContract) Decode(li interface{}) (res map[string]interface{}, err error) {
	log, ok := li.(icon.TxnEventLog)
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
