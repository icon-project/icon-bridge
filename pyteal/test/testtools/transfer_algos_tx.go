package testtools

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

func TransferAlgosTx(t *testing.T, txParams types.SuggestedParams, from types.Address, to types.Address, amount uint64) (txn types.Transaction) {
	t.Helper()

	var err error
	txn, err = transaction.MakePaymentTxn(
		from.String(),
		to.String(),
		uint64(txParams.Fee),
		amount,
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		[]byte(nil),
		"",
		txParams.GenesisID,
		txParams.GenesisHash)
	if err != nil {
		t.Fatalf("Cannot create algo transfer tx: %s", err)
	}
	return
}
