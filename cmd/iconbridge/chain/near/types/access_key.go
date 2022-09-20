package types

type AccessKeyResponse struct {
	Nonce       int64  `json:"nonce"`
	Permission  string `json:"permission"`
	BlockHeight int64  `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	Error       string `json:"error"`
}
