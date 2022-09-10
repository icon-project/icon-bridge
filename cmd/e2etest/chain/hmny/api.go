//go:build hmny
// +build hmny

package hmny

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func NewApi(l log.Logger, cfg *chain.Config) (chain.ChainAPI, error) {
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
		log:      l,
		fd:       NewFinder(l, cfg.ContractAddresses),
		sinkChan: make(chan *chain.EventLogInfo),
		errChan:  make(chan error),
	}

	r.par, err = NewParser(cfg.URL, cfg.ContractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "newParser ")
	}
	r.requester, err = newRequestAPI(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "newRequestAPI %v", err)
	}
	return r, err
}

type api struct {
	*hmny.ReceiverCore
	log       log.Logger
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	par       *parser
	fd        *finder
	requester *requestAPI
}

func (a *api) client() *hmny.Client {
	return a.ReceiverCore.Cls[0]
}

// Options for a new block notifications channel

func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	height, err := a.client().GetBlockNumber()
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetBlockNumber ")
	}
	a.log.Infof("Subscribe Start Height %v", height)
	go func() {
		lastHeight := height - 1
		if err := a.ReceiverCore.ReceiveLoop(ctx,
			&hmny.BnOptions{
				StartHeight: height,
				Concurrency: monitorBlockMaxConcurrency,
			},
			func(v *hmny.BlockNotification) error {
				if v.Height.Int64()%100 == 0 {
					a.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")
				}
				if v.Height.Uint64() != lastHeight+1 {
					a.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, v.Height.Uint64())
					return fmt.Errorf(
						"block notification: expected=%d, got=%d",
						lastHeight+1, v.Height.Uint64())
				}
				if len(v.Receipts) > 0 {
					for _, sev := range v.Receipts {
						for _, txnLog := range sev.Logs {
							res, evtType, err := a.par.Parse(txnLog)
							if err != nil {
								a.log.Trace(errors.Wrap(err, "Parse "))
								err = nil
								continue
							}
							nel := &chain.EventLogInfo{ContractAddress: txnLog.Address.String(), EventType: evtType, EventLog: res}
							//r.Log.Infof("HFirst %+v", nel)
							//r.Log.Infof("HSecond %+v", nel.EventLog)
							if a.fd.Match(nel) {
								//a.log.Infof("Matched %+v", el)
								a.sinkChan <- nel
							}
						}
					}
				}
				lastHeight++
				return nil
			}); err != nil {
			a.log.Errorf("receiveLoop terminated: %+v", err)
			a.errChan <- err
		}
	}()

	return a.sinkChan, a.errChan, nil
}

func (a *api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".hmny") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	return a.requester.getCoinBalance(coinName, address)
}

func (a *api) Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
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
		txnHash, err = a.requester.transferIntraChain(coinName, senderKey, recepientAddress, amount)
	} else {
		txnHash, err = a.requester.transferInterChain(coinName, senderKey, recepientAddress, amount)
	}
	return
}

func (a *api) TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
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
		txnHash, err = a.requester.transferBatch(coinNames, senderKey, recepientAddress, amounts)
	}
	return
}

func (a *api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.approveCoin(coinName, ownerKey, amount)
	return
}
func (a *api) Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.reclaim(coinName, ownerKey, amount)
	return
}

func (a *api) WaitForTxnResult(ctx context.Context, hash string) (*chain.TxnResult, error) {
	txRes, err := a.requester.waitForResults(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "waitForResults ")
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.Logs {
		decodedLog, eventType, err := a.par.ParseEth(v)
		if err != nil {
			a.log.Trace(errors.Wrap(err, "ParseEth "))
			err = nil
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: v.Address.String(), EventType: eventType, EventLog: decodedLog})
	}
	return &chain.TxnResult{StatusCode: int(txRes.Status), ElInfo: plogs, Raw: txRes}, nil
}

func (a *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + a.requester.networkID + ".hmny/" + addr
	return fullAddr
}

func (a *api) NativeCoin() string {
	return a.requester.nativeCoin
}

func (a *api) NativeTokens() []string {
	return a.requester.nativeTokens
}

func (a *api) WatchForTransferStart(id uint64, seq int64) error {
	return a.fd.watchFor(chain.TransferStart, id, seq)
}

func (a *api) WatchForTransferReceived(id uint64, seq int64) error {
	return a.fd.watchFor(chain.TransferReceived, id, seq)
}

func (a *api) WatchForTransferEnd(id uint64, seq int64) error {
	return a.fd.watchFor(chain.TransferEnd, id, seq)
}

func (a *api) WatchForAddToBlacklistRequest(ID uint64, seq int64) error {
	return errors.New("not implemented")
}
func (a *api) WatchForRemoveFromBlacklistRequest(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) WatchForSetTokenLmitRequest(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) WatchForBlacklistResponse(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) WatchForSetTokenLmitResponse(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) GetConfigRequestEvent(evtType chain.EventLogType, hash string) (*chain.EventLogInfo, error) {
	return nil, errors.New("not implemented")
}

func (a *api) GetKeyPairFromKeystore(walFile string, password string) (privKey, pubKey string, err error) {
	keyReader, err := os.Open(walFile)
	if err != nil {
		err = errors.Wrapf(err, "os.Open file %v", walFile)
		return
	}
	defer keyReader.Close()

	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		err = errors.Wrapf(err, "ioutil.ReadAll %v", walFile)
		return
	}
	key, err := keystore.DecryptKey(keyStore, password)
	if err != nil {
		err = errors.Wrapf(err, "keystore.Decrypt %v", walFile)
		return
	}
	privBytes := crypto.FromECDSA(key.PrivateKey)
	privKey = hex.EncodeToString(privBytes)
	pubKey = crypto.PubkeyToAddress(key.PrivateKey.PublicKey).String()
	return
}

func (a *api) GetKeyPairs(num int) ([][2]string, error) {
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

func (a *api) SetTokenLimit(ownerKey string, coinNames []string, tokenLimits []*big.Int) (txnHash string, err error) {
	return "", errors.New("not implemented")
}
func (a *api) AddBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	return "", errors.New("not implemented")
}

func (a *api) RemoveBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	return "", errors.New("not implemented")
}

func (a *api) ChangeRestriction(ownerKey string, enable bool) (txnHash string, err error) {
	return "", errors.New("not implemented")
}

func (a *api) IsUserBlackListed(net, addr string) (response bool, err error) {
	return false, errors.New("not implemented")
}

func (a *api) GetTokenLimit(coinName string) (tokenLimit *big.Int, err error) {
	return nil, errors.New("not implemented")
}

func (a *api) IsBTSOwner(addr string) (response bool, err error) {
	return false, errors.New("not implemented")
}

func (a *api) GetTokenLimitStatus(net, coinName string) (response bool, err error) {
	return false, errors.New("not implemented")
}

func (a *api) GetBlackListedUsers(net string, startCursor, endCursor int) (addrs []string, err error) {
	return nil, errors.New("not implemented")
}

func (a *api) ChargedGasFee(txnHash string) (*big.Int, error) {
	return nil, errors.New("not implemented")
}

func (a *api) SuggestGasPrice() *big.Int {
	return nil
}
