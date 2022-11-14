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
					Verifier: &types.VerifierConfig{
						BlockHeight:       377825,
						PreviousBlockHash: types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						CurrentBpsHash:    types.NewCryptoHash("C4zVnMf27hRJYoWEC816Pttyz122TWZN7zjUMoZCNkuw"),
						CurrentEpochId:    types.NewCryptoHash("FtrJuAXqH5oXDVADh6QkUyacf2MGmLHYbHCHKSZ8C7KS"),
						NextEpochId:       types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						NextBpsHash:       types.NewCryptoHash("5QouG4ceHjyARjVTaySWXcXdsQduDqExVRKdwLjeANi"),
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
				mockApi.On("BlockProducers", mock.MockParam).Return(mockApi.BlockProducersFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: true,
			},
		},
		{
			Description: "Invalid previous hash",
			Input: struct {
				Offset  uint64
				Options types.ReceiverOptions
			}{
				Offset: 377826,
				Options: types.ReceiverOptions{
					Verifier: &types.VerifierConfig{
						BlockHeight:       377825,
						PreviousBlockHash: types.NewCryptoHash("74toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						CurrentBpsHash:    types.NewCryptoHash("G2TyLP33XfqndppUzipoTWTs6XnKjmUhCQg1tH44isAG"),
						CurrentEpochId:    types.NewCryptoHash("FtrJuAXqH5oXDVADh6QkUyacf2MGmLHYbHCHKSZ8C7KS"),
						NextEpochId:       types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						NextBpsHash:       types.NewCryptoHash("5QouG4ceHjyARjVTaySWXcXdsQduDqExVRKdwLjeANi"),
					},
				},
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
					BlockProducersMap: map[string]mock.Response{
						"74toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95": {
							Reponse: []byte(`[]`),
						},
					},
				})

				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				mockApi.On("BlockProducers", mock.MockParam).Return(mockApi.BlockProducersFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Fail: func(err error) bool {
					return err.Error() == string("expected hash: 5YqjrSoiQjqrmrMHQJVZB25at7yQ2BZEC2exweLFmc6w, got hash: E78zFTZ21jN4nFaiygYqZNLw8iKMSYxrTNFG8X7eAJhE for block: 377826")
				},
			},
		},
	}

	RegisterTest("ValidateHeader", VerifierValidateHeaderTest{
		description: "Test ValidateHeader",
		testData:    testData,
	})
}
