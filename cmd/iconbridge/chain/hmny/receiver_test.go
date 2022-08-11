//go:build hmny
// +build hmny

package hmny

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

var (
	emptyReceiptsRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	nid               = "0x63564c40"
	net               = "0x63564c40.hmny"
	rpc_uri           = "https://api.harmony.one"
	block_height      = 28070290
	hmny_bmc          = "btp://0x63564c40.hmny/0x293b2D1B12393c70fCFcA0D9cb99889fFD4A23a8"
	icon_bmc          = "btp://0x1.icon/cx06f42ea934731b4867fca00d37c25aa30bc3e3d7"
)

func newTestReceiver(t *testing.T) chain.Receiver {
	url := rpc_uri
	receiver, err := NewReceiver(chain.BTPAddress(hmny_bmc), chain.BTPAddress(icon_bmc), []string{url}, nil, log.New())
	require.NoError(t, err)
	return receiver
}

func TestSubscribeMessage(t *testing.T) {
	cl, _ := newTestClient(t, rpc_uri)

	recv := newTestReceiver(t).(*receiver)

	height := uint64(block_height)
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
	recv.opts.SyncConcurrency = 10

	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := recv.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    0,
			Height: height,
		})
	require.NoError(t, err, "failed to subscribe message")

	for {
		select {
		case err := <-srcErrCh:
			t.Logf("subscription closed: %v", err)
			t.FailNow()
		case msg := <-srcMsgCh:

			// validate receipts height matches block height
			if len(msg.Receipts) > 0 && msg.Receipts[0].Height == 28070299 {
				//found expected block
				return
			}
		}
	}

}
