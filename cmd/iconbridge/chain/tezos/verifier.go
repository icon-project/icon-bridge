package tezos

import (
	"fmt"
	"sync"

	"context"
	"strconv"

	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/tezos/types"
)

type IVerifier interface {
	Next() int64
	Verify(ctx context.Context, lbn *types.BlockNotification) error
	Update(ctx context.Context, lbn *types.BlockNotification) error
	ParentHash() tezos.BlockHash
	IsValidator(proposer tezos.Address, height int64) bool
	Height() int64
	LastVerifiedBn() *types.BlockNotification
}

const (
	threshold = 4667
)

type Verifier struct {
	chainID             uint32
	mu                  sync.RWMutex
	validators          map[tezos.Address]bool
	validatorsPublicKey map[tezos.Address]tezos.Key
	next                int64
	parentHash          tezos.BlockHash
	parentFittness      int64
	height              int64
	cycle               int64
	lastVerifiedBn      *types.BlockNotification
	cl                  *Client
}

func (vr *Verifier) Next() int64 {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.next
}

func (vr *Verifier) Verify(ctx context.Context, lbn *types.BlockNotification) error {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	blockFittness := lbn.Header.Fitness
	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	if currentFittness < vr.parentFittness {
		return fmt.Errorf("invalid block fittness %d", currentFittness)
	}

	previousHashInBlock := lbn.Block.Header.Predecessor

	if previousHashInBlock.String() != vr.parentHash.String() {
		return fmt.Errorf("invalid block hash %d", lbn.Header.Level)
	}

	isValidSignature, _ := vr.VerifySignature(ctx, lbn)

	if !isValidSignature {
		return fmt.Errorf("invalid block signature. signature mismatch")
	}

	err = vr.verifyEndorsement(lbn.Block, lbn.Header.ChainId)
	if err != nil {
		return fmt.Errorf("invlid endorsement")
	}

	return nil
}

func (vr *Verifier) Update(ctx context.Context, lbn *types.BlockNotification) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	header := lbn.Header
	block := lbn.Block
	blockFittness := header.Fitness

	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	vr.parentFittness = currentFittness
	vr.parentHash = block.Hash
	vr.height = header.Level
	vr.next = header.Level + 1

	// if vr.cycle != block.Metadata.LevelInfo.Cycle {
	// 	vr.updateValidatorsAndCycle(ctx, block.Header.Level, block.Metadata.LevelInfo.Cycle)
	// }

	vr.lastVerifiedBn = lbn
	return nil
}

func (vr *Verifier) ParentHash() tezos.BlockHash {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.parentHash
}

func (vr *Verifier) LastVerifiedBn() *types.BlockNotification {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.lastVerifiedBn
}

func (vr *Verifier) Height() int64 {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.height
}

func (vr *Verifier) IsValidator(proposer tezos.Address, height int64) bool {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return true
}

func (vr *Verifier) VerifySignature(ctx context.Context, lbn *types.BlockNotification) (bool, error) {
	header := lbn.Block.Header

	blockHeader := codec.BlockHeader{
		Level:            int32(header.Level),
		Proto:            byte(header.Proto),
		Predecessor:      header.Predecessor,
		Timestamp:        header.Timestamp,
		ValidationPass:   byte(header.ValidationPass),
		OperationsHash:   header.OperationsHash,
		Fitness:          header.Fitness,
		Context:          header.Context,
		PayloadHash:      header.PayloadHash,
		PayloadRound:     header.PayloadRound,
		ProofOfWorkNonce: header.ProofOfWorkNonce,
		LbToggleVote:     header.LbVote(),
		SeedNonceHash:    lbn.Block.Metadata.NonceHash,
		ChainId:          &lbn.Block.ChainId,
	}

	digestedHash := blockHeader.Digest()

	err := vr.validatorsPublicKey[lbn.Block.Metadata.Baker].Verify(digestedHash[:], header.Signature)

	if err != nil {
		panic("signature failed")
		// return false, err
	}

	return true, nil
}

func (vr *Verifier) updateValidatorsAndCycle(ctx context.Context, blockHeight int64, cycle int64) error {
	PrintSync()
	validatorsList, err := vr.cl.Cl.ListEndorsingRights(ctx, rpc.BlockLevel(blockHeight))
	var validatorsPublicKey tezos.Key
	if err != nil {
		return err
	}
	// remove all validators
	for a := range vr.validators {
		delete(vr.validators, a)
	}
	vr.validators = make(map[tezos.Address]bool)
	// add new validators
	for _, validator := range validatorsList {
		vr.validators[validator.Delegate] = true
		validatorsPublicKey, err = vr.cl.GetConsensusKey(ctx, validator.Delegate)
		if err != nil {
			return err
		}
		vr.validatorsPublicKey[validator.Delegate] = validatorsPublicKey
	}
	vr.cycle = cycle
	return nil
}

func (vr *Verifier) verifyEndorsement(block *rpc.Block, chainID tezos.ChainIdHash) error {
	endorsementPower := 0
	endorsers := make(map[tezos.Address]bool)
	op := block.Operations
	for i := 0; i < len(op); i++ {
		for j := 0; j < len(op[i]); j++ {
			for _, operation := range op[i][j].Contents {
				signature := op[i][j].Signature
				branch := op[i][j].Branch
				switch operation.Kind() {
				case tezos.OpTypeEndorsement:
					tx := operation.(*rpc.Endorsement)

					if _, isDelegate := vr.validators[tx.Metadata.Delegate]; isDelegate {
						endorsement := codec.TenderbakeEndorsement{
							Slot:             int16(tx.Slot),
							Level:            int32(tx.Level),
							Round:            int32(tx.Round),
							BlockPayloadHash: tx.PayloadHash,
						}
						digested := codec.NewOp().WithContentsFront(&endorsement).WithChainId(block.ChainId).WithBranch(branch).Digest()

						managerKey := vr.validatorsPublicKey[tx.Metadata.Delegate]

						err := managerKey.Verify(digested[:], signature)

						if err != nil {
							panic("signature unverified")
							// return err
						}

						if _, ok := endorsers[tx.Metadata.Delegate]; !ok {
							endorsers[tx.Metadata.Delegate] = true
							endorsementPower += tx.Metadata.EndorsementPower
						}
					}
				}
			}
		}
	}
	if endorsementPower > int(threshold) { // && len(endorsers)*100/len(vr.validators) >= 66 {
		return nil
	}
	panic("threshold didnot meet")
	// return errors.New("endorsement verification failed")

}

type VerifierOptions struct {
	BlockHeight int64           `json:"blockHeight"`
	BlockHash   tezos.BlockHash `json:"parentHash"`
}
