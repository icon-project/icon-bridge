package icon

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
)

// const (
// 	btp_icon_demo_wallet_address = "hx51ecae93216cb6d58bbdc51e2b5d790da94f738a"
// 	btp_icon_god_wallet_address  = "hxff0ea998b84ab9955157ab27915a9dc1805edd35"
// 	btp_icon_step_limit          = 3500000000
// 	btp_icon_nid                 = 0x5b9a77
// 	Version                      = 3
// 	btp_icon_irc2                = "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f"
// wallet_addr_2_file   = "/home/manish/go/src/work/icon-bridge/cmd/endpointtest/chain/icon/tmp/icon_demo.json"
// wallet_addr_2_secret = "secret"
// wallet_addr_2        = "hx51ecae93216cb6d58bbdc51e2b5d790da94f738a"
// wallet_netaddr_2     = "btp://0x5b9a77.con/51ecae93216cb6d58bbdc51e2b5d790da94f738a"
// )

const (
	god_wallet_file        = "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.god.wallet.json"
	god_wallet_secret      = "gochain"
	god_wallet_priv_key    = ""
	god_wallet_addr        = "hxff0ea998b84ab9955157ab27915a9dc1805edd35"
	wallet_addr_1_file     = "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.demo.wallet.json"
	wallet_addr_1_secret   = "1234"
	wallet_addr_1_priv_key = ""
	wallet_addr_1          = "hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	wallet_netaddr_1       = "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	// btp_icon_irc2           = "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f"
	// btp_icon_nativecoin_bsh = "cx6fa46cc92fcf8e3135bc6645ef3a258b0ce73602"
	// btp_icon_token_bsh      = "cx2194782b7951d8abf26ea88204e45862e9821bbc"
)

func newLocalClient() *client {
	const URL = "http://127.0.0.1:9080/api/v3/default"
	l := log.New()
	log.SetGlobalLogger(l)
	cMap := map[string]string{
		"btp_icon_irc2":           "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f",
		"btp_icon_irc2_tradeable": "cx7831ba8969943d96c375261f8245e6e9964389c9",
		"btp_icon_nativecoin_bsh": "cx6fa46cc92fcf8e3135bc6645ef3a258b0ce73602",
		"btp_icon_token_bsh":      "cx2194782b7951d8abf26ea88204e45862e9821bbc",
	}
	cAddr := &contractAddress{}
	cAddr.FromMap(cMap)
	client, _ := newClient(URL, l, cAddr, "0x5b9a77")
	return client
}

func TestGetICXBalance(t *testing.T) {
	cl := newLocalClient()
	addrs := []string{god_wallet_addr, wallet_addr_1}
	for _, addr := range addrs {
		v, err := cl.GetICXBalance(addr)
		if err != nil {
			log.Fatal(err)
		}
		log.Warnf("%v: %v", addr, v)
	}
}

func TestTransferICX(t *testing.T) {
	// stepLimit := NewHexInt(3500000000)
	// networkID := NewHexInt(0x5b9a77)
	amount := big.NewInt(200000000000)
	cl := newLocalClient()
	_, err := cl.TransferICX(god_wallet_priv_key, *amount, wallet_addr_1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetIrc2(t *testing.T) {
	cl := newLocalClient()

	addrs := []string{god_wallet_addr, wallet_addr_1}
	for _, addr := range addrs {
		//args := map[string]string{"_owner": addr}
		val, err := cl.GetIrc2Balance(addr)
		if err != nil {
			log.Fatal(err)
		}
		log.Warnf("%v: %v", addr, *val)
	}
}

func TestTransferIrc2(t *testing.T) {
	cl := newLocalClient()
	// stepLimit := NewHexInt(3500000000)
	// networkID := NewHexInt(0x5b9a77)
	//args := map[string]string{"_to": wallet_addr_1, "_value": "0x39c121a270000"}
	_, err := cl.TransferIrc2(god_wallet_priv_key, *big.NewInt(0x39c121a270000), wallet_addr_1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestTransferICXToHarmony(t *testing.T) {
	cl := newLocalClient()
	// stepLimit := NewHexInt(3500000000)
	// networkID := NewHexInt(0x5b9a77)
	amount := big.NewInt(1016026492239872)
	hmnyAddr := "btp://0x6357d2e0.hmny/0x8fc668275b4fa032342ea3039653d841f069a83b"
	//args := map[string]string{"_to": hmnyAddr}
	_, err := cl.TransferICXToHarmony(
		wallet_addr_1_priv_key, *amount,
		hmnyAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func TestApproveIconNativeCoinBSHToAccessHmnyOne(t *testing.T) {
	cl := newLocalClient()
	// stepLimit := NewHexInt(3500000000)
	// networkID := NewHexInt(0x5b9a77)
	amount := big.NewInt(0x39c121a270000)
	//approveArgs := map[string]string{"spender": cl.contractAddress.btp_icon_nativecoin_bsh, "amount": amount}
	//allowArgs := map[string]string{"owner": wallet_addr_1, "spender": cl.contractAddress.btp_icon_nativecoin_bsh}
	_, _, err := cl.ApproveIconNativeCoinBSHToAccessHmnyOne(wallet_addr_1_priv_key, *amount)
	if err != nil {
		log.Fatal(err)
	}
}

func TestTransferWrappedOneFromIconToHmny(t *testing.T) {
	cl := newLocalClient()
	// stepLimit := NewHexInt(3500000000)
	// networkID := NewHexInt(0x5b9a77)
	amount := big.NewInt(0x39c121a270000)
	hmnyAddr := "btp://0x6357d2e0.hmny/0x8fc668275b4fa032342ea3039653d841f069a83b"
	//args := map[string]string{"_coinName": "ONE", "_value": amount, "_to": hmnyAddr}
	_, err := cl.TransferWrappedOneFromIconToHmny(wallet_addr_1_priv_key, *amount, hmnyAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func TestTransferIrc2ToHmny(t *testing.T) {
	cl := newLocalClient()
	amount := big.NewInt(0x39c121a270000)
	hmnyAddr := "btp://0x6357d2e0.hmny/0x8fc668275b4fa032342ea3039653d841f069a83b"
	//ircArgs := map[string]string{"_to": cl.contractAddress.btp_icon_token_bsh, "_value": amount}
	//tokenArgs := map[string]string{"tokenName": "ETH", "value": amount, "to": hmnyAddr}
	_, _, err := cl.TransferIrc2ToHmny(wallet_addr_1_priv_key, *amount, hmnyAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStringToBigInt(t *testing.T) {
	amount := new(big.Int)
	amount.SetString("100000000000000000000000", 10)
	a := fmt.Sprintf("0x%x", amount)
	fmt.Println(a, HexInt(a), amount.String(), intconv.FormatBigInt(amount))
}

func TestHexBytes(t *testing.T) {
	var a HexBytes = "0x123abc"
	ab, _ := a.Value()
	fmt.Println(a, hexutil.Encode(ab[:]))
}

func TestHexInt(t *testing.T) {
	nid := "0x61235"
	fmt.Println(HexInt(nid))
}
