//go:build hmny
// +build hmny

package hmny

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/bmr/common/log"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T) *client {
	url := "https://rpc.s0.b.hmny.io"
	cls, err := newClients([]string{url}, "", log.New())
	require.NoError(t, err)
	return cls[0]
}

func getDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Minute)
}

func TestGetTransactionRevertReason(t *testing.T) {
	cl := newTestClient(t)
	txh := common.HexToHash("0x04c3009eb637b8871cfc3732bfe6c23bca1b6e850a6e8bb47dd32ac521d7af7b")

	ctx, cancel := getDefaultContext()
	defer cancel()
	tx, _, err := cl.eth.TransactionByHash(ctx, txh)
	require.NoError(t, err)

	ctx, cancel = getDefaultContext()
	defer cancel()
	txr, err := cl.eth.TransactionReceipt(ctx, txh)
	require.NoError(t, err)

	if txr.Status == 0 {
		callMsg := ethereum.CallMsg{
			From:       common.HexToAddress("0x5f7043477705a4b4a5cb612c76715aec35c26afc"),
			To:         tx.To(),
			Gas:        tx.Gas(),
			GasPrice:   tx.GasPrice(),
			Value:      tx.Value(),
			AccessList: tx.AccessList(),
			Data:       tx.Data(),
		}

		ctx, cancel = getDefaultContext()
		defer cancel()
		data, err := cl.eth.CallContract(ctx, callMsg, txr.BlockNumber)
		require.NoError(t, err)

		t.Logf("revert reason: %v", revertReason(data))
	}
}

func TestRevertReason(t *testing.T) {
	str := "08c379a0" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"000000000000000000000000000000000000000000000000000000000000002b" +
		"526576657274496e76616c696452785365713a2065762e736571203e206578706563746564207278536571000000000000000000000000000000000000000000"
	reason := revertReason(common.Hex2Bytes(str))
	require.Equal(t, "RevertInvalidRxSeq: ev.seq > expected rxSeq", reason, "revert reason should match")
}

func TestBlockAndHeaderHashMatch(t *testing.T) {
	n := int64(1000000) // block number
	cl := newTestClient(t)
	b, err := cl.GetHmyV2BlockByHeight(big.NewInt(n))
	require.NoError(t, err, "failed to get block by height")

	h, err := cl.GetHmyV2HeaderByHeight(big.NewInt(n))
	require.NoError(t, err, "failed to get header by height")

	require.Equal(t, h.Hash(), b.Hash)
}

func TestNewVerifier(t *testing.T) {
	n := uint64(1000000)
	cl := newTestClient(t)

	next, err := cl.GetHmyV2HeaderByHeight((&big.Int{}).SetUint64(n + 1))
	require.NoError(t, err, "failed to fetch next header")

	_, err = cl.newVerifier(&VerifierOptions{
		BlockHeight:     n,
		CommitBitmap:    next.LastCommitBitmap,
		CommitSignature: next.LastCommitSignature,
	})
	require.NoError(t, err, "failed to initialize verifier")
}

