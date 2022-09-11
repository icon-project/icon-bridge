package executor

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

/*

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

*/

var TransferBiDirection Script = Script{
	Name:        "TransferBiDirectionWithApprove",
	Type:        "Flow",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf("Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, tmpOk := ts.clsPerChain[srcChain]
		if !tmpOk {
			return nil, fmt.Errorf("Chain %v not found", srcChain)
		}
		dst, tmpOk := ts.clsPerChain[dstChain]
		if !tmpOk {
			return nil, fmt.Errorf("Chain %v not found", srcChain)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		fmt.Println("Accounts ", srcAddr, dstAddr)
		// How much tokens do we need on src and dst accounts ?
		tokenAmountAfterFeeChargeOnDst := big.NewInt(1)
		tokenAmountBeforeFeeChargeOnDst, err := ts.getAmountBeforeFeeCharge(dstChain, coinName, tokenAmountAfterFeeChargeOnDst)
		if err != nil {
			return nil, errors.Wrapf(err, "getAmountBeforeFeeCharge %v", err)
		}
		tokenAmountBeforeFeeChargeOnSrc, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, tokenAmountBeforeFeeChargeOnDst)
		if err != nil {
			return nil, errors.Wrapf(err, "getAmountBeforeFeeCharge %v", err)
		}
		fmt.Println("Tokens ", tokenAmountAfterFeeChargeOnDst, tokenAmountBeforeFeeChargeOnDst, tokenAmountBeforeFeeChargeOnSrc)
		// How much native coins do we need to cover gas fee ?
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)
		gasLimitOnDst := big.NewInt(int64(ts.cfgPerChain[dstChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != dst.NativeCoin() {
			gasLimitOnDst.Add(gasLimitOnDst, big.NewInt(int64(ts.cfgPerChain[dstChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnDst := (&big.Int{}).Mul(dst.SuggestGasPrice(), gasLimitOnDst)
		fmt.Println("Gas ", gasFeeOnDst, gasFeeOnSrc)
		// Fund source side with the required tokens
		// These tokens should be something a newly deployed god address can have
		if err := ts.Fund(srcChain, srcAddr, tokenAmountBeforeFeeChargeOnSrc, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(dstChain, dstAddr, gasFeeOnDst, dst.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		fmt.Println("Init ", initSrcBalance, tokenAmountBeforeFeeChargeOnSrc)
		fmt.Println("Init ", initDstBalance)
		// ApproveToken On Source
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, tokenAmountBeforeFeeChargeOnSrc); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, tokenAmountBeforeFeeChargeOnSrc)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if err := ts.ValidateTransactionResultAndEvents(ctx, srcChain, transferHashOnSrc, []string{coinName}, srcAddr, dstAddr, []*big.Int{tokenAmountBeforeFeeChargeOnSrc}); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}

		// Wait For Events
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				startEvt.From = src.GetBTPAddress(startEvt.From)
				if startEvt.From != srcAddr {
					return fmt.Errorf("Expected Same Value for SrcAddr; Got startEvt.From: %v srcAddr: %v", startEvt.From, srcAddr)
				}
				if startEvt.To != dstAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got startEvt.To: %v dstAddr: %v", startEvt.To, dstAddr)
				}
				if len(startEvt.Assets) != 1 {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(startEvt.Assets))
				}
				if startEvt.Assets[0].Name != coinName {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Assets.Name: %v coinName %v", startEvt.Assets[0].Name, coinName)
				}
				if startEvt.Assets[0].Value.Cmp(tokenAmountBeforeFeeChargeOnDst) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", startEvt.Assets[0].Value, tokenAmountBeforeFeeChargeOnDst)
				}
				sum := (&big.Int{}).Add(startEvt.Assets[0].Value, startEvt.Assets[0].Fee)
				if sum.Cmp(tokenAmountBeforeFeeChargeOnSrc) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value+Fee: %v AmountBeforeFeeCharge %v", sum, tokenAmountBeforeFeeChargeOnSrc)
				}
				return nil
			},
			chain.TransferReceived: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				receivedEvt, ok := ev.msg.EventLog.(*chain.TransferReceivedEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferReceivedEvent. Got %T", ev.msg.EventLog)
				}
				receivedEvt.To = dst.GetBTPAddress(receivedEvt.To)
				if receivedEvt.To != dstAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got receivedEvt.To: %v dstAddr: %v", receivedEvt.To, dstAddr)
				}
				if len(receivedEvt.Assets) != 1 {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(receivedEvt.Assets))
				}
				if receivedEvt.Assets[0].Name != coinName {
					return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Assets.Name: %v coinName %v", receivedEvt.Assets[0].Name, coinName)
				}
				if receivedEvt.Assets[0].Value.Cmp(tokenAmountBeforeFeeChargeOnDst) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Value: %v AmountAfterFeeCharge %v", receivedEvt.Assets[0].Value, tokenAmountBeforeFeeChargeOnDst)
				}
				return nil
			},
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

		// Intermediate Tally
		intermediateSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		intermediateSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		intermediateDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			return nil, errors.Wrapf(err, "ChargedGasFee For Transfer Err: %v", err)
		}
		if coinName != src.NativeCoin() {
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				return nil, errors.Wrapf(err, "ChargedGasFee for Approve Err: %v", err)
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, intermediateSrcBalance.UserBalance)
			if tokenAmountBeforeFeeChargeOnSrc.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountBeforeFeeChargeOnSrc, tmpDiff)
			}
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, intermediateSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for src nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
			}
		} else {
			tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnSrc, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, intermediateSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
			}
		}
		tmpDiff := (&big.Int{}).Sub(intermediateDstBalance.UserBalance, initDstBalance.UserBalance)
		if tokenAmountBeforeFeeChargeOnDst.Cmp(tmpDiff) != 0 {
			return nil, fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst, tmpDiff)
		}

		// Start Transfer On Opposite Direction
		fmt.Println("IntermediateSrc ", intermediateSrcBalance)
		fmt.Println("IntermediateDst ", intermediateDstBalance, tokenAmountBeforeFeeChargeOnDst)

		// ApproveToken On Destination
		if coinName != dst.NativeCoin() {
			if approveHash, err = dst.Approve(coinName, dstKey, tokenAmountBeforeFeeChargeOnDst); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, dstChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer On Destination
		transferHashOnDst, err := dst.Transfer(coinName, dstKey, srcAddr, tokenAmountBeforeFeeChargeOnDst)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnDst ", transferHashOnDst)
		if err := ts.ValidateTransactionResultAndEvents(ctx, dstChain, transferHashOnDst, []string{coinName}, dstAddr, srcAddr, []*big.Int{tokenAmountBeforeFeeChargeOnDst}); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}

		// Wait For Events
		err = ts.WaitForEvents(ctx, dstChain, srcChain, transferHashOnDst, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				startEvt.From = dst.GetBTPAddress(startEvt.From)
				if startEvt.From != dstAddr {
					return fmt.Errorf("Expected Same Value for SrcAddr; Got startEvt.From: %v srcAddr: %v", startEvt.From, dstAddr)
				}
				if startEvt.To != srcAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got startEvt.To: %v dstAddr: %v", startEvt.To, srcAddr)
				}
				if len(startEvt.Assets) != 1 {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(startEvt.Assets))
				}
				if startEvt.Assets[0].Name != coinName {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Assets.Name: %v coinName %v", startEvt.Assets[0].Name, coinName)
				}
				if startEvt.Assets[0].Value.Cmp(tokenAmountAfterFeeChargeOnDst) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", startEvt.Assets[0].Value, tokenAmountAfterFeeChargeOnDst)
				}
				sum := (&big.Int{}).Add(startEvt.Assets[0].Value, startEvt.Assets[0].Fee)
				if sum.Cmp(tokenAmountBeforeFeeChargeOnDst) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value+Fee: %v AmountBeforeFeeCharge %v", sum, tokenAmountBeforeFeeChargeOnDst)
				}
				return nil
			},
			chain.TransferReceived: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				receivedEvt, ok := ev.msg.EventLog.(*chain.TransferReceivedEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferReceivedEvent. Got %T", ev.msg.EventLog)
				}
				receivedEvt.To = src.GetBTPAddress(receivedEvt.To)
				if receivedEvt.To != srcAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got receivedEvt.To: %v dstAddr: %v", receivedEvt.To, srcAddr)
				}
				if len(receivedEvt.Assets) != 1 {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(receivedEvt.Assets))
				}
				if receivedEvt.Assets[0].Name != coinName {
					return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Assets.Name: %v coinName %v", receivedEvt.Assets[0].Name, coinName)
				}
				if receivedEvt.Assets[0].Value.Cmp(tokenAmountAfterFeeChargeOnDst) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Value: %v AmountAfterFeeCharge %v", receivedEvt.Assets[0].Value, tokenAmountAfterFeeChargeOnDst)
				}
				return nil
			},
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

		// Final Tally
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		gasSpentOnTxn, err = dst.ChargedGasFee(transferHashOnDst)
		if err != nil {
			return nil, errors.Wrapf(err, "GetGasUsed Err: %v", err)
		}
		if coinName != src.NativeCoin() {
			gasSpentOnApprove, err := dst.ChargedGasFee(approveHash)
			if err != nil {
				return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff = (&big.Int{}).Sub(intermediateDstBalance.UserBalance, finalDstBalance.UserBalance)
			if tokenAmountBeforeFeeChargeOnDst.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst, tmpDiff)
			}
			tmpDiff = (&big.Int{}).Sub(initDstNativeCoinBalance.UserBalance, finalDstNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for dst nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
			}
		} else {
			tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnDst, gasSpentOnTxn)
			tmpDiff = (&big.Int{}).Sub(intermediateDstBalance.UserBalance, finalDstBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
			}
		}
		tmpDiff = (&big.Int{}).Sub(finalSrcBalance.UserBalance, intermediateSrcBalance.UserBalance)
		if tokenAmountAfterFeeChargeOnDst.Cmp(tmpDiff) != 0 {
			return nil, fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountAfterFeeChargeOnDst, tmpDiff)
		}
		fmt.Println("Pass")
		return nil, nil
	},
}

