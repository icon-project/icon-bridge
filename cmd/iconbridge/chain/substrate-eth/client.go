package substrate_eth

import (
	"context"
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/substrate-eth/abi"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	subEthTypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/substrate-eth/types"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/icon-project/icon-bridge/common/log"
)

func newClients(urls []string, bmc string, l log.Logger) (cls []IClient, bmcs []*abi.BMC, err error) {
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Errorf("failed to create snow rpc client: url=%v, %v", url, err)
			return nil, nil, err
		}
		cleth := ethclient.NewClient(clrpc)
		clbmc, err := abi.NewBMC(common.HexToAddress(bmc), cleth)
		if err != nil {
			l.Errorf("failed to create bmc binding to snow ethclient: url=%v, %v", url, err)
			return nil, nil, err
		}
		bmcs = append(bmcs, clbmc)
		cl := &Client{
			log: l,
			rpc: clrpc,
			eth: cleth,
		}
		cl.chainID, err = cleth.ChainID(context.Background())
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cleth.ChainID %v", err)
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

type IClient interface {
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	GetBlockNumber() (uint64, error)
	//GetBlockByHash(hash common.Hash) (*snowTypes.Block, error)
	GetHeaderByHeight(height *big.Int) (*subEthTypes.Header, error)
	//GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error)
	GetBlockReceiptsFromHeight(height *big.Int) (ethTypes.Receipts, bool, error)
	GetChainID() *big.Int
	GetEthClient() *ethclient.Client
	Log() log.Logger
}

func (cl *Client) GetBalance(ctx context.Context, hexAddr string) (*big.Int, error) {
	if !common.IsHexAddress(hexAddr) {
		return nil, fmt.Errorf("invalid hex address: %v", hexAddr)
	}
	return cl.eth.BalanceAt(ctx, common.HexToAddress(hexAddr), nil)
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

func (cl *Client) GetHeaderByHeight(height *big.Int) (*subEthTypes.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	return headerByNumber(cl, ctx, height)
}

func headerByNumber(cl *Client, ctx context.Context, number *big.Int) (*subEthTypes.Header, error) {
	var head *subEthTypes.Header
	err := cl.rpc.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)

	if err == nil && head == nil {
		err = ethereum.NotFound
	}

	return head, err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func (c *Client) GetMedianGasPriceForBlock(ctx context.Context) (gasPrice *big.Int, gasHeight *big.Int, err error) {
	gasPrice = big.NewInt(0)
	header, err := c.eth.HeaderByNumber(ctx, nil)
	if err != nil {
		err = errors.Wrapf(err, "GetHeaderByNumber(height:latest) Err: %v", err)
		return
	}
	height := header.Number
	txnCount, err := c.eth.TransactionCount(ctx, header.Hash())
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
	txnS, err := c.eth.TransactionInBlock(ctx, header.Hash(), uint(math.Floor(float64(txnCount)/2)))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, txnCount-1, err)
	}
	gasPrice = txnS.GasPrice()
	gasHeight = header.Number
	return
}

func (cl *Client) GetBlockReceiptsFromHeight(height *big.Int) (ethTypes.Receipts, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb, err := cl.eth.BlockByNumber(ctx, height)
	if err != nil {
		return nil, false, err
	}
	if hb.GasUsed() == 0 || len(hb.Transactions()) == 0 {
		return nil, false, nil
	}
	txhs := []string{}
	isEIP1559 := false
	for _, v := range hb.Transactions() {
		if v.Type() == 2 {
			isEIP1559 = true
		}
		txhs = append(txhs, v.Hash().String())
	}
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
				return nil, isEIP1559, q.err
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
				q.v, err = cl.eth.TransactionReceipt(ctx, common.HexToHash(q.txh))
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
	return receipts, isEIP1559, nil
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
