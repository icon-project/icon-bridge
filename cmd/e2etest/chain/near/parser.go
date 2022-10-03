package near

import (
	"encoding/json"
	"fmt"
	// "strings"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type AssetDetails struct{}

type Event struct {
	name            chain.EventLogType `json:"event"`
	code            types.BigInt       `json:"code"`
	senderAddress   types.AccountId    `json:"sender_address,omitempty"`
	serialNumber    types.BigInt       `json:"serial_number"`
	receiverAddress chain.BTPAddress   `json:"receiver_address,omitempty"`
	assets          []AssetDetails     `json:"assets,omitempty"`
	tokenName       string             `json:"token_name,omitempty"`
	message         string             `json:"message,omitempty"`
	tokenAccount    types.AccountId    `json:"token_account,omitempty"`
}

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

func (p *parser) Parse(log string) (resLog interface{}, eventType chain.EventLogType, err error) {
	var event Event
	err = json.Unmarshal([]byte(log), &event)
	if err != nil {
		return nil, "", err
	}

	if event.name == chain.TransferStart {
		resLog = chain.TransferStartEvent{}
	} else if event.name == chain.TransferReceived {
		resLog = chain.TransferReceivedEvent{}
	} else if event.name == chain.TransferEnd {
		resLog = chain.TransferEndEvent{}
	} else {
		err = fmt.Errorf("no matching signature for event log for %v found", eventType)
	}
	return
}
