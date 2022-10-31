package substrate_eth

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

const endpoint = ""

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
