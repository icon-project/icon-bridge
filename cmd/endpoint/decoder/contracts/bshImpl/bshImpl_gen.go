// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bshImpl

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TypesAsset is an auto generated low-level Go binding around an user-defined struct.
type TypesAsset struct {
	Name  string
	Value *big.Int
	Fee   *big.Int
}

// TypesTransferAssets is an auto generated low-level Go binding around an user-defined struct.
type TypesTransferAssets struct {
	From  string
	To    string
	Asset []TypesAsset
}

// BshImplABI is the input ABI used to generate the binding from.
const BshImplABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"HandleBTPMessageEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"}],\"name\":\"ResponseUnknownType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_response\",\"type\":\"string\"}],\"name\":\"TransferEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"TransferStart\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bmc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bshProxy\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_serviceName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_toFA\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"}],\"name\":\"handleFeeGathering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"handleBTPMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"asset\",\"type\":\"tuple[]\"}],\"internalType\":\"structTypes.TransferAssets\",\"name\":\"transferAssets\",\"type\":\"tuple\"}],\"name\":\"handleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"checkParseAddress\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_src\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"handleBTPError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"sendServiceMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// BshImpl is an auto generated Go binding around an Ethereum contract.
type BshImpl struct {
	BshImplCaller     // Read-only binding to the contract
	BshImplTransactor // Write-only binding to the contract
	BshImplFilterer   // Log filterer for contract events
}

