package icon

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"

	"github.com/haltingstate/secp256k1-go"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gocrypto "github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/icon-project/icon-bridge/common/intconv"
)

const (
	DefaultSendTransactionRetryInterval        = 3 * time.Second         //3sec
	DefaultGetTransactionResultPollingInterval = 1500 * time.Millisecond //1.5sec
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	cl                    *icon.Client
	nativeCoin            string
	wrappedCoinsAddr      map[string]string
	nativeTokensAddr      map[string]string
	gasLimits             map[chain.GasLimitType]uint64
}

func newRequestAPI(cl *icon.Client, cfg *chain.Config) (req *requestAPI, err error) {
	var defaultMapForDifferentGasLimits = map[chain.GasLimitType]uint64{
		chain.TransferNativeCoinIntraChainGasLimit: 150000,
		chain.TransferTokenIntraChainGasLimit:      300000,
		chain.ApproveTokenInterChainGasLimit:       800000,
		chain.TransferCoinInterChainGasLimit:       2500000,
		chain.TransferBatchCoinInterChainGasLimit:  4000000,
		chain.DefaultGasLimit:                      5000000,
	}

	if !strings.Contains(cfg.NetworkID, ".icon") {
		return nil, fmt.Errorf("Expected cfg.NetwrkID=0xnid.icon Got %v", cfg.NetworkID)
	}
	req = &requestAPI{
		networkID:             strings.Split(cfg.NetworkID, ".")[0],
		contractNameToAddress: cfg.ContractAddresses,
		cl:                    cl,
		nativeCoin:            cfg.NativeCoin,
		gasLimits:             cfg.GasLimit,
	}
	req.nativeTokensAddr, req.wrappedCoinsAddr, err = req.getCoinAddresses(cfg.NativeTokens, cfg.WrappedCoins)
	for k, v := range defaultMapForDifferentGasLimits {
		if _, ok := req.gasLimits[k]; !ok {
			req.gasLimits[k] = v
		}
	}
	return req, err
}

