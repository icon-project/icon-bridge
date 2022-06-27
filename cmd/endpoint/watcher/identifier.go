package watcher

import (
	"errors"
	"math/big"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	ctr "github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/nativeHmy"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/nativeIcon"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/tokenHmy"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/tokenIcon"
)

// If txn and txr matches a condition, spawn and instance of an identifier
// This identifier object will check the eventLog for a matching criteria
// So, for a single transaction there can be multiple identifiers
// All of these identifiers will send to a single channel
// This single channel will act as sink for multiple identifiers of a single txn
// A test unit can have multiple such txns: all results should be sent to that single channel
// Identifier Group maps inputCriterion to an instance of txn

type eventLogInfo struct {
	sourceChain  chain.ChainType
	contractName ctr.ContractName
	eventType    string
	eventLog     interface{}
}

type args struct {
	req     *chainAPI.RequestParam
	initRes interface{}
}

type identifierGroup struct {
	name        string
	description string
	init        func(args []eventLogInfo, req *chainAPI.RequestParam) (interface{}, bool)
	idfs        []identifier
}

type identifier struct {
	preRun func(args args, info eventLogInfo) bool
	run    func(args args, info eventLogInfo) (bool, error)
}

var DefaultIdentifierGroup = []identifierGroup{
	{
		name:        "InterChainI2H",
		description: "Transactions from icon to hmny chain for all tokens",
		init: func(logs []eventLogInfo, req *chainAPI.RequestParam) (interface{}, bool) {
			if req.FromChain == chain.ICON && req.ToChain == chain.HMNY {
				var argSn *big.Int = nil
				for _, al := range logs {
					if al.sourceChain == chain.ICON && al.eventType == "TransferStart" {
						if ad, aok := al.eventLog.(*nativeIcon.NativeIconTransferStart); aok {
							argSn = ad.Sn
							return argSn, true
						} else {
							return errors.New("Interface should have been *nativeIcon.NativeIconTransferStart "), false
						}
					}
				}
			}
			return nil, false
		},
		idfs: []identifier{ // all identifiers should match
			{ // identifier for TransferStart for all tokens
				preRun: func(args args, info eventLogInfo) bool { // TransferStart msg on icon chain
					if info.sourceChain == chain.ICON && info.eventType == "TransferStart" &&
						(info.contractName == ctr.NativeIcon || info.contractName == ctr.TokenIcon) {
						return true
					}
					return false
				},
				run: func(args args, info eventLogInfo) (bool, error) {
					if info.contractName == ctr.NativeIcon {
						if el, ok := info.eventLog.(*nativeIcon.NativeIconTransferStart); ok {
							seq := args.initRes.(*big.Int)
							if el.From == args.req.FromAddress && args.req.ToAddress == el.To && el.Sn.Int64() == seq.Int64() {
								return true, nil
							}
						} else {
							return false, errors.New("Expected *nativeIcon.NativeIconTransferStart")
						}
					} else if info.contractName == ctr.TokenIcon {
						if el, ok := info.eventLog.(*tokenIcon.TokenIconTransferStart); ok {
							seq := args.initRes.(*big.Int)
							if el.From == args.req.FromAddress && args.req.ToAddress == el.To && el.Sn.Int64() == seq.Int64() {
								return true, nil
							}
						} else {
							return false, errors.New("Expected *tokenIcon.TokenIconTransferStart")
						}
					}
					return false, nil
				},
			},
			{
				preRun: func(args args, info eventLogInfo) bool { // TransferStart msg on icon chain
					if info.sourceChain == chain.ICON && info.eventType == "TransferEnd" &&
						(info.contractName == ctr.NativeIcon || info.contractName == ctr.TokenIcon) {
						return true
					}
					return false
				},
				run: func(args args, info eventLogInfo) (bool, error) {
					if info.contractName == ctr.NativeIcon {
						if el, ok := info.eventLog.(*nativeIcon.NativeIconTransferEnd); ok {
							seq := args.initRes.(*big.Int)
							if el.From == args.req.FromAddress && el.Sn.Int64() == seq.Int64() {
								return true, nil
							}
						} else {
							return false, errors.New("Expected *nativeIcon.NativeIconTransferEnd")
						}
					} else if info.contractName == ctr.TokenIcon {
						if el, ok := info.eventLog.(*tokenIcon.TokenIconTransferEnd); ok {
							seq := args.initRes.(*big.Int)
							if el.From == args.req.FromAddress && el.Sn.Int64() == seq.Int64() {
								return true, nil
							}
						} else {
							return false, errors.New("Expected *tokenIcon.NativeIconTransferEnd")
						}
					}
					return false, nil
				},
			},
			{
				preRun: func(args args, info eventLogInfo) bool { // TransferStart msg on icon chain
					if info.sourceChain == chain.HMNY && info.eventType == "TransferReceived" &&
						(info.contractName == ctr.NativeHmy || info.contractName == ctr.TokenHmy) {
						return true
					}
					return false
				},
				run: func(args args, info eventLogInfo) (bool, error) {
					if info.contractName == ctr.NativeHmy {
						if el, ok := info.eventLog.(*nativeHmy.NativeHmyTransferReceived); ok {
							seq := args.initRes.(*big.Int)
							splts := strings.Split(args.req.ToAddress, "/")
							if el.To.Hex() == splts[len(splts)-1] && el.Sn.Int64() == seq.Int64() {
								return true, nil
							}
						} else {
							return false, errors.New("Expected *nativeHmy.NativeHmyTransferReceived")
						}
					} else if info.contractName == ctr.TokenHmy {
						if el, ok := info.eventLog.(*tokenHmy.TokenHmyTransferReceived); ok {
							seq := args.initRes.(*big.Int)
							splts := strings.Split(args.req.ToAddress, "/")
							if el.To.Hex() == splts[len(splts)-1] && el.Sn.Int64() == seq.Int64() {
								return true, nil
							}
						} else {
							return false, errors.New("Expected *tokenHmy.TokenHmyTransferReceived")
						}
					}
					return false, nil
				},
			},
		},
	},
}
