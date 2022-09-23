package executor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

type aggReq struct {
	id   uint64
	msgs []*txnErrPlusRecord
}

type txnErrPlusRecord struct {
	res *txnRecord
	err error
}

type transferReq struct {
	scr *Script
	pt  *transferPoint
	msg *txnErrPlusRecord
	id  uint64
}

type configureReq struct {
	scr *ConfigureScript
	pt  *configPoint
	msg *txnErrPlusRecord
	id  uint64
}

var transferScripts = []*Script{
	&TransferUniDirection,
	&TransferBiDirection,
	&TransferBatchBiDirection,
	&TransferToBlackListedDstAddress,
	&TransferFromBlackListedSrcAddress,
	&TransferEqualToFee,
	&TransferLessThanFee,
	&TransferToZeroAddress,
	&TransferToUnknownNetwork,
	&TransferWithoutApprove,
}

var configScripts = []*ConfigureScript{
	&ConfigureFeeChange,
	&ConfigureTokenLimit,
}

var ignoreableErrorMap = map[string]*struct{}{
	InsufficientWrappedCoin.Error(): nil,
	InsufficientUnknownCoin.Error(): nil,
	UnsupportedCoinArgs.Error():     nil,
	IgnoreableError.Error():         nil,
}

func (ex *executor) RunFlowTest(ctx context.Context, maxTasks int) error {
	// Generator
	ex.log.Info("Start FlowTest ", maxTasks)
	batchSize := 10 // if unspecified, run 10 transfer tasks
	if maxTasks > 0 && maxTasks < 50 {
		batchSize = maxTasks
	} else if maxTasks > 50 {
		batchSize = 50
	}
	tg := pointGenerator{
		cfgPerChain:  ex.cfgPerChain,
		clsPerChain:  ex.clientsPerChain,
		maxBatchSize: &batchSize,
		transferFilter: func(tp []*transferPoint) []*transferPoint {
			ntp := []*transferPoint{}
			for _, v := range tp {
				if true {
					ntp = append(ntp, v)
				}
			}
			return ntp
		},
		configFilter: nil,
	}

	cpts, err := tg.GenerateConfigPoints()
	if err != nil {
		return errors.Wrapf(err, "GenerateConfigPoints %v", err)
	}

	ts, err := ex.getTestSuite()
	if err != nil {
		return errors.Wrapf(err, "getTestSuite %v", err)
	}
	defer ex.removeChan(ts.id)

	bgJobRecords := []*txnErrPlusRecord{}
	tmu := sync.Mutex{}
	extractRecords := func() (ret []*txnErrPlusRecord) {
		tmu.Lock()
		defer tmu.Unlock()
		ret = []*txnErrPlusRecord{}
		for _, v := range bgJobRecords {
			ret = append(ret, v)
		}
		bgJobRecords = []*txnErrPlusRecord{}
		return
	}

	appendToRecords := func(rec *txnErrPlusRecord) {
		tmu.Lock()
		defer tmu.Unlock()
		bgJobRecords = append(bgJobRecords, rec)
	}

	stopJob, waitForPendingBGMsg, err := ex.startBackgroundJob(ctx, ts, appendToRecords)
	if err != nil {
		return err
	}
	defer stopJob()

	ex.log.Info("Number of Batch to Process: ", len(cpts))
	for cpti, cpt := range cpts {
		ex.log.Info("Processing Batch ", cpti+1)
		initFeeAccumulatedPerChain, err := ex.getFeeAccumulatedPerChain()
		if err != nil {
			return errors.Wrapf(err, "getFeeAccumulatedPerChain %v", err)
		}
		initAggregatedFees, err := ex.getAggregatedFees()
		if err != nil {
			return errors.Wrapf(err, "getAggregatedFees %v", err)
		}
		ex.log.Info("Setting up Configuration for this batch ", cpti+1)
		confResponse, err := ex.processConfigurePoint(ctx, cpt)
		if err != nil {
			return errors.Wrapf(err, "processConfigurePoint %v", err)
		}
		tpts, err := tg.GenerateTransferPoints()
		if err != nil {
			return errors.Wrapf(err, "GenerateTransferPoints %v", err)
		}
		ex.log.Infof("Number of Tasks in current batch: %v. Processing...\n", len(tpts))
		transResponse, err := ex.processTransferPoints(ctx, tpts)
		if err != nil {
			return errors.Wrapf(err, "processTransferPoints %v", err)
		}
		ex.log.Info("Processed all tasks in current batch")
		waitForPendingBGMsg()
		bgResponse := extractRecords()
		ex.log.Info("Post Process current batch")
		ex.postProcessBatch(confResponse, transResponse, initFeeAccumulatedPerChain, initAggregatedFees, &aggReq{msgs: bgResponse, id: ts.id})
	}
	ex.log.Info("Completed Flow Test")
	return nil
}

