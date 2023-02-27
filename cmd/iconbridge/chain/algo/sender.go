package algo

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/future"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

// TODO review consts
const (
	defaultSendTxTimeout = 15 * time.Second
	defaultReadTimeout   = 50 * time.Second
	atomicTxnLimit       = 16
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
	s      *sender
	round  uint64
	svcs   []AbiFunc
	txIDs  []string
	height uint64
}

type bmcLink struct {
	TxSeq    uint64 `json:"tx_seq"`
	RxSeq    uint64 `json:"rx_seq"`
	RxHeight uint64 `json:"rx_height"`
	TxHeight uint64 `json:"tx_height"`
}

func (opts *senderOptions) unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	return getStatus()
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

	if len(msg.Receipts) == 0 {
		return nil, msg, nil
	}

	newMsg = &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts,
	}

	abiFuncs := make([]AbiFunc, atomicTxnLimit)

	// egment messages to fit the 16 atc limit and process all events in the same abi call
	for i, receipt := range msg.Receipts {
		if len(abiFuncs)+len(receipt.Events) >= cap(abiFuncs) {
			newMsg.Receipts = msg.Receipts[i:]
			break
		}
		for _, event := range receipt.Events {
			svcName, svcArgs, err := DecodeRelayMessage(hex.EncodeToString(event.Message))
			if err != nil {
				return nil, nil, fmt.Errorf("Error decoding event message: %w", err)
			}
			abiFuncs = append(abiFuncs, AbiFunc{svcName, []interface{}{svcArgs}})
		}
	}
	newTx := &relayTx{
		s:      s,
		svcs:   abiFuncs,
		height: msg.Receipts[0].Height,
	}
	return newTx, newMsg, nil
}

func (tx relayTx) Send(ctx context.Context) (err error) {
	tx.s.cl.Log().WithFields(log.Fields{"prev": tx.s.wallet}).Debug("handleRelayMessage: send tx")

	ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()

	res, err := tx.s.callAbi(ctx, tx.svcs...)
	if err != nil {
		return fmt.Errorf("Error calling abi to execute relay txn: %w", err)
	}
	tx.round = res.ConfirmedRound
	tx.txIDs = res.TxIDs
	return nil
}

// Increment sequeence number when a new message gets to the Algorand BMC
func (tx relayTx) Receipt(ctx context.Context) (blockNumber uint64, err error) {
	err = incrementSeq("rx_seq")
	if err != nil {
		return 0, err
	}
	err = updateHeight("rx_height", tx.height)
	if err != nil {
		return 0, err
	}
	return tx.round, nil
}

func (tx relayTx) ID() interface{} {
	return tx.txIDs
}
