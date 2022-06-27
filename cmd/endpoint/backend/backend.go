package backend

import (
	"context"

	"github.com/icon-project/icon-bridge/common/log"

	capi "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	wtchr "github.com/icon-project/icon-bridge/cmd/endpoint/watcher"
)

type backend struct {
	log      log.Logger
	chainapi capi.ChainAPI
	wtch     wtchr.Watcher
}

type Backend interface {
	Start(ctx context.Context) error
	Transfer(param *capi.RequestParam) (txHash string, err error)
}

func New(l log.Logger, configPerChain map[chain.ChainType]*chain.ChainConfig) (Backend, error) {
	var err error
	be := &backend{log: l}
	be.chainapi, err = capi.New(l, configPerChain)
	if err != nil {
		return nil, err
	}
	be.wtch, err = wtchr.New(l, configPerChain, be.chainapi.SubEventChan(), be.chainapi.ErrChan())
	if err != nil {
		return nil, err
	}
	return be, nil
}

func (be *backend) Start(ctx context.Context) (err error) {
	if err = be.wtch.Start(ctx); err != nil {
		return
	}
	if err = be.chainapi.StartSubscription(ctx); err != nil {
		return
	}
	return
}

func (be *backend) Transfer(param *capi.RequestParam) (txHash string, err error) {
	txHash, logs, err := be.chainapi.Transfer(param)
	if err != nil {
		return
	}
	err = be.wtch.ProcessTxn(param, logs)
	return
}
