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
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestReceiver(t *testing.T) {

	recv, err := getNewApi()
	if err != nil {
		panic(err)
	}
	recv.WatchForTransferStart(1, "ICX", 14)
	recv.WatchForTransferReceived(1, "ONE", 15)
	recv.WatchForTransferEnd(1, "ICX", 14)

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
	time.Sleep(time.Second * 30)
}

func getNewApi() (chain.ChainAPI, error) {
	srcAddress := "btp://0x5b9a77.icon/cx0f011b8b10f2c0d850d5135ef57ea42120452003"
	dstAddress := "btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91"
	srcEndpoint := "http://localhost:9080/api/v3/default"

	addrToName := map[chain.ContractName]string{
		chain.TokenBSHIcon:          "cxd029bba56f72e2ced0c88fd2fc289f4ae4dcd31f",
		chain.NativeBSHIcon:         "cxabbcd08546141646dd169ae70170da87b9296778",
		chain.Irc2Icon:              "cx70053f7c2d0d985c5b342886c5fe8f5e4db1fb1b",
		chain.Irc2TradeableIcon:     "cxd5faca679820dd974245eaceca3fb74536815f96",
		chain.TokenBSHImplHmy:       "0x8283e3bE7ac5f6dB332Df605f20E2B4c9977c662",
		chain.NativeBSHPeripheryHmy: "0xfad748a1063a40FF447B5D766331904d9bedDC26",
		chain.Erc20Hmy:              "0xb54f5e97972AcF96470e02BE0456c8DB2173f33a",
		chain.NativeBSHCoreHmy:      "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95",
		chain.TokenBSHProxyHmy:      "0x48cacC89f023f318B4289A18aBEd44753a127782",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	networkID := "0x5b9a77"
	return icon.NewApi(l, &chain.ChainConfig{Name: chain.ICON, URL: srcEndpoint, ConftractAddresses: addrToName, Src: chain.BTPAddress(srcAddress), Dst: chain.BTPAddress(dstAddress), NetworkID: networkID})
}

func TestGodWalletTransfer(t *testing.T) {
	getKeyPairFromFile := func(walFile string, password string) (pair [2]string, err error) {
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
	godKeyPair, err := getKeyPairFromFile("/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.god.wallet.json", "gochain")
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
	demoKeyPair, err := api.GetKeyPairs(1)
	if err != nil {
		t.Fatal(err)
	}
	amount := new(big.Int)
	amount.SetString("100000000000000000000", 10)
	t.Logf("Demo KeyPair %v", demoKeyPair)
	txnHash, err := api.Transfer("ETH", godKeyPair[0], api.GetBTPAddress(demoKeyPair[0][1]), *amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, elInfo, err := api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range elInfo {
		t.Logf("Log %+v ", lin)
	}
	if val, err := api.GetCoinBalance("ETH", demoKeyPair[0][1]); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Balance %v", val.String())
	}
	//[[13a135a25373970bfd77c21c6f8e5c38f73a8bf55fa0c691f6b4629546e74a23 hxe98400c26c64bcdf6dc26cf1c3d3da5160e76dc3]]
	return
}

func TestTransferAcross(t *testing.T) {
	senderKey := "f5b501cb7527c39ad064313c24fd9e0ee2dc443baa6a1dfa50f97cc7ab88aee0"
	senderAddress := "hx2b66a78bf1ebc8d34133058ea648b243be099267"
	rxAddress := "btp://0x6357d2e0.hmny/0x5ce1c8b80020cE82054d114e0440117470d3611F"
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	if val, err := api.GetCoinBalance("ETH", senderAddress); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Balance %v", val.String())
	}

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)
	txnHash, err := api.Approve("ICX", senderKey, *amount)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	txnHash, err = api.Transfer("ICX", senderKey, rxAddress, *amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, elInfo, err := api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range elInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	}
}

func TestIconEventParse(t *testing.T) {
	btp_icon_token_bsh := "cx5924a147ae30091ed9c6fe0c153ef77de4132902"
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