func (ex *executor) postProcessBatch(confResponse []*configureReq, transResponse []*transferReq, initFeePerChain map[chain.ChainType]map[string]*big.Int, initAggFee map[string]*big.Int, bgRecords *aggReq) {
	ex.refundGodWithLeftOverAmounts(confResponse, transResponse)
	ex.showErrorMessage(confResponse, transResponse)
	if ex.enableExperimentalFeatures {
		ex.log.Info("Using Experimental Feature")
		if err := ex.checkWhetherFeeAddsUp(confResponse, transResponse, initFeePerChain, initAggFee, bgRecords); err != nil {
			ex.log.Error("checkWhetherFeeAddsUp %v", err)
		}
	}
}

func (ex *executor) checkWhetherFeeAddsUp(confResponse []*configureReq, transResponse []*transferReq, initAccumulatedFeePerChain map[chain.ChainType]map[string]*big.Int, initAggFee map[string]*big.Int, aggRequest *aggReq) (err error) {
	finalAccumulatedFeePerChain, err := ex.getFeeAccumulatedPerChain()
	if err != nil {
		return errors.Wrapf(err, "getFeeAccumulatedPerChain %v", err)
	}
	finalAggregatedFee, err := ex.getAggregatedFees()
	if err != nil {
		return errors.Wrapf(err, "getAggregatedFees %v", err)
	}
	packets := prepareFeePacket(confResponse, transResponse, aggRequest)

	totalCalculatedFeePerCoin := map[string]*big.Int{}
	for _, pkts := range packets {
		for _, pkt := range pkts {
			for coinName, pktAmt := range pkt.feePerCoin {
				_, ok := totalCalculatedFeePerCoin[coinName]
				if !ok {
					totalCalculatedFeePerCoin[coinName] = big.NewInt(0)
				}
				if !pkt.isFeeAgg {
					totalCalculatedFeePerCoin[coinName].Add(totalCalculatedFeePerCoin[coinName], pktAmt)
					continue
				}
			}
		}
	}
	/*
		for chainName, pkts := range packets {
			sort.SliceStable(pkts, func(i, j int) bool {
				return pkts[i].sn < pkts[j].sn
			})
			initAccFeePerCoin, ok := initAccumulatedFeePerChain[chainName]
			if !ok {
				return fmt.Errorf("initAccumulatedFeePerChain does not include fee for chain %v", chainName)
			}
			finalAccFeePerCoin, ok := finalAccumulatedFeePerChain[chainName]
			if !ok {
				return fmt.Errorf("finalAccumulatedFeePerChain does not include fee for chain %v", chainName)
			}
			cumulativeCalculatedFeePerCoin := map[string]*big.Int{}
			subTractedAmountPerCoin := map[string]*big.Int{}
			for k := range initAccFeePerCoin {
				cumulativeCalculatedFeePerCoin[k] = big.NewInt(0)
			}
			for _, pkt := range pkts {
				fmt.Printf("head %v %v %v %v \n", pkt.sn, pkt.isFeeAgg, pkt.pid, pkt.err)
				for coinName, pktAmt := range pkt.feePerCoin {
					fmt.Printf("%v %v \n", coinName, pktAmt)
					cumulativeCalculatedFee, ok := cumulativeCalculatedFeePerCoin[coinName]
					if !ok {
						return fmt.Errorf("cumulativeCalculatedFeePerChainPerCoin does not include fee for coin %v on chain %v", coinName, chainName)
					}
					if !pkt.isFeeAgg {
						cumulativeCalculatedFee.Add(cumulativeCalculatedFee, pktAmt)
						continue
					}
					if _, ok := subTractedAmountPerCoin[coinName]; !ok {
						subTractedAmountPerCoin[coinName] = big.NewInt(0)
					}
					if pkt.err == nil {
						if pktAmt.Cmp((&big.Int{}).Sub((&big.Int{}).Add(cumulativeCalculatedFee, initAccFeePerCoin[coinName]), subTractedAmountPerCoin[coinName])) != 0 {
							return fmt.Errorf("Expected same. Got Different Sn %v Chain %v Coin %v cumulativeCalculatedFee  %v initAccFeePerCoin %v, AggTxFee %v SumAggTxFee %v", pkt.sn, chainName, coinName, cumulativeCalculatedFee, initAccFeePerCoin[coinName], pktAmt, subTractedAmountPerCoin[coinName])
						}
						subTractedAmountPerCoin[coinName].Add(subTractedAmountPerCoin[coinName], pktAmt)
					} else {
						// TODO
					}
				}
			}
			for coin, value := range subTractedAmountPerCoin {
				if finalAccFeePerCoin[coin].Cmp((&big.Int{}).Sub((&big.Int{}).Add(cumulativeCalculatedFeePerCoin[coin], initAccFeePerCoin[coin]), value)) != 0 {
					return fmt.Errorf("Expected same. Got Different Chain %v Coin %v finalAcc  %v AggTxFee %v SumAggTxFee %v", chainName, coin, finalAccFeePerCoin[coin], initAccFeePerCoin[coin], cumulativeCalculatedFeePerCoin[coin])
				}
			}
			for coin, cFee := range cumulativeCalculatedFeePerCoin {
				if _, ok := totalCalculatedFeePerCoin[coin]; !ok {
					totalCalculatedFeePerCoin[coin] = big.NewInt(0)
				}
				totalCalculatedFeePerCoin[coin] = (&big.Int{}).Add(totalCalculatedFeePerCoin[coin], cFee)
			}
		}
	*/
	totalInitialAccumulatedFeePerCoin := map[string]*big.Int{}
	totalFinalAccumulatedFeePerCoin := map[string]*big.Int{}
	for _, fpc := range initAccumulatedFeePerChain {
		for c, f := range fpc {
			if _, ok := totalInitialAccumulatedFeePerCoin[c]; !ok {
				totalInitialAccumulatedFeePerCoin[c] = big.NewInt(0)
			}
			totalInitialAccumulatedFeePerCoin[c] = totalInitialAccumulatedFeePerCoin[c].Add(totalInitialAccumulatedFeePerCoin[c], f)
		}
	}
	for _, fpc := range finalAccumulatedFeePerChain {
		for c, f := range fpc {
			if _, ok := totalFinalAccumulatedFeePerCoin[c]; !ok {
				totalFinalAccumulatedFeePerCoin[c] = big.NewInt(0)
			}
			totalFinalAccumulatedFeePerCoin[c] = totalFinalAccumulatedFeePerCoin[c].Add(totalFinalAccumulatedFeePerCoin[c], f)
		}
	}
	for coin, iAggFee := range initAggFee {
		fAggFee := finalAggregatedFee[coin]
		iAccFee := totalInitialAccumulatedFeePerCoin[coin]
		fAccFee := totalFinalAccumulatedFeePerCoin[coin]
		diff := totalCalculatedFeePerCoin[coin]
		if iAccFee == nil || iAggFee == nil || fAccFee == nil || fAggFee == nil || diff == nil {
			ex.log.Debug(fmt.Errorf("Got nil. Coin %v initAggFee %v initAccFee %v calculatedFee %v finalAggFee %v finalAccFee %v", coin, iAggFee, iAccFee, diff, fAggFee, fAccFee))
			continue
		}
		tmpA := (&big.Int{}).Add(iAggFee, iAccFee)
		tmpA = tmpA.Add(tmpA, diff)
		tmpB := (&big.Int{}).Add(fAccFee, fAggFee)
		if tmpA.Cmp(tmpB) != 0 {
			return fmt.Errorf("Expected Same got different. Coin %v initAggFee %v initAccFee %v calculatedFee %v finalAggFee %v finalAccFee %v", coin, iAggFee, iAccFee, diff, fAggFee, fAccFee)
		}
	}
	return
}

