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

	"github.com/gorilla/websocket"
	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
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
		return nil, errors.Wrap(err, "Unmarshal Error: ")
	}
	fmt.Println(*recvOpts.Verifier)

	hv := headerValidator{}
	if vHash, err := base64.StdEncoding.DecodeString(recvOpts.Verifier.ValidatorHash); err != nil {
		return nil, errors.Wrap(err, "Base64Decode recvOpts.Verifier.ValidatorHash; Err: ")
	} else {
		if vs, err := getValidatorsFromHash(client, vHash); err != nil {
			return nil, errors.Wrap(err, "getValidatorsFromHash; ")
		} else {
			hv.validatorHash = vHash
			hv.validators = vs
			hv.height = recvOpts.Verifier.BlockHeight
		}
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

func (r *receiver) receiveLoop(ctx context.Context, height HexInt, seq uint64, callback func(rs []*chain.Receipt) error) error {
	if err := r.syncVerifier(height); err != nil {
		return errors.Wrap(err, "ReceiveLoop; ")
	}
	return r.cl.MonitorBlock(ctx, r.evtReq,
		func(conn *websocket.Conn, v *BlockNotification) error {
			if rps, err := r.verify(v); err != nil {
				return errors.Wrap(err, "ReceiveLoop; Verify: ")
			} else {
				htNum, hterr := v.Height.Value()
				if hterr != nil {
					return errors.Wrapf(err, "ReceiveLoop; Conversion Error at Height %v ", v.Height)
				}
				if len(rps) > 0 {
					r.log.WithFields(log.Fields{"Height": v.Height, "Length": len(rps)}).Debug("Receipt Verified")
					if err := r.verifySequence(rps); err != nil {
						return errors.Wrap(err, "ReceiveLoop; Verify Sequence: ")
					} else {
						callback(rps)
					}
				}
				// update height state to point to the next block; update even if len(rps) == 0; no receipts in BlockNotification
				r.evtReq.Height = NewHexInt(int64(htNum + 1))
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

func (r *receiver) verifySequence(receipts []*chain.Receipt) error {
	for _, receipt := range receipts {
		newEvents := []*chain.Event{}
		for _, event := range receipt.Events {
			switch {
			case event.Sequence == r.evtLogRawFilter.seq:
				newEvents = append(newEvents, event)
				r.evtLogRawFilter.seq++
			case event.Sequence > r.evtLogRawFilter.seq:
				return errors.New("Invalid Sequence")
			}
		}
		receipt.Events = newEvents
	}
	return nil
}

func (r *receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	if opts.Height < 1 {
		return nil, errors.New("Height of BlockChain should be positive number")
	}
	r.evtReq.Height = NewHexInt(int64(opts.Height))
	r.evtLogRawFilter.seq = opts.Seq //common.NewHexInt(int64(opts.Seq)).Bytes()

	_errCh := make(chan error)
	go func() {
		defer close(_errCh)
		cb := func(rs []*chain.Receipt) error {
			msgCh <- &chain.Message{Receipts: rs}
			return nil
		}
	RetryIfEOF:
		if err := r.receiveLoop(ctx, r.evtReq.Height, r.evtLogRawFilter.seq, cb); err != nil {
			if isUnexpectedEOFError(err) && r.retries < MAX_RETRY {
				r.retries++
				r.log.WithFields(log.Fields{"Retry Count ": r.retries, "Resuming Height": r.evtReq.Height}).Info("Retrying Websocket Connection")
				goto RetryIfEOF
			} else {
				_errCh <- err
			}
		}
	}()
	return _errCh, nil
}
