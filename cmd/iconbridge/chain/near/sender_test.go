package near

import (
	"context"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNearSender(t *testing.T) {
	if test, err := testdata.GetTest("GetBmcLinkStatus", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					mockApi := NewMockApi(testData.MockStorage)
					client := &Client{
						api: &mockApi,
					}

					links, Ok := (testData.Input).([]chain.BTPAddress)
					assert.True(f, Ok)

					sender, err := newMockSender(links[1], links[0], client, nil, nil, nil)
					assert.Nil(f, err)
					
					status, err := sender.Status(context.Background())
					assert.Nil(f, err)
					assert.NotNil(f, status)
					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(int)
						assert.True(f, Ok)
						assert.Equal(f, uint64(expected), status.RotateHeight)
					} else {
						assert.Error(f, err)
					}

				})
			}
		})
	}
}
