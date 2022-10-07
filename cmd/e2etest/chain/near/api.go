package near

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	common "github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/reactivex/rxgo/v2"

	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

type api struct {
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
		return nil, errors.New("Expected URL for chain NEAR. Got ")
	} else if cfg.Name != chain.NEAR {
		return nil, fmt.Errorf("expected cfg.Name=NEAR Got %v", cfg.Name)
	}
	Client, err := NewClient(cfg.URL, l)
	if err != nil {
		fmt.Println(err)
	}

	recvr := &api{
		sinkChan: make(chan *chain.EventLogInfo),
		errChan:  make(chan error),
		fd:       NewFinder(l, cfg.ContractAddresses),
	}
	recvr.par, err = NewParser(cfg.ContractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "NewParser ")
	}
	recvr.requester, err = newRequestAPI(Client, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "newRequestAPI %v", err)
	}
	return recvr, err
}

// Approve implements chain.ChainAPI
func (a *api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	// Deposit
	return "", nil
}

// GetBTPAddress implements chain.ChainAPI
func (a *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + a.requester.networkID + "/" + addr
	return fullAddr
}

// GetCoinBalance implements chain.ChainAPI
func (a *api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".near") {
		return nil, fmt.Errorf("address should be BTP address of account in native chain. Got %v", addr)
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

	privKey, err := wallet.DecryptNearKeyStore([]byte(wal), []byte(secret))
	if err != nil {
		err = errors.Wrapf(err, "wallet.DecryptKeyStore %v", err)
	}

	w, err := wallet.NewNearwalletFromPrivateKey(privKey)
	if err != nil {
		err = errors.Wrapf(err, "wallet.DecryptKeyStore %v", err)
	}

	priv = hex.EncodeToString([]byte(*privKey))
	pub = hex.EncodeToString(w.PublicKey())
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
	return a.requester.networkID
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
	return txnHash, err
}

// Subscribe implements chain.ChainAPI
func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	var opts common.SubscribeOptions
	opts.Seq++

	height, _ := a.requester.cl.GetLatestBlockHeight()

	go func() {
		defer close(errChan)
		if err := a.receiveTransactions(uint64(height), func(blockNotification *BlockNotification) error {
			log.WithFields(log.Fields{"height": blockNotification.Block().Height()}).Debug("block notification")
			for _, tx := range blockNotification.transactions {
				// filter out bmc, bsh
				if tx.ReceiverId == types.AccountId(a.cfg.ContractAddresses[chain.BMC]) || tx.ReceiverId == types.AccountId(a.cfg.ContractAddresses[chain.BTS]) {
					tx, err := a.requester.cl.GetTransactionResult(tx.Txid, tx.ReceiverId)
					if err != nil {
						return err
					}
					for _, outcome := range tx.ReceiptsOutcome {
						for _, log := range outcome.Outcome.Logs {
							res, evtType, err := a.par.Parse(log)
							if err != nil {
								err = nil
								continue
							}

							nel := &chain.EventLogInfo{ContractAddress: string(tx.Transaction.ReceiverId), EventType: evtType, EventLog: res}
							if a.fd.Match(nel) {
								a.sinkChan <- nel
							}
						}
					}
				}

			}
			for _, receipt := range blockNotification.Receipts() {
				for _, event := range receipt.Events {
					res, evtType, err := a.par.ParseMessage(event.Message)
					if err != nil {
						err = nil
						continue
					}
					nel := &chain.EventLogInfo{ContractAddress: a.cfg.ContractAddresses[chain.BMC], EventType: evtType, EventLog: res}
					if a.fd.Match(nel) {
						a.sinkChan <- nel
					}
				}
			}
			return nil
		}); err != nil {
			a.requester.cl.Logger().Errorf("receiveTransactions terminated: %v", err)
			a.errChan <- err
		}
	}()

	return sinkChan, errChan, nil
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
		err = fmt.Errorf("batch Transfers are supported for inter chain transfers only")
	} else {
		txnHash, err = a.requester.transferBatch(coinNames, senderKey, recepientAddress, amounts)
	}
	return
}

