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
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature      = "Message(str,int,bytes)"
	EventIndexSignature = 0
	EventIndexNext      = 1
	EventIndexSequence  = 2
)

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

type receiver struct {
	log  log.Logger
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts receiverOptions
	cl   *client

	// src
	evtLogRawFilter struct {
		addr      []byte
		signature []byte
		next      []byte
		seq       []byte
	}
	evtReq             *BlockRequest
	isFoundOffsetBySeq bool

	//cached data that verifier uses and updates
	validatorHash []byte
	validators    [][]byte
}

// NewReceiver ...
// returns a new receiver client for harmony
func NewReceiver(
	src, dst chain.BTPAddress, urls []string,
	opts map[string]interface{}, l log.Logger) (chain.Receiver, error) {
	cl := &receiver{
		log: l,
		src: src,
		dst: dst,
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}
	if err := cl.opts.Unmarshal(opts); err != nil {
		return nil, err
	}
	cl.cl = newClient(urls[0], l)
	return cl, nil
}

func (r *receiver) getBlockHeader(height HexInt) (*BlockHeader, error) {
	p := &BlockHeightParam{Height: height}
	b, err := r.cl.GetBlockHeaderByHeight(p)
	if err != nil {
		return nil, mapError(err)
	}
	var bh BlockHeader
	_, err = codec.RLP.UnmarshalFromBytes(b, &bh)
	if err != nil {
		return nil, err
	}
	bh.serialized = b
	return &bh, nil
}

func (r *receiver) receiveLoop(
	ctx context.Context, height, seq uint64,
	callback func(rs []*chain.Receipt) error) error {
	s := r.dst.String()
	ef := &EventFilter{
		Addr:      Address(r.src.ContractAddress()),
		Signature: EventSignature,
		Indexed:   []*string{&s},
	}
	r.evtReq = &BlockRequest{
		Height:       NewHexInt(int64(height)),
		EventFilters: []*EventFilter{ef},
	}

	if height < 1 {
		return fmt.Errorf("cannot catchup from zero height")
	}
	heightTracker := height // track processed block height

	var err error

	vHash := r.opts.Verifier.ValidatorsHash
	vBytes, err := r.cl.GetDataByHash(&DataHashParam{Hash: NewHexBytes(vHash)})
	if err != nil {
		return errors.Wrap(err, "receiveLoop; GetDataByHash Validators; Err: ")
	}
	var vs [][]byte
	_, err = vlcodec.BC.UnmarshalFromBytes(vBytes, &vs)
	if err != nil {
		return errors.Wrap(err, "receiveLoop; Unmarshal ValidatorBytes; Err: ")
	}
	r.validatorHash = vHash
	r.validators = vs

	if seq < 1 {
		r.isFoundOffsetBySeq = true
	}
	if r.evtLogRawFilter.addr, err = ef.Addr.Value(); err != nil {
		r.log.Panicf("ef.Addr.Value() err:%+v", err)
	}
	r.evtLogRawFilter.signature = []byte(EventSignature)
	r.evtLogRawFilter.next = []byte(s)
	r.evtLogRawFilter.seq = common.NewHexInt(int64(seq)).Bytes()

	retryCount := 0
	retryLimit := 3
Retry:
	monErr := r.cl.MonitorBlock(r.evtReq,
		func(conn *websocket.Conn, v *BlockNotification) error {
			var err error
			var rps []*chain.Receipt
			if rps, err = r.blockVerification(v); err != nil {
				return err
			} else {
				if r.isFoundOffsetBySeq && rps != nil && len(rps) > 0 {
					r.log.WithFields(log.Fields{"Height": v.Height, "Length": len(rps)}).Debug("Receipt Verified")
					callback(rps)
				}
				if ht, hterr := v.Height.Value(); hterr == nil {
					heightTracker = uint64(ht) // save block height fetched so far; so it can be used later in case of reconnection attempts
				} else {
					r.log.WithFields(log.Fields{"heightTracker": heightTracker}).Error(errors.Wrap(hterr, "Conversion Error: "))
				}
			}
			return nil
		},
		func(conn *websocket.Conn) {
			r.log.WithFields(log.Fields{"local": conn.LocalAddr().String()}).Debug("connected")
			if retryCount > 0 {
				r.log.WithFields(log.Fields{"Previous Retry Count": retryCount}).Debug("Reset to zero")
				retryCount = 0
			}
		},
		func(conn *websocket.Conn, err error) {
			r.log.WithFields(log.Fields{"error": err, "local": conn.LocalAddr().String()}).Debug("disconnected")
			_ = conn.Close()
		})
	if isUnexpectedEOFError(monErr) && retryCount < retryLimit {
		retryCount++
		r.evtReq.Height = NewHexInt(int64(heightTracker) + 1) // Resume from the block height processed so far
		r.log.WithFields(log.Fields{"Retry Count ": retryCount, "Resuming Height": r.evtReq.Height}).Info("Retrying Websocket Connection")
		goto Retry
	} // If connection is not re-established even after retrying retryLimit times in a row for an EOF error type, return

	return monErr
}

