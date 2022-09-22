package near

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haltingstate/secp256k1-go"
	gocrypto "github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/intconv"
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	cl                    near.IClient
	stepLimit             int64
	nativeCoin            string
	wrappedCoinsAddr      map[string]string
	nativeTokensAddr      map[string]string
}

type Metadata struct {
	CoinMetadata CoinMetadata `json:"metadata,omitempty"`
}

type CoinMetadata struct {
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Uri          string `json:"uri"`
	Network      string `json:"network"`
	FeeNumerator string `json:"fee_numerator"`
	FixedFee     string `json:"fixed_fee"`
	Extras       Extras `json:"extras"`
}

type Extras struct {
	Spec          string      `json:"spec"`
	Icon          interface{} `json:"icon,omitempty"`
	Refrence      interface{} `json:"refrence,omitempty"`
	ReferenceHash interface{} `json:"reference_hash,omitempty"`
	Decimals      int64       `json:"decimals"`
}

type coinNames struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Network string `json:"network"`
}

func newRequestAPI(cl near.IClient, cfg *chain.Config) (req *requestAPI, err error) {
	if !strings.Contains(cfg.NetworkID, ".near") {
		return nil, fmt.Errorf("expected cfg.NetwrkID=0xnid.near Got %v", cfg.NetworkID)
	}
	req = &requestAPI{
		networkID:             strings.Split(cfg.NetworkID, ".")[0],
		contractNameToAddress: cfg.ContractAddresses,
		cl:                    cl,
		stepLimit:             cfg.GasLimit,
		nativeCoin:            cfg.NativeCoin,
	}
	req.nativeTokensAddr, req.wrappedCoinsAddr, err = req.getCoinAddresses(cfg.NativeTokens, cfg.WrappedCoins)
	return req, err
}

