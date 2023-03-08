package tests

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/algorand"
	"appliedblockchain.com/icon-bridge/config"
	contracts "appliedblockchain.com/icon-bridge/contracts"
	contracttools "appliedblockchain.com/icon-bridge/contracts/tools"
	tools "appliedblockchain.com/icon-bridge/testtools"
	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

var client *algod.Client
var deployer, algoMinter crypto.Account
var asaId uint64
var txParams types.SuggestedParams
var bmcAppId uint64
var bmcContract *abi.Contract
var bmcMcp future.AddMethodCallParams
var escrowAppId uint64
var escrowContract *abi.Contract
var escrowMcp future.AddMethodCallParams
var escrowAddress types.Address
var err error

const dummyBTPMessage = "Hello Algorand"
const dummyToAddress = "btp://0x1.icon/0x12333"
const dummyServiceName = "wtt"
const dummyIconAddress = "hx3a94e17b282e5a8718c5e4a91010be7901d3b271"

func Test_Init(t *testing.T) {
	client, deployer, txParams = tools.Init(t)
	
	algoMinter = tools.GetAccount(t, 1)

	bmcAppId = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)
	escrowAppId = tools.EscrowTestInit(t, client, config.EscrowTealDir, deployer, txParams)

	bmcContract, bmcMcp, err = contracts.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmcAppId)

	if err != nil {
		t.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	escrowContract, escrowMcp, err = contracts.InitABIContract(client, deployer, filepath.Join(config.EscrowTealDir, "contract.json"), escrowAppId)
	escrowMcp.ForeignApps = []uint64{bmcAppId}

	if err != nil {
		t.Fatalf("Failed to init Escrow ABI contract: %+v", err)
	}

	escrowAddress = crypto.GetApplicationAddress(escrowAppId)
	txnIds := tools.TransferAlgos(t, client, txParams, deployer, []types.Address{escrowAddress}, 614000)
	tools.WaitForConfirmationsT(t, client, txnIds)
}

func Test_DeployASA(t *testing.T) {
	mintTx, err := algorand.MintTx(txParams, algoMinter.Address, 1000000000000, 0, "ABC", "AB Coin",
		"http://example.com/", "abcd")
	
	if err != nil {
		t.Fatalf("Could not generate asset creation transaction: %s", err)
	}

	mintTxId := tools.SendTransaction(t, client, algoMinter.PrivateKey, mintTx)
	res := tools.WaitForConfirmationsT(t, client, []string{mintTxId})

	asaId = res[0].AssetIndex
	
	log.Print(asaId)
}

func Test_InitEscrow(t *testing.T) {
	escrowMcp.ForeignAssets = []uint64{asaId}
	_, err = contracts.CallAbiMethod(client, escrowContract, escrowMcp, "init", []interface{}{bmcAppId, dummyToAddress, asaId})

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}
}

func Test_RegisterEscrowContract(t *testing.T) {
	escrowAddress := crypto.GetApplicationAddress(escrowAppId)

	bmcMcp.ForeignAccounts = []string{escrowAddress.String()}
	
	_, err = contracts.CallAbiMethod(client, bmcContract, bmcMcp, "registerBSHContract", []interface{}{escrowAddress, dummyServiceName})

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	info, err := client.AccountApplicationInformation(escrowAddress.String(), bmcAppId).Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to get application information: %+v", err)
	}

	fmt.Printf("%+v\n", info)
}

func Test_CallSendMessageFromEscrow(t *testing.T) {
  iconAddrBytes, err := hex.DecodeString(dummyIconAddress[2:])
	if err != nil {
		t.Fatalf("Failed to decode hex to byte slice: %+v \n", err)
	}
	
	var atc = future.AtomicTransactionComposer{}
	signer := future.BasicAccountTransactionSigner{Account: algoMinter}

	escrowMcp.Sender = algoMinter.Address
	escrowMcp.Signer = signer
	
	err = atc.AddMethodCall(contracttools.CombineMethod(escrowMcp, contracttools.GetMethod(escrowContract, "deposit"), []interface{}{1, false, iconAddrBytes}))

	if err != nil {
		t.Fatalf("Failed to add method sendServiceMessage call: %+v \n", err)
		return
	}

	assetTxn, err := algorand.TransferAssetTx(txParams, algoMinter.Address, escrowAddress, asaId, 1)

	if err != nil {
		t.Fatalf("Cannot create asset transfer transaction: %s\n", err)
	}

	assetTxnWithSigner := future.TransactionWithSigner{
    Txn:    assetTxn,
    Signer: signer,
	}
	
	atc.AddTransaction(assetTxnWithSigner)

	_, err = atc.Execute(client, context.Background(), config.TransactionWaitRounds)

	if err != nil {
		t.Fatalf("Failed to execute call: %+v \n", err)
	}

	round := tools.GetLatestRound(t, client)

	newBlock := tools.GetBlock(t, client, round)

	for _, stxn := range newBlock.Payset {
		for _, innertxn := range stxn.EvalDelta.InnerTxns {
			if innertxn.EvalDelta.Logs[0] != dummyServiceName {
				t.Fatal("Service name is not valid")
			}
		}
	}
}

func Test_CallHandleRelayMessageUsingRelayerAsSender(t *testing.T) {
	_, err = contracts.CallAbiMethod(client, bmcContract, bmcMcp, "handleRelayMessage", []interface{}{escrowAppId, dummyServiceName, []byte(dummyBTPMessage)})

	if err != nil {
		t.Fatalf("Failed to add call handleRelayMessage method: %+v", err)
	}

	lastReceivedMesage := tools.GetGlobalStateByKey(t, client, escrowAppId, "last_received_message")

	if !bytes.Equal(lastReceivedMesage, []byte(dummyBTPMessage)) {
		t.Fatal("Failed to update last_received_message")
	}
}