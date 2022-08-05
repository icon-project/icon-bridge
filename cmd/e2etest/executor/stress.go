package executor

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

func modifyAddress(addr string) (string, error) {
	modifiers := []func(string) (string, error){
		func(inputStr string) (string, error) {
			//fmt.Println("zero address")
			splits := strings.Split(inputStr, "/")
			if len(splits) != 4 {
				return "", errors.New("Unexpected length")
			}
			lenSplits := len(splits)
			if len(splits[lenSplits-1]) > 2 {
				prefix := splits[lenSplits-1][0:2]
				postFix := splits[lenSplits-1][2:]
				byteArr := make([]byte, len(postFix)/2)
				for i := 0; i < len(postFix)/2; i++ {
					byteArr[i] = 0
				}
				splits[lenSplits-1] = prefix + hex.EncodeToString(byteArr)
			}
			return strings.Join(splits, "/"), nil
		},
		func(inputStr string) (string, error) {
			//fmt.Println("wrong net")
			splits := strings.Split(inputStr, "/")
			if len(splits) != 4 {
				return "", errors.New("Unexpected length")
			}
			network := splits[2]
			networkSplits := strings.Split(network, ".")
			if len(networkSplits) != 2 {
				return "", errors.New("Unexpected length")
			}
			networkSplits[1] += "s"
			splits[2] = strings.Join(networkSplits, ".")
			return strings.Join(splits, "/"), nil
		},
		func(inputStr string) (string, error) {
			return inputStr + "1", nil
		},
	}
	i := rand.Intn(len(modifiers) * 10) // equal chance of working and non-working address
	if i < len(modifiers) {
		return modifiers[i](addr)
	}
	// fmt.Println("correct address")
	return addr, nil
}

