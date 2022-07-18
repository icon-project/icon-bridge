package types

type BmcStatus struct {
	TxSeq            uint64        `json:"tx_seq"`
	RxSeq            uint64        `json:"rx_seq"`
	Verifier         AccountId     `json:"verifier"`
	BMRs             []RelayStatus `json:"relays"`
	BMRIndex         uint          `json:"relay_index"`
	RotateHeight     uint64        `json:"rotate_height"`
	RotateTerm       uint          `json:"rotate_term"`
	DelayLimit       uint          `json:"delay_limit"`
	MaxAggregation   uint          `json:"max_aggregation"`
	CurrentHeight    uint64        `json:"current_height"`
	RxHeight         uint64        `json:"rx_height"`
	RxHeightSrc      uint64        `json:"rx_height_src"`
	BlockIntervalSrc uint          `json:"block_interval_src"`
	BlockIntervalDst uint          `json:"block_interval_dst"`
}
