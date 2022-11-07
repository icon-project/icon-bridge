package near

import (
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type Verifier struct {
	blockHeight    uint64
	blockHash      types.CryptoHash
	nextEpochId    types.CryptoHash
	blockProducers []*types.BlockProducer
}

func newVerifier(blockHeight uint64, blockHash types.CryptoHash, nextEpochId types.CryptoHash, blockProducers []*types.BlockProducer) *Verifier {
	return &Verifier{
		blockHeight,
		blockHash,
		nextEpochId,
		blockProducers,
	}
}

func (v *Verifier) validateHeader(blockNotification *types.BlockNotification) (err error) {
	if blockNotification.Block().Header.EpochId == v.nextEpochId {
		v.blockProducers = blockNotification.BlockProducers()
	}

	v.blockHeight = uint64(blockNotification.Offset())
	v.blockHash, err = blockNotification.Block().ComputeHash(v.blockHash, blockNotification.Block().InnerLite(), blockNotification.Block().InnerRest())
	if err != nil {
		return
	}

	if v.blockHash != *blockNotification.Block().Hash() {
		return fmt.Errorf("expected hash: %v, got hash: %v", blockNotification.Block().Hash().Base58Encode(), v.blockHash.Base58Encode())
	}



	return err
}

func (v *Verifier) validateReceiptsOutcome() error {
	return nil
}
