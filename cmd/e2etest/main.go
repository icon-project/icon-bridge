package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"

	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
)

func main() {
	var cfgFile string
	var maxTasks int
	flag.StringVar(&cfgFile, "config", "", "e2e.config.json file")
	flag.IntVar(&maxTasks, "maxTasks", 0, "maximum number of tasks to run")
	flag.Parse()
	cfg, err := executor.LoadConfig(cfgFile)
	if err != nil {
		log.Error(errors.Wrap(err, "loadConfig "))
		return
	}
	l := executor.SetLogger(cfg)
	l.Info("Initializing service")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigCh)
		cancel()
	}()

	go func() {
		select {
		case <-sigCh: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-sigCh // second signal, hard exit
		os.Exit(2)
	}()

	l.Info("Running Service")
	ex, err := executor.New(l, cfg)
	if err != nil {
		log.Error(errors.Wrap(err, "executor.New "))
		return
	}
	ex.Subscribe(ctx)
	time.Sleep(5) // wait for subscription to start
	err = ex.RunFlowTest(ctx, maxTasks)
	if err != nil {
		log.Errorf("Error Executing Flow Test %+v", err)
	}
	time.Sleep(time.Second * 2)
	log.Warn("Exiting Service...")
}
