package mock

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/MuhammedIrfan/testify-mock/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type MockApi struct {
	mock.Mock
	*Storage
}

func (m *MockApi) Block(param interface{}) (types.Block, error) {
	ret := m.Called(param)

	var r0 types.Block
	if rf, ok := ret.Get(0).(func(interface{}) types.Block); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.Block)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) BroadcastTxAsync(param interface{}) (types.CryptoHash, error) {
	ret := m.Called(param)

	var r0 types.CryptoHash
	if rf, ok := ret.Get(0).(func(interface{}) types.CryptoHash); ok {
		r0 = rf(param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.CryptoHash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) BroadcastTxCommit(param interface{}) (types.TransactionResult, error) {
	ret := m.Called(param)

	var r0 types.TransactionResult
	if rf, ok := ret.Get(0).(func(interface{}) types.TransactionResult); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.TransactionResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) CallFunction(param interface{}) (types.CallFunctionResponse, error) {
	ret := m.Called(param)

	var r0 types.CallFunctionResponse
	if rf, ok := ret.Get(0).(func(interface{}) types.CallFunctionResponse); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.CallFunctionResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) Changes(param interface{}) (types.ContractStateChange, error) {
	ret := m.Called(param)

	var r0 types.ContractStateChange
	if rf, ok := ret.Get(0).(func(interface{}) types.ContractStateChange); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.ContractStateChange)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) LightClientProof(param interface{}) (types.ReceiptProof, error) {
	ret := m.Called(param)

	var r0 types.ReceiptProof
	if rf, ok := ret.Get(0).(func(interface{}) types.ReceiptProof); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.ReceiptProof)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) Status(param interface{}) (types.ChainStatus, error) {
	ret := m.Called(param)

	var r0 types.ChainStatus
	if rf, ok := ret.Get(0).(func(interface{}) types.ChainStatus); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.ChainStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) Transaction(param interface{}) (types.TransactionResult, error) {
	ret := m.Called(param)

	var r0 types.TransactionResult
	if rf, ok := ret.Get(0).(func(interface{}) types.TransactionResult); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.TransactionResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) ViewAccessKey(param interface{}) (types.AccessKeyResponse, error) {
	ret := m.Called(param)

	var r0 types.AccessKeyResponse
	if rf, ok := ret.Get(0).(func(interface{}) types.AccessKeyResponse); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.AccessKeyResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockApi) ViewAccount(param interface{}) (types.Account, error) {
	ret := m.Called(param)

	var r0 types.Account
	if rf, ok := ret.Get(0).(func(interface{}) types.Account); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func NewMockApi(storage Storage) *MockApi {
	var emptyResponse Response
	var defaults = Default()

	if storage.LatestChainStatus == emptyResponse {
		storage.LatestChainStatus = defaults.LatestChainStatus
	}

	if storage.TransactionHash == emptyResponse {
		storage.TransactionHash = defaults.TransactionHash
	}

	for key, value := range storage.TransactionResultMap {
		defaults.TransactionResultMap[key] = value
	}
	storage.TransactionResultMap = defaults.TransactionResultMap

	for key, value := range storage.BlockByHeightMap {
		defaults.BlockByHeightMap[key] = value
	}
	storage.BlockByHeightMap = defaults.BlockByHeightMap

	for key, value := range storage.BlockByHashMap {
		defaults.BlockByHashMap[key] = value
	}
	storage.BlockByHashMap = defaults.BlockByHashMap

	for key, value := range storage.NonceMap {
		defaults.NonceMap[key] = value
	}
	storage.NonceMap = defaults.NonceMap

	for key, value := range storage.BmcLinkStatusMap {
		defaults.BmcLinkStatusMap[key] = value
	}
	storage.BmcLinkStatusMap = defaults.BmcLinkStatusMap

	for key, value := range storage.BmvStatusMap {
		defaults.BmvStatusMap[key] = value
	}
	storage.BmvStatusMap = defaults.BmvStatusMap

	for key, value := range storage.ContractStateChangeMap {
		defaults.ContractStateChangeMap[key] = value
	}
	storage.ContractStateChangeMap = defaults.ContractStateChangeMap

	for key, value := range storage.AccountMap {
		defaults.AccountMap[key] = value
	}
	storage.AccountMap = defaults.AccountMap

	mockApi := &MockApi{
		Storage: &storage,
	}

	mockApi.On("Block", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		var block types.Block

		blockHashParam, isBlockHashParam := args.Get(0).(struct {
			BlockId string `json:"block_id"`
		})

		blockHeightParam, isBlockHeightParam := args.Get(0).(struct {
			BlockId int64 `json:"block_id"`
		})

		if isBlockHashParam && storage.BlockByHashMap[blockHashParam.BlockId] != emptyResponse {
			if storage.BlockByHashMap[blockHashParam.BlockId].Error != nil {
				return []interface{}{block, storage.BlockByHashMap[blockHashParam.BlockId].Error}
			}

			if response, Ok := (storage.BlockByHashMap[blockHashParam.BlockId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &block)
				if err != nil {
					return []interface{}{block, err}
				}

				return []interface{}{block, nil}
			}

		}

		if isBlockHeightParam && storage.BlockByHeightMap[blockHeightParam.BlockId] != emptyResponse {
			if storage.BlockByHeightMap[blockHeightParam.BlockId].Error != nil {
				return []interface{}{block, storage.BlockByHeightMap[blockHeightParam.BlockId].Error}
			}
			if response, Ok := (storage.BlockByHeightMap[blockHeightParam.BlockId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &block)
				if err != nil {
					return []interface{}{block, err}
				}
				return []interface{}{block, nil}
			}

		}

		return []interface{}{block, errors.New("invalid Param")}
	})

	//BroadcastTxAsync
	mockApi.On("BroadcastTxAsync", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		return []interface{}{nil}
	})

	//BroadcastTxCommit
	mockApi.On("BroadcastTxCommit", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		return []interface{}{nil}
	})

	mockApi.On("CallFunction", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		var response types.CallFunctionResponse

		param, Ok := args.Get(0).(types.CallFunction)

		if Ok && param.MethodName == "get_status" {
			var linkParam struct {
				Link chain.BTPAddress `json:"link"`
			}
			data, err := base64.URLEncoding.DecodeString(param.ArgumentsB64)
			if err != nil {
				return []interface{}{response, err}
			}

			err = json.Unmarshal(data, &linkParam)
			if err != nil {
				return []interface{}{response, err}
			}

			if storage.BmcLinkStatusMap[linkParam.Link.ContractAddress()] != emptyResponse {
				if storage.BmcLinkStatusMap[linkParam.Link.ContractAddress()].Error != nil {
					return []interface{}{response, storage.BmcLinkStatusMap[linkParam.Link.ContractAddress()].Error}
				}

				if data, Ok := (storage.BmcLinkStatusMap[linkParam.Link.ContractAddress()].Reponse).([]byte); Ok {
					err := json.Unmarshal(data, &response)
					if err != nil {
						return []interface{}{response, err}
					}
					return []interface{}{response, nil}
				}
			}

		}

		return []interface{}{response, errors.New("invalid Param")}
	})

	mockApi.On("Changes", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		var changes types.ContractStateChange

		param, Ok := args.Get(0).(struct {
			ChangeType string   `json:"changes_type"`
			AccountIds []string `json:"account_ids"`
			KeyPrefix  string   `json:"key_prefix_base64"`
			BlockId    int64    `json:"block_id"`
		})

		if Ok && storage.ContractStateChangeMap[param.BlockId] != emptyResponse {
			if storage.ContractStateChangeMap[param.BlockId].Error != nil {
				return []interface{}{changes, storage.ContractStateChangeMap[param.BlockId].Error}
			}

			if response, Ok := (storage.ContractStateChangeMap[param.BlockId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &changes)
				if err != nil {
					return []interface{}{changes, err}
				}

				return []interface{}{changes, nil}
			}
		}

		return []interface{}{changes, errors.New("invalid Param")}
	})

	mockApi.On("LightClientProof", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		return []interface{}{nil}
	})

	mockApi.On("Status", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		var chainStatus types.ChainStatus

		if storage.LatestChainStatus != emptyResponse {
			if storage.LatestChainStatus.Error != nil {
				return []interface{}{chainStatus, storage.LatestChainStatus.Error}
			}

			if response, Ok := (storage.LatestChainStatus.Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &chainStatus)
				if err != nil {
					return []interface{}{chainStatus, err}
				}
				return []interface{}{chainStatus, nil}
			}
		}

		return []interface{}{chainStatus, errors.New("invalid Param")}
	})

	//Transaction
	mockApi.On("Transaction", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		return []interface{}{nil}
	})

	//ViewAccessKey
	mockApi.On("ViewAccessKey", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		return []interface{}{nil}
	})

	//ViewAccount
	mockApi.On("ViewAccount", mock.Anything).Return(func(args mock.Arguments) mock.Arguments {
		return []interface{}{nil}
	})

	return mockApi
}
