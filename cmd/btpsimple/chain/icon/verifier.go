package icon

import (
	"bytes"
	"encoding/base64"

	"github.com/pkg/errors"

	"github.com/icon-project/goloop/common"
	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/icon-project/icon-bridge/common/log"
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

type headerValidator struct {
	validators    [][]byte
	validatorHash []byte
	height        uint64
}

func (r *receiver) syncVerifier(hexHeight HexInt) error {
	ht, hterr := hexHeight.Value()
	if hterr != nil {
		return errors.Wrapf(hterr, "syncVerifier; HexInt Conversion Error at Height %v ", hexHeight)
	}
	targetHeight := uint64(ht)

	if targetHeight < r.hv.height {
		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorHeight": NewHexInt(int64(r.hv.height))}).Error("SyncVerifier; TargetHeight is less than known validator height")
		return errors.New("SyncVerifier; TargetHeight is less than height for which we know the validator hash ")
	} else if targetHeight == r.hv.height {
		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorHeight": NewHexInt(int64(r.hv.height))}).Error("SyncVerifier; Same Height so already in sync ")
		return nil
	}
	r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorHeight": NewHexInt(int64(r.hv.height)), "ValidatorHash": base64.StdEncoding.EncodeToString(r.hv.validatorHash)}).Info("Sync Verifier; Start Sync ")
	for ht := r.hv.height; ht < targetHeight; ht++ {
		if header, err := r.getVerifiedHeaderForHeight(ht); err != nil {
			return errors.Wrap(err, "syncVerifier; ")
		} else {
			if !bytes.Equal(header.NextValidatorsHash, r.hv.validatorHash) { // should update
				r.log.WithFields(log.Fields{"ValidatorHeight": NewHexInt(int64(ht)), "NewValidatorHash": base64.StdEncoding.EncodeToString(header.NextValidatorsHash), "OldValidatorHash": base64.StdEncoding.EncodeToString(r.hv.validatorHash)}).Info("Sync Verifier; Updating Validator Hash ")
				if vs, err := getValidatorsFromHash(r.cl, header.NextValidatorsHash); err != nil {
					return errors.Wrap(err, "syncVerifier; ")
				} else {
					r.hv.validatorHash = header.NextValidatorsHash
					r.hv.validators = vs
				}
			}
			r.hv.height = ht + 1
		}
		if ht%50 == 0 {
			r.log.WithFields(log.Fields{"ValidatorHeight": NewHexInt(int64(ht)), "TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorHash": base64.StdEncoding.EncodeToString(r.hv.validatorHash)}).Info("Sync Verifier; In Progress ")
		}
	}
	r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorHash": base64.StdEncoding.EncodeToString(r.hv.validatorHash)}).Info("Sync Verifier; Complete ")
	return nil
}

func (r *receiver) verify(v *BlockNotification) (*BlockHeader, []*chain.Receipt, error) {
	r.log.WithFields(log.Fields{"Height": v.Height, "Hash": v.Hash}).Debug("verify Start")
	ht, err := v.Height.Value()
	if err != nil {
		return nil, nil, errors.Wrap(err, "verifyHeader; HexInt Conversion Error ")
	}
	header, err := r.getVerifiedHeaderForHeight(uint64(ht))
	if err != nil {
		return nil, nil, errors.Wrap(err, "verify; ")
	} else if header == nil {
		return nil, nil, errors.New("verify; returned nil header")
	}

	rs, err := r.verifyReceipt(header, v)
	if err != nil {
		return nil, nil, errors.Wrap(err, "verify; ")
	}

	return header, rs, nil
}

func (r *receiver) verifyReceipt(header *BlockHeader, v *BlockNotification) ([]*chain.Receipt, error) {
	// Update blockHeaders
	if len(v.Indexes) == 0 || len(v.Events) == 0 {
		r.log.Debug("verifyReceipt; Events and Indexes are empty; Skipping their verification", v.Height)
		return []*chain.Receipt{}, nil
	}

	var headerResult Result
	vlcodec.RLP.MustUnmarshalFromBytes(header.Result, &headerResult)

	r.log.WithFields(log.Fields{"Height": v.Height, "ReceiptHash": headerResult.ReceiptHash, "NumIndexes": len(v.Indexes[0])}).Debug("Receipt Verification Start")
	rps := make([]*chain.Receipt, 0)

	for i, index := range v.Indexes[0] {
		p := &ProofEventsParam{BlockHash: v.Hash, Index: index, Events: v.Events[0][i]}
		proofs, err := r.cl.GetProofForEvents(p)
		if err != nil {
			return nil, errors.Wrap(mapError(err), "verifyReceipt; GetProofForEvents; Err:")
		}
		if len(proofs) != 1+len(p.Events) { // returned proofs should be for all of the requested Events plus 1 for the receipt
			return nil, errors.New("verifyReceipt; Proof Not returned for all requested events")
		}

		// Processing receipt index
		serializedReceipt, err := mptProve(index, proofs[0], headerResult.ReceiptHash)
		if err != nil {
			return nil, errors.Wrap(err, "verifyReceipt; MPTProve Receipt; Err:")
		}
		var receipt Receipt
		_, err = vlcodec.RLP.UnmarshalFromBytes(serializedReceipt, &receipt)
		if err != nil {
			return nil, errors.Wrap(err, "verifyReceipt; Unmarshal Receipt; Err:")
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
				seqGot.SetBytes(el.Indexed[EventIndexSequence])
				evt := &chain.Event{
					Next:     chain.BTPAddress(el.Indexed[EventIndexNext]),
					Sequence: seqGot.Uint64(),
					Message:  el.Data[0],
				}
				rp.Events = append(rp.Events, evt)
			} else {
				if !bytes.Equal(el.Addr, r.evtLogRawFilter.addr) {
					r.log.WithFields(log.Fields{
						"Height":   v.Height,
						"got":      common.HexBytes(el.Addr),
						"expected": common.HexBytes(r.evtLogRawFilter.addr)}).Error("invalid event: cannot match add")
				}
				if !bytes.Equal(el.Indexed[EventIndexSignature], r.evtLogRawFilter.signature) {
					r.log.WithFields(log.Fields{
						"Height":   v.Height,
						"got":      common.HexBytes(el.Indexed[EventIndexSignature]),
						"expected": common.HexBytes(r.evtLogRawFilter.signature)}).Error("invalid event: cannot match sig")
				}
				if !bytes.Equal(el.Indexed[EventIndexNext], r.evtLogRawFilter.next) {
					r.log.WithFields(log.Fields{
						"Height":   v.Height,
						"got":      common.HexBytes(el.Indexed[EventIndexNext]),
						"expected": common.HexBytes(r.evtLogRawFilter.next)}).Error("invalid event: cannot match next")
				}
				return nil, errors.New("verifyReceipt; Invalid event")
			}
		}
		if len(rp.Events) > 0 && len(rp.Events) == len(p.Events) { //Only add if all the events were verified
			rps = append(rps, rp)
		} else if len(rp.Events) > 0 && len(rp.Events) != len(p.Events) {
			r.log.WithFields(log.Fields{
				"Height":       v.Height,
				"ReceiptIndex": index,
				"got":          len(rp.Events),
				"expected":     len(p.Events)}).Info(" Not all events were verified for receipt ")
			return nil, errors.New("verifyReceipt; Not all events were verified for receipt")
		}
	}
	return rps, nil
}

