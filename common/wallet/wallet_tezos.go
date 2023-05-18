package wallet

import (
	"blockwatch.cc/tzgo/tezos"
)

type TezosWallet struct {
	Skey tezos.PrivateKey
	Pkey tezos.Key
}

func (w *TezosWallet) Address() string {
	return w.Pkey.Address().ContractAddress()
}

func (w *TezosWallet) Sign(data []byte) ([]byte, error) {
	skData, err := w.Skey.MarshalText()
	if err != nil {
		return nil, err
	}
	return skData, nil 
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
