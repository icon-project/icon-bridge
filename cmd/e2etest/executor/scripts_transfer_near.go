package executor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

var TransferUniDirectionNear Script = Script{
	Name:        "TransferUniDirection",
	Type:        "Transfer",
	Description: "Transfer Fixed Amount of coin and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		srcChain := tp.SrcChain
		dstChain := tp.DstChain
		coinNames := tp.CoinNames

		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}

		if len(coinNames) != 1 {
			errs = UnsupportedCoinArgs
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
			gasSpentOnTxn.Add(gasSpentOnTxn, big.NewInt(0)) // Passing Zero as we don't spend gas on approve
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
