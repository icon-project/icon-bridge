/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package substrate_eth

import (
	"encoding/hex"
	"fmt"
	subEthTypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/substrate-eth/types"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventLog struct {
	Addr    []byte
	Indexed [][]byte
	Data    [][]byte
}

type TransactionResult struct {
	To                 Address `json:"to"`
	CumulativeStepUsed HexInt  `json:"cumulativeStepUsed"`
	StepUsed           HexInt  `json:"stepUsed"`
	StepPrice          HexInt  `json:"stepPrice"`
	EventLogs          []struct {
		Addr    Address  `json:"scoreAddress"`
		Indexed []string `json:"indexed"`
		Data    []string `json:"data"`
	} `json:"eventLogs"`
	LogsBloom HexBytes `json:"logsBloom"`
	Status    HexInt   `json:"status"`
	Failure   *struct {
		CodeValue    HexInt `json:"code"`
		MessageValue string `json:"message"`
	} `json:"failure,omitempty"`
	SCOREAddress Address  `json:"scoreAddress,omitempty"`
	BlockHash    HexBytes `json:"blockHash" validate:"required,t_hash"`
	BlockHeight  HexInt   `json:"blockHeight" validate:"required,t_int"`
	TxIndex      HexInt   `json:"txIndex" validate:"required,t_int"`
	TxHash       HexBytes `json:"txHash" validate:"required,t_int"`
}

type TransactionParam struct {
	FromAddress string      `json:"from" validate:"required,t_addr_eoa"`
	ToAddress   string      `json:"to" validate:"required,t_addr"`
	NetworkID   HexInt      `json:"nid" validate:"required,t_int"`
	Params      interface{} `json:"params,omitempty"`
	TransactOpt *bind.TransactOpts
}

type BMCRelayMethodParams struct {
	Prev     string `json:"_prev"`
	Messages string `json:"_msg"`
}

type BMCStatus struct {
	TxSeq            HexInt `json:"tx_seq"`
	RxSeq            HexInt `json:"rx_seq"`
	BMRIndex         HexInt `json:"relay_idx"`
	RotateHeight     HexInt `json:"rotate_height"`
	RotateTerm       HexInt `json:"rotate_term"`
	DelayLimit       HexInt `json:"delay_limit"`
	MaxAggregation   HexInt `json:"max_agg"`
	CurrentHeight    HexInt `json:"cur_height"`
	RxHeight         HexInt `json:"rx_height"`
	RxHeightSrc      HexInt `json:"rx_height_src"`
	BlockIntervalSrc HexInt `json:"block_interval_src"`
	BlockIntervalDst HexInt `json:"block_interval_dst"`
}

type TransactionHashParam struct {
	Hash common.Hash
}

type BlockRequest struct {
	Height             *big.Int       `json:"height"`
	EventFilters       []*EventFilter `json:"eventFilters,omitempty"`
	SrcContractAddress common.Address `json:"srcContractAddress,omitempty"`
}

type EventFilter struct {
	Addr      Address   `json:"addr,omitempty"`
	Signature string    `json:"event"`
	Indexed   []*string `json:"indexed,omitempty"`
	Data      []*string `json:"data,omitempty"`
}

type BlockNotification struct {
	Hash          common.Hash
	Height        *big.Int
	Header        *subEthTypes.Header
	Receipts      types.Receipts
	HasBTPMessage *bool
}

type RelayMessage struct {
	ReceiptProofs [][]byte
	//
	height        int64
	eventSequence int64
	numberOfEvent int
}

type ReceiptProof struct {
	Index  int
	Events []byte
	Height int64
}

type EVMLog struct {
	Address     string
	Topics      [][]byte
	Data        []byte
	BlockNumber uint64
	TxHash      []byte
	TxIndex     uint
	BlockHash   []byte
	Index       uint
	Removed     bool
}

func MakeLog(log *types.Log) *EVMLog {
	topics := make([][]byte, 0)

	for _, topic := range log.Topics {
		topics = append(topics, topic.Bytes())
	}

	return &EVMLog{
		Address:     log.Address.String(),
		Topics:      topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.Bytes(),
		TxIndex:     log.TxIndex,
		BlockHash:   log.BlockHash.Bytes(),
		Index:       log.Index,
		Removed:     log.Removed,
	}
}

type Receipt struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Type              uint8
	PostState         []byte
	Status            uint64
	CumulativeGasUsed uint64
	Bloom             []byte
	Logs              []*EVMLog

	TxHash          common.Hash
	ContractAddress common.Address
	GasUsed         uint64

	BlockHash        common.Hash
	BlockNumber      uint64
	TransactionIndex uint
}

func MakeReceipt(receipt *types.Receipt) *Receipt {
	logs := make([]*EVMLog, len(receipt.Logs))

	for _, log := range receipt.Logs {
		logs = append(logs, MakeLog(log))
	}

	return &Receipt{
		Type:              receipt.Type,
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		Bloom:             receipt.Bloom.Bytes(),
		Logs:              logs,
		TxHash:            receipt.TxHash,
		ContractAddress:   receipt.ContractAddress,
		GasUsed:           receipt.GasUsed,
		BlockHash:         receipt.BlockHash,
		BlockNumber:       receipt.BlockNumber.Uint64(),
		TransactionIndex:  receipt.TransactionIndex,
	}
}

func HexToAddress(s string) common.Address {
	return common.HexToAddress(s)
}

// HexBytes T_BIN_DATA, T_HASH
type HexBytes string

func (hs HexBytes) Value() ([]byte, error) {
	if hs == "" {
		return nil, nil
	}
	return hex.DecodeString(string(hs[2:]))
}

// HexInt T_INT
type HexInt string

func (i HexInt) Value() (int64, error) {
	s := string(i)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	return strconv.ParseInt(s, 16, 64)
}

func (i HexInt) Int() (int, error) {
	s := string(i)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	v, err := strconv.ParseInt(s, 16, 32)
	return int(v), err
}

// Address T_ADDR_EOA, T_ADDR_SCORE
type Address string

func (a Address) Value() ([]byte, error) {
	var b [21]byte
	switch a[:2] {
	case "cx":
		b[0] = 1
	case "hx":
	default:
		return nil, fmt.Errorf("invalid prefix %s", a[:2])
	}
	n, err := hex.Decode(b[1:], []byte(a[2:]))
	if err != nil {
		return nil, err
	}
	if n != 20 {
		return nil, fmt.Errorf("invalid length %d", n)
	}
	return b[:], nil
}

// Signature T_SIG
type Signature string
