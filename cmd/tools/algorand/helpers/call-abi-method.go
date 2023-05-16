package helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"
)

func CallAbiMethod(client *algod.Client, contract *abi.Contract, mcp future.AddMethodCallParams, name string, args []interface{}) (ret future.ExecuteResult, err error) {
	var atc = future.AtomicTransactionComposer{}

	method, err := contract.GetMethodByName(name)

	if err != nil {
		log.Fatalf("No method named: %s", name)
	}

	mcp.Method = method
	mcp.MethodArgs = args

	err = atc.AddMethodCall(mcp)

	if err != nil {
		fmt.Printf("Failed to add method %s call: %+v \n", name, err)
		return
	}

	ret, err = atc.Execute(client, context.Background(), 2)

	if err != nil {
		fmt.Printf("Failed to execute call: %+v \n", err)
		return
	}

	return
}