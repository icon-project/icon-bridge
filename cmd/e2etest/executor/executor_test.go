package executor_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"

	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
)

func getConfig() (*executor.Config, error) {
	loadConfig := func(file string) (*executor.Config, error) {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		cfg := &executor.Config{}
		err = json.NewDecoder(f).Decode(cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	var err error
	cfg, err := loadConfig("../example-config.json")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func TestExecutor(t *testing.T) {
	cfg, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}
	l := log.New()
	log.SetGlobalLogger(l)
	ex, err := executor.New(l, cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	ex.Subscribe(ctx)
	time.Sleep(5 * time.Second)
	err = ex.RunFlowTest(ctx)
	if err != nil {
		t.Error(err)
	}
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
