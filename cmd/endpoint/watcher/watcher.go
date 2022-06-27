package watcher

import (
	"context"

	"github.com/icon-project/icon-bridge/common/log"

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
}

type watcher struct {
	log           log.Logger
	subChan       <-chan *chain.SubscribedEvent
	errChan       <-chan error
	ctrAddrToName map[string]ctr.ContractName
	dec           decoder.Decoder
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
	w := &watcher{log: log, subChan: subChan, errChan: errChan, ctrAddrToName: resMap, dec: dec}
	return w, nil
}

func (w *watcher) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.log.Warn("Watcher; Context Cancelled")
				return
			case msg := <-w.subChan:
				if decLogs, err := w.decodeEventLog(msg); err != nil {
					w.log.Error(err)
				} else {
					for dli, dl := range decLogs {
						w.log.Info(dli, dl)
					}
				}
			case err := <-w.errChan:
				w.log.Error(err)
				return
			}
		}
	}()
	return nil
}
