package testtools

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
)

func GetAccount(t *testing.T, accountIndex int) crypto.Account {
	accounts, err := GetAccounts()

	if err != nil {
		t.Fatalf("Failed to get accounts: %+v", err)
	}

	if (len(accounts) <= accountIndex) {
		t.Fatal("Failed to get account by specified index")
	}

	return accounts[accountIndex]
}
