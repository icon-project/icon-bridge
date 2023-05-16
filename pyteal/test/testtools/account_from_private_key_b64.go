package testtools

import (
	"encoding/base64"
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
)

func AccountFromPrivateKeyB64(t *testing.T, privateKeyB64 string) (acc crypto.Account) {
	t.Helper()

	var err error
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		t.Fatalf("Could not decode private key base64 string: %s", err)
	}
	acc = AccountFromPrivateKey(t, privateKey)
	return
}
