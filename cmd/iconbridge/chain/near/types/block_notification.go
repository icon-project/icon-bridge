package types

import "github.com/icon-project/icon-bridge/cmd/iconbridge/chain"

type BlockNotification struct {
	offset          int64
	block           Block
	receipts        []*chain.Receipt
}

func NewBlockNotification(offset int64) *BlockNotification {
	return &BlockNotification{
		offset: offset,
	}
}

func (bn *BlockNotification) Offset() int64 {
	return bn.offset
}

func (bn *BlockNotification) Block() *Block {
	return &bn.block
}

func (bn *BlockNotification) Receipts() []*chain.Receipt {
	return bn.receipts
}

func (bn *BlockNotification) SetBlock(block Block) {
	bn.block = block
}

func (bn *BlockNotification) SetReceipts(receipts []*chain.Receipt) {
	bn.receipts = receipts
}
