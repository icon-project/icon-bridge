package testtools

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

const TestnetAlgodAddress = "http://localhost:4001"
const TestnetAlgodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func MakeAlgodClient(t *testing.T) (client *algod.Client) {
	t.Helper()

	var err error
	client, err = algod.MakeClient(TestnetAlgodAddress, TestnetAlgodToken)
	if err != nil {
		t.Fatalf("Algod client could not be created: %s\n", err)
	}
	return
}
