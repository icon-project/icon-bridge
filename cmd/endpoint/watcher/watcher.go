package watcher

import (
	"context"
	"errors"
	"sync"

	"github.com/icon-project/icon-bridge/common/log"

	capi "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder"
	ctr "github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
)

var nameMap = map[string]ctr.ContractName{
	"btp_icon_token_bsh":                ctr.TokenIcon,
	"btp_icon_nativecoin_bsh":           ctr.NativeIcon,
	"btp_icon_irc2":                     ctr.Irc2Icon,
	"btp_icon_irc2_tradeable":           ctr.Irc2TradeableIcon,
	"btp_icon_bmc":                      ctr.BmcIcon,
	"btp_hmny_token_bsh_impl":           ctr.TokenHmy,
	"btp_hmny_nativecoin_bsh_periphery": ctr.NativeHmy,
	"btp_hmny_erc20":                    ctr.Erc20Hmy,
	"btp_hmny_erc20_tradeable":          ctr.Erc20TradeableHmy,
	"btp_hmny_bmc_periphery":            ctr.BmcHmy,
	"btp_hmny_nativecoin_bsh_core":      ctr.OwnerNativeHmy,
	"btp_hmny_token_bsh_proxy":          ctr.OwnerTokenHmy,
}

type Watcher interface {
	Start(ctx context.Context) error
	ProcessTxn(reqParam *capi.RequestParam, logs interface{}) error
	GetOutputChan() <-chan eventLogInfo
}

type watcher struct {
	log                log.Logger
	ctrAddrToName      map[string]ctr.ContractName
	dec                decoder.Decoder
	inputSubChan       <-chan *chain.SubscribedEvent
	inputErrChan       <-chan error
	outputEventLogChan chan eventLogInfo
	idGroups           []identifierGroup
	runCache           *runnableCache
}

func New(log log.Logger, cfgPerChain map[chain.ChainType]*chain.ChainConfig, subChan <-chan *chain.SubscribedEvent, errChan <-chan error) (Watcher, error) {
	resMap := map[string]ctr.ContractName{}
	endpointPerChain := map[chain.ChainType]string{}
	for name, cfg := range cfgPerChain {
		endpointPerChain[name] = cfg.URL
		for cName, cAddr := range cfg.ConftractAddresses {
			if v, ok := nameMap[cName]; ok {
				resMap[cAddr] = v
			} else {
				log.Errorf("Contract isn't mentioned under watchlist %v", cName)
			}
		}
	}
	dec, err := decoder.New(endpointPerChain, resMap)
	if err != nil {
		return nil, err
	}
	w := &watcher{
		log:          log,
		inputSubChan: subChan, inputErrChan: errChan,
		ctrAddrToName: resMap, dec: dec,
		idGroups: DefaultIdentifierGroup, runCache: &runnableCache{mem: []*runnable{}, mtx: sync.RWMutex{}},
		outputEventLogChan: make(chan eventLogInfo),
	}
	return w, nil
}

func (w *watcher) Start(ctx context.Context) error {
	go func() {
		defer close(w.outputEventLogChan)
		for {
			select {
			case <-ctx.Done():
				w.log.Warn("Watcher; Context Cancelled")
				return
			case msg := <-w.inputSubChan:
				if decLogs, err := w.decodeSubscribedMessage(msg); err != nil {
					w.log.Error(err)
				} else {
					for _, dl := range decLogs {
						if matchedIDs := w.lookupCache(dl); len(matchedIDs) > 0 {
							w.outputEventLogChan <- dl
							w.removeFromFromRunCache(matchedIDs)
						}
					}
				}
			case err := <-w.inputErrChan:
				w.log.Error(err)
				return
			}
		}
	}()
	return nil
}

func (w *watcher) ProcessTxn(reqParam *capi.RequestParam, logs interface{}) error {
	elArr, err := w.decodeTransferReceipt(reqParam, logs)
	if err != nil {
		return err
	}
	for _, idg := range w.idGroups {
		res, ok := idg.init(elArr, reqParam)
		if !ok {
			if res != nil {
				return res.(error)
			}
			return errors.New("Did not match")
		}
		args := args{req: reqParam, initRes: res}
		for _, idf := range idg.idfs {
			w.addToRunCache(&runnable{args: args, idf: idf})
		}
		w.log.Warnf("Added to CacheLen %d; Seq %v ", len(w.runCache.mem), res)
	}
	return nil
}

type runnable struct {
	args args
	idf  identifier
}

type runnableCache struct {
	mem []*runnable
	mtx sync.RWMutex
}

func (w *watcher) addToRunCache(r *runnable) {
	w.runCache.mtx.Lock()
	defer w.runCache.mtx.Unlock()
	for i, v := range w.runCache.mem {
		if v == nil { // fill void if any
			w.runCache.mem[i] = r
			return
		}
	}
	w.runCache.mem = append(w.runCache.mem, r)
}

func (w *watcher) removeFromFromRunCache(ids []int) {
	w.runCache.mtx.Lock()
	defer w.runCache.mtx.Unlock()
	for _, id := range ids {
		w.log.Warnf("Removing %d", id)
		w.runCache.mem[id] = nil
	}
}

func (w *watcher) lookupCache(elInfo eventLogInfo) []int {
	w.runCache.mtx.RLock()
	defer w.runCache.mtx.RUnlock()
	matchedIDs := []int{}
	for runid, runP := range w.runCache.mem {
		if runP == nil { // nil is set for removed runnable. See removeFromFromRunCache
			continue
		}
		if runP.idf.preRun(runP.args, elInfo) {
			match, err := runP.idf.run(runP.args, elInfo)
			if match {
				w.log.Warn("Match RunID ", runid)
				matchedIDs = append(matchedIDs, runid)
			} else if !match && err != nil {
				w.log.Error("Non Match ", err)
			}
		}
	}
	return matchedIDs
}

func (w *watcher) display(el eventLogInfo) {
	w.log.Info("Matched ")
	w.log.Info("Prelim ", el.sourceChain, el.contractName, el.eventType)
	w.log.Infof("Extra Info %+v", el.eventLog)
}

func (w *watcher) GetOutputChan() <-chan eventLogInfo {
	return w.outputEventLogChan
}
