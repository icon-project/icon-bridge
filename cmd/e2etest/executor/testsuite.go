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

const DENOMINATOR = 10000
const MAX_UINT256 = "115792089237316195423570985008687907853269984665640564039457584007913129639935"

type testSuite struct {
	id      uint64
	logger  log.Logger
	subChan <-chan *evt

	report string

	clsPerChain     map[chain.ChainType]chain.ChainAPI
	godKeysPerChain map[chain.ChainType]keypair
	cfgPerChain     map[chain.ChainType]*chain.Config
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
	return
}

func (ts *testSuite) GetFullConfigAPI() (cfgSrc chain.FullConfigAPI, err error) {
	cfgSrc, ok := ts.clsPerChain[chain.ICON]
	if !ok {
		err = fmt.Errorf("Chain %v not found", chain.ICON)
	}
	return
}

func (ts *testSuite) GetStandardConfigAPI(chainName chain.ChainType) (cfgAPI chain.StandardConfigAPI, err error) {
	cfgAPI, ok := ts.clsPerChain[chainName]
	if !ok {
		err = fmt.Errorf("Chain %v not found", chainName)
	}
	return
}

func (ts *testSuite) GetStandardConfigAPIOwnerKey(dstChain chain.ChainType) string {
	return ts.godKeysPerChain[dstChain].PrivKey
}

func (ts *testSuite) FullConfigAPIsOwner() string {
	return ts.godKeysPerChain[chain.ICON].PrivKey
}

func (ts *testSuite) FullConfigAPIChain() chain.ChainType {
	return chain.ICON
}

func (ts *testSuite) NetAddr(btpAddr string) (net string, addr string) {
	splts := strings.Split(btpAddr, "/")
	return splts[len(splts)-2], splts[len(splts)-1]
}

