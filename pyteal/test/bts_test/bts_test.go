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
var bts_app_id uint64
var bmc_app_id uint64

func TestBtsSendServiceMessage(t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bts_app_id = tools.BtsTestInit(t, client, config.BtsTealDir, deployer, txParams)
	bmc_app_id = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)

	bts_contract, bts_mcp, err := internalABI.InitABIContract(client, deployer, filepath.Join(config.BtsTealDir, "contract.json"), bts_app_id)

	if err != nil {
		log.Fatalf("Failed to init BTS ABI contract: %+v", err)
	}

	bmc_contract, bmc_mcp, err := internalABI.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmc_app_id)

	if err != nil {
		log.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bmcApp, err := client.GetApplicationByID(bmc_app_id).Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get BMC contract application id: %+v", err)
	}

	_, err = bmcMethods.RegisterBSHContract(client, bts_app_id, bmc_contract, bmc_mcp)

	if err != nil {
		log.Fatalf("Failed to execute RegisterBSHContract: %+v", err)
	}

	_, err = bshMethods.SendServiceMessage(client, bmcApp.Id, bts_contract, bts_mcp)

	if err != nil {
		log.Fatalf("Failed to execute SendServiceMessage: %+v", err)
	}
}

func TestGetMessagePushedFromBmcToRelayer(t *testing.T) {
	round := tools.GetLatestRound(t, client)

	newBlock := tools.GetBlock(t, client, round)

	txns := algorand.GetTxns(&newBlock, bts_app_id)

	if txns == nil {
		log.Fatalf("No txns containing btp msgs")
	}

	for _, txn := range *txns {
		log.Printf("%+v\n", txn.EvalDelta)
	}
}