package tests

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/config"
	contracts "appliedblockchain.com/icon-bridge/contracts"
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
var bshContract *abi.Contract
var bshMcp future.AddMethodCallParams
var err error

const dummyBTPMessage = "Hello Algorand"
const dummyToAddress = "btp://0x1.icon/0x12333"
const dummyServiceName = "dbsh"

func Test_Init(t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bmcAppId = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)
	bshAppId = tools.BshTestInit(t, client, config.BshTealDir, deployer, txParams)

	bmcContract, bmcMcp, err = contracts.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmcAppId)

	if err != nil {
		t.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bshContract, bshMcp, err = contracts.InitABIContract(client, deployer, filepath.Join(config.BshTealDir, "contract.json"), bshAppId)
	bshMcp.ForeignApps = []uint64{bmcAppId}

	if err != nil {
		t.Fatalf("Failed to init BSH ABI contract: %+v", err)
	}

	bshAddress := crypto.GetApplicationAddress(bshAppId)
	txnIds := tools.TransferAlgos(t, client, txParams, deployer, []types.Address{bshAddress}, 514000)
	tools.WaitForConfirmationsT(t, client, txnIds)
}

func Test_RelayerAsDeployer(t *testing.T) {
	appRelayerAddress := tools.GetGlobalStateByKey(t, client, bmcAppId, "relayer_acc_address")

	if !bytes.Equal(appRelayerAddress, deployer.Address[:]) {
		t.Fatal("Failed to align relayer address to address in global state of BMC application")
	}
}

func Test_CallSendMessageWithoutInitBsh(t *testing.T) {
	_, err = contracts.CallAbiMethod(client, bshContract, bshMcp, "sendServiceMessage", []interface{}{})

	if err == nil {
		t.Fatal("Should throw exception, that BSH is not initialized properly")
	}
}

func Test_InitBsh(t *testing.T) {
	_, err = contracts.CallAbiMethod(client, bshContract, bshMcp, "init", []interface{}{bmcAppId, dummyToAddress})

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}
}

func Test_InitMoreThenOnce(t *testing.T) {
	_, err = contracts.CallAbiMethod(client, bshContract, bshMcp, "init", []interface{}{bmcAppId, dummyToAddress})

	if err == nil {
		t.Fatal("Should assert if init was called before")
	}
}

func Test_CallSendMessageFromUnregisteredBsh(t *testing.T) {
	_, err = contracts.CallAbiMethod(client, bshContract, bshMcp, "sendServiceMessage", []interface{}{})

	if err == nil {
		t.Fatal("Should throw exception, that BSH is not registered")
	}
}

func Test_RegisterBSHContract(t *testing.T) {
	bshAddress := crypto.GetApplicationAddress(bshAppId)

	bmcMcp.ForeignAccounts = []string{bshAddress.String()}
	
	_, err = contracts.CallAbiMethod(client, bmcContract, bmcMcp, "registerBSHContract", []interface{}{bshAddress, dummyServiceName})

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	info, err := client.AccountApplicationInformation(bshAddress.String(), bmcAppId).Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to get application information: %+v", err)
	}

	fmt.Printf("%+v\n", info)
}

func Test_CallSendMessageFromBsh(t *testing.T) {
	_, err = contracts.CallAbiMethod(client, bshContract, bshMcp, "sendServiceMessage", []interface{}{})

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
	_, err = contracts.CallAbiMethod(client, bmcContract, bmcMcp, "handleRelayMessage", []interface{}{bshAppId, dummyServiceName, dummyBTPMessage})

	if err != nil {
		t.Fatalf("Failed to add call handleRelayMessage method: %+v", err)
	}

	lastReceivedMesage := tools.GetGlobalStateByKey(t, client, bshAppId, "last_received_message")

	if !bytes.Equal(lastReceivedMesage, []byte(dummyBTPMessage)) {
		t.Fatal("Failed to update last_received_message")
	}
}