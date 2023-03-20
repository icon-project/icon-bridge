package tezos

import (
	"fmt"
	"sync"

	"strconv"

	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
)

type IVerifier interface {
	Next() int64
	Verify(header *rpc.BlockHeader, proposer *tezos.Address) error
	Update(header *rpc.BlockHeader) error
	ParentHash() tezos.BlockHash
	IsValidator(proposer tezos.Address, height int64) bool	
}

type Verifier struct{
	chainID 		uint32
	mu 				sync.RWMutex
	validators 		[]tezos.Address
	next 			int64
	parentHash 		tezos.BlockHash
	parentFittness	int64
}

func (vr *Verifier) Next() int64{
	vr.mu.RLock()
	return vr.next
}

func (vr *Verifier) Verify(header *rpc.BlockHeader, proposer *tezos.Address) error {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	blockFittness := header.Fitness
	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	fmt.Print("Current fittness: ")
	fmt.Println(currentFittness)

	fmt.Print("Parent fittness")
	fmt.Println(vr.parentFittness)

	if currentFittness < vr.parentFittness {
		return fmt.Errorf("Invalid block fittness")
	}

	previousHashInBlock := header.Predecessor

	fmt.Print("Current fittness: ")
	fmt.Println(previousHashInBlock)

	fmt.Print("Parent fittness")
	fmt.Println(vr.parentHash)


	if previousHashInBlock.String() != vr.parentHash.String() {
		return fmt.Errorf("Invalid block hash")
	}
	fmt.Println("Block is verified")
	fmt.Println("*******         *******")
	fmt.Println("  *******     *******")
	fmt.Println("    ******* *******")


	return nil 
}

func (vr *Verifier) Update(header *rpc.BlockHeader) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	blockFittness := header.Fitness

	currentFittness, err := strconv.ParseInt(string(blockFittness[1].String()), 16, 64)
	if err != nil {
		return err
	}

	vr.parentFittness = currentFittness

	vr.parentHash = header.Hash
	fmt.Println(header.Hash)
	fmt.Println("updated")
	return nil 
}

func (vr *Verifier) ParentHash() tezos.BlockHash {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return vr.parentHash
}

func (vr *Verifier) IsValidator(proposer tezos.Address, height int64) bool {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return true
}