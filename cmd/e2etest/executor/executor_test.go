package executor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"

	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
)

func TestExecutor(t *testing.T) {
	cfg, err := executor.LoadConfig("../example-config.json")
	if err != nil {
		log.Error(errors.Wrap(err, "loadConfig "))
		return
	}
	l := executor.SetLogger(cfg)

	ex, err := executor.New(l, cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	ex.Subscribe(ctx)
	time.Sleep(5 * time.Second)
	err = ex.RunFlowTest(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Main Context Cancel")
	cancel()

	time.Sleep(time.Second * 5)
	fmt.Println("Exit")
	// defer func() {
	// 	cancel()
	// }()
	// <-ex.Done()
	// cancel()
	// time.Sleep(time.Second * 3)
	// l.Info("Exit")
}
