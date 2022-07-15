package icon_test

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	GodKey      = "a8779264ee269028cd4d997ff8a49f0875972ccf9faa753ca6edbcdc99060528"
	GodAddr     = "btp://0xdf6463.icon/hxbf007df04a21bdbc1462eb610d596bf711620fda"
	DemoSrcKey  = "f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5"
	DemoSrcAddr = "btp://0xdf6463.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	DemoDstAddr = "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
)

func TestTransferIntraChain(t *testing.T) {
	// godKeyPair, err := getKeyPairFromFile("../../icon.god.wallet.json", "d7b864bc6b02cc30")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	// t.Logf("God KeyPair %v", godKeyPair)

	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}

	amount := new(big.Int)
	amount.SetString("100000000000000000000", 10)

	for _, coinName := range []string{"ICX", "TICX"} {
		txnHash, err := api.Transfer(coinName, GodKey, DemoSrcAddr, *amount)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Transaction Hash %v", txnHash)
		res, err := api.WaitForTxnResult(context.TODO(), txnHash)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Receipt %+v", res)
		for _, lin := range res.ElInfo {
			t.Logf("Log %+v ", lin)
		}

		if val, err := api.GetCoinBalance(coinName, DemoSrcAddr); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Balance %v", val.String())
		}
	}
	return
}

func TestAllowance(t *testing.T) {
	for _, coin := range []string{"TBNB", "ICX", "TICX"} {
		rpi, err := getNewApi()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if allowanceAmt, err := rpi.GetAllowance(coin, DemoSrcAddr); err != nil {
			t.Fatalf("%+v", err)
		} else {
			t.Logf("Allowance %v: %v", coin, allowanceAmt)
		}
	}
}

func TestApprove(t *testing.T) {
	showBalance(DemoSrcAddr)
	coin := "TICX"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("20000000000000000000", 10)
	approveHash, err := rpi.Approve(coin, DemoSrcKey, *amt)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	res, err := rpi.WaitForTxnResult(context.TODO(), approveHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Hash %v Receipt %+v", approveHash, res.Raw)
	showBalance(DemoSrcAddr)
}

func TestTransferInterChain(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	coin := "TICX"
	if val, err := api.GetCoinBalance(coin, DemoSrcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial Balance %v", val.String())
	}

	amount := new(big.Int)
	amount.SetString("20000000000000000000", 10)

	txnHash, err := api.Transfer(coin, DemoSrcKey, DemoDstAddr, *amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash  %v", txnHash)
	res, err := api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range res.ElInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	}
	if val, err := api.GetCoinBalance(coin, DemoSrcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Final Balance %v", val.String())
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
	for _, coinName := range []string{"ICX", "TICX", "BNB", "TBNB"} {
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
	srcEndpoint := "http://localhost:9080/api/v3/icon"
	addrToName := map[chain.ContractName]string{
		chain.BTSIcon:  "cxa1da4ba07a3fcf2ee8027ffba022102ca2f8d321",
		chain.TICXIcon: "cx13f080e39ca30fb111465376953efc3f24690442", //irc2
	}
	l := log.New()
	log.SetGlobalLogger(l)
	networkID := "0xdf6463"
	return icon.NewApi(l, &chain.ChainConfig{Name: chain.ICON, URL: srcEndpoint, ConftractAddresses: addrToName, NetworkID: networkID})
}

/*
func TestIconEventParse(t *testing.T) {
	m := map[chain.ContractName]string{
		chain.TokenBSHIcon: btp_icon_token_bsh,
	}
	parser, err := icon.NewParser(m)
	if err != nil {
		t.Fatal(err)
	}
	log := &icon.TxnEventLog{
		Addr:    icon.Address("cx5924a147ae30091ed9c6fe0c153ef77de4132902"),
		Indexed: []string{"TransferStart(Address,str,int,bytes)", "hx4a707b2ecbb5f40a8d761976d99244f53575eeb6"},
		Data:    []string{"btp://0x6357d2e0.hmny/0x8BE8641225CC0Afdb24499409863E8E3f6557C32", "0x25", "0xd6d583455448880dbd2fc137a30000872386f26fc10000"},
	}
	res, eventType, err := parser.Parse(log)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("EventType %v  Res %+v", eventType, res)
}
*/
