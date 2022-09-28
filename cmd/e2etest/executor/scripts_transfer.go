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

/*

var TransferExceedingBTSBalance Script = Script{
	Name:        "TransferExceedingContractsBalance",
	Type:        "Transfer",
	Description: "Transfer Native Tokens, which are of fixed supply, in different amounts. The Token should be native for both chains",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
	srcChain := tp.SrcChain
	dstChain := tp.DstChain
	coinNames := tp.CoinNames

		if len(coinNames) == 0 {
			errs =errors.New("Should specify at least one coinname, got zero")
 ts.logger.Error(errs)
 return
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs =errors.Wrapf(err, "GetChainPair %v", err)
 ts.logger.Error(errs)
 return
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
			errs =nil
 ts.logger.Error(errs)
 return
		}

		btsAddr, ok := ts.btsAddressPerChain[dstChain]
		if !ok {
			errs =errors.Wrapf(err, "dst.GetBTPAddressOfBTS %v", err)
 ts.logger.Error(errs)
 return
		}
		btsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			errs =errors.Wrapf(err, "dst.getCoinBalance %v", err)
 ts.logger.Error(errs)
 return
		}

		// prepare accounts
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs =errors.Wrapf(err, "GetKeyPairs %v", err)
 ts.logger.Error(errs)
 return
		}
 txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs =errors.Wrapf(err, "GetKeyPairs %v", err)
 ts.logger.Error(errs)
 return
		}
 txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		amt := ts.withFeeAdded(btsBalance.UserBalance)
		amt.Add(amt, big.NewInt(MINIMUM_BALANCE)) //exceed balance by 100g
		if errs = ts.Fund(srcAddr, amt, coinName); errs != nil {
			errs =errors.Wrapf(err, "Fund %v", err)
 ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
 return
		}

		// how much is necessary as gas cost
		if errs = ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); errs != nil {
			errs =errors.Wrapf(err, "AddGasFee %v", err)
 ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
 return
		}

		// approve
		if approveHash, err := src.Approve(coinName, srcKey, amt); err != nil {
			errs =errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
 ts.logger.Error(errs)
 return
		} else {
			if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
				errs =errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
 ts.logger.Error(errs)
 return
			}
		}
		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, amt)
		if err != nil {
			errs =errors.Wrapf(err, "Transfer %v", err)
 ts.logger.Error(errs)
 return
		}
		if err := ts.ValidateTransactionResult(ctx, hash, []string{coinName}, srcAddr, dstAddr, []*big.Int{amt}); err != nil {
			errs =errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
 ts.logger.Error(errs)
 return
		}
		err = ts.WaitForTransferEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			// chain.TransferReceived: func(e *evt) error {
			// 	ts.logger.Debug("Got TransferReceived")
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
			errs =errors.Wrapf(err, "WaitForTransferEvents %v", err)
 ts.logger.Error(errs)
 return
		}
		// finalBtsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		// if err != nil {
		// 	errs =errors.Wrapf(err, "dst.getCoinBalance %v", err)
 ts.logger.Error(errs)
 return
		// }
		// if finalBtsBalance.UserBalance.Cmp(btsBalance.UserBalance) != 0 {
		// 	return fmt.Errorf("BTS Balance should have been same since txn does not succeed. Init %v  Final %v", btsBalance.UserBalance.String(), finalBtsBalance.UserBalance.String())
		// }
		errs =nil
 ts.logger.Error(errs)
 return
	},
}

var TransferAllBTSBalance Script = Script{
	Name:        "TransferAllBTSBalance",
	Type:        "Transfer",
	Description: "Transfer Native Tokens, which are of fixed supply, in different amounts. The Token should be native for both chains",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
	srcChain := tp.SrcChain
	dstChain := tp.DstChain
	coinNames := tp.CoinNames

		if len(coinNames) == 0 {
			errs =errors.New("Should specify at least one coinname, got zero")
 ts.logger.Error(errs)
 return
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs =errors.Wrapf(err, "GetChainPair %v", err)
 ts.logger.Error(errs)
 return
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
			errs =nil
 ts.logger.Error(errs)
 return
		}

		btsAddr, ok := ts.btsAddressPerChain[dstChain]
		if !ok {
			errs =errors.Wrapf(err, "dst.GetBTPAddressOfBTS %v", err)
 ts.logger.Error(errs)
 return
		}
		btsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			errs =errors.Wrapf(err, "dst.getCoinBalance %v", err)
 ts.logger.Error(errs)
 return
		}

		// prepare accounts
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs =errors.Wrapf(err, "GetKeyPairs %v", err)
 ts.logger.Error(errs)
 return
		}
 txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs =errors.Wrapf(err, "GetKeyPairs %v", err)
 ts.logger.Error(errs)
 return
		}
 txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		amt := ts.withFeeAdded(btsBalance.UserBalance)
		if errs = ts.Fund(srcAddr, amt, coinName); errs != nil {
			errs =errors.Wrapf(err, "Fund %v", err)
 ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
 return
		}

		// how much is necessary as gas cost
		if errs = ts.Fund(srcAddr, ts.SuggestGasPrice(), src.NativeCoin()); errs != nil {
			errs =errors.Wrapf(err, "AddGasFee %v", err)
 ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
 return
		}

		// approve
		if approveHash, err := src.Approve(coinName, srcKey, amt); err != nil {
			errs =errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
 ts.logger.Error(errs)
 return
		} else {
			if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
				errs =errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
 ts.logger.Error(errs)
 return
			}
		}
		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, amt)
		if err != nil {
			errs =errors.Wrapf(err, "Transfer %v", err)
 ts.logger.Error(errs)
 return
		}
		if err := ts.ValidateTransactionResult(ctx, hash, []string{coinName}, srcAddr, dstAddr, []*big.Int{amt}); err != nil {
			errs =errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
 ts.logger.Error(errs)
 return
		}
		err = ts.WaitForTransferEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
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
			errs =errors.Wrapf(err, "WaitForTransferEvents %v", err)
 ts.logger.Error(errs)
 return
		}
		// finalBtsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		// if err != nil {
		// 	errs =errors.Wrapf(err, "dst.getCoinBalance %v", err)
 ts.logger.Error(errs)
 return
		// }
		// if finalBtsBalance.UserBalance.Cmp(btsBalance.UserBalance) != 0 {
		// 	return fmt.Errorf("BTS Balance should have been same since txn does not succeed. Init %v  Final %v", btsBalance.UserBalance.String(), finalBtsBalance.UserBalance.String())
		// }
		errs =nil
 ts.logger.Error(errs)
 return
	},
}

*/

