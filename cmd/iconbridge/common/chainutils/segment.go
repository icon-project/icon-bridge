package chainutils

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/codec"
)

func SegmentByTxDataSize(msg *chain.Message, txDataSizeLimit uint64) (relayMsg *chain.RelayMessage, newMsg *chain.Message, err error) {

	relayMsg = &chain.RelayMessage{
		Receipts: make([][]byte, 0),
	}

	newMsg = &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts,
	}

	var msgSize uint64

rloop:
	for i, receipt := range msg.Receipts {

	eloop:
		for j := range receipt.Events {
			// try all events first
			// if it exceeds limit, try again by removing last event
			events := receipt.Events[:len(receipt.Events)-j]

			rlpEvents, err := codec.RLP.MarshalToBytes(events)
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
			if newMsgSize <= txDataSizeLimit {

				msgSize = newMsgSize
				if len(events) == len(receipt.Events) { // all events added
					newMsg.Receipts = msg.Receipts[i+1:]
				} else { // save remaining events in this receipt
					newMsg.Receipts = make([]*chain.Receipt, len(msg.Receipts[i:]))
					copy(newMsg.Receipts, msg.Receipts[i:]) // make a copy to not mutate original receipt
					newMsg.Receipts[0] = &chain.Receipt{
						Index:  receipt.Index,
						Height: receipt.Height,
						Events: receipt.Events[len(events):],
					}
				}
				relayMsg.Receipts = append(relayMsg.Receipts, rlpReceipt)
				break eloop

			} else if len(events) == 1 {
				// stop iterating over receipts when adding even a single event
				// exceeds tx size limit
				break rloop
			}

		}
	}

	return relayMsg, newMsg, nil

}
