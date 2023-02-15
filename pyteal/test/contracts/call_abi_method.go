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

func CallAbiMethod(client *algod.Client, contract *abi.Contract, mcp future.AddMethodCallParams, name string, args []interface{}) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	err = atc.AddMethodCall(tools.CombineMethod(mcp, tools.GetMethod(contract, name), args))

	if err != nil {
		fmt.Printf("Failed to add method %s call: %+v \n", name, err)
		return
	}

	ret, err = atc.Execute(client, context.Background(), config.TransactionWaitRounds)

	if err != nil {
		fmt.Printf("Failed to execute call: %+v \n", err)
		return
	}

	return
}
