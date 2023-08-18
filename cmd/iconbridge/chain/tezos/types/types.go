package types 

import (
	"math/big"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	

)

type BlockNotification struct {
	Hash          	tezos.BlockHash
	Height        	*big.Int
	Header        	*rpc.BlockHeader
	Receipts      	[]*chain.Receipt
	HasBTPMessage 	*bool
	Proposer 		tezos.Address
	Block 			*rpc.Block
	BakerConsensusKey tezos.Key
}

