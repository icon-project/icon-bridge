package hmny

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

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
	btp_hmny_token_bsh_impl := "0x8283e3bE7ac5f6dB332Df605f20E2B4c9977c662"
	btp_hmny_nativecoin_bsh_periphery := "0xfad748a1063a40FF447B5D766331904d9bedDC26"
	addrToName := map[chain.ContractName]string{
		chain.TokenBSHImplHmy:       btp_hmny_token_bsh_impl,
		chain.NativeBSHPeripheryHmy: btp_hmny_nativecoin_bsh_periphery,
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

func TestTransferAcross(t *testing.T) {
	senderKey := "05bfc351f5ec0e81d88ccda3df3108d993415d92c1d41e09afa6ea24ba8a4307"
	senderAddress := "0x80d1f81A5E541cA370308571AAbD096cCA6C901c"
	rxAddress := "btp://0x5b9a77.icon/hx3fdf6ff1c0e747f7573b365d2890c84bed107162"
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	amount := new(big.Int)
	amount.SetString("2000000000000000000", 10)
	txnHash, err := api.Transfer(&chain.RequestParam{
		FromChain:   chain.HMNY,
		ToChain:     chain.ICON,
		SenderKey:   senderKey,
		FromAddress: senderAddress,
		ToAddress:   rxAddress,
		Amount:      *amount,
		Token:       chain.ONEToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, elInfo, err := api.WaitForTxnResult(txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range elInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	}
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
	amount.SetString("250000000000000000000", 10)
	t.Logf("Demo KeyPair %v", demoKeyPair)
	txnHash, err := api.Transfer(&chain.RequestParam{
		FromChain:   chain.HMNY,
		ToChain:     chain.HMNY,
		SenderKey:   godKeyPair[0],
		FromAddress: godKeyPair[1],
		ToAddress:   demoKeyPair[0][1],
		Amount:      *amount,
		Token:       chain.ONEToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Transaction Hash %v", txnHash)
	res, elInfo, err := api.WaitForTxnResult(txnHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Receipt %+v", res)
	for _, lin := range elInfo {
		seq, _ := lin.GetSeq()
		t.Logf("Log %+v and Seq %v", lin, seq)
	} //[[05bfc351f5ec0e81d88ccda3df3108d993415d92c1d41e09afa6ea24ba8a4307 0x80d1f81A5E541cA370308571AAbD096cCA6C901c]]
	return
}

func TestReceiver(t *testing.T) {
	rx, err := getNewApi()
	if err != nil {
		t.Fatal(err)
	}
	startHeight := 15000
	rx.WatchFor(chain.TransferStart, 7, "0xfad748a1063a40FF447B5D766331904d9bedDC26")
	rx.WatchFor(chain.TransferEnd, 7, "0xfad748a1063a40FF447B5D766331904d9bedDC26")
	rx.WatchFor(chain.TransferReceived, 5, "0xfad748a1063a40FF447B5D766331904d9bedDC26")
	sinkChan, errChan, err := rx.Subscribe(context.TODO(), uint64(startHeight))
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
