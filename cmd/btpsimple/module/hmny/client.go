package hmny

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/btp/common/errors"
	"github.com/icon-project/btp/common/log"
	"github.com/icon-project/btp/common/wallet"
)

const (
	BlockRetryInterval               = time.Second * 2
	DefaultReadTimeout               = 15 * time.Second
	BlockNotificationSyncConcurrency = 100 // number of concurrent requests to synchronize older blocks from source chain
)

type Client struct {
	log            log.Logger
	stopAllMonitor chan struct{}

	rpcs []*rpc.Client
	eths []*ethclient.Client
	bmcs []*BMC
}

func NewClient(urls []string, bmcAddr string, l log.Logger) *Client {
	cl := &Client{log: l, stopAllMonitor: make(chan struct{})}
	if len(urls) == 0 {
		l.Fatalf("invalid client urls: %v", urls)
	}
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Fatalf("failed to create harmony rpc client: %v", err)
		}
		cleth := ethclient.NewClient(clrpc)
		clbmc, err := NewBMC(HexToAddress(bmcAddr), cleth)
		if err != nil {
			l.Fatalf("failed to create bmc binding to ethclient: %v", err)
		}
		cl.rpcs = append(cl.rpcs, clrpc)
		cl.eths = append(cl.eths, cleth)
		cl.bmcs = append(cl.bmcs, clbmc)
	}
	return cl
}

func (cl *Client) rpc() *rpc.Client       { return cl.rpcs[rand.Intn(len(cl.rpcs))] }
func (cl *Client) eth() *ethclient.Client { return cl.eths[rand.Intn(len(cl.eths))] }
func (cl *Client) bmc() *BMC              { return cl.bmcs[rand.Intn(len(cl.bmcs))] }

func (cl *Client) newTransactOpts(w Wallet) *bind.TransactOpts {
	ew := w.(*wallet.EvmWallet)
	context := context.Background()
	chainID, _ := cl.eth().ChainID(context)
	txOpts, _ := bind.NewKeyedTransactorWithChainID(ew.Skey, chainID)
	txOpts.GasPrice, _ = cl.eth().SuggestGasPrice(context)
	txOpts.Context = context
	return txOpts
}

func (cl *Client) GetChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	return cl.eth().ChainID(ctx)
}

func (cl *Client) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	bn, err := cl.eth().BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

func (cl *Client) GetTransaction(hash common.Hash) (*ethtypes.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	tx, pending, err := cl.eth().TransactionByHash(ctx, hash)
	if err != nil {
		return nil, pending, err
	}
	return tx, pending, err
}

func (cl *Client) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	tr := new(types.Receipt)
	err := cl.rpc().CallContext(ctx, tr, "hmy_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (cl *Client) GetEpochLastBlock(epoch *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	lbn := big.NewInt(0)
	err := cl.rpc().CallContext(ctx, lbn, "hmyv2_epochLastBlock", epoch)
	if err != nil {
		return nil, err
	}
	return lbn, nil
}

func (cl *Client) GetHmyBlockByHeight(height *big.Int) (*Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	hb := new(Block)
	err := cl.rpc().CallContext(ctx, hb, "hmyv2_getBlockByNumber", height, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (cl *Client) rpcConsensusCall(
	threshold float64,
	method string,
	valfn func() interface{},
	keyfn func(val interface{}) interface{},
	args ...interface{}) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()

	total := len(cl.rpcs)

	ech := make(chan error, total)
	vch := make(chan interface{}, total)
	for _, clr := range cl.rpcs {
		go func(clr *rpc.Client) {
			val := valfn()
			err := clr.CallContext(ctx, val, method, args...)
			if err != nil {
				val = nil
			}
			ech <- err
			vch <- val
		}(clr)
	}

	counts := make(map[interface{}]int, total)
	lookup := make(map[interface{}]interface{}, total)
	for i := 0; i < total; i++ {
		val := <-vch
		if val != nil {
			key := keyfn(val)
			counts[key]++
			lookup[key] = val
		}
	}

	mk, mc := interface{}(nil), 0
	for k, c := range counts {
		if c > mc {
			mk, mc = k, c
		}
	}

	consensus := float64(mc) / float64(total)
	if consensus < threshold {
		return nil, fmt.Errorf("consensus failure: %.2f/%.2f", consensus, threshold)
	}

	return lookup[mk], nil
}

func (cl *Client) GetHmyHeaderByHeight(height *big.Int, consensusThreshold float64) (*Header, error) {
	h, err := cl.rpcConsensusCall(
		consensusThreshold,
		"hmyv2_getFullHeader",
		func() interface{} { return &Header{} },
		func(val interface{}) interface{} {
			if val, ok := val.(*Header); ok {
				return val.Hash()
			}
			return nil
		},
		height)
	if err != nil {
		return nil, err
	}
	return h.(*Header), nil
}

func (cl *Client) GetBlockReceipts(hash common.Hash) ([]*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultReadTimeout)
	defer cancel()
	receipts := make([]*types.Receipt, 0)
	err := cl.rpc().CallContext(ctx, &receipts, "hmy_getBlockReceipts", hash)
	if err != nil {
		return nil, err
	}
	return receipts, nil

}

