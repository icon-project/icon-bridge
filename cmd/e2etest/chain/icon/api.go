package icon

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	gocommon "github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	"github.com/icon-project/icon-bridge/common"

	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature             = "Message(str,int,bytes)"
	MonitorBlockMaxConcurrency = 1
)

type api struct {
	*icon.ReceiverCore
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	par       *parser
	fd        *finder
	requester *requestAPI
	cfg       *chain.Config
}

func NewApi(l log.Logger, cfg *chain.Config) (chain.ChainAPI, error) {
	var err error
	if len(cfg.URL) == 0 {
		return nil, errors.New("Expected URL for chain ICON. Got ")
	} else if cfg.Name != chain.ICON {
		return nil, fmt.Errorf("Expected cfg.Name=ICON Got %v", cfg.Name)
	}
	client := icon.NewClient(cfg.URL, l)

	btsIconAddr, ok := cfg.ContractAddresses[chain.BTS]
	if !ok {
		return nil, errors.New("cfg.ConftractAddresses does not include chain.BTS")
	}
	bmcIconAddr, ok := cfg.ContractAddresses[chain.BMC]
	if !ok {
		return nil, errors.New("cfg.ConftractAddresses does not include chain.BMC")
	}

	evtReq := types.BlockRequest{
		EventFilters: []*types.EventFilter{
			{
				Addr:      types.Address(btsIconAddr),
				Signature: "TransferStart(Address,str,int,bytes)",
				Indexed:   []*string{},
			},
			{
				Addr:      types.Address(btsIconAddr),
				Signature: "TransferReceived(str,Address,int,bytes)",
				Indexed:   []*string{},
			},
			{
				Addr:      types.Address(btsIconAddr),
				Signature: "TransferEnd(Address,int,int,bytes)",
				Indexed:   []*string{},
			},
			{
				Addr:      icon.Address(bmcIconAddr),
				Signature: "Message(str,int,bytes)",
				Indexed:   []*string{},
			},
		},
	}
	recvr := &api{
		ReceiverCore: &icon.ReceiverCore{
			Log:      l,
			Cl:       client,
			BlockReq: evtReq,
			Opts:     icon.ReceiverOptions{},
		},
		sinkChan: make(chan *chain.EventLogInfo),
		errChan:  make(chan error),
		fd:       NewFinder(l, cfg.ContractAddresses),
	}
	recvr.par, err = NewParser(cfg.ContractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "NewParser ")
	}
	recvr.requester, err = newRequestAPI(client, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "newRequestAPI %v", err)
	}
	return recvr, err
}

func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	blk, err := a.Cl.GetLastBlock()
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetLastBlock ")
	}
	height := uint64(blk.Height)
	a.Log.Infof("Subscribe Start Height %v", height)
	// _errCh := make(chan error)
	go func() {
		// defer close(_errCh)
		err := a.ReceiveLoop(ctx, height, 0, func(txnLogs []*icon.TxResult) error {
			for _, txnLog := range txnLogs {
				for _, el := range txnLog.EventLogs {
					res, evtType, err := a.par.Parse(&el)
					if err != nil {
						//a.Log.Trace(errors.Wrapf(err, "Parseth %v", err))
						err = nil
						continue
					}
					nel := &chain.EventLogInfo{ContractAddress: common.NewAddress(el.Addr).String(), EventType: evtType, EventLog: res}

					a.Log.Debugf("IFirst %+v", nel)
					a.Log.Debugf("ISecond %+v", nel.EventLog)
					if a.fd.Match(nel) { //el.IDs is updated by match if matched
						//a.Log.Infof("Matched %+v", el)
						a.sinkChan <- nel
					}
				}
			}
			return nil
		})
		if err != nil {
			a.Log.Errorf("receiveLoop terminated: %v", err)
			a.errChan <- err
		}
	}()
	return a.sinkChan, a.errChan, nil
}

