package backend

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
		if ch.Name == chain.ICON {
			cfgPerMap[ch.Name] = ch
		}
	}
	be, err := New(l, cfgPerMap)
	if err != nil {
		log.Fatal(err)
	}
	be.Start(context.TODO())
	log.Warn("Wait")
	time.Sleep(time.Hour)
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

func TestDecodeRLP(t *testing.T) {

	v := &[]AssetDetailNative{}
	a, err := rlpDecodeHex("0x5472616e736665722053756363657373", v)
	if err != nil {
		log.Fatal(err)
	}
	v, ok := a.(*[]AssetDetailNative)

	fmt.Println(v, ok)
	m := &[]AssetDetailToken{}
	tb, err := rlpDecodeHex("0xd6d583455448880dbd2fc137a30000872386f26fc10000", m)
	if err != nil {
		log.Fatal(err)
	}
	w, ok := tb.(*[]AssetDetailToken)
	fmt.Println(w, ok)
	aa, ab, _, err := rlpDecodeData("0xd6d583455448880dbd2fc137a30000872386f26fc10000", 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*aa, *ab)
}
