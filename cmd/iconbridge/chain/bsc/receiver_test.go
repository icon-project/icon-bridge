package bsc

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

const (
	ICON_BMC          = "btp://0x7.icon/cx8a6606d526b96a16e6764aee5d9abecf926689df"
	BSC_BMC_PERIPHERY = "btp://0x61.bsc/0xB4fC4b3b4e3157448B7D279f06BC8e340d63e2a9"
	BlockHeight       = 21447824
)

func newTestReceiver(t *testing.T, src, dst chain.BTPAddress) chain.Receiver {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	receiver, _ := NewReceiver(src, dst, []string{url}, nil, log.New())
	return receiver
}

func newTestClient(t *testing.T, bmcAddr string) *Client {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	cls, _, err := newClients([]string{url}, bmcAddr, log.New())
	require.NoError(t, err)
	return cls[0]
}

func TestMedianGasPrice(t *testing.T) {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	cls, _, err := newClients([]string{url}, BSC_BMC_PERIPHERY, log.New())
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = cls[0].GetMedianGasPriceForBlock()
	require.NoError(t, err)
}

func TestSubscribeMessage(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set(BSC_BMC_PERIPHERY)
	err = dst.Set(ICON_BMC)
	if err != nil {
		fmt.Println(err)
	}

	recv := newTestReceiver(t, src, dst).(*receiver)

	ctx, cancel := context.Background(), func() {}
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()
	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := recv.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    75,
			Height: uint64(BlockHeight),
		})
	require.NoError(t, err, "failed to subscribe")

	for {
		defer cancel()
		select {
		case err := <-srcErrCh:
			t.Logf("subscription closed: %v", err)
			t.FailNow()
		case msg := <-srcMsgCh:
			if len(msg.Receipts) > 0 && msg.Receipts[0].Height == 21447824 {
				// received event exit
				return
			}
		}
	}
}

func TestReceiver_GetReceiptProofs(t *testing.T) {
	cl := newTestClient(t, BSC_BMC_PERIPHERY)
	header, err := cl.GetHeaderByHeight(big.NewInt(BlockHeight))
	require.NoError(t, err)
	hash := header.Hash()
	receipts, err := cl.GetBlockReceipts(hash)
	require.NoError(t, err)
	receiptsRoot := ethTypes.DeriveSha(receipts, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptsRoot.Bytes(), header.ReceiptHash.Bytes()) {
		err = fmt.Errorf(
			"invalid receipts: remote=%v, local=%v",
			header.ReceiptHash, receiptsRoot)
		require.NoError(t, err)
	}
}
