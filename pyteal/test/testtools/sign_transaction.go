package testtools

import (
	"crypto/ed25519"
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func SignTransaction(t *testing.T, privateKey ed25519.PrivateKey, tx types.Transaction) (txid string, stx []byte) {
	t.Helper()

	var err error
	txid, stx, err = crypto.SignTransaction(privateKey, tx)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %s\n", err)
	}
	return
}
