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
	"fmt"
	"sort"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

const (
	EventSignature      = "Message(str,int,bytes)"
	EventIndexSignature = 0
	EventIndexNext      = 1
	EventIndexSequence  = 2
	RPCCallRetry        = 5
)

const RECONNECT_ON_UNEXPECTED_HEIGHT = "Unexpected Block Height. Should Reconnect"
const (
	MonitorBlockMaxConcurrency = 300
)

type ReceiverOptions struct {
	SyncConcurrency uint64                 `json:"syncConcurrency"`
	Verifier        *types.VerifierOptions `json:"verifier"`
}

type eventLogRawFilter struct {
	addr      []byte
	signature []byte
	next      []byte
	seq       uint64
}

type Receiver struct {
	log       log.Logger
	src       chain.BTPAddress
	dst       chain.BTPAddress
	Client    IClient
	opts      ReceiverOptions
	blockReq  types.BlockRequest
	logFilter eventLogRawFilter
}

type verifierBlockResponse struct {
	Height         int64
	Header         *types.BlockHeader
	Votes          []byte
	NextValidators []common.Address
	err            error
}

type verifierBlockRequest struct {
	height   int64
	err      error
	retry    int
	response *verifierBlockResponse
}

type btpBlockResponse struct {
	Height         int64
	Hash           common.HexHash
	Header         *types.BlockHeader
	Votes          []byte
	NextValidators []common.Address
	Receipts       []*chain.Receipt
}

type btpBlockRequest struct {
	height   int64
	hash     types.HexBytes
	indexes  [][]types.HexInt
	events   [][][]types.HexInt
	err      error
	retry    int
	response *btpBlockResponse
}

func NewReceiver(src, dst chain.BTPAddress,
	urls []string, rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error) {
	var recvOpts ReceiverOptions
	if err := json.Unmarshal(rawOpts, &recvOpts); err != nil {
		return nil, errors.Wrapf(err, "recvOpts.Unmarshal: %v", err)
	}

	dstAddr := dst.String()
	eventFilter := &types.EventFilter{
		Addr:      types.Address(src.ContractAddress()),
		Signature: EventSignature,
		Indexed:   []*string{&dstAddr},
	}
	evtReq := types.BlockRequest{
		EventFilters: []*types.EventFilter{eventFilter},
	} // fill height later

	efAddr, err := eventFilter.Addr.Value()
	if err != nil {
		return nil, errors.Wrapf(err, "eventFilter.Addr.Value: %v", err)
	}

	if recvOpts.SyncConcurrency < 1 {
		recvOpts.SyncConcurrency = 1
	} else if recvOpts.SyncConcurrency > MonitorBlockMaxConcurrency {
		recvOpts.SyncConcurrency = MonitorBlockMaxConcurrency
	}

	var client IClient
	if len(urls) > 0 {
		client = NewClient(urls[0], l)
	}

	recvr := &Receiver{
		log:      l,
		src:      src,
		dst:      dst,
		Client:   client,
		opts:     recvOpts,
		blockReq: evtReq,
		logFilter: eventLogRawFilter{
			addr:      efAddr,
			signature: []byte(EventSignature),
			next:      []byte(dstAddr),
		}, // fill seq later
	}

	return recvr, nil
}

func (r *Receiver) newVerifier(opts *types.VerifierOptions) (*Verifier, error) {
	validators, err := r.Client.GetValidatorsByHash(opts.ValidatorsHash)
	if err != nil {
		return nil, err
	}
	vr := Verifier{
		next:               int64(opts.BlockHeight),
		nextValidatorsHash: opts.ValidatorsHash,
		validators: map[string][]common.Address{
			opts.ValidatorsHash.String(): validators,
		},
	}
	header, err := r.Client.GetBlockHeaderByHeight(vr.next)
	if err != nil {
		return nil, err
	}
	votes, err := r.Client.GetVotesByHeight(
		&types.BlockHeightParam{Height: types.NewHexInt(vr.next)})
	if err != nil {
		return nil, err
	}

	ok, err := vr.Verify(header, votes)
	if !ok {
		err = errors.New("verification failed")
	}
	if err != nil {
		return nil, err
	}
	return &vr, nil
}

