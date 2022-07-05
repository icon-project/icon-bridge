package hmny

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
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
	bshcore "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/bsh/bshcore"
	erc20 "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/bsh/erc20tradable"
	bep20tkn "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/tokenbsh/bep20tkn"
	bshproxy "github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny/abi/tokenbsh/bshproxy"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

// func (r *requestAPI) GetAddressFromPrivKey(key string) (*string, error) {
// 	return getAddressFromPrivKey(key)
// }

const (
	DefaultGasLimit = 80000000
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	ethCl                 *ethclient.Client
	log                   log.Logger
	bshc                  *bshcore.Bshcore
	erc                   *erc20.Erc20tradable
	bep                   *bep20tkn.BEP
	tokbsh                *bshproxy.TokenBSH
}

func newRequestAPI(url string, l log.Logger, contractNameToAddress map[chain.ContractName]string, networkID string) (*requestAPI, error) {

	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, errors.Wrapf(err, "rpc.Dial(%v)", url)
	}
	cleth := ethclient.NewClient(clrpc)

	caddr, ok := contractNameToAddress[chain.NativeBSHCoreHmy]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include %v", chain.NativeBSHCoreHmy)
	}
	bshc, err := bshcore.NewBshcore(common.HexToAddress(caddr), cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewBshcore ")
	}
	coinAddress, err := bshc.CoinId(&bind.CallOpts{Pending: false, Context: nil}, "ICX")
	if err != nil {
		return nil, errors.Wrap(err, "bshc.CoinId ")
	}
	caddr, ok = contractNameToAddress[chain.Erc20Hmy]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include %v", chain.Erc20Hmy)
	}
	bep, err := bep20tkn.NewBEP(common.HexToAddress(caddr), cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewBEP")
	}
	erc, err := erc20.NewErc20tradable(coinAddress, cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewErc20tradable")
	}
	caddr, ok = contractNameToAddress[chain.TokenBSHProxyHmy]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include %v", chain.TokenBSHProxyHmy)
	}
	tokbsh, err := bshproxy.NewTokenBSH(common.HexToAddress(caddr), cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewTokenBSH")
	}
	a := &requestAPI{
		log:                   l,
		contractNameToAddress: contractNameToAddress,
		networkID:             networkID,
		ethCl:                 cleth,
		bshc:                  bshc,
		erc:                   erc,
		bep:                   bep,
		tokbsh:                tokbsh,
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

func getAddressFromPrivKey(pKey string) (*string, error) {
	privBytes, err := hex.DecodeString(pKey)
	if err != nil {
		return nil, errors.Wrap(err, "DecodeString ")
	}
	privKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, errors.Wrap(err, "ToECDSA ")
	}
	addr := crypto.PubkeyToAddress(privKey.PublicKey).String()
	return &addr, nil
}

func (r *requestAPI) getHmnyBalance(addr string) (*big.Int, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return r.ethCl.BalanceAt(ctx, common.HexToAddress(addr), nil)
}

func (r *requestAPI) getHmnyErc20Balance(addr string) (val *big.Int, err error) {
	//common.HexToAddress("0x8fc668275b4fa032342ea3039653d841f069a83b")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := bind.CallOpts{Pending: false, Context: ctx}
	return r.bep.BalanceOf(&opts, common.HexToAddress(addr))
}

func (r *requestAPI) getHmnyWrappedICX(addr string) (val *big.Int, err error) {
	coinName := "ICX"
	v, err := r.bshc.GetBalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, common.HexToAddress(addr), coinName)
	if err != nil {
		err = errors.Wrap(err, "bshc.GetBalanceOf ")
		return
	}
	return v.UsableBalance, nil
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

func (r *requestAPI) transferNativeWithin(senderKey, recepientAddress string, amount big.Int) (txnHash string, logs interface{}, err error) {
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
	// rcpts, err := r.waitForResults(context.TODO(), signedTx.Hash())
	// logs = rcpts.Logs
	return
}

func (r *requestAPI) transferTokenWithin(senderKey, recepientAddress string, amount big.Int) (txnHash string, logs interface{}, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txn, err := r.bep.Transfer(txo, common.HexToAddress(recepientAddress), &amount)
	if err != nil {
		err = errors.Wrap(err, "bep.Transfer ")
		return
	}
	txnHash = txn.Hash().String()
	// rcpts, err := r.waitForResults(context.TODO(), txn.Hash())
	// logs = rcpts.Logs
	return
}

