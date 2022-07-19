package executor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

var TransferExceedingContractsBalance Script = Script{
	Name:        "TransferExceedingContractsBalance",
	Description: "Transfer Native Tokens, which are of fixed supply, in different amounts",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinName string, ts *testSuite) error {
		src, dst, err := ts.GetChainPair(srcChain, dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetChainPair %v", err)
		}
		btsAddr, err := dst.GetBTPAddressOfBTS()
		if err != nil {
			return errors.Wrapf(err, "dst.GetBTPAddressOfBTS %v", err)
		}
		btsBalance, err := dst.GetCoinBalance(coinName, btsAddr)
		if err != nil {
			return errors.Wrapf(err, "dst.getCoinBalance %v", err)
		}
		fmt.Printf("Init %+v \n", btsBalance)
		// prepare accounts
		amt := big.NewInt(1).Mul(btsBalance.Total, big.NewInt(2))
		fmt.Printf("Transferring %+v \n", amt.String())
		return nil
		_, dstAddr, err := ts.GetKeyPairs(dstChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		srcKey, srcAddr, err := ts.GetKeyPairs(srcChain)
		if err != nil {
			return errors.Wrapf(err, "GetKeyPairs %v", err)
		}
		if err := ts.Fund(srcAddr, amt, coinName); err != nil {
			return errors.Wrapf(err, "Fund %v", err)
		}
		if coinName != src.NativeCoin() {
			gasFee := new(big.Int)
			gasFee.SetString("10000000000000000000", 10)
			if err := ts.Fund(srcAddr, gasFee, src.NativeCoin()); err != nil {
				return errors.Wrapf(err, "Fund %v", err)
			}
		}

		// approve
		if approveHash, err := src.Approve(coinName, srcKey, *amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		} else {
			if _, err := ts.ValidateTransactionResult(ctx, approveHash); err != nil {
				return errors.Wrapf(err, "Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
			}
		}
		// Transfer
		hash, err := src.Transfer(coinName, srcKey, dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer %v", err)
		}
		if err := ts.ValidateTransactionResultEvents(ctx, hash, coinName, srcAddr, dstAddr, amt); err != nil {
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
				if endEvt.Code.Cmp(big.NewInt(1)) == 0 && endEvt.Response == "TransferFailed" {
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
		if finalBtsBalance.Usable.Cmp(btsBalance.Usable) != 0 {
			return fmt.Errorf("BTS Balance should have been same since txn does not succeed. Init %v  Final %v", btsBalance.Usable.String(), finalBtsBalance.Usable.String())
		}
		return nil
	},
}

var Transfer Script = Script{
	Name:        "Transfer",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, srcChain, dstChain chain.ChainType, coinName string, ts *testSuite) error {
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

		amt := new(big.Int)
		amt.SetString("1000000000000000000", 10)
		if err := ts.Fund(srcAddr, amt, coinName); err != nil {
			return errors.Wrapf(err, "Fund %v", err)
		}
		if coinName != src.NativeCoin() {
			if err := ts.Fund(srcAddr, amt, src.NativeCoin()); err != nil {
				return errors.Wrapf(err, "Fund %v", err)
			}
		}

		initSrcBalance, err := src.GetCoinBalance(coinName, srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}

		// initDstBalance, err := dst.GetCoinBalance(coinName, dstAddr)
		// if err != nil {
		// 	return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		// }
		if initSrcBalance.Usable.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 1000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if coinName != src.NativeCoin() {
			if approveHash, err := src.Approve(coinName, srcKey, *amt); err != nil {
				return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			}
		}

		hash, err := src.Transfer(coinName, srcKey, dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		if err := ts.ValidateTransactionResultEvents(ctx, hash, coinName, srcAddr, dstAddr, amt); err != nil {
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
