package executor_test

import (
	"context"
	"encoding/json"
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
	cfg, err := loadConfig("/home/manish/go/src/work/icon-bridge/cmd/e2etest/example-config.json")
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
		t.Fatalf("%+v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	ex.Subscribe(ctx)
	time.Sleep(5 * time.Second)
	go func() {
		err = ex.Execute(ctx, chain.BSC, chain.ICON, "TBNB", executor.TransferToIncorrectAddress)
		if err != nil {
			log.Errorf("%+v", err)
		}
	}()
	<-ex.Done()
	cancel()
	time.Sleep(time.Second * 2)

	// defer func() {
	// 	cancel()
	// }()
	// <-ex.Done()
	// cancel()
	// time.Sleep(time.Second * 3)
	// l.Info("Exit")
}
