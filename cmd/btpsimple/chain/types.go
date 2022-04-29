package chain

import "context"

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
	Send(ctx context.Context) (err error)
	Receipt(ctx context.Context) (txr interface{}, err error)
}

type Receiver interface {
	// SubscribeMessage ...
	// subscribes to BTP messages and block headers of the src chain
	SubscribeMessage(ctx context.Context, height, seq uint64) (msgCh <-chan *Message, err error)
}

type Sender interface {
	// Status ...
	// returns current BMCLinkStatus of the dst chain
	Status(ctx context.Context) (link *BMCLinkStatus, err error)

	// Segment ...
	// returns a "tx" Tx object including events upto "txSizeLimit" bytes and
	// returns rest of the "msg" Message as "newMsg"
	Segment(ctx context.Context, msg *Message, txSizeLimit uint64) (tx RelayTx, newMsg *Message, err error)
}

type Relayer interface {
	Sender
	Receiver
}
