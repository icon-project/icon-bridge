package tezos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"context"
	"strconv"

	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
)

type IVerifier interface {
	Next() int64
	Verify(ctx context.Context, header *rpc.BlockHeader, block *rpc.Block, proposer tezos.Address, c *rpc.Client, nextHeader *rpc.BlockHeader) error
	Update(ctx context.Context, header *rpc.BlockHeader, block *rpc.Block) error
	ParentHash() tezos.BlockHash
	IsValidator(proposer tezos.Address, height int64) bool
	Height() int64
}

type Verifier struct{
	chainID 		uint32
	mu 				sync.RWMutex
	validators 		map[tezos.Address]bool
	next 			int64
	parentHash 		tezos.BlockHash
	parentFittness	int64
	height 			int64
	cycle 	int64
	c *rpc.Client
}

func (vr *Verifier) Next() int64 {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.next
}

func (vr *Verifier) Verify(ctx context.Context, header *rpc.BlockHeader, block *rpc.Block, proposer tezos.Address, c *rpc.Client, nextHeader *rpc.BlockHeader) error {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	blockFittness := header.Fitness
	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	if currentFittness < vr.parentFittness {
		return fmt.Errorf("Invalid block fittness", currentFittness)
	}

	previousHashInBlock := header.Predecessor

	if previousHashInBlock.String() != vr.parentHash.String() {
		return fmt.Errorf("Invalid block hash", header.Level)
	}

	isValidSignature, err := vr.VerifySignature(ctx, proposer, header.Signature, header.Level, header, c)

	if !isValidSignature {
		return fmt.Errorf("Invalid block hash. Signature mismatch")
	}

	// err = vr.verifyEndorsement(block.Operations, vr.c, block.GetLevel())

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (vr *Verifier) Update(ctx context.Context, header *rpc.BlockHeader, block *rpc.Block) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	fmt.Println("updating for block ????", header.Level)
	blockFittness := header.Fitness

	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	vr.parentFittness = currentFittness

	vr.parentHash = header.Hash
	vr.height = header.Level
	vr.next = header.Level + 1

	if vr.cycle != block.Metadata.LevelInfo.Cycle {
		fmt.Println("reached in updating validators and cycle")
		vr.updateValidatorsAndCycle(ctx, block.Header.Level, block.Metadata.LevelInfo.Cycle)
	}

	return nil
}

func (vr *Verifier) ParentHash() tezos.BlockHash {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.parentHash
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

func (vr *Verifier) VerifySignature(ctx context.Context, proposer tezos.Address, signature tezos.Signature, blockLevel int64, header *rpc.BlockHeader, c *rpc.Client) (bool, error) {
	exposedPublicKey, err := vr.GetConsensusKey(ctx, c, proposer)
	
	if err != nil {
		return false, err
	}

	// c.ListBakingRights()

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
		// SeedNonceHash: 		block.Metadata.NonceHash,
		ChainId: &header.ChainId,
	}

	digestedHash := blockHeader.Digest()

	err = exposedPublicKey.Verify(digestedHash[:], header.Signature)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (vr *Verifier) GetConsensusKey(ctx context.Context, c *rpc.Client, bakerConsensusKey tezos.Address) (tezos.Key, error){
	url := c.BaseURL.String() + "/chains/main/blocks/head/context/raw/json/contracts/index/" + bakerConsensusKey.String() + "/consensus_key/active"

	fmt.Println(c.BaseURL)

	resp, err := http.Get(url)
	fmt.Println(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tezos.Key{}, err
	}
	//Convert the body to type string
	sb := string(body)

	exposedPublicKey := tezos.MustParseKey(sb)
	return exposedPublicKey, nil 
}

func (vr *Verifier) updateValidatorsAndCycle(ctx context.Context, blockHeight int64, cycle int64) error {

	if true {
		return nil
	}

	fmt.Println("reached update validators")
	validatorsList, err := vr.c.ListEndorsingRights(ctx, rpc.BlockLevel(blockHeight))
	if err != nil {
		fmt.Println("error here?", err)
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
	}
	vr.cycle = cycle
	fmt.Println("reached to updating cycle")
	return nil
}

func (vr *Verifier) verifyEndorsement(op [][]*rpc.Operation, c *rpc.Client, blockHeight int64) error {
	endorsementPower := 0

	threshold := 7000 * float32(2) / float32(3)
	endorsers := make(map[tezos.Address]bool)
	for i := 0; i < len(op); i++ {
		for j := 0; j < len(op[i]); j++ {
			for _, operation := range op[i][j].Contents {
				switch operation.Kind() {
				case tezos.OpTypeEndorsement:
					tx := operation.(*rpc.Endorsement)
					if _, isDelegate := vr.validators[tx.Metadata.Delegate]; isDelegate {
						if _, ok := endorsers[tx.Metadata.Delegate]; !ok {
							endorsers[tx.Metadata.Delegate] = true
							endorsementPower += tx.Metadata.EndorsementPower
						}
					} else {
						fmt.Println(vr.validators[tx.Metadata.Delegate])
					}
				}
			}
		}
	}
	fmt.Println(len(endorsers))

	if endorsementPower > int(threshold) && len(endorsers)*100/len(vr.validators) > 66 {
		return nil
	}
	return errors.New("endorsement verification failed")

}

type VerifierOptions struct {
	BlockHeight int64           `json:"blockHeight"`
	BlockHash   tezos.BlockHash `json:"parentHash"`
}
