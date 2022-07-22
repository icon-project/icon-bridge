package executor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

var TransferToIncorrectAddress Script = Script{
	Name:        "TransferToIncorrectAddress",
	Description: "Transfer Address that does not comply to the format that the recipeient accepts",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) error {
		if len(coinNames) == 0 {
			return errors.New("Should specify at least one coinname, got zero")
		}
		coinName := coinNames[0]
		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		//Make incorrect by adding a string at the end
		dstAddr += "1"

		amt := new(big.Int)
		amt.SetString("1000000000000000000", 10)
		if err := ts.Fund(srcAddr, amt, coinName); err != nil {
			return errors.Wrapf(err, "Fund %v", err)
		}
		if coinName != src.NativeCoin() {
			gasFee := new(big.Int)
			gasFee.SetString("1000000000000000000", 10)
			if err := ts.Fund(srcAddr, gasFee, src.NativeCoin()); err != nil {
				return errors.Wrapf(err, "Fund %v", err)
			}
		}

		// Approve
		if coinName != src.NativeCoin() {
			if approveHash, err := src.Approve(coinName, srcKey, amt); err != nil {
				return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			} else {
				if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
					return errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
				}
			}
		}

		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		if err := ts.ValidateTransactionResultEvents(ctx, hash, []string{coinName}, srcAddr, dstAddr, []*big.Int{amt}); err != nil {
			return errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferEnd: func(e *evt) error {
				endEvt, ok := e.msg.EventLog.(chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent. Got %T", e.msg.EventLog)
				}
				if endEvt.Code.Cmp(big.NewInt(1)) == 0 { //&& endEvt.Response == "InvalidAddress" {
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			return errors.Wrapf(err, "WaitForEvents %v", err)
		}

		return nil
	},
}

var TransferExceedingContractsBalance Script = Script{
	Name:        "TransferExceedingContractsBalance",
	Description: "Transfer Native Tokens, which are of fixed supply, in different amounts. The Token should be native for both chains",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) error {
		if len(coinNames) == 0 {
			return errors.New("Should specify at least one coinname, got zero")
		}
		coinName := coinNames[0]
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetChainPair %v", err)
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
			return fmt.Errorf("Token %v does not exist on both chains %v and %v", coinName, srcChain, dstChain)
		}

		btsAddr, err := dst.GetBTPAddressOfBTS()
		if err != nil {
			return errors.Wrapf(err, "dst.GetBTPAddressOfBTS %v", err)
		}
		btsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			return errors.Wrapf(err, "dst.getCoinBalance %v", err)
		}

		// prepare accounts
		amt := big.NewInt(1).Mul(btsBalance.UserBalance, big.NewInt(2))
		srcKey, srcAddr, err := ts.GetGodKeyPairs(srcChain)
		if err != nil {
			return errors.Wrapf(err, "GetGodKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		// approve
		if approveHash, err := src.Approve(coinName, srcKey, amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		} else {
			if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
				return errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
			}
		}
		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer %v", err)
		}
		if err := ts.ValidateTransactionResultEvents(ctx, hash, []string{coinName}, srcAddr, dstAddr, []*big.Int{amt}); err != nil {
			return errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			// chain.TransferReceived: func(e *evt) error {
			// 	ts.logger.Info("Got TransferReceived")
			// 	return nil
			// },
			chain.TransferEnd: func(e *evt) error {
				endEvt, ok := e.msg.EventLog.(chain.TransferEndEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferEndEvent. Got %T", e.msg.EventLog)
				}
				if endEvt.Code.Cmp(big.NewInt(1)) == 0 { //&& endEvt.Response == "TransferFailed" {
					return nil
				}
				return fmt.Errorf("Unexpected code %v and response %v", endEvt.Code, endEvt.Response)
			},
		})
		if err != nil {
			return errors.Wrapf(err, "WaitForEvents %v", err)
		}
		finalBtsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			return errors.Wrapf(err, "dst.getCoinBalance %v", err)
		}
		if finalBtsBalance.UserBalance.Cmp(btsBalance.UserBalance) != 0 {
			return fmt.Errorf("BTS Balance should have been same since txn does not succeed. Init %v  Final %v", btsBalance.UserBalance.String(), finalBtsBalance.UserBalance.String())
		}
		return nil
	},
}

var BasicTransfer Script = Script{
	Name:        "BasicTransfer",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) error {
		if len(coinNames) == 0 {
			return errors.New("Should specify at least one coinname, got zero")
		}

		src, _, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetChainPair %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}

		tmpAmt := new(big.Int)
		tmpAmt.SetString("1000000000000000000", 10)
		amts := make([]*big.Int, len(coinNames))
		for i := 0; i < len(coinNames); i++ {
			amts[i] = tmpAmt
		}

		hasNative := false
		for i, coinName := range coinNames {
			if err := ts.Fund(srcAddr, amts[i], coinName); err != nil {
				return errors.Wrapf(err, "Fund %v", err)
			}
			if coinName == src.NativeCoin() {
				hasNative = true
			}
		}
		if !hasNative {
			gasFee := new(big.Int)
			gasFee.SetString("1000000000000000000", 10)
			if err := ts.Fund(srcAddr, gasFee, src.NativeCoin()); err != nil {
				return errors.Wrapf(err, "Fund %v", err)
			}
		}

		// Approve
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err := src.Approve(coinName, srcKey, amts[i]); err != nil {
					return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
				} else {
					if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						return errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
					}
				}
			}
		}

		var hash string
		if len(coinNames) == 1 {
			hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
			if err != nil {
				return errors.Wrapf(err, "Transfer Err: %v", err)
			}
		} else {
			hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
			if err != nil {
				return errors.Wrapf(err, "Transfer Err: %v", err)
			}
		}

		if err := ts.ValidateTransactionResultEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
			return errors.Wrapf(err, "ValidateTransactionResultEvents %v", err)
		}
		err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
			chain.TransferReceived: func(e *evt) error {
				return nil
			},
			chain.TransferEnd: func(e *evt) error {
				return nil
			},
		})
		if err != nil {
			return errors.Wrapf(err, "WaitForEvents %v", err)
		}

		// finalSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		// if err != nil {
		// 	return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		// }
		// finalDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		// if err != nil {
		// 	return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		// }
		// //args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
		// //args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
		// if finalDstBalance.Usable.Cmp(initDstBalance.Usable) != 1 {
		// 	return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
		// }
		// if finalSrcBalance.Usable.Cmp(initSrcBalance.Usable) != -1 {
		// 	return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
		// }
		return nil
	},
}
