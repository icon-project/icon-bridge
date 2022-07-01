// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokenHmy

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

// TokenHmyABI is the input ABI used to generate the binding from.
const TokenHmyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"HandleBTPMessageEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"}],\"name\":\"ResponseUnknownType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_response\",\"type\":\"string\"}],\"name\":\"TransferEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"TransferStart\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bmc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bshProxy\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_serviceName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_toFA\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"}],\"name\":\"handleFeeGathering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"handleBTPMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"asset\",\"type\":\"tuple[]\"}],\"internalType\":\"structTypes.TransferAssets\",\"name\":\"transferAssets\",\"type\":\"tuple\"}],\"name\":\"handleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"checkParseAddress\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_src\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"handleBTPError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"sendServiceMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenHmy is an auto generated Go binding around an Ethereum contract.
type TokenHmy struct {
	TokenHmyCaller     // Read-only binding to the contract
	TokenHmyTransactor // Write-only binding to the contract
	TokenHmyFilterer   // Log filterer for contract events
}

// TokenHmyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenHmyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenHmyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenHmyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenHmyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenHmyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenHmySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenHmySession struct {
	Contract     *TokenHmy         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenHmyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenHmyCallerSession struct {
	Contract *TokenHmyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TokenHmyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenHmyTransactorSession struct {
	Contract     *TokenHmyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TokenHmyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenHmyRaw struct {
	Contract *TokenHmy // Generic contract binding to access the raw methods on
}

// TokenHmyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenHmyCallerRaw struct {
	Contract *TokenHmyCaller // Generic read-only contract binding to access the raw methods on
}

// TokenHmyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenHmyTransactorRaw struct {
	Contract *TokenHmyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenHmy creates a new instance of TokenHmy, bound to a specific deployed contract.
func NewTokenHmy(address common.Address, backend bind.ContractBackend) (*TokenHmy, error) {
	contract, err := bindTokenHmy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenHmy{TokenHmyCaller: TokenHmyCaller{contract: contract}, TokenHmyTransactor: TokenHmyTransactor{contract: contract}, TokenHmyFilterer: TokenHmyFilterer{contract: contract}}, nil
}

// NewTokenHmyCaller creates a new read-only instance of TokenHmy, bound to a specific deployed contract.
func NewTokenHmyCaller(address common.Address, caller bind.ContractCaller) (*TokenHmyCaller, error) {
	contract, err := bindTokenHmy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenHmyCaller{contract: contract}, nil
}

// NewTokenHmyTransactor creates a new write-only instance of TokenHmy, bound to a specific deployed contract.
func NewTokenHmyTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenHmyTransactor, error) {
	contract, err := bindTokenHmy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenHmyTransactor{contract: contract}, nil
}

// NewTokenHmyFilterer creates a new log filterer instance of TokenHmy, bound to a specific deployed contract.
func NewTokenHmyFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenHmyFilterer, error) {
	contract, err := bindTokenHmy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenHmyFilterer{contract: contract}, nil
}

