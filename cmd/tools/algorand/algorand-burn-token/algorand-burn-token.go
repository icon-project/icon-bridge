package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	algodAddress := helpers.GetEnvVar("ALGOD_ADDRESS")
	algodToken := helpers.GetEnvVar("ALGOD_TOKEN")
	sender := helpers.GetAccountFromPrivateKey()
	tealDir := os.Args[1]
	bmcId := helpers.GetUint64FromArgs(2, "bmc id")
	reserveId := helpers.GetUint64FromArgs(3, "reserve app id")
	dstAddress := os.Args[4]

	iconAddrBytes, err := hex.DecodeString(dstAddress[2:])
	if err != nil {
		log.Fatalf("Failed to decode hex to byte slice: %+v \n", err)
	}

	asaId := helpers.GetUint64FromArgs(5, "asset id")
	amount := helpers.GetUint64FromArgs(6, "amount")
	
	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		log.Fatalf("Failed to create algod client: %v", err)
	}

	txParams, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested tx params: %s\n", err)
	}

	var atc = future.AtomicTransactionComposer{}
	signer := future.BasicAccountTransactionSigner{Account: sender}


	contract, mcp, err := helpers.InitABIContract(client, sender, filepath.Join(tealDir, "contract.json"), reserveId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	mcp.ForeignApps = []uint64{bmcId}
	mcp.SuggestedParams.Fee = 2000

	err = atc.AddMethodCall(helpers.CombineMethod(mcp, helpers.GetMethod(contract, "burn"), []interface{}{amount, false, iconAddrBytes}))

	if err != nil {
		log.Fatalf("Failed to add method burn call: %+v \n", err)
		return
	}

	reserveAddress := crypto.GetApplicationAddress(reserveId)

	assetTxn, err := transaction.MakeAssetTransferTxn(
		sender.Address.String(),
		reserveAddress.String(),
		"",
		amount,
		uint64(txParams.Fee),
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		[]byte(nil),
		txParams.GenesisID,
		base64.StdEncoding.EncodeToString(txParams.GenesisHash),
		asaId,
	)

	if err != nil {
		log.Fatalf("Cannot create asset transfer transaction: %s\n", err)
	}

	assetTxnWithSigner := future.TransactionWithSigner{
		Txn:    assetTxn,
		Signer: signer,
	}

	atc.AddTransaction(assetTxnWithSigner)

	_, err = atc.Execute(client, context.Background(), 2)
	
	if err != nil {
		log.Fatalf("Failed to execute call: %+v \n", err)
	}
}