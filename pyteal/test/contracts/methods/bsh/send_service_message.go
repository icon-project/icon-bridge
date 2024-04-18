package contracts

import (
	"context"
	"fmt"

	"appliedblockchain.com/icon-bridge/config"
	tools "appliedblockchain.com/icon-bridge/contracts/tools"
	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"
)

func SendServiceMessage(client *algod.Client, bshContract *abi.Contract, mcp future.AddMethodCallParams, args []interface{}) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	err = atc.AddMethodCall(tools.CombineMethod(mcp, tools.GetMethod(bshContract, "sendServiceMessage"), args))

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