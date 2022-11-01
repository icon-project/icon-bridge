package substrate_eth

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc/abi/btscore"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

const endpoint = "wss://arctic-rpc.icenetwork.io:9944"

func newTestClient(t *testing.T, bmcAddr string) IClient {
	url := endpoint
	cls, _, err := newClients([]string{url}, "", log.New())
	require.NoError(t, err)
	return cls[0]
}

func TestReceiver_GetReceiptProofs(t *testing.T) {
	cl := newTestClient(t, "")
	for _, height := range []int64{244743, 244826, 244834, 245138, 245141, 245179, 245204, 245207} {
		if height%1000 == 0 {
			fmt.Println("checkpoint ", height)
		}
		h, err := cl.GetHeaderByHeight(big.NewInt(height))
		require.NoError(t, err)
		receipts, isEIP1559, err := cl.GetBlockReceiptsFromHeight(big.NewInt(height))
		require.NoError(t, err)
		fmt.Println(isEIP1559)
		receiptsRoot := ethTypes.DeriveSha(receipts, trie.NewStackTrie(nil))
		if bytes.Equal(receiptsRoot.Bytes(), h.ReceiptHash.Bytes()) {
			err = fmt.Errorf(
				"invalid receipts: remote=%v, local=%v",
				h.ReceiptHash, receiptsRoot)
			require.NoError(t, err)
		} else {
			fmt.Println(height)
		}
	}
}

func TestReceiver_GetReceiptProofsType(t *testing.T) {
	cl := newTestClient(t, "")
	initHeight := 244743
	for i := 0; i < 1; i++ {
		height := int64(initHeight + i)
		if height%1000 == 0 {
			fmt.Println("checkpoint ", height)
		}
		b, err := cl.GetEthClient().BlockByNumber(context.TODO(), big.NewInt(height))
		require.NoError(t, err)
		for _, txn := range b.Transactions() {
			fmt.Println(txn.Hash())
			rcpts, err := cl.GetEthClient().TransactionReceipt(context.TODO(), common.HexToHash(txn.Hash().String()))
			require.NoError(t, err)
			fmt.Println(rcpts)
		}
	}
}

func TestParallelFetch(t *testing.T) {
	start := 245204
	end := 245208
	concurrency := 10
	err := parallelFetch(start, end, concurrency)
	if err != nil {
		t.Fatal(fmt.Errorf("parallelTransfers %v", err))
	}
}

func parallelFetch(start, end, concurrency int) error {
	url := endpoint
	cls, _, err := newClients([]string{url}, "", log.New())
	if err != nil {
		return err
	}
	cl := cls[0]
	zeroInt := 0
	reqCursor := 0
	lenRequests := end - start
	type req struct {
		height  int
		err     error
		res     *int
		txnType int
	}

	for reqCursor < lenRequests {
		rqch := make(chan *req, concurrency)
		for i := reqCursor; len(rqch) < cap(rqch) && i < lenRequests; i++ {
			rqch <- &req{height: start + i, err: nil, res: nil}
			reqCursor++
		}
		sres := make([]*req, 0, len(rqch))
		for q := range rqch {
			switch {
			case q.err != nil || q.res != nil:
				sres = append(sres, q)
				if len(sres) == cap(sres) {
					close(rqch)
				}
			default:
				go func(q *req) {
					defer func() {
						time.Sleep(time.Millisecond * 20)
						//fmt.Println("reqCursor ", q.height, q.err, *q.res, q.txnType)
						rqch <- q
					}()
					q.res = &zeroInt
					h, err := cl.GetHeaderByHeight(big.NewInt(int64(q.height)))
					if err != nil {
						q.err = fmt.Errorf("GetHeaderByHeight %v err %v", q.height, err)
						return
					}
					rcpts, _, err := cl.GetBlockReceiptsFromHeight(big.NewInt(int64(q.height)))
					if err != nil {
						q.err = fmt.Errorf("GetBlockReceiptsFromHeight %v err %v", q.height, err)
						return
					}
					receiptsRoot := ethTypes.DeriveSha(rcpts, trie.NewStackTrie(nil))
					if receiptsRoot != h.ReceiptHash {
						q.res = &q.height
						body, err := cl.GetEthClient().BlockByNumber(context.TODO(), big.NewInt(int64(q.height)))
						if err != nil {
							q.err = fmt.Errorf("BlockByNumber %v err %v", q.height, err)
							return
						}
						for _, txn := range body.Transactions() {
							if txn.Type() != 0 {
								q.txnType = int(txn.Type())
								break
							}
						}
					}
				}(q)
			}
		}
		for _, sr := range sres {
			if sr.err != nil || *sr.res != 0 {
				fmt.Println("Mark ", sr.txnType, sr.height, sr.err)
			}
		}
		if reqCursor%1000 == 0 {
			fmt.Println("cursor ", reqCursor)
		}
	}
	return nil
}