var TransferUniDirection Script = Script{
	Name:        "TransferUniDirection",
	Type:        "Transfer",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, tmpOk := ts.clsPerChain[srcChain]
		if !tmpOk {
			errs = fmt.Errorf("Chain %v not found", srcChain)
			ts.logger.Error(errs)
			return
		}
		dst, tmpOk := ts.clsPerChain[dstChain]
		if !tmpOk {
			errs = fmt.Errorf("Chain %v not found", srcChain)
			ts.logger.Error(errs)
			return
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}

		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		ts.logger.Debug("Accounts ", srcAddr, dstAddr)
		// How much tokens do we need on src and dst accounts ?
		tokenAmountAfterFeeChargeOnSrc := big.NewInt(1)
		tokenAmountBeforeFeeChargeOnSrc, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, tokenAmountAfterFeeChargeOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "getAmountBeforeFeeCharge %v", err)
			ts.logger.Error(errs)
			return
		}
		// How much native coins do we need to cover gas fee ?
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)

		// Fund source side with the required tokens
		// These tokens should be something a newly deployed god address can have
		if errs = ts.Fund(srcChain, srcAddr, tokenAmountBeforeFeeChargeOnSrc, coinName); errs != nil {
			ts.logger.Debug(errs)
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			ts.logger.Debug(errs)
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("initSrcNativeCoinBalance ", initSrcNativeCoinBalance)
		// ApproveToken On Source
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, tokenAmountBeforeFeeChargeOnSrc); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, tokenAmountBeforeFeeChargeOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}

		// Wait For Events
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})
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
				if startEvt.Assets[0].Value.Cmp(tokenAmountAfterFeeChargeOnSrc) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got startEvt.Value: %v AmountAfterFeeCharge %v", startEvt.Assets[0].Value, tokenAmountAfterFeeChargeOnSrc)
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
				if receivedEvt.Assets[0].Value.Cmp(tokenAmountAfterFeeChargeOnSrc) != 0 {
					return fmt.Errorf("Expected same value for coinName; Got receivedEvt.Value: %v AmountAfterFeeCharge %v", receivedEvt.Assets[0].Value, tokenAmountAfterFeeChargeOnSrc)
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
					ts.logger.Debug("Got Transfer End")
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		// Intermediate Tally
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "ChargedGasFee For Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		if coinName != src.NativeCoin() {
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				errs = errors.Wrapf(err, "ChargedGasFee for Approve Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			if tokenAmountBeforeFeeChargeOnSrc.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountBeforeFeeChargeOnSrc, tmpDiff)
				ts.logger.Error(errs)
				return
			}
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		} else {
			tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnSrc, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}
		tmpDiff := (&big.Int{}).Sub(finalDstBalance.UserBalance, initDstBalance.UserBalance)
		if tokenAmountAfterFeeChargeOnSrc.Cmp(tmpDiff) != 0 {
			errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountAfterFeeChargeOnSrc, tmpDiff)
			ts.logger.Error(errs)
			return
		}
		return
	},
}

