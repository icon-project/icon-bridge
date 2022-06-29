package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/executor"
	"github.com/icon-project/icon-bridge/common/log"
)

func init() {

}

const NUM_PARALLEL_DEMOS = 1

func main() {
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
	ug, err := executor.New(l, cfgPerMap)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < NUM_PARALLEL_DEMOS; i++ {
		log.Info("Register Process ", i)
		err = ug.Execute([]chain.ChainType{chain.ICON, chain.HMNY}, executor.DemoSubCallback)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * time.Duration(5))
	}
	fmt.Println("Wait")
	time.Sleep(time.Second * time.Duration(3000))
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