func (ts *testSuite) getAmountBeforeFeeCharge(chainName chain.ChainType, coinName string, outputBalance *big.Int) (*big.Int, error) {
	/*
		What is the input amount that we must have so that the net transferrable amount
		after fee charged on chainName is equal to outputBalance for coinName ?
		feeCharged = inputBalance * ratio + fixedFee
		outputBalance = inputBalance - feeCharged
		inputBalance = (outputBalance + fixed) / (1 - ratio)
		inputBalance = (outputBalance + fixed) * deniminator / (denominator - numerator)

	*/
	coinDetails := ts.cfgPerChain[chainName].CoinDetails
	for i := 0; i < len(coinDetails); i++ {
		if coinDetails[i].Name == coinName {
			fixedFee, _ := (&big.Int{}).SetString(coinDetails[i].FixedFee, 10)
			bplusf := (&big.Int{}).Add(outputBalance, fixedFee)
			bplusf.Mul(bplusf, big.NewInt(DENOMINATOR))
			dminusn := new(big.Int).Sub(big.NewInt(DENOMINATOR), big.NewInt(int64(coinDetails[i].FeeNumerator)))
			bplusf.Div(bplusf, dminusn)
			return bplusf, nil
		}
	}
	return nil, fmt.Errorf("Coin %v Not Found in coinDetails", coinName)
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

func (ts *testSuite) Fund(chainName chain.ChainType, addr string, amount *big.Int, coinName string) error {
	// IntraChain Transfer of Tokens from God to an address
	srcCl, ok := ts.clsPerChain[chainName]
	if !ok {
		return fmt.Errorf("Chain %v not found", chainName)
	}
	godKey, ok := ts.godKeysPerChain[chainName]
	if !ok {
		return fmt.Errorf("GodKeys %v not found", chainName)
	}
	if strings.Contains(addr, godKey.PubKey) {
		return nil // Sender == Receiver; so skip
	}
	ts.logger.Infof("Fund coin %v addr %v amt %v ", coinName, addr, amount.String())
	hash, err := srcCl.Transfer(coinName, godKey.PrivKey, addr, amount)
	if err != nil {
		return errors.Wrapf(err, "srcCl.Transfer err=%v", err)
	}
	_, err = ts.ValidateTransactionResult(context.TODO(), chainName, hash)
	return err
}

func (ts *testSuite) ValidateTransactionResult(ctx context.Context, chainName chain.ChainType, hash string) (res *chain.TxnResult, err error) {
	srcCl, ok := ts.clsPerChain[chainName]
	if !ok {
		err = fmt.Errorf("Chain %v not found", chainName)
		return
	}
	tctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	res, err = srcCl.WaitForTxnResult(tctx, hash)
	if err != nil {
		err = errors.Wrapf(err, "WaitForTxnResult(%v) Err %v", hash, err)
	} else if res == nil {
		err = fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
	} else if res != nil && res.StatusCode != 1 {
		err = errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", res.StatusCode, hash)
		err = StatusCodeZero
		return
	}
	return
}

func (ts *testSuite) WaitForConfigResponse(ctx context.Context, reqEvent, responseEvent chain.EventLogType, responseChainName chain.ChainType, reqHash string, cbPerEvent map[chain.EventLogType]func(event *evt) error) (err error) {
	if responseChainName == ts.FullConfigAPIChain() {
		return nil
	}
	fCfg, err := ts.GetFullConfigAPI()
	if err != nil {
		return errors.Wrapf(err, "GetFullConfigAPI %v", err)
	}
	reqEvtInfo, err := fCfg.GetConfigRequestEvent(reqEvent, reqHash)
	if err != nil {
		return errors.Wrapf(err, "%v %v", reqEvent, err)
	}
	if cb, ok := cbPerEvent[reqEvent]; ok && cb != nil {
		if err = cb(&evt{msg: reqEvtInfo, chainType: chain.ICON}); err != nil {
			return errors.Wrapf(err, "CallBack(%v) %v", reqEvent, err)
		}
	}
	// RegisterWatchEvents
	dstCfg, err := ts.GetStandardConfigAPI(responseChainName)
	if err != nil {
		return errors.Wrapf(err, "GetStandardConfigAPI %v", err)
	}
	if reqEvent == chain.AddToBlacklistRequest && responseEvent == chain.BlacklistResponse {
		res, ok := reqEvtInfo.EventLog.(*chain.AddToBlacklistRequestEvent)
		if !ok {
			return errors.Wrapf(err, "Expected *chain.AddToBlacklistRequestEvent Got %T", reqEvtInfo.EventLog)
		}
		if err := dstCfg.WatchForBlacklistResponse(ts.id, res.Sn.Int64()); err != nil {
			return errors.Wrapf(err, "WatchForBlacklistResponse %v", err)
		}
	} else if reqEvent == chain.RemoveFromBlacklistRequest && responseEvent == chain.BlacklistResponse {
		res, ok := reqEvtInfo.EventLog.(*chain.RemoveFromBlacklistRequestEvent)
		if !ok {
			return errors.Wrapf(err, "Expected *chain.RemoveFromBlacklistRequestEvent Got %T", reqEvtInfo.EventLog)
		}
		if err := dstCfg.WatchForBlacklistResponse(ts.id, res.Sn.Int64()); err != nil {
			return errors.Wrapf(err, "WatchForBlacklistResponse %v", err)
		}
	} else if reqEvent == chain.TokenLimitRequest && responseEvent == chain.TokenLimitResponse {
		res, ok := reqEvtInfo.EventLog.(*chain.TokenLimitRequestEvent)
		if !ok {
			return errors.Wrapf(err, "Expected *chain.TokenLimitRequestEvent Got %T", reqEvtInfo.EventLog)
		}
		if err := dstCfg.WatchForSetTokenLmitResponse(ts.id, res.Sn.Int64()); err != nil {
			return errors.Wrapf(err, "WatchForSetTokenLmitResponse %v", err)
		}
	}

	// Listen to result from watchEvents
	newCtx := context.Background()
	timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*180)

	for {
		defer timedContextCancel()
		select {
		case <-timedContext.Done():
			ts.report += "Context Timeout Exiting task"
			return errors.New("Context Timeout Exiting task----------------")
		case <-ctx.Done():
			ts.report += "Context Cancelled. Return from Callback watch"
			return errors.New("Context Cancelled. Return from Callback watch---------------")
		case ev := <-ts.subChan:
			if cb, ok := cbPerEvent[ev.msg.EventType]; ok && ev.msg.EventType == responseEvent {
				if cb != nil {
					if err := cb(ev); err != nil {
						ts.report += fmt.Sprintf("CallBackPerEvent %v Err:%v \n", ev.msg.EventType, err)
						return errors.Wrapf(err, "Callback (%v) %v", responseEvent, err)
					}
				}
				ts.report += "All events found. Exiting \n"
				return
			}
		}
	}
	return nil
}

