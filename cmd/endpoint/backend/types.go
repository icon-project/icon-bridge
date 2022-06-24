package backend

import (
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
)

type TokenType string

const (
	ICXToken        TokenType = "ICX"
	IRC2Token       TokenType = "IRC2"
	ONEWrappedToken TokenType = "OneWrapped"
	ONEToken        TokenType = "ONE"
	ERC20Token      TokenType = "ERC20"
	ICXWrappedToken TokenType = "ICXWrapped"
)

// type ChainType string

// const (
// 	ICONChain ChainType = "ICON"
// 	HMNYChain ChainType = "HMNY"
// )

type RequestParam struct {
	FromChain chain.ChainType
	ToChain   chain.ChainType
	SenderKey string
	ToAddress string
	Amount    big.Int
	Token     TokenType
}

type ApproveParam struct {
	Chain       chain.ChainType
	OwnerKey    string
	Amount      big.Int
	AccessToken TokenType
}

type ReceiptEvent struct {
}

type Response struct {
}

type RequestMatchingFn func(req *RequestParam) bool
type EventMatchingFn func(*ReceiptEvent, *RequestParam) bool
type MatchingFn struct {
	ReqFn   RequestMatchingFn
	EventFn EventMatchingFn
}

type DecodedEvent struct {
	TxHash      string           `json:"txHash" validate:"required,t_int"`
	From        string           `json:"from"`
	To          string           `json:"to"`
	EventLogs   []chain.EventLog `json:"eventLogs"`
	Status      int64            `json:"status"`
	BlockHeight int64            `json:"blockHeight"`
}
