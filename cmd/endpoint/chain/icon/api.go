package icon

import (
	"context"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
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
	} // TODO: Remove these contract name constants
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

func NewRequestAPI(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (*requestAPI, error) {
	cAddr := &contractAddress{}
	cAddr.FromMap(contractAddrsMap)

	return newARequestPI(url, l, cAddr, networkID)
}

func (r *requestAPI) Transfer(param *chain.RequestParam) (txnHash string, err error) {
	if param.FromChain != chain.ICON {
		err = errors.New("Source Chan should be Icon")
		return
	}
	if param.ToChain == chain.ICON {
		if param.Token == chain.ICXToken {
			txnHash, _, err = r.transferICX(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.IRC2Token {
			txnHash, _, err = r.transferIrc2(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else if param.ToChain == chain.HMNY {
		if param.Token == chain.ICXToken {
			txnHash, _, err = r.TransferICXToHarmony(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.IRC2Token {
			txnHash, _, err = r.transferIrc2ToHmny(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.ONEToken {
			txnHash, _, err = r.transferWrappedOneFromIconToHmny(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else {
		err = errors.New("Unsupport Transaction Parameters ")
	}
	return
}

func (r *requestAPI) GetCoinBalance(addr string, coinType chain.TokenType) (*big.Int, error) {
	if coinType == chain.ICXToken {
		return r.getICXBalance(addr)
	} else if coinType == chain.IRC2Token {
		return r.getIrc2Balance(addr)
	} else if coinType == chain.ONEToken {
		return r.getIconWrappedOne(addr)
	}
	return nil, errors.New("Unsupported Token Type ")
}

func (r *requestAPI) WaitForTxnResult(hash string) (txr interface{}, err error) {
	_, txr, err = r.cl.waitForResults(context.TODO(), &TransactionHashParam{Hash: HexBytes(hash)})
	return
}

func (r *requestAPI) Approve(ownerKey string, amount big.Int) (txnHash string, err error) {
	txnHash, _, _, err = r.approveIconNativeCoinBSHToAccessHmnyOne(ownerKey, amount)
	return
}

func (r *requestAPI) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + r.networkID + ".icon/" + addr
	return &fullAddr
}

func (r *requestAPI) GetKeyPairs(num int) ([][2]string, error) {
	var err error
	res := make([][2]string, num)
	for i := 0; i < num; i++ {
		res[i], err = generateKeyPair()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
