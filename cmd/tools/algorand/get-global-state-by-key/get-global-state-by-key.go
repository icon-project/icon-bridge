package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")

	appId, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("Invalid App id %s\n", err)
	}

	key := os.Args[2]

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Failed to create algod client: %v", err)
	}

	app, err := client.GetApplicationByID(appId).Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get application with id: %v. details: %+v", appId, err)
	}

	encodedKey := base64.StdEncoding.EncodeToString([]byte(key))

	globalState := app.Params.GlobalState

	for _, modelTealKeyValue := range globalState {
		if modelTealKeyValue.Key == encodedKey {
			decodedValue, err := base64.StdEncoding.DecodeString(modelTealKeyValue.Value.Bytes)

			if err != nil {
				log.Fatalf("Failed to decode base64 string: %+v", err)
			}

			fmt.Println(string(decodedValue))
			return
		}
	}

	log.Fatalf("Failed to find %s in global state of application by id: %v", key, appId)
}
