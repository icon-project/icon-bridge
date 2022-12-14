package algo

import (
	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/future"
)

const contractDir = "/contract/bmc.json"

func getMethod(c *abi.Contract, name string) (abi.Method, error) {
	m, err := c.GetMethodByName(name)
	if err != nil {
		return abi.Method{}, err
	}
	return m, nil
}

func combine(mcp future.AddMethodCallParams, m abi.Method,
	a []interface{}) future.AddMethodCallParams {
	mcp.Method = m
	mcp.MethodArgs = a
	return mcp
}
