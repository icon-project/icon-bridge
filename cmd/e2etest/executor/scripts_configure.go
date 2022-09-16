package executor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

var ConfigureFeeChange ConfigureScript = ConfigureScript{
	Name:        "ChangeFee",
	Description: "Change Fee",
	Type:        "Configure",
	Callback: func(ctx context.Context, conf *configPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		if len(conf.Fee) == 0 {
			errs = IgnoreableError
			return
		}
		stdCfg, err := ts.GetStandardConfigAPI(conf.chainName)
		if err != nil {
			errs = errors.Wrapf(err, "GetStandardConfigAPI %v", err)
			ts.logger.Error(errs)
			return
		}
		for coinName, fee := range conf.Fee {
			setFeeHash, err := stdCfg.SetFeeRatio(ts.GetStandardConfigAPIOwnerKey(conf.chainName), coinName, fee[0], fee[1])
			if err != nil {
				errs = errors.Wrapf(err, "SetFeeRatio %v", err)
				ts.logger.Error(errs)
				return
			}
			if _, errs = ts.ValidateTransactionResult(ctx, conf.chainName, setFeeHash); errs != nil {
				ts.logger.Error("ValidateTransactionResult SetFeeRatio Hash %v Err %v", setFeeHash, errs)
				return
			}
		}
		for coinName, fee := range conf.Fee {
			newFeeNumerator, newFixedFee, err := stdCfg.GetFeeRatio(coinName)
			if err != nil {
				errs = errors.Wrapf(err, "GetFeeRatio %v", err)
				ts.logger.Error(errs)
				return
			}
			if newFeeNumerator.Cmp(fee[0]) != 0 {
				errs = fmt.Errorf("Expected same. Got newFeeNumerator %v feeNumerator %v", newFeeNumerator, fee[0])
				ts.logger.Error(errs)
				return
			}
			if newFixedFee.Cmp(fee[1]) != 0 {
				errs = fmt.Errorf("Expected same. Got newFeeNumerator %v feeNumerator %v", newFixedFee, fee[1])
				ts.logger.Error(errs)
				return
			}
		}
		return
	},
}

var ConfigureTokenLimit ConfigureScript = ConfigureScript{
	Name:        "ConfigureTokenLimit",
	Description: "Configure Token Limit",
	Type:        "Configure",
	Callback: func(ctx context.Context, conf *configPoint, ts *testSuite) (txnRec *txnRecord, errs error) {
		txnRec = &txnRecord{
			feeRecords: []*feeRecord{},
			addresses:  make(map[chain.ChainType][]keypair),
		}
		if len(conf.TokenLimits) == 0 {
			errs = IgnoreableError
			return
		}
		fCfg, err := ts.GetFullConfigAPI()
		if err != nil {
			errs = errors.Wrapf(err, "GetFullConfigAPI %v", err)
			ts.logger.Error(errs)
			return
		}
		coinNames := []string{}
		tokenLimits := []*big.Int{}
		for k, v := range conf.TokenLimits {
			coinNames = append(coinNames, k)
			tokenLimits = append(tokenLimits, v)
		}
		setTokenHash, err := fCfg.SetTokenLimit(ts.FullConfigAPIsOwner(), coinNames, tokenLimits)
		if err != nil {
			errs = errors.Wrapf(err, "setTokenHash %v", err)
			ts.logger.Error(errs)
			return
		}
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), setTokenHash); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResult %v", err)
			ts.logger.Error(errs)
			return
		}
		responseChains := ts.StdConfigAPIChains()
		if len(responseChains) == 0 {
			errs = fmt.Errorf("Expected finite standard ConfigAPI chains.Got zero")
		}
		// TODO: Used a single chain, should wait for response from all chains ?
		err = ts.WaitForConfigResponse(ctx, chain.TokenLimitRequest, chain.TokenLimitResponse, responseChains[0], setTokenHash,
			map[chain.EventLogType]func(event *evt) error{
				chain.TokenLimitRequest: func(ev *evt) error {
					if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
						return errors.New("Got nil value for event ")
					}
					reqEvt, ok := ev.msg.EventLog.(*chain.TokenLimitRequestEvent)
					if !ok {
						return fmt.Errorf("Expected *chain.TokenLimitRequestEvent. Got %T", ev.msg.EventLog)
					}
					txnRec.feeRecords = append(txnRec.feeRecords, &feeRecord{
						ChainName: ts.FullConfigAPIChain(),
						Sn:        reqEvt.Sn,
						Fee:       map[string]*big.Int{},
					})
					if len(reqEvt.CoinNames) != len(conf.TokenLimits) && len(reqEvt.TokenLimits) != len(conf.TokenLimits) {
						return fmt.Errorf("Expected same len reqEvt.CoinNames %v reqEvt.TokenLimits %v conf.TokenLimits %v", len(reqEvt.CoinNames), len(reqEvt.TokenLimits), len(conf.TokenLimits))
					}
					for i := 0; i < len(reqEvt.CoinNames); i++ {
						if confTokenLimit, ok := conf.TokenLimits[reqEvt.CoinNames[i]]; !ok {
							return fmt.Errorf("Unexpected coinName reqEvt.CoinName %v", reqEvt.CoinNames[i])
						} else if ok && confTokenLimit.Cmp(reqEvt.TokenLimits[i]) != 0 {
							return fmt.Errorf("Expected same. Got Different reqEvt.TokenLimit %v conf.TokenLimit %v", reqEvt.TokenLimits[i], confTokenLimit)
						}
					}
					return nil
				},
				chain.TokenLimitResponse: func(ev *evt) error {
					if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
						return errors.New("Got nil value for event ")
					}
					resEvt, ok := ev.msg.EventLog.(*chain.TokenLimitResponseEvent)
					if !ok {
						return fmt.Errorf("Expected *chain.TokenLimitResponseEvent. Got %T", ev.msg.EventLog)
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
		// TODO: Check TokenLimitStatus
		for k, v := range conf.TokenLimits {
			if nv, err := fCfg.GetTokenLimit(k); err != nil {
				errs = errors.Wrapf(err, "GetTokenLimit %v", err)
				ts.logger.Error(errs)
				return
			} else if err == nil && v.Cmp(nv) != 0 {
				errs = errors.Wrapf(err, "Expected same token limit. Got different input %v output %v ", v, nv)
				return
			}
		}
		return
	},
}