func stressTransferInterChain(
	ctx context.Context,
	srcChain, dstChain chain.ChainType,
	srcKeyPair, dstKeyPair keypair,
	coinNames []string,
	ts *testSuite) (response *txnRecord, err error) {
	response = &txnRecord{}
	if len(coinNames) == 0 {
		err = errors.New("Should specify at least one coinname, got zero")
		return
	}
	if srcChain == dstChain {
		response.msg = "Expected different chains"
		return
	}
	src, dst, err := ts.GetChainPair(srcChain, dstChain)
	if err != nil {
		err = errors.Wrapf(err, "GetChainPair %v", err)
		return
	}
	srcKey := srcKeyPair.PrivKey
	srcAddr := src.GetBTPAddress(srcKeyPair.PubKey)
	dstAddr := dst.GetBTPAddress(dstKeyPair.PubKey)

	dstAddr, err = modifyAddress(dstAddr)
	if err != nil {
		err = errors.Wrapf(err, "modifyAddress %v", err)
		return
	}
	getAmountToTransfer := func(fee *big.Int, _coin string) *big.Int {
		checkIfCommonTokenExists := func() bool {
			exists := false
			for _, stkn := range src.NativeTokens() {
				if stkn == _coin {
					for _, dtkn := range dst.NativeTokens() {
						if dtkn == _coin {
							exists = true
							break
						}
					}
					break
				}
			}
			return exists
		}
		modifiers := []func(*big.Int) *big.Int{
			func(_fee *big.Int) *big.Int {
				// tx=> fee+1, rx=>1, Pass
				_fee.Mul(_fee, big.NewInt(10))
				return _fee.Add(_fee, big.NewInt(1))
			},
			func(_fee *big.Int) *big.Int {
				// tx=>fee-1, Fail
				return _fee.Sub(_fee, big.NewInt(1))
			},
			func(_fee *big.Int) *big.Int {
				if checkIfCommonTokenExists() {
					if btsBalance, err := dst.GetCoinBalance(_coin, ts.btsAddressPerChain[dstChain]); err == nil {
						// tx=>fee+bts, rx=>bts, Pass
						return _fee.Add(_fee, btsBalance.UserBalance)
					}
				}
				// tx => fee+1, rx => 1, Pass
				return _fee.Add(_fee, big.NewInt(1))
			},
			func(_fee *big.Int) *big.Int {
				// tx => fee, rx => 0, Fail
				return _fee
			},
			func(_fee *big.Int) *big.Int {
				if checkIfCommonTokenExists() {
					if btsBalance, err := dst.GetCoinBalance(_coin, ts.btsAddressPerChain[dstChain]); err == nil {
						// tx=> fee+bts+1, rx=>bts+1, Fail
						_fee.Add(_fee, big.NewInt(1))
						return _fee.Add(_fee, btsBalance.UserBalance)
					}
				}
				// tx=> fee-1, Fail
				return _fee.Sub(_fee, big.NewInt(1))
			},
		}
		i := rand.Intn(len(modifiers) * 10)
		if i < len(modifiers) {
			return modifiers[i](fee)
		}
		addedFee := fee.Add(fee, big.NewInt(1))
		return addedFee.Mul(addedFee, big.NewInt(3))
	}

	amts := make([]*big.Int, len(coinNames))
	for i := 0; i < len(coinNames); i++ {
		amts[i] = getAmountToTransfer(ts.withFeeAdded(big.NewInt(0)), coinNames[i])
	}
	//ts.logger.Infof("Src %+v Amts %+v dstaddr %+v coin %v", srcAddr, amts[0], dstAddr, coinNames[0])
	// Approve
	if rand.Int()%10 != 0 {
		var approveHash string
		for i, coinName := range coinNames {
			if coinName != src.NativeCoin() {
				if approveHash, err = src.Approve(coinName, srcKey, amts[i]); err != nil {
					response.msg = fmt.Sprintf("Approve Err: %v Hash %v", err, approveHash)
					return
				} else {
					if _, err = ts.ValidateTransactionResult(ctx, approveHash); err != nil {
						response.msg = fmt.Sprintf("Approve ValidateTransactionResult Err: %v Hash %v", err, approveHash)
						return
					}
				}
			}
		}
	}
	//ts.logger.Info("Transfer Now")
	var hash string
	if len(coinNames) == 1 {
		hash, err = src.Transfer(coinNames[0], srcKey, dstAddr, amts[0])
		if err != nil {
			response.msg = fmt.Sprintf("Transfer Err: %v", err)
			return
		}
	} else {
		hash, err = src.TransferBatch(coinNames, srcKey, dstAddr, amts)
		if err != nil {
			response.msg = fmt.Sprintf("TransferBatch Err: %v", err)
			return
		}
	}
	//ts.logger.Info("ValidateEvents Now")
	if err = ts.ValidateTransactionResultAndEvents(ctx, hash, coinNames, srcAddr, dstAddr, amts); err != nil {
		response.msg = fmt.Sprintf("ValidateTransactionResultAndEvents Unexpected error %v", err)
		return
	}
	//ts.logger.Info("WaitForEvents Now")
	err = ts.WaitForEvents(ctx, hash, map[chain.EventLogType]func(*evt) error{
		chain.TransferStart: func(ev *evt) error {
			if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
				return errors.New("Got nil value for event ")
			}
			startEvt, ok := ev.msg.EventLog.(*chain.TransferStartEvent)
			if !ok {
				return fmt.Errorf("Expected *chain.TransferStartEvent Got %T", ev.msg.EventLog)
			}
			response.startEvent = startEvt
			return nil
		},
		chain.TransferEnd: func(ev *evt) error {
			if ev == nil || (ev != nil && ev.msg == nil) || (ev != nil && ev.msg != nil && ev.msg.EventLog == nil) {
				return errors.New("Got nil value for event ")
			}
			endEvt, ok := ev.msg.EventLog.(*chain.TransferEndEvent)
			if !ok {
				return fmt.Errorf("Expected *chain.TransferEndEvent Got %T", ev.msg.EventLog)
			}
			response.endEvent = endEvt
			return nil
		},
	})
	if err != nil {
		err = errors.Wrapf(err, "WaitForEvents %v", err)
		return
	}
	return
}

