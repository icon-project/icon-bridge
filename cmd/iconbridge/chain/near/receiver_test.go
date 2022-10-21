package near

import (
	"context"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearReceiver(t *testing.T) {
	if test, err := tests.GetTest("ReceiverReceiveBlocks", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						Offset      uint64
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					require.True(f, Ok)

					receiver, err := NewReceiver(ReceiverConfig{source: input.Source, destination: input.Destination}, log.New(), client)
					require.Nil(f, err)

					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(struct {
							Hash   string
							Height uint64
						})
						require.True(f, Ok)

						err = receiver.ReceiveBlocks(input.Offset, input.Source.ContractAddress(), func(blockNotification *types.BlockNotification) {
							if expected.Height == uint64(blockNotification.Offset()) {
								assert.Equal(f, expected.Hash, blockNotification.Block().Hash().Base58Encode())

								receiver.StopReceivingBlocks()
							}
						})
						assert.Nil(f, err)
					} else {
						assert.Error(f, err)
					}
				})
			}
		})
	}

	if test, err := tests.GetTest("ReceiverSubscribe", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						Seq         uint64
						Offset      uint64
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					require.True(f, Ok)

					receiver, err := NewReceiver(ReceiverConfig{source: input.Source, destination: input.Destination}, log.New(), client)
					require.Nil(f, err)
					srcMsgCh := make(chan *chain.Message)
					deadline, _ := f.Deadline()
					ctx, cancel := context.WithDeadline(context.Background(), deadline)
					defer cancel()
					errCh, err := receiver.Subscribe(ctx,
						srcMsgCh,
						chain.SubscribeOptions{
							Seq:    input.Seq,
							Height: input.Offset,
						})

					if testData.Expected.Success != nil {
						assert.True(f, testData.Expected.Success.(func(chan *chain.Message) bool)(srcMsgCh))
						assert.Nil(f, err)
					} else {
						assert.True(f, testData.Expected.Fail.(func(<-chan error, error) bool)(errCh, err))
					}
				})
			}
		})
	}
}