func (r *requestAPI) transferNativeCrossChain(senderKey string, recepientAddress string, amount big.Int) (txnHash string, logs interface{}, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Value = &amount
	txo.Context = context.Background()
	txn, err := r.bshc.TransferNativeCoin(txo, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "bshc.TransferNativeCoin ")
		return
	}
	txnHash = txn.Hash().String()
	// rcpts, err := r.waitForResults(context.TODO(), txn.Hash())
	// logs = rcpts.Logs
	return
}

func (r *requestAPI) transferWrappedCrossChain(coinName string, senderKey, recepientAddress string, amount big.Int) (txnHash string, logs interface{}, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txn, err := r.bshc.Transfer(txo, coinName, &amount, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "bshc.Transfer ")
		return
	}
	txnHash = txn.Hash().String()
	// rcpts, err := r.waitForResults(context.TODO(), txn.Hash())
	// logs = rcpts.Logs
	return
}

func (r *requestAPI) approveToken(coinName, senderKey string, amount big.Int) (txnHash string, logs interface{}, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	caddr, ok := r.contractNameToAddress[chain.TokenBSHProxyHmy]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include %v", chain.TokenBSHProxyHmy)
		return
	}
	txo.Context = context.Background()

	approveTxn, err := r.bep.Approve(txo, common.HexToAddress(caddr), &amount)
	if err != nil {
		err = errors.Wrapf(err, "bep.Approve %v  %v  %v", amount.String(), txo.GasPrice, txo.GasLimit)
		return
	}
	txnHash = approveTxn.Hash().String()
	return
}

func (r *requestAPI) transferTokenCrossChain(coinName string, senderKey, recepientAddress string, amount big.Int) (transferTxnHash string, transferLogs interface{}, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	transferTxn, err := r.tokbsh.Transfer(txo, coinName, &amount, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "tokbsh.Transfer ")
		return
	}
	transferTxnHash = transferTxn.Hash().String()
	return
}

func (r *requestAPI) approveCrossNativeCoin(coinName, senderKey string, amount big.Int) (approveTxnHash string, logs interface{}, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	caddr, ok := r.contractNameToAddress[chain.NativeBSHCoreHmy]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include %v ", chain.NativeBSHCoreHmy)
		return
	}
	txo.Context = context.Background()
	approveTxn, err := r.erc.Approve(txo, common.HexToAddress(caddr), &amount)
	if err != nil {
		err = errors.Wrap(err, "erc.Approve ")
		return
	}
	approveTxnHash = approveTxn.Hash().String()
	return
	// rcpts, err := r.waitForResults(context.TODO(), approveTxn.Hash())
	// if err != nil {
	// 	return
	// }
	// logs = rcpts.Logs
	// ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()
	// ownerWallet, _, err := GetWalletFromPrivKey(senderKey)
	// if err != nil {
	// 	err = errors.Wrap(err, "GetWalletFromPrivKey ")
	// 	return
	// }
	// caddr, ok = r.contractNameToAddress[chain.NativeBSHCoreHmy]
	// if !ok {
	// 	err = fmt.Errorf("contractNameToAddress doesn't include %v ", chain.NativeBSHCoreHmy)
	// 	return
	// }
	// allowanceAmount, err = r.erc.Allowance(&bind.CallOpts{Pending: false, Context: ctx}, common.HexToAddress(ownerWallet.Address()), common.HexToAddress(caddr))
	// if err != nil {
	// 	err = errors.Wrap(err, "erc.Allowance ")
	// }
	// return
}

// func GetWalletAndPrivKey(walFile string, password string) (wal *wallet.EvmWallet, pKey *ecdsr.PrivateKey, err error) {
// 	keyReader, err := os.Open(walFile)
// 	if err != nil {
// 		return
// 	}
// 	defer keyReader.Close()

// 	keyStore, err := ioutil.ReadAll(keyReader)
// 	if err != nil {
// 		return
// 	}

// 	pKey, err = wallet.DecryptEvmKeyStore(keyStore, []byte(password))
// 	wal = &wallet.EvmWallet{
// 		Skey: pKey,
// 		Pkey: &pKey.PublicKey,
// 	}
// 	return
// }