// WaitForTxnResult implements chain.ChainAPI
func (a *api) WaitForTxnResult(ctx context.Context, hash string) (*chain.TxnResult, error) {
	txRes, err := a.requester.cl.GetTransactionResult(types.NewCryptoHash(hash), types.AccountId(a.requester.contractNameToAddress[chain.BTS])) // need to update and verify
	if err != nil {
		return nil, errors.Wrapf(err, "WaitForTxnResult(%v)", hash)
	}

	plogs := []*chain.EventLogInfo{}
	plogs = append(plogs, &chain.EventLogInfo{ContractAddress: string(txRes.Transaction.ReceiverId), EventType: "", EventLog: txRes.TransactionOutcome.Outcome.Logs})
	// statusCode, err := strconv.Atoi(txRes.Status.SuccessValue)  // Need to Implement statusCode
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "GetStatusCode err=%v", err)
	// }
	return &chain.TxnResult{StatusCode: 200, ElInfo: plogs, Raw: txRes}, nil
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

// AddBlackListAddress implements chain.ChainAPI
func (a *api) AddBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_net": net, "_addresses": addrs}, "addBlacklistAddress", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
}

// ChangeRestriction implements chain.ChainAPI
func (a *api) ChangeRestriction(ownerKey string, enable bool) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	if enable {
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{}, "addRestriction", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{}, "disableRestrictions", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
}

// ChargedGasFee implements chain.ChainAPI
func (a *api) ChargedGasFee(txnHash string) (*big.Int, error) {
	txr, err := a.requester.cl.GetTransactionResult(types.NewCryptoHash(txnHash), types.AccountId(a.requester.contractNameToAddress[chain.BTS])) // check if Acc ID is correct
	if err != nil {
		return nil, errors.Wrapf(err, "TransactionByHash %v", err)
	}
	gasUsed := txr.TransactionOutcome.Outcome.GasBurnt
	gasUsedInBigInt := big.NewInt(int64(gasUsed)) // check if this is correct
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "BigInt Conversion %v", err)
	// }
	return (&big.Int{}).Mul(big.NewInt(12500000000), gasUsedInBigInt), nil
}

// GetAccumulatedFees implements chain.ChainAPI
func (a *api) GetAccumulatedFees() (map[string]*big.Int, error) {
	return a.requester.getAccumulatedFees()
}

// GetBlackListedUsers implements chain.ChainAPI
func (a *api) GetBlackListedUsers(net string, startCursor int, endCursor int) (users []string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_net": net, "_start": "0x0", "_end": "0x64"}, "get_blacklisted_user")
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, errors.New("getBlackListedUsers result is nil")
	}
	intArr, ok := res.([]interface{})
	if !ok {
		return nil, fmt.Errorf("getBlackListedUsers Response Expected []interface Got %T", res)
	} else if ok && len(intArr) == 0 {
		return
	}
	users = make([]string, len(intArr))
	for i, v := range intArr {
		users[i], ok = v.(string)
		if !ok {
			return nil, fmt.Errorf("getBlackListedUsers; response element Expected string Got %T", v)
		}
	}
	return
}

