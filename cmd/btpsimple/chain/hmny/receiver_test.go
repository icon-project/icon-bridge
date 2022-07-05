//go:build hmny
// +build hmny

// TODO add more receiver tests
package hmny

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

var (
	emptyReceiptsRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

func newTestReceiver(t *testing.T) chain.Receiver {
	url := "https://rpc.s0.b.hmny.io"
	receiver, _ := NewReceiver("", "", []string{url}, nil, log.New())
	return receiver
}

func TestSubscribeMessage(t *testing.T) {
	cl := newTestClient(t)
	recv := newTestReceiver(t).(*receiver)

	height := uint64(1000000)
	next, err := cl.GetHmyV2HeaderByHeight((&big.Int{}).SetUint64(height + 1))
	require.NoError(t, err, "failed to fetch next header")

	ctx, cancel := context.Background(), func() {}
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()

	recv.opts.Verifier = &VerifierOptions{
		BlockHeight:     height,
		CommitBitmap:    next.LastCommitBitmap,
		CommitSignature: next.LastCommitSignature,
	}

	height += 100

	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := recv.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    0,
			Height: height,
		})
	require.NoError(t, err, "failed to subscribe message")

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
