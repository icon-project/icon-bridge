package near

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	RPC_URI      = "https://rpc.testnet.near.org"
	TokenGodKey  = ""
	TokenGodAddr = ""
	GodKey       = ""
	GodAddr      = ""
	DemoSrcKey   = ""
	DemoSrcAddr  = ""
	DemoDstAddr  = ""
	GodDstAddr   = ""
	NID          = "0x1.near"
	BtsOwner     = ""
)

func TestGetCoinNames(t *testing.T) {
	_, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestGetOwners(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
		return
	}
	owner, err := api.CallBTS("get_owners", nil)
	if err != nil {
		t.Fatalf("%+v", err)
		return
	}
	if data, ok := (owner).(types.CallFunctionResponse); ok {
		var r []string
		err = json.Unmarshal(data.Result, &r)
		fmt.Println(data.BlockHash)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(r)

		// assert.Equal(t, 1, 0)
	}

}

func TestIsUserBlackListed(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.CallBTS(chain.IsUserBlackListed, []interface{}{
		"0x61.bsc",
		GodDstAddr,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}

func TestTransferIntraChain(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	amount := new(big.Int)
	amount.SetString("10000000000000000000", 10)
	srckey := TokenGodKey
	dstaddr := DemoSrcAddr
	for _, coinName := range []string{"ICX", "bnUSD"} { // need to change
		txnHash, err := api.Transfer(coinName, srckey, dstaddr, amount)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 3)
		t.Logf("Transaction Hash %v", txnHash)
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

func getNewApi() (chain.ChainAPI, error) {
	srcEndpoint := RPC_URI
	addrToName := map[chain.ContractName]string{
		chain.BTS: "bts.iconbridge.testnet",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	rx, err := NewApi(l, &chain.Config{
		Name:              chain.NEAR,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NativeTokens:      []string{},
		WrappedCoins:      []string{"btp-0x2.icon-ICX"},
		NativeCoin:        "btp-0x1.near-NEAR",
		NetworkID:         "0x1.near",
		GasLimit:          300000000000000,
	})
	if err != nil {
		return nil, err
	}
	return rx, nil
}
