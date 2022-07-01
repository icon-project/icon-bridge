package hmny

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func getNewApi() (chain.ChainAPI, error) {
	const (
		src = "btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91"
		dst = "btp://0x5b9a77.icon/cx0f011b8b10f2c0d850d5135ef57ea42120452003"
		url = "http://localhost:9500"
	)

	l := log.New()
	log.SetGlobalLogger(l)

	addrToName := map[chain.ContractName]string{
		chain.TokenBSHImplHmy:       "0x8283e3bE7ac5f6dB332Df605f20E2B4c9977c662",
		chain.NativeBSHPeripheryHmy: "0xfad748a1063a40FF447B5D766331904d9bedDC26",
		chain.Erc20Hmy:              "0xb54f5e97972AcF96470e02BE0456c8DB2173f33a",
		chain.NativeBSHCoreHmy:      "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95",
		chain.TokenBSHProxyHmy:      "0x48cacC89f023f318B4289A18aBEd44753a127782",
	}
	rx, err := NewApi(l, &chain.ChainConfig{Name: chain.HMNY, URL: url, Src: chain.BTPAddress(src), Dst: chain.BTPAddress(dst), ConftractAddresses: addrToName, NetworkID: "0x6357d2e0"})
	if err != nil {
		log.Fatal((err))
	}
	return rx, nil
}

var getKeyPairFromFile = func(walFile string, password string) (pair [2]string, err error) {
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

func TestGodWalletTransfer(t *testing.T) {

	godKeyPair, err := getKeyPairFromFile("/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.god.wallet.json", "")
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
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	} //7957f548a6b2ea3fff4d698a79dbeecaf2f911fce06a5b8a3f4a8c476be4f0e2 0x94efF7b6e91C08195395AA7F5aDF295A62d50Dc3
	if val, err := api.GetCoinBalance("ETH", demoKeyPair[0][1]); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Balance %v", val.String())
	}
	return
}

func TestTransferAcross(t *testing.T) {
	senderKey := "9915b304cfcd7bd2b8ae0232a1d1ca432d5237c167a63a8530d69594c755519e"
	senderAddress := "0x5ce1c8b80020cE82054d114e0440117470d3611F"
	rxAddress := "btp://0x5b9a77.icon/hx2b66a78bf1ebc8d34133058ea648b243be099267"
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
	txnHash, err := api.Approve("ETH", senderKey, *amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, elInfo, err := api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	txnHash, err = api.Transfer("ETH", senderKey, rxAddress, *amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, elInfo, err = api.WaitForTxnResult(context.TODO(), txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range elInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	}

}
func TestReceiver(t *testing.T) {
	recv, err := getNewApi()
	if err != nil {
		t.Fatal(err)
	}
	startHeight := 18000
	recv.WatchForTransferStart(1, "ONE", 15)
	recv.WatchForTransferReceived(1, "ICX", 14)
	recv.WatchForTransferEnd(1, "ONE", 15)
	sinkChan, errChan, err := recv.Subscribe(context.TODO(), uint64(startHeight))
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
