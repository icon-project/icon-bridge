package executor

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

const MINIMUM_BALANCE = 1

var TransferToUnparseableAddress Script = Script{
	Name:        "TransferToUnparseableAddress",
	Type:        "Flow",
	Description: "Transfer to unparseable address",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		dstAddr += "1"
		// how much you want receiving end to get

		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = ts.withFeeAdded(big.NewInt(MINIMUM_BALANCE))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}

		if err = ts.ValidateTransactionResultAndEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Unexpected error %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "1" {
					return fmt.Errorf("Expected error code (1) Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		return nil, err
	},
}

var TransferToZeroAddress Script = Script{
	Name:        "TransferToZeroAddress",
	Description: "Transfer to zero address",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, tmpAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		convertToZeroAddress := func(inputStr string) string {
			splits := strings.Split(inputStr, "/")
			lenSplits := len(splits)
			if len(splits[lenSplits-1]) > 2 {
				prefix := splits[lenSplits-1][0:2]
				postFix := splits[lenSplits-1][2:]
				byteArr := make([]byte, len(postFix)/2)
				for i := 0; i < len(postFix)/2; i++ {
					byteArr[i] = 0
				}
				splits[lenSplits-1] = prefix + hex.EncodeToString(byteArr)
			}
			return strings.Join(splits, "/")
		}
		dstAddr := convertToZeroAddress(tmpAddr)
		// how much you want receiving end to get
		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = ts.withFeeAdded(big.NewInt(MINIMUM_BALANCE))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}

		if err = ts.ValidateTransactionResultAndEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Unexpected error %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "1" {
					return fmt.Errorf("Expected error code (1) Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		return nil, err
	},
}

var TransferToUnknownNetwork Script = Script{
	Name:        "TransferToUnknownNetwork",
	Description: "Transfer to unknown bmc network",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, tmpAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		changeBMCNetwork := func(inputStr string) (string, error) {
			splits := strings.Split(inputStr, "/")
			if len(splits) != 4 {
				return "", errors.New("Unexpected length")
			}
			network := splits[2]
			networkSplits := strings.Split(network, ".")
			if len(networkSplits) != 2 {
				return "", errors.New("Unexpected length")
			}
			networkSplits[1] += "s"
			splits[2] = strings.Join(networkSplits, ".")
			return strings.Join(splits, "/"), nil
		}
		dstAddr, err := changeBMCNetwork(tmpAddr)
		if err != nil {
			return nil, err
		}
		// how much you want receiving end to get
		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = ts.withFeeAdded(big.NewInt(MINIMUM_BALANCE))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}
		if err = ts.ValidateTransactionResultAndEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Unexpected error %v", err)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		return nil, err
	},
}

var TransferExceedingBTSBalance Script = Script{
	Name:        "TransferExceedingContractsBalance",
	Type:        "Flow",
	Description: "Transfer Native Tokens, which are of fixed supply, in different amounts. The Token should be native for both chains",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// coinName should be a token common on both chains
		tokenExists := false
		for _, stkn := range src.NativeTokens() {
			if stkn == coinName {
				for _, dtkn := range dst.NativeTokens() {
					if dtkn == coinName {
						tokenExists = true
						break
					}
				}
				break
			}
		}
		if !tokenExists {
			ts.logger.Errorf("Token %v does not exist on both chains %v and %v", coinName, srcChain, dstChain)
			return nil, nil
		}

		btsAddr, ok := ts.btsAddressPerChain[dstChain]
		if !ok {
			return nil, errors.Wrapf(err, "dst.GetBTPAddressOfBTS %v", err)
		}
		btsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "dst.getCoinBalance %v", err)
		}

		// prepare accounts
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		amt := ts.withFeeAdded(btsBalance.UserBalance)
		amt.Add(amt, big.NewInt(MINIMUM_BALANCE)) //exceed balance by 100g
		if err := ts.Fund(srcAddr, amt, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund %v", err)
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// approve
		if approveHash, err := src.Approve(coinName, srcKey, amt); err != nil {
			return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		} else {
			if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
				return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
			}
		}
		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, amt)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer %v", err)
		}
		if err := ts.ValidateTransactionResultAndEvents(ctx, hash, []string{coinName}, srcAddr, dstAddr, []*big.Int{amt}); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			// chain.TransferReceived: func(e *evt) error {
			// 	ts.logger.Info("Got TransferReceived")
			// 	return nil
			// },
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent. Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() == "1" { //&& endEvt.Response == "TransferFailed" {
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		// finalBtsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "dst.getCoinBalance %v", err)
		// }
		// if finalBtsBalance.UserBalance.Cmp(btsBalance.UserBalance) != 0 {
		// 	return fmt.Errorf("BTS Balance should have been same since txn does not succeed. Init %v  Final %v", btsBalance.UserBalance.String(), finalBtsBalance.UserBalance.String())
		// }
		return nil, nil
	},
}

