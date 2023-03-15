package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

// create test for this main
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

	// Check if the file exists
	if _, err := os.Stat("cache/algo_btp_addr"); os.IsNotExist(err) {
		// If the file doesn't exist, create it and write the bmc address there
		bmcAddr := crypto.GetApplicationAddress(res.ApplicationIndex)
		file, err := os.Create("cache/algo_btp_addr")
		if err != nil {
			log.Fatalf("Failed to create file: %s\n", err)
		}
		defer file.Close()
		file.WriteString("btp://0x14.algo/" + bmcAddr.String())
	}
}
