package executor

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

var Transfer Script = Script{
	Name:        "Transfer",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("1000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		initDstBalance, err := args.dst.GetCoinBalance(args.coinName, args.dstAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Usable.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 1000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if args.coinName != args.src.NativeCoinName() {
			if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
				return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			}
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		} else if res != nil && res.StatusCode != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", res.StatusCode, hash)
		} else if res != nil && len(res.ElInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		startEvent := &chain.TransferStartEvent{}
		tmpOk := false
		for _, el := range res.ElInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, tmpOk = el.EventLog.(*chain.TransferStartEvent)
			if !tmpOk {
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
					if finalDstBalance.Usable.Cmp(startEvent.Assets[0].Value) != 0 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; Got Unequal Destination Balance %v EventAssets Value %v", finalDstBalance.String(), startEvent.Assets[0].Value.String())
					}
					if finalDstBalance.Usable.Cmp(initDstBalance.Usable) != 1 {
						return fmt.Errorf("Balance Compare after Transfer; Dst; final Balance should have been greater than initial balance; Got Final %v Initial %v", finalDstBalance.String(), initDstBalance.String())
					}
					if finalSrcBalance.Usable.Cmp(initSrcBalance.Usable) != -1 {
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

var StressTransfer Script = Script{
	Name:        "StressTransfer",
	Description: "Transfer Fixed Amount of coin with approve and monitor eventlogs TransferReceived and TransferEnd",
	Callback: func(ctx context.Context, args *args) error {
		amt := new(big.Int)
		amt.SetString("1000000000000000000", 10)
		initSrcBalance, err := args.src.GetCoinBalance(args.coinName, args.srcAddr)
		if err != nil {
			return errors.Wrapf(err, "GetCoinBalance Err: %v", err)
		}
		if initSrcBalance.Usable.Cmp(amt) == -1 {
			return fmt.Errorf("Initial Balance %v is less than 1000000000000000000. Expected greater.", initSrcBalance.String())
		}
		if args.coinName != args.src.NativeCoinName() {
			if approveHash, err := args.src.Approve(args.coinName, args.srcKey, *amt); err != nil {
				return errors.Wrapf(err, "Approve Err: %v Hash %v", err, approveHash)
			}
		}
		hash, err := args.src.Transfer(args.coinName, args.srcKey, args.dstAddr, *amt)
		if err != nil {
			return errors.Wrapf(err, "Transfer Err: %v", err)
		}
		res, err := args.src.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			return errors.Wrapf(err, "WaitForTxnResult Coin %v Hash %v", args.coinName, hash)
		} else if res == nil {
			return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
		} else if res != nil && res.StatusCode != 1 {
			return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", res.StatusCode, hash)
		} else if res != nil && len(res.ElInfo) == 0 {
			return fmt.Errorf("WaitForTxnResult; Got zero parsed event logs. Hash %v", hash)
		}

		evtFound := false
		gotEventTypes := []chain.EventLogType{}
		tmpOk := false
		startEvent := &chain.TransferStartEvent{}
		for _, el := range res.ElInfo {
			gotEventTypes = append(gotEventTypes, el.EventType)
			if el.EventType != chain.TransferStart {
				continue
			}
			evtFound = true
			startEvent, tmpOk = el.EventLog.(*chain.TransferStartEvent)
			if !tmpOk {
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
