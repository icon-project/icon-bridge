package tezos

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/icon-project/icon-bridge/common/codec"

	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

const (
	txMaxDataSize        = 1024 // 1 KB
	txOverheadScale      = 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale) // with the rlp overhead
	defaultSendTxTimeOut = 30 * time.Second                      // 30 seconds is the block time for tezos
	maxEventPropagation = 5
)

var (
	originalRxSeq = big.NewInt(0)
	statusFlag = false 
)

type senderOptions struct {
	StepLimit        uint64 `json:"step_limit"`
	TxDataSizeLimit  uint64 `json:"tx_data_size_limit"`
	BalanceThreshold uint64 `json:"balance_threshold"`
	BMCManagment     string `json:"bmcManagement"`
}

type sender struct {
	log        log.Logger
	src        chain.BTPAddress
	dst        tezos.Address
	connection *contract.Contract
	parameters micheline.Parameters
	cls        *Client
	blockLevel int64
	opts       senderOptions
	w          wallet.Wallet
}

func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {
	var err error
	// srcAddr := tezos.MustParseAddress(src.ContractAddress())
	dstAddr := tezos.MustParseAddress(dst.ContractAddress())
	s := &sender{
		log: l,
		src: src,
		dst: dstAddr,
		w:   w,
	}

	json.Unmarshal(rawOpts, &s.opts)

	if len(urls) == 0 {
		return nil, fmt.Errorf("Empty url")
	}

	bmcManaement := tezos.MustParseAddress(s.opts.BMCManagment)

	s.cls, err = NewClient(urls[0], dstAddr, bmcManaement, l)
	if err != nil {
		return nil, err
	}

	return s, nil

}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	address := tezos.MustParseAddress(s.w.Address())
	balance, err = s.cls.GetBalance(ctx, s.cls.Cl, address, s.cls.blockLevel)
	if err != nil {
		return nil, nil, err
	}

	return balance, big.NewInt(0), nil
}

func (s *sender) Segment(ctx context.Context, msg *chain.Message) (tx chain.RelayTx, newMsg *chain.Message, err error) {

	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}

	if s.opts.TxDataSizeLimit == 0 {
		limit := defaultTxSizeLimit
		s.opts.TxDataSizeLimit = uint64(limit)
	}

	if len(msg.Receipts) == 0 {
		return nil, msg, nil
	}
	rm := &chain.RelayMessage{
		Receipts: make([][]byte, 0),
	}

	var msgSize uint64

	newMsg = &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts,
	}

	var newEvent []*chain.Event
	var newReceipt *chain.Receipt
	var newReceipts []*chain.Receipt

	for i, receipt := range msg.Receipts {
		fmt.Println("from segment of tezos: ", receipt.Events[0].Message)
		fmt.Println("from segment of tezos: ", receipt.Events[0].Sequence)
		fmt.Println("from segment of tezos: ", receipt.Events[0].Next)	
		fmt.Println("len of events", len(receipt.Events))
		fmt.Println("msg.receipts", len(msg.Receipts))

		if len(receipt.Events) > maxEventPropagation {
			newEvent = receipt.Events[maxEventPropagation:]
			receipt.Events = receipt.Events[:maxEventPropagation]
		}

		rlpEvents, err := codec.RLP.MarshalToBytes(receipt.Events) //json.Marshal(receipt.Events) // change to rlp bytes
		if err != nil {
			return nil, nil, err
		}

		rlpReceipt, err := codec.RLP.MarshalToBytes(&chain.RelayReceipt{
			Index:  receipt.Index,
			Height: receipt.Height,
			Events: rlpEvents,
		}) //json.Marshal(chainReceipt) // change to rlp bytes
		if err != nil {
			return nil, nil, err
		}

		fmt.Println("Message size is initially", msgSize)

		newMsgSize := msgSize + uint64(len(rlpReceipt))
		fmt.Println(newMsgSize)
		if newMsgSize > s.opts.TxDataSizeLimit {
			fmt.Println("limit is", s.opts.TxDataSizeLimit)
			fmt.Println("The value of i is", i)
			newMsg.Receipts = msg.Receipts[i:]
			break
		}
		msgSize = newMsgSize
		fmt.Println("message size", msgSize)
		rm.Receipts = append(rm.Receipts, rlpReceipt)

		if newEvent != nil {
			newReceipt = receipt
			newReceipt.Events = newEvent
			newReceipts = append(newReceipts, newReceipt)
			newReceipts = append(newReceipts, msg.Receipts...)
			msg.Receipts = newReceipts
			break
		}
	}
	message, err := codec.RLP.MarshalToBytes(rm) // json.Marshal(rm)
	if err != nil {
		return nil, nil, err
	}

	tx, err = s.newRelayTx(ctx, msg.From.String(), message)
	if err != nil {
		return nil, nil, err
	}

	return tx, newMsg, nil
}

