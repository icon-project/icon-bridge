package near

import (
	"context"
	"fmt"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/reactivex/rxgo/v2"
)

type client struct {
	near.IClient
	concurrency int
	MonitorClosed bool
}

type BlockNotification struct {
	*types.BlockNotification
	transactions []types.Transaction
}

func NewBlockNotification(offset int64) BlockNotification {
	return BlockNotification{
		BlockNotification: types.NewBlockNotification(offset),
	}
}

func (bn *BlockNotification) AddTransactions(transactions []types.Transaction) {
	bn.transactions = append(bn.transactions, transactions...)
}

func (c *client) getChunk(id types.CryptoHash) (types.ChunkHeader, error) {
	param := struct {
		ChunkId string `json:"chunk_id"`
	}{
		ChunkId: id.Base58Encode(),
	}

	chuckHeaderResponse, err := c.Api().Chunk(param)
	if err != nil {
		return types.ChunkHeader{}, err
	}

	return chuckHeaderResponse, nil
}

func (c *client) fetchBlockTransactions(context context.Context, bn interface{}) (interface{}, error) {
	if bn, Ok := (bn).(*BlockNotification); Ok {
		if *bn.Block().Hash() != [32]byte{} {
			for _, chunk := range bn.Block().Chunks {
				chunk, err := c.getChunk(chunk.ChunkHash)
				if err != nil {
					return nil, fmt.Errorf("failed to get chunk for block %v", bn.Offset())
				}

				bn.AddTransactions(chunk.Transactions)
			}
		}
		return bn, nil
	}

	return nil, fmt.Errorf("expected BlockNotification but received %v", bn)
}

func (c *client) MonitorTransactions(height uint64, callback func(rxgo.Observable) error) error {
	return callback(c.MonitorBlockHeight(int64(height)).Map(func(_ context.Context, offset interface{}) (interface{}, error) {
		if offset, Ok := (offset).(int64); Ok {
			block, err := c.GetBlockByHeight(offset)
			bn := NewBlockNotification(offset)

			if err != nil {
				return &bn, nil // TODO: Handle Error
			}

			bn.SetBlock(block)
			return &bn, nil
		}
		return nil, fmt.Errorf("error casting offset to int64")
	}, rxgo.WithPool(c.concurrency)).Map(c.fetchBlockTransactions,
		rxgo.WithPool(c.concurrency)).Serialize(int(height), func(_bn interface{}) int {
		bn := _bn.(*BlockNotification)

		return int(bn.Offset())
	}, rxgo.WithPool(c.concurrency), rxgo.WithErrorStrategy(rxgo.ContinueOnError)).TakeUntil(func(i interface{}) bool {
		return c.IsMonitorClosed()
	}))
}

func (c *client) IsMonitorClosed() bool {
	return c.MonitorClosed
}

func (c *client) CloseMonitor() {
	c.MonitorClosed = true
}

func NewClient(url string, logger log.Logger) (*client, error) {
	c, err := near.NewClient(url, logger)
	if err != nil {
		return nil, err
	}
	return &client{
		IClient:     c,
		concurrency: 100,
	}, nil
}
