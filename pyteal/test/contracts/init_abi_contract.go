package contracts

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

func InitABIContract(client *algod.Client, deployer crypto.Account, contractDir string, appId uint64) (contract *abi.Contract, mcp future.AddMethodCallParams, err error) {
	b, err := ioutil.ReadFile(contractDir)
	if err != nil {
		fmt.Printf("Failed to open contract file: %+v", err)
		return
	}

	contract = &abi.Contract{}
	if err = json.Unmarshal(b, contract); err != nil {
		fmt.Printf("Failed to marshal contract: %+v", err)
		return
	}

	sp, err := client.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get suggeted params: %+v", err)
		return
	}

	sp.Fee = 1000

	signer := future.BasicAccountTransactionSigner{Account: deployer}

	mcp = future.AddMethodCallParams{
		AppID:           appId,
		Sender:          deployer.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
	}

	return
}
