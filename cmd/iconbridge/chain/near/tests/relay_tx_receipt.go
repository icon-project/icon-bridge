package tests

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
)

type RelayTxReceiptTest struct {
	description string
	testData    []TestData
}

func (t RelayTxReceiptTest) Description() string {
	return t.description
}

func (t RelayTxReceiptTest) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "RelayTx Receipt Success",
			Input: struct {
				PrivateKey      string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				PrivateKey:      "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("BroadcastTxAsync", mock.MockParam).Return(mockApi.BroadcastTxAsyncFactory())
				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("Transaction", mock.MockParam).Return(mockApi.TransactionFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: 377842,
			},
		},
	}

	RegisterTest("RelayTxReceipt", RelayTxReceiptTest{
		description: "Test RelayTx Receipt",
		testData:    testData,
	})
}
