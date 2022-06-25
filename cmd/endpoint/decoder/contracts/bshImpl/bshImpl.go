package bshImpl

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type bshImplContract struct {
	backend       bind.ContractBackend
	genObj        *BshImpl
	eventIDToName map[common.Hash]string
}

func (b *bshImplContract) GetName() contracts.ContractName {
	return contracts.BSHImpl
}

func NewContract(url string, cAddr common.Address) (contracts.Contract, error) {
	var err error
	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	ctr := &bshImplContract{backend: ethclient.NewClient(clrpc)}

	ctr.genObj, err = NewBshImpl(cAddr, ctr.backend)
	if err != nil {
		return nil, err
	}
	ctr.eventIDToName, err = contracts.EventIDToName(BshImplABI)
	if err != nil {
		return nil, err
	}
	return ctr, nil
}
func (b *bshImplContract) Decode(log types.Log) (res map[string]interface{}, err error) {
	res = map[string]interface{}{}
	for _, tid := range log.Topics {
		topicName, ok := b.eventIDToName[tid]
		if ok {
			fmt.Println("Okay", topicName, tid)
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
