package chain

import (
	"context"
	"math/big"
)

type ChainType string

const (
	ICON ChainType = "ICON"
	HMNY ChainType = "HMNY"
)

type TokenType string

const (
	ICXToken   TokenType = "ICX"
	IRC2Token  TokenType = "IRC2"
	ONEToken   TokenType = "ONE"
	ERC20Token TokenType = "ERC20"
)

type ContractName string

const (
	TokenHmy   ContractName = "TokenHmy"
	NativeHmy  ContractName = "NativeHmy"
	TokenIcon  ContractName = "TokenIcon"
	NativeIcon ContractName = "NativeIcon"
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

type RequestParam struct {
	FromChain   ChainType
	ToChain     ChainType
	SenderKey   string
	FromAddress string
	ToAddress   string
	Amount      big.Int
	Token       TokenType
}

// type RequestAPI interface {
// 	GetCoinBalance(addr string) (*big.Int, error)
// 	GetEthToken(addr string) (val *big.Int, err error)
// 	GetWrappedCoin(addr string) (val *big.Int, err error)
// 	TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
// 	TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
// 	TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
// 	TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
// 	TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, approveLogs interface{}, transferTxnHash string, transferLogs interface{}, err error)
// 	ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, logs interface{}, allowanceAmount *big.Int, err error)
// 	GetAddressFromPrivKey(key string) (*string, error)
// 	GetBTPAddress(addr string) *string
// }
// type EnvVariables struct {
// 	Client       RequestAPI
// 	GodKeys      [2]string
// 	AccountsKeys [][2]string
// }

type SubscriptionAPI interface {
	Start(ctx context.Context) error
	OutputChan() <-chan *SubscribedEvent
	ErrChan() <-chan error
}

type ChainConfig struct {
	Name               ChainType         `json:"name"`
	URL                string            `json:"url"`
	ConftractAddresses map[string]string `json:"contract_addresses"`
	GodWallet          GodWallet         `json:"god_wallet"`
	NetworkID          string            `json:"network_id"`
	Subscriber         SubscriberConfig  `json:"subscriber"`
}

type SubscriberConfig struct {
	Src  BTPAddress             `json:"src"`
	Dst  BTPAddress             `json:"dst"`
	Opts map[string]interface{} `json:"options"`
}

type ContractAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type GodWallet struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

type SubscribedEvent struct {
	Res       interface{}
	ChainName ChainType
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
