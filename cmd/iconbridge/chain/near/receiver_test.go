package near

import (
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
					mockApi := NewMockApi(testData.MockStorage)
					client := &Client{
						api: &mockApi,
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
							Hash string
							Height uint64
						})
						assert.True(f, Ok)

						err = receiver.receiveBlocks(input.Offset, func(block *types.Block) error {
							if expected.Height == uint64(block.Height()) {
								assert.Equal(f, expected.Hash, block.Hash().Base58Encode())

								receiver.StopReceivingBlocks()
							}
							return nil
						})
						assert.Nil(f, err)
					} else {
						assert.Error(f, err)
					}
				})
			}
		})
	}
}