// BshImplCaller is an auto generated read-only Go binding around an Ethereum contract.
type BshImplCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshImplTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BshImplTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshImplFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BshImplFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshImplSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BshImplSession struct {
	Contract     *BshImpl          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BshImplCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BshImplCallerSession struct {
	Contract *BshImplCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BshImplTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BshImplTransactorSession struct {
	Contract     *BshImplTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BshImplRaw is an auto generated low-level Go binding around an Ethereum contract.
type BshImplRaw struct {
	Contract *BshImpl // Generic contract binding to access the raw methods on
}

// BshImplCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BshImplCallerRaw struct {
	Contract *BshImplCaller // Generic read-only contract binding to access the raw methods on
}

// BshImplTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BshImplTransactorRaw struct {
	Contract *BshImplTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBshImpl creates a new instance of BshImpl, bound to a specific deployed contract.
func NewBshImpl(address common.Address, backend bind.ContractBackend) (*BshImpl, error) {
	contract, err := bindBshImpl(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BshImpl{BshImplCaller: BshImplCaller{contract: contract}, BshImplTransactor: BshImplTransactor{contract: contract}, BshImplFilterer: BshImplFilterer{contract: contract}}, nil
}

// NewBshImplCaller creates a new read-only instance of BshImpl, bound to a specific deployed contract.
func NewBshImplCaller(address common.Address, caller bind.ContractCaller) (*BshImplCaller, error) {
	contract, err := bindBshImpl(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BshImplCaller{contract: contract}, nil
}

// NewBshImplTransactor creates a new write-only instance of BshImpl, bound to a specific deployed contract.
func NewBshImplTransactor(address common.Address, transactor bind.ContractTransactor) (*BshImplTransactor, error) {
	contract, err := bindBshImpl(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BshImplTransactor{contract: contract}, nil
}

// NewBshImplFilterer creates a new log filterer instance of BshImpl, bound to a specific deployed contract.
func NewBshImplFilterer(address common.Address, filterer bind.ContractFilterer) (*BshImplFilterer, error) {
	contract, err := bindBshImpl(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BshImplFilterer{contract: contract}, nil
}

// bindBshImpl binds a generic wrapper to an already deployed contract.
func bindBshImpl(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BshImplABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BshImpl *BshImplRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BshImpl.Contract.BshImplCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BshImpl *BshImplRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BshImpl.Contract.BshImplTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BshImpl *BshImplRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BshImpl.Contract.BshImplTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BshImpl *BshImplCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BshImpl.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BshImpl *BshImplTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BshImpl.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BshImpl *BshImplTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BshImpl.Contract.contract.Transact(opts, method, params...)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_BshImpl *BshImplCaller) CheckParseAddress(opts *bind.CallOpts, _to string) error {
	var out []interface{}
	err := _BshImpl.contract.Call(opts, &out, "checkParseAddress", _to)

	if err != nil {
		return err
	}

	return err

}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_BshImpl *BshImplSession) CheckParseAddress(_to string) error {
	return _BshImpl.Contract.CheckParseAddress(&_BshImpl.CallOpts, _to)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_BshImpl *BshImplCallerSession) CheckParseAddress(_to string) error {
	return _BshImpl.Contract.CheckParseAddress(&_BshImpl.CallOpts, _to)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_BshImpl *BshImplCaller) HasPendingRequest(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BshImpl.contract.Call(opts, &out, "hasPendingRequest")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_BshImpl *BshImplSession) HasPendingRequest() (bool, error) {
	return _BshImpl.Contract.HasPendingRequest(&_BshImpl.CallOpts)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_BshImpl *BshImplCallerSession) HasPendingRequest() (bool, error) {
	return _BshImpl.Contract.HasPendingRequest(&_BshImpl.CallOpts)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_BshImpl *BshImplCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	var out []interface{}
	err := _BshImpl.contract.Call(opts, &out, "requests", arg0)

	outstruct := new(struct {
		From string
		To   string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.From = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.To = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_BshImpl *BshImplSession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _BshImpl.Contract.Requests(&_BshImpl.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_BshImpl *BshImplCallerSession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _BshImpl.Contract.Requests(&_BshImpl.CallOpts, arg0)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_BshImpl *BshImplCaller) ServiceName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BshImpl.contract.Call(opts, &out, "serviceName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_BshImpl *BshImplSession) ServiceName() (string, error) {
	return _BshImpl.Contract.ServiceName(&_BshImpl.CallOpts)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_BshImpl *BshImplCallerSession) ServiceName() (string, error) {
	return _BshImpl.Contract.ServiceName(&_BshImpl.CallOpts)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_BshImpl *BshImplTransactor) HandleBTPError(opts *bind.TransactOpts, _src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _BshImpl.contract.Transact(opts, "handleBTPError", _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_BshImpl *BshImplSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleBTPError(&_BshImpl.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_BshImpl *BshImplTransactorSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleBTPError(&_BshImpl.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_BshImpl *BshImplTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BshImpl.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_BshImpl *BshImplSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleBTPMessage(&_BshImpl.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_BshImpl *BshImplTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleBTPMessage(&_BshImpl.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _toFA, string _svc) returns()
func (_BshImpl *BshImplTransactor) HandleFeeGathering(opts *bind.TransactOpts, _toFA string, _svc string) (*types.Transaction, error) {
	return _BshImpl.contract.Transact(opts, "handleFeeGathering", _toFA, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _toFA, string _svc) returns()
func (_BshImpl *BshImplSession) HandleFeeGathering(_toFA string, _svc string) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleFeeGathering(&_BshImpl.TransactOpts, _toFA, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _toFA, string _svc) returns()
func (_BshImpl *BshImplTransactorSession) HandleFeeGathering(_toFA string, _svc string) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleFeeGathering(&_BshImpl.TransactOpts, _toFA, _svc)
}

// HandleRequest is a paid mutator transaction binding the contract method 0x898b83e7.
//
// Solidity: function handleRequest((string,string,(string,uint256,uint256)[]) transferAssets) returns()
func (_BshImpl *BshImplTransactor) HandleRequest(opts *bind.TransactOpts, transferAssets TypesTransferAssets) (*types.Transaction, error) {
	return _BshImpl.contract.Transact(opts, "handleRequest", transferAssets)
}

// HandleRequest is a paid mutator transaction binding the contract method 0x898b83e7.
//
// Solidity: function handleRequest((string,string,(string,uint256,uint256)[]) transferAssets) returns()
func (_BshImpl *BshImplSession) HandleRequest(transferAssets TypesTransferAssets) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleRequest(&_BshImpl.TransactOpts, transferAssets)
}

// HandleRequest is a paid mutator transaction binding the contract method 0x898b83e7.
//
// Solidity: function handleRequest((string,string,(string,uint256,uint256)[]) transferAssets) returns()
func (_BshImpl *BshImplTransactorSession) HandleRequest(transferAssets TypesTransferAssets) (*types.Transaction, error) {
	return _BshImpl.Contract.HandleRequest(&_BshImpl.TransactOpts, transferAssets)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshProxy, string _serviceName) returns()
func (_BshImpl *BshImplTransactor) Initialize(opts *bind.TransactOpts, _bmc common.Address, _bshProxy common.Address, _serviceName string) (*types.Transaction, error) {
	return _BshImpl.contract.Transact(opts, "initialize", _bmc, _bshProxy, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshProxy, string _serviceName) returns()
func (_BshImpl *BshImplSession) Initialize(_bmc common.Address, _bshProxy common.Address, _serviceName string) (*types.Transaction, error) {
	return _BshImpl.Contract.Initialize(&_BshImpl.TransactOpts, _bmc, _bshProxy, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshProxy, string _serviceName) returns()
func (_BshImpl *BshImplTransactorSession) Initialize(_bmc common.Address, _bshProxy common.Address, _serviceName string) (*types.Transaction, error) {
	return _BshImpl.Contract.Initialize(&_BshImpl.TransactOpts, _bmc, _bshProxy, _serviceName)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0x5c436dbe.
//
// Solidity: function sendServiceMessage(address _from, string _to, (string,uint256,uint256)[] _assets) returns()
func (_BshImpl *BshImplTransactor) SendServiceMessage(opts *bind.TransactOpts, _from common.Address, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _BshImpl.contract.Transact(opts, "sendServiceMessage", _from, _to, _assets)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0x5c436dbe.
//
// Solidity: function sendServiceMessage(address _from, string _to, (string,uint256,uint256)[] _assets) returns()
func (_BshImpl *BshImplSession) SendServiceMessage(_from common.Address, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _BshImpl.Contract.SendServiceMessage(&_BshImpl.TransactOpts, _from, _to, _assets)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0x5c436dbe.
//
// Solidity: function sendServiceMessage(address _from, string _to, (string,uint256,uint256)[] _assets) returns()
func (_BshImpl *BshImplTransactorSession) SendServiceMessage(_from common.Address, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _BshImpl.Contract.SendServiceMessage(&_BshImpl.TransactOpts, _from, _to, _assets)
}

// BshImplHandleBTPMessageEventIterator is returned from FilterHandleBTPMessageEvent and is used to iterate over the raw logs and unpacked data for HandleBTPMessageEvent events raised by the BshImpl contract.
type BshImplHandleBTPMessageEventIterator struct {
	Event *BshImplHandleBTPMessageEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BshImplHandleBTPMessageEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshImplHandleBTPMessageEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BshImplHandleBTPMessageEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BshImplHandleBTPMessageEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshImplHandleBTPMessageEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshImplHandleBTPMessageEvent represents a HandleBTPMessageEvent event raised by the BshImpl contract.
type BshImplHandleBTPMessageEvent struct {
	Sn   *big.Int
	Code *big.Int
	Msg  string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterHandleBTPMessageEvent is a free log retrieval operation binding the contract event 0x356868e4a05430bccb6aa9c954e410ab0792c5a5baa7b973b03e1d4c03fa1366.
//
// Solidity: event HandleBTPMessageEvent(uint256 _sn, uint256 _code, string _msg)
func (_BshImpl *BshImplFilterer) FilterHandleBTPMessageEvent(opts *bind.FilterOpts) (*BshImplHandleBTPMessageEventIterator, error) {

	logs, sub, err := _BshImpl.contract.FilterLogs(opts, "HandleBTPMessageEvent")
	if err != nil {
		return nil, err
	}
	return &BshImplHandleBTPMessageEventIterator{contract: _BshImpl.contract, event: "HandleBTPMessageEvent", logs: logs, sub: sub}, nil
}

// WatchHandleBTPMessageEvent is a free log subscription operation binding the contract event 0x356868e4a05430bccb6aa9c954e410ab0792c5a5baa7b973b03e1d4c03fa1366.
//
// Solidity: event HandleBTPMessageEvent(uint256 _sn, uint256 _code, string _msg)
func (_BshImpl *BshImplFilterer) WatchHandleBTPMessageEvent(opts *bind.WatchOpts, sink chan<- *BshImplHandleBTPMessageEvent) (event.Subscription, error) {

	logs, sub, err := _BshImpl.contract.WatchLogs(opts, "HandleBTPMessageEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshImplHandleBTPMessageEvent)
				if err := _BshImpl.contract.UnpackLog(event, "HandleBTPMessageEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHandleBTPMessageEvent is a log parse operation binding the contract event 0x356868e4a05430bccb6aa9c954e410ab0792c5a5baa7b973b03e1d4c03fa1366.
//
// Solidity: event HandleBTPMessageEvent(uint256 _sn, uint256 _code, string _msg)
func (_BshImpl *BshImplFilterer) ParseHandleBTPMessageEvent(log types.Log) (*BshImplHandleBTPMessageEvent, error) {
	event := new(BshImplHandleBTPMessageEvent)
	if err := _BshImpl.contract.UnpackLog(event, "HandleBTPMessageEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshImplResponseUnknownTypeIterator is returned from FilterResponseUnknownType and is used to iterate over the raw logs and unpacked data for ResponseUnknownType events raised by the BshImpl contract.
type BshImplResponseUnknownTypeIterator struct {
	Event *BshImplResponseUnknownType // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BshImplResponseUnknownTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshImplResponseUnknownType)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BshImplResponseUnknownType)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BshImplResponseUnknownTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshImplResponseUnknownTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshImplResponseUnknownType represents a ResponseUnknownType event raised by the BshImpl contract.
type BshImplResponseUnknownType struct {
	From string
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResponseUnknownType is a free log retrieval operation binding the contract event 0x64f88365dae9c547bfcae6a186a7827fdc2613baffd8d5164dc59a74f55fbeba.
//
// Solidity: event ResponseUnknownType(string _from, uint256 _sn)
func (_BshImpl *BshImplFilterer) FilterResponseUnknownType(opts *bind.FilterOpts) (*BshImplResponseUnknownTypeIterator, error) {

	logs, sub, err := _BshImpl.contract.FilterLogs(opts, "ResponseUnknownType")
	if err != nil {
		return nil, err
	}
	return &BshImplResponseUnknownTypeIterator{contract: _BshImpl.contract, event: "ResponseUnknownType", logs: logs, sub: sub}, nil
}

// WatchResponseUnknownType is a free log subscription operation binding the contract event 0x64f88365dae9c547bfcae6a186a7827fdc2613baffd8d5164dc59a74f55fbeba.
//
// Solidity: event ResponseUnknownType(string _from, uint256 _sn)
func (_BshImpl *BshImplFilterer) WatchResponseUnknownType(opts *bind.WatchOpts, sink chan<- *BshImplResponseUnknownType) (event.Subscription, error) {

	logs, sub, err := _BshImpl.contract.WatchLogs(opts, "ResponseUnknownType")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshImplResponseUnknownType)
				if err := _BshImpl.contract.UnpackLog(event, "ResponseUnknownType", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseResponseUnknownType is a log parse operation binding the contract event 0x64f88365dae9c547bfcae6a186a7827fdc2613baffd8d5164dc59a74f55fbeba.
//
// Solidity: event ResponseUnknownType(string _from, uint256 _sn)
func (_BshImpl *BshImplFilterer) ParseResponseUnknownType(log types.Log) (*BshImplResponseUnknownType, error) {
	event := new(BshImplResponseUnknownType)
	if err := _BshImpl.contract.UnpackLog(event, "ResponseUnknownType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshImplTransferEndIterator is returned from FilterTransferEnd and is used to iterate over the raw logs and unpacked data for TransferEnd events raised by the BshImpl contract.
type BshImplTransferEndIterator struct {
	Event *BshImplTransferEnd // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BshImplTransferEndIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshImplTransferEnd)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BshImplTransferEnd)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BshImplTransferEndIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshImplTransferEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshImplTransferEnd represents a TransferEnd event raised by the BshImpl contract.
type BshImplTransferEnd struct {
	From     common.Address
	Sn       *big.Int
	Code     *big.Int
	Response string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferEnd is a free log retrieval operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_BshImpl *BshImplFilterer) FilterTransferEnd(opts *bind.FilterOpts, _from []common.Address) (*BshImplTransferEndIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshImpl.contract.FilterLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return &BshImplTransferEndIterator{contract: _BshImpl.contract, event: "TransferEnd", logs: logs, sub: sub}, nil
}

// WatchTransferEnd is a free log subscription operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_BshImpl *BshImplFilterer) WatchTransferEnd(opts *bind.WatchOpts, sink chan<- *BshImplTransferEnd, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshImpl.contract.WatchLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshImplTransferEnd)
				if err := _BshImpl.contract.UnpackLog(event, "TransferEnd", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferEnd is a log parse operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_BshImpl *BshImplFilterer) ParseTransferEnd(log types.Log) (*BshImplTransferEnd, error) {
	event := new(BshImplTransferEnd)
	if err := _BshImpl.contract.UnpackLog(event, "TransferEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshImplTransferReceivedIterator is returned from FilterTransferReceived and is used to iterate over the raw logs and unpacked data for TransferReceived events raised by the BshImpl contract.
type BshImplTransferReceivedIterator struct {
	Event *BshImplTransferReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BshImplTransferReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshImplTransferReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BshImplTransferReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BshImplTransferReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshImplTransferReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshImplTransferReceived represents a TransferReceived event raised by the BshImpl contract.
type BshImplTransferReceived struct {
	From         common.Hash
	To           common.Address
	Sn           *big.Int
	AssetDetails []TypesAsset
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferReceived is a free log retrieval operation binding the contract event 0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_BshImpl *BshImplFilterer) FilterTransferReceived(opts *bind.FilterOpts, _from []string, _to []common.Address) (*BshImplTransferReceivedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _BshImpl.contract.FilterLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &BshImplTransferReceivedIterator{contract: _BshImpl.contract, event: "TransferReceived", logs: logs, sub: sub}, nil
}

// WatchTransferReceived is a free log subscription operation binding the contract event 0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_BshImpl *BshImplFilterer) WatchTransferReceived(opts *bind.WatchOpts, sink chan<- *BshImplTransferReceived, _from []string, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _BshImpl.contract.WatchLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshImplTransferReceived)
				if err := _BshImpl.contract.UnpackLog(event, "TransferReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferReceived is a log parse operation binding the contract event 0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_BshImpl *BshImplFilterer) ParseTransferReceived(log types.Log) (*BshImplTransferReceived, error) {
	event := new(BshImplTransferReceived)
	if err := _BshImpl.contract.UnpackLog(event, "TransferReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshImplTransferStartIterator is returned from FilterTransferStart and is used to iterate over the raw logs and unpacked data for TransferStart events raised by the BshImpl contract.
type BshImplTransferStartIterator struct {
	Event *BshImplTransferStart // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BshImplTransferStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshImplTransferStart)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BshImplTransferStart)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BshImplTransferStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshImplTransferStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshImplTransferStart represents a TransferStart event raised by the BshImpl contract.
type BshImplTransferStart struct {
	From   common.Address
	To     string
	Sn     *big.Int
	Assets []TypesAsset
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferStart is a free log retrieval operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assets)
func (_BshImpl *BshImplFilterer) FilterTransferStart(opts *bind.FilterOpts, _from []common.Address) (*BshImplTransferStartIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshImpl.contract.FilterLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return &BshImplTransferStartIterator{contract: _BshImpl.contract, event: "TransferStart", logs: logs, sub: sub}, nil
}

// WatchTransferStart is a free log subscription operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assets)
func (_BshImpl *BshImplFilterer) WatchTransferStart(opts *bind.WatchOpts, sink chan<- *BshImplTransferStart, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshImpl.contract.WatchLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshImplTransferStart)
				if err := _BshImpl.contract.UnpackLog(event, "TransferStart", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferStart is a log parse operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assets)
func (_BshImpl *BshImplFilterer) ParseTransferStart(log types.Log) (*BshImplTransferStart, error) {
	event := new(BshImplTransferStart)
	if err := _BshImpl.contract.UnpackLog(event, "TransferStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