func (r *requestAPI) transactWithContract(senderKey string, contractAddress string,
	amount *big.Int, args map[string]interface{}, method string, stepLimit int64) (txHash string, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}

	param := types.TransactionParam{
		Version:     types.NewHexInt(types.JsonrpcApiVersion),
		ToAddress:   types.Address(contractAddress),
		Value:       types.HexInt(intconv.FormatBigInt(amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: types.Address(senderWallet.Address().String()),
		StepLimit:   types.NewHexInt(stepLimit),
		Timestamp:   types.NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   types.HexInt(r.networkID),
		DataType:    "call",
	}
	argMap := map[string]interface{}{}
	argMap["method"] = method
	argMap["params"] = args
	param.Data = argMap

	if err = SignTransactionParam(senderWallet, &param); err != nil {
		err = errors.Wrap(err, "SignTransactionParam ")
		return
	}
	txH, err := r.cl.SendTransaction(&param)
	if err != nil {
		err = errors.Wrap(err, "SendTransaction ")
		return
	}
	txBytes, err := txH.Value()
	if err != nil {
		err = errors.Wrap(err, "HexBytes.Value() ")
		return
	}
	txHash = hexutil.Encode(txBytes[:])
	return
}

func (r *requestAPI) callContract(contractAddress string, args map[string]interface{}, method string) (interface{}, error) {
	param := &types.CallParam{
		ToAddress: types.Address(contractAddress),
		DataType:  "call",
	}
	argMap := map[string]interface{}{}
	argMap["method"] = method
	argMap["params"] = args
	param.Data = argMap
	var res interface{}
	err := r.cl.Call(param, &res)
	if err != nil {
		return nil, errors.Wrap(err, "Call ")
	}
	return res, nil
}

func (r *requestAPI) transferIntraChain(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if coinName == r.nativeCoin {
		txnHash, err = r.transferNativeIntraChain(senderKey, recepientAddress, amount)
	} else if caddr, ok := r.nativeTokensAddr[coinName]; ok {
		txnHash, err = r.transferTokenIntraChain(senderKey, recepientAddress, amount, caddr)
	} else if _, ok := r.wrappedCoinsAddr[coinName]; ok {
		err = fmt.Errorf("IntraChain transfers not supported for wrapped coins. Got %v", coinName)
	} else {
		err = fmt.Errorf("Coin %v not amongst registered coins", coinName)
	}
	return
}

func (r *requestAPI) transferTokenIntraChain(senderKey, recepientAddress string, amount *big.Int, caddr string) (txHash string, err error) {
	args := map[string]interface{}{"_to": recepientAddress, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(senderKey, caddr, big.NewInt(0), args, "transfer", int64(r.gasLimits[chain.TransferTokenIntraChainGasLimit]))
}

func (r *requestAPI) transferNativeIntraChain(senderKey, recepientAddress string, amount *big.Int) (txHash string, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}
	param := types.TransactionParam{
		Version:     types.NewHexInt(types.JsonrpcApiVersion),
		ToAddress:   types.Address(recepientAddress),
		Value:       types.HexInt(intconv.FormatBigInt(amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: types.Address(senderWallet.Address().String()),
		StepLimit:   types.NewHexInt(int64(r.gasLimits[chain.TransferNativeCoinIntraChainGasLimit])),
		Timestamp:   types.NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   types.HexInt(r.networkID),
	}
	if err = SignTransactionParam(senderWallet, &param); err != nil {
		err = errors.Wrap(err, "SignTransactionParam ")
		return
	}
	txH, err := r.cl.SendTransaction(&param)
	if err != nil {
		err = errors.Wrap(err, "SendTransaction ")
		return
	}
	txBytes, err := txH.Value()
	if err != nil {
		err = errors.Wrap(err, "HexBytes.Value() ")
		return
	}
	txHash = hexutil.Encode(txBytes[:])
	return
}
func (r *requestAPI) transferInterChain(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if coinName == r.nativeCoin {
		txnHash, err = r.transferNativeCrossChain(senderKey, recepientAddress, amount)
	} else {
		_, tok := r.nativeTokensAddr[coinName]
		_, wok := r.wrappedCoinsAddr[coinName]
		if wok || tok {
			txnHash, err = r.transferTokensCrossChain(coinName, senderKey, recepientAddress, amount)
		} else {
			err = fmt.Errorf("Coin %v not among registered ", coinName)
		}
	}
	return
}

func (r *requestAPI) transferNativeCrossChain(senderKey, recepientAddress string, amount *big.Int) (txHash string, err error) {
	args := map[string]interface{}{"_to": recepientAddress}
	caddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	return r.transactWithContract(senderKey, caddr, amount, args, "transferNativeCoin", int64(r.gasLimits[chain.TransferCoinInterChainGasLimit]))
}

func (r *requestAPI) transferTokensCrossChain(coinName, senderKey, recepientAddress string, amount *big.Int) (string, error) {
	args := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount), "_to": recepientAddress}
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}
	return r.transactWithContract(senderKey, btsaddr, big.NewInt(0), args, "transfer", int64(r.gasLimits[chain.TransferCoinInterChainGasLimit]))
}

func (r *requestAPI) approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	if addr, ok := r.nativeTokensAddr[coinName]; ok {
		txnHash, err = r.approveToken(coinName, ownerKey, amount, addr)
	} else if coinName == r.nativeCoin {
		err = fmt.Errorf("Native Coin %v does not need to be approved", coinName)
	} else if addr, ok := r.wrappedCoinsAddr[coinName]; ok {
		txnHash, err = r.approveCrossNativeCoin(coinName, ownerKey, amount, addr)
	} else {
		err = fmt.Errorf("Coin %v not amongst registered coins", coinName)
	}
	return
}

func (r *requestAPI) approveCrossNativeCoin(coinName string, ownerKey string, amount *big.Int, coinAddress string) (approveTxnHash string, err error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	approveArgs := map[string]interface{}{"spender": btsaddr, "amount": intconv.FormatBigInt(amount)}
	approveTxnHash, err = r.transactWithContract(ownerKey, coinAddress, big.NewInt(0), approveArgs, "approve", int64(r.gasLimits[chain.ApproveTokenInterChainGasLimit]))
	if err != nil {
		err = errors.Wrapf(err, "transactWithContract %v", coinAddress)
		return
	}
	return
}

