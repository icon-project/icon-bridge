package executor

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type keypair struct {
	PrivKey string
	PubKey  string
}

type testSuite struct {
	id              uint64
	clsPerChain     map[chain.ChainType]chain.ChainAPI
	godKeysPerChain map[chain.ChainType]keypair
	logger          log.Logger
	src             chain.ChainType
	dst             chain.ChainType
	subChan         <-chan *evt
	report          string
}

func (ts *testSuite) GetChainPair(srcChain, dstChain chain.ChainType) (src chain.SrcAPI, dst chain.DstAPI, err error) {
	ok := false
	src, ok = ts.clsPerChain[srcChain]
	if !ok {
		err = fmt.Errorf("Chain %v not found", srcChain)
		return
	}
	dst, ok = ts.clsPerChain[dstChain]
	if !ok {
		err = fmt.Errorf("Chain %v not found", dstChain)
	}
	ts.src = srcChain
	ts.dst = dstChain
	return
}

func (ts *testSuite) GetKeyPairs(chainName chain.ChainType) (key, addr string, err error) {
	cl, ok := ts.clsPerChain[chainName]
	if !ok {
		err = fmt.Errorf("Chain %v not found", chainName)
		return
	}
	keyPairs, err := cl.GetKeyPairs(1)
	if err != nil {
		err = errors.Wrapf(err, "GetKeyPairs %v", err)
		return
	}
	key = keyPairs[0][0]
	addr = cl.GetBTPAddress(keyPairs[0][1])
	return
}

func (ts *testSuite) GetGodKeyPairs(chainName chain.ChainType) (key, addr string, err error) {
	godkeyPair, ok := ts.godKeysPerChain[chainName]
	if !ok {
		err = errors.Wrapf(err, "GetKeyPairs %v", err)
		return
	}
	cl, ok := ts.clsPerChain[chainName]
	if !ok {
		err = fmt.Errorf("Chain %v not found", chainName)
		return
	}
	key = godkeyPair.PrivKey
	addr = cl.GetBTPAddress(godkeyPair.PubKey)
	return
}

func (ts *testSuite) Fund(addr string, amount *big.Int, coinName string) error {
	// If coin is a native, intrachain-tranfer,
	// else if it's wrapped inter-chain, else not found
	srcCl, ok := ts.clsPerChain[ts.src]
	if !ok {
		return fmt.Errorf("Chain %v not found", ts.src)
	}
	srcKeys, ok := ts.godKeysPerChain[ts.src]
	if !ok {
		return fmt.Errorf("GodKeys %v not found", ts.src)
	}

	for _, scoin := range append(srcCl.NativeTokens(), srcCl.NativeCoin()) {
		if scoin != coinName {
			continue
		}
		hash, err := srcCl.Transfer(coinName, srcKeys.PrivKey, addr, *amount)
		if err != nil {
			return errors.Wrapf(err, "srcCl.Transfer err=%v", err)
		}
		_, err = ts.ValidateTransactionResult(context.TODO(), hash)
		return err
	}
	dstCl, ok := ts.clsPerChain[ts.dst]
	if !ok {
		return fmt.Errorf("Chain %v not found", ts.dst)
	}
	dstKeys, ok := ts.godKeysPerChain[ts.dst]
	if !ok {
		return fmt.Errorf("GodKeys %v not found", ts.dst)
	}

	for _, dcoin := range append(dstCl.NativeTokens(), dstCl.NativeCoin()) {
		if dcoin != coinName {
			continue
		}
		if coinName != dstCl.NativeCoin() {
			if _, err := dstCl.Approve(coinName, dstKeys.PrivKey, *amount); err != nil {
				return errors.Wrapf(err, "dstCl.Approve %v", err)
			}
		}
		hash, err := dstCl.Transfer(coinName, dstKeys.PrivKey, addr, *amount)
		if err != nil {
			return errors.Wrapf(err, "dstCl.Transfer err=%v", err)
		}
		_, err = ts.ValidateTransactionResult(context.TODO(), hash)
		return err
	}
	return fmt.Errorf("coinName %v not found amongst coins used by chains %v and %v", coinName, ts.src, ts.dst)
}

func (ts *testSuite) ValidateTransactionResult(ctx context.Context, hash string) (res *chain.TxnResult, err error) {
	srcCl, ok := ts.clsPerChain[ts.src]
	if !ok {
		err = fmt.Errorf("Chain %v not found", ts.src)
		return
	}
	tctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	res, err = srcCl.WaitForTxnResult(tctx, hash)
	if err != nil {
		err = errors.Wrapf(err, "WaitForTxnResult Hash %v", hash)
	} else if res == nil {
		err = fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
	} else if res != nil && res.StatusCode != 1 {
		err = errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", res.StatusCode, hash)
	}
	return
}

