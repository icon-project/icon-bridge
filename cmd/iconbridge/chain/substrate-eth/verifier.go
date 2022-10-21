package substrate_eth

import (
	"bytes"
	"fmt"
	"math/big"
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/common"
)

var (
	big1 = big.NewInt(1)
)

type VerifierOptions struct {
	BlockHeight uint64          `json:"blockHeight"`
	BlockHash   common.HexBytes `json:"parentHash"`
}

// next points to height whose parentHash is expected
// parentHash of height h is got from next-1's hash
type Verifier struct {
	mu         sync.RWMutex
	next       *big.Int
	parentHash ethCommon.Hash
}

func (vr *Verifier) Next() *big.Int {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return (&big.Int{}).Set(vr.next)
}

func (vr *Verifier) Verify(h *types.Header, newHeader *types.Header) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	if newHeader.Number.Cmp((&big.Int{}).Add(h.Number, big1)) != 0 {
		return fmt.Errorf("Different height between successive header: Prev %v New %v", h.Number, newHeader.Number)
	}
	if !bytes.Equal(h.Hash().Bytes(), newHeader.ParentHash.Bytes()) {
		return fmt.Errorf("Different hash between successive header: (%v): Prev %v New %v", h.Number, h.Hash(), newHeader.ParentHash)
	}
	if vr.next.Cmp(h.Number) != 0 {
		return fmt.Errorf("Unexpected height: Got %v Expected %v", h.Number, vr.next)
	}
	if !bytes.Equal(h.ParentHash.Bytes(), vr.parentHash.Bytes()) {
		return fmt.Errorf("Unexpected Hash(%v): Got %v Expected %v", h.Number, h.ParentHash, vr.parentHash)
	}
	vr.parentHash = h.Hash()
	vr.next.Add(h.Number, big1)
	return nil
}

// func (vr *Verifier) Update(h *types.Header) error {
// 	vr.mu.Lock()
// 	defer vr.mu.Unlock()
// 	// next height should have vr.parentHash as parentHash
// 	return nil
// }
