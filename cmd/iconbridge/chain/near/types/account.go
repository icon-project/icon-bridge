package types

type Account struct {
	Amount        BigInt     `json:"amount"`
	Locked        BigInt     `json:"locked"`
	CodeHash      CryptoHash `json:"code_hash"`
	StorageUsage  int64      `json:"storage_usage"`
	StoragePaidAt int64      `json:"storage_paid_at"`
	Height        int64      `json:"block_height"`
	BlockHash     CryptoHash `json:"block_hash"`
}
