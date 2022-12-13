package wallet

import (
	"crypto/ed25519"
	"log"

	"github.com/algorand/go-algorand-sdk/crypto"
)

type AvmWallet struct {
	account crypto.Account
}

func (w *AvmWallet) PublicKey() []byte {
	pubKey := w.account.PrivateKey.Public().(ed25519.PublicKey)
	return pubKey
}

func (w *AvmWallet) Address() string {
	return w.account.Address.String()
}

func (w *AvmWallet) Sign(data []byte) ([]byte, error) {

	signedTxn, err := crypto.SignBytes(w.account.PrivateKey, data)
	if err != nil {
		log.Fatalf("Cannot sign transaction: %s", err)
	}
	return signedTxn, nil
}

func (w *AvmWallet) ECDH(pubkey []byte) ([]byte, error) {
	//Need to be implemnted
	return nil, nil
}

func NewAvmWalletFromPrivateKey(sk *ed25519.PrivateKey) (*AvmWallet, error) {
	acc, err := crypto.AccountFromPrivateKey(*sk)
	if err != nil {
		log.Fatalf("Cannot create wallet from SK: %s", err)
	}
	return &AvmWallet{
		acc,
	}, nil
}
