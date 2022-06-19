package icon

import (
	"fmt"
	"sync"

	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/common/crypto"
)

type VerifierOptions struct {
	BlockHeight    uint64          `json:"blockHeight"`
	ValidatorsHash common.HexBytes `json:"validatorsHash"`
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

type BlockHeaderResult struct {
	StateHash        []byte
	PatchReceiptHash []byte
	ReceiptHash      common.HexBytes
	ExtensionData    []byte
}

type TxResult struct {
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

type Verifier struct {
	mu                 sync.RWMutex
	next               int64
	nextValidatorsHash common.HexBytes
	validators         map[string][]common.HexBytes
}

func (vr *Verifier) Next() int64 { return vr.next }

func (vr *Verifier) Verify(h *BlockHeader, votes []byte) (ok bool, err error) {
	vr.mu.RLock()
	defer vr.mu.RUnlock()

	nvh := vr.nextValidatorsHash
	validators, ok := vr.validators[nvh.String()]
	if !ok {
		return false, fmt.Errorf("no validators for hash=%v", nvh)
	}

	cvl := &commitVoteList{}
	_, err = codec.BC.UnmarshalFromBytes(votes, cvl)
	if err != nil {
		return false, fmt.Errorf("invalid votes: %v; err=%v", common.HexBytes(votes), err)
	}

	hash := crypto.SHA3Sum256(codec.BC.MustMarshalToBytes(h))
	vote := &vote{
		voteBase: voteBase{
			_HR: _HR{
				Height: h.Height,
				Round:  cvl.Round,
			},
			Type:           VoteTypePrecommit,
			BlockID:        hash,
			BlockPartSetID: cvl.BlockPartSetID,
		},
	}

	votesCount := 0
	for _, item := range cvl.Items {
		vote.Timestamp = item.Timestamp
		pub, err := item.Signature.RecoverPublicKey(crypto.SHA3Sum256(codec.BC.MustMarshalToBytes(vote)))
		if err != nil {
			continue
		}
		address := common.NewAccountAddressFromPublicKey(pub)
		if listContains(validators, address.Bytes()) {
			votesCount++
		}
	}
	if votesCount < (2*len(validators))/3 {
		return false, fmt.Errorf("insufficient votes")
	}

	return true, nil
}

func (vr *Verifier) Update(h *BlockHeader, nextValidators []common.HexBytes) (err error) {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	nvh := common.HexBytes(h.NextValidatorsHash)
	if _, ok := vr.validators[nvh.String()]; !ok {
		vr.validators[nvh.String()] = nextValidators
	}
	vr.next = h.Height + 1
	vr.nextValidatorsHash = h.NextValidatorsHash
	return nil
}

func (vr *Verifier) Validators(nvh common.HexBytes) []common.HexBytes {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	validators, ok := vr.validators[nvh.String()]
	if ok {
		return validators
	}
	return nil
}

// func (r *receiver) syncVerifier(hexHeight HexInt) error {
// 	ht, hterr := hexHeight.Value()
// 	if hterr != nil {
// 		return errors.Wrapf(hterr, "syncVerifier; HexInt Conversion Error at Height %v ", hexHeight)
// 	}
// 	targetHeight := uint64(ht)

// 	if targetHeight == r.hv.height {
// 		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorsHeight": NewHexInt(int64(r.hv.height))}).Info("SyncVerifier; Same Height so already in sync ")
// 		return nil
// 	} else if targetHeight < r.hv.height {
// 		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorsHeight": NewHexInt(int64(r.hv.height)), "ValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; Start Syncing Backward ")
// 		if targetHeight == 1 {
// 			return errors.New("Cannot find validator list for block of height 1")
// 		}

// 		verifiedHeader, err := r.getVerifiedHeaderForHeight(r.hv.height)
// 		if err != nil {
// 			return errors.Wrap(err, "syncVerifier; ")
// 		}
// 		for ht := r.hv.height - 1; ht >= targetHeight-1; ht-- {
// 			if verifiedHeader, err = r.verifyWithPreviousHeader(verifiedHeader); err != nil {
// 				return err
// 			}
// 			if ht%50 == 0 {
// 				r.log.WithFields(log.Fields{"ValidatorsHeight": NewHexInt(int64(r.hv.height)), "TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; In Progress ")
// 			}
// 		}
// 		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(r.hv.height)), "ValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; Complete ")
// 	} else {
// 		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorsHeight": NewHexInt(int64(r.hv.height)), "ValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; Start Syncing Forward ")
// 		for ht := r.hv.height; ht < targetHeight; ht++ {
// 			if header, err := r.getVerifiedHeaderForHeight(ht); err != nil {
// 				return errors.Wrap(err, "syncVerifier; ")
// 			} else {
// 				if !bytes.Equal(header.NextValidatorsHash, r.hv.validatorsHash) { // should update
// 					if vs, err := getValidatorsFromHash(r.cl, header.NextValidatorsHash); err != nil {
// 						return errors.Wrap(err, "syncVerifier; ")
// 					} else {
// 						r.log.WithFields(log.Fields{"ValidatorsHeight": NewHexInt(int64(ht + 1)), "NewValidatorsHash": common.HexBytes(header.NextValidatorsHash), "OldValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; Updating Validator Hash ")
// 						r.hv.validatorsHash = header.NextValidatorsHash
// 						r.hv.validators = vs
// 					}
// 				}
// 				r.hv.height = ht + 1
// 			}
// 			if ht%50 == 0 {
// 				r.log.WithFields(log.Fields{"ValidatorsHeight": NewHexInt(int64(r.hv.height)), "TargetHeight": NewHexInt(int64(targetHeight)), "ValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; In Progress ")
// 			}
// 		}
// 		r.log.WithFields(log.Fields{"TargetHeight": NewHexInt(int64(r.hv.height)), "ValidatorsHeight": NewHexInt(int64(r.hv.height)), "ValidatorsHash": r.hv.validatorsHash}).Info("Sync Verifier; Complete ")
// 	}
// 	return nil
// }