var TransferAllBTSBalance Script = Script{
	Name:        "TransferAllBTSBalance",
	Type:        "Flow",
	Description: "Transfer Native Tokens, which are of fixed supply, in different amounts. The Token should be native for both chains",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// coinName should be a token common on both chains
		tokenExists := false
		for _, stkn := range src.NativeTokens() {
			if stkn == coinName {
				for _, dtkn := range dst.NativeTokens() {
					if dtkn == coinName {
						tokenExists = true
						break
					}
				}
				break
			}
		}
		if !tokenExists {
			ts.logger.Errorf("Token %v does not exist on both chains %v and %v", coinName, srcChain, dstChain)
			return nil, nil
		}

		btsAddr, ok := ts.btsAddressPerChain[dstChain]
		if !ok {
			return nil, errors.Wrapf(err, "dst.GetBTPAddressOfBTS %v", err)
		}
		btsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "dst.getCoinBalance %v", err)
		}

		// prepare accounts
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		amt := ts.withFeeAdded(btsBalance.UserBalance)
		if err := ts.Fund(srcAddr, amt, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund %v", err)
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// approve
		if approveHash, err := src.Approve(coinName, srcKey, amt); err != nil {
			return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		} else {
			if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
				return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
			}
		}
		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, amt)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer %v", err)
		}
		if err := ts.ValidateTransactionResultAndEvents(ctx, hash, []string{coinName}, srcAddr, dstAddr, []*big.Int{amt}); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferReceived: nil,
			chain.TransferEnd: func(e *evt) error {
				endEvt, ok := e.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent. Got %T", e.msg.EventLog)
				}
				if endEvt.Code.String() == "0" {
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		// finalBtsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "dst.getCoinBalance %v", err)
		// }
		// if finalBtsBalance.UserBalance.Cmp(btsBalance.UserBalance) != 0 {
		// 	return fmt.Errorf("BTS Balance should have been same since txn does not succeed. Init %v  Final %v", btsBalance.UserBalance.String(), finalBtsBalance.UserBalance.String())
		// }
		return nil, nil
	},
}

var TransferWithoutApprove Script = Script{
	Name:        "TransferWithoutApprove",
	Type:        "Flow",
	Description: "Transfer without approving tokens",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		if len(coinNames) == 1 && coinNames[0] == src.NativeCoin() {
			ts.logger.Info("Test valid where there is at least one token") // to not approve
			return nil, nil
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// how much you want receiving end to get

		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			existingBalance, err := src.GetCoinBalance(coinNames[i], srcAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance %v", err)
			}
			amts[i] = ts.withFeeAdded(existingBalance.UsableBalance.Add(existingBalance.UsableBalance, big.NewInt(MINIMUM_BALANCE)))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}
		// Skipping approve
		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}
		if _, err = ts.ValidateTransactionResult(ctx, hash); err != nil {
			if err.Error() == StatusCodeZero.Error() { // Failed as expected
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
		} else {
			err = errors.New("Expected event to fail but it did not ")
		}
		return nil, err
	},
}

var TransferWithApprove Script = Script{
	Name:        "TransferWithApprove",
	Type:        "Flow",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = ts.withFeeAdded(big.NewInt(MINIMUM_BALANCE))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}
		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}

		if err := ts.ValidateTransactionResultAndEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferReceived: nil,
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent. Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() == "0" {
					ts.logger.Info("Got Transfer End")
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}

		// finalSrcBalance, err := src.GetCoinBalance(coinNames[0], srcAddr)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		// }
		// finalDstBalance, err := dst.GetCoinBalance(coinNames[0], dstAddr)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		// }
		// ts.logger.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
		// ts.logger.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
		// if finalDstBalance.Usable.Cmp(initDstBalance.Usable) != 1 {
		// 	return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
		// }
		// if finalSrcBalance.Usable.Cmp(initSrcBalance.Usable) != -1 {
		// 	return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
		// }
		return nil, nil
	},
}

var TransferLessThanFee Script = Script{
	Name:        "TransferLessThanFee",
	Type:        "Flow",
	Description: "Transfer less than fee charged by BTS",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// how much you want receiving end to get
		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = ts.withFeeAdded(big.NewInt(-1))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// Skipping approve
		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}

		if _, err = ts.ValidateTransactionResult(ctx, hash); err != nil {
			if err.Error() == StatusCodeZero.Error() { // Failed as expected
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		} else {
			err = errors.New("Expected event to fail but it did not")
		}
		return nil, err
	},
}

var TransferEqualToFee Script = Script{
	Name:        "TransferEqualToFee",
	Type:        "Flow",
	Description: "Transfer equal to fee charged by BTS",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, errors.New("Should specify at least one coinname, got zero")
		}
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// how much you want receiving end to get
		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = ts.withFeeAdded(big.NewInt(0))
			if err := ts.Fund(srcAddr, amts[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund %v", err)
			}
		}

		// how much is necessary as gas cost
		if err := ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "AddGasFee %v", err)
		}

		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return nil, errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}

		if err = ts.ValidateTransactionResultAndEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Unexpected error %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "1" {
					return fmt.Errorf("Expected error code (1) Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		return nil, err
	},
}
