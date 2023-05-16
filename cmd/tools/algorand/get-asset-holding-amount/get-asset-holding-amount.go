package main

import (
	"context"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	assetHolder := helpers.GetAccountFromPrivateKey()

	asaId := helpers.GetUint64FromArgs(1, "asset id")

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
	}

	assetInfo, err := client.AccountAssetInformation(assetHolder.Address.String(), asaId).Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get Asset information method: %+v", err)
	}

	fmt.Println(assetInfo.AssetHolding.Amount)
}