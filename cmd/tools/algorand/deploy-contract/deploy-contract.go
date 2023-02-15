package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func compileTeal(client *algod.Client, teal []byte) (program []byte, err error) {
	var response models.CompileResponse
	response, err = client.TealCompile(teal).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("compilation failed: %+v\n ", err)
		return
	}

	program, err = base64.StdEncoding.DecodeString(response.Result)
	if err != nil {
		err = fmt.Errorf("failed to base64 decode compiled program: %s", err)
		return
	}
	return
}

func main () {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")
	
	tealDir := os.Args[1]

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	deployer, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Cannot create deployer account: %s", err)
	}

	approvalSourceCode, err := os.ReadFile(filepath.Join(tealDir, "approval.teal"))
	if err != nil {
		log.Fatalf("Failed to read file: %s\n", err)
	}

	clearStateSourceCode, err := os.ReadFile(filepath.Join(tealDir, "clear.teal"))
	if err != nil {
		log.Fatalf("Failed to read file: %s\n", err)
	}

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
	}

	approvalProgram, err := compileTeal(client, approvalSourceCode)
	if err != nil {
		log.Fatalf("Compilation failed: %+v\n", err)
	}

	clearStateProgram, err := compileTeal(client, clearStateSourceCode)
	if err != nil {
		log.Fatalf("Compilation failed: %+v\n", err)
	}

	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	tx, err := future.MakeApplicationCreateTx(
		false,
		approvalProgram,
		clearStateProgram,
		types.StateSchema{NumUint: 4, NumByteSlice: 4},
		types.StateSchema{NumUint: 4, NumByteSlice: 4},
		[][]byte{},
		nil,
		nil,
		nil,
		txParams,
		deployer.Address,
		nil,
		types.Digest{},
		[32]byte{},
		types.Address{},
	)

	if err != nil {
		log.Fatalf("Failed to make application creation transaction: %s\n", err)
	}

	_, stx, err := crypto.SignTransaction(privateKey, tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %s\n", err)
	}

	txid, err := client.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("Could not send transaction: %s", err)
	}

	fmt.Println(txid)
}