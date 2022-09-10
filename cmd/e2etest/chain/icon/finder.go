package icon

import (
	"fmt"
	"sync"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

type finder struct {
	log           log.Logger
	runCache      *runnableCache
	nameToAddrMap map[chain.ContractName]string
}

type args struct {
	id              uint64
	eventType       chain.EventLogType
	seq             int64
	contractAddress string
}

type callBackFunc func(args args, info *chain.EventLogInfo) (bool, error)

type runnable struct {
	args     args
	callback callBackFunc
}

type runnableCache struct {
	mem []*runnable
	mtx sync.RWMutex
}

func (f *finder) watchFor(eventType chain.EventLogType, id uint64, seq int64) error {
	btsAddr, ok := f.nameToAddrMap[chain.BTS]
	if !ok {
		return fmt.Errorf("watchFor; Contract %v not found on map", chain.BTS)
	}
	bmcAddr, ok := f.nameToAddrMap[chain.BMC]
	if !ok {
		return fmt.Errorf("watchFor; Contract %v not found on map", chain.BTS)
	}
	if eventType == chain.TransferStart {
		args := args{id: id, eventType: chain.TransferStart, seq: seq, contractAddress: btsAddr}
		f.addToRunCache(&runnable{args: args, callback: transferStartCB})
	} else if eventType == chain.TransferReceived {
		args := args{id: id, eventType: chain.TransferReceived, seq: seq, contractAddress: btsAddr}
		f.addToRunCache(&runnable{args: args, callback: transferReceivedCB})
	} else if eventType == chain.TransferEnd {
		args := args{id: id, eventType: chain.TransferEnd, seq: seq, contractAddress: btsAddr}
		f.addToRunCache(&runnable{args: args, callback: transferEndCB})
	} else if eventType == chain.AddToBlacklistRequest {
		args := args{id: id, eventType: chain.AddToBlacklistRequest, seq: seq, contractAddress: bmcAddr}
		f.addToRunCache(&runnable{args: args, callback: addedToBlacklistCB})
	} else if eventType == chain.RemoveFromBlacklistRequest {
		args := args{id: id, eventType: chain.RemoveFromBlacklistRequest, seq: seq, contractAddress: bmcAddr}
		f.addToRunCache(&runnable{args: args, callback: removedFromBlacklistCB})
	} else if eventType == chain.TokenLimitRequest {
		args := args{id: id, eventType: chain.TokenLimitRequest, seq: seq, contractAddress: bmcAddr}
		f.addToRunCache(&runnable{args: args, callback: tokenLimitSetCB})
	} else {
		return fmt.Errorf("watchFor; EventType %v not among supported ones", eventType)
	}
	return nil
}

func (f *finder) Match(elinfo *chain.EventLogInfo) bool {
	if matchedIndex, matchedIDs := f.lookupCache(elinfo); matchedIndex >= 0 {
		elinfo.ID = matchedIDs
		f.removeFromFromRunCache(matchedIndex)
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

func (f *finder) removeFromFromRunCache(id int) {
	f.runCache.mtx.Lock()
	defer f.runCache.mtx.Unlock()
	//f.log.Tracef("Removing %d", id)
	f.runCache.mem[id] = nil
}

func (f *finder) lookupCache(elInfo *chain.EventLogInfo) (int, uint64) {
	f.runCache.mtx.RLock()
	defer f.runCache.mtx.RUnlock()
	var matchedIndex int = -1
	var matchedIDs uint64
	for runid, runP := range f.runCache.mem {
		if runP == nil { // nil is set for removed runnable. See removeFromFromRunCache
			continue
		}
		match, err := runP.callback(runP.args, elInfo)
		if match {
			//f.log.Warn("Match RunID ", runid)
			matchedIndex = runid
			matchedIDs = runP.args.id
			break
		} else if !match && err != nil {
			//f.log.Error("Non Match ", err)
		}

	}
	return matchedIndex, matchedIDs
}

func NewFinder(l log.Logger, nameToAddrMap map[chain.ContractName]string) *finder {
	return &finder{
		log:           l,
		runCache:      &runnableCache{mem: []*runnable{}, mtx: sync.RWMutex{}},
		nameToAddrMap: nameToAddrMap,
	}
}

var transferStartCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferStartEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TransferStartEvent", elInfo.EventLog)
	}
	if elInfo.EventType == chain.TransferStart &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var transferReceivedCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferReceivedEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TransferReceivedEvent", elInfo.EventLog)
	}
	if elInfo.EventType == chain.TransferReceived &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var transferEndCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferEndEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TransferEndEvent", elInfo.EventLog)
	}
	if elInfo.EventType == chain.TransferEnd &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var addedToBlacklistCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.AddToBlacklistRequestEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *AddToBlacklistRequestEvent", elInfo.EventLog)
	}
	if elInfo.EventType == chain.AddToBlacklistRequest &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var removedFromBlacklistCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.RemoveFromBlacklistRequestEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *RemoveFromBlacklistRequestEvent", elInfo.EventLog)
	}
	if elInfo.EventType == chain.RemoveFromBlacklistRequest &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}

var tokenLimitSetCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TokenLimitRequestEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TokenLimitRequestEvent", elInfo.EventLog)
	}
	if elInfo.EventType == chain.TokenLimitRequest &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == args.seq {
		return true, nil
	}
	return false, nil
}
