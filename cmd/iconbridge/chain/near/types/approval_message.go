package types

import (
	"bytes"
	"fmt"
	"github.com/near/borsh-go"
)

const (
	ApprovalEndorsement = 0
	ApprovalSkip        = 1
)

type ApprovalMessage struct {
	Type                [1]byte
	PreviousBlockHash   CryptoHash
	PreviousBlockHeight uint64
	TargetHeight        uint64
}

func (a ApprovalMessage) BorshSerialize() ([]byte, error) {
	serialized := new(bytes.Buffer)
	serialized.Write(a.Type[:])

	if a.Type == [1]byte{ApprovalEndorsement} {
		serialized.Write(a.PreviousBlockHash[:])
	} else {
		PreviousBlockHeight, err := borsh.Serialize(a.PreviousBlockHeight)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize ApprovalMessage: PreviousBlockHeight")
		}

		serialized.Write(PreviousBlockHeight)
	}

	targetHeight, err := borsh.Serialize(a.TargetHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ApprovalMessage: TargetHeight")
	}
	serialized.Write(targetHeight)

	return serialized.Bytes(), nil
}

func (a ApprovalMessage) Verify(p PublicKey, s Signature) (bool, error)  {
	_, err := a.BorshSerialize()
	if err != nil {
		return false, fmt.Errorf("failed to Verify ApprovalMessage: %v", err)
	}

	return false, nil
}
