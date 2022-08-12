//go:build hmny
// +build hmny

package hmny

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewClients(urls []string, l log.Logger) (cls []*Client, err error) {
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Errorf("failed to create hmny rpc client: url=%v, %v", url, err)
			return nil, err
		}
		cleth := ethclient.NewClient(clrpc)
		cls = append(cls, &Client{
			log: l,
			rpc: clrpc,
			eth: cleth,
			//bmc: clbmc,
		})
	}
	return cls, nil
}

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
		cls = append(cls, &Client{
			log: l,
			rpc: clrpc,
			eth: cleth,
			//bmc: clbmc,
		})
	}
	return cls, bmcs, nil
}

// grouped rpc api clients
type Client struct {
	log log.Logger
	rpc *rpc.Client
	eth *ethclient.Client
	//bmc *BMC
}

func (cl *Client) newVerifier(opts *VerifierOptions) (Verifier, error) {
	h, err := cl.GetHmyV2HeaderByHeight((&big.Int{}).SetUint64(opts.BlockHeight))
	if err != nil {
		return nil, errors.Wrapf(err, "cl.GetHeaderByHeight(%d): %v", opts.BlockHeight, err)
	}
	ssh := h // shard state header
	if ssh.Epoch.Cmp(bigZero) <= 0 {
		if ssh.Number.Cmp(bigZero) > 0 {
			ssh, err = cl.GetHmyV2HeaderByHeight(bigZero)
			if err != nil {
				return nil, errors.Wrapf(err, "cl.GetHeaderByHeight(%d): %v", 0, err)
			}
		}
	} else {
		epoch := new(big.Int).Sub(ssh.Epoch, bigOne)
		elb, err := cl.GetEpochLastBlock(epoch)
		if err != nil {
			return nil, errors.Wrapf(err, "cl.GetEpochLastBlock(%d): %v", epoch, err)
		}
		ssh, err = cl.GetHmyV2HeaderByHeight(elb)
		if err != nil {
			return nil, errors.Wrapf(err, "cl.GetHeaderByHeight(%d): %v", elb, err)
		}
	}
	vr := NewVerifier()
	if err = vr.Update(ssh); err != nil {
		return nil, errors.Wrapf(err, "verifier.Update: %v", err)
	}
	ok, err := vr.Verify(h, opts.CommitBitmap, opts.CommitSignature)
	if !ok || err != nil {
		return nil, errors.Wrapf(err, "invalid signature: %v", err)
	}
	return vr, nil
}

func (cl *Client) syncVerifier(vr Verifier, height uint64) (err error) {
	h, err := cl.GetHmyV2HeaderByHeight((&big.Int{}).SetUint64(height))
	if err != nil {
		return err
	}
	cl.log.WithFields(log.Fields{"epoch": vr.Epoch()}).Debugf("syncVerifier: start")
	for epoch := vr.Epoch(); epoch < h.Epoch.Uint64(); epoch++ {
		elb, err := cl.GetEpochLastBlock((&big.Int{}).SetUint64(epoch))
		if err != nil {
			return errors.Wrapf(err, "cl.GetEpochLastBlock: %v", err)
		}
		elh, err := cl.GetHmyV2HeaderByHeight(elb)
		if err != nil {
			return errors.Wrapf(err, "cl.GetHmyHeaderByHeight(elh): %v", err)
		}
		elx, err := cl.GetHmyV2HeaderByHeight((&big.Int{}).Add(elb, bigOne))
		if err != nil {
			return errors.Wrapf(err, "cl.GetHmyHeaderByHeight(elx): %v", err)
		}
		ok, err := vr.Verify(elh, elx.LastCommitBitmap, elx.LastCommitSignature)
		if !ok || err != nil {
			return errors.Wrapf(err, "vr.Verify: invalid signature: %v", err)
		}
		if err = vr.Update(elh); err != nil {
			return errors.Wrapf(err, "vr.Update: %v", err)
		}
		cl.log.WithFields(log.Fields{
			"epoch": vr.Epoch(), "height": elb.Uint64()}).Debugf("syncVerifier: syncing")
	}
	cl.log.WithFields(log.Fields{"epoch": vr.Epoch()}).Debugf("syncVerifier: complete")
	return nil
}

func (cl *Client) GetChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	return cl.eth.ChainID(ctx)
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

func (cl *Client) GetTransaction(hash common.Hash) (*ethtypes.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	tx, pending, err := cl.eth.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, pending, err
	}
	return tx, pending, err
}

func (cl *Client) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	tr := new(types.Receipt)
	err := cl.rpc.CallContext(ctx, tr, "hmy_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (cl *Client) GetEpochLastBlock(epoch *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	lbn := big.NewInt(0)
	err := cl.rpc.CallContext(ctx, lbn, "hmyv2_epochLastBlock", epoch)
	if err != nil {
		return nil, err
	}
	return lbn, nil
}

func (cl *Client) GetHmyBlockByHeight(height *big.Int) (*BlockWithTxHash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(BlockWithTxHash)
	err := cl.rpc.CallContext(ctx, hb, "hmy_getBlockByNumber", height, false)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *Client) GetHmyV2BlockByHeight(height *big.Int) (*BlockV2WithTxHash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(BlockV2WithTxHash)
	err := cl.rpc.CallContext(ctx, hb, "hmyv2_getBlockByNumber", height, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *Client) GetHmyBlockByHash(hash common.Hash) (*BlockWithTxHash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(BlockWithTxHash)
	err := cl.rpc.CallContext(ctx, hb, "hmy_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *Client) GetHmyV2BlockByHash(hash common.Hash) (*BlockV2WithTxHash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(BlockV2WithTxHash)
	err := cl.rpc.CallContext(ctx, hb, "hmyv2_getBlockByHash", hash, map[string]interface{}{"inclStaking": true})
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *Client) GetHmyV2HeaderByHeight(height *big.Int) (*Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	hb := new(Header)
	err := cl.rpc.CallContext(ctx, hb, "hmyv2_getFullHeader", height)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *Client) GetBlockReceipts(hash common.Hash) (types.Receipts, error) {
	receipts, err := cl.getHmyBlockReceipts(hash)
	if err != nil {
		return cl.getHmyTxnReceiptsByBlockHash(hash)
	}
	return receipts, nil
}

func (cl *Client) getHmyBlockReceipts(hash common.Hash) (types.Receipts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	receipts := make([]*types.Receipt, 0)
	err := cl.rpc.CallContext(ctx, &receipts, "hmy_getBlockReceipts", hash)
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func (cl *Client) getHmyTxnReceiptsByBlockHash(hash common.Hash) (types.Receipts, error) {
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
		qch <- &rcq{txh, nil, nil, RPCCallRetry}
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
