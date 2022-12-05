package testtools

import (
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func BtsTestInit(t *testing.T, client *algod.Client, btsTealDir string, deployer crypto.Account,
	txParams types.SuggestedParams,
) uint64 {
	t.Helper()

	btsAppCreationTx := MakeBtsDeployTx(t, client, btsTealDir, deployer, txParams)
	btsCreationTxId := SendTransaction(t, client, deployer.PrivateKey, btsAppCreationTx)
	deployRes := WaitForConfirmationsT(t, client, []string{btsCreationTxId})

	fmt.Printf("App ID: %d \n", deployRes[0].ApplicationIndex)
	return deployRes[0].ApplicationIndex
}
