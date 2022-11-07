package types

type VerifierConfig struct {
	PreviousBlockHeight uint64           `json:"previousBlockHeight"`
	PreviousBlockHash   CryptoHash       `json:"previousBlockHash"`
	NextEpochId         CryptoHash       `json:"nextEpoch"`
	BlockProducers      []*BlockProducer `json:"blockProducers"`
	NextBpHash          CryptoHash       `json:"nextBpHash"`
}

type ReceiverOptions struct {
	SyncConcurrency uint           `json:"syncConcurrency"`
	Verifier        VerifierConfig `json:"verifier"`
}
