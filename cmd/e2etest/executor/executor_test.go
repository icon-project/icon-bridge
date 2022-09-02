package executor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
	"github.com/icon-project/icon-bridge/common/log"

	_ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
	// _ "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
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
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	ex.Subscribe(ctx)
	time.Sleep(5 * time.Second)
	ex.RunFlowTest(ctx, "ICON", "BSC", []string{"bnUSD"})
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

func TestStress(t *testing.T) {
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
	for i := 0; i < 1; i++ {
		fmt.Println("Epochs ", i)
		ex.RunStressTest(ctx, "ICON", "BSC", []string{"sICX", "bnUSD"})
	}
	// <-ex.Done()
	cancel()
	time.Sleep(time.Second * 2)
}

func TestKeystore(t *testing.T) {
	// prepare accounts
	// there is a path for each chain; GetKeyPairFromFile

	cfgPerMap, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}
	l := log.New()
	log.SetGlobalLogger(l)
	ex, err := executor.New(l, cfgPerMap)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	keyMap, err := getKeystores()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	tmpFile, err := ioutil.TempFile("./", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	for ch, keys := range keyMap {
		clsMap := ex.Clients()
		if cl, ok := clsMap[ch]; !ok {
			continue
		} else {
			for _, key := range keys {
				_, pub, err := cl.GetKeyPairFromKeystore(key, tmpFile.Name())
				if err != nil {
					t.Fatal(err)
				}
				bal, err := cl.GetCoinBalance("sICX", cl.GetBTPAddress(pub))
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(pub, "  ", bal.UserBalance)
			}
		}
	}

}

func getKeystores() (map[chain.ChainType][]string, error) {
	pathPerChain := map[chain.ChainType]string{
		chain.ICON: "../../../devnet/docker/icon-bsc/_ixh/wallets/icon",
		chain.BSC:  "../../../devnet/docker/icon-bsc/_ixh/wallets/bsc",
	}
	keystores := map[chain.ChainType][]string{
		chain.ICON: {},
		chain.BSC:  {},
	}
	for ch, dir := range pathPerChain {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			keystores[ch] = append(keystores[ch], filepath.Join(dir, f.Name()))
		}
	}
	return keystores, nil
}
