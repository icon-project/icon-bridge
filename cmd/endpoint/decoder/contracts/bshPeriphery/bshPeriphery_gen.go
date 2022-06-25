// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bshPeriphery

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
	CoinName string
	Value    *big.Int
}

// TypesAssetTransferDetail is an auto generated low-level Go binding around an user-defined struct.
type TypesAssetTransferDetail struct {
	CoinName string
	Value    *big.Int
	Fee      *big.Int
}

// BshPeripheryABI is the input ABI used to generate the binding from.
const BshPeripheryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_response\",\"type\":\"string\"}],\"name\":\"TransferEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.AssetTransferDetail[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"}],\"name\":\"UnknownResponse\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bmc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bshCore\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_serviceName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"}],\"name\":\"sendServiceMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"handleBTPMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"handleBTPError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"handleRequestService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_fa\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"}],\"name\":\"handleFeeGathering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"checkParseAddress\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// BshPeriphery is an auto generated Go binding around an Ethereum contract.
type BshPeriphery struct {
	BshPeripheryCaller     // Read-only binding to the contract
	BshPeripheryTransactor // Write-only binding to the contract
	BshPeripheryFilterer   // Log filterer for contract events
}

// BshPeripheryCaller is an auto generated read-only Go binding around an Ethereum contract.
type BshPeripheryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshPeripheryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BshPeripheryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshPeripheryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BshPeripheryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshPeripherySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BshPeripherySession struct {
	Contract     *BshPeriphery     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BshPeripheryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BshPeripheryCallerSession struct {
	Contract *BshPeripheryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BshPeripheryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BshPeripheryTransactorSession struct {
	Contract     *BshPeripheryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BshPeripheryRaw is an auto generated low-level Go binding around an Ethereum contract.
type BshPeripheryRaw struct {
	Contract *BshPeriphery // Generic contract binding to access the raw methods on
}

// BshPeripheryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BshPeripheryCallerRaw struct {
	Contract *BshPeripheryCaller // Generic read-only contract binding to access the raw methods on
}

// BshPeripheryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BshPeripheryTransactorRaw struct {
	Contract *BshPeripheryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBshPeriphery creates a new instance of BshPeriphery, bound to a specific deployed contract.
func NewBshPeriphery(address common.Address, backend bind.ContractBackend) (*BshPeriphery, error) {
	contract, err := bindBshPeriphery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BshPeriphery{BshPeripheryCaller: BshPeripheryCaller{contract: contract}, BshPeripheryTransactor: BshPeripheryTransactor{contract: contract}, BshPeripheryFilterer: BshPeripheryFilterer{contract: contract}}, nil
}

// NewBshPeripheryCaller creates a new read-only instance of BshPeriphery, bound to a specific deployed contract.
func NewBshPeripheryCaller(address common.Address, caller bind.ContractCaller) (*BshPeripheryCaller, error) {
	contract, err := bindBshPeriphery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BshPeripheryCaller{contract: contract}, nil
}

// NewBshPeripheryTransactor creates a new write-only instance of BshPeriphery, bound to a specific deployed contract.
func NewBshPeripheryTransactor(address common.Address, transactor bind.ContractTransactor) (*BshPeripheryTransactor, error) {
	contract, err := bindBshPeriphery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BshPeripheryTransactor{contract: contract}, nil
}

// NewBshPeripheryFilterer creates a new log filterer instance of BshPeriphery, bound to a specific deployed contract.
func NewBshPeripheryFilterer(address common.Address, filterer bind.ContractFilterer) (*BshPeripheryFilterer, error) {
	contract, err := bindBshPeriphery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BshPeripheryFilterer{contract: contract}, nil
}

// bindBshPeriphery binds a generic wrapper to an already deployed contract.
func bindBshPeriphery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BshPeripheryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BshPeriphery *BshPeripheryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BshPeriphery.Contract.BshPeripheryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BshPeriphery *BshPeripheryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BshPeriphery.Contract.BshPeripheryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BshPeriphery *BshPeripheryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BshPeriphery.Contract.BshPeripheryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BshPeriphery *BshPeripheryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BshPeriphery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BshPeriphery *BshPeripheryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BshPeriphery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BshPeriphery *BshPeripheryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BshPeriphery.Contract.contract.Transact(opts, method, params...)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_BshPeriphery *BshPeripheryCaller) CheckParseAddress(opts *bind.CallOpts, _to string) error {
	var out []interface{}
	err := _BshPeriphery.contract.Call(opts, &out, "checkParseAddress", _to)

	if err != nil {
		return err
	}

	return err

}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_BshPeriphery *BshPeripherySession) CheckParseAddress(_to string) error {
	return _BshPeriphery.Contract.CheckParseAddress(&_BshPeriphery.CallOpts, _to)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_BshPeriphery *BshPeripheryCallerSession) CheckParseAddress(_to string) error {
	return _BshPeriphery.Contract.CheckParseAddress(&_BshPeriphery.CallOpts, _to)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_BshPeriphery *BshPeripheryCaller) HasPendingRequest(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BshPeriphery.contract.Call(opts, &out, "hasPendingRequest")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_BshPeriphery *BshPeripherySession) HasPendingRequest() (bool, error) {
	return _BshPeriphery.Contract.HasPendingRequest(&_BshPeriphery.CallOpts)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_BshPeriphery *BshPeripheryCallerSession) HasPendingRequest() (bool, error) {
	return _BshPeriphery.Contract.HasPendingRequest(&_BshPeriphery.CallOpts)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_BshPeriphery *BshPeripheryCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	var out []interface{}
	err := _BshPeriphery.contract.Call(opts, &out, "requests", arg0)

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
func (_BshPeriphery *BshPeripherySession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _BshPeriphery.Contract.Requests(&_BshPeriphery.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_BshPeriphery *BshPeripheryCallerSession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _BshPeriphery.Contract.Requests(&_BshPeriphery.CallOpts, arg0)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_BshPeriphery *BshPeripheryCaller) ServiceName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BshPeriphery.contract.Call(opts, &out, "serviceName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_BshPeriphery *BshPeripherySession) ServiceName() (string, error) {
	return _BshPeriphery.Contract.ServiceName(&_BshPeriphery.CallOpts)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_BshPeriphery *BshPeripheryCallerSession) ServiceName() (string, error) {
	return _BshPeriphery.Contract.ServiceName(&_BshPeriphery.CallOpts)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_BshPeriphery *BshPeripheryTransactor) HandleBTPError(opts *bind.TransactOpts, arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _BshPeriphery.contract.Transact(opts, "handleBTPError", arg0, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_BshPeriphery *BshPeripherySession) HandleBTPError(arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleBTPError(&_BshPeriphery.TransactOpts, arg0, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_BshPeriphery *BshPeripheryTransactorSession) HandleBTPError(arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleBTPError(&_BshPeriphery.TransactOpts, arg0, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_BshPeriphery *BshPeripheryTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BshPeriphery.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_BshPeriphery *BshPeripherySession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleBTPMessage(&_BshPeriphery.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_BshPeriphery *BshPeripheryTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleBTPMessage(&_BshPeriphery.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_BshPeriphery *BshPeripheryTransactor) HandleFeeGathering(opts *bind.TransactOpts, _fa string, _svc string) (*types.Transaction, error) {
	return _BshPeriphery.contract.Transact(opts, "handleFeeGathering", _fa, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_BshPeriphery *BshPeripherySession) HandleFeeGathering(_fa string, _svc string) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleFeeGathering(&_BshPeriphery.TransactOpts, _fa, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_BshPeriphery *BshPeripheryTransactorSession) HandleFeeGathering(_fa string, _svc string) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleFeeGathering(&_BshPeriphery.TransactOpts, _fa, _svc)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_BshPeriphery *BshPeripheryTransactor) HandleRequestService(opts *bind.TransactOpts, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _BshPeriphery.contract.Transact(opts, "handleRequestService", _to, _assets)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_BshPeriphery *BshPeripherySession) HandleRequestService(_to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleRequestService(&_BshPeriphery.TransactOpts, _to, _assets)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_BshPeriphery *BshPeripheryTransactorSession) HandleRequestService(_to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _BshPeriphery.Contract.HandleRequestService(&_BshPeriphery.TransactOpts, _to, _assets)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshCore, string _serviceName) returns()
func (_BshPeriphery *BshPeripheryTransactor) Initialize(opts *bind.TransactOpts, _bmc common.Address, _bshCore common.Address, _serviceName string) (*types.Transaction, error) {
	return _BshPeriphery.contract.Transact(opts, "initialize", _bmc, _bshCore, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshCore, string _serviceName) returns()
func (_BshPeriphery *BshPeripherySession) Initialize(_bmc common.Address, _bshCore common.Address, _serviceName string) (*types.Transaction, error) {
	return _BshPeriphery.Contract.Initialize(&_BshPeriphery.TransactOpts, _bmc, _bshCore, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshCore, string _serviceName) returns()
func (_BshPeriphery *BshPeripheryTransactorSession) Initialize(_bmc common.Address, _bshCore common.Address, _serviceName string) (*types.Transaction, error) {
	return _BshPeriphery.Contract.Initialize(&_BshPeriphery.TransactOpts, _bmc, _bshCore, _serviceName)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_BshPeriphery *BshPeripheryTransactor) SendServiceMessage(opts *bind.TransactOpts, _from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _BshPeriphery.contract.Transact(opts, "sendServiceMessage", _from, _to, _coinNames, _values, _fees)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_BshPeriphery *BshPeripherySession) SendServiceMessage(_from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _BshPeriphery.Contract.SendServiceMessage(&_BshPeriphery.TransactOpts, _from, _to, _coinNames, _values, _fees)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_BshPeriphery *BshPeripheryTransactorSession) SendServiceMessage(_from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _BshPeriphery.Contract.SendServiceMessage(&_BshPeriphery.TransactOpts, _from, _to, _coinNames, _values, _fees)
}

// BshPeripheryTransferEndIterator is returned from FilterTransferEnd and is used to iterate over the raw logs and unpacked data for TransferEnd events raised by the BshPeriphery contract.
type BshPeripheryTransferEndIterator struct {
	Event *BshPeripheryTransferEnd // Event containing the contract specifics and raw log

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
func (it *BshPeripheryTransferEndIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshPeripheryTransferEnd)
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
		it.Event = new(BshPeripheryTransferEnd)
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
func (it *BshPeripheryTransferEndIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshPeripheryTransferEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshPeripheryTransferEnd represents a TransferEnd event raised by the BshPeriphery contract.
type BshPeripheryTransferEnd struct {
	From     common.Address
	Sn       *big.Int
	Code     *big.Int
	Response string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferEnd is a free log retrieval operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_BshPeriphery *BshPeripheryFilterer) FilterTransferEnd(opts *bind.FilterOpts, _from []common.Address) (*BshPeripheryTransferEndIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshPeriphery.contract.FilterLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return &BshPeripheryTransferEndIterator{contract: _BshPeriphery.contract, event: "TransferEnd", logs: logs, sub: sub}, nil
}

// WatchTransferEnd is a free log subscription operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_BshPeriphery *BshPeripheryFilterer) WatchTransferEnd(opts *bind.WatchOpts, sink chan<- *BshPeripheryTransferEnd, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshPeriphery.contract.WatchLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshPeripheryTransferEnd)
				if err := _BshPeriphery.contract.UnpackLog(event, "TransferEnd", log); err != nil {
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
func (_BshPeriphery *BshPeripheryFilterer) ParseTransferEnd(log types.Log) (*BshPeripheryTransferEnd, error) {
	event := new(BshPeripheryTransferEnd)
	if err := _BshPeriphery.contract.UnpackLog(event, "TransferEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshPeripheryTransferReceivedIterator is returned from FilterTransferReceived and is used to iterate over the raw logs and unpacked data for TransferReceived events raised by the BshPeriphery contract.
type BshPeripheryTransferReceivedIterator struct {
	Event *BshPeripheryTransferReceived // Event containing the contract specifics and raw log

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
func (it *BshPeripheryTransferReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshPeripheryTransferReceived)
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
		it.Event = new(BshPeripheryTransferReceived)
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
func (it *BshPeripheryTransferReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshPeripheryTransferReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshPeripheryTransferReceived represents a TransferReceived event raised by the BshPeriphery contract.
type BshPeripheryTransferReceived struct {
	From         common.Hash
	To           common.Address
	Sn           *big.Int
	AssetDetails []TypesAsset
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferReceived is a free log retrieval operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_BshPeriphery *BshPeripheryFilterer) FilterTransferReceived(opts *bind.FilterOpts, _from []string, _to []common.Address) (*BshPeripheryTransferReceivedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _BshPeriphery.contract.FilterLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &BshPeripheryTransferReceivedIterator{contract: _BshPeriphery.contract, event: "TransferReceived", logs: logs, sub: sub}, nil
}

// WatchTransferReceived is a free log subscription operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_BshPeriphery *BshPeripheryFilterer) WatchTransferReceived(opts *bind.WatchOpts, sink chan<- *BshPeripheryTransferReceived, _from []string, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _BshPeriphery.contract.WatchLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshPeripheryTransferReceived)
				if err := _BshPeriphery.contract.UnpackLog(event, "TransferReceived", log); err != nil {
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

// ParseTransferReceived is a log parse operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_BshPeriphery *BshPeripheryFilterer) ParseTransferReceived(log types.Log) (*BshPeripheryTransferReceived, error) {
	event := new(BshPeripheryTransferReceived)
	if err := _BshPeriphery.contract.UnpackLog(event, "TransferReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshPeripheryTransferStartIterator is returned from FilterTransferStart and is used to iterate over the raw logs and unpacked data for TransferStart events raised by the BshPeriphery contract.
type BshPeripheryTransferStartIterator struct {
	Event *BshPeripheryTransferStart // Event containing the contract specifics and raw log

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
func (it *BshPeripheryTransferStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshPeripheryTransferStart)
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
		it.Event = new(BshPeripheryTransferStart)
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
func (it *BshPeripheryTransferStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshPeripheryTransferStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshPeripheryTransferStart represents a TransferStart event raised by the BshPeriphery contract.
type BshPeripheryTransferStart struct {
	From         common.Address
	To           string
	Sn           *big.Int
	AssetDetails []TypesAssetTransferDetail
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferStart is a free log retrieval operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_BshPeriphery *BshPeripheryFilterer) FilterTransferStart(opts *bind.FilterOpts, _from []common.Address) (*BshPeripheryTransferStartIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshPeriphery.contract.FilterLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return &BshPeripheryTransferStartIterator{contract: _BshPeriphery.contract, event: "TransferStart", logs: logs, sub: sub}, nil
}

// WatchTransferStart is a free log subscription operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_BshPeriphery *BshPeripheryFilterer) WatchTransferStart(opts *bind.WatchOpts, sink chan<- *BshPeripheryTransferStart, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _BshPeriphery.contract.WatchLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshPeripheryTransferStart)
				if err := _BshPeriphery.contract.UnpackLog(event, "TransferStart", log); err != nil {
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
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_BshPeriphery *BshPeripheryFilterer) ParseTransferStart(log types.Log) (*BshPeripheryTransferStart, error) {
	event := new(BshPeripheryTransferStart)
	if err := _BshPeriphery.contract.UnpackLog(event, "TransferStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshPeripheryUnknownResponseIterator is returned from FilterUnknownResponse and is used to iterate over the raw logs and unpacked data for UnknownResponse events raised by the BshPeriphery contract.
type BshPeripheryUnknownResponseIterator struct {
	Event *BshPeripheryUnknownResponse // Event containing the contract specifics and raw log

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
func (it *BshPeripheryUnknownResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshPeripheryUnknownResponse)
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
		it.Event = new(BshPeripheryUnknownResponse)
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
func (it *BshPeripheryUnknownResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshPeripheryUnknownResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshPeripheryUnknownResponse represents a UnknownResponse event raised by the BshPeriphery contract.
type BshPeripheryUnknownResponse struct {
	From string
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUnknownResponse is a free log retrieval operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_BshPeriphery *BshPeripheryFilterer) FilterUnknownResponse(opts *bind.FilterOpts) (*BshPeripheryUnknownResponseIterator, error) {

	logs, sub, err := _BshPeriphery.contract.FilterLogs(opts, "UnknownResponse")
	if err != nil {
		return nil, err
	}
	return &BshPeripheryUnknownResponseIterator{contract: _BshPeriphery.contract, event: "UnknownResponse", logs: logs, sub: sub}, nil
}

// WatchUnknownResponse is a free log subscription operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_BshPeriphery *BshPeripheryFilterer) WatchUnknownResponse(opts *bind.WatchOpts, sink chan<- *BshPeripheryUnknownResponse) (event.Subscription, error) {

	logs, sub, err := _BshPeriphery.contract.WatchLogs(opts, "UnknownResponse")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshPeripheryUnknownResponse)
				if err := _BshPeriphery.contract.UnpackLog(event, "UnknownResponse", log); err != nil {
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

// ParseUnknownResponse is a log parse operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_BshPeriphery *BshPeripheryFilterer) ParseUnknownResponse(log types.Log) (*BshPeripheryUnknownResponse, error) {
	event := new(BshPeripheryUnknownResponse)
	if err := _BshPeriphery.contract.UnpackLog(event, "UnknownResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
