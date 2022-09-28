package bsc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	bmcp "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc/abi/bmcperiphery"
	btsp "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc/abi/btsperiphery"
	"github.com/pkg/errors"
)

const BLACKLIST_MESSAGE = 3
const CHANGE_TOKEN_LIMIT = 4

type parser struct {
	backend               bind.ContractBackend
	genBtsObj             *btsp.Btsperiphery
	genBmcObj             *bmcp.Bmcperiphery
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
	bmcperiAddr, ok := nameToAddr[chain.BMCPeriphery]
	if !ok {
		return nil, fmt.Errorf("nameToAddr doesn't include %v ", chain.BMCPeriphery)
	}

	p.genBtsObj, err = btsp.NewBtsperiphery(common.HexToAddress(btsperiAddr), p.backend)
	if err != nil {
		err = errors.Wrap(err, "nativeHmy.NewNativeHmy ")
		return nil, err
	}
	p.genBmcObj, err = bmcp.NewBmcperiphery(common.HexToAddress(bmcperiAddr), p.backend)
	if err != nil {
		err = errors.Wrap(err, "nativeHmy.NewNativeHmy ")
		return nil, err
	}

	p.eventIDToName, err = eventIDToName(btsp.BtsperipheryABI)
	if err != nil {
		err = errors.Wrap(err, "eventIDToName ")
		return nil, err
	}
	bmcEventIDToName, err := eventIDToName(bmcp.BmcperipheryABI)
	if err != nil {
		err = errors.Wrap(err, "eventIDToName ")
		return nil, err
	}
	for k, v := range bmcEventIDToName {
		p.eventIDToName[k] = v
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

func (p *parser) Parse(log *ethTypes.Log) (resLog interface{}, eventType chain.EventLogType, err error) {

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
	} else if eventType == chain.Message {
		resLog, eventType, err = p.parseMessage(log) // handles BlackList & TokenLimit request/response parsing
	} else {
		err = fmt.Errorf("Unexpected eventType. Got %v ", eventType)
	}
	return
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

func (p *parser) parseMessage(hlog *ethTypes.Log) (res interface{}, evt chain.EventLogType, err error) {
	out, err := p.genBmcObj.ParseMessage(*hlog)
	if err != nil {
		err = errors.Wrapf(err, "genBmcObj.ParseMessage %v", err)
		return
	}
	bmcMsg := BMCMessage{}
	if err = rlpDecodeHex(hex.EncodeToString(out.Msg), &bmcMsg); err != nil {
		err = errors.Wrapf(err, "rlpDecodeHex %v", err)
		return
	}
	msgSn := (&big.Int{}).SetBytes(bmcMsg.Sn)

	svcMessage := ServiceMessage{}
	if err = rlpDecodeHex(hex.EncodeToString(bmcMsg.Message), &svcMessage); err != nil {
		err = errors.Wrapf(err, "rlpDecodeHex %v", err)
		return
	}
	svcMessagePayload := ServiceMessagePayload{}
	if err = rlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcMessagePayload); err != nil {
		err = errors.Wrapf(err, "rlpDecodeHex %v", err)
		return
	}
	svcTypeNum := (&big.Int{}).SetBytes(svcMessage.ServiceType).Int64()
	svcPayloadCode := (&big.Int{}).SetBytes(svcMessagePayload.Code).Int64()
	if svcTypeNum == BLACKLIST_MESSAGE {
		return &chain.BlacklistResponseEvent{
			Sn:   msgSn,
			Code: svcPayloadCode,
			Msg:  string(svcMessagePayload.Msg),
		}, chain.BlacklistResponse, nil
	} else if svcTypeNum == CHANGE_TOKEN_LIMIT {
		return &chain.TokenLimitResponseEvent{
			Sn:   msgSn,
			Code: svcPayloadCode,
			Msg:  string(svcMessagePayload.Msg),
		}, chain.TokenLimitResponse, nil
	}
	return nil, "", fmt.Errorf("Unexpected Message Type SvcTypeNum %v", svcTypeNum)
}

func rlpDecodeHex(str string, out interface{}) error {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	input, err := hex.DecodeString(str)
	if err != nil {
		return errors.Wrap(err, "hex.DecodeString ")
	}
	err = rlp.Decode(bytes.NewReader(input), out)
	if err != nil {
		return errors.Wrap(err, "rlp.Decode ")
	}
	return nil
}

type BMCMessage struct {
	Src     string //  an address of BMC (i.e. btp://1234.PARA/0x1234)
	Dst     string //  an address of destination BMC
	Svc     string //  service name of BSH
	Sn      []byte //  sequence number of BMC
	Message []byte //  serialized Service Message from BSH
}

type ServiceMessage struct {
	ServiceType []byte
	Payload     []byte
}

type ServiceMessagePayload struct {
	Code []byte
	Msg  []byte
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