func stressTransferIntraChain(
	ctx context.Context,
	srcChain, dstChain chain.ChainType,
	srcKeyPair, dstKeyPair keypair,
	coinNames []string,
	ts *testSuite) (response *txnRecord, err error) {
	response = &txnRecord{}
	if len(coinNames) == 0 {
		err = errors.New("Should specify only one coinname, got zero")
		return
	}
	if srcChain != dstChain {
		response.msg = "Expected same chain"
		return
	}

	if srcKeyPair.PubKey == dstKeyPair.PubKey {
		response.msg = "Expected different addresses"
		return
	}

	cl, _, err := ts.GetChainPair(srcChain, dstChain)
	if err != nil {
		err = errors.Wrapf(err, "GetChainPair %v", err)
		return
	}

	srcKey := srcKeyPair.PrivKey
	srcAddr := cl.GetBTPAddress(srcKeyPair.PubKey)
	dstAddr := cl.GetBTPAddress(dstKeyPair.PubKey)

	amt := ts.withFeeAdded(big.NewInt(1))
	amt.Mul(amt, big.NewInt(9))
	var hash string
	hash, err = cl.Transfer(coinNames[0], srcKey, dstAddr, amt)
	if err != nil {
		response.msg = fmt.Sprintf("Transfer Err: %v", err)
		return
	}
	if _, err = ts.ValidateTransactionResult(ctx, hash); err != nil {
		response.msg = fmt.Sprintf("ValidateTransactionResultAndEvents Unexpected error %v", err)
	} else {
		response.startEvent = &chain.TransferStartEvent{
			From: srcAddr,
			To:   dstAddr,
			Sn:   big.NewInt(-1),
			Assets: []chain.AssetTransferDetails{
				{
					Name:  coinNames[0],
					Value: amt,
					Fee:   big.NewInt(0),
				},
			},
		}
	}
	return
}

