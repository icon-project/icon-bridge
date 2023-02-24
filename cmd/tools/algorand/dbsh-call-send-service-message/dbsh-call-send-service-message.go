package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")

	tealDir := os.Args[1]

	bmcId, err := strconv.ParseUint(os.Args[2], 10, 64)

	if err != nil {
		log.Fatalf("Invalid BMC Id %s\n", err)
	}

	dbshId, err := strconv.ParseUint(os.Args[3], 10, 64)

	if err != nil {
		log.Fatalf("Invalid Dummy BSH Id %s\n", err)
	}
	
	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Failed to create algod client: %v", err)
	}

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	deployer, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to create account from private key: %v", err)
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