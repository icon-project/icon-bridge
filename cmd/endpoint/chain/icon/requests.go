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

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/haltingstate/secp256k1-go"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gocrypto "github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	cl                    *client
}

func newRequestAPI(url string, l log.Logger, contractNameToAddress map[chain.ContractName]string, networkID string) (*requestAPI, error) {
	cl, err := newClient(url, l)
	if err != nil {
		return nil, errors.Wrap(err, "newClient ")
	}
	return &requestAPI{networkID: networkID, contractNameToAddress: contractNameToAddress, cl: cl}, nil
}

func SignTransactionParam(wallet module.Wallet, param *TransactionParam) error {
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
	param := TransactionParam{
		Version:     NewHexInt(JsonrpcApiVersion),
		ToAddress:   Address(contractAddress),
		Value:       HexInt(intconv.FormatBigInt(&amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: Address(senderWallet.Address().String()),
		StepLimit:   NewHexInt(StepLimit),
		Timestamp:   NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   HexInt(r.networkID),
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
	// _, txr, err := r.cl.waitForResults(context.TODO(), &TransactionHashParam{Hash: *txH})
	// plogs := make([]*TxnEventLog, len(txr.EventLogs))
	// for i := range txr.EventLogs {
	// 	plogs[i] = &txr.EventLogs[i]
	// }
	// logs = plogs
	return
}

func (r *requestAPI) callContract(contractAddress string, args map[string]string, method string) (interface{}, error) {
	param := &CallParam{
		ToAddress: Address(contractAddress),
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

func (r *requestAPI) getICXBalance(addr string) (*big.Int, error) {
	// curl -X POST 'http://127.0.0.1:9080/api/v3/default' -H 'Content-Type:application/json' -d '{"id":"1001", "jsonrpc":"2.0", "method": "icx_getBalance", "params":{"address":"hxff0ea998b84ab9955157ab27915a9dc1805edd35"} }'
	return r.cl.GetBalance(&AddressParam{Address: Address(addr)})
}

func (r *requestAPI) transferICX(senderKey string, amount big.Int, recepientAddress string) (txHash string, logs interface{}, err error) {
	//goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" sendtx transfer --to "hx267ed8d02bae84ada9f6ab486d4557aa4763b33a" --value "20" --key_store devnet/docker/icon-hmny/src/icon.god.wallet.json --key_password "gochain" --nid "6003319" --step_limit "3500000000"
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "GetWalletFromPrivKey ")
		return
	}
	param := TransactionParam{
		Version:     NewHexInt(JsonrpcApiVersion),
		ToAddress:   Address(recepientAddress),
		Value:       HexInt(intconv.FormatBigInt(&amount)), //NewHexInt(amount.Int64()) Using Int64() can overflow for large amounts
		FromAddress: Address(senderWallet.Address().String()),
		StepLimit:   NewHexInt(StepLimit),
		Timestamp:   NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   HexInt(r.networkID),
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
	// _, txr, err := r.cl.waitForResults(context.TODO(), &TransactionHashParam{Hash: *txH})
	// plogs := make([]*TxnEventLog, len(txr.EventLogs))
	// for i := range txr.EventLogs {
	// 	plogs[i] = &txr.EventLogs[i]
	// }
	// logs = plogs
	return
}

func (r *requestAPI) getIrc2Balance(addr string) (*big.Int, error) {
	//goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" call --to "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f" --method "balanceOf" --param _owner=hx51ecae93216cb6d58bbdc51e2b5d790da94f738a
	args := map[string]string{"_owner": addr}
	caddr, ok := r.contractNameToAddress[chain.Irc2Icon]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.Irc2Icon)
	}
	res, err := r.callContract(caddr, args, "balanceOf")
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
	n.SetString(resStr[2:], 16) //remove 0x
	return n, nil
}

func (r *requestAPI) getIconWrappedOne(addr string) (*big.Int, error) {
	//goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" call --to "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f" --method "balanceOf" --param _owner=hx51ecae93216cb6d58bbdc51e2b5d790da94f738a
	args := map[string]string{"_owner": addr}
	caddr, ok := r.contractNameToAddress[chain.Irc2TradeableIcon]
	if !ok {
		return nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.Irc2TradeableIcon)
	}
	res, err := r.callContract(caddr, args, "balanceOf")
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

func (r *requestAPI) transferIrc2(senderKey string, amount big.Int, recepientAddress string) (txHash string, logs interface{}, err error) {
	args := map[string]string{"_to": recepientAddress, "_value": intconv.FormatBigInt(&amount)}
	caddr, ok := r.contractNameToAddress[chain.Irc2Icon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.Irc2Icon)
		return
	}
	return r.transactWithContract(senderKey, caddr, *big.NewInt(0), args, "transfer", "call")
}

func (r *requestAPI) TransferICXToHarmony(senderKey string, amount big.Int, recepientAddress string) (txHash string, logs interface{}, err error) {
	args := map[string]string{"_to": recepientAddress} //"btp://$btp_hmny_net/$btp_hmny_demo_wallet_address"}
	caddr, ok := r.contractNameToAddress[chain.NativeBSHIcon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.NativeBSHIcon)
		return
	}
	return r.transactWithContract(senderKey, caddr, amount, args, "transferNativeCoin", "call")
}

