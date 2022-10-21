package snow

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	btscore "github.com/icon-project/icon-bridge/cmd/e2etest/chain/snow/abi/btscore"
	btsperi "github.com/icon-project/icon-bridge/cmd/e2etest/chain/snow/abi/btsperiphery"
	erc20tradeable "github.com/icon-project/icon-bridge/cmd/e2etest/chain/snow/abi/erc20tradable"
	"github.com/icon-project/icon-bridge/common/wallet"
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	ethCl                 *ethclient.Client
	gasLimits             map[chain.GasLimitType]uint64
	nativeCoin            string
	nativeTokens          []string
	btsc                  *btscore.Btscore
	btsp                  *btsperi.Btsperiphery
	ercPerCoin            map[string]*erc20tradeable.Erc20tradable
}

func newRequestAPI(cfg *chain.Config) (*requestAPI, error) {
	if !strings.Contains(cfg.NetworkID, ".snow") {
		return nil, fmt.Errorf("Expected cfg.NetwrkID=0xnid.snow Got %v", cfg.NetworkID)
	}
	clrpc, err := rpc.Dial(cfg.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "rpc.Dial(%v)", cfg.URL)
	}
	cleth := ethclient.NewClient(clrpc)

	caddr, ok := cfg.ContractAddresses[chain.BTS]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include %v", chain.BTS)
	}
	btscore, err := btscore.NewBtscore(common.HexToAddress(caddr), cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewBtscore")
	}
	caddr, ok = cfg.ContractAddresses[chain.BTSPeriphery]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include %v", chain.BTSPeriphery)
	}
	btsp, err := btsperi.NewBtsperiphery(common.HexToAddress(caddr), cleth)
	if err != nil {
		return nil, errors.Wrap(err, "NewBtscore")
	}
	defaultMapForDifferentGasLimits := map[chain.GasLimitType]uint64{
		chain.TransferNativeCoinIntraChainGasLimit: 5000000,
		chain.TransferTokenIntraChainGasLimit:      5000000,
		chain.ApproveTokenInterChainGasLimit:       5000000,
		chain.TransferCoinInterChainGasLimit:       5000000,
		chain.TransferBatchCoinInterChainGasLimit:  5000000,
		chain.DefaultGasLimit:                      5000000,
	}
	req := &requestAPI{
		contractNameToAddress: cfg.ContractAddresses,
		networkID:             strings.Split(cfg.NetworkID, ".")[0],
		ethCl:                 cleth,
		btsc:                  btscore,
		btsp:                  btsp,
		gasLimits:             cfg.GasLimit,
		nativeCoin:            cfg.NativeCoin,
		nativeTokens:          cfg.NativeTokens,
	}
	req.ercPerCoin, err = req.getCoinAddresses(append(cfg.NativeTokens, cfg.WrappedCoins...))
	for k, v := range defaultMapForDifferentGasLimits {
		if _, ok := req.gasLimits[k]; !ok {
			req.gasLimits[k] = v
		}
	}
	return req, err
}

func (r *requestAPI) getCoinAddresses(inputCoins []string) (ercMap map[string]*erc20tradeable.Erc20tradable, err error) {
	coinNames, err := r.btsc.CoinNames(&bind.CallOpts{Pending: false, Context: nil})
	if err != nil {
		err = errors.Wrap(err, "btsc.CoinId ")
		return
	}
	exists := func(arr []string, val string) bool {
		for _, a := range arr {
			if a == val {
				return true
			}
		}
		return false
	}
	// all registered coins have to be given in input config
	for _, coinName := range coinNames {
		if coinName == r.nativeCoin {
			continue
		}
		if !exists(inputCoins, coinName) {
			err = fmt.Errorf("Registered coin %v not provided in input config ", coinName)
			return
		}
	}
	// all coins given in input config have to have been registered
	for _, inputCoin := range inputCoins {
		if !exists(coinNames, inputCoin) {
			err = fmt.Errorf("Input coin %v does not exist among registered coins ", inputCoin)
			return
		}
	}
	ercMap = map[string]*erc20tradeable.Erc20tradable{}
	for _, coinName := range coinNames {
		if coinName == r.nativeCoin {
			continue
		}
		coinAddress, errs := r.btsc.CoinId(&bind.CallOpts{Pending: false, Context: nil}, coinName)
		if err != nil {
			err = errors.Wrap(errs, "btsc.CoinId ")
			return
		}
		ercMap[coinName], err = erc20tradeable.NewErc20tradable(coinAddress, r.ethCl)
		if err != nil {
			err = errors.Wrap(errs, "NewErc20tradable")
			return
		}
	}
	return
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
	txo.GasPrice.Mul(txo.GasPrice, big.NewInt(3))
	txo.GasLimit = r.gasLimits[chain.DefaultGasLimit] // max gas limit
	return txo, nil
}

