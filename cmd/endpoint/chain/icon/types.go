package icon

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
)

type TransactionResult struct {
	To                 Address       `json:"to"`
	CumulativeStepUsed HexInt        `json:"cumulativeStepUsed"`
	StepUsed           HexInt        `json:"stepUsed"`
	StepPrice          HexInt        `json:"stepPrice"`
	EventLogs          []TxnEventLog `json:"eventLogs"`
	LogsBloom          HexBytes      `json:"logsBloom"`
	Status             HexInt        `json:"status"`
	Failure            *struct {
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
	Version     HexInt      `json:"version" validate:"required,t_int"`
	FromAddress Address     `json:"from" validate:"required,t_addr_eoa"`
	ToAddress   Address     `json:"to" validate:"required,t_addr"`
	Value       HexInt      `json:"value,omitempty" validate:"optional,t_int"`
	StepLimit   HexInt      `json:"stepLimit" validate:"required,t_int"`
	Timestamp   HexInt      `json:"timestamp" validate:"required,t_int"`
	NetworkID   HexInt      `json:"nid" validate:"required,t_int"`
	Nonce       HexInt      `json:"nonce,omitempty" validate:"optional,t_int"`
	Signature   string      `json:"signature" validate:"required,t_sig"`
	DataType    string      `json:"dataType,omitempty" validate:"optional,call|deploy|message"`
	Data        interface{} `json:"data,omitempty"`
	TxHash      HexBytes    `json:"-"`
}

type AddressParam struct {
	Address Address `json:"address" validate:"required,t_addr"`
	Height  HexInt  `json:"height,omitempty" validate:"optional,t_int"`
}

type CallParam struct {
	FromAddress Address     `json:"from" validate:"optional,t_addr_eoa"`
	ToAddress   Address     `json:"to" validate:"required,t_addr_score"`
	DataType    string      `json:"dataType" validate:"required,call"`
	Data        interface{} `json:"data"`
}

type TransactionHashParam struct {
	Hash HexBytes `json:"txHash" validate:"required,t_hash"`
}

type BlockHeightParam struct {
	Height HexInt `json:"height" validate:"required,t_int"`
}

type Block struct {
	//BlockHash              HexBytes  `json:"block_hash" validate:"required,t_hash"`
	//Version                HexInt    `json:"version" validate:"required,t_int"`
	Height    int64 `json:"height" validate:"required,t_int"`
	Timestamp int64 `json:"time_stamp" validate:"required,t_int"`
	//Proposer               HexBytes  `json:"peer_id" validate:"optional,t_addr_eoa"`
	//PrevID                 HexBytes  `json:"prev_block_hash" validate:"required,t_hash"`
	//NormalTransactionsHash HexBytes  `json:"merkle_tree_root_hash" validate:"required,t_hash"`
	NormalTransactions []struct {
		TxHash HexBytes `json:"txHash"`
		//Version   HexInt   `json:"version"`
		From Address `json:"from"`
		To   Address `json:"to"`
		//Value     HexInt   `json:"value,omitempty" `
		//StepLimit HexInt   `json:"stepLimit"`
		//TimeStamp HexInt   `json:"timestamp"`
		//NID       HexInt   `json:"nid,omitempty"`
		//Nonce     HexInt   `json:"nonce,omitempty"`
		//Signature HexBytes `json:"signature"`
		//DataType  string          `json:"dataType,omitempty"`
		//Data json.RawMessage `json:"data,omitempty"`
	} `json:"confirmed_transaction_list"`
	//Signature              HexBytes  `json:"signature" validate:"optional,t_hash"`
}

type BlockRequest struct {
	Height       HexInt         `json:"height"`
	EventFilters []*EventFilter `json:"eventFilters,omitempty"`
}

type EventFilter struct {
	Addr      Address   `json:"addr,omitempty"`
	Signature string    `json:"event"`
	Indexed   []*string `json:"indexed,omitempty"`
	Data      []*string `json:"data,omitempty"`
}

type BlockNotification struct {
	Hash    HexBytes     `json:"hash"`
	Height  HexInt       `json:"height"`
	Indexes [][]HexInt   `json:"indexes,omitempty"`
	Events  [][][]HexInt `json:"events,omitempty"`
}

type EventRequest struct {
	EventFilter
	Height HexInt `json:"height"`
}

type TxnEventLog struct {
	Addr    Address  `json:"scoreAddress"`
	Indexed []string `json:"indexed"`
	Data    []string `json:"data"`
}

type TxnLog struct {
	TxHash      HexBytes      `json:"txHash" validate:"required,t_int"`
	From        Address       `json:"from"`
	To          Address       `json:"to"`
	EventLogs   []TxnEventLog `json:"eventLogs"`
	Status      HexInt        `json:"status"`
	BlockHeight int64         `json:"blockHeight"`
}

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}

const (
	JsonrpcApiVersion                                = 3
	JsonrpcErrorCodeSystem         jsonrpc.ErrorCode = -31000
	JsonrpcErrorCodeTxPoolOverflow jsonrpc.ErrorCode = -31001
	JsonrpcErrorCodePending        jsonrpc.ErrorCode = -31002
	JsonrpcErrorCodeExecuting      jsonrpc.ErrorCode = -31003
	JsonrpcErrorCodeNotFound       jsonrpc.ErrorCode = -31004
	JsonrpcErrorLackOfResource     jsonrpc.ErrorCode = -31005
	JsonrpcErrorCodeTimeout        jsonrpc.ErrorCode = -31006
	JsonrpcErrorCodeSystemTimeout  jsonrpc.ErrorCode = -31007
	JsonrpcErrorCodeScore          jsonrpc.ErrorCode = -30000
)

const (
	DuplicateTransactionError = iota + 2000
	TransactionPoolOverflowError
	ExpiredTransactionError
	FutureTransactionError
	TransitionInterruptedError
	InvalidTransactionError
	InvalidQueryError
	InvalidResultError
	NoActiveContractError
	NotContractAddressError
	InvalidPatchDataError
	CommittedTransactionError
)

const (
	ResultStatusSuccess           = "0x1"
	ResultStatusFailureCodeRevert = 32
	ResultStatusFailureCodeEnd    = 99
)

//T_BIN_DATA, T_HASH
type HexBytes string

func (hs HexBytes) Value() ([]byte, error) {
	if hs == "" {
		return nil, nil
	}
	return hex.DecodeString(string(hs[2:]))
}
func NewHexBytes(b []byte) HexBytes {
	return HexBytes("0x" + hex.EncodeToString(b))
}

//T_INT
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

func (i HexInt) BigInt() (*big.Int, error) {
	bi := new(big.Int)
	if err := intconv.ParseBigInt(bi, string(i)); err != nil {
		return nil, err
	} else {
		return bi, nil
	}
}

func NewHexInt(v int64) HexInt {
	return HexInt("0x" + strconv.FormatInt(v, 16))
}

//T_ADDR_EOA, T_ADDR_SCORE
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

func NewAddress(b []byte) Address {
	if len(b) != 21 {
		return ""
	}
	switch b[0] {
	case 1:
		return Address("cx" + hex.EncodeToString(b[1:]))
	case 0:
		return Address("hx" + hex.EncodeToString(b[1:]))
	default:
		return ""
	}
}