func (r *requestAPI) approveIconNativeCoinBSHToAccessHmnyOne(ownerKey string, amount big.Int) (approveTxnHash string, logs interface{}, allowanceAmount *big.Int, err error) {

	btpHmnyNativecoinSymbol := "ONE"
	coinAddressArgs := map[string]string{"_coinName": btpHmnyNativecoinSymbol}
	caddr, ok := r.contractNameToAddress[chain.NativeBSHIcon]
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.NativeBSHIcon)
		return
	}
	res, err := r.callContract(caddr, coinAddressArgs, "coinAddress")
	if err != nil {
		err = errors.Wrap(err, "callContract coinAddress ")
		return
	}
	coinAddress := res.(string)

	approveArgs := map[string]string{"spender": caddr, "amount": intconv.FormatBigInt(&amount)}
	approveTxnHash, logs, err = r.transactWithContract(ownerKey, coinAddress, *big.NewInt(0), approveArgs, "approve", "call")
	if err != nil {
		err = errors.Wrapf(err, "transactWithContract %v", coinAddress)
		return
	}

	var ownerWallet module.Wallet
	ownerWallet, err = GetWalletFromPrivKey(ownerKey)
	allowArgs := map[string]string{"owner": ownerWallet.Address().String(), "spender": caddr}
	res, err = r.callContract(coinAddress, allowArgs, "allowance")
	if err != nil {
		err = errors.Wrap(err, "callContract allowance ")
		return
	}
	if resStr, ok := res.(string); ok {
		allowanceAmount = new(big.Int)
		allowanceAmount.SetString(resStr[2:], 16)
	} else {
		err = fmt.Errorf("Expected type string; Got %T", res)
	}
	return
}

func (r *requestAPI) transferWrappedOneFromIconToHmny(senderKey string, amount big.Int, recepientAddress string) (string, interface{}, error) {
	args := map[string]string{"_coinName": "ONE", "_value": intconv.FormatBigInt(&amount), "_to": recepientAddress}
	caddr, ok := r.contractNameToAddress[chain.NativeBSHIcon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.NativeBSHIcon)
	}
	return r.transactWithContract(senderKey, caddr, *big.NewInt(0), args, "transfer", "call")
}

func (r *requestAPI) transferIrc2ToHmny(senderKey string, amount big.Int, recepientAddress string) (string, interface{}, error) {
	caddrbsh, ok := r.contractNameToAddress[chain.TokenBSHIcon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.TokenBSHIcon)
	}
	caddrirc2, ok := r.contractNameToAddress[chain.Irc2Icon]
	if !ok {
		return "", nil, fmt.Errorf("contractNameToAddress doesn't include name %v", chain.Irc2Icon)
	}

	arg1 := map[string]string{"_to": caddrbsh, "_value": intconv.FormatBigInt(&amount)}
	_, _, err := r.transactWithContract(senderKey, caddrirc2, *big.NewInt(0), arg1, "transfer", "call")
	if err != nil {
		return "", nil, errors.Wrapf(err, "transactWithContract %v", caddrirc2)
	}

	arg2 := map[string]string{"tokenName": "ETH", "value": intconv.FormatBigInt(&amount), "to": recepientAddress}
	return r.transactWithContract(senderKey, caddrbsh, *big.NewInt(0), arg2, "transfer", "call")
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

func CreateKeyStore(password string) (*string, error) {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return nil, errors.Wrap(err, "keystore.NewAccount ")
	}
	addr := account.Address.Hex()
	return &addr, nil
}

func getAddressFromPrivKey(pKey string) (*string, error) {
	privBytes, err := hex.DecodeString(pKey)
	if err != nil {
		return nil, errors.Wrap(err, "hex.DecodeString ")
	}
	pubkeyBytes := secp256k1.PubkeyFromSeckey(privBytes)
	pubKey, err := crypto.ParsePublicKey(pubkeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "crypto.ParsePublicKey ")
	}
	addr := common.NewAccountAddressFromPublicKey(pubKey).String()
	return &addr, nil
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
