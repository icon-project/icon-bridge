package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	privateKeyStr := helpers.GetEnvVar("PRIVATE_KEY")

	privateKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		log.Fatalf("Cannot base64-decode private key seed: %s\n", err)
	}

	deployer, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Cannot create deployer account: %s", err)
	}

	encodedString := hex.EncodeToString(deployer.PublicKey)

	fmt.Println(encodedString)
}