package testtools

import (
	"crypto/ed25519"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/types"
)

func SendTransaction(t *testing.T, algodClient *algod.Client, privateKey ed25519.PrivateKey, tx types.Transaction) (txId string) {
	t.Helper()

	_, stx := SignTransaction(t, privateKey, tx)
	txId = SendRawTransaction(t, algodClient, stx)
	return
}
