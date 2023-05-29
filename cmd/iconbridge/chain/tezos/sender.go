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

	// "github.com/icon-project/icon-bridge/common/wallet"

	"github.com/icon-project/icon-bridge/common/codec"

	// "blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

const (
	txMaxDataSize        = 32 * 1024 * 4 // 8 KB
	txOverheadScale      = 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale)
	defaultSendTxTimeOut = 30 * time.Second // 30 seconds is the block time for tezos
)

type senderOptions struct {
	StepLimit        uint64 `json:"step_limit"`
	TxDataSizeLimit  uint64 `json:"tx_data_size_limit"`
	BalanceThreshold uint64 `json:"balance_threshold"`
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
	fmt.Println(src.ContractAddress())
	fmt.Println(dst.ContractAddress())
	// srcAddr := tezos.MustParseAddress(src.ContractAddress())
	dstAddr := tezos.MustParseAddress(dst.ContractAddress())
	s := &sender{
		log: l,
		src: src,
		dst: dstAddr,
		w:   w,
	}
	PrintPlus()
	fmt.Println(w.Address())
	if len(urls) == 0 {
		return nil, fmt.Errorf("Empty url")
	}
	s.cls, err = NewClient(urls[0], dstAddr, l)
	if err != nil {
		return nil, err
	}

	return s, nil

}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	fmt.Println("reached in balance of tezos")
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
	fmt.Println("reached upto here")
	if len(msg.Receipts) == 0 {
		fmt.Println("Probably gone from here")
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
	
	for i, receipt := range msg.Receipts {
		rlpEvents, err := codec.RLP.MarshalToBytes(receipt.Events) //json.Marshal(receipt.Events) // change to rlp bytes
		if err != nil {
			return nil, nil, err
		}

		Print()

		fmt.Println(receipt.Index)
		fmt.Println(receipt.Height)
		fmt.Println(receipt.Events)

		rlpReceipt, err := codec.RLP.MarshalToBytes(&chain.RelayReceipt{
			Index:  receipt.Index,
			Height: receipt.Height,
			Events: rlpEvents,
		}) //json.Marshal(chainReceipt) // change to rlp bytes
		if err != nil {
			return nil, nil, err
		}

		newMsgSize := msgSize + uint64(len(rlpReceipt))
		if newMsgSize > s.opts.TxDataSizeLimit {
			newMsg.Receipts = msg.Receipts[i:]
			break
		}
		msgSize = newMsgSize
		rm.Receipts = append(rm.Receipts, rlpReceipt)
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

	fmt.Println("reached in new relaytx")
	client := s.cls

	return &relayTx{
		Prev:    prev,
		Message: message,
		cl:      client,
		w:       s.w,
	}, nil
}

type relayTx struct {
	Prev    string `json:"_prev"`
	Message []byte `json:"_msg"`

	cl      *Client
	receipt *rpc.Receipt
	w       wallet.Wallet
}

func (tx *relayTx) ID() interface{} {
	return nil
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	fmt.Println("reached in sender of tezos")
	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeOut)
	defer cancel()

	prim := micheline.Prim{}
	messageHex := hex.EncodeToString(tx.Message)

	fmt.Println("Previous is: ", tx.Prev)

	in := "{ \"prim\": \"Pair\", \"args\": [ { \"bytes\": \"" + messageHex + "\" }, { \"string\": \"" + tx.Prev + "\" } ] }"
	fmt.Println(in)
	
	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
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
	fmt.Println("memory signer is : ", opts.Signer)
	opts.TTL = 3

	from := tezos.MustParseAddress(tx.w.Address()) // pubk

	argument := args.WithSource(from).WithDestination(tx.cl.Contract.Address())

	fmt.Println("The message is", messageHex)
	receipt, err := tx.cl.HandleRelayMessage(_ctx, argument, &opts)

	if err != nil {
		return nil
	}

	tx.receipt = receipt

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

func Print() {
	for i := 0; i < 100; i++ {
		fmt.Println("*")
	}
}

func PrintPlus() {
	for i := 0; i < 100; i++ {
		fmt.Println("__")
	}
}
