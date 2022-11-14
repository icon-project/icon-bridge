package types

import (
	"fmt"
	"crypto/sha256"
	"encoding/json"
	"github.com/near/borsh-go"
)

type ValidatorStakeStructVersion borsh.Enum

const (
	ValidatorStakeStructVersion1 ValidatorStakeStructVersion = iota
	ValidatorStakeStructVersion2
  )

func (vs *ValidatorStakeStructVersion) UnmarshalJSON(p []byte) error {
	var validatorStakeStructVersion string
	err := json.Unmarshal(p, &validatorStakeStructVersion)
	if err != nil {
		return err
	}

	if validatorStakeStructVersion == "" {
		*vs = ValidatorStakeStructVersion1
		return nil
	}

	switch validatorStakeStructVersion {
	case "V1":
		*vs = ValidatorStakeStructVersion1
	default:
		return fmt.Errorf("not supported validator struct")
	}
	return nil
}

type BlockProducer struct {
	ValidatorStakeStructVersion ValidatorStakeStructVersion `json:"validator_stake_struct_version"`
	AccountId                   AccountId                   `json:"account_id"`
	PublicKey                   PublicKey                   `json:"public_key"`
	Stake                       BigInt                      `json:"stake"`
}

type BlockProducers []BlockProducer

func (bps *BlockProducers) BorshSerialize() ([]byte, error) {
	return borsh.Serialize(*bps)
}

func (bps *BlockProducers) Hash() (CryptoHash, error) {
	serializedBps, err := bps.BorshSerialize()

	if err != nil {
		return CryptoHash{}, err
	}

	return sha256.Sum256(serializedBps), nil
}
