package wallet

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"log"
)

type AvmWallet struct {
	Skey *ed25519.PrivateKey
	Pkey *ed25519.PublicKey
}

func (w *AvmWallet) PublicKey() []byte {
	pubKey := w.Skey.Public().(ed25519.PublicKey)
	return pubKey
}

func (w *AvmWallet) Address() string {
	pubBytes := w.PublicKey()
	if len(pubBytes) != 32 {
		log.Panic("pubkey is incorrect size")
	}
	address := hex.EncodeToString(pubBytes)
	return address

}

func (w *AvmWallet) Sign(data []byte) ([]byte, error) {
	signature, err := w.Skey.Sign(rand.Reader, data, crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (w *AvmWallet) ECDH(pubkey []byte) ([]byte, error) {
	//Need to be implemnted
	return nil, nil
}

func NewAvmWalletFromPrivateKey(sk *ed25519.PrivateKey) (*AvmWallet, error) {
	pkey := sk.Public().(ed25519.PublicKey)
	return &AvmWallet{
		Skey: sk,
		Pkey: &pkey,
	}, nil
}
