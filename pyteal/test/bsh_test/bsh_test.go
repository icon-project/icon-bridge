package tests

import (
	"context"
	"log"
	"path/filepath"
	"testing"

	"appliedblockchain.com/icon-bridge/algorand"
	"appliedblockchain.com/icon-bridge/config"
	contracts "appliedblockchain.com/icon-bridge/contracts"
	bmcmethods "appliedblockchain.com/icon-bridge/contracts/methods/bmc"
	bshmethods "appliedblockchain.com/icon-bridge/contracts/methods/bsh"
	tools "appliedblockchain.com/icon-bridge/testtools"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

var client *algod.Client
var deployer crypto.Account
var txParams types.SuggestedParams
var bshAppId uint64
var bmcAppId uint64

func Test_BshSendServiceMessage(t *testing.T) {
	client, deployer, txParams = tools.Init(t)

	bshAppId = tools.BshTestInit(t, client, config.BshTealDir, deployer, txParams)
	bmcAppId = tools.BmcTestInit(t, client, config.BmcTealDir, deployer, txParams)

	bshContract, bshMcp, err := contracts.InitABIContract(client, deployer, filepath.Join(config.BshTealDir, "contract.json"), bshAppId)

	if err != nil {
		t.Fatalf("Failed to init BSH ABI contract: %+v", err)
	}

	bmcContract, bmcMcp, err := contracts.InitABIContract(client, deployer, filepath.Join(config.BmcTealDir, "contract.json"), bmcAppId)

	if err != nil {
		t.Fatalf("Failed to init BMC ABI contract: %+v", err)
	}

	bmcApp, err := client.GetApplicationByID(bmcAppId).Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to get BMC contract application id: %+v", err)
	}

	_, err = bmcmethods.RegisterBSHContract(client, bshAppId, bmcContract, bmcMcp)

	if err != nil {
		t.Fatalf("Failed to execute RegisterBSHContract: %+v", err)
	}

	_, err = bshmethods.SendServiceMessage(client, bmcApp.Id, bshContract, bshMcp)

	if err != nil {
		t.Fatalf("Failed to execute SendServiceMessage: %+v", err)
	}
}

func Test_GetMessagePushedFromBmcToRelayer(t *testing.T) {
	round := tools.GetLatestRound(t, client)

	newBlock := tools.GetBlock(t, client, round)

	txns := algorand.GetTxns(&newBlock, bshAppId)

	if txns == nil {
		t.Fatalf("No txns containing btp msgs")
	}

	for _, txn := range *txns {
		log.Printf("%+v\n", txn.EvalDelta)
	}
}