func (r *Receiver) syncVerifier(verifier IVerifier, height int64) error {
	if height == verifier.Next() {
		return nil
	}
	if verifier.Next() > height {
		return fmt.Errorf(
			"invalid target height: verifier height (%d) > target height (%d)",
			verifier.Next(), height)
	}

	r.log.WithFields(log.Fields{"height": verifier.Next(), "target": height}).Info("syncVerifier: start")

	for verifier.Next() < height {
		requestCh := make(chan *verifierBlockRequest, r.opts.SyncConcurrency)
		for i := verifier.Next(); len(requestCh) < cap(requestCh); i++ {
			requestCh <- &verifierBlockRequest{height: i, retry: 3}
		}

		responses := handleVerifierBlockRequests(requestCh, r.Client, verifier, r.log)

		// filter nil
		_sres, responses := responses, responses[:0]
		for _, resp := range _sres {
			if resp != nil && resp.err == nil {
				responses = append(responses, resp)
			}
		}

		// sort and forward notifications
		if len(responses) > 0 {
			sort.SliceStable(responses, func(i, j int) bool {
				return responses[i].Height < responses[j].Height
			})

			for _, response := range responses {
				if verifier.Next() == response.Height {
					ok, err := verifier.Verify(response.Header, response.Votes)
					if err != nil {
						return errors.Wrapf(err, "syncVerifier: Verify: height=%d, error=%v", response.Height, err)
					}
					if !ok {
						return fmt.Errorf("syncVerifier: invalid header: height=%d", response.Height)
					}

					err = verifier.Update(response.Header, response.NextValidators)
					if err != nil {
						return errors.Wrapf(err, "syncVerifier: Update: %v", err)
					}
				}
			}
			r.log.WithFields(log.Fields{"height": verifier.Next(), "target": height}).Debug("syncVerifier: syncing")
		}
	}

	r.log.WithFields(log.Fields{"height": verifier.Next()}).Info("syncVerifier: complete")
	return nil
}

func handleVerifierBlockRequests(requestCh chan *verifierBlockRequest, client IClient, verifier IVerifier, logger log.Logger) []*verifierBlockResponse {
	responseCh := make([]*verifierBlockResponse, 0, len(requestCh))

	for req := range requestCh {
		switch {
		case req.err != nil:
			if req.retry > 1 {
				req.retry--
				req.response, req.err = nil, nil
				requestCh <- req
				continue
			}
			logger.WithFields(log.Fields{"height": req.height, "error": req.err.Error()}).
				Debug("syncVerifier: request error")

			responseCh = append(responseCh, &verifierBlockResponse{err: req.err})
			if len(responseCh) == cap(responseCh) {
				close(requestCh)
			}

		case req.response != nil:
			responseCh = append(responseCh, req.response)
			if len(responseCh) == cap(responseCh) {
				close(requestCh)
			}

		default:
			go func(req *verifierBlockRequest) {
				defer func() {
					time.Sleep(500 * time.Millisecond)
					requestCh <- req
				}()
				if req.response == nil {
					req.response = &verifierBlockResponse{}
				}
				req.response.Height = req.height
				req.response.Header, req.err = client.GetBlockHeaderByHeight(req.height)
				if req.err != nil {
					req.err = errors.Wrapf(req.err, "syncVerifier: block height: %v, getBlockHeader: %v", req.height, req.err)
					return
				}
				req.response.Votes, req.err = client.GetVotesByHeight(
					&types.BlockHeightParam{Height: types.NewHexInt(req.height)})
				if req.err != nil {
					req.err = errors.Wrapf(req.err, "syncVerifier: block height: %v, GetVotesByHeight: %v", req.height, req.err)
					return
				}
				if len(verifier.Validators(req.response.Header.NextValidatorsHash)) == 0 {
					req.response.NextValidators, req.err = client.GetValidatorsByHash(req.response.Header.NextValidatorsHash)
					if req.err != nil {
						req.err = errors.Wrapf(req.err, "syncVerifier: block height: %v, GetValidatorsByHash: %v", req.height, req.err)
						return
					}
				}
			}(req)
		}
	}

	return responseCh
}

