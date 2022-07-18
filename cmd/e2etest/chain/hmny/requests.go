package hmny

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	btscore "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/btscore"
	erc20tradeable "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/erc20tradable"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	DefaultGasLimit = 80000000
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	ethCl                 *ethclient.Client
	log                   log.Logger
	btsc                  *btscore.Btscore
	ercPerCoin            sync.Map
}

func newRequestAPI(url string, l log.Logger, contractNameToAddress map[chain.ContractName]string, networkID string, tokenNameToAddr map[string]string) (*requestAPI, error) {

	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, errors.Wrapf(err, "rpc.Dial(%v)", url)
	}
	cleth := ethclient.NewClient(clrpc)

	caddr, ok := contractNameToAddress[chain.BTSCoreHmny]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include %v", chain.BTSCoreHmny)
	}
	btscore, err := btscore.NewBtscore(common.HexToAddress(caddr), cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewBtscore")
	}

	a := &requestAPI{
		log:                   l,
		contractNameToAddress: contractNameToAddress,
		networkID:             networkID,
		ethCl:                 cleth,
		btsc:                  btscore,
		ercPerCoin:            sync.Map{},
	}
	for name, coinAddress := range tokenNameToAddr {
		ercp, err := erc20tradeable.NewErc20tradable(common.HexToAddress(coinAddress), cleth)
		if err != nil {
			return nil, errors.Wrap(err, "NewErc20tradable")
		}
		a.ercPerCoin.Store(name, ercp)
	}
	return a, nil
}

func GetWalletFromPrivKey(privKey string) (wal *wallet.EvmWallet, pKey *ecdsa.PrivateKey, err error) {
	privBytes, err := hex.DecodeString(privKey)
	if err != nil {
		err = errors.Wrap(err, "hex.DecodeString ")
		return
	}
	ethPrivKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		err = errors.Wrap(err, "ToECDSA ")
		return
	}
	wal = &wallet.EvmWallet{
		Skey: ethPrivKey,
		Pkey: &ethPrivKey.PublicKey,
	}
	return wal, ethPrivKey, nil
}

func generateKeyPair() ([2]string, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return [2]string{}, errors.Wrap(err, "GenerateKey")
	}
	privStr := hex.EncodeToString(crypto.FromECDSA(privKey))
	if err != nil {
		return [2]string{}, errors.Wrap(err, "EncodeToString")
	}
	pubAddress := crypto.PubkeyToAddress(privKey.PublicKey).String()
	return [2]string{privStr, pubAddress}, nil
}

func (r *requestAPI) getTransactionRequest(senderKey string) (*bind.TransactOpts, error) {
	_, senderPrivKey, err := GetWalletFromPrivKey(senderKey)
	if err != nil {
		return nil, errors.Wrap(err, "GetWalletFromPrivKey")
	}
	chainID, err := r.ethCl.ChainID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "ChainID ")
	}
	txo, err := bind.NewKeyedTransactorWithChainID(senderPrivKey, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "NewKeyedTransactorWithChainID ")
	}
	txo.GasPrice, err = r.ethCl.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "SuggestGasPrice ")
	}
	txo.GasLimit = uint64(DefaultGasLimit)
	return txo, nil
}

func (r *requestAPI) waitForResults(ctx context.Context, txHash common.Hash) (txr *types.Receipt, err error) {
	const DefaultGetTransactionResultPollingInterval = 1500 * time.Millisecond //1.5sec
	ticker := time.NewTicker(time.Duration(DefaultGetTransactionResultPollingInterval) * time.Nanosecond)
	retryLimit := 10
	retryCounter := 0
	for {
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			err = errors.New("Context Cancelled. ResultWait Exiting ")
			return
		case <-ticker.C:
			if retryCounter >= retryLimit {
				err = errors.New("Retry Limit Exceeded while waiting for results of transaction")
				return
			}
			retryCounter++
			//r.log.Debugf("GetTransactionResult Attempt: %d", retryCounter)
			txr, err = r.ethCl.TransactionReceipt(context.Background(), txHash)
			if err != nil && err == ethereum.NotFound {
				r.log.Trace(errors.Wrap(err, "waitForResults "))
				err = nil
				continue
			}
			//r.log.Debugf("GetTransactionResult hash:%v, txr:%+v, err:%+v", txHash, txr, err)
			return
		}
	}
}

func (r *requestAPI) transferNativeIntraChain(senderKey, recepientAddress string, amount big.Int) (txnHash string, err error) {
	senderWallet, senderPrivKey, err := GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}
	nonce, err := r.ethCl.PendingNonceAt(context.Background(), common.HexToAddress(senderWallet.Address()))
	if err != nil {
		err = errors.Wrap(err, "PendingNonceAt ")
		return
	}
	gasPrice, err := r.ethCl.SuggestGasPrice(context.Background())
	if err != nil {
		err = errors.Wrap(err, "SuggestGasPrice ")
		return
	}
	chainID, err := r.ethCl.ChainID(context.Background())
	if err != nil {
		err = errors.Wrap(err, "ChainID ")
		return
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(recepientAddress), &amount, uint64(DefaultGasLimit), gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), senderPrivKey)
	if err != nil {
		err = errors.Wrap(err, "SignTx ")
		return
	}
	if err = r.ethCl.SendTransaction(context.TODO(), signedTx); err != nil {
		err = errors.Wrap(err, "SendTransaction ")
		return
	}
	txnHash = signedTx.Hash().String()
	return
}

