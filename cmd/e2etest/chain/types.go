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
	BTS          ContractName = "BTS"
	BTSPeriphery ContractName = "BTSPeriphery"
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

// type ContractCallMethodName string

// const (

// )

// type ContractTransactMethodName string

// const (
// 	SetTokenLimit          ContractTransactMethodName = "SetTokenLimit"
// 	AddBlackListAddress    ContractTransactMethodName = "AddBlackListAddress"
// 	RemoveBlackListAddress ContractTransactMethodName = "RemoveBlackListAddress"
// 	AddRestriction         ContractTransactMethodName = "AddRestriction"
// 	DisableRestrictions    ContractTransactMethodName = "DisableRestrictions"
// )

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
	Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error)

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
	// Subscription
	Subscribe(ctx context.Context) (sinkChan chan *EventLogInfo, errChan chan error, err error)

	// Account
	GetKeyPairs(num int) ([][2]string, error)
	GetKeyPairFromKeystore(keystoreFile, secretFile string) (string, string, error)

	// Transfer
	TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error)
	Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error)
	WaitForTxnResult(ctx context.Context, hash string) (txnr *TxnResult, err error)
	WatchForTransferStart(ID uint64, seq int64) error
	WatchForTransferReceived(ID uint64, seq int64) error
	WatchForTransferEnd(ID uint64, seq int64) error
	Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error)
	GetCoinBalance(coinName string, addr string) (*CoinBalance, error)
	Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error)

	// Query
	NativeCoin() string
	NativeTokens() []string
	GetBTPAddress(addr string) string

	// Configure
	SetTokenLimit(ownerKey string, coinNames []string, tokenLimits []*big.Int) (txnHash string, err error)
	AddBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error)
	RemoveBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error)
	ChangeRestriction(ownerKey string, enable bool) (txnHash string, err error)
	IsUserBlackListed(net, addr string) (response bool, err error)
	GetTokenLimit(coinName string) (tokenLimit *big.Int, err error)
	IsBTSOwner(addr string) (response bool, err error)
	GetTokenLimitStatus(net, coinName string) (response bool, err error)
	GetBlackListedUsers(net string, startCursor, endCursor int) (addrs []string, err error)
}

type Config struct {
	Name                   ChainType               `json:"name"`
	URL                    string                  `json:"url"`
	ContractAddresses      map[ContractName]string `json:"contract_addresses"`
	NativeCoin             string                  `json:"native_coin"`
	NativeTokens           []string                `json:"native_tokens"`
	WrappedCoins           []string                `json:"wrapped_coins"`
	GodWalletKeystorePath  string                  `json:"god_wallet_keystore_path"`
	GodWalletSecretPath    string                  `json:"god_wallet_secret_path"`
	DemoWalletKeystorePath string                  `json:"demo_wallet_keystore_path"`
	NetworkID              string                  `json:"network_id"`
	GasLimit               int64                   `json:"gas_limit"`
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

// type AssetDetails struct {
// 	Name  string
// 	Value *big.Int
// }

type TransferEndEvent struct {
	From     string
	Sn       *big.Int
	Code     *big.Int
	Response string
}
