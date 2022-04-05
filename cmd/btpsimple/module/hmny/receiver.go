package hmny

import (
	"bytes"
	"encoding/json"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/btp/cmd/btpsimple/module"
	"github.com/icon-project/btp/common/errors"
	"github.com/icon-project/btp/common/log"
)

type receiver struct {
	c     *Client
	src   module.BtpAddress
	dst   module.BtpAddress
	log   log.Logger
	rxSeq uint64
	opts  struct{}
}

func NewReceiver(src, dst module.BtpAddress, endpoints []string, opt map[string]interface{}, l log.Logger) module.Receiver {
	r := &receiver{
		src: src,
		dst: dst,
		log: l,
	}
	b, err := json.Marshal(opt)
	if err != nil {
		l.Panicf("fail to marshal opt:%#v err:%+v", opt, err)
	}
	if err = json.Unmarshal(b, &r.opts); err != nil {
		l.Panicf("fail to unmarshal opt:%#v err:%+v", opt, err)
	}
	r.c = NewClient(endpoints, src.ContractAddress(), l)
	return r
}

func (r *receiver) newReceiptProofs(v *BlockNotification) ([]*module.ReceiptProof, error) {
	sc := HexToAddress(r.src.ContractAddress())
	rps := make([]*module.ReceiptProof, 0, len(v.Receipts))
	var events []*module.Event
	for i, receipt := range v.Receipts {
		events = events[:0]
		for _, log := range receipt.Logs {
			if !bytes.Equal(log.Address.Bytes(), sc.Bytes()) {
				continue
			}
			msg, err := r.c.bmc().ParseMessage(ethtypes.Log{
				Data: log.Data, Topics: log.Topics,
			})
			if err == nil {
				events = append(events, &module.Event{
					Message:  msg.Msg,
					Next:     module.BtpAddress(msg.Next),
					Sequence: msg.Seq.Int64(),
				})
			}
		}
		if len(events) > 0 {
			rp := &module.ReceiptProof{}
			rp.Index, rp.Height = i, v.Height.Int64()
			rp.Events = append(rp.Events, events...)
			rps = append(rps, rp)
			r.log.Debugf("received event: h=%d: sc=%v", rp.Height, sc)
		}
	}
	return rps, nil
}

func (r *receiver) ReceiveLoop(height int64, seq int64, cb module.ReceiveCallback, scb func()) error {
	r.rxSeq = uint64(seq)
	var v *BlockNotification
	if err := r.c.MonitorBlock(uint64(height), true, func(next *BlockNotification) error {
		r.log.Debugf("receive loop: block notification: height=%d", next.Height)
		if v != nil {
			if next.Height.Int64() != v.Height.Int64()+1 {
				return fmt.Errorf(
					"receive loop: next.Height (%d) != v.Height (%d)",
					next.Height.Int64(), v.Height.Int64())
			}
			rps, err := r.newReceiptProofs(v)
			if err != nil {
				return errors.Wrapf(err, "receipt proofs: %v", err)
			}
			cb(rps)
		}
		v = next
		return nil
	}); err != nil {
		return errors.Wrapf(err, "receive loop: terminated: %v", err)
	}
	return nil
}

func (r *receiver) StopReceiveLoop() {
	r.c.CloseAllMonitor()
}
