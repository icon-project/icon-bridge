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

func SendServiceMessage(client *algod.Client, bmc_id uint64, bsh_contract *abi.Contract, mcp future.AddMethodCallParams) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	err = atc.AddMethodCall(toolsABI.CombineMethod(mcp, toolsABI.GetMethod(bsh_contract, "sendServiceMessage"), []interface{}{bmc_id, "ICON", "TOKEN_TRANSFER_SERVICE", 3}))

	if err != nil {
		fmt.Printf("Failed to add method SendServiceMessage call into BSH contract: %+v \n", err)
		return
	}

	ret, err = atc.Execute(client, context.Background(), config.TransactionWaitRounds)

	if err != nil {
		fmt.Printf("Failed to execute call: %+v \n", err)
		return
	}

	return
}
