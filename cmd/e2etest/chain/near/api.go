package near

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"strings"

	gocommon "github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near"
	"github.com/icon-project/icon-bridge/common"

	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

type api struct {
	receiver  *near.Receiver
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	par       *parser
	fd        *finder
	requester *requestAPI
}

func NewApi(l log.Logger, cfg *chain.Config) (chain.ChainAPI, error) {
	var err error
	if len(cfg.URL) == 0 {
		return nil, errors.New("Expected URL for chain NEAR. Got ")
	} else if cfg.Name != chain.NEAR {
		return nil, fmt.Errorf("Expected cfg.Name=NEAR Got %v", cfg.Name)
	}
	Clients := near.NewClients([]string{cfg.URL}, l)

	recvr := &api{
		receiver: &near.Receiver{
			Clients: Clients,
			Logger:  l,
		},
		sinkChan: make(chan *chain.EventLogInfo),
		errChan:  make(chan error),
		fd:       NewFinder(l, cfg.ContractAddresses),
	}
	recvr.par, err = NewParser(cfg.ContractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "NewParser ")
	}
	client := Clients[rand.Intn(len(Clients))]
	recvr.requester, err = newRequestAPI(client, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "newRequestAPI %v", err)
	}
	return recvr, err
}

// Approve implements chain.ChainAPI
func (a *api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.approve(coinName, ownerKey, amount)
	return
}

// CallBTS implements chain.ChainAPI
func (api *api) CallBTS(method chain.ContractCallMethodName, args []interface{}) (response interface{}, err error) {
	// Tokens ..
	btsAddr, ok := api.requester.contractNameToAddress[chain.BTS]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}

	res, err := api.requester.callContract(btsAddr, map[string]interface{}{}, string(method))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetBTPAddress implements chain.ChainAPI
func (a *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + a.requester.networkID + ".near/" + addr
	return fullAddr
}

// GetCoinBalance implements chain.ChainAPI
func (a *api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".icon") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	return a.requester.getCoinBalance(coinName, address)
}

// GetKeyPairFromKeystore implements chain.ChainAPI
func (a *api) GetKeyPairFromKeystore(keystoreFile, secretFile string) (priv string, pub string, err error) {
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

	privKey, err := wallet.DecryptKeyStore([]byte(wal), []byte(secret))
	if err != nil {
		err = errors.Wrapf(err, "wallet.DecryptKeyStore %v", err)
	}
	priv = hex.EncodeToString(privKey.Bytes())
	pub = gocommon.NewAccountAddressFromPublicKey(privKey.PublicKey()).String()
	return
}

// GetKeyPairs implements chain.ChainAPI
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

// GetNetwork implements chain.ChainAPI
func (a *api) GetNetwork() string {
	return a.requester.networkID + ".near"
}

// NativeCoin implements chain.ChainAPI
func (a *api) NativeCoin() string {
	return a.requester.nativeCoin
}

// NativeTokens implements chain.ChainAPI
func (a *api) NativeTokens() []string {
	nativeTokens := []string{}
	for name := range a.requester.nativeTokensAddr {
		nativeTokens = append(nativeTokens, name)
	}
	return nativeTokens
}

// Reclaim implements chain.ChainAPI
func (a *api) Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.reclaim(coinName, ownerKey, amount)
	return
}

// func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
// 	return nil, nil, nil
// }

// Subscribe implements chain.ChainAPI
func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	height, err := a.requester.cl.CallLatestBlockHeight()
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetLastBlock ")
	}
	//height := uint64(blk.Height)
	a.receiver.Logger.Infof("Subscribe Start Height %v", height)
	// _errCh := make(chan error)
	go func() {
		// defer close(_errCh)
		err := a.ReceiveLoop(ctx, height, 0, func(txnLogs []*near.TxResult) error {
			for _, txnLog := range txnLogs {
				a.receiver.Logger.Info("height ", txnLog.BlockHeight)
				for _, el := range txnLog.EventLogs {
					res, evtType, err := a.par.Parse(&el)
					if err != nil {
						a.receiver.Logger.Debug(errors.Wrap(err, "Parse "))
						err = nil
						continue
					}
					nel := &chain.EventLogInfo{ContractAddress: common.NewAddress(el.Addr).String(), EventType: evtType, EventLog: res}

					a.receiver.Logger.Infof("IFirst %+v", nel)
					a.receiver.Logger.Infof("ISecond %+v", nel.EventLog)
					if a.fd.Match(nel) { //el.IDs is updated by match if matched
						//a.Log.Infof("Matched %+v", el)
						a.sinkChan <- nel
					}
				}
			}
			return nil
		})
		if err != nil {
			a.receiver.Logger.Errorf("receiveLoop terminated: %v", err)
			a.errChan <- err
		}
	}()
	return a.sinkChan, a.errChan, nil
}

// TransactWithBTS implements chain.ChainAPI
func (*api) TransactWithBTS(ownerKey string, method chain.ContractTransactMethodName, args []interface{}) (txnHash string, err error) {
	return "", nil
}

// Transfer implements chain.ChainAPI
func (a *api) Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp:") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".near") {
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

// TransferBatch implements chain.ChainAPI
func (a *api) TransferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp:") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".near") {
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

// WaitForTxnResult implements chain.ChainAPI
func (a *api) WaitForTxnResult(ctx context.Context, hash string) (*chain.TxnResult, error) {
	txRes, err := a.requester.cl.CallgetTransactionResult("transaction ID", "sender") // need to update and verify
	if err != nil {
		return nil, errors.Wrapf(err, "waitForResults(%v)", hash)
	}

	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.EventLogs {
		decodedLog, eventType, err := a.par.ParseTxn(&TxnEventLog{Addr: icon.Address(v.Addr), Indexed: v.Indexed, Data: v.Data})
		if err != nil {
			a.receiver.Logger.Trace(errors.Wrapf(err, "waitForResults.Parse %v", err))
			err = nil
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: string(v.Addr), EventType: eventType, EventLog: decodedLog})
	}
	statusCode, err := txRes.Status.Value()
	if err != nil {
		return nil, errors.Wrapf(err, "GetStatusCode err=%v", err)
	}
	return &chain.TxnResult{StatusCode: int(statusCode), ElInfo: plogs, Raw: txRes}, nil
}

// WatchForTransferEnd implements chain.ChainAPI
func (a *api) WatchForTransferEnd(id uint64, seq int64) error {
	return a.fd.watchFor(chain.TransferEnd, id, seq)
}

// WatchForTransferReceived implements chain.ChainAPI
func (a *api) WatchForTransferReceived(id uint64, seq int64) error {
	return a.fd.watchFor(chain.TransferReceived, id, seq)
}

// WatchForTransferStart implements chain.ChainAPI
func (a *api) WatchForTransferStart(id uint64, seq int64) error {
	return a.fd.watchFor(chain.TransferStart, id, seq)
}
