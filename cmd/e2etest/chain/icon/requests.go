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
	StepLimit                                  = 3500000000
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
	amount big.Int, args map[string]string, method string, dataType string) (txHash string, logs interface{}, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}
	param := icon.TransactionParam{
		Version:     icon.NewHexInt(icon.JsonrpcApiVersion),
		ToAddress:   icon.Address(contractAddress),
		Value:       icon.HexInt(intconv.FormatBigInt(&amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: icon.Address(senderWallet.Address().String()),
		StepLimit:   icon.NewHexInt(StepLimit),
		Timestamp:   icon.NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   icon.HexInt(r.networkID),
		DataType:    dataType,
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

func (r *requestAPI) getICXBalance(addr string) (*big.Int, error) {
	return r.cl.GetBalance(&icon.AddressParam{Address: icon.Address(addr)})
}

func (r *requestAPI) transferTokenIntraChain(senderKey, recepientAddress string, amount big.Int) (txHash string, logs interface{}, err error) {
	args := map[string]string{"_to": recepientAddress, "_value": intconv.FormatBigInt(&amount)}
	caddr, ok := r.contractNameToAddress[chain.TICXIcon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.TICXIcon)
		return
	}
	return r.transactWithContract(senderKey, caddr, *big.NewInt(0), args, "transfer", "call")
}

func (r *requestAPI) transferNativeIntraChain(senderKey, recepientAddress string, amount big.Int) (txHash string, logs interface{}, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}
	param := icon.TransactionParam{
		Version:     icon.NewHexInt(icon.JsonrpcApiVersion),
		ToAddress:   icon.Address(recepientAddress),
		Value:       icon.HexInt(intconv.FormatBigInt(&amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
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

func (r *requestAPI) transferNativeCrossChain(senderKey, recepientAddress string, amount big.Int) (txHash string, logs interface{}, err error) {
	args := map[string]string{"_to": recepientAddress}
	caddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
		return
	}
	return r.transactWithContract(senderKey, caddr, amount, args, "transferNativeCoin", "call")
}

func (r *requestAPI) approveCrossNativeCoin(coinName string, ownerKey string, amount big.Int) (approveTxnHash string, logs interface{}, err error) {
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

	approveArgs := map[string]string{"spender": btsaddr, "amount": intconv.FormatBigInt(&amount)}
	approveTxnHash, logs, err = r.transactWithContract(ownerKey, coinAddress, *big.NewInt(0), approveArgs, "approve", "call")
	if err != nil {
		err = errors.Wrapf(err, "transactWithContract %v", coinAddress)
		return
	}
	return
}

func (r *requestAPI) transferWrappedCrossChain(coinName, senderKey, recepientAddress string, amount big.Int) (string, interface{}, error) {
	args := map[string]string{"_coinName": coinName, "_value": intconv.FormatBigInt(&amount), "_to": recepientAddress}
	btsaddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	return r.transactWithContract(senderKey, btsaddr, *big.NewInt(0), args, "transfer", "call")
}

func (r *requestAPI) approveToken(coinName, senderKey string, amount big.Int) (hash string, log interface{}, err error) {
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	ticxAddr, ok := r.contractNameToAddress[chain.TICXIcon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.TICXIcon)
	}
	arg1 := map[string]string{"_to": btsAddr, "_value": intconv.FormatBigInt(&amount)}
	return r.transactWithContract(senderKey, ticxAddr, *big.NewInt(0), arg1, "transfer", "call")
}

func (r *requestAPI) transferTokenCrossChain(coinName, senderKey, recepientAddress string, amount big.Int) (string, interface{}, error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	arg2 := map[string]string{"_coinName": coinName, "_value": intconv.FormatBigInt(&amount), "_to": recepientAddress}
	return r.transactWithContract(senderKey, btsaddr, *big.NewInt(0), arg2, "transfer", "call")
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

func (r *requestAPI) getWrappedCoinBalance(coinName, addr string) (*big.Int, error) {
	args := map[string]string{"_coinName": coinName}
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	res, err := r.callContract(btsAddr, args, "coinAddress")
	if err != nil {
		return nil, errors.Wrap(err, "callContract coinAddress ")
	} else if res == nil {
		return nil, errors.New("callContract returned nil value ")
	}
	coinAddress, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string Got %T", res)
	}

	args = map[string]string{"_owner": addr}
	res, err = r.callContract(coinAddress, args, "balanceOf")
	if err != nil {
		return nil, errors.Wrap(err, "callContract balanceOf ")
	} else if res == nil {
		return nil, errors.New("callContract returned nil value ")
	}
	resStr, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string Got %T", res)
	}
	n := new(big.Int)
	n.SetString(resStr[2:], 16)
	return n, nil
}

func (r *requestAPI) getAllowanceForNativeTokens(coinName string, addr string) (*big.Int, error) {
	args := map[string]string{"_coinName": coinName, "_owner": addr}
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	res, err := r.callContract(btsAddr, args, "balanceOf")
	if err != nil {
		return nil, errors.Wrap(err, "callContract coinAddress ")
	} else if res == nil {
		return nil, errors.New("callContract returned nil value ")
	}
	// fmt.Println(res)
	balanceMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected type map[string]interface{} Got %T", res)
	}
	usableBalance, ok := balanceMap["usable"]
	if !ok {
		return nil, fmt.Errorf("")
	}
	usableBalanceStr, ok := usableBalance.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string Got %T", usableBalance)
	}
	n := new(big.Int)
	n.SetString(usableBalanceStr[2:], 16)
	return n, nil
}

func (r *requestAPI) getAllowanceForWrappedCoins(coinName string, addr string) (*big.Int, error) {
	args := map[string]string{"_coinName": coinName}
	btsAddr, ok := r.contractNameToAddress[chain.BTSIcon]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTSIcon)
	}
	res, err := r.callContract(btsAddr, args, "coinAddress")
	if err != nil {
		return nil, errors.Wrap(err, "callContract coinAddress ")
	} else if res == nil {
		return nil, errors.New("callContract returned nil value ")
	}
	coinAddress, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string Got %T", res)
	}

	args = map[string]string{"owner": addr, "spender": btsAddr}
	res, err = r.callContract(coinAddress, args, "allowance")
	if err != nil {
		return nil, errors.Wrap(err, "callContract allowance ")
	} else if res == nil {
		return nil, errors.New("callContract returned nil value ")
	}
	resStr, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string Got %T", res)
	}
	n := new(big.Int)
	n.SetString(resStr[2:], 16)
	return n, nil
}
