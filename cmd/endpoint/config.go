package endpoint

import (
	"encoding/json"
	"os"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
)

type Config struct {
	Chains []chain.ChainConfig `json:"chains"`
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
