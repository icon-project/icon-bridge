package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/executor"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

func init() {

}

const NUM_PARALLEL_DEMOS = 2

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
	ug, err := executor.New(l, cfgPerMap)
	if err != nil {
		log.Error(errors.Wrap(err, "executor.New "))
		return
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	startHeight := uint64(100)
	ug.Start(ctx, startHeight)
	for i := 0; i < NUM_PARALLEL_DEMOS; i++ {
		log.Info("Register Process ", i)
		err = ug.Execute(ctx, []chain.ChainType{chain.ICON, chain.HMNY}, executor.DemoSubCallback)
		if err != nil {
			log.Error(err, errors.Wrap(err, "executor.Execute"))
		}
		time.Sleep(time.Second * time.Duration(15))
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigCh)
	}()

	<-sigCh // second signal, hard exit
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