func stressReclaim(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (response *txnRecord, err error) {
	if len(coinNames) == 0 {
		err = errors.New("Should specify only one coinname, got zero")
		return
	}
	src, _, err := ts.GetChainPair(srcChain, dstChain)
	if err != nil {
		err = errors.Wrapf(err, "GetChainPair %v", err)
		return
	}
	ownerKey, ownerAddr, err := ts.GetDemoKeyPairs(srcChain)
	if err != nil {
		err = errors.Wrapf(err, "GetDemoKeyPairs %v", err)
		return
	}
	bal, err := src.GetCoinBalance(coinNames[0], ownerAddr)
	if err != nil {
		err = errors.Wrapf(err, "GetCoinBalance %v", err)
		return
	}
	var hash string
	hash, err = src.Reclaim(coinNames[0], ownerKey, bal.RefundableBalance)
	if err != nil {
		response.msg = fmt.Sprintf("Reclaim Err: %v", err)
		return
	}

	if _, err = ts.ValidateTransactionResult(ctx, hash); err != nil {
		response.msg = fmt.Sprintf("ValidateTransactionResultAndEvents Unexpected error %v", err)
	} else {
		response.startEvent = &chain.TransferStartEvent{
			From: ts.btsAddressPerChain[srcChain],
			To:   ownerAddr,
			Sn:   big.NewInt(-1),
			Assets: []chain.AssetTransferDetails{
				{
					Name:  coinNames[0],
					Value: bal.RefundableBalance,
					Fee:   big.NewInt(0),
				},
			},
		}
	}
	return
}

const (
	SRC_POS  = 0
	DST_POS  = 1
	COIN_POS = 2
	FUNC_POS = 3
)

const (
	INTER_FUNC       = 0
	INTRA_SRC_FUNC   = 1
	INTRA_DST_FUNC   = 2
	RECLAIM_SRC_FUNC = 3
	RECLAIM_DST_FUNC = 4
)

func (ex *executor) RunStressTest(ctx context.Context, srcChainName, dstChainName chain.ChainType, coinNames []string) error {
	if srcChainName == dstChainName {
		return fmt.Errorf("Src and Dst Chain should be different")
	}
	srcCl, ok := ex.clientsPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("Client for chain %v not found", srcChainName)
	}
	dstCl, ok := ex.clientsPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("Client for chain %v not found", dstChainName)
	}
	srcGod, ok := ex.godKeysPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("GodKeys for chain %v not found", srcChainName)
	}
	dstGod, ok := ex.godKeysPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("GodKeys for chain %v not found", dstChainName)
	}
	srcDemo, ok := ex.demoKeysPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("DemoKeys for chain %v not found", srcChainName)
	}
	srcDemo = append(srcDemo, srcGod)
	dstDemo, ok := ex.demoKeysPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("DemoKeys for chain %v not found", dstChainName)
	}
	dstDemo = append(dstDemo, dstGod)
	srcCfg, ok := ex.cfgPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
	}
	dstCfg, ok := ex.cfgPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
	}
	btsAddressPerChain := map[chain.ChainType]string{
		srcChainName: srcCfg.ContractAddresses[chain.BTS],
		dstChainName: dstCfg.ContractAddresses[chain.BTS],
	}
	gasLimitPerChain := map[chain.ChainType]int64{
		srcChainName: srcCfg.GasLimit,
		dstChainName: dstCfg.GasLimit,
	}

	ts := &testSuite{
		logger:             ex.log,
		env:                ex.env,
		btsAddressPerChain: btsAddressPerChain,
		gasLimitPerChain:   gasLimitPerChain,
		clsPerChain:        map[chain.ChainType]chain.ChainAPI{srcChainName: srcCl, dstChainName: dstCl},
		godKeysPerChain:    map[chain.ChainType]keypair{srcChainName: srcGod, dstChainName: dstGod},
		demoKeysPerChain:   map[chain.ChainType][]keypair{srcChainName: srcDemo, dstChainName: dstDemo},
		fee:                fee{numerator: big.NewInt(FEE_NUMERATOR), denominator: big.NewInt(FEE_DENOMINATOR), fixed: big.NewInt(FIXED_PRICE)},
	}
	coinBalanceSrc, err := getBalanceForAddress(srcCl, srcDemo, coinNames)
	if err != nil {
		return errors.Wrapf(err, "getBalanceForSrcAddress %v", err)
	}
	coinBalanceDst, err := getBalanceForAddress(dstCl, dstDemo, coinNames)
	if err != nil {
		return errors.Wrapf(err, "getBalanceForDstAddress %v", err)
	}

	type response struct {
		itr int
		req *req
		err error
		rec *txnRecord
	}
	requests := getRequests(len(ts.demoKeysPerChain[srcChainName]), len(ts.demoKeysPerChain[dstChainName]), len(coinNames))
	//log.Info("Len Requests ", len(requests))
	concurrency := len(requests)
	idx := 0
	for {
		select {
		case <-ctx.Done():
			log.Info("Inf Loop context cancelled ")
			return nil
		default:
			if idx == len(requests) {
				ex.log.Info("All requests have been served for this epoch. Exiting loop")
				return nil
			}
			qch := make(chan *response, concurrency)
			for i := idx; i < len(requests) && len(qch) < cap(qch); i++ {
				qch <- &response{itr: idx, req: requests[i], err: nil, rec: nil}
				time.Sleep(time.Millisecond)
				idx++
			}
			batchRes := make([]*response, 0, len(qch))
			for q := range qch {
				switch {
				case q.rec != nil || q.err != nil:
					batchRes = append(batchRes, q)
					if len(batchRes) == cap(batchRes) {
						close(qch)
					}
				default:
					go func(q *response) {
						defer func() {
							qch <- q
						}()
						rand.Seed(time.Now().UnixNano())
						if q.req.funcIdx == INTER_FUNC {
							id, err := ex.getID()
							if err != nil {
								q.err = errors.Wrap(err, "getID ")
								return
							}
							sinkChan := make(chan *evt)
							if err := ex.addChan(id, sinkChan); err != nil {
								q.err = err
								return
							}
							defer ex.removeChan(id)
							tsf := &testSuite{
								id:                 id,
								logger:             ex.log.WithFields(log.Fields{"pid": id}),
								env:                "testnet",
								subChan:            sinkChan,
								btsAddressPerChain: ts.btsAddressPerChain,
								gasLimitPerChain:   ts.gasLimitPerChain,
								clsPerChain:        ts.clsPerChain,
								godKeysPerChain:    ts.godKeysPerChain,
								demoKeysPerChain:   ts.demoKeysPerChain,
								fee:                ts.fee,
							}

							_v, _err := stressTransferInterChain(ctx, srcChainName, dstChainName, tsf.demoKeysPerChain[srcChainName][q.req.srcIdx], tsf.demoKeysPerChain[dstChainName][q.req.dstIdx], []string{coinNames[q.req.coinIdx]}, tsf)
							if _err != nil {
								q.err = errors.Wrapf(_err, "stressTransferInterChain %v", _err)
							}
							q.rec = _v
							q.err = _err
						} else if q.req.funcIdx == INTRA_SRC_FUNC {
							id, err := ex.getID()
							if err != nil {
								q.err = errors.Wrap(err, "getID ")
								return
							}
							tss := &testSuite{
								id:                 id,
								logger:             ex.log.WithFields(log.Fields{"pid": id}),
								env:                "testnet",
								btsAddressPerChain: ts.btsAddressPerChain,
								gasLimitPerChain:   ts.gasLimitPerChain,
								clsPerChain:        ts.clsPerChain,
								godKeysPerChain:    ts.godKeysPerChain,
								demoKeysPerChain:   ts.demoKeysPerChain,
								fee:                ts.fee,
							}
							_v, _err := stressTransferIntraChain(ctx, srcChainName, srcChainName, tss.demoKeysPerChain[srcChainName][q.req.srcIdx], tss.demoKeysPerChain[srcChainName][q.req.dstIdx], []string{coinNames[q.req.coinIdx]}, tss)
							if _err != nil {
								q.err = errors.Wrapf(_err, "stressTransferInterChain %v", _err)
							}
							q.rec = _v
							q.err = _err
						} else {
							q.err = errors.New("not implemented")
						}

					}(q)
				}
			}

			srcChangedTokens := map[int][]int{}
			dstChangedTokens := map[int][]int{}
			for _, res := range batchRes {
				if res.err == nil && len(res.rec.msg) == 0 {
					if res.req.funcIdx == INTER_FUNC && res.rec.startEvent != nil && res.rec.endEvent != nil && res.rec.endEvent.Code.String() == "0" {
						for _, as := range res.rec.startEvent.Assets {
							if as.Name == coinNames[res.req.coinIdx] {
								srcBal := coinBalanceSrc[res.req.srcIdx][res.req.coinIdx].UserBalance
								dstBal := coinBalanceDst[res.req.dstIdx][res.req.coinIdx].UserBalance
								srcBal = srcBal.Sub(srcBal, as.Fee)
								srcBal = srcBal.Sub(srcBal, as.Value)
								dstBal = dstBal.Add(dstBal, as.Value)
								coinBalanceSrc[res.req.srcIdx][res.req.coinIdx].UserBalance = srcBal
								coinBalanceDst[res.req.dstIdx][res.req.coinIdx].UserBalance = dstBal
								srcChangedTokens[res.req.srcIdx] = append(srcChangedTokens[res.req.srcIdx], res.req.coinIdx)
								dstChangedTokens[res.req.dstIdx] = append(dstChangedTokens[res.req.dstIdx], res.req.coinIdx)
							}
						}
					} else if res.req.funcIdx == INTRA_SRC_FUNC && res.rec.startEvent != nil {
						for _, as := range res.rec.startEvent.Assets {
							if as.Name == coinNames[res.req.coinIdx] {
								srcBal := coinBalanceSrc[res.req.srcIdx][res.req.coinIdx].UserBalance
								dstBal := coinBalanceSrc[res.req.dstIdx][res.req.coinIdx].UserBalance
								srcBal = srcBal.Sub(srcBal, as.Fee)
								srcBal = srcBal.Sub(srcBal, as.Value)
								dstBal = dstBal.Add(dstBal, as.Value)
								coinBalanceSrc[res.req.srcIdx][res.req.coinIdx].UserBalance = srcBal
								coinBalanceSrc[res.req.dstIdx][res.req.coinIdx].UserBalance = dstBal
								srcChangedTokens[res.req.srcIdx] = append(srcChangedTokens[res.req.srcIdx], res.req.coinIdx)
								srcChangedTokens[res.req.dstIdx] = append(srcChangedTokens[res.req.dstIdx], res.req.coinIdx)
							}
						}
					}
				}
			}
			if err := compareChangedTokens(srcCl, srcChangedTokens, coinBalanceSrc, coinNames, ts.demoKeysPerChain[srcChainName]); err != nil {
				return errors.Wrapf(err, "compareChangedTokens %v", err)
			}
			if err := compareChangedTokens(dstCl, dstChangedTokens, coinBalanceDst, coinNames, ts.demoKeysPerChain[dstChainName]); err != nil {
				return errors.Wrapf(err, "compareChangedTokens %v", err)
			}

		}
	}
	return nil
}

