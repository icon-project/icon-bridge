package chain

import (
	"context"
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

type CoinBalance struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
	UserBalance       *big.Int
}

func (cb *CoinBalance) String() string {
	return "Usable " + cb.UsableBalance.String() +
		" Locked " + cb.LockedBalance.String() + " Refundable " + cb.RefundableBalance.String() +
		" UserBalance " + cb.UserBalance.String()
}

type SrcAPI interface {
	Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error)
	TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error)
	WaitForTxnResult(ctx context.Context, hash string) (txnr *TxnResult, err error)
	WatchForTransferStart(requestID uint64, seq int64) error
	WatchForTransferEnd(ID uint64, seq int64) error
	Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error)
	GetCoinBalance(coinName string, addr string) (*CoinBalance, error)

	NativeCoin() string
	NativeTokens() []string
	GetBTPAddress(addr string) string
}

type DstAPI interface {
	GetCoinBalance(coinName string, addr string) (*CoinBalance, error)
	WatchForTransferReceived(requestID uint64, seq int64) error
	GetBTPAddress(addr string) string
	NativeTokens() []string
}

type TxnResult struct {
	StatusCode int
	ElInfo     []*EventLogInfo
	Raw        interface{}
}

type ChainAPI interface {
	Subscribe(ctx context.Context) (sinkChan chan *EventLogInfo, errChan chan error, err error)
	GetKeyPairs(num int) ([][2]string, error)
	GetKeyPairFromKeystore(keystore, secret string) (string, string, error)

	TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error)
	Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error)
	WaitForTxnResult(ctx context.Context, hash string) (txnr *TxnResult, err error)
	WatchForTransferStart(ID uint64, seq int64) error
	WatchForTransferReceived(ID uint64, seq int64) error
	WatchForTransferEnd(ID uint64, seq int64) error
	Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error)
	GetCoinBalance(coinName string, addr string) (*CoinBalance, error)

	NativeCoin() string
	NativeTokens() []string
	GetBTPAddress(addr string) string
}

type Config struct {
	Name                  ChainType               `json:"name"`
	URL                   string                  `json:"url"`
	ContractAddresses     map[ContractName]string `json:"contract_addresses"`
	NativeCoin            string                  `json:"native_coin"`
	NativeTokens          []string                `json:"native_tokens"`
	WrappedCoins          []string                `json:"wrapped_coins"`
	GodWalletKeystorePath string                  `json:"god_wallet_keystore_path"`
	GodWalletSecretPath   string                  `json:"god_wallet_secret_path"`
	NetworkID             string                  `json:"network_id"`
	GasLimit              int64                   `json:"gasLimit"`
}

type EventLogInfo struct {
	IDs             []uint64
	ContractAddress string
	EventType       EventLogType
	EventLog        interface{}
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
	From     string
	Sn       *big.Int
	Code     *big.Int
	Response string
}
