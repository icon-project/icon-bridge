package near

import (
	"context"
	"fmt"
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
	TokenGodAddr = "btp://0x7.icon/hx23552e15bfe0cf3d8166b809b344329c2e20feaa"
	GodKey       = ""
	GodAddr      = "btp://0x7.icon/hx077ada6dd02f63b02650c5861f9f41166e45d9f1"
	DemoSrcKey   = ""
	DemoSrcAddr  = "btp://0x2.near/alice.testnet"
	DemoDstAddr  = "btp://0x7.icon/hx96f2c7524c0557f8b56d461205443367cb731e83"
	GodDstAddr   = "btp://0x2.near/e072b70f2caa18b9e8e795ce970ec48f67368055d489f14174b779594dd6a5aa"
	NID          = "0x1.near"
	BtsOwner     = "btp://0x2.near/bts.iconbridge-6.testnet"
)

func TestGetCoinNames(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	i := api.NativeCoin()
	t.Log(i)
	//assert.Equal(t, 1, 0)
}

func TestTransferIntraChain(t *testing.T) {
	_api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	amount := new(big.Int)
	amount.SetString("1000000000", 10)
	privKey := "2TPaWkb7zjkF6PoHFtWcSG55Ckc3V9qPtbyww7HHzytwo5EEbkZAVeaUdwjxvpFLt6DhmSqZAmJAMdew1V5rk9fb"
	srckey := "alice.testnet"
	dstaddr := DemoDstAddr
	btpaddr := _api.GetBTPAddress(srckey)

	// go func() {
	// 	if sinkChan, errChan, err := _api.Subscribe(context.Background()); err != nil {
	// 		panic(err)
	// 	} else {
	// 		for {
	// 			select {
	// 			case err := <-errChan:
	// 				panic(err)
	// 			case msg := <-sinkChan:
	// 				t.Logf("\nMessage %+v\n", msg)
	// 				_api.(*api).StopSubscriptionMethod()
	// 			}
	// 		}
	// 	}
	// }()

	for _, coinName := range []string{"btp-0x2.near-NEAR"} {
		txnHash, err := _api.Transfer(coinName, privKey, dstaddr, amount)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 3)
		t.Logf("Transaction Hash %v %v", coinName, txnHash)

		time.Sleep(time.Second * 3)
		res, err := _api.WaitForTxnResult(context.TODO(), txnHash)
		if err != nil {
			t.Log(err)
		}
		// t.Logf("Receipt %+v", res)
		for _, lin := range res.ElInfo {
			t.Logf("Log %+v ", lin)
		}
		if val, err := _api.GetCoinBalance(coinName, btpaddr); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Balance %v", val)
		}
	}
	// assert.Equal(t, 1, 0)
}

func TestGetCoinBalance(t *testing.T) {
	if err := showBalance(DemoSrcAddr); err != nil {
		t.Fatalf(" %+v", err)
	}
	//assert.Equal(t, 1, 0)

}

func showBalance(addr string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}
	for _, coinName := range []string{"btp-0x2.near-NEAR"} {
		res, err := api.GetCoinBalance(coinName, addr)
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res.String())
	}
	return nil
}

// func TestReceiver(t *testing.T) {
// 	recv, err := getNewApi()
// 	if err != nil {
// 		panic(err)
// 	}
// 	// recv.WatchForTransferStart(1, "ICX", 10)
// 	// recv.WatchForTransferReceived(1, "TONE", 8)
// 	// recv.WatchForTransferEnd(1, "ICX", 10)

// 	go func() {
// 		if sinkChan, errChan, err := recv.Subscribe(context.Background()); err != nil {
// 			panic(err)
// 		} else {
// 			for {
// 				select {
// 				case err := <-errChan:
// 					panic(err)
// 				case msg := <-sinkChan:
// 					t.Logf("\nMessage %+v\n", msg)
// 					recv.(*api).StopSubscriptionMethod()
// 				}
// 			}
// 		}
// 	}()
// 	time.Sleep(time.Second * 3000)
// }

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
	assert.Equal(t, 1, 0)
}

// Need to check GetBlackListedUsers method
func TestGetBlackListedUsers(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.GetBlackListedUsers("0x2.near", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
	// assert.Equal(t, 1, 0)
}

// Same as GetBlackListedUsers, type conversion is not working
func TestIsUserBlackListed(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.IsUserBlackListed("0x2.near", "e072b70f2caa18b9e8e795ce970ec48f67368055d489f14174b779594dd6a5aa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}

func TestGetFeeRatio(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fn, ff, err := rpi.GetFeeRatio("btp-0x2.near-NEAR")
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(fn, ff)
	// assert.Equal(t, 1, 0)
}

func TestGetAccumulatedFees(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fees, err := rpi.GetAccumulatedFees()
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(fees)
	// assert.Equal(t, 1, 0)
}

// type conversion fails
func TestGetTokenLimit(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, coin := range []string{"btp-0x2.icon-bnUSD", "btp-0x2.icon-sICX", "btp-0x2.icon-ICX"} {
		res, err := rpi.GetTokenLimit(coin)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Res coin %v %v", coin, res)
	}
}

func getNewApi() (chain.ChainAPI, error) {
	srcEndpoint := RPC_URI
	addrToName := map[chain.ContractName]string{
		chain.BTS: "bts.iconbridge-6.testnet",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	return NewApi(l, &chain.Config{
		Name:              chain.NEAR,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NativeTokens:      []string{},
		WrappedCoins:      []string{"btp-0x7.icon-ICX"},
		NativeCoin:        "btp-0x2.near-NEAR",
		NetworkID:         NID,
		// GasLimit:          300000000000000,
	})
}
