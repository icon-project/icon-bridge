package algo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

// TODO adjust settings for algo
const (
	MonitorBlockMaxConcurrency = 300 // number of concurrent requests to synchronize older blocks
)

type receiver struct {
	log  log.Logger
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts ReceiverOptions
	cl   *Client
	vr   Verifier
}

type VerifierOptions struct {
	Round     uint64   `json:"Round"`
	BlockHash [32]byte `json:"BlockHash"`
}

type Verifier struct {
	Round     uint64
	BlockHash [32]byte
}

func NewReceiver(
	src, dst chain.BTPAddress,
	algodAccess []string,
	rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error) {
	r := &receiver{
		log: l,
		src: src,
		dst: dst,
	}

	if len(algodAccess) < 2 {
		return nil, fmt.Errorf("Invalid algorand credentials")
	}

	err := json.Unmarshal(rawOpts, &r.opts)
	if err != nil {
		return nil, err
	}
	if r.opts.SyncConcurrency < 1 {
		r.opts.SyncConcurrency = 1
	} else if r.opts.SyncConcurrency > MonitorBlockMaxConcurrency {
		r.opts.SyncConcurrency = MonitorBlockMaxConcurrency
	}

	r.vr = Verifier{
		Round:     r.opts.Verifier.Round,
		BlockHash: r.opts.Verifier.BlockHash,
	}
	r.cl, err = newClient(algodAccess, r.log)
	if err != nil {
		return nil, err
	}
	return r, nil

}

type ReceiverOptions struct {
	SyncConcurrency uint64           `json:"syncConcurrency"`
	Verifier        *VerifierOptions `json:"verifier"`
}

func (opts *ReceiverOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

func (r *receiver) Subscribe(
	ctx context.Context, msgCh chan<- *chain.Message,
	subOpts chain.SubscribeOptions) (errCh <-chan error, err error) {

	subOpts.Seq++
	_errCh := make(chan error)

	if subOpts.Seq <= 0 || subOpts.Height <= 0 {
		return _errCh, errors.New("receiveLoop: invalid options: <nil>")
	}

	latestRound, err := r.cl.GetLatestRound(ctx)
	if err != nil {
		r.log.WithFields(log.Fields{"error": err}).Error(
			"receiveLoop: error failed to getLatestRound-")
		return _errCh, err
	}

	if err != nil {
		r.log.WithFields(log.Fields{"error": err}).Error("receiveLoop: failed to GetLatestRound")
		return _errCh, err
	}

	go func() {
		defer func() {
			_errCh <- errors.New("aborting receiveloop")
			close(_errCh)
		}()
	receiveLoop:
		for {
			select {
			case <-ctx.Done():
				break receiveLoop
			default:
				if r.vr.Round >= latestRound {
					time.Sleep(500 * time.Millisecond)

					latestRound, err = r.cl.GetLatestRound(ctx)
					if err != nil {
						r.log.WithFields(log.Fields{"error": err}).Error(
							"receiveLoop: error failed to getLatestRound")
						_errCh <- err
					}
					continue
				}
				//Check the latest block for txns addressed to this BMC
				r.inspectBlock(ctx, r.vr.Round, &subOpts, msgCh, _errCh)

			}
		}
	}()
	return _errCh, err
}

// Inspects the latest block created for new relay messages
func (r *receiver) inspectBlock(ctx context.Context, round uint64, subOpts *chain.SubscribeOptions,
	msgCh chan<- *chain.Message, _errCh chan error) {
	newBlock, err := r.cl.GetBlockbyRound(ctx, round)
	if err != nil {
		_errCh <- err
		return
	}

	if bytes.Equal(newBlock.BlockHeader.Branch[:], r.vr.BlockHash[:]) {
		r.vr.BlockHash = EncodeBlockHash(newBlock)
		r.vr.Round++
	} else {
		_errCh <- fmt.Errorf("Block at round %d does not have a valid parent hash.", round)
		return
	}

	// Don't start inspecting blocks until the subscribed round
	if round <= subOpts.Height {
		return
	}

	bmcTxns := r.getBMCTxns(newBlock)
	if len(*bmcTxns) <= 0 {
		return
	}

	relayRcps, err := r.getRelayReceipts(bmcTxns, round)
	if err != nil {
		_errCh <- err
		return
	}

	err = r.validateEvents(&relayRcps, subOpts)
	if err != nil {
		_errCh <- err
		return
	}
	msgCh <- &chain.Message{Receipts: relayRcps}
}

// Check if the new block has any transaction meant to be sent across the relayer
func (r *receiver) getBMCTxns(block *types.Block) *[]types.SignedTxnWithAD {
	txns := make([]types.SignedTxnWithAD, 0)
	for _, signedTxnInBlock := range block.Payset {
		signedTxnWithAD := signedTxnInBlock.SignedTxnWithAD

		if signedTxnWithAD.SignedTxn.AuthAddr.String() == r.src.ContractAddress() {
			txns = append(txns, signedTxnWithAD)
		}
	}
	return &txns
}

func (r *receiver) getRelayReceipts(txns *[]types.SignedTxnWithAD, round uint64) (
	[]*chain.Receipt, error) {
	var receipts []*chain.Receipt
	var events []*chain.Event
	for i, txn := range *txns {
		events := events[:0]
		for _, log := range txn.ApplyData.EvalDelta.Logs {
			if txn.Txn.Header.Sender.String() != r.src.ContractAddress() {
				continue
			}
			decodedMsg, err := r.cl.DecodeBtpMsg(log)
			if err == nil {
				events = append(events, decodedMsg)
			}
		}
		if len(events) > 0 {
			rcp := &chain.Receipt{}
			rcp.Index, rcp.Height = uint64(i), round
			rcp.Events = append(rcp.Events, events...)
			receipts = append(receipts, rcp)
		}
	}
	if len(receipts) <= 0 {
		return receipts, errors.New("Couldn't retrieve any receipt from the new block")
	}
	return receipts, nil
}

func (r *receiver) validateEvents(rcps *[]*chain.Receipt, subOpts *chain.SubscribeOptions) error {
	for _, receipt := range *rcps {
		events := receipt.Events[:0]
		for _, event := range receipt.Events {
			switch {
			case event.Sequence == subOpts.Seq:
				events = append(events, event)
				subOpts.Seq++
			case event.Sequence > subOpts.Seq:
				r.log.WithFields(log.Fields{
					"seq": log.Fields{"got": event.Sequence, "expected": subOpts.Seq},
				}).Error("invalid event seq")
				return fmt.Errorf("invalid event seq")
			}
		}
		receipt.Events = events
	}
	return nil
}