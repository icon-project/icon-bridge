package tezos

import (
	"fmt"
	"sync"

	"strconv"
	"context"

	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"blockwatch.cc/tzgo/codec"
)

type IVerifier interface {
	Next() int64
	Verify(ctx context.Context, header *rpc.BlockHeader, proposer tezos.Address, c *rpc.Client, nextHeader *rpc.BlockHeader) error
	Update(header *rpc.BlockHeader) error
	ParentHash() tezos.BlockHash
	IsValidator(proposer tezos.Address, height int64) bool
	Height() int64 
}

type Verifier struct{
	chainID 		uint32
	mu 				sync.RWMutex
	validators 		[]tezos.Address
	next 			int64
	parentHash 		tezos.BlockHash
	parentFittness	int64
	height 			int64
}

func (vr *Verifier) Next() int64{
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.next
}

func (vr *Verifier) Verify(ctx context.Context, header *rpc.BlockHeader, proposer tezos.Address, c *rpc.Client, nextHeader *rpc.BlockHeader) error {
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

	
	// isValidSignature, err := vr.VerifySignature(ctx, proposer, header.Signature, header.Level, header, c)

	// if !isValidSignature {
	// 	return fmt.Errorf("Invalid block hash. Signature mismatch")
	// }

	fmt.Println(nextHeader.ValidationPass)

	fmt.Println(true)
	return nil 
}

func (vr *Verifier) Update(header *rpc.BlockHeader) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	fmt.Println("updating????")
	blockFittness := header.Fitness

	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	vr.parentFittness = currentFittness

	vr.parentHash = header.Hash
	vr.height = header.Level
	vr.next = header.Level + 1
	fmt.Println(header.Hash)
	fmt.Println("updated")
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
	exposedPublicKey, err := c.GetManagerKey(ctx, proposer, rpc.BlockLevel(blockLevel))
	if err != nil {
		return false, err 
	}

	// c.ListBakingRights()

	blockHeader := codec.BlockHeader{
		Level: 				int32(header.Level),     
		Proto: 				byte(header.Proto),
		Predecessor: 		header.Predecessor,
		Timestamp: 			header.Timestamp,
		ValidationPass: 	byte(header.ValidationPass),
		OperationsHash: 	header.OperationsHash,
		Fitness: 			header.Fitness,
		Context: 			header.Context,
		PayloadHash: 		header.PayloadHash,
		PayloadRound: 		header.PayloadRound,
		ProofOfWorkNonce: 	header.ProofOfWorkNonce,
		LbToggleVote: 		header.LbVote(),
		// SeedNonceHash: 		block.Metadata.NonceHash,		
		ChainId: 			&header.ChainId,	
	}


	digestedHash := blockHeader.Digest()


	err = exposedPublicKey.Verify(digestedHash[:], header.Signature)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil 
}

type VerifierOptions struct {
	BlockHeight 		int64 				`json:"blockHeight"`
	BlockHash 			tezos.BlockHash 	`json:"parentHash"`
}