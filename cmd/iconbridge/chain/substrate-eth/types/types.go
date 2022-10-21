package types

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}
