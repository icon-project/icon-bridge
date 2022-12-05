package wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//TODO find which cryptography protocol Algorand uses

type AvmWallet struct {
	Skey *ecdsa.PrivateKey
	Pkey *ecdsa.PublicKey
}

func (w *AvmWallet) Address() string {
	pubBytes := w.PublicKey()
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:]).Hex()
}

func (w *AvmWallet) Sign(data []byte) ([]byte, error) {
	//TODO: Not implemented yet
	return nil, errors.New("Not implemented yet")
}

func (w *AvmWallet) PublicKey() []byte {
	return crypto.FromECDSAPub(w.Pkey)
}

func (w *AvmWallet) ECDH(pubKey []byte) ([]byte, error) {
	//TODO: Not implemented yet
	return nil, nil
}

func NewAvmWalletFromPrivateKey(sk *ecdsa.PrivateKey) (*AvmWallet, error) {
	return &AvmWallet{
		Skey: sk,
		Pkey: &sk.PublicKey,
	}, nil
}
