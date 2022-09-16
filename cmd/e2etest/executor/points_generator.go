package executor

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
)

func (gen *pointGenerator) singlePointGenerator(pts []*transferPoint) []*transferPoint {
	chains := make([]chain.ChainType, 0)
	for k := range gen.cfgPerChain {
		chains = append(chains, k)
	}
	for _, coinDetail := range gen.cfgPerChain[chains[0]].CoinDetails {
		pts = append(pts, &transferPoint{SrcChain: chains[0], DstChain: chains[1], CoinNames: []string{coinDetail.Name}, Amounts: []*big.Int{big.NewInt(1)}})
	}
	for _, coinDetail := range gen.cfgPerChain[chains[1]].CoinDetails {
		pts = append(pts, &transferPoint{SrcChain: chains[1], DstChain: chains[0], CoinNames: []string{coinDetail.Name}, Amounts: []*big.Int{big.NewInt(1)}})
	}
	return pts
}

func (gen *pointGenerator) batchPointGenerator(pts []*transferPoint) []*transferPoint {
	chains := make([]chain.ChainType, 0)
	for k := range gen.cfgPerChain {
		chains = append(chains, k)
	}
	for _, pair := range [][2]int{{0, 1}, {1, 0}} {
		pts = append(pts, &transferPoint{
			SrcChain: chains[pair[0]],
			DstChain: chains[pair[1]],
			CoinNames: []string{
				gen.cfgPerChain[chains[pair[0]]].NativeCoin,
				gen.cfgPerChain[chains[pair[0]]].NativeTokens[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].NativeTokens))],
			},
			Amounts: []*big.Int{
				big.NewInt(1),
				big.NewInt(2),
			},
		})

		pts = append(pts, &transferPoint{
			SrcChain: chains[pair[0]],
			DstChain: chains[pair[1]],
			CoinNames: []string{
				gen.cfgPerChain[chains[pair[0]]].NativeCoin,
				gen.cfgPerChain[chains[pair[0]]].WrappedCoins[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].WrappedCoins))],
			},
			Amounts: []*big.Int{
				big.NewInt(1),
				big.NewInt(2),
			},
		})

		pts = append(pts, &transferPoint{
			SrcChain: chains[pair[0]],
			DstChain: chains[pair[1]],
			CoinNames: []string{
				gen.cfgPerChain[chains[pair[0]]].NativeTokens[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].NativeTokens))],
				gen.cfgPerChain[chains[pair[0]]].NativeTokens[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].NativeTokens))],
			},
			Amounts: []*big.Int{
				big.NewInt(1),
				big.NewInt(2),
			},
		})

		pts = append(pts, &transferPoint{
			SrcChain: chains[pair[0]],
			DstChain: chains[pair[1]],
			CoinNames: []string{
				gen.cfgPerChain[chains[pair[0]]].WrappedCoins[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].WrappedCoins))],
				gen.cfgPerChain[chains[pair[0]]].WrappedCoins[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].WrappedCoins))],
			},
			Amounts: []*big.Int{
				big.NewInt(1),
				big.NewInt(2),
			},
		})

		pts = append(pts, &transferPoint{
			SrcChain: chains[pair[0]],
			DstChain: chains[pair[1]],
			CoinNames: []string{
				gen.cfgPerChain[chains[pair[0]]].NativeTokens[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].NativeTokens))],
				gen.cfgPerChain[chains[pair[0]]].WrappedCoins[rand.Intn(len(gen.cfgPerChain[chains[pair[0]]].WrappedCoins))],
			},
			Amounts: []*big.Int{
				big.NewInt(1),
				big.NewInt(2),
			},
		})
	}
	return pts
}

