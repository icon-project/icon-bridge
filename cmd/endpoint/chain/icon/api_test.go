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
	recv.WatchFor(chain.TransferStart, 5, "cxe4b60a773c63961aa2303961483c3c95b9de3360")
	recv.WatchFor(chain.TransferEnd, 5, "cxe4b60a773c63961aa2303961483c3c95b9de3360")
	recv.WatchFor(chain.TransferReceived, 7, "cxe4b60a773c63961aa2303961483c3c95b9de3360")

	startHeight := 15000
	go func() {
		if sinkChan, errChan, err := recv.Subscribe(context.Background(), uint64(startHeight)); err != nil {
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
	srcEndpoint := []string{"http://localhost:9080/api/v3/default"}

	btp_icon_token_bsh := "cx3e836c763af780392a00a9ac2fc6e0471c95cb50"
	btp_icon_nativecoin_bsh := "cxe4b60a773c63961aa2303961483c3c95b9de3360"
	addrToName := map[chain.ContractName]string{
		chain.TokenBSHIcon:      btp_icon_token_bsh,
		chain.NativeBSHIcon:     btp_icon_nativecoin_bsh,
		chain.Irc2Icon:          "cx10129552153ad5899eb841baf03be5105801bd9a",
		chain.Irc2TradeableIcon: "cx1017c1beb68b5d5c8706c530c164bba91970a6eb",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	networkID := "0x5b9a77"
	api, err := icon.NewApi(chain.BTPAddress(srcAddress), chain.BTPAddress(dstAddress), srcEndpoint, l, addrToName, networkID)
	if err != nil {
		return nil, err
	}
	return api, nil
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
	amount.SetString("250000000000000000000", 10)
	t.Logf("Demo KeyPair %v", demoKeyPair)
	txnHash, err := api.Transfer(&chain.RequestParam{
		FromChain:   chain.ICON,
		ToChain:     chain.ICON,
		SenderKey:   godKeyPair[0],
		FromAddress: godKeyPair[1],
		ToAddress:   demoKeyPair[0][1],
		Amount:      *amount,
		Token:       chain.ICXToken,
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
		t.Logf("Log %+v ", lin)
	} //[[2f3467b8eb43c733fca4387ca8669d569527c26afc3f4d0b18a4cccf44bd44d3 hx3fdf6ff1c0e747f7573b365d2890c84bed107162]]
	return
}

func TestTransferAcross(t *testing.T) {
	senderKey := "2f3467b8eb43c733fca4387ca8669d569527c26afc3f4d0b18a4cccf44bd44d3"
	senderAddress := "hx3fdf6ff1c0e747f7573b365d2890c84bed107162"
	rxAddress := "btp://0x6357d2e0.hmny/0x80d1f81A5E541cA370308571AAbD096cCA6C901c"
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	amount := new(big.Int)
	amount.SetString("2000000000000000000", 10)
	txnHash, err := api.Transfer(&chain.RequestParam{
		FromChain:   chain.ICON,
		ToChain:     chain.HMNY,
		SenderKey:   senderKey,
		FromAddress: senderAddress,
		ToAddress:   rxAddress,
		Amount:      *amount,
		Token:       chain.ICXToken,
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