func (r *requestAPI) getCoinAddresses(nativeTokens, wrappedCoins []string) (tokenAddrMap, wrappedAddrMap map[string]string, err error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	var coin_names []coinNames
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := r.callContract(btsaddr, map[string]interface{}{}, "coins")
	if err != nil {
		err = errors.Wrap(err, "callContract coinNames ")
		return
	} else if res == nil {
		err = fmt.Errorf("call to Method %v returned nil", "coinNames")
		return
	}

	resArr := res.(types.CallFunctionResponse).Result
	err = json.Unmarshal(resArr, &coin_names)
	if err != nil {
		err = fmt.Errorf("for method coinNames, Expected Type []interface{} Got %T", err)
		return
	}
	coinNames := []string{}
	for _, re := range coin_names {
		c := re.Name
		if c == r.nativeCoin {
			continue
		}
		coinNames = append(coinNames, c)
	}
	println(coinNames)
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
			err = fmt.Errorf("registered coin %v not provided in input config ", coinName)
			return
		}
	}
	// all coins given in input config have to have been registered
	for _, inputCoin := range allInputCoins {
		if !exists(coinNames, inputCoin) {
			err = fmt.Errorf("input coin %v does not exist among registered coins ", inputCoin)
			return
		}
	}
	getAddr := func(coin string) (coinId string, err error) {
		var res interface{}

		res, err = r.callContract(btsaddr, map[string]interface{}{"coin_name": coin}, "coin")
		if err != nil {
			err = errors.Wrap(err, "callContract coin ")
			return
		} else if res == nil {
			err = fmt.Errorf("call to Method %v returned nil for _coinName=%v", "coin", coin)
			return
		}
		resArr := res.(types.CallFunctionResponse).Result
		var coin_metadata Metadata
		err = json.Unmarshal(resArr, &coin_metadata)
		if err != nil {
			err = errors.Wrap(err, "callContract coin ")
			return
		}
		coinId = base64.StdEncoding.EncodeToString([]byte(coin_metadata.CoinMetadata.Uri))
		if !ok {
			err = fmt.Errorf("for method coinId, Expected Type string Got %T", res)
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

func (r *requestAPI) callContract(contractAddress string, args map[string]interface{}, method string) (interface{}, error) {
	methodParam, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	param := &types.CallFunction{
		RequestType:  "call_function",
		Finality:     "final",
		AccountId:    types.AccountId(contractAddress),
		MethodName:   method,
		ArgumentsB64: base64.URLEncoding.EncodeToString(methodParam),
	}

	var res types.CallFunctionResponse
	res, err = r.cl.Api().CallFunction(param)
	if err != nil {
		return nil, errors.Wrap(err, "Call ")
	}
	return res, nil
}

func (r *requestAPI) approve(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	if addr, ok := r.nativeTokensAddr[coinName]; ok {
		txnHash, err = r.approveToken(coinName, ownerKey, amount, addr)
	} else if coinName == r.nativeCoin {
		err = fmt.Errorf("native Coin %v does not need to be approved", coinName)
	} else if addr, ok := r.wrappedCoinsAddr[coinName]; ok {
		txnHash, err = r.approveCrossNativeCoin(coinName, ownerKey, amount, addr)
	} else {
		err = fmt.Errorf("coin %v not amongst registered coins", coinName)
	}
	return
}

func (r *requestAPI) approveToken(coinName, senderKey string, amount *big.Int, caddr string) (hash string, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}
	arg1 := map[string]interface{}{"_to": btsAddr, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(senderKey, caddr, big.NewInt(0), arg1, "transfer")
}

func (r *requestAPI) approveCrossNativeCoin(coinName string, ownerKey string, amount *big.Int, coinAddress string) (approveTxnHash string, err error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	approveArgs := map[string]interface{}{"spender": btsaddr, "amount": intconv.FormatBigInt(amount)}
	approveTxnHash, err = r.transactWithContract(ownerKey, coinAddress, big.NewInt(0), approveArgs, "approve")
	if err != nil {
		err = errors.Wrapf(err, "transactWithContract %v", coinAddress)
		return
	}
	return
}

func SignTransactionParam(wallet module.Wallet, param *types.Transaction) error {
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
	var sign types.Signature
	err = json.Unmarshal(sig, &sign)
	if err != nil {
		fmt.Println(err)
	}
	param.Signature = sign
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
			return nil, fmt.Errorf("expected type string Got %T", tmp)
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
		return nil, fmt.Errorf("expected type map[string]interface{} Got %T", res)
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
	return
}

func (r *requestAPI) getNativeCoinBalance(coinName, addr string) (bal *chain.CoinBalance, err error) {
	zeroBalance := big.NewInt(0)
	bal = &chain.CoinBalance{UsableBalance: zeroBalance, RefundableBalance: zeroBalance, LockedBalance: zeroBalance, UserBalance: new(big.Int)}
	// Native
	bal.UserBalance, err = r.cl.GetBalance(types.AccountId(addr))
	if err != nil {
		return nil, errors.Wrapf(err, "%v", err)
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

func (r *requestAPI) reclaim(coinName string, ownerKey string, amount *big.Int) (txnHash string, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}
	arg1 := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(ownerKey, btsAddr, big.NewInt(0), arg1, "reclaim")
}

func (r *requestAPI) transferIntraChain(coinName, senderKey, recepientAddress string, amount *big.Int) (txnHash string, err error) {
	if coinName == r.nativeCoin {
		txnHash, err = r.transferNativeIntraChain(senderKey, recepientAddress, amount)
	} else if caddr, ok := r.nativeTokensAddr[coinName]; ok {
		txnHash, err = r.transferTokenIntraChain(senderKey, recepientAddress, amount, caddr)
	} else if _, ok := r.wrappedCoinsAddr[coinName]; ok {
		err = fmt.Errorf("IntraChain transfers not supported for wrapped coins. Got %v", coinName)
	} else {
		err = fmt.Errorf("coin %v not amongst registered coins", coinName)
	}
	return
}

func (r *requestAPI) transferTokenIntraChain(senderKey, recepientAddress string, amount *big.Int, caddr string) (txHash string, err error) {
	args := map[string]interface{}{"_to": recepientAddress, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(senderKey, caddr, big.NewInt(0), args, "transfer")
}

func (r *requestAPI) transferNativeIntraChain(senderKey, recepientAddress string, amount *big.Int) (txHash string, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}

	nonce, err := r.cl.GetNonce(types.NewPublicKeyFromED25519(senderWallet.PublicKey()), senderWallet.Address().String())
	if err != nil {
		fmt.Println(err)
	}

	blockHash, err := r.cl.GetLatestBlockHash()
	if err != nil {
		fmt.Println(err)
	}
	actions := []types.Action{
		{
			Enum: 2,
			FunctionCall: types.FunctionCall{
				// MethodName: method,
				// Args:    data,
				Gas:     uint64(r.stepLimit),
				Deposit: *big.NewInt(0),
			},
		},
	}
	param := types.Transaction{
		SignerId:   types.AccountId(senderWallet.Address().ID()),
		PublicKey:  types.NewPublicKeyFromED25519(senderWallet.PublicKey()),
		Nonce:      int(nonce),
		ReceiverId: types.AccountId(recepientAddress),
		BlockHash:  blockHash,
		Actions:    actions,
		Txid:       blockHash,
	}
	if err = SignTransactionParam(senderWallet, &param); err != nil {
		err = errors.Wrap(err, "SignTransactionParam ")
		return
	}
	txH, err := r.cl.SendTransaction(param.Signature.Base58Encode())
	if err != nil {
		err = errors.Wrap(err, "SendTransaction ")
		return
	}
	txHash = hexutil.Encode(txH[:])
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
			err = fmt.Errorf("coin %v not among registered ", coinName)
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
	return r.transactWithContract(senderKey, caddr, amount, args, "transferNativeCoin")
}

func (r *requestAPI) transferTokensCrossChain(coinName, senderKey, recepientAddress string, amount *big.Int) (string, error) {
	args := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount), "_to": recepientAddress}
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
	}
	return r.transactWithContract(senderKey, btsaddr, big.NewInt(0), args, "transfer")
}

