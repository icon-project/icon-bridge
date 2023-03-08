package testtools

import (
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func BshTestInit(t *testing.T, client *algod.Client, bshTealDir string, deployer crypto.Account,
	txParams types.SuggestedParams,
) uint64 {
	t.Helper()

	bshAppCreationTx := MakeBshDeployTx(t, client, bshTealDir, deployer, txParams)
	bshCreationTxId := SendTransaction(t, client, deployer.PrivateKey, bshAppCreationTx)
	deployRes := WaitForConfirmationsT(t, client, []string{bshCreationTxId})

	fmt.Printf("BSH App ID: %d \n", deployRes[0].ApplicationIndex)
	return deployRes[0].ApplicationIndex
}