var TransferToZeroAddress Script = Script{
	Name:        "TransferToZeroAddress",
	Description: "Transfer to zero address",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf("Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, tmpDstAddr, err := ts.GetKeyPairs(dstChain)
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
		dstAddr := convertToZeroAddress(tmpDstAddr)

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if err := ts.ValidateTransactionResultAndEvents(ctx, srcChain, transferHashOnSrc, []string{coinName}, srcAddr, dstAddr, []*big.Int{userSuppliedAmount}); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}

		// Wait For Events
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				startEvt.From = src.GetBTPAddress(startEvt.From)
				if startEvt.From != srcAddr {
					return fmt.Errorf("Expected Same Value for SrcAddr; Got startEvt.From: %v srcAddr: %v", startEvt.From, srcAddr)
				}
				if startEvt.To != dstAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got startEvt.To: %v dstAddr: %v", startEvt.To, dstAddr)
				}
				if len(startEvt.Assets) != 1 {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(startEvt.Assets))
				}
				if startEvt.Assets[0].Name != coinName {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Assets.Name: %v coinName %v", startEvt.Assets[0].Name, coinName)
				}
				if startEvt.Assets[0].Value.Cmp(netTransferrableAmount) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", startEvt.Assets[0].Value, netTransferrableAmount)
				}
				sum := (&big.Int{}).Add(startEvt.Assets[0].Value, startEvt.Assets[0].Fee)
				if sum.Cmp(userSuppliedAmount) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value+Fee: %v AmountBeforeFeeCharge %v", sum, userSuppliedAmount)
				}
				return nil
			},
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

		// Final Tally
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			return nil, errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
		}
		if coinName != src.NativeCoin() {
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			if tmpDiff.Cmp(feeCharged) != 0 {
				return nil, fmt.Errorf("Expected same value; Got different feeCharged %v BalanceDiff %v", feeCharged, tmpDiff)
			}
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
			}
		} else {
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			tmpNativeCoinUsed := (&big.Int{}).Add(feeCharged, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same, Got Different. NativeCoinUsed %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
			}
		}
		if initDstBalance.UserBalance.Cmp(finalDstBalance.UserBalance) != 0 {
			return nil, fmt.Errorf("Epected same; Got Different. initDstBalance %v finalDstBalance %v", initDstBalance, finalDstBalance)
		}
		fmt.Println("pass", err)
		return nil, err
	},
}