// Websocket connection is closed by peer abruptly with EOF message. The function checks and verifies if the error thrown is unexpected EOF
func isUnexpectedEOFError(err error) bool {
	//websocket/conn.go 	errUnexpectedEOF       = &CloseError{Code: CloseAbnormalClosure, Text: io.ErrUnexpectedEOF.Error()}
	if cErr, ok := err.(*websocket.CloseError); ok && cErr.Code == websocket.CloseAbnormalClosure && cErr.Text == io.ErrUnexpectedEOF.Error() {
		return true
	} else if err.Error() == io.ErrUnexpectedEOF.Error() {
		return true
	}
	return false
}

func (r *receiver) Subscribe(
	ctx context.Context, msgCh chan<- *chain.Message,
	opts chain.SubscribeOptions) (errCh <-chan error, err error) {

	_errCh := make(chan error)
	go func() {
		defer close(_errCh)
		if err := r.receiveLoop(ctx,
			opts.Height, opts.Seq,
			func(rs []*chain.Receipt) error {
				msgCh <- &chain.Message{Receipts: rs}
				return nil
			}); err != nil {
			// TODO decide whether to ignore or handle err
			r.log.Errorf("receiveLoop terminated: %v", err)
			_errCh <- err
		}
	}()
	return _errCh, nil
}

