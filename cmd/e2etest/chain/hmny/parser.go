//go:build hmny
// +build hmny

package hmny

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/harmony-one/harmony/accounts/abi"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	btsp "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/btsperiphery"
	"github.com/pkg/errors"
)

type parser struct {
	backend               bind.ContractBackend
	genBtsObj             *btsp.Btsperiphery
	eventIDToName         map[common.Hash]string
	addressToContractName map[string]chain.ContractName
}

func NewParser(url string, nameToAddr map[chain.ContractName]string) (*parser, error) {
	var err error
	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, errors.Wrapf(err, "rpc.Dial(%v) ", url)
	}
	p := &parser{backend: ethclient.NewClient(clrpc), addressToContractName: make(map[string]chain.ContractName)}
	btsperiAddr, ok := nameToAddr[chain.BTSPeriphery]
	if !ok {
		return nil, fmt.Errorf("nameToAddr doesn't include %v ", chain.BTSPeriphery)
	}

	p.genBtsObj, err = btsp.NewBtsperiphery(common.HexToAddress(btsperiAddr), p.backend)
	if err != nil {
		err = errors.Wrap(err, "nativeHmy.NewNativeHmy ")
		return nil, err
	}

	p.eventIDToName, err = eventIDToName(btsp.BtsperipheryABI)
	if err != nil {
		err = errors.Wrap(err, "eventIDToName ")
		return nil, err
	}

	for name, addr := range nameToAddr {
		p.addressToContractName[addr] = name
	}

	return p, nil
}

func findTopic(topics []common.Hash, eventIDToName map[common.Hash]string) *string {
	for _, tid := range topics {
		topicName, ok := eventIDToName[tid]
		if !ok {
			continue
		}
		return &topicName
	}
	return nil
}

func (p *parser) ParseEth(log *ethTypes.Log) (resLog interface{}, eventType chain.EventLogType, err error) {

	tres := findTopic(log.Topics, p.eventIDToName)
	if tres == nil {
		err = errors.New("log.Topics not among p.eventIDToNameNative ")
		return
	}
	eventType = chain.EventLogType(*tres)
	if eventType == chain.TransferStart {
		resLog, err = p.parseTransferStart(log)
	} else if eventType == chain.TransferReceived {
		resLog, err = p.parseTransferReceived(log)
	} else if eventType == chain.TransferEnd {
		resLog, err = p.parseTransferEnd(log)
	} else {
		err = fmt.Errorf("Unexpected eventType. Got %v ", eventType)
	}
	return
}

func (p *parser) Parse(hlog *types.Log) (resLog interface{}, eventType chain.EventLogType, err error) {
	log := &ethTypes.Log{
		Address:     hlog.Address,
		Topics:      hlog.Topics,
		Data:        hlog.Data,
		BlockNumber: hlog.BlockNumber,
		TxHash:      hlog.TxHash,
		TxIndex:     hlog.TxIndex,
		BlockHash:   hlog.BlockHash,
		Index:       hlog.Index,
		Removed:     hlog.Removed,
	}
	return p.ParseEth(log)
}

func (p *parser) parseTransferStart(hlog *ethTypes.Log) (*chain.TransferStartEvent, error) {
	log := ethTypes.Log{
		Address:     hlog.Address,
		Topics:      hlog.Topics,
		Data:        hlog.Data,
		BlockNumber: hlog.BlockNumber,
		TxHash:      hlog.TxHash,
		TxIndex:     hlog.TxIndex,
		BlockHash:   hlog.BlockHash,
		Index:       hlog.Index,
		Removed:     hlog.Removed,
	}
	out, err := p.genBtsObj.ParseTransferStart(log)
	if err != nil {
		return nil, errors.Wrap(err, "genBtsObj.ParseTransferStart ")
	}
	newAssetDetails := make([]chain.AssetTransferDetails, len(out.AssetDetails))
	for i, v := range out.AssetDetails {
		newAssetDetails[i].Name = v.CoinName
		newAssetDetails[i].Value = v.Value
		newAssetDetails[i].Fee = v.Fee
	}
	return &chain.TransferStartEvent{
		From:   out.From.String(),
		To:     out.To,
		Sn:     out.Sn,
		Assets: newAssetDetails,
	}, nil
}

func (p *parser) parseTransferReceived(hlog *ethTypes.Log) (*chain.TransferReceivedEvent, error) {
	log := ethTypes.Log{
		Address:     hlog.Address,
		Topics:      hlog.Topics,
		Data:        hlog.Data,
		BlockNumber: hlog.BlockNumber,
		TxHash:      hlog.TxHash,
		TxIndex:     hlog.TxIndex,
		BlockHash:   hlog.BlockHash,
		Index:       hlog.Index,
		Removed:     hlog.Removed,
	}
	out, err := p.genBtsObj.ParseTransferReceived(log)
	if err != nil {
		return nil, errors.Wrap(err, "genBtsObj.ParseTransferReceived ")
	}
	newAssetDetails := make([]chain.AssetTransferDetails, len(out.AssetDetails))
	for i, v := range out.AssetDetails {
		newAssetDetails[i].Name = v.CoinName
		newAssetDetails[i].Value = v.Value
	}
	return &chain.TransferReceivedEvent{
		From:   out.From.String(),
		To:     out.To.String(),
		Sn:     out.Sn,
		Assets: newAssetDetails,
	}, nil
}

func (p *parser) parseTransferEnd(hlog *ethTypes.Log) (*chain.TransferEndEvent, error) {
	log := ethTypes.Log{
		Address:     hlog.Address,
		Topics:      hlog.Topics,
		Data:        hlog.Data,
		BlockNumber: hlog.BlockNumber,
		TxHash:      hlog.TxHash,
		TxIndex:     hlog.TxIndex,
		BlockHash:   hlog.BlockHash,
		Index:       hlog.Index,
		Removed:     hlog.Removed,
	}
	out, err := p.genBtsObj.ParseTransferEnd(log)
	if err != nil {
		return nil, errors.Wrap(err, "genBtsObj.ParseTransferEnd ")
	}
	return &chain.TransferEndEvent{
		From: out.From.String(),
		Sn:   out.Sn,
		Code: out.Code,
	}, nil
}

func eventIDToName(abiStr string) (map[common.Hash]string, error) {
	resMap := map[common.Hash]string{}
	abi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.Wrap(err, "abi.JSON ")
	}
	for _, a := range abi.Events {
		resMap[a.ID] = a.Name
	}
	return resMap, nil
}
