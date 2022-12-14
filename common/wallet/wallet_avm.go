package wallet

import (
	"crypto/ed25519"
	"encoding/json"
	"log"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

type IWallet interface {
	Address() string
	Sign(data []byte) ([]byte, error)
	PublicKey() []byte
	ECDH(pubKey []byte) ([]byte, error)
	TypedAddress() types.Address
	SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error)
	Equals(other future.TransactionSigner) bool
}

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

func (w *AvmWallet) TypedAddress() types.Address {
	return w.account.Address
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

func (w AvmWallet) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		_, stxBytes, err := crypto.SignTransaction(w.account.PrivateKey, txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}
	return stxs, nil
}

func (w AvmWallet) Equals(other future.TransactionSigner) bool {
	if castedSigner, ok := other.(AvmWallet); ok {
		otherJson, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}
		selfJson, err := json.Marshal(w)
		if err != nil {
			return false
		}
		return string(otherJson) == string(selfJson)
	}
	return false
}
