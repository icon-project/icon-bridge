package types

type CallFunction struct {
	RequestType  string    `json:"request_type"`
	Finality     string    `json:"finality"`
	AccountId    AccountId `json:"account_id"`
	MethodName   string    `json:"method_name"`
	ArgumentsB64 string    `json:"args_base64"`
}

type CallFunctionResponse struct {
	Result      []byte `json:"result"`
	BlockHeight int64  `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}