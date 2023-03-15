package helpers

import (
	"encoding/base64"
	"log"

	"github.com/algorand/go-algorand-sdk/crypto"
)

func GetAccountFromPrivateKey() crypto.Account {
	privateKeyStr := GetEnvVar("PRIVATE_KEY")

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	account, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Cannot create account: %s", err)
	}

	return account
}
