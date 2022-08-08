package executor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

func (ex *executor) RunFlowTest(ctx context.Context, srcChainName, dstChainName chain.ChainType, coinNames []string) error {
	if srcChainName == dstChainName {
		return fmt.Errorf("Src and Dst Chain should be different")
	}
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
	srcDemo, ok := ex.demoKeysPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("DemoKeys for chain %v not found", srcChainName)
	}
	srcDemo = append(srcDemo, srcGod)
	dstDemo, ok := ex.demoKeysPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("DemoKeys for chain %v not found", dstChainName)
	}
	dstDemo = append(dstDemo, dstGod)
	srcCfg, ok := ex.cfgPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
	}
	dstCfg, ok := ex.cfgPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
	}
	btsAddressPerChain := map[chain.ChainType]string{
		srcChainName: srcCfg.ContractAddresses[chain.BTS],
		dstChainName: dstCfg.ContractAddresses[chain.BTS],
	}
	gasLimitPerChain := map[chain.ChainType]int64{
		srcChainName: srcCfg.GasLimit,
		dstChainName: dstCfg.GasLimit,
	}

	id, err := ex.getID()
	if err != nil {
		return errors.Wrap(err, "getID ")
	}
	log := ex.log.WithFields(log.Fields{"pid": id})
	sinkChan := make(chan *evt)
	ex.addChan(id, sinkChan)
	defer ex.removeChan(id)

	ts := &testSuite{
		id:                 id,
		logger:             log,
		env:                ex.env,
		subChan:            sinkChan,
		btsAddressPerChain: btsAddressPerChain,
		gasLimitPerChain:   gasLimitPerChain,
		clsPerChain:        map[chain.ChainType]chain.ChainAPI{srcChainName: srcCl, dstChainName: dstCl},
		godKeysPerChain:    map[chain.ChainType]keypair{srcChainName: srcGod, dstChainName: dstGod},
		demoKeysPerChain:   map[chain.ChainType][]keypair{srcChainName: srcDemo, dstChainName: dstDemo},
		fee:                fee{numerator: big.NewInt(FEE_NUMERATOR), denominator: big.NewInt(FEE_DENOMINATOR), fixed: big.NewInt(FIXED_PRICE)},
	}
	for _, coin := range coinNames {
		for _, cb := range []Script{
			TransferWithApprove,
			// TransferWithoutApprove,
			// TransferToZeroAddress,
			// TransferToUnknownNetwork,
			// TransferToUnparseableAddress,
			// TransferLessThanFee,
			// TransferEqualToFee,
			// TransferExceedingBTSBalance,
		} {
			if cb.Callback != nil {
				_, err := cb.Callback(ctx, srcChainName, dstChainName, []string{coin}, ts)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