var TransferBiDirection Script = Script{
	Name:        "TransferBiDirectionWithApprove",
	Type:        "Transfer",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, tmpOk := ts.clsPerChain[srcChain]
		if !tmpOk {
			errs = fmt.Errorf("Chain %v not found", srcChain)
			ts.logger.Error(errs)
			return
		}
		dst, tmpOk := ts.clsPerChain[dstChain]
		if !tmpOk {
			errs = fmt.Errorf("Chain %v not found", srcChain)
			ts.logger.Error(errs)
			return
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}

		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}
		ts.logger.Debug("Accounts ", srcAddr, dstAddr)
		// How much tokens do we need on src and dst accounts ?
		tokenAmountAfterFeeChargeOnDst := big.NewInt(1)
		tokenAmountBeforeFeeChargeOnDst, err := ts.getAmountBeforeFeeCharge(dstChain, coinName, tokenAmountAfterFeeChargeOnDst)
		if err != nil {
			errs = errors.Wrapf(err, "getAmountBeforeFeeCharge %v", err)
			ts.logger.Error(errs)
			return
		}
		tokenAmountBeforeFeeChargeOnSrc, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, tokenAmountBeforeFeeChargeOnDst)
		if err != nil {
			errs = errors.Wrapf(err, "getAmountBeforeFeeCharge %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("Tokens ", tokenAmountAfterFeeChargeOnDst, tokenAmountBeforeFeeChargeOnDst, tokenAmountBeforeFeeChargeOnSrc)
		// How much native coins do we need to cover gas fee ?
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gp := src.SuggestGasPrice()
		ts.logger.Debug("GP ", gp, gasLimitOnSrc)
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)
		gasLimitOnDst := big.NewInt(int64(ts.cfgPerChain[dstChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != dst.NativeCoin() {
			gasLimitOnDst.Add(gasLimitOnDst, big.NewInt(int64(ts.cfgPerChain[dstChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gp = dst.SuggestGasPrice()
		gasFeeOnDst := (&big.Int{}).Mul(gp, gasLimitOnDst)

		ts.logger.Debug("Gas fee on dst and src ", gasFeeOnDst, gasFeeOnSrc)
		// Fund source side with the required tokens
		// These tokens should be something a newly deployed god address can have
		if errs = ts.Fund(srcChain, srcAddr, tokenAmountBeforeFeeChargeOnSrc, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(err, "srcChain %v, srcAddr %v, amount %v, coinName %v err %v", srcChain, srcAddr, tokenAmountBeforeFeeChargeOnSrc, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(err, "srcChain %v, srcAddr %v, amount %v, coinName %v err %v", srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin(), errs))
			return
		}
		if errs = ts.Fund(dstChain, dstAddr, gasFeeOnDst, dst.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(err, "srcChain %v, srcAddr %v, amount %v, coinName %v err %v", srcChain, srcAddr, gasFeeOnDst, dst.NativeCoin(), errs))
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("InitSrc ", initSrcBalance, tokenAmountBeforeFeeChargeOnSrc)
		ts.logger.Debug("InitSrcNative ", initSrcNativeCoinBalance, tokenAmountBeforeFeeChargeOnSrc)
		ts.logger.Debug("InitDst ", initDstBalance)
		ts.logger.Debug("InitDstNative ", initDstNativeCoinBalance)
		// ApproveToken On Source
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, tokenAmountBeforeFeeChargeOnSrc); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, tokenAmountBeforeFeeChargeOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}

		// Wait For Events
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})
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
					ts.logger.Debug("Got Transfer End")
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		// Intermediate Tally
		intermediateSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		intermediateSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		intermediateDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "ChargedGasFee For Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		if coinName != src.NativeCoin() {
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				errs = errors.Wrapf(err, "ChargedGasFee for Approve Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, intermediateSrcBalance.UserBalance)
			if tokenAmountBeforeFeeChargeOnSrc.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountBeforeFeeChargeOnSrc, tmpDiff)
				ts.logger.Error(errs)
				return
			}
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, intermediateSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		} else {
			tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnSrc, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, intermediateSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}
		tmpDiff := (&big.Int{}).Sub(intermediateDstBalance.UserBalance, initDstBalance.UserBalance)
		if tokenAmountBeforeFeeChargeOnDst.Cmp(tmpDiff) != 0 {
			errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst, tmpDiff)
			ts.logger.Error(errs)
			return
		}

		// Start Transfer On Opposite Direction
		ts.logger.Debug("IntermediateSrc ", intermediateSrcBalance)
		ts.logger.Debug("IntermediateDst ", intermediateDstBalance, tokenAmountBeforeFeeChargeOnDst)

		// ApproveToken On Destination
		if coinName != dst.NativeCoin() {
			if approveHash, err = dst.Approve(coinName, dstKey, tokenAmountBeforeFeeChargeOnDst); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, dstChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Destination
		transferHashOnDst, err := dst.Transfer(coinName, dstKey, srcAddr, tokenAmountBeforeFeeChargeOnDst)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnDst ", transferHashOnDst)
		if _, err := ts.ValidateTransactionResult(ctx, dstChain, transferHashOnDst); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}

		// Wait For Events
		err = ts.WaitForTransferEvents(ctx, dstChain, srcChain, transferHashOnDst, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: dstChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})
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
					ts.logger.Debug("Got Transfer End")
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		// Final Tally
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		gasSpentOnTxn, err = dst.ChargedGasFee(transferHashOnDst)
		if err != nil {
			errs = errors.Wrapf(err, "GetGasUsed Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		if coinName != dst.NativeCoin() {
			gasSpentOnApprove, err := dst.ChargedGasFee(approveHash)
			if err != nil {
				errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff = (&big.Int{}).Sub(intermediateDstBalance.UserBalance, finalDstBalance.UserBalance)
			if tokenAmountBeforeFeeChargeOnDst.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst, tmpDiff)
				ts.logger.Error(errs)
				return
			}
			tmpDiff = (&big.Int{}).Sub(initDstNativeCoinBalance.UserBalance, finalDstNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for dst nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		} else {
			tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnDst, gasSpentOnTxn)
			tmpDiff = (&big.Int{}).Sub(intermediateDstBalance.UserBalance, finalDstBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}
		tmpDiff = (&big.Int{}).Sub(finalSrcBalance.UserBalance, intermediateSrcBalance.UserBalance)
		if tokenAmountAfterFeeChargeOnDst.Cmp(tmpDiff) != 0 {
			errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountAfterFeeChargeOnDst, tmpDiff)
			ts.logger.Error(errs)
			return
		}
		return
	},
}

var TransferToZeroAddress Script = Script{
	Name:        "TransferToZeroAddress",
	Description: "Transfer to zero address",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		_, tmpDstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
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

		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}

		// Wait For Events
		errs = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})
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
		if errs != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		// Final Tally
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		if coinName != src.NativeCoin() {
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			if tmpDiff.Cmp(feeCharged) != 0 {
				errs = fmt.Errorf("Expected same value; Got different feeCharged %v BalanceDiff %v", feeCharged, tmpDiff)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		} else {
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			tmpNativeCoinUsed := (&big.Int{}).Add(feeCharged, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same, Got Different. NativeCoinUsed %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}
		if initDstBalance.UserBalance.Cmp(finalDstBalance.UserBalance) != 0 {
			errs = fmt.Errorf("Epected same; Got Different. initDstBalance %v finalDstBalance %v", initDstBalance, finalDstBalance)
			ts.logger.Error(errs)
			return
		}
		return
	},
}

var TransferToUnknownNetwork Script = Script{
	Name:        "TransferToUnknownNetwork",
	Description: "Transfer to unknow network",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}
		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		_, tmpDstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
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

		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, errs = ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); errs != nil {
			//ts.logger.Error(errs)
			if errs.Error() == StatusCodeZero.Error() {
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					errs = errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				if coinName != src.NativeCoin() {
					if initSrcBalance.TotalBalance.Cmp(finalSrcBalance.TotalBalance) != 0 {
						errs = fmt.Errorf("Expected Same, Got Different. finalSrcBalance %v initialSrcBalance %v", finalSrcBalance.TotalBalance, initSrcBalance.TotalBalance)
						ts.logger.Error(errs)
						return
					}
					gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
					if err != nil {
						errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
						ts.logger.Error(errs)
						return
					}
					gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
						ts.logger.Error(errs)
						return
					}
				} else {
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
						ts.logger.Error(errs)
						return
					}
				}
				errs = nil
				return
			}
			errs = errors.Wrapf(err, "ValidateTransactionResult Got Unexpected Error: %v", err)
			ts.logger.Error(errs)
			return
		}
		errs = fmt.Errorf("Expected Transaction to fail but it did not")
		ts.logger.Error(errs)
		return
	},
}