func (ex *executor) showErrorMessage(confResponse []*configureReq, transResponse []*transferReq) {
	ex.log.Info("Showing Error Message If Any")
	nonIgnoreErrorCount := 0
	for _, t := range confResponse {
		if t.msg.err == nil {
			continue
		}
		errs := t.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("ERROR: SNo %v, PID %v, Type Configuration, Function %v, Input %+v Err %+v", nonIgnoreErrorCount, t.id, t.scr.Name, t.pt, errs)
		}
	}
	for _, t := range transResponse {
		if t.msg.err == nil {
			continue
		}
		errs := t.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("ERROR: SNo %v, PID %v, Type Transfer, Function %v, Input %+v Err %+v", nonIgnoreErrorCount, t.id, t.scr.Name, t.pt, errs)
		}
	}
}

func (ex *executor) refundGodWithLeftOverAmounts(confResponse []*configureReq, transResponse []*transferReq) {
	ex.log.Info("Reclaiming leftover native coins if any")
	for _, c := range confResponse {
		if c.msg.err != nil || c.msg.res == nil {
			continue
		}
		ex.refund(c.msg.res.addresses)
	}
	for _, t := range transResponse {
		if t.msg.err != nil || t.msg.res == nil {
			continue
		}
		ex.refund(t.msg.res.addresses)
	}
	return
}

