package main

import (
	"context"
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <algod-address> <algod-token>\n", os.Args[0])
		return
	}

	algodAddress := os.Args[1]
	algodToken := os.Args[2]

	// Create an Algod client
	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("Failed to make client: %v\n", err)
		return
	}

	// Get the last round
	status, err := client.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get status: %v\n", err)
		return
	}

	fmt.Printf("%d", status.LastRound)
}
