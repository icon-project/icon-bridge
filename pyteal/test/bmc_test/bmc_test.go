package tests

import (
	"bytes"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/algorand"
	"appliedblockchain.com/icon-bridge/config"
	contracts "appliedblockchain.com/icon-bridge/contracts"
	bmcmethods "appliedblockchain.com/icon-bridge/contracts/methods/bmc"
	tools "appliedblockchain.com/icon-bridge/testtools"
	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

var client *algod.Client
var deployer crypto.Account
var txParams types.SuggestedParams
var bshAppId uint64
var bmcAppId uint64
var bmcContract *abi.Contract
var bmcMcp future.AddMethodCallParams
var err error

const dummyBTPMessage = "btp message"

func Test_Init(t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bshAppId = tools.BshTestInit(t, client, config.BshTealDir, deployer, txParams)
	bmcAppId = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)

	bmcContract, bmcMcp, err = contracts.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmcAppId)

	if err != nil {
		t.Fatalf("Failed to init ABI contract: %+v", err)
	}
}

func Test_RelayerAsDeployer(t *testing.T) {
	appRelayerAddress := tools.GetGlobalStateByKey(t, client, bmcAppId, "relayer_acc_address")

	if !bytes.Equal(appRelayerAddress, deployer.Address[:]) {
		t.Fatal("Failed to align relayer address to address in global state of BMC application")
	}
}

func Test_SetRelayer(t *testing.T) {
	relayer := crypto.GenerateAccount()

	_, err = bmcmethods.SetRelayer(client, relayer.Address, bmcContract, bmcMcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	appRelayerAddress := tools.GetGlobalStateByKey(t, client, bmcAppId, "relayer_acc_address")

	if !bytes.Equal(appRelayerAddress, relayer.Address[:]) {
		t.Fatal("Failed to align relayer address to address in global state of BMC application")
	}
}

func Test_RegisterBSHContract(t *testing.T) {
	_, err = bmcmethods.RegisterBSHContract(client, bshAppId, bmcContract, bmcMcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}
	
	bshAddress := crypto.GetApplicationAddress(bshAppId)
	globalBshAddress := tools.GetGlobalStateByKey(t, client, bmcAppId, "bsh_app_address")

	if !bytes.Equal(globalBshAddress, bshAddress[:]) {
		t.Fatal("Failed to align BSH address to address in global state of BMC application")
	}
}

func Test_CallSendMessageFromOutsideOfBsh(t *testing.T) {
	_, err = bmcmethods.SendMessage(client, bmcContract, bmcMcp)

	if err == nil {
		t.Fatal("SendMessage should throw error, as it's not been called from BSH contract")
	}
}

func Test_CallHandleRelayMessageUsingRelayerAsSender(t *testing.T) {
	_, err = bmcmethods.SetRelayer(client, deployer.Address, bmcContract, bmcMcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	ret, err := bmcmethods.HandleRelayMessage(client, bshAppId, dummyBTPMessage, bmcContract, bmcMcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	for _, r := range ret.MethodResults {
		if r.ReturnValue != "event:start handleBTPMessage" {
			t.Fatal("Failed to get event after running handleBTPMessage")
		}
	}
}

func Test_GetMessagePushedFromRelayerToBmc(t *testing.T) {
	round := tools.GetLatestRound(t, client)

	newBlock := tools.GetBlock(t, client, round)

	txns := algorand.GetTxns(&newBlock, bmcAppId)

	if txns == nil {
		t.Fatalf("No txns containing btp msgs")
	}

	for _, txn := range *txns {
		for _, innerTxn := range txn.EvalDelta.InnerTxns {
			if innerTxn.EvalDelta.Logs[0] != dummyBTPMessage {
				t.Fatal("Failed to get BTP message pushed from relayer to BMC")
			}
		}
	}
}
