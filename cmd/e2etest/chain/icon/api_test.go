package icon

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
	RPC_URI      = "https://lisbon.net.solidwallet.io/api/v3/icon_dex"
	TokenGodKey  = ""
	TokenGodAddr = "btp://0x2.icon/"
	GodKey       = ""
	GodAddr      = "btp://0x2.icon/hx077ada6dd02f63b02650c5861f9f41166e45d9f1"
	DemoSrcKey   = ""
	DemoSrcAddr  = "btp://0x2.icon/hx96f2c7524c0557f8b56d461205443367cb731e83"
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
	amount.SetString("70750343697585467", 10)
	srckey := TokenGodKey
	dstaddr := DemoSrcAddr
	for _, coinName := range []string{"btp-0x228.snow-ICZ"} {
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
	coin := "btp-0x2.icon-ICX"
	srcKey := TokenGodKey
	srcAddr := TokenGodAddr
	dstAddr := "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81201"
	if val, err := api.GetCoinBalance(coin, srcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial Balance %v", val)
	}

	amount := new(big.Int)
	amount.SetString("8750741761044791348", 10)

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
	coins := []string{"btp-0x2.icon-sICX"}

	amount := new(big.Int)
	amount.SetString("4500000000000000000", 10)
	amounts := []*big.Int{amount}
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
	if err := showBalance(DemoSrcAddr); err != nil {
		t.Fatalf(" %+v", err)
	}

}

func showBalance(addr string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}
	for _, coinName := range []string{"btp-0x61.bsc-BUSD", "btp-0x61.bsc-USDT", "btp-0x61.bsc-USDC", "btp-0x61.bsc-BTCB", "btp-0x61.bsc-ETH", "btp-0x2.icon-ICX", "btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD", "btp-0x61.bsc-BNB", "btp-0x228.snow-ICZ"} {
		res, err := api.GetCoinBalance(coinName, addr)
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res.String())
	}
	return nil
}

func TestFeeGatheringTerm(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
	}
	bmcOwnerKey, _, err := api.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hash, err := api.SetFeeGatheringTerm(bmcOwnerKey, 1200)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Log("Hash ", hash)
	res, err := api.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 1 {
		t.Fatalf("StatusCode not 1 for Batch Transfer \n %+v", res.Raw)
	}
	interval, err := api.GetFeeGatheringTerm()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Get ", interval)
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
		chain.BMC: "cx8059df76efcd0b076c7493756b5baf6a5bfe03c4",
		chain.BTS: "cx9b16a2374e6fd35f223bb902137ce34013c8a5f2",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	return NewApi(l, &chain.Config{
		Name:              chain.ICON,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NetworkID:         NID,
		NativeCoin:        "btp-0x2.icon-ICX",
		NativeTokens:      []string{"btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD"},
		WrappedCoins:      []string{"btp-0x61.bsc-BNB", "btp-0x61.bsc-BUSD", "btp-0x61.bsc-USDT", "btp-0x61.bsc-USDC", "btp-0x61.bsc-BTCB", "btp-0x61.bsc-ETH", "btp-0x228.snow-ICZ", "btp-0x2.near-NEAR"},
		GasLimit:          make(map[chain.GasLimitType]uint64),
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
	priv, pub, err := api.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh_arxiv/_ixh_icon-bsc/keystore/icon.god.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh_arxiv/_ixh_icon-bsc/keystore/icon.god.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fmt.Println(priv, "   ", pub)
}

func TestAddToBlackList(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	btsOwnerKey, _, err := rpi.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	hash, err := rpi.AddBlackListAddress(btsOwnerKey,
		"0x61.bsc",
		[]string{"0xb4c1be63C9260A52B7C3eaa8Bc143ff4c66b81206"},
	)
	t.Log("Hash ", hash)
	res, err := rpi.WaitForTxnResult(context.TODO(), hash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Res StatusCode %v Raw %v", res.StatusCode, res.Raw)
}

func TestGetBlackListedUsers(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.GetBlackListedUsers("0x61.bsc", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}

func TestIsUserBlackListed(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.IsUserBlackListed("0x61.bsc", "0x94ACCD1f12cF6FF25Aaeb483605536918D7760b5")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Res ", res)
}

func TestSetTokenLimit(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	btsOwnerKey, _, err := rpi.GetKeyPairFromKeystore("../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.json", "../../../../devnet/docker/icon-bsc/_ixh/keystore/icon.god.wallet.secret")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	coin := "btp-0x2.icon-bnUSD"
	amount := new(big.Int)
	amount.SetString("115792089237316195423570985", 10)
	hash, err := rpi.SetTokenLimit(btsOwnerKey, []string{coin}, []*big.Int{amount})
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

func TestBMCParseMessage(t *testing.T) {
	rcpt := &TxnEventLog{
		Addr:    "cx3fbcd7c25be9aac63bd944a0b8d40becdaa35425",
		Indexed: []string{"Message(str,int,bytes)", "btp://0x61.bsc/0x855183AefDA98Cd10Cf31488e26fF599685625a7", "0x1e"},
		Data:    []string{"0xf893b8396274703a2f2f307836312e6273632f307838353531383341656644413938436431304366333134383865323666463539393638353632356137b8396274703a2f2f3078322e69636f6e2f6378336662636437633235626539616163363362643934346130623864343062656364616133353432358362747381c496d52893496e76616c69642075696e74206e756d626572"},
	}

	rcpt.Data = []string{"0xf8d0b8396274703a2f2f3078322e69636f6e2f637833666263643763323562653961616336336264393434613062386434306265636461613335343235b8396274703a2f2f307836312e6273632f30783835353138334165664441393843643130436633313438386532366646353939363835363235613783626d6300b853f8518c466565476174686572696e67b842f840b8396274703a2f2f3078322e69636f6e2f687832373531333661633936396161646135656466623534333939313264623564636431303734626466c483627473"}

	addrToName := map[chain.ContractName]string{
		chain.BTS: "cx9d4f2ecaa38d3c34d508bd9c8df73c52dbef60dd",
		chain.BMC: "cx3fbcd7c25be9aac63bd944a0b8d40becdaa35425",
	}
	parser, err := NewParser(addrToName)
	if err != nil {
		t.Fatal(err)
	}

	res, _, err := parser.ParseTxn(rcpt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
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
}

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
*/
