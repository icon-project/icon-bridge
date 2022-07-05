package executor_test

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/executor"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestExecutor(t *testing.T) {
	type Config struct {
		Chains []*chain.ChainConfig `json:"chains"`
	}
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
		t.Fatal(err)
	}
	cfgPerMap := map[chain.ChainType]*chain.ChainConfig{}
	for _, ch := range cfg.Chains {
		cfgPerMap[ch.Name] = ch
	}

	l := log.New()
	log.SetGlobalLogger(l)
	ex, err := executor.New(l, cfgPerMap)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)

	ex.Subscribe(ctx)

	amount := new(big.Int)
	amount.SetString("10000000000000000000", 10)
	for tsi, ts := range executor.TestScripts {
		l.Info("Running TestScript SN.", tsi)
		go func() {
			err = ex.Execute(ctx, chain.ICON, chain.HMNY, amount, ts)
			if err != nil {
				log.Errorf("%+v", err)
			}
		}()
		time.Sleep(time.Second * 5)
	}
	defer func() {
		cancel()
	}()
	<-ex.Done()
	cancel()
	time.Sleep(time.Second * 3)
	l.Info("Exit")
}

/*
func (ug *executor) createAccounts(accountMap map[chain.ChainType]int) (map[chain.ChainType][]string, error) {
	resMap := map[chain.ChainType][]string{}
	for name, count := range accountMap {
		resMap[name] = make([]string, count)
		for i := 0; i < count; i++ {
			privKey, err := ethcrypto.GenerateKey()
			if err != nil {
				return nil, err
			}
			resMap[name][i] = hex.EncodeToString(ethcrypto.FromECDSA(privKey))
			// pubKey, _ := crypto.ParsePublicKey(pub)
			// addr := common.NewAccountAddressFromPublicKey(pubKey).String()
			if err != nil {
				return nil, errors.Wrap(err, "Unmarshal Public Key")
			}
		}
	}
	return resMap, nil
}
*/