/*
func toEventLog(proof [][]byte) (*EventLog, error) {
	mp, err := mpt.NewMptProof(proof)
	if err != nil {
		return nil, err
	}
	el := &EventLog{}
	if _, err := codec.RLP.UnmarshalFromBytes(mp.Leaf().Data, el); err != nil {
		return nil, fmt.Errorf("fail to parse EventLog on leaf err:%+v", err)
	}
	return el, nil
}


func (r *receiver) toEvent(proof [][]byte) (*chain.Event, error) {
	el, err := toEventLog(proof)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(el.Addr, r.evtLogRawFilter.addr) &&
		bytes.Equal(el.Indexed[EventIndexSignature], r.evtLogRawFilter.signature) &&
		bytes.Equal(el.Indexed[EventIndexNext], r.evtLogRawFilter.next) {
		var i common.HexInt
		i.SetBytes(el.Indexed[EventIndexSequence])
		evt := &chain.Event{
			Next:     chain.BTPAddress(el.Indexed[EventIndexNext]),
			Sequence: i.Uint64(),
			Message:  el.Data[0],
		}
		return evt, nil
	}
	return nil, fmt.Errorf("invalid event")
}

func (r *receiver) getReceipts(v *BlockNotification) ([]*chain.Receipt, error) {
	nextEp := 0
	rps := make([]*chain.Receipt, 0)
	fmt.Println("Get Receipts ", len(v.Indexes))
	if len(v.Indexes) > 0 {
		l := v.Indexes[0]
	RpLoop:
		for i, index := range l {
			p := &ProofEventsParam{BlockHash: v.Hash, Index: index, Events: v.Events[0][i]}
			proofs, err := r.cl.GetProofForEvents(p)
			if err != nil {
				return nil, mapError(err)
			}
			if !r.isFoundOffsetBySeq {
			EpLoop:
				for j := 0; j < len(p.Events); j++ {
					if el, err := toEventLog(proofs[j+1]); err != nil {
						return nil, err
					} else if bytes.Equal(el.Addr, r.evtLogRawFilter.addr) &&
						bytes.Equal(el.Indexed[EventIndexSignature], r.evtLogRawFilter.signature) &&
						bytes.Equal(el.Indexed[EventIndexNext], r.evtLogRawFilter.next) &&
						bytes.Equal(el.Indexed[EventIndexSequence], r.evtLogRawFilter.seq) {
						r.isFoundOffsetBySeq = true
						if (j + 1) < len(p.Events) {
							nextEp = j + 1
							break EpLoop
						}
					} else {
						r.log.WithFields(log.Fields{
							"addr": log.Fields{
								"got":      common.HexBytes(el.Addr),
								"expected": common.HexBytes(r.evtLogRawFilter.addr),
							},
							"sig": log.Fields{
								"got":      common.HexBytes(el.Indexed[EventIndexSignature]),
								"expected": common.HexBytes(r.evtLogRawFilter.signature),
							},
							"next": log.Fields{
								"got":      common.HexBytes(el.Indexed[EventIndexNext]),
								"expected": common.HexBytes(r.evtLogRawFilter.next),
							},
							"seq": log.Fields{
								"got":      common.HexBytes(el.Indexed[EventIndexSequence]),
								"expected": common.HexBytes(r.evtLogRawFilter.seq),
							},
						}).Error("invalid event: cannot match addr/sig/next/seq")
					}
				}
				if nextEp == 0 {
					continue RpLoop
				}
			}
			idx, _ := index.Value()
			rp := &chain.Receipt{
				Index: uint64(idx),
			}
			rp.Height = hexInt2Uint64(v.Height)
			for k := nextEp; k < len(p.Events); k++ {
				var evt *chain.Event
				if evt, err = r.toEvent(proofs[k+1]); err != nil {
					return nil, err
				}
				rp.Events = append(rp.Events, evt)
			}
			rps = append(rps, rp)
			nextEp = 0
		}
	}
	return rps, nil
}

func (r *receiver) receiveLoop(
	ctx context.Context, height, seq uint64,
	callback func(rs []*chain.Receipt) error) error {
	s := r.dst.String()
	ef := &EventFilter{
		Addr:      Address(r.src.ContractAddress()),
		Signature: EventSignature,
		Indexed:   []*string{&s},
	}
	r.evtReq = &BlockRequest{
		Height:       NewHexInt(int64(height)),
		EventFilters: []*EventFilter{ef},
	}

	if height < 1 {
		return fmt.Errorf("cannot catchup from zero height")
	}
	var err error
	if r.bh, err = r.getBlockHeader(NewHexInt(int64(height) - 1)); err != nil {
		return err
	}
	if seq < 1 {
		r.isFoundOffsetBySeq = true
	}
	if r.evtLogRawFilter.addr, err = ef.Addr.Value(); err != nil {
		r.log.Panicf("ef.Addr.Value() err:%+v", err)
	}
	r.evtLogRawFilter.signature = []byte(EventSignature)
	r.evtLogRawFilter.next = []byte(s)
	r.evtLogRawFilter.seq = common.NewHexInt(int64(seq)).Bytes()
	return r.cl.MonitorBlock(r.evtReq,
		func(conn *websocket.Conn, v *BlockNotification) error {
			var blockNum, _ = v.Height.Value()
			r.log.WithFields(log.Fields{"height": blockNum}).Debug("block notification")
			var err error
			var rps []*chain.Receipt
			if rps, err = r.getReceipts(v); err != nil {
				return err
			} else if r.isFoundOffsetBySeq {
				callback(rps)
			}
			return nil
		},
		func(conn *websocket.Conn) {
			r.log.WithFields(log.Fields{"local": conn.LocalAddr().String()}).Debug("connected")
		},
		func(conn *websocket.Conn, err error) {
			r.log.WithFields(log.Fields{"error": err, "local": conn.LocalAddr().String()}).Debug("disconnected")
			_ = conn.Close()
		})
}

*/
