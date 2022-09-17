package executor

import (
	"context"
	"fmt"
	"math/big"
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
	&TransferBiDirection,
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

	for _, cpt := range cpts {
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
		ex.postProcessBatch(confResponse, transResponse, initFeeAccumulatedPerChain, initAggregatedFees)
	}
	return nil
}

func (ex *executor) postProcessBatch(confResponse []*configureReq, transResponse []*transferReq, initFeePerChain map[chain.ChainType]map[string]*big.Int, initAggFee map[string]*big.Int) {
	ex.refundGodWithLeftOverAmounts(confResponse, transResponse)
	ex.showErrorMessage(confResponse, transResponse)
	if err := ex.checkWhetherFeeAddsUp(confResponse, transResponse, initFeePerChain, initAggFee); err != nil {
		ex.log.Error("checkWhetherFeeAddsUp %v", err)
	}
}

func (ex *executor) checkWhetherFeeAddsUp(confResponse []*configureReq, transResponse []*transferReq, initAccumulatedFeePerChain map[chain.ChainType]map[string]*big.Int, initAggFee map[string]*big.Int) (err error) {
	finalAccumulatedFeePerChain, err := ex.getFeeAccumulatedPerChain()
	if err != nil {
		return errors.Wrapf(err, "getFeeAccumulatedPerChain %v", err)
	}
	finalAggregatedFee, err := ex.getAggregatedFees()
	if err != nil {
		return errors.Wrapf(err, "getAggregatedFees %v", err)
	}
	fmt.Println("finalAccumulatedFeePerChain")
	for chain, feePerCoin := range finalAccumulatedFeePerChain {
		fmt.Printf("Chain %v \n", chain)
		for coin, fee := range feePerCoin {
			fmt.Printf("%v %v \n", coin, fee)
		}
	}
	fmt.Println("initialAccumulatedFeePerChain")
	for chain, feePerCoin := range initAccumulatedFeePerChain {
		fmt.Printf("Chain %v \n", chain)
		for coin, fee := range feePerCoin {
			fmt.Printf("%v %v \n", coin, fee)
		}
	}
	fmt.Println("initFeeAggregation")
	for coin, amt := range initAggFee {
		fmt.Printf("%v %v \n", coin, amt)
	}
	fmt.Println("finalAggregation")
	for coin, amt := range finalAggregatedFee {
		fmt.Printf("%v %v \n", coin, amt)
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

	fmt.Println("Packets")
	for chain, pkts := range packets {
		fmt.Printf("Chain %v \n", chain)
		for _, pkt := range pkts {
			fmt.Printf("head %v %v %v %v \n", pkt.sn, pkt.isFeeAgg, pkt.id, pkt.err)
			for coin, amt := range pkt.feePerCoin {
				fmt.Printf("%v %v \n", coin, amt)
			}
		}
	}
	totalCalculatedFeePerCoin := map[string]*big.Int{}
	for _, pkts := range packets {
		for _, pkt := range pkts {
			for coinName, pktAmt := range pkt.feePerCoin {
				if _, ok := totalCalculatedFeePerCoin[coinName]; !ok {
					totalCalculatedFeePerCoin[coinName] = big.NewInt(0)
				}
				totalCalculatedFeePerCoin[coinName].Add(totalCalculatedFeePerCoin[coinName], pktAmt) // add charged fee to cumulativeTxnFee
			}
		}
	}
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
		tmpA := (&big.Int{}).Add(iAggFee, iAccFee)
		tmpA = tmpA.Add(tmpA, diff)
		tmpB := (&big.Int{}).Add(fAccFee, fAggFee)
		if tmpA.Cmp(tmpB) != 0 {
			return fmt.Errorf("Expected Same got different. Coin %v initAggFee %v initAccFee %v calculatedFee %v finalAggFee %v finalAccFee %v", coin, iAggFee, iAccFee, diff, fAggFee, fAccFee)
		}
	}
	fmt.Println("Done with feeagg checks")
	return
}

func (ex *executor) showErrorMessage(confResponse []*configureReq, transResponse []*transferReq) {
	nonIgnoreErrorCount := 0
	for _, t := range transResponse {
		if t.msg.err == nil {
			continue
		}
		errs := t.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("ERROR: SNo %v, PID %v, Type Transfer, Function %v, Input %+v", nonIgnoreErrorCount, t.id, t.scr.Name, t.pt)
		}
	}
	for _, t := range transResponse {
		if t.msg.err == nil {
			continue
		}
		errs := t.msg.err
		if _, ignore := ignoreableErrorMap[errs.Error()]; !ignore {
			nonIgnoreErrorCount++
			ex.log.Errorf("ERROR: SNo %v, PID %v, Type Configuration, Function %v, Input %+v", nonIgnoreErrorCount, t.id, t.scr.Name, t.pt)
		}
	}
}

func (ex *executor) refundGodWithLeftOverAmounts(confResponse []*configureReq, transResponse []*transferReq) {
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
		q := &configureReq{scr: script, pt: pt, msg: &txnMsg{}}
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
					q.id = ts.id
					ts.logger.Debugf("%v %v %v %v %v \n", q.id, q.scr.Name, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames)
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