func (r *requestAPI) transferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int) (txnHash string, err error) {
	if len(amounts) != len(coinNames) {
		return "", fmt.Errorf("amount and CoinNames len should be same; Got %v and %v", len(amounts), len(coinNames))
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

	txnHash, err = r.transactWithContract(senderKey, btsaddr, nativeAmount, args, "transferBatch")
	return
}

func (r *requestAPI) transactWithContract(senderKey string, contractAddress string,
	amount *big.Int, args map[string]interface{}, method string) (txHash string, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}
	data, err := json.Marshal(args)
	if err != nil {
		return "", err
	}

	nonce, err := r.cl.GetNonce(types.NewPublicKeyFromED25519(senderWallet.PublicKey()), senderWallet.Address().String())
	if err != nil {
		fmt.Println(err)
	}
	blockHash, err := r.cl.GetLatestBlockHash()
	if err != nil {
		fmt.Println(err)
	}
	actions := []types.Action{
		{
			Enum: 2,
			FunctionCall: types.FunctionCall{
				MethodName: method,
				Args:       data,
				Gas:        uint64(r.stepLimit),
				Deposit:    *big.NewInt(0),
			},
		},
	}
	param := types.Transaction{
		SignerId:   types.AccountId(senderWallet.Address().ID()),
		PublicKey:  types.NewPublicKeyFromED25519(senderWallet.PublicKey()),
		Nonce:      int(nonce),
		ReceiverId: types.AccountId(contractAddress),
		BlockHash:  blockHash,
		Actions:    actions,
		Txid:       blockHash,
	}

	if err = SignTransactionParam(senderWallet, &param); err != nil {
		err = errors.Wrap(err, "SignTransactionParam ")
		return
	}
	txH, err := r.cl.SendTransaction(param.Signature.Base58Encode())
	if err != nil {
		err = errors.Wrap(err, "SendTransaction ")
		return
	}

	txHash = hexutil.Encode(txH[:])
	return
}