package algo

import (
	"os"
	"testing"

	"github.com/icon-project/icon-bridge/common/log"
)

var (
	algodAddress = os.Getenv("ALGO_TEST_ADR")
	algodToken   = os.Getenv("ALGO_TEST_TOK")
)

func Test_GetLatestBlock(t *testing.T) {
	urls := []string{algodAddress, algodToken}
	c, err := newClient(urls, log.New())
	if err != nil {
		t.Logf("Error creating algorand client: %v", err)
		t.FailNow()
	}

	latestRound, err := c.GetLatestRound()
	if err != nil {
		t.Logf("Error getting latest round: %v", err)
		t.FailNow()
	}

	latestBlock, err := c.GetBlockbyRound(latestRound)
	if err != nil {
		t.Logf("Error getting latest block: %v", err)
		t.FailNow()
	}

	if uint64(latestBlock.BlockHeader.Round) != latestRound {
		t.Logf("Wrong block fetched. Expected:%v Got:%v", uint64(latestBlock.BlockHeader.Round), latestRound)
		t.Fail()
	}
}
