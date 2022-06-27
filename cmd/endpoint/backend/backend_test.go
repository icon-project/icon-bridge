package backend

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestBackend(t *testing.T) {
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
	be, err := New(l, cfgPerMap)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	be.Start(ctx)
	time.Sleep(time.Second * 10)
	t.Log("Send Cancel")
	cancel()
	time.Sleep(time.Second * 5)
}

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
