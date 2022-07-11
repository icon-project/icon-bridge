package icon

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature             = "Message(str,int,bytes)"
	MonitorBlockMaxConcurrency = 50
)

type api struct {
	*icon.ReceiverCore
	networkID string
	sinkChan  chan *chain.EventLogInfo
	errChan   chan error
	par       *parser
	fd        *finder
	requester *requestAPI
}

func NewApi(l log.Logger, cfg *chain.ChainConfig) (chain.ChainAPI, error) {
	var err error
	if len(cfg.URL) == 0 {
		return nil, errors.New("List of Urls is empty ")
	}
	client := icon.NewClient(cfg.URL, l)

	btsIconAddr, ok := cfg.ConftractAddresses[chain.BTSIcon]
	if !ok {
		return nil, errors.New("cfg.ConftractAddresses does not include chain.BTSIcon")
	}

	evtReq := icon.BlockRequest{
		EventFilters: []*icon.EventFilter{
			{
				Addr:      icon.Address(btsIconAddr),
				Signature: "TransferStart(Address,str,int,bytes)",
				Indexed:   []*string{},
			},
			{
				Addr:      icon.Address(btsIconAddr),
				Signature: "TransferReceived(str,Address,int,bytes)",
				Indexed:   []*string{},
			},
			{
				Addr:      icon.Address(btsIconAddr),
				Signature: "TransferEnd(Address,int,int,bytes)",
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
		sinkChan:  make(chan *chain.EventLogInfo),
		errChan:   make(chan error),
		fd:        NewFinder(l, cfg.ConftractAddresses),
		networkID: cfg.NetworkID,
	}
	recvr.par, err = NewParser(cfg.ConftractAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "NewParser ")
	}
	recvr.requester, err = newRequestAPI(client, cfg.ConftractAddresses, cfg.NetworkID)
	return recvr, nil
}

func (r *api) Subscribe(ctx context.Context) (sinkChan chan *chain.EventLogInfo, errChan chan error, err error) {
	blk, err := r.Cl.GetLastBlock()
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetLastBlock ")
	}
	height := uint64(blk.Height)
	r.Log.Infof("Subscribe Start Height %v", height)
	// _errCh := make(chan error)
	go func() {
		// defer close(_errCh)
		err := r.ReceiveLoop(ctx, height, 0, func(txnLogs []*icon.TxResult) error {
			for _, txnLog := range txnLogs {
				for _, el := range txnLog.EventLogs {
					res, evtType, err := r.par.Parse(&el)
					if err != nil {
						r.Log.Trace(errors.Wrap(err, "Parse "))
						err = nil
						continue
					}

					el := &chain.EventLogInfo{ContractAddress: common.NewAddress(el.Addr).String(), EventType: evtType, EventLog: res}
					if r.fd.Match(el) { //el.IDs is updated by match if matched
						//r.Log.Infof("Matched %+v", el)
						r.sinkChan <- el
					}
				}
			}
			return nil
		})
		if err != nil {
			r.Log.Errorf("receiveLoop terminated: %v", err)
			r.errChan <- err
		}
	}()
	return r.sinkChan, r.errChan, nil
}

func (r *api) Transfer(coinName, senderKey, recepientAddress string, amount big.Int) (txnHash string, err error) {
	if !strings.Contains(recepientAddress, "btp://") {
		return "", errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	within := false
	if strings.Contains(recepientAddress, ".icon") {
		within = true
		splts := strings.Split(recepientAddress, "/")
		recepientAddress = splts[len(splts)-1]
	}
	if within {
		if coinName == "ICX" {
			txnHash, _, err = r.requester.transferNativeIntraChain(senderKey, recepientAddress, amount)
		} else if coinName == "TICX" {
			txnHash, _, err = r.requester.transferTokenIntraChain(senderKey, recepientAddress, amount)
		} else {
			err = fmt.Errorf("IntraChain transfers are supported for coins ICX and TICX only")
		}
	} else {
		if coinName == "ICX" {
			txnHash, _, err = r.requester.transferNativeCrossChain(senderKey, recepientAddress, amount)
		} else { // ONE, TONE, TICX
			txnHash, _, err = r.requester.transferWrappedCrossChain(coinName, senderKey, recepientAddress, amount)
		}
	}
	return
}

func (r *api) GetCoinBalance(coinName string, addr string) (*big.Int, error) {
	if !strings.Contains(addr, "btp://") {
		return nil, errors.New("Address should be BTP address. Use GetBTPAddress(hexAddr)")
	}
	if !strings.Contains(addr, ".icon") {
		return nil, fmt.Errorf("Address should be BTP address of account in native chain. Got %v", addr)
	}
	splts := strings.Split(addr, "/")
	address := splts[len(splts)-1]
	if coinName == "ICX" {
		return r.requester.getICXBalance(address)
	}
	return r.requester.getWrappedCoinBalance(coinName, address)
}

func (r *api) WaitForTxnResult(ctx context.Context, hash string) (interface{}, []*chain.EventLogInfo, error) {
	_, txRes, err := r.Cl.WaitForResults(ctx, &icon.TransactionHashParam{Hash: icon.HexBytes(hash)})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "waitForResults(%v)", hash)
	}
	plogs := []*chain.EventLogInfo{}
	for _, v := range txRes.EventLogs {
		decodedLog, eventType, err := r.par.ParseTxn(&TxnEventLog{Addr: icon.Address(v.Addr), Indexed: v.Indexed, Data: v.Data})
		if err != nil {
			r.Log.Trace(errors.Wrap(err, "waitForResults.Parse "))
			err = nil
			continue
			//return nil, nil, err
		}
		plogs = append(plogs, &chain.EventLogInfo{ContractAddress: string(v.Addr), EventType: eventType, EventLog: decodedLog})
	}
	return txRes, plogs, nil
}

func (r *api) Approve(coinName string, ownerKey string, amount big.Int) (txnHash string, err error) {
	if coinName == "BNB" || coinName == "TBNB" {
		txnHash, _, err = r.requester.approveCrossNativeCoin(coinName, ownerKey, amount)
	} else if coinName == "TICX" {
		txnHash, _, err = r.requester.approveToken(coinName, ownerKey, amount)
	} else {
		err = errors.Wrapf(err, "CoinName not among accepted Values ONE, ETH. Got %v", coinName)
	}
	return
}

func (r *api) GetChainType() chain.ChainType {
	return chain.ICON
}

func (r *api) GetBTPAddress(addr string) string {
	fullAddr := "btp://" + r.networkID + ".icon/" + addr
	return fullAddr
}

func (r *api) GetNativeCoin() string {
	return "ICX"
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