var TransferToUnknownNetwork Script = Script{
	Name:        "TransferToUnknownNetwork",
	Description: "Transfer to unknow network",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf("Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, tmpDstAddr, err := ts.GetKeyPairs(dstChain)
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
		dstAddr, err := changeBMCNetwork(tmpDstAddr)

		// Funds
		netTransferrableAmount := big.NewInt(-1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					return nil, errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
				}
				if coinName != src.NativeCoin() {
					if initSrcBalance.TotalBalance.Cmp(finalSrcBalance.TotalBalance) != 0 {
						return nil, fmt.Errorf("Expected Same, Got Different. finalSrcBalance %v initialSrcBalance %v", finalSrcBalance.TotalBalance, initSrcBalance.TotalBalance)
					}
					gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
					if err != nil {
						return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
					}
					gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					}
				} else {
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					}
				}
				fmt.Println("Pass 1")
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
		}
		return nil, err
	},
}

var TransferWithoutApprove Script = Script{
	Name:        "TransferWithoutApprove",
	Description: "Transfer Without Approve",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf(" Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]

		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		if coinName == src.NativeCoin() {
			ts.logger.Warnf("Expected non-native coin; Got native coin %v", src.NativeCoin())
			return nil, nil // not returning an error here
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)

		if _, err = ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			if err.Error() == StatusCodeZero.Error() { // Failed as expected
				// Final Tally
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					return nil, errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
				}
				if initSrcBalance.UserBalance.Cmp(finalSrcBalance.UserBalance) != 0 {
					return nil, fmt.Errorf("Expected same value; Got different initSrcBalance %v finalSrcbalance %v", initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
				}
				tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
				if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
				}
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
		}
		return nil, errors.New("Expected event to fail but it did not ")
	},
}

var TransferLessThanFee Script = Script{
	Name:        "TransferLessThanFee",
	Description: "Transfer to unknow network",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf("Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// Funds
		netTransferrableAmount := big.NewInt(-1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					return nil, errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
				}
				if coinName != src.NativeCoin() {
					if initSrcBalance.TotalBalance.Cmp(finalSrcBalance.TotalBalance) != 0 {
						return nil, fmt.Errorf("Expected Same, Got Different. finalSrcBalance %v initialSrcBalance %v", finalSrcBalance.TotalBalance, initSrcBalance.TotalBalance)
					}
					gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
					if err != nil {
						return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
					}
					gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					}
				} else {
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					}
				}
				fmt.Println("Pass 1")
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
		}
		return nil, err
	},
}

var TransferEqualToFee Script = Script{
	Name:        "TransferEqualToFee",
	Description: "Transfer equal to fee",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf("Should specify a single coinName, got %v", len(coinNames))
		}

		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// Funds
		netTransferrableAmount := big.NewInt(0)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			if err.Error() == StatusCodeZero.Error() {
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					return nil, errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
				}
				if coinName != src.NativeCoin() {
					if initSrcBalance.TotalBalance.Cmp(finalSrcBalance.TotalBalance) != 0 {
						return nil, fmt.Errorf("Expected Same, Got Different. finalSrcBalance %v initialSrcBalance %v", finalSrcBalance.TotalBalance, initSrcBalance.TotalBalance)
					}
					gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
					if err != nil {
						return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
					}
					gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					}
				} else {
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					}
				}
				fmt.Println("Pass 1")
				return nil, nil
			}
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
		}
		// Wait For Events
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				startEvt.From = src.GetBTPAddress(startEvt.From)
				if startEvt.From != srcAddr {
					return fmt.Errorf("Expected Same Value for SrcAddr; Got startEvt.From: %v srcAddr: %v", startEvt.From, srcAddr)
				}
				if startEvt.To != dstAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got startEvt.To: %v dstAddr: %v", startEvt.To, dstAddr)
				}
				if len(startEvt.Assets) != 1 {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(startEvt.Assets))
				}
				if startEvt.Assets[0].Name != coinName {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Assets.Name: %v coinName %v", startEvt.Assets[0].Name, coinName)
				}
				if startEvt.Assets[0].Value.Cmp(netTransferrableAmount) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", startEvt.Assets[0].Value, netTransferrableAmount)
				}
				sum := (&big.Int{}).Add(startEvt.Assets[0].Value, startEvt.Assets[0].Fee)
				if sum.Cmp(userSuppliedAmount) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value+Fee: %v AmountBeforeFeeCharge %v", sum, userSuppliedAmount)
				}
				return nil
			},
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
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			return nil, errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
		}
		if coinName != src.NativeCoin() {
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			if tmpDiff.Cmp(feeCharged) != 0 {
				return nil, fmt.Errorf("Expected same value; Got different feeCharged %v BalanceDiff %v", feeCharged, tmpDiff)
			}
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
			}
		} else {
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			tmpNativeCoinUsed := (&big.Int{}).Add(feeCharged, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same, Got Different. NativeCoinUsed %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
			}
		}
		if initDstBalance.UserBalance.Cmp(finalDstBalance.UserBalance) != 0 {
			return nil, fmt.Errorf("Epected same; Got Different. initDstBalance %v finalDstBalance %v", initDstBalance, finalDstBalance)
		}
		fmt.Println("pass 2")
		return nil, err
	},
}

