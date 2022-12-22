package algo

import (
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/common/log"
	"golang.org/x/net/context"
)

func Test_GetLatestBlock(t *testing.T) {
	c, err := newClient(sandboxAccess, log.New())
	if err != nil {
		t.Logf("Error creating algorand client: %v", err)
		t.FailNow()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	latestRound, err := c.GetLatestRound(ctx)
	if err != nil {
		t.Logf("Error getting latest round: %v", err)
		t.FailNow()
	}

	latestBlock, err := c.GetBlockbyRound(ctx, latestRound)
	if err != nil {
		t.Logf("Error getting latest block: %v", err)
		t.FailNow()
	}

	if uint64(latestBlock.BlockHeader.Round) != latestRound {
		t.Logf("Wrong block fetched. Expected:%v Got:%v", uint64(latestBlock.BlockHeader.Round), latestRound)
		t.Fail()
	}
}
