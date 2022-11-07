package types

import (
	"math/big"

	"github.com/near/borsh-go"
)

type Header struct {
	Height                int64              `json:"height"`
	PreviousHeight        int64              `json:"prev_height"`
	EpochId               CryptoHash         `json:"epoch_id"`
	NextEpochId           CryptoHash         `json:"next_epoch_id"`
	Hash                  CryptoHash         `json:"hash"`
	PreviousBlockHash     CryptoHash         `json:"prev_hash"`
	PreviousStateRoot     CryptoHash         `json:"prev_state_root"`
	ChunkReceiptsRoot     CryptoHash         `json:"chunk_receipts_root"`
	ChunkHeadersRoot      CryptoHash         `json:"chunk_headers_root"`
	ChunkTransactionRoot  CryptoHash         `json:"chunk_tx_root"`
	OutcomeRoot           CryptoHash         `json:"outcome_root"`
	ChunksIncluded        uint8              `json:"chunks_included"`
	ChallengesRoot        CryptoHash         `json:"challenges_root"`
	Timestamp             Timestamp          `json:"timestamp_nanosec"`
	RandomValue           CryptoHash         `json:"random_value"`
	ValidatorProposals    []BlockProducer    `json:"validator_proposals"`
	ChunkMask             []bool             `json:"chunk_mask"`
	GasPrice              BigInt             `json:"gas_price"`
	BlockOrdinal          uint64             `json:"block_ordinal"`
	TotalSupply           BigInt             `json:"total_supply"`
	ChallengesResult      []SlashedValidator `json:"challenges_result"`
	LastFinalBlock        CryptoHash         `json:"last_final_block"`
	LastDSFinalBlock      CryptoHash         `json:"last_ds_final_block"`
	NextBlockProducerHash CryptoHash         `json:"next_bp_hash"`
	BlockMerkleRoot       CryptoHash         `json:"block_merkle_root"`
	EpochSyncDataHash     *CryptoHash        `json:"epoch_sync_data_hash"`
	Approvals             []*Signature       `json:"approvals"`
	Signature             Signature          `json:"signature"`
	LatestProtocolVersion uint32             `json:"latest_protocol_version"`
}

type HeaderInnerLite struct {
	Height                uint64
	EpochId               [32]byte
	NextEpochId           [32]byte
	PreviousStateRoot     [32]byte
	OutcomeRoot           [32]byte
	Timestamp             uint64
	NextBlockProducerHash [32]byte
	BlockMerkleRoot       [32]byte
}

func (h HeaderInnerLite) BorshSerialize() ([]byte, error) {
	return borsh.Serialize(h)
}

type HeaderInnerRest struct {
	ChunkReceiptsRoot     [32]byte
	ChunkHeadersRoot      [32]byte
	ChunkTransactionRoot  [32]byte
	ChallengesRoot        [32]byte
	RandomValue           [32]byte
	ValidatorProposals    []BlockProducer
	ChunkMask             []bool
	GasPrice              BigInt
	TotalSupply           BigInt
	ChallengesResult      []SlashedValidator
	LastFinalBlock        [32]byte
	LastDSFinalBlock      [32]byte
	BlockOrdinal          uint64
	PreviousHeight        uint64
	EpochSyncDataHash     *CryptoHash
	Approvals             []*Signature
	LatestProtocolVersion uint32
}

func (h HeaderInnerRest) BorshSerialize() ([]byte, error) {
	return borsh.Serialize(struct {
		ChunkReceiptsRoot     CryptoHash
		ChunkHeadersRoot      CryptoHash
		ChunkTransactionRoot  CryptoHash
		ChallengesRoot        CryptoHash
		RandomValue           CryptoHash
		ValidatorProposals    []BlockProducer
		ChunkMask             []bool
		GasPrice              big.Int
		TotalSupply           big.Int
		ChallengesResult      []SlashedValidator
		LastFinalBlock        CryptoHash
		LastDSFinalBlock      CryptoHash
		BlockOrdinal          uint64
		PreviousHeight        uint64
		EpochSyncDataHash     *CryptoHash
		Approvals             []*Signature
		LatestProtocolVersion uint32
	}{
		ChunkReceiptsRoot:     h.ChunkReceiptsRoot,
		ChunkHeadersRoot:      h.ChunkHeadersRoot,
		ChunkTransactionRoot:  h.ChunkTransactionRoot,
		ChallengesRoot:        h.ChallengesRoot,
		RandomValue:           h.RandomValue,
		ValidatorProposals:    h.ValidatorProposals,
		ChunkMask:             h.ChunkMask,
		GasPrice:              big.Int(h.GasPrice),
		TotalSupply:           big.Int(h.TotalSupply),
		ChallengesResult:      h.ChallengesResult,
		LastFinalBlock:        h.LastFinalBlock,
		LastDSFinalBlock:      h.LastDSFinalBlock,
		BlockOrdinal:          h.BlockOrdinal,
		PreviousHeight:        h.PreviousHeight,
		EpochSyncDataHash:     h.EpochSyncDataHash,
		Approvals:             h.Approvals,
		LatestProtocolVersion: h.LatestProtocolVersion,
	})
}
