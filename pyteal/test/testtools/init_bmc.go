package testtools

import (
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func BmcTestInit(t *testing.T, client *algod.Client, bmcTealDir string, deployer crypto.Account,
	txParams types.SuggestedParams,
) uint64 {
	t.Helper()

	bmcAppCreationTx := MakeBmcDeployTx(t, client, bmcTealDir, deployer, txParams)
	bmcCreationTxId := SendTransaction(t, client, deployer.PrivateKey, bmcAppCreationTx)
	deployRes := WaitForConfirmationsT(t, client, []string{bmcCreationTxId})

	fmt.Printf("BMC App ID: %d \n", deployRes[0].ApplicationIndex)
	return deployRes[0].ApplicationIndex
}
