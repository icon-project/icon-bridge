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
		ex.initflowTransfer()
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
		ex.postProcessBatch(confResponse, transResponse, bgResponse)
	}

	stopJob()

	return nil

}

func (ex *executor) postProcessBatch(confResponse []*configureReq, transResponse []*transferReq, bgResponse []*txnMsg) {

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
	for _, t := range bgResponse {
		if t.err == nil {
			continue
		}
		errs := t.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("SN %v, PID %v, Type BackgroundProcess, Function %v, Input %+v", nonIgnoreErrorCount, "<>", "FeeGathering", "<>")
		}
	}

	for _, t := range transResponse {
		if t.msg.res == nil {
			continue
		}
		for _, fres := range t.msg.res.feeRecords {
			fmt.Printf("FeeResponse %+v\n", fres)
		}
		ex.refund(t.msg.res.addresses)
	}

	////////////////////////////
	fmt.Println("CONFIGURATION++++++++++")
	ignoreErrorCount = 0
	nonIgnoreErrorCount = 0
	for _, res := range confResponse {
		if res.msg.err == nil {
			continue
		}
		errs := res.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			fmt.Println(errs.Error())
			nonIgnoreErrorCount++
		} else {
			ignoreErrorCount++
		}
	}
	fmt.Println("Len IgnoreErr ", ignoreErrorCount)
	fmt.Println("Len Err ", nonIgnoreErrorCount)
	fmt.Println("Len Response ", len(confResponse))
	for _, t := range confResponse {
		if t.msg.res == nil {
			continue
		}
		for _, fres := range t.msg.res.feeRecords {
			fmt.Printf("FeeResponse %+v\n", fres)
		}
		//ex.refund(t.msg.res.addresses)
	}

	////////////////////////////
	fmt.Println("AGGREGATION++++++++++")
	ignoreErrorCount = 0
	nonIgnoreErrorCount = 0
	for _, t := range bgResponse {
		if t.err == nil {
			continue
		}
		errs := t.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			fmt.Println(errs.Error())
			nonIgnoreErrorCount++
		} else {
			ignoreErrorCount++
		}
	}
	fmt.Println("Len IgnoreErr ", ignoreErrorCount)
	fmt.Println("Len Err ", nonIgnoreErrorCount)
	fmt.Println("Len Response ", len(bgResponse))
	for _, t := range bgResponse {
		if t.res == nil {
			continue
		}
		for _, fres := range t.res.feeRecords {
			fmt.Printf("FeeResponse %+v\n", fres)
		}
		//ex.refund(t.msg.res.addresses)
	}
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

func (ex *executor) initflowTransfer() (feePerChain map[chain.ChainType]map[string]*big.Int, errs error) {
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
