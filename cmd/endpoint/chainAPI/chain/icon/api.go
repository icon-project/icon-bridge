package icon

import (
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

type requestAPI struct {
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

func newARequestPI(uri string, l log.Logger, cAddr *contractAddress, networkID string) (*requestAPI, error) {
	cl, err := newClient(uri, l)
	if err != nil {
		return nil, err
	}
	return &requestAPI{networkID: networkID, contractAddress: cAddr, cl: cl}, nil
}

func NewRequestAPI(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (chain.RequestAPI, error) {
	cAddr := &contractAddress{}
	cAddr.FromMap(contractAddrsMap)

	return newARequestPI(url, l, cAddr, networkID)
}

func (r *requestAPI) GetCoinBalance(addr string) (*big.Int, error) {
	return r.GetICXBalance(addr)
}

func (r *requestAPI) GetEthToken(addr string) (val *big.Int, err error) {
	return r.GetIrc2Balance(addr)
}

func (r *requestAPI) GetWrappedCoin(addr string) (val *big.Int, err error) {
	return r.GetIconWrappedOne(addr)
}

func (r *requestAPI) TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferICX(senderKey, amount, recepientAddress)
}

func (r *requestAPI) TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferIrc2(senderKey, amount, recepientAddress)
}

func (r *requestAPI) TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferICXToHarmony(senderKey, amount, recepientAddress)
}

func (r *requestAPI) TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferWrappedOneFromIconToHmny(senderKey, amount, recepientAddress)
}

func (r *requestAPI) TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	return r.TransferIrc2ToHmny(senderKey, amount, recepientAddress)
}

func (r *requestAPI) ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	return r.ApproveIconNativeCoinBSHToAccessHmnyOne(ownerKey, amount)
}

func (r *requestAPI) GetAddressFromPrivKey(key string) (*string, error) {
	return getAddressFromPrivKey(key)
}

func (r *requestAPI) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + r.networkID + ".icon/" + addr
	return &fullAddr
}
