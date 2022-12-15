package algo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

const contractDir = "bmc/contract.json"

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

func (s *sender) initAbi() error {
	rawBmc, err := ioutil.ReadFile(contractDir)
	if err != nil {
		return fmt.Errorf("Failed to open contract file: %w", err)
	}
	abiBmc := &abi.Contract{}
	if err = json.Unmarshal(rawBmc, abiBmc); err != nil {
		return fmt.Errorf("Failed to marshal abi contract: %w", err)
	}
	s.bmc = abiBmc
	sp, err := s.cl.algod.SuggestedParams().Do(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get suggeted params: %w", err)
	}
	sp.Fee = 0
	s.mcp = &future.AddMethodCallParams{
		AppID:           s.opts.AppId,
		Sender:          s.wallet.TypedAddress(),
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          s.wallet,
	}
	return nil
}

func (s *sender) callAbi(name string, args []interface{}) (future.ExecuteResult, error) {
	var atc = future.AtomicTransactionComposer{}
	method, err := getMethod(s.bmc, name)
	if err != nil {
		return future.ExecuteResult{}, fmt.Errorf("Failed to get %s method from json contract: %w",
			name, err)
	}
	err = atc.AddMethodCall(combine(*s.mcp, method, args))
	if err != nil {
		return future.ExecuteResult{}, fmt.Errorf("Failed to add %s method to atc: %w", name, err)
	}
	ret, err := atc.Execute(s.cl.algod, context.Background(), 5)
	if err != nil {
		return future.ExecuteResult{}, fmt.Errorf("Failed to execute atc: %w", err)
	}
	return ret, nil
}
