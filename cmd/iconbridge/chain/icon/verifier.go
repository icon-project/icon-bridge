package icon

import (
	"fmt"
	"sync"

	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/common/crypto"
)

type VerifierOptions struct {
	BlockHeight    uint64         `json:"blockHeight"`
	ValidatorsHash common.HexHash `json:"validatorsHash"`
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
	TxIndex            HexInt
	BlockHeight        HexInt
}

type Verifier struct {
	mu                 sync.RWMutex
	next               int64
	nextValidatorsHash common.HexHash
	validators         map[string][]common.Address // convert this to lru cache
}

func (vr *Verifier) Next() int64 { return vr.next }

func (vr *Verifier) Verify(blockHeader *BlockHeader, votes []byte) (ok bool, err error) {
	vr.mu.RLock()
	defer vr.mu.RUnlock()

	nextValidatorsHash := vr.nextValidatorsHash
	listValidators, ok := vr.validators[nextValidatorsHash.String()]
	if !ok {
		return false, fmt.Errorf("no validators for hash=%v", nextValidatorsHash)
	}

	requiredVotes := (2 * len(listValidators)) / 3
	if requiredVotes < 1 {
		requiredVotes = 1
	}

	cvl := &commitVoteList{}
	_, err = codec.BC.UnmarshalFromBytes(votes, cvl)
	if err != nil {
		return false, fmt.Errorf("invalid votes: %v; err=%v", common.HexBytes(votes), err)
	}

	hash := crypto.SHA3Sum256(codec.BC.MustMarshalToBytes(blockHeader))
	vote := &vote{
		voteBase: voteBase{
			_HR: _HR{
				Height: blockHeader.Height,
				Round:  cvl.Round,
			},
			Type:           VoteTypePrecommit,
			BlockID:        hash,
			BlockPartSetID: cvl.BlockPartSetID,
		},
	}

	numVotes := 0
	validators := make(map[common.Address]struct{})
	for _, val := range listValidators {
		validators[val] = struct{}{}
	}

	for _, item := range cvl.Items {
		vote.Timestamp = item.Timestamp
		pub, err := item.Signature.RecoverPublicKey(crypto.SHA3Sum256(codec.BC.MustMarshalToBytes(vote)))
		if err != nil {
			continue // skip error
		}
		address := common.NewAccountAddressFromPublicKey(pub)
		if address == nil {
			continue
		}
		if _, ok := validators[*address]; !ok {
			continue // already voted or invalid validator
		}
		delete(validators, *address)
		if numVotes++; numVotes >= requiredVotes {
			return true, nil
		}
	}

	return false, fmt.Errorf("insufficient votes")
}

func (vr *Verifier) Update(blockHeader *BlockHeader, nextValidators []common.Address) (err error) {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	nextValidatorsHash := common.HexBytes(blockHeader.NextValidatorsHash)

	if _, ok := vr.validators[nextValidatorsHash.String()]; !ok {
		vr.validators[nextValidatorsHash.String()] = nextValidators
	}

	vr.next = blockHeader.Height + 1
	vr.nextValidatorsHash = blockHeader.NextValidatorsHash
	return nil
}

func (vr *Verifier) Validators(nextValidatorsHash common.HexBytes) []common.Address {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	validators, ok := vr.validators[nextValidatorsHash.String()]
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