func (r *requestAPI) approveToken(coinName, senderKey string, amount *big.Int, caddr string) (hash string, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}
	arg1 := map[string]interface{}{"_to": btsAddr, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(senderKey, caddr, big.NewInt(0), arg1, "transfer", int64(r.gasLimits[chain.ApproveTokenInterChainGasLimit]))
}

func (r *requestAPI) transferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
	if len(amounts) != len(coinNames) {
		return "", fmt.Errorf("Amount and CoinNames len should be same; Got %v and %v", len(amounts), len(coinNames))
	}
	nativeAmount := big.NewInt(0)
	filterNames := []string{}
	filterAmounts := []string{}
	for i := 0; i < len(amounts); i++ {
		if coinNames[i] == r.nativeCoin {
			nativeAmount = amounts[i]
			continue
		}
		filterAmounts = append(filterAmounts, intconv.FormatBigInt(amounts[i]))
		filterNames = append(filterNames, coinNames[i])
	}

	args := map[string]interface{}{"_coinNames": filterNames, "_values": filterAmounts, "_to": recepientAddress}
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}

	txnHash, err = r.transactWithContract(senderKey, btsaddr, nativeAmount, args, "transferBatch", int64(r.gasLimits[chain.TransferBatchCoinInterChainGasLimit]))
	return
}

func (r *requestAPI) getCoinBalance(coinName, addr string) (bal *chain.CoinBalance, err error) {
	if coinName == r.nativeCoin {
		return r.getNativeCoinBalance(coinName, addr)
	}
	getBalanceOfType := func(balanceMap map[string]interface{}, key string) (bal *big.Int, err error) {
		tmp, ok := balanceMap[key]
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

	// Tokens ..
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}

	// BTS BALANCEOF
	res, err := r.callContract(btsAddr, map[string]interface{}{"_coinName": coinName, "_owner": addr}, "balanceOf")
	if err != nil {
		return nil, errors.Wrap(err, "callContract balanceOf ")
	} else if res == nil {
		return nil, errors.New("callContract returned nil value ")
	}
	balanceMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected type map[string]interface{} Got %T", res)
	}
	bal = &chain.CoinBalance{}
	bal.UsableBalance, err = getBalanceOfType(balanceMap, "usable")
	if err != nil {
		return nil, errors.Wrapf(err, "getBalanceOfType(usable) %v", err)
	}
	bal.LockedBalance, err = getBalanceOfType(balanceMap, "locked")
	if err != nil {
		return nil, errors.Wrapf(err, "getBalanceOfType(locked) %v", err)
	}
	bal.RefundableBalance, err = getBalanceOfType(balanceMap, "refundable")
	if err != nil {
		return nil, errors.Wrapf(err, "getBalanceOfType(refundable) %v", err)
	}
	bal.UserBalance, err = getBalanceOfType(balanceMap, "userBalance")
	if err != nil {
		return nil, errors.Wrapf(err, "getBalanceOfType(userBalance) %v", err)
	}
	bal.TotalBalance = (&big.Int{}).Add(bal.UserBalance, bal.UsableBalance)
	return
}

