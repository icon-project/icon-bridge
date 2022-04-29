package relay

import (
	"encoding/json"

	"github.com/icon-project/btp/cmd/btpsimple/chain"
	"github.com/icon-project/btp/common/wallet"
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
	Address  chain.BTPAddress       `json:"address"`
	Endpoint []string               `json:"endpoint"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type SrcConfig struct {
	ChainConfig `json:",squash"`
	Offset      uint64 `json:"offset"`
}

type DstConfig struct {
	ChainConfig `json:",squash"`

	KeyStore    json.RawMessage `json:"key_store"`
	KeyPassword string          `json:"key_password"`

	// TxSizeLimit
	// is the maximum size of a transaction in bytes
	TxDataSizeLimit uint64 `json:"tx_data_size_limit"`
}

func (cfg *DstConfig) Wallet() (wallet.Wallet, error) {
	pw, err := cfg.resolvePassword()
	if err != nil {
		return nil, err
	}
	return wallet.DecryptKeyStore(cfg.KeyStore, pw)
}

func (cfg *DstConfig) resolvePassword() ([]byte, error) {
	if cfg.KeyPassword == "" {
		return []byte(DefaultKeyPassword), nil
	}
	return []byte(cfg.KeyPassword), nil
}
