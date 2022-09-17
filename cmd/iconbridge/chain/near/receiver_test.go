package near

import (
	"context"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/assert"
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
					assert.True(f, Ok)

					receiver, err := newMockReceiver(input.Source, input.Destination, client, nil, nil, log.New())
					assert.Nil(f, err)

					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(struct {
							Hash   string
							Height uint64
						})
						assert.True(f, Ok)

						err = receiver.receiveBlocks(input.Offset, input.Source.ContractAddress(), func(blockNotification *types.BlockNotification) {
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
						Offset      uint64
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					assert.True(f, Ok)

					receiver, err := newMockReceiver(input.Source, input.Destination, client, nil, nil, log.New())
					assert.Nil(f, err)

					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(struct {
							From chain.BTPAddress
						})
						assert.True(f, Ok)

						srcMsgCh := make(chan *chain.Message)

						deadline, _ := f.Deadline()
						ctx, cancel := context.WithDeadline(context.Background(), deadline)
						defer cancel()
						_, err := receiver.Subscribe(ctx,
							srcMsgCh,
							chain.SubscribeOptions{
								Seq:    1,
								Height: input.Offset,
							})

						for msg := range srcMsgCh {
							f.Log(msg)
							assert.Equal(f, msg.From, expected.From)
							break
						}

						assert.Nil(f, err)
					} else {
						assert.Error(f, err)
					}
				})
			}
		})
	}
}
