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
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	txMaxDataSize                 = 524288 //512 * 1024 // 512kB
	txOverheadScale               = 0.37   //base64 encoding overhead 0.36, rlp and other fields 0.01
	defaultTxSizeLimit            = txMaxDataSize / (1 + txOverheadScale)
	defaultGetRelayResultInterval = time.Second
	defaultRelayReSendInterval    = time.Second
	defaultStepLimit              = 13610920010
)

// NewSender ...
// returns a new sender client for icon
func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {
	s := &sender{
		log: l,
		w:   w,
		src: src,
		dst: dst,
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	if err := unmarshalOpt(rawOpts, &s.opts); err != nil {
		return nil, err
	}
	s.cl = NewClient(urls[0], l)
	return s, nil
}

type senderOptions struct {
	StepLimit        uint64  `json:"step_limit"`
	TxDataSizeLimit  uint64  `json:"tx_data_size_limit"`
	BalanceThreshold big.Int `json:"balance_threshold"`
}

func unmarshalOpt(data []byte, opts *senderOptions) error {
	type SenderOptionsTemp struct {
		StepLimit        uint64  `json:"step_limit"`
		TxDataSizeLimit  uint64  `json:"tx_data_size_limit"`
		BalanceThreshold string `json:"balance_threshold"`
	}
	var senderOptionsObj SenderOptionsTemp

	if err := json.Unmarshal(data, &senderOptionsObj); err != nil {
		return err
	}

	opts.StepLimit = senderOptionsObj.StepLimit
	opts.TxDataSizeLimit = senderOptionsObj.TxDataSizeLimit

	threshold := new(big.Int)
	valueInt, ok := threshold.SetString(senderOptionsObj.BalanceThreshold, 10)
	if !ok {
		return errors.New("Can't parse field Balance Threshold")
	} else{
		opts.BalanceThreshold = *valueInt
	}

	return nil
}

func (opts *senderOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

type sender struct {
	log  log.Logger
	w    wallet.Wallet
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts senderOptions
	cl   *Client
}

func hexInt2Uint64(hi HexInt) uint64 {
	v, _ := hi.Value()
	return uint64(v)
}

// BMCLinkStatus
// Returns the BMCLinkStatus for "src" link
func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	p := &CallParam{
		FromAddress: Address(s.w.Address()),
		ToAddress:   Address(s.dst.ContractAddress()),
		DataType:    "call",
		Data: CallData{
			Method: BMCGetStatusMethod,
			Params: BMCStatusParams{
				Target: s.src.String(),
			},
		},
	}
	bs := &BMCStatus{}
	err := mapError(s.cl.Call(p, bs))
	if err != nil {
		return nil, err
	}
	ls := &chain.BMCLinkStatus{}
	ls.TxSeq = hexInt2Uint64(bs.TxSeq)
	ls.RxSeq = hexInt2Uint64(bs.RxSeq)
	ls.BMRIndex = uint(hexInt2Uint64(bs.BMRIndex))
	ls.RotateHeight = hexInt2Uint64(bs.RotateHeight)
	ls.RotateTerm = uint(hexInt2Uint64(bs.RotateTerm))
	ls.DelayLimit = uint(hexInt2Uint64(bs.DelayLimit))
	ls.MaxAggregation = uint(hexInt2Uint64(bs.MaxAggregation))
	ls.CurrentHeight = hexInt2Uint64(bs.CurrentHeight)
	ls.RxHeight = hexInt2Uint64(bs.RxHeight)
	ls.RxHeightSrc = hexInt2Uint64(bs.RxHeightSrc)
	return ls, nil
}

func (s *sender) Segment(
	ctx context.Context, msg *chain.Message,
) (tx chain.RelayTx, newMsg *chain.Message, err error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
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

	tx, err = s.newRelayTx(ctx, msg.From.String(), message)
	if err != nil {
		return nil, nil, err
	}

	return tx, newMsg, nil
}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	bal, err := s.cl.GetBalance(&AddressParam{Address: Address(s.w.Address())})
	return bal, &s.opts.BalanceThreshold, err
}

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte) (*relayTx, error) {
	txParam := &TransactionParam{
		Version:     NewHexInt(JsonrpcApiVersion),
		FromAddress: Address(s.w.Address()),
		ToAddress:   Address(s.dst.ContractAddress()),
		NetworkID:   HexInt(s.dst.NetworkID()),
		StepLimit:   NewHexInt(int64(defaultStepLimit)),
		DataType:    "call",
		Data: CallData{
			Method: BMCRelayMethod,
			Params: BMCRelayMethodParams{
				Prev:     prev,
				Messages: base64.URLEncoding.EncodeToString(message),
			},
		},
	}
	if s.opts.StepLimit > 0 {
		txParam.StepLimit = NewHexInt(int64(s.opts.StepLimit))
	}
	return &relayTx{
		Prev:    prev,
		Message: message,
		txParam: txParam,
		cl:      s.cl,
		w:       s.w,
	}, nil
}

