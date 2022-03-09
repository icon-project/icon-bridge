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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/btp/cmd/btpsimple/module/bsc/binding"

	"github.com/icon-project/btp/cmd/btpsimple/module"
	"github.com/icon-project/btp/common"
	"github.com/icon-project/btp/common/codec"
	"github.com/icon-project/btp/common/jsonrpc"
	"github.com/icon-project/btp/common/log"
)

const (
	txMaxDataSize                 = 32768 //512 * 1024 // 512kB
	txOverheadScale               = 0.37  //base64 encoding overhead 0.36, rlp and other fields 0.01
	txSizeLimit                   = txMaxDataSize / (1 + txOverheadScale)
	DefaultGetRelayResultInterval = time.Second
	DefaultRelayReSendInterval    = time.Second
)

type sender struct {
	c   *Client
	src module.BtpAddress
	dst module.BtpAddress
	w   Wallet
	l   log.Logger
	opt struct {
	}

	bmc *binding.BMC

	evtLogRawFilter struct {
		addr      []byte
		signature []byte
		next      []byte
		seq       []byte
	}
	evtReq             *BlockRequest
	isFoundOffsetBySeq bool
	cb                 module.ReceiveCallback

	mutex sync.Mutex
}

func (s *sender) newTransactionParam(prev string, rm *RelayMessage) (*TransactionParam, error) {
	b, err := codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, err
	}
	rmp := BMCRelayMethodParams{
		Prev: prev,
		//Messages: base64.URLEncoding.EncodeToString(b),
		Messages: string(b[:]),
	}
	p := &TransactionParam{
		Params: rmp,
	}
	return p, nil
}

func (s *sender) Segment(rm *module.RelayMessage, height int64) ([]*module.Segment, error) {
	segments := make([]*module.Segment, 0)
	var err error
	msg := &RelayMessage{
		BlockUpdates:  make([][]byte, 0),
		ReceiptProofs: make([][]byte, 0),
	}
	//size := 0
	var b []byte
	for _, rp := range rm.ReceiptProofs {
		//TODO: segment for the events
		/* if s.isOverLimit(len(rp.Proof)) {
			return nil, fmt.Errorf("invalid ReceiptProof.Proof size")
		} */
		/* if len(msg.BlockUpdates) == 0 {
			size += len(bp)
			msg.BlockProof = bp
			msg.height = rm.BlockProof.BlockWitness.Height
		} */
		//size += len(rp.Proof)
		var eventBytes []byte
		if eventBytes, err = codec.RLP.MarshalToBytes(rp.Events); err != nil {
			return nil, err
		}
		trp := &ReceiptProof{
			Index:  rp.Index,
			Events: eventBytes,
		}
		/* for j, ep := range rp.EventProofs {
			if s.isOverLimit(len(ep.Proof)) {
				return nil, fmt.Errorf("invalid EventProof.Proof size")
			}
			size += len(ep.Proof)
			if s.isOverLimit(size) {
				if j == 0 && len(msg.BlockUpdates) == 0 {
					return nil, fmt.Errorf("BlockProof + ReceiptProof + EventProof > limit")
				}
				//
				segment := &module.Segment{
					Height:              msg.height,
					NumberOfBlockUpdate: msg.numberOfBlockUpdate,
					EventSequence:       msg.eventSequence,
					NumberOfEvent:       msg.numberOfEvent,
				}
				if segment.TransactionParam, err = s.newTransactionParam(rm.From.String(), msg); err != nil {
					return nil, err
				}
				segments = append(segments, segment)

				msg = &RelayMessage{
					BlockUpdates:  make([][]byte, 0),
					ReceiptProofs: make([][]byte, 0),
					BlockProof:    bp,
				}
				size = len(ep.Proof)
				size += len(rp.Proof)
				size += len(bp)

				trp = &ReceiptProof{
					Index:       rp.Index,
					Proof:       rp.Proof,
					EventProofs: make([]*module.EventProof, 0),
				}
			}
			trp.EventProofs = append(trp.EventProofs, ep)
			msg.eventSequence = rp.Events[j].Sequence
			msg.numberOfEvent += 1
		} */

		if b, err = codec.RLP.MarshalToBytes(trp); err != nil {
			return nil, err
		}
		msg.ReceiptProofs = append(msg.ReceiptProofs, b)

	}
	//
	segment := &module.Segment{
		Height:              msg.height,
		NumberOfBlockUpdate: msg.numberOfBlockUpdate,
		EventSequence:       msg.eventSequence,
		NumberOfEvent:       msg.numberOfEvent,
	}
	if segment.TransactionParam, err = s.newTransactionParam(rm.From.String(), msg); err != nil {
		return nil, err
	}
	s.l.Debugf("Segmentation Done")
	segments = append(segments, segment)
	return segments, nil
}

