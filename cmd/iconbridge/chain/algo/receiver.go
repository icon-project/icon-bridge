package algo

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/algorand/go-algorand-sdk/abi"
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

	blockHash, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(r.opts.Verifier.BlockHash)
	if err != nil {
		return nil, err
	}

	var arr [32]byte
	copy(arr[:], blockHash)

	r.vr = Verifier{
		Round:     r.opts.Verifier.Round,
		BlockHash: arr,
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
				//Wait for new blocks to be created
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
		// identify transactions sent from the algorand BMC
		if signedTxnInBlock.SignedTxnWithAD.SignedTxn.Txn.Header.Sender.String() == testAddress &&
			signedTxnInBlock.EvalDelta.Logs != nil {
			// there could be multiple logs sent from each transaction
			for _, txnLog := range signedTxnInBlock.EvalDelta.Logs {
				bmcMsg, err := extractMsg(txnLog)
				if err != nil {
					return nil, fmt.Errorf("Error extracting message from log: %s", err)
				}
				if chain.BTPAddress(bmcMsg.Dst) != r.dst {
					return nil, fmt.Errorf("Unexpected destination %s, expected %s. - Block %d",
						bmcMsg.Dst, r.dst, block.Round)
				}
				var sn uint64
				binary.Read(bytes.NewReader(bmcMsg.Sn), binary.BigEndian, &sn)
				events = append(events, &chain.Event{
					Next:     chain.BTPAddress(bmcMsg.Dst),
					Sequence: sn,
					Message:  bmcMsg.Message,
				})
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

func extractMsg(txnLog string) (BMCMessage, error) {
	var bmcMessage BMCMessage

	tupleType, err := abi.TypeOf("(string,string,string,uint64,byte[])")

	if err != nil {
		return bmcMessage, fmt.Errorf("Failed to get tuple type: %+v", err)
	}
	decoded, err := tupleType.Decode([]byte(txnLog))
	if err != nil {
		return bmcMessage, fmt.Errorf("Failed to decode tuple type: %+v", err)
	}

	if val, ok := decoded.([]interface{}); ok {
		var msgBytes []byte
		for _, v := range val[4].([]interface{}) {
			msgBytes = append(msgBytes, byte(v.(uint8)))
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, val[3].(uint64))

		bmcMessage.Src = val[0].(string)
		bmcMessage.Dst = val[1].(string)
		bmcMessage.Svc = val[2].(string)
		bmcMessage.Sn = buf.Bytes()
		bmcMessage.Message = msgBytes
	} else {
		return bmcMessage, fmt.Errorf("Decoded tuple had unexpected type %+v", err)
	}
	return bmcMessage, nil
}
