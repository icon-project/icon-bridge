package watcher

import (
	"errors"

	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
)

func (w *watcher) decodeEventLog(evt *chain.SubscribedEvent) ([]eventLogInfo, error) {
	elInfoList := []eventLogInfo{}
	if evt.ChainName == chain.ICON {
		el, ok := evt.Res.([]*icon.TxnEventLog)
		if !ok {
			return nil, errors.New("Subscribed EventLog of wrong type; Expected *icon.TxnLog")
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
					sourceChain:  evt.ChainName,
					contractName: ctrName,
					eventType:    eventType,
					eventLog:     decLog,
				}
				elInfoList = append(elInfoList, elInfo)
			}
		}
	} else if evt.ChainName == chain.HMNY {
		el, ok := evt.Res.([]*types.Log)
		if !ok {
			return nil, errors.New("Subscribed EventLog of wrong type; Expected *types.Receipt")
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
					sourceChain:  evt.ChainName,
					contractName: ctrName,
					eventType:    eventType,
					eventLog:     decLog,
				}
				elInfoList = append(elInfoList, elInfo)
			}
		}
	} else {
		return nil, errors.New("Unknow Chain Type")
	}
	return elInfoList, nil
}
