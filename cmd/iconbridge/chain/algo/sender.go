package algo

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/common/codec"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

// TODO review consts
const (
	txMaxDataSize        = 8 * 1024 // 8 KB
	txOverheadScale      = 0.01     // base64 encoding overhead 0.36, rlp and other fields 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale)
	defaultSendTxTimeout = 15 * time.Second
	defaultGasPrice      = 18000000000
	maxGasPriceBoost     = 10.0
	defaultReadTimeout   = 50 * time.Second //
	DefaultGasLimit      = 25000000
)

func NewSender(
	src, dst chain.BTPAddress,
	algodAccess []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {
	s := &sender{
		log:    l,
		wallet: w,
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
	return s, nil
}

type senderOptions struct {
	GasLimit         uint64         `json:"gas_limit"`
	BoostGasPrice    float64        `json:"boost_gas_price"`
	TxDataSizeLimit  uint64         `json:"tx_data_size_limit"`
	BalanceThreshold intconv.BigInt `json:"balance_threshold"`
}

type sender struct {
	log    log.Logger
	wallet wallet.Wallet
	src    chain.BTPAddress
	dst    chain.BTPAddress
	opts   senderOptions
	cl     *Client
}

// TODO review relayTx and all the methods using it
type relayTx struct {
	wallet *wallet.AvmWallet
	txn    types.Transaction
	txId   string
	cl     IClient
}

func (opts *senderOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	return s.cl.GetBmcStatus(ctx)
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
	for i, _ := range msg.Receipts {

		//TODO create decoding methad similar to RLP to convert relay receipts into a byte array
		var rlpReceipt []byte

		newMsgSize := msgSize + uint64(len(rlpReceipt))
		if newMsgSize > s.opts.TxDataSizeLimit {
			newMsg.Receipts = msg.Receipts[i:]
			break
		}
		msgSize = newMsgSize
		rm.Receipts = append(rm.Receipts, rlpReceipt)
	}
	_, err = codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, nil, err
	}

	newTx := &relayTx{
		wallet: (s.wallet).(*wallet.AvmWallet),
		cl:     s.cl,
	}

	return newTx, newMsg, nil
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	tx.cl.Log().WithFields(log.Fields{
		"prev": tx.wallet.Pkey}).Debug("handleRelayMessage: send tx")

	ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer func() {
		cancel()
		if !tx.txn.Empty() {
			txBytes := tx.txn.Note
			tx.cl.Log().WithFields(log.Fields{
				"tx": string(txBytes)}).Debug("handleRelayMessage: tx sent")
		}
	}()

	/* this func should make bmc call to execute its HRM method and get new tx and txId */
	_, tx.txId, err = tx.cl.HandleRelayMessage(ctx, []byte(tx.wallet.Address()), tx.txn.Header.Note, tx.txn.Header.Sender)
	if err != nil {
		tx.cl.Log().WithFields(log.Fields{
			"error": err}).Debug("handleRelayMessage: send tx")
		if err.Error() == "insufficient funds for gas * price + value" {
			return chain.ErrInsufficientBalance
		}
		return err
	}
	return nil
}

// Waits for txn to be confirmed and gets its receipt
func (tx *relayTx) Receipt(ctx context.Context) (blockNumber uint64, err error) {

	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	if tx.txn.Empty() {
		return 0, fmt.Errorf("no pending tx")
	}

	confirmedTxn, err := tx.cl.WaitForTransaction(ctx, tx.txId)

	if err != nil {
		return 0, err
	} else {
		return confirmedTxn.ConfirmedRound, nil
	}
}

func (tx *relayTx) ID() interface{} {
	if !tx.txn.Empty() {
		//TODO check if lease is the same as transaction id
		return tx.txn.Lease
	} else {
		return nil
	}
}
