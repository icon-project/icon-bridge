package testtools

import (
	"context"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

func GetLatestRound(t *testing.T, client *algod.Client) uint64 {
	sta, err := client.Status().Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to getting last round: %+v", err)
	}

	return sta.LastRound
}
