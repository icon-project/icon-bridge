package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func transferAlgos(
	client *algod.Client,
	from crypto.Account,
	to types.Address,
	amount uint64,
) {
	txParams, err := client.SuggestedParams().Do(context.Background())

	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	txn, err := transaction.MakePaymentTxn(
		from.Address.String(),
		to.String(),
		uint64(txParams.Fee),
		amount,
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		[]byte(nil),
		"",
		txParams.GenesisID,
		txParams.GenesisHash,
	)

	if err != nil {
		log.Fatalf("Cannot create algo transfer tx: %s", err)
	}

	_, stx, err := crypto.SignTransaction(from.PrivateKey, txn)

	if err != nil {
		log.Fatalf("Failed to sign transaction: %s\n", err)
	}

	txId, err := client.SendRawTransaction(stx).Do(context.Background())

	if err != nil {
		log.Fatalf("Could not send transaction: %s", err)
	}

	_, err = future.WaitForConfirmation(client, txId, 4, context.Background())

	if err != nil {
		log.Fatalf("While waiting for confirmation: %s\n", err)
	}

}

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")

	tealDir := os.Args[1]
	cacheDir := os.Args[2]

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	deployer, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Cannot create deployer account: %s", err)
	}

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Algod client could not be created: %s\n", err)
	}

	bmcId, err := strconv.ParseUint(helpers.GetFileVar(cacheDir, "bmc_app_id"), 10, 64)

	if err != nil {
		log.Fatalf("Invalid BMC Id %s\n", err)
	}

	dbshId, err := strconv.ParseUint(helpers.GetFileVar(cacheDir, "dbsh_app_id"), 10, 64)

	if err != nil {
		log.Fatalf("Invalid Dummy BSH Id %s\n", err)
	}

	bshAddress := crypto.GetApplicationAddress(dbshId)

	bshContract, bshMcp, err := helpers.InitABIContract(client, deployer, filepath.Join(tealDir, "bsh", "contract.json"), dbshId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bshMcp.ForeignApps = []uint64{bmcId}

	transferAlgos(client, deployer, bshAddress, 514000)

	_, err = helpers.CallAbiMethod(client, bshContract, bshMcp, "init", []interface{}{bmcId, helpers.GetFileVar(cacheDir, "icon_btp_addr")})

	if err != nil {
		log.Fatalf("Failed to call init method for bsh %+v", err)
	}

	bmcContract, bmcMcp, err := helpers.InitABIContract(client, deployer, filepath.Join(tealDir, "bmc", "contract.json"), bmcId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bmcMcp.ForeignAccounts = []string{bshAddress.String()}
	_, err = helpers.CallAbiMethod(client, bmcContract, bmcMcp, "registerBSHContract", []interface{}{bshAddress, "dbsh"})

	if err != nil {
		log.Fatalf("Failed to add method call: %+v", err)
	}

	info, err := client.AccountApplicationInformation(bshAddress.String(), bmcId).Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get application information: %+v", err)
	}

	fmt.Printf("%+v\n", info)
}
