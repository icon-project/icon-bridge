// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nativeHmy

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

// NativeHmyABI is the input ABI used to generate the binding from.
const NativeHmyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_response\",\"type\":\"string\"}],\"name\":\"TransferEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.AssetTransferDetail[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"}],\"name\":\"UnknownResponse\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bmc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bshCore\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_serviceName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"}],\"name\":\"sendServiceMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"handleBTPMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"handleBTPError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"handleRequestService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_fa\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"}],\"name\":\"handleFeeGathering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"checkParseAddress\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// NativeHmy is an auto generated Go binding around an Ethereum contract.
type NativeHmy struct {
	NativeHmyCaller     // Read-only binding to the contract
	NativeHmyTransactor // Write-only binding to the contract
	NativeHmyFilterer   // Log filterer for contract events
}

// NativeHmyCaller is an auto generated read-only Go binding around an Ethereum contract.
type NativeHmyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NativeHmyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NativeHmyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NativeHmyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NativeHmyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NativeHmySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NativeHmySession struct {
	Contract     *NativeHmy        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NativeHmyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NativeHmyCallerSession struct {
	Contract *NativeHmyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// NativeHmyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NativeHmyTransactorSession struct {
	Contract     *NativeHmyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// NativeHmyRaw is an auto generated low-level Go binding around an Ethereum contract.
type NativeHmyRaw struct {
	Contract *NativeHmy // Generic contract binding to access the raw methods on
}

// NativeHmyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NativeHmyCallerRaw struct {
	Contract *NativeHmyCaller // Generic read-only contract binding to access the raw methods on
}

// NativeHmyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NativeHmyTransactorRaw struct {
	Contract *NativeHmyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNativeHmy creates a new instance of NativeHmy, bound to a specific deployed contract.
func NewNativeHmy(address common.Address, backend bind.ContractBackend) (*NativeHmy, error) {
	contract, err := bindNativeHmy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NativeHmy{NativeHmyCaller: NativeHmyCaller{contract: contract}, NativeHmyTransactor: NativeHmyTransactor{contract: contract}, NativeHmyFilterer: NativeHmyFilterer{contract: contract}}, nil
}

// NewNativeHmyCaller creates a new read-only instance of NativeHmy, bound to a specific deployed contract.
func NewNativeHmyCaller(address common.Address, caller bind.ContractCaller) (*NativeHmyCaller, error) {
	contract, err := bindNativeHmy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NativeHmyCaller{contract: contract}, nil
}

// NewNativeHmyTransactor creates a new write-only instance of NativeHmy, bound to a specific deployed contract.
func NewNativeHmyTransactor(address common.Address, transactor bind.ContractTransactor) (*NativeHmyTransactor, error) {
	contract, err := bindNativeHmy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NativeHmyTransactor{contract: contract}, nil
}

// NewNativeHmyFilterer creates a new log filterer instance of NativeHmy, bound to a specific deployed contract.
func NewNativeHmyFilterer(address common.Address, filterer bind.ContractFilterer) (*NativeHmyFilterer, error) {
	contract, err := bindNativeHmy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NativeHmyFilterer{contract: contract}, nil
}

// bindNativeHmy binds a generic wrapper to an already deployed contract.
func bindNativeHmy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NativeHmyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NativeHmy *NativeHmyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NativeHmy.Contract.NativeHmyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NativeHmy *NativeHmyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeHmy.Contract.NativeHmyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NativeHmy *NativeHmyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NativeHmy.Contract.NativeHmyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NativeHmy *NativeHmyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NativeHmy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NativeHmy *NativeHmyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeHmy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NativeHmy *NativeHmyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NativeHmy.Contract.contract.Transact(opts, method, params...)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_NativeHmy *NativeHmyCaller) CheckParseAddress(opts *bind.CallOpts, _to string) error {
	var out []interface{}
	err := _NativeHmy.contract.Call(opts, &out, "checkParseAddress", _to)

	if err != nil {
		return err
	}

	return err

}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_NativeHmy *NativeHmySession) CheckParseAddress(_to string) error {
	return _NativeHmy.Contract.CheckParseAddress(&_NativeHmy.CallOpts, _to)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_NativeHmy *NativeHmyCallerSession) CheckParseAddress(_to string) error {
	return _NativeHmy.Contract.CheckParseAddress(&_NativeHmy.CallOpts, _to)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_NativeHmy *NativeHmyCaller) HasPendingRequest(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _NativeHmy.contract.Call(opts, &out, "hasPendingRequest")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_NativeHmy *NativeHmySession) HasPendingRequest() (bool, error) {
	return _NativeHmy.Contract.HasPendingRequest(&_NativeHmy.CallOpts)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_NativeHmy *NativeHmyCallerSession) HasPendingRequest() (bool, error) {
	return _NativeHmy.Contract.HasPendingRequest(&_NativeHmy.CallOpts)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_NativeHmy *NativeHmyCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	var out []interface{}
	err := _NativeHmy.contract.Call(opts, &out, "requests", arg0)

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
func (_NativeHmy *NativeHmySession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _NativeHmy.Contract.Requests(&_NativeHmy.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_NativeHmy *NativeHmyCallerSession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _NativeHmy.Contract.Requests(&_NativeHmy.CallOpts, arg0)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_NativeHmy *NativeHmyCaller) ServiceName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NativeHmy.contract.Call(opts, &out, "serviceName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_NativeHmy *NativeHmySession) ServiceName() (string, error) {
	return _NativeHmy.Contract.ServiceName(&_NativeHmy.CallOpts)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_NativeHmy *NativeHmyCallerSession) ServiceName() (string, error) {
	return _NativeHmy.Contract.ServiceName(&_NativeHmy.CallOpts)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_NativeHmy *NativeHmyTransactor) HandleBTPError(opts *bind.TransactOpts, arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _NativeHmy.contract.Transact(opts, "handleBTPError", arg0, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_NativeHmy *NativeHmySession) HandleBTPError(arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleBTPError(&_NativeHmy.TransactOpts, arg0, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_NativeHmy *NativeHmyTransactorSession) HandleBTPError(arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleBTPError(&_NativeHmy.TransactOpts, arg0, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_NativeHmy *NativeHmyTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _NativeHmy.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_NativeHmy *NativeHmySession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleBTPMessage(&_NativeHmy.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_NativeHmy *NativeHmyTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleBTPMessage(&_NativeHmy.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_NativeHmy *NativeHmyTransactor) HandleFeeGathering(opts *bind.TransactOpts, _fa string, _svc string) (*types.Transaction, error) {
	return _NativeHmy.contract.Transact(opts, "handleFeeGathering", _fa, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_NativeHmy *NativeHmySession) HandleFeeGathering(_fa string, _svc string) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleFeeGathering(&_NativeHmy.TransactOpts, _fa, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_NativeHmy *NativeHmyTransactorSession) HandleFeeGathering(_fa string, _svc string) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleFeeGathering(&_NativeHmy.TransactOpts, _fa, _svc)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_NativeHmy *NativeHmyTransactor) HandleRequestService(opts *bind.TransactOpts, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _NativeHmy.contract.Transact(opts, "handleRequestService", _to, _assets)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_NativeHmy *NativeHmySession) HandleRequestService(_to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleRequestService(&_NativeHmy.TransactOpts, _to, _assets)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_NativeHmy *NativeHmyTransactorSession) HandleRequestService(_to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _NativeHmy.Contract.HandleRequestService(&_NativeHmy.TransactOpts, _to, _assets)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshCore, string _serviceName) returns()
func (_NativeHmy *NativeHmyTransactor) Initialize(opts *bind.TransactOpts, _bmc common.Address, _bshCore common.Address, _serviceName string) (*types.Transaction, error) {
	return _NativeHmy.contract.Transact(opts, "initialize", _bmc, _bshCore, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshCore, string _serviceName) returns()
func (_NativeHmy *NativeHmySession) Initialize(_bmc common.Address, _bshCore common.Address, _serviceName string) (*types.Transaction, error) {
	return _NativeHmy.Contract.Initialize(&_NativeHmy.TransactOpts, _bmc, _bshCore, _serviceName)
}

// Initialize is a paid mutator transaction binding the contract method 0x4571e3a6.
//
// Solidity: function initialize(address _bmc, address _bshCore, string _serviceName) returns()
func (_NativeHmy *NativeHmyTransactorSession) Initialize(_bmc common.Address, _bshCore common.Address, _serviceName string) (*types.Transaction, error) {
	return _NativeHmy.Contract.Initialize(&_NativeHmy.TransactOpts, _bmc, _bshCore, _serviceName)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_NativeHmy *NativeHmyTransactor) SendServiceMessage(opts *bind.TransactOpts, _from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _NativeHmy.contract.Transact(opts, "sendServiceMessage", _from, _to, _coinNames, _values, _fees)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_NativeHmy *NativeHmySession) SendServiceMessage(_from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _NativeHmy.Contract.SendServiceMessage(&_NativeHmy.TransactOpts, _from, _to, _coinNames, _values, _fees)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_NativeHmy *NativeHmyTransactorSession) SendServiceMessage(_from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _NativeHmy.Contract.SendServiceMessage(&_NativeHmy.TransactOpts, _from, _to, _coinNames, _values, _fees)
}

// NativeHmyTransferEndIterator is returned from FilterTransferEnd and is used to iterate over the raw logs and unpacked data for TransferEnd events raised by the NativeHmy contract.
type NativeHmyTransferEndIterator struct {
	Event *NativeHmyTransferEnd // Event containing the contract specifics and raw log

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
func (it *NativeHmyTransferEndIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeHmyTransferEnd)
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
		it.Event = new(NativeHmyTransferEnd)
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
func (it *NativeHmyTransferEndIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NativeHmyTransferEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NativeHmyTransferEnd represents a TransferEnd event raised by the NativeHmy contract.
type NativeHmyTransferEnd struct {
	From     common.Address
	Sn       *big.Int
	Code     *big.Int
	Response string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferEnd is a free log retrieval operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_NativeHmy *NativeHmyFilterer) FilterTransferEnd(opts *bind.FilterOpts, _from []common.Address) (*NativeHmyTransferEndIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _NativeHmy.contract.FilterLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return &NativeHmyTransferEndIterator{contract: _NativeHmy.contract, event: "TransferEnd", logs: logs, sub: sub}, nil
}

// WatchTransferEnd is a free log subscription operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_NativeHmy *NativeHmyFilterer) WatchTransferEnd(opts *bind.WatchOpts, sink chan<- *NativeHmyTransferEnd, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _NativeHmy.contract.WatchLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NativeHmyTransferEnd)
				if err := _NativeHmy.contract.UnpackLog(event, "TransferEnd", log); err != nil {
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
func (_NativeHmy *NativeHmyFilterer) ParseTransferEnd(log types.Log) (*NativeHmyTransferEnd, error) {
	event := new(NativeHmyTransferEnd)
	if err := _NativeHmy.contract.UnpackLog(event, "TransferEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NativeHmyTransferReceivedIterator is returned from FilterTransferReceived and is used to iterate over the raw logs and unpacked data for TransferReceived events raised by the NativeHmy contract.
type NativeHmyTransferReceivedIterator struct {
	Event *NativeHmyTransferReceived // Event containing the contract specifics and raw log

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
func (it *NativeHmyTransferReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeHmyTransferReceived)
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
		it.Event = new(NativeHmyTransferReceived)
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
func (it *NativeHmyTransferReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NativeHmyTransferReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NativeHmyTransferReceived represents a TransferReceived event raised by the NativeHmy contract.
type NativeHmyTransferReceived struct {
	From         common.Hash
	To           common.Address
	Sn           *big.Int
	AssetDetails []TypesAsset
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferReceived is a free log retrieval operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_NativeHmy *NativeHmyFilterer) FilterTransferReceived(opts *bind.FilterOpts, _from []string, _to []common.Address) (*NativeHmyTransferReceivedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _NativeHmy.contract.FilterLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &NativeHmyTransferReceivedIterator{contract: _NativeHmy.contract, event: "TransferReceived", logs: logs, sub: sub}, nil
}

// WatchTransferReceived is a free log subscription operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_NativeHmy *NativeHmyFilterer) WatchTransferReceived(opts *bind.WatchOpts, sink chan<- *NativeHmyTransferReceived, _from []string, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _NativeHmy.contract.WatchLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NativeHmyTransferReceived)
				if err := _NativeHmy.contract.UnpackLog(event, "TransferReceived", log); err != nil {
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
func (_NativeHmy *NativeHmyFilterer) ParseTransferReceived(log types.Log) (*NativeHmyTransferReceived, error) {
	event := new(NativeHmyTransferReceived)
	if err := _NativeHmy.contract.UnpackLog(event, "TransferReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NativeHmyTransferStartIterator is returned from FilterTransferStart and is used to iterate over the raw logs and unpacked data for TransferStart events raised by the NativeHmy contract.
type NativeHmyTransferStartIterator struct {
	Event *NativeHmyTransferStart // Event containing the contract specifics and raw log

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
func (it *NativeHmyTransferStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeHmyTransferStart)
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
		it.Event = new(NativeHmyTransferStart)
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
func (it *NativeHmyTransferStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NativeHmyTransferStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NativeHmyTransferStart represents a TransferStart event raised by the NativeHmy contract.
type NativeHmyTransferStart struct {
	From         common.Address
	To           string
	Sn           *big.Int
	AssetDetails []TypesAssetTransferDetail
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferStart is a free log retrieval operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_NativeHmy *NativeHmyFilterer) FilterTransferStart(opts *bind.FilterOpts, _from []common.Address) (*NativeHmyTransferStartIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _NativeHmy.contract.FilterLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return &NativeHmyTransferStartIterator{contract: _NativeHmy.contract, event: "TransferStart", logs: logs, sub: sub}, nil
}

// WatchTransferStart is a free log subscription operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_NativeHmy *NativeHmyFilterer) WatchTransferStart(opts *bind.WatchOpts, sink chan<- *NativeHmyTransferStart, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _NativeHmy.contract.WatchLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NativeHmyTransferStart)
				if err := _NativeHmy.contract.UnpackLog(event, "TransferStart", log); err != nil {
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
func (_NativeHmy *NativeHmyFilterer) ParseTransferStart(log types.Log) (*NativeHmyTransferStart, error) {
	event := new(NativeHmyTransferStart)
	if err := _NativeHmy.contract.UnpackLog(event, "TransferStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NativeHmyUnknownResponseIterator is returned from FilterUnknownResponse and is used to iterate over the raw logs and unpacked data for UnknownResponse events raised by the NativeHmy contract.
type NativeHmyUnknownResponseIterator struct {
	Event *NativeHmyUnknownResponse // Event containing the contract specifics and raw log

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
func (it *NativeHmyUnknownResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeHmyUnknownResponse)
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
		it.Event = new(NativeHmyUnknownResponse)
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
func (it *NativeHmyUnknownResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NativeHmyUnknownResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NativeHmyUnknownResponse represents a UnknownResponse event raised by the NativeHmy contract.
type NativeHmyUnknownResponse struct {
	From string
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUnknownResponse is a free log retrieval operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_NativeHmy *NativeHmyFilterer) FilterUnknownResponse(opts *bind.FilterOpts) (*NativeHmyUnknownResponseIterator, error) {

	logs, sub, err := _NativeHmy.contract.FilterLogs(opts, "UnknownResponse")
	if err != nil {
		return nil, err
	}
	return &NativeHmyUnknownResponseIterator{contract: _NativeHmy.contract, event: "UnknownResponse", logs: logs, sub: sub}, nil
}

// WatchUnknownResponse is a free log subscription operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_NativeHmy *NativeHmyFilterer) WatchUnknownResponse(opts *bind.WatchOpts, sink chan<- *NativeHmyUnknownResponse) (event.Subscription, error) {

	logs, sub, err := _NativeHmy.contract.WatchLogs(opts, "UnknownResponse")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NativeHmyUnknownResponse)
				if err := _NativeHmy.contract.UnpackLog(event, "UnknownResponse", log); err != nil {
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
func (_NativeHmy *NativeHmyFilterer) ParseUnknownResponse(log types.Log) (*NativeHmyUnknownResponse, error) {
	event := new(NativeHmyUnknownResponse)
	if err := _NativeHmy.contract.UnpackLog(event, "UnknownResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
