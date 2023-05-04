package main

import (
	"encoding/hex"
	"fmt"

	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	deployer := helpers.GetAccountFromPrivateKey()

	encodedString := hex.EncodeToString(deployer.PublicKey)

	fmt.Println(encodedString)
}