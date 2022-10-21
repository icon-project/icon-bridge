// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package btsperiphery

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

// BtsperipheryABI is the input ABI used to generate the binding from.
const BtsperipheryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_response\",\"type\":\"string\"}],\"name\":\"TransferEnd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTypes.AssetTransferDetail[]\",\"name\":\"_assetDetails\",\"type\":\"tuple[]\"}],\"name\":\"TransferStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"}],\"name\":\"UnknownResponse\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"blacklist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"serviceName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"tokenLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bmc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_btsCore\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasPendingRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_address\",\"type\":\"string[]\"}],\"name\":\"addToBlacklist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_address\",\"type\":\"string[]\"}],\"name\":\"removeFromBlacklist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_tokenLimits\",\"type\":\"uint256[]\"}],\"name\":\"setTokenLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"}],\"name\":\"sendServiceMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_from\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"handleBTPMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"handleBTPError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"}],\"name\":\"handleRequestService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_fa\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"}],\"name\":\"handleFeeGathering\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"checkParseAddress\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"checkTransferRestrictions\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]"

// Btsperiphery is an auto generated Go binding around an Ethereum contract.
type Btsperiphery struct {
	BtsperipheryCaller     // Read-only binding to the contract
	BtsperipheryTransactor // Write-only binding to the contract
	BtsperipheryFilterer   // Log filterer for contract events
}

// BtsperipheryCaller is an auto generated read-only Go binding around an Ethereum contract.
type BtsperipheryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtsperipheryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BtsperipheryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtsperipheryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BtsperipheryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtsperipherySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BtsperipherySession struct {
	Contract     *Btsperiphery     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BtsperipheryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BtsperipheryCallerSession struct {
	Contract *BtsperipheryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BtsperipheryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BtsperipheryTransactorSession struct {
	Contract     *BtsperipheryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BtsperipheryRaw is an auto generated low-level Go binding around an Ethereum contract.
type BtsperipheryRaw struct {
	Contract *Btsperiphery // Generic contract binding to access the raw methods on
}

// BtsperipheryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BtsperipheryCallerRaw struct {
	Contract *BtsperipheryCaller // Generic read-only contract binding to access the raw methods on
}

// BtsperipheryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BtsperipheryTransactorRaw struct {
	Contract *BtsperipheryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBtsperiphery creates a new instance of Btsperiphery, bound to a specific deployed contract.
func NewBtsperiphery(address common.Address, backend bind.ContractBackend) (*Btsperiphery, error) {
	contract, err := bindBtsperiphery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Btsperiphery{BtsperipheryCaller: BtsperipheryCaller{contract: contract}, BtsperipheryTransactor: BtsperipheryTransactor{contract: contract}, BtsperipheryFilterer: BtsperipheryFilterer{contract: contract}}, nil
}

// NewBtsperipheryCaller creates a new read-only instance of Btsperiphery, bound to a specific deployed contract.
func NewBtsperipheryCaller(address common.Address, caller bind.ContractCaller) (*BtsperipheryCaller, error) {
	contract, err := bindBtsperiphery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BtsperipheryCaller{contract: contract}, nil
}

// NewBtsperipheryTransactor creates a new write-only instance of Btsperiphery, bound to a specific deployed contract.
func NewBtsperipheryTransactor(address common.Address, transactor bind.ContractTransactor) (*BtsperipheryTransactor, error) {
	contract, err := bindBtsperiphery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BtsperipheryTransactor{contract: contract}, nil
}

// NewBtsperipheryFilterer creates a new log filterer instance of Btsperiphery, bound to a specific deployed contract.
func NewBtsperipheryFilterer(address common.Address, filterer bind.ContractFilterer) (*BtsperipheryFilterer, error) {
	contract, err := bindBtsperiphery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BtsperipheryFilterer{contract: contract}, nil
}

// bindBtsperiphery binds a generic wrapper to an already deployed contract.
func bindBtsperiphery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BtsperipheryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Btsperiphery *BtsperipheryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Btsperiphery.Contract.BtsperipheryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Btsperiphery *BtsperipheryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Btsperiphery.Contract.BtsperipheryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Btsperiphery *BtsperipheryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Btsperiphery.Contract.BtsperipheryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Btsperiphery *BtsperipheryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Btsperiphery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Btsperiphery *BtsperipheryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Btsperiphery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Btsperiphery *BtsperipheryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Btsperiphery.Contract.contract.Transact(opts, method, params...)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Btsperiphery *BtsperipheryCaller) Blacklist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "blacklist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Btsperiphery *BtsperipherySession) Blacklist(arg0 common.Address) (bool, error) {
	return _Btsperiphery.Contract.Blacklist(&_Btsperiphery.CallOpts, arg0)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Btsperiphery *BtsperipheryCallerSession) Blacklist(arg0 common.Address) (bool, error) {
	return _Btsperiphery.Contract.Blacklist(&_Btsperiphery.CallOpts, arg0)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_Btsperiphery *BtsperipheryCaller) CheckParseAddress(opts *bind.CallOpts, _to string) error {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "checkParseAddress", _to)

	if err != nil {
		return err
	}

	return err

}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_Btsperiphery *BtsperipherySession) CheckParseAddress(_to string) error {
	return _Btsperiphery.Contract.CheckParseAddress(&_Btsperiphery.CallOpts, _to)
}

