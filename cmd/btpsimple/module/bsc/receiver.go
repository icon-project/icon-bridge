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
	consensusStates    ConsensusStates
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

	//receiptTrie, err := trieFromReceipts(receipts) // receiptTrie.Hash() == block.ReceiptHash

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

			/* proof, err := codec.RLP.MarshalToBytes(*MakeLog(eventLog))
			if err != nil {
				return nil, err
			}
			rp.EventProofs = append(rp.EventProofs, &module.EventProof{
				Index: int(eventLog.Index),
				Proof: proof,
			}) */
		}

		if len(rp.Events) > 0 {
			/* key, err := rlp.EncodeToBytes(receipt.TransactionIndex)
			r.log.Debugf("newReceiptProofs: height:%d hash:%s key:%d", v.Height, block.ReceiptHash(), key)
			proofs, err := receiptProof(receiptTrie, key)
			if err != nil {
				return nil, err
			} */
			rp.Index = int(receipt.TransactionIndex)
			//rp.Proof, err = codec.RLP.MarshalToBytes(proofs)
			/* if err != nil {
				return nil, err
			} */
			rps = append(rps, rp)
		}
	}
	return rps, nil
}

/*
func trieFromReceipts(receipts []*types.Receipt) (*trie.Trie, error) {
	tr, _ := trie.New(common.Hash{}, trie.NewDatabase(memorydb.New()))

	for i, r := range receipts {
		path, err := rlp.EncodeToBytes(uint(i))

		if err != nil {
			return nil, err
		}

		rawReceipt, err := rlp.EncodeToBytes(r)
		if err != nil {
			return nil, err
		}

		tr.Update(path, rawReceipt)
	}

	_, err := tr.Commit(nil)
	if err != nil {
		return nil, err
	}

	return tr, nil
} */
/*
func receiptProof(receiptTrie *trie.Trie, key []byte) ([][]byte, error) {
	proofSet := light.NewNodeSet()
	err := receiptTrie.Prove(key, 0, proofSet)
	if err != nil {
		return nil, err
	}
	proofs := make([][]byte, 0)
	for _, node := range proofSet.NodeList() {
		fmt.Println(hexutil.Encode(node))
		proofs = append(proofs, node)
	}
	return proofs, nil
}
*/
func (r *receiver) ReceiveLoop(height int64, seq int64, cb module.ReceiveCallback, scb func()) error {
	r.log.Debugf("ReceiveLoop connected")
	br := &BlockRequest{
		Height: big.NewInt(height),
	}
	var err error
	if seq < 1 {
		r.isFoundOffsetBySeq = true
	}
	r.consensusStates, err = r.c.GetLatestConsensusState()
	if err != nil {
		r.log.Errorf(err.Error())
	}
	return r.c.MonitorBlock(br,
		func(v *BlockNotification) error {
			var bu *module.BlockUpdate
			var rps []*module.ReceiptProof
			if rps, err = r.newReceiptProofs(v); err != nil {
				return err
			} else if r.isFoundOffsetBySeq {
				cb(bu, rps)
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
