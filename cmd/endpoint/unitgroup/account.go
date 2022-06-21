package unitgroup

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/pkg/errors"
)

func (ug *unitgroup) createAccounts(accountMap map[chain.ChainType]int) (map[chain.ChainType][]string, error) {
	resMap := map[chain.ChainType][]string{}
	for name, count := range accountMap {
		resMap[name] = make([]string, count)
		for i := 0; i < count; i++ {
			privKey, err := ethcrypto.GenerateKey()
			if err != nil {
				return nil, err
			}
			resMap[name][i] = hex.EncodeToString(ethcrypto.FromECDSA(privKey))
			// pubKey, _ := crypto.ParsePublicKey(pub)
			// addr := common.NewAccountAddressFromPublicKey(pubKey).String()
			if err != nil {
				return nil, errors.Wrap(err, "Unmarshal Public Key")
			}
		}
	}
	return resMap, nil
}

// ethcrypto used for ICON's keystore
// TODO: check if there can be differences
func GetKeyPairFromFile(walFile string, password string) (pair [2]string, err error) {
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
