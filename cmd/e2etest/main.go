package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"

	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
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

	ex, err := executor.New(l, cfg)
	if err != nil {
		log.Error(errors.Wrap(err, "executor.New "))
		return
	}
	ex.Subscribe(ctx)
	time.Sleep(5) // wait for subscription to start
	if !testCfg.FlowTest.Disable {
		log.Info("Starting Flow Test ....")
		for _, fts := range testCfg.FlowTest.Chains {
			for _, coin := range fts.CoinNames {
				err = ex.RunFlowTest(ctx, fts.SrcChain, fts.DstChain, []string{coin})
				if err != nil {
					log.Errorf("%+v", err)
				}
			}
		}
	}
	if !testCfg.StressTest.Disable {
		log.Info("Starting Stress Test ....")
		for _, fts := range testCfg.FlowTest.Chains {
			err = ex.RunStressTest(ctx, fts.SrcChain, fts.DstChain, fts.CoinNames)
			if err != nil {
				log.Errorf("%+v", err)
			}
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

func loadConfig(file string) (*executor.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open file %v", file)
	}
	cfg := &executor.Config{}
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
	Disable bool `json:"disable"`
}
