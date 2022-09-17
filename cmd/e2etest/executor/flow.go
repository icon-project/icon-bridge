package executor

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
)

type txnMsg struct {
	res *txnRecord
	err error
}

type transferReq struct {
	scr *Script
	pt  *transferPoint
	msg *txnMsg
	id  uint64
}

type configureReq struct {
	scr *ConfigureScript
	pt  *configPoint
	msg *txnMsg
	id  uint64
}

type aggReq struct {
	id   uint64
	msgs []*txnMsg
}

var transferScripts = []*Script{
	&TransferUniDirection,
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

func (ex *executor) RunFlowTest(ctx context.Context) error {
	// Generator
	tg := pointGenerator{
		cfgPerChain:  ex.cfgPerChain,
		maxBatchSize: nil,
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

	bgJobRecords := []*txnMsg{}
	tmu := sync.Mutex{}
	extractRecords := func() (ret []*txnMsg) {
		tmu.Lock()
		defer tmu.Unlock()
		ret = []*txnMsg{}
		for _, v := range bgJobRecords {
			ret = append(ret, v)
		}
		fmt.Printf("Extract from records %v", ret)
		bgJobRecords = []*txnMsg{}
		return
	}

	appendToRecords := func(rec *txnMsg) {
		tmu.Lock()
		defer tmu.Unlock()
		fmt.Printf("Append to records %v", rec)
		bgJobRecords = append(bgJobRecords, rec)
	}

	stopJob, err := ex.startBackgroundJob(ctx, ts, appendToRecords)
	if err != nil {
		return err
	}
	defer stopJob()

	for _, cpt := range cpts {
		// TODO: start with clean slate
		initFeeAccumulatedPerChain, err := ex.getFeeAccumulatedPerChain()
		if err != nil {
			return errors.Wrapf(err, "getFeeAccumulatedPerChain %v", err)
		}
		initAggregatedFees, err := ex.getAggregatedFees()
		if err != nil {
			return errors.Wrapf(err, "getAggregatedFees %v", err)
		}
		confResponse, err := ex.processConfigurePoint(ctx, cpt)
		if err != nil {
			return errors.Wrapf(err, "processConfigurePoint %v", err)
		}
		tpts, err := tg.GenerateTransferPoints(cpt)
		if err != nil {
			return errors.Wrapf(err, "GenerateTransferPoints %v", err)
		}
		transResponse, err := ex.processTransferPoints(ctx, tpts)
		if err != nil {
			return errors.Wrapf(err, "processTransferPoints %v", err)
		}
		bgResponse := extractRecords()
		ex.postProcessBatch(confResponse, transResponse, &aggReq{id: ts.id, msgs: bgResponse}, initFeeAccumulatedPerChain, initAggregatedFees)
	}
	stopJob()
	return nil
}

func (ex *executor) postProcessBatch(confResponse []*configureReq, transResponse []*transferReq, bgResponse *aggReq, initFeePerChain map[chain.ChainType]map[string]*big.Int, initAggFee map[string]*big.Int) {
	ex.refundGodWithLeftOverAmounts(confResponse, transResponse)
	ex.showErrorMessage(confResponse, transResponse, bgResponse)
	ex.checkWhetherSnAligns(confResponse, transResponse, bgResponse)
	ex.checkWhetherFeeAddsUp(confResponse, transResponse, bgResponse, initFeePerChain, initAggFee)

}

func (ex *executor) checkWhetherFeeAddsUp(confResponse []*configureReq, transResponse []*transferReq, bgResponse *aggReq, initFeePerChain map[chain.ChainType]map[string]*big.Int, initAggFee map[string]*big.Int) (err error) {
	finalFeePerChain, err := ex.getFeeAccumulatedPerChain()
	if err != nil {
		return errors.Wrapf(err, "getFeeAccumulatedPerChain %v", err)
	}
	finalAggFee, err := ex.getAggregatedFees()
	if err != nil {
		return errors.Wrapf(err, "getAggregatedFees %v", err)
	}
	type feePacket struct {
		isFeeAgg   bool
		sn         uint64
		feePerCoin map[string]*big.Int
		id         uint16
		err        error
	}
	packets := map[chain.ChainType][]*feePacket{}
	for _, c := range confResponse {
		if c.msg.res == nil {
			continue
		}
		for _, cf := range c.msg.res.feeRecords {
			packets[cf.ChainName] = append(packets[cf.ChainName], &feePacket{sn: cf.Sn.Uint64(), id: uint16(c.id), err: c.msg.err, feePerCoin: cf.Fee})
		}
	}
	for _, c := range transResponse {
		if c.msg.res == nil {
			continue
		}
		for _, cf := range c.msg.res.feeRecords { // if c.err != nil, then it's probably on the last cf, TODO: handle it properly
			packets[cf.ChainName] = append(packets[cf.ChainName], &feePacket{sn: cf.Sn.Uint64(), id: uint16(c.id), err: c.msg.err, feePerCoin: cf.Fee})
		}
	}
	for _, c := range bgResponse.msgs {
		if c.res == nil {
			continue
		}
		for _, cf := range c.res.feeRecords {
			packets[cf.ChainName] = append(packets[cf.ChainName], &feePacket{isFeeAgg: true, sn: cf.Sn.Uint64(), id: uint16(bgResponse.id), err: c.err, feePerCoin: cf.Fee})
		}
	}
	totalFeePerCoin := map[string]*big.Int{}
	for chainName, pkts := range packets {
		sort.SliceStable(pkts, func(i, j int) bool {
			return pkts[i].sn < pkts[j].sn
		})
		initFeePerCoin, ok := initFeePerChain[chainName]
		if !ok {
			return fmt.Errorf("InitFeePerChain does not include fee for chain %v", chainName)
		}
		finalFeePerCoin, ok := finalFeePerChain[chainName]
		if !ok {
			return fmt.Errorf("FinalFeePerChain does not include fee for chain %v", chainName)
		}
		cumulativeFeePerCoin := map[string]*big.Int{}
		subTractedAmountPerCoin := map[string]*big.Int{}
		for k, v := range initFeePerCoin {
			cumulativeFeePerCoin[k] = (&big.Int{}).Set(v)
			subTractedAmountPerCoin[k] = big.NewInt(0)
		}
		for _, pkt := range pkts {
			for coinName, pktAmt := range pkt.feePerCoin {
				cumulativeFee, ok := cumulativeFeePerCoin[coinName]
				if !ok {
					return fmt.Errorf("cumulativeFeePerChainPerCoin does not include fee for coin %v on chain %v", coinName, chainName)
				}
				if !pkt.isFeeAgg {
					cumulativeFee.Add(cumulativeFee, pktAmt)
					continue
				}
				if pkt.err == nil {
					if pktAmt.Cmp((&big.Int{}).Sub(cumulativeFee, subTractedAmountPerCoin[coinName])) != 0 {
						return fmt.Errorf("Expected same. Got Different Sn %v cumulativeFee  %v AggTxFee %v", pkt.sn, cumulativeFee, pktAmt)
					}
					subTractedAmountPerCoin[coinName].Add(subTractedAmountPerCoin[coinName], pktAmt)
				} else {
					// TODO
				}
			}
		}
		for coin, value := range finalFeePerCoin {
			if value.Cmp((&big.Int{}).Sub(cumulativeFeePerCoin[coin], subTractedAmountPerCoin[coin])) != 0 {
				return fmt.Errorf("Expected Same Got Different finalAccumulatedFee(%v,%v) %v CululativeFee %v SubAmount %v", chainName, coin, value, cumulativeFeePerCoin[coin], subTractedAmountPerCoin[coin])
			}
		}
		for coin, cFee := range cumulativeFeePerCoin {
			if _, ok := totalFeePerCoin[coin]; !ok {
				totalFeePerCoin[coin] = big.NewInt(0)
			}
			totalFeePerCoin[coin] = (&big.Int{}).Add(totalFeePerCoin[coin], cFee)
		}
	}
	for coin, iAggFee := range initAggFee {
		fAggFee := finalAggFee[coin]
		diff := totalFeePerCoin[coin]
		if ((&big.Int{}).Add(iAggFee, diff)).Cmp(fAggFee) != 0 {
			return fmt.Errorf("Expected Same got different. Coin %v initAggFee %v cumulativeCalculatedFee %v finalAggFee %v", coin, iAggFee, diff, fAggFee)
		}
	}
	return
}

func (ex *executor) showErrorMessage(confResponse []*configureReq, transResponse []*transferReq, bgResponse *aggReq) {
	nonIgnoreErrorCount := 0
	for _, t := range transResponse {
		if t.msg.err == nil {
			continue
		}
		errs := t.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("SN %v, PID %v, Type Transfer, Function %v, Input %+v", nonIgnoreErrorCount, t.id, t.scr.Name, t.pt)
		}
	}
	for _, t := range transResponse {
		if t.msg.err == nil {
			continue
		}
		errs := t.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("SN %v, PID %v, Type Configuration, Function %v, Input %+v", nonIgnoreErrorCount, t.id, t.scr.Name, t.pt)
		}
	}
	for _, t := range bgResponse.msgs {
		if t.err == nil {
			continue
		}
		errs := t.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("SN %v, PID %v, Type BackgroundProcess, Function %v, Input %+v", nonIgnoreErrorCount, bgResponse.id, "FeeGathering", "<>")
		}
	}
}