func TestBMCMessageDecode(t *testing.T) {
	cl := newTestClient(t)

	receiptWithBMCMessage := `{
		"blockHash": "0xb5261bf0156a310b2de99d4e30bd69cf5e28bd7c501c9313abff6a62d4fd955c",
		"blockNumber": "0x179ef4d",
		"contractAddress": "0x0000000000000000000000000000000000000000",
		"cumulativeGasUsed": "0x2ec0b6",
		"from": "one1tacyx3mhqkjtffwtvyk8vu26as6uy6hu8khl77",
		"gasUsed": "0x2ec0b6",
		"logs": [
			{
				"address": "0x233909be3797bbd135a837ac945bcde3cb078969",
				"blockHash": "0xb5261bf0156a310b2de99d4e30bd69cf5e28bd7c501c9313abff6a62d4fd955c",
				"blockNumber": "0x179ef4d",
				"data": "0x00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000d9500000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000396274703a2f2f3078372e69636f6e2f6378396535633061373439656539346330316665626530343730323138343030326137366138346638340000000000000000000000000000000000000000000000000000000000000000000000000000aaf8a8b8406274703a2f2f307836333537643265302e686d6e792f307832333339303962453337393742426431333541383337414339343562436445336342303738393639b8396274703a2f2f3078372e69636f6e2f6378396535633061373439656539346330316665626530343730323138343030326137366138346638349a576f6e6465726c616e64546f6b656e53616c6553657276696365028ecd028bca088802bda9341f23c00000000000000000000000000000000000000000000000",
				"logIndex": "0x11",
				"removed": false,
				"topics": [
					"0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b"
				],
				"transactionHash": "0x3d7ee415b25de559e0eb062e9d57f29e2ca67f01de893836806f001bea413dcb",
				"transactionIndex": "0x0"
			}
		],
		"logsBloom": "0x00000006000000000040000000001020000000084000020000000000000000000000000000000000000000080000000000000000000001000201000000840000000000000000000000004069000008000000000000000000000000000010002000100000020001000400000000000800000000000000000000000010000000040000000000000000004000000000100000000000000000000000000000000000000000000008000044000000000000000000000000000000100000000008010000210802020400001010002000000000000100000000000000000000800020000000000000000020000002000004000200000000000020000000008004200200",
		"root": "0x",
		"shardID": 0,
		"status": "0x1",
		"to": "one1yvusn03hj7aazddgx7kfgk7du09s0ztf9akn06",
		"transactionHash": "0x3d7ee415b25de559e0eb062e9d57f29e2ca67f01de893836806f001bea413dcb",
		"transactionIndex": "0x0"
	}`

	r := types.Receipt{}

	err := json.Unmarshal([]byte(receiptWithBMCMessage), &r)
	require.NoError(t, err, "failed to unmarshal receipt from json")

	for _, log := range r.Logs {
		ethlog := ethtypes.Log{
			Data:   log.Data,
			Topics: log.Topics,
		}
		msg, err := cl.bmc.ParseMessage(ethlog)
		require.NoError(t, err, "failed to parse btp log")

		json.NewEncoder(os.Stdout).Encode(msg)
		require.Equal(t, uint64(0xd95), msg.Seq.Uint64())
		require.Equal(t, "btp://0x7.icon/cx9e5c0a749ee94c01febe04702184002a76a84f84", msg.Next)

		fmt.Println(common.Bytes2Hex(msg.Msg))
		var bmcMsg TypesBMCMessage
		err = rlp.DecodeBytes(msg.Msg, &bmcMsg)
		require.NoError(t, err, "failed to decode rlp into underlying bmc message")

		require.Equal(t, "btp://0x6357d2e0.hmny/0x233909bE3797BBd135A837AC945bCdE3cB078969", bmcMsg.Src)
		require.Equal(t, "btp://0x7.icon/cx9e5c0a749ee94c01febe04702184002a76a84f84", bmcMsg.Dst)
		require.Equal(t, "WonderlandTokenSaleService", bmcMsg.Svc)
		require.Equal(t, uint64(2), bmcMsg.Sn.Uint64())

		var svcMsg TypesBMCService
		err = rlp.DecodeBytes(bmcMsg.Message, &svcMsg)
		require.NoError(t, err, "failed to decode rlp into underlying bmc service")

		require.Equal(t, "\x02", svcMsg.ServiceType)

		json.NewEncoder(os.Stdout).Encode(bmcMsg)
	}
}

func TestGetBlockReceiptsByBlockHash(t *testing.T) {
	cl := newTestClient(t)

	// TODO generate transactions and note their block numbers
	s, e := 5, 21

	// validate the receipt roots
	for i := int64(s); i <= int64(e); i++ {
		b, err := cl.GetHmyBlockByHeight(big.NewInt(i))
		require.NoError(t, err, "failed to get block by height")

		receipts, err := cl.GetBlockReceipts(b.Hash)
		require.NoError(t, err, "failed to get block receipts")

		require.Equal(t, b.ReceiptsRoot, types.DeriveSha(receipts))
	}

}

func TestGetBlockReceiptsByHeaderHash(t *testing.T) {
	cl := newTestClient(t)

	s, err := cl.GetBlockNumber()
	require.NoError(t, err, "failed to get block number")

	// validate the receipt roots
	for i := int64(s - 10); i <= int64(s); i++ {
		h, err := cl.GetHmyV2HeaderByHeight(big.NewInt(i))
		if err != nil {
			i--
			t.Logf("failed to get block header: h=%d, err=%v", i, err)
			continue
		}
		hash := h.Hash()
		receipts, err := cl.GetBlockReceipts(hash)
		if err != nil {
			i--
			t.Logf("failed to get block receipts: h=%d, v=%v, err=%v", i, hash, err)
			continue
		}
		require.Equal(t, h.ReceiptsRoot, types.DeriveSha(receipts))
	}
}
