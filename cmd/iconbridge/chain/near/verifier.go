package near

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/reactivex/rxgo/v2"
)

type Verifier struct {
	mu                sync.RWMutex
	blockHeight       uint64
	previousBlockHash types.CryptoHash
	currentEpochId    types.CryptoHash
	nextEpochId       types.CryptoHash
	currentBpsHash    types.CryptoHash
	nextBpsHash       types.CryptoHash
	blockProducers    types.BlockProducers
	SyncConcurrency   int
	client            IClient
}

func NewVerifier(blockHeight uint64, previousBlockHash, currentEpochId, nextEpochId, currentBpsHash, nextBpsHash types.CryptoHash, SyncConcurrency int, client IClient) (*Verifier, error) {
	v := &Verifier{
		blockHeight:       blockHeight,
		previousBlockHash: previousBlockHash,
		client:            client,
		currentEpochId:    currentEpochId,
		nextEpochId:       nextEpochId,
		nextBpsHash:       nextBpsHash,
		currentBpsHash:    currentBpsHash,
		SyncConcurrency:   SyncConcurrency,
	}

	bps, err := client.GetBlockProducers(previousBlockHash)
	if err != nil {
		return nil, err
	}

	bpsHash, err := bps.Hash()
	if err != nil {
		return nil, err
	}

	if currentBpsHash != bpsHash {
		return nil, fmt.Errorf("expected block producers hash: %v, got block producers hash: %v for epoch: %v", currentBpsHash.Base58Encode(), bpsHash.Base58Encode(), currentEpochId.Base58Encode())
	}

	v.blockProducers = bps

	return v, nil
}

func (v *Verifier) SyncHeader(wg *sync.WaitGroup, target uint64) error {
	defer wg.Done()

	err := v.client.MonitorBlocks(v.blockHeight, target, v.SyncConcurrency, func(observable rxgo.Observable) error {
		result := observable.Serialize(int(v.blockHeight),
			func(bn interface{}) int {
				return int(bn.(*types.BlockNotification).Offset())
			},
		).Filter(v.client.FilterUnknownBlocks).TakeUntil(
			func(bn interface{}) bool {
				return bn.(*types.BlockNotification).Block().Height() >= int64(target)
			},
		).Observe()

		for item := range result {
			if err := item.E; err != nil {
				return err
			}

			bn, _ := item.V.(*types.BlockNotification)
			v.ValidateHeader(bn)

			v.client.Logger().WithFields(log.Fields{"height": bn.Block().Height()}).Debug("syncing verifier")
		}

		return nil
	})

	return err
}

func (v *Verifier) ValidateHeader(blockNotification *types.BlockNotification) (err error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	totalStake := big.NewInt(0)
	approvedStake := big.NewInt(0)
	approvalMessage := types.ApprovalMessage{}

	if uint64(blockNotification.Block().Header.PreviousHeight+1) == uint64(blockNotification.Offset()) {
		approvalMessage = types.ApprovalMessage{
			Type:              [1]byte{types.ApprovalEndorsement},
			PreviousBlockHash: v.previousBlockHash,
		}
	} else {
		approvalMessage = types.ApprovalMessage{
			Type:                [1]byte{types.ApprovalSkip},
			PreviousBlockHeight: v.blockHeight,
		}
	}

	approvalMessage.TargetHeight = uint64(blockNotification.Offset())

	currBlockHash, err := blockNotification.Block().ComputeHash(v.previousBlockHash, blockNotification.Block().InnerLite(), blockNotification.Block().InnerRest())
	if err != nil {
		return err
	}

	if currBlockHash != *blockNotification.Block().Hash() {
		return fmt.Errorf("expected hash: %v, got hash: %v for block: %v", blockNotification.Block().Hash().Base58Encode(), currBlockHash.Base58Encode(), blockNotification.Block().Height())
	}

	v.blockHeight = uint64(blockNotification.Offset())
	v.previousBlockHash = currBlockHash

	if blockNotification.Block().Header.EpochId != v.currentEpochId && blockNotification.Block().Header.EpochId != v.nextEpochId {
		return fmt.Errorf("block in invalid epoch")
	}

	if blockNotification.Block().Header.EpochId == v.nextEpochId {
		v.client.Logger().WithFields(log.Fields{"block": blockNotification.Block().Hash().Base58Encode()}).Debug("fetching block producers")
		bps, err := v.client.GetBlockProducers(*blockNotification.Block().Hash())
		if err != nil {
			return err
		}

		bpsHash, err := bps.Hash()
		if err != nil {
			return err
		}

		if bpsHash != v.nextBpsHash {
			return fmt.Errorf("expected block producers hash: %v, got block producers hash: %v for epoch: %v", v.nextBpsHash.Base58Encode(), bpsHash.Base58Encode(), blockNotification.Block().Header.EpochId.Base58Encode())
		}

		v.blockProducers = bps
		v.currentEpochId = blockNotification.Block().Header.EpochId
		v.nextEpochId = blockNotification.Block().Header.NextEpochId
		v.nextBpsHash = blockNotification.Block().Header.NextBlockProducerHash
	}


	for i, approval := range blockNotification.Block().Header.Approvals {
		if i >= len(v.blockProducers) {
			break
		}

		totalStake = totalStake.Add(totalStake, (*big.Int)(&v.blockProducers[i].Stake))

		if approval == nil {
			continue
		}

		approvedStake = approvedStake.Add(approvedStake, (*big.Int)(&v.blockProducers[i].Stake))

		if err := approvalMessage.Verify(&v.blockProducers[i].PublicKey, approval); err != nil {
			return err
		}
	}

	threshold := big.NewInt(0).Div(totalStake.Mul(totalStake, big.NewInt(2)), big.NewInt(3))
	if approvedStake.Cmp(threshold) <= 0 {
		return fmt.Errorf("approved stake: %v below threshold: %v", approvedStake, threshold)
	}

	return
}
