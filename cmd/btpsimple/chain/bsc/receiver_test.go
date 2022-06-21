package bsc

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

func newTestReceiver(t *testing.T) chain.Receiver {
	url := "http://localhost:8545"
	receiver, _ := NewReceiver("", "", []string{url}, nil, log.New())
	return receiver
}

func newTestClient(t *testing.T) *client {
	url := "http://localhost:8545"
	cls, err := NewClient([]string{url}, "", log.New())
	require.NoError(t, err)
	return cls
}

func TestSubscribeMessage(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set("btp://0x97.icon/0xAaFc8EeaEE8d9C8bD3262CCE3D73E56DeE3FB776")
	err = dst.Set("btp://0xf8aac3.icon/cxea19a7d6e9a926767d1d05eea467299fe461c0eb")
	if err != nil {
		fmt.Println(err)
	}

	recv := newTestReceiver(t).(*receiver)
	recv.src = src
	recv.dst = dst
	height := uint64(614)

	ctx, cancel := context.Background(), func() {}
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()

	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := recv.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    12,
			Height: height,
		})
	require.NoError(t, err, "failed to subscribe")

	startHeight := height
	for {
		select {
		case err := <-srcErrCh:
			t.Logf("subscription closed: %v", err)
			t.FailNow()
		case msg := <-srcMsgCh:
			t.Logf("received block: %d", height)

			// validate receipts height matches block height
			if len(msg.Receipts) > 0 {
				require.Equal(t,
					msg.Receipts[0].Height, height,
					"receipts height should match block height")
			}

			// terminate the test after 10 blocks
			height++
			if height > startHeight+10 {
				break
			}
		}
	}
}

func TestReceiver_GetReceiptProofs(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set("btp://0x97.icon/0xAaFc8EeaEE8d9C8bD3262CCE3D73E56DeE3FB776")
	err = dst.Set("btp://0xf8aac3.icon/cxea19a7d6e9a926767d1d05eea467299fe461c0eb")
	if err != nil {
		fmt.Println(err)
	}

	r, _ := NewReceiver(src, dst, []string{"http://localhost:8545"}, nil, log.New())

	/* blockNotification := &BlockNotification{Height: big.NewInt(191)}
	receiptProofs, err := r.(*receiver).newReceiptProofs(blockNotification)

	//fmt.Println(receiptProofs[0].Proof)

	var bytes [][]byte
	_, err = codec.RLP.UnmarshalFromBytes(receiptProofs[0].Proof, &bytes)

	if err != nil {
		return
	} */

	block, err := r.(*receiver).cl.GetBlockByHeight(big.NewInt(191))
	fmt.Println(block.ReceiptHash())
	//fmt.Println(block.Hash())
	//fmt.Println(receiptProofs[0])
	//fmt.Println(len(bytes))
	//for _, proof := range bytes {
	//	fmt.Println(hexutil.Encode(proof))
	//}
}
