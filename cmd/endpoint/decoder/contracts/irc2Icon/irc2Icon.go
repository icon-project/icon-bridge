package irc2Icon

import (
	"errors"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type irc2IconContract struct {
	name contracts.ContractName
}

func (ti *irc2IconContract) GetName() contracts.ContractName {
	return ti.name
}

func NewContract(name contracts.ContractName) (contracts.Contract, error) {
	tic := &irc2IconContract{name: name}
	return tic, nil
}

func (ti *irc2IconContract) Decode(li interface{}) (res map[string]interface{}, err error) {
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
	if sign == "Transfer" {
		res[sign], err = parseTransfer(log)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
