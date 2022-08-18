package chain

import (
	"context"
	"math/big"
)

// RelayMessage is encoded
type RelayMessage struct {
	Receipts [][]byte
}

type RelayReceipt struct {
	Index  uint64
	Events []byte
	Height uint64
}

type Event struct {
	Next     BTPAddress
	Sequence uint64
	Message  []byte
}

type Receipt struct {
	Index  uint64
	Events []*Event
	Height uint64
}

type Message struct {
	From     BTPAddress
	Receipts []*Receipt
	// Headers  []interface{}
}

type BMCLinkStatus struct {
	TxSeq            uint64
	RxSeq            uint64
	BMRIndex         uint
	RotateHeight     uint64
	RotateTerm       uint
	DelayLimit       uint
	MaxAggregation   uint
	CurrentHeight    uint64
	RxHeight         uint64
	RxHeightSrc      uint64
	BlockIntervalSrc uint
	BlockIntervalDst uint
}

// RelayTx ...
type RelayTx interface {
	ID() interface{}
	Send(ctx context.Context) (err error)
	Receipt(ctx context.Context) (blockHeight uint64, err error)
}

type SubscribeOptions struct {
	Seq    uint64
	Height uint64
}

type Receiver interface {
	// Subscribe ...
	// subscribes to BTP messages and block headers on `msgCh` of the src chain
	// and returns an `errCh` that sends any error during subscription and terminates
	// the subscription by closing `errCh`
	Subscribe(ctx context.Context, msgCh chan<- *Message, opts SubscribeOptions) (errCh <-chan error, err error)
}

type Sender interface {
	// Status ...
	// returns current BMCLinkStatus of the dst chain
	Status(ctx context.Context) (link *BMCLinkStatus, err error)

	// Segment ...
	// returns a "tx" Tx object including events upto "txSizeLimit" bytes and
	// returns rest of the "msg" Message as "newMsg"
	Segment(ctx context.Context, msg *Message) (tx RelayTx, newMsg *Message, err error)

	// Returns the current relayer balance
	Balance(ctx context.Context) (balance, threshold *big.Int, err error)
}

type Relayer interface {
	Sender
	Receiver
}
