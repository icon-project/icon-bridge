package hmny

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/bmr/common/errors"
	"github.com/icon-project/icon-bridge/bmr/common/log"
)

func newClients(urls []string, bmc string, l log.Logger) (cls []*client, err error) {
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Errorf("failed to create hmny rpc client: url=%v, %v", url, err)
			return nil, err
		}
		cleth := ethclient.NewClient(clrpc)
		cls = append(cls, &client{
			log: l,
			rpc: clrpc,
			eth: cleth,
		})
	}
	return cls, nil
}

// grouped rpc api clients
type client struct {
	log log.Logger
	rpc *rpc.Client
	eth *ethclient.Client
}

func (cl *client) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	bn, err := cl.eth.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

func (cl *client) GetHmyBlockByHash(hash common.Hash) (*BlockWithTxHash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(BlockWithTxHash)
	err := cl.rpc.CallContext(ctx, hb, "hmy_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *client) GetHmyV2HeaderByHeight(height *big.Int) (*Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(Header)
	err := cl.rpc.CallContext(ctx, hb, "hmyv2_getFullHeader", height)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *client) GetBlockReceipts(hash common.Hash) (types.Receipts, error) {
	receipts, err := cl.getHmyBlockReceipts(hash)
	if err != nil {
		return cl.getHmyTxnReceiptsByBlockHash(hash)
	}
	return receipts, nil
}

func (cl *client) getHmyBlockReceipts(hash common.Hash) (types.Receipts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	receipts := make([]*types.Receipt, 0)
	err := cl.rpc.CallContext(ctx, &receipts, "hmy_getBlockReceipts", hash)
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func (cl *client) getHmyTxnReceiptsByBlockHash(hash common.Hash) (types.Receipts, error) {
	b, err := cl.GetHmyBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	txhs := append(b.Transactions, b.StakingTxs...)
	if b.GasUsed == 0 || len(txhs) == 0 {
		return nil, nil
	}
	// fetch all txn receipts concurrently
	type rcq struct {
		txh   common.Hash
		v     *types.Receipt
		err   error
		retry int
	}
	qch := make(chan *rcq, len(txhs))
	for _, txh := range txhs {
		qch <- &rcq{txh, nil, nil, 3}
	}
	rmap := make(map[common.Hash]*types.Receipt)
	for q := range qch {
		switch {
		case q.err != nil:
			if q.retry == 0 {
				return nil, q.err
			}
			q.retry--
			q.err = nil
			qch <- q
		case q.v != nil:
			rmap[q.txh] = q.v
			if len(rmap) == cap(qch) {
				close(qch)
			}
		default:
			go func(q *rcq) {
				defer func() { qch <- q }()
				ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
				defer cancel()
				if q.v == nil {
					q.v = &types.Receipt{}
				}
				q.err = cl.rpc.CallContext(ctx, q.v, "hmy_getTransactionReceipt", q.txh)
				if q.err != nil {
					q.err = errors.Wrapf(q.err, "hmy_getTransactionReceipt: %v", q.err)
				}
			}(q)
		}
	}
	receipts := make(types.Receipts, 0, len(txhs))
	for _, txh := range txhs {
		if r, ok := rmap[txh]; ok {
			receipts = append(receipts, r)
		}
	}
	return receipts, nil
}
