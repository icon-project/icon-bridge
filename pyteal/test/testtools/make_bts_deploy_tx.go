package testtools

import (
	"path/filepath"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

func MakeBtsDeployTx(
	t *testing.T,
	client *algod.Client,
	tealDirPath string,
	deployer crypto.Account,
	txParams types.SuggestedParams) (tx types.Transaction) {
	approvalSourceCode := ReadFileT(t, filepath.Join(tealDirPath, "approval.teal"))
	clearStateSourceCode := ReadFileT(t, filepath.Join(tealDirPath, "clear.teal"))

	approvalProgram := CompileT(t, client, approvalSourceCode)
	clearStateProgram := CompileT(t, client, clearStateSourceCode)

	tx, err := future.MakeApplicationCreateTx(
		false,
		approvalProgram,
		clearStateProgram,
		types.StateSchema{NumUint: 4, NumByteSlice: 4},
		types.StateSchema{NumUint: 4, NumByteSlice: 4},
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
		t.Fatalf("Failed to make BTS application creation transaction: %s\n", err)
	}
	return
}