func TestReceiver_GetGasPriceUsed(t *testing.T) {
	cls := newTestClient(t, "")
	startBlockNum := 251992 - 50
	endBlockNum := 251992 + 50
	for i := startBlockNum; i < endBlockNum; i++ {
		height := int64(i)
		blk, err := cls.GetEthClient().BlockByNumber(context.TODO(), big.NewInt(height))
		require.NoError(t, err)
		for txni, txn := range blk.Transactions() {
			fmt.Println(i, txni, txn.Hash(), txn.GasPrice())
		}
		fmt.Println("------------------------------------------")
	}
}

func TestParallelTransactions(t *testing.T) {
	for i := 0; i < 1; i++ {
		gp := big.NewInt(int64(1200)) // retry with less than 8
		hash := nativeCoinTransferRequest(t, context.TODO(), gp)
		fmt.Println(i, " ", hash)
		time.Sleep(time.Second)
	}
}

//0x4B9c58976F89f211A5B1079d79cEa25bc5C1e4ED
func nativeCoinTransferRequest(t *testing.T, ctx context.Context, gasPrice *big.Int) string {
	privKeyString := "76069f7df129ea8c1cf8aa1f117ecd223b51e724f4df9e3bdbc733f7b74279f4"
	btscoreAddr := "0x01d8d7802F41FE2DFa962f5807427C9267E390b6"
	recepientAddress := "btp://0x2.icon/hx2637472e23df38a3b90d644017d0c4c973142e72"

	clrpc, err := rpc.Dial(endpoint)
	require.NoError(t, err)
	ethcl := ethclient.NewClient(clrpc)
	btscore, err := btscore.NewBtscore(common.HexToAddress(btscoreAddr), ethcl)
	require.NoError(t, err)

	privBytes, err := hex.DecodeString(privKeyString)
	require.NoError(t, err)
	senderPrivKey, err := crypto.ToECDSA(privBytes)
	require.NoError(t, err)
	txo, err := bind.NewKeyedTransactorWithChainID(senderPrivKey, big.NewInt(552))
	require.NoError(t, err)
	txo.GasLimit = 1400000
	txo.Value = big.NewInt(5000000000000000000)
	txo.Context = ctx
	txo.GasPrice = gasPrice
	nonce, err := ethcl.NonceAt(ctx, common.HexToAddress("0x4B9c58976F89f211A5B1079d79cEa25bc5C1e4ED"), nil)
	require.NoError(t, err)
	txo.Nonce = big.NewInt(int64(nonce))
	// txn, err := btscore.TransferNativeCoin(txo, recepientAddress)
	// require.NoError(t, err)
	// fmt.Println("Init ", txn.Hash().String())
	for i := 0; i < 5; i++ {
		txn, err := btscore.TransferNativeCoin(txo, recepientAddress)
		if err != nil {
			fmt.Println("retrywith error ", err)
		} else {
			fmt.Println("itr ", txn.Hash().String())
		}
	}

	txo.GasPrice = (&big.Int{}).Add(txo.GasPrice, big.NewInt(1))
	txn, err := btscore.TransferNativeCoin(txo, recepientAddress)
	if err != nil {
		fmt.Println("retrywith error ", err)
	} else {
		fmt.Println("final ", txn.Hash().String())
	}

	return "Done"
}
