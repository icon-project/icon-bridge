package executor

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
)

var TransferWithoutApproveFromICON Script = Script{
	Name:        "TransferWithoutApproveFromICON",
	Description: "Transfer Fixed Amount of coin without approve",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		iconTxResult, ok := res.(*icon.TransactionResult)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
		}
		if status, err := iconTxResult.Status.Int(); err != nil {
			return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
		} else if status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
		}

		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}
		time.Sleep(time.Second * 10) // Waiting For Relay
		finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
		// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
		if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
		}
		if finalDstBalance.Cmp(initDstBalance) != 1 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
		}
		if finalSrcBalance.Cmp(initSrcBalance) != -1 {
			return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var TransferWithApproveFromICON Script = Script{
	Name:        "TransferWithApproveFromICON",
	Description: "Transfer Fixed Amount of coin with approve",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v  Hash %v", err, hash)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		iconTxResult, ok := res.(*icon.TransactionResult)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
		}
		if status, err := iconTxResult.Status.Int(); err != nil {
			return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
		} else if status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
		}

		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}
		time.Sleep(time.Second * 10) // Waiting For Relay
		finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
		// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
		if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
		}
		if finalDstBalance.Cmp(initDstBalance) != 1 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
		}
		if finalSrcBalance.Cmp(initSrcBalance) != -1 {
			return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
		}

		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var TransferWithoutApproveFromHMNY Script = Script{
	Name:        "TransferWithoutApproveFromHMNY",
	Description: "Transfer Fixed Amount of coin without approve",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		txResult, ok := res.(*ethTypes.Receipt)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
		} else if ok && txResult.Status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
		}
		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}
		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}
		time.Sleep(time.Second * 15) // Waiting For Relay
		finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
		// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
		if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
		}
		if finalDstBalance.Cmp(initDstBalance) != 1 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
		}
		if finalSrcBalance.Cmp(initSrcBalance) != -1 {
			return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var TransferWithApproveFromHMNY Script = Script{
	Name:        "TransferWithApproveFromHMNY",
	Description: "Transfer Fixed Amount of coin with approve",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		txResult, ok := res.(*ethTypes.Receipt)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
		} else if ok && txResult.Status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
		}
		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}
		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}
		time.Sleep(time.Second * 15) // Waiting For Relay
		finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
		// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
		if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
		}
		if finalDstBalance.Cmp(initDstBalance) != 1 {
			return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
		}
		if finalSrcBalance.Cmp(initSrcBalance) != -1 {
			return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var MonitorTransferWithoutApproveFromICON Script = Script{
	Name:        "MonitorTransferWithoutApproveFromICON",
	Description: "Transfer Fixed Amount of coin without approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		iconTxResult, ok := res.(*icon.TransactionResult)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
		}
		if status, err := iconTxResult.Status.Int(); err != nil {
			return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
		} else if status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
		}

		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}

		// WaitFor Events
		numExpectedEvents := 2
		args.src.WatchForTransferEnd(args.watchRequestID, startEvent.Sn.Int64())
		args.dst.WatchForTransferReceived(args.watchRequestID, startEvent.Sn.Int64())
		newCtx := context.Background()
		timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)
		for {
			defer timedContextCancel()
			select {
			case <-timedContext.Done():
				return errors.New("Context Timeout Exiting task")
			case <-ctx.Done():
				args.log.Warn("Context Cancelled. Return from Callback watch")
				return nil
			case el := <-args.sinkChan:
				if el.msg.EventType == chain.TransferReceived && el.chainType == args.dst.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				} else if el.msg.EventType == chain.TransferEnd && el.chainType == args.src.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				}
				if numExpectedEvents == 0 {
					finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
					// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
					if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
					}
					if finalDstBalance.Cmp(initDstBalance) != 1 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
					}
					if finalSrcBalance.Cmp(initSrcBalance) != -1 {
						return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
					}
					return nil
				}
			}
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var MonitorTransferWithApproveFromICON Script = Script{
	Name:        "MonitorTransferWithApproveFromICON",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		iconTxResult, ok := res.(*icon.TransactionResult)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
		}
		if status, err := iconTxResult.Status.Int(); err != nil {
			return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
		} else if status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
		}

		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}

		// WaitFor Events
		numExpectedEvents := 2
		args.src.WatchForTransferEnd(args.watchRequestID, startEvent.Sn.Int64())
		args.dst.WatchForTransferReceived(args.watchRequestID, startEvent.Sn.Int64())
		newCtx := context.Background()
		timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)
		for {
			defer timedContextCancel()
			select {
			case <-timedContext.Done():
				return errors.New("Context Timeout Exiting task")
			case <-ctx.Done():
				args.log.Warn("Context Cancelled. Return from Callback watch")
				return nil
			case el := <-args.sinkChan:
				if el.msg.EventType == chain.TransferReceived && el.chainType == args.dst.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				} else if el.msg.EventType == chain.TransferEnd && el.chainType == args.src.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				}
				if numExpectedEvents == 0 {
					finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					//args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
					//args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
					if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
					}
					if finalDstBalance.Cmp(initDstBalance) != 1 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
					}
					if finalSrcBalance.Cmp(initSrcBalance) != -1 {
						return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
					}
					return nil
				}
			}
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var MonitorTransferWithoutApproveFromHMNY Script = Script{
	Name:        "MonitorTransferWithoutApproveFromHMNY",
	Description: "Transfer Fixed Amount of coin without approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		txResult, ok := res.(*ethTypes.Receipt)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
		} else if ok && txResult.Status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
		}
		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}

		// WaitFor Events
		numExpectedEvents := 2
		args.src.WatchForTransferEnd(args.watchRequestID, startEvent.Sn.Int64())
		args.dst.WatchForTransferReceived(args.watchRequestID, startEvent.Sn.Int64())
		newCtx := context.Background()
		timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)
		for {
			defer timedContextCancel()
			select {
			case <-timedContext.Done():
				return errors.New("Context Timeout Exiting task")
			case <-ctx.Done():
				args.log.Warn("Context Cancelled. Return from Callback watch")
				return nil
			case el := <-args.sinkChan:
				if el.msg.EventType == chain.TransferReceived && el.chainType == args.dst.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				} else if el.msg.EventType == chain.TransferEnd && el.chainType == args.src.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				}
				if numExpectedEvents == 0 {
					finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					//args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
					//args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
					if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
					}
					if finalDstBalance.Cmp(initDstBalance) != 1 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
					}
					if finalSrcBalance.Cmp(initSrcBalance) != -1 {
						return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
					}
					return nil
				}
			}
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var MonitorTransferWithApproveFromHMNY Script = Script{
	Name:        "MonitorTransferWithApproveFromHMNY",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		txResult, ok := res.(*ethTypes.Receipt)
		if !ok {
			return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
		} else if ok && txResult.Status != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
		}
		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}

		// WaitFor Events
		numExpectedEvents := 2
		args.src.WatchForTransferEnd(args.watchRequestID, startEvent.Sn.Int64())
		args.dst.WatchForTransferReceived(args.watchRequestID, startEvent.Sn.Int64())
		newCtx := context.Background()
		timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)
		for {
			defer timedContextCancel()
			select {
			case <-timedContext.Done():
				return errors.New("Context Timeout Exiting task")
			case <-ctx.Done():
				args.log.Warn("Context Cancelled. Return from Callback watch")
				return nil
			case el := <-args.sinkChan:
				if el.msg.EventType == chain.TransferReceived && el.chainType == args.dst.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				} else if el.msg.EventType == chain.TransferEnd && el.chainType == args.src.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				}
				if numExpectedEvents == 0 {
					finalSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					finalDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
					if err != nil {
						return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
					}
					//args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
					//args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
					if finalDstBalance.Cmp(startEvent.Assets[0].Value) != 0 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
					}
					if finalDstBalance.Cmp(initDstBalance) != 1 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
					}
					if finalSrcBalance.Cmp(initSrcBalance) != -1 {
						return fmt.Errorf("Balance Compare after Transfer; Src; final Balance should have been less than initial balance; Got Final %v Initial %v", finalSrcBalance.String(), initSrcBalance.String())
					}
					return nil
				}
			}
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var StressTransferWithoutApprove Script = Script{
	Name:        "StressTransferWithoutApprove",
	Description: "Transfer Fixed Amount of coin without approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		if args.src.GetChainType() == chain.ICON {
			iconTxResult, ok := res.(*icon.TransactionResult)
			if !ok {
				return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
			}
			if status, err := iconTxResult.Status.Int(); err != nil {
				return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
			} else if status != 1 {
				return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
			}
		} else if args.src.GetChainType() == chain.HMNY {
			txResult, ok := res.(*ethTypes.Receipt)
			if !ok {
				return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
			} else if ok && txResult.Status != 1 {
				return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
			}
		}
		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			ok := false
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}

		// WaitFor Events
		numExpectedEvents := 2
		args.src.WatchForTransferEnd(args.watchRequestID, startEvent.Sn.Int64())
		args.dst.WatchForTransferReceived(args.watchRequestID, startEvent.Sn.Int64())
		newCtx := context.Background()
		timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)
		for {
			defer timedContextCancel()
			select {
			case <-timedContext.Done():
				return errors.New("Context Timeout Exiting task")
			case <-ctx.Done():
				args.log.Warn("Context Cancelled. Return from Callback watch")
				return nil
			case el := <-args.sinkChan:
				if el.msg.EventType == chain.TransferReceived && el.chainType == args.dst.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				} else if el.msg.EventType == chain.TransferEnd && el.chainType == args.src.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				}
				if numExpectedEvents == 0 {
					// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
					// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
					return nil
				}
			}
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}

