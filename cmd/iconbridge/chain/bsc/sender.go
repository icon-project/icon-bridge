/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bsc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

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

/*
   type sender struct {
	   c   *Client
	   src module.BtpAddress
	   dst module.BtpAddress
	   w   Wallet
	   l   log.Logger
	   opt struct {
	   }

	   bmc *binding.BMC

	   evtLogRawFilter struct {
		   addr      []byte
		   signature []byte
		   next      []byte
		   seq       []byte
	   }
	   evtReq             *BlockRequest
	   isFoundOffsetBySeq bool
	   cb                 module.ReceiveCallback

	   mutex sync.Mutex
   }
*/

type senderOptions struct {
	GasLimit         uint64  `json:"gas_limit"`
	TxDataSizeLimit  uint64  `json:"tx_data_size_limit"`
	BoostGasPrice    float64 `json:"boost_gas_price"`
	BalanceThreshold big.Int `json:"balance_threshold"`
}

type sender struct {
	log          log.Logger
	w            *wallet.EvmWallet
	src          chain.BTPAddress
	dst          chain.BTPAddress
	opts         senderOptions
	cls          []*Client
	bmcs         []*BMC
	prevGasPrice *big.Int
}

func (s *sender) jointClient() (*Client, *BMC) {
	randInt := rand.Intn(len(s.cls))
	return s.cls[randInt], s.bmcs[randInt]
}

