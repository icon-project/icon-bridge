package algo

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
)

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

type ReceiverOptions struct {
	SyncConcurrency uint64           `json:"syncConcurrency"`
	AppID           uint64           `json:"appID"`
	Verifier        *VerifierOptions `json:"verifier"`
}
type VerifierOptions struct {
	Round     uint64 `json:"round"`
	BlockHash string `json:"blockHash"`
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

	r.cl, err = newClient(algodAccess, r.log)
	if err != nil {
		return nil, err
	}
	hashStr, err := r.cl.GetBlockHash(context.Background(), r.opts.Verifier.Round-1)
	if err != nil {
		return nil, err
	}

	blockHash, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(hashStr)
	if err != nil {
		return nil, err
	}

	var hashBytes [32]byte
	copy(hashBytes[:], blockHash)

	r.vr = Verifier{
		Round:     r.opts.Verifier.Round,
		BlockHash: hashBytes,
	}
	return r, nil
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

	if subOpts.Seq < 0 || subOpts.Height < 0 {
		return _errCh, errors.New("receiveLoop: invalid options: <nil>")
	}

	latestRound, err := r.cl.GetLatestRound(ctx)

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

	receipts, err := r.getRelayReceipts(newBlock, &subOpts.Seq)
	if err != nil {
		_errCh <- fmt.Errorf("Error getting relay receipts: %s", err)
	} else if len(receipts) > 0 {
		msgCh <- &chain.Message{Receipts: receipts}
	}
}

// Check if the new block has any transaction meant to be sent across the relayer
// If so read its logs to produce an event to be forwarded
func (r *receiver) getRelayReceipts(block *types.Block, seq *uint64) (
	[]*chain.Receipt, error) {
	receipts := make([]*chain.Receipt, 0)
	events := make([]*chain.Event, 0)
	var index uint64

	for _, signedTxnInBlock := range block.Payset {
		// identify transactions sent from the algorand BMC containign logged messages
		for _, innerTxn := range signedTxnInBlock.EvalDelta.InnerTxns {
			if innerTxn.SignedTxn.Txn.ApplicationFields.ApplicationCallTxnFields.ApplicationID == types.AppIndex(r.opts.AppID) &&
				len(innerTxn.ApplyData.EvalDelta.Logs) > 0 {
				r.log.Debug("New message from algorand BMC")
				err := incrementSeq("tx_seq")
				if err != nil {
					r.log.WithFields(log.Fields{"error": err}).Error(
						"getRelayReceipts: error incrementing tx_seq")
					return nil, err
				}
				err = updateHeight("tx_height", uint64(block.Round))
				if err != nil {
					r.log.WithFields(log.Fields{"error": err}).Error(
						"getRelayReceipts: error updating tx_height")
					return nil, err
				}
				args := innerTxn.SignedTxn.Txn.ApplicationFields.ApplicationArgs
				event, err := r.getEventFromMsg(innerTxn.ApplyData.EvalDelta.Logs[0], args)
				if err != nil {
					r.log.WithFields(log.Fields{"error": err}).Error(
						"getRelayReceipts: error extracting event from relay message")
					return nil, err
				}
				events = append(events, event)
			}
			// sort txn events in case they came out of order
			sort.Slice(events, func(i, j int) bool {
				return events[i].Sequence < events[j].Sequence
			})
			// check if event sequece of each event increments starting at the current sequence
			for _, event := range events {
				if event.Sequence != *seq {
					return nil, fmt.Errorf("Unexpected sequece, got %d, expected %d. - Block %d",
						event.Sequence, *seq, block.Round)
				}
				*seq++
			}
			receipts = append(receipts, &chain.Receipt{
				Index:  index,
				Events: events,
				Height: uint64(block.Round),
			})
			events = nil
			index++
		}
	}
	return receipts, nil
}

func (r *receiver) getEventFromMsg(txnLog string, appArgs [][]uint8) (*chain.Event, error) {
	btpIndex := strings.Index(string(appArgs[1]), "btp")
	var dst string
	if btpIndex != -1 {
		dst = string(appArgs[1])[btpIndex:]
	}
	link, err := getStatus()
	if err != nil {
		return &chain.Event{}, fmt.Errorf("Failed to get status: %v", err)
	}

	bmcMsg := BMCMessageAlgo{
		Src:     r.src.String(),
		Dst:     strings.TrimRight(dst, "\n"),
		Svc:     txnLog,
		Sn:      link.TxSeq,
		Message: appArgs[3],
	}

	if chain.BTPAddress(bmcMsg.Dst) != r.dst {
		return &chain.Event{}, fmt.Errorf("Unexpected msg destination %s, expected %s.", bmcMsg.Dst, r.dst)
	}

	rlpMsg, err := RlpEncodeHex(bmcMsg)

	if err != nil {
		return &chain.Event{}, fmt.Errorf("Failed to rlp encode BMC message: %v", err)
	}

	event := &chain.Event{
		Next:     chain.BTPAddress(bmcMsg.Dst),
		Sequence: bmcMsg.Sn,
		Message:  rlpMsg,
	}
	return event, nil
}
