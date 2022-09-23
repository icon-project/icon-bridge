package bsc

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethCommon "github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewApi(l log.Logger, cfg *chain.Config) (chain.ChainAPI, error) {
	if len(cfg.URL) == 0 {
		return nil, errors.New("empty urls")
	}
	clrpc, err := rpc.Dial(cfg.URL)
	if err != nil {
		l.Errorf("failed to create bsc rpc client: url=%v, %v", cfg.URL, err)
		return nil, err
	}
	r := &api{
		log:      l,
		mu:       sync.Mutex{},
		fd:       NewFinder(l, cfg.ContractAddresses),
		sinkChan: make(chan *chain.EventLogInfo),
		errChan:  make(chan error),
		ReceiverCore: &ReceiverCore{
			Log:  l,
			Opts: ReceiverOptions{},
			Cls:  []*ethclient.Client{ethclient.NewClient(clrpc)},
			BlockReq: ethereum.FilterQuery{
				Addresses: []ethCommon.Address{
					ethCommon.HexToAddress(cfg.ContractAddresses[chain.BTSPeriphery]),
					ethCommon.HexToAddress(cfg.ContractAddresses[chain.BMCPeriphery]),
				},
			},
		},
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
	*ReceiverCore
	log       log.Logger
	par       *parser
	requester *requestAPI
	fd        *finder
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	mu        sync.Mutex
}

func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	height, err := a.client().BlockNumber(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetBlockNumber ")
	}
	a.log.Infof("Subscribe Start Height %v", height)
	go func() {
		lastHeight := height - 1
		if err := a.ReceiverCore.ReceiveLoop(ctx,
			&BnOptions{
				StartHeight: height,
				Concurrency: monitorBlockMaxConcurrency,
			},
			func(v *BlockNotification) error {
				if v.Height.Int64()%100 == 0 {
					a.log.WithFields(log.Fields{"height": v.Height}).Debug("block notification")
				}
				if v.Height.Uint64() != lastHeight+1 {
					a.log.Errorf("expected v.Height == %d, got %d", lastHeight+1, v.Height.Uint64())
					return fmt.Errorf(
						"block notification: expected=%d, got=%d",
						lastHeight+1, v.Height.Uint64())
				}
				if len(v.Logs) > 0 {
					a.log.Debugf("Height %v", v.Height)
					for _, txnLog := range v.Logs {
						res, evtType, err := a.par.Parse(&txnLog)
						if err != nil {
							//a.log.Trace(errors.Wrap(err, "Parse "))
							err = nil
							continue
						}
						nel := &chain.EventLogInfo{ContractAddress: txnLog.Address.String(), EventType: evtType, EventLog: res}
						a.log.Debugf("BFirst  %+v", nel)
						a.log.Debugf("BSecond  %+v", nel.EventLog)
						if a.fd.Match(nel) {
							//a.log.Infof("Matched %+v", el)
							a.sinkChan <- nel
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
	if !strings.Contains(addr, ".bsc") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain bsc. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	return a.requester.getCoinBalance(coinName, address)
}

func (a *api) Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp:") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".bsc") {
		within = true
		splts := strings.Split(recepientAddress, "/")
		recepientAddress = splts[len(splts)-1]
	}
	var nonce uint64
	if within {
		txnHash, nonce, err = a.requester.transferIntraChain(coinName, senderKey, recepientAddress, amount, nil)
		if err != nil && strings.Contains(err.Error(), "replacement transaction underpriced") {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Millisecond * time.Duration(100*rand.Intn(5))) // wait random times and retry
				newNonce := a.requester.getNonce(senderKey, nonce)
				txnHash, nonce, err = a.requester.transferIntraChain(coinName, senderKey, recepientAddress, amount, &newNonce)
				if err != nil && strings.Contains(err.Error(), "replacement transaction underpriced") {
					err = errors.Wrapf(err, "Retry %v Nonce %v Err %v", i, nonce, err)
					continue
				}
				break
			}
		}
	} else {
		txnHash, nonce, err = a.requester.transferInterChain(coinName, senderKey, recepientAddress, amount, nil)
		if err != nil && strings.Contains(err.Error(), "replacement transaction underpriced") {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Millisecond * time.Duration(100*rand.Intn(5))) // wait random times and retry
				newNonce := a.requester.getNonce(senderKey, nonce)
				txnHash, nonce, err = a.requester.transferInterChain(coinName, senderKey, recepientAddress, amount, &newNonce)
				if err != nil && strings.Contains(err.Error(), "replacement transaction underpriced") {
					err = errors.Wrapf(err, "Retry %v Nonce %v Err %v", i, nonce, err)
					continue
				}
				break
			}
		}
	}
	return
}

