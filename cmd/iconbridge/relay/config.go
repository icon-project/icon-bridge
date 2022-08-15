package relay

import (
	"encoding/json"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	DefaultKeyPassword = "gochain"
)

type Config struct {
	Relays []*RelayConfig `json:"relays"`
}

type RelayConfig struct {
	Name string    `json:"name"`
	Src  SrcConfig `json:"src"`
	Dst  DstConfig `json:"dst"`
}

type ChainConfig struct {
	Address  chain.BTPAddress `json:"address"`
	Endpoint []string         `json:"endpoint"`
	Options  json.RawMessage  `json:"options,omitempty"`
	// Options  map[string]interface{} `json:"options,omitempty"`
}

type SrcConfig struct {
	ChainConfig `json:",squash"`
	Offset      uint64 `json:"offset"`
}

type DstConfig struct {
	ChainConfig `json:",squash"`

	KeyStore    json.RawMessage `json:"key_store"`
	KeyPassword string          `json:"key_password"`

	// AWS
	AWSSecretName string `json:"aws_secret_name,omitempty"`
	AWSRegion     string `json:"aws_region,omitempty"`

	// TxSizeLimit
	// is the maximum size of a transaction in bytes
	TxDataSizeLimit uint64 `json:"tx_data_size_limit"`
}

func (cfg *DstConfig) Wallet() (wallet.Wallet, error) {
	keyStore, password, err := cfg.resolveKeyStore()
	if err != nil {
		return nil, err
	}
	return wallet.DecryptKeyStore(keyStore, password)
}

func (cfg *DstConfig) resolveKeyStore() (json.RawMessage, []byte, error) {
	if cfg.AWSSecretName != "" && cfg.AWSRegion != "" {
		result, err := wallet.GetSecret(cfg.AWSSecretName, cfg.AWSRegion)
		if err != nil {
			return nil, nil, err
		}
		if result != "" {
			var w struct {
				KeyStore json.RawMessage `json:"key_store"`
				Secret   string          `json:"secret"`
			}
			err = json.Unmarshal([]byte(result), &w)
			if err != nil {
				return nil, nil, err
			}
			return w.KeyStore, []byte(w.Secret), nil
		}
	}
	if cfg.KeyPassword == "" {
		return cfg.KeyStore, []byte(DefaultKeyPassword), nil
	}
	return cfg.KeyStore, []byte(cfg.KeyPassword), nil
}
