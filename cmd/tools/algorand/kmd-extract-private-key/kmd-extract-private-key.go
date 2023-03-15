package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/icon-project/icon-bridge/cmd/tools/algorand/helpers"
)

func main() {
	kmdAddress := helpers.GetEnvVar("KMD_ADDRESS")
	kmdToken := helpers.GetEnvVar("KMD_TOKEN")

	accountIndex, err := strconv.Atoi(os.Args[1])
	
	if err != nil {
		log.Fatalf("Invalid account index %s\n", err)
	}

	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)

	if err != nil {
		log.Fatalf("KMD client could not be created: %s\n", err)
	}

	walletsResponse, err := kmdClient.ListWallets()

	if err != nil {
		log.Fatalf("cannot list wallets: %s", err)
	}

	if len(walletsResponse.Wallets) == 0 {
		log.Fatal("no wallets")
	}

	walletID := walletsResponse.Wallets[0].ID
	initResponse, err := kmdClient.InitWalletHandle(walletID, "")
	if err != nil {
		log.Fatalf("initWalletHandle failed: %s", err)
	}
	
	keysResponse, err := kmdClient.ListKeys(initResponse.WalletHandleToken)
	if err != nil {
		log.Fatalf("listKeys failed: %s", err)
	}
	if len(keysResponse.Addresses) == 0 {
		log.Fatal("no accounts in wallet")
	}
	if len(keysResponse.Addresses) < accountIndex+1 {
		log.Fatal("not enough accounts in wallet")
	}

	keyResponse, err := kmdClient.ExportKey(initResponse.WalletHandleToken, "", keysResponse.Addresses[accountIndex])
	if err != nil {
		log.Fatalf("exportKey failed: %s", err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(keyResponse.PrivateKey))
}