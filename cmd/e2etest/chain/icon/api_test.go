package icon_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	RPC_URI      = "https://lisbon.net.solidwallet.io/api/v3/icon_dex"
	TokenGodKey  = "c6e4954a60ed41a76d96b4b5eebead9d13c697de176553d96fdbd6faebd01838"
	TokenGodAddr = "btp://0x2.icon/hx077ada6dd02f63b02650c5861f9f41166e45d9f1"
	GodKey       = ""
	GodAddr      = "btp://0x2.icon/hx077ada6dd02f63b02650c5861f9f41166e45d9f1"
	DemoSrcKey   = ""
	DemoSrcAddr  = "btp://0x2.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	DemoDstAddr  = "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	GodDstAddr   = "btp://0x61.bsc/0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
	NID          = "0x2.icon"
	BtsOwner     = "btp://0x2.icon/hx1a2aeb3a100f2179846307095b82aa8ace43ca9d"
)

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
	for _, coinName := range []string{"btp-0x2.icon-ICX", "btp-0x2.icon-bnUSD"} {
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

	return
}

func TestReclaim(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hash, err := rpi.Reclaim("ETH", DemoSrcKey, big.NewInt(1333035204))
	if err != nil {
		t.Fatal(err)
	}
	res, err := rpi.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Hash %v Receipt %+v", hash, res.Raw)
	showBalance(DemoSrcAddr)
}

func TestApprove(t *testing.T) {
	ownerKey := TokenGodKey
	ownerAddr := TokenGodAddr
	showBalance(ownerAddr)
	coin := "btp-0x61.bsc-BUSD"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("4300000000000000000", 10)
	approveHash, err := rpi.Approve(coin, ownerKey, amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	time.Sleep(time.Second * 2)
	res, err := rpi.WaitForTxnResult(context.TODO(), approveHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Hash %v Receipt %+v", approveHash, res.Raw)
	showBalance(ownerAddr)
}

func TestTransferInterChain(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	coin := "btp-0x2.icon-sICX"
	srcKey := TokenGodKey
	srcAddr := TokenGodAddr
	dstAddr := "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81201"
	if val, err := api.GetCoinBalance(coin, srcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial Balance %v", val)
	}

	amount := new(big.Int)
	amount.SetString("4300000000000000000", 10)

	txnHash, err := api.Transfer(coin, srcKey, dstAddr, amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash  %v", txnHash)

	if val, err := api.GetCoinBalance(coin, srcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Intermediate Demo Balance %v", val)
	}
	time.Sleep(time.Second * 3)
	res, err := api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range res.ElInfo {
		// seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq", lin)
	}
	if val, err := api.GetCoinBalance(coin, srcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Final Balance %v", val)
	}
}

func TestBatchTransfer(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	coins := []string{"btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD", "btp-0x2.icon-ICX"}

	amount := new(big.Int)
	amount.SetString("4500000000000000000", 10)
	amounts := []*big.Int{amount, amount, amount}
	for i, coin := range coins {

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
	if err := showBalance(TokenGodAddr); err != nil {
		t.Fatalf(" %+v", err)
	}

}

func showBalance(addr string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}
	for _, coinName := range []string{"btp-0x61.bsc-BUSD", "btp-0x61.bsc-USDT", "btp-0x61.bsc-USDC", "btp-0x61.bsc-BTCB", "btp-0x61.bsc-ETH", "btp-0x2.icon-ICX", "btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD", "btp-0x61.bsc-BNB"} {
		res, err := api.GetCoinBalance(coinName, addr)
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res.String())
	}
	return nil
}

func TestReceiver(t *testing.T) {
	recv, err := getNewApi()
	if err != nil {
		panic(err)
	}
	// recv.WatchForTransferStart(1, "ICX", 10)
	// recv.WatchForTransferReceived(1, "TONE", 8)
	// recv.WatchForTransferEnd(1, "ICX", 10)

	go func() {
		if sinkChan, errChan, err := recv.Subscribe(context.Background()); err != nil {
			panic(err)
		} else {
			for {
				select {
				case err := <-errChan:
					panic(err)
				case msg := <-sinkChan:
					t.Logf("\nMessage %+v\n", msg)
				}
			}
		}
	}()
	time.Sleep(time.Second * 3000)
}

func getNewApi() (chain.ChainAPI, error) {
	srcEndpoint := RPC_URI
	addrToName := map[chain.ContractName]string{
		chain.BTS: "cx220b39946c06487027b6fbbedc8eea58899a73a2",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	return icon.NewApi(l, &chain.Config{
		Name:              chain.ICON,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NetworkID:         NID,
		NativeCoin:        "btp-0x2.icon-ICX",
		NativeTokens:      []string{"btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD"},
		WrappedCoins:      []string{"btp-0x61.bsc-BNB", "btp-0x61.bsc-BUSD", "btp-0x61.bsc-USDT", "btp-0x61.bsc-USDC", "btp-0x61.bsc-BTCB", "btp-0x61.bsc-ETH"},
	})
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

func TestGetKeyFromFile(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	priv, pub, err := api.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fmt.Println(priv, "   ", pub)
}

/*
func TestIsOwner(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.CallBTS(chain.IsOwner, []interface{}{BtsOwner})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Res ", res)
}

func TestSetTokenLimit(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	btsOwnerKey, _, err := rpi.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.bts.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.bts.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	coin := "sICX"
	amount := new(big.Int)
	amount.SetString("100000000000000", 10)
	hash, err := rpi.TransactWithBTS(btsOwnerKey, chain.SetTokenLimit, []interface{}{[]string{coin}, []*big.Int{amount}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Hash ", hash)
	res, err := rpi.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Res StatusCode %v Raw %+v", res.StatusCode, res.Raw)
}

func TestGetTokenLimit(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, coin := range []string{"bnUSD", "sICX", "ICX"} {
		res, err := rpi.CallBTS(chain.GetTokenLimit, []interface{}{coin})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Res coin ", res)
	}
}

func TestGetTokenLimitStatus(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for _, coin := range []string{"bnUSD", "sICX", "ICX"} {
		res, err := rpi.CallBTS(chain.GetTokenLimitStatus, []interface{}{"0x61.bsc", coin})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Res coin ", res)
	}
}
func TestCheckTransferRestrictions(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amount := new(big.Int)
	amount.SetString("2000000000000000000000", 10)
	res, err := rpi.CallBTS(chain.CheckTransferRestrictions, []interface{}{
		"0x61.bsc",
		"sICX",
		DemoDstAddr,
		amount,
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fmt.Println("Res ", res)
}

func TestChangeRestrictions(t *testing.T) {
	funcn := chain.AddRestriction
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	btsOwnerKey, _, err := rpi.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.bts.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.bts.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hash, err := rpi.TransactWithBTS(btsOwnerKey, funcn, []interface{}{})
	t.Log("Hash ", hash)
	res, err := rpi.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Res StatusCode %v Raw %+v", res.StatusCode, res.Raw)
}

func TestChangeBlackList(t *testing.T) {
	funcn := chain.AddBlackListAddress
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	btsOwnerKey, _, err := rpi.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.bts.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.bts.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hash, err := rpi.TransactWithBTS(btsOwnerKey, funcn, []interface{}{
		"0x61.bsc",
		[]string{"0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202", "0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"},
	})
	t.Log("Hash ", hash)
	res, err := rpi.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Res StatusCode %v Raw %+v", res.StatusCode, res.Raw)
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

func TestGetBlackListedUsers(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.CallBTS(chain.GetBlackListedUsers, []interface{}{
		"0x61.bsc",
		"0x0",
		"0x5",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}
*/