var TransferToBlackListedDstAddress Script = Script{
	Name:        "TransferToBlackListedDstAddress",
	Description: "Transfer to BlackListed Destination Address",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf(" Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}

		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}

		fmt.Println("Start Add to blacklist")

		// Add To BlackList
		fCfg, err := ts.GetFullConfigAPI()
		if err != nil {
			return nil, errors.Wrapf(err, "GetFullConfigAPI %v", err)
		}
		fCfgOwnerKey := ts.FullConfigAPIsOwner()
		stdCfg, err := ts.GetStandardConfigAPI(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetStandardConfigAPI %v", err)
		}
		blDstNet, blDstAddr := ts.NetAddr(dstAddr)
		blackListAddHash, err := fCfg.AddBlackListAddress(fCfgOwnerKey, blDstNet, []string{blDstAddr})
		if err != nil {
			return nil, errors.Wrapf(err, "AddBlackListAddress %v", err)
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListAddHash); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		isBlackListed, err := fCfg.IsUserBlackListed(blDstNet, blDstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "isBlackListed %v", err)
		} else if err == nil && !isBlackListed {
			return nil, fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blDstNet, blDstAddr)
		}
		if dstChain != ts.FullConfigAPIChain() { // for interchain-configurations
			err = ts.WaitForConfigResponse(ctx, chain.AddToBlacklistRequest, chain.BlacklistResponse, dstChain, blackListAddHash,
				map[chain.EventLogType]func(event *evt) error{
					chain.AddToBlacklistRequest: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						reqEvt, ok := ev.msg.EventLog.(*chain.AddToBlacklistRequestEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.AddToBlacklistRequestEvent. Got %T", ev.msg.EventLog)
						}
						if reqEvt.Net != blDstNet {
							return fmt.Errorf("Expected same; Got different reqEvt.Net %v DstNet %v", reqEvt.Net, blDstNet)
						}
						if len(reqEvt.Addrs) != 1 {
							return fmt.Errorf("Expected reqEvt.AddrsLen 1; Got %v", len(reqEvt.Addrs))
						}
						if strings.ToLower(reqEvt.Addrs[0]) != strings.ToLower(blDstAddr) {
							return fmt.Errorf("Expected same; Got Different reqEvt.Addrs %v DstAddr %v", reqEvt.Addrs[0], blDstAddr)
						}
						return nil
					},
					chain.BlacklistResponse: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						resEvt, ok := ev.msg.EventLog.(*chain.BlacklistResponseEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.BlacklistResponseEvent. Got %T", ev.msg.EventLog)
						}
						if resEvt.Code != 0 {
							return fmt.Errorf("Expected Code 0; Got Sn %v Code %v Msg %v", resEvt.Sn, resEvt.Code, resEvt.Msg)
						}
						return nil
					},
				},
			)
			if err != nil {
				return nil, errors.Wrapf(err, "WaitForConfigResponse %v", err)
			}
			isBlackListed, err = stdCfg.IsUserBlackListed(blDstNet, blDstAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "isBlackListed %v", err)
			} else if err == nil && !isBlackListed {
				return nil, fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blDstNet, blDstAddr)
			}
		}

		// Send After BlackListing
		fmt.Println("Send to blacklist")
		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}

		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err == nil {
			// Wait For Events
			err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
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
		} else {
			if err.Error() != StatusCodeZero.Error() {
				return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
			}
		}

		// Remove From BlackList
		fmt.Println("Remove From BlackList")

		blackListRemoveHash, err := fCfg.RemoveBlackListAddress(fCfgOwnerKey, blDstNet, []string{blDstAddr})
		if err != nil {
			return nil, errors.Wrapf(err, "RemoveBlackListAddress %v", err)
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListRemoveHash); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResult %v", err)
		}
		isBlackListed, err = fCfg.IsUserBlackListed(blDstNet, blDstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "isBlackListed %v", err)
		} else if err == nil && isBlackListed {
			return nil, fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blDstNet, blDstAddr)
		}
		if dstChain != ts.FullConfigAPIChain() { // for interchain-configurations
			err = ts.WaitForConfigResponse(ctx, chain.RemoveFromBlacklistRequest, chain.BlacklistResponse, dstChain, blackListRemoveHash,
				map[chain.EventLogType]func(event *evt) error{
					chain.RemoveFromBlacklistRequest: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						reqEvt, ok := ev.msg.EventLog.(*chain.RemoveFromBlacklistRequestEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.RemoveFromBlacklistRequestEvent. Got %T", ev.msg.EventLog)
						}
						if reqEvt.Net != blDstNet {
							return fmt.Errorf("Expected same; Got different reqEvt.Net %v DstNet %v", reqEvt.Net, blDstNet)
						}
						if len(reqEvt.Addrs) != 1 {
							return fmt.Errorf("Expected reqEvt.AddrsLen 1; Got %v", len(reqEvt.Addrs))
						}
						if strings.ToLower(reqEvt.Addrs[0]) != strings.ToLower(blDstAddr) {
							return fmt.Errorf("Expected same; Got Different reqEvt.Addrs %v DstAddr %v", reqEvt.Addrs[0], dstAddr)
						}
						return nil
					},
					chain.BlacklistResponse: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						resEvt, ok := ev.msg.EventLog.(*chain.BlacklistResponseEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.BlacklistResponseEvent. Got %T", ev.msg.EventLog)
						}
						if resEvt.Code != 0 {
							return fmt.Errorf("Expected Code 0; Got Sn %v Code %v Msg %v", resEvt.Sn, resEvt.Code, resEvt.Msg)
						}
						return nil
					},
				},
			)
			if err != nil {
				return nil, errors.Wrapf(err, "WaitForConfigResponse %v", err)
			}
			isBlackListed, err = stdCfg.IsUserBlackListed(blDstNet, blDstAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "isBlackListed %v", err)
			} else if err == nil && isBlackListed {
				return nil, fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blDstNet, blDstAddr)
			}
		}

		fmt.Println("Final Send Should Succeed")
		// Final Send Should Succeed
		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}

		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		fmt.Println("pass")
		return nil, nil
	},
}

