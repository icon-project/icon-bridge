package icon

import (
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

type api struct {
	contractAddress *contractAddress
	networkID       string
	cl              *client
}
type contractAddress struct {
	btp_icon_irc2           string
	btp_icon_irc2_tradeable string
	btp_icon_nativecoin_bsh string
	btp_icon_token_bsh      string
}

func (cAddr *contractAddress) FromMap(contractAddrsMap map[string]string) {
	if contractAddrsMap == nil {
		return
	}
	cAddr.btp_icon_irc2 = contractAddrsMap["btp_icon_irc2"]
	cAddr.btp_icon_irc2_tradeable = contractAddrsMap["btp_icon_irc2_tradeable"]
	cAddr.btp_icon_nativecoin_bsh = contractAddrsMap["btp_icon_nativecoin_bsh"]
	cAddr.btp_icon_token_bsh = contractAddrsMap["btp_icon_token_bsh"]
}

func newAPI(uri string, l log.Logger, cAddr *contractAddress, networkID string) (*api, error) {
	cl, err := newClient(uri, l)
	if err != nil {
		return nil, err
	}
	return &api{networkID: networkID, contractAddress: cAddr, cl: cl}, nil
}

func New(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (chain.API, error) {
	cAddr := &contractAddress{}
	cAddr.FromMap(contractAddrsMap)

	return newAPI(url, l, cAddr, networkID)
}

func (a *api) GetCoinBalance(addr string) (*big.Int, error) {
	return a.GetICXBalance(addr)
}

func (a *api) GetEthToken(addr string) (val *big.Int, err error) {
	return a.GetIrc2Balance(addr)
}

func (a *api) GetWrappedCoin(addr string) (val *big.Int, err error) {
	return a.GetIconWrappedOne(addr)
}

func (a *api) TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferICX(senderKey, amount, recepientAddress)
}

func (a *api) TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferIrc2(senderKey, amount, recepientAddress)
}

func (a *api) TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferICXToHarmony(senderKey, amount, recepientAddress)
}

func (a *api) TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferWrappedOneFromIconToHmny(senderKey, amount, recepientAddress)
}

func (a *api) TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	return a.TransferIrc2ToHmny(senderKey, amount, recepientAddress)
}

func (a *api) ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	return a.ApproveIconNativeCoinBSHToAccessHmnyOne(ownerKey, amount)
}

func (a *api) GetAddressFromPrivKey(key string) (*string, error) {
	return getAddressFromPrivKey(key)
}

func (a *api) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + a.networkID + ".icon/" + addr
	return &fullAddr
}