// CheckParseAddress is a free data retrieval call binding the contract method 0xc7a6d7fe.
//
// Solidity: function checkParseAddress(string _to) pure returns()
func (_Btsperiphery *BtsperipheryCallerSession) CheckParseAddress(_to string) error {
	return _Btsperiphery.Contract.CheckParseAddress(&_Btsperiphery.CallOpts, _to)
}

// CheckTransferRestrictions is a free data retrieval call binding the contract method 0xb148f625.
//
// Solidity: function checkTransferRestrictions(string _coinName, address _user, uint256 _value) view returns()
func (_Btsperiphery *BtsperipheryCaller) CheckTransferRestrictions(opts *bind.CallOpts, _coinName string, _user common.Address, _value *big.Int) error {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "checkTransferRestrictions", _coinName, _user, _value)

	if err != nil {
		return err
	}

	return err

}

// CheckTransferRestrictions is a free data retrieval call binding the contract method 0xb148f625.
//
// Solidity: function checkTransferRestrictions(string _coinName, address _user, uint256 _value) view returns()
func (_Btsperiphery *BtsperipherySession) CheckTransferRestrictions(_coinName string, _user common.Address, _value *big.Int) error {
	return _Btsperiphery.Contract.CheckTransferRestrictions(&_Btsperiphery.CallOpts, _coinName, _user, _value)
}

// CheckTransferRestrictions is a free data retrieval call binding the contract method 0xb148f625.
//
// Solidity: function checkTransferRestrictions(string _coinName, address _user, uint256 _value) view returns()
func (_Btsperiphery *BtsperipheryCallerSession) CheckTransferRestrictions(_coinName string, _user common.Address, _value *big.Int) error {
	return _Btsperiphery.Contract.CheckTransferRestrictions(&_Btsperiphery.CallOpts, _coinName, _user, _value)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_Btsperiphery *BtsperipheryCaller) HasPendingRequest(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "hasPendingRequest")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_Btsperiphery *BtsperipherySession) HasPendingRequest() (bool, error) {
	return _Btsperiphery.Contract.HasPendingRequest(&_Btsperiphery.CallOpts)
}

