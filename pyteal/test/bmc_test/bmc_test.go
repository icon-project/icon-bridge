package tests

import (
	"bytes"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/config"
	"appliedblockchain.com/icon-bridge/internalABI"
	bmcMethods "appliedblockchain.com/icon-bridge/internalABI/methods/bmc"
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
var bts_app_id uint64
var bmc_app_id uint64
var bmc_contract *abi.Contract
var bmc_mcp future.AddMethodCallParams
var err error

func Test_Init (t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bts_app_id = tools.BtsTestInit(t, client, config.BtsTealDir, deployer, txParams)
	bmc_app_id = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)
	
	bmc_contract, bmc_mcp, err = internalABI.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmc_app_id)

	if err != nil {
		t.Fatalf("Failed to init ABI contract: %+v", err)
	}
}

func Test_CallSendMessageFromOutsideOfBts(t *testing.T) {
	_, err = bmcMethods.RegisterBSHContract(client, bts_app_id, bmc_contract, bmc_mcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	_, err = bmcMethods.SendMessage(client, bmc_contract, bmc_mcp)

	if err == nil {
		t.Fatal("SendMessage should throw error, as it's not been called from BTS contract")
	}
}

func Test_RegisterRelayer(t *testing.T) {
	_, err = bmcMethods.RegisterRelayer(client, deployer.Address, bmc_contract, bmc_mcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	appRelayerAddress := tools.GetGlobalStateByKey(t, client, bmc_app_id, "relayer_acc_address")

	if !bytes.Equal(appRelayerAddress, deployer.Address[:]) {
		t.Fatal("Failed to align relayer address to address in global state of BMC application")
	}
}

func Test_CallHandleRelayMessageUsingRelayerAsSender(t *testing.T) {
	ret, err := bmcMethods.HandleRelayMessage(client, bts_app_id, "btp message", bmc_contract, bmc_mcp)

	if err != nil {
		t.Fatalf("Failed to add method call: %+v", err)
	}

	for _, r := range ret.MethodResults {
		if r.ReturnValue != "event:start handleBTPMessage" {
			t.Fatal("Failed to get event after running handleBTPMessage")
		}
	}
}