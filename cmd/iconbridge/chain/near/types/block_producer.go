package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/near/borsh-go"
)

type ValidatorStakeStructVersion []byte

func (vs *ValidatorStakeStructVersion) UnmarshalJSON(p []byte) error {
	var validatorStakeStructVersion string
	err := json.Unmarshal(p, &validatorStakeStructVersion)
	if err != nil {
		return err
	}

	if validatorStakeStructVersion == "" {
		*vs = nil
		return nil
	}

	switch validatorStakeStructVersion {
	case "V1":
		*vs = []byte{0}
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

type BlockProducers []*BlockProducer

func (bps *BlockProducers) UnmarshalJSON(p []byte) error {
	var response struct {
		BlockProducers []*BlockProducer `json:"next_bps"`
	}
	err := json.Unmarshal(p, &response)
	if err != nil {
		return err
	}

	*bps = BlockProducers(response.BlockProducers)
	return nil
}

func (bps *BlockProducers) Hash() (CryptoHash, error) {
	serializedBps, err := borsh.Serialize(*bps)

	if err != nil {
		return CryptoHash{}, err
	}

	return sha256.Sum256(serializedBps), nil
}


