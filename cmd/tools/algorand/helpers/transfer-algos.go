package helpers

import (
	"context"
	"log"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
)

func TransferAlgos(
	client *algod.Client,
	from crypto.Account,
	to string,
	amount uint64,
) {
	txParams, err := client.SuggestedParams().Do(context.Background())

	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	txn, err := transaction.MakePaymentTxn(
		from.Address.String(),
		to,
		uint64(txParams.Fee),
		amount,
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		[]byte(nil),
		"",
		txParams.GenesisID,
		txParams.GenesisHash,
	)

	if err != nil {
		log.Fatalf("Cannot create algo transfer tx: %s", err)
	}

	_, stx, err := crypto.SignTransaction(from.PrivateKey, txn)

	if err != nil {
		log.Fatalf("Failed to sign transaction: %s\n", err)
	}

	txId, err := client.SendRawTransaction(stx).Do(context.Background())

	if err != nil {
		log.Fatalf("Could not send transaction: %s", err)
	}

	_, err = future.WaitForConfirmation(client, txId, 4, context.Background())

	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

}