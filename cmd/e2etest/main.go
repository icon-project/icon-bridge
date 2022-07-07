package main

import (
	"context"
	"encoding/json"
	"math/big"
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
	cfgPerMap := map[chain.ChainType]*chain.ChainConfig{}
	for _, ch := range cfg.Chains {
		cfgPerMap[ch.Name] = ch
	}
	ex, err := executor.New(l, cfgPerMap)
	if err != nil {
		log.Error(errors.Wrap(err, "executor.New "))
		return
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ex.Subscribe(ctx)

	fundAmount := new(big.Int)
	fundAmount.SetString("10000000000000000000", 10)

	for _, fts := range testCfg.FlowTests {
		if fts.SrcChain == chain.ICON {
			for _, coinName := range fts.CoinNames {
				go func(coinName string) {
					script := executor.MonitorTransferWithApproveFromICON
					if coinName == "ICX" {
						script = executor.MonitorTransferWithoutApproveFromICON
					}
					err = ex.Execute(ctx, fts.SrcChain, fts.DstChain, coinName, fundAmount, script)
					if err != nil {
						log.Errorf("%+v", err)
					}
				}(coinName)
				time.Sleep(time.Second * 5)
			}
		} else if fts.SrcChain == chain.HMNY {
			for _, coinName := range fts.CoinNames {
				go func(coinName string) {
					script := executor.MonitorTransferWithApproveFromHMNY
					if coinName == "ONE" {
						script = executor.MonitorTransferWithoutApproveFromHMNY
					}
					err = ex.Execute(ctx, fts.SrcChain, fts.DstChain, coinName, fundAmount, script)
					if err != nil {
						log.Errorf("%+v", err)
					}
				}(coinName)
				time.Sleep(time.Second * 5)
			}
		}
	}

	defer func() {
		cancel()
	}()
	<-ex.Done()
	cancel()
	time.Sleep(time.Second * 2)
	log.Warn("Exit...")

}

type Config struct {
	Chains []*chain.ChainConfig `json:"chains"`
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

type TestConfig struct {
	FlowTests []*FlowTestConfig `json:"flowTests"`
}

type FlowTestConfig struct {
	SrcChain  chain.ChainType `json:"srcChain"`
	DstChain  chain.ChainType `json:"dstChain"`
	CoinNames []string        `json:"coins"`
}
