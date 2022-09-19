package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"

	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
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
	if lv, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.Panicf("Invalid log_level=%s", cfg.LogLevel)
	} else {
		l.SetConsoleLevel(lv)
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
	log.Info("Starting Flow Test ....")
	err = ex.RunFlowTest(ctx)
	if err != nil {
		log.Errorf("%+v", err)
	}
	cancel()
	time.Sleep(time.Second * 2)
	log.Warn("Exit...")
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
