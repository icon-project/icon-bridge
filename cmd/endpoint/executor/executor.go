package executor

import (
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

type executor struct {
	godKeysPerChain map[chain.ChainType][2]string
	cfgPerChain     map[chain.ChainType]*chain.ChainConfig
	log             log.Logger
	counter         int
}

func New(l log.Logger, cfgPerChain map[chain.ChainType]*chain.ChainConfig) (ug *executor, err error) {
	getKeyPairFromFile := func(walFile string, password string) (pair [2]string, err error) {
		keyReader, err := os.Open(walFile)
		if err != nil {
			return
		}
		defer keyReader.Close()

		keyStore, err := ioutil.ReadAll(keyReader)
		if err != nil {
			return
		}
		key, err := keystore.DecryptKey(keyStore, password)
		if err != nil {
			return
		}
		privBytes := ethcrypto.FromECDSA(key.PrivateKey)
		privString := hex.EncodeToString(privBytes)
		addr := ethcrypto.PubkeyToAddress(key.PrivateKey.PublicKey)
		pair = [2]string{privString, addr.String()}
		return
	}
	ug = &executor{
		log:             l,
		cfgPerChain:     cfgPerChain,
		godKeysPerChain: make(map[chain.ChainType][2]string),
		counter:         0,
	}
	for name, cfg := range cfgPerChain {
		if pair, err := getKeyPairFromFile(cfg.GodWallet.Path, cfg.GodWallet.Password); err != nil {
			return nil, err
		} else {
			ug.godKeysPerChain[name] = pair
		}
	}
	return
}

func (ug *executor) Execute(chains []chain.ChainType, cb callBackFunc) (err error) {
	newCfg := map[chain.ChainType]*chain.ChainConfig{}
	for name, cfg := range ug.cfgPerChain {
		newCfg[name] = cfg
	}
	godKeys := map[chain.ChainType][2]string{}
	for name, keys := range ug.godKeysPerChain {
		godKeys[name] = keys
	}
	ug.log.Info("Creating clients")

	args, err := newArgs(
		ug.log.WithFields(log.Fields{"id": rand.Intn(100)}),
		newCfg, godKeys,
	)
	if err != nil {
		return
	}
	go cb(args)
	return
}
