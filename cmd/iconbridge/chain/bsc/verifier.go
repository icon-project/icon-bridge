package bsc

import (
	"bytes"
	"math/big"
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/common"
	"github.com/pkg/errors"
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
	h := new(big.Int)
	h.SetString(vr.next.String(), 10)
	return h
}

func (vr *Verifier) Verify(h *types.Header) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	if vr.next.Cmp(h.Number) != 0 {
		return errors.New("Unexpected height")
	}
	if !bytes.Equal(h.ParentHash.Bytes(), vr.parentHash.Bytes()) {
		return errors.New("Unexpected Hash")
	}
	vr.parentHash = h.Hash()
	vr.next.Add(vr.next, big.NewInt(1)) // next height should have vr.parentHash as parentHash
	return nil
}

// func (vr *Verifier) Update(h *types.Header) error {
// 	return vr.Verify(h)
// }
