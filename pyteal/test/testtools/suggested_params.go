package testtools

import (
	"context"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/types"
)

func SuggestedParams(t *testing.T, client *algod.Client) (params types.SuggestedParams) {
	t.Helper()

	var err error
	params, err = client.SuggestedParams().Do(context.Background())
	if err != nil {
		t.Fatalf("Error getting suggested tx params: %s\n", err)
	}
	return
}
