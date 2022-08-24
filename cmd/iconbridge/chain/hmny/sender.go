//go:build hmny
// +build hmny

package hmny

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"
	"errors"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	txMaxDataSize        = 64 * 1024 // 64 KB
	txOverheadScale      = 0.01      // base64 encoding overhead 0.36, rlp and other fields 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale)
	defaultSendTxTimeout = 15 * time.Second
	defaultGasLimit      = 8e7
	defaultGasPrice      = 3e10
	maxGasPriceBoost     = 10.0
)

func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w wallet.Wallet,
	rawOpts json.RawMessage, l log.Logger) (chain.Sender, error) {
	s := &sender{
		log: l,
		w:   w.(*wallet.EvmWallet),
		src: src,
		dst: dst,
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}
	err := unmarshalOpt(rawOpts, &s.opts)
	if err != nil {
		return nil, err
	}

	if s.opts.BoostGasPrice < 1.0 {
		s.opts.BoostGasPrice = 1.0
	}
	if s.opts.BoostGasPrice > maxGasPriceBoost {
		s.opts.BoostGasPrice = maxGasPriceBoost
	}

	s.cls, s.bmcs, err = newClients(urls, dst.ContractAddress(), s.log)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type senderOptions struct {
	GasLimit         uint64  `json:"gas_limit"`
	BoostGasPrice    float64 `json:"boost_gas_price"`
	TxDataSizeLimit  uint64  `json:"tx_data_size_limit"`
	BalanceThreshold big.Int `json:"balance_threshold"`
}

func (opts *senderOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

func unmarshalOpt(data []byte, opts *senderOptions) error {
	type SenderOptionsTemp struct {
		GasLimit         uint64  `json:"gas_limit"`
		BoostGasPrice    float64 `json:"boost_gas_price"`
		TxDataSizeLimit  uint64  `json:"tx_data_size_limit"`
		BalanceThreshold string `json:"balance_threshold"`
	}
	var senderOptionsObj SenderOptionsTemp

	if err := json.Unmarshal(data, &senderOptionsObj); err != nil {
		return err
	}

	opts.GasLimit = senderOptionsObj.GasLimit
	opts.BoostGasPrice = senderOptionsObj.BoostGasPrice
	opts.TxDataSizeLimit = senderOptionsObj.TxDataSizeLimit

	threshold := new(big.Int)
	valueInt, ok := threshold.SetString(senderOptionsObj.BalanceThreshold, 10)
	if !ok {
		return errors.New("Can't parse field Balance Threshold")
	} else{
		opts.BalanceThreshold = *valueInt
	}

	return nil
}

type sender struct {
	log  log.Logger
	w    *wallet.EvmWallet
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts senderOptions
	cls  []*Client
	bmcs []*BMC
}

func (s *sender) jointClient() (*Client, *BMC) {
	randInt := rand.Intn(len(s.cls))
	return s.cls[randInt], s.bmcs[randInt]
}

// BMCLinkStatus ...
// returns the BMCLinkStatus for "src" link
func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	_, bmcCl := s.jointClient()
	status, err := bmcCl.GetStatus(&bind.CallOpts{Context: ctx}, s.src.String())
	if err != nil {
		return nil, err
	}
	ls := &chain.BMCLinkStatus{}
	ls.TxSeq = status.TxSeq.Uint64()
	ls.RxSeq = status.RxSeq.Uint64()
	// ls.BMRIndex = uint(status.RelayIdx.Uint64())
	// ls.RotateHeight = status.RotateHeight.Uint64()
	// ls.RotateTerm = uint(status.RotateTerm.Uint64())
	// ls.DelayLimit = uint(status.DelayLimit.Uint64())
	// ls.MaxAggregation = uint(status.MaxAggregation.Uint64())
	ls.CurrentHeight = status.CurrentHeight.Uint64()
	ls.RxHeight = status.RxHeight.Uint64()
	// ls.RxHeightSrc = status.RxHeightSrc.Uint64()
	return ls, nil
}

// Segment ...
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

	for i, receipt := range msg.Receipts {
		rlpEvents, err := codec.RLP.MarshalToBytes(receipt.Events)
		if err != nil {
			return nil, nil, err
		}
		rlpReceipt, err := codec.RLP.MarshalToBytes(&chain.RelayReceipt{
			Index:  receipt.Index,
			Height: receipt.Height,
			Events: rlpEvents,
		})
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

	message, err := codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, nil, err
	}

	tx, err = s.newRelayTx(ctx, msg.From.String(), message)
	if err != nil {
		return nil, nil, err
	}

	return tx, newMsg, nil
}

