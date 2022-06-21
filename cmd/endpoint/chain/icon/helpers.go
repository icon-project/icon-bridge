package icon

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gocrypto "github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
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

func (c *client) TransactWithContract(senderKey string, contractAddress string,
	amount big.Int, args map[string]string, method string, dataType string) (txHash string, err error) {
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
		NetworkID:   HexInt(c.networkID),
		DataType:    dataType,
	}
	argMap := map[string]interface{}{}
	argMap["method"] = method
	argMap["params"] = args
	param.Data = argMap

	if err = SignTransactionParam(senderWallet, &param); err != nil {
		return
	}
	txH, err := c.SendTransaction(&param)
	if err != nil {
		return
	}

	txBytes, err := txH.Value()
	if err != nil {
		return
	}
	txHash = hexutil.Encode(txBytes[:])
	c.waitForResults(context.TODO(), &TransactionHashParam{Hash: *txH})
	return
}

func (c *client) CallContract(contractAddress string, args map[string]string, method string) (interface{}, error) {
	param := &CallParam{
		ToAddress: Address(contractAddress),
		DataType:  "call",
	}
	argMap := map[string]interface{}{}
	argMap["method"] = method
	argMap["params"] = args
	param.Data = argMap
	var res interface{}
	err := c.Call(param, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) GetICXBalance(addr string) (*big.Int, error) {
	// curl -X POST 'http://127.0.0.1:9080/api/v3/default' -H 'Content-Type:application/json' -d '{"id":"1001", "jsonrpc":"2.0", "method": "icx_getBalance", "params":{"address":"hxff0ea998b84ab9955157ab27915a9dc1805edd35"} }'
	return c.GetBalance(&AddressParam{Address: Address(addr)})
}

func (c *client) TransferICX(senderKey string, amount big.Int, recepientAddress string) (txHash string, err error) {
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
		NetworkID:   HexInt(c.networkID),
	}
	if err = SignTransactionParam(senderWallet, &param); err != nil {
		return
	}
	txH, err := c.SendTransaction(&param)
	if err != nil {
		return
	}
	txBytes, err := txH.Value()
	if err != nil {
		return
	}
	txHash = hexutil.Encode(txBytes[:])
	c.waitForResults(context.TODO(), &TransactionHashParam{Hash: *txH})
	return
}