func (r *requestAPI) transferTokenIntraChain(senderKey, recepientAddress string, amount big.Int, coinName string) (txnHash string, err error) {
	erc, err := r.getERC(coinName)
	if err != nil {
		err = fmt.Errorf("GetERC %v", err)
		return
	}

	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txn, err := erc.Transfer(txo, common.HexToAddress(recepientAddress), &amount)
	if err != nil {
		err = errors.Wrap(err, "hrc.Transfer ")
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) transferNativeCrossChain(senderKey string, recepientAddress string, amount big.Int) (txnHash string, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Value = &amount
	txo.Context = context.Background()
	txn, err := r.btsc.TransferNativeCoin(txo, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "btsc.TransferNativeCoin ")
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) transferTokensCrossChain(coinName string, senderKey, recepientAddress string, amount big.Int) (txnHash string, err error) {

	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txn, err := r.btsc.Transfer(txo, coinName, &amount, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "btsc.Transfer ")
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) getERC(coinName string) (erc *erc20tradeable.Erc20tradable, err error) {
	erc = &erc20tradeable.Erc20tradable{}
	res, ok := r.ercPerCoin.Load(coinName)
	if !ok {
		r.log.Debugf("Registering Input coinName %v ", coinName)
		coinAddress, errs := r.btsc.CoinId(&bind.CallOpts{Pending: false, Context: nil}, coinName)
		if err != nil {
			err = errors.Wrap(errs, "btsc.CoinId ")
			return
		}
		if erc, err = erc20tradeable.NewErc20tradable(coinAddress, r.ethCl); err != nil {
			err = errors.Wrap(err, "NewErc20tradable")
			return
		}
		res = erc
		r.ercPerCoin.Store(coinName, erc)
	} else if ok && res == nil {
		err = fmt.Errorf("ercPerCoin includes coin %v but value is nil", coinName)
		return
	}
	erc, ok = res.(*erc20tradeable.Erc20tradable)
	if !ok {
		err = fmt.Errorf("Expected type *erc20.Erc20tradable; Got %T", res)
		return
	}
	return erc, nil
}

func (r *requestAPI) approveCoin(coinName, senderKey string, amount big.Int) (approveTxnHash string, err error) {
	erc, err := r.getERC(coinName)
	if err != nil {
		err = fmt.Errorf("GetERC %v", err)
		return
	}
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	btscaddr, ok := r.contractNameToAddress[chain.BTSCoreHmny]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include %v ", chain.BTSCoreHmny)
		return
	}
	txo.Context = context.Background()
	approveTxn, err := erc.Approve(txo, common.HexToAddress(btscaddr), &amount)
	if err != nil {
		err = errors.Wrap(err, "erc.Approve ")
		return
	}
	approveTxnHash = approveTxn.Hash().String()
	return
}

func (r *requestAPI) getCoinBalance(coinName, addr string) (bal *chain.CoinBalance, err error) {
	b, err := r.btsc.BalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, common.HexToAddress(addr), coinName)
	if err != nil {
		err = errors.Wrap(err, "btsc.GetBalanceOf ")
		return
	}
	erc, err := r.getERC(coinName)
	if err != nil {
		err = fmt.Errorf("GetERC %v", err)
		return
	}
	btscaddr, ok := r.contractNameToAddress[chain.BTSCoreHmny]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include %v ", chain.BTSCoreHmny)
		return
	}
	allowance, err := erc.Allowance(&bind.CallOpts{Pending: false, Context: context.TODO()}, common.HexToAddress(addr), common.HexToAddress(btscaddr))
	if err != nil {
		err = fmt.Errorf("Allowance; err: %v", err)
	}

	bal = &chain.CoinBalance{
		Approved:   allowance,
		Usable:     b.UsableBalance,
		Locked:     b.LockedBalance,
		Refundable: b.RefundableBalance,
	}
	bal.Total = bal.Total.Add(bal.Locked, b.UsableBalance)
	bal.Total = bal.Total.Add(bal.Total, bal.Refundable)
	return bal, nil
}

/*
func (r *requestAPI) getHmnyBalance(addr string) (*big.Int, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return r.ethCl.BalanceAt(ctx, common.HexToAddress(addr), nil)
}

func (r *requestAPI) getWrappedCoinBalance(coinName string, addr string) (val *big.Int, err error) {
	v, err := r.btsc.BalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, common.HexToAddress(addr), coinName)
	if err != nil {
		err = errors.Wrap(err, "btsc.GetBalanceOf ")
		return
	}
	return v.UsableBalance, nil
}
func (r *requestAPI) getAllowance(coinName, ownerAddr string) (amont *big.Int, err error) {
	erc := &erc20.Erc20tradable{}
	res, ok := r.ercPerCoin.Load(coinName)
	if !ok {
		err = fmt.Errorf("ercPerCoin does not includes coin %v", coinName)
		return
	} else if ok && res == nil {
		err = fmt.Errorf("ercPerCoin includes coin %v but value is nil", coinName)
		return
	}
	if erc, ok = res.(*erc20.Erc20tradable); !ok {
		err = fmt.Errorf("Expected type *erc20.Erc20tradable; Got %T", res)
		return
	}
	btscaddr, ok := r.contractNameToAddress[chain.BTSCoreHmny]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include %v ", chain.BTSCoreHmny)
		return
	}
	return erc.Allowance(&bind.CallOpts{Pending: false, Context: context.TODO()}, common.HexToAddress(ownerAddr), common.HexToAddress(btscaddr))
}

*/
