package tests

import (
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
)

type ReceiverSubscribe struct {
	description string
	testData    []TestData
}

func (t ReceiverSubscribe) Description() string {
	return t.description
}

func (t ReceiverSubscribe) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "With Btp Messages",
			Input: struct {
				Seq         uint64
				Offset      uint64
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				Seq:         0,
				Offset:      377825,
				Source:      chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Destination: chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315"),
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
				Success: func(ouput chan *chain.Message) bool {
					for msg := range ouput {
						if len(msg.Receipts) > 0 {
							break
						}
					}

					return true
				},
				Fail: nil,
			},
		},
		{
			Description: "With Btp Messages only for different destination",
			Input: struct {
				Seq         uint64
				Offset      uint64
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				Seq:         0,
				Offset:      377825,
				Source:      chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Destination: chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315"),
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
					ContractStateChangeMap: map[int64]mock.Response{
						377825: {
							Reponse: []byte(`{ "block_hash": "DDbjZ12VbmV36trcJDPxAAHsDWTtGEC9DB6ZSVLE9N1c", "changes": [ { "cause": { "receipt_hash": "2VWWEfyg5BzyDBVRsHRXApQJdZ37Bdtj8GtkH7UvNm7G", "type": "receipt_processing" }, "change": { "account_id": "alice.node1", "key_base64": "bWVzc2FnZQ==", "value_base64": "2AEAAHsibmV4dCI6ImJ0cDovLzB4Mi5pY29uL2N4YjM1ZGFhNmNkZmZjMWMzNjQ2ZGI2NTgzNDVlZDVlNTNlNDA5YjMwNSIsInNlcXVlbmNlIjoiMSIsIm1lc3NhZ2UiOiIrUUVXdUU5aWRIQTZMeTh3ZURJdWJtVmhjaTgzTWpjd1lUYzVZbVUzT0Rsa056Y3daakprWlRBeE5UQTBOelk0TkdVeU9EQTJOVGszWldWbFpUazJaV1V6WTJFNE4ySXhOemxqTmpNNU9XUmxZV0ZtdURsaWRIQTZMeTh3ZURJdWFXTnZiaTlqZUdJek5XUmhZVFpqWkdabVl6RmpNelkwTm1SaU5qVTRNelExWldRMVpUVXpaVFF3T1dJek1EV0RZbTFqQUxpRCtJR0VTVzVwZExoNitIajRkcmc1WW5Sd09pOHZNSGczTG1samIyNHZZM2d4WVdRMlptTmpORFkxWkRGaU9EWTBOR05oTXpjMVpqbGxNVEJpWVdKbFpXRTBZek00TXpFMXVEbGlkSEE2THk4d2VESXVhV052Ymk5amVHSXpOV1JoWVRaalpHWm1ZekZqTXpZME5tUmlOalU0TXpRMVpXUTFaVFV6WlRRd09XSXpNRFU9In0=" }, "type": "data_update" } ] }`),
							Error:   nil,
						},
						377826: {
							Reponse: []byte(`{ "block_hash": "DDbjZ12VbmV36trcJDPxAAHsDWTtGEC9DB6ZSVLE9N1c", "changes": [ { "cause": { "receipt_hash": "2VWWEfyg5BzyDBVRsHRXApQJdZ37Bdtj8GtkH7UvNm7G", "type": "receipt_processing" }, "change": { "account_id": "alice.node1", "key_base64": "bWVzc2FnZQ==", "value_base64": "2AEAAHsibmV4dCI6ImJ0cDovLzB4Mi5pY29uL2N4YjM1ZGFhNmNkZmZjMWMzNjQ2ZGI2NTgzNDVlZDVlNTNlNDA5YjMwNSIsInNlcXVlbmNlIjoiMSIsIm1lc3NhZ2UiOiIrUUVXdUU5aWRIQTZMeTh3ZURJdWJtVmhjaTgzTWpjd1lUYzVZbVUzT0Rsa056Y3daakprWlRBeE5UQTBOelk0TkdVeU9EQTJOVGszWldWbFpUazJaV1V6WTJFNE4ySXhOemxqTmpNNU9XUmxZV0ZtdURsaWRIQTZMeTh3ZURJdWFXTnZiaTlqZUdJek5XUmhZVFpqWkdabVl6RmpNelkwTm1SaU5qVTRNelExWldRMVpUVXpaVFF3T1dJek1EV0RZbTFqQUxpRCtJR0VTVzVwZExoNitIajRkcmc1WW5Sd09pOHZNSGczTG1samIyNHZZM2d4WVdRMlptTmpORFkxWkRGaU9EWTBOR05oTXpjMVpqbGxNVEJpWVdKbFpXRTBZek00TXpFMXVEbGlkSEE2THk4d2VESXVhV052Ymk5amVHSXpOV1JoWVRaalpHWm1ZekZqTXpZME5tUmlOalU0TXpRMVpXUTFaVFV6WlRRd09XSXpNRFU9In0=" }, "type": "data_update" } ] }`),
							Error:   nil,
						},
					},
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
				Success: func(ouput chan *chain.Message) bool {
					result := true
					go func() {
						for msg := range ouput {
							if len(msg.Receipts) > 0 {
								result = false
							}
						}
					}()

					time.Sleep(time.Second * 3)
					return result
				},
				Fail: nil,
			},
		},
		{
			Description: "With Btp Messages for destination",
			Input: struct {
				Seq         uint64
				Offset      uint64
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				Seq:         0,
				Offset:      377825,
				Source:      chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Destination: chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315"),
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
				Success: func(ouput chan *chain.Message) bool {
				msgloop:
					for msg := range ouput {
						for _, j := range msg.Receipts {
							for _, e := range j.Events {
								if e.Next == chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315") {
									break msgloop
								}
							}
						}
					}
					return true
				},
				Fail: nil,
			},
		},
		{
			Description: "With Btp Messages inlcuding messages for different destination",
			Input: struct {
				Seq         uint64
				Offset      uint64
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				Seq:         0,
				Offset:      377825,
				Source:      chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Destination: chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315"),
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
					ContractStateChangeMap: map[int64]mock.Response{
						377825: {
							Reponse: []byte(`{ "block_hash": "DDbjZ12VbmV36trcJDPxAAHsDWTtGEC9DB6ZSVLE9N1c", "changes": [ { "cause": { "receipt_hash": "2VWWEfyg5BzyDBVRsHRXApQJdZ37Bdtj8GtkH7UvNm7G", "type": "receipt_processing" }, "change": { "account_id": "alice.node1", "key_base64": "bWVzc2FnZQ==", "value_base64": "2AEAAHsibmV4dCI6ImJ0cDovLzB4Mi5pY29uL2N4YjM1ZGFhNmNkZmZjMWMzNjQ2ZGI2NTgzNDVlZDVlNTNlNDA5YjMwNSIsInNlcXVlbmNlIjoiMSIsIm1lc3NhZ2UiOiIrUUVXdUU5aWRIQTZMeTh3ZURJdWJtVmhjaTgzTWpjd1lUYzVZbVUzT0Rsa056Y3daakprWlRBeE5UQTBOelk0TkdVeU9EQTJOVGszWldWbFpUazJaV1V6WTJFNE4ySXhOemxqTmpNNU9XUmxZV0ZtdURsaWRIQTZMeTh3ZURJdWFXTnZiaTlqZUdJek5XUmhZVFpqWkdabVl6RmpNelkwTm1SaU5qVTRNelExWldRMVpUVXpaVFF3T1dJek1EV0RZbTFqQUxpRCtJR0VTVzVwZExoNitIajRkcmc1WW5Sd09pOHZNSGczTG1samIyNHZZM2d4WVdRMlptTmpORFkxWkRGaU9EWTBOR05oTXpjMVpqbGxNVEJpWVdKbFpXRTBZek00TXpFMXVEbGlkSEE2THk4d2VESXVhV052Ymk5amVHSXpOV1JoWVRaalpHWm1ZekZqTXpZME5tUmlOalU0TXpRMVpXUTFaVFV6WlRRd09XSXpNRFU9In0=" }, "type": "data_update" } ] }`),
							Error:   nil,
						},
					},
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
				Success: func(ouput chan *chain.Message) bool {
				msgloop:
					for msg := range ouput {
						for _, j := range msg.Receipts {
							for _, e := range j.Events {
								if e.Next == chain.BTPAddress("btp://0x2.icon/cxb35daa6cdffc1c3646db658345ed5e53e409b305") {
									return false
								} else if e.Next == chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315") {
									break msgloop
								}

							}
						}
					}
					return true
				},
				Fail: nil,
			},
		},
		{
			Description: "With Btp Message sequence higher than sequence provided in config",
			Input: struct {
				Seq         uint64
				Offset      uint64
				Source      chain.BTPAddress
				Destination chain.BTPAddress
			}{
				Seq:         0,
				Offset:      377825,
				Source:      chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Destination: chain.BTPAddress("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315"),
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
					ContractStateChangeMap: map[int64]mock.Response{
						377825: {
							Reponse: []byte(`{ "block_hash": "DDbjZ12VbmV36trcJDPxAAHsDWTtGEC9DB6ZSVLE9N1c", "changes": [ { "cause": { "receipt_hash": "2VWWEfyg5BzyDBVRsHRXApQJdZ37Bdtj8GtkH7UvNm7G", "type": "receipt_processing" }, "change": { "account_id": "alice.node1", "key_base64": "bWVzc2FnZQ==", "value_base64": "uQEAAHsibmV4dCI6ImJ0cDovLzB4Ny5pY29uL2N4MWFkNmZjYzQ2NWQxYjg2NDRjYTM3NWY5ZTEwYmFiZWVhNGMzODMxNSIsInNlcXVlbmNlIjoiNDkiLCJtZXNzYWdlIjoiK1ArNFQySjBjRG92THpCNE1pNXVaV0Z5THpjeU56QmhOemxpWlRjNE9XUTNOekJtTW1SbE1ERTFNRFEzTmpnMFpUSTRNRFkxT1RkbFpXVmxPVFpsWlROallUZzNZakUzT1dNMk16azVaR1ZoWVdhNE9XSjBjRG92THpCNE55NXBZMjl1TDJONE1XRmtObVpqWXpRMk5XUXhZamcyTkRSallUTTNOV1k1WlRFd1ltRmlaV1ZoTkdNek9ETXhOWU5pZEhNZHVHejRhZ0M0Wi9obGptSjBjQzB5Tnk1MFpYTjBibVYwcW1oNFltUXlORGswWmpaa1pUWXlOVEJpT0dRNE16Z3pNV0UxTXpoak0yVmxZbVEwTnpZNVpHVmhadXJwa1dKMGNDMHdlREl1Ym1WaGNpMU9SVUZTaXdDeDRIM0NNVUo5QUFBQWlpSGhuZ3licXlRQUFBQT0ifQ==" }, "type": "data_update" } ] }`),
							Error:   nil,
						},
					},
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
				Success: nil,
				Fail: func(srcMsg chan *chain.Message, output <-chan error, err error) bool {
					result := false

					if err != nil {
						return true
					}
					go func() {
						for {
							select {
							case <-srcMsg:
							case err := <-output:
								if err.Error() == "invalid event seq" {
									result = true
								}
							}
						}
					}()

					time.Sleep(time.Second * 3)
					return result
				},
			},
		},
	}

	RegisterTest("ReceiverSubscribe", ReceiverSubscribe{
		description: "Test Receiver Subscribe",
		testData:    testData,
	})
}