func (r *requestAPI) getCoinAddresses(nativeTokens, wrappedCoins []string) (tokenAddrMap, wrappedAddrMap map[string]string, err error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := r.callContract(btsaddr, map[string]interface{}{}, "coinNames")
	if err != nil {
		err = errors.Wrap(err, "callContract coinNames ")
		return
	} else if res == nil {
		err = fmt.Errorf("Call to Method %v returned nil", "coinNames")
		return
	}
	resArr, ok := res.([]interface{})
	if !ok {
		err = fmt.Errorf("For method coinNames, Expected Type []interface{} Got %T", resArr)
		return
	}
	coinNames := []string{}
	for _, re := range resArr {
		c, ok := re.(string)
		if !ok {
			err = fmt.Errorf("Expected Type string Got %T", re)
			return
		}
		if c == r.nativeCoin {
			continue
		}
		coinNames = append(coinNames, c)
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
	allInputCoins := append(nativeTokens, wrappedCoins...)
	for _, coinName := range coinNames {
		if !exists(allInputCoins, coinName) {
			err = fmt.Errorf("Registered coin %v not provided in input config ", coinName)
			return
		}
	}
	// all coins given in input config have to have been registered
	for _, inputCoin := range allInputCoins {
		if !exists(coinNames, inputCoin) {
			err = fmt.Errorf("Input coin %v does not exist among registered coins ", inputCoin)
			return
		}
	}
	getAddr := func(coin string) (coinId string, err error) {
		var res interface{}
		res, err = r.callContract(btsaddr, map[string]interface{}{"_coinName": coin}, "coinId")
		if err != nil {
			err = errors.Wrap(err, "callContract coinId ")
			return
		} else if res == nil {
			err = fmt.Errorf("Call to Method %v returned nil for _coinName=%v", "coinId", coin)
			return
		}
		coinId, ok := res.(string)
		if !ok {
			err = fmt.Errorf("For method coinId, Expected Type string Got %T", res)
			return
		}
		return coinId, nil
	}

	tokenAddrMap = map[string]string{}
	for _, coin := range nativeTokens {
		tokenAddrMap[coin], err = getAddr(coin)
		if err != nil {
			return
		}
	}
	wrappedAddrMap = map[string]string{}
	for _, coin := range wrappedCoins {
		wrappedAddrMap[coin], err = getAddr(coin)
		if err != nil {
			return
		}
	}
	return
}

func generateKeyPair() ([2]string, error) {
	pubkeyBytes, priv := secp256k1.GenerateKeyPair()
	pubKey, err := crypto.ParsePublicKey(pubkeyBytes)
	if err != nil {
		return [2]string{}, errors.Wrap(err, "crypto.ParsePublicKey ")
	}
	addr := common.NewAccountAddressFromPublicKey(pubKey).String()
	return [2]string{hex.EncodeToString(priv), addr}, nil
}

func (r *requestAPI) getNativeCoinBalance(coinName, addr string) (bal *chain.CoinBalance, err error) {
	zeroBalance := big.NewInt(0)
	bal = &chain.CoinBalance{UsableBalance: zeroBalance, RefundableBalance: zeroBalance, LockedBalance: zeroBalance, UserBalance: new(big.Int)}
	// Native
	bal.UserBalance, err = r.cl.GetBalance(&types.AddressParam{Address: types.Address(addr)})
	if err != nil {
		return nil, errors.Wrapf(err, "%v", err)
	}
	bal.TotalBalance = (&big.Int{}).Set(bal.UserBalance)
	return
}

func SignTransactionParam(wallet module.Wallet, param *types.TransactionParam) error {
	js, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, "jsonMarshal ")
	}
	var txSerializeExcludes = map[string]bool{"signature": true}
	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		return errors.Wrap(err, "tx.SerializeJSON ")
	}
	bs = append([]byte("icx_sendTransaction."), bs...)
	sig, err := wallet.Sign(gocrypto.SHA3Sum256(bs))
	if err != nil {
		return errors.Wrap(err, "wallet.Sign ")
	}
	param.Signature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func GetWalletFromPrivKey(privKey string) (module.Wallet, error) {
	privBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "DecodeString ")
	}
	pKey, err := gocrypto.ParsePrivateKey(privBytes)
	if err != nil {
		return nil, errors.Wrap(err, "crypto.ParsePrivateKey ")
	}
	wal, err := wallet.NewFromPrivateKey(pKey)
	if err != nil {
		return nil, errors.Wrap(err, "crypto.NewFromPrivateKey ")
	}
	return wal, nil
}

func (r *requestAPI) reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}
	arg1 := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(ownerKey, btsAddr, big.NewInt(0), arg1, "reclaim", int64(r.gasLimits[chain.DefaultGasLimit]))
}

// Configure API
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
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_coinNames": coinNames, "_tokenLimits": strTokenLimits}, "setTokenLimit", int64(a.requester.gasLimits[chain.DefaultGasLimit]))
}

func (a *api) AddBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_net": net, "_addresses": addrs}, "addBlacklistAddress", int64(a.requester.gasLimits[chain.DefaultGasLimit]))
}

func (a *api) RemoveBlackListAddress(ownerKey string, net string, addrs []string) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_net": net, "_addresses": addrs}, "removeBlacklistAddress", int64(a.requester.gasLimits[chain.DefaultGasLimit]))
}

