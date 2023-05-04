package testtools

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

func ReserveTestInit(t *testing.T, client *algod.Client, tealDirPath string, deployer crypto.Account,
	txParams types.SuggestedParams,
) uint64 {
	t.Helper()

	approvalSourceCode := ReadFileT(t, filepath.Join(tealDirPath, "approval.teal"))
	clearStateSourceCode := ReadFileT(t, filepath.Join(tealDirPath, "clear.teal"))

	approvalProgram := CompileT(t, client, approvalSourceCode)
	clearStateProgram := CompileT(t, client, clearStateSourceCode)

	txn, err := future.MakeApplicationCreateTx(
		false,
		approvalProgram,
		clearStateProgram,
		types.StateSchema{NumUint: 3, NumByteSlice: 1},
		types.StateSchema{NumUint: 0, NumByteSlice: 0},
		[][]byte{},
		nil,
		nil,
		nil,
		txParams,
		deployer.Address,
		nil,
		types.Digest{},
		[32]byte{},
		types.Address{},
	)

	if err != nil {
		t.Fatalf("Failed to make Reserve application creation transaction: %s\n", err)
	}
	
	txId := SendTransaction(t, client, deployer.PrivateKey, txn)
	deployRes := WaitForConfirmationsT(t, client, []string{txId})

	fmt.Printf("Reserve App ID: %d \n", deployRes[0].ApplicationIndex)
	return deployRes[0].ApplicationIndex
}
