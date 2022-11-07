package near

import (
	"crypto/ed25519"
	"fmt"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type Verifier struct {
	blockHeight    uint64
	blockHash      types.CryptoHash
	nextEpochId    types.CryptoHash
	nextBpHash     types.CryptoHash
	blockProducers types.BlockProducers
}

func newVerifier(blockHeight uint64, blockHash types.CryptoHash, nextEpochId types.CryptoHash, nextBpHash types.CryptoHash, blockProducers []*types.BlockProducer) *Verifier {
	return &Verifier{
		blockHeight,
		blockHash,
		nextEpochId,
		nextBpHash,
		blockProducers,
	}
}

func (v *Verifier) validateHeader(blockNotification *types.BlockNotification) (err error) {
	totalStake := big.NewInt(0)
	approvedStake := big.NewInt(0)

	approvalMessage, err := blockNotification.ApprovalMessage().BorshSerialize()
	if err != nil {
		return err
	}

	if blockNotification.Block().Header.EpochId == v.nextEpochId {
		nextBpHash, err := blockNotification.BlockProducers().Hash()
		if err != nil {
			return err
		}

		if nextBpHash != v.nextBpHash {
			return fmt.Errorf("invalid block producers")
		}

		v.nextBpHash = blockNotification.Block().Header.NextBlockProducerHash
		v.blockProducers = *blockNotification.BlockProducers()
		v.nextEpochId = blockNotification.Block().Header.NextEpochId
	}

	v.blockHeight = uint64(blockNotification.Offset())
	v.blockHash, err = blockNotification.Block().ComputeHash(v.blockHash, blockNotification.Block().InnerLite(), blockNotification.Block().InnerRest())
	if err != nil {
		return
	}

	if v.blockHash != *blockNotification.Block().Hash() {
		return fmt.Errorf("expected hash: %v, got hash: %v", blockNotification.Block().Hash().Base58Encode(), v.blockHash.Base58Encode())
	}

	if len(v.blockProducers) != len(blockNotification.Block().Header.Approvals) {
		return fmt.Errorf("invalid length of block producers")
	}

	for i, approval := range blockNotification.Block().Header.Approvals {
		totalStake = totalStake.Add(totalStake, (*big.Int)(&v.blockProducers[i].Stake))

		if approval == nil {
			continue
		}

		approvedStake = approvedStake.Add(approvedStake, (*big.Int)(&v.blockProducers[i].Stake))

		if !ed25519.Verify(v.blockProducers[i].PublicKey.Data[:], approvalMessage, approval.Data[:]) {
			return fmt.Errorf("invalid signature:%v for block producer: %v", approval.Base58Encode(), v.blockProducers[i].PublicKey.Base58Encode())
		}
	}

	threshold := big.NewInt(0).Div(totalStake.Mul(totalStake, big.NewInt(2)), big.NewInt(3))
	if (approvedStake.Cmp(threshold) <= 0) {
		return fmt.Errorf("invalid stake")
	}

	return
}

func (v *Verifier) validateReceiptsOutcome() error {
	return nil
}
