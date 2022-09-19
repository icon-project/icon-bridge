package near

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNearClient(t *testing.T) {
	if test, err := tests.GetTest("SendTransaction", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api: testData.MockApi,
					}

					input, Ok := (testData.Input).(string)
					assert.True(f, Ok)

					transactionHash, err := client.SendTransaction(input)
					assert.Nil(f, err)
					assert.NotNil(f, transactionHash)
					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(string)
						assert.True(f, Ok)
						assert.Equal(f, expected, transactionHash.Base58Encode())
					} else {
						assert.Error(f, err)
					}

				})
			}
		})
	}
}
