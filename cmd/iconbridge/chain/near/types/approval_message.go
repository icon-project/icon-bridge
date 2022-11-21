package types

import (
	"bytes"
	"crypto/ed25519"
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

func (a *ApprovalMessage) BorshSerialize() ([]byte, error) {
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

func (a *ApprovalMessage) Verify(p *PublicKey, s *Signature) error {
	approvalMessage, err := a.BorshSerialize()

	if !ed25519.Verify(p.Data[:], approvalMessage, s.Data[:]) {
		return fmt.Errorf("invalid signature: %v for block producer: %v", s.Base58Encode(), p.Base58Encode())
	}

	if err != nil {
		return fmt.Errorf("failed to Verify ApprovalMessage: %v", err)
	}

	return nil
}