var TransferFromBlackListedSrcAddress Script = Script{
	Name:        "TransferFromBlackListedSrcAddress",
	Description: "Transfer from BlackListed Source Address",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf(" Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)

		if err := ts.Fund(srcChain, srcAddr, (&big.Int{}).Mul(userSuppliedAmount, big.NewInt(2)), coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, (&big.Int{}).Mul(gasFeeOnSrc, big.NewInt(2)), src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}

		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}

		fmt.Println("Start Add to blacklist")

		// Add To BlackList
		fCfg, err := ts.GetFullConfigAPI()
		if err != nil {
			return nil, errors.Wrapf(err, "GetFullConfigAPI %v", err)
		}
		fCfgOwnerKey := ts.FullConfigAPIsOwner()
		stdCfg, err := ts.GetStandardConfigAPI(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetStandardConfigAPI %v", err)
		}
		blsSrcNet, blsSrcAddr := ts.NetAddr(srcAddr)
		blackListAddHash, err := fCfg.AddBlackListAddress(fCfgOwnerKey, blsSrcNet, []string{blsSrcAddr})
		if err != nil {
			return nil, errors.Wrapf(err, "AddBlackListAddress %v", err)
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListAddHash); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		isBlackListed, err := fCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "isBlackListed %v", err)
		} else if err == nil && !isBlackListed {
			return nil, fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blsSrcNet, blsSrcAddr)
		}
		if srcChain != ts.FullConfigAPIChain() { // for interchain-configurations
			err = ts.WaitForConfigResponse(ctx, chain.AddToBlacklistRequest, chain.BlacklistResponse, srcChain, blackListAddHash,
				map[chain.EventLogType]func(event *evt) error{
					chain.AddToBlacklistRequest: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						reqEvt, ok := ev.msg.EventLog.(*chain.AddToBlacklistRequestEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.AddToBlacklistRequestEvent. Got %T", ev.msg.EventLog)
						}
						if reqEvt.Net != blsSrcNet {
							return fmt.Errorf("Expected same; Got different reqEvt.Net %v DstNet %v", reqEvt.Net, blsSrcNet)
						}
						if len(reqEvt.Addrs) != 1 {
							return fmt.Errorf("Expected reqEvt.AddrsLen 1; Got %v", len(reqEvt.Addrs))
						}
						if strings.ToLower(reqEvt.Addrs[0]) != strings.ToLower(blsSrcAddr) {
							return fmt.Errorf("Expected same; Got Different reqEvt.Addrs %v DstAddr %v", reqEvt.Addrs[0], blsSrcAddr)
						}
						return nil
					},
					chain.BlacklistResponse: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						resEvt, ok := ev.msg.EventLog.(*chain.BlacklistResponseEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.BlacklistResponseEvent. Got %T", ev.msg.EventLog)
						}
						if resEvt.Code != 0 {
							return fmt.Errorf("Expected Code 0; Got Sn %v Code %v Msg %v", resEvt.Sn, resEvt.Code, resEvt.Msg)
						}
						return nil
					},
				},
			)
			if err != nil {
				return nil, errors.Wrapf(err, "WaitForConfigResponse %v", err)
			}
			isBlackListed, err = stdCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "isBlackListed %v", err)
			} else if err == nil && !isBlackListed {
				return nil, fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blsSrcNet, blsSrcAddr)
			}
		}

		// Send After BlackListing
		fmt.Println("Send to blacklist")

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil && err.Error() != StatusCodeZero.Error() {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}

		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err == nil {
			// Wait For Events
			err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
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
		} else {
			if err.Error() != StatusCodeZero.Error() {
				return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents Got Unexpected Error: %v", err)
			}
		}

		// Remove From BlackList
		fmt.Println("Remove From BlackList")

		blackListRemoveHash, err := fCfg.RemoveBlackListAddress(fCfgOwnerKey, blsSrcNet, []string{blsSrcAddr})
		if err != nil {
			return nil, errors.Wrapf(err, "RemoveBlackListAddress %v", err)
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListRemoveHash); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResult %v", err)
		}
		isBlackListed, err = fCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "isBlackListed %v", err)
		} else if err == nil && isBlackListed {
			return nil, fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blsSrcNet, blsSrcAddr)
		}
		if srcChain != ts.FullConfigAPIChain() { // for interchain-configurations
			err = ts.WaitForConfigResponse(ctx, chain.RemoveFromBlacklistRequest, chain.BlacklistResponse, dstChain, blackListRemoveHash,
				map[chain.EventLogType]func(event *evt) error{
					chain.RemoveFromBlacklistRequest: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						reqEvt, ok := ev.msg.EventLog.(*chain.RemoveFromBlacklistRequestEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.RemoveFromBlacklistRequestEvent. Got %T", ev.msg.EventLog)
						}
						if reqEvt.Net != blsSrcNet {
							return fmt.Errorf("Expected same; Got different reqEvt.Net %v DstNet %v", reqEvt.Net, blsSrcNet)
						}
						if len(reqEvt.Addrs) != 1 {
							return fmt.Errorf("Expected reqEvt.AddrsLen 1; Got %v", len(reqEvt.Addrs))
						}
						if strings.ToLower(reqEvt.Addrs[0]) != strings.ToLower(blsSrcAddr) {
							return fmt.Errorf("Expected same; Got Different reqEvt.Addrs %v DstAddr %v", reqEvt.Addrs[0], dstAddr)
						}
						return nil
					},
					chain.BlacklistResponse: func(ev *evt) error {
						if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
							return errors.New("Got nil value for event ")
						}
						resEvt, ok := ev.msg.EventLog.(*chain.BlacklistResponseEvent)
						if !ok {
							return fmt.Errorf("Expected *chain.BlacklistResponseEvent. Got %T", ev.msg.EventLog)
						}
						if resEvt.Code != 0 {
							return fmt.Errorf("Expected Code 0; Got Sn %v Code %v Msg %v", resEvt.Sn, resEvt.Code, resEvt.Msg)
						}
						return nil
					},
				},
			)
			if err != nil {
				return nil, errors.Wrapf(err, "WaitForConfigResponse %v", err)
			}
			isBlackListed, err = stdCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "isBlackListed %v", err)
			} else if err == nil && isBlackListed {
				return nil, fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blsSrcNet, blsSrcAddr)
			}
		}

		fmt.Println("Final Send Should Succeed")
		// Final Send Should Succeed
		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}

		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "WaitForEvents %v", err)
		}
		fmt.Println("pass")
		return nil, nil
	},
}

