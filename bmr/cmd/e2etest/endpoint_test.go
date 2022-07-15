package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/haltingstate/secp256k1-go"
	gocrypto "github.com/icon-project/goloop/common/crypto"
	gowallet "github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/icon-bridge/bmr/common/wallet"
)

func TestEncryption(t *testing.T) {
	_, priv := secp256k1.GenerateKeyPair()
	// _, err := crypto.ToECDSA(priv)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// addr := common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:]).Hex()

	h := hex.EncodeToString(priv)
	b, err := hex.DecodeString(h)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(priv, h)
	fmt.Println(b, bytes.Equal(priv, b))
}

func TestPubUnmarshal(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privStr := hex.EncodeToString(crypto.FromECDSA(privKey))
	// pubKey, _ := crypto.ParsePublicKey(pub)
	// addr := common.NewAccountAddressFromPublicKey(pubKey).String()
	if err != nil {
		log.Fatal(err)
	}
	pubAddress := crypto.PubkeyToAddress(privKey.PublicKey).String()

	wal := &wallet.EvmWallet{
		Skey: privKey,
		Pkey: &privKey.PublicKey,
	}
	fmt.Println(wal.Address(), wal.PublicKey(), pubAddress)
	newWallet, err := GetWalletFromPrivKey(privStr)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(newWallet.Address(), newWallet.PublicKey())
	}
}

func GetWalletFromPrivKey(privKey string) (module.Wallet, error) {
	privBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	pKey, err := gocrypto.ParsePrivateKey(privBytes)
	if err != nil {
		return nil, err
	}
	wal, err := gowallet.NewFromPrivateKey(pKey)
	if err != nil {
		return nil, err
	}
	return wal, nil
}
