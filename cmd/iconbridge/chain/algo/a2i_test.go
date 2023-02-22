package algo

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
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

const (
	cacheDir = "../../../../devnet/docker/icon-algorand/cache/"
)

// This test is not final, will need changes, but for now sends a transaction to the local network
func Test_SendDummyMessage(t *testing.T) {
	algod, err := algod.MakeClient(algodAddress, algoToken)
	if err != nil {
		t.Errorf("Failed to create algod client: %v", err)
	}

	abiPath, err := filepath.Abs(contractDir + "contract.json")
	if err != nil {
		t.Errorf("Couldn't retrieve abi file: %v", err)
	}
	rawBmc, err := ioutil.ReadFile(abiPath)
	if err != nil {
		t.Errorf("Failed to open contract file: %v", err)
	}
	abiDbsh := &abi.Contract{}
	if err = json.Unmarshal(rawBmc, abiDbsh); err != nil {
		t.Errorf("Failed to marshal abi contract: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	sp, err := algod.SuggestedParams().Do(ctx)
	if err != nil {
		t.Errorf("Failed to get suggeted params: %v", err)
	}

	dbshId, _ := strconv.ParseUint(getFileVar("dbsh_app_id"), 10, 64)
	bmcId, _ := strconv.ParseUint(getFileVar("bmc_app_id"), 10, 64)

	privateKeyStr := getFileVar("algo_private_key")

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	deployer, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		t.Errorf("Failed to create account from private key: %v", err)
	}

	signer := future.BasicAccountTransactionSigner{Account: deployer}
	sp.Fee = 2000

	mcp := future.AddMethodCallParams{
		AppID:           dbshId,
		Sender:          deployer.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	mcp.ForeignApps = []uint64{bmcId}

	sendMsgFunc, err := getMethod(abiDbsh, "sendServiceMessage")
	if err != nil {
		t.Errorf("Failed to get sendMessage method: %v", err)
	}

	var atc = future.AtomicTransactionComposer{}

	err = atc.AddMethodCall(combine(mcp, sendMsgFunc, []interface{}{}))
	if err != nil {
		t.Errorf("Failed to add sendMessage method call: %v", err)
	}

	ret, err := atc.Execute(algod, ctx, 5)
	if err != nil {
		log.Fatalf("Failed to execute call: %+v", err)
	}

	for _, r := range ret.MethodResults {
		log.Printf("%s returned %+v", r.Method.Name, r.ReturnValue)
	}
}

// create func to read file and return a string with the contents
func getFileVar(filename string) string {
	// open file
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
