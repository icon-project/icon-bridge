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

package bsc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	bscTypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/bsc/types"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	txMaxDataSize        = 8 * 1024 // 8 KB
	txOverheadScale      = 0.01     // base64 encoding overhead 0.36, rlp and other fields 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale)
	defaultSendTxTimeout = 15 * time.Second
	defaultGasPrice      = 18000000000
	maxGasPriceBoost     = 10.0
	defaultReadTimeout   = 50 * time.Second //
	DefaultGasLimit      = 25000000
)


type senderOptions struct {
	GasLimit         uint64         `json:"gas_limit"`
	TxDataSizeLimit  uint64         `json:"tx_data_size_limit"`
	BoostGasPrice    float64        `json:"boost_gas_price"`
	BalanceThreshold intconv.BigInt `json:"balance_threshold"`
}

type sender struct {
	log          log.Logger
	w            *wallet.EvmWallet
	src          chain.BTPAddress
	dst          chain.BTPAddress
	opts         senderOptions
	cls          []IClient
	prevGasPrice *big.Int
}

func (s *sender) client() IClient {
	randInt := rand.Intn(len(s.cls))
	return s.cls[randInt]
}

func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {
	s := &sender{
		log:          l,
		w:            w.(*wallet.EvmWallet),
		src:          src,
		dst:          dst,
		prevGasPrice: big.NewInt(defaultGasPrice),
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}
	err := json.Unmarshal(rawOpts, &s.opts)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal opt:%v err:%+v", rawOpts, err)
	}
	if s.opts.BoostGasPrice < 1.0 {
		s.opts.BoostGasPrice = 1.0
	}
	if s.opts.BoostGasPrice > maxGasPriceBoost {
		s.opts.BoostGasPrice = maxGasPriceBoost
	}
	s.cls, err = newClients(urls, dst.ContractAddress(), s.log)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// BMCLinkStatus ...
// returns the BMCLinkStatus for "src" link
func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	status, err := s.client().GetStatus(&bind.CallOpts{Context: ctx}, s.src.String())
	if err != nil {
		s.log.Error("GetStatus", "err", err)
		return nil, err
	}
	ls := &chain.BMCLinkStatus{}
	ls.TxSeq = status.TxSeq.Uint64()
	ls.RxSeq = status.RxSeq.Uint64()
	ls.RxHeight = status.RxHeight.Uint64()
	ls.CurrentHeight = status.CurrentHeight.Uint64()
	return ls, nil
}

// Segment ...
func (s *sender) Segment(
	ctx context.Context, msg *chain.Message,
) (tx chain.RelayTx, newMsg *chain.Message, err error) {
	if ctx.Err() != nil {
		return nil, msg, ctx.Err()
	}

	if s.opts.TxDataSizeLimit == 0 {
		limit := defaultTxSizeLimit
		s.opts.TxDataSizeLimit = uint64(limit)
	}
	if len(msg.Receipts) == 0 {
		return nil, msg, nil
	}

	rm := &chain.RelayMessage{
		Receipts: make([][]byte, 0),
	}

	var msgSize uint64

	newMsg = &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts,
	}
	for i, receipt := range msg.Receipts {
		rlpEvents, err := codec.RLP.MarshalToBytes(receipt.Events)
		if err != nil {
			return nil, nil, err
		}
		rlpReceipt, err := codec.RLP.MarshalToBytes(&chain.RelayReceipt{
			Index:  receipt.Index,
			Height: receipt.Height,
			Events: rlpEvents,
		})
		if err != nil {
			return nil, nil, err
		}
		newMsgSize := msgSize + uint64(len(rlpReceipt))
		if newMsgSize > s.opts.TxDataSizeLimit {
			newMsg.Receipts = msg.Receipts[i:]
			break
		}
		msgSize = newMsgSize
		rm.Receipts = append(rm.Receipts, rlpReceipt)
	}
	message, err := codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, nil, err
	}
	gasPrice, gasHeight, err := s.client().GetMedianGasPriceForBlock(ctx)
	if err != nil || gasPrice.Int64() == 0 {
		s.log.Infof("GetMedianGasPriceForBlock(%v) Msg: %v. Using default value for gas price \n", gasHeight.String(), err)
		gasPrice = s.prevGasPrice
	} else {
		s.prevGasPrice = gasPrice
		s.log.Infof("GetMedianGasPriceForBlock(%v) price: %v \n", gasHeight.String(), gasPrice.String())
	}
	boostedGasPrice, _ := (&big.Float{}).Mul(
		(&big.Float{}).SetInt64(gasPrice.Int64()),
		(&big.Float{}).SetFloat64(s.opts.BoostGasPrice),
	).Int(nil)
	tx, err = s.newRelayTx(ctx, msg.From.String(), message, boostedGasPrice)
	if err != nil {
		return nil, nil, err
	}

	return tx, newMsg, nil
}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	bal, err := s.client().GetBalance(ctx, s.w.Address())
	return bal, &s.opts.BalanceThreshold.Int, err
}

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte, gasPrice *big.Int) (*relayTx, error) {
	client := s.client()

	newTransactOpts := func(w bscTypes.Wallet) (*bind.TransactOpts, error) {
		txo, err := bind.NewKeyedTransactorWithChainID(w.(*wallet.EvmWallet).Skey, client.GetChainID())
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
		defer cancel()
		txo.GasPrice, _ = client.SuggestGasPrice(ctx)
		txo.GasLimit = uint64(DefaultGasLimit)
		return txo, nil
	}

	txOpts, err := newTransactOpts(s.w)
	if err != nil {
		return nil, err
	}
	txOpts.Context = ctx
	if s.opts.GasLimit > 0 {
		txOpts.GasLimit = s.opts.GasLimit
	}
	txOpts.GasPrice = gasPrice
	return &relayTx{
		Prev:    prev,
		Message: message, // base64.URLEncoding.EncodeToString(rlpCrm),
		opts:    txOpts,
		cl:      client,
	}, nil
}