func (r *Receiver) receiveLoop(ctx context.Context, startHeight, startSeq uint64, callback func(rs []*chain.Receipt) error) (err error) {
	blockReq, logFilter := r.blockReq, r.logFilter // copy

	blockReq.Height, logFilter.seq = types.NewHexInt(int64(startHeight)), startSeq

	var vr IVerifier
	if r.opts.Verifier != nil {
		vr, err = r.newVerifier(r.opts.Verifier)
		if err != nil {
			return err
		}
	}

	errCh := make(chan error)                                                      // error channel
	reconnectCh := make(chan struct{}, 1)                                          // reconnect channel
	btpBlockNotifCh := make(chan *types.BlockNotification, r.opts.SyncConcurrency) // block notification channel
	btpBlockRespCh := make(chan *btpBlockResponse, cap(btpBlockNotifCh))           // block result channel

	reconnect := func() {
		select {
		case reconnectCh <- struct{}{}:
		default:
		}
		for len(btpBlockRespCh) > 0 || len(btpBlockNotifCh) > 0 {
			select {
			case <-btpBlockRespCh: // clear block result channel
			case <-btpBlockNotifCh: // clear block notification channel
			}
		}
	}

	next := int64(startHeight) // next block height to process

	// subscribe to monitor block
	ctxMonitorBlock, cancelMonitorBlock := context.WithCancel(ctx)
	reconnect()

loop:
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-errCh:
			return err

		// reconnect channel
		case <-reconnectCh:
			cancelMonitorBlock()
			ctxMonitorBlock, cancelMonitorBlock = context.WithCancel(ctx)

			// start new monitor loop
			go func(ctx context.Context, cancel context.CancelFunc) {
				defer cancel()
				blockReq.Height = types.NewHexInt(next)
				err := r.Client.MonitorBlock(ctx, &blockReq,
					func(conn *websocket.Conn, v *types.BlockNotification) error {
						if !errors.Is(ctx.Err(), context.Canceled) {
							btpBlockNotifCh <- v
						}
						return nil
					},
					func(conn *websocket.Conn) {},
					func(c *websocket.Conn, err error) {})
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					time.Sleep(time.Second * 5)
					reconnect()
					r.log.WithFields(log.Fields{"error": err}).Error("reconnect: monitor block error")
				}
			}(ctxMonitorBlock, cancelMonitorBlock)

			// sync verifier disabled
			if vr != nil {
				if err := r.syncVerifier(vr, next); err != nil {
					return errors.Wrapf(err, "sync verifier: %v", err)
				}
			}

		case blockResponse := <-btpBlockRespCh:

			err = handleBTPBlockResponse(blockResponse, vr, &next, reconnect, callback, btpBlockRespCh, r.log)
			if err != nil {
				return err
			}

		default:
			select {
			default:
			case bn := <-btpBlockNotifCh:

				requestCh := make(chan *btpBlockRequest, cap(btpBlockNotifCh))
				for i := int64(0); bn != nil; i++ {
					height, err := bn.Height.Value()
					if err != nil {
						panic(err)
					} else if height != next+i {
						r.log.WithFields(log.Fields{
							"height": log.Fields{"got": height, "expected": next + i},
						}).Error("reconnect: missing block notification")
						reconnect()
						continue loop
					}
					requestCh <- &btpBlockRequest{
						height:  height,
						hash:    bn.Hash,
						indexes: bn.Indexes,
						events:  bn.Events,
						retry:   RPCCallRetry,
					} // fill requestCh with requests
					if bn = nil; len(btpBlockNotifCh) > 0 && len(requestCh) < cap(requestCh) {
						bn = <-btpBlockNotifCh
					}
				}

				brs := make([]*btpBlockResponse, 0, len(requestCh))
				for request := range requestCh {
					switch {
					case request.err != nil:
						if request.retry > 0 {
							request.retry--
							request.response, request.err = nil, nil
							requestCh <- request
							continue
						}
						r.log.WithFields(log.Fields{"height": request.height, "error": request.err}).Debug("receiveLoop: request error")
						brs = append(brs, nil)
						if len(brs) == cap(brs) {
							close(requestCh)
						}

					case request.response != nil:
						brs = append(brs, request.response)
						if len(brs) == cap(brs) {
							close(requestCh)
						}

					default:
						go handleBTPBlockRequest(request, requestCh, vr, r.Client, logFilter, r.log)

					}
				}
				// filter nil
				_brs, brs := brs, brs[:0]
				for _, v := range _brs {
					if v != nil {
						brs = append(brs, v)
					}
				}
				// sort and forward notifications
				if len(brs) > 0 {
					sort.SliceStable(brs, func(i, j int) bool {
						return brs[i].Height < brs[j].Height
					})
					for i, d := range brs {
						if d.Height == next+int64(i) {
							btpBlockRespCh <- d
						}
					}
				}
			}
		}
	}

}