// HasPendingRequest is a free data retrieval call binding the contract method 0x6bf39c09.
//
// Solidity: function hasPendingRequest() view returns(bool)
func (_Btsperiphery *BtsperipheryCallerSession) HasPendingRequest() (bool, error) {
	return _Btsperiphery.Contract.HasPendingRequest(&_Btsperiphery.CallOpts)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_Btsperiphery *BtsperipheryCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "requests", arg0)

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
func (_Btsperiphery *BtsperipherySession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _Btsperiphery.Contract.Requests(&_Btsperiphery.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(string from, string to)
func (_Btsperiphery *BtsperipheryCallerSession) Requests(arg0 *big.Int) (struct {
	From string
	To   string
}, error) {
	return _Btsperiphery.Contract.Requests(&_Btsperiphery.CallOpts, arg0)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_Btsperiphery *BtsperipheryCaller) ServiceName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "serviceName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_Btsperiphery *BtsperipherySession) ServiceName() (string, error) {
	return _Btsperiphery.Contract.ServiceName(&_Btsperiphery.CallOpts)
}

// ServiceName is a free data retrieval call binding the contract method 0x9fdc7bc4.
//
// Solidity: function serviceName() view returns(string)
func (_Btsperiphery *BtsperipheryCallerSession) ServiceName() (string, error) {
	return _Btsperiphery.Contract.ServiceName(&_Btsperiphery.CallOpts)
}

// TokenLimit is a free data retrieval call binding the contract method 0x9b61557d.
//
// Solidity: function tokenLimit(string ) view returns(uint256)
func (_Btsperiphery *BtsperipheryCaller) TokenLimit(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _Btsperiphery.contract.Call(opts, &out, "tokenLimit", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenLimit is a free data retrieval call binding the contract method 0x9b61557d.
//
// Solidity: function tokenLimit(string ) view returns(uint256)
func (_Btsperiphery *BtsperipherySession) TokenLimit(arg0 string) (*big.Int, error) {
	return _Btsperiphery.Contract.TokenLimit(&_Btsperiphery.CallOpts, arg0)
}

// TokenLimit is a free data retrieval call binding the contract method 0x9b61557d.
//
// Solidity: function tokenLimit(string ) view returns(uint256)
func (_Btsperiphery *BtsperipheryCallerSession) TokenLimit(arg0 string) (*big.Int, error) {
	return _Btsperiphery.Contract.TokenLimit(&_Btsperiphery.CallOpts, arg0)
}

// AddToBlacklist is a paid mutator transaction binding the contract method 0x4b716a84.
//
// Solidity: function addToBlacklist(string[] _address) returns()
func (_Btsperiphery *BtsperipheryTransactor) AddToBlacklist(opts *bind.TransactOpts, _address []string) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "addToBlacklist", _address)
}

// AddToBlacklist is a paid mutator transaction binding the contract method 0x4b716a84.
//
// Solidity: function addToBlacklist(string[] _address) returns()
func (_Btsperiphery *BtsperipherySession) AddToBlacklist(_address []string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.AddToBlacklist(&_Btsperiphery.TransactOpts, _address)
}

// AddToBlacklist is a paid mutator transaction binding the contract method 0x4b716a84.
//
// Solidity: function addToBlacklist(string[] _address) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) AddToBlacklist(_address []string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.AddToBlacklist(&_Btsperiphery.TransactOpts, _address)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_Btsperiphery *BtsperipheryTransactor) HandleBTPError(opts *bind.TransactOpts, arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "handleBTPError", arg0, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_Btsperiphery *BtsperipherySession) HandleBTPError(arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleBTPError(&_Btsperiphery.TransactOpts, arg0, _svc, _sn, _code, _msg)
}