func (ex *executor) refundGodWithLeftOverAmounts(confResponse []*configureReq, transResponse []*transferReq) {
	for _, c := range confResponse {
		if c.msg.err != nil {
			continue
		}
		ex.refund(c.msg.res.addresses)
	}
	for _, t := range transResponse {
		if t.msg.err != nil {
			continue
		}
		ex.refund(t.msg.res.addresses)
	}
	return
}

func (ex *executor) checkWhetherSnAligns(confResponse []*configureReq, transResponse []*transferReq, bgResponse *aggReq) (errs error) {
	snID := map[chain.ChainType][][2]uint64{}
	for _, c := range confResponse {
		if c.msg.res == nil {
			continue
		}
		for _, cf := range c.msg.res.feeRecords {
			snID[cf.ChainName] = append(snID[cf.ChainName], [2]uint64{c.id, cf.Sn.Uint64()})
		}
	}
	for _, t := range transResponse {
		if t.msg.res == nil {
			continue
		}
		for _, tf := range t.msg.res.feeRecords {
			snID[tf.ChainName] = append(snID[tf.ChainName], [2]uint64{t.id, tf.Sn.Uint64()})
		}
	}
	for _, b := range bgResponse.msgs {
		if b.res == nil {
			continue
		}
		for _, bf := range b.res.feeRecords {
			snID[bf.ChainName] = append(snID[bf.ChainName], [2]uint64{bgResponse.id, bf.Sn.Uint64()})
		}
	}

	for cName, v := range snID {
		sort.SliceStable(v, func(i, j int) bool {
			return v[i][1] < v[j][1]
		})
		for i := 1; i < len(v); i++ {
			if v[i][1]-v[i-1][1] != 1 {
				errs = fmt.Errorf("Not in serialized order Chain %v, Pairs(id,sn) (%v) and (%v)", cName, v[i], v[i-1])
				return
			}
		}
	}
	return
}

