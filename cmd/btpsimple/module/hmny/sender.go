package hmny

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/btp/cmd/btpsimple/module"
	"github.com/icon-project/btp/common/codec"
	"github.com/icon-project/btp/common/errors"
	"github.com/icon-project/btp/common/log"
)

const (
	txMaxDataSize                 = 524288 //512 * 1024 // 512kB
	txOverheadScale               = 0.37   //base64 encoding overhead 0.36, rlp and other fields 0.01
	txSizeLimit                   = txMaxDataSize / (1 + txOverheadScale)
	DefaultGetRelayResultInterval = time.Second
	DefaultRelayReSendInterval    = time.Second
	defaultGasLimit               = 10000000
)

type sender struct {
	c    *Client
	src  module.BtpAddress
	dst  module.BtpAddress
	w    Wallet
	log  log.Logger
	opts struct {
		GasLimit uint64 `json:"gasLimit"`
	}

	rmu sync.Mutex // relay mutex
}

func NewSender(src, dst module.BtpAddress, w Wallet, endpoints []string, opt map[string]interface{}, l log.Logger) module.Sender {
	s := &sender{
		src: src,
		dst: dst,
		w:   w,
		log: l,
	}
	b, err := json.Marshal(opt)
	if err != nil {
		l.Panicf("fail to marshal opt:%#v err:%+v", opt, err)
	}
	if err = json.Unmarshal(b, &s.opts); err != nil {
		l.Panicf("fail to unmarshal opt:%#v err:%+v", opt, err)
	}
	s.c = NewClient(endpoints, dst.ContractAddress(), l)
	return s
}

func (s *sender) newTransactionParam(prev string, rm *RelayMessage) (*BMCRelayMethodParams, error) {
	b, err := codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, err
	}
	rmp := &BMCRelayMethodParams{
		Prev:     prev,
		Messages: b,
		// Messages: base64.URLEncoding.EncodeToString(b),
	}
	return rmp, nil
}

func (s *sender) Segment(rm *module.RelayMessage, height int64) ([]*module.Segment, error) {
	segments := make([]*module.Segment, 0)
	var err error
	msg := &RelayMessage{
		ReceiptProofs: make([][]byte, 0),
	}
	//size := 0
	var b []byte
	for _, rp := range rm.ReceiptProofs {
		//TODO: segment for the events
		var eventBytes []byte
		if eventBytes, err = codec.RLP.MarshalToBytes(rp.Events); err != nil {
			return nil, err
		}
		trp := &ReceiptProof{
			Index:  rp.Index,
			Events: eventBytes,
			Height: rp.Height,
		}

		if b, err = codec.RLP.MarshalToBytes(trp); err != nil {
			return nil, err
		}
		msg.ReceiptProofs = append(msg.ReceiptProofs, b)

	}
	//
	segment := &module.Segment{
		Height:        msg.height,
		EventSequence: msg.eventSequence,
		NumberOfEvent: msg.numberOfEvent,
	}
	if segment.TransactionParam, err = s.newTransactionParam(rm.From.String(), msg); err != nil {
		return nil, err
	}
	s.log.Debugf("Segmentation Done")
	segments = append(segments, segment)
	return segments, nil
}

func (s *sender) Relay(segment *module.Segment) (module.GetResultParam, error) {
	s.rmu.Lock()
	defer s.rmu.Unlock()
	rmp, ok := segment.TransactionParam.(*BMCRelayMethodParams)
	if !ok {
		return nil, fmt.Errorf("casting failure")
	}
	t := s.c.newTransactOpts(s.w)
	if s.opts.GasLimit > 0 {
		t.GasLimit = s.opts.GasLimit
	} else {
		t.GasLimit = defaultGasLimit
	}
	s.log.Debugf("HandleRelayMessage prev %s, msg: %s", rmp.Prev, common.Bytes2Hex(rmp.Messages))
	tx, err := s.c.bmc().HandleRelayMessage(t, rmp.Prev, rmp.Messages)
	if err != nil {
		s.log.Errorf("relay: bmc.handleRelayMessage: %v", err)
		return nil, err
	}
	return &TransactionHashParam{Hash: tx.Hash()}, nil
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
	var status TypesLinkStats
	status, err := s.c.bmc().GetStatus(nil, s.src.String())

	if err != nil {
		s.log.Errorf("Error retrieving relay status from BMC")
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
	if err := s.c.MonitorBlock(uint64(height),
		false, func(next *BlockNotification) error {
			s.log.Debugf("monitor loop: block notification: height=%d", next.Height)
			return cb(next.Height.Int64())
		}); err != nil {
		return errors.Wrapf(err, "monitor loop terminated: %v", err)
	}
	return nil
}

func (s *sender) StopMonitorLoop() {
	s.c.CloseAllMonitor()
}

func (s *sender) FinalizeLatency() int {
	//on-the-next
	return 1
}
