package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")

	totalIssuance, err := strconv.ParseUint(os.Args[1], 10, 64)

	if err != nil {
		log.Fatalf("Invalid total issuance %s\n", err)
	}

	decimals, err := strconv.ParseUint(os.Args[2], 10, 32)

	if err != nil {
		log.Fatalf("Invalid decimals %s\n", err)
	}

	unitName := os.Args[3]
	assetName := os.Args[4]
	assetURL := os.Args[5]
	assetMetadataHash := os.Args[6]

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Failed to create algod client: %v", err)
	}

	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	minter, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to create account from private key: %v", err)
	}

	defaultFrozen := false
	manager := ""
	freeze := ""
	clawback := ""
	note := []byte(nil)

	tx, err := transaction.MakeAssetCreateTxn(
		minter.Address.String(),
		uint64(txParams.Fee),
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		note,
		txParams.GenesisID,
		base64.StdEncoding.EncodeToString(txParams.GenesisHash),
		totalIssuance,
		uint32(decimals),
		defaultFrozen,
		manager,
		minter.Address.String(), // reserve
		freeze,
		clawback,
		unitName,
		assetName,
		assetURL,
		assetMetadataHash,
	)
	
	if err != nil {
		log.Fatalf("Could not generate asset creation transaction: %s", err)
	}

	_, stx, err := crypto.SignTransaction(privateKey, tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %s\n", err)
	}

	txId, err := client.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("Could not send transaction: %s", err)
	}

	res, err := future.WaitForConfirmation(client, txId, 4, context.Background())

	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

	asaId := res.AssetIndex

	fmt.Println(asaId)
}