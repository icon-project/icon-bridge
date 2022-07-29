package icon_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	TokenGodKey  = "1613b2f469bc7e33673846ee538aae4e7377f853deab7e5e97bf60503cc8ac42"
	TokenGodAddr = "btp://0x2.icon/hx6bbea9df10c329c0c6b428788070ca6b4d589b80"
	NID          = "0x2.icon"
	RPC_URI      = "https://lisbon.net.solidwallet.io/api/v3/icon_dex"
	GodKey       = "596294cb42cddb975b9730113097c2cf370a5819d7a58baaaa4a3d0e9af065fa"
	GodAddr      = "btp://0x2.icon/hxc86452374f94bd8db99f703bb1fc3fad2f7b2024"
	GodDstAddr   = "btp://0x61.bsc/0xDf9e6205Ac201c8a11082842857C6f7673a8246e"
	BtsAddr      = "btp://0x2.icon/cx69774ba6f0d2718bef41065227345529a11b57f1"

	DemoDstAddr = "btp://0x61.bsc/0xcf5BC0BD5aEdf6cd216f7288c2Fd704a397F453d"
	DemoSrcKey  = "41528d2ae0a203914f39584c6b2ace17b61c1208492be5952177ed8b16b1b99f"
	DemoSrcAddr = "btp://0x2.icon/hx6bbea9df10c329c0c6b428788070ca6b4d589b80"
)

// const (
// 	RPC_URI     = "http://localhost:9080/api/v3/icon"
// 	GodKey      = "c4a15fbef04e99892caaa11374b115795c182d290d5d8bd7821a9ef16f4a9bcf"
// 	GodAddr     = "btp://0x613f17.icon/hxad8eec2e167c24020600ddf1acd4d03673d3f49b"
// 	DemoSrcKey  = "f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5"
// 	DemoSrcAddr = "btp://0x613f17.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
// 	DemoDstAddr = "btp://0x61.bsc/0x0000000000000000000000000000000000000000"
// 	GodDstAddr  = "btp://0x61.bsc/0x20E789D2f5D469eA30e0525DbfDD5515d6EAd30D"
// 	NID         = "0x613f17.icon"
// )

func TestTransferIntraChain(t *testing.T) {
	// godKeyPair, err := getKeyPairFromFile("//home/manish/go/src/work/icon-bridge/lisbon/wallets/icon.bmr.wallet.json", "1234")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	// t.Log(godKeyPair)
	// return

	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}

	amount := new(big.Int)
	amount.SetString("100000000000000", 10)
	srckey := GodKey
	dstaddr := BtsAddr
	for _, coinName := range []string{"DUM"} {
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
	ownerKey := GodKey
	ownerAddr := GodAddr
	showBalance(ownerAddr)
	coin := "sICX"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("100000000000000", 10)
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
	coin := "sICX"
	srcKey := GodKey
	srcAddr := GodAddr
	dstAddr := GodDstAddr
	if val, err := api.GetCoinBalance(coin, srcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial Balance %v", val)
	}

	amount := new(big.Int)
	amount.SetString("100000000000000", 10)

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
	coins := []string{"TICX", "ETH", "ICX"}
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
	if err := showBalance(GodAddr); err != nil {
		t.Fatalf(" %+v", err)
	}
}

func showBalance(addr string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}

	for _, coinName := range []string{"ICX", "sICX", "bnUSD", "DUM", "BNB", "BUSD", "USDT", "USDC", "BTCB", "ETH"} {
		res, err := api.GetCoinBalance(coinName, addr)
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res.String())
	}
	return nil
}

func getKeyPairFromFile(walFile string, password string) (pair [2]string, err error) {
	keyReader, err := os.Open(walFile)
	if err != nil {
		return
	}
	defer keyReader.Close()

	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		return
	}
	key, err := keystore.DecryptKey(keyStore, password)
	if err != nil {
		return
	}
	privBytes := ethcrypto.FromECDSA(key.PrivateKey)
	privString := hex.EncodeToString(privBytes)
	addr := ethcrypto.PubkeyToAddress(key.PrivateKey.PublicKey)
	pair = [2]string{privString, addr.String()}
	return
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
		chain.BTS: "cx69774ba6f0d2718bef41065227345529a11b57f1",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	return icon.NewApi(l, &chain.Config{
		Name:              chain.ICON,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NetworkID:         NID,
		NativeCoin:        "ICX",
		NativeTokens:      []string{"sICX", "bnUSD", "DUM"},
		WrappedCoins:      []string{"BNB", "BUSD", "USDT", "USDC", "BTCB", "ETH"},
		GasLimit:          8000000,
	})
}

func TestConverToZeroAddress(t *testing.T) {
	addr := DemoSrcAddr
	splits := strings.Split(addr, "/")
	if len(splits) != 4 {
		return
	}
	network := splits[2]
	networkSplits := strings.Split(network, ".")
	if len(networkSplits) != 2 {
		return
	}
	networkSplits[1] += "s"
	splits[2] = strings.Join(networkSplits, ".")
	joinStr := strings.Join(splits, "/")
	fmt.Println(joinStr)
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