// GetConfigRequestEvent implements chain.ChainAPI
func (a *api) GetConfigRequestEvent(evtType chain.EventLogType, hash string) (*chain.EventLogInfo, error) {
	txRes, err := a.requester.cl.GetTransactionResult(types.NewCryptoHash(hash), types.AccountId(a.requester.contractNameToAddress[chain.BTS]))
	if err != nil {
		return nil, errors.Wrapf(err, "GetTransactionResult %v", err)
	}
	// if txRes.Status != types.NewHexInt(1) {                                              // How do we check status in Near
	// 	return nil, errors.Wrapf(err, "Expected Status Code 1. Got %v", txRes.Status)
	// }

	for _, outcome := range txRes.ReceiptsOutcome {
		for _, log := range outcome.Outcome.Logs {
			tmpRes, eventType, err := a.par.Parse(log)
			if eventType != evtType {
				continue
			}
			if err != nil {
				return nil, errors.Wrapf(err, "ParseTxn(%v) %v", eventType, err)
			}
			if eventType == chain.AddToBlacklistRequest {
				res, ok := tmpRes.(*chain.AddToBlacklistRequestEvent)
				if !ok {
					return nil, fmt.Errorf("expected *chain.AddToBlacklistRequestEvent; Got %T", tmpRes)
				}
				return &chain.EventLogInfo{ContractAddress: string(txRes.Transaction.ReceiverId), EventType: eventType, EventLog: res}, nil
			} else if eventType == chain.RemoveFromBlacklistRequest {
				res, ok := tmpRes.(*chain.RemoveFromBlacklistRequestEvent)
				if !ok {
					return nil, fmt.Errorf("expected *chain.RemoveFromBlacklistRequestEvent; Got %T", tmpRes)
				}
				return &chain.EventLogInfo{ContractAddress: string(txRes.Transaction.ReceiverId), EventType: eventType, EventLog: res}, nil
			} else if eventType == chain.TokenLimitRequest {
				res, ok := tmpRes.(*chain.TokenLimitRequestEvent)
				if !ok {
					return nil, fmt.Errorf("expected *chain.TokenLimitRequestEvent; Got %T", tmpRes)
				}
				return &chain.EventLogInfo{ContractAddress: string(txRes.Transaction.ReceiverId), EventType: eventType, EventLog: res}, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find %v; NumEventLogs %v", evtType, len(txRes.ReceiptsOutcome))
}

// GetFeeGatheringTerm implements chain.ChainAPI
func (a *api) GetFeeGatheringTerm() (interval uint64, err error) {
	bmcAddr, ok := a.requester.contractNameToAddress[chain.BMC]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BMC)
		return
	}
	res, err := a.requester.callContract(bmcAddr, map[string]interface{}{}, "getFeeGatheringTerm")
	if err != nil {
		return 0, errors.Wrap(err, "callContract getFeeGatheringTerm ")
	} else if res == nil {
		return 0, errors.New("callContract getFeeGatheringTerm returned nil value ")
	}
	tmpStr, ok := res.(string)
	if !ok {
		return 0, fmt.Errorf("expected type string Got %T", res)
	}
	return hexutil.DecodeUint64(tmpStr)
}

// GetFeeRatio implements chain.ChainAPI
func (a *api) GetFeeRatio(coinName string) (feeNumerator *big.Int, fixedFee *big.Int, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_name": coinName}, "feeRatio")
	if err != nil {
		return nil, nil, errors.Wrap(err, "callContract feeRatio ")
	} else if res == nil {
		return nil, nil, errors.New("callContract returned nil value ")
	}
	feeMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("Expected type map[string]interface{} Got %T", res)
	}
	getFeeOfType := func(feeMap map[string]interface{}, key string) (bal *big.Int, err error) {
		tmp, ok := feeMap[key]
		if !ok {
			return nil, fmt.Errorf("")
		}
		tmpStr, ok := tmp.(string)
		if !ok {
			return nil, fmt.Errorf("Expected type string Got %T", tmp)
		}
		bal = new(big.Int)
		bal.SetString(tmpStr[2:], 16)
		return
	}
	feeNumerator, err = getFeeOfType(feeMap, "feeNumerator")
	if err != nil {
		return nil, nil, errors.Wrapf(err, "getFeeOfType(feeNumerator) %v", err)
	}
	fixedFee, err = getFeeOfType(feeMap, "fixedFee")
	if err != nil {
		return nil, nil, errors.Wrapf(err, "getFeeOfType(fixedFee) %v", err)
	}
	return
}

// GetTokenLimit implements chain.ChainAPI
func (a *api) GetTokenLimit(coinName string) (tokenLimit *big.Int, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_name": coinName}, "get_token_limit")
	res = res.(types.CallFunctionResponse).Result
	if err != nil {
		err = errors.Wrapf(err, "CallContract %v", err)
		return
	} else if res == nil {
		err = errors.New("getTokenLimit result is nil")
		return
	}
	tmpStr, ok := res.(string)
	if !ok {
		err = fmt.Errorf("expected type string Got %T", res)
		return
	}
	tokenLimit = new(big.Int)
	tokenLimit.SetString(tmpStr[2:], 16)
	return
}

// GetTokenLimitStatus implements chain.ChainAPI
func (a *api) GetTokenLimitStatus(net string, coinName string) (response bool, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_net": net, "_coinName": coinName}, "tokenLimitStatus")
	if err != nil {
		return false, err
	} else if res == nil {
		return false, errors.New("tokenLimitStatus result is nil")
	}
	resStr, ok := res.(string)
	if !ok {
		return false, fmt.Errorf("tokenLimitStatus Response Expected string Got %T", res)
	}
	response = true
	if resStr == "0x0" {
		response = false
	}
	return
}

