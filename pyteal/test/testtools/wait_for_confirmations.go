package testtools

import (
	"context"
	"testing"

	"appliedblockchain.com/icon-bridge/config"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/future"
)

func WaitForConfirmationsT(
	t *testing.T, client *algod.Client, txIDs []string,
) (result []*models.PendingTransactionInfoResponse) {
	t.Helper()
	result = make([]*models.PendingTransactionInfoResponse, len(txIDs))

	for i, txID := range txIDs {
		res, err := future.WaitForConfirmation(client, txID, config.ConfirmationWaitRounds, context.Background())

		if err != nil {
			t.Fatalf("While waiting for confirmation: %s\n", err)
		}

		result[i] = &res
	}

	return
}
