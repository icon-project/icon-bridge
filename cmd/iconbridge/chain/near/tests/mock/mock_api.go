package mock

import (
	"github.com/MuhammedIrfan/testify-mock/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

const (
	MockParam = "mock.Anything"
)
var emptyResponse Response

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

func (m *MockApi) Chunk(param interface{}) (types.ChunkHeader, error) {
	ret := m.Called(param)

	var r0 types.ChunkHeader
	if rf, ok := ret.Get(0).(func(interface{}) types.ChunkHeader); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(types.ChunkHeader)
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

	for key, value := range storage.AccessKeyMap {
		defaults.AccessKeyMap[key] = value
	}
	storage.AccessKeyMap = defaults.AccessKeyMap

	for key, value := range storage.BmcLinkStatusMap {
		defaults.BmcLinkStatusMap[key] = value
	}
	storage.BmcLinkStatusMap = defaults.BmcLinkStatusMap

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

	return mockApi
}