func handleBTPBlockRequest(
	request *btpBlockRequest, requestCh chan *btpBlockRequest, vr IVerifier, client IClient, logFilter eventLogRawFilter, logger log.Logger) {
	defer func() {
		time.Sleep(500 * time.Millisecond)
		requestCh <- request
	}()

	if request.response == nil {
		request.response = &btpBlockResponse{}
	}
	request.response.Height = request.height
	request.response.Hash, request.err = request.hash.Value()
	if request.err != nil {
		request.err = errors.Wrapf(request.err,
			"invalid hash: height=%v, hash=%v, %v", request.height, request.hash, request.err)
		return
	}

	request.response.Header, request.err = client.GetBlockHeaderByHeight(request.height)
	if request.err != nil {
		request.err = errors.Wrapf(request.err, "getBlockHeader: %v", request.err)
		return
	}
	// fetch votes, next validators only if verifier exists
	if vr != nil {
		request.response.Votes, request.err = client.GetVotesByHeight(
			&types.BlockHeightParam{Height: types.NewHexInt(request.height)})
		if request.err != nil {
			request.err = errors.Wrapf(request.err, "GetVotesByHeight: %v", request.err)
			return
		}
		if len(vr.Validators(request.response.Header.NextValidatorsHash)) == 0 {
			request.response.NextValidators, request.err = client.GetValidatorsByHash(request.response.Header.NextValidatorsHash)
			if request.err != nil {
				request.err = errors.Wrapf(request.err, "GetValidatorsByHash: %v", request.err)
				return
			}
		}
	}

	if len(request.indexes) > 0 && len(request.events) > 0 {
		var hr BlockHeaderResult
		_, err := codec.RLP.UnmarshalFromBytes(request.response.Header.Result, &hr)
		if err != nil {
			request.err = errors.Wrapf(err, "BlockHeaderResult.UnmarshalFromBytes: %v", request.err)
			return
		}
		for i, index := range request.indexes[0] {
			p := &types.ProofEventsParam{
				Index:     index,
				BlockHash: request.hash,
				Events:    request.events[0][i],
			}
			proofs, err := client.GetProofForEvents(p)
			if err != nil {
				request.err = errors.Wrapf(err, "GetProofForEvents: %v", err)
				return
			}

			if len(proofs) != 1+len(p.Events) { // num_receipt + num_events
				var err error
				if request.err != nil {
					err = request.err
				} else {
					err = errors.New("Proof does not include all events")
				}

				request.err = errors.Wrapf(err,
					"Proof does not include all events: len(proofs)=%d, expected=%d",
					len(proofs), len(p.Events)+1)
				return
			}

			// Processing receipt index
			serializedReceipt, err := mptProve(index, proofs[0], hr.ReceiptHash)
			if err != nil {
				request.err = errors.Wrapf(err, "MPTProve Receipt: %v", err)
				return
			}
			var result TxResult
			_, err = codec.RLP.UnmarshalFromBytes(serializedReceipt, &result)
			if err != nil {
				request.err = errors.Wrapf(err, "Unmarshal Receipt: %v", err)
				return
			}

			idx, _ := index.Value()
			receipt := &chain.Receipt{
				Index:  uint64(idx),
				Height: uint64(request.height),
			}
			for j := 0; j < len(p.Events); j++ {
				// nextEP is pointer to event where sequence has caught up
				serializedEventLog, err := mptProve(
					p.Events[j], proofs[j+1], common.HexBytes(result.EventLogsHash))
				if err != nil {
					request.err = errors.Wrapf(err, "event.MPTProve: %v", err)
					return
				}
				var el types.EventLog
				_, err = codec.RLP.UnmarshalFromBytes(serializedEventLog, &el)
				if err != nil {
					request.err = errors.Wrapf(err, "event.UnmarshalFromBytes: %v", err)
					return
				}

				if bytes.Equal(el.Addr, logFilter.addr) &&
					bytes.Equal(el.Indexed[EventIndexSignature], logFilter.signature) &&
					bytes.Equal(el.Indexed[EventIndexNext], logFilter.next) {
					var seqGot common.HexInt
					seqGot.SetBytes(el.Indexed[EventIndexSequence])
					evt := &chain.Event{
						Next:     chain.BTPAddress(el.Indexed[EventIndexNext]),
						Sequence: seqGot.Uint64(),
						Message:  el.Data[0],
					}
					receipt.Events = append(receipt.Events, evt)
				} else {
					if !bytes.Equal(el.Addr, logFilter.addr) {
						logger.WithFields(log.Fields{
							"height":   request.height,
							"got":      common.HexBytes(el.Addr),
							"expected": common.HexBytes(logFilter.addr)}).Error("invalid event: cannot match addr")
					}
					if !bytes.Equal(el.Indexed[EventIndexSignature], logFilter.signature) {
						logger.WithFields(log.Fields{
							"height":   request.height,
							"got":      common.HexBytes(el.Indexed[EventIndexSignature]),
							"expected": common.HexBytes(logFilter.signature)}).Error("invalid event: cannot match sig")
					}
					if !bytes.Equal(el.Indexed[EventIndexNext], logFilter.next) {
						logger.WithFields(log.Fields{
							"height":   request.height,
							"got":      common.HexBytes(el.Indexed[EventIndexNext]),
							"expected": common.HexBytes(logFilter.next)}).Error("invalid event: cannot match next")
					}
					request.err = errors.New("invalid event")
					return
				}
			}
			if len(receipt.Events) > 0 {
				if len(receipt.Events) == len(p.Events) {
					request.response.Receipts = append(request.response.Receipts, receipt)
				} else {
					logger.WithFields(log.Fields{
						"height":              request.height,
						"receipt_index":       index,
						"got_num_events":      len(receipt.Events),
						"expected_num_events": len(p.Events)}).Error("failed to verify all events for the receipt")
					request.err = errors.New("failed to verify all events for the receipt")
					return
				}
			}
		}
	}
}

