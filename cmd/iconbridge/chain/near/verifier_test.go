package near

import (
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

					verifier := newVerifier(input.Options.Verifier.PreviousBlockHeight, input.Options.Verifier.PreviousBlockHash, input.Options.Verifier.NextEpochId, input.Options.Verifier.BlockProducers)

					bn := types.NewBlockNotification(int64(input.Offset))
					block, err := client.GetBlockByHeight(int64(input.Offset))
					require.Nil(f, err)

					bn.SetBlock(block)
					if verifier.blockHeight+1 == uint64(bn.Offset()) {
						bn.SetApprovalMessage(types.ApprovalMessage{
							Type:              [1]byte{types.ApprovalEndorsement},
							PreviousBlockHash: verifier.blockHash,
						})
					} else {
						bn.SetApprovalMessage(types.ApprovalMessage{
							Type:                [1]byte{types.ApprovalSkip},
							PreviousBlockHeight: verifier.blockHeight,
						})
					}


					err = verifier.validateHeader(bn)
					if testData.Expected.Success != nil {
						assert.Nil(f, err)
					} else {
						assert.Error(f, err)
					}

				})
			}
		})
	}
}
