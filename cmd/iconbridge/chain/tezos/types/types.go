package types 

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"blockwatch.cc/tzgo/rpc"


)
type BlockNotification struct {
	Hash          common.Hash
	Height        *big.Int
	Header        *types.Header
	Receipts      []rpc.Receipt
	HasBTPMessage *bool
}