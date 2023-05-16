package testtools

import (
	"testing"

	"appliedblockchain.com/icon-bridge/algorand"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

func CompileT(t *testing.T, client *algod.Client, teal []byte) (program []byte) {
	t.Helper()
	var err error
	program, err = algorand.Compile(client, teal)
	if err != nil {
		t.Fatalf("Compilation failed: %+v\n", err)
	}
	return
}
