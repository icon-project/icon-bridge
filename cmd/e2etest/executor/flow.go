package executor

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type transferReq struct {
	scr Script
	pt  *transferPoint
	res *txnRecord
	err error
	id  uint64
}

func (ex *executor) RunFlowTest(ctx context.Context) error {
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

	pts, err := tg.GenerateTransferPoints()
	if err != nil {
		return errors.Wrapf(err, "GenerateTransferPoints %v", err)
	}
	scripts := []Script{
		TransferToBlackListedDstAddress,
	}
	batchRequest := make([]*transferReq, len(scripts)*len(pts))
	tmpi := 0
	for _, pt := range pts {
		for _, scr := range scripts {
			batchRequest[tmpi] = &transferReq{scr: scr, pt: pt}
			tmpi++
		}
	}
	fmt.Println("Len Request ", len(batchRequest))
	batchResponse, batchError, err := ex.processBatchForFlowTest(ctx, batchRequest)
	if err != nil {
		return errors.Wrapf(err, "processBatchForFlowTest %v", err)
	}
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

	return nil
}

func (ex *executor) processBatchForFlowTest(ctx context.Context, requests []*transferReq) (totalResponse []*txnRecord, totalError []error, err error) {

	concurrency := len(requests)
	if concurrency <= 1 {
		concurrency = 2
	}
	totalResponse = []*txnRecord{}
	totalError = []error{}
	reqCursor := 0
	for reqCursor < len(requests) {
		rqch := make(chan *transferReq, concurrency)
		for i := reqCursor; len(rqch) < cap(rqch) && i < len(requests); i++ {
			rqch <- requests[i]
			reqCursor++
		}
		sres := make([]*txnRecord, 0, len(rqch))
		for q := range rqch {
			switch {
			case q.err != nil || q.res != nil:
				if q.err != nil {
					totalError = append(totalError, q.err)
				}
				sres = append(sres, q.res)
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
					if q.scr.Callback == nil {
						q.err = errors.New("Callback nil")
						return
					}
					id, err := ex.getID()
					if err != nil {
						q.err = errors.Wrap(err, "getID ")
						return
					}
					q.id = id
					log := ex.log.WithFields(log.Fields{"pid": id})
					sinkChan := make(chan *evt)
					ex.addChan(id, sinkChan)
					defer ex.removeChan(id)

					ts := &testSuite{
						id:                   id,
						logger:               log,
						subChan:              sinkChan,
						clsPerChain:          ex.clientsPerChain,
						godKeysPerChain:      ex.godKeysPerChain,
						cfgPerChain:          ex.cfgPerChain,
						feeAggregatorAddress: ex.feeAggregatorAddress,
					}
					ts.logger.Debug("%v %v %v %v \n", q.scr.Name, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames)
					q.res, err = q.scr.Callback(ctx, q.pt.SrcChain, q.pt.DstChain, q.pt.CoinNames, ts)
					if err != nil {
						q.err = err
						return
					}
				}(q)

			}
		}
		for _, rs := range sres {
			if rs == nil {
				continue
			}
			totalResponse = append(totalResponse, rs)
		}
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
			} else {
				ex.log.Infof("Insufficient to redeem addr %v gasFeeOnSrc %v UserBalance %v", addr, gasFeeOnSrc, bal.UserBalance)
			}

		}

	}
}
