package chain

import "math/big"

type ChainType string

const (
	ICON ChainType = "ICON"
	HMNY ChainType = "HMNY"
)

type Client interface {
	GetCoinBalance(addr string) (*big.Int, error)
	GetEthToken(addr string) (val *big.Int, err error)
	GetWrappedCoin(addr string) (val *big.Int, err error)
	TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error)
	TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error)
	TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error)
	TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error)
	TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error)
	ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error)
	GetAddressFromPrivKey(key string) (*string, error)
	GetFullAddress(addr string) *string
}

type ChainConfig struct {
	Name               ChainType         `json:"name"`
	URL                string            `json:"url"`
	ConftractAddresses map[string]string `json:"contract_addresses"`
	GodWallet          GodWallet         `json:"god_wallet"`
	NetworkID          string            `json:"network_id"`
}

type ContractAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type GodWallet struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

type EnvVariables struct {
	Client       Client
	GodKeys      [2]string
	AccountsKeys [][2]string
}
