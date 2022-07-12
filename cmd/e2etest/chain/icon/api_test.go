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
	//ICONDemo [f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5 btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee]
	//HmnyDemo [564971a566ce839535681eef81ccd44005944b98f7409cb5c0f5684ae862a530 btp://0x6357d2e0.hmny/0x8Fc668275b4fA032342eA3039653D841f069a83b]

	// srcAddress := "btp://0x5b9a77.icon/cx7db813639e4b3be5f66a05addbbbea7958ba5247"
	// dstAddress := "btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91"
	srcEndpoint := "http://localhost:9080/api/v3/icon"

	addrToName := map[chain.ContractName]string{
		chain.BTSIcon:  "cxf9a4556e7049bf81bf4fb3ffb4f5c23691d3aef6",
		chain.TICXIcon: "cxc39fce2d84ad7a49c07f967f08341900023f1566", //irc2
	}

	l := log.New()
	log.SetGlobalLogger(l)
	networkID := "0xdf6463"
	return icon.NewApi(l, &chain.ChainConfig{Name: chain.ICON, URL: srcEndpoint, ConftractAddresses: addrToName, NetworkID: networkID})
}

func TestTransferIntraChain(t *testing.T) {
	godKeyPair, err := getKeyPairFromFile("../../icon.god.wallet.json", "d7b864bc6b02cc30")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("God KeyPair %v", godKeyPair)

	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	// demoKeyPair, err := api.GetKeyPairs(1)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	demoKeyPair := [][2]string{{"f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5", "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"}}

	amount := new(big.Int)
	amount.SetString("100000000000000000", 10)
	t.Logf("Demo KeyPair %v", demoKeyPair)

	for _, coinName := range []string{"ICX", "TICX"} {
		txnHash, err := api.Transfer(coinName, godKeyPair[0], api.GetBTPAddress(demoKeyPair[0][1]), *amount)
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

		if val, err := api.GetCoinBalance(coinName, demoKeyPair[0][1]); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Balance %v", val.String())
		}
	}
	return
}

func TestGetCoinBalance(t *testing.T) {
	demoKeyPair := [][2]string{{"f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5", "btp://0xdf6463.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"}}
	if err := showBalance(demoKeyPair[0]); err != nil {
		t.Fatalf(" %+v", err)
	}

}

func TestTransferInterChain(t *testing.T) {
	senderKey := "f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5"
	senderAddress := "btp://0xdf6463.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	rxAddress := "btp://0x61.bsc/0x54a1be6CB9260A52B7E2e988Bc143e4c66b81202"
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	if val, err := api.GetCoinBalance("TICX", senderAddress); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial Balance %v", val.String())
	}

	amount := new(big.Int)
	amount.SetString("10000000000000000", 10)
	_, err = api.Approve("TICX", senderKey, *amount)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	txnHash, err := api.Transfer("TICX", senderKey, rxAddress, *amount)
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
	time.Sleep(5 * time.Second)
	if val, err := api.GetCoinBalance("TICX", senderAddress); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Final Balance %v", val.String())
	}
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

func showBalance(demoKeyPair [2]string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}
	for _, coinName := range []string{"ICX", "TICX", "BNB", "TBNB"} {
		res, err := api.GetCoinBalance(coinName, demoKeyPair[1])
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res.String())
	}
	return nil
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