func handleBTPBlockResponse(blockResponse *btpBlockResponse, vr IVerifier, next *int64,
	reconnect func(), callback func(rs []*chain.Receipt) error,
	blockResponseCh chan *btpBlockResponse, logger log.Logger) error {

	for ; blockResponse != nil; *next++ {
		log.WithFields(log.Fields{"height": blockResponse.Height}).Debug("block notification")

		if vr != nil {
			ok, err := vr.Verify(blockResponse.Header, blockResponse.Votes)
			if !ok || err != nil {
				if err != nil {
					logger.WithFields(log.Fields{"height": blockResponse.Height, "error": err}).Error("receiveLoop: verification error")
				} else if !ok {
					logger.WithFields(log.Fields{"height": blockResponse.Height, "hash": blockResponse.Hash}).Error("receiveLoop: invalid header")
				}

				reconnect() // reconnect websocket
				logger.WithFields(log.Fields{"height": blockResponse.Height, "hash": blockResponse.Hash}).Error("reconnect: verification failed")
				break
			}
			if err := vr.Update(blockResponse.Header, blockResponse.NextValidators); err != nil {
				return errors.Wrapf(err, "receiveLoop: update verifier: %v", err)
			}
		}
		if err := callback(blockResponse.Receipts); err != nil {
			return errors.Wrapf(err, "receiveLoop: callback: %v", err)
		}
		if blockResponse = nil; len(blockResponseCh) > 0 {
			blockResponse = <-blockResponseCh
		}
	}

	// remove unprocessed block responses
	for len(blockResponseCh) > 0 {
		<-blockResponseCh
	}

	return nil
}

func (r *Receiver) Subscribe(
	ctx context.Context, msgCh chan<- *chain.Message,
	opts chain.SubscribeOptions) (errCh <-chan error, err error) {

	opts.Seq++

	if opts.Height < 1 {
		opts.Height = 1
	}

	_errCh := make(chan error)
	go func() {
		defer close(_errCh)
		err := r.receiveLoop(ctx, opts.Height, opts.Seq, func(receipts []*chain.Receipt) error {
			for _, receipt := range receipts {
				events := receipt.Events[:0]
				for _, event := range receipt.Events {
					switch {
					case event.Sequence == opts.Seq:
						events = append(events, event)
						opts.Seq++
					case event.Sequence > opts.Seq:
						r.log.WithFields(log.Fields{
							"seq": log.Fields{"got": event.Sequence, "expected": opts.Seq},
						}).Error("invalid event seq")
						return fmt.Errorf("invalid event seq")
					}
				}
				receipt.Events = events
			}
			if len(receipts) > 0 {
				msgCh <- &chain.Message{Receipts: receipts}
			}
			return nil
		})
		if err != nil {
			r.log.Errorf("receiveLoop terminated: %v", err)
			_errCh <- err
		}
	}()
	return _errCh, nil
}
