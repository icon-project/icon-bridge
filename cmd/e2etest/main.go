package main

import (
	"context"
	"encoding/json"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

func main() {
	l := log.New()
	log.SetGlobalLogger(l)
	cfg, err := loadConfig("./example-config.json")
	if err != nil {
		log.Error(errors.Wrap(err, "loadConfig "))
		return
	}
	testCfg, err := loadTestConfig("./test-config.json")
	if err != nil {
		log.Error(errors.Wrap(err, "loadConfig "))
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	cfgPerMap := map[chain.ChainType]*chain.ChainConfig{}
	for _, ch := range cfg.Chains {
		cfgPerMap[ch.Name] = ch
	}
	ex, err := executor.New(l, cfgPerMap)
	if err != nil {
		log.Error(errors.Wrap(err, "executor.New "))
		return
	}
	ex.Subscribe(ctx)

	if !testCfg.FlowTest.Disable {
		log.Info("Starting Flow Test ....")
		fundAmount := new(big.Int)
		fundAmount.SetString("10000000000000000000", 10)
		for _, fts := range testCfg.FlowTest.Chains {
			for _, coinName := range fts.CoinNames {
				go func(coinName string) {
					err = ex.Execute(ctx, fts.SrcChain, fts.DstChain, coinName, fundAmount, executor.Transfer)
					if err != nil {
						log.Errorf("%+v", err)
					}
				}(coinName)
				time.Sleep(time.Second * 5)
			}
		}
		<-ex.Done()
	}

	if !testCfg.StressTest.Disable {
		log.Info("Starting Stress Test ....")
		if len(testCfg.StressTest.AddressMap) <= 1 {
			log.Error("Require at least two chains for inter chain tests")
		}
		log.Info("Fund addresses ....")
		if addrsPerChain, err := ex.GetFundedAddresses(testCfg.StressTest.AddressMap); err != nil {
			log.Errorf("%v", err)
			return
		} else {
			cns := []chain.ChainType{}
			for cn := range addrsPerChain {
				cns = append(cns, cn)
			}
			// TODO
			allCoins := []string{"ICX", "TICX", "BNB", "TBNB"}
			log.Error("Run Jobs")
			for j := 0; j < int(testCfg.StressTest.JobsCount); j++ {
				rand.Seed(time.Now().UnixNano())
				go func() {
					srcChainType, dstChainType := getRandomChains(cns)
					coin := allCoins[rand.Intn(len(allCoins))]
					srcAddr := addrsPerChain[srcChainType][rand.Intn(len(addrsPerChain[srcChainType]))]
					dstAddr := addrsPerChain[dstChainType][rand.Intn(len(addrsPerChain[dstChainType]))]
					if err := ex.ExecuteOnAddr(ctx, srcChainType, dstChainType, coin, srcAddr, dstAddr, executor.StressTransfer); err != nil {
						log.Errorf("%v", err)
					}
				}()
				time.Sleep(time.Second * 5)
			}
			<-ex.Done()
		}
	}

	cancel()
	time.Sleep(time.Second * 2)
	log.Warn("Exit...")

}

func getRandomChains(cns []chain.ChainType) (chain.ChainType, chain.ChainType) {
	count := len(cns)
	if count == 1 {
		return cns[0], cns[0]
	}
	first := rand.Intn(count)
	for i := 0; i < 10; i++ { // try at max 10 times to get a pair
		second := rand.Intn(count)
		if second != first {
			return cns[first], cns[second]
		}
	}
	return cns[0], cns[count-1]
}

func loadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open file %v", file)
	}
	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "json.Decode file %v", file)
	}
	return cfg, nil
}

func loadTestConfig(file string) (*TestConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open file %v", file)
	}
	cfg := &TestConfig{}
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "json.Decode file %v", file)
	}
	return cfg, nil
}

type Config struct {
	Chains []*chain.ChainConfig `json:"chains"`
}

type TestConfig struct {
	FlowTest   *FlowTestConfig   `json:"flowTest"`
	StressTest *StressTestConfig `json:"stressTest"`
}

type FlowTestConfig struct {
	Disable bool               `json:"disable"`
	Chains  []*FlowChainConfig `json:"chains"`
}

type FlowChainConfig struct {
	SrcChain  chain.ChainType `json:"srcChain"`
	DstChain  chain.ChainType `json:"dstChain"`
	CoinNames []string        `json:"coins"`
}

type StressTestConfig struct {
	Disable    bool                     `json:"disable"`
	AddressMap map[chain.ChainType]uint `json:"addresses"`
	JobsCount  uint                     `json:"jobs"`
}
