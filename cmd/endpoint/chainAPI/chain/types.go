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

type RequestAPI interface {
	GetCoinBalance(addr string) (*big.Int, error)
	GetEthToken(addr string) (val *big.Int, err error)
	GetWrappedCoin(addr string) (val *big.Int, err error)
	TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
	TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
	TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
	TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, logs interface{}, err error)
	TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, approveLogs interface{}, transferTxnHash string, transferLogs interface{}, err error)
	ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, logs interface{}, allowanceAmount *big.Int, err error)
	GetAddressFromPrivKey(key string) (*string, error)
	GetBTPAddress(addr string) *string
}

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

type EnvVariables struct {
	Client       RequestAPI
	GodKeys      [2]string
	AccountsKeys [][2]string
}

type EventLog struct {
	Addr    string   `json:"scoreAddress"`
	Indexed []string `json:"indexed"`
	Data    []string `json:"data"`
}

type SubscribedEvent struct {
	Res       interface{}
	ChainName ChainType
}

type Event struct {
	Next     BTPAddress
	Sequence uint64
	Message  []byte
}

type Receipt struct {
	Index  uint64
	Events []*Event
	Height uint64
}

type SubscribeOptions struct {
	Seq    uint64
	Height uint64
}
