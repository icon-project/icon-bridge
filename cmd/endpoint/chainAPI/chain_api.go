package chainAPI

import (
	"context"
	"math/big"
	"reflect"

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
	log     log.Logger
	req     map[chain.ChainType]chain.RequestAPI
	sub     map[chain.ChainType]chain.SubscriptionAPI
	subChan chan *chain.SubscribedEvent
	errChan chan error
}

type ChainAPI interface {
	Transfer(param *RequestParam) (txHash string, err error)
	StartSubscription(ctx context.Context) (err error)
	SubEventChan() <-chan *chain.SubscribedEvent
	ErrChan() <-chan error
}

func New(l log.Logger, configPerChain map[chain.ChainType]*chain.ChainConfig) (ChainAPI, error) {
	cAPI := chainAPI{
		log:     l,
		req:     make(map[chain.ChainType]chain.RequestAPI),
		sub:     make(map[chain.ChainType]chain.SubscriptionAPI),
		subChan: make(chan *chain.SubscribedEvent),
		errChan: make(chan error),
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
			cAPI.req[name] = req
			cAPI.sub[name] = sub
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
			cAPI.req[name] = req
			cAPI.sub[name] = sub
		} else {
			return nil, errors.New("Unknown Chain Type supplied from config: " + string(name))
		}
	}
	return &cAPI, nil
}

func (capi *chainAPI) Transfer(param *RequestParam) (txHash string, err error) {
	reqApi := capi.req[param.FromChain]
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

func (capi *chainAPI) StartSubscription(ctx context.Context) (err error) {
	go func() {
		defer func() {
			close(capi.subChan)
			close(capi.errChan)
			capi.log.Warn("Closing Supscription API")
		}()
		cases := make([]reflect.SelectCase, 1+len(capi.sub)*2) // context + (sub + err)
		i := 0
		for _, sub := range capi.sub {
			if err := sub.Start(ctx); err != nil {
				return
			}
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sub.OutputChan())}
			i++
		}
		i = len(capi.sub)
		for _, sub := range capi.sub {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sub.ErrChan())}
			i++
		}
		cases[len(cases)-1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}
		for {
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				if chosen == len(cases)-1 { // Context channel has been closed
					capi.log.Error("Context cancelled")
					return
				}
				cases[chosen].Chan = reflect.ValueOf(nil)
				continue
			}
			if chosen < len(capi.sub) { // [0, lenCapi-1] is message
				res, dok := value.Interface().(*chain.SubscribedEvent)
				if !dok {
					capi.log.Error("Wrong interface; Expected *SubscribedEvent")
					break
				}
				capi.subChan <- res
			} else if chosen >= len(capi.sub) && chosen < 2*len(cases) {
				res, eok := value.Interface().(error)
				if !eok {
					capi.log.Error("Wrong interface; Expected errorType")
					break
				}
				capi.errChan <- res
			} else { // last element is context
				capi.log.Error("Context cancelled")
				return
			}
		}
	}()
	return
}

func (capi *chainAPI) SubEventChan() <-chan *chain.SubscribedEvent {
	return capi.subChan
}

func (capi *chainAPI) ErrChan() <-chan error {
	return capi.errChan
}