func (s *sender) Relay(segment *module.Segment) (module.GetResultParam, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p, ok := segment.TransactionParam.(*TransactionParam)
	if !ok {
		return nil, fmt.Errorf("casting failure")
	}
	t, err := s.c.newTransactOpts(s.w)
	if err != nil {
		return nil, err
	}
	rmp := p.Params.(BMCRelayMethodParams)
	var tx *types.Transaction
	tx, err = s.bmc.HandleRelayMessage(t, rmp.Prev, rmp.Messages)
	if err != nil {
		s.l.Errorf("handleRelayMessage: ", err.Error())
		return nil, err
	}
	thp := &TransactionHashParam{}
	thp.Hash = tx.Hash()
	s.l.Debugf("HandleRelayMessage tx hash:%s, prev %s, msg: %s", thp.Hash, rmp.Prev, base64.URLEncoding.EncodeToString([]byte(rmp.Messages)))
	return thp, nil
}

func (s *sender) GetResult(p module.GetResultParam) (module.TransactionResult, error) {
	if txh, ok := p.(*TransactionHashParam); ok {
		for {
			_, pending, err := s.c.GetTransaction(txh.Hash)
			if err != nil {
				return nil, err
			}
			if pending {
				<-time.After(DefaultGetRelayResultInterval)
				continue
			}
			tx, err := s.c.GetTransactionReceipt(txh.Hash)
			if err != nil {
				return nil, err
			}
			return tx, nil //mapErrorWithTransactionResult(&types.Receipt{}, err) // TODO: map transaction.js result error
		}
	} else {
		return nil, fmt.Errorf("fail to casting TransactionHashParam %T", p)
	}
}

func (s *sender) GetStatus() (*module.BMCLinkStatus, error) {
	var status binding.TypesLinkStats
	status, err := s.bmc.GetStatus(nil, s.src.String())

	if err != nil {
		s.l.Errorf("Error retrieving relay status from BMC")
		return nil, err
	}

	ls := &module.BMCLinkStatus{}
	ls.TxSeq = status.TxSeq.Int64()
	ls.RxSeq = status.RxSeq.Int64()
	ls.BMRIndex = int(status.RelayIdx.Int64())
	ls.RotateHeight = status.RotateHeight.Int64()
	ls.RotateTerm = int(status.RotateTerm.Int64())
	ls.DelayLimit = int(status.DelayLimit.Int64())
	ls.MaxAggregation = int(status.MaxAggregation.Int64())
	ls.CurrentHeight = status.CurrentHeight.Int64()
	ls.RxHeight = status.RxHeight.Int64()
	ls.RxHeightSrc = status.RxHeightSrc.Int64()
	return ls, nil
}

func (s *sender) isOverLimit(size int) bool {
	return txSizeLimit < float64(size)
}

func (s *sender) MonitorLoop(height int64, cb module.MonitorCallback, scb func()) error {
	s.l.Debugf("MonitorLoop (sender) connected")
	br := &BlockRequest{
		Height: big.NewInt(height),
	}
	return s.c.MonitorBlock(br,
		func(v *BlockNotification) error {
			return cb(v.Height.Int64())
		})
}

func (s *sender) StopMonitorLoop() {
	s.c.CloseAllMonitor()
}
func (s *sender) FinalizeLatency() int {
	//on-the-next
	return 1
}

func NewSender(src, dst module.BtpAddress, w Wallet, endpoint string, opt map[string]interface{}, l log.Logger) module.Sender {
	s := &sender{
		src: src,
		dst: dst,
		w:   w,
		l:   l,
	}
	b, err := json.Marshal(opt)
	if err != nil {
		l.Panicf("fail to marshal opt:%#v err:%+v", opt, err)
	}
	if err = json.Unmarshal(b, &s.opt); err != nil {
		l.Panicf("fail to unmarshal opt:%#v err:%+v", opt, err)
	}
	s.c = NewClient(endpoint, l)

	s.bmc, _ = binding.NewBMC(HexToAddress(s.dst.ContractAddress()), s.c.ethClient)

	return s
}

func mapError(err error) error {
	if err != nil {
		switch re := err.(type) {
		case *jsonrpc.Error:
			//fmt.Printf("jrResp.Error:%+v", re)
			switch re.Code {
			case JsonrpcErrorCodeTxPoolOverflow:
				return module.ErrSendFailByOverflow
			case JsonrpcErrorCodeSystem:
				if subEc, err := strconv.ParseInt(re.Message[1:5], 0, 32); err == nil {
					//TODO return JsonRPC Error
					switch subEc {
					case ExpiredTransactionError:
						return module.ErrSendFailByExpired
					case FutureTransactionError:
						return module.ErrSendFailByFuture
					case TransactionPoolOverflowError:
						return module.ErrSendFailByOverflow
					}
				}
			case JsonrpcErrorCodePending, JsonrpcErrorCodeExecuting:
				return module.ErrGetResultFailByPending
			}
		case *common.HttpError:
			fmt.Printf("*common.HttpError:%+v", re)
			return module.ErrConnectFail
		case *url.Error:
			if common.IsConnectRefusedError(re.Err) {
				//fmt.Printf("*url.Error:%+v", re)
				return module.ErrConnectFail
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
			err = module.NewRevertError(int(fc - ResultStatusFailureCodeRevert))
		}
	}
	return err
}
