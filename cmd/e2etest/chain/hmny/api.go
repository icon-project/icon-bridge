package hmny

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/hmny"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	BlockInterval              = 2 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	defaultReadTimeout         = 15 * time.Second
	monitorBlockMaxConcurrency = 50 // number of concurrent requests to synchronize older blocks from source chain
)

// const (
// 	NativeCoinName = "ONE"
// 	TokenName      = "TONE"
// )

func NewApi(l log.Logger, cfg *chain.ChainConfig) (chain.ChainAPI, error) {
	if len(cfg.URL) == 0 {
		return nil, errors.New("empty urls")
	}
	cls, err := hmny.NewClients([]string{cfg.URL}, l)
	if err != nil {
		return nil, errors.Wrap(err, "newCient ")
	}

	r := &api{
		ReceiverCore: &hmny.ReceiverCore{
			Log: l, Opts: hmny.ReceiverOptions{}, Cls: cls,
		},
		log:                l,
		networkID:          cfg.NetworkID,
		fd:                 NewFinder(l, cfg.ContractAddresses),
		sinkChan:           make(chan *chain.EventLogInfo),
		errChan:            make(chan error),
		nativeCoin:         cfg.NativeCoin,
		tokenNameToAddr:    cfg.NativeTokenAddresses,
		contractNameToAddr: cfg.ContractAddresses,
	}

	r.par, err = NewParser(cfg.URL, cfg.ContractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "newParser ")
	}
	r.requester, err = newRequestAPI(cfg.URL, l, cfg.ContractAddresses, cfg.NetworkID, r.tokenNameToAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "newRequestAPI %v", err)
	}
	return r, nil
}

type api struct {
	*hmny.ReceiverCore
	log                log.Logger
	sinkChan           chan *chain.EventLogInfo
	errChan            chan error
	par                *parser
	fd                 *finder
	requester          *requestAPI
	networkID          string
	nativeCoin         string
	tokenNameToAddr    map[string]string
	contractNameToAddr map[chain.ContractName]string
}

func (r *api) client() *hmny.Client {
	return r.ReceiverCore.Cls[0]
}

// Options for a new block notifications channel

func (r *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	height, err := r.client().GetBlockNumber()
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetBlockNumber ")
	}
	r.log.Infof("Subscribe Start Height %v", height)
	go func() {
		lastHeight := height - 1
		if err := r.ReceiverCore.ReceiveLoop(ctx,
			&hmny.BnOptions{
				StartHeight: height,
				Concurrency: monitorBlockMaxConcurrency,
			},
			func(v *hmny.BlockNotification) error {
				if v.Height.Int64()%100 == 0 {
					r.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")
				}
				if v.Height.Uint64() != lastHeight+1 {
					r.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, v.Height.Uint64())
					return fmt.Errorf(
						"block notification: expected=%d, got=%d",
						lastHeight+1, v.Height.Uint64())
				}
				if len(v.Receipts) > 0 {
					for _, sev := range v.Receipts {
						for _, txnLog := range sev.Logs {
							res, evtType, err := r.par.Parse(txnLog)
							if err != nil {
								r.log.Trace(errors.Wrap(err, "Parse "))
								err = nil
								continue
							}
							nel := &chain.EventLogInfo{ContractAddress: txnLog.Address.String(), EventType: evtType, EventLog: res}
							//r.Log.Infof("HFirst %+v", nel)
							//r.Log.Infof("HSecond %+v", nel.EventLog)
							if r.fd.Match(nel) {
								//r.log.Infof("Matched %+v", el)
								r.sinkChan <- nel
							}
						}
					}
				}
				lastHeight++
				return nil
			}); err != nil {
			r.log.Errorf("receiveLoop terminated: %+v", err)
			r.errChan <- err
		}
	}()

	return r.sinkChan, r.errChan, nil
}

func (r *api) GetChainType() chain.ChainType {
	return chain.HMNY
}

