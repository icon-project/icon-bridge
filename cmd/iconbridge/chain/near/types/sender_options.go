package types

type SenderOptions struct {
	TxDataSizeLimit  uint64 `json:"tx_data_size_limit"`
	BalanceThreshold BigInt `json:"balance_threshold"`
}
