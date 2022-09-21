package near

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"

	gocommon "github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	common "github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"

	"github.com/icon-project/icon-bridge/common/intconv"
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
		return nil, fmt.Errorf("expected cfg.Name=NEAR Got %v", cfg.Name)
	}
	Clients, err := near.NewClient(cfg.URL, l)
	if err != nil {
		fmt.Println(err)
	}
	Receiver, err := near.NewReceiver(near.ReceiverConfig{}, l, Clients)
	if err != nil {
		fmt.Println(err)
	}

	recvr := &api{
		receiver: Receiver,
		sinkChan: make(chan *chain.EventLogInfo),
		errChan:  make(chan error),
		fd:       NewFinder(l, cfg.ContractAddresses),
	}
	recvr.par, err = NewParser(cfg.ContractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "NewParser ")
	}
	fmt.Println(Clients)
	recvr.requester, err = newRequestAPI(Clients, cfg)
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
	if !strings.Contains(addr, ".near") {
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
	return txnHash, err
}

// Subscribe implements chain.ChainAPI
func (a *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	var opts common.SubscribeOptions
	opts.Seq++

	go func() {
		defer close(errChan)

		if err := a.receiver.ReceiveBlocks(opts.Height, "", func(blockNotification *types.BlockNotification) {
			log.WithFields(log.Fields{"height": blockNotification.Block().Height()}).Debug("block notification")
			receipts := blockNotification.Receipts()

			for _, receipt := range receipts {
				events := receipt.Events[:0]
				for _, event := range receipt.Events {
					switch {

					case event.Sequence == opts.Seq:
						events = append(events, event)
						opts.Seq++

					case event.Sequence > opts.Seq:
						log.WithFields(log.Fields{
							"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
						}).Error("invalid event seq")

						errChan <- fmt.Errorf("invalid event seq")
					}

					nel := &chain.EventLogInfo{ContractAddress: event.Next.ContractAddress(), EventType: chain.EventLogType(event.Next.Type()),
						EventLog: event.Next.BlockChain()}
					sinkChan <- nel
				}
			}
		}); err != nil {
			errChan <- err
		}
	}()

	return sinkChan, errChan, nil
}

// TransactWithBTS implements chain.ChainAPI
func (a *api) TransactWithBTS(ownerKey string, method chain.ContractTransactMethodName, args []interface{}) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	if args == nil {
		err = errors.New("Got nil args")
		return
	}
	if method == chain.AddBlackListAddress {
		if len(args) != 2 {
			return "", fmt.Errorf("expected 2 args _net, _addresses. Got %v", len(args))
		}
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_net": args[0], "_addresses": args[1]}, "addBlacklistAddress")
	} else if method == chain.RemoveBlackListAddress {
		if len(args) != 2 {
			return "", fmt.Errorf("expected 2 args _net, _addresses. Got %v", len(args))
		}
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_net": args[0], "_addresses": args[1]}, "removeBlacklistAddress")
	} else if method == chain.AddRestriction {
		if len(args) != 0 {
			return "", fmt.Errorf("expected 0 args. Got %v", len(args))
		}
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{}, "addRestriction")
	} else if method == chain.DisableRestrictions {
		if len(args) != 0 {
			return "", fmt.Errorf("expected 0 args. Got %v", len(args))
		}
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{}, "disableRestrictions")
	} else if method == chain.SetTokenLimit {
		if len(args) != 2 {
			return "", fmt.Errorf("expected 2 args for _coinNames, _tokenLimits. Got %v", len(args))
		}
		resArr, ok := args[1].([]*big.Int)
		if !ok {
			return "", fmt.Errorf("expected second arg _tokenLimits field of type []interface{}; Got %T", args[1])
		}
		_tokenLimits := make([]string, len(resArr))
		for i := 0; i < len(resArr); i++ {
			_tokenLimits[i] = intconv.FormatBigInt(resArr[i])
		}
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_coinNames": args[0], "_tokenLimits": _tokenLimits}, "setTokenLimit")
	}
	return "", fmt.Errorf("method %v not supported", method)
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
	var iapi near.IApi
	txRes, err := iapi.Transaction(hash) // need to update and verify
	if err != nil {
		return nil, errors.Wrapf(err, "waitForResults(%v)", hash)
	}

	plogs := []*chain.EventLogInfo{}
	plogs = append(plogs, &chain.EventLogInfo{ContractAddress: string(txRes.Transaction.ReceiverId), EventType: "", EventLog: txRes.TransactionOutcome.Outcome.Logs})
	statusCode, err := strconv.Atoi(txRes.Status.SuccessValue)
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