func (r *receiver) getVerifiedHeaderForHeight(ht uint64) (*BlockHeader, error) {
	height := NewHexInt(int64(ht))
	header, err := getBlockHeader(r.cl, height)
	if err != nil {
		return nil, err
	}
	votesBytes, err := r.cl.GetVotesByHeight(&BlockHeightParam{Height: height})
	if err != nil {
		return nil, errors.Wrap(mapError(err), "verifyHeader; GetVotesByHeight; Err: ")
	}
	cvl := &commitVoteList{}
	_, err = vlcodec.BC.UnmarshalFromBytes(votesBytes, cvl)
	if err != nil {
		return nil, errors.Wrap(mapError(err), "verifyHeader; UnmarshalFromBytes of Votes; Err: ")
	}
	blockHash := crypto.SHA3Sum256(vlcodec.BC.MustMarshalToBytes(header))
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
			err = nil
			continue
		}
		address := common.NewAccountAddressFromPublicKey(pub)

		if addressesContains(address.Bytes(), r.hv.validators) {
			validCount++
		}
	}
	if validCount < int(2*len(r.hv.validators)/3) {
		return nil, errors.New("verifyHeader; Block not validated by >= 2/3 of validators")
	}
	//r.log.WithFields(log.Fields{"NumValidators": len(r.hv.validators), "NumValidAddresses": validCount, "Height": ht}).Debug("Block Header Verified by Votes")
	return header, nil
}
