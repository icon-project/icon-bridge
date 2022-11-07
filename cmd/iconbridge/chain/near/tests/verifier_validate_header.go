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
						NextBpHash:          types.NewCryptoHash("5QouG4ceHjyARjVTaySWXcXdsQduDqExVRKdwLjeANi"),
						NextEpochId:         types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						BlockProducers: []*types.BlockProducer{
							{
								ValidatorStakeStructVersion: []byte{0},
								AccountId:                   "node2",
								PublicKey:                   types.NewPublicKeyFromString("ed25519:GkDv7nSMS3xcqA45cpMvFmfV1o4fRF6zYo1JRR6mNqg5"),
								Stake:                       types.NewBigInt("50902386756263328030239719089112"),
							},
							{
								ValidatorStakeStructVersion: []byte{0},
								AccountId:                   "node1",
								PublicKey:                   types.NewPublicKeyFromString("ed25519:6DSjZ8mvsRZDvFqFxo8tCKePG96omXW7eVYVSySmDk8e"),
								Stake:                       types.NewBigInt("50879053856734837790009183916723"),
							},
							{
								ValidatorStakeStructVersion: []byte{0},
								AccountId:                   "node0",
								PublicKey:                   types.NewPublicKeyFromString("ed25519:7PGseFbWxvYVgZ89K1uTJKYoKetWs7BJtbyXDzfbAcqX"),
								Stake:                       types.NewBigInt("50868735427469424893279346843466"),
							},
						},
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
		{
			Description: "Valid invalid previous hash",
			Input: struct {
				Offset  uint64
				Options types.ReceiverOptions
			}{
				Offset: 377826,
				Options: types.ReceiverOptions{
					Verifier: types.VerifierConfig{
						PreviousBlockHeight: 377825,
						PreviousBlockHash:   types.NewCryptoHash("FDbjZ12VbmV36trcJDPxAAHsDWTtGEC9DB6ZSVLE9N1c"),
						NextBpHash:          types.NewCryptoHash("5QouG4ceHjyARjVTaySWXcXdsQduDqExVRKdwLjeANi"),
						NextEpochId:         types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						BlockProducers: []*types.BlockProducer{
							{
								ValidatorStakeStructVersion: []byte{0},
								AccountId:                   "node2",
								PublicKey:                   types.NewPublicKeyFromString("ed25519:GkDv7nSMS3xcqA45cpMvFmfV1o4fRF6zYo1JRR6mNqg5"),
								Stake:                       types.NewBigInt("50902386756263328030239719089112"),
							},
							{
								ValidatorStakeStructVersion: []byte{0},
								AccountId:                   "node1",
								PublicKey:                   types.NewPublicKeyFromString("ed25519:6DSjZ8mvsRZDvFqFxo8tCKePG96omXW7eVYVSySmDk8e"),
								Stake:                       types.NewBigInt("50879053856734837790009183916723"),
							},
							{
								ValidatorStakeStructVersion: []byte{0},
								AccountId:                   "node0",
								PublicKey:                   types.NewPublicKeyFromString("ed25519:7PGseFbWxvYVgZ89K1uTJKYoKetWs7BJtbyXDzfbAcqX"),
								Stake:                       types.NewBigInt("50868735427469424893279346843466"),
							},
						},
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
				Fail: func(err error) bool {
					return err.Error() == string("expected hash: 5YqjrSoiQjqrmrMHQJVZB25at7yQ2BZEC2exweLFmc6w, got hash: F7NwcVd8LMzi3PGzNXKz5v9Hq7k1P1UC7ys1JBrHV8m9")
				},
			},
		},
	}

	RegisterTest("ValidateHeader", VerifierValidateHeaderTest{
		description: "Test ValidateHeader",
		testData:    testData,
	})
}
