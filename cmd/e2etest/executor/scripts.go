package executor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

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
