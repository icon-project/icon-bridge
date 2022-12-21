package algo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

// TODO review consts
const (
	blockSizeLimit       = 1000000
	defaultSendTxTimeout = 15 * time.Second
	defaultReadTimeout   = 50 * time.Second //
)

func NewSender(
	src, dst chain.BTPAddress,
	algodAccess []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {

	s := &sender{
		log:    l,
		wallet: w.(*wallet.AvmWallet),
		src:    src,
		dst:    dst,
	}
	if len(algodAccess) < 2 {
		return nil, fmt.Errorf("Invalid algorand credentials")
	}

	err := json.Unmarshal(rawOpts, &s.opts)
	if err != nil {
		return nil, err
	}

	s.cl, err = newClient(algodAccess, s.log)
	if err != nil {
		return nil, err
	}

	err = s.initAbi()
	if err != nil {
		return nil, err
	}
	return s, nil
}

type senderOptions struct {
	AppId            uint64         `json:"app_id"`
	BlockSizeLimit   uint64         `json:"tx_data_size_limit"`
	BalanceThreshold intconv.BigInt `json:"balance_threshold"`
}

type sender struct {
	log    log.Logger
	wallet *wallet.AvmWallet
	src    chain.BTPAddress
	dst    chain.BTPAddress
	opts   senderOptions
	cl     *Client
	bmc    *abi.Contract
	mcp    *future.AddMethodCallParams
}

type relayTx struct {
	s    *sender
	txId string
	msg  []byte
}

func (opts *senderOptions) unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	return s.GetBmcStatus(ctx)
}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	bal, err := s.cl.GetBalance(ctx, s.wallet.Address())
	return bal, &s.opts.BalanceThreshold.Int, err
}

func (s *sender) Segment(
	ctx context.Context, msg *chain.Message,
) (tx chain.RelayTx, newMsg *chain.Message, err error) {
	if ctx.Err() != nil {
		return nil, msg, ctx.Err()
	}
	if s.opts.BlockSizeLimit == 0 {
		limit := blockSizeLimit
		s.opts.BlockSizeLimit = uint64(limit)
	}
	if len(msg.Receipts) == 0 {
		return nil, msg, nil
	}

	tmpReceipts := make([][]byte, 0)
	var msgSize uint64

	newMsg = &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts,
	}
	for i, receipt := range msg.Receipts {
		var encReceipt = make([]byte, 1e3)
		base64.StdEncoding.Encode(encReceipt, msgpack.Encode(receipt))

		msgSize += uint64(len(encReceipt))
		if msgSize > s.opts.BlockSizeLimit {
			newMsg.Receipts = msg.Receipts[i:]
			break
		}

		tmpReceipts = append(tmpReceipts, encReceipt)
	}
	rm := &chain.RelayMessage{
		Receipts: tmpReceipts,
	}

	var encMessage = make([]byte, 10000)
	base64.StdEncoding.Encode(encMessage, msgpack.Encode(rm))

	newTx := &relayTx{
		s:    s,
		txId: "",
		msg:  encMessage,
	}

	return newTx, newMsg, nil
}

func (tx relayTx) Send(ctx context.Context) (err error) {
	tx.s.cl.Log().WithFields(log.Fields{
		"prev": tx.s.wallet}).Debug("handleRelayMessage: send tx")
	ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()
	tx.txId, err = tx.s.HandleRelayMessage(ctx, []byte(tx.s.wallet.Address()), tx.msg)
	if err != nil {
		tx.s.cl.Log().WithFields(log.Fields{
			"error": err}).Debug("handleRelayMessage: send tx")
		return err
	}
	return nil
}

// Waits for txn to be confirmed and gets its receipt
func (tx relayTx) Receipt(ctx context.Context) (blockNumber uint64, err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	if tx.msg == nil {
		return 0, fmt.Errorf("Can't get receipt from tx: Empty relay message")
	}
	confirmedTxn, err := tx.s.cl.WaitForTransaction(ctx, tx.txId)
	if err != nil {
		return 0, fmt.Errorf("Can't get receipt from tx: %w", err)

	} else {
		return confirmedTxn.ConfirmedRound, nil
	}
}

func (tx relayTx) ID() interface{} {
	return tx.txId
}
