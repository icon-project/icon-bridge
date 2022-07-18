package bsc

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewApi(l log.Logger, cfg *chain.ChainConfig) (chain.ChainAPI, error) {
	if len(cfg.URL) == 0 {
		return nil, errors.New("empty urls")
	}
	clrpc, err := rpc.Dial(cfg.URL)
	if err != nil {
		l.Errorf("failed to create bsc rpc client: url=%v, %v", cfg.URL, err)
		return nil, err
	}
	r := &api{
		log:             l,
		fd:              NewFinder(l, cfg.ContractAddresses),
		sinkChan:        make(chan *chain.EventLogInfo),
		errChan:         make(chan error),
		networkID:       cfg.NetworkID,
		nativeCoin:      cfg.NativeCoin,
		tokenNameToAddr: cfg.NativeTokenAddresses,
		ReceiverCore: &ReceiverCore{
			Log:  l,
			Opts: ReceiverOptions{},
			Cls:  []*ethclient.Client{ethclient.NewClient(clrpc)},
			BlockReq: ethereum.FilterQuery{
				Addresses: []ethCommon.Address{
					ethCommon.HexToAddress(cfg.ContractAddresses[chain.BTSPeripheryBsc]),
				},
			},
		},
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
	*ReceiverCore
	log             log.Logger
	par             *parser
	requester       *requestAPI
	networkID       string
	fd              *finder
	sinkChan        chan *chain.EventLogInfo
	errChan         chan error
	nativeCoin      string
	tokenNameToAddr map[string]string
}

func (r *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	height, err := r.client().BlockNumber(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetBlockNumber ")
	}
	r.log.Infof("Subscribe Start Height %v", height)
	go func() {
		lastHeight := height - 1
		if err := r.ReceiverCore.ReceiveLoop(ctx,
			&BnOptions{
				StartHeight: height,
				Concurrency: monitorBlockMaxConcurrency,
			},
			func(v *BlockNotification) error {
				if v.Height.Int64()%100 == 0 {
					r.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")
				}
				if v.Height.Uint64() != lastHeight+1 {
					r.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, v.Height.Uint64())
					return fmt.Errorf(
						"block notification: expected=%d, got=%d",
						lastHeight+1, v.Height.Uint64())
				}
				if len(v.Logs) > 0 {
					for _, txnLog := range v.Logs {
						res, evtType, err := r.par.Parse(&txnLog)
						if err != nil {
							r.log.Trace(errors.Wrap(err, "Parse "))
							err = nil
							continue
						}
						nel := &chain.EventLogInfo{ContractAddress: txnLog.Address.String(), EventType: evtType, EventLog: res}
						r.Log.Infof("BFirst  %+v", nel)
						r.Log.Infof("BSecond  %+v", nel.EventLog)
						if r.fd.Match(nel) {
							//r.log.Infof("Matched %+v", el)
							r.sinkChan <- nel
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
	return chain.BSC
}

func (r *api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".bsc") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain bsc. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	if coinName == r.nativeCoin {
		return r.requester.getNativeCoinBalance(coinName, address)
	}
	return r.requester.getCoinBalance(coinName, address)
}

func (r *api) Transfer(coinName, senderKey, recepientAddress string, amount big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp:") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".bsc") {
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
		} else { // TBNB,ICX.TICX
			txnHash, err = r.requester.transferTokensCrossChain(coinName, senderKey, recepientAddress, amount)
		}
	}

	return
}

func (r *api) Approve(coinName string, ownerKey string, amount big.Int) (txnHash string, err error) {
	txnHash, err = r.requester.approveCoin(coinName, ownerKey, amount)
	return
}

func (r *api) WaitForTxnResult(ctx context.Context, hash string) (*chain.TxnResult, error) {
	txRes, err := r.requester.waitForResults(ctx, ethCommon.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "waitForResults ")
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.Logs {
		decodedLog, eventType, err := r.par.Parse(v)
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
	fullAddr := "btp://" + r.networkID + ".bsc/" + addr
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

func exists(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
