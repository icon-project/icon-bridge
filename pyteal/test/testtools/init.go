package testtools

import (
	"log"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func Init(t *testing.T) (client *algod.Client, deployer crypto.Account, txParams types.SuggestedParams) {
	client = MakeAlgodClient(t)

	fundingAccounts, err := GetAccounts()

	if err != nil {
		log.Fatalf("Failed to get accounts: %+v", err)
	}

	fundingAccount := fundingAccounts[1]

	deployer = crypto.GenerateAccount()

	txParams = SuggestedParams(t, client)

	txnIds := TransferAlgos(t, client, txParams, fundingAccount, []types.Address{deployer.Address}, 10000000)
	WaitForConfirmationsT(t, client, txnIds)

	return
}
