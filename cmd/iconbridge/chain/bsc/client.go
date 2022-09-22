package bsc

import (
	"context"
	"fmt"
	"math"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc/abi/bmcperiphery"
	bscTypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/bsc/types"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/icon-project/icon-bridge/common/log"
)

func newClients(urls []string, bmc string, l log.Logger) (cls []IClient, err error) {
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Errorf("failed to create bsc rpc client: url=%v, %v", url, err)
			return nil, err
		}
		cleth := ethclient.NewClient(clrpc)
		clbmc, err := bmcperiphery.NewBmcperiphery(common.HexToAddress(bmc), cleth)
		if err != nil {
			l.Errorf("failed to create bmc binding to bsc ethclient: url=%v, %v", url, err)
			return nil, err
		}
		cl := &Client{
			log: l,
			rpc: clrpc,
			eth: cleth,
			bmc: clbmc,
		}
		cl.chainID, err = cleth.ChainID(context.Background())
		if err != nil {
			return nil, errors.Wrapf(err, "cleth.ChainID %v", err)
		}
		cls = append(cls, cl)
	}
	return cls, nil
}

// grouped rpc api clients
type Client struct {
	log     log.Logger
	rpc     *rpc.Client
	eth     *ethclient.Client
	chainID *big.Int
	bmc     *bmcperiphery.Bmcperiphery
}

type IClient interface {
	Log() log.Logger
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	GetBlockNumber() (uint64, error)
	GetBlockByHash(hash common.Hash) (*bscTypes.Block, error)
	GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error)
	GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error)
	GetMedianGasPriceForBlock(ctx context.Context) (gasPrice *big.Int, gasHeight *big.Int, err error)
	GetChainID() *big.Int

	// ethClient
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error)

	// bmcClient
	ParseMessage(log ethTypes.Log) (*bmcperiphery.BmcperipheryMessage, error)
	HandleRelayMessage(opts *bind.TransactOpts, _prev string, _msg []byte) (*ethTypes.Transaction, error)
	GetStatus(opts *bind.CallOpts, _link string) (bmcperiphery.TypesLinkStats, error)
}

func (cl *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return cl.eth.NonceAt(ctx, account, blockNumber)
}

func (cl *Client) HandleRelayMessage(opts *bind.TransactOpts, _prev string, _msg []byte) (*ethTypes.Transaction, error) {
	return cl.bmc.HandleRelayMessage(opts, _prev, _msg)
}

func (cl *Client) GetStatus(opts *bind.CallOpts, _link string) (bmcperiphery.TypesLinkStats, error) {
	return cl.bmc.GetStatus(opts, _link)
}

func (cl *Client) GetBMCClient() *bmcperiphery.Bmcperiphery {
	return cl.bmc
}

func (cl *Client) ParseMessage(log ethTypes.Log) (*bmcperiphery.BmcperipheryMessage, error) {
	return cl.bmc.ParseMessage(log)
}

func (cl *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	return cl.eth.TransactionCount(ctx, blockHash)
}

func (cl *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error) {
	return cl.eth.TransactionInBlock(ctx, blockHash, index)
}

func (cl *Client) TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error) {
	return cl.eth.TransactionByHash(ctx, blockHash)
}

func (cl *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {
	return cl.eth.TransactionReceipt(ctx, txHash)
}

func (cl *Client) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return cl.eth.CallContract(ctx, msg, blockNumber)
}

func (cl *Client) GetBalance(ctx context.Context, hexAddr string) (*big.Int, error) {
	if !common.IsHexAddress(hexAddr) {
		return nil, fmt.Errorf("invalid hex address: %v", hexAddr)
	}
	return cl.eth.BalanceAt(ctx, common.HexToAddress(hexAddr), nil)
}

func (cl *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
	return cl.eth.FilterLogs(ctx, q)
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

func (cl *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return cl.eth.SuggestGasPrice(ctx)
}

func (cl *Client) GetBlockByHash(hash common.Hash) (*bscTypes.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	var hb bscTypes.Block
	err := cl.rpc.CallContext(ctx, &hb, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	}
	return &hb, nil
}

func (cl *Client) GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	return cl.eth.HeaderByNumber(ctx, height)
}

func (cl *Client) GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error) {
	hb, err := cl.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	if hb.GasUsed == "0x0" || len(hb.Transactions) == 0 {
		return nil, nil
	}
	txhs := hb.Transactions
	// fetch all txn receipts concurrently
	type rcq struct {
		txh   string
		v     *ethTypes.Receipt
		err   error
		retry int
	}
	qch := make(chan *rcq, len(txhs))
	for _, txh := range txhs {
		qch <- &rcq{txh, nil, nil, RPCCallRetry}
	}
	rmap := make(map[string]*ethTypes.Receipt)
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
					q.v = &ethTypes.Receipt{}
				}
				q.v, err = cl.TransactionReceipt(ctx, common.HexToHash(q.txh))
				if q.err != nil {
					q.err = errors.Wrapf(q.err, "getTranasctionReceipt: %v", q.err)
				}
			}(q)
		}
	}
	receipts := make(ethTypes.Receipts, 0, len(txhs))
	for _, txh := range txhs {
		if r, ok := rmap[txh]; ok {
			receipts = append(receipts, r)
		}
	}
	return receipts, nil
}

func (c *Client) GetMedianGasPriceForBlock(ctx context.Context) (gasPrice *big.Int, gasHeight *big.Int, err error) {
	gasPrice = big.NewInt(0)
	header, err := c.GetHeaderByHeight(ctx, nil)
	if err != nil {
		err = errors.Wrapf(err, "GetHeaderByNumber(height:latest) Err: %v", err)
		return
	}
	height := header.Number
	txnCount, err := c.TransactionCount(ctx, header.Hash())
	if err != nil {
		err = errors.Wrapf(err, "GetTransactionCount(height:%v, headerHash: %v) Err: %v", height, header.Hash(), err)
		return
	} else if err == nil && txnCount == 0 {
		return nil, nil, fmt.Errorf("TransactionCount is zero for height(%v, headerHash %v)", height, header.Hash())
	}
	// txnF, err := c.eth.TransactionInBlock(ctx, header.Hash(), 0)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, 0, err)
	// }
	txnS, err := c.TransactionInBlock(ctx, header.Hash(), uint(math.Floor(float64(txnCount)/2)))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, txnCount-1, err)
	}

	gasPrice = txnS.GasPrice()
	gasHeight = header.Number
	return
}

func (c *Client) GetChainID() *big.Int {
	return c.chainID
}

func (c *Client) GetEthClient() *ethclient.Client {
	return c.eth
}

func (c *Client) Log() log.Logger {
	return c.log
}
