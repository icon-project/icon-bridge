package executor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
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
	&TransferToBlackListedDstAddress,
}

var configScripts = []*ConfigureScript{
	&ConfigureFeeChange,
	&ConfigureTokenLimit,
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

	bts, err := ex.getTestSuite()
	if err != nil {
		return err
	}
	defer ex.removeChan(bts.id)

	// Callback
	bgJobRecords := []*txnMsg{}
	tmu := sync.RWMutex{}
	extractRecords := func() (ret []*txnMsg) {
		tmu.Lock()
		defer tmu.Unlock()
		ret = bgJobRecords
		bgJobRecords = []*txnMsg{}
		return
	}
	appendToRecords := func(rec *txnMsg) {
		tmu.Lock()
		defer tmu.Unlock()
		bgJobRecords = append(bgJobRecords, rec)
	}

	stopJob, err := ex.startBackgroundJob(ctx, bts, appendToRecords)
	if err != nil {
		return err
	}

	// Interate
	for _, cpt := range cpts {
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
	if stopJob != nil {
		fmt.Println("StopJobs ")
		stopJob()
	}

	return nil

}

func (ex *executor) postProcessBatch(confResponse []*configureReq, transResponse []*transferReq, bgResponse []*txnMsg) {
	/*
		ignoreableError := map[string]*struct{}{
			InsufficientWrappedCoin.Error(): nil,
			InsufficientUnknownCoin.Error(): nil,
			UnsupportedCoinArgs.Error():     nil,
			IgnoreableError.Error():         nil,
		}
		ignoreErrorCount := 0
		nonIgnoreErrorCount := 0
		for _, err := range batchError {
			if _, ignore := ignoreableError[err.Error()]; !ignore {
				fmt.Println(err.Error())
				nonIgnoreErrorCount++
			} else {
				ignoreErrorCount++
			}
		}
		fmt.Println("Len IgnoreErr ", ignoreErrorCount)
		fmt.Println("Len Err ", nonIgnoreErrorCount)
		fmt.Println("Len Response ", batchResponse)
		for _, res := range batchResponse {
			for _, fres := range res.feeRecords {
				fmt.Printf("FeeResponse %+v\n", fres)
			}
			ex.refund(res.addresses)
		}
	*/
}

func (ex *executor) processConfigurePoint(ctx context.Context, pt *configPoint) (totalResponse []*configureReq, errs error) {
	totalResponse = []*configureReq{}
	for _, script := range configScripts {
		q := &configureReq{scr: script, pt: pt}
		ts, err := ex.getTestSuite()
		if err != nil {
			q.msg.err = err
		}
		defer ex.removeChan(ts.id)
		q.id = ts.id
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

					ts.logger.Debug("%v %v %v %v \n", q.scr.Name, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames)
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
	err := watchFeeGatheringInBackground(ctx, newCtx, ts, cb, 150)
	if err != nil {
		fmt.Println("Got Error ", err)
		newCancel()
		return nil, err
	}
	fmt.Println("return channel")
	return newCancel, nil
}

func (ex *executor) getTestSuite() (ts *testSuite, err error) {
	id, err := ex.getID()
	if err != nil {
		return nil, errors.Wrapf(err, "getID %v", err)
	}
	log := ex.log.WithFields(log.Fields{"pid": id})
	sinkChan := make(chan *evt)
	ex.addChan(id, sinkChan)

	ts = &testSuite{
		id:                   id,
		logger:               log,
		subChan:              sinkChan,
		clsPerChain:          ex.clientsPerChain,
		godKeysPerChain:      ex.godKeysPerChain,
		cfgPerChain:          ex.cfgPerChain,
		feeAggregatorAddress: ex.feeAggregatorAddress,
	}

	return ts, nil
}
