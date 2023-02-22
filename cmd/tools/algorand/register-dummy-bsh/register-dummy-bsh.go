package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func initABIContract(client *algod.Client, deployer crypto.Account, contractDir string, appId uint64) (contract *abi.Contract, mcp future.AddMethodCallParams, err error) {
	b, err := ioutil.ReadFile(contractDir)
	if err != nil {
		fmt.Printf("Failed to open contract file: %+v", err)
		return
	}

	contract = &abi.Contract{}
	if err = json.Unmarshal(b, contract); err != nil {
		fmt.Printf("Failed to marshal contract: %+v", err)
		return
	}

	sp, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get suggeted params: %+v", err)
		return
	}

	sp.Fee = 1000

	signer := future.BasicAccountTransactionSigner{Account: deployer}

	mcp = future.AddMethodCallParams{
		AppID:           appId,
		Sender:          deployer.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	return
}

func callAbiMethod(client *algod.Client, contract *abi.Contract, mcp future.AddMethodCallParams, name string, args []interface{}) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	method, err := contract.GetMethodByName(name)

	if err != nil {
		log.Fatalf("No method named: %s", name)
	}

	mcp.Method = method
	mcp.MethodArgs = args

	err = atc.AddMethodCall(mcp)

	if err != nil {
		fmt.Printf("Failed to add method %s call: %+v \n", name, err)
		return
	}

	ret, err = atc.Execute(client, context.Background(), 2)

	if err != nil {
		fmt.Printf("Failed to execute call: %+v \n", err)
		return
	}

	return
}

func getFileVar(filename string) string {
	// open file
	cacheDir := os.Args[2]

	file, err := os.Open(cacheDir + filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	// read file contents as byte slice
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// convert byte slice to string
	return string(byteValue)
}

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

	bmcId, err := strconv.ParseUint(getFileVar("bmc_app_id"), 10, 64)

	if err != nil {
		log.Fatalf("Invalid BMC Id %s\n", err)
	}

	dbshId, err := strconv.ParseUint(getFileVar("dbsh_app_id"), 10, 64)

	if err != nil {
		log.Fatalf("Invalid Dummy BSH Id %s\n", err)
	}

	bshAddress := crypto.GetApplicationAddress(dbshId)

	bshContract, bshMcp, err := initABIContract(client, deployer, filepath.Join(tealDir, "bsh", "contract.json"), dbshId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bshMcp.ForeignApps = []uint64{bmcId}

	transferAlgos(client, deployer, bshAddress, 514000)

	_, err = callAbiMethod(client, bshContract, bshMcp, "init", []interface{}{bmcId, getFileVar("icon_btp_addr")})

	if err != nil {
		log.Fatalf("Failed to call init method for bsh %+v", err)
	}

	bmcContract, bmcMcp, err := initABIContract(client, deployer, filepath.Join(tealDir, "bmc", "contract.json"), bmcId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bmcMcp.ForeignAccounts = []string{bshAddress.String()}
	_, err = callAbiMethod(client, bmcContract, bmcMcp, "registerBSHContract", []interface{}{bshAddress, "dbsh"})

	if err != nil {
		log.Fatalf("Failed to add method call: %+v", err)
	}

	info, err := client.AccountApplicationInformation(bshAddress.String(), bmcId).Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get application information: %+v", err)
	}

	fmt.Printf("%+v\n", info)
}
