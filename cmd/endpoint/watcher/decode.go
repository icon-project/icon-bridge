package watcher

import (
	"errors"

	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
)

func (w *watcher) decodeIconEventLog(res interface{}) ([]eventLogInfo, error) {
	elInfoList := []eventLogInfo{}
	el, ok := res.([]*icon.TxnEventLog)
	if !ok {
		w.log.Errorf("\nExpected []*icon.TxnEventLog; Got Type: %T\n", res)
		return nil, errors.New("Subscribed EventLog of wrong type; Expected []*icon.TxnEventLog")
	}
	for _, l := range el {
		decEvt, err := w.dec.DecodeEventLogData(l, string(l.Addr))
		if err != nil {
			return nil, err
		}
		ctrName, ok := w.ctrAddrToName[string(l.Addr)]
		if !ok {
			return nil, errors.New("Watcher's ctrAddrToName doesn't include address " + string(l.Addr))
		}
		for eventType, decLog := range decEvt {
			elInfo := eventLogInfo{
				sourceChain:  chain.ICON,
				contractName: ctrName,
				eventType:    eventType,
				eventLog:     decLog,
			}
			elInfoList = append(elInfoList, elInfo)
		}
	}
	return elInfoList, nil
}

func (w *watcher) decodeHmnyEventLog(res interface{}) ([]eventLogInfo, error) {
	elInfoList := []eventLogInfo{}
	el, ok := res.([]*types.Log)
	if !ok {
		w.log.Error("Expected []*types.Log; Got Type: %T", res)
		return nil, errors.New("Subscribed EventLog of wrong type; Expected []*types.Log")
	}
	for _, l := range el {
		decEvt, err := w.dec.DecodeEventLogData(l, l.Address.Hex())
		if err != nil {
			return nil, err
		}
		ctrName, ok := w.ctrAddrToName[l.Address.Hex()]
		if !ok {
			return nil, errors.New("Watcher's ctrAddrToName doesn't include address " + l.Address.Hex())
		}
		for eventType, decLog := range decEvt {
			elInfo := eventLogInfo{
				sourceChain:  chain.HMNY,
				contractName: ctrName,
				eventType:    eventType,
				eventLog:     decLog,
			}
			elInfoList = append(elInfoList, elInfo)
		}
	}
	return elInfoList, nil
}

func (w *watcher) decodeEventLog(evt *chain.SubscribedEvent) ([]eventLogInfo, error) {
	if evt.ChainName == chain.ICON {
		return w.decodeIconEventLog(evt.Res)
	} else if evt.ChainName == chain.HMNY {
		return w.decodeHmnyEventLog(evt.Res)
	}
	return nil, errors.New("Unknow Chain Type")
}
