package relay

import (
	"context"
	"errors"
	"time"

	"github.com/icon-project/btp/cmd/btpsimple/chain"
	"github.com/icon-project/btp/common/log"
)

const (
	relayTickerInterval       = 5 * time.Second
	relayTriggerReceiptsCount = 20
)

type Relay interface {
	Start(ctx context.Context) (err error)
}

func NewRelay(cfg *RelayConfig, src chain.Receiver, dst chain.Sender, log log.Logger) (Relay, error) {
	r := &relay{
		cfg: cfg,
		log: log,
		src: src,
		dst: dst,
	}
	return r, nil
}

type relay struct {
	cfg *RelayConfig
	log log.Logger
	src chain.Receiver
	dst chain.Sender
}

func (r *relay) rxHeight(linkRxHeight uint64) uint64 {
	height := linkRxHeight
	if r.cfg.Src.Offset > height {
		height = r.cfg.Src.Offset
	}
	return height
}

func (r *relay) createMessage() *chain.Message {
	return &chain.Message{
		From: r.cfg.Src.Address,
	}
}

func (r *relay) Start(ctx context.Context) error {

	link, err := r.dst.Status(ctx)
	if err != nil {
		return err
	}
	r.log.Infof("init: link.rxSeq=%d, link.rxHeight=%d", link.RxSeq, link.RxHeight)

	srcMessageCh, err := r.src.
		SubscribeMessage(ctx, r.rxHeight(link.RxHeight), link.RxSeq)
	if err != nil {
		return err
	}

	srcMsg := r.createMessage()

	removeProcessedMessages := func(rxHeight, rxSeq uint64) {
		receipts := srcMsg.Receipts[:0]
		for _, receipt := range srcMsg.Receipts {
			if receipt.Height < rxHeight {
				continue
			}
			events := receipt.Events[:0]
			for _, event := range receipt.Events {
				if event.Sequence > rxSeq {
					events = append(events, event)
				}
			}
			receipt.Events = events
			if len(receipt.Events) > 0 {
				receipts = append(receipts, receipt)
			}
		}
		srcMsg.Receipts = receipts
	}

	relayCh := make(chan struct{}, 1)
	relayTicker := time.NewTicker(relayTickerInterval)
	defer relayTicker.Stop()
	relaySignal := func() {
		select {
		case relayCh <- struct{}{}:
		default:
		}
		relayTicker.Reset(relayTickerInterval)
		r.log.Debug("relaySignal")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-relayTicker.C:
			relaySignal()

		case msg := <-srcMessageCh:

			var seqBegin, seqEnd uint64
			receipts := msg.Receipts[:0]
			for _, receipt := range msg.Receipts {
				if len(receipt.Events) > 0 {
					if seqBegin == 0 {
						seqBegin = receipt.Events[0].Sequence
					}
					seqEnd = receipt.Events[len(receipt.Events)-1].Sequence
					receipts = append(receipts, receipt)
				}
			}
			msg.Receipts = receipts

			if len(msg.Receipts) > 0 {
				r.log.WithFields(log.Fields{
					"seq": []uint64{seqBegin, seqEnd}}).Info("srcMsg added")
				srcMsg.Receipts = append(srcMsg.Receipts, msg.Receipts...)
				if len(srcMsg.Receipts) > relayTriggerReceiptsCount {
					relaySignal()
				}
			}

		case <-relayCh:

			link, err = r.dst.Status(ctx)
			if err != nil {
				r.log.WithFields(log.Fields{"error": err}).Debug("dst.Status: failed")
				if errors.Is(err, context.Canceled) {
					r.log.WithFields(log.Fields{"error": err}).Error("dst.Status: failed")
					return err
				}
				// TODO decide whether to ignore error or not
				continue
			}

			removeProcessedMessages(link.RxHeight, link.RxSeq)

			tx, newMsg, err := r.dst.Segment(ctx, srcMsg, r.cfg.Dst.TxDataSizeLimit)
			if err != nil {
				return err
			}
			if tx == nil { // ignore if tx is nil
				continue
			}
			for i, err := 1, tx.Send(ctx); err != nil; i, err = i+1, tx.Send(ctx) {
				if errors.Is(err, context.Canceled) {
					r.log.WithFields(log.Fields{"error": err}).Error("tx.Send: failed", i)
					return err
				}
				r.log.WithFields(log.Fields{"error": err}).Debugf("tx.Send: retrying=%d", i)
				time.Sleep(time.Second / 2) // wait before sending tx
			}
			for txr, err := tx.Receipt(ctx); true; txr, err = tx.Receipt(ctx) {
				if err == nil {
					newMsg.From, srcMsg = srcMsg.From, newMsg
					break
				}
				if errors.Is(err, context.Canceled) {
					r.log.WithFields(log.Fields{"error": err}).Error("tx.Receipt: failed")
					return err
				}
				r.log.WithFields(log.Fields{"error": err}).Debug("tx.Receipt: retrying")
				if txr != nil {
					break
				}
				time.Sleep(time.Second) // wait before asking for receipt
			}

		}

	}
}
