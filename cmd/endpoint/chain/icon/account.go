package icon

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/haltingstate/secp256k1-go"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
)

func GetWalletFromFile(walFile string, password string) (module.Wallet, error) {
	keyReader, err := os.Open(walFile)
	if err != nil {
		return nil, err
	}
	defer keyReader.Close()

	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		return nil, err
	}
	w, err := wallet.NewFromKeyStore(keyStore, []byte(password))
	if err != nil {
		return nil, err
	}
	return w, nil
}

func GetWalletFromPrivKey(privKey string) (module.Wallet, error) {
	privBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	pKey, err := crypto.ParsePrivateKey(privBytes)
	if err != nil {
		return nil, err
	}
	wal, err := wallet.NewFromPrivateKey(pKey)
	if err != nil {
		return nil, err
	}
	return wal, nil
}

func CreateKeyStore(password string) (*string, error) {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return nil, err
	}
	addr := account.Address.Hex()
	return &addr, nil
}

func getAddressFromPrivKey(pKey string) (*string, error) {
	privBytes, err := hex.DecodeString(pKey)
	if err != nil {
		return nil, err
	}
	pubkeyBytes := secp256k1.PubkeyFromSeckey(privBytes)
	pubKey, err := crypto.ParsePublicKey(pubkeyBytes)
	if err != nil {
		return nil, err
	}
	addr := common.NewAccountAddressFromPublicKey(pubKey).String()
	return &addr, nil
}

func generateKeyPair() ([2]string, error) {
	pubkeyBytes, priv := secp256k1.GenerateKeyPair()
	pubKey, err := crypto.ParsePublicKey(pubkeyBytes)
	if err != nil {
		return [2]string{}, err
	}
	addr := common.NewAccountAddressFromPublicKey(pubKey).String()
	return [2]string{hex.EncodeToString(priv), addr}, nil
}

// func (r *requestAPI) GetAddressFromPrivKey(key string) (*string, error) {
// 	return getAddressFromPrivKey(key)
// }