func (gen *pointGenerator) GenerateTransferPoints(cpt *configPoint) (pts []*transferPoint, err error) {
	if len(gen.cfgPerChain) != 2 {
		err = fmt.Errorf("Expected a pair of chains. Got %v", len(gen.cfgPerChain))
		return
	}
	pts = []*transferPoint{}
	pts = gen.singlePointGenerator(pts)
	pts = gen.batchPointGenerator(pts)
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

func (gen *pointGenerator) tokenLimitsGenerator(pts []*configPoint) {
	chains := make([]chain.ChainType, 0)
	for k := range gen.cfgPerChain {
		chains = append(chains, k)
	}
	zero := big.NewInt(0)
	one := big.NewInt(1)
	maxUint256, _ := (&big.Int{}).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	maxUint256MinusOne, _ := (&big.Int{}).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639934", 10)
	batch := []*big.Int{zero, one, maxUint256, maxUint256MinusOne}
	for _, chain := range chains {
		chainNativeCoin := gen.cfgPerChain[chain].NativeCoin
		chainNativeTokens := gen.cfgPerChain[chain].NativeTokens
		if zeroNativeBeforeFee, err := gen.getAmountBeforeFeeCharge(chain, chainNativeCoin, big.NewInt(0)); err != nil {
			batch = append(batch, zeroNativeBeforeFee)
		}
		if oneNativeBeforeFee, err := gen.getAmountBeforeFeeCharge(chain, chainNativeCoin, big.NewInt(1)); err != nil {
			batch = append(batch, oneNativeBeforeFee)
		}
		if zeroTokenBeforeFee, err := gen.getAmountBeforeFeeCharge(chain, chainNativeTokens[rand.Intn(len(chainNativeTokens))], big.NewInt(1)); err != nil {
			batch = append(batch, zeroTokenBeforeFee)
		}
		if oneTokenBeforeFee, err := gen.getAmountBeforeFeeCharge(chain, chainNativeTokens[rand.Intn(len(chainNativeTokens))], big.NewInt(1)); err != nil {
			batch = append(batch, oneTokenBeforeFee)
		}
		for _, amt := range []*big.Int{zero, one, maxUint256, maxUint256MinusOne} {
			pts = append(pts, &configPoint{
				chainName: chain,
				TokenLimits: map[string]*big.Int{
					chainNativeCoin: (&big.Int{}).Set(amt),
				},
			})
			pts = append(pts, &configPoint{
				chainName: chain,
				TokenLimits: map[string]*big.Int{
					chainNativeTokens[rand.Intn(len(chainNativeTokens))]: (&big.Int{}).Set(amt),
				},
			})
			pts = append(pts, &configPoint{
				chainName: chain,
				TokenLimits: map[string]*big.Int{
					chainNativeTokens[rand.Intn(len(chainNativeTokens))]: (&big.Int{}).Set(amt),
					chainNativeTokens[rand.Intn(len(chainNativeTokens))]: (&big.Int{}).Set(amt),
				},
			})
			pts = append(pts, &configPoint{
				chainName: chain,
				TokenLimits: map[string]*big.Int{
					chainNativeCoin: (&big.Int{}).Set(amt),
					chainNativeTokens[rand.Intn(len(chainNativeTokens))]: (&big.Int{}).Set(amt),
				},
			})
			// pts = append(pts, &configPoint{
			// 	TokenLimits: map[string]*big.Int{
			// 		"apple": (&big.Int{}).Set(amt),
			// 		"ball":  (&big.Int{}).Set(amt),
			// 	},
			// })
		}
	}

}

func (gen *pointGenerator) feeGenerator(pts []*configPoint) {

}

func (gen *pointGenerator) GenerateConfigPoints() (pts []*configPoint, err error) {
	// if len(gen.cfgPerChain) != 2 {
	// 	err = fmt.Errorf("Expected a pair of chains. Got %v", len(gen.cfgPerChain))
	// 	return
	// }
	// pts = []*configPoint{}
	// gen.tokenLimitsGenerator(pts)
	// gen.feeGenerator(pts)
	// if gen.transferFilter != nil {
	// 	pts = gen.configFilter(pts)
	// }

	// arrLen := len(pts)
	// for i := 0; i < arrLen; i++ {
	// 	a := rand.Intn(arrLen)
	// 	b := rand.Intn(arrLen)
	// 	tmp := pts[a]
	// 	pts[a] = pts[b]
	// 	pts[b] = tmp
	// }

	// if gen.maxBatchSize != nil && len(pts) > *gen.maxBatchSize {
	// 	truncPts := make([]*configPoint, *gen.maxBatchSize)
	// 	for i := 0; i < *gen.maxBatchSize; i++ {
	// 		truncPts[i] = pts[i]
	// 	}
	// 	return truncPts, nil
	// }
	// return
	pts = []*configPoint{
		{
			chainName: chain.BSC,
			TokenLimits: map[string]*big.Int{
				"btp-0x2.icon-ICX":  big.NewInt(9000000000000000000),
				"btp-0x2.icon-BUSD": big.NewInt(7000000000000000000),
			},
			Fee: map[string][2]*big.Int{
				"btp-0x2.icon-sICX": {
					big.NewInt(100), big.NewInt(3900000000000000000),
				},
				"btp-0x2.icon-BTCB": {
					big.NewInt(100), big.NewInt(62500000000000),
				},
			},
		},
	}
	return
}

func (ts *pointGenerator) getAmountBeforeFeeCharge(chainName chain.ChainType, coinName string, outputBalance *big.Int) (*big.Int, error) {
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
