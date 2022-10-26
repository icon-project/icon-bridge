package substrate_eth

import (
	"fmt"
	subEthTypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/substrate-eth/types"
	"math/big"
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"
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

func (vr *Verifier) Verify(h *subEthTypes.Header, newHeader *subEthTypes.Header) error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	if newHeader.Number.Cmp((&big.Int{}).Add(h.Number, big1)) != 0 {
		return fmt.Errorf("Different height between successive header: Prev %v New %v", h.Number, newHeader.Number)
	}
	if h.Hash != newHeader.ParentHash {
		return fmt.Errorf("Different hash between successive header: (%v): Prev %v New %v", h.Number, h.Hash, newHeader.ParentHash)
	}
	if vr.next.Cmp(h.Number) != 0 {
		return fmt.Errorf("Unexpected height: Got %v Expected %v", h.Number, vr.next)
	}
	if h.ParentHash != vr.parentHash {
		return fmt.Errorf("Unexpected Hash(%v): Got %v Expected %v", h.Number, h.ParentHash, vr.parentHash)
	}
	vr.parentHash = h.Hash
	vr.next.Add(h.Number, big1)
	return nil
}