// IsBTSOwner implements chain.ChainAPI
func (a *api) IsBTSOwner(addr string) (response bool, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	if strings.Contains(addr, "btp:") {
		splts := strings.Split(addr, "/")
		addr = splts[len(splts)-1]
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_addr": addr}, "isOwner")
	if err != nil {
		return false, err
	} else if res == nil {
		return false, errors.New("isOwner result is nil")
	}
	resStr, ok := res.(string)
	if !ok {
		return false, fmt.Errorf("isOwner Response Expected string Got %T", res)
	}
	response = true
	if resStr == "0x0" {
		response = false
	}
	return response, err
}

// IsUserBlackListed implements chain.ChainAPI
func (a *api) IsUserBlackListed(net string, addr string) (response bool, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	if strings.Contains(addr, "btp:") {
		splts := strings.Split(addr, "/")
		addr = splts[len(splts)-1]
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_net": net, "_address": addr}, "isUserBlackListed")
	if err != nil {
		err = errors.Wrapf(err, "CallContract %v", err)
		return
	} else if res == nil {
		err = errors.New("isUserBlackListed result is nil")
		return
	}
	resStr, ok := res.(string)
	if !ok {
		return false, fmt.Errorf("isUserBlackListed Response Expected string Got %T", res)
	}
	response = true
	if resStr == "0x0" {
		response = false
	}
	return
}

// RemoveBlackListAddress implements chain.ChainAPI
func (a *api) RemoveBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_net": net, "_addresses": addrs}, "removeBlacklistAddress", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
}

// SetFeeGatheringTerm implements chain.ChainAPI
func (a *api) SetFeeGatheringTerm(ownerKey string, interval uint64) (hash string, err error) {
	bmcAddr, ok := a.requester.contractNameToAddress[chain.BMC]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BMC)
		return
	}
	return a.requester.transactWithContract(ownerKey, bmcAddr, big.NewInt(0), map[string]interface{}{"_value": hexutil.EncodeUint64(interval)}, "setFeeGatheringTerm", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
}

// SetFeeRatio implements chain.ChainAPI
func (a *api) SetFeeRatio(ownerKey string, coinName string, feeNumerator *big.Int, fixedFee *big.Int) (hash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	_feeNumerator := intconv.FormatBigInt(feeNumerator)
	_fixedFee := intconv.FormatBigInt(fixedFee)
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_name": coinName, "_feeNumerator": _feeNumerator, "_fixedFee": _fixedFee}, "setFeeRatio", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
}

// SetTokenLimit implements chain.ChainAPI
func (a *api) SetTokenLimit(ownerKey string, coinNames []string, tokenLimits []*big.Int) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	strTokenLimits := make([]string, len(tokenLimits))
	for i, v := range tokenLimits {
		strTokenLimits[i] = intconv.FormatBigInt(v)
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_coinNames": coinNames, "_tokenLimits": strTokenLimits}, "setTokenLimit", int64(a.requester.gasLimit[chain.DefaultGasLimit]))
}

// SuggestGasPrice implements chain.ChainAPI
func (*api) SuggestGasPrice() *big.Int {
	return big.NewInt(12500000000) // check if this value is correct
}

// WatchForAddToBlacklistRequest implements chain.ChainAPI
func (a *api) WatchForAddToBlacklistRequest(ID uint64, seq int64) error {
	return a.fd.watchFor(chain.AddToBlacklistRequest, ID, seq)
}

// WatchForBlacklistResponse implements chain.ChainAPI
func (*api) WatchForBlacklistResponse(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

// WatchForFeeGatheringRequest implements chain.ChainAPI
func (a *api) WatchForFeeGatheringRequest(ID uint64, addr string) error {
	return a.fd.watchFor(chain.FeeGatheringRequest, ID, addr)
}

// WatchForFeeGatheringTransferStart implements chain.ChainAPI
func (*api) WatchForFeeGatheringTransferStart(ID uint64, addr string) error {
	return errors.New("not implemented")
}

// WatchForRemoveFromBlacklistRequest implements chain.ChainAPI
func (a *api) WatchForRemoveFromBlacklistRequest(ID uint64, seq int64) error {
	return a.fd.watchFor(chain.RemoveFromBlacklistRequest, ID, seq)
}

// WatchForSetTokenLmitRequest implements chain.ChainAPI
func (a *api) WatchForSetTokenLmitRequest(ID uint64, seq int64) error {
	return a.fd.watchFor(chain.TokenLimitRequest, ID, seq)
}

// WatchForSetTokenLmitResponse implements chain.ChainAPI
func (*api) WatchForSetTokenLmitResponse(ID uint64, seq int64) error {
	return errors.New("not implemented")
}

func (a *api) receiveTransactions(height uint64, processBlockNotification func(blockNotification *BlockNotification) error) error {
	return a.requester.cl.MonitorTransactions(height, func(observable rxgo.Observable) error {
		result := observable.Observe()
		for item := range result {
			if err := item.E; err != nil {
				return err
			}

			bn, _ := item.V.(*BlockNotification)

			if *bn.Block().Hash() != [32]byte{} {
				return processBlockNotification(bn)
			}
		}
		return nil
	})
}

func (a *api) StopSubscriptionMethod() {
	a.requester.cl.CloseMonitor()
}
