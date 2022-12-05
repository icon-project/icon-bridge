package internalABI

import (
	"context"
	"fmt"

	"appliedblockchain.com/icon-bridge/config"
	toolsABI "appliedblockchain.com/icon-bridge/internalABI/tools"
	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"
)

func SendMessage(client *algod.Client, bmc_contract *abi.Contract, mcp future.AddMethodCallParams) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	err = atc.AddMethodCall(toolsABI.CombineMethod(mcp, toolsABI.GetMethod(bmc_contract, "sendMessage"), []interface{}{"ICON", "TOKEN_TRANSFER_SERVICE", 3}))

	if err != nil {
		fmt.Printf("Failed to add method SendMessage call into BMC contract: %+v", err)
		return
	}

	ret, err = atc.Execute(client, context.Background(), config.TransactionWaitRounds)

	if err != nil {
		fmt.Printf("Failed to execute call: %+v", err)
		return
	}

	return
}
