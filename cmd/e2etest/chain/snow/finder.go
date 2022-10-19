package snow

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
	data            interface{}
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

func (f *finder) watchFor(eventType chain.EventLogType, id uint64, data interface{}) error {
	btsAddress, ok := f.nameToAddrMap[chain.BTSPeriphery]
	if !ok {
		return fmt.Errorf("watchFor; Contract %v not found on map", chain.BTSPeriphery)
	}
	bmcAddress, ok := f.nameToAddrMap[chain.BMCPeriphery]
	if !ok {
		return fmt.Errorf("watchFor; Contract %v not found on map", chain.BTSPeriphery)
	}
	if eventType == chain.TransferStart {
		args := args{id: id, eventType: chain.TransferStart, data: data, contractAddress: btsAddress}
		f.addToRunCache(&runnable{args: args, callback: transferStartCB})
	} else if eventType == chain.TransferReceived {
		args := args{id: id, eventType: chain.TransferReceived, data: data, contractAddress: btsAddress}
		f.addToRunCache(&runnable{args: args, callback: transferReceivedCB})
	} else if eventType == chain.TransferEnd {
		args := args{id: id, eventType: chain.TransferEnd, data: data, contractAddress: btsAddress}
		f.addToRunCache(&runnable{args: args, callback: transferEndCB})
	} else if eventType == chain.BlacklistResponse {
		args := args{id: id, eventType: chain.BlacklistResponse, data: data, contractAddress: bmcAddress}
		f.addToRunCache(&runnable{args: args, callback: blackListResponseCB})
	} else if eventType == chain.TokenLimitResponse {
		args := args{id: id, eventType: chain.TokenLimitResponse, data: data, contractAddress: bmcAddress}
		f.addToRunCache(&runnable{args: args, callback: tokenLimitResponseCB})
	} else {
		return fmt.Errorf("EventType %v not among supported ones", eventType)
	}

	return nil
}

func (f *finder) Match(elinfo *chain.EventLogInfo) bool {
	if matchedIndex, matchedIDs := f.lookupCache(elinfo); matchedIndex >= 0 {
		elinfo.PID = matchedIDs
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
	argsdataNum, ok := args.data.(int64)
	if ok {
		if elInfo.EventType == chain.TransferStart &&
			elInfo.ContractAddress == args.contractAddress &&
			elog.Sn.Int64() == argsdataNum {
			return true, nil
		}
		return false, nil
	}
	argsdataArr, ok := args.data.([2]string) // FeeAgg
	if ok {
		if elInfo.EventType == chain.TransferStart &&
			elInfo.ContractAddress == args.contractAddress &&
			elog.From == argsdataArr[0] && // BTSCoreAddress
			elog.To == argsdataArr[1] { // FeeAggBTPAddress
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("Unexpected Type %T", args.data)
}

var transferReceivedCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferReceivedEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TransferReceivedEvent", elInfo.EventLog)
	}
	argsdata, ok := args.data.(int64)
	if !ok {
		return false, fmt.Errorf("Expected int64. Got %T", args.data)
	}
	if elInfo.EventType == chain.TransferReceived &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == argsdata {
		return true, nil
	}
	return false, nil
}

var transferEndCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TransferEndEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TransferEndEvent", elInfo.EventLog)
	}
	argsdata, ok := args.data.(int64)
	if !ok {
		return false, fmt.Errorf("Expected int64. Got %T", args.data)
	}
	if elInfo.EventType == chain.TransferEnd &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == argsdata {
		return true, nil
	}
	return false, nil
}

var blackListResponseCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.BlacklistResponseEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *BlacklistResponseEvent", elInfo.EventLog)
	}
	argsdata, ok := args.data.(int64)
	if !ok {
		return false, fmt.Errorf("Expected int64. Got %T", args.data)
	}
	if elInfo.EventType == chain.BlacklistResponse &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == argsdata {
		return true, nil
	}
	return false, nil
}

var tokenLimitResponseCB callBackFunc = func(args args, elInfo *chain.EventLogInfo) (bool, error) {
	elog, ok := (elInfo.EventLog).(*chain.TokenLimitResponseEvent)
	if !ok {
		return false, fmt.Errorf("Unexpected Type. Gor %T. Expected *TokenLimitResponseEvent", elInfo.EventLog)
	}
	argsdata, ok := args.data.(int64)
	if !ok {
		return false, fmt.Errorf("Expected int64. Got %T", args.data)
	}
	if elInfo.EventType == chain.TokenLimitResponse &&
		elInfo.ContractAddress == args.contractAddress &&
		elog.Sn.Int64() == argsdata {
		return true, nil
	}
	return false, nil
}
