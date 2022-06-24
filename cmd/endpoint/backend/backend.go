package backend

import (
	"context"
	"sync"

	"github.com/icon-project/icon-bridge/common/log"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
)

const EventListenerCapacity = 1000

type backend struct {
	log                       log.Logger
	availableCBMutex          sync.RWMutex
	availableCBList           []MatchingFn
	eventListenerCursor       uint64
	eventListenerMutex        sync.RWMutex
	eventListenerList         []*eventListener
	requestAPIPerChain        map[chain.ChainType]chain.RequestAPI
	subscriptionAPIPerChain   map[chain.ChainType]chain.SubscritionAPI
	subscriptionReceivingChan chan *chain.SubscribedEvent
	subscriptionErrChan       chan error
}

type eventListener struct {
	fn EventMatchingFn
	ip *RequestParam
}

type Backend interface {
	Start(ctx context.Context)
	Transfer(param *RequestParam) (txHash string, err error)
	Accounts(num int) error
	RegisterCriterion(fn MatchingFn) error
	// Approve(param *TxnApproveParam) (txnHash string, err error)
}

func New(l log.Logger, configPerChain map[chain.ChainType]*chain.ChainConfig) (Backend, error) {
	var err error
	be := &backend{
		log:                       l,
		availableCBList:           DefaultCBs,
		availableCBMutex:          sync.RWMutex{},
		eventListenerCursor:       0,
		eventListenerList:         make([]*eventListener, EventListenerCapacity),
		eventListenerMutex:        sync.RWMutex{},
		requestAPIPerChain:        map[chain.ChainType]chain.RequestAPI{},
		subscriptionAPIPerChain:   map[chain.ChainType]chain.SubscritionAPI{},
		subscriptionReceivingChan: make(chan *chain.SubscribedEvent),
		subscriptionErrChan:       make(chan error),
	}
	for name, cfg := range configPerChain {
		if name == chain.HMNY {
			be.requestAPIPerChain[name], err = hmny.NewRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
			be.subscriptionAPIPerChain[name], err = hmny.NewSubscriptionAPI(l, cfg.Subscriber, cfg.URL) // config
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
		} else if name == chain.ICON {
			be.requestAPIPerChain[name], err = icon.NewRequestAPI(cfg.URL, l, cfg.ConftractAddresses, cfg.NetworkID)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
			be.subscriptionAPIPerChain[name], err = icon.NewSubscriptionAPI(l, cfg.Subscriber, cfg.URL)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
		} else {
			return nil, errors.New("Unknown Chain Type supplied from config: " + string(name))
		}
	}
	return be, nil
}

func (be *backend) Start(ctx context.Context) {
	for _, sub := range be.subscriptionAPIPerChain {
		err := sub.Start(ctx, be.subscriptionReceivingChan, be.subscriptionErrChan)
		if err != nil {
			be.log.Error(err)
			return
		}
	}
	go func() {
		defer func() {
			if be.subscriptionErrChan != nil {
				close(be.subscriptionErrChan)
			}
			if be.subscriptionAPIPerChain != nil {
				close(be.subscriptionReceivingChan)
			}
			be.log.Warn("Exiting Backend's processing function")
		}()
		for {
			select {
			case <-ctx.Done():
				be.log.Error("Context cancelled")
				return
			case el := <-be.subscriptionReceivingChan:
				if el.ChainName == chain.ICON {
					iconTxn := (el.Res).(*icon.TxnLog)
					decodeIconLog(iconTxn)
					// be.log.Warn(*devt)
				}
			// if index, ok := be.checkIfLogMatchesListener(el); ok {
			// 	be.removeFromEventListenerQueue(index)
			// }
			case err := <-be.subscriptionErrChan:
				be.log.Error(err)
				return
			}
		}
	}()
}

var DefaultCBs = []MatchingFn{
	{
		ReqFn: func(req *RequestParam) bool {
			return true
		},
		EventFn: func(*ReceiptEvent, *RequestParam) bool {
			return true
		},
	},
}
