package bmcIcon

import (
	"errors"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type bmcIconContract struct {
	name contracts.ContractName
}

func (ti *bmcIconContract) GetName() contracts.ContractName {
	return ti.name
}

func NewContract(name contracts.ContractName) (contracts.Contract, error) {
	tic := &bmcIconContract{name: name}
	return tic, nil
}

func (ti *bmcIconContract) Decode(li interface{}) (res map[string]interface{}, err error) {
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
	if sign == "Message" {
		res[sign], err = parseMessage(log)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
