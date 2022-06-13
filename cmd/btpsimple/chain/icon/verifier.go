package icon

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"

	"github.com/icon-project/btp/cmd/btpsimple/chain"
	"github.com/icon-project/btp/common/codec"
	"github.com/icon-project/btp/common/crypto"
	"github.com/icon-project/btp/common/log"
	"github.com/icon-project/goloop/common"
	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/common/trie/ompt"
)

type VerifierOptions struct {
	BlockHeight   uint64 `json:"blockHeight"`
	ValidatorHash string `json:"validatorHash"`
}

type commitVoteItem struct {
	Timestamp int64
	Signature common.Signature
}

type commitVoteList struct {
	Round          int32
	BlockPartSetID *PartSetID
	Items          []commitVoteItem
}

type PartSetID struct {
	Count uint16
	Hash  []byte
}

type _HR struct {
	Height int64
	Round  int32
}

type voteBase struct {
	_HR
	Type           VoteType
	BlockID        []byte
	BlockPartSetID *PartSetID
}

type vote struct {
	voteBase
	Timestamp int64
}

type VoteType byte

const (
	VoteTypePrevote VoteType = iota
	VoteTypePrecommit
	numberOfVoteTypes
)

type Result struct {
	StateHash        []byte
	PatchReceiptHash []byte
	ReceiptHash      common.HexBytes
	ExtensionData    []byte
}

type Receipt struct {
	Status             int64
	To                 []byte
	CumulativeStepUsed []byte
	StepUsed           []byte
	StepPrice          []byte
	LogsBloom          []byte
	EventLogs          []EventLog
	ScoreAddress       []byte
	EventLogsHash      common.HexBytes
}

func (r *receiver) blockVerification(v *BlockNotification) ([]*chain.Receipt, error) {
	header, err := r.verifyHeader(v)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, errors.New("Err: verifyHeader returned nil header")
	}
	return r.verifyReceipt(header, v)
}

func (r *receiver) verifyHeader(v *BlockNotification) (*BlockHeader, error) {
	r.log.WithFields(log.Fields{"Height": v.Height, "Hash": v.Hash}).Debug("verifyHeader Start")
	header, err := r.getBlockHeader(v.Height)
	if err != nil {
		return nil, errors.Wrap(err, "verifyHeader; getBlockHeader; Err: ")
	}
	// Hash
	blockHash, err := v.Hash.Value()
	if err != nil {
		return nil, errors.Wrap(err, "verifyHeader; GetHashValue; Err: ")
	}
	if !bytes.Equal(blockHash, crypto.SHA3Sum256(header.serialized)) {
		return nil, fmt.Errorf("verifyHeader; mismatch block hash with BlockNotification")
	}

	// Votes
	votesBytes, err := r.cl.GetVotesByHeight(&BlockHeightParam{Height: v.Height})
	if err != nil {
		return nil, errors.Wrap(mapError(err), "verifyHeader; GetVotesByHeight; Err: ")
	}
	cvl := &commitVoteList{}
	_, err = vlcodec.BC.UnmarshalFromBytes(votesBytes, cvl)
	if err != nil {
		return nil, errors.Wrap(mapError(err), "verifyHeader; UnmarshalFromBytes of Votes; Err: ")
	}

	vote := &vote{
		voteBase: voteBase{
			_HR: _HR{
				Height: header.Height,
				Round:  cvl.Round,
			},
			Type:           VoteTypePrecommit,
			BlockID:        blockHash,
			BlockPartSetID: cvl.BlockPartSetID,
		},
	}

	validCount := 0
	for _, item := range cvl.Items {
		vote.Timestamp = item.Timestamp
		pub, err := item.Signature.RecoverPublicKey(crypto.SHA3Sum256(vlcodec.BC.MustMarshalToBytes(vote)))
		if err != nil {
			r.log.Error(errors.Wrap(mapError(err), "verifyHeader; UnmarshalFromBytes of Validators; Err: "))
			err = nil
			continue
		}
		address := common.NewAccountAddressFromPublicKey(pub)

		if addressesContains(address.Bytes(), r.validators) {
			validCount++
		}
	}
	if validCount < int(2*len(r.validators)/3) {
		return nil, errors.New("verifyHeader; Block not validated by >= 2/3 of validators")
	}

	r.log.WithFields(log.Fields{"NumValidators": len(r.validators), "NumVotes": validCount}).Debug("Verified Votes")

	// Update validators
	if !bytes.Equal(header.NextValidatorsHash, r.validatorHash) { // should update
		r.log.WithFields(log.Fields{"Height": v.Height, "NewHash": base64.StdEncoding.EncodeToString(header.NextValidatorsHash), "OldHash": base64.StdEncoding.EncodeToString(r.validatorHash)}).Info("Updating Validators ")
		vBytes, err := r.cl.GetDataByHash(&DataHashParam{Hash: NewHexBytes(header.NextValidatorsHash)})
		if err != nil {
			return nil, errors.Wrap(err, "verifyHeader; GetDataByHash Validators; Err: ")
		}
		var vs [][]byte
		_, err = vlcodec.BC.UnmarshalFromBytes(vBytes, &vs)
		if err != nil {
			return nil, errors.Wrap(err, "verifyHeader; Unmarshal Validators; Err: ")
		}
		r.validatorHash = header.NextValidatorsHash
		r.validators = vs
	}
	return header, nil
}

