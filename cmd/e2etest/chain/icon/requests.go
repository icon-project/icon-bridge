package icon

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

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
	StepLimit                                  = 35000000
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	cl                    *icon.Client
}

func newRequestAPI(cl *icon.Client, contractNameToAddress map[chain.ContractName]string, networkID string) (*requestAPI, error) {
	return &requestAPI{networkID: networkID, contractNameToAddress: contractNameToAddress, cl: cl}, nil
}

func SignTransactionParam(wallet module.Wallet, param *icon.TransactionParam) error {
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

func (r *requestAPI) transactWithContract(senderKey string, contractAddress string,
	amount *big.Int, args map[string]interface{}, method string) (txHash string, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}

	param := icon.TransactionParam{
		Version:     icon.NewHexInt(icon.JsonrpcApiVersion),
		ToAddress:   icon.Address(contractAddress),
		Value:       icon.HexInt(intconv.FormatBigInt(amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: icon.Address(senderWallet.Address().String()),
		StepLimit:   icon.NewHexInt(StepLimit),
		Timestamp:   icon.NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   icon.HexInt(r.networkID),
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

func (r *requestAPI) callContract(contractAddress string, args map[string]string, method string) (interface{}, error) {
	param := &icon.CallParam{
		ToAddress: icon.Address(contractAddress),
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
	} else if res == nil {
		return nil, fmt.Errorf("%v Method Call to ContractAdddress returned nil", method)
	}
	return res, nil
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
	param := icon.TransactionParam{
		Version:     icon.NewHexInt(icon.JsonrpcApiVersion),
		ToAddress:   icon.Address(recepientAddress),
		Value:       icon.HexInt(intconv.FormatBigInt(amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: icon.Address(senderWallet.Address().String()),
		StepLimit:   icon.NewHexInt(StepLimit),
		Timestamp:   icon.NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   icon.HexInt(r.networkID),
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

func (r *requestAPI) transferNativeCrossChain(senderKey, recepientAddress string, amount *big.Int) (txHash string, err error) {
	args := map[string]interface{}{"_to": recepientAddress}
	caddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
		return
	}
	return r.transactWithContract(senderKey, caddr, amount, args, "transferNativeCoin")
}

func (r *requestAPI) approveCrossNativeCoin(coinName string, ownerKey string, amount *big.Int) (approveTxnHash string, err error) {
	coinAddressArgs := map[string]string{"_coinName": coinName}
	btsaddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
		return
	}
	res, err := r.callContract(btsaddr, coinAddressArgs, "coinAddress")
	if err != nil {
		err = errors.Wrap(err, "callContract coinAddress ")
		return
	} else if res == nil {
		err = fmt.Errorf("Approve Call returned nil for coin %v", coinName)
		return
	}
	coinAddress := res.(string)

	approveArgs := map[string]interface{}{"spender": btsaddr, "amount": intconv.FormatBigInt(amount)}
	approveTxnHash, err = r.transactWithContract(ownerKey, coinAddress, big.NewInt(0), approveArgs, "approve")
	if err != nil {
		err = errors.Wrapf(err, "transactWithContract %v", coinAddress)
		return
	}
	return
}

func (r *requestAPI) transferTokensCrossChain(coinName, senderKey, recepientAddress string, amount *big.Int) (string, error) {
	args := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount), "_to": recepientAddress}
	btsaddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	return r.transactWithContract(senderKey, btsaddr, big.NewInt(0), args, "transfer")
}

func (r *requestAPI) transferBatch(coinNames []string, senderKey, recepientAddress string, amounts []*big.Int, nativeCoin string) (txnHash string, err error) {
	if len(amounts) != len(coinNames) {
		return "", fmt.Errorf("Amount and CoinNames len should be same; Got %v and %v", len(amounts), len(coinNames))
	}
	nativeAmount := big.NewInt(0)
	filterNames := []string{}
	filterAmounts := []string{}
	for i := 0; i < len(amounts); i++ {
		if coinNames[i] == nativeCoin {
			nativeAmount = amounts[i]
			continue
		}
		filterAmounts = append(filterAmounts, intconv.FormatBigInt(amounts[i]))
		filterNames = append(filterNames, coinNames[i])
	}

	args := map[string]interface{}{"_coinNames": filterNames, "_values": filterAmounts, "_to": recepientAddress}
	btsaddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}

	txnHash, err = r.transactWithContract(senderKey, btsaddr, nativeAmount, args, "transferBatch")
	return
}

func (r *requestAPI) approveToken(coinName, senderKey string, amount *big.Int, caddr string) (hash string, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	arg1 := map[string]interface{}{"_to": btsAddr, "_value": intconv.FormatBigInt(amount)}
	return r.transactWithContract(senderKey, caddr, big.NewInt(0), arg1, "transfer")
}

func (r *requestAPI) transferTokenCrossChain(coinName, senderKey, recepientAddress string, amount *big.Int) (string, error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	arg2 := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount), "_to": recepientAddress}
	return r.transactWithContract(senderKey, btsaddr, big.NewInt(0), arg2, "transfer")
}

func GetWalletFromFile(walFile string, password string) (module.Wallet, error) {
	keyReader, err := os.Open(walFile)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open(%v)", walFile)
	}
	defer keyReader.Close()

	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll ")
	}
	w, err := wallet.NewFromKeyStore(keyStore, []byte(password))
	if err != nil {
		return nil, errors.Wrap(err, "wallet.NewFromKeyStore")
	}
	return w, nil
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

func generateKeyPair() ([2]string, error) {
	pubkeyBytes, priv := secp256k1.GenerateKeyPair()
	pubKey, err := crypto.ParsePublicKey(pubkeyBytes)
	if err != nil {
		return [2]string{}, errors.Wrap(err, "crypto.ParsePublicKey ")
	}
	addr := common.NewAccountAddressFromPublicKey(pubKey).String()
	return [2]string{hex.EncodeToString(priv), addr}, nil
}

func (r *requestAPI) reclaim() {
	coinName := "TICX"
	amount := big.NewInt(10100000000000000)
	senderKey := "f4e8307da2b4fb7ff89bd984cd0613cfcfacac53abe3a1fd5b7378222bafa5b5"
	args := map[string]interface{}{"_coinName": coinName, "_value": intconv.FormatBigInt(amount)}
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		fmt.Printf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
		return
	}
	hash, err := r.transactWithContract(senderKey, btsAddr, big.NewInt(0), args, "reclaim")
	if err != nil {
		fmt.Println(errors.Wrap(err, "callContract coinAddress "))
		return
	}
	fmt.Println("Hash ", hash)
	return
}

func (r *requestAPI) getNativeCoinBalance(coinName, addr string) (bal *chain.CoinBalance, err error) {
	zeroBalance := big.NewInt(0)
	bal = &chain.CoinBalance{UsableBalance: zeroBalance, RefundableBalance: zeroBalance, LockedBalance: zeroBalance, UserBalance: zeroBalance}
	// Native
	bal.UserBalance, err = r.cl.GetBalance(&icon.AddressParam{Address: icon.Address(addr)})
	if err != nil {
		return nil, errors.Wrapf(err, "%v", err)
	}
	return

}

func (r *requestAPI) getCoinBalance(coinName, addr string, isToken bool) (bal *chain.CoinBalance, err error) {
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
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	// BTS BALANCEOF
	res, err := r.callContract(btsAddr, map[string]string{"_coinName": coinName, "_owner": addr}, "balanceOf")
	if err != nil {
		return nil, errors.Wrap(err, "callContract coinAddress ")
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

	return
}
