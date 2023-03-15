package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	deployer := helpers.GetAccountFromPrivateKey()
	tealDir := os.Args[1]
	bmcId := helpers.GetUint64FromArgs(2, "bmc id")
	dbshId := helpers.GetUint64FromArgs(3, "dbsh id")
	
	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Failed to create algod client: %v", err)
	}

	bshContract, bshMcp, err := helpers.InitABIContract(client, deployer, filepath.Join(tealDir, "bsh", "contract.json"), dbshId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bshMcp.ForeignApps = []uint64{bmcId}
	bshMcp.SuggestedParams.Fee = 2000

	_, err = helpers.CallAbiMethod(client, bshContract, bshMcp, "sendServiceMessage", []interface{}{})

	if err != nil {
		log.Fatalf("Failed to call sendServiceMessage method for bsh %+v", err)
	}
}