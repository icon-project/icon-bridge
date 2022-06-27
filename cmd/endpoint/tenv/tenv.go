package tenv

import (
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

/*

-Initial Setup of Test Unit ? Accounts and such maybe
-What are the sources of variability for a test ? Amount, Sender, Receiver
-How are these varaibles defined ? Do I send an amount greater than my holdings ? Do I send to multiple receivers ?
-How are the processes carried out ? Do the processes follow a certain sequence and make necessary checks during the stages ?
-What is the output of a process and how much of the result do I need to know OR what must be the format of result ?
-How is the result validated ? Expected Outcome vs Observed Outcomes and other comparables ?
-Closure of Test ? Return Result
*/
type tenv struct {
	l                log.Logger
	clientsPerChain  map[chain.ChainType]chain.RequestAPI
	accountsPerChain map[chain.ChainType][]string
	godKeysPerChain  map[chain.ChainType][2]string
}

type TEnv interface {
	GetClient(name chain.ChainType) (chain.RequestAPI, error)
	GetAccounts(name chain.ChainType) ([][2]string, error)
	GetGodKeyPair(name chain.ChainType) ([2]string, error)
	GetEnvVariables(name chain.ChainType) (*chain.EnvVariables, error)
	Logger() log.Logger
}

func New(l log.Logger, clientsPerChain map[chain.ChainType]*chain.ChainConfig, accountsPerChain map[chain.ChainType][]string, godKeysPerChain map[chain.ChainType][2]string) (t TEnv, err error) {
	tu := &tenv{l: l,
		accountsPerChain: accountsPerChain,
		clientsPerChain:  map[chain.ChainType]chain.RequestAPI{},
		godKeysPerChain:  godKeysPerChain,
	}
	for name, cfg := range clientsPerChain {
		if name == chain.HMNY {
			tu.clientsPerChain[name], err = hmny.NewRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
			}
		} else if name == chain.ICON {
			tu.clientsPerChain[name], err = icon.NewRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
			}
		} else {
			return nil, errors.New("Unknown Chain Type supplied from config: " + string(name))
		}
	}
	return tu, nil
}

func (tv *tenv) GetClient(name chain.ChainType) (chain.RequestAPI, error) {
	if cl, ok := tv.clientsPerChain[name]; ok {
		return cl, nil
	}
	return nil, errors.New("Clients of chain type not found " + string(name))
}

func (tv *tenv) GetAccounts(name chain.ChainType) ([][2]string, error) {
	pKeys, ok := tv.accountsPerChain[name]
	if !ok {
		return nil, errors.New("Client of Chain type not found " + string(name))
	}
	keyPairs := make([][2]string, len(pKeys))
	for i, pKey := range pKeys {
		addr, err := tv.getAddressFromPrivKey(pKey, name)
		if err != nil {
			return nil, err
		} else {
			keyPairs[i] = [2]string{pKey, *addr}
		}
	}

	return keyPairs, nil
}

func (tv *tenv) getAddressFromPrivKey(pKey string, name chain.ChainType) (*string, error) {
	cl, err := tv.GetClient(name)
	if err != nil {
		return nil, errors.Wrap(err, " Need client api to get account address ")
	}
	return cl.GetAddressFromPrivKey(pKey)

}

func (tv *tenv) Logger() log.Logger {
	return tv.l
}

func (tv *tenv) GetGodKeyPair(name chain.ChainType) (pair [2]string, err error) {
	var ok bool
	if pair, ok = tv.godKeysPerChain[name]; !ok {
		err = errors.New("God Keypair for chain not found")
	}
	return
}

func (tv *tenv) GetEnvVariables(name chain.ChainType) (env *chain.EnvVariables, err error) {
	env = &chain.EnvVariables{}
	if env.GodKeys, err = tv.GetGodKeyPair(name); err != nil {
		return
	}
	if env.AccountsKeys, err = tv.GetAccounts(name); err != nil {
		return
	}
	if env.Client, err = tv.GetClient(name); err != nil {
		return
	}
	return
}

/*
func getHmnyAddressFromPrivKey(pKey string) (*string, error) {
	privBytes, err := hex.DecodeString(pKey)
	if err != nil {
		return nil, err
	}
	privKey, err := ethCrypto.ToECDSA(privBytes)
	if err != nil {
		return nil, err
	}
	addr := ethCrypto.PubkeyToAddress(privKey.PublicKey).String()
	return &addr, nil
}

func getIconAddressFromPrivKey(pKey string) (*string, error) {
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
*/
