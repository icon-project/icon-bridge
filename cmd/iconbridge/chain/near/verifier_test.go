package near

import (
	"sync"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearVerifier(t *testing.T) {
	if test, err := tests.GetTest("ValidateHeader", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						Offset  uint64
						Options types.ReceiverOptions
					})
					require.True(f, Ok)

					verifier, err := NewVerifier(input.Options.Verifier.BlockHeight, input.Options.Verifier.PreviousBlockHash, input.Options.Verifier.CurrentEpochId, input.Options.Verifier.NextEpochId, input.Options.Verifier.CurrentBpsHash, input.Options.Verifier.NextBpsHash, 100, client)
					require.Nil(f, err)

					wg := new(sync.WaitGroup)
					wg.Add(1)

					err = verifier.SyncHeader(wg, input.Offset-1)
					require.Nil(f, err)

					wg.Wait()

					bn := types.NewBlockNotification(int64(input.Offset))
					block, err := client.GetBlockByHeight(int64(input.Offset))
					require.Nil(f, err)

					bn.SetBlock(block)

					err = verifier.ValidateHeader(bn)
					if testData.Expected.Success != nil {
						assert.Nil(f, err)
					} else {
						assert.True(f, testData.Expected.Fail.(func(error) bool)(err))
					}

				})
			}
		})
	}
}
