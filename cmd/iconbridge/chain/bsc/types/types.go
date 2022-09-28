package types

type Block struct {
	Transactions []string `json:"transactions"`
	GasUsed      string   `json:"gasUsed"`
}

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}
