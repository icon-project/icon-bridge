package algo

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func Test_Subscribe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	c, err := newClient(testnetAccess, log.New())
	if err != nil {
		t.Log("Couldn't create client %w", err)
		t.FailNow()
	}
	curRound, err := c.GetLatestRound(ctx)
	if err != nil {
		t.Log("Couldn't retrieve latest round: %w", err)
		t.FailNow()
	}
	blk, err := c.GetBlockbyRound(ctx, curRound-11)

	if err != nil {
		t.Log("Couldn't retrieve block: %w", err)
		t.FailNow()
	}

	// start receiver 10 rounds late to test that it can update until the current round
	rcv, err := createTestReceiver(testnetAccess, curRound-10, EncodeBlockHash(blk))
	if err != nil {
		t.Logf("NewReceiver error: %v", err)
		t.FailNow()
	}

	msgCh := make(chan *chain.Message)

	subOpts := chain.SubscribeOptions{
		Seq:    777,
		Height: curRound,
	}

	errCh, err := rcv.Subscribe(ctx, msgCh, subOpts)

	if err != nil {
		t.Logf("Couldn't Subscribe. Error: %v", err)
		t.FailNow()
	}

	// create a sender to send a call from the bmc that the receiver will be monotoring
	s, err := createTestSender(testnetAccess)
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}

	_, err = s.(*sender).callAbi(ctx, AbiFunc{"sendMessage",
		[]interface{}{"this", "string", "hhll", []byte{0x01, 0x02, 0x03}}})

	if err != nil {
		t.Logf("Couldn't call sendMessage. Error: %v", err)
		t.FailNow()
	}

	//Expect receive error msg when a block with an ApplicationTxCall does not contain receipts
	select {
	case <-ctx.Done():
		t.Error("Test timed out with no blocks")
		t.FailNow()
	case err := <-errCh:
		t.Logf("Received error: %v", err)
		t.FailNow()
	case msg := <-msgCh:
		t.Logf("Received message: %v", msg)
	}

	//Expect goroutine to close error chanel after its ctx aborts
	cancel()
	select {
	case err := <-errCh:
		t.Log(err)
	}
}

func xTest_GetHash(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	cl, err := newClient(testnetAccess, log.New())
	if err != nil {
		t.Log("Couldn't create client %w", err)
		t.FailNow()
	}
	curRound, err := cl.GetLatestRound(ctx)
	if err != nil {
		t.Log("Couldn't retrieve latest round")
		t.FailNow()
	}

	curBlock, err := cl.GetBlockbyRound(ctx, curRound)
	if err != nil {
		t.Logf("Current block error: %v", err)
		t.FailNow()
	}

	prvBlock, err := cl.GetBlockbyRound(ctx, curRound-1)
	if err != nil {
		t.Logf("Previous block error: %v", err)
		t.FailNow()
	}

	prvHash := EncodeBlockHash(prvBlock)
	curHash := curBlock.Branch
	if !bytes.Equal(prvHash[:], curHash[:]) {
		t.Errorf("Error: expected %v, got %v", prvHash, curHash)
	}
}

func Test_newTestAccount(t *testing.T) {
	genAlgoAccount()
}
