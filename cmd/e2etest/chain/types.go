package chain

import (
	"context"
	"fmt"
	"math/big"
)

type ChainType string

const (
	ICON ChainType = "ICON"
	HMNY ChainType = "HMNY"
	BSC  ChainType = "BSC"
)

type ContractName string

const (
	TBNBBsc          ContractName = "TBNBBsc"
	TONEHmny         ContractName = "TONEHmny"
	TICXIcon         ContractName = "TICXIcon"
	BTSIcon          ContractName = "BTSIcon"
	BTSCoreHmny      ContractName = "BTSCoreHmny"
	BTSPeripheryHmny ContractName = "BTSPeripheryHmny"
	BTSCoreBsc       ContractName = "BTSCoreBsc"
	BTSPeripheryBsc  ContractName = "BTSPeripheryBsc"
)

type EventLogType string

const (
	TransferStart    EventLogType = "TransferStart"
	TransferReceived EventLogType = "TransferReceived"
	TransferEnd      EventLogType = "TransferEnd"
)

func (e EventLogType) String() string {
	return string(e)
}

type SrcAPI interface {
	Transfer(coinName, senderKey, recepientAddress string, amount big.Int) (txnHash string, err error)
	WaitForTxnResult(ctx context.Context, hash string) (txnr *TxnResult, err error)
	WatchForTransferStart(requestID uint64, seq int64) error
	WatchForTransferEnd(ID uint64, seq int64) error
	Approve(coinName string, ownerKey string, amount big.Int) (txnHash string, err error)
	GetCoinBalance(coinName string, addr string) (*big.Int, error)
	GetChainType() ChainType
	NativeCoinName() string
	TokenName() string
	GetAllowance(coinName, ownerAddr string) (amont *big.Int, err error)
}

type DstAPI interface {
	GetCoinBalance(coinName string, addr string) (*big.Int, error)
	WatchForTransferReceived(requestID uint64, seq int64) error
	GetChainType() ChainType
}

type TxnResult struct {
	StatusCode int
	ElInfo     []*EventLogInfo
	Raw        interface{}
}

type ChainAPI interface {
	Subscribe(ctx context.Context) (sinkChan chan *EventLogInfo, errChan chan error, err error)
	GetKeyPairs(num int) ([][2]string, error)
	GetBTPAddress(addr string) string

	Transfer(coinName, senderKey, recepientAddress string, amount big.Int) (txnHash string, err error)
	WaitForTxnResult(ctx context.Context, hash string) (txnr *TxnResult, err error)
	WatchForTransferStart(ID uint64, seq int64) error
	WatchForTransferReceived(ID uint64, seq int64) error
	WatchForTransferEnd(ID uint64, seq int64) error
	Approve(coinName string, ownerKey string, amount big.Int) (txnHash string, err error)
	GetCoinBalance(coinName string, addr string) (*big.Int, error)
	GetChainType() ChainType
	NativeCoinName() string
	TokenName() string
	GetAllowance(coinName, ownerAddr string) (amont *big.Int, err error)
}

type ChainConfig struct {
	Name               ChainType               `json:"name"`
	URL                string                  `json:"url"`
	ConftractAddresses map[ContractName]string `json:"contract_addresses"`
	GodWallet          GodWallet               `json:"god_wallet"`
	NetworkID          string                  `json:"network_id"`
}

type GodWallet struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

type EventLogInfo struct {
	IDs             []uint64
	ContractAddress string
	EventType       EventLogType
	EventLog        interface{}
}

func (e *EventLogInfo) GetSeq() (seq int64, err error) {
	if e.EventType == TransferStart {
		st, ok := e.EventLog.(*TransferStartEvent)
		if !ok {
			err = fmt.Errorf("Expected *TransferStartEvent. Got %v", e.EventLog)
		}
		seq = st.Sn.Int64()
	} else if e.EventType == TransferReceived {
		st, ok := e.EventLog.(*TransferReceivedEvent)
		if !ok {
			err = fmt.Errorf("Expected *TransferReceivedEvent. Got %v", e.EventLog)
		}
		seq = st.Sn.Int64()
	} else if e.EventType == TransferEnd {
		st, ok := e.EventLog.(*TransferEndEvent)
		if !ok {
			err = fmt.Errorf("Expected *TransferEndEvent. Got %v", e.EventLog)
		}
		seq = st.Sn.Int64()
	}
	return
}

type TransferStartEvent struct {
	From   string
	To     string
	Sn     *big.Int
	Assets []AssetTransferDetails
}

type TransferReceivedEvent struct {
	From   string
	To     string
	Sn     *big.Int
	Assets []AssetTransferDetails
}

type AssetTransferDetails struct {
	Name  string
	Value *big.Int
	Fee   *big.Int
}

type AssetDetails struct {
	Name  string
	Value *big.Int
}

type TransferEndEvent struct {
	From string
	Sn   *big.Int
	Code *big.Int
}