type req struct {
	srcIdx  int
	dstIdx  int
	coinIdx int
	funcIdx int
}

func getRequests(srcDemoLen, dstDemoLen, coinsLen int) []*req {
	intraTransfersSrcCount := srcDemoLen * (srcDemoLen - 1)
	intraTransfersDstCount := dstDemoLen * (dstDemoLen - 1)
	reclaimSrcCount := srcDemoLen
	reclaimDstCount := dstDemoLen
	interTransfersCount := srcDemoLen * dstDemoLen
	totalCount := (intraTransfersSrcCount + intraTransfersDstCount +
		reclaimDstCount + reclaimSrcCount +
		interTransfersCount) * coinsLen
	requests := make([]*req, totalCount)
	//interChain
	idx := 0
	for c := 0; c < coinsLen; c++ {
		for i := 0; i < srcDemoLen; i++ {
			for j := 0; j < dstDemoLen; j++ {
				requests[idx] = &req{i, j, c, INTER_FUNC}
				idx++
			}
		}
	}
	//intraChain
	for c := 0; c < coinsLen; c++ {
		for i := 0; i < srcDemoLen; i++ {
			for j := 0; j < srcDemoLen; j++ {
				if i == j {
					continue
				}
				requests[idx] = &req{i, j, c, INTRA_SRC_FUNC}
				idx++
			}
		}
	}
	for c := 0; c < coinsLen; c++ {
		for i := 0; i < dstDemoLen; i++ {
			for j := 0; j < dstDemoLen; j++ {
				if i == j {
					continue
				}
				requests[idx] = &req{i, j, c, INTRA_DST_FUNC}
				idx++
			}
		}
	}
	//reclaim
	for c := 0; c < coinsLen; c++ {
		for i := 0; i < reclaimSrcCount; i++ {
			_f := rand.Intn(srcDemoLen)
			requests[idx] = &req{_f, _f, c, RECLAIM_SRC_FUNC}
			idx++
		}
	}
	for c := 0; c < coinsLen; c++ {
		for i := 0; i < reclaimDstCount; i++ {
			_f := rand.Intn(dstDemoLen)
			requests[idx] = &req{_f, _f, c, RECLAIM_DST_FUNC}
			idx++
		}
	}
	//shuffle
	for i := 0; i < len(requests); i++ {
		_to := rand.Intn(len(requests))
		tmp := requests[i]
		requests[i] = requests[_to]
		requests[_to] = tmp
	}
	return requests
}