func (a *api) Transfer(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp:") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".icon") {
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
	if strings.Contains(recepientAddress, ".icon") {
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

func (a *api) WaitForTxnResult(ctx context.Context, hash string) (*chain.TxnResult, error) {
	_, txRes, err := a.Cl.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.HexBytes(hash)})
	if err != nil {
		return nil, errors.Wrapf(err, "waitForResults(%v)", hash)
	}

	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.EventLogs {
		decodedLog, eventType, err := a.par.ParseTxn(&TxnEventLog{Addr: types.Address(v.Addr), Indexed: v.Indexed, Data: v.Data})
		if err != nil {
			a.Log.Trace(errors.Wrapf(err, "waitForResults.Parse %v", err))
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

func (a *api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.approve(coinName, ownerKey, amount)
	return
}

func (a *api) Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txnHash, err = a.requester.reclaim(coinName, ownerKey, amount)
	return
}

func (a *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + a.requester.networkID + ".icon/" + addr
	return fullAddr
}

func (a *api) NativeCoin() string {
	return a.requester.nativeCoin
}

func (a *api) NativeTokens() []string {
	nativeTokens := []string{}
	for name := range a.requester.nativeTokensAddr {
		nativeTokens = append(nativeTokens, name)
	}
	return nativeTokens
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
	return a.fd.watchFor(chain.AddToBlacklistRequest, ID, seq)
}
func (a *api) WatchForRemoveFromBlacklistRequest(ID uint64, seq int64) error {
	return a.fd.watchFor(chain.RemoveFromBlacklistRequest, ID, seq)
}

func (a *api) WatchForSetTokenLmitRequest(ID uint64, seq int64) error {
	return a.fd.watchFor(chain.TokenLimitRequest, ID, seq)
}

func (a *api) WatchForFeeGatheringRequest(ID uint64, addr string) error {
	return a.fd.watchFor(chain.FeeGatheringRequest, ID, addr)
}

func (a *api) WatchForFeeGatheringTransferStart(ID uint64, addr string) error {
	return errors.New("not implemented")
}

func (a *api) WatchForBlacklistResponse(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) WatchForSetTokenLmitResponse(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) SetFeeRatio(ownerKey string, coinName string, feeNumerator, fixedFee *big.Int) (string, error) {
	return a.requester.setFeeRatio(ownerKey, coinName, feeNumerator, fixedFee)
}

func (a *api) GetFeeRatio(coinName string) (feeNumerator *big.Int, fixedFee *big.Int, err error) {
	return a.requester.getFeeRatio(coinName)
}

func (a *api) GetAccumulatedFees() (map[string]*big.Int, error) {
	return a.requester.getAccumulatedFees()
}

func (a *api) SetFeeGatheringTerm(ownerKey string, interval uint64) (hash string, err error) {
	return a.requester.setFeeGatheringTerm(ownerKey, interval)
}

func (a *api) GetFeeGatheringTerm() (interval uint64, err error) {
	return a.requester.getFeeGatheringTerm()
}

func (a *api) GetConfigRequestEvent(evtType chain.EventLogType, hash string) (*chain.EventLogInfo, error) {
	txRes, err := a.Cl.GetTransactionResult(&icon.TransactionHashParam{Hash: icon.HexBytes(hash)})
	if err != nil {
		return nil, errors.Wrapf(err, "GetTransactionResult %v", err)
	}
	if txRes.Status != icon.NewHexInt(1) {
		return nil, errors.Wrapf(err, "Expected Status Code 1. Got %v", txRes.Status)
	}

	for _, log := range txRes.EventLogs {
		tmpRes, eventType, err := a.par.ParseTxn(&TxnEventLog{Addr: log.Addr, Indexed: log.Indexed, Data: log.Data})
		if eventType != evtType {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "ParseTxn(%v) %v", eventType, err)
		}
		if eventType == chain.AddToBlacklistRequest {
			res, ok := tmpRes.(*chain.AddToBlacklistRequestEvent)
			if !ok {
				return nil, fmt.Errorf("Expected *chain.AddToBlacklistRequestEvent; Got %T", tmpRes)
			}
			return &chain.EventLogInfo{ContractAddress: string(log.Addr), EventType: eventType, EventLog: res}, nil
		} else if eventType == chain.RemoveFromBlacklistRequest {
			res, ok := tmpRes.(*chain.RemoveFromBlacklistRequestEvent)
			if !ok {
				return nil, fmt.Errorf("Expected *chain.RemoveFromBlacklistRequestEvent; Got %T", tmpRes)
			}
			return &chain.EventLogInfo{ContractAddress: string(log.Addr), EventType: eventType, EventLog: res}, nil
		} else if eventType == chain.TokenLimitRequest {
			res, ok := tmpRes.(*chain.TokenLimitRequestEvent)
			if !ok {
				return nil, fmt.Errorf("Expected *chain.TokenLimitRequestEvent; Got %T", tmpRes)
			}
			return &chain.EventLogInfo{ContractAddress: string(log.Addr), EventType: eventType, EventLog: res}, nil
		}
	}
	return nil, fmt.Errorf("Unable to find %v; NumEventLogs %v", evtType, len(txRes.EventLogs))
}

func (a *api) GetAddToBlacklistRequestEvent(hash string) (*chain.AddToBlacklistRequestEvent, error) {
	txRes, err := a.Cl.GetTransactionResult(&icon.TransactionHashParam{Hash: icon.HexBytes(hash)})
	if err != nil {
		return nil, errors.Wrapf(err, "GetTransactionResult %v", err)
	}
	if txRes.Status != icon.NewHexInt(0) {
		return nil, errors.Wrapf(err, "Expected Status Code 0. Got %v", txRes.Status)
	}
	for _, log := range txRes.EventLogs {
		tmpRes, eventType, err := a.par.ParseTxn(&TxnEventLog{Addr: log.Addr, Indexed: log.Indexed, Data: log.Data})
		if eventType != chain.AddToBlacklistRequest {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "ParseTxn(%v) %v", eventType, err)
		}
		res, ok := tmpRes.(*chain.AddToBlacklistRequestEvent)
		if !ok {
			return nil, fmt.Errorf("Expected *chain.AddToBlacklistRequestEvent; Got %T", tmpRes)
		}
		return res, nil
	}
	return nil, fmt.Errorf("Unable to find *chain.AddToBlacklistRequestEvent; NumEventLogs %v", len(txRes.EventLogs))
}