func (ex *executor) processConfigurePoint(ctx context.Context, pt *configPoint) (totalResponse []*configureReq, errs error) {
	totalResponse = []*configureReq{}
	for _, script := range configScripts {
		q := &configureReq{scr: script, pt: pt, msg: &txnErrPlusRecord{}}
		ts, err := ex.getTestSuite()
		if err != nil {
			q.msg.err = err
		}
		defer ex.removeChan(ts.id)
		q.id = ts.id
		ts.logger.Debugf("%v %v %v \n", q.id, q.scr.Name, q.pt.Fee)
		q.msg.res, q.msg.err = script.Callback(ctx, pt, ts)
		totalResponse = append(totalResponse, q)
	}
	return
}

func (ex *executor) processTransferPoints(ctx context.Context, tpts []*transferPoint) (totalResponse []*transferReq, err error) {
	numDivisions := (len(transferScripts) * len(tpts)) / 10
	if numDivisions == 0 {
		numDivisions = 1
	}
	requests := make([]*transferReq, len(transferScripts)*len(tpts))
	tmpi := 0
	for _, tpt := range tpts {
		for _, scr := range transferScripts {
			requests[tmpi] = &transferReq{scr: scr, pt: tpt}
			tmpi++
		}
	}

	concurrency := len(requests)
	if concurrency <= 1 {
		concurrency = 2
	}
	totalResponse = []*transferReq{}
	reqCursor := 0
	for reqCursor < len(requests) {
		rqch := make(chan *transferReq, concurrency)
		for i := reqCursor; len(rqch) < cap(rqch) && i < len(requests); i++ {
			rqch <- requests[i]
			reqCursor++
		}
		sres := make([]*transferReq, 0, len(rqch))
		for q := range rqch {
			switch {
			case q.msg != nil:
				sres = append(sres, q)
				if (len(sres)+len(totalResponse))%numDivisions == 0 {
					ex.log.Infof("Status: Processed Tasks %v out of %v \n", len(totalResponse)+len(sres), len(requests))
				}
				if len(sres) == cap(sres) {
					close(rqch)
				}
			default:
				//time.Sleep(time.Second * 5)
				go func(q *transferReq) {
					defer func() {
						time.Sleep(time.Millisecond * 100)
						rqch <- q
					}()
					q.msg = &txnErrPlusRecord{}
					if q.scr.Callback == nil {
						q.msg.err = errors.New("Callback nil")
						return
					}
					ts, err := ex.getTestSuite()
					if err != nil {
						q.msg.err = err
					}
					defer ex.removeChan(ts.id)
					q.id = ts.id
					ts.logger.Debugf("%v %v %v %v %v \n", q.id, q.scr.Name, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames)
					q.msg.res, q.msg.err = q.scr.Callback(ctx, q.pt, ts)
				}(q)

			}
		}
		totalResponse = append(totalResponse, sres...)
	}
	return
}

func (ex *executor) refund(addrMap map[chain.ChainType][]keypair) {
	for chainName, addrs := range addrMap {
		cl, ok := ex.clientsPerChain[chainName]
		if !ok {
			ex.log.Warn(fmt.Errorf("Client %v does not exist ", chainName))
			continue
		}
		for _, addr := range addrs {
			bal, err := cl.GetCoinBalance(cl.NativeCoin(), cl.GetBTPAddress(addr.PubKey))
			if err != nil {
				ex.log.Warn(errors.Wrapf(err, "GetCoinBalance %v", err))
				continue
			}
			gasLimitOnSrc := big.NewInt(int64(ex.cfgPerChain[chainName].GasLimit[chain.TransferNativeCoinIntraChainGasLimit]))
			gasFeeOnSrc := (&big.Int{}).Mul(cl.SuggestGasPrice(), gasLimitOnSrc)

			if bal.UserBalance.Cmp((&big.Int{}).Mul(gasFeeOnSrc, big.NewInt(1))) > 0 {
				transferrableAmount := (&big.Int{}).Sub(bal.UserBalance, (&big.Int{}).Mul(gasFeeOnSrc, big.NewInt(1)))
				_, err := cl.Transfer(cl.NativeCoin(), addr.PrivKey, cl.GetBTPAddress(ex.godKeysPerChain[chainName].PubKey), transferrableAmount)
				if err != nil {
					ex.log.Warn(errors.Wrapf(err, "Transfer %v", err))
					continue
				}
				ex.log.Debugf("Redeem %v %v", chainName, transferrableAmount)
			}
			// else {
			// 	ex.log.Infof("Insufficient to redeem addr %v gasFeeOnSrc %v UserBalance %v", addr, gasFeeOnSrc, bal.UserBalance)
			// }
		}
	}
}

