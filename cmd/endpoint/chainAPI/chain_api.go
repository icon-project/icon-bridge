package chainAPI

import (
	"math/big"

	chain "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type TokenType string

const (
	ICXToken        TokenType = "ICX"
	IRC2Token       TokenType = "IRC2"
	ONEWrappedToken TokenType = "OneWrapped"
	ONEToken        TokenType = "ONE"
	ERC20Token      TokenType = "ERC20"
	ICXWrappedToken TokenType = "ICXWrapped"
)

type RequestParam struct {
	FromChain   chain.ChainType
	ToChain     chain.ChainType
	SenderKey   string
	FromAddress string
	ToAddress   string
	Amount      big.Int
	Token       TokenType
}

type chainAPI struct {
	Req map[chain.ChainType]chain.RequestAPI
	Sub map[chain.ChainType]chain.SubscritionAPI
}

type ChainAPI interface {
}

func New(l log.Logger, configPerChain map[chain.ChainType]*chain.ChainConfig) (ChainAPI, error) {
	cAPI := chainAPI{
		Req: make(map[chain.ChainType]chain.RequestAPI),
		Sub: make(map[chain.ChainType]chain.SubscritionAPI),
	}
	for name, cfg := range configPerChain {
		if name == chain.HMNY {
			req, err := hmny.NewRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
			sub, err := hmny.NewSubscriptionAPI(l, cfg.Subscriber, cfg.URL) // config
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
			cAPI.Req[name] = req
			cAPI.Sub[name] = sub
		} else if name == chain.ICON {
			req, err := icon.NewRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
			sub, err := icon.NewSubscriptionAPI(l, cfg.Subscriber, cfg.URL)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
			cAPI.Req[name] = req
			cAPI.Sub[name] = sub
		} else {
			return nil, errors.New("Unknown Chain Type supplied from config: " + string(name))
		}
	}
	return &cAPI, nil
}

func (capi *chainAPI) Transfer(param *RequestParam) (txHash string, err error) {
	reqApi := capi.Req[param.FromChain]
	if param.Token == ICXToken || param.Token == ONEToken { // Native Coin
		if param.FromChain == param.ToChain {
			txHash, err = reqApi.TransferCoin(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			txHash, err = reqApi.TransferCoinCrossChain(param.SenderKey, param.Amount, param.ToAddress)
		}
	} else if param.Token == IRC2Token || param.Token == ERC20Token { // EthToken
		if param.FromChain == param.ToChain {
			txHash, err = reqApi.TransferEthToken(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			_, txHash, err = reqApi.TransferEthTokenCrossChain(param.SenderKey, param.Amount, param.ToAddress)
		}
	} else if param.Token == ONEWrappedToken || param.Token == ICXWrappedToken { // WrappedToken
		if param.FromChain != param.ToChain {
			txHash, err = reqApi.TransferWrappedCoinCrossChain(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("Wrapped Coins can be transferred only cross chain")
		}
	} else {
		err = errors.New("Unrecognized token type ")
	}
	return
}