func (r *api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".hmny") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	return r.requester.getCoinBalance(coinName, address)
}

func (r *api) Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp://") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".hmny") {
		within = true
		splts := strings.Split(recepientAddress, "/")
		recepientAddress = splts[len(splts)-1]
	}
	if within {
		if coinName == r.nativeCoin {
			txnHash, err = r.requester.transferNativeIntraChain(senderKey, recepientAddress, amount)
		} else if _, ok := r.tokenNameToAddr[coinName]; ok {
			txnHash, err = r.requester.transferTokenIntraChain(senderKey, recepientAddress, amount, coinName)
		} else {
			err = fmt.Errorf("IntraChain transfers are supported for coins ONE and TONE only")
		}
	} else {
		if coinName == r.nativeCoin {
			txnHash, err = r.requester.transferNativeCrossChain(senderKey, recepientAddress, amount)
		} else { // TONE,ICX.TICX
			txnHash, err = r.requester.transferTokensCrossChain(coinName, senderKey, recepientAddress, amount)
		}
	}

	return
}

func (r *api) TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp:") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".hmny") {
		within = true
		splts := strings.Split(recepientAddress, "/")
		recepientAddress = splts[len(splts)-1]
	}
	if within {
		err = fmt.Errorf("Batch Transfers are supported for inter chain transfers only")
	} else {
		txnHash, err = r.requester.transferBatch(coinNames, senderKey, recepientAddress, amounts, r.nativeCoin)
	}
	return
}

func (r *api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = r.requester.approveCoin(coinName, ownerKey, amount)
	return
}

func (r *api) WaitForTxnResult(ctx context.Context, hash string) (*chain.TxnResult, error) {
	txRes, err := r.requester.waitForResults(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "waitForResults ")
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.Logs {
		decodedLog, eventType, err := r.par.ParseEth(v)
		if err != nil {
			r.log.Trace(errors.Wrap(err, "ParseEth "))
			err = nil
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: v.Address.String(), EventType: eventType, EventLog: decodedLog})
	}
	return &chain.TxnResult{StatusCode: int(txRes.Status), ElInfo: plogs, Raw: txRes}, nil
}

func (r *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + r.networkID + ".hmny/" + addr
	return fullAddr
}

func (r *api) NativeCoin() string {
	return r.nativeCoin
}

func (r *api) NativeTokens() []string {
	nativeTokens := []string{}
	for nt := range r.tokenNameToAddr {
		nativeTokens = append(nativeTokens, nt)
	}
	return nativeTokens
}

func (r *api) GetKeyPairs(num int) ([][2]string, error) {
	var err error
	res := make([][2]string, num)
	for i := 0; i < num; i++ {
		res[i], err = generateKeyPair()
		if err != nil {
			return nil, errors.Wrap(err, "generateKeyPair ")
		}
	}
	return res, nil
}

func (r *api) WatchForTransferStart(id uint64, seq int64) error {
	return r.fd.watchFor(chain.TransferStart, id, seq)
}

func (r *api) WatchForTransferReceived(id uint64, seq int64) error {
	return r.fd.watchFor(chain.TransferReceived, id, seq)
}

func (r *api) WatchForTransferEnd(id uint64, seq int64) error {
	return r.fd.watchFor(chain.TransferEnd, id, seq)
}

func (r *api) GetBTPAddressOfBTS() (btpaddr string, err error) {
	addr, ok := r.contractNameToAddr[chain.BTSCoreHmny]
	if !ok {
		err = fmt.Errorf("Contract %v does not exist ", chain.BTSCoreHmny)
		return
	}
	btpaddr = r.GetBTPAddress(addr)
	return
}

func (r *api) GetPubKey(privkey string) (string, error) {
	w, _, err := GetWalletFromPrivKey(privkey)
	if err != nil {
		return "", errors.Wrapf(err, "GetWalletFromPrivKey %v", err)
	}
	pubKey := w.PublicKey()
	return string(pubKey), nil
}
