package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	sender := helpers.GetAccountFromPrivateKey()

	asaId := helpers.GetUint64FromArgs(1, "asset id")
	receiverAddress := os.Args[2]
	amount := helpers.GetUint64FromArgs(3, "amount")

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
	}

	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	assetTxn, err :=  transaction.MakeAssetTransferTxn(
		sender.Address.String(),
		receiverAddress,
		"",
		amount,
		uint64(txParams.Fee),
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		[]byte(nil),
		txParams.GenesisID,
		base64.StdEncoding.EncodeToString(txParams.GenesisHash),
		asaId,
	)

	if err != nil {
		log.Fatalf("Could not generate asset transfer transaction: %s", err)
	}

	assetTxId := helpers.SendTransaction(client, sender.PrivateKey, assetTxn)
	res, err := future.WaitForConfirmation(client, assetTxId, 4, context.Background())
	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

	fmt.Printf("%+v\n", res)
}
