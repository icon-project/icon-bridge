package tests

import (
	"log"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/config"
	"appliedblockchain.com/icon-bridge/internalABI"
	bmcMethods "appliedblockchain.com/icon-bridge/internalABI/methods/bmc"
	tools "appliedblockchain.com/icon-bridge/testtools"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

var client *algod.Client
var deployer crypto.Account
var txParams types.SuggestedParams
var bts_app_id uint64
var bmc_app_id uint64

func TestCallSendMessageFromOutsideOfBts(t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bts_app_id = tools.BtsTestInit(t, client, config.BtsTealDir, deployer, txParams)
	bmc_app_id = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)

	contract, mcp, err := internalABI.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmc_app_id)

	if err != nil {
		log.Fatalf("Failed to init ABI contract: %+v", err)
	}

	_, err = bmcMethods.RegisterBSHContract(client, bts_app_id, contract, mcp)

	if err != nil {
		log.Fatalf("Failed to add method call: %+v", err)
	}

	_, err = bmcMethods.SendMessage(client, contract, mcp)

	if err == nil {
		log.Fatal("SendMessage should throw error, as it's not been called from BTS contract")
	}
}
