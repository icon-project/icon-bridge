package bsc

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
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
					r.log.Info("Height %v", v.Height)
					for _, txnLog := range v.Logs {
						res, evtType, err := r.par.Parse(&txnLog)
						if err != nil {
							r.log.Warn(errors.Wrap(err, "Parse "))
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

func (r *api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".bsc") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain bsc. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	return r.requester.getCoinBalance(coinName, address)
}

func (r *api) Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
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
		txnHash, err = r.requester.transferIntraChain(coinName, senderKey, recepientAddress, amount)
	} else {
		txnHash, err = r.requester.transferInterChain(coinName, senderKey, recepientAddress, amount)
	}
	return
}

func (r *api) TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
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
		txnHash, err = r.requester.transferBatch(coinNames, senderKey, recepientAddress, amounts)
	}
	return
}

func (r *api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = r.requester.approveCoin(coinName, ownerKey, amount)
	return
}

func (a *api) Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.reclaim(coinName, ownerKey, amount)
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
	fullAddr := "btp://" + r.requester.networkID + ".bsc/" + addr
	return fullAddr
}

func (r *api) NativeCoin() string {
	return r.requester.nativeCoin
}

func (r *api) NativeTokens() []string {
	return r.requester.nativeTokens
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

func (r *api) GetKeyPairFromKeystore(keystoreFile string, secretFile string) (privKey, pubKey string, err error) {
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
