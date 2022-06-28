/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	DefaultSendTransactionRetryInterval        = 3 * time.Second         //3sec
	DefaultGetTransactionResultPollingInterval = 1500 * time.Millisecond //1.5sec
	StepLimit                                  = 3500000000
)

type client struct {
	*jsonrpc.Client
	log   log.Logger
	conns map[string]*websocket.Conn
	mtx   sync.Mutex
}

func newClient(uri string, l log.Logger) (*client, error) {

	tr := &http.Transport{MaxIdleConnsPerHost: 1000}
	c := &client{
		Client: jsonrpc.NewJsonRpcClient(&http.Client{Transport: tr}, uri),
		log:    l,
		conns:  make(map[string]*websocket.Conn),
	}
	opts := IconOptions{}
	opts.SetBool(IconOptionsDebug, true)
	c.CustomHeader[HeaderKeyIconOptions] = opts.ToHeaderValue()
	return c, nil
}

// func NewClient(url string, l log.Logger) (*client, error) {
// 	return newClient(url, l)
// }

func (c *client) SendTransaction(p *TransactionParam) (*HexBytes, error) {
	var result HexBytes
	if _, err := c.Do("icx_sendTransaction", p, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) SendTransactionAndWait(p *TransactionParam) (*HexBytes, error) {
	var result HexBytes
	if _, err := c.Do("icx_sendTransactionAndWait", p, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *client) GetTransactionResult(p *TransactionHashParam) (*TransactionResult, error) {
	tr := &TransactionResult{}
	if _, err := c.Do("icx_getTransactionResult", p, tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *client) GetBalance(param *AddressParam) (*big.Int, error) {
	var result HexInt
	_, err := c.Do("icx_getBalance", param, &result)
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
	_, err := c.Do("icx_call", p, r)
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

var txSerializeExcludes = map[string]bool{"signature": true}

func (c *client) SignTransaction(w Wallet, p *TransactionParam) error {
	p.Timestamp = NewHexInt(time.Now().UnixNano() / int64(time.Microsecond))
	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	bs, err := SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		return err
	}
	bs = append([]byte("icx_sendTransaction."), bs...)
	txHash := crypto.SHA3Sum256(bs)
	p.TxHash = NewHexBytes(txHash)
	sig, err := w.Sign(txHash)
	if err != nil {
		return err
	}
	p.Signature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func (c *client) GetBlockByHeight(p *BlockHeightParam) (*Block, error) {
	result := &Block{}
	if _, err := c.Do("icx_getBlockByHeight", p, &result); err != nil {
		return nil, err
	}
	return result, nil
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
