package types

type VerifierConfig struct {
	BlockHeight       uint64     `json:"block_height"`
	PreviousBlockHash CryptoHash `json:"previous_block_hash"`
	CurrentEpochId    CryptoHash `json:"current_epoch_id"`
	NextEpochId       CryptoHash `json:"next_epoch_id"`
	NextBpsHash       CryptoHash `json:"next_bps_hash"`
	CurrentBpsHash    CryptoHash `json:"current_bps_hash"`
}

type ReceiverOptions struct {
	SyncConcurrency int            `json:"sync_concurrency"`
	Verifier        *VerifierConfig `json:"verifier"`
}