func (r *receiver) verifyReceipt(header *BlockHeader, v *BlockNotification) ([]*chain.Receipt, error) {
	// Update blockHeaders
	if len(v.Indexes) == 0 || len(v.Events) == 0 {
		r.log.Debug("ReceiptVerification; Events and Indexes are empty; Skipping their verification", v.Height)
		return nil, nil
	}

	var headerResult Result
	vlcodec.RLP.MustUnmarshalFromBytes(header.Result, &headerResult)

	r.log.WithFields(log.Fields{"Height": v.Height, "ReceiptHash": headerResult.ReceiptHash, "NumIndexes": len(v.Indexes[0])}).Debug("Receipt Verification Start")
	rps := make([]*chain.Receipt, 0)

	for i, index := range v.Indexes[0] {
		p := &ProofEventsParam{BlockHash: v.Hash, Index: index, Events: v.Events[0][i]}
		proofs, err := r.cl.GetProofForEvents(p)
		if err != nil {
			return nil, errors.Wrap(mapError(err), "ReceiptVerification; GetProofForEvents; Err:")
		}
		if len(proofs) != 1+len(p.Events) { // returned proofs should be for all of the requested Events plus 1 for the receipt
			return nil, errors.New("ReceiptVerification; Proof Not returned for all requested events")
		}

		// Processing receipt index
		serializedReceipt, err := mptProve(index, proofs[0], headerResult.ReceiptHash)
		if err != nil {
			return nil, errors.Wrap(err, "ReceiptVerification; MPTProve Receipt; Err:")
		}
		var receipt Receipt
		_, err = vlcodec.RLP.UnmarshalFromBytes(serializedReceipt, &receipt)
		if err != nil {
			return nil, errors.Wrap(err, "ReceiptVerification; Unmarshal Receipt; Err:")
		}

		idx, _ := index.Value()
		rp := &chain.Receipt{
			Index:  uint64(idx),
			Height: hexInt2Uint64(v.Height),
		}
		for j := 0; j < len(p.Events); j++ { // nextEP is pointer to event where sequence has caught up
			serializedEventLog, err := mptProve(p.Events[j], proofs[j+1], common.HexBytes(receipt.EventLogsHash))
			if err != nil {
				return nil, errors.Wrap(err, "ReceiptVerification MPTPrice Events; Err:")
			}
			var el EventLog
			_, err = codec.RLP.UnmarshalFromBytes(serializedEventLog, &el)
			if err != nil {
				return nil, errors.Wrap(err, "ReceiptVerification; Unmarshal Events; Err:")
			}

			if bytes.Equal(el.Addr, r.evtLogRawFilter.addr) && bytes.Equal(el.Indexed[EventIndexSignature], r.evtLogRawFilter.signature) && bytes.Equal(el.Indexed[EventIndexNext], r.evtLogRawFilter.next) {
				var seqGot common.HexInt
				var seqExpected common.HexInt
				seqGot.SetBytes(el.Indexed[EventIndexSequence])
				seqExpected.SetBytes(r.evtLogRawFilter.seq)
				if !r.isFoundOffsetBySeq && seqGot.Uint64() < seqExpected.Uint64() { // If sequence has not been found and this is not the one; continue searching
					r.log.WithFields(log.Fields{"Height": v.Height, "got": common.HexBytes(el.Indexed[EventIndexSequence]), "expected": common.HexBytes(r.evtLogRawFilter.seq)}).Info("Searching for matching sequence...")
					continue
				}
				r.isFoundOffsetBySeq = true
				evt := &chain.Event{
					Next:     chain.BTPAddress(el.Indexed[EventIndexNext]),
					Sequence: seqGot.Uint64(),
					Message:  el.Data[0],
				}
				rp.Events = append(rp.Events, evt)
			} else {
				if !bytes.Equal(el.Addr, r.evtLogRawFilter.addr) {
					r.log.WithFields(log.Fields{"Height": v.Height, "got": common.HexBytes(el.Addr), "expected": common.HexBytes(r.evtLogRawFilter.addr)}).Error("invalid event: cannot match add")
				}
				if !bytes.Equal(el.Indexed[EventIndexSignature], r.evtLogRawFilter.signature) {
					r.log.WithFields(log.Fields{"Height": v.Height, "got": common.HexBytes(el.Indexed[EventIndexSignature]), "expected": common.HexBytes(r.evtLogRawFilter.signature)}).Error("invalid event: cannot match sig")
				}
				if !bytes.Equal(el.Indexed[EventIndexNext], r.evtLogRawFilter.next) {
					r.log.WithFields(log.Fields{"Height": v.Height, "got": common.HexBytes(el.Indexed[EventIndexNext]), "expected": common.HexBytes(r.evtLogRawFilter.next)}).Error("invalid event: cannot match next")
				}
			}
		}
		if len(rp.Events) > 0 && len(rp.Events) == len(p.Events) { //Only add if all the events were verified
			rps = append(rps, rp)
		} else if len(rp.Events) > 0 && len(rp.Events) != len(p.Events) {
			r.log.WithFields(log.Fields{"Height": v.Height, "ReceiptIndex": index, "got": len(rp.Events), "expected": len(p.Events)}).Info(" Not all events were verified for receipt ")
		}
	}
	return rps, nil
}

func mptProve(key HexInt, proofs [][]byte, hash []byte) ([]byte, error) {
	db := db.NewMapDB()
	defer db.Close()
	index, err := key.Value()
	if err != nil {
		return nil, err
	}
	indexKey, err := vlcodec.RLP.MarshalToBytes(index)
	if err != nil {
		return nil, err
	}
	mpt := ompt.NewMPTForBytes(db, hash)
	trie, err1 := mpt.Prove(indexKey, proofs)
	if err1 != nil {
		return nil, err1

	}
	return trie, nil
}

func addressesContains(data []byte, list [][]byte) bool {
	for _, current := range list {
		if bytes.Equal(data, current) {
			return true
		}
	}
	return false
}
