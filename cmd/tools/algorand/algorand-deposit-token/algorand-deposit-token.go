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
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")

	tealDir := os.Args[1]
	bmcId := helpers.GetUint64FromArgs(2, "bmc id")
	escrowId := helpers.GetUint64FromArgs(3, "escrow app id")
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

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	minter, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to create account from private key: %v", err)
	}

	var atc = future.AtomicTransactionComposer{}
	signer := future.BasicAccountTransactionSigner{Account: minter}


	contract, mcp, err := helpers.InitABIContract(client, minter, filepath.Join(tealDir, "contract.json"), escrowId)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	mcp.ForeignApps = []uint64{bmcId}
	mcp.SuggestedParams.Fee = 2000

	err = atc.AddMethodCall(helpers.CombineMethod(mcp, helpers.GetMethod(contract, "deposit"), []interface{}{amount, false, iconAddrBytes}))

	if err != nil {
		log.Fatalf("Failed to add method sendServiceMessage call: %+v \n", err)
		return
	}

	escrowAddress := crypto.GetApplicationAddress(escrowId)

	assetTxn, err := transaction.MakeAssetTransferTxn(
		minter.Address.String(),
		escrowAddress.String(),
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



