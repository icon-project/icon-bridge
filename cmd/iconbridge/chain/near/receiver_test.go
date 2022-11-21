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
						Options     types.ReceiverOptions
					})
					require.True(f, Ok)

					receiver, err := NewReceiver(ReceiverConfig{input.Source, input.Destination, input.Options}, log.New(), client)
					require.Nil(f, err)

					if testData.Expected.Success != nil {
						err = receiver.ReceiveBlocks(input.Offset, input.Source.ContractAddress(), func(blockNotification *types.BlockNotification) {
							assert.True(f, testData.Expected.Success.(func(*types.BlockNotification, func()) bool)(blockNotification, receiver.StopReceivingBlocks))
						})
						assert.Nil(f, err)
					} else {
						err = receiver.ReceiveBlocks(input.Offset, input.Source.ContractAddress(), func(blockNotification *types.BlockNotification) {
							if err != nil {
								assert.True(f, testData.Expected.Fail.(func(error) bool)(err))
							}
						})
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
						Options     types.ReceiverOptions
					})
					require.True(f, Ok)

					receiver, err := NewReceiver(ReceiverConfig{input.Source, input.Destination, input.Options}, log.New(), client)
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
