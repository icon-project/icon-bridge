package testtools

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
)

func AccountFromPrivateKey(t *testing.T, privateKey []byte) (acc crypto.Account) {
	t.Helper()

	var err error
	acc, err = crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to create account from private key: %+v", err)
	}
	return
}