func (cl *Client) CloseAllMonitor() error {
	close(cl.stopAllMonitor)
	return nil
}

func (cl *Client) NewValidator(height uint64) (*Validator, error) {
	return newValidator(cl, new(big.Int).SetUint64(height))
}

// MonitorBlock subscribes to the new head notifications
func (cl *Client) MonitorBlock(startHeight uint64, fetchReceipts bool, cb func(v *BlockNotification) error) error {
	if startHeight == 0 {
		startHeight = 1
	}
	vl, err := cl.NewValidator(startHeight)
	if err != nil {
		return errors.Wrapf(err, "monitor block: NewValidator: %v", err)
	}

	// block notification channel (buffered: to avoid deadlock)
	bns := make(chan *BlockNotification, BlockNotificationSyncConcurrency) // increase this for faster sync

	latest, next := int64(0), int64(startHeight)

	poll := time.NewTicker(time.Second)
	defer poll.Stop()

	// last few unverified block notifications
	lbns := make([]*BlockNotification, 0, 1)

	// start monitor loop
	for {
		select {
		case <-cl.stopAllMonitor:
			return nil

		case <-poll.C:
			n, err := cl.GetBlockNumber()
			if err != nil {
				return errors.Wrapf(err,
					"monitor block: poll block number: %v", err)
			}
			if int64(n) <= latest {
				continue
			}
			latest = int64(n)
			if next > latest {
				cl.log.Debugf(
					"monitor block: skipping; latest=%d, next=%d",
					latest, next)
			}

		case bn := <-bns:
			for ; bn != nil; next++ { // empty bns channel: process all notifications
				if len(lbns) > 0 {
					lbn := lbns[len(lbns)-1]
					ok, err := vl.verify(lbn.Header,
						bn.Header.LastCommitSignature, bn.Header.LastCommitBitmap)
					if err != nil || !ok {
						cl.log.Errorf("monitor block: invalid header: n=%v, err=%v", lbn.Header.Number, err)
						break
					}
					if err := cb(lbn); err != nil {
						return errors.Wrapf(err, "monitor block: callback: %v", err)
					}
					if err := vl.update(lbn.Header); err != nil {
						return errors.Wrapf(err, "monitor block: update validator: %v", err)
					}
				}
				lbns = lbns[:1]
				if lbns[0], bn = bn, nil; len(bns) > 0 {
					bn = <-bns
				}
			}

		default:
			if next > latest {
				time.Sleep(time.Millisecond)
				continue
			}
			type bnq struct {
				h     int64
				v     *BlockNotification
				err   error
				retry int
			}

			bch := make(chan *bnq, cap(bns))
			for i := next; i <= latest &&
				len(bch) < cap(bch); i++ {
				bch <- &bnq{i, nil, nil, 3} // fill bch with requests
			}

			_bns := make([]*BlockNotification, 0, len(bch))

			for r := range bch {
				switch {
				case r.err != nil:
					if r.retry == 0 {
						return errors.Wrapf(r.err,
							"monitor block: h=%d, %v", r.h, r.err)
					}
					r.retry--
					r.err = nil
					bch <- r
				case r.v != nil:
					_bns = append(_bns, r.v)
					if len(_bns) == cap(_bns) {
						close(bch)
					}
				default:
					go func(r *bnq, fr bool) {
						var err error
						v := &BlockNotification{Height: big.NewInt(r.h)}
						v.Header, err = cl.GetHmyHeaderByHeight(v.Height, 0)
						if err != nil {
							r.err = errors.Wrapf(err, "GetHmyHeaderByHeight: %v", err)
							bch <- r
							return
						}
						v.Hash = v.Header.Hash()
						if fr && v.Header.GasUsed > 0 {
							v.Receipts, err = cl.GetBlockReceipts(v.Hash)
							if err == nil {
								var tr *trie.Trie
								tr, err = receiptsTrie(v.Receipts)
								if err == nil {
									if !bytes.Equal(tr.Hash().Bytes(), v.Header.ReceiptsRoot.Bytes()) {
										err = fmt.Errorf(
											"invalid receipts: root=%v, trie=%v",
											v.Header.ReceiptsRoot, tr.Hash())
									}
								}
							}
							if err != nil {
								r.err = errors.Wrapf(err, "GetBlockReceipts: %v", err)
								bch <- r
								return
							}
						}
						r.v = v
						bch <- r
					}(r, fetchReceipts)
				}
			}

			// sort and forward notifications
			sort.SliceStable(_bns, func(i, j int) bool {
				return _bns[i].Height.Int64() < _bns[j].Height.Int64()
			})
			for _, v := range _bns {
				bns <- v
			}
		}

	}
}
