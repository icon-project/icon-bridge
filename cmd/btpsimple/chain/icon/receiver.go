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
	"bytes"
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature      = "Message(str,int,bytes)"
	EventIndexSignature = 0
	EventIndexNext      = 1
	EventIndexSequence  = 2
)
const MAX_RETRY = 3

type receiverOptions struct {
	Verifier *VerifierOptions `json:"verifier"`
}

func (opts *receiverOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

type eventLogRawFilter struct {
	addr      []byte
	signature []byte
	next      []byte
	seq       uint64
}

type receiver struct {
	log             log.Logger
	src             chain.BTPAddress
	dst             chain.BTPAddress
	cl              *client
	evtLogRawFilter *eventLogRawFilter
	evtReq          *BlockRequest
	hv              *headerValidator
	retries         int
}

func NewReceiver(src, dst chain.BTPAddress, urls []string, opts map[string]interface{}, l log.Logger) (chain.Receiver, error) {
	if len(urls) == 0 {
		return nil, errors.New("List of Urls is empty")
	}
	client := newClient(urls[0], l)

	var recvOpts receiverOptions
	if err := recvOpts.Unmarshal(opts); err != nil {
		return nil, errors.Wrap(err, "Unmarshal recvOpts Error: ")
	}

	hv := headerValidator{}
	vHash := recvOpts.Verifier.ValidatorsHash
	if vs, err := getValidatorsFromHash(client, vHash); err != nil {
		return nil, errors.Wrap(err, "getValidatorsFromHash; ")
	} else {
		hv.validatorsHash = vHash
		hv.validators = vs
		hv.height = recvOpts.Verifier.BlockHeight
	}

	dstr := dst.String()
	ef := &EventFilter{Addr: Address(src.ContractAddress()), Signature: EventSignature, Indexed: []*string{&dstr}}
	evtReq := &BlockRequest{EventFilters: []*EventFilter{ef}} // fill height later

	efAddr, err := ef.Addr.Value()
	if err != nil {
		return nil, errors.Wrap(err, "Get Value from EventFilter.Addr; ")
	}

	recvr := &receiver{
		log:             l,
		src:             src,
		dst:             dst,
		cl:              client,
		hv:              &hv,
		retries:         0,
		evtReq:          evtReq,
		evtLogRawFilter: &eventLogRawFilter{addr: efAddr, signature: []byte(EventSignature), next: []byte(dstr)}, // fill seq later
	}

	return recvr, nil
}

func (r *receiver) receiveLoop(ctx context.Context, height HexInt, seq uint64, sendCallback func(rs []*chain.Receipt) error) error {
	if err := r.syncVerifier(height); err != nil {
		return errors.Wrap(err, "ReceiveLoop; ")
	}

	return r.cl.MonitorBlock(ctx, r.evtReq,
		func(conn *websocket.Conn, v *BlockNotification) error {
			if header, rps, err := r.verify(v); err != nil {
				return errors.Wrap(err, "ReceiveLoop; Verify: ")
			} else {
				htNum, hterr := v.Height.Value()
				if hterr != nil {
					return errors.Wrapf(err, "ReceiveLoop; Conversion Error at Height %v ", v.Height)
				}
				validatorForNextBlock, err := r.getNewValidatorState(header)
				if err != nil {
					return errors.Wrap(err, "ReceiveLoop; ")
				}
				seqForNextEventLog := r.evtLogRawFilter.seq
				if len(rps) > 0 {
					seqForNextEventLog, err = r.getVerifiedSequenceNum(r.evtLogRawFilter.seq, rps)
					if err != nil {
						return errors.Wrap(err, "ReceiveLoop; Verify Sequence: ")
					}
					if err := sendCallback(rps); err != nil {
						return errors.Wrap(err, "ReceiveLoop; sendCallback: ")
					}
					r.log.WithFields(log.Fields{"CurHeight": r.evtReq.Height, "ReceiptLength": len(rps), "CurSeq": r.evtLogRawFilter.seq}).Info(" Receipts Sent ")
				}
				// Update state now that receipts (if any) has been sent
				// Since there isn't any error in the next code segment, state update will happen for sure
				// As such, there won't be the case when receipt is sent but state is not updated
				r.hv = validatorForNextBlock
				r.evtLogRawFilter.seq = seqForNextEventLog
				r.evtReq.Height = NewHexInt(int64(htNum + 1))

				r.log.WithFields(log.Fields{"NextHeight": r.evtReq.Height, "NextSeq": r.evtLogRawFilter.seq}).Debug(" Done: Verified Receipt and Updated State")
			}
			return nil
		},
		func(conn *websocket.Conn) {
			r.log.WithFields(log.Fields{"local": conn.LocalAddr().String()}).Debug("connected")
			if r.retries > 0 {
				r.log.WithFields(log.Fields{"Previous Retry Count": r.retries}).Debug("Reset to zero")
				r.retries = 0
			}
		},
		func(conn *websocket.Conn, err error) {
			r.log.WithFields(log.Fields{"error": err, "local": conn.LocalAddr().String()}).Info("disconnected")
			_ = conn.Close()
		})
}

func (r *receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	if opts.Height < 1 {
		return nil, errors.New("Height of BlockChain should be positive number")
	}
	r.evtReq.Height = NewHexInt(int64(opts.Height))
	r.evtLogRawFilter.seq = opts.Seq //common.NewHexInt(int64(opts.Seq)).Bytes()

	_errCh := make(chan error)
	go func() {
		var err error
		defer close(_errCh)
		cb := func(rs []*chain.Receipt) error {
			msgCh <- &chain.Message{Receipts: rs}
			return nil
		}
	RetryIfEOF:
		if err = r.receiveLoop(ctx, r.evtReq.Height, r.evtLogRawFilter.seq, cb); err != nil {
			if isUnexpectedEOFError(err) && r.retries < MAX_RETRY {
				r.retries++
				r.log.WithFields(log.Fields{
					"Retry Count ":       r.retries,
					"EventRequestHeight": r.evtReq.Height, "EventSequence": r.evtLogRawFilter.seq,
					"ValidatorHash": r.hv.validatorsHash, "ValidatorHeight": NewHexInt(int64(r.hv.height)),
				}).Warn("Retrying Websocket Connection")
				goto RetryIfEOF
			} else {
				r.log.WithFields(log.Fields{
					"EventRequestHeight": r.evtReq.Height, "EventSequence": r.evtLogRawFilter.seq,
					"ValidatorHash": r.hv.validatorsHash, "ValidatorHeight": NewHexInt(int64(r.hv.height)),
				}).Warn("State Info Before returning error ")
				_errCh <- err
			}
		}
		r.log.Warnf("Receive Loop Terminated; err %+v", err)
	}()
	return _errCh, nil
}

func (r *receiver) getNewValidatorState(header *BlockHeader) (*headerValidator, error) {
	nhv := &headerValidator{
		validators:     r.hv.validators,
		validatorsHash: r.hv.validatorsHash,
		height:         r.hv.height + 1, // point to the next block
	}
	if bytes.Equal(header.NextValidatorsHash, r.hv.validatorsHash) { // If same validatorHash, only update height to point to the next block
		return nhv, nil
	}
	r.log.WithFields(log.Fields{"Height": NewHexInt(header.Height), "NewValidatorHash": common.HexBytes(header.NextValidatorsHash), "OldValidatorHash": r.hv.validatorsHash}).Info(" Updating Validator Hash ")
	if vs, err := getValidatorsFromHash(r.cl, header.NextValidatorsHash); err != nil {
		return nil, errors.Wrap(err, "verifyHeader; ")
	} else {
		nhv.validatorsHash = header.NextValidatorsHash
		nhv.validators = vs
	}
	return nhv, nil
}

func (r *receiver) getVerifiedSequenceNum(expectedSeq uint64, receipts []*chain.Receipt) (uint64, error) {
	for _, receipt := range receipts {
		newEvents := []*chain.Event{}
		for _, event := range receipt.Events {
			switch {
			case event.Sequence == expectedSeq:
				newEvents = append(newEvents, event)
				expectedSeq++
			case event.Sequence > expectedSeq: // event.sequencce - expectedSeq has not been considered or is missed ?
				r.log.WithFields(log.Fields{"Expected": expectedSeq, "Got": event.Sequence}).Error("Current event log sequence higher than expected")
				return expectedSeq, errors.New("Invalid Sequence for event log of receipt ")
			}
		}
		receipt.Events = newEvents
	}
	return expectedSeq, nil
}
