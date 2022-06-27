package chainAPI_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI"
	chain "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

type Config struct {
	Chains []*chain.ChainConfig `json:"chains"`
}

func loadConfig(file string) (*Config, error) {
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

func TestAPISubscription(t *testing.T) {
	l := log.New()
	log.SetGlobalLogger(l)
	cfg, err := loadConfig("/home/manish/go/src/work/icon-bridge/cmd/endpoint/example-config.json")
	if err != nil {
		log.Error(err)
		return
	}
	cfgPerMap := map[chain.ChainType]*chain.ChainConfig{}
	for _, ch := range cfg.Chains {
		cfgPerMap[ch.Name] = ch
	}
	capi, err := chainAPI.New(l, cfgPerMap)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case msg := <-capi.SubEventChan():
				fmt.Println("Message ", msg)
			case errmsg := <-capi.ErrChan():
				fmt.Println("Err ", errmsg)
			case <-ctx.Done():
				fmt.Println("Context cancelled from outside")
				return
			}
		}
	}()

	capi.StartSubscription(ctx)
	log.Warn("Wait for cancel in 20 seconds")
	time.Sleep(time.Second * 4)
	cancel()
	log.Warn("Should have cancelled")
	time.Sleep(time.Second * 5)
}
