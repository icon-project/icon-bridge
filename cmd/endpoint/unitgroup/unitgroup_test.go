package unitgroup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestUnitGroup(t *testing.T) {
	type Config struct {
		Chains []*chain.ChainConfig `json:"chains"`
	}
	fmt.Println("Start")
	loadConfig := func(file string) (*Config, error) {
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
	var err error
	cfg, err := loadConfig("/home/manish/go/src/work/icon-bridge/cmd/endpoint/example-config.json")
	if err != nil {
		log.Fatal(err)
	}
	cfgPerMap := map[chain.ChainType]*chain.ChainConfig{}
	for _, ch := range cfg.Chains {
		cfgPerMap[ch.Name] = ch
	}

	l := log.New()
	log.SetGlobalLogger(l)
	fmt.Println("New")
	ug, err := New(l, map[chain.ChainType]int{chain.HMNY: 1, chain.ICON: 1}, cfgPerMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UG Start")
	ug.Start(context.TODO())
	fmt.Println("register")
	tf := DefaultTaskFunctions["DemoTransaction"]
	if err := ug.RegisterTestUnit(map[chain.ChainType]int{chain.ICON: 1, chain.HMNY: 1}, tf, false); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * time.Duration(10))
	if err := ug.RegisterTestUnit(map[chain.ChainType]int{chain.ICON: 1, chain.HMNY: 1}, tf, false); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wait")
	time.Sleep(time.Second * time.Duration(3000))
}
