package relay

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	relayTickerInterval                  = 5 * time.Second
	relayTriggerReceiptsCount            = 20
	relayTxSendWaitInterval              = time.Second / 2
	relayTxReceiptWaitInterval           = time.Second
	relayInsufficientBalanceWaitInterval = 30 * time.Second
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

func (r *relay) Start(ctx context.Context) error {

	link, err := r.dst.Status(ctx)
	if err != nil {
		return err
	}
	r.log.WithFields(log.Fields{
		"rxSeq":         link.RxSeq,
		"rxHeight":      link.RxHeight,
		"currentHeight": link.CurrentHeight,
	}).Info("link status")

	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := r.src.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    link.RxSeq,
			Height: r.rxHeight(link.RxHeight),
		})
	if err != nil {
		return err
	}

	srcMsg := &chain.Message{
		From: r.cfg.Src.Address,
	}

	filterSrcMsg := func(rxHeight, rxSeq uint64) (missingRxSeq uint64) {
		receipts := srcMsg.Receipts[:0]
		for _, receipt := range srcMsg.Receipts {
			if receipt.Height < rxHeight {
				continue
			}
			events := receipt.Events[:0]
			for _, event := range receipt.Events {
				if event.Sequence > rxSeq {
					rxSeq++
					if event.Sequence != rxSeq {
						return rxSeq
					}
					events = append(events, event)
				}
			}
			receipt.Events = events
			if len(receipt.Events) > 0 {
				receipts = append(receipts, receipt)
			}
		}
		srcMsg.Receipts = receipts
		return 0
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

	txBlockHeight := link.CurrentHeight

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-relayTicker.C:
			relaySignal()

		case err := <-srcErrCh:
			return err

		case msg := <-srcMsgCh:

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

			if link.CurrentHeight < txBlockHeight {
				continue // skip until dst.Status is updated
			}

			if missing := filterSrcMsg(link.RxHeight, link.RxSeq); missing > 0 {
				r.log.WithFields(log.Fields{"rxSeq": missing}).Error("missing event sequence")
				return fmt.Errorf("missing event sequence")
			}

			tx, newMsg, err := r.dst.Segment(ctx, srcMsg)
			if err != nil {
				return err
			} else if tx == nil { // ignore if tx is nil
				continue
			}

		sendLoop:
			for i, err := 1, tx.Send(ctx); true; i, err = i+1, tx.Send(ctx) {
				switch {
				case err == nil:
					break sendLoop
				case errors.Is(err, context.Canceled):
					r.log.WithFields(log.Fields{"id": tx.ID(), "error": err}).Error("tx.Send failed")
					return err
				case errors.Is(err, chain.ErrInsufficientBalance):
					r.log.WithFields(log.Fields{"error": err}).Error(
						"add balance to relay account: waiting for %v", relayInsufficientBalanceWaitInterval)
					time.Sleep(relayInsufficientBalanceWaitInterval)
				default:
					time.Sleep(relayTxSendWaitInterval) // wait before sending tx
					r.log.WithFields(log.Fields{"error": err}).Debugf("tx.Send: retry=%d", i)
				}
			}

			retryCount := 0
		waitLoop:
			for blockHeight, err := tx.Receipt(ctx); retryCount < 30; _, err = tx.Receipt(ctx) {
				switch {
				case err == nil:
					newMsg.From = srcMsg.From
					srcMsg = newMsg
					txBlockHeight = blockHeight
					break waitLoop
				case errors.Is(err, context.Canceled):
					r.log.WithFields(log.Fields{"error": err}).Error("tx.Receipt failed")
					return err
				case errors.Is(err, chain.ErrGasLimitExceeded):
					// increase transaction gas limit
				case errors.Is(err, chain.ErrBlockGasLimitExceeded):
					// reduce batch size
				case errors.Is(err, chain.ErrBMCRevertInvalidSeqNumber):
					// messages skipped; refetch from source

				default:
					time.Sleep(relayTxReceiptWaitInterval) // wait before asking for receipt
					r.log.WithFields(log.Fields{"error": err, "retry": retryCount + 1}).Debug("tx.Receipt: retry")
				}
				retryCount++
			}

		}

	}
}
