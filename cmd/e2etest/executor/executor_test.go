package executor_test

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestExecutor(t *testing.T) {
	type Config struct {
		Chains []*chain.ChainConfig `json:"chains"`
	}
	loadConfig := func(file string) (*Config, error) {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		cfg := &Config{}
		err = json.NewDecoder(f).Decode(cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	var err error
	cfg, err := loadConfig("../example-config.json")
	if err != nil {
		t.Fatal(err)
	}
	cfgPerMap := map[chain.ChainType]*chain.ChainConfig{}
	for _, ch := range cfg.Chains {
		cfgPerMap[ch.Name] = ch
	}

	l := log.New()
	log.SetGlobalLogger(l)
	ex, err := executor.New(l, cfgPerMap)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)

	ex.Subscribe(ctx)

	amount := new(big.Int)
	amount.SetString("10000000000000000000", 10)

	go func() {
		err = ex.Execute(ctx, chain.ICON, chain.BSC, "ICX", amount, executor.TransferWithoutApproveFromICON)
		if err != nil {
			log.Errorf("%+v", err)
		}
	}()
	time.Sleep(time.Second * 5)

	defer func() {
		cancel()
	}()
	<-ex.Done()
	cancel()
	time.Sleep(time.Second * 3)
	l.Info("Exit")
}
