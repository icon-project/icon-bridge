package hmny

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/btp/cmd/btpsimple/chain"
	btpcommon "github.com/icon-project/btp/common"
	"github.com/icon-project/btp/common/codec"
	"github.com/icon-project/btp/common/log"
	"github.com/icon-project/btp/common/wallet"
)

const (
	txMaxDataSize        = 64 * 1024 // 64 KB
	txOverheadScale      = 0.01      // base64 encoding overhead 0.36, rlp and other fields 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale)
	defaultGasLimit      = 10000000
	defaultSendTxTimeout = 15 * time.Second
)

func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w *wallet.EvmWallet,
	opts map[string]interface{}, l log.Logger) (chain.Sender, error) {
	s := &sender{
		log: l,
		w:   w,
		src: src,
		dst: dst,
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}
	err := s.opts.Unmarshal(opts)
	if err != nil {
		return nil, err
	}
	s.cls, err = newClients(urls, dst.ContractAddress(), s.log)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type senderOptions struct {
	GasLimit uint64 `json:"gas_limit"`
}

func (opts *senderOptions) Unmarshal(v map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, opts)
}

type sender struct {
	log  log.Logger
	w    *wallet.EvmWallet
	src  chain.BTPAddress
	dst  chain.BTPAddress
	opts senderOptions
	cls  []*client
}

func (s *sender) client() *client {
	return s.cls[rand.Intn(len(s.cls))]
}

// BMCLinkStatus ...
// returns the BMCLinkStatus for "src" link
func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	status, err := s.client().bmc.GetStatus(&bind.CallOpts{Context: ctx}, s.src.String())
	if err != nil {
		return nil, err
	}
	ls := &chain.BMCLinkStatus{}
	ls.TxSeq = status.TxSeq.Uint64()
	ls.RxSeq = status.RxSeq.Uint64()
	ls.BMRIndex = uint(status.RelayIdx.Uint64())
	ls.RotateHeight = status.RotateHeight.Uint64()
	ls.RotateTerm = uint(status.RotateTerm.Uint64())
	ls.DelayLimit = uint(status.DelayLimit.Uint64())
	ls.MaxAggregation = uint(status.MaxAggregation.Uint64())
	ls.CurrentHeight = status.CurrentHeight.Uint64()
	ls.RxHeight = status.RxHeight.Uint64()
	ls.RxHeightSrc = status.RxHeightSrc.Uint64()
	return ls, nil
}

// Segment ...
func (s *sender) Segment(
	ctx context.Context,
	msg *chain.Message, txSizeLimit uint64,
) (tx chain.RelayTx, newMsg *chain.Message, err error) {
	if ctx.Err() != nil {
		return nil, msg, ctx.Err()
	}

	if txSizeLimit == 0 {
		limit := defaultTxSizeLimit
		txSizeLimit = uint64(limit)
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
		if newMsgSize > txSizeLimit {
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

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte) (*relayTx, error) {
	client := s.client()
	chainID, err := client.eth.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(s.w.Skey, chainID)
	if err != nil {
		return nil, err
	}
	txOpts.Context = ctx
	txOpts.GasPrice, err = client.eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	txOpts.GasLimit = defaultGasLimit
	if s.opts.GasLimit > 0 {
		txOpts.GasLimit = s.opts.GasLimit
	}
	return &relayTx{
		Prev:    prev,
		Message: message, // base64.URLEncoding.EncodeToString(rlpCrm),
		opts:    txOpts,
		cl:      client,
	}, nil
}

type relayTx struct {
	Prev    string `json:"_prev"`
	Message []byte `json:"_msg"`

	opts      *bind.TransactOpts
	pendingTx *ethtypes.Transaction
	cl        *client
}

func (tx *relayTx) Send(ctx context.Context) (err error) {
	tx.cl.log.WithFields(log.Fields{"prev": tx.Prev}).Debug("handleRelayMessage: sending tx")
	tx.cl.log.WithFields(log.Fields{"msg": btpcommon.HexBytes(tx.Message)}).Debug("handleRelayMessage: sending tx")

	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()
	txOpts := *tx.opts
	txOpts.Context = _ctx

	tx.pendingTx, err = tx.cl.bmc.HandleRelayMessage(&txOpts, tx.Prev, tx.Message)
	if err != nil {
		tx.cl.log.WithFields(log.Fields{"error": err}).Debug("handleRelayMessage: failed to send tx")
		return err
	}

	tx.cl.log.WithFields(log.Fields{"txh": tx.pendingTx.Hash()}).Debug("handleRelayMessage: tx sent")
	return nil
}

func (tx *relayTx) Receipt(ctx context.Context) (receipt interface{}, err error) {
	if tx.pendingTx == nil {
		return nil, fmt.Errorf("no pending tx")
	}

	for isPending := true; isPending; {
		_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
		defer cancel()
		_, isPending, err = tx.cl.eth.TransactionByHash(_ctx, tx.pendingTx.Hash())
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Second)
	}

	_ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	txr, err := tx.cl.eth.TransactionReceipt(_ctx, tx.pendingTx.Hash())
	if err != nil {
		return nil, err
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
			return nil, err
		}

		return nil, chain.RevertError(revertReason(data))
	}

	return txr, nil
}

func revertReason(data []byte) string {
	if len(data) < 4+32+32 {
		return ""
	}
	data = data[4+32:] // ignore method and index
	length := binary.BigEndian.Uint64(data[24:32])
	return string(data[32 : 32+length])
}