func getBalanceForAddress(cl chain.ChainAPI, addrs []keypair, coinNames []string) (bals [][]*chain.CoinBalance, err error) {
	bals = make([][]*chain.CoinBalance, len(addrs))
	for _ia, addr := range addrs {
		bals[_ia] = make([]*chain.CoinBalance, len(coinNames))
		for _ic, coinName := range coinNames {
			if bal, errs := cl.GetCoinBalance(coinName, cl.GetBTPAddress(addr.PubKey)); err != nil {
				err = errors.Wrapf(errs, "getCoinBalance %v", errs)
				return
			} else {
				bals[_ia][_ic] = bal
			}
		}
	}
	return
}

func compareChangedTokens(srcCl chain.ChainAPI, srcChangedTokens map[int][]int, origBalance [][]*chain.CoinBalance, coinNames []string, keyPairs []keypair) error {
	for addrIdx, coinsIdx := range srcChangedTokens {
		uniqCoin := map[int]bool{}
		for _, coinIdx := range coinsIdx {
			if _, ok := uniqCoin[coinIdx]; ok {
				continue
			}
			if coinNames[coinIdx] == srcCl.NativeCoin() { // skipping native for now
				continue
			}
			uniqCoin[coinIdx] = true
			addr := srcCl.GetBTPAddress(keyPairs[addrIdx].PubKey)
			res, err := srcCl.GetCoinBalance(coinNames[coinIdx], addr)
			if err != nil {
				return errors.Wrapf(err, "GetCoinBalance %v", err)
			}
			if res.UserBalance.Add(res.UserBalance, res.UsableBalance).String() != origBalance[addrIdx][coinIdx].UserBalance.Add(origBalance[addrIdx][coinIdx].UserBalance, origBalance[addrIdx][coinIdx].UsableBalance).String() {
				//fmt.Println("--------------- ", addr, " New ", res.String())
				//fmt.Println("--------------- ", " Old ", origBalance[addrIdx][coinIdx].String())
			} else {
				//fmt.Println("+++++++++++++++ All good")
			}
		}
	}
	return nil
}
