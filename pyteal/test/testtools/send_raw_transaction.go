package testtools

import (
	"context"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

func SendRawTransaction(t *testing.T, client *algod.Client, stx []byte) (txID string) {
	t.Helper()

	var err error
	txID, err = client.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		t.Fatalf("Could not send transaction: %s", err)
	}
	return
}
