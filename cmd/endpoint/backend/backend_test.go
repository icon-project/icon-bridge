package backend

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestBackendTransferRequests(t *testing.T) {
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
	defer cancel()
	amt := new(big.Int)
	amt.SetString("100000000000000000", 10)
	txHash, err := be.Transfer(&chainAPI.RequestParam{
		FromChain:   chain.ICON,
		ToChain:     chain.HMNY,
		SenderKey:   "89053eee5a5f524097bf449b9e544dcc78066e548b9143bf48e88b1495d1bae8",
		FromAddress: "hx27407bd352d28064b20dde9411486f52849bf82d",
		ToAddress:   "0x9C35e844b0e3c3d6d50e7FEe2E77F6b7D0Ed4ADB",
		Amount:      *amt,
		Token:       chainAPI.ICXToken,
	})
	if err != nil {
		t.Fatalf("Err %+v", err)
	}
	t.Logf("TxHash of transaction is %v", txHash)
}

func TestBackendSubscription(t *testing.T) {
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
	defer cancel()

	if err := be.Start(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	t.Log("Send Transaction Now++++++++++++++++++++++++++++++")

	amt := new(big.Int)
	amt.SetString("100000000000000000", 10)
	txHash, err := be.Transfer(&chainAPI.RequestParam{
		FromChain:   chain.ICON,
		ToChain:     chain.HMNY,
		SenderKey:   "89053eee5a5f524097bf449b9e544dcc78066e548b9143bf48e88b1495d1bae8",
		FromAddress: "hx27407bd352d28064b20dde9411486f52849bf82d",
		ToAddress:   "0x9C35e844b0e3c3d6d50e7FEe2E77F6b7D0Ed4ADB",
		Amount:      *amt,
		Token:       chainAPI.ICXToken,
	})
	if err != nil {
		t.Fatalf("Err %+v", err)
	}
	t.Log("Transaction hash that was sent ", txHash)

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
