package algo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

const contractDir = "../../../../pyteal/teal/bmc/"
const waitRounds = 5

type AbiFunc struct {
	name string
	args []interface{}
}

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
	abiPath, err := filepath.Abs(contractDir + "contract.json")
	if err != nil {
		return fmt.Errorf("Couldn't retrieve abi file: %w", err)
	}
	rawBmc, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return fmt.Errorf("Failed to open contract file: %w", err)
	}
	abiBmc := &abi.Contract{}
	if err = json.Unmarshal(rawBmc, abiBmc); err != nil {
		return fmt.Errorf("Failed to marshal abi contract: %w", err)
	}
	s.bmc = abiBmc

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	sp, err := s.cl.algod.SuggestedParams().Do(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get suggeted params: %w", err)
	}
	s.mcp = &future.AddMethodCallParams{
		AppID:           s.opts.AppId,
		Sender:          s.wallet.TypedAddress(),
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          s.wallet,
	}
	return nil
}

func (s *sender) callAbi(ctx context.Context, abiFuncs ...AbiFunc) (future.ExecuteResult, error) {
	var atc = future.AtomicTransactionComposer{}
	for _, abiFunc := range abiFuncs {
		method, err := getMethod(s.bmc, abiFunc.name)
		if err != nil {
			return future.ExecuteResult{}, fmt.Errorf("Failed to get %s method from json contract: %w",
				abiFunc.name, err)
		}
		err = atc.AddMethodCall(combine(*s.mcp, method, abiFunc.args))
		if err != nil {
			return future.ExecuteResult{}, fmt.Errorf("Failed to add %s method to atc: %w", abiFunc.name, err)
		}
	}
	ret, err := atc.Execute(s.cl.algod, ctx, waitRounds)
	if err != nil {
		return future.ExecuteResult{}, fmt.Errorf("Failed to execute atc: %w", err)
	}
	return ret, nil
}

// EncodeMessage encodes a receipt into a byte array,
// to then be decoded inside the bmc
func encodeReceipt(relayMsg *chain.Receipt) ([]byte, error) {
	next, err := abi.TypeOf("string")
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of next bmc: %w", err)
	}
	seq, err := abi.TypeOf("uint64")
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of sequence: %w", err)
	}
	msg, err := abi.TypeOf("byte[]")
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of message: %w", err)
	}
	event, err := abi.MakeTupleType([]abi.Type{next, seq, msg})
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of event: %w", err)
	}

	encEvent, err := event.Encode([]interface{}{string(relayMsg.Events[0].Next), relayMsg.Events[0].Sequence, relayMsg.Events[0].Message})
	if err != nil {
		return nil, fmt.Errorf("Failed to encode event: %w", err)
	}

	index, err := abi.TypeOf("uint64")
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of index: %w", err)
	}

	height, err := abi.TypeOf("uint64")
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of height: %w", err)
	}

	receipt, err := abi.MakeTupleType([]abi.Type{index, msg, height})
	if err != nil {
		return nil, fmt.Errorf("Failed to encode receipt: %w", err)
	}

	encMessage, err := receipt.Encode([]interface{}{relayMsg.Index, encEvent, relayMsg.Height})
	if err != nil {
		return nil, fmt.Errorf("Failed to create type of relay message: %w", err)
	}
	return encMessage, nil
}
