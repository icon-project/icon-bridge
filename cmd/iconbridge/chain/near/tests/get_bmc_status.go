package tests

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
)

type GetBmcLinkStatus struct {
	description string
	testData    []TestData
}

func (t GetBmcLinkStatus) Description() string {
	return t.description
}

func (t GetBmcLinkStatus) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "GetBmcStatus Sucess",
			Input: struct {
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				Source:      chain.BTPAddress("btp://0x1.icon/c294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("CallFunction", mock.MockParam).Return(mockApi.CallFunctionFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: 98855025,
				Fail:    nil,
			},
		},
	}

	RegisterTest("GetBmcLinkStatus", GetBmcLinkStatus{
		description: "Test_GetBmcStatus",
		testData:    testData,
	})
}
