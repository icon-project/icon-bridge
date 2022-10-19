package executor

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
)

type pointGenerator struct {
	cfgPerChain    map[chain.ChainType]*chain.Config
	clsPerChain    map[chain.ChainType]chain.ChainAPI
	maxBatchSize   *int
	transferFilter func([]*transferPoint) []*transferPoint
	configFilter   func([]*configPoint) []*configPoint
}

type transferPoint struct {
	SrcChain  chain.ChainType
	DstChain  chain.ChainType
	CoinNames []string
	Amounts   []*big.Int
}

type configPoint struct {
	chainName   chain.ChainType
	TokenLimits map[string]*big.Int
	Fee         map[string][2]*big.Int
}

type tmpCfg struct {
	numerator *big.Int
	baseFee   *big.Int
	limit     *big.Int
}

func (gen *pointGenerator) singlePointGenerator(pts []*transferPoint, cfgPerCoinPerChain map[chain.ChainType]map[string]*tmpCfg) []*transferPoint {
	chains := make([]chain.ChainType, 0)
	for k := range gen.cfgPerChain {
		chains = append(chains, k)
	}
	srcChain := chains[0]
	dstChain := chains[1]
	for coinName, cfg := range cfgPerCoinPerChain[srcChain] {
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{
			gen.getAmountBeforeFeeCharge(srcChain, coinName, cfg.baseFee, cfg.numerator, big.NewInt(1)),
		}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{
			gen.getAmountBeforeFeeCharge(srcChain, coinName, cfg.baseFee, cfg.numerator, big.NewInt(-1)),
		}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{
			gen.getAmountBeforeFeeCharge(srcChain, coinName, cfg.baseFee, cfg.numerator, big.NewInt(0)),
		}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{(&big.Int{}).Add(cfg.limit, big.NewInt(1))}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{(&big.Int{}).Add(cfg.limit, big.NewInt(0))}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{(&big.Int{}).Add(cfg.limit, big.NewInt(-1))}})
	}
	srcChain = chains[1]
	dstChain = chains[0]
	for coinName, cfg := range cfgPerCoinPerChain[srcChain] {
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{
			gen.getAmountBeforeFeeCharge(srcChain, coinName, cfg.baseFee, cfg.numerator, big.NewInt(1)),
		}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{
			gen.getAmountBeforeFeeCharge(srcChain, coinName, cfg.baseFee, cfg.numerator, big.NewInt(-1)),
		}})
		pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{
			gen.getAmountBeforeFeeCharge(srcChain, coinName, cfg.baseFee, cfg.numerator, big.NewInt(0)),
		}})
		//pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{(&big.Int{}).Add(cfg.limit, big.NewInt(1))}})
		//pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{(&big.Int{}).Add(cfg.limit, big.NewInt(0))}})
		//pts = append(pts, &transferPoint{SrcChain: srcChain, DstChain: dstChain, CoinNames: []string{coinName}, Amounts: []*big.Int{(&big.Int{}).Add(cfg.limit, big.NewInt(-1))}})
	}
	return pts
}

func (gen *pointGenerator) batchPointGenerator(pts []*transferPoint, cfgPerCoinPerChain map[chain.ChainType]map[string]*tmpCfg) []*transferPoint {
	chains := make([]chain.ChainType, 0)
	for k := range gen.cfgPerChain {
		chains = append(chains, k)
	}
	for _, pair := range [][2]int{{0, 1}, {1, 0}} {
		chainA := chains[pair[0]]
		chainB := chains[pair[1]]
		tp := &transferPoint{
			SrcChain: chainA,
			DstChain: chainB,
			CoinNames: []string{
				gen.cfgPerChain[chainA].NativeCoin,
				gen.cfgPerChain[chainA].NativeTokens[rand.Intn(len(gen.cfgPerChain[chainA].NativeTokens))],
			},
		}
		for _, c := range tp.CoinNames {
			tp.Amounts = append(tp.Amounts, gen.getAmountBeforeFeeCharge(chainA, c, cfgPerCoinPerChain[chainA][c].baseFee, cfgPerCoinPerChain[chainA][c].numerator, big.NewInt(1)))
		}
		pts = append(pts, tp)

		tp = &transferPoint{
			SrcChain: chainA,
			DstChain: chainB,
			CoinNames: []string{
				gen.cfgPerChain[chainA].NativeCoin,
				gen.cfgPerChain[chainA].WrappedCoins[rand.Intn(len(gen.cfgPerChain[chainA].WrappedCoins))],
			},
		}
		for _, c := range tp.CoinNames {
			tp.Amounts = append(tp.Amounts, gen.getAmountBeforeFeeCharge(chainA, c, cfgPerCoinPerChain[chainA][c].baseFee, cfgPerCoinPerChain[chainA][c].numerator, big.NewInt(1)))
		}
		pts = append(pts, tp)

		pts = append(pts, &transferPoint{
			SrcChain: chainA,
			DstChain: chainB,
			CoinNames: []string{
				gen.cfgPerChain[chainA].NativeTokens[rand.Intn(len(gen.cfgPerChain[chainA].NativeTokens))],
				gen.cfgPerChain[chainA].WrappedCoins[rand.Intn(len(gen.cfgPerChain[chainA].WrappedCoins))],
			},
		})
		for _, c := range tp.CoinNames {
			tp.Amounts = append(tp.Amounts, gen.getAmountBeforeFeeCharge(chainA, c, cfgPerCoinPerChain[chainA][c].baseFee, cfgPerCoinPerChain[chainA][c].numerator, big.NewInt(1)))
		}
		pts = append(pts, tp)
	}
	return pts
}

