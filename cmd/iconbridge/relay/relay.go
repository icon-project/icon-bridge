package relay

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	relayTickerInterval                  = 5 * time.Second
	relayBalanceCheckInterval            = 60 * time.Second
	relayTriggerReceiptsCount            = 20
	relayTxSendWaitInterval              = time.Second / 2
	relayTxReceiptWaitInterval           = time.Second
	relayInsufficientBalanceWaitInterval = 30 * time.Second
	retryWarnThreshold                   = 15
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
		fmt.Println("reached to srcMsg. receipts")
		receipts := srcMsg.Receipts[:0]
		for _, receipt := range srcMsg.Receipts {
			fmt.Println("receipt.height", receipt.Height)
			fmt.Println("rx_height", rxHeight)
			if receipt.Height < rxHeight {
				continue
			}
			events := receipt.Events[:0]
			for _, event := range receipt.Events {
				fmt.Println("event.seq: ", event.Sequence)
				fmt.Println("rx_seq:", rxSeq)
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
		fmt.Println(len(srcMsg.Receipts))
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

	relayBalanceCheckTicker := time.NewTicker(relayBalanceCheckInterval)
	defer relayBalanceCheckTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-relayTicker.C:
			relaySignal()

		case <-relayBalanceCheckTicker.C:
			go func() {
				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				bal, thres, err := r.dst.Balance(ctx)
				l := r.log.WithFields(log.Fields{"balance": bal, "threshold": thres})
				if err != nil {
					l.Error("failed to fetch relay wallet balance")
				} else if bal.Cmp(thres) <= 0 {
					l.Warn("relay wallet balance below threshold")
				}
			}()

		case err := <-srcErrCh:
			return err

		case msg := <-srcMsgCh:

			var seqBegin, seqEnd uint64
			receipts := msg.Receipts[:0]
			for _, receipt := range msg.Receipts {
				if len(receipt.Events) > 0 {
					if seqBegin == 0 {
						fmt.Println(receipt.Events[0], receipt.Height, receipt.Index)
						seqBegin = receipt.Events[0].Sequence
					}
					seqEnd = receipt.Events[len(receipt.Events)-1].Sequence
					receipts = append(receipts, receipt)
				}
			}
			msg.Receipts = receipts
			fmt.Println("length of msg.Receipts is ", len(msg.Receipts))
			if len(msg.Receipts) > 0 {
				r.log.WithFields(log.Fields{
					"seq": []uint64{seqBegin, seqEnd}}).Debug("srcMsg added")
				srcMsg.Receipts = append(srcMsg.Receipts, msg.Receipts...)
				fmt.Println(len(srcMsg.Receipts))
				fmt.Println(srcMsg.Receipts[0].Height)
				if len(srcMsg.Receipts) > relayTriggerReceiptsCount {
					relaySignal()
				}
			}

		case <-relayCh:

			fmt.Println("reached in status")
			fmt.Println(r.cfg.Name)

			link, err = r.dst.Status(ctx)
			fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXx")
			fmt.Println(link)
			if err != nil {
				r.log.WithFields(log.Fields{"error": err}).Debug("dst.Status: failed")
				if errors.Is(err, context.Canceled) {
					r.log.WithFields(log.Fields{"error": err}).Error("dst.Status: failed, Context Cancelled")
					return err
				}
				fmt.Println("continued from getting status")
				// TODO decide whether to ignore error or not
				continue
			}

			if link.CurrentHeight < txBlockHeight {
				fmt.Println("continued from here")
				continue // skip until dst.Status is updated
			}

			fmt.Println("before filtering the message")
			fmt.Println(len(srcMsg.Receipts))

			if missing := filterSrcMsg(link.RxHeight, link.RxSeq); missing > 0 {
				fmt.Println("did this filter the messages")
				r.log.WithFields(log.Fields{"rxSeq": missing}).Error("missing event sequence")
				return fmt.Errorf("missing event sequence")
			}

			fmt.Println("reached before sequence")
			fmt.Println("*****************************************************************")
			fmt.Println(len(srcMsg.Receipts))
			// fmt.Println(srcMsg.Receipts[0].Height)
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
					r.log.WithFields(log.Fields{"error": err}).Errorf(
						"add balance to relay account: waiting for %v", relayInsufficientBalanceWaitInterval)
					time.Sleep(relayInsufficientBalanceWaitInterval)
				default:
					time.Sleep(relayTxSendWaitInterval) // wait before sending tx
					if i > retryWarnThreshold {
						r.log.WithFields(log.Fields{"error": err}).Warnf("tx.Send: retry=%d", i)
					} else {
						r.log.WithFields(log.Fields{"error": err}).Debugf("tx.Send: retry=%d", i)
					}
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
					if retryCount > retryWarnThreshold {
						r.log.WithFields(log.Fields{"error": err, "retry": retryCount + 1}).Warn("tx.Receipt: ")
					} else {
						if strings.Contains(err.Error(), "not found") {
							r.log.WithFields(log.Fields{"retry": retryCount + 1}).Debug("tx.Receipt: ")
						} else {
							r.log.WithFields(log.Fields{"error": err, "retry": retryCount + 1}).Debug("tx.Receipt: ")
						}

					}
				}
				retryCount++
			}

		}

	}
}