var TransferWithoutApprove Script = Script{
	Name:        "TransferWithoutApprove",
	Description: "Transfer Without Approve",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}
		if len(coinNames) != 1 {
			errs = fmt.Errorf(" Should specify a single coinName, got %v", len(coinNames))
			ts.logger.Error(errs)
			return
		}
		coinName := coinNames[0]

		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		if coinName == src.NativeCoin() {
			ts.logger.Debugf("Expected non-native coin; Got native coin %v", src.NativeCoin())
			errs = UnsupportedCoinArgs // not returning an error here
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)

		if _, errs = ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); errs != nil {
			if errs.Error() == StatusCodeZero.Error() { // Failed as expected
				// Final Tally
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					errs = errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				if initSrcBalance.UserBalance.Cmp(finalSrcBalance.UserBalance) != 0 {
					errs = fmt.Errorf("Expected same value; Got different initSrcBalance %v finalSrcbalance %v", initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
					ts.logger.Error(errs)
					return
				}
				tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
				if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
					ts.logger.Error(errs)
					return
				}
				errs = nil
				return
			}
			errs = errors.Wrapf(err, "ValidateTransactionResult Got Unexpected Error: %v", err)
			ts.logger.Error(errs)
			return
		}
		errs = fmt.Errorf("Expected event to fail but it did not ")
		ts.logger.Error(errs)
		return
	},
}

var TransferLessThanFee Script = Script{
	Name:        "TransferLessThanFee",
	Description: "Transfer to unknow network",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}
		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		// Funds
		netTransferrableAmount := big.NewInt(-1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, errs = ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); errs != nil {
			if errs.Error() == StatusCodeZero.Error() {
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					errs = errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				if coinName != src.NativeCoin() {
					if initSrcBalance.TotalBalance.Cmp(finalSrcBalance.TotalBalance) != 0 {
						errs = fmt.Errorf("Expected Same, Got Different. finalSrcBalance %v initialSrcBalance %v", finalSrcBalance.TotalBalance, initSrcBalance.TotalBalance)
						ts.logger.Error(errs)
						return
					}
					gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
					if err != nil {
						errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
						ts.logger.Error(errs)
						return
					}
					gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
						ts.logger.Error(errs)
						return
					}
				} else {
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
						ts.logger.Error(errs)
						return
					}
				}
				errs = nil
				return
			}
			errs = errors.Wrapf(err, "ValidateTransactionResult Got Unexpected Error: %v", err)
			ts.logger.Error(errs)
			return
		}
		errs = fmt.Errorf("Expected event to fail but it did not ")
		ts.logger.Error(errs)
		return
	},
}

var TransferEqualToFee Script = Script{
	Name:        "TransferEqualToFee",
	Description: "Transfer equal to fee",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}

		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		// Funds
		netTransferrableAmount := big.NewInt(0)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		srcGasPrice := src.SuggestGasPrice()
		gasFeeOnSrc := (&big.Int{}).Mul(srcGasPrice, gasLimitOnSrc)

		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Record Initial Balance
		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}

		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, errs = ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); errs != nil {
			if errs.Error() == StatusCodeZero.Error() {
				finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
				if err != nil {
					errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
				if err != nil {
					errs = errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				if coinName != src.NativeCoin() {
					if initSrcBalance.TotalBalance.Cmp(finalSrcBalance.TotalBalance) != 0 {
						errs = fmt.Errorf("Expected Same, Got Different. finalSrcBalance %v initialSrcBalance %v", finalSrcBalance.TotalBalance, initSrcBalance.TotalBalance)
						ts.logger.Error(errs)
						return
					}
					gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
					if err != nil {
						errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
						ts.logger.Error(errs)
						return
					}
					gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
						ts.logger.Error(errs)
						return
					}
				} else {
					tmpDiff := (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
					if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
						errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
						ts.logger.Error(errs)
						return
					}
				}
				errs = nil
				return
			}
			errs = errors.Wrapf(err, "ValidateTransactionResult Got Unexpected Error: %v", err)
			ts.logger.Error(errs)
			return
		}
		// Wait For Events
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})
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
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}
		finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		finalSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "GetGasUsed For Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		if coinName != src.NativeCoin() {
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			if tmpDiff.Cmp(feeCharged) != 0 {
				errs = fmt.Errorf("Expected same value; Got different feeCharged %v BalanceDiff %v", feeCharged, tmpDiff)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnApprove, err := src.ChargedGasFee(approveHash)
			if err != nil {
				errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, finalSrcNativeCoinBalance.UserBalance)
			if gasSpentOnTxn.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value; Got different GasSpent %v NativeCoinBalanceDiff %v", gasSpentOnTxn, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		} else {
			feeCharged := (&big.Int{}).Sub(userSuppliedAmount, netTransferrableAmount)
			tmpNativeCoinUsed := (&big.Int{}).Add(feeCharged, gasSpentOnTxn)
			tmpDiff := (&big.Int{}).Sub(initSrcBalance.UserBalance, finalSrcBalance.UserBalance)
			if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same, Got Different. NativeCoinUsed %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}
		if initDstBalance.UserBalance.Cmp(finalDstBalance.UserBalance) != 0 {
			errs = fmt.Errorf("Epected same; Got Different. initDstBalance %v finalDstBalance %v", initDstBalance, finalDstBalance)
			ts.logger.Error(errs)
			return
		}
		errs = err
		ts.logger.Error(errs)
		return
	},
}

