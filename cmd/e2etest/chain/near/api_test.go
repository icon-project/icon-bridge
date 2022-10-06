package near

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	RPC_URI      = "https://rpc.testnet.near.org"
	TokenGodKey  = ""
	TokenGodAddr = "btp://0x2.icon/hx23552e15bfe0cf3d8166b809b344329c2e20feaa"
	GodKey       = ""
	GodAddr      = "btp://0x2.icon/hx077ada6dd02f63b02650c5861f9f41166e45d9f1"
	DemoSrcKey   = ""
	DemoSrcAddr  = "btp://0x1.near/alice.testnet"
	DemoDstAddr  = "btp://0x2.icon/hx96f2c7524c0557f8b56d461205443367cb731e83"
	GodDstAddr   = "btp://0x1.near/e072b70f2caa18b9e8e795ce970ec48f67368055d489f14174b779594dd6a5aa"
	NID          = "0x1.near"
	BtsOwner     = "btp://0x1.near/bts.iconbridge.testnet"
)

func TestGetCoinNames(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	i := api.NativeCoin()
	t.Log(i)
	// assert.Equal(t, 1, 0)
}

// func TestGetOwners(t *testing.T) {
// 	api, err := getNewApi()
// 	if err != nil {
// 		t.Fatalf("%+v", err)
// 		return
// 	}
// 	owner, err := api.CallBTS("get_owners", nil)
// 	if err != nil {
// 		t.Fatalf("%+v", err)
// 		return
// 	}
// 	if data, ok := (owner).(types.CallFunctionResponse); ok {
// 		var r []string
// 		err = json.Unmarshal(data.Result, &r)
// 		fmt.Println(data.BlockHash)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(r)

// 		// assert.Equal(t, 1, 0)
// 	}

// }

// func TestIsUserBlackListed(t *testing.T) {
// 	rpi, err := getNewApi()
// 	if err != nil {
// 		t.Fatalf("%+v", err)
// 	}
// 	res, err := rpi.CallBTS(chain.IsUserBlackListed, []interface{}{
// 		"0x61.bsc",
// 		GodDstAddr,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println("Res ", res)
// }

func TestTransferIntraChain(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	amount := new(big.Int)
	amount.SetString("1000", 10)
	srckey := ""
	dstaddr := DemoDstAddr
	for _, coinName := range []string{"btp-0x1.near-NEAR"} {
		txnHash, err := api.Transfer(coinName, srckey, dstaddr, amount)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 3)
		t.Logf("Transaction Hash %v %v", coinName, txnHash)
		res, err := api.WaitForTxnResult(context.TODO(), txnHash)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Receipt %+v", res)
		for _, lin := range res.ElInfo {
			t.Logf("Log %+v ", lin)
		}
		if val, err := api.GetCoinBalance(coinName, dstaddr); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Balance %v", val)
		}
	}
}

func TestGetCoinBalance(t *testing.T) {
	if err := showBalance(DemoSrcAddr); err != nil {
		t.Fatalf(" %+v", err)
	}
	assert.Equal(t, 1, 0)

}

func showBalance(addr string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}
	for _, coinName := range []string{"btp-0x1.near-NEAR"} {
		res, err := api.GetCoinBalance(coinName, addr)
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res.String())
	}
	return nil
}

func getNewApi() (chain.ChainAPI, error) {
	srcEndpoint := RPC_URI
	addrToName := map[chain.ContractName]string{
		chain.BTS: "bts.iconbridge.testnet",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	return NewApi(l, &chain.Config{
		Name:              chain.NEAR,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NativeTokens:      []string{},
		WrappedCoins:      []string{"btp-0x2.icon-ICX", "btp-0x1.bsc-BNB"},
		NativeCoin:        "btp-0x1.near-NEAR",
		NetworkID:         NID,
		// GasLimit:          300000000000000,
	})
}
