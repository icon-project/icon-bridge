package tests

import (
	"context"
	"log"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/algorand"
	"appliedblockchain.com/icon-bridge/config"
	"appliedblockchain.com/icon-bridge/internalABI"
	bmcMethods "appliedblockchain.com/icon-bridge/internalABI/methods/bmc"
	bshMethods "appliedblockchain.com/icon-bridge/internalABI/methods/bsh"
	tools "appliedblockchain.com/icon-bridge/testtools"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

var client *algod.Client
var deployer crypto.Account
var txParams types.SuggestedParams
var bsh_app_id uint64
var bmc_app_id uint64

func Test_BshSendServiceMessage(t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bsh_app_id = tools.BshTestInit(t, client, config.BshTealDir, deployer, txParams)
	bmc_app_id = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)

	bsh_contract, bsh_mcp, err := internalABI.InitABIContract(client, deployer, filepath.Join(config.BshTealDir, "contract.json"), bsh_app_id)

	if err != nil {
		t.Fatalf("Failed to init BSH ABI contract: %+v", err)
	}

	bmc_contract, bmc_mcp, err := internalABI.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmc_app_id)

	if err != nil {
		t.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bmcApp, err := client.GetApplicationByID(bmc_app_id).Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to get BMC contract application id: %+v", err)
	}

	_, err = bmcMethods.RegisterBSHContract(client, bsh_app_id, bmc_contract, bmc_mcp)

	if err != nil {
		t.Fatalf("Failed to execute RegisterBSHContract: %+v", err)
	}

	_, err = bshMethods.SendServiceMessage(client, bmcApp.Id, bsh_contract, bsh_mcp)

	if err != nil {
		t.Fatalf("Failed to execute SendServiceMessage: %+v", err)
	}
}

func Test_GetMessagePushedFromBmcToRelayer(t *testing.T) {
	round := tools.GetLatestRound(t, client)

	newBlock := tools.GetBlock(t, client, round)

	txns := algorand.GetTxns(&newBlock, bsh_app_id)

	if txns == nil {
		t.Fatalf("No txns containing btp msgs")
	}

	for _, txn := range *txns {
		log.Printf("%+v\n", txn.EvalDelta)
	}
}
