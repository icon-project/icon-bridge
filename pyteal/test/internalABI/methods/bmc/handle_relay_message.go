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

func HandleRelayMessage(client *algod.Client, bsh_id uint64, msg string, bmc_contract *abi.Contract, mcp future.AddMethodCallParams) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	err = atc.AddMethodCall(toolsABI.CombineMethod(mcp, toolsABI.GetMethod(bmc_contract, "handleRelayMessage"), []interface{}{bsh_id, msg}))

	if err != nil {
		fmt.Printf("Failed to add method handleRelayMessage call into BMC contract: %+v \n", err)
		return
	}

	ret, err = atc.Execute(client, context.Background(), config.TransactionWaitRounds)

	if err != nil {
		fmt.Printf("Failed to execute call: %+v \n", err)
		return
	}

	return
}
