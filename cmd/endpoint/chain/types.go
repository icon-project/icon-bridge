package chain

import (
	"context"
	"errors"
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
	TokenBSHImplHmy       ContractName = "TokenBSHImplHmy" //TokenHmy
	NativeBSHPeripheryHmy ContractName = "NativeBSHPeripheryHmy"
	Erc20Hmy              ContractName = "Erc20Hmy"
	NativeBSHCoreHmy      ContractName = "NativeBSHCoreHmy"
	TokenBSHProxyHmy      ContractName = "TokenBSHProxyHmy"
	TokenBSHIcon          ContractName = "TokenBSHIcon"
	NativeBSHIcon         ContractName = "NativeBSHIcon"
	Irc2Icon              ContractName = "Irc2Icon"
	Irc2TradeableIcon     ContractName = "Irc2TradeableIcon"
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

type ChainAPI interface {
	Subscribe(ctx context.Context, height uint64) (sinkChan chan *EventLogInfo, errChan chan error, err error)
	Transfer(param *RequestParam) (txnHash string, err error)
	GetCoinBalance(addr string, coinType TokenType) (*big.Int, error)
	WaitForTxnResult(hash string) (txr interface{}, elInfo []*EventLogInfo, err error)
	Approve(ownerKey string, amount big.Int) (txnHash string, err error)
	GetBTPAddress(addr string) *string
	GetKeyPairs(num int) ([][2]string, error)
	WatchFor(ID uint64, eventType EventLogType, seq int64, contractAddress string) error
}

type ChainConfig struct {
	Name               ChainType               `json:"name"`
	URL                string                  `json:"url"`
	ConftractAddresses map[ContractName]string `json:"contract_addresses"`
	GodWallet          GodWallet               `json:"god_wallet"`
	NetworkID          string                  `json:"network_id"`
	Src                BTPAddress              `json:"src"`
	Dst                BTPAddress              `json:"dst"`
}

// type ContractAddress struct {
// 	Name    string `json:"name"`
// 	Address string `json:"address"`
// }

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
			err = errors.New("EventLg type is not *TransferStartEvent")
		}
		seq = st.Sn.Int64()
	} else if e.EventType == TransferReceived {
		st, ok := e.EventLog.(*TransferReceivedEvent)
		if !ok {
			err = errors.New("EventLg type is not *TransferReceivedEvent")
		}
		seq = st.Sn.Int64()
	} else if e.EventType == TransferEnd {
		st, ok := e.EventLog.(*TransferEndEvent)
		if !ok {
			err = errors.New("EventLg type is not *TransferEndEvent")
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
