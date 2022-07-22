package hmny_test

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	GodKey      = "1f84c95ac16e6a50f08d44c7bde7aff8742212fda6e4321fde48bf83bef266dc"
	GodAddr     = "btp://0x6357d2e0.hmny/0xA5241513DA9F4463F1d4874b548dFBAC29D91f34"
	DemoSrcKey  = "564971a566ce839535681eef81ccd44005944b98f7409cb5c0f5684ae862a530"
	DemoSrcAddr = "btp://0x6357d2e0.hmny/0x8fc668275b4fa032342ea3039653d841f069a83b"
	DemoDstAddr = "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee1"
	GodDstAddr  = "btp://0x5b9a77.icon/hxff0ea998b84ab9955157ab27915a9dc1805edd35"
	//BtsAddr     = "btp://0x6357d2e0.hmny/0x05AcF27495FAAf9A178e316B9Da2f330983b9B95"
)

func getNewApi() (chain.ChainAPI, error) {
	//ICONDemo [f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5 btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee]
	//HmnyDemo [564971a566ce839535681eef81ccd44005944b98f7409cb5c0f5684ae862a530 btp://0x6357d2e0.hmny/0x8fc668275b4fa032342ea3039653d841f069a83b]

	coinMap := map[string]string{
		"TONE": "0xB20CCD2a42e5486054AE3439f2bDa95DC75d9B75",
	}
	l := log.New()
	log.SetGlobalLogger(l)

	addrToName := map[chain.ContractName]string{
		chain.BTSCoreHmny:      "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95",
		chain.BTSPeripheryHmny: "0xfad748a1063a40FF447B5D766331904d9bedDC26",
	}
	rx, err := hmny.NewApi(l, &chain.ChainConfig{
		Name:                 chain.HMNY,
		URL:                  "http://localhost:9500",
		ContractAddresses:    addrToName,
		NetworkID:            "0x6357d2e0",
		NativeCoin:           "ONE",
		NativeTokenAddresses: coinMap,
	})
	if err != nil {
		return nil, err
	}
	return rx, nil
}

func TestApprove(t *testing.T) {
	showBalance(DemoSrcAddr)
	coin := "TONE"
	rpi, err := getNewApi()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	amt := new(big.Int)
	amt.SetString("27000000000000000000", 10)
	approveHash, err := rpi.Approve(coin, DemoSrcKey, amt)
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

func TestGetCoinBalance(t *testing.T) {
	// demoKeyPair, err := getKeyPairFromFile("../../../../devnet/docker/icon-hmny/src/hmny.god.wallet.json", "")
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
	// fmt.Println(demoKeyPair)
	// return
	if err := showBalance(DemoSrcAddr); err != nil {
		t.Fatalf("%+v ", err)
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

func showBalance(addr string) error {
	api, err := getNewApi()
	if err != nil {
		return err
	}
	for _, coinName := range []string{"ICX", "ONE", "TICX", "TONE"} {
		res, err := api.GetCoinBalance(coinName, addr)
		if err != nil {
			return err
		}
		log.Infof("coin %v amount %v", coinName, res)
	}
	return nil
}

func TestTransferIntraChain(t *testing.T) {

	// godKeyPair, err := getKeyPairFromFile("../../../../devnet/docker/icon-hmny/src/hmny.god.wallet.json", "")
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
	amount.SetString("90000000000000000000", 10)
	for _, coin := range []string{"TONE", "ONE"} {
		txnHash, err := api.Transfer(coin, GodKey, DemoSrcAddr, amount)
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
			seq, _ := lin.GetSeq()
			t.Logf("Log %+v and Seq %v", lin, seq)
		}
		if val, err := api.GetCoinBalance(coin, DemoSrcAddr); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Balance %v", val)
		}
	}
	return
}

//ICONDemo [f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5 btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee]
//HmnyDemo [564971a566ce839535681eef81ccd44005944b98f7409cb5c0f5684ae862a530 btp://0x6357d2e0.hmny/0x8Fc668275b4fA032342eA3039653D841f069a83b]

func TestTransferInterChain(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	if val, err := api.GetCoinBalance("TONE", DemoSrcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial  Balance %v", val)
	}
	amount := new(big.Int)
	amount.SetString("9000000000000000000", 10)

	txnHash, err := api.Transfer("TONE", DemoSrcKey, DemoDstAddr, amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v   %v   %v", txnHash, DemoSrcAddr, DemoDstAddr)
	res, err := api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range res.ElInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	}
	if val, err := api.GetCoinBalance("TONE", DemoSrcAddr); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Final Balance %v", val)
	}

}
func TestReceiver(t *testing.T) {
	recv, err := getNewApi()
	if err != nil {
		t.Fatal(err)
	}

	sinkChan, errChan, err := recv.Subscribe(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-errChan:
			t.Fatalf("%+v", err)

		case msg := <-sinkChan:
			t.Logf("\nMessage %+v\n", msg)
		}
	}
}

/*
func TestBshEvent(t *testing.T) {

	btp_hmny_token_bsh_impl := "0xfAC8B63F77d8056A9BB45175b3DEd7D316D868D4"
	btp_hmny_nativecoin_bsh_periphery := "0xfEe5c5B2bc2f927335C60879d78304e4305CdBaC"
	contractUsed := btp_hmny_nativecoin_bsh_periphery
	m := map[chain.ContractName]string{
		chain.TokenBSHImplHmy:       btp_hmny_token_bsh_impl,
		chain.NativeBSHPeripheryHmy: btp_hmny_nativecoin_bsh_periphery,
	}
	p, err := NewParser("http://127.0.0.1:9500", m)
	if err != nil {
		t.Fatal(err)
	}
	const transferStartHex = "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000001d00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000003e6274703a2f2f30783562396137372e69636f6e2f68786663303136306133306565373033393861303134393438303531623331383963326563373261306200000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000dbd2fc137a30000000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000000034554480000000000000000000000000000000000000000000000000000000000"
	txStartBytes, err := hex.DecodeString(transferStartHex)
	if err != nil {
		t.Fatal(err)
	}

	log := &types.Log{
		Address:     common.HexToAddress(contractUsed),
		Topics:      []common.Hash{common.HexToHash("0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a"), common.HexToHash("0x00000000000000000000000052c08a9a3a457e9ec8db545793ab9f0630dec4b4")},
		Data:        txStartBytes,
		BlockNumber: 100,
		TxHash:      common.HexToHash("0x123"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x456"),
		Index:       2,
	}

	res, eventType, err := p.Parse(log)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("EventType %v  Res %+v", eventType, res)
}
*/