// HandleBTPError is a paid mutator transaction binding the contract method 0x0a823dea.
//
// Solidity: function handleBTPError(string , string _svc, uint256 _sn, uint256 _code, string _msg) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) HandleBTPError(arg0 string, _svc string, _sn *big.Int, _code *big.Int, _msg string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleBTPError(&_Btsperiphery.TransactOpts, arg0, _svc, _sn, _code, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_Btsperiphery *BtsperipheryTransactor) HandleBTPMessage(opts *bind.TransactOpts, _from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "handleBTPMessage", _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_Btsperiphery *BtsperipherySession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleBTPMessage(&_Btsperiphery.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleBTPMessage is a paid mutator transaction binding the contract method 0xb70eeb8d.
//
// Solidity: function handleBTPMessage(string _from, string _svc, uint256 _sn, bytes _msg) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) HandleBTPMessage(_from string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleBTPMessage(&_Btsperiphery.TransactOpts, _from, _svc, _sn, _msg)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_Btsperiphery *BtsperipheryTransactor) HandleFeeGathering(opts *bind.TransactOpts, _fa string, _svc string) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "handleFeeGathering", _fa, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_Btsperiphery *BtsperipherySession) HandleFeeGathering(_fa string, _svc string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleFeeGathering(&_Btsperiphery.TransactOpts, _fa, _svc)
}

// HandleFeeGathering is a paid mutator transaction binding the contract method 0x3842888c.
//
// Solidity: function handleFeeGathering(string _fa, string _svc) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) HandleFeeGathering(_fa string, _svc string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleFeeGathering(&_Btsperiphery.TransactOpts, _fa, _svc)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_Btsperiphery *BtsperipheryTransactor) HandleRequestService(opts *bind.TransactOpts, _to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "handleRequestService", _to, _assets)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_Btsperiphery *BtsperipherySession) HandleRequestService(_to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleRequestService(&_Btsperiphery.TransactOpts, _to, _assets)
}

// HandleRequestService is a paid mutator transaction binding the contract method 0xdd129575.
//
// Solidity: function handleRequestService(string _to, (string,uint256)[] _assets) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) HandleRequestService(_to string, _assets []TypesAsset) (*types.Transaction, error) {
	return _Btsperiphery.Contract.HandleRequestService(&_Btsperiphery.TransactOpts, _to, _assets)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _bmc, address _btsCore) returns()
func (_Btsperiphery *BtsperipheryTransactor) Initialize(opts *bind.TransactOpts, _bmc common.Address, _btsCore common.Address) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "initialize", _bmc, _btsCore)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _bmc, address _btsCore) returns()
func (_Btsperiphery *BtsperipherySession) Initialize(_bmc common.Address, _btsCore common.Address) (*types.Transaction, error) {
	return _Btsperiphery.Contract.Initialize(&_Btsperiphery.TransactOpts, _bmc, _btsCore)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _bmc, address _btsCore) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) Initialize(_bmc common.Address, _btsCore common.Address) (*types.Transaction, error) {
	return _Btsperiphery.Contract.Initialize(&_Btsperiphery.TransactOpts, _bmc, _btsCore)
}

// RemoveFromBlacklist is a paid mutator transaction binding the contract method 0xc925d633.
//
// Solidity: function removeFromBlacklist(string[] _address) returns()
func (_Btsperiphery *BtsperipheryTransactor) RemoveFromBlacklist(opts *bind.TransactOpts, _address []string) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "removeFromBlacklist", _address)
}

// RemoveFromBlacklist is a paid mutator transaction binding the contract method 0xc925d633.
//
// Solidity: function removeFromBlacklist(string[] _address) returns()
func (_Btsperiphery *BtsperipherySession) RemoveFromBlacklist(_address []string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.RemoveFromBlacklist(&_Btsperiphery.TransactOpts, _address)
}

// RemoveFromBlacklist is a paid mutator transaction binding the contract method 0xc925d633.
//
// Solidity: function removeFromBlacklist(string[] _address) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) RemoveFromBlacklist(_address []string) (*types.Transaction, error) {
	return _Btsperiphery.Contract.RemoveFromBlacklist(&_Btsperiphery.TransactOpts, _address)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_Btsperiphery *BtsperipheryTransactor) SendServiceMessage(opts *bind.TransactOpts, _from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "sendServiceMessage", _from, _to, _coinNames, _values, _fees)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_Btsperiphery *BtsperipherySession) SendServiceMessage(_from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _Btsperiphery.Contract.SendServiceMessage(&_Btsperiphery.TransactOpts, _from, _to, _coinNames, _values, _fees)
}