func (ts *testSuite) ValidateTransactionResultAndEvents(ctx context.Context, chainName chain.ChainType, hash string, coinNames []string, srcAddr, dstAddr string, amts []*big.Int) (err error) {
	srcCl, ok := ts.clsPerChain[chainName]
	if !ok {
		return fmt.Errorf("Chain %v not found", chainName)
	}
	tctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	res, err := srcCl.WaitForTxnResult(tctx, hash)
	if err != nil {
		return errors.Wrapf(err, "WaitForTxnResult Hash %v", hash)
	} else if res == nil {
		return fmt.Errorf("WaitForTxnResult; Transaction Result is nil. Hash %v", hash)
	} else if res != nil && res.StatusCode != 1 {
		err = errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", res.StatusCode, hash)
		return StatusCodeZero
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
			for i := 0; i < len(coinNames); i++ {
				sum := big.NewInt(0)
				sum.Add(startEvent.Assets[i].Value, startEvent.Assets[i].Fee)
				if startEvent.Assets[i].Name != coinNames[i] || sum.Cmp(amts[i]) != 0 {
					return fmt.Errorf("EventLog; Expected Name %v, Amount %v Got Len(assets) %v Name %v Value %v Fee %v. Hash %v",
						coinNames[i], amts[i].String(),
						len(startEvent.Assets), startEvent.Assets[i].Name, startEvent.Assets[i].Value.String(), startEvent.Assets[i].Fee.String(),
						hash)
				}
			}
		}
	}
	if !evtFound {
		return fmt.Errorf("Transfer Start Event Not Found. Got %v Hash %v", gotEventTypes, hash)
	}
	return
}

func (ts *testSuite) WaitForEvents(ctx context.Context, srcChainName, dstChainName chain.ChainType, hash string, cbPerEvent map[chain.EventLogType]func(event *evt) error) (err error) {
	res, err := ts.ValidateTransactionResult(ctx, srcChainName, hash)
	if err != nil {
		return
	}
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
			if err := startCb(&evt{chainType: srcChainName, msg: el}); err != nil {
				return err
				//ts.report += fmt.Sprintf("CallBackPerEvent %v Err:%v \n", "TransferStart", err)
			}
		}
		break
	}
	if !tmpOk { // if no start event, tmpOk is not set
		return fmt.Errorf("TransferStart event not found in txn result Hash=%v", hash)
	}

	// Register WatchEvents
	srcCl, ok := ts.clsPerChain[srcChainName]
	if !ok {
		err = fmt.Errorf("Client for chain %v not found", srcChainName)
		return
	}
	dstCl, ok := ts.clsPerChain[dstChainName]
	if !ok {
		err = fmt.Errorf("Client for chain %v not found", dstChainName)
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
	timedContext, timedContextCancel := context.WithTimeout(newCtx, time.Second*180)

	for {
		defer timedContextCancel()
		select {
		case <-timedContext.Done():
			ts.report += "Context Timeout Exiting task"
			return errors.New("Context Timeout Exiting task----------------")
		case <-ctx.Done():
			ts.report += "Context Cancelled. Return from Callback watch"
			return errors.New("Context Cancelled. Return from Callback watch---------------")
		case ev := <-ts.subChan:
			if cb, ok := cbPerEvent[ev.msg.EventType]; ok {
				numExpectedEvents--
				if cb != nil {
					if err := cb(ev); err != nil {
						return err
						ts.report += fmt.Sprintf("CallBackPerEvent %v Err:%v \n", ev.msg.EventType, err)
					}
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