func (s *sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	cl, _ := s.jointClient()
	bal, err := cl.GetBalance(ctx, s.w.Address())
	return bal, &s.opts.BalanceThreshold, err

}

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte) (*relayTx, error) {
	client, bmcClient := s.jointClient()
	chainID, err := client.eth.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(s.w.Skey, chainID)
	if err != nil {
		return nil, err
	}
	txOpts.Context = ctx
	txOpts.GasPrice, _ = (&big.Float{}).Mul(
		(&big.Float{}).SetInt64(defaultGasPrice),
		(&big.Float{}).SetFloat64(s.opts.BoostGasPrice),
	).Int(nil)

	txOpts.GasLimit = defaultGasLimit
	if s.opts.GasLimit > 0 {
		txOpts.GasLimit = s.opts.GasLimit
	}
	return &relayTx{
		Prev:    prev,
		Message: message,
		opts:    txOpts,
		cl:      client,
		bmcCl:   bmcClient,
	}, nil
}

type relayTx struct {
	Prev    string `json:"_prev"`
	Message []byte `json:"_msg"`

	opts      *bind.TransactOpts
	pendingTx *ethtypes.Transaction
	cl        *Client
	bmcCl     *BMC
}

func (tx *relayTx) ID() interface{} {
	if tx.pendingTx != nil {
		return tx.pendingTx.Hash()
	}
	return nil
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	tx.cl.log.WithFields(log.Fields{
		"prev": tx.Prev}).Debug("handleRelayMessage: send tx")

	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()
	txOpts := *tx.opts
	txOpts.Context = _ctx

	nonce, err := tx.cl.eth.NonceAt(ctx, txOpts.From, nil)
	if err != nil {
		return err
	}
	txOpts.Nonce = (&big.Int{}).SetUint64(nonce)
	defer func() {
		if tx.pendingTx != nil {
			txBytes, _ := tx.pendingTx.MarshalJSON()
			tx.cl.log.WithFields(log.Fields{
				"tx": string(txBytes)}).Debug("handleRelayMessage: tx sent")
		}
	}()
	tx.pendingTx, err = tx.bmcCl.HandleRelayMessage(&txOpts, tx.Prev, tx.Message)
	if err != nil {
		tx.cl.log.WithFields(log.Fields{
			"error": err}).Debug("handleRelayMessage: send tx")
		if err.Error() == "insufficient funds for gas * price + value" {
			return chain.ErrInsufficientBalance
		}
		return err
	}

	// tx.cl.log.WithFields(log.Fields{
	// 	"txh": tx.pendingTx.Hash(),
	// 	"msg": btpcommon.HexBytes(tx.Message)}).Debug("handleRelayMessage: tx sent")

	return nil
}

func (tx *relayTx) Receipt(ctx context.Context) (blockHeight uint64, err error) {
	if tx.pendingTx == nil {
		return 0, fmt.Errorf("no pending tx")
	}

	for i, isPending := 0, true; i < 5 && (isPending || err == ethereum.NotFound); i++ {
		time.Sleep(time.Second)
		_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
		defer cancel()
		_, isPending, err = tx.cl.eth.TransactionByHash(_ctx, tx.pendingTx.Hash())
	}
	if err != nil {
		return 0, err
	}

	_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	txr, err := tx.cl.eth.TransactionReceipt(_ctx, tx.pendingTx.Hash())
	if err != nil {
		return 0, err
	}

	if txr.Status == 0 {
		callMsg := ethereum.CallMsg{
			From:       tx.opts.From,
			To:         tx.pendingTx.To(),
			Gas:        tx.pendingTx.Gas(),
			GasPrice:   tx.pendingTx.GasPrice(),
			Value:      tx.pendingTx.Value(),
			AccessList: tx.pendingTx.AccessList(),
			Data:       tx.pendingTx.Data(),
		}

		_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
		defer cancel()
		data, err := tx.cl.eth.CallContract(_ctx, callMsg, txr.BlockNumber)
		if err != nil {
			return 0, err
		}

		if txr.GasUsed >= tx.pendingTx.Gas()*63/64 { // gas limit exceeded
			if txr.GasUsed == txr.CumulativeGasUsed { // block gas limit exceeded
				return 0, chain.ErrBlockGasLimitExceeded
			}
			return 0, chain.ErrGasLimitExceeded
		}

		return 0, chain.RevertError(revertReason(data))
	}

	tx.cl.log.WithFields(log.Fields{
		"txh": tx.pendingTx.Hash()}).Debug("handleRelayMessage: success")

	return txr.BlockNumber.Uint64(), nil
}

func revertReason(data []byte) string {
	if len(data) < 4+32+32 {
		return ""
	}
	data = data[4+32:] // ignore method and index
	length := binary.BigEndian.Uint64(data[24:32])
	return string(data[32 : 32+length])
}
