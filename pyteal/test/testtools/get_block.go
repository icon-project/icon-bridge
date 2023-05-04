package testtools

import (
	"context"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/types"
)

func GetBlock(t *testing.T, client *algod.Client, round uint64) types.Block {
	t.Helper()

	block, err := client.Block(round).Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to get block by round: %+v\n", err)
	}

	return block
}
