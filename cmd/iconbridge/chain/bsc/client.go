package bsc

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

func newClients(urls []string, bmc string, l log.Logger) (cls []*Client, bmcs []*BMC, err error) {
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Errorf("failed to create hmny rpc client: url=%v, %v", url, err)
			return nil, nil, err
		}
		cleth := ethclient.NewClient(clrpc)
		clbmc, err := NewBMC(common.HexToAddress(bmc), cleth)
		if err != nil {
			l.Errorf("failed to create bmc binding to hmny ethclient: url=%v, %v", url, err)
			return nil, nil, err
		}
		bmcs = append(bmcs, clbmc)
		cl := &Client{
			log: l,
			rpc: clrpc,
			eth: cleth,
			//bmc: clbmc,
		}
		cl.chainID, err = cl.GetChainID()
		if err != nil {
			return nil, nil, errors.Wrapf(err, "GetChainID %v", err)
		}
		cls = append(cls, cl)
	}
	return cls, bmcs, nil
}

// grouped rpc api clients
type Client struct {
	log     log.Logger
	rpc     *rpc.Client
	eth     *ethclient.Client
	chainID *big.Int
	//bmc *BMC
}

func (cl *Client) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	bn, err := cl.eth.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

func (cl *Client) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	return cl.eth.BlockByHash(ctx, hash)
}

func (cl *Client) GetHeaderByHeight(height *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	return cl.eth.HeaderByNumber(ctx, height)
}

func (cl *Client) GetBlockReceipts(hash common.Hash) (types.Receipts, error) {
	txnCount, err := cl.eth.TransactionCount(context.TODO(), hash)
	if err != nil {
		return nil, err
	} else if txnCount == 0 {
		return nil, nil
	}
	type rcq struct {
		txIdx int
		v     *types.Receipt
		err   error
		retry int
	}
	qch := make(chan *rcq, txnCount)
	for i := 0; i < int(txnCount); i++ {
		qch <- &rcq{i, nil, nil, 3}
	}
	rmap := make(map[int]*types.Receipt)
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
			rmap[q.txIdx] = q.v
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
				txn, err := cl.eth.TransactionInBlock(ctx, hash, uint(q.txIdx))
				if err != nil {
					q.err = errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, Index: %v) Err: %v", hash, q.txIdx, err)
				} else {
					q.v, err = cl.eth.TransactionReceipt(ctx, txn.Hash())
					if q.err != nil {
						q.err = errors.Wrapf(q.err, "getTranasctionReceipt %v", q.err)
					}
				}

			}(q)
		}
	}
	receipts := make(types.Receipts, 0, txnCount)
	for i := 0; i < int(txnCount); i++ {
		if r, ok := rmap[i]; ok {
			receipts = append(receipts, r)
		}
	}
	return receipts, nil
}

func (c *Client) GetMedianGasPriceForBlock() (gasPrice *big.Int, gasHeight *big.Int, err error) {
	gasPrice = big.NewInt(0)
	header, err := c.eth.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		err = errors.Wrapf(err, "GetHeaderByNumber(height:latest) Err: %v", err)
		return
	}
	height := header.Number
	txnCount, err := c.eth.TransactionCount(context.TODO(), header.Hash())
	if err != nil {
		err = errors.Wrapf(err, "GetTransactionCount(height:%v, headerHash: %v) Err: %v", height, header.Hash(), err)
		return
	} else if err == nil && txnCount == 0 {
		return nil, nil, fmt.Errorf("TransactionCount is zero for height(%v, headerHash %v)", height, header.Hash())
	}
	// txnF, err := c.eth.TransactionInBlock(context.TODO(), header.Hash(), 0)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, 0, err)
	// }
	txnS, err := c.eth.TransactionInBlock(context.TODO(), header.Hash(), uint(math.Floor(float64(txnCount)/2)))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, txnCount-1, err)
	}
	gasPrice = txnS.GasPrice()
	gasHeight = header.Number
	return
}

func (c *Client) newTransactOpts(w Wallet) (*bind.TransactOpts, error) {
	txo, err := bind.NewKeyedTransactorWithChainID(w.(*wallet.EvmWallet).Skey, c.chainID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	txo.GasPrice, _ = c.eth.SuggestGasPrice(ctx)
	txo.GasLimit = uint64(DefaultGasLimit)
	return txo, nil
}

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}

func (c *Client) GetChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	return c.eth.ChainID(ctx)
}
