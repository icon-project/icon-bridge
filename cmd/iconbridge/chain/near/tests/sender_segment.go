package tests

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SenderSegmentTest struct {
	description string
	testData    []TestData
}

func (t SenderSegmentTest) Description() string {
	return t.description
}

func (t SenderSegmentTest) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "Without any receipts",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From:     chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: []*chain.Receipt{},
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "Within Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3988),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f, len(relayMessage) > 3988)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 1 Event above Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 5000),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f, len(relayMessage) > 5000)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 2 Events within Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 1000),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 1000),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 2000)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 2 Events, and 1 Event above Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3000),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3988),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 3000)
					assert.Equal(f, 1, len(nextMessage.Receipts))
					assert.Equal(f, 3988, len(nextMessage.Receipts[0].Events[0].Message))
				},
			},
		},
		{
			Description: "With 2 Receipts, and Events within Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 500),
								},
							},
							Height: uint64(1),
						},
							&chain.Receipt{
								Index: uint64(0),
								Events: []*chain.Event{
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(0),
										Message:  make([]byte, 500),
									},
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(0),
										Message:  make([]byte, 500),
									},
								},
								Height: uint64(1),
							})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 2000)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 2 Receipts, and 1 Receipt have Events above Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(1),
									Message:  make([]byte, 500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(2),
									Message:  make([]byte, 500),
								},
							},
							Height: uint64(1),
						},
							&chain.Receipt{
								Index: uint64(1),
								Events: []*chain.Event{
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(3),
										Message:  make([]byte, 500),
									},
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(4),
										Message:  make([]byte, 4000),
									},
								},
								Height: uint64(1),
							})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 1500)
					assert.Equal(f, 1, len(nextMessage.Receipts))
					assert.Equal(f, 4000, len(nextMessage.Receipts[0].Events[0].Message))
				},
			},
		},
		{
			Description: "With 2 Receipts, and 1 Receipt have 2 Events above Default Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(1),
									Message:  make([]byte, 500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(2),
									Message:  make([]byte, 500),
								},
							},
							Height: uint64(1),
						},
							&chain.Receipt{
								Index: uint64(1),
								Events: []*chain.Event{
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(3),
										Message:  make([]byte, 4000),
									},
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(4),
										Message:  make([]byte, 4000),
									},
								},
								Height: uint64(1),
							})

						return receipts
					}(),
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 1000)
					assert.Equal(f, 1, len(nextMessage.Receipts))
					assert.Equal(f, 4000, len(nextMessage.Receipts[0].Events[0].Message))
					assert.Equal(f, 4000, len(nextMessage.Receipts[0].Events[1].Message))
				},
			},
		},
		{
			Description: "Within Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3988),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 3988)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 1 Event above Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 9000),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 9000)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 2 Events within Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3500),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 7000)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 2 Events, and 1 Event above Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 7000),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 3988),
								},
							},
							Height: uint64(1),
						})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 7000)
					assert.Equal(f, 1, len(nextMessage.Receipts))
					assert.Equal(f, 3988, len(nextMessage.Receipts[0].Events[0].Message))
				},
			},
		},
		{
			Description: "With 2 Receipts, and Events within Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 1500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(0),
									Message:  make([]byte, 1500),
								},
							},
							Height: uint64(1),
						},
							&chain.Receipt{
								Index: uint64(0),
								Events: []*chain.Event{
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(0),
										Message:  make([]byte, 1500),
									},
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(0),
										Message:  make([]byte, 1500),
									},
								},
								Height: uint64(1),
							})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 6000)
					assert.Equal(f, 0, len(nextMessage.Receipts))
				},
			},
		},
		{
			Description: "With 2 Receipts, and 1 Receipt have Events above Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(1),
									Message:  make([]byte, 1500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(2),
									Message:  make([]byte, 1500),
								},
							},
							Height: uint64(1),
						},
							&chain.Receipt{
								Index: uint64(1),
								Events: []*chain.Event{
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(3),
										Message:  make([]byte, 1500),
									},
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(4),
										Message:  make([]byte, 10000),
									},
								},
								Height: uint64(1),
							})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 4500)
					assert.Equal(f, 1, len(nextMessage.Receipts))
					assert.Equal(f, 10000, len(nextMessage.Receipts[0].Events[0].Message))
				},
			},
		},
		{
			Description: "With 2 Receipts, and 1 Receipt have 2 Events above Custom Limit",
			Input: struct {
				PrivateKey  string
				Source      chain.BTPAddress
				Destination chain.BTPAddress
				Message     *chain.Message
				Options     types.SenderOptions
			}{
				PrivateKey:  "22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J",
				Source:      chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
				Destination: chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
				Message: &chain.Message{
					From: chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"),
					Receipts: func() []*chain.Receipt {
						receipts := make([]*chain.Receipt, 0)
						receipts = append(receipts, &chain.Receipt{
							Index: uint64(0),
							Events: []*chain.Event{
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(1),
									Message:  make([]byte, 1500),
								},
								{
									Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
									Sequence: uint64(2),
									Message:  make([]byte, 1500),
								},
							},
							Height: uint64(1),
						},
							&chain.Receipt{
								Index: uint64(1),
								Events: []*chain.Event{
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(3),
										Message:  make([]byte, 10000),
									},
									{
										Next:     chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598"),
										Sequence: uint64(4),
										Message:  make([]byte, 10000),
									},
								},
								Height: uint64(1),
							})

						return receipts
					}(),
				},
				Options: types.SenderOptions{
					TxDataSizeLimit: 8000,
				},
			},
			MockApi: func() *mock.MockApi {
				mockApi := mock.NewMockApi(mock.Storage{})

				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				mockApi.On("ViewAccessKey", mock.MockParam).Return(mockApi.ViewAccessKeyFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: func(f *testing.T, relayMessage []byte, nextMessage *chain.Message) {
					assert.True(f,len(relayMessage) > 3000)
					assert.Equal(f, 1, len(nextMessage.Receipts))
					assert.Equal(f, 10000, len(nextMessage.Receipts[0].Events[0].Message))
					assert.Equal(f, 10000, len(nextMessage.Receipts[0].Events[1].Message))
				},
			},
		},
	}

	RegisterTest("SenderSegment", SenderSegmentTest{
		description: "Test Sender Segment",
		testData:    testData,
	})
}
