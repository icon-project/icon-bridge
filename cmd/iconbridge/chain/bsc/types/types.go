package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Block struct {
	Transactions []string `json:"transactions"`
	GasUsed      string   `json:"gasUsed"`
}

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}

type BlockNotification struct {
	Hash          common.Hash
	Height        *big.Int
	Header        *types.Header
	Receipts      types.Receipts
	HasBTPMessage *bool
}