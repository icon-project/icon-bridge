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
	GodKey      = "1deb607f38b0bd1390df3b312a1edc11a00a34f248b5d53f4157de054f3c71ae"
	GodAddr     = "btp://0x61.bsc/0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
	DemoSrcKey  = "ce69f928c68b0b7bc198824b081cfbde60d6b1e0f1695d5aaa9d8564bb35dcb3"
	DemoSrcAddr = "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	DemoDstAddr = "btp://0x613f17.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	GodDstAddr  = "btp://0x613f17.icon/hxad8eec2e167c24020600ddf1acd4d03673d3f49b"
	BtsAddr     = "btp://0x61.bsc/0x71a1520bBb7e6072Bbf3682A60c73D63b693690A"
)

func TestApprove(t *testing.T) {

	coin := "BNB"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("1000000000000", 10)
	approveHash, err := rpi.Approve(coin, DemoSrcKey, amt)
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

	for _, coin := range []string{"BNB", "TBNB"} {
		amt := new(big.Int)
		amt.SetString("178199999998000000", 10)
		hash, err := rpi.Transfer(coin, GodKey, DemoSrcAddr, amt)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("Hash %v", hash)
	}
}

func TestTransferInterChain(t *testing.T) {

	coin := "ETH"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, coin := range []string{coin} {
		res, err := rpi.GetCoinBalance(coin, GodAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res.String())
	}
	amt := new(big.Int)
	amt.SetString("2039709999998000000", 10)
	txnHash, err := rpi.Transfer(coin, GodKey, DemoDstAddr, amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, err := rpi.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range res.ElInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v %v", lin, seq)
	}
	for _, coin := range []string{coin} {
		res, err := rpi.GetCoinBalance(coin, GodAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res.String())
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

//1018700000000000000
//2019601000000000000
//99998980000000000000000
func TestGetCoinBalance(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, coin := range []string{"TBNB", "BNB", "ETH", "ICX", "TICX"} {
		res, err := rpi.GetCoinBalance(coin, GodAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res.String())
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
	demoKeyPair, err := api.GetKeyPairs(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v ", demoKeyPair)
}

func getNewApi() (chain.ChainAPI, error) {
	ctrMap := map[chain.ContractName]string{
		chain.BTSCoreBsc:      "0x71a1520bBb7e6072Bbf3682A60c73D63b693690A",
		chain.BTSPeripheryBsc: "0x3abC8DFF0C95B8982399daCf6ED5bD7b94a40068",
		// chain.TBNBBsc:         "0xBA34F3c6893b12fF4115ACf1b4712C6E2783aD83",
	}
	coinMap := map[string]string{
		"ETH":  "0x81C0094F73123EeBd250Ab4ee1e8aA6e82A7cA6F",
		"TBNB": "0xBA34F3c6893b12fF4115ACf1b4712C6E2783aD83",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	rx, err := NewApi(l, &chain.ChainConfig{
		Name:                 chain.BSC,
		URL:                  "http://localhost:8545",
		ContractAddresses:    ctrMap,
		NativeTokenAddresses: coinMap,
		NativeCoin:           "BNB",
		NetworkID:            "0x61.bsc",
	})
	if err != nil {
		return nil, err
	}
	return rx, nil
}

func getNewClient() (*Client, error) {
	l := log.New()
	log.SetGlobalLogger(l)
	cls, err := NewClients([]string{"http://localhost:8545"}, l)
	if err != nil {
		return nil, err
	}
	return cls[0], nil
}