func (ex *executor) getFeeAccumulatedPerChain() (feePerChain map[chain.ChainType]map[string]*big.Int, errs error) {
	ts, err := ex.getTestSuite()
	if err != nil {
		errs = fmt.Errorf("getTestSuite %v", err)
		return
	}
	defer ex.removeChan(ts.id)
	feePerChain, err = getAccumulatedFees(ts)
	if err != nil {
		errs = fmt.Errorf("getAccumulatedFees %v", err)
	}
	return
}

func (ex *executor) getAggregatedFees() (aggregatedAmountPerCoin map[string]*big.Int, errs error) {
	ts, err := ex.getTestSuite()
	if err != nil {
		errs = fmt.Errorf("getTestSuite %v", err)
		return
	}
	defer ex.removeChan(ts.id)
	aggregatedAmountPerCoin = map[string]*big.Int{}
	cl, ok := ex.clientsPerChain[ts.FullConfigAPIChain()]
	if !ok {
		errs = fmt.Errorf("Expected chain %v not found in clientsPerChain", ts.FullConfigAPIChain())
		return
	}
	cfg, ok := ex.cfgPerChain[ts.FullConfigAPIChain()]
	if !ok {
		errs = fmt.Errorf("Expected chain %v not found in cfgPerChain", ts.FullConfigAPIChain())
		return
	}
	for _, coinName := range append(append(cfg.NativeTokens, cfg.NativeCoin), cfg.WrappedCoins...) {
		bal, err := cl.GetCoinBalance(coinName, cl.GetBTPAddress(ts.feeAggregatorAddress))
		if err != nil {
			errs = fmt.Errorf("GetCoinBalance %v", err)
			return
		}
		aggregatedAmountPerCoin[coinName] = bal.UserBalance
	}
	return
}

func (ex *executor) startBackgroundJob(ctx context.Context, ts *testSuite, cb func(txn *txnErrPlusRecord)) (context.CancelFunc, func(), error) {
	newCtx, newCancel := context.WithCancel(context.Background())
	waitForPendingMsg, err := watchFeeGatheringInBackground(ctx, newCtx, ts, cb, 90)
	if err != nil {
		newCancel()
		return nil, nil, errors.Wrapf(err, "watchFeeGatheringInBackground %v", err)
	}
	return newCancel, waitForPendingMsg, nil
}

type feePacket struct {
	isFeeAgg   bool
	sn         uint64
	feePerCoin map[string]*big.Int
	pid        uint64
	err        error
}

func prepareFeePacket(confResponse []*configureReq, transResponse []*transferReq, aggRequest *aggReq) (packets map[chain.ChainType][]*feePacket) {
	packets = map[chain.ChainType][]*feePacket{}
	for _, c := range confResponse {
		if c.msg.res == nil {
			continue
		}
		for _, cf := range c.msg.res.feeRecords {
			packets[cf.ChainName] = append(packets[cf.ChainName], &feePacket{sn: cf.Sn.Uint64(), pid: c.id, err: c.msg.err, feePerCoin: cf.Fee})
		}
	}
	for _, c := range transResponse {
		if c.msg.res == nil {
			continue
		}
		for _, cf := range c.msg.res.feeRecords { // if c.err != nil, then it's probably on the last cf, TODO: handle it properly
			packets[cf.ChainName] = append(packets[cf.ChainName], &feePacket{sn: cf.Sn.Uint64(), pid: c.id, err: c.msg.err, feePerCoin: cf.Fee})
		}
	}

	for _, c := range aggRequest.msgs {
		if c.res == nil {
			continue
		}
		for _, cf := range c.res.feeRecords {
			packets[cf.ChainName] = append(packets[cf.ChainName], &feePacket{isFeeAgg: true, sn: cf.Sn.Uint64(), pid: aggRequest.id, err: c.err, feePerCoin: cf.Fee})
		}
	}
	return
}
