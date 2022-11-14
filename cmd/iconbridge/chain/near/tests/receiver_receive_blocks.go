package tests

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type ReceiverReceiveBlocks struct {
	description string
	testData    []TestData
}

func (t ReceiverReceiveBlocks) Description() string {
	return t.description
}

func (t ReceiverReceiveBlocks) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "Receiver Receive Blocks Success",
			Input: struct {
				Offset      uint64
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Options     types.ReceiverOptions
			}{
				Offset:      377825,
				Source:      chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Destination: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Options:     types.ReceiverOptions{},
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})
				
				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
				})

				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				mockApi.On("Changes", mock.MockParam).Return(mockApi.ChangesFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: struct {
					Hash   string
					Height uint64
				}{
					Hash:   "4CSFBudwkUAgoHHQNh6UnVx78EP9ubeUhbc9ZHoC5w4u",
					Height: 377828,
				},
				Fail: nil,
			},
		},
	}

	RegisterTest("ReceiverReceiveBlocks", ReceiverReceiveBlocks{
		description: "Test Receiver Receive Blocks",
		testData:    testData,
	})
}