func (ts *testSuite) ValidateTransactionResultEvents(ctx context.Context, hash, coinName, srcAddr, dstAddr string, amt *big.Int) (err error) {
	srcCl, ok := ts.clsPerChain[ts.src]
	if !ok {
		return fmt.Errorf("Chain %v not found", ts.src)
	}
	tctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	res, err := srcCl.WaitForTxnResult(tctx, hash)
	if err != nil {
		return errors.Wrapf(err, "WaitForTxnResult Hash %v", hash)
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
		srcAddrSplts := strings.Split(srcAddr, "/")
		if srcAddrSplts[len(srcAddrSplts)-1] != startEvent.From {
			return fmt.Errorf("EventLog; Expected Source Address %v Got %v Hash %v", srcAddrSplts[len(srcAddrSplts)-1], startEvent.From, hash)
		} else if dstAddr != startEvent.To {
			return fmt.Errorf("EventLog; Expected Destination Address %v Got %v Hash %v", dstAddr, startEvent.To, hash)
		} else if len(startEvent.Assets) == 0 {
			return fmt.Errorf("EventLog; Got zero Asset Details")
		} else if len(startEvent.Assets) > 0 {
			sum := big.NewInt(0)
			sum.Add(startEvent.Assets[0].Value, startEvent.Assets[0].Fee)
			if startEvent.Assets[0].Name != coinName || sum.Cmp(amt) != 0 {
				return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
					coinName, amt.String(),
					len(startEvent.Assets), startEvent.Assets[0].Name, startEvent.Assets[0].Value.String(), startEvent.Assets[0].Fee.String(),
					hash)
			}
		}
	}
	if !evtFound {
		return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
	}
	return
}

func (ts *testSuite) WaitForEvents(ctx context.Context, hash string, cbPerEvent map[chain.EventLogType]func(event *evt) error) (err error) {
	res, err := ts.ValidateTransactionResult(ctx, hash)
	startEvent := &chain.TransferStartEvent{}
	tmpOk := false
	for _, el := range res.ElInfo {
		if el.EventType != chain.TransferStart {
			continue
		}
		startEvent, tmpOk = el.EventLog.(*chain.TransferStartEvent)
		if !tmpOk {
			return fmt.Errorf("EventLog; Execpted *chain.TransferStartEvent. Got %T Hash %v", el.EventLog, hash)
		}
		if startCb, ok := cbPerEvent[chain.TransferStart]; ok {
			if err := startCb(&evt{chainType: ts.src, msg: el}); err != nil {
				ts.report += fmt.Sprintf("CallBackPerEvent %v Err:%v \n", "TransferStart", err)
			}
		}
		break
	}
	if !tmpOk { // if no start event, tmpOk is not set
		return fmt.Errorf("TransferStart event not found in txn result Hash=%v", hash)
	}

	// Register WatchEvents
	srcCl, ok := ts.clsPerChain[ts.src]
	if !ok {
		err = fmt.Errorf("Client for chain %v not found", ts.src)
		return
	}
	dstCl, ok := ts.clsPerChain[ts.dst]
	if !ok {
		err = fmt.Errorf("Client for chain %v not found", ts.dst)
		return
	}
	numExpectedEvents := 0
	for ev := range cbPerEvent {
		if ev == chain.TransferStart {
			// Trasfer Start event is not watched as it is premise for other watches and as such
			// has already been known to be true. i.e. startEvent got above and callback called if given
		} else if ev == chain.TransferReceived {
			if err := dstCl.WatchForTransferReceived(ts.id, startEvent.Sn.Int64()); err != nil {
				return errors.Wrapf(err, "WatchForTransferStart Err=%v", err)
			}
			numExpectedEvents++
		} else if ev == chain.TransferEnd {
			if err := srcCl.WatchForTransferEnd(ts.id, startEvent.Sn.Int64()); err != nil {
				return errors.Wrapf(err, "WatchForTransferStart Err=%v", err)
			}
			numExpectedEvents++
		} else {
			ts.report += fmt.Sprintf("Event %v not available. Skipping it.", ev)
		}
	}

	// Listen to result from watchEvents
	newCtx := context.Background()
	timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*60)

	for {
		defer timedContextCancel()
		select {
		case <-timedContext.Done():
			ts.report += "Context Timeout Exiting task"
			return errors.New("Context Timeout Exiting task")
		case <-ctx.Done():
			ts.report += "Context Cancelled. Return from Callback watch"
			return errors.New("Context Cancelled. Return from Callback watch")
		case ev := <-ts.subChan:
			if cb, ok := cbPerEvent[ev.msg.EventType]; ok {
				if (ev.msg.EventType == chain.TransferReceived && ev.chainType == dstCl.GetChainType()) ||
					(ev.msg.EventType == chain.TransferEnd && ev.chainType == srcCl.GetChainType()) {
					if err := cb(ev); err != nil {
						ts.report += fmt.Sprintf("CallBackPerEvent %v Err:%v \n", ev.msg.EventType, err)
					}
					numExpectedEvents--
				}
			}
			if numExpectedEvents == 0 {
				ts.report += "All events found. Exiting \n"
				return
			}
		}
	}
	return nil
}