func (ex *executor) processConfigurePoint(ctx context.Context, pt *configPoint) (totalResponse []*configureReq, errs error) {
	totalResponse = []*configureReq{}
	for _, script := range configScripts {
		q := &configureReq{scr: script, pt: pt, msg: &txnMsg{}}
		ts, err := ex.getTestSuite()
		if err != nil {
			q.msg.err = err
		}
		defer ex.removeChan(ts.id)
		q.id = ts.id
		ts.logger.Debugf("%v %v \n", q.scr.Name, q.pt.Fee)
		q.msg.res, q.msg.err = script.Callback(ctx, pt, ts)
		totalResponse = append(totalResponse, q)
	}
	return
}

func (ex *executor) processTransferPoints(ctx context.Context, tpts []*transferPoint) (totalResponse []*transferReq, err error) {
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
				if len(sres) == cap(sres) {
					close(rqch)
				}
			default:
				time.Sleep(time.Second * 5)
				go func(q *transferReq) {
					defer func() {
						time.Sleep(time.Millisecond * 100)
						rqch <- q
					}()
					q.msg = &txnMsg{}
					if q.scr.Callback == nil {
						q.msg.err = errors.New("Callback nil")
						return
					}
					ts, err := ex.getTestSuite()
					if err != nil {
						q.msg.err = err
					}
					defer ex.removeChan(ts.id)

					ts.logger.Debugf("%v %v %v %v \n", q.scr.Name, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames)
					q.msg.res, q.msg.err = q.scr.Callback(ctx, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames, ts)
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
				ex.log.Infof("Redeem %v %v", chainName, transferrableAmount)
			}
			// else {
			// 	ex.log.Infof("Insufficient to redeem addr %v gasFeeOnSrc %v UserBalance %v", addr, gasFeeOnSrc, bal.UserBalance)
			// }
		}
	}
}

func (ex *executor) startBackgroundJob(ctx context.Context, ts *testSuite, cb func(txn *txnMsg)) (context.CancelFunc, error) {

	newCtx, newCancel := context.WithCancel(context.Background())
	fmt.Println("watchFeeGatheringInBackground")
	err := watchFeeGatheringInBackground(ctx, newCtx, ts, cb, 90)
	if err != nil {
		fmt.Println("Got Error ", err)
		newCancel()
		return nil, err
	}
	return newCancel, nil
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
	for _, cd := range cfg.CoinDetails {
		bal, err := cl.GetCoinBalance(cd.Name, cl.GetBTPAddress(ts.feeAggregatorAddress))
		if err != nil {
			errs = fmt.Errorf("GetCoinBalance %v", err)
			return
		}
		aggregatedAmountPerCoin[cd.Name] = bal.UserBalance
	}
	return
}