var TransferToBlackListedDstAddress Script = Script{
	Name:        "TransferToBlackListedDstAddress",
	Description: "Transfer to BlackListed Destination Address",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)

		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})
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
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("Start Add to blacklist")

		// Add To BlackList
		fCfg, err := ts.GetFullConfigAPI()
		if err != nil {
			errs = errors.Wrapf(err, "GetFullConfigAPI %v", err)
			ts.logger.Error(errs)
			return
		}
		fCfgOwnerKey := ts.FullConfigAPIsOwner()
		stdCfg, err := ts.GetStandardConfigAPI(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetStandardConfigAPI %v", err)
			ts.logger.Error(errs)
			return
		}
		blDstNet, blDstAddr := ts.NetAddr(dstAddr)
		blackListAddHash, err := fCfg.AddBlackListAddress(fCfgOwnerKey, blDstNet, []string{blDstAddr})
		if err != nil {
			errs = errors.Wrapf(err, "AddBlackListAddress %v", err)
			ts.logger.Error(errs)
			return
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListAddHash); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		isBlackListed, err := fCfg.IsUserBlackListed(blDstNet, blDstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "isBlackListed %v", err)
			ts.logger.Error(errs)
			return
		} else if err == nil && !isBlackListed {
			errs = fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blDstNet, blDstAddr)
			ts.logger.Error(errs)
			return
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
						txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
							ChainName: ts.FullConfigAPIChain(),
							Sn:        reqEvt.Sn,
							Fee:       map[string]*big.Int{},
						})
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
				errs = errors.Wrapf(err, "WaitForConfigResponse %v", err)
				ts.logger.Error(errs)
				return
			}
			isBlackListed, err = fCfg.IsUserBlackListed(blDstNet, blDstAddr)
			if err != nil {
				errs = errors.Wrapf(err, "isBlackListed %v", err)
				ts.logger.Error(errs)
				return
			} else if err == nil && !isBlackListed {
				errs = fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blDstNet, blDstAddr)
				ts.logger.Error(errs)
				return
			}
		}

		// Send After BlackListing
		ts.logger.Debug("Send to blacklist")
		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err == nil {
			// Wait For Events
			err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
				chain.TransferStart: func(ev *evt) error {
					if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
						return errors.New("Got nil value for event ")
					}
					startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
					if !ok {
						return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
					}
					txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
						ChainName: srcChain,
						Sn:        startEvt.Sn,
						Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
					})
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
				errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
				ts.logger.Error(errs)
				return
			}
		} else {
			if err.Error() != StatusCodeZero.Error() {
				errs = errors.Wrapf(err, "ValidateTransactionResult Got Unexpected Error: %v", err)
				ts.logger.Error(errs)
				return
			}
		}

		// Remove From BlackList
		ts.logger.Debug("Remove From BlackList")

		blackListRemoveHash, err := fCfg.RemoveBlackListAddress(fCfgOwnerKey, blDstNet, []string{blDstAddr})
		if err != nil {
			errs = errors.Wrapf(err, "RemoveBlackListAddress %v", err)
			ts.logger.Error(errs)
			return
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListRemoveHash); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		isBlackListed, err = fCfg.IsUserBlackListed(blDstNet, blDstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "isBlackListed %v", err)
			ts.logger.Error(errs)
			return
		} else if err == nil && isBlackListed {
			errs = fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blDstNet, blDstAddr)
			ts.logger.Error(errs)
			return
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
						txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
							ChainName: ts.FullConfigAPIChain(),
							Sn:        reqEvt.Sn,
							Fee:       map[string]*big.Int{},
						})
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
				errs = errors.Wrapf(err, "WaitForConfigResponse %v", err)
				ts.logger.Error(errs)
				return
			}
			isBlackListed, err = stdCfg.IsUserBlackListed(blDstNet, blDstAddr)
			if err != nil {
				errs = errors.Wrapf(err, "isBlackListed %v", err)
				ts.logger.Error(errs)
				return
			} else if err == nil && isBlackListed {
				errs = fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blDstNet, blDstAddr)
				ts.logger.Error(errs)
				return
			}
		}

		ts.logger.Debug("Final Send Should Succeed")
		// Final Send Should Succeed
		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})

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
				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}

				return nil
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("pass")
		return
	},
}