type relayTx struct {
	Prev    string `json:"_prev"`
	Message []byte `json:"_msg"`

	txParam     *TransactionParam
	txHashParam *TransactionHashParam
	cl          *Client
	w           wallet.Wallet
}

func (tx *relayTx) ID() interface{} {
	if tx.txHashParam != nil {
		return tx.txHashParam.Hash
	}
	return nil
}

func (tx *relayTx) Send(ctx context.Context) error {
	tx.cl.log.WithFields(log.Fields{
		"prev": tx.Prev}).Debug("handleRelayMessage: send tx")

SignLoop:
	for {
		if err := tx.cl.SignTransaction(tx.w, tx.txParam); err != nil {
			return err
		}
	SendLoop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			txh, err := tx.cl.SendTransaction(tx.txParam)
			if txh != nil {
				tx.txHashParam = &TransactionHashParam{*txh}
				// tx.cl.log.WithFields(log.Fields{
				// 	"txh": tx.txHashParam.Hash,
				// 	"msg": common.HexBytes(tx.Message)}).Debug("handleRelayMessage: tx sent")
				txBytes, _ := json.Marshal(tx.txParam)
				tx.cl.log.WithFields(log.Fields{
					"txh": tx.txHashParam.Hash,
					"tx":  string(txBytes)}).Debug("handleRelayMessage: tx sent")

			}
			if err != nil {
				tx.cl.log.WithFields(log.Fields{
					"error": err}).Debug("handleRelayMessage: send tx")
				if je, ok := err.(*jsonrpc.Error); ok {
					switch je.Code {
					case JsonrpcErrorCodeTxPoolOverflow:
						<-time.After(defaultRelayReSendInterval)
						continue SendLoop
					case JsonrpcErrorCodeSystem:
						if subEc, err := strconv.ParseInt(je.Message[1:5], 0, 32); err == nil {
							switch subEc {
							case DuplicateTransactionError:
								return nil
							case ExpiredTransactionError:
								continue SignLoop
							}
						}
					}
				}
				return mapError(err)
			}
			return nil
		}
	}
}

func (tx *relayTx) Receipt(ctx context.Context) (blockHeight uint64, err error) {
	if tx.txHashParam == nil {
		return 0, fmt.Errorf("no pending tx")
	}
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		txr, err := tx.cl.GetTransactionResult(tx.txHashParam)
		if err != nil {
			if je, ok := err.(*jsonrpc.Error); ok {
				switch je.Code {
				case JsonrpcErrorCodePending, JsonrpcErrorCodeExecuting:
					time.Sleep(defaultGetRelayResultInterval)
					continue
				}
			}
			return 0, mapErrorWithTransactionResult(txr, err)
		}
		tx.cl.log.WithFields(log.Fields{
			"txh": tx.txHashParam.Hash}).Debug("handleRelayMessage: success")
		height, _ := txr.BlockHeight.Value()
		return uint64(height), nil
	}
}

func mapError(err error) error {
	if err != nil {
		switch re := err.(type) {
		case *jsonrpc.Error:
			//fmt.Printf("jrResp.Error:%+v", re)
			switch re.Code {
			case JsonrpcErrorCodeTxPoolOverflow:
				return ErrSendFailByOverflow
			case JsonrpcErrorCodeSystem:
				if subEc, err := strconv.ParseInt(re.Message[1:5], 0, 32); err == nil {
					//TODO return JsonRPC Error
					switch subEc {
					case ExpiredTransactionError:
						return ErrSendFailByExpired
					case FutureTransactionError:
						return ErrSendFailByFuture
					case TransactionPoolOverflowError:
						return ErrSendFailByOverflow
					}
				}
			case JsonrpcErrorCodePending, JsonrpcErrorCodeExecuting:
				return ErrGetResultFailByPending
			}
		case *common.HttpError:
			fmt.Printf("*common.HttpError:%+v", re)
			return ErrConnectFail
		case *url.Error:
			if common.IsConnectRefusedError(re.Err) {
				fmt.Printf("*url.Error:%+v", re)
				return ErrConnectFail
			}
		}
	}
	return err
}

func mapErrorWithTransactionResult(txr *TransactionResult, err error) error {
	err = mapError(err)
	if err == nil && txr != nil && txr.Status != ResultStatusSuccess {
		fc, _ := txr.Failure.CodeValue.Value()
		if fc < ResultStatusFailureCodeRevert || fc > ResultStatusFailureCodeEnd {
			err = fmt.Errorf("failure with code:%s, message:%s",
				txr.Failure.CodeValue, txr.Failure.MessageValue)
		} else {
			err = NewRevertError(int(fc - ResultStatusFailureCodeRevert))
		}
	}
	return err
}
