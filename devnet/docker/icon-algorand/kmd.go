package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

type KeyResponse struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: go run main.go <kmd-address> <kmd-token> <wallet-name>")
		return
	}

	kmdAddress := os.Args[1]
	kmdToken := os.Args[2]
	walletName := os.Args[3]

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}

	// Check if the wallet already exists
	walletsResponse, err := kmdClient.ListWallets()
	if err != nil {
		fmt.Printf("Error listing wallets: %s\n", err)
		return
	}

	var exampleWalletID string
	for _, wallet := range walletsResponse.Wallets {
		if wallet.Name == walletName {
			exampleWalletID = wallet.ID
			break
		}
	}

	// If the wallet doesn't exist, create it without a password
	if exampleWalletID == "" {
		cwResponse, err := kmdClient.CreateWallet(walletName, "", kmd.DefaultWalletDriver, types.MasterDerivationKey{})
		if err != nil {
			fmt.Printf("error creating wallet: %s\n", err)
			return
		}
		exampleWalletID = cwResponse.Wallet.ID
	}

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err := kmdClient.InitWalletHandle(exampleWalletID, "")
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken := initResponse.WalletHandleToken

	// Generate a new address from the wallet handle
	genResponse, err := kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}

	address := genResponse.Address

	keyResponse, err := kmdClient.ExportKey(initResponse.WalletHandleToken, "", address)
	if err != nil {
		fmt.Printf("Error exporting key: %s\n", err)
		return
	}
	deployer, err := crypto.AccountFromPrivateKey(keyResponse.PrivateKey)
	if err != nil {
		fmt.Printf("Cannot create deployer account: %s", err)
	}
	key := KeyResponse{
		Address:    deployer.Address.String(),
		PrivateKey: base64.StdEncoding.EncodeToString(keyResponse.PrivateKey),
	}

	jsonData, err := json.Marshal(key)
	if err != nil {
		fmt.Printf("Failed to marshal key data: %s\n", err)
		return
	}

	fmt.Println(string(jsonData))
}