var TransferFromBlackListedSrcAddress Script = Script{
	Name:        "TransferFromBlackListedSrcAddress",
	Description: "Transfer from BlackListed Source Address",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs // fmt.Errorf(" Should specify at least one coinName, got zero")
			ts.logger.Debug(errs)
			return
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetChainPair %v", err)
			ts.logger.Error(errs)
			return
		}
		// Account
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}

		// Funds
		netTransferrableAmount := big.NewInt(1)
		userSuppliedAmount, err := ts.getAmountBeforeFeeCharge(srcChain, coinName, netTransferrableAmount)
		gasLimitOnSrc := big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.TransferCoinInterChainGasLimit]))
		if coinName != src.NativeCoin() {
			gasLimitOnSrc.Add(gasLimitOnSrc, big.NewInt(int64(ts.cfgPerChain[srcChain].GasLimit[chain.ApproveTokenInterChainGasLimit])))
		}
		gasFeeOnSrc := (&big.Int{}).Mul(src.SuggestGasPrice(), gasLimitOnSrc)

		if errs = ts.Fund(srcChain, srcAddr, (&big.Int{}).Mul(userSuppliedAmount, big.NewInt(2)), coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, (&big.Int{}).Mul(gasFeeOnSrc, big.NewInt(2)), src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Approve
		var approveHash string
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err := src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})

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

				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("Start Add to blacklist")

		// Add To BlackList
		fCfg, err := ts.GetFullConfigAPI()
		if err != nil {
			errs = errors.Wrapf(err, "GetFullConfigAPI %v", err)
			ts.logger.Error(errs)
			return
		}
		fCfgOwnerKey := ts.FullConfigAPIsOwner()
		stdCfg, err := ts.GetStandardConfigAPI(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetStandardConfigAPI %v", err)
			ts.logger.Error(errs)
			return
		}
		blsSrcNet, blsSrcAddr := ts.NetAddr(srcAddr)
		blackListAddHash, err := fCfg.AddBlackListAddress(fCfgOwnerKey, blsSrcNet, []string{blsSrcAddr})
		if err != nil {
			errs = errors.Wrapf(err, "AddBlackListAddress %v", err)
			ts.logger.Error(errs)
			return
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListAddHash); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		isBlackListed, err := fCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "isBlackListed %v", err)
			ts.logger.Error(errs)
			return
		} else if err == nil && !isBlackListed {
			errs = fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blsSrcNet, blsSrcAddr)
			ts.logger.Error(errs)
			return
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
						txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
							ChainName: ts.FullConfigAPIChain(),
							Sn:        reqEvt.Sn,
							Fee:       map[string]*big.Int{},
						})

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
				errs = errors.Wrapf(err, "WaitForConfigResponse %v", err)
				ts.logger.Error(errs)
				return
			}
			isBlackListed, err = stdCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
			if err != nil {
				errs = errors.Wrapf(err, "isBlackListed %v", err)
				ts.logger.Error(errs)
				return
			} else if err == nil && !isBlackListed {
				errs = fmt.Errorf("Expected addr ( %v , %v ) to be blacklisted, but was not", blsSrcNet, blsSrcAddr)
				ts.logger.Error(errs)
				return
			}
		}
		// Send After BlackListing
		ts.logger.Debug("Send to blacklist")

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil && err.Error() != StatusCodeZero.Error() {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err == nil {
			// Wait For Events
			err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
				chain.TransferStart: func(ev *evt) error {
					if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
						return errors.New("Got nil value for event ")
					}
					startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
					if !ok {
						return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
					}
					txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
						ChainName: srcChain,
						Sn:        startEvt.Sn,
						Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
					})

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
				errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
				ts.logger.Error(errs)
				return
			}
		} else {
			if err.Error() != StatusCodeZero.Error() {
				errs = errors.Wrapf(err, "ValidateTransactionResult Got Unexpected Error: %v", err)
				ts.logger.Error(errs)
				return
			}
		}

		// Remove From BlackList
		ts.logger.Debug("Remove From BlackList")

		blackListRemoveHash, err := fCfg.RemoveBlackListAddress(fCfgOwnerKey, blsSrcNet, []string{blsSrcAddr})
		if err != nil {
			errs = errors.Wrapf(err, "RemoveBlackListAddress %v", err)
			ts.logger.Error(errs)
			return
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), blackListRemoveHash); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		isBlackListed, err = fCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "isBlackListed %v", err)
			ts.logger.Error(errs)
			return
		} else if err == nil && isBlackListed {
			errs = fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blsSrcNet, blsSrcAddr)
			ts.logger.Error(errs)
			return
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
						txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
							ChainName: ts.FullConfigAPIChain(),
							Sn:        reqEvt.Sn,
							Fee:       map[string]*big.Int{},
						})

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
				errs = errors.Wrapf(err, "WaitForConfigResponse %v", err)
				ts.logger.Error(errs)
				return
			}
			if isBlackListed, err = fCfg.IsUserBlackListed(blsSrcNet, blsSrcAddr); err != nil {
				errs = errors.Wrapf(err, "isBlackListed %v", err)
				ts.logger.Error(errs)
				return
			} else if err == nil && isBlackListed {
				errs = fmt.Errorf("Expected addr ( %v , %v ) to not be blacklisted, but was blackListed", blsSrcNet, blsSrcAddr)
				ts.logger.Error(errs)
				return
			}
		}

		ts.logger.Debug("Final Send Should Succeed")
		// Final Send Should Succeed
		if errs = ts.Fund(srcChain, srcAddr, userSuppliedAmount, coinName); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(errs, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinName, errs))
			return
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err = src.Approve(coinName, srcKey, userSuppliedAmount); err != nil {
				errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				ts.logger.Error(errs)
				return
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash); err != nil {
					errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					ts.logger.Error(errs)
					return
				}
			}
		}
		// Transfer On Source
		transferHashOnSrc, err = src.Transfer(coinName, srcKey, dstAddr, userSuppliedAmount)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{startEvt.Assets[0].Name: startEvt.Assets[0].Fee},
				})

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

				if endEvt.Code.String() != "0" {
					return fmt.Errorf("Expected code 0 Got %v", endEvt.Code.String())
				}
				return nil
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("pass")
		return
	},
}