func (c *client) GetIrc2Balance(addr string) (*big.Int, error) {
	//goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" call --to "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f" --method "balanceOf" --param _owner=hx51ecae93216cb6d58bbdc51e2b5d790da94f738a
	args := map[string]string{"_owner": addr}
	res, err := c.CallContract(c.contractAddress.btp_icon_irc2, args, "balanceOf")
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

func (c *client) GetIconWrappedOne(addr string) (*big.Int, error) {
	//goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" call --to "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f" --method "balanceOf" --param _owner=hx51ecae93216cb6d58bbdc51e2b5d790da94f738a
	args := map[string]string{"_owner": addr}
	res, err := c.CallContract(c.contractAddress.btp_icon_irc2_tradeable, args, "balanceOf")
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

func (c *client) TransferIrc2(senderKey string, amount big.Int, recepientAddress string) (txHash string, err error) {
	args := map[string]string{"_to": recepientAddress, "_value": intconv.FormatBigInt(&amount)}
	return c.TransactWithContract(senderKey, c.contractAddress.btp_icon_irc2, *big.NewInt(0), args, "transfer", "call")
}

func (c *client) TransferICXToHarmony(senderKey string, amount big.Int, recepientAddress string) (txHash string, err error) {
	args := map[string]string{"_to": recepientAddress} //"btp://$btp_hmny_net/$btp_hmny_demo_wallet_address"}
	return c.TransactWithContract(senderKey, c.contractAddress.btp_icon_nativecoin_bsh, amount, args, "transferNativeCoin", "call")
}

func (c *client) ApproveIconNativeCoinBSHToAccessHmnyOne(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {

	btpHmnyNativecoinSymbol := "ONE"
	coinAddressArgs := map[string]string{"_coinName": btpHmnyNativecoinSymbol}
	res, err := c.CallContract(c.contractAddress.btp_icon_nativecoin_bsh, coinAddressArgs, "coinAddress")
	if err != nil {
		return
	}
	coinAddress := res.(string)

	approveArgs := map[string]string{"spender": c.contractAddress.btp_icon_nativecoin_bsh, "amount": intconv.FormatBigInt(&amount)}
	approveTxnHash, err = c.TransactWithContract(ownerKey, coinAddress, *big.NewInt(0), approveArgs, "approve", "call")
	if err != nil {
		return
	}

	var ownerWallet module.Wallet
	ownerWallet, err = GetWalletFromPrivKey(ownerKey)
	allowArgs := map[string]string{"owner": ownerWallet.Address().String(), "spender": c.contractAddress.btp_icon_nativecoin_bsh}
	res, err = c.CallContract(coinAddress, allowArgs, "allowance")
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

func (c *client) TransferWrappedOneFromIconToHmny(senderKey string, amount big.Int, recepientAddress string) (string, error) {
	args := map[string]string{"_coinName": "ONE", "_value": intconv.FormatBigInt(&amount), "_to": recepientAddress}
	return c.TransactWithContract(senderKey, c.contractAddress.btp_icon_nativecoin_bsh, *big.NewInt(0), args, "transfer", "call")
}

func (c *client) TransferIrc2ToHmny(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {

	arg1 := map[string]string{"_to": c.contractAddress.btp_icon_token_bsh, "_value": intconv.FormatBigInt(&amount)}
	approveTxnHash, err = c.TransactWithContract(senderKey, c.contractAddress.btp_icon_irc2, *big.NewInt(0), arg1, "transfer", "call")
	if err != nil {
		return
	}

	arg2 := map[string]string{"tokenName": "ETH", "value": intconv.FormatBigInt(&amount), "to": recepientAddress}
	transferTxnHash, err = c.TransactWithContract(senderKey, c.contractAddress.btp_icon_token_bsh, *big.NewInt(0), arg2, "transfer", "call")
	if err != nil {
		return
	}
	return
}

func (c *client) SendTransaction(p *TransactionParam) (*HexBytes, error) {
	var result HexBytes
	if _, err := c.rpcClient.Do("icx_sendTransaction", p, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) SendTransactionAndWait(p *TransactionParam) (*HexBytes, error) {
	var result HexBytes
	if _, err := c.rpcClient.Do("icx_sendTransactionAndWait", p, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) GetTransactionResult(p *TransactionHashParam) (*TransactionResult, error) {
	tr := &TransactionResult{}
	if _, err := c.rpcClient.Do("icx_getTransactionResult", p, tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *client) GetBalance(param *AddressParam) (*big.Int, error) {
	var result HexInt
	_, err := c.rpcClient.Do("icx_getBalance", param, &result)
	if err != nil {
		return nil, err
	}
	bInt, err := result.BigInt()
	if err != nil {
		return nil, err
	}
	return bInt, nil
}

func (c *client) Call(p *CallParam, r interface{}) error {
	_, err := c.rpcClient.Do("icx_call", p, r)
	return err
}

func (c *client) SendTransactionAndGetResult(p *TransactionParam) (*HexBytes, *TransactionResult, error) {
	thp := &TransactionHashParam{}
txLoop:
	for {
		txh, err := c.SendTransaction(p)
		if err != nil {
			switch err {
			case ErrSendFailByOverflow:
				//TODO Retry max
				time.Sleep(DefaultSendTransactionRetryInterval)
				c.log.Debugf("Retry SendTransaction")
				continue txLoop
			default:
				switch re := err.(type) {
				case *jsonrpc.Error:
					switch re.Code {
					case JsonrpcErrorCodeSystem:
						if subEc, err := strconv.ParseInt(re.Message[1:5], 0, 32); err == nil {
							switch subEc {
							case 2000: //DuplicateTransactionError
								//Ignore
								c.log.Debugf("DuplicateTransactionError txh:%v", txh)
								thp.Hash = *txh
								break txLoop
							}
						}
					}
				}
			}
			c.log.Debugf("fail to SendTransaction hash:%v, err:%+v", txh, err)
			return &thp.Hash, nil, err
		}
		thp.Hash = *txh
		break txLoop
	}

txrLoop:
	for {
		time.Sleep(DefaultGetTransactionResultPollingInterval)
		txr, err := c.GetTransactionResult(thp)
		if err != nil {
			switch re := err.(type) {
			case *jsonrpc.Error:
				switch re.Code {
				case JsonrpcErrorCodePending, JsonrpcErrorCodeExecuting:
					//TODO Retry max
					c.log.Debugln("Retry GetTransactionResult", thp)
					continue txrLoop
				}
			}
		}
		c.log.Debugf("GetTransactionResult hash:%v, txr:%+v, err:%+v", thp.Hash, txr, err)
		return &thp.Hash, txr, err
	}
}

func (c *client) waitForResults(ctx context.Context, thp *TransactionHashParam) (txh *HexBytes, txr *TransactionResult, err error) {
	ticker := time.NewTicker(time.Duration(DefaultGetTransactionResultPollingInterval) * time.Nanosecond)
	retryLimit := 10
	retryCounter := 0
	txh = &thp.Hash
	for {
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			err = errors.New("Context Cancelled")
			return
		case <-ticker.C:
			if retryCounter >= retryLimit {
				err = errors.New("Retry Limit Exceeded while waiting for results of transaction")
				return
			}
			retryCounter++
			//c.log.Debugf("GetTransactionResult Attempt: %d", retryCounter)
			txr, err = c.GetTransactionResult(thp)
			if err != nil {
				switch re := err.(type) {
				case *jsonrpc.Error:
					switch re.Code {
					case JsonrpcErrorCodePending, JsonrpcErrorCodeExecuting:
						continue
					}
				}
			}
			//c.log.Debugf("GetTransactionResult hash:%v, txr:%+v, err:%+v", thp.Hash, txr, err)
			return
		}
	}
}

const (
	HeaderKeyIconOptions = "Icon-Options"
	IconOptionsDebug     = "debug"
	IconOptionsTimeout   = "timeout"
)

type IconOptions map[string]string

func (opts IconOptions) Set(key, value string) {
	opts[key] = value
}

func (opts IconOptions) Get(key string) string {
	if opts == nil {
		return ""
	}
	v := opts[key]
	if len(v) == 0 {
		return ""
	}
	return v
}

func (opts IconOptions) Del(key string) {
	delete(opts, key)
}

func (opts IconOptions) SetBool(key string, value bool) {
	opts.Set(key, strconv.FormatBool(value))
}

func (opts IconOptions) GetBool(key string) (bool, error) {
	return strconv.ParseBool(opts.Get(key))
}

func (opts IconOptions) SetInt(key string, v int64) {
	opts.Set(key, strconv.FormatInt(v, 10))
}

func (opts IconOptions) GetInt(key string) (int64, error) {
	return strconv.ParseInt(opts.Get(key), 10, 64)
}

func (opts IconOptions) ToHeaderValue() string {
	if opts == nil {
		return ""
	}
	strs := make([]string, len(opts))
	i := 0
	for k, v := range opts {
		strs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return strings.Join(strs, ",")
}

func NewIconOptionsByHeader(h http.Header) IconOptions {
	s := h.Get(HeaderKeyIconOptions)
	if s != "" {
		kvs := strings.Split(s, ",")
		m := make(map[string]string)
		for _, kv := range kvs {
			if kv != "" {
				idx := strings.Index(kv, "=")
				if idx > 0 {
					m[kv[:idx]] = kv[(idx + 1):]
				} else {
					m[kv] = ""
				}
			}
		}
		return m
	}
	return nil
}
