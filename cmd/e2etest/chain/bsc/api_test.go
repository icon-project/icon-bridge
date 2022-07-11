package bsc

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

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
	//[ce69f928c68b0b7bc198824b081cfbde60d6b1e0f1695d5aaa9d8564bb35dcb3 0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202]
}

func TestGetCoinBalance(t *testing.T) {
	//GOD 1deb607f38b0bd1390df3b312a1edc11a00a34f248b5d53f4157de054f3c71ae 0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D
	rpi, err := getNewRequestAPI()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.getWrappedCoinBalance("TBNB", "0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf(" %v", res)
}

func getNewRequestAPI() (*requestAPI, error) {
	l := log.New()
	log.SetGlobalLogger(l)
	ctrMap := map[chain.ContractName]string{
		chain.BTSCoreBsc:      "0x71a1520bBb7e6072Bbf3682A60c73D63b693690A",
		chain.BTSPeripheryBsc: "0x3abC8DFF0C95B8982399daCf6ED5bD7b94a40068",
		chain.TBNBBsc:         "0xBA34F3c6893b12fF4115ACf1b4712C6E2783aD83",
	}
	rpi, err := newRequestAPI("http://localhost:8545", l, ctrMap, "0x61.bsc")
	if err != nil {
		return nil, err
	}
	return rpi, nil
}

func getNewApi() (chain.ChainAPI, error) {
	ctrMap := map[chain.ContractName]string{
		chain.BTSCoreBsc:      "0x71a1520bBb7e6072Bbf3682A60c73D63b693690A",
		chain.BTSPeripheryBsc: "0x3abC8DFF0C95B8982399daCf6ED5bD7b94a40068",
		chain.TBNBBsc:         "0xBA34F3c6893b12fF4115ACf1b4712C6E2783aD83",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	rx, err := NewApi(l, &chain.ChainConfig{Name: chain.HMNY, URL: "http://localhost:8545", ConftractAddresses: ctrMap, NetworkID: "0x61.bsc"})
	if err != nil {
		return nil, err
	}
	return rx, nil
}

func TestTransferIntraChain(t *testing.T) {
	senderKey := "1deb607f38b0bd1390df3b312a1edc11a00a34f248b5d53f4157de054f3c71ae"
	//senderAddress := "0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
	recepientAddress := "0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	rpi, err := getNewRequestAPI()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("50000000000000000", 10)
	hash, err := rpi.transferNativeIntraChain(senderKey, recepientAddress, *amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("Hash %v", hash)
}

func TestTransferInterChain(t *testing.T) {
	senderKey := "ce69f928c68b0b7bc198824b081cfbde60d6b1e0f1695d5aaa9d8564bb35dcb3"
	//senderAddress := "0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	recepientAddress := "btp://0xdf6463.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	rpi, err := getNewRequestAPI()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("1000000000000000", 10)
	_, _, err = rpi.approveCoin("ICX", senderKey, *amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hash, _, err := rpi.transferWrappedCrossChain("ICX", senderKey, recepientAddress, *amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("Hash  %v", hash)
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
