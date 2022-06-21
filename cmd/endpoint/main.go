package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/unitgroup"
	"github.com/icon-project/icon-bridge/common/log"
)

func init() {

}

const NUM_PARALLEL_DEMOS = 3

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

	ug, err := unitgroup.New(l, cfgPerMap)
	if err != nil {
		log.Fatal(err)
	}
	ug.Start(context.TODO())
	for i := 0; i < NUM_PARALLEL_DEMOS; i++ {
		tf := unitgroup.DefaultTaskFunctions["DemoTransaction"]
		if err := ug.RegisterTestUnit(map[chain.ChainType]int{chain.ICON: 1, chain.HMNY: 1}, tf, false); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * time.Duration(2))
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
