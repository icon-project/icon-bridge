package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	
	txId := os.Args[1]

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
	}

	res, err := future.WaitForConfirmation(client, txId, 4, context.Background())

	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

	fmt.Println(res.ApplicationIndex)
}