package bsc

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

const (
	RPC_URI      = "https://data-seed-prebsc-1-s1.binance.org:8545"
	TokenGodKey  = ""
	TokenGodAddr = "btp://0x61.bsc/59d1d3450c1275ebf4ca477bf49fbcf910676e62"
	GodKey       = ""
	GodAddr      = "btp://0x61.bsc/0x70E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
	DemoSrcKey   = ""
	DemoSrcAddr  = "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	DemoDstAddr  = "btp://0x2.icon/hx0000000000000000000000000000000000000000"
	GodDstAddr   = "btp://0x2.icon/hx077ada6dd02f63b02650c5861f9f41166e45d9f1"
	BtsAddr      = "btp://0x61.bsc/0x71a1520bBb7e6072Bbf3682A60c73D63b693690A"
)

func TestApprove(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("5000000000000000000", 10)
	for _, coin := range []string{"btp-0x61.bsc-ETH"} {
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
	dstAddr := "btp://0x61.bsc/0x59d1d3450c1275ebf4ca477bf49fbcf910676e62"
	amt := new(big.Int)
	amt.SetString("500000000000000000", 10)
	for _, coin := range []string{"btp-0x61.bsc-BNB"} {
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
	coin := "btp-0x61.bsc-ETH"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	srcKey := TokenGodKey
	srcAddr := TokenGodAddr
	dstAddr := GodDstAddr
	for _, coin := range []string{coin} {
		res, err := rpi.GetCoinBalance(coin, srcAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res)
	}

	amt := new(big.Int)
	amt.SetString("5000000000000000000", 10)
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

func TestBatchTransfer(t *testing.T) {
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	coins := []string{"btp-0x61.bsc-USDC"}

	largeAmt := new(big.Int)
	largeAmt.SetString("5000000000000000000", 10)
	amounts := []*big.Int{largeAmt}
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

	for _, coin := range []string{"btp-0x61.bsc-BUSD", "btp-0x61.bsc-USDT", "btp-0x61.bsc-USDC", "btp-0x61.bsc-BTCB", "btp-0x61.bsc-ETH", "btp-0x2.icon-ICX", "btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD", "btp-0x61.bsc-BNB"} {
		res, err := rpi.GetCoinBalance(coin, TokenGodAddr)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		t.Logf("%v %v", coin, res)
	}
}

func TestReceiver(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
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
		chain.BTS:          "0xb319825645c90ad8De3881c1ec45fC76e46C3080",
		chain.BTSPeriphery: "0x0f0CaF4BE74A1fBB09E5FBc8E9BbBbb7b279598D",
		chain.BMCPeriphery: "0x04576539eB7C9f811a4f436ca84Fe2B47D60E4C8",
	}

	l := log.New()
	log.SetGlobalLogger(l)
	rx, err := NewApi(l, &chain.Config{
		Name:              chain.BSC,
		URL:               RPC_URI,
		ContractAddresses: ctrMap,
		NativeTokens:      []string{"btp-0x61.bsc-BUSD", "btp-0x61.bsc-USDT", "btp-0x61.bsc-USDC", "btp-0x61.bsc-BTCB", "btp-0x61.bsc-ETH"},
		WrappedCoins:      []string{"btp-0x2.icon-ICX", "btp-0x2.icon-sICX", "btp-0x2.icon-bnUSD"},
		NativeCoin:        "btp-0x61.bsc-BNB",
		NetworkID:         "0x61.bsc",
	})
	if err != nil {
		return nil, err
	}
	return rx, nil
}

func TestBMCParseMessage(t *testing.T) {
	rcptJSON := "7b22726f6f74223a223078222c22737461747573223a22307831222c2263756d756c617469766547617355736564223a223078343431303437222c226c6f6773426c6f6f6d223a2230783030303030303032303030303030303030303030303034303030303030303030303430303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030323030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303034303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030323030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030222c226c6f6773223a5b7b2261646472657373223a22307838353531383361656664613938636431306366333134383865323666663539393638353632356137222c22746f70696373223a5b22307833376265333533663231366366376533333633393130316664363130633534326536613063303130393137336661316331643862303464333465646237633162225d2c2264617461223a223078303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303036303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030323630303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030306330303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303033393632373437303361326632663330373833323265363936333666366532663633373833333636363236333634333736333332333536323635333936313631363333363333363236343339333433343631333036323338363433343330363236353633363436313631333333353334333233353030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303039346638393262383339363237343730336132663266333037383336333132653632373336333266333037383338333533353331333833333431363536363434343133393338343336343331333034333636333333313334333833383635333233363636343633353339333933363338333533363332333536313337623833393632373437303361326632663330373833323265363936333666366532663633373833333636363236333634333736333332333536323635333936313631363333363333363236343339333433343631333036323338363433343330363236353633363436313631333333353334333233353833363237343733313039366435303339336432303039303431363436343635363435343666343236633631363336623663363937333734303030303030303030303030303030303030303030303030222c22626c6f636b4e756d626572223a22307831356137323762222c227472616e73616374696f6e48617368223a22307839666339393334366561353735646564643939646239393836346663303634643439393336373735336631636632613431323561343564623431316266653463222c227472616e73616374696f6e496e646578223a2230783139222c22626c6f636b48617368223a22307838343432336339323265613764623632636364363935393164643738306537333663313932356631323732366663633230626162656439396539363939326466222c226c6f67496e646578223a2230783163222c2272656d6f766564223a66616c73657d5d2c227472616e73616374696f6e48617368223a22307839666339393334366561353735646564643939646239393836346663303634643439393336373735336631636632613431323561343564623431316266653463222c22636f6e747261637441646472657373223a22307830303030303030303030303030303030303030303030303030303030303030303030303030303030222c2267617355736564223a2230783735353237222c22626c6f636b48617368223a22307838343432336339323265613764623632636364363935393164643738306537333663313932356631323732366663633230626162656439396539363939326466222c22626c6f636b4e756d626572223a22307831356137323762222c227472616e73616374696f6e496e646578223a2230783139227d"
	rcpt := &ethTypes.Receipt{}
	rcptBytes, err := hex.DecodeString(rcptJSON)
	require.NoError(t, err)
	err = rcpt.UnmarshalJSON(rcptBytes)
	require.NoError(t, err)
	ctrMap := map[chain.ContractName]string{
		chain.BTS:          "0x049a7Ef08e4f4D54Db764cA141f585b79C0310c7",
		chain.BTSPeriphery: "0x5003101271147Fcb6A4c79EE2807C2571e8c6dC7",
		chain.BMCPeriphery: "0x855183AefDA98Cd10Cf31488e26fF599685625a7",
	}
	parser, err := NewParser(RPC_URI, ctrMap)
	if err != nil {
		t.Fatal(err)
	}
	res, _, err := parser.parseMessage(rcpt.Logs[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
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