// SendServiceMessage is a paid mutator transaction binding the contract method 0xd7c37995.
//
// Solidity: function sendServiceMessage(address _from, string _to, string[] _coinNames, uint256[] _values, uint256[] _fees) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) SendServiceMessage(_from common.Address, _to string, _coinNames []string, _values []*big.Int, _fees []*big.Int) (*types.Transaction, error) {
	return _Btsperiphery.Contract.SendServiceMessage(&_Btsperiphery.TransactOpts, _from, _to, _coinNames, _values, _fees)
}

// SetTokenLimit is a paid mutator transaction binding the contract method 0x52647fc4.
//
// Solidity: function setTokenLimit(string[] _coinNames, uint256[] _tokenLimits) returns()
func (_Btsperiphery *BtsperipheryTransactor) SetTokenLimit(opts *bind.TransactOpts, _coinNames []string, _tokenLimits []*big.Int) (*types.Transaction, error) {
	return _Btsperiphery.contract.Transact(opts, "setTokenLimit", _coinNames, _tokenLimits)
}

// SetTokenLimit is a paid mutator transaction binding the contract method 0x52647fc4.
//
// Solidity: function setTokenLimit(string[] _coinNames, uint256[] _tokenLimits) returns()
func (_Btsperiphery *BtsperipherySession) SetTokenLimit(_coinNames []string, _tokenLimits []*big.Int) (*types.Transaction, error) {
	return _Btsperiphery.Contract.SetTokenLimit(&_Btsperiphery.TransactOpts, _coinNames, _tokenLimits)
}

// SetTokenLimit is a paid mutator transaction binding the contract method 0x52647fc4.
//
// Solidity: function setTokenLimit(string[] _coinNames, uint256[] _tokenLimits) returns()
func (_Btsperiphery *BtsperipheryTransactorSession) SetTokenLimit(_coinNames []string, _tokenLimits []*big.Int) (*types.Transaction, error) {
	return _Btsperiphery.Contract.SetTokenLimit(&_Btsperiphery.TransactOpts, _coinNames, _tokenLimits)
}

// BtsperipheryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Btsperiphery contract.
type BtsperipheryInitializedIterator struct {
	Event *BtsperipheryInitialized // Event containing the contract specifics and raw log

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
func (it *BtsperipheryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtsperipheryInitialized)
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
		it.Event = new(BtsperipheryInitialized)
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
func (it *BtsperipheryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtsperipheryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtsperipheryInitialized represents a Initialized event raised by the Btsperiphery contract.
type BtsperipheryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Btsperiphery *BtsperipheryFilterer) FilterInitialized(opts *bind.FilterOpts) (*BtsperipheryInitializedIterator, error) {

	logs, sub, err := _Btsperiphery.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BtsperipheryInitializedIterator{contract: _Btsperiphery.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Btsperiphery *BtsperipheryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BtsperipheryInitialized) (event.Subscription, error) {

	logs, sub, err := _Btsperiphery.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtsperipheryInitialized)
				if err := _Btsperiphery.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Btsperiphery *BtsperipheryFilterer) ParseInitialized(log types.Log) (*BtsperipheryInitialized, error) {
	event := new(BtsperipheryInitialized)
	if err := _Btsperiphery.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtsperipheryTransferEndIterator is returned from FilterTransferEnd and is used to iterate over the raw logs and unpacked data for TransferEnd events raised by the Btsperiphery contract.
type BtsperipheryTransferEndIterator struct {
	Event *BtsperipheryTransferEnd // Event containing the contract specifics and raw log

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
func (it *BtsperipheryTransferEndIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtsperipheryTransferEnd)
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
		it.Event = new(BtsperipheryTransferEnd)
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
func (it *BtsperipheryTransferEndIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtsperipheryTransferEndIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtsperipheryTransferEnd represents a TransferEnd event raised by the Btsperiphery contract.
type BtsperipheryTransferEnd struct {
	From     common.Address
	Sn       *big.Int
	Code     *big.Int
	Response string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferEnd is a free log retrieval operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_Btsperiphery *BtsperipheryFilterer) FilterTransferEnd(opts *bind.FilterOpts, _from []common.Address) (*BtsperipheryTransferEndIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _Btsperiphery.contract.FilterLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return &BtsperipheryTransferEndIterator{contract: _Btsperiphery.contract, event: "TransferEnd", logs: logs, sub: sub}, nil
}

// WatchTransferEnd is a free log subscription operation binding the contract event 0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2.
//
// Solidity: event TransferEnd(address indexed _from, uint256 _sn, uint256 _code, string _response)
func (_Btsperiphery *BtsperipheryFilterer) WatchTransferEnd(opts *bind.WatchOpts, sink chan<- *BtsperipheryTransferEnd, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _Btsperiphery.contract.WatchLogs(opts, "TransferEnd", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtsperipheryTransferEnd)
				if err := _Btsperiphery.contract.UnpackLog(event, "TransferEnd", log); err != nil {
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
func (_Btsperiphery *BtsperipheryFilterer) ParseTransferEnd(log types.Log) (*BtsperipheryTransferEnd, error) {
	event := new(BtsperipheryTransferEnd)
	if err := _Btsperiphery.contract.UnpackLog(event, "TransferEnd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtsperipheryTransferReceivedIterator is returned from FilterTransferReceived and is used to iterate over the raw logs and unpacked data for TransferReceived events raised by the Btsperiphery contract.
type BtsperipheryTransferReceivedIterator struct {
	Event *BtsperipheryTransferReceived // Event containing the contract specifics and raw log

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
func (it *BtsperipheryTransferReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtsperipheryTransferReceived)
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
		it.Event = new(BtsperipheryTransferReceived)
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
func (it *BtsperipheryTransferReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtsperipheryTransferReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtsperipheryTransferReceived represents a TransferReceived event raised by the Btsperiphery contract.
type BtsperipheryTransferReceived struct {
	From         common.Hash
	To           common.Address
	Sn           *big.Int
	AssetDetails []TypesAsset
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferReceived is a free log retrieval operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_Btsperiphery *BtsperipheryFilterer) FilterTransferReceived(opts *bind.FilterOpts, _from []string, _to []common.Address) (*BtsperipheryTransferReceivedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Btsperiphery.contract.FilterLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &BtsperipheryTransferReceivedIterator{contract: _Btsperiphery.contract, event: "TransferReceived", logs: logs, sub: sub}, nil
}

// WatchTransferReceived is a free log subscription operation binding the contract event 0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680.
//
// Solidity: event TransferReceived(string indexed _from, address indexed _to, uint256 _sn, (string,uint256)[] _assetDetails)
func (_Btsperiphery *BtsperipheryFilterer) WatchTransferReceived(opts *bind.WatchOpts, sink chan<- *BtsperipheryTransferReceived, _from []string, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _Btsperiphery.contract.WatchLogs(opts, "TransferReceived", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtsperipheryTransferReceived)
				if err := _Btsperiphery.contract.UnpackLog(event, "TransferReceived", log); err != nil {
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
func (_Btsperiphery *BtsperipheryFilterer) ParseTransferReceived(log types.Log) (*BtsperipheryTransferReceived, error) {
	event := new(BtsperipheryTransferReceived)
	if err := _Btsperiphery.contract.UnpackLog(event, "TransferReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtsperipheryTransferStartIterator is returned from FilterTransferStart and is used to iterate over the raw logs and unpacked data for TransferStart events raised by the Btsperiphery contract.
type BtsperipheryTransferStartIterator struct {
	Event *BtsperipheryTransferStart // Event containing the contract specifics and raw log

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
func (it *BtsperipheryTransferStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtsperipheryTransferStart)
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
		it.Event = new(BtsperipheryTransferStart)
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
func (it *BtsperipheryTransferStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtsperipheryTransferStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtsperipheryTransferStart represents a TransferStart event raised by the Btsperiphery contract.
type BtsperipheryTransferStart struct {
	From         common.Address
	To           string
	Sn           *big.Int
	AssetDetails []TypesAssetTransferDetail
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferStart is a free log retrieval operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_Btsperiphery *BtsperipheryFilterer) FilterTransferStart(opts *bind.FilterOpts, _from []common.Address) (*BtsperipheryTransferStartIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _Btsperiphery.contract.FilterLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return &BtsperipheryTransferStartIterator{contract: _Btsperiphery.contract, event: "TransferStart", logs: logs, sub: sub}, nil
}

// WatchTransferStart is a free log subscription operation binding the contract event 0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a.
//
// Solidity: event TransferStart(address indexed _from, string _to, uint256 _sn, (string,uint256,uint256)[] _assetDetails)
func (_Btsperiphery *BtsperipheryFilterer) WatchTransferStart(opts *bind.WatchOpts, sink chan<- *BtsperipheryTransferStart, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _Btsperiphery.contract.WatchLogs(opts, "TransferStart", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtsperipheryTransferStart)
				if err := _Btsperiphery.contract.UnpackLog(event, "TransferStart", log); err != nil {
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
func (_Btsperiphery *BtsperipheryFilterer) ParseTransferStart(log types.Log) (*BtsperipheryTransferStart, error) {
	event := new(BtsperipheryTransferStart)
	if err := _Btsperiphery.contract.UnpackLog(event, "TransferStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtsperipheryUnknownResponseIterator is returned from FilterUnknownResponse and is used to iterate over the raw logs and unpacked data for UnknownResponse events raised by the Btsperiphery contract.
type BtsperipheryUnknownResponseIterator struct {
	Event *BtsperipheryUnknownResponse // Event containing the contract specifics and raw log

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
func (it *BtsperipheryUnknownResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtsperipheryUnknownResponse)
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
		it.Event = new(BtsperipheryUnknownResponse)
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
func (it *BtsperipheryUnknownResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtsperipheryUnknownResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtsperipheryUnknownResponse represents a UnknownResponse event raised by the Btsperiphery contract.
type BtsperipheryUnknownResponse struct {
	From string
	Sn   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUnknownResponse is a free log retrieval operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_Btsperiphery *BtsperipheryFilterer) FilterUnknownResponse(opts *bind.FilterOpts) (*BtsperipheryUnknownResponseIterator, error) {

	logs, sub, err := _Btsperiphery.contract.FilterLogs(opts, "UnknownResponse")
	if err != nil {
		return nil, err
	}
	return &BtsperipheryUnknownResponseIterator{contract: _Btsperiphery.contract, event: "UnknownResponse", logs: logs, sub: sub}, nil
}

// WatchUnknownResponse is a free log subscription operation binding the contract event 0x0e2e04e992df368a336276b84416f1b66f8aaca143ea47e284b229cc9f10a889.
//
// Solidity: event UnknownResponse(string _from, uint256 _sn)
func (_Btsperiphery *BtsperipheryFilterer) WatchUnknownResponse(opts *bind.WatchOpts, sink chan<- *BtsperipheryUnknownResponse) (event.Subscription, error) {

	logs, sub, err := _Btsperiphery.contract.WatchLogs(opts, "UnknownResponse")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtsperipheryUnknownResponse)
				if err := _Btsperiphery.contract.UnpackLog(event, "UnknownResponse", log); err != nil {
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
func (_Btsperiphery *BtsperipheryFilterer) ParseUnknownResponse(log types.Log) (*BtsperipheryUnknownResponse, error) {
	event := new(BtsperipheryUnknownResponse)
	if err := _Btsperiphery.contract.UnpackLog(event, "UnknownResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
