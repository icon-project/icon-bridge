package near

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
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
func (*api) Approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	panic("unimplemented")
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
func (*api) GetBTPAddress(addr string) string {
	return ""
}

// GetCoinBalance implements chain.ChainAPI
func (*api) GetCoinBalance(coinName string, addr string) (*chain.CoinBalance, error) {
	return nil, nil
}

// GetKeyPairFromKeystore implements chain.ChainAPI
func (*api) GetKeyPairFromKeystore(keystoreFile string, secretFile string) (string, string, error) {
	return "", "", nil
}

// GetKeyPairs implements chain.ChainAPI
func (*api) GetKeyPairs(num int) ([][2]string, error) {
	return nil, nil
}

// GetNetwork implements chain.ChainAPI
func (*api) GetNetwork() string {
	return ""
}

// NativeCoin implements chain.ChainAPI
func (*api) NativeCoin() string {
	return ""
}

// NativeTokens implements chain.ChainAPI
func (*api) NativeTokens() []string {
	return []string{""}
}

// Reclaim implements chain.ChainAPI
func (*api) Reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	return "", nil
}

// Subscribe implements chain.ChainAPI
func (*api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	return nil, nil, nil
}

// TransactWithBTS implements chain.ChainAPI
func (*api) TransactWithBTS(ownerKey string, method chain.ContractTransactMethodName, args []interface{}) (txnHash string, err error) {
	return "", nil
}

// Transfer implements chain.ChainAPI
func (*api) Transfer(coinName string, senderKey string, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	return "", nil
}

// TransferBatch implements chain.ChainAPI
func (*api) TransferBatch(coinNames []string, senderKey string, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
	return "", nil
}

// WaitForTxnResult implements chain.ChainAPI
func (*api) WaitForTxnResult(ctx context.Context, hash string) (txnr *chain.TxnResult, err error) {
	return nil, nil
}

// WatchForTransferEnd implements chain.ChainAPI
func (*api) WatchForTransferEnd(ID uint64, seq int64) error {
	return nil
}

// WatchForTransferReceived implements chain.ChainAPI
func (*api) WatchForTransferReceived(ID uint64, seq int64) error {
	return nil
}

// WatchForTransferStart implements chain.ChainAPI
func (*api) WatchForTransferStart(ID uint64, seq int64) error {
	return nil
}