func (r *requestAPI) waitForResults(ctx context.Context, txHash common.Hash) (txr *types.Receipt, err error) {
	const DefaultGetTransactionResultPollingInterval = 300 * time.Second //5mins
	ticker := time.NewTicker(time.Second * 2)
	retryLimit := 150
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
			txr, err = r.ethCl.TransactionReceipt(context.Background(), txHash)
			if err != nil && err == ethereum.NotFound {
				//r.log.Trace(errors.Wrap(err, "waitForResults "))
				time.Sleep(time.Second * 2)
				err = nil
				continue
			}
			//r.log.Debugf("GetTransactionResult hash:%v, txr:%+v, err:%+v", txHash, txr, err)
			return
		}
	}
}

func (r *requestAPI) transferIntraChain(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if coinName == r.nativeCoin {
		return r.transferNativeIntraChain(senderKey, recepientAddress, amount)
	}
	return r.transferTokenIntraChain(senderKey, recepientAddress, amount, coinName)
}

func (r *requestAPI) transferNativeIntraChain(senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
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
	gasPrice.Mul(gasPrice, big.NewInt(3))
	chainID, err := r.ethCl.ChainID(context.Background())
	if err != nil {
		err = errors.Wrap(err, "ChainID ")
		return
	}
	tx := types.NewTransaction(nonce, common.HexToAddress(recepientAddress), amount, r.gasLimits[chain.TransferNativeCoinIntraChainGasLimit], gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), senderPrivKey)
	if err != nil {
		err = errors.Wrap(err, "SignTx ")
		return
	}

	if err = r.ethCl.SendTransaction(context.TODO(), signedTx); err != nil {
		err = errors.Wrap(err, "SendNativeTransaction ")
		return
	}
	txnHash = signedTx.Hash().String()
	return
}