func getAccumulatedFees(ts *testSuite) (feePerChain map[chain.ChainType]map[string]*big.Int, errs error) {
	feePerChain = make(map[chain.ChainType]map[string]*big.Int)
	for chainName := range ts.cfgPerChain {
		api, err := ts.GetStandardConfigAPI(chainName)
		if err != nil {
			errs = fmt.Errorf("GetStandardConfigAPI(%v) Err: %v", chainName, err)
			return
		}
		feePerCoin, err := api.GetAccumulatedFees()
		if err != nil {
			errs = fmt.Errorf("GetAccumulatedFees(%v) Err: %v", chainName, err)
			return
		}
		feePerChain[chainName] = feePerCoin
	}
	return
}

func watchFeeGatheringInBackground(ctx context.Context, stopCtx context.Context, ts *testSuite, saveCb func(txn *txnMsg), feeGatheringInterval uint64) (errs error) {
	msg := &txnMsg{
		res: &txnRecord{
			feeRecords: []*feeRecord{},
			addresses: map[chain.ChainType][]keypair{
				chain.BSC: {{"priv", "pub"}},
			},
		},
		err: nil,
	}

	fCfg, err := ts.GetFullConfigAPI()
	if err != nil {
		errs = errors.Wrapf(err, "GetFullConfigAPI %v", err)
		ts.logger.Error(errs)
		return
	}

	// check if feeGatheringInterval Updates
	fCfgOwnerKey := ts.FullConfigAPIsOwner()
	if setFeeHash, err := fCfg.SetFeeGatheringTerm(fCfgOwnerKey, feeGatheringInterval); err != nil {
		errs = errors.Wrapf(err, "SetFeeGatheringTerm %v ", err)
		ts.logger.Error(errs)
		return
	} else {
		if _, err := ts.ValidateTransactionResult(ctx, ts.FullConfigAPIChain(), setFeeHash); err != nil {
			errs = errors.Wrapf(err, "ValidateTransactionResultAndEvents %v", err)
			ts.logger.Error(errs)
			return
		}
		term, err := fCfg.GetFeeGatheringTerm()
		if err != nil {
			errs = errors.Wrapf(err, "GetFeeGatheringTerm %v", err)
			ts.logger.Error(errs)
			return
		}
		if term != feeGatheringInterval {
			errs = fmt.Errorf("Expected same. Got Different GetFeeGatheringTerm(%v) SetFeeGatheringTerm(%v) ", term, feeGatheringInterval)
			ts.logger.Error(errs)
			return
		}
	}
	// TODO:
	responseChains := ts.StdConfigAPIChains()
	if len(responseChains) == 0 {
		errs = fmt.Errorf("Expected finite standard ConfigAPI chains.Got zero")
		ts.logger.Error(errs)
		return
	}
	resMap := map[*big.Int]*txnMsg{}
	go func() {
		msg.err = ts.WaitForFeeGathering(ctx, stopCtx, responseChains[0], map[chain.EventLogType]func(event *evt) error{
			chain.TransferStart: func(ev *evt) error {
				if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
					return errors.New("Got nil value for event ")
				}
				startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
				if !ok {
					return fmt.Errorf("Expected *chain.TransferStartEvent. Got %T", ev.msg.EventLog)
				}
				msgNew := &txnMsg{
					res: &txnRecord{
						feeRecords: []*feeRecord{{
							ChainName: responseChains[0],
							Sn:        startEvt.Sn,
							Fee:       map[string]*big.Int{},
						}},
					},
				}
				for i := 0; i < len(startEvt.Assets); i++ {
					msgNew.res.feeRecords[len(msgNew.res.feeRecords)-1].Fee[startEvt.Assets[i].Name] = startEvt.Assets[i].Value
				}
				resMap[startEvt.Sn] = msgNew
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
				if msg, ok := resMap[endEvt.Sn]; !ok {
					return fmt.Errorf("EndEvt.Sn %v does not exist on map ", endEvt.Sn)
				} else {
					if endEvt.Code.String() != "0" { // remove fee saved if it was unsuccessful response
						ts.logger.Warn("Received Event Code %v Removing saved feeRecords")
						for i := range msg.res.feeRecords {
							msg.res.feeRecords[i].Fee = map[string]*big.Int{}
						}
					}
					saveCb(msg)
					delete(resMap, endEvt.Sn)
				}
				return nil
			},
		})
		ts.logger.Debug("msg err ", msg.err)
		if msg.err != nil && (msg.err.Error() == ExternalContextCancelled.Error() || msg.err.Error() != NilEventReceived.Error()) {
			return
		} else {
			if len(resMap) == 0 {
				saveCb(msg)
				return
			}
			for _, resMsg := range resMap {
				resMsg.err = msg.err
				saveCb(resMsg) // save partial response present in resMap and the error got while waiting for full response
			}
		}
	}()
	return
}