func (s *sender) Status(ctx context.Context) (link *chain.BMCLinkStatus, err error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	status, err := s.cls.GetStatus(ctx, s.cls.Contract, s.src.String())
	if err != nil {
		return nil, err
	}

	ls := &chain.BMCLinkStatus{}

	ls.TxSeq = status.TxSeq.Uint64()
	ls.RxSeq = status.RxSeq.Uint64()
	ls.CurrentHeight = status.CurrentHeight.Uint64()
	ls.RxHeight = status.RxHeight.Uint64()

	return ls, nil
}

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte) (*relayTx, error) {
	client := s.cls

	return &relayTx{
		Prev:    prev,
		Message: message,
		cl:      client,
		w:       s.w,
		link:	s.src.String(),
	}, nil
}

type relayTx struct {
	Prev    string `json:"_prev"`
	Message []byte `json:"_msg"`

	cl      *Client
	receipt *rpc.Receipt
	w       wallet.Wallet
	link 	string
}

func (tx *relayTx) ID() interface{} {
	return nil
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeOut)
	defer cancel()

	prim := micheline.Prim{}
	messageHex := hex.EncodeToString(tx.Message)

	fmt.Println("starting ma status flag is ", statusFlag)

	status, err := tx.cl.GetStatus(ctx, tx.cl.Contract, tx.link)
	if err != nil {
		return err
	}

	if !statusFlag {
		originalRxSeq = status.RxSeq
		statusFlag = true
	}

	if status.RxSeq.Cmp(originalRxSeq) > 0 {
		statusFlag = false
		hash, err := tx.cl.GetBlockByHeight(ctx, tx.cl.Cl, status.CurrentHeight.Int64())
		if err != nil {
			return err
		}

		tx.receipt = &rpc.Receipt{
			Pos: 0,
			List: 0,
			Block: hash.Hash,
		}
		return nil
	} 

	fmt.Println("status flag", statusFlag)

	fmt.Println(messageHex)

	in := "{ \"prim\": \"Pair\", \"args\": [ { \"bytes\": \"" + messageHex + "\" }, { \"string\": \"" + tx.Prev + "\" } ] }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		return err
	}

	args := contract.NewTxArgs()
	args.WithParameters(micheline.Parameters{Entrypoint: "handle_relay_message", Value: prim})

	opts := rpc.DefaultOptions

	w, ok := tx.w.(*wallet.TezosWallet)
	if !ok {
		return errors.New("not a tezos wallet")
	}

	opts.Signer = w.Signer()
	opts.TTL = 3

	from := tezos.MustParseAddress(tx.w.Address()) // pubk

	argument := args.WithSource(from).WithDestination(tx.cl.Contract.Address())

	receipt, err := tx.cl.HandleRelayMessage(_ctx, argument, &opts)

	if err != nil {
		return err
	}

	tx.receipt = receipt
	statusFlag = false 
	return nil
}

func (tx *relayTx) Receipt(ctx context.Context) (blockHeight uint64, err error) {
	if tx.receipt == nil {
		return 0, fmt.Errorf("couldnot get receipt")
	}

	_, err = tx.cl.GetOperationByHash(ctx, tx.cl.Cl, tx.receipt.Block, tx.receipt.List, tx.receipt.Pos)
	if err != nil {
		return 0, err
	}

	blockHeight, err = tx.cl.GetBlockHeightByHash(ctx, tx.cl.Cl, tx.receipt.Block)
	if err != nil {
		return 0, err
	}
	return blockHeight, nil
}
