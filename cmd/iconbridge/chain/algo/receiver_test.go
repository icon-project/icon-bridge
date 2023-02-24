package algo

import (
	"testing"
)

// func Test_Subscribe(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	c, err := newClient(sandboxAccess, log.New())
// 	if err != nil {
// 		t.Log("Couldn't create client %w", err)
// 		t.FailNow()
// 	}
// 	curRound, err := c.GetLatestRound(ctx)
// 	if err != nil {
// 		t.Log("Couldn't retrieve latest round: %w", err)
// 		t.FailNow()
// 	}
// 	hash, err := c.GetBlockHash(ctx, curRound-11)

// 	if err != nil {
// 		t.Log("Couldn't retrieve hash: %w", err)
// 		t.FailNow()
// 	}

// 	// start receiver 10 rounds late to test that it can update until the current round
// 	rcv, err := createTestReceiver(sandboxAccess, curRound-10, hash)
// 	if err != nil {
// 		t.Logf("NewReceiver error: %v", err)
// 		t.FailNow()
// 	}

// 	msgCh := make(chan *chain.Message)

// 	subOpts := chain.SubscribeOptions{
// 		Seq:    777,
// 		Height: curRound,
// 	}

// 	errCh, err := rcv.Subscribe(ctx, msgCh, subOpts)

// 	if err != nil {
// 		t.Logf("Couldn't Subscribe. Error: %v", err)
// 		t.FailNow()
// 	}

// 	// create a sender to send a call from the bmc that the receiver will be monotoring
// 	s, err := createTestSender(sandboxAccess)
// 	if err != nil {
// 		t.Logf("Failed creting new sender:%v", err)
// 		t.FailNow()
// 	}

// 	_, err = s.(*sender).callAbi(ctx, AbiFunc{"sendMessage",
// 		[]interface{}{"btp://0x14.algo/0x293b2D1B12393c70fCFcA0D9cb99889fFD4A23a8",
// 			"btp://0x2.icon/cx04d4cc5ee639aa2fc5f2ededa7b50df6044dd325",
// 			"tokentransfer", 778, []byte{0x01, 0x02, 0x03}}})

// 	if err != nil {
// 		t.Logf("Couldn't call sendMessage. Error: %v", err)
// 		t.FailNow()
// 	}

// 	//Expect receive error msg when a block with an ApplicationTxCall does not contain receipts
// 	select {
// 	case <-ctx.Done():
// 		t.Error("Test timed out with no blocks")
// 		t.FailNow()
// 	case err := <-errCh:
// 		t.Logf("Received error: %v", err)
// 		t.FailNow()
// 	case msg := <-msgCh:
// 		t.Logf("Received message: %v", msg)
// 	}

// 	//Expect goroutine to close error chanel after its ctx aborts
// 	cancel()
// 	select {
// 	case err := <-errCh:
// 		t.Log(err)
// 	}
// }

// func Test_GetHash(t *testing.T) {
// 	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
// 	cl, err := newClient(testnetAccess, log.New())
// 	if err != nil {
// 		t.Log("Couldn't create client %w", err)
// 		t.FailNow()
// 	}
// 	curRound, err := cl.GetLatestRound(ctx)
// 	if err != nil {
// 		t.Log("Couldn't retrieve latest round")
// 		t.FailNow()
// 	}

// 	curBlock, err := cl.GetBlockbyRound(ctx, curRound)
// 	if err != nil {
// 		t.Logf("Current block error: %v", err)
// 		t.FailNow()
// 	}

// 	prvBlock, err := cl.GetBlockbyRound(ctx, curRound-1)
// 	if err != nil {
// 		t.Logf("Previous block error: %v", err)
// 		t.FailNow()
// 	}

// 	prvHash := EncodeBlockHash(prvBlock)
// 	curHash := curBlock.Branch
// 	if !bytes.Equal(prvHash[:], curHash[:]) {
// 		t.Errorf("Error: expected %v, got %v", prvHash, curHash)
// 	}
// }

func Test_newTestAccount(t *testing.T) {
	genAlgoAccount()
}
