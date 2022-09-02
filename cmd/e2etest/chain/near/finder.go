package near

import (
	"sync"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

type finder struct {
	log           log.Logger
	runCache      *runnableCache
	nameToAddrMap map[chain.ContractName]string
}

type runnableCache struct {
	mem []*runnable
	mtx sync.RWMutex
}

type runnable struct {
	args     args
	callback callBackFunc
}

type args struct {
	id              uint64
	eventType       chain.EventLogType
	seq             int64
	contractAddress string
}

type callBackFunc func(args args, info *chain.EventLogInfo) (bool, error)

func NewFinder(l log.Logger, nameToAddrMap map[chain.ContractName]string) *finder {
	return &finder{
		log:           l,
		runCache:      &runnableCache{mem: []*runnable{}, mtx: sync.RWMutex{}},
		nameToAddrMap: nameToAddrMap,
	}
}
