package tests

import (
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
	"github.com/shopspring/decimal"
)

type GetBalanceTest struct {
	description string
	testData    []TestData
}

func (t GetBalanceTest) Description() string {
	return t.description
}

func (t GetBalanceTest) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "GetBalance Pass",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("ViewAccount", mock.MockParam).Return(mockApi.ViewAccountFactory())
				
				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func() *big.Int {
					decimal, _ := decimal.NewFromString("399992611103597728750000000")
					return decimal.BigInt()
				}(),
			},
		},
	}

	RegisterTest("GetBalance", GetBalanceTest{
		description: "Test GetBalance",
		testData:    testData,
	})
}
