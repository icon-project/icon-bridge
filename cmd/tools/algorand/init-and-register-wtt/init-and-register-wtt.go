package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	deployer := helpers.GetAccountFromPrivateKey()
	tealDir := os.Args[1]
	contractName := os.Args[2]
	bmcId := helpers.GetUint64FromArgs(3, "bmc id")
	bshId := helpers.GetUint64FromArgs(4, "bsh id")
	iconBtpAddress := os.Args[5]
	serviceName := os.Args[6]
	asaId := helpers.GetUint64FromArgs(7, "asset id")

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
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
