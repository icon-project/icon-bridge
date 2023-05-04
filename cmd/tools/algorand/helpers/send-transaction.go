package helpers

import (
	"context"
	"crypto/ed25519"
	"log"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func SendTransaction(client *algod.Client, privateKey ed25519.PrivateKey, tx types.Transaction) (txId string) {
	_, stx, err := crypto.SignTransaction(privateKey, tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %s\n", err)
	}

	txId, err = client.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("Could not send transaction: %s", err)
	}

	return
}
