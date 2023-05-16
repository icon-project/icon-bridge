package testtools

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func TransferAlgos(
	t *testing.T,
	client *algod.Client,
	txParams types.SuggestedParams,
	from crypto.Account,
	recipients []types.Address,
	amount uint64,
) (txIDs []string) {
	t.Helper()

	txIDs = make([]string, len(recipients))
	for i, address := range recipients {
		tx := TransferAlgosTx(t, txParams, from.Address, address, amount)
		txIDs[i] = SendTransaction(t, client, from.PrivateKey, tx)
	}
	return
}