// bindTokenHmy binds a generic wrapper to an already deployed contract.
func bindTokenHmy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenHmyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenHmy *TokenHmyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenHmy.Contract.TokenHmyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenHmy *TokenHmyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenHmy.Contract.TokenHmyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenHmy *TokenHmyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenHmy.Contract.TokenHmyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenHmy *TokenHmyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenHmy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenHmy *TokenHmyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenHmy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenHmy *TokenHmyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenHmy.Contract.contract.Transact(opts, method, params...)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_TokenHmy *TokenHmyCaller) CheckParseAddress(opts *bind.CallOpts, _to string) error {
	var out []interface{}
	err := _TokenHmy.contract.Call(opts, &out, "checkParseAddress", _to)

	if err != nil {
		return err
	}

	return err

}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_TokenHmy *TokenHmySession) CheckParseAddress(_to string) error {
	return _TokenHmy.Contract.CheckParseAddress(&_TokenHmy.CallOpts, _to)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_TokenHmy *TokenHmyCallerSession) CheckParseAddress(_to string) error {
	return _TokenHmy.Contract.CheckParseAddress(&_TokenHmy.CallOpts, _to)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_TokenHmy *TokenHmyCaller) HasPendingRequest(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenHmy.contract.Call(opts, &out, "hasPendingRequest")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_TokenHmy *TokenHmySession) HasPendingRequest() (bool, error) {
	return _TokenHmy.Contract.HasPendingRequest(&_TokenHmy.CallOpts)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_TokenHmy *TokenHmyCallerSession) HasPendingRequest() (bool, error) {
	return _TokenHmy.Contract.HasPendingRequest(&_TokenHmy.CallOpts)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_TokenHmy *TokenHmyCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	var out []interface{}
	err := _TokenHmy.contract.Call(opts, &out, "requests", arg0)

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
func (_TokenHmy *TokenHmySession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _TokenHmy.Contract.Requests(&_TokenHmy.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_TokenHmy *TokenHmyCallerSession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _TokenHmy.Contract.Requests(&_TokenHmy.CallOpts, arg0)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_TokenHmy *TokenHmyCaller) ServiceName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenHmy.contract.Call(opts, &out, "serviceName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_TokenHmy *TokenHmySession) ServiceName() (string, error) {
	return _TokenHmy.Contract.ServiceName(&_TokenHmy.CallOpts)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_TokenHmy *TokenHmyCallerSession) ServiceName() (string, error) {
	return _TokenHmy.Contract.ServiceName(&_TokenHmy.CallOpts)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_TokenHmy *TokenHmyTransactor) HandleBTPError(opts *bind.TransactOpts, _src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _TokenHmy.contract.Transact(opts, "handleBTPError", _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_TokenHmy *TokenHmySession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleBTPError(&_TokenHmy.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string _src, string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_TokenHmy *TokenHmyTransactorSession) HandleBTPError(_src string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleBTPError(&_TokenHmy.TransactOpts, _src, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_TokenHmy *TokenHmyTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _TokenHmy.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_TokenHmy *TokenHmySession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleBTPMessage(&_TokenHmy.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_TokenHmy *TokenHmyTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleBTPMessage(&_TokenHmy.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _toFA, string _svc) returns()
func (_TokenHmy *TokenHmyTransactor) HandleFeeGathering(opts *bind.TransactOpts, _toFA string, _svc string) (*types.Transaction, error) {
	return _TokenHmy.contract.Transact(opts, "handleFeeGathering", _toFA, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _toFA, string _svc) returns()
func (_TokenHmy *TokenHmySession) HandleFeeGathering(_toFA string, _svc string) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleFeeGathering(&_TokenHmy.TransactOpts, _toFA, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _toFA, string _svc) returns()
func (_TokenHmy *TokenHmyTransactorSession) HandleFeeGathering(_toFA string, _svc string) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleFeeGathering(&_TokenHmy.TransactOpts, _toFA, _svc)
}

// HandleRequest is a paid mutator transaction binding the contract method 0x898b83e7.
//
// Solidity: function handleRequest((string,string,(string,uint256,uint256)[]) transferAssets) returns()
func (_TokenHmy *TokenHmyTransactor) HandleRequest(opts *bind.TransactOpts, transferAssets TypesTransferAssets) (*types.Transaction, error) {
	return _TokenHmy.contract.Transact(opts, "handleRequest", transferAssets)
}

// HandleRequest is a paid mutator transaction binding the contract method 0x898b83e7.
//
// Solidity: function handleRequest((string,string,(string,uint256,uint256)[]) transferAssets) returns()
func (_TokenHmy *TokenHmySession) HandleRequest(transferAssets TypesTransferAssets) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleRequest(&_TokenHmy.TransactOpts, transferAssets)
}

// HandleRequest is a paid mutator transaction binding the contract method 0x898b83e7.
//
// Solidity: function handleRequest((string,string,(string,uint256,uint256)[]) transferAssets) returns()
func (_TokenHmy *TokenHmyTransactorSession) HandleRequest(transferAssets TypesTransferAssets) (*types.Transaction, error) {
	return _TokenHmy.Contract.HandleRequest(&_TokenHmy.TransactOpts, transferAssets)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshProxy, string _serviceName) returns()
func (_TokenHmy *TokenHmyTransactor) Initialize(opts *bind.TransactOpts, _bmc common.Address, _bshProxy common.Address, _serviceName string) (*types.Transaction, error) {
	return _TokenHmy.contract.Transact(opts, "initialize", _bmc, _bshProxy, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshProxy, string _serviceName) returns()
func (_TokenHmy *TokenHmySession) Initialize(_bmc common.Address, _bshProxy common.Address, _serviceName string) (*types.Transaction, error) {
	return _TokenHmy.Contract.Initialize(&_TokenHmy.TransactOpts, _bmc, _bshProxy, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshProxy, string _serviceName) returns()
func (_TokenHmy *TokenHmyTransactorSession) Initialize(_bmc common.Address, _bshProxy common.Address, _serviceName string) (*types.Transaction, error) {
	return _TokenHmy.Contract.Initialize(&_TokenHmy.TransactOpts, _bmc, _bshProxy, _serviceName)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0x5c436dbe.
//
// Solidity: function sendServiceMessage(address _from, string _to, (string,uint256,uint256)[] _assets) returns()
func (_TokenHmy *TokenHmyTransactor) SendServiceMessage(opts *bind.TransactOpts, _from common.Address, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _TokenHmy.contract.Transact(opts, "sendServiceMessage", _from, _to, _assets)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0x5c436dbe.
//
// Solidity: function sendServiceMessage(address _from, string _to, (string,uint256,uint256)[] _assets) returns()
func (_TokenHmy *TokenHmySession) SendServiceMessage(_from common.Address, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _TokenHmy.Contract.SendServiceMessage(&_TokenHmy.TransactOpts, _from, _to, _assets)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0x5c436dbe.
//
// Solidity: function sendServiceMessage(address _from, string _to, (string,uint256,uint256)[] _assets) returns()
func (_TokenHmy *TokenHmyTransactorSession) SendServiceMessage(_from common.Address, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _TokenHmy.Contract.SendServiceMessage(&_TokenHmy.TransactOpts, _from, _to, _assets)
}

// TokenHmyHandleBTPMessageEventIterator is returned from FilterHandleBTPMessageEvent and is used to iterate over the raw logs and unpacked data for HandleBTPMessageEvent events raised by the TokenHmy contract.
type TokenHmyHandleBTPMessageEventIterator struct {
	Event *TokenHmyHandleBTPMessageEvent // Event containing the contract specifics and raw log

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
func (it *TokenHmyHandleBTPMessageEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHmyHandleBTPMessageEvent)
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
		it.Event = new(TokenHmyHandleBTPMessageEvent)
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
func (it *TokenHmyHandleBTPMessageEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHmyHandleBTPMessageEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHmyHandleBTPMessageEvent represents a HandleBTPMessageEvent event raised by the TokenHmy contract.
type TokenHmyHandleBTPMessageEvent struct {
	Sn   *big.Int
	Code *big.Int
	Msg  string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterHandleBTPMessageEvent is a free log retrieval operation binding the contract event 0x356868e4a05430bccb6aa9c954e410ab0792c5a5baa7b973b03e1d4c03fa1366.
//
// Solidity: event HandleBTPMessageEvent(uint256 _sn, uint256 _code, string _msg)
func (_TokenHmy *TokenHmyFilterer) FilterHandleBTPMessageEvent(opts *bind.FilterOpts) (*TokenHmyHandleBTPMessageEventIterator, error) {

	logs, sub, err := _TokenHmy.contract.FilterLogs(opts, "HandleBTPMessageEvent")
	if err != nil {
		return nil, err
	}
	return &TokenHmyHandleBTPMessageEventIterator{contract: _TokenHmy.contract, event: "HandleBTPMessageEvent", logs: logs, sub: sub}, nil
}

// WatchHandleBTPMessageEvent is a free log subscription operation binding the contract event 0x356868e4a05430bccb6aa9c954e410ab0792c5a5baa7b973b03e1d4c03fa1366.
//
// Solidity: event HandleBTPMessageEvent(uint256 _sn, uint256 _code, string _msg)
func (_TokenHmy *TokenHmyFilterer) WatchHandleBTPMessageEvent(opts *bind.WatchOpts, sink chan<- *TokenHmyHandleBTPMessageEvent) (event.Subscription, error) {

	logs, sub, err := _TokenHmy.contract.WatchLogs(opts, "HandleBTPMessageEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHmyHandleBTPMessageEvent)
				if err := _TokenHmy.contract.UnpackLog(event, "HandleBTPMessageEvent", log); err != nil {
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
func (_TokenHmy *TokenHmyFilterer) ParseHandleBTPMessageEvent(log types.Log) (*TokenHmyHandleBTPMessageEvent, error) {
	event := new(TokenHmyHandleBTPMessageEvent)
	if err := _TokenHmy.contract.UnpackLog(event, "HandleBTPMessageEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenHmyResponseUnknownTypeIterator is returned from FilterResponseUnknownType and is used to iterate over the raw logs and unpacked data for ResponseUnknownType events raised by the TokenHmy contract.
type TokenHmyResponseUnknownTypeIterator struct {
	Event *TokenHmyResponseUnknownType // Event containing the contract specifics and raw log

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
func (it *TokenHmyResponseUnknownTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHmyResponseUnknownType)
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
		it.Event = new(TokenHmyResponseUnknownType)
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
func (it *TokenHmyResponseUnknownTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHmyResponseUnknownTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHmyResponseUnknownType represents a ResponseUnknownType event raised by the TokenHmy contract.
type TokenHmyResponseUnknownType struct {
	From string
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResponseUnknownType is a free log retrieval operation binding the contract event 0x64f88365dae9c547bfcae6a186a7827fdc2613baffd8d5164dc59a74f55fbeba.
//
// Solidity: event ResponseUnknownType(string _from, uint256 _sn)
func (_TokenHmy *TokenHmyFilterer) FilterResponseUnknownType(opts *bind.FilterOpts) (*TokenHmyResponseUnknownTypeIterator, error) {

	logs, sub, err := _TokenHmy.contract.FilterLogs(opts, "ResponseUnknownType")
	if err != nil {
		return nil, err
	}
	return &TokenHmyResponseUnknownTypeIterator{contract: _TokenHmy.contract, event: "ResponseUnknownType", logs: logs, sub: sub}, nil
}

// WatchResponseUnknownType is a free log subscription operation binding the contract event 0x64f88365dae9c547bfcae6a186a7827fdc2613baffd8d5164dc59a74f55fbeba.
//
// Solidity: event ResponseUnknownType(string _from, uint256 _sn)
func (_TokenHmy *TokenHmyFilterer) WatchResponseUnknownType(opts *bind.WatchOpts, sink chan<- *TokenHmyResponseUnknownType) (event.Subscription, error) {

	logs, sub, err := _TokenHmy.contract.WatchLogs(opts, "ResponseUnknownType")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHmyResponseUnknownType)
				if err := _TokenHmy.contract.UnpackLog(event, "ResponseUnknownType", log); err != nil {
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
func (_TokenHmy *TokenHmyFilterer) ParseResponseUnknownType(log types.Log) (*TokenHmyResponseUnknownType, error) {
	event := new(TokenHmyResponseUnknownType)
	if err := _TokenHmy.contract.UnpackLog(event, "ResponseUnknownType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenHmyTransferEndIterator is returned from FilterTransferEnd and is used to iterate over the raw logs and unpacked data for TransferEnd events raised by the TokenHmy contract.
type TokenHmyTransferEndIterator struct {
	Event *TokenHmyTransferEnd // Event containing the contract specifics and raw log

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
func (it *TokenHmyTransferEndIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHmyTransferEnd)
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
		it.Event = new(TokenHmyTransferEnd)
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
func (it *TokenHmyTransferEndIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHmyTransferEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHmyTransferEnd represents a TransferEnd event raised by the TokenHmy contract.
type TokenHmyTransferEnd struct {
	From     common.Address
	Sn       *big.Int
	Code     *big.Int
	Response string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferEnd is a free log retrieval operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_TokenHmy *TokenHmyFilterer) FilterTransferEnd(opts *bind.FilterOpts, _from []common.Address) (*TokenHmyTransferEndIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _TokenHmy.contract.FilterLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return &TokenHmyTransferEndIterator{contract: _TokenHmy.contract, event: "TransferEnd", logs: logs, sub: sub}, nil
}

// WatchTransferEnd is a free log subscription operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_TokenHmy *TokenHmyFilterer) WatchTransferEnd(opts *bind.WatchOpts, sink chan<- *TokenHmyTransferEnd, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _TokenHmy.contract.WatchLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHmyTransferEnd)
				if err := _TokenHmy.contract.UnpackLog(event, "TransferEnd", log); err != nil {
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
func (_TokenHmy *TokenHmyFilterer) ParseTransferEnd(log types.Log) (*TokenHmyTransferEnd, error) {
	event := new(TokenHmyTransferEnd)
	if err := _TokenHmy.contract.UnpackLog(event, "TransferEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenHmyTransferReceivedIterator is returned from FilterTransferReceived and is used to iterate over the raw logs and unpacked data for TransferReceived events raised by the TokenHmy contract.
type TokenHmyTransferReceivedIterator struct {
	Event *TokenHmyTransferReceived // Event containing the contract specifics and raw log

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
func (it *TokenHmyTransferReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHmyTransferReceived)
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
		it.Event = new(TokenHmyTransferReceived)
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
func (it *TokenHmyTransferReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHmyTransferReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHmyTransferReceived represents a TransferReceived event raised by the TokenHmy contract.
type TokenHmyTransferReceived struct {
	From         common.Hash
	To           common.Address
	Sn           *big.Int
	AssetDetails []TypesAsset
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferReceived is a free log retrieval operation binding the contract event 0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_TokenHmy *TokenHmyFilterer) FilterTransferReceived(opts *bind.FilterOpts, _from []string, _to []common.Address) (*TokenHmyTransferReceivedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _TokenHmy.contract.FilterLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &TokenHmyTransferReceivedIterator{contract: _TokenHmy.contract, event: "TransferReceived", logs: logs, sub: sub}, nil
}

// WatchTransferReceived is a free log subscription operation binding the contract event 0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_TokenHmy *TokenHmyFilterer) WatchTransferReceived(opts *bind.WatchOpts, sink chan<- *TokenHmyTransferReceived, _from []string, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _TokenHmy.contract.WatchLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHmyTransferReceived)
				if err := _TokenHmy.contract.UnpackLog(event, "TransferReceived", log); err != nil {
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
func (_TokenHmy *TokenHmyFilterer) ParseTransferReceived(log types.Log) (*TokenHmyTransferReceived, error) {
	event := new(TokenHmyTransferReceived)
	if err := _TokenHmy.contract.UnpackLog(event, "TransferReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenHmyTransferStartIterator is returned from FilterTransferStart and is used to iterate over the raw logs and unpacked data for TransferStart events raised by the TokenHmy contract.
type TokenHmyTransferStartIterator struct {
	Event *TokenHmyTransferStart // Event containing the contract specifics and raw log

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
func (it *TokenHmyTransferStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHmyTransferStart)
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
		it.Event = new(TokenHmyTransferStart)
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
func (it *TokenHmyTransferStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHmyTransferStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHmyTransferStart represents a TransferStart event raised by the TokenHmy contract.
type TokenHmyTransferStart struct {
	From   common.Address
	To     string
	Sn     *big.Int
	Assets []TypesAsset
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferStart is a free log retrieval operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assets)
func (_TokenHmy *TokenHmyFilterer) FilterTransferStart(opts *bind.FilterOpts, _from []common.Address) (*TokenHmyTransferStartIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _TokenHmy.contract.FilterLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return &TokenHmyTransferStartIterator{contract: _TokenHmy.contract, event: "TransferStart", logs: logs, sub: sub}, nil
}

// WatchTransferStart is a free log subscription operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assets)
func (_TokenHmy *TokenHmyFilterer) WatchTransferStart(opts *bind.WatchOpts, sink chan<- *TokenHmyTransferStart, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _TokenHmy.contract.WatchLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHmyTransferStart)
				if err := _TokenHmy.contract.UnpackLog(event, "TransferStart", log); err != nil {
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
func (_TokenHmy *TokenHmyFilterer) ParseTransferStart(log types.Log) (*TokenHmyTransferStart, error) {
	event := new(TokenHmyTransferStart)
	if err := _TokenHmy.contract.UnpackLog(event, "TransferStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
