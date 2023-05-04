package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")
	minterPrivateKeyStr := helpers.GetEnvVar("MINTER_PRIVATE_KEY")
	receiverPrivateKeyStr := helpers.GetEnvVar("RECEIVER_PRIVATE_KEY")

	tealDir := os.Args[1]
	contractName := os.Args[2]
	bmcId := helpers.GetUint64FromArgs(3, "bmc id")
	bshId := helpers.GetUint64FromArgs(4, "bsh id")
	iconBtpAddress := os.Args[5]
	serviceName := os.Args[6]
	asaId := helpers.GetUint64FromArgs(7, "asset id")
	
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}
	
	minterPrivateKey, err := base64.StdEncoding.DecodeString(minterPrivateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode minter private key seed: %s\n", err)
	}
	
	receiverPrivateKey, err := base64.StdEncoding.DecodeString(receiverPrivateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode minter private key seed: %s\n", err)
	}

	deployer, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Cannot create deployer account: %s", err)
	}

	minter, err := crypto.AccountFromPrivateKey(minterPrivateKey)
	if err != nil {
		log.Fatalf("Cannot create minter account: %s", err)
	}

	receiver, err := crypto.AccountFromPrivateKey(receiverPrivateKey)
	if err != nil {
		log.Fatalf("Cannot create receiver account: %s", err)
	}

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
	}

	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	bshAddress := crypto.GetApplicationAddress(bshId)

	bshContract, bshMcp, err := helpers.InitABIContract(client, deployer, filepath.Join(tealDir, contractName, "contract.json"), bshId)

	if err != nil {
		log.Fatalf("Failed to init %s ABI contract: %+v", contractName, err)
	}

	bshMcp.ForeignApps = []uint64{bmcId}

	helpers.TransferAlgos(client, deployer, bshAddress.String(), 614000)

	bshMcp.ForeignAssets = []uint64{asaId}
	_, err = helpers.CallAbiMethod(client, bshContract, bshMcp, "init", []interface{}{bmcId, iconBtpAddress, asaId})

	if err != nil {
		log.Fatalf("Failed to call init method for bsh %+v", err)
	}

	assetTxn, err := helpers.TransferAssetTx(txParams, minter.Address, bshAddress, asaId, 100000000000)
	if err != nil {
		log.Fatalf("Could not generate asset transfer transaction: %s", err)
	}
	assetTxId := helpers.SendTransaction(client, minter.PrivateKey, assetTxn)
	_, err = future.WaitForConfirmation(client, assetTxId, 4, context.Background())
	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

	optInTx, err := helpers.TransferAssetTx(txParams, receiver.Address, receiver.Address, asaId, 0)
	if err != nil {
		log.Fatalf("Could not generate asset optin transaction: %s", err)
	}
	optInTxId := helpers.SendTransaction(client, receiver.PrivateKey, optInTx)
	_, err = future.WaitForConfirmation(client, optInTxId, 4, context.Background())
	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

	bmcContract, bmcMcp, err := helpers.InitABIContract(client, deployer, filepath.Join(tealDir, "bmc", "contract.json"), bmcId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bmcMcp.ForeignAccounts = []string{bshAddress.String()}
	_, err = helpers.CallAbiMethod(client, bmcContract, bmcMcp, "registerBSHContract", []interface{}{bshAddress, serviceName})

	if err != nil {
		log.Fatalf("Failed to add method call: %+v", err)
	}

	info, err := client.AccountApplicationInformation(bshAddress.String(), bmcId).Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get application information: %+v", err)
	}

	fmt.Printf("%+v\n", info)
}
