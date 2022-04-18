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
	"encoding/json"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/btp/cmd/btpsimple/module"
	"github.com/icon-project/btp/cmd/btpsimple/module/bsc/binding"

	"math/big"

	"github.com/icon-project/btp/common/log"
)

const (
	EPOCH = 200
)

type receiver struct {
	c   *Client
	src module.BtpAddress
	dst module.BtpAddress
	log log.Logger
	opt struct {
	}
	evtReq             *BlockRequest
	isFoundOffsetBySeq bool
}

func (r *receiver) newReceiptProofs(v *BlockNotification) ([]*module.ReceiptProof, error) {
	rps := make([]*module.ReceiptProof, 0)

	block, err := r.c.GetBlockByHeight(v.Height)
	if err != nil {
		return nil, err
	}

	if len(block.Transactions()) == 0 {
		return rps, nil
	}

	receipts, err := r.c.GetBlockReceipts(block)
	if err != nil {
		return nil, err
	}

	if block.GasUsed() == 0 {
		r.log.Println("Block %s has 0 gas", block.Number(), len(block.Transactions()))
		return rps, nil
	}

	srcContractAddress := HexToAddress(r.src.ContractAddress())

	for _, receipt := range receipts {
		rp := &module.ReceiptProof{}

		for _, eventLog := range receipt.Logs {
			if eventLog.Address != srcContractAddress {
				continue
			}

			if bmcMsg, err := binding.UnpackEventLog(eventLog.Data); err == nil {
				rp.Events = append(rp.Events, &module.Event{
					Message:  bmcMsg.Msg,
					Next:     module.BtpAddress(bmcMsg.Next),
					Sequence: bmcMsg.Seq.Int64(),
				})
			}
		}

		if len(rp.Events) > 0 {
			rp.Index = int(receipt.TransactionIndex)
			rp.Height = v.Height.Int64()
			rps = append(rps, rp)
			r.log.Debugf("event found for height & address:", rp.Height, srcContractAddress)
			r.isFoundOffsetBySeq = true
		}
	}
	return rps, nil
}

func (r *receiver) newBTPMessage(v *BlockNotification) ([]*module.ReceiptProof, error) {
	rps := make([]*module.ReceiptProof, 0)

	srcContractAddress := HexToAddress(r.src.ContractAddress())

	query := ethereum.FilterQuery{
		FromBlock: v.Height,
		ToBlock:   v.Height,
		Addresses: []common.Address{
			srcContractAddress,
		},
	}

	logs, err := r.c.FilterLogs(query)
	if err != nil {
		return nil, err
	}

	for _, vLog := range logs {
		rp := &module.ReceiptProof{}
		if bmcMsg, err := binding.UnpackEventLog(vLog.Data); err == nil {
			rp.Events = append(rp.Events, &module.Event{
				Message:  bmcMsg.Msg,
				Next:     module.BtpAddress(bmcMsg.Next),
				Sequence: bmcMsg.Seq.Int64(),
			})
		}

		if len(rp.Events) > 0 {
			rp.Index = int(vLog.TxIndex)
			rp.Height = int64(vLog.BlockNumber)
			rps = append(rps, rp)
			r.log.Debugf("event found for height & address:", rp.Height, srcContractAddress)
			r.isFoundOffsetBySeq = true
		}
	}
	return rps, nil
}

func (r *receiver) ReceiveLoop(height int64, seq int64, cb module.ReceiveCallback, scb func()) error {
	r.log.Debugf("ReceiveLoop connected")
	br := &BlockRequest{
		Height: big.NewInt(height),
	}
	var err error
	if seq < 1 {
		r.isFoundOffsetBySeq = true
	}
	if err != nil {
		r.log.Errorf(err.Error())
	}
	return r.c.MonitorBlock(br,
		func(v *BlockNotification) error {
			r.log.Debugf("onBlockOfSrc BSC %d", v.Height.Int64())
			var rps []*module.ReceiptProof
			if rps, err = r.newBTPMessage(v); err != nil {
				return err
			} else if r.isFoundOffsetBySeq {
				cb(rps)
			}
			return nil
		},
	)
}

func (r *receiver) StopReceiveLoop() {
	r.c.CloseAllMonitor()
}

func NewReceiver(src, dst module.BtpAddress, endpoint string, opt map[string]interface{}, l log.Logger) module.Receiver {
	r := &receiver{
		src: src,
		dst: dst,
		log: l,
	}
	b, err := json.Marshal(opt)
	if err != nil {
		l.Panicf("fail to marshal opt:%#v err:%+v", opt, err)
	}
	if err = json.Unmarshal(b, &r.opt); err != nil {
		l.Panicf("fail to unmarshal opt:%#v err:%+v", opt, err)
	}
	r.c = NewClient(endpoint, l)
	return r
}
