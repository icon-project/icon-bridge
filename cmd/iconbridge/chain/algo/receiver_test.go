package algo

import (
	"context"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

// This function should receive a msg chanel as input, to which it shall forward a new msg as soon
// as it detects valid events in txn from new blocks
func Test_Subscribe(t *testing.T) {
	rcv, err := createTestReceiver(testnetAccess)
	if err != nil {
		t.Logf("NewReceiver error: %v", err)
		t.FailNow()
	}

	msgCh := make(chan *chain.Message)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	c, err := newClient(testnetAccess, log.New())

	if err != nil {
		t.Log("Couldn't create client %w", err)
		t.FailNow()
	}

	curRound, err := c.GetLatestRound(ctx)

	if err != nil {
		t.Log("Couldn't retrieve latest round")
		t.FailNow()
	}

	subOpts := chain.SubscribeOptions{
		Seq:    777,
		Height: curRound,
	}

	errCh, err := rcv.Subscribe(ctx, msgCh, subOpts)

	if err != nil {
		t.Log("Couldn't Subscribe")
		t.FailNow()
	}

	//Expect receive error msg when a block with an ApplicationTxCall does not contain receipts
	select {
	case <-ctx.Done():
		t.Error("Test timed out with no blocks")
		t.FailNow()
	case err := <-errCh:
		t.Log(err)
	}
	cancel()

	//Expect goroutine to close error chanel after its ctx aborts
	select {
	case err := <-errCh:
		t.Log(err)
	}
	//TODO add case for successful message once BMC is working
}