var StressTransferWithApprove Script = Script{
	Name:        "StressTransferWithApprove",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 2000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
			return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, elInfo, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		}
		if args.src.GetChainType() == chain.ICON {
			iconTxResult, ok := res.(*icon.TransactionResult)
			if !ok {
				return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
			}
			if status, err := iconTxResult.Status.Int(); err != nil {
				return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
			} else if status != 1 {
				return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
			}
		} else if args.src.GetChainType() == chain.HMNY {
			txResult, ok := res.(*ethTypes.Receipt)
			if !ok {
				return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
			} else if ok && txResult.Status != 1 {
				return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
			}
		}
		if len(elInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		for _, el := range elInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			ok := false
			startEvent, ok = el.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
			}
			srcAddrSplts := strings.Split(args.srcAddr, "/")
			if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
				return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
			} else if args.dstAddr != startEvent.To {
				return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", args.dstAddr, startEvent.To, hash)
			} else if len(startEvent.Assets) == 0 {
				return fmt.Errorf("EventLog; Got zero Asset Details")
			} else if len(startEvent.Assets) > 0 {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
				if startEvent.Assets[0].Name != args.coinName || sum.Cmp(amt) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						args.coinName, amt.String(),
						len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
						hash)
				}
			}
		}
		if !evtFound {
			return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
		}

		// WaitFor Events
		numExpectedEvents := 2
		args.src.WatchForTransferEnd(args.watchRequestID, startEvent.Sn.Int64())
		args.dst.WatchForTransferReceived(args.watchRequestID, startEvent.Sn.Int64())
		newCtx := context.Background()
		timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)
		for {
			defer timedContextCancel()
			select {
			case <-timedContext.Done():
				return errors.New("Context Timeout Exiting task")
			case <-ctx.Done():
				args.log.Warn("Context Cancelled. Return from Callback watch")
				return nil
			case el := <-args.sinkChan:
				if el.msg.EventType == chain.TransferReceived && el.chainType == args.dst.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				} else if el.msg.EventType == chain.TransferEnd && el.chainType == args.src.GetChainType() {
					//args.log.Infof("%+v", el.msg.EventLog)
					numExpectedEvents--
				}
				if numExpectedEvents == 0 {
					// args.log.Infof("Src Balance Init %v Final %v ", initSrcBalance.String(), finalSrcBalance.String())
					// args.log.Infof("Dst Balance Init %v Final %v ", initDstBalance.String(), finalDstBalance.String())
					return nil
				}
			}
		}
		// tests on gas price not included; fee = (intialSrc - finalSrc) > 0
		return nil
	},
}
