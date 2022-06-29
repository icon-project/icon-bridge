package icon

import (
	"sync"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

type finder struct {
	log      log.Logger
	runCache *runnableCache
}

type args struct {
	eventType       chain.EventLogType
	seq             int64
	contractAddress string
}

type callBackFunc func(args args, info chain.EventLogInfo) (bool, error)

type runnable struct {
	args     args
	callback callBackFunc
}

type runnableCache struct {
	mem []*runnable
	mtx sync.RWMutex
}

func (f *finder) WatchFor(args args) (err error) {
	if args.eventType == chain.TransferStart {
		f.addToRunCache(&runnable{args: args, callback: transferStartCB})
	} else if args.eventType == chain.TransferReceived {
		f.addToRunCache(&runnable{args: args, callback: transferReceivedCB})
	} else if args.eventType == chain.TransferEnd {
		f.addToRunCache(&runnable{args: args, callback: transferEndCB})
	} else {
		err = errors.New("args.EventType not among supported EventLog Types ")
	}
	return
}

func (f *finder) Match(elinfo chain.EventLogInfo) bool {
	if matchedIDs := f.lookupCache(elinfo); len(matchedIDs) > 0 {
		f.removeFromFromRunCache(matchedIDs)
		return true
	}
	return false
}

func (f *finder) addToRunCache(r *runnable) {
	f.runCache.mtx.Lock()
	defer f.runCache.mtx.Unlock()
	for i, v := range f.runCache.mem {
		if v == nil { // fill void if any
			f.runCache.mem[i] = r
			return
		}
	}
	f.runCache.mem = append(f.runCache.mem, r)
}

func (f *finder) removeFromFromRunCache(ids []int) {
	f.runCache.mtx.Lock()
	defer f.runCache.mtx.Unlock()
	for _, id := range ids {
		//f.log.Warnf("Removing %d", id)
		f.runCache.mem[id] = nil
	}
}

func (f *finder) lookupCache(elInfo chain.EventLogInfo) []int {
	f.runCache.mtx.RLock()
	defer f.runCache.mtx.RUnlock()
	matchedIDs := []int{}
	for runid, runP := range f.runCache.mem {
		if runP == nil { // nil is set for removed runnable. See removeFromFromRunCache
			continue
		}
		match, err := runP.callback(runP.args, elInfo)
		if match {
			//f.log.Warn("Match RunID ", runid)
			matchedIDs = append(matchedIDs, runid)
		} else if !match && err != nil {
			//f.log.Error("Non Match ", err)
		}

	}
	return matchedIDs
}

func NewFinder(l log.Logger) *finder {
	return &finder{
		log:      l,
		runCache: &runnableCache{mem: []*runnable{}, mtx: sync.RWMutex{}},
	}
}

var transferStartCB callBackFunc = func(args args, elInfo chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferStartEvent)
	if !ok {
		return false, errors.New("Should be TransferStart event log")
	}
	if elInfo.EventType == chain.TransferStart &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var transferReceivedCB callBackFunc = func(args args, elInfo chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferReceivedEvent)
	if !ok {
		return false, errors.New("Should be transferReceivedCB event log")
	}
	if elInfo.EventType == chain.TransferReceived &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var transferEndCB callBackFunc = func(args args, elInfo chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferEndEvent)
	if !ok {
		return false, errors.New("Should be transferEndCB event log")
	}
	if elInfo.EventType == chain.TransferEnd &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}
