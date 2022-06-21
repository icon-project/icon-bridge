package endpoint

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/common/log"
)

type endpoint struct {
	iconClient chain.Client
	hmnyClient chain.Client
	cfg        *Config
}

func NewService() (*endpoint, error) {
	cfg, err := loadConfig("/home/manish/go/src/work/icon-bridge/cmd/endpoint/example-config.json")
	if err != nil {
		return nil, err
	}
	l := log.New()
	log.SetGlobalLogger(l)

	e := &endpoint{cfg: cfg}
	for _, chain := range cfg.Chains {
		if chain.Name == "ICON" {
			e.iconClient, _ = icon.New(chain.URL, l, chain.ConftractAddresses, chain.NetworkID)
		} else if chain.Name == "HMNY" {
			e.hmnyClient, _ = hmny.New(chain.URL, l, chain.ConftractAddresses, chain.NetworkID)
		} else {
			return nil, errors.New("Chain name not among supported ones (icon, hmny)")
		}
	}
	return e, nil
}

func (e *endpoint) showBalance(iconAddress, hmnyAddress string) error {
	icx, err := e.iconClient.GetCoinBalance(iconAddress)
	if err != nil {
		return err
	}
	ethIrc2, err := e.iconClient.GetEthToken(iconAddress)
	if err != nil {
		return err
	}
	oneWrapped, err := e.iconClient.GetWrappedCoin(iconAddress)
	if err != nil {
		return err
	}

	one, err := e.hmnyClient.GetCoinBalance(hmnyAddress)
	if err != nil {
		return err
	}
	ethErc20, err := e.hmnyClient.GetEthToken(hmnyAddress)
	if err != nil {
		return err
	}
	icxWrapped, err := e.hmnyClient.GetWrappedCoin(hmnyAddress)
	if err != nil {
		return err
	}
	fmt.Printf(`Balance:
	ICON:
		ICX: %d
		ONE(Wrapped): %d
		ETH(IRC2): %d 
	HMNY:
		ONE: %d
		ICX(Wrapped): %d
		ETH(ERC20): %d
	`, icx.Uint64(), oneWrapped.Uint64(), ethIrc2.Uint64(), one.Uint64(), icxWrapped.Uint64(), ethErc20.Uint64())
	return nil
}

func (e *endpoint) runDemo() error {
	//iconDemoWalletFile := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.demo.wallet.json"
	//iconDemoWalletPass := "1234"
	// iconGodWalletFile := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.god.wallet.json"
	// iconGodWalletFilePass := "gochain"
	amount := new(big.Int)
	amount.SetString("1000000000000000", 10)
	// recepientAddress := "hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	// hmnyAddress := "btp://0x6357d2e0.hmny/0x8fc668275b4fa032342ea3039653d841f069a83b"
	// _, err := e.iconClient.TransferCoinCrossChain(iconDemoWalletFile, iconDemoWalletPass, *amount, hmnyAddress)
	// if err != nil {
	// 	return err
	// }

	// txnHash, err := e.iconClient.TransferCoin(iconGodWalletFile, iconGodWalletFilePass, *amount, recepientAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Warn(txnHash)
	// hmnyDemoWalletFile := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.demo.wallet.json"
	// hmnyDemoWalletPass := "1234"
	hmnyDemoWalletPrivKey := ""
	iconAddress := "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	txnHash, err := e.hmnyClient.TransferCoinCrossChain(hmnyDemoWalletPrivKey, *amount, iconAddress)
	if err != nil {
		return err
	}
	fmt.Println(txnHash)

	// hash, all, err := e.hmnyClient.ApproveContractToAccessCrossCoin(hmnyDemoWalletFile, hmnyDemoWalletPass, *amount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Hash ", hash)
	// fmt.Println("All ", all)
	return nil
}
