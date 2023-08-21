package wallet

import (
	"errors"

	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
)

type TezosWallet struct {
	Skey tezos.PrivateKey
	Pkey tezos.Key
}

func (w *TezosWallet) Address() string {
	return w.Pkey.Address().String()
}

func (w *TezosWallet) Sign(data []byte) ([]byte, error) {
	return nil, errors.New("not allowed, use Signer instead")
}

func (w *TezosWallet) Signer() *signer.MemorySigner {
	return signer.NewFromKey(w.Skey)
}

func (w *TezosWallet) PublicKey() []byte {
	return w.Pkey.Bytes()
}

func (w *TezosWallet) ECDH(pubKey []byte) ([]byte, error) {
	//TODO: Not implemented yet
	return nil, nil
}

func NewTezosWalletFromPrivateKey(sk tezos.PrivateKey) (*TezosWallet, error) {
	return &TezosWallet{
		Skey: sk,
		Pkey: sk.Public(),
	}, nil
}