type relayTx struct {
	Prev    string `json:"_prev"`
	Message []byte `json:"_msg"`

	opts      *bind.TransactOpts
	pendingTx *ethtypes.Transaction
	cl        IClient
}

func (tx *relayTx) ID() interface{} {
	if tx.pendingTx != nil {
		return tx.pendingTx.Hash()
	}
	return nil
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	tx.cl.Log().WithFields(log.Fields{
		"prev": tx.Prev}).Debug("handleRelayMessage: send tx")

	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()
	txOpts := *tx.opts
	txOpts.Context = _ctx
	nonce, err := tx.cl.NonceAt(ctx, txOpts.From, nil)
	if err != nil {
		return err
	}
	txOpts.Nonce = (&big.Int{}).SetUint64(nonce)
	defer func() {
		if tx.pendingTx != nil {
			txBytes, _ := tx.pendingTx.MarshalJSON()
			tx.cl.Log().WithFields(log.Fields{
				"tx": string(txBytes)}).Debug("handleRelayMessage: tx sent")
		}
	}()
	tx.pendingTx, err = tx.cl.HandleRelayMessage(&txOpts, tx.Prev, tx.Message)
	if err != nil {
		tx.cl.Log().WithFields(log.Fields{
			"error": err}).Debug("handleRelayMessage: send tx")
		if err.Error() == "insufficient funds for gas * price + value" {
			return chain.ErrInsufficientBalance
		}
		return err
	}
	// tx.cl.Log().WithFields(log.Fields{
	// 	"txh": tx.pendingTx.Hash(),
	// 	"msg": btpcommon.HexBytes(tx.Message)}).Debug("handleRelayMessage: tx sent")
	return nil
}

func (tx *relayTx) Receipt(ctx context.Context) (blockNumber uint64, err error) {
	if tx.pendingTx == nil {
		return 0, fmt.Errorf("no pending tx")
	}

	for i, isPending := 0, true; i < 5 && (isPending || err == ethereum.NotFound); i++ {
		time.Sleep(time.Second)
		_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
		defer cancel()
		_, isPending, err = tx.cl.TransactionByHash(_ctx, tx.pendingTx.Hash())
	}
	if err != nil {
		return 0, err
	}
	_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	txr, err := tx.cl.TransactionReceipt(_ctx, tx.pendingTx.Hash())
	if err != nil {
		return 0, err
	}

	if txr.Status == 0 {
		callMsg := ethereum.CallMsg{
			From:       tx.opts.From,
			To:         tx.pendingTx.To(),
			Gas:        tx.pendingTx.Gas(),
			GasPrice:   tx.pendingTx.GasPrice(),
			Value:      tx.pendingTx.Value(),
			AccessList: tx.pendingTx.AccessList(),
			Data:       tx.pendingTx.Data(),
		}

		_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
		defer cancel()
		data, err := tx.cl.CallContract(_ctx, callMsg, txr.BlockNumber)
		if err != nil {
			return 0, err
		}

		return 0, chain.RevertError(revertReason(data))
	}

	tx.cl.Log().WithFields(log.Fields{
		"txh": tx.pendingTx.Hash()}).Debug("handleRelayMessage: success")

	return txr.BlockNumber.Uint64(), nil
}

func revertReason(data []byte) string {
	if len(data) < 4+32+32 {
		return ""
	}
	data = data[4+32:] // ignore method and index
	length := binary.BigEndian.Uint64(data[24:32])
	return string(data[32 : 32+length])
}
