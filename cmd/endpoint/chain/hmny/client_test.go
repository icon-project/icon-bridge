package hmny

import (
	"math/big"
	"testing"

	"github.com/icon-project/icon-bridge/common/log"
)

const (
	god_wallet_file              = "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.god.wallet.json"
	god_wallet_secret            = ""
	god_wallet_addr              = "0xa5241513da9f4463f1d4874b548dfbac29d91f34"
	god_wallet_priv_key          = ""
	wallet_addr_1_file           = "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.demo.wallet.json"
	wallet_addr_1_secret         = "1234"
	wallet_addr_1                = "0x8fc668275b4fa032342ea3039653d841f069a83b"
	wallet_addr_1_priv_key       = ""
	wallet_addr_2_file           = "/home/manish/go/src/work/icon-bridge/cmd/endpointtest/chain/hmny/tmp/hmny_demo.json"
	wallet_addr_2_secret         = "secret"
	wallet_addr_2                = "0x606f95a0d893ab26aa3e7dd9ce33530bca0e6dbf"
	wallet_addr_2_priv_key       = ""
	btp_icon_demo_wallet_address = "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	networkID                    = "0x6357d2e0"
)

func newLocalClient() (*client, error) {
	const (
		URL                          = "http://127.0.0.1:9500"
		btp_hmny_erc20               = "0xb54f5e97972AcF96470e02BE0456c8DB2173f33a"
		btp_hmny_nativecoin_bsh_core = "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95"
		btp_hmny_token_bsh_proxy     = "0x48cacC89f023f318B4289A18aBEd44753a127782"
	)
	l := log.New()
	log.SetGlobalLogger(l)
	cAddr := &contractAddress{
		btp_hmny_erc20:               btp_hmny_erc20,
		btp_hmny_nativecoin_bsh_core: btp_hmny_nativecoin_bsh_core,
		btp_hmny_token_bsh_proxy:     btp_hmny_token_bsh_proxy,
	}
	return newClient(URL, l, cAddr, networkID)
}

func TestGetWalletAndPrivKey(t *testing.T) {
	w, _, err := GetWalletFromPrivKey(wallet_addr_2_priv_key)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn(w.Address())
}

func TestGetHmnyBalance(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	addrs := []string{god_wallet_addr, wallet_addr_1, wallet_addr_2}
	for _, addr := range addrs {
		if val, err := cl.GetHmnyBalance(addr); err != nil {
			log.Fatal(addr, err)
		} else {
			log.Warnf("%v: %v", addr, val)
		}
	}
}

func TestTransferHmnyOne(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(500000000000000)
	hash, err := cl.TransferHmnyOne(wallet_addr_1_priv_key, *amount, wallet_addr_2)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("Hash ", hash)
}

func TestGetHmnyErc20Balance(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	addrs := []string{god_wallet_addr, wallet_addr_1, wallet_addr_2}
	for _, addr := range addrs {
		if val, err := cl.GetHmnyErc20Balance(addr); err != nil {
			log.Fatal(addr, err)
		} else {
			log.Warnf("%v: %v", addr, val)
		}
	}
}

func TestTransferErc20(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(500000000000000)
	hash, err := cl.TransferErc20(wallet_addr_1_priv_key, *amount, wallet_addr_2)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("Hash ", hash)
}

func TestTransferOneToIcon(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(1500000000000000)
	hash, err := cl.TransferOneToIcon(wallet_addr_1_priv_key, btp_icon_demo_wallet_address, *amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("Hash ", hash)
}

func TestApproveHmnyNativeBSHCoreToAccessICX(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(500000000000000)
	hash, amt, err := cl.ApproveHmnyNativeBSHCoreToAccessICX(wallet_addr_1_priv_key, *amount)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("Hash ", hash, " Amount ", *amt)
}

func TestTransferWrappedICXFromHmnyToIcon(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(500000000000000)
	hash, err := cl.TransferWrappedICXFromHmnyToIcon(wallet_addr_1_priv_key, *amount, btp_icon_demo_wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("Hash ", hash)
}

func TestTransferERC20ToIcon(t *testing.T) {
	cl, err := newLocalClient()
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(500000000000000)
	hash, hash2, err := cl.TransferERC20ToIcon(wallet_addr_1_priv_key, *amount, btp_icon_demo_wallet_address)
	if err != nil {
		log.Fatal(err)
	}
	log.Warn("ApproveHash ", hash, "  TransferHash ", hash2)
}