func NewSender(
	src, dst chain.BTPAddress,
	urls []string, w wallet.Wallet,
	opts map[string]interface{}, l log.Logger) (chain.Sender, error) {
	s := &sender{
		log:          l,
		w:            w.(*wallet.EvmWallet),
		src:          src,
		dst:          dst,
		prevGasPrice: big.NewInt(defaultGasPrice),
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	b, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal opt:%#v err:%+v", opts, err)
	}
	if err = json.Unmarshal(b, &s.opts); err != nil {
		return nil, fmt.Errorf("fail to unmarshal opt:%#v err:%+v", opts, err)
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

// BMCLinkStatus ...
// returns the BMCLinkStatus for "src" link
func (s *sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	_, bmcCl := s.jointClient()
	status, err := bmcCl.GetStatus(&bind.CallOpts{Context: ctx}, s.src.String())
	if err != nil {
		s.log.Error("GetStatus", "err", err)
		return nil, err
	}
	ls := &chain.BMCLinkStatus{}
	ls.TxSeq = status.TxSeq.Uint64()
	ls.RxSeq = status.RxSeq.Uint64()
	ls.RxHeight = status.RxHeight.Uint64()
	ls.CurrentHeight = status.CurrentHeight.Uint64()
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
	cl, _ := s.jointClient()
	gasPrice, gasHeight, err := cl.GetMedianGasPriceForBlock()
	if err != nil || gasPrice.Int64() == 0 {
		s.log.Infof("GetMedianGasPriceForBlock(%v) Msg: %v. Using default value for gas price \n", gasHeight.String(), err)
		gasPrice = s.prevGasPrice
	} else {
		s.prevGasPrice = gasPrice
		s.log.Infof("GetMedianGasPriceForBlock(%v) price: %v \n", gasHeight.String(), gasPrice.String())
	}
	boostedGasPrice, _ := (&big.Float{}).Mul(
		(&big.Float{}).SetInt64(gasPrice.Int64()),
		(&big.Float{}).SetFloat64(s.opts.BoostGasPrice),
	).Int(nil)
	tx, err = s.newRelayTx(ctx, msg.From.String(), message, boostedGasPrice)
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

func (s *sender) newRelayTx(ctx context.Context, prev string, message []byte, gasPrice *big.Int) (*relayTx, error) {
	client, bmcClient := s.jointClient()
	txOpts, err := client.newTransactOpts(s.w)
	if err != nil {
		return nil, err
	}
	txOpts.Context = ctx
	if s.opts.GasLimit > 0 {
		txOpts.GasLimit = s.opts.GasLimit
	}
	txOpts.GasPrice = gasPrice
	return &relayTx{
		Prev:    prev,
		Message: message, // base64.URLEncoding.EncodeToString(rlpCrm),
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

func (tx *relayTx) Receipt(ctx context.Context) (blockNumber uint64, err error) {
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

/*
   func (s *sender) newTransactionParam(prev string, rm *RelayMessage) (*TransactionParam, error) {
	   b, err := codec.RLP.MarshalToBytes(rm)
	   if err != nil {
		   return nil, err
	   }
	   rmp := BMCRelayMethodParams{
		   Prev: prev,
		   //Messages: base64.URLEncoding.EncodeToString(b[:]),
		   Messages: string(b[:]),
	   }
	   s.l.Debugf("HandleRelayMessage msg: %s", base64.URLEncoding.EncodeToString(b))
	   p := &TransactionParam{
		   Params: rmp,
	   }
	   return p, nil
   }

   func (s *sender) Relay(segment *module.Segment) (module.GetResultParam, error) {
	   s.mutex.Lock()
	   defer s.mutex.Unlock()
	   p, ok := segment.TransactionParam.(*TransactionParam)
	   if !ok {
		   return nil, fmt.Errorf("casting failure")
	   }
	   t, err := s.c.newTransactOpts(s.w)
	   if err != nil {
		   return nil, err
	   }
	   rmp := p.Params.(BMCRelayMethodParams)
	   var tx *types.Transaction
	   tx, err = s.bmc.HandleRelayMessage(t, rmp.Prev, rmp.Messages)
	   if err != nil {
		   s.l.Errorf("handleRelayMessage: ", err.Error())
		   return nil, err
	   }
	   thp := &TransactionHashParam{}
	   thp.Hash = tx.Hash()
	   s.l.Debugf("HandleRelayMessage tx hash:%s, prev %s, msg: %s", thp.Hash, rmp.Prev, base64.URLEncoding.EncodeToString([]byte(rmp.Messages)))
	   return thp, nil
   }

   func (s *sender) GetResult(p module.GetResultParam) (module.TransactionResult, error) {
	   if txh, ok := p.(*TransactionHashParam); ok {
		   for {
			   _, pending, err := s.c.GetTransaction(txh.Hash)
			   if err != nil {
				   return nil, err
			   }
			   if pending {
				   <-time.After(DefaultGetRelayResultInterval)
				   continue
			   }
			   tx, err := s.c.GetTransactionReceipt(txh.Hash)
			   if err != nil {
				   return nil, err
			   }
			   return tx, nil //mapErrorWithTransactionResult(&types.Receipt{}, err) // TODO: map transaction.js result error
		   }
	   } else {
		   return nil, fmt.Errorf("fail to casting TransactionHashParam %T", p)
	   }
   }

   func (s *sender) GetStatus() (*module.BMCLinkStatus, error) {
	   var status binding.TypesLinkStats
	   status, err := s.bmc.GetStatus(nil, s.src.String())

	   if err != nil {
		   s.l.Errorf("Error retrieving relay status from BMC")
		   return nil, err
	   }

	   ls := &module.BMCLinkStatus{}
	   ls.TxSeq = status.TxSeq.Int64()
	   ls.RxSeq = status.RxSeq.Int64()
	   ls.BMRIndex = int(status.RelayIdx.Int64())
	   ls.RotateHeight = status.RotateHeight.Int64()
	   ls.RotateTerm = int(status.RotateTerm.Int64())
	   ls.DelayLimit = int(status.DelayLimit.Int64())
	   ls.MaxAggregation = int(status.MaxAggregation.Int64())
	   ls.CurrentHeight = status.CurrentHeight.Int64()
	   ls.RxHeight = status.RxHeight.Int64()
	   ls.RxHeightSrc = status.RxHeightSrc.Int64()
	   return ls, nil
   }

   func (s *sender) isOverLimit(size int) bool {
	   return txSizeLimit < float64(size)
   }

   func (s *sender) MonitorLoop(height int64, cb module.MonitorCallback, scb func()) error {
	   s.l.Debugf("MonitorLoop (sender) connected")
	   br := &BlockRequest{
		   Height: big.NewInt(height),
	   }
	   return s.c.MonitorBlock(br,
		   func(v *BlockNotification) error {
			   return cb(v.Height.Int64())
		   })
   }

   func (s *sender) StopMonitorLoop() {
	   s.c.CloseAllMonitor()
   }
   func (s *sender) FinalizeLatency() int {
	   //on-the-next
	   return 1
   }

   func NewSender(src, dst module.BtpAddress, w Wallet, endpoints []string, opt map[string]interface{}, l log.Logger) module.Sender {
	   s := &sender{
		   src: src,
		   dst: dst,
		   w:   w,
		   l:   l,
	   }
	   b, err := json.Marshal(opt)
	   if err != nil {
		   l.Panicf("fail to marshal opt:%#v err:%+v", opt, err)
	   }
	   if err = json.Unmarshal(b, &s.opt); err != nil {
		   l.Panicf("fail to unmarshal opt:%#v err:%+v", opt, err)
	   }
	   s.c = NewClient(endpoints, l)

	   s.bmc, _ = binding.NewBMC(HexToAddress(s.dst.ContractAddress()), s.c.ethcl)

	   return s
   }
*/
