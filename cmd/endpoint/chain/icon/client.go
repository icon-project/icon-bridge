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
	"math/big"
	"net/http"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	DefaultSendTransactionRetryInterval        = 3 * time.Second         //3sec
	DefaultGetTransactionResultPollingInterval = 1500 * time.Millisecond //1.5sec
	StepLimit                                  = 3500000000
)

type client struct {
	rpcClient       *jsonrpc.Client
	log             log.Logger
	contractAddress *contractAddress
	networkID       string
}

type contractAddress struct {
	btp_icon_irc2           string
	btp_icon_irc2_tradeable string
	btp_icon_nativecoin_bsh string
	btp_icon_token_bsh      string
}

func (cAddr *contractAddress) FromMap(contractAddrsMap map[string]string) {
	if contractAddrsMap == nil {
		return
	}
	cAddr.btp_icon_irc2 = contractAddrsMap["btp_icon_irc2"]
	cAddr.btp_icon_irc2_tradeable = contractAddrsMap["btp_icon_irc2_tradeable"]
	cAddr.btp_icon_nativecoin_bsh = contractAddrsMap["btp_icon_nativecoin_bsh"]
	cAddr.btp_icon_token_bsh = contractAddrsMap["btp_icon_token_bsh"]
}

func newClient(uri string, l log.Logger, cAddr *contractAddress, networkID string) (*client, error) {

	tr := &http.Transport{MaxIdleConnsPerHost: 1000}
	c := &client{
		rpcClient:       jsonrpc.NewJsonRpcClient(&http.Client{Transport: tr}, uri),
		log:             l,
		contractAddress: cAddr,
		networkID:       networkID,
	}
	opts := IconOptions{}
	opts.SetBool(IconOptionsDebug, true)
	c.rpcClient.CustomHeader[HeaderKeyIconOptions] = opts.ToHeaderValue()
	return c, nil
}

func New(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (chain.Client, error) {
	cAddr := &contractAddress{}
	cAddr.FromMap(contractAddrsMap)

	return newClient(url, l, cAddr, networkID)
}

func (c *client) GetCoinBalance(addr string) (*big.Int, error) {
	return c.GetICXBalance(addr)
}

func (c *client) GetEthToken(addr string) (val *big.Int, err error) {
	return c.GetIrc2Balance(addr)
}

func (c *client) GetWrappedCoin(addr string) (val *big.Int, err error) {
	return c.GetIconWrappedOne(addr)
}

func (c *client) TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferICX(senderKey, amount, recepientAddress)
}

func (c *client) TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferIrc2(senderKey, amount, recepientAddress)
}

func (c *client) TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferICXToHarmony(senderKey, amount, recepientAddress)
}

func (c *client) TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferWrappedOneFromIconToHmny(senderKey, amount, recepientAddress)
}

func (c *client) TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	return c.TransferIrc2ToHmny(senderKey, amount, recepientAddress)
}

func (c *client) ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	return c.ApproveIconNativeCoinBSHToAccessHmnyOne(ownerKey, amount)
}

func (c *client) GetAddressFromPrivKey(key string) (*string, error) {
	return getAddressFromPrivKey(key)
}

func (c *client) GetFullAddress(addr string) *string {
	fullAddr := "btp://" + c.networkID + ".icon/" + addr
	return &fullAddr
}
