package tezos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"

	// "blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
)

const (
	txMaxDataSize 			= 32 * 1024 // 8 KB
	defaultSendTxTimeOut 	= 30 * time.Second // 30 seconds is the block time for tezos 
)

type senderOptions struct {
	StepLimit        uint64         `json:"step_limit"`
	TxDataSizeLimit  uint64         `json:"tx_data_size_limit"`
	BalanceThreshold uint64 `json:"balance_threshold"`
}

type sender struct {
	log log.Logger
	src tezos.Address
	dst tezos.Address
	connection *contract.Contract
	parameters micheline.Parameters
	cls *Client
	blockLevel int64
	opts senderOptions
}

func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {
		var err error
		srcAddr := tezos.MustParseAddress(src.String())

		dstAddr := tezos.MustParseAddress(dst.String())

		s := &sender {
			log: l,
			src: srcAddr,
			dst: dstAddr,
		}

		if len(urls) == 0 {
			return nil, fmt.Errorf("Empty url")
		}
		s.cls, err = NewClient(urls[0], srcAddr, l)
		if err != nil {
			return nil, err 
		}

		return s, nil

}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error){
	balance, err = s.cls.GetBalance(ctx, s.cls.Cl, s.src, s.cls.blockLevel)
	if err != nil {
		return nil, nil, err
	}

	return balance, big.NewInt(0), nil 
}

func (s *sender) Segment(ctx context.Context, msg *chain.Message) (tx chain.RelayTx, newMsg *chain.Message, err error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}

	if s.opts.TxDataSizeLimit == 0{
		s.opts.TxDataSizeLimit = uint64(txMaxDataSize)
	}

	if len(msg.Receipts) == 0 {
		return nil, msg, nil
	}

	rm := &chain.RelayMessage{
		Receipts: make([][]byte, 0),
	}

	var msgSize uint64

	newMsg = &chain.Message{
		From: msg.From,
		Receipts: msg.Receipts,
	}

	for i , receipt := range msg.Receipts{
		rlpEvents := receipt.Events

		chainReceipt := &chain.Receipt{
			Index: receipt.Index,
			Height: receipt.Height,
			Events: rlpEvents,
		}

		rlpReceipt, err := json.Marshal(chainReceipt)
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
	message, err := json.Marshal(rm)
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
	return nil, nil
}

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte) (*relayTx, error) {
	client := s.cls

	return &relayTx{
		Prev: prev,
		Message: message,
		cl: client,
	}, nil 
}

type relayTx struct {
	Prev    	string `json:"_prev"`
	Message 	[]byte `json:"_msg"`

	cl        	*Client
	receipt 	*rpc.Receipt
}

func (tx *relayTx) ID() interface{}{
	return nil
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeOut)
	defer cancel()

	prim := micheline.Prim{}

	michJsonStrMsg := "{}" // add message here 

	if err := prim.UnmarshalJSON([]byte(michJsonStrMsg)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return err
	}

	args := contract.NewTxArgs()
	args.WithParameters(micheline.Parameters{Entrypoint: "handleRelayMessage", Value: prim})

	opts := rpc.DefaultOptions

	opts.Signer = signer.NewFromKey(tezos.MustParsePrivateKey("")) // pk 

	from := tezos.MustParseAddress("") // pubk

	argument := args.WithSource(from).WithDestination(tx.cl.Contract.Address())

	receipt, err := tx.cl.HandleRelayMessage(_ctx, argument)

	if err != nil {
		return nil
	}

	tx.receipt = receipt

	return nil
}

func (tx *relayTx) Receipt(ctx context.Context) (blockHeight uint64, err error) {
	return uint64(0), nil
}