func (r *requestAPI) transferTokenIntraChain(senderKey, recepientAddress string, amount *big.Int, coinName string) (txnHash string, err error) {
	erc, ok := r.ercPerCoin[coinName]
	if !ok {
		err = fmt.Errorf("coin %v not registered", coinName)
		return
	}

	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txo.GasLimit = r.gasLimits[chain.TransferTokenIntraChainGasLimit]
	txn, err := erc.Transfer(txo, common.HexToAddress(recepientAddress), amount)
	if err != nil {
		err = errors.Wrap(err, "hrc.Transfer ")
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) transferInterChain(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if coinName == r.nativeCoin {
		return r.transferNativeCrossChain(senderKey, recepientAddress, amount)
	}
	return r.transferTokensCrossChain(coinName, senderKey, recepientAddress, amount)
}

func (r *requestAPI) transferNativeCrossChain(senderKey string, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Value = amount
	txo.Context = context.Background()
	txo.GasLimit = r.gasLimits[chain.TransferCoinInterChainGasLimit]
	txn, err := r.btsc.TransferNativeCoin(txo, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "btsc.TransferNativeCoin ")
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) transferTokensCrossChain(coinName string, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	_, ok := r.ercPerCoin[coinName]
	if !ok {
		err = fmt.Errorf("coin %v not registered", coinName)
		return
	}
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txo.GasLimit = r.gasLimits[chain.TransferCoinInterChainGasLimit]
	txn, err := r.btsc.Transfer(txo, coinName, amount, recepientAddress)
	if err != nil {
		err = errors.Wrap(err, "btsc.Transfer ")
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) transferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
	if len(amounts) != len(coinNames) {
		return "", fmt.Errorf("Amount and CoinNames len should be same; Got %v and %v", len(amounts), len(coinNames))
	}
	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txo.GasLimit = r.gasLimits[chain.TransferBatchCoinInterChainGasLimit]
	filterNames := []string{}
	filterAmounts := []*big.Int{}
	for i := 0; i < len(amounts); i++ {
		if coinNames[i] == r.nativeCoin {
			txo.Value = amounts[i]
			continue
		} else if _, ok := r.ercPerCoin[coinNames[i]]; !ok {
			err = fmt.Errorf("coin %v not registered", coinNames[i])
			return
		}
		filterAmounts = append(filterAmounts, amounts[i])
		filterNames = append(filterNames, coinNames[i])
	}
	txn, err := r.btsc.TransferBatch(txo, filterNames, filterAmounts, recepientAddress)
	if err != nil {
		err = errors.Wrapf(err, "btsc.TransferBatch %v", err)
		return
	}
	txnHash = txn.Hash().String()
	return
}

func (r *requestAPI) approveCoin(coinName, senderKey string, amount *big.Int) (approveTxnHash string, err error) {
	if coinName == r.nativeCoin {
		err = fmt.Errorf("Native Coin %v does not need to be approved", coinName)
		return
	}
	erc, ok := r.ercPerCoin[coinName]
	if !ok {
		err = fmt.Errorf("coin %v not registered", coinName)
		return
	}

	txo, err := r.getTransactionRequest(senderKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.GasLimit = r.gasLimits[chain.ApproveTokenInterChainGasLimit]
	btscaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include %v ", chain.BTS)
		return
	}
	txo.Context = context.Background()
	approveTxn, err := erc.Approve(txo, common.HexToAddress(btscaddr), amount)
	if err != nil {
		err = errors.Wrap(err, "erc.Approve ")
		return
	}
	approveTxnHash = approveTxn.Hash().String()
	return
}

func (r *requestAPI) setFeeRatio(ownerKey string, coinName string, feeNumerator, fixedFee *big.Int) (hash string, err error) {
	txo, err := r.getTransactionRequest(ownerKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	txo.Context = context.Background()
	txn, err := r.btsc.SetFeeRatio(txo, coinName, feeNumerator, fixedFee)
	if err != nil {
		err = errors.Wrap(err, "setFeeRatio ")
		return
	}
	hash = txn.Hash().String()
	return
}

func (r *requestAPI) getCoinBalance(coinName, addr string) (bal *chain.CoinBalance, err error) {
	b, err := r.btsc.BalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, common.HexToAddress(addr), coinName)
	if err != nil {
		err = errors.Wrap(err, "btsc.GetBalanceOf ")
		return
	}
	bal = &chain.CoinBalance{
		UsableBalance:     b.UsableBalance,
		LockedBalance:     b.LockedBalance,
		RefundableBalance: b.RefundableBalance,
		UserBalance:       b.UserBalance,
		TotalBalance:      b.UserBalance,
	}
	return bal, nil
}

func (r *requestAPI) reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	txo, err := r.getTransactionRequest(ownerKey)
	if err != nil {
		err = errors.Wrap(err, "getTransactionRequest ")
		return
	}
	// Add GasLimit
	txn, err := r.btsc.Reclaim(txo, coinName, amount)
	txnHash = txn.Hash().String()
	return
}

func (r *api) SetTokenLimit(ownerKey string, coinNames []string, tokenLimits []*big.Int) (txnHash string, err error) {
	return "", errors.New("not implemented")
}
func (r *api) AddBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	return "", errors.New("not implemented")
}

func (r *api) RemoveBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	return "", errors.New("not implemented")
}

func (r *api) ChangeRestriction(ownerKey string, enable bool) (txnHash string, err error) {
	return "", errors.New("not implemented")
}

func (r *api) IsUserBlackListed(net, addr string) (response bool, err error) {
	if strings.Contains(addr, "btp:") {
		splts := strings.Split(addr, "/")
		addr = splts[len(splts)-1]
	}
	response, err = r.requester.btsp.Blacklist(&bind.CallOpts{Pending: false, Context: context.TODO()}, ethCommon.HexToAddress(addr))
	return
}

func (r *api) GetTokenLimit(coinName string) (tokenLimit *big.Int, err error) {
	tokenLimit, err = r.requester.btsp.TokenLimit(&bind.CallOpts{Pending: false, Context: context.TODO()}, coinName)
	return
}

func (r *api) IsBTSOwner(addr string) (response bool, err error) {
	if strings.Contains(addr, "btp:") {
		splts := strings.Split(addr, "/")
		addr = splts[len(splts)-1]
	}
	response, err = r.requester.btsc.IsOwner(&bind.CallOpts{Pending: false, Context: context.TODO()}, ethCommon.HexToAddress(addr))
	return
}

func (r *api) GetTokenLimitStatus(net, coinName string) (response bool, err error) {
	return false, errors.New("not implemented")
}

func (r *api) GetBlackListedUsers(net string, startCursor, endCursor int) (addrs []string, err error) {
	return nil, errors.New("not implemented")
}
