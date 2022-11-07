package tests

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type VerifierValidateHeaderTest struct {
	description string
	testData    []TestData
}

func (t VerifierValidateHeaderTest) Description() string {
	return t.description
}

func (t VerifierValidateHeaderTest) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "Valid BlockHeader",
			Input: struct {
				Offset  uint64
				Options types.ReceiverOptions
			}{
				Offset: 377826,
				Options: types.ReceiverOptions{
					Verifier: types.VerifierConfig{
						PreviousBlockHeight: 377825,
						PreviousBlockHash:   types.NewCryptoHash("DDbjZ12VbmV36trcJDPxAAHsDWTtGEC9DB6ZSVLE9N1c"),
						BlockProducers: []*types.BlockProducer{
							{
								AccountId: types.AccountId(""),
								PublicKey: types.NewPublicKeyFromString(""),
								Stake:     types.NewBigInt(""),
							},
						},
						NextEpochId: types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
					},
				},
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
				})

				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: true,
			},
		},
	}

	RegisterTest("ValidateHeader", VerifierValidateHeaderTest{
		description: "Test ValidateHeader",
		testData:    testData,
	})
}
