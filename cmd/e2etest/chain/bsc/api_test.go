package bsc

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	RPC_URI      = "http://localhost:8545"
	TokenGodKey  = ""
	TokenGodAddr = "btp://0x61.bsc/0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
	GodKey       = ""
	GodAddr      = "btp://0x61.bsc/0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
	DemoSrcKey   = ""
	DemoSrcAddr  = "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	DemoDstAddr  = "btp://0x5b9a77.icon/hx0000000000000000000000000000000000000000"
	GodDstAddr   = "btp://0x5b9a77.icon/hxad8eec2e167c24020600ddf1acd4d03673d3f49b"
	BtsAddr      = "btp://0x61.bsc/0x71a1520bBb7e6072Bbf3682A60c73D63b693690A"
)

func TestApprove(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("100000000000000000", 10)
	for _, coin := range []string{"BUSD", "USDT", "USDC", "BTCB", "ETH"} {
		// coin := "USDC"
		approveHash, err := rpi.Approve(coin, TokenGodKey, amt)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		res, err := rpi.WaitForTxnResult(context.TODO(), approveHash)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Hash %v Receipt %+v", approveHash, res.Raw)
	}
}

func TestTransferIntraChain(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	senderKey := TokenGodKey
	dstAddr := "btp://0x61.bsc/0x8Bde22A645051B8772E4d6d9125Bb0B77EE2Ca0d"
	amt := new(big.Int)
	amt.SetString("5000000000000000000", 10)
	for _, coin := range []string{"BNB"} {
		hash, err := rpi.Transfer(coin, senderKey, dstAddr, amt)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("Hash %v", hash)
		time.Sleep(time.Second * 3)
		res, err := rpi.WaitForTxnResult(context.TODO(), hash)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Receipt %+v", res)
		for _, lin := range res.ElInfo {
			//seq, _ := lin.GetSeq()
			t.Logf("Log %+v ", lin)
		}
	}

}

func TestTransferInterChain(t *testing.T) {
	//"BUSD", "USDT", "USDC", "BTCB", "ETH"
	coin := "BNB"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	srcKey := TokenGodKey
	srcAddr := TokenGodAddr
	dstAddr := DemoSrcAddr
	for _, coin := range []string{coin} {
		res, err := rpi.GetCoinBalance(coin, srcAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res)
	}

	amt := new(big.Int)
	amt.SetString("100000000000000000", 10)
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
	coins := []string{"BUSD", "USDT", "USDC", "BTCB", "ETH"}

	largeAmt := new(big.Int)
	largeAmt.SetString("1000000000000000000000", 10)
	amounts := []*big.Int{largeAmt, largeAmt, largeAmt, largeAmt, largeAmt}
	for i, coin := range coins {
		fmt.Println("coin", coin)
		if coin == rpi.NativeCoin() {
			continue
		}
		approveHash, err := rpi.Approve(coin, TokenGodKey, amounts[i])
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
	hash, err := rpi.TransferBatch(coins, TokenGodKey, GodDstAddr, amounts)
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
		res, err := rpi.GetCoinBalance(coin, DemoSrcAddr)
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
	demoKeyPair, err := api.GetKeyPairs(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v ", demoKeyPair)
}

func getNewApi() (chain.ChainAPI, error) {
	ctrMap := map[chain.ContractName]string{
		chain.BTS:          "0x71a1520bBb7e6072Bbf3682A60c73D63b693690A",
		chain.BTSPeriphery: "0x3abC8DFF0C95B8982399daCf6ED5bD7b94a40068",
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
		GasLimit:          8000000,
	})
	if err != nil {
		return nil, err
	}
	return rx, nil
}

/*
func TestIsOwner(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.CallBTS(chain.IsOwner, []interface{}{"0x8Bde22A645051B8772E4d6d9125Bb0B77EE2Ca0d"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Res ", res)
}

func TestGetTokenLimit(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, coin := range []string{"BNB", "BUSD", "USDT", "USDC", "BTCB", "ETH", "DUM", "ICX", "sICX", "bnUSD"} {
		res, err := rpi.CallBTS(chain.GetTokenLimit, []interface{}{coin})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Res coin ", res)
	}
}

func TestIsUserBlackListed(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.CallBTS(chain.IsUserBlackListed, []interface{}{
		"0x61.bsc",
		"0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}

func TestCheckTransferRestrictions(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10)
	res, err := rpi.CallBTS(chain.CheckTransferRestrictions, []interface{}{
		"0x61.bsc",
		"sICX",
		DemoSrcAddr,
		amount,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}
*/
