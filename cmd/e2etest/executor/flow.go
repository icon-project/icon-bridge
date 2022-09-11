package executor

import (
	"context"
	"fmt"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

func (ex *executor) RunFlowTest(ctx context.Context, srcChainName, dstChainName chain.ChainType, coinNames []string) error {
	if srcChainName == dstChainName {
		return fmt.Errorf("Src and Dst Chain should be different")
	}

	for _, coin := range coinNames {
		for _, cb := range []Script{
			TransferBatchBiDirection,
		} {
			if cb.Callback != nil {
				id, err := ex.getID()
				if err != nil {
					return errors.Wrap(err, "getID ")
				}
				log := ex.log.WithFields(log.Fields{"pid": id})
				sinkChan := make(chan *evt)
				ex.addChan(id, sinkChan)
				defer ex.removeChan(id)

				ts := &testSuite{
					id:              id,
					logger:          log,
					subChan:         sinkChan,
					clsPerChain:     ex.clientsPerChain,
					godKeysPerChain: ex.godKeysPerChain,
					cfgPerChain:     ex.cfgPerChain,
				}
				fmt.Printf("%v %v %v %v \n", cb.Name, srcChainName, dstChainName, coin)
				_, err = cb.Callback(ctx, srcChainName, dstChainName, coinNames, ts)
				if err != nil {
					return fmt.Errorf("%v Err: %v ", cb.Name, err)
				}
			}
		}
		return nil
	}

	return nil
}

/*
	srcCl, ok := ex.clientsPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("Client for chain %v not found", srcChainName)
	}
	dstCl, ok := ex.clientsPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("Client for chain %v not found", dstChainName)
	}
	srcGod, ok := ex.godKeysPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("GodKeys for chain %v not found", srcChainName)
	}
	dstGod, ok := ex.godKeysPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("GodKeys for chain %v not found", dstChainName)
	}
	srcCfg, ok := ex.cfgPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
	}
	dstCfg, ok := ex.cfgPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
	}
*/
