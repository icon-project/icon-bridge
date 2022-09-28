package chainutils

import (
	"testing"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/stretchr/testify/require"
)

func getTestMsg() *chain.Message {
	return &chain.Message{
		From: "",
		Receipts: []*chain.Receipt{
			{Index: 0, Height: 1, Events: []*chain.Event{
				{Next: "", Sequence: 0, Message: []byte{}},
				{Next: "", Sequence: 0, Message: []byte{}},
			}},
			{Index: 2, Height: 6, Events: []*chain.Event{
				{Next: "", Sequence: 0, Message: []byte{}},
				{Next: "", Sequence: 0, Message: []byte{}},
				{Next: "", Sequence: 0, Message: []byte{}},
				{Next: "", Sequence: 0, Message: []byte{}},
			}},
			{Index: 5, Height: 7, Events: []*chain.Event{
				{Next: "", Sequence: 0, Message: []byte{}},
			}},
			{Index: 0, Height: 10, Events: []*chain.Event{
				{Next: "", Sequence: 0, Message: []byte("msg")},
				{Next: "", Sequence: 0, Message: []byte{}},
				{Next: "", Sequence: 0, Message: []byte{}},
			}},
		},
	}
}

func TestSegmentByTxDataSize(t *testing.T) {

	/* byte sizes
	receipt:
		with 1 event: 9 bytes
		with more events: 4 bytes each
	*/

	srcMsg := getTestMsg()
	_, msg, err := SegmentByTxDataSize(srcMsg, 5)
	require.NoError(t, err)
	require.EqualValues(t, 4, len(msg.Receipts))
	require.EqualValues(t, 1, msg.Receipts[0].Height)
	require.EqualValues(t, 2, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 9) // encode first receipt with 1 event
	require.NoError(t, err)
	require.EqualValues(t, 4, len(msg.Receipts))
	require.EqualValues(t, 1, msg.Receipts[0].Height)
	require.EqualValues(t, 1, len(msg.Receipts[0].Events))
	require.EqualValues(t, 2, len(srcMsg.Receipts[0].Events), "events in original receipt should not mutate")

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 12) // encode first receipt with 1 event
	require.NoError(t, err)
	require.EqualValues(t, 4, len(msg.Receipts))
	require.EqualValues(t, 1, msg.Receipts[0].Height)
	require.EqualValues(t, 1, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 15) // encode first receipt with 2 events
	require.NoError(t, err)
	require.EqualValues(t, 3, len(msg.Receipts))
	require.EqualValues(t, 6, msg.Receipts[0].Height)
	require.EqualValues(t, 4, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 22) // encode first receipt with 2 events and second receipt with 1 event
	require.NoError(t, err)
	require.EqualValues(t, 3, len(msg.Receipts))
	require.EqualValues(t, 6, msg.Receipts[0].Height)
	require.EqualValues(t, 3, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 35) // encode first receipt with 2 events and second receipt with 4 event
	require.NoError(t, err)
	require.EqualValues(t, 2, len(msg.Receipts))
	require.EqualValues(t, 7, msg.Receipts[0].Height)
	require.EqualValues(t, 1, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 43) // encode first receipt with 2 events and second receipt with 4 events and third receipt with 1 event
	require.NoError(t, err)
	require.EqualValues(t, 1, len(msg.Receipts))
	require.EqualValues(t, 10, msg.Receipts[0].Height)
	require.EqualValues(t, 3, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 52) // encode first receipt with 2 events and second receipt with 4 events and third receipt with 1 event
	require.NoError(t, err)
	require.EqualValues(t, 1, len(msg.Receipts))
	require.EqualValues(t, 10, msg.Receipts[0].Height)
	require.EqualValues(t, 3, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 55) // encode first receipt with 2 events and second receipt with 4 events and third receipt with 1 event and fourth receipt with 1 event with message body
	require.NoError(t, err)
	require.EqualValues(t, 1, len(msg.Receipts))
	require.EqualValues(t, 10, msg.Receipts[0].Height)
	require.EqualValues(t, 2, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 60) // encode first receipt with 2 events and second receipt with 4 events and third receipt with 1 event and fourth receipt with 2 events with one with message body
	require.NoError(t, err)
	require.EqualValues(t, 1, len(msg.Receipts))
	require.EqualValues(t, 10, msg.Receipts[0].Height)
	require.EqualValues(t, 1, len(msg.Receipts[0].Events))

	srcMsg = getTestMsg()
	_, msg, err = SegmentByTxDataSize(srcMsg, 70) // encode all receipts
	require.NoError(t, err)
	require.EqualValues(t, 0, len(msg.Receipts))

}
