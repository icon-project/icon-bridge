package bshPeriphery

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

type bshPeripheryContract struct {
	backend       bind.ContractBackend
	genObj        *BshPeriphery
	eventIDToName map[common.Hash]string
}

func (b *bshPeripheryContract) GetName() contracts.ContractName {
	return contracts.BSHPeriphery
}

func NewContract(url string, cAddr common.Address) (contracts.Contract, error) {
	var err error
	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	ctr := &bshPeripheryContract{backend: ethclient.NewClient(clrpc)}

	ctr.genObj, err = NewBshPeriphery(cAddr, ctr.backend)
	if err != nil {
		return nil, err
	}
	ctr.eventIDToName, err = contracts.EventIDToName(BshPeripheryABI)
	if err != nil {
		return nil, err
	}
	return ctr, nil
}

func (b *bshPeripheryContract) Decode(log types.Log) (res map[string]interface{}, err error) {
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
