package icon

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gocrypto "github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/icon-project/icon-bridge/common/intconv"
)

func SignTransactionParam(wallet module.Wallet, param *TransactionParam) error {
	js, err := json.Marshal(param)
	if err != nil {
		return err
	}
	var txSerializeExcludes = map[string]bool{"signature": true}
	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		return err
	}
	bs = append([]byte("icx_sendTransaction."), bs...)
	sig, err := wallet.Sign(gocrypto.SHA3Sum256(bs))
	if err != nil {
		return err
	}
	param.Signature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func (r *requestAPI) transactWithContract(senderKey string, contractAddress string,
	amount big.Int, args map[string]string, method string, dataType string) (txHash string, logs interface{}, err error) {
	var senderWallet module.Wallet
	senderWallet, err = GetWalletFromPrivKey(senderKey)
	if err != nil {
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
		return
	}
	txH, err := r.cl.SendTransaction(&param)
	if err != nil {
		return
	}

	txBytes, err := txH.Value()
	if err != nil {
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
		return nil, err
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
		err = errors.Wrap(err, "Get Wallet ")
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
		return
	}
	txH, err := r.cl.SendTransaction(&param)
	if err != nil {
		return
	}
	txBytes, err := txH.Value()
	if err != nil {
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
	res, err := r.callContract(r.contractAddress.btp_icon_irc2, args, "balanceOf")
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, errors.New("Nil value")
	}
	resStr, ok := res.(string)
	if !ok {
		return nil, errors.New("Unexpected type")
	}
	n := new(big.Int)
	n.SetString(resStr[2:], 16) //remove 0x
	return n, nil
}

func (r *requestAPI) getIconWrappedOne(addr string) (*big.Int, error) {
	//goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" call --to "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f" --method "balanceOf" --param _owner=hx51ecae93216cb6d58bbdc51e2b5d790da94f738a
	args := map[string]string{"_owner": addr}
	res, err := r.callContract(r.contractAddress.btp_icon_irc2_tradeable, args, "balanceOf")
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, errors.New("Nil value")
	}
	resStr, ok := res.(string)
	if !ok {
		return nil, errors.New("Unexpected type")
	}
	n := new(big.Int)
	n.SetString(resStr[2:], 16)
	return n, nil
}

func (r *requestAPI) transferIrc2(senderKey string, amount big.Int, recepientAddress string) (txHash string, logs interface{}, err error) {
	args := map[string]string{"_to": recepientAddress, "_value": intconv.FormatBigInt(&amount)}
	return r.transactWithContract(senderKey, r.contractAddress.btp_icon_irc2, *big.NewInt(0), args, "transfer", "call")
}

func (r *requestAPI) TransferICXToHarmony(senderKey string, amount big.Int, recepientAddress string) (txHash string, logs interface{}, err error) {
	args := map[string]string{"_to": recepientAddress} //"btp://$btp_hmny_net/$btp_hmny_demo_wallet_address"}
	return r.transactWithContract(senderKey, r.contractAddress.btp_icon_nativecoin_bsh, amount, args, "transferNativeCoin", "call")
}

func (r *requestAPI) approveIconNativeCoinBSHToAccessHmnyOne(ownerKey string, amount big.Int) (approveTxnHash string, logs interface{}, allowanceAmount *big.Int, err error) {

	btpHmnyNativecoinSymbol := "ONE"
	coinAddressArgs := map[string]string{"_coinName": btpHmnyNativecoinSymbol}
	res, err := r.callContract(r.contractAddress.btp_icon_nativecoin_bsh, coinAddressArgs, "coinAddress")
	if err != nil {
		return
	}
	coinAddress := res.(string)

	approveArgs := map[string]string{"spender": r.contractAddress.btp_icon_nativecoin_bsh, "amount": intconv.FormatBigInt(&amount)}
	approveTxnHash, logs, err = r.transactWithContract(ownerKey, coinAddress, *big.NewInt(0), approveArgs, "approve", "call")
	if err != nil {
		return
	}

	var ownerWallet module.Wallet
	ownerWallet, err = GetWalletFromPrivKey(ownerKey)
	allowArgs := map[string]string{"owner": ownerWallet.Address().String(), "spender": r.contractAddress.btp_icon_nativecoin_bsh}
	res, err = r.callContract(coinAddress, allowArgs, "allowance")
	if err != nil {
		return
	}
	if resStr, ok := res.(string); ok {
		allowanceAmount = new(big.Int)
		allowanceAmount.SetString(resStr[2:], 16)
	} else {
		err = errors.New("allowance is not expected type ")
	}
	return
}

func (r *requestAPI) transferWrappedOneFromIconToHmny(senderKey string, amount big.Int, recepientAddress string) (string, interface{}, error) {
	args := map[string]string{"_coinName": "ONE", "_value": intconv.FormatBigInt(&amount), "_to": recepientAddress}
	return r.transactWithContract(senderKey, r.contractAddress.btp_icon_nativecoin_bsh, *big.NewInt(0), args, "transfer", "call")
}

func (r *requestAPI) transferIrc2ToHmny(senderKey string, amount big.Int, recepientAddress string) (string, interface{}, error) {

	if recepientAddress == r.contractAddress.btp_icon_token_bsh {
		arg1 := map[string]string{"_to": recepientAddress, "_value": intconv.FormatBigInt(&amount)}
		return r.transactWithContract(senderKey, r.contractAddress.btp_icon_irc2, *big.NewInt(0), arg1, "transfer", "call")
	}

	arg2 := map[string]string{"tokenName": "ETH", "value": intconv.FormatBigInt(&amount), "to": recepientAddress}
	return r.transactWithContract(senderKey, r.contractAddress.btp_icon_token_bsh, *big.NewInt(0), arg2, "transfer", "call")
}
