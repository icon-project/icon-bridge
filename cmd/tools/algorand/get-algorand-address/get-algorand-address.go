package main

import (
	"fmt"

	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	deployer := helpers.GetAccountFromPrivateKey()
	fmt.Println(deployer.Address.String())
}