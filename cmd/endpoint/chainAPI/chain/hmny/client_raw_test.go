package hmny

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/icon-project/icon-bridge/common/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/rpc"
	bshcore "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/bsh/bshcore"
	erc20 "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/bsh/erc20tradable"
	bep20tkn "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/tokenbsh/bep20tkn"

	"github.com/icon-project/icon-bridge/common/wallet"
)

func TestGet_hmny_erc20_balance(t *testing.T) {
	url := "http://127.0.0.1:9500"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	cleth := ethclient.NewClient(clrpc)

	// Smart contract addresses
	btp_hmny_erc20 := "0xb54f5e97972AcF96470e02BE0456c8DB2173f33a" //Contract address of Token BSH

	bep, err := bep20tkn.NewBEP(common.HexToAddress(btp_hmny_erc20), cleth)
	if err != nil {
		log.Fatal("Error setting BEP  ", err)
	} else {
		opts := bind.CallOpts{Pending: false, Context: ctx}
		if val, err := bep.BalanceOf(&opts, common.HexToAddress("0x8fc668275b4fa032342ea3039653d841f069a83b")); err != nil {
			log.Fatal("Error retrieving balance  ", err)
		} else {
			fmt.Printf("Initial Value %d \n", val)
		}
	}

	// Transfer
	fmt.Println("Transfer balance")
	btp_hmny_god_wallet := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.god.wallet.json"
	btp_hmny_demo_wallet_address := "0x8fc668275b4fa032342ea3039653d841f069a83b"
	keyReader, err := os.Open(btp_hmny_god_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()
	btp_hmny_god_wallet_password := ""
	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}
	w, _ := DecryptHMYKeyStore(keyStore, []byte(btp_hmny_god_wallet_password))
	chainID, _ := cleth.ChainID(context.Background())

	txo, err := bind.NewKeyedTransactorWithChainID(w.Skey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewInt(234452180000000000)
	const DefaultGasLimit = 80000000
	txo.GasPrice, _ = cleth.SuggestGasPrice(context.Background())
	txo.GasLimit = uint64(DefaultGasLimit)
	//txo.Value = amount

	if txn, err := bep.Transfer(txo, common.HexToAddress(btp_hmny_demo_wallet_address), amount); err != nil {
		log.Fatal(err)
	} else {
		v, _ := txn.MarshalJSON()
		fmt.Println(string(v))
	}

	opts := bind.CallOpts{Pending: false, Context: ctx}
	if val, err := bep.BalanceOf(&opts, common.HexToAddress("0x8fc668275b4fa032342ea3039653d841f069a83b")); err != nil {
		log.Fatal("Error retrieving balance  ", err)
	} else {
		fmt.Printf("Final Value %d \n", val)
	}
}

func DecryptHMYKeyStore(data []byte, pw []byte) (*wallet.EvmWallet, error) {
	key, err := wallet.DecryptEvmKeyStore(data, pw)
	if err != nil {
		return nil, err
	}
	return &wallet.EvmWallet{
		Skey: key,
		Pkey: &key.PublicKey,
	}, nil
}

func TestHmnyTransferNativeCoin(t *testing.T) {
	/*
	   WALLET=$btp_hmny_demo_wallet \
	       PASSWORD=$btp_hmny_demo_wallet_password \
	       run_exec hmnyTransferNativeCoin \
	       $h2i_nativecoin_transfer_amount \
	       "btp://$btp_icon_net/$btp_icon_demo_wallet_address"
	*/
	url := "http://127.0.0.1:9500"

	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	cleth := ethclient.NewClient(clrpc)

	btp_hmny_god_wallet := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.god.wallet.json"
	//btp_hmny_demo_wallet_address := "0x8fc668275b4fa032342ea3039653d841f069a83b"
	btp_icon_demo_wallet_address := "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	keyReader, err := os.Open(btp_hmny_god_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()
	btp_hmny_god_wallet_password := ""
	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}
	w, _ := DecryptHMYKeyStore(keyStore, []byte(btp_hmny_god_wallet_password))
	chainID, _ := cleth.ChainID(context.Background())

	txOps, err := bind.NewKeyedTransactorWithChainID(w.Skey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewInt(4452180000000000)
	const DefaultGasLimit = 80000000
	txOps.GasPrice, _ = cleth.SuggestGasPrice(context.Background())
	txOps.GasLimit = uint64(DefaultGasLimit)
	txOps.Value = amount

	btp_hmny_nativecoin_bsh_core := "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95"
	bshc, err := bshcore.NewBshcore(common.HexToAddress(btp_hmny_nativecoin_bsh_core), cleth)
	if err != nil {
		log.Fatal(err)
	} else {
		if txn, err := bshc.TransferNativeCoin(txOps, btp_icon_demo_wallet_address); err != nil {
			log.Fatal(err)
		} else {
			v, _ := txn.MarshalJSON()
			fmt.Println(string(v))
		}
	}
}

func TestHmnyBSHApprove(t *testing.T) {
	/*
	   echo "Approve HMNY BSHCore to access $btp_icon_nativecoin_symbol"
	   WALLET=$btp_hmny_demo_wallet \
	       PASSWORD=$btp_hmny_demo_wallet_password \
	       run_exec hmnyBSHApprove "$btp_icon_nativecoin_symbol" \
	       "$btp_hmny_nativecoin_bsh_core" 100000000000000000000000 >/dev/null # 100000
	*/
	btp_hmny_nativecoin_bsh_core := "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95"
	coinName := "ICX"
	url := "http://127.0.0.1:9500"

	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	cleth := ethclient.NewClient(clrpc)

	amount := 100000000000

	bshc, err := bshcore.NewBshcore(common.HexToAddress(btp_hmny_nativecoin_bsh_core), cleth)
	if err != nil {
		log.Fatal(err)
	}
	copts := &bind.CallOpts{Pending: false, Context: nil}
	coinAddress, err := bshc.CoinId(copts, coinName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(coinAddress)

	btp_hmny_demo_wallet := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.demo.wallet.json"
	keyReader, err := os.Open(btp_hmny_demo_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()
	btp_hmny_demo_wallet_password := "1234"
	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}
	w, _ := DecryptHMYKeyStore(keyStore, []byte(btp_hmny_demo_wallet_password))
	chainID, _ := cleth.ChainID(context.Background())

	txOps, err := bind.NewKeyedTransactorWithChainID(w.Skey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	const DefaultGasLimit = 80000000
	txOps.GasPrice, _ = cleth.SuggestGasPrice(context.Background())
	txOps.GasLimit = uint64(DefaultGasLimit)

	erc, err := erc20.NewErc20tradable(coinAddress, cleth)
	if err != nil {
		log.Fatal(err)
	}
	txn, err := erc.Approve(txOps, common.HexToAddress(btp_hmny_nativecoin_bsh_core), big.NewInt(int64(amount)))
	if err != nil {
		log.Fatal(err)
	} else {
		v, _ := txn.MarshalJSON()
		fmt.Println("TXN   ", string(v))
	}

}

func TestHmnyBSHAllowance(t *testing.T) {
	/*
	   echo "Allowance: $(format_token $(hex2dec 0x$(run_exec hmnyBSHAllowance \
	       $btp_icon_nativecoin_symbol $btp_hmny_demo_wallet_address \
	       $btp_hmny_nativecoin_bsh_core)))"
	*/
	btp_icon_nativecoin_symbol := "ICX"
	btp_hmny_demo_wallet_address := "0x8fc668275b4fa032342ea3039653d841f069a83b"
	btp_hmny_nativecoin_bsh_core := "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95"
	coinName := btp_icon_nativecoin_symbol
	owner := btp_hmny_demo_wallet_address
	spender := btp_hmny_nativecoin_bsh_core

	url := "http://127.0.0.1:9500"

	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	cleth := ethclient.NewClient(clrpc)

	bshc, err := bshcore.NewBshcore(common.HexToAddress(btp_hmny_nativecoin_bsh_core), cleth)
	if err != nil {
		log.Fatal(err)
	}
	copts := &bind.CallOpts{Pending: false, Context: nil}
	coinAddress, err := bshc.CoinId(copts, coinName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(coinAddress)

	erc, err := erc20.NewErc20tradable(coinAddress, cleth)
	if err != nil {
		log.Fatal(err)
	}
	txn, err := erc.Allowance(&bind.CallOpts{Pending: false, Context: context.Background()}, common.HexToAddress(owner), common.HexToAddress(spender))
	if err != nil {
		log.Fatal(err)
	} else {
		v, _ := txn.MarshalJSON()
		fmt.Println("TXN   ", string(v))
	}

}

func TestHmnyTransferWrappedCoin(t *testing.T) {
	/*
	   h2i_wrapped_ICX_transfer_amount=1000000000000000000 # 1 ICX
	   echo "Transfer Wrapped ICX (HMNY -> ICON):"
	   echo "    amount=$(format_token $h2i_wrapped_ICX_transfer_amount)"
	   echo -n "    "
	   WALLET=$btp_hmny_demo_wallet \
	       PASSWORD=$btp_hmny_demo_wallet_password \
	       run_exec hmnyTransferWrappedCoin \
	       $btp_icon_nativecoin_symbol \
	       $h2i_wrapped_ICX_transfer_amount \
	       "btp://$btp_icon_net/$btp_icon_demo_wallet_address" >/dev/null

	*/
	url := "http://127.0.0.1:9500"

	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	cleth := ethclient.NewClient(clrpc)
	btp_hmny_nativecoin_bsh_core := "0x05AcF27495FAAf9A178e316B9Da2f330983b9B95"
	bshc, err := bshcore.NewBshcore(common.HexToAddress(btp_hmny_nativecoin_bsh_core), cleth)
	if err != nil {
		log.Fatal(err)
	}

	btp_hmny_demo_wallet := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.demo.wallet.json"
	keyReader, err := os.Open(btp_hmny_demo_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()
	btp_hmny_demo_wallet_password := "1234"
	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}
	w, _ := DecryptHMYKeyStore(keyStore, []byte(btp_hmny_demo_wallet_password))
	chainID, _ := cleth.ChainID(context.Background())

	txOps, err := bind.NewKeyedTransactorWithChainID(w.Skey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(1000000000)
	const DefaultGasLimit = 80000000
	txOps.GasPrice, _ = cleth.SuggestGasPrice(context.Background())
	txOps.GasLimit = uint64(DefaultGasLimit)

	btp_icon_demo_wallet_address := "btp://0x5b9a77.icon/hx691ead88bd5945a43c8a1da331ff6dd80e2936ee"
	txn, err := bshc.Transfer(txOps, "ICX", amount, btp_icon_demo_wallet_address)
	if err != nil {
		log.Fatal(err)
	} else {
		v, _ := txn.MarshalJSON()
		fmt.Println("TXN   ", string(v))
	}
}

func TestBEPTKNApprove(t *testing.T) {
	/*
	   WALLET=$btp_hmny_demo_wallet \
	       PASSWORD=$btp_hmny_demo_wallet_password \
	       run_sol >/dev/null \
	       TokenBSH.BEP20TKN.approve \
	       "'$btp_hmny_token_bsh_proxy','$h2i_erc20_ETH_transfer_amount'"
	*/

	url := "http://127.0.0.1:9500"

	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	cleth := ethclient.NewClient(clrpc)
	btp_hmny_erc20 := "0xb54f5e97972AcF96470e02BE0456c8DB2173f33a" //Contract address of Token BSH
	beptkn, err := bep20tkn.NewBEP(common.HexToAddress(btp_hmny_erc20), cleth)
	if err != nil {
		log.Fatal(err)
	}

	btp_hmny_demo_wallet := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.demo.wallet.json"
	keyReader, err := os.Open(btp_hmny_demo_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()
	btp_hmny_demo_wallet_password := "1234"
	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}
	w, _ := DecryptHMYKeyStore(keyStore, []byte(btp_hmny_demo_wallet_password))
	chainID, _ := cleth.ChainID(context.Background())
	btp_hmny_token_bsh_proxy := "0x48cacC89f023f318B4289A18aBEd44753a127782"

	txOps, err := bind.NewKeyedTransactorWithChainID(w.Skey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	amount := big.NewInt(1000000000)
	const DefaultGasLimit = 80000000
	txOps.GasPrice, _ = cleth.SuggestGasPrice(context.Background())
	txOps.GasLimit = uint64(DefaultGasLimit)

	txn, err := beptkn.Approve(txOps, common.HexToAddress(btp_hmny_token_bsh_proxy), amount)
	if err != nil {
		log.Fatal(err)
	} else {
		v, _ := txn.MarshalJSON()
		fmt.Println("TXN   ", string(v))
	}
}

func TestHmnyTransfer(t *testing.T) {
	url := "http://127.0.0.1:9500"

	clrpc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	btp_hmny_god_wallet_address := "0xa5241513da9f4463f1d4874b548dfbac29d91f34"
	btp_hmny_demo_wallet_address := "0x606f95a0d893ab26aa3e7dd9ce33530bca0e6dbf"
	cleth := ethclient.NewClient(clrpc)
	nonce, err := cleth.PendingNonceAt(context.Background(), common.HexToAddress(btp_hmny_god_wallet_address))
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewInt(100000000000)
	const DefaultGasLimit = 80000000
	gasPrice, _ := cleth.SuggestGasPrice(context.Background())
	gasLimit := uint64(DefaultGasLimit)

	tx := types.NewTransaction(nonce, common.HexToAddress(btp_hmny_demo_wallet_address), amount, gasLimit, gasPrice, []byte{})
	chainID, err := cleth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	btp_hmny_god_wallet := "/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/hmny.god.wallet.json"

	keyReader, err := os.Open(btp_hmny_god_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()
	btp_hmny_god_wallet_password := ""
	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}

	key, err := wallet.DecryptEvmKeyStore(keyStore, []byte(btp_hmny_god_wallet_password))
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	if err != nil {
		log.Fatal(err)
	}
	if err := cleth.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatal(err)
	}
}

func TestIconCreateWallet(t *testing.T) {
	priv, _ := crypto.GenerateKeyPair()

	pass := "1234"
	w, err := wallet.EncryptKeyAsKeyStore(priv, []byte(pass))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(w))
}

func TestCreateKeystore(t *testing.T) {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3
}
