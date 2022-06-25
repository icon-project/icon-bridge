package tokenHmy

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type tokenHmyContract struct {
	backend       bind.ContractBackend
	genObj        *TokenHmy
	eventIDToName map[common.Hash]string
}

func (b *tokenHmyContract) GetName() contracts.ContractName {
	return contracts.TokenHmy
}

func NewContract(url string, cAddr common.Address) (contracts.Contract, error) {
	var err error
	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	ctr := &tokenHmyContract{backend: ethclient.NewClient(clrpc)}

	ctr.genObj, err = NewTokenHmy(cAddr, ctr.backend)
	if err != nil {
		return nil, err
	}
	ctr.eventIDToName, err = contracts.EventIDToName(TokenHmyABI)
	if err != nil {
		return nil, err
	}
	return ctr, nil
}
func (b *tokenHmyContract) Decode(l interface{}) (res map[string]interface{}, err error) {
	log, ok := l.(types.Log)
	if !ok {
		return nil, errors.New("Log of wrong type. Expected types.Log")
	}
	res = map[string]interface{}{}
	for _, tid := range log.Topics {
		topicName, ok := b.eventIDToName[tid]
		if !ok {
			continue
		}
		if topicName == "TransferStart" {
			if out, err := b.genObj.ParseTransferStart(log); err != nil {
				return nil, err
			} else {
				res[topicName] = out
			}
		} else if topicName == "TransferReceived" {
			if out, err := b.genObj.ParseTransferReceived(log); err != nil {
				return nil, err
			} else {
				res[topicName] = out
			}
		} else if topicName == "TransferEnd" {
			if out, err := b.genObj.ParseTransferEnd(log); err != nil {
				return nil, err
			} else {
				res[topicName] = out
			}
		}
	}
	return res, nil
}