func (gen *pointGenerator) GenerateTransferPoints() (pts []*transferPoint, errs error) {
	if len(gen.cfgPerChain) != 2 {
		errs = fmt.Errorf("Expected a pair of chains. Got %v", len(gen.cfgPerChain))
		return
	}
	cfgPerCoinPerChain := map[chain.ChainType]map[string]*tmpCfg{}
	for chain, cl := range gen.clsPerChain {
		cfgPerCoinPerChain[chain] = map[string]*tmpCfg{}
		cfg := gen.cfgPerChain[chain]
		for _, coinName := range append(append(cfg.NativeTokens, cfg.NativeCoin), cfg.WrappedCoins...) {
			fNum, fBase, err := cl.GetFeeRatio(coinName)
			if err != nil {
				err = fmt.Errorf("GetFeeRatio %v", err)
				return
			}
			limit, err := cl.GetTokenLimit(coinName)
			if err != nil {
				err = fmt.Errorf("GetTokenLimit %v", err)
			}
			cfgPerCoinPerChain[chain][coinName] = &tmpCfg{numerator: fNum, baseFee: fBase, limit: limit}
		}
	}

	pts = []*transferPoint{}
	pts = gen.singlePointGenerator(pts, cfgPerCoinPerChain)
	pts = gen.batchPointGenerator(pts, cfgPerCoinPerChain)
	if gen.transferFilter != nil {
		pts = gen.transferFilter(pts)
	}
	arrLen := len(pts)
	for i := 0; i < arrLen; i++ {
		a := rand.Intn(arrLen)
		b := rand.Intn(arrLen)
		tmp := pts[a]
		pts[a] = pts[b]
		pts[b] = tmp
	}

	if gen.maxBatchSize != nil && len(pts) > *gen.maxBatchSize {
		truncPts := make([]*transferPoint, *gen.maxBatchSize)
		for i := 0; i < *gen.maxBatchSize; i++ {
			truncPts[i] = pts[i]
		}
		return truncPts, nil
	}
	return
}

func (gen *pointGenerator) GenerateConfigPoints() (pts []*configPoint, err error) {
	highValue := new(big.Int)
	highValue.SetString("11579208923731619542357098500868790785326998466564056403", 10)
	pts = []*configPoint{
		{
			chainName: chain.ICON,
			TokenLimits: map[string]*big.Int{
				"btp-0x2.icon-ICX":  highValue,
				"btp-0x2.icon-BUSD": highValue,
			},
			Fee: map[string][2]*big.Int{
				"btp-0x2.icon-ICX": {
					big.NewInt(100), big.NewInt(4300000000000000000),
				},
			},
		},
		{
			chainName: chain.ICON,
			TokenLimits: map[string]*big.Int{
				"btp-0x2.icon-sICX": highValue,
			},
			Fee: map[string][2]*big.Int{
				"btp-0x2.icon-sICX": {
					big.NewInt(100), big.NewInt(3900000000000000000),
				},
			},
		},
	}
	return
}

func (ts *pointGenerator) getAmountBeforeFeeCharge(chainName chain.ChainType, coinName string, fixedFee, feeNumerator, outputBalance *big.Int) *big.Int {
	/*
		What is the input amount that we must have so that the net transferrable amount
		after fee charged on chainName is equal to outputBalance for coinName ?
		feeCharged = inputBalance * ratio + fixedFee
		outputBalance = inputBalance - feeCharged
		inputBalance = (outputBalance + fixed) / (1 - ratio)
		inputBalance = (outputBalance + fixed) * deniminator / (denominator - numerator)

	*/
	bplusf := (&big.Int{}).Add(outputBalance, fixedFee)
	bplusf.Mul(bplusf, big.NewInt(DENOMINATOR))
	dminusn := new(big.Int).Sub(big.NewInt(DENOMINATOR), feeNumerator)
	bplusf.Div(bplusf, dminusn)
	return bplusf
}
