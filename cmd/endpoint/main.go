package main

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/executor"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

func main() {
	l := log.New()
	log.SetGlobalLogger(l)
	cfg, err := loadConfig("/home/manish/go/src/work/icon-bridge/cmd/endpoint/example-config.json")
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

	amount := new(big.Int)
	amount.SetString("10000000000000000000", 10)
	for tsi, ts := range executor.TestScripts {
		l.Info("Running TestScript SN.", tsi)
		go func() {
			err = ex.Execute(ctx, chain.ICON, chain.HMNY, amount, ts)
			if err != nil {
				log.Errorf("%+v", err)
			}
		}()
		time.Sleep(time.Second * 5)
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