func (r *requestAPI) getNonce(addr string, input uint64) uint64 {
	r.nonceMapMutex.Lock()
	defer r.nonceMapMutex.Unlock()
	existing, ok := r.nonceMap[addr]
	if !ok || (ok && input > existing) {
		r.nonceMap[addr] = input + 1
		return input + 1
	}
	// ok && input < existing
	r.nonceMap[addr] = existing + 1
	return existing + 1
}

func (a *api) TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
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
		err = fmt.Errorf("Batch Transfers are supported for inter chain transfers only")
	} else {
		var nonce uint64
		txnHash, nonce, err = a.requester.transferBatch(coinNames, senderKey, recepientAddress, amounts, nil)
		if err != nil && strings.Contains(err.Error(), "replacement transaction underpriced") {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Millisecond * time.Duration(100*rand.Intn(5))) // wait random times and retry
				newNonce := a.requester.getNonce(senderKey, nonce)
				txnHash, nonce, err = a.requester.transferBatch(coinNames, senderKey, recepientAddress, amounts, &newNonce)
				if err != nil && strings.Contains(err.Error(), "replacement transaction underpriced") {
					err = errors.Wrapf(err, "Retry %v Nonce %v Err %v", i, nonce, err)
					continue
				}
				break
			}
		}
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
	txRes, err := a.requester.waitForResults(ctx, ethCommon.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrap(err, "waitForResults ")
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.Logs {
		decodedLog, eventType, err := a.par.Parse(v)
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
	fullAddr := "btp://" + a.requester.networkID + ".bsc/" + addr
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
	return a.fd.watchFor(chain.BlacklistResponse, ID, seq)
}

func (a *api) WatchForSetTokenLmitResponse(ID uint64, seq int64) error {
	return a.fd.watchFor(chain.TokenLimitResponse, ID, seq)
}

func (a *api) GetConfigRequestEvent(evtType chain.EventLogType, hash string) (*chain.EventLogInfo, error) {
	return nil, errors.New("not implemented")
}

func (a *api) SetFeeRatio(ownerKey string, coinName string, feeNumerator, fixedFee *big.Int) (string, error) {
	return a.requester.setFeeRatio(ownerKey, coinName, feeNumerator, fixedFee)
}

func (a *api) GetFeeRatio(coinName string) (feeNumerator *big.Int, fixedFee *big.Int, err error) {
	res, err := a.requester.btsc.FeeRatio(&bind.CallOpts{Pending: false, Context: context.Background()}, coinName)
	return res.FeeNumerator, res.FixedFee, err
}

func (a *api) GetAccumulatedFees() (map[string]*big.Int, error) {
	res, err := a.requester.btsc.GetAccumulatedFees(&bind.CallOpts{Pending: false, Context: context.Background()})
	if err != nil {
		return nil, err
	}
	resMap := map[string]*big.Int{}
	for _, v := range res {
		resMap[v.CoinName] = v.Value
	}
	return resMap, nil
}

func (a *api) SetFeeGatheringTerm(ownerKey string, interval uint64) (hash string, err error) {
	return "", errors.New("not implemented")
}

func (a *api) GetFeeGatheringTerm() (interval uint64, err error) {
	return 0, errors.New("not implemented")
}

func (a *api) GetKeyPairs(num int) ([][2]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var err error
	res := make([][2]string, num)
	for i := 0; i < num; i++ {
		res[i], err = generateKeyPair() // crypto library is not thread safe, so mutex locked at start
		if err != nil {
			return nil, errors.Wrap(err, "generateKeyPair ")
		}
	}
	return res, nil
}

func (a *api) GetKeyPairFromKeystore(keystoreFile string, secretFile string) (privKey, pubKey string, err error) {
	readFile := func(file string) (string, error) {
		f, err := os.Open(file)
		if err != nil {
			return "", err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}

	secret, err := readFile(secretFile)
	if err != nil {
		err = errors.Wrapf(err, "readPassFromFile(%v) %v", secretFile, err)
		return
	}
	wal, err := readFile(keystoreFile)
	if err != nil {
		err = errors.Wrapf(err, "readKeystoreFromFile(%v) %v", keystoreFile, err)
		return
	}

	key, err := keystore.DecryptKey([]byte(wal), secret)
	if err != nil {
		err = errors.Wrapf(err, "keystore.Decrypt %v", err)
		return
	}
	privBytes := crypto.FromECDSA(key.PrivateKey)
	privKey = hex.EncodeToString(privBytes)
	pubKey = crypto.PubkeyToAddress(key.PrivateKey.PublicKey).String()
	return
}

func (a *api) SuggestGasPrice() (gasPrice *big.Int) {
	ctx := context.TODO()
	cleth := a.client()

	gasPrice = big.NewInt(20000000000) // default Gas Price
	header, err := cleth.HeaderByNumber(ctx, nil)
	if err != nil {
		err = errors.Wrapf(err, "GetHeaderByNumber(height:latest) Err: %v", err)
		return
	}
	height := header.Number
	txnCount, err := cleth.TransactionCount(ctx, header.Hash())
	if err != nil {
		err = errors.Wrapf(err, "GetTransactionCount(height:%v, headerHash: %v) Err: %v", height, header.Hash(), err)
		a.log.Error(err)
		return
	} else if err == nil && txnCount == 0 {
		err = fmt.Errorf("TransactionCount is zero for height(%v, headerHash %v)", height, header.Hash())
		a.log.Error(err)
		return
	}

	txnS, err := cleth.TransactionInBlock(ctx, header.Hash(), uint(math.Floor(float64(txnCount)/2)))
	if err != nil {
		err = errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, txnCount-1, err)
		a.log.Error(err)
		return
	}
	gasPrice = (&big.Int{}).Mul(txnS.GasPrice(), big.NewInt(112))
	gasPrice = gasPrice.Div(gasPrice, big.NewInt(100))
	suggested, err := a.Cls[0].SuggestGasPrice(ctx)
	if err != nil && suggested.Cmp(gasPrice) > 0 {
		a.log.Debug("Using Suggested ", suggested, " instead of calculated ", gasPrice)
		gasPrice = suggested
	}
	if gasPrice.Int64() == 0 {
		a.log.Debug("Calculated Gas Price was zero++++++++")
		return big.NewInt(20000000000)
	}
	return
}

func (a *api) ChargedGasFee(txnHash string) (*big.Int, error) {
	txr, err := a.requester.ethCl.TransactionReceipt(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		return nil, errors.Wrapf(err, "TransactionByHash %v", err)
	}
	txh, _, err := a.requester.ethCl.TransactionByHash(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		return nil, errors.Wrapf(err, "TransactionByHash %v", err)
	}
	ret := (&big.Int{}).Mul(big.NewInt(int64(txr.GasUsed)), txh.GasPrice())
	return ret, nil
}

func (a *api) WatchForFeeGatheringRequest(ID uint64, addr string) error {
	return errors.New("not implemented")
}

func (a *api) WatchForFeeGatheringTransferStart(ID uint64, addr string) error {
	return a.fd.watchFor(chain.TransferStart, ID, [2]string{a.fd.nameToAddrMap[chain.BTS], addr})
}