var TransferBatchBiDirection Script = Script{
	Name:        "TransferBatchBiDirection",
	Description: "Transfer batch bi-direction",
	Type:        "Transfer",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) == 0 {
			errs = UnsupportedCoinArgs //
			ts.logger.Debug(fmt.Errorf(" Should specify at least one coinName, got zero"))
			return
		}
		src, tmpOk := ts.clsPerChain[srcChain]
		if !tmpOk {
			errs = fmt.Errorf("Chain %v not found", srcChain)
			ts.logger.Error(errs)
			return
		}
		dst, tmpOk := ts.clsPerChain[dstChain]
		if !tmpOk {
			errs = fmt.Errorf("Chain %v not found", srcChain)
			ts.logger.Error(errs)
			return
		}
		if len(coinNames) == 1 && coinNames[0] == src.NativeCoin() {
			errs = UnsupportedCoinArgs //
			ts.logger.Error(fmt.Errorf("A single src.NativeCoin %v has been used", coinNames[0]))
			return
		}
		if len(coinNames) == 1 && coinNames[0] == dst.NativeCoin() {
			errs = UnsupportedCoinArgs //
			ts.logger.Error(fmt.Errorf("A single dst.NativeCoin %v has been used", coinNames[0]))
			return
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[srcChain] = []keypair{{PrivKey: srcKey, PubKey: srcAddr}}
		dstKey, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			errs = errors.Wrapf(err, "GetKeyPairs %v", err)
			ts.logger.Error(errs)
			return
		}
		txnRec.addresses[dstChain] = []keypair{{PrivKey: dstKey, PubKey: dstAddr}}
		ts.logger.Debug("Calculate Transfer Amounts")
		// Calculate Transfer Amounts
		tokenAmountAfterFeeChargeOnDst := big.NewInt(1)
		tokenAmountBeforeFeeChargeOnDst := make([]*big.Int, len(coinNames))
		nativeCoinAmountBeforeFeeChargeOnDst := big.NewInt(0)
		for i := 0; i < len(coinNames); i++ {
			tokenAmountBeforeFeeChargeOnDst[i], err = ts.getAmountBeforeFeeCharge(dstChain, coinNames[i], tokenAmountAfterFeeChargeOnDst)
			if err != nil {
				errs = errors.Wrapf(err, "getAmountBeforeFeeCharge(%v) %v", coinNames[i], err)
				ts.logger.Error(errs)
				return
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
				errs = errors.Wrapf(err, "getAmountBeforeFeeCharge(%v) %v", coinNames[i], err)
				ts.logger.Error(errs)
				return
			}
			if coinNames[i] == src.NativeCoin() {
				nativeCoinAmountBeforeFeeChargeOnSrc.Set(tokenAmountBeforeFeeChargeOnSrc[i])
			}
		}
		ts.logger.Debug("Calculate Gas Fees")
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
		//gasFeeOnDst.Add(gasFeeOnDst, big.NewInt(500000000000000))
		ts.logger.Debug("Fund")
		//Fund
		for i := 0; i < len(coinNames); i++ {
			if errs = ts.Fund(srcChain, srcAddr, tokenAmountBeforeFeeChargeOnSrc[i], coinNames[i]); errs != nil {
				// errs = errors.Wrapf(err, "Fund Token %v", err)
				ts.logger.Debug(errors.Wrapf(err, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, coinNames[i], err))
				return
			}
		}
		if errs = ts.Fund(srcChain, srcAddr, gasFeeOnSrc, src.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(err, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, src.NativeCoin(), err))
			return
		}
		if errs = ts.Fund(dstChain, dstAddr, gasFeeOnDst, dst.NativeCoin()); errs != nil {
			// errs = errors.Wrapf(err, "Fund Token %v", err)
			ts.logger.Debug(errors.Wrapf(err, "srcChain %v, srcAddr %v, coinName %v err %v", srcChain, srcAddr, dst.NativeCoin(), err))
			return
		}

		ts.logger.Debug("Record Initial Balance")
		// Record Initial Balance
		initSrcBalance := make([]*chain.CoinBalance, len(coinNames))
		initDstBalance := make([]*chain.CoinBalance, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			initSrcBalance[i], err = src.GetCoinBalance(coinNames[i], srcAddr)
			if err != nil {
				errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			initDstBalance[i], err = dst.GetCoinBalance(coinNames[i], dstAddr)
			if err != nil {
				errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				ts.logger.Error(errs)
				return
			}
		}
		initSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		initDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("ApproveToken On Source")
		// ApproveToken On Source
		approveHash := map[string]string{}
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash[coinName], err = src.Approve(coinName, srcKey, tokenAmountBeforeFeeChargeOnSrc[i]); err != nil {
					errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash[coinName])
					ts.logger.Error(errs)
					return
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, srcChain, approveHash[coinName]); err != nil {
						errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash[coinName])
						ts.logger.Error(errs)
						return
					}
				}
			}
		}

		ts.logger.Debug("Transfer On Source")
		// Transfer On Source
		transferHashOnSrc, err := src.TransferBatch(coinNames, srcKey, dstAddr, tokenAmountBeforeFeeChargeOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnSrc ", transferHashOnSrc)
		if _, err := ts.ValidateTransactionResult(ctx, srcChain, transferHashOnSrc); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("Wait For Events")
		// Wait For Events
		err = ts.WaitForTransferEvents(ctx, srcChain, dstChain, transferHashOnSrc, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: srcChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{},
				})
				for i := 0; i < len(startEvt.Assets); i++ {
					txnRec.feeRecords[len(txnRec.feeRecords)-1].Fee[startEvt.Assets[i].Name] = startEvt.Assets[i].Fee
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
					ts.logger.Debug("Got Transfer End")
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("Intermediate Tally")
		// Intermediate Tally
		intermediateSrcBalance := make([]*chain.CoinBalance, len(coinNames))
		intermediateDstBalance := make([]*chain.CoinBalance, len(coinNames))
		for i, coinName := range coinNames {
			intermediateSrcBalance[i], err = src.GetCoinBalance(coinName, srcAddr)
			if err != nil {
				errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			intermediateDstBalance[i], err = dst.GetCoinBalance(coinName, dstAddr)
			if err != nil {
				errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				ts.logger.Error(errs)
				return
			}
		}
		intermediateSrcNativeCoinBalance, err := src.GetCoinBalance(src.NativeCoin(), srcAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}

		gasSpentOnTxn, err := src.ChargedGasFee(transferHashOnSrc)
		if err != nil {
			errs = errors.Wrapf(err, "ChargedGasFee For Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		for _, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				gasSpentOnApprove, err := src.ChargedGasFee(approveHash[coinName])
				if err != nil {
					errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			}
		}
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				tmpDiff := (&big.Int{}).Sub(initSrcBalance[i].UserBalance, intermediateSrcBalance[i].UserBalance)
				if tokenAmountBeforeFeeChargeOnSrc[i].Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountBeforeFeeChargeOnSrc[i], tmpDiff)
					ts.logger.Error(errs)
					return
				}
				tmpDiff = (&big.Int{}).Sub(initSrcNativeCoinBalance.UserBalance, intermediateSrcNativeCoinBalance.UserBalance)
				tmpNativeCoinUsed := (&big.Int{}).Add(nativeCoinAmountBeforeFeeChargeOnSrc, gasSpentOnTxn)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value for src nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
					ts.logger.Error(errs)
					return
				}
			} else {
				tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnSrc[i], gasSpentOnTxn)
				tmpDiff := (&big.Int{}).Sub(initSrcBalance[i].UserBalance, intermediateSrcBalance[i].UserBalance)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
					ts.logger.Error(errs)
					return
				}
			}
			tmpDiff := (&big.Int{}).Sub(intermediateDstBalance[i].UserBalance, initDstBalance[i].UserBalance)
			if tokenAmountBeforeFeeChargeOnDst[i].Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst[i], tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}

		ts.logger.Debug("ApproveToken On Destination")
		// ApproveToken On Destination
		approveHash = map[string]string{}
		for i, coinName := range coinNames {
			if coinName != dst.NativeCoin() {
				if approveHash[coinName], err = dst.Approve(coinName, dstKey, tokenAmountBeforeFeeChargeOnDst[i]); err != nil {
					errs = errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash[coinName])
					ts.logger.Error(errs)
					return
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, dstChain, approveHash[coinName]); err != nil {
						errs = errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash[coinName])
						ts.logger.Error(errs)
						return
					}
				}
			}
		}

		ts.logger.Debug("Transfer On Destination")
		// Transfer On Destination
		transferHashOnDst, err := dst.TransferBatch(coinNames, dstKey, srcAddr, tokenAmountBeforeFeeChargeOnDst)
		if err != nil {
			errs = errors.Wrapf(err, "Transfer Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		ts.logger.Debug("transferHashOnDst ", transferHashOnDst)
		if _, err := ts.ValidateTransactionResult(ctx, dstChain, transferHashOnDst); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("Wait For Events")
		// Wait For Events
		err = ts.WaitForTransferEvents(ctx, dstChain, srcChain, transferHashOnDst, map[chain.EventLogType]func(*evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
					ChainName: dstChain,
					Sn:        startEvt.Sn,
					Fee:       map[string]*big.Int{},
				})

				for i := 0; i < len(startEvt.Assets); i++ {
					txnRec.feeRecords[len(txnRec.feeRecords)-1].Fee[startEvt.Assets[i].Name] = startEvt.Assets[i].Fee
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
					ts.logger.Debug("Got Transfer End")
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			errs = errors.Wrapf(err, "WaitForTransferEvents %v", err)
			ts.logger.Error(errs)
			return
		}

		ts.logger.Debug("final tally")
		// Final Tally
		finalSrcBalance := make([]*chain.CoinBalance, len(coinNames))
		finalDstBalance := make([]*chain.CoinBalance, len(coinNames))
		for i, coinName := range coinNames {
			finalSrcBalance[i], err = src.GetCoinBalance(coinName, srcAddr)
			if err != nil {
				errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				ts.logger.Error(errs)
				return
			}
			finalDstBalance[i], err = dst.GetCoinBalance(coinName, dstAddr)
			if err != nil {
				errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
				ts.logger.Error(errs)
				return
			}
		}
		finalDstNativeCoinBalance, err := dst.GetCoinBalance(dst.NativeCoin(), dstAddr)
		if err != nil {
			errs = errors.Wrapf(err, "GetCoinBalance Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		gasSpentOnTxn, err = dst.ChargedGasFee(transferHashOnDst)
		if err != nil {
			errs = errors.Wrapf(err, "GetGasUsed Err: %v", err)
			ts.logger.Error(errs)
			return
		}
		for _, coinName := range coinNames {
			if coinName != dst.NativeCoin() {
				gasSpentOnApprove, err := dst.ChargedGasFee(approveHash[coinName])
				if err != nil {
					errs = errors.Wrapf(err, "GetGasUsed for Approve Err: %v", err)
					ts.logger.Error(errs)
					return
				}
				gasSpentOnTxn.Add(gasSpentOnTxn, gasSpentOnApprove)
			}
		}
		for i, coinName := range coinNames {
			if coinName != dst.NativeCoin() {
				tmpDiff := (&big.Int{}).Sub(intermediateDstBalance[i].UserBalance, finalDstBalance[i].UserBalance)
				if tokenAmountBeforeFeeChargeOnDst[i].Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tokenAmountBeforeFeeChargeOnDst[i], tmpDiff)
					ts.logger.Error(errs)
					return
				}
				tmpNativeCoinUsed := (&big.Int{}).Add(nativeCoinAmountBeforeFeeChargeOnDst, gasSpentOnTxn)
				tmpDiff = (&big.Int{}).Sub(initDstNativeCoinBalance.UserBalance, finalDstNativeCoinBalance.UserBalance)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value for dst nativeCoin balance after txn; Got GasSpentOnTxn %v srcDiffAmt %v", gasSpentOnTxn, tmpDiff)
					ts.logger.Error(errs)
					return
				}
			} else {
				tmpNativeCoinUsed := (&big.Int{}).Add(tokenAmountBeforeFeeChargeOnDst[i], gasSpentOnTxn)
				tmpDiff := (&big.Int{}).Sub(intermediateDstBalance[i].UserBalance, finalDstBalance[i].UserBalance)
				if tmpNativeCoinUsed.Cmp(tmpDiff) != 0 {
					errs = fmt.Errorf("Expected same value for dst balance After transfer, Got TransferAmt %v DstDiffAmt %v", tmpNativeCoinUsed, tmpDiff)
					ts.logger.Error(errs)
					return
				}
			}
			tmpDiff := (&big.Int{}).Sub(finalSrcBalance[i].UserBalance, intermediateSrcBalance[i].UserBalance)
			if tokenAmountAfterFeeChargeOnDst.Cmp(tmpDiff) != 0 {
				errs = fmt.Errorf("Expected same value for src balance After transfer, Got TransferAmt %v SrcDiffAmt %v", tokenAmountAfterFeeChargeOnDst, tmpDiff)
				ts.logger.Error(errs)
				return
			}
		}
		ts.logger.Debug("Pass")
		return
	},
}
