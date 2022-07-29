package bsc

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	TokenGodKey  = "d4901a43dc4cee775fed483636a57ad7b807e77c62134d3cb2d1411d90b072dc"
	TokenGodAddr = "btp://0x61.bsc/0x210730B1f5B9C4A02dF0808093aC5E72676cF70c"
	NID          = "0x61.bsc"
	RPC_URI      = "https://data-seed-prebsc-1-s1.binance.org:8545"
	GodKey       = "541a205a7d3119e9b617b1023d9c874db572134d50b5f1ef2590bc5e5143dc2c"
	GodAddr      = "btp://0x61.bsc/0xDf9e6205Ac201c8a11082842857C6f7673a8246e"
	BtsAddr      = "btp://0x61.bsc/0x9F90806DBDaA783766483d2D24b431CFFB793eEb"
	GodDstAddr   = "btp://0x2.icon/hxc86452374f94bd8db99f703bb1fc3fad2f7b2024"

	DemoDstAddr = "btp://0x2.icon/hx6d338536ac11a0a2db06fb21fe8903e617a6764d"
	DemoSrcKey  = "a851faf7310664601b9396e2e3e45e36456f5052c537a8354229ec9059255d59"
	DemoSrcAddr = "btp://0x61.bsc/0xDf9e6205Ac201c8a11082842857C6f7673a8246e"
)

// const (
// 	RPC_URI     = "http://localhost:8545"
// 	GodKey      = "1deb607f38b0bd1390df3b312a1edc11a00a34f248b5d53f4157de054f3c71ae"
// 	GodAddr     = "btp://0x61.bsc/0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
// 	DemoSrcKey  = "ce69f928c68b0b7bc198824b081cfbde60d6b1e0f1695d5aaa9d8564bb35dcb3"
// 	DemoSrcAddr = "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
// 	DemoDstAddr = "btp://0x613f17.icon/hx0000000000000000000000000000000000000000"
// 	GodDstAddr  = "btp://0x613f17.icon/hxad8eec2e167c24020600ddf1acd4d03673d3f49b"
// 	BtsAddr     = "btp://0x61.bsc/0x71a1520bBb7e6072Bbf3682A60c73D63b693690A"
// )

func TestApprove(t *testing.T) {

	coin := "USDC"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("995000", 10)
	approveHash, err := rpi.Approve(coin, GodKey, amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.WaitForTxnResult(context.TODO(), approveHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Hash %v Receipt %+v", approveHash, res.Raw)
}

func TestTransferIntraChain(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, coin := range []string{"DUM"} {
		amt := new(big.Int)
		amt.SetString("100000000000000", 10)
		hash, err := rpi.Transfer(coin, GodKey, BtsAddr, amt)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("Hash %v", hash)
	}
}

func TestTransferInterChain(t *testing.T) {

	coin := "USDC"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	srcKey := GodKey
	srcAddr := GodAddr
	dstAddr := GodDstAddr
	for _, coin := range []string{coin} {
		res, err := rpi.GetCoinBalance(coin, srcAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res)
	}

	amt := new(big.Int)
	amt.SetString("995000", 10)
	txnHash, err := rpi.Transfer(coin, srcKey, dstAddr, amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	time.Sleep(time.Second * 2)
	res, err := rpi.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range res.ElInfo {
		//seq, _ := lin.GetSeq()
		t.Logf("Log %+v ", lin)
	}
	for _, coin := range []string{coin} {
		res, err := rpi.GetCoinBalance(coin, srcAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res)
	}
}

func TestBatchTransfer(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	coins := []string{"TBNB", "ETH", "BNB"}
	amount := big.NewInt(100000000000000)
	largeAmt := new(big.Int)
	largeAmt.SetString("26271926117961986739", 10)
	amounts := []*big.Int{amount, largeAmt, amount}
	for i, coin := range coins {

		if coin == rpi.NativeCoin() {
			continue
		}
		approveHash, err := rpi.Approve(coin, GodKey, amounts[i])
		if err != nil {
			t.Fatalf("%+v", err)
		}
		res, err := rpi.WaitForTxnResult(context.TODO(), approveHash)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 1 {
			t.Fatalf("Approve StatusCode not 1 for %vth coin %v \n %v", i, coin, res.Raw)
		}
	}
	hash, err := rpi.TransferBatch(coins, GodKey, GodDstAddr, amounts)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Hash ", hash)
	res, err := rpi.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 1 {
		t.Fatalf("StatusCode not 1 for Batch Transfer \n %+v", res.Raw)
	}
}

func TestGetCoinBalance(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, coin := range []string{"BNB", "BUSD", "USDT", "USDC", "BTCB", "ETH", "DUM", "ICX", "sICX", "bnUSD"} {
		res, err := rpi.GetCoinBalance(coin, GodAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res)
	}
}

func TestReceiver(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", api)
	}
	_, _, err = api.Subscribe(context.TODO())
	if err != nil {
		t.Fatalf("%+v", api)
	}
	time.Sleep(time.Hour)
}

func TestGetKeyPair(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	demoKeyPair, err := api.GetKeyPairs(3)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v ", demoKeyPair)
}

func getNewApi() (chain.ChainAPI, error) {
	ctrMap := map[chain.ContractName]string{
		chain.BTS:          "0x9F90806DBDaA783766483d2D24b431CFFB793eEb",
		chain.BTSPeriphery: "0x94D9842507AAbB4D7ce010206f662b44efA8496F",
	}

	l := log.New()
	log.SetGlobalLogger(l)
	rx, err := NewApi(l, &chain.Config{
		Name:              chain.BSC,
		URL:               RPC_URI,
		ContractAddresses: ctrMap,
		NativeTokens:      []string{"BUSD", "USDT", "USDC", "BTCB", "ETH", "DUM"},
		WrappedCoins:      []string{"ICX", "sICX", "bnUSD"},
		NativeCoin:        "BNB",
		NetworkID:         "0x61.bsc",
		GasLimit:          5000000,
	})
	if err != nil {
		return nil, err
	}
	return rx, nil
}