func (a *api) ChangeRestriction(ownerKey string, enable bool) (txnHash string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	if enable {
		return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{}, "addRestriction", int64(a.requester.gasLimits[chain.DefaultGasLimit]))
	}
	return a.requester.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{}, "disableRestrictions", int64(a.requester.gasLimits[chain.DefaultGasLimit]))
}

func (r *requestAPI) setFeeGatheringTerm(ownerKey string, interval uint64) (hash string, err error) {
	bmcAddr, ok := r.contractNameToAddress[chain.BMC]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BMC)
		return
	}
	return r.transactWithContract(ownerKey, bmcAddr, big.NewInt(0), map[string]interface{}{"_value": hexutil.EncodeUint64(interval)}, "setFeeGatheringTerm", int64(r.gasLimits[chain.DefaultGasLimit]))
}

func (r *requestAPI) getFeeGatheringTerm() (interval uint64, err error) {
	bmcAddr, ok := r.contractNameToAddress[chain.BMC]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BMC)
		return
	}
	res, err := r.callContract(bmcAddr, map[string]interface{}{}, "getFeeGatheringTerm")
	if err != nil {
		return 0, errors.Wrap(err, "callContract getFeeGatheringTerm ")
	} else if res == nil {
		return 0, errors.New("callContract getFeeGatheringTerm returned nil value ")
	}
	tmpStr, ok := res.(string)
	if !ok {
		return 0, fmt.Errorf("Expected type string Got %T", res)
	}
	return hexutil.DecodeUint64(tmpStr)
}

func (r *requestAPI) setFeeRatio(ownerKey string, coinName string, feeNumerator, fixedFee *big.Int) (hash string, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	_feeNumerator := intconv.FormatBigInt(feeNumerator)
	_fixedFee := intconv.FormatBigInt(fixedFee)
	return r.transactWithContract(ownerKey, btsAddr, big.NewInt(0), map[string]interface{}{"_name": coinName, "_feeNumerator": _feeNumerator, "_fixedFee": _fixedFee}, "setFeeRatio", int64(r.gasLimits[chain.DefaultGasLimit]))
}

func (r *requestAPI) getAccumulatedFees() (ret map[string]*big.Int, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := r.callContract(btsAddr, map[string]interface{}{}, "getAccumulatedFees")
	if err != nil {
		return nil, errors.Wrap(err, "callContract getAccumulatedFees ")
	} else if res == nil {
		return nil, errors.New("callContract getAccumulatedFees returned nil value ")
	}
	resMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected type map[string]interface{} Got %T", res)
	}
	ret = map[string]*big.Int{}
	for k, v := range resMap {
		tmpStr, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("Expected type string Got %T", v)
		}
		bal := new(big.Int)
		bal.SetString(tmpStr[2:], 16)
		ret[k] = bal
	}
	return
}

func (r *requestAPI) getFeeRatio(coinName string) (feeNumerator, fixedFee *big.Int, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := r.callContract(btsAddr, map[string]interface{}{"_name": coinName}, "feeRatio")
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

func (a *api) IsUserBlackListed(net, addr string) (response bool, err error) {
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

func (a *api) GetTokenLimit(coinName string) (tokenLimit *big.Int, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_name": coinName}, "getTokenLimit")
	if err != nil {
		err = errors.Wrapf(err, "CallContract %v", err)
		return
	} else if res == nil {
		err = errors.New("getTokenLimit result is nil")
		return
	}
	tmpStr, ok := res.(string)
	if !ok {
		err = fmt.Errorf("Expected type string Got %T", res)
		return
	}
	tokenLimit = new(big.Int)
	tokenLimit.SetString(tmpStr[2:], 16)
	return
}

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

func (a *api) GetTokenLimitStatus(net, coinName string) (response bool, err error) {
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

func (a *api) GetBlackListedUsers(net string, startCursor, endCursor int) (users []string, err error) {
	btsAddr, ok := a.requester.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := a.requester.callContract(btsAddr, map[string]interface{}{"_net": net, "_start": "0x0", "_end": "0x64"}, "getBlackListedUsers")
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