var ChangeFee Script = Script{
	Name:        "ChangeFee",
	Description: "Change Fee",
	Type:        "Configure",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) != 1 {
			return nil, fmt.Errorf(" Should specify a single coinName, got %v", len(coinNames))
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetChainPair %v", err)
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// Fund
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, big.NewInt(1))
		if err != nil {
			return nil, errors.Wrapf(err, "getAmountBeforeFeeCharge %v", err)
		}
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)
		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}
		// Transfer Amount greater than tokenLimit
		transferHash, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHash); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResult Err: %v Hash %v", err, approveHash)
		}
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHash, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
				}
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})

		// Update Fee
		stdCfg, err := ts.GetStandardConfigAPI(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetStandardConfigAPI %v", err)
		}
		feeNumerator, fixedFee, err := stdCfg.GetFeeRatio(coinName)
		if err != nil {
			return nil, errors.Wrapf(err, "GetFeeRatio %v", err)
		}
		feeNumerator.Sub(feeNumerator, big.NewInt(2))
		_, err = stdCfg.SetFeeRatio(ts.GetStandardConfigAPIOwnerKey(srcChain), coinName, feeNumerator, fixedFee)
		if err != nil {
			return nil, errors.Wrapf(err, "SetFeeRatio %v", err)
		}
		time.Sleep(time.Second * 5) // TODO: Wait for status true on txn hash
		newFeeNumerator, _, err := stdCfg.GetFeeRatio(coinName)
		if err != nil {
			return nil, errors.Wrapf(err, "GetFeeRatio %v", err)
		}
		if newFeeNumerator.Cmp(feeNumerator) != 0 {
			return nil, fmt.Errorf("Expected same. Got newFeeNumerator %v feeNumerator %v", newFeeNumerator, feeNumerator)
		}

		//Fund
		if err := ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		// Transfer Again
		transferHash, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), transferHash); err.Error() != StatusCodeZero.Error() {
			return nil, errors.Wrapf(err, "Expected zero code Got Err %v", err)
		}

		fmt.Println("Pass")
		return nil, nil
	},
}

