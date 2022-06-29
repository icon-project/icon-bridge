package hmny

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/harmony-one/harmony/accounts/abi"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	nativeHmy "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/bsh/bshPeriphery"
	tokenHmy "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/tokenbsh/bshimpl"
	"github.com/pkg/errors"
)

type parser struct {
	backend               bind.ContractBackend
	genNativeObj          *nativeHmy.NativeHmy
	genTokenObj           *tokenHmy.TokenHmy
	eventIDToNameNative   map[common.Hash]string
	eventIDToNameToken    map[common.Hash]string
	addressToContractName map[string]chain.ContractName
}

func find(m map[string]chain.ContractName, key chain.ContractName) *string {
	for cAddr, cName := range m {
		if cName == key {
			return &cAddr
		}
	}
	return nil
}

func NewParser(url string, nameToAddr map[chain.ContractName]string) (*parser, error) {
	var err error
	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	addrToContractName := map[string]chain.ContractName{}
	for name, addr := range nameToAddr {
		addrToContractName[addr] = name
	}
	p := &parser{backend: ethclient.NewClient(clrpc)}
	nativeAddr := find(addrToContractName, chain.NativeBSHPeripheryHmy)
	if nativeAddr == nil {
		return nil, errors.New("Did not find native hmy contract in input map")
	}
	tokenAddr := find(addrToContractName, chain.TokenBSHImplHmy)
	if tokenAddr == nil {
		return nil, errors.New("Did not find token hmy contract in input map")
	}
	p.genNativeObj, err = nativeHmy.NewNativeHmy(common.HexToAddress(*nativeAddr), p.backend)
	if err != nil {
		return nil, err
	}
	p.genTokenObj, err = tokenHmy.NewTokenHmy(common.HexToAddress(*tokenAddr), p.backend)
	if err != nil {
		return nil, err
	}
	p.eventIDToNameNative, err = eventIDToName(nativeHmy.NativeHmyABI)
	if err != nil {
		return nil, err
	}
	p.eventIDToNameToken, err = eventIDToName(tokenHmy.TokenHmyABI)
	if err != nil {
		return nil, err
	}
	p.addressToContractName = addrToContractName
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
	cName, ok := p.addressToContractName[log.Address.String()]
	if !ok {
		err = errors.New("Couldn't find contract matching the log")
		return
	}
	if cName == chain.NativeBSHPeripheryHmy {
		tres := findTopic(log.Topics, p.eventIDToNameNative)
		if tres == nil {
			err = errors.New("Topic not among mentioned ones")
			return
		}
		eventType = chain.EventLogType(*tres)
		if eventType == chain.TransferStart {
			resLog, err = p.parseTransferStartNativeCoin(log)
		} else if eventType == chain.TransferReceived {
			resLog, err = p.parseTransferReceivedNativeCoin(log)
		} else if eventType == chain.TransferEnd {
			resLog, err = p.parseTransferEndNativeCoin(log)
		} else {
			err = errors.New("No matching signature ")
		}
	} else if cName == chain.TokenBSHImplHmy {
		tres := findTopic(log.Topics, p.eventIDToNameNative)
		if tres == nil {
			err = errors.New("Topic not among mentioned ones")
			return
		}
		eventType = chain.EventLogType(*tres)
		if eventType == chain.TransferStart {
			resLog, err = p.parseTransferStartToken(log)
		} else if eventType == chain.TransferReceived {
			resLog, err = p.parseTransferReceivedToken(log)
		} else if eventType == chain.TransferEnd {
			resLog, err = p.parseTransferEndToken(log)
		} else {
			err = errors.New("No matching signature ")
		}
	} else {
		err = errors.New("Contract not amongst processed ones")
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

func (p *parser) parseTransferStartNativeCoin(hlog *ethTypes.Log) (*chain.TransferStartEvent, error) {
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
	out, err := p.genNativeObj.ParseTransferStart(log)
	if err != nil {
		return nil, err
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

func (p *parser) parseTransferReceivedNativeCoin(hlog *ethTypes.Log) (*chain.TransferReceivedEvent, error) {
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
	out, err := p.genNativeObj.ParseTransferReceived(log)
	if err != nil {
		return nil, err
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

func (p *parser) parseTransferEndNativeCoin(hlog *ethTypes.Log) (*chain.TransferEndEvent, error) {
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
	out, err := p.genNativeObj.ParseTransferEnd(log)
	if err != nil {
		return nil, err
	}
	return &chain.TransferEndEvent{
		From: out.From.String(),
		Sn:   out.Sn,
		Code: out.Code,
	}, nil
}

func (p *parser) parseTransferStartToken(hlog *ethTypes.Log) (*chain.TransferStartEvent, error) {
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
	out, err := p.genTokenObj.ParseTransferStart(log)
	if err != nil {
		return nil, err
	}
	newAssetDetails := make([]chain.AssetTransferDetails, len(out.Assets))
	for i, v := range out.Assets {
		newAssetDetails[i].Name = v.Name
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

func (p *parser) parseTransferReceivedToken(hlog *ethTypes.Log) (*chain.TransferReceivedEvent, error) {
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
	out, err := p.genTokenObj.ParseTransferReceived(log)
	if err != nil {
		return nil, err
	}
	newAssetDetails := make([]chain.AssetTransferDetails, len(out.AssetDetails))
	for i, v := range out.AssetDetails {
		newAssetDetails[i].Name = v.Name
		newAssetDetails[i].Value = v.Value
		newAssetDetails[i].Fee = v.Fee
	}
	return &chain.TransferReceivedEvent{
		From:   out.From.String(),
		To:     out.To.String(),
		Sn:     out.Sn,
		Assets: newAssetDetails,
	}, nil
}

func (p *parser) parseTransferEndToken(hlog *ethTypes.Log) (*chain.TransferEndEvent, error) {
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
	out, err := p.genTokenObj.ParseTransferEnd(log)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	for _, a := range abi.Events {
		resMap[a.ID] = a.Name
	}
	return resMap, nil
}
