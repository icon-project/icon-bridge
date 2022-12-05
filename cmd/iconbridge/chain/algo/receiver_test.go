package algo

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

var (
	block_height = 28070290
	algo_bmc     = "btp://0x14.algo/0x293b2D1B12393c70fCFcA0D9cb99889fFD4A23a8"
	icon_bmc     = "btp://0x1.icon/cx06f42ea934731b4867fca00d37c25aa30bc3e3d7"
)

// This function should receive a msg chanel as input, to which it shall forward a new msg as soon
// as it detects valid events in txn from new blocks
func Test_Subscribe(t *testing.T) {
	algodAccess := []string{algodAddress, algodToken}
	opts := map[string]interface{}{"syncConcurrency": 2}
	rawOpts, err := json.Marshal(opts)
	if err != nil {
		t.Logf("Marshalling opts: %v", err)
		t.FailNow()
	}

	rcv, err := NewReceiver(chain.BTPAddress(icon_bmc), chain.BTPAddress(algo_bmc),
		algodAccess, rawOpts, log.New())
	if err != nil {
		t.Logf("NewReceiver error: %v", err)
		t.FailNow()
	}

	msgCh := make(chan *chain.Message)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	c, err := newClient(algodAccess, log.New())

	curRound, err := c.GetLatestRound()

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