var TransferBatchBiDirection Script = Script{
	Name:        "TransferBatchBiDirection",
	Description: "Transfer batch bi-direction",
	Type:        "Flow",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (*txnRecord, error) {
		if len(coinNames) == 0 {
			return nil, fmt.Errorf(" Should specify at least one coinName, got zero")
		}
		src, tmpOk := ts.clsPerChain[srcChain]
		if !tmpOk {
			return nil, fmt.Errorf("Chain %v not found", srcChain)
		}
		dst, tmpOk := ts.clsPerChain[dstChain]
		if !tmpOk {
			return nil, fmt.Errorf("Chain %v not found", srcChain)
		}
		if len(coinNames) == 1 && coinNames[0] == src.NativeCoin() {
			return nil, fmt.Errorf("A single src.NativeCoin %v has been used", coinNames[0])
		}
		if len(coinNames) == 1 && coinNames[0] == dst.NativeCoin() {
			return nil, fmt.Errorf("A single dst.NativeCoin %v has been used", coinNames[0])
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		fmt.Println("Calculate Transfer Amounts")
		// Calculate Transfer Amounts
		tokenAmountAfterFeeChargeOnDst := big.NewInt(1)
		tokenAmountBeforeFeeChargeOnDst := make([]*big.Int, len(coinNames))
		nativeCoinAmountBeforeFeeChargeOnDst := big.NewInt(0)
		for i := 0; i < len(coinNames); i++ {
			tokenAmountBeforeFeeChargeOnDst[i], err = ts.getAmountBeforeFeeCharge(dstChain, coinNames[i], tokenAmountAfterFeeChargeOnDst)
			if err != nil {
				return nil, errors.Wrapf(err, "getAmountBeforeFeeCharge(%v) %v", coinNames[i], err)
			}
			if coinNames[i] == dst.NativeCoin() {
				nativeCoinAmountBeforeFeeChargeOnDst.Set(tokenAmountBeforeFeeChargeOnDst[i])
			}
		}
		tokenAmountBeforeFeeChargeOnSrc := make([]*big.Int, len(coinNames))
		nativeCoinAmountBeforeFeeChargeOnSrc := big.NewInt(0)
		for i := 0; i < len(coinNames); i++ {
			tokenAmountBeforeFeeChargeOnSrc[i], err = ts.getAmountBeforeFeeCharge(dstChain, coinNames[i], tokenAmountBeforeFeeChargeOnDst[i])
			if err != nil {
				return nil, errors.Wrapf(err, "getAmountBeforeFeeCharge(%v) %v", coinNames[i], err)
			}
			if coinNames[i] == src.NativeCoin() {
				nativeCoinAmountBeforeFeeChargeOnSrc.Set(tokenAmountBeforeFeeChargeOnSrc[i])
			}
		}
		fmt.Println("Calculate Gas Fees")
		// Calculate Gas Fees
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferBatchCoinInterChainGasLimit]))
		for _, coinName := range coinNames {
			if coinName == src.NativeCoin() {
				continue
			}
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)
		gasLimitOnDst := big.NewInt(int64(ts.cfgPerChain[dstChain].GasLimit[chain.TransferBatchCoinInterChainGasLimit]))
		for _, coinName := range coinNames {
			if coinName == src.NativeCoin() {
				continue
			}
			gasLimitOnDst.Add(gasLimitOnDst, big.NewInt(int64(ts.cfgPerChain[dstChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}

		gasFeeOnDst := (&big.Int{}).Mul(dst.SuggestGasPrice(), gasLimitOnDst)
		gasFeeOnDst.Add(gasFeeOnDst, big.NewInt(500000000000000))
		fmt.Println("Fund")
		//Fund
		for i := 0; i < len(coinNames); i++ {
			if err := ts.Fund(srcChain, srcAddr, tokenAmountBeforeFeeChargeOnSrc[i], coinNames[i]); err != nil {
				return nil, errors.Wrapf(err, "Fund Token %v", err)
			}
		}
		if err := ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}
		if err := ts.Fund(dstChain, dstAddr, gasFeeOnDst, dst.NativeCoin()); err != nil {
			return nil, errors.Wrapf(err, "Fund Token %v", err)
		}

		fmt.Println("Record Initial Balance")
		// Record Initial Balance
		initSrcBalance := make([]*chain.CoinBalance, len(coinNames))
		initDstBalance := make([]*chain.CoinBalance, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			initSrcBalance[i], err = src.GetCoinBalance(coinNames[i], srcAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			}
			initDstBalance[i], err = dst.GetCoinBalance(coinNames[i], dstAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			}
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		fmt.Println("ApproveToken On Source")
		// ApproveToken On Source
		approveHash := map[string]string{}
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash[coinName], err = src.Approve(coinName, srcKey, tokenAmountBeforeFeeChargeOnSrc[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash[coinName])
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash[coinName]); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash[coinName])
					}
				}
			}
		}

		fmt.Println("Transfer On Source")
		// Transfer On Source
		transferHashOnSrc, err := src.TransferBatch(coinNames, srcKey, dstAddr, tokenAmountBeforeFeeChargeOnSrc)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnSrc ", transferHashOnSrc)
		if err := ts.ValidateTransactionResultAndEvents(ctx, srcChain, transferHashOnSrc, coinNames, srcAddr, dstAddr, tokenAmountBeforeFeeChargeOnSrc); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}

		fmt.Println("Wait For Events")
		// Wait For Events
		err = ts.WaitForEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				startEvt.From = src.GetBTPAddress(startEvt.From)
				if startEvt.From != srcAddr {
					return fmt.Errorf("Expected Same Value for SrcAddr; Got startEvt.From: %v srcAddr: %v", startEvt.From, srcAddr)
				}
				if startEvt.To != dstAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got startEvt.To: %v dstAddr: %v", startEvt.To, dstAddr)
				}
				if len(startEvt.Assets) != len(coinNames) {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(startEvt.Assets))
				}
				for _, assets := range startEvt.Assets {
					index := -1
					for i, coinName := range coinNames {
						if coinName == assets.Name {
							index = i
							break
						}
					}
					if index == -1 {
						return fmt.Errorf("Asset name %v not on coinNames list", assets.Name)
					}

					if assets.Value.Cmp(tokenAmountBeforeFeeChargeOnDst[index]) != 0 {
						return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", assets.Value, tokenAmountBeforeFeeChargeOnDst[index])
					}
					sum := (&big.Int{}).Add(assets.Value, assets.Fee)
					if sum.Cmp(tokenAmountBeforeFeeChargeOnSrc[index]) != 0 {
						return fmt.Errorf("Expected same value for coinName; Got startEvt.Value+Fee: %v AmountBeforeFeeCharge %v", sum, tokenAmountBeforeFeeChargeOnSrc[index])
					}
				}
				return nil
			},
			chain.TransferReceived: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				receivedEvt, ok := ev.msg.EventLog.(*chain.TransferReceivedEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferReceivedEvent. Got %T", ev.msg.EventLog)
				}
				receivedEvt.To = dst.GetBTPAddress(receivedEvt.To)
				if receivedEvt.To != dstAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got receivedEvt.To: %v dstAddr: %v", receivedEvt.To, dstAddr)
				}
				if len(receivedEvt.Assets) != len(coinNames) {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(receivedEvt.Assets))
				}
				for _, assets := range receivedEvt.Assets {
					index := -1
					for i, coinName := range coinNames {
						if coinName == assets.Name {
							index = i
							break
						}
					}
					if index == -1 {
						return fmt.Errorf("Asset name %v not on coinNames list", assets.Name)
					}
					if assets.Value.Cmp(tokenAmountBeforeFeeChargeOnDst[index]) != 0 {
						return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Value: %v AmountAfterFeeCharge %v", assets.Value, tokenAmountBeforeFeeChargeOnDst[index])
					}
				}
				return nil
			},
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
		fmt.Println("Intermediate Tally")
		// Intermediate Tally
		intermediateSrcBalance := make([]*chain.CoinBalance, len(coinNames))
		intermediateDstBalance := make([]*chain.CoinBalance, len(coinNames))
		for i, coinName := range coinNames {
			intermediateSrcBalance[i], err = src.GetCoinBalance(coinName, srcAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			}
			intermediateDstBalance[i], err = dst.GetCoinBalance(coinName, dstAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			}
		}
		intermediateSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			return nil, errors.Wrapf(err, "ChargedGasFee For Transfer Err: %v", err)
		}
		for _, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				gasSpentOnApprove, err := src.ChargedGasFee(approveHash[coinName])
				if err != nil {
					return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
				}
				gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			}
		}
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				tmpDiff := (&big.Int{}).Sub(initSrcBalance[i].UserBalance, intermediateSrcBalance[i].UserBalance)
				if tokenAmountBeforeFeeChargeOnSrc[i].Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountBeforeFeeChargeOnSrc[i], tmpDiff)
				}
				tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, intermediateSrcNativeCoinBalance.UserBalance)
				tmpNativeCoinUsed := (&big.Int{}).Add(nativeCoinAmountBeforeFeeChargeOnSrc, gasSpentOnTxn)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value for src nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
				}
			} else {
				tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnSrc[i], gasSpentOnTxn)
				tmpDiff := (&big.Int{}).Sub(initSrcBalance[i].UserBalance, intermediateSrcBalance[i].UserBalance)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				}
			}
			tmpDiff := (&big.Int{}).Sub(intermediateDstBalance[i].UserBalance, initDstBalance[i].UserBalance)
			if tokenAmountBeforeFeeChargeOnDst[i].Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst[i], tmpDiff)
			}
		}

		fmt.Println("ApproveToken On Destination")
		// ApproveToken On Destination
		approveHash = map[string]string{}
		for i, coinName := range coinNames {
			if coinName != dst.NativeCoin() {
				if approveHash[coinName], err = dst.Approve(coinName, dstKey, tokenAmountBeforeFeeChargeOnDst[i]); err != nil {
					return nil, errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash[coinName])
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, dstChain, approveHash[coinName]); err != nil {
						return nil, errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash[coinName])
					}
				}
			}
		}

		fmt.Println("Transfer On Destination")
		// Transfer On Destination
		transferHashOnDst, err := dst.TransferBatch(coinNames, dstKey, srcAddr, tokenAmountBeforeFeeChargeOnDst)
		if err != nil {
			return nil, errors.Wrapf(err, "Transfer Err: %v", err)
		}
		fmt.Println("transferHashOnDst ", transferHashOnDst)
		if err := ts.ValidateTransactionResultAndEvents(ctx, dstChain, transferHashOnDst, coinNames, dstAddr, srcAddr, tokenAmountBeforeFeeChargeOnDst); err != nil {
			return nil, errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
		}

		fmt.Println("Wait For Events")
		// Wait For Events
		err = ts.WaitForEvents(ctx, dstChain, srcChain, transferHashOnDst, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				startEvt.From = dst.GetBTPAddress(startEvt.From)
				if startEvt.From != dstAddr {
					return fmt.Errorf("Expected Same Value for SrcAddr; Got startEvt.From: %v srcAddr: %v", startEvt.From, dstAddr)
				}
				if startEvt.To != srcAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got startEvt.To: %v dstAddr: %v", startEvt.To, srcAddr)
				}
				if len(startEvt.Assets) != len(coinNames) {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(startEvt.Assets))
				}
				for _, assets := range startEvt.Assets {
					index := -1
					for i, coinName := range coinNames {
						if coinName == assets.Name {
							index = i
							break
						}
					}
					if index == -1 {
						return fmt.Errorf("Asset name %v not on coinNames list", assets.Name)
					}
					if assets.Value.Cmp(tokenAmountAfterFeeChargeOnDst) != 0 {
						return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", assets.Value, tokenAmountAfterFeeChargeOnDst)
					}
					sum := (&big.Int{}).Add(assets.Value, assets.Fee)
					if sum.Cmp(tokenAmountBeforeFeeChargeOnDst[index]) != 0 {
						return fmt.Errorf("Expected same value for coinName; Got startEvt.Value+Fee: %v AmountBeforeFeeCharge %v", sum, tokenAmountBeforeFeeChargeOnDst[index])
					}
				}
				return nil
			},
			chain.TransferReceived: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				receivedEvt, ok := ev.msg.EventLog.(*chain.TransferReceivedEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferReceivedEvent. Got %T", ev.msg.EventLog)
				}
				receivedEvt.To = src.GetBTPAddress(receivedEvt.To)
				if receivedEvt.To != srcAddr {
					return fmt.Errorf("Expected Same Value for DstAddr; Got receivedEvt.To: %v dstAddr: %v", receivedEvt.To, srcAddr)
				}
				if len(receivedEvt.Assets) != len(coinNames) {
					return fmt.Errorf("For single token transfer; Expected single asset; Got %v", len(receivedEvt.Assets))
				}
				for _, assets := range receivedEvt.Assets {
					index := -1
					for i, coinName := range coinNames {
						if coinName == assets.Name {
							index = i
							break
						}
					}
					if index == -1 {
						return fmt.Errorf("Asset name %v not on coinNames list", assets.Name)
					}
					if assets.Value.Cmp(tokenAmountAfterFeeChargeOnDst) != 0 {
						return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Value: %v AmountAfterFeeCharge %v", receivedEvt.Assets[0].Value, tokenAmountAfterFeeChargeOnDst)
					}
				}
				return nil
			},
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

		fmt.Println("final tally")
		// Final Tally
		finalSrcBalance := make([]*chain.CoinBalance, len(coinNames))
		finalDstBalance := make([]*chain.CoinBalance, len(coinNames))
		for i, coinName := range coinNames {
			finalSrcBalance[i], err = src.GetCoinBalance(coinName, srcAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			}
			finalDstBalance[i], err = dst.GetCoinBalance(coinName, dstAddr)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			}
		}
		finalDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		gasSpentOnTxn, err = dst.ChargedGasFee(transferHashOnDst)
		if err != nil {
			return nil, errors.Wrapf(err, "GetGasUsed Err: %v", err)
		}
		for _, coinName := range coinNames {
			if coinName != dst.NativeCoin() {
				gasSpentOnApprove, err := dst.ChargedGasFee(approveHash[coinName])
				if err != nil {
					return nil, errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
				}
				gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			}
		}
		for i, coinName := range coinNames {
			if coinName != dst.NativeCoin() {
				tmpDiff := (&big.Int{}).Sub(intermediateDstBalance[i].UserBalance, finalDstBalance[i].UserBalance)
				if tokenAmountBeforeFeeChargeOnDst[i].Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst[i], tmpDiff)
				}
				tmpNativeCoinUsed := (&big.Int{}).Add(nativeCoinAmountBeforeFeeChargeOnDst, gasSpentOnTxn)
				tmpDiff = (&big.Int{}).Sub(initDstNativeCoinBalance.UserBalance, finalDstNativeCoinBalance.UserBalance)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value for dst nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
				}
			} else {
				tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnDst[i], gasSpentOnTxn)
				tmpDiff := (&big.Int{}).Sub(intermediateDstBalance[i].UserBalance, finalDstBalance[i].UserBalance)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					return nil, fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				}
			}
			tmpDiff := (&big.Int{}).Sub(finalSrcBalance[i].UserBalance, intermediateSrcBalance[i].UserBalance)
			if tokenAmountAfterFeeChargeOnDst.Cmp(tmpDiff) != 0 {
				return nil, fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountAfterFeeChargeOnDst, tmpDiff)
			}
		}
		fmt.Println("Pass")
		return nil, nil
	},
}