func (a *api) GetRemoveFromBlacklistRequestEvent(hash string) (*chain.RemoveFromBlacklistRequestEvent, error) {
	txRes, err := a.Cl.GetTransactionResult(&icon.TransactionHashParam{Hash: icon.HexBytes(hash)})
	if err != nil {
		return nil, errors.Wrapf(err, "GetTransactionResult %v", err)
	}
	if txRes.Status != icon.NewHexInt(0) {
		return nil, errors.Wrapf(err, "Expected Status Code 0. Got %v", txRes.Status)
	}
	for _, log := range txRes.EventLogs {
		tmpRes, eventType, err := a.par.ParseTxn(&TxnEventLog{Addr: log.Addr, Indexed: log.Indexed, Data: log.Data})
		if eventType != chain.RemoveFromBlacklistRequest {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "ParseTxn(%v) %v", eventType, err)
		}
		res, ok := tmpRes.(*chain.RemoveFromBlacklistRequestEvent)
		if !ok {
			return nil, fmt.Errorf("Expected *chain.RemoveFromBlacklistRequestEvent; Got %T", tmpRes)
		}
		return res, nil
	}
	return nil, fmt.Errorf("Unable to find *chain.RemoveFromBlacklistRequestEvent; NumEventLogs %v", len(txRes.EventLogs))
}

func (a *api) GetTokenLimitRequestEvent(hash string) (*chain.TokenLimitRequestEvent, error) {
	txRes, err := a.Cl.GetTransactionResult(&icon.TransactionHashParam{Hash: icon.HexBytes(hash)})
	if err != nil {
		return nil, errors.Wrapf(err, "GetTransactionResult %v", err)
	}
	if txRes.Status != icon.NewHexInt(0) {
		return nil, errors.Wrapf(err, "Expected Status Code 0. Got %v", txRes.Status)
	}
	for _, log := range txRes.EventLogs {
		tmpRes, eventType, err := a.par.ParseTxn(&TxnEventLog{Addr: log.Addr, Indexed: log.Indexed, Data: log.Data})
		if eventType != chain.TokenLimitRequest {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "ParseTxn(%v) %v", eventType, err)
		}
		res, ok := tmpRes.(*chain.TokenLimitRequestEvent)
		if !ok {
			return nil, fmt.Errorf("Expected *chain.TokenLimitRequestEvent; Got %T", tmpRes)
		}
		return res, nil
	}
	return nil, fmt.Errorf("Unable to find *chain.TokenLimitRequestEvent; NumEventLogs %v", len(txRes.EventLogs))
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

func (a *api) ChargedGasFee(txnHash string) (*big.Int, error) {
	txr, err := a.Cl.GetTransactionResult(&icon.TransactionHashParam{Hash: icon.HexBytes(txnHash)})
	if err != nil {
		return nil, errors.Wrapf(err, "TransactionByHash %v", err)
	}
	gasUsed, err := txr.StepUsed.BigInt()
	if err != nil {
		return nil, errors.Wrapf(err, "BigInt Conversion %v", err)
	}
	return (&big.Int{}).Mul(big.NewInt(12500000000), gasUsed), nil
}

func (a *api) SuggestGasPrice() *big.Int {
	return big.NewInt(12500000000)
}
