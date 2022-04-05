package hmny

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/btp/common/log"
	"github.com/stretchr/testify/require"
)

func newTestClient() *Client {
	url := "http://localnets:9500"
	// url := "http://44.192.123.4:9500"
	// url := "https://rpc.s0.b.hmny.io"
	return NewClient([]string{url, url, url, url, url, url}, "", log.New())
}

func TestGetBlockReceipts(t *testing.T) {
	cl := newTestClient()

	// TODO generate transactions and note their block numbers
	s, e := 5, 21

	// validate the receipt roots
	for i := int64(s); i <= int64(e); i++ {
		b, err := cl.GetHmyBlockByHeight(big.NewInt(i))
		require.NoError(t, err, "failed to get block by height")

		receipts, err := cl.GetBlockReceipts(b.Hash)
		require.NoError(t, err, "failed to get block receipts")

		tr, err := receiptsTrie(receipts)
		require.NoError(t, err, "failed to create trie from receipts")

		require.Equal(t, b.ReceiptsRoot, tr.Hash())
	}

}

func TestBlockAndHeaderHashMatch(t *testing.T) {
	n := int64(1) // block number
	cl := newTestClient()
	b, err := cl.GetHmyBlockByHeight(big.NewInt(n))
	require.NoError(t, err, "failed to get block by height")

	h, err := cl.GetHmyHeaderByHeight(big.NewInt(n), 0)
	require.NoError(t, err, "failed to get header by height")

	require.Equal(t, h.Hash(), b.Hash)
}

func TestValidator(t *testing.T) {
	n := int64(50)
	cl := newTestClient()
	vl, err := cl.NewValidator(uint64(n))
	require.NoError(t, err, "failed to get validator")

	h, err := cl.GetHmyHeaderByHeight(big.NewInt(n), 0)
	require.NoError(t, err, "failed to fetch header")

	next, err := cl.GetHmyHeaderByHeight(big.NewInt(n+1), 0)
	require.NoError(t, err, "failed to fetch next header")

	ok, err := vl.verify(h, next)
	require.NoError(t, err, "failed to verify commit signature")

	require.True(t, ok)
}

func TestMonitorBlock(t *testing.T) {
	cl := newTestClient()
	err := cl.MonitorBlock(10, true, func(v *BlockNotification) error {
		return json.NewEncoder(os.Stdout).Encode(v.Header)
	})
	require.NoError(t, err)
}

func TestBMCMessageDecode(t *testing.T) {
	cl := newTestClient()

	receiptWithBMCMessage := `{
		"root": "0x",
		"status": "0x1",
		"cumulativeGasUsed": "0x6980c",
		"logsBloom": "0x00000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000004000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000",
		"logs": [
			{
				"address": "0x33b02a85cc1a88071168eb7f527f940baf6f680f",
				"topics": [
					"0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b"
				],
				"data": "0x0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000003e6274703a2f2f30783562396137372e69636f6e2f6378653037303961666331326236303661373161653064323935336364303236336162636237396330300000000000000000000000000000000000000000000000000000000000000000008cf88ab8396274703a2f2f3078322e686d6e792f307833334230326138356363314138383037313136384542374635323746393430424146366636383066b83e6274703a2f2f30783562396137372e69636f6e2f63786530373039616663313262363036613731616530643239353363643032363361626362373963303083626d630089c884496e697482c1c00000000000000000000000000000000000000000",
				"blockNumber": "0x370a",
				"transactionHash": "0x342b5a86d60a30aedcc6dd938e58ab229cafcc6ca5bc36a6acd39cb0598bb2a9",
				"transactionIndex": "0x0",
				"blockHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"logIndex": "0x0",
				"removed": false
			}
		],
		"transactionHash": "0x342b5a86d60a30aedcc6dd938e58ab229cafcc6ca5bc36a6acd39cb0598bb2a9",
		"contractAddress": "0x0000000000000000000000000000000000000000",
		"gasUsed": "0x6980c"
	}`

	r := types.Receipt{}

	err := json.Unmarshal([]byte(receiptWithBMCMessage), &r)
	require.NoError(t, err, "failed to unmarshal receipt from json")

	for _, log := range r.Logs {
		ethlog := ethtypes.Log{
			Data:   log.Data,
			Topics: log.Topics,
		}
		msg, err := cl.bmc().ParseMessage(ethlog)
		require.NoError(t, err, "failed to parse btp log")

		json.NewEncoder(os.Stdout).Encode(msg)
		require.Equal(t, int64(1), msg.Seq.Int64())
		require.Equal(t, "btp://0x5b9a77.icon/cxe0709afc12b606a71ae0d2953cd0263abcb79c00", msg.Next)

		var bmcMsg TypesBMCMessage
		err = rlp.DecodeBytes(msg.Msg, &bmcMsg)
		require.NoError(t, err, "failed to decode rlp into underlying bmc message")

		require.Equal(t, "btp://0x2.hmny/0x33B02a85cc1A88071168EB7F527F940BAF6f680f", bmcMsg.Src)
		require.Equal(t, "btp://0x5b9a77.icon/cxe0709afc12b606a71ae0d2953cd0263abcb79c00", bmcMsg.Dst)
		require.Equal(t, "bmc", bmcMsg.Svc)
		require.Equal(t, int64(0), bmcMsg.Sn.Int64())

		var svcMsg TypesBMCService
		err = rlp.DecodeBytes(bmcMsg.Message, &svcMsg)
		require.NoError(t, err, "failed to decode rlp into underlying bmc service")

		require.Equal(t, "Init", svcMsg.ServiceType)

		json.NewEncoder(os.Stdout).Encode(bmcMsg)
	}
}
