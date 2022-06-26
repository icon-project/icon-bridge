// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ownerTokenHmy

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

// OwnerTokenHmyABI is the input ABI used to generate the binding from.
const OwnerTokenHmyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"Register\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"formerOwner\",\"type\":\"address\"}],\"name\":\"RemoveOwnership\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"promoter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwnership\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bshImpl\",\"type\":\"address\"}],\"name\":\"updateBSHImplementation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"}],\"name\":\"setFeeRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_names\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_tokenName\",\"type\":\"string\"}],\"name\":\"getBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_usableBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lockedBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_refundableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccumulatedFees\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"collectedFees\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_tokenName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_tokenName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"calculateTransferFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_toFA\",\"type\":\"string\"}],\"name\":\"handleFeeTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"}],\"name\":\"handleResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_toAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_tokenName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"handleTransferRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_tokenName\",\"type\":\"string\"}],\"name\":\"isTokenRegisterd\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_registered\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OwnerTokenHmy is an auto generated Go binding around an Ethereum contract.
type OwnerTokenHmy struct {
	OwnerTokenHmyCaller     // Read-only binding to the contract
	OwnerTokenHmyTransactor // Write-only binding to the contract
	OwnerTokenHmyFilterer   // Log filterer for contract events
}

// OwnerTokenHmyCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnerTokenHmyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerTokenHmyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnerTokenHmyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerTokenHmyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnerTokenHmyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerTokenHmySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnerTokenHmySession struct {
	Contract     *OwnerTokenHmy    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnerTokenHmyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnerTokenHmyCallerSession struct {
	Contract *OwnerTokenHmyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// OwnerTokenHmyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnerTokenHmyTransactorSession struct {
	Contract     *OwnerTokenHmyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// OwnerTokenHmyRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnerTokenHmyRaw struct {
	Contract *OwnerTokenHmy // Generic contract binding to access the raw methods on
}

// OwnerTokenHmyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnerTokenHmyCallerRaw struct {
	Contract *OwnerTokenHmyCaller // Generic read-only contract binding to access the raw methods on
}

// OwnerTokenHmyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnerTokenHmyTransactorRaw struct {
	Contract *OwnerTokenHmyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnerTokenHmy creates a new instance of OwnerTokenHmy, bound to a specific deployed contract.
func NewOwnerTokenHmy(address common.Address, backend bind.ContractBackend) (*OwnerTokenHmy, error) {
	contract, err := bindOwnerTokenHmy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmy{OwnerTokenHmyCaller: OwnerTokenHmyCaller{contract: contract}, OwnerTokenHmyTransactor: OwnerTokenHmyTransactor{contract: contract}, OwnerTokenHmyFilterer: OwnerTokenHmyFilterer{contract: contract}}, nil
}

// NewOwnerTokenHmyCaller creates a new read-only instance of OwnerTokenHmy, bound to a specific deployed contract.
func NewOwnerTokenHmyCaller(address common.Address, caller bind.ContractCaller) (*OwnerTokenHmyCaller, error) {
	contract, err := bindOwnerTokenHmy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmyCaller{contract: contract}, nil
}

// NewOwnerTokenHmyTransactor creates a new write-only instance of OwnerTokenHmy, bound to a specific deployed contract.
func NewOwnerTokenHmyTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnerTokenHmyTransactor, error) {
	contract, err := bindOwnerTokenHmy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmyTransactor{contract: contract}, nil
}

// NewOwnerTokenHmyFilterer creates a new log filterer instance of OwnerTokenHmy, bound to a specific deployed contract.
func NewOwnerTokenHmyFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnerTokenHmyFilterer, error) {
	contract, err := bindOwnerTokenHmy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmyFilterer{contract: contract}, nil
}

// bindOwnerTokenHmy binds a generic wrapper to an already deployed contract.
func bindOwnerTokenHmy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnerTokenHmyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerTokenHmy *OwnerTokenHmyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerTokenHmy.Contract.OwnerTokenHmyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerTokenHmy *OwnerTokenHmyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.OwnerTokenHmyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerTokenHmy *OwnerTokenHmyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.OwnerTokenHmyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerTokenHmy *OwnerTokenHmyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerTokenHmy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerTokenHmy *OwnerTokenHmyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerTokenHmy *OwnerTokenHmyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.contract.Transact(opts, method, params...)
}

// CalculateTransferFee is a free data retrieval call binding the contract method 0x173652b0.
//
// Solidity: function calculateTransferFee(address token_addr, uint256 _value) view returns(uint256 value, uint256 fee)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) CalculateTransferFee(opts *bind.CallOpts, token_addr common.Address, _value *big.Int) (struct {
	Value *big.Int
	Fee   *big.Int
}, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "calculateTransferFee", token_addr, _value)

	outstruct := new(struct {
		Value *big.Int
		Fee   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Value = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CalculateTransferFee is a free data retrieval call binding the contract method 0x173652b0.
//
// Solidity: function calculateTransferFee(address token_addr, uint256 _value) view returns(uint256 value, uint256 fee)
func (_OwnerTokenHmy *OwnerTokenHmySession) CalculateTransferFee(token_addr common.Address, _value *big.Int) (struct {
	Value *big.Int
	Fee   *big.Int
}, error) {
	return _OwnerTokenHmy.Contract.CalculateTransferFee(&_OwnerTokenHmy.CallOpts, token_addr, _value)
}

// CalculateTransferFee is a free data retrieval call binding the contract method 0x173652b0.
//
// Solidity: function calculateTransferFee(address token_addr, uint256 _value) view returns(uint256 value, uint256 fee)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) CalculateTransferFee(token_addr common.Address, _value *big.Int) (struct {
	Value *big.Int
	Fee   *big.Int
}, error) {
	return _OwnerTokenHmy.Contract.CalculateTransferFee(&_OwnerTokenHmy.CallOpts, token_addr, _value)
}

// FeeCollector is a free data retrieval call binding the contract method 0x243b14cf.
//
// Solidity: function feeCollector(string ) view returns(uint256)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) FeeCollector(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "feeCollector", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeCollector is a free data retrieval call binding the contract method 0x243b14cf.
//
// Solidity: function feeCollector(string ) view returns(uint256)
func (_OwnerTokenHmy *OwnerTokenHmySession) FeeCollector(arg0 string) (*big.Int, error) {
	return _OwnerTokenHmy.Contract.FeeCollector(&_OwnerTokenHmy.CallOpts, arg0)
}

// FeeCollector is a free data retrieval call binding the contract method 0x243b14cf.
//
// Solidity: function feeCollector(string ) view returns(uint256)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) FeeCollector(arg0 string) (*big.Int, error) {
	return _OwnerTokenHmy.Contract.FeeCollector(&_OwnerTokenHmy.CallOpts, arg0)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256,uint256)[] collectedFees)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) GetAccumulatedFees(opts *bind.CallOpts) ([]TypesAsset, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "getAccumulatedFees")

	if err != nil {
		return *new([]TypesAsset), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesAsset)).(*[]TypesAsset)

	return out0, err

}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256,uint256)[] collectedFees)
func (_OwnerTokenHmy *OwnerTokenHmySession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _OwnerTokenHmy.Contract.GetAccumulatedFees(&_OwnerTokenHmy.CallOpts)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256,uint256)[] collectedFees)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _OwnerTokenHmy.Contract.GetAccumulatedFees(&_OwnerTokenHmy.CallOpts)
}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _tokenName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) GetBalanceOf(opts *bind.CallOpts, _owner common.Address, _tokenName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "getBalanceOf", _owner, _tokenName)

	outstruct := new(struct {
		UsableBalance     *big.Int
		LockedBalance     *big.Int
		RefundableBalance *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UsableBalance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LockedBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RefundableBalance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _tokenName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_OwnerTokenHmy *OwnerTokenHmySession) GetBalanceOf(_owner common.Address, _tokenName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	return _OwnerTokenHmy.Contract.GetBalanceOf(&_OwnerTokenHmy.CallOpts, _owner, _tokenName)
}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _tokenName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) GetBalanceOf(_owner common.Address, _tokenName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	return _OwnerTokenHmy.Contract.GetBalanceOf(&_OwnerTokenHmy.CallOpts, _owner, _tokenName)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_OwnerTokenHmy *OwnerTokenHmyCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_OwnerTokenHmy *OwnerTokenHmySession) GetOwners() ([]common.Address, error) {
	return _OwnerTokenHmy.Contract.GetOwners(&_OwnerTokenHmy.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) GetOwners() ([]common.Address, error) {
	return _OwnerTokenHmy.Contract.GetOwners(&_OwnerTokenHmy.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) IsOwner(opts *bind.CallOpts, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "isOwner", _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_OwnerTokenHmy *OwnerTokenHmySession) IsOwner(_owner common.Address) (bool, error) {
	return _OwnerTokenHmy.Contract.IsOwner(&_OwnerTokenHmy.CallOpts, _owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) IsOwner(_owner common.Address) (bool, error) {
	return _OwnerTokenHmy.Contract.IsOwner(&_OwnerTokenHmy.CallOpts, _owner)
}

// IsTokenRegisterd is a free data retrieval call binding the contract method 0x8dd79d51.
//
// Solidity: function isTokenRegisterd(string _tokenName) view returns(bool _registered)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) IsTokenRegisterd(opts *bind.CallOpts, _tokenName string) (bool, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "isTokenRegisterd", _tokenName)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTokenRegisterd is a free data retrieval call binding the contract method 0x8dd79d51.
//
// Solidity: function isTokenRegisterd(string _tokenName) view returns(bool _registered)
func (_OwnerTokenHmy *OwnerTokenHmySession) IsTokenRegisterd(_tokenName string) (bool, error) {
	return _OwnerTokenHmy.Contract.IsTokenRegisterd(&_OwnerTokenHmy.CallOpts, _tokenName)
}

// IsTokenRegisterd is a free data retrieval call binding the contract method 0x8dd79d51.
//
// Solidity: function isTokenRegisterd(string _tokenName) view returns(bool _registered)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) IsTokenRegisterd(_tokenName string) (bool, error) {
	return _OwnerTokenHmy.Contract.IsTokenRegisterd(&_OwnerTokenHmy.CallOpts, _tokenName)
}

// TokenNames is a free data retrieval call binding the contract method 0x188e7852.
//
// Solidity: function tokenNames() view returns(string[] _names)
func (_OwnerTokenHmy *OwnerTokenHmyCaller) TokenNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _OwnerTokenHmy.contract.Call(opts, &out, "tokenNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// TokenNames is a free data retrieval call binding the contract method 0x188e7852.
//
// Solidity: function tokenNames() view returns(string[] _names)
func (_OwnerTokenHmy *OwnerTokenHmySession) TokenNames() ([]string, error) {
	return _OwnerTokenHmy.Contract.TokenNames(&_OwnerTokenHmy.CallOpts)
}

// TokenNames is a free data retrieval call binding the contract method 0x188e7852.
//
// Solidity: function tokenNames() view returns(string[] _names)
func (_OwnerTokenHmy *OwnerTokenHmyCallerSession) TokenNames() ([]string, error) {
	return _OwnerTokenHmy.Contract.TokenNames(&_OwnerTokenHmy.CallOpts)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) AddOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "addOwner", _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.AddOwner(&_OwnerTokenHmy.TransactOpts, _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.AddOwner(&_OwnerTokenHmy.TransactOpts, _owner)
}

// HandleFeeTransfer is a paid mutator transaction binding the contract method 0xc7e07860.
//
// Solidity: function handleFeeTransfer(string _toFA) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) HandleFeeTransfer(opts *bind.TransactOpts, _toFA string) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "handleFeeTransfer", _toFA)
}

// HandleFeeTransfer is a paid mutator transaction binding the contract method 0xc7e07860.
//
// Solidity: function handleFeeTransfer(string _toFA) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) HandleFeeTransfer(_toFA string) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.HandleFeeTransfer(&_OwnerTokenHmy.TransactOpts, _toFA)
}

// HandleFeeTransfer is a paid mutator transaction binding the contract method 0xc7e07860.
//
// Solidity: function handleFeeTransfer(string _toFA) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) HandleFeeTransfer(_toFA string) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.HandleFeeTransfer(&_OwnerTokenHmy.TransactOpts, _toFA)
}

// HandleResponse is a paid mutator transaction binding the contract method 0xf82d7eb1.
//
// Solidity: function handleResponse(address _caller, (string,uint256,uint256)[] _assets, uint256 _code) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) HandleResponse(opts *bind.TransactOpts, _caller common.Address, _assets []TypesAsset, _code *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "handleResponse", _caller, _assets, _code)
}

// HandleResponse is a paid mutator transaction binding the contract method 0xf82d7eb1.
//
// Solidity: function handleResponse(address _caller, (string,uint256,uint256)[] _assets, uint256 _code) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) HandleResponse(_caller common.Address, _assets []TypesAsset, _code *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.HandleResponse(&_OwnerTokenHmy.TransactOpts, _caller, _assets, _code)
}

// HandleResponse is a paid mutator transaction binding the contract method 0xf82d7eb1.
//
// Solidity: function handleResponse(address _caller, (string,uint256,uint256)[] _assets, uint256 _code) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) HandleResponse(_caller common.Address, _assets []TypesAsset, _code *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.HandleResponse(&_OwnerTokenHmy.TransactOpts, _caller, _assets, _code)
}

// HandleTransferRequest is a paid mutator transaction binding the contract method 0x5d883186.
//
// Solidity: function handleTransferRequest(address _toAddress, string _tokenName, uint256 _amount) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) HandleTransferRequest(opts *bind.TransactOpts, _toAddress common.Address, _tokenName string, _amount *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "handleTransferRequest", _toAddress, _tokenName, _amount)
}

// HandleTransferRequest is a paid mutator transaction binding the contract method 0x5d883186.
//
// Solidity: function handleTransferRequest(address _toAddress, string _tokenName, uint256 _amount) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) HandleTransferRequest(_toAddress common.Address, _tokenName string, _amount *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.HandleTransferRequest(&_OwnerTokenHmy.TransactOpts, _toAddress, _tokenName, _amount)
}

// HandleTransferRequest is a paid mutator transaction binding the contract method 0x5d883186.
//
// Solidity: function handleTransferRequest(address _toAddress, string _tokenName, uint256 _amount) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) HandleTransferRequest(_toAddress common.Address, _tokenName string, _amount *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.HandleTransferRequest(&_OwnerTokenHmy.TransactOpts, _toAddress, _tokenName, _amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 _feeNumerator) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) Initialize(opts *bind.TransactOpts, _feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "initialize", _feeNumerator)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 _feeNumerator) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) Initialize(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Initialize(&_OwnerTokenHmy.TransactOpts, _feeNumerator)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 _feeNumerator) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) Initialize(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Initialize(&_OwnerTokenHmy.TransactOpts, _feeNumerator)
}

// Register is a paid mutator transaction binding the contract method 0xf63327ee.
//
// Solidity: function register(string _name, string _symbol, uint256 _decimals, uint256 _feeNumerator, address _addr) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) Register(opts *bind.TransactOpts, _name string, _symbol string, _decimals *big.Int, _feeNumerator *big.Int, _addr common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "register", _name, _symbol, _decimals, _feeNumerator, _addr)
}

// Register is a paid mutator transaction binding the contract method 0xf63327ee.
//
// Solidity: function register(string _name, string _symbol, uint256 _decimals, uint256 _feeNumerator, address _addr) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) Register(_name string, _symbol string, _decimals *big.Int, _feeNumerator *big.Int, _addr common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Register(&_OwnerTokenHmy.TransactOpts, _name, _symbol, _decimals, _feeNumerator, _addr)
}

// Register is a paid mutator transaction binding the contract method 0xf63327ee.
//
// Solidity: function register(string _name, string _symbol, uint256 _decimals, uint256 _feeNumerator, address _addr) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) Register(_name string, _symbol string, _decimals *big.Int, _feeNumerator *big.Int, _addr common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Register(&_OwnerTokenHmy.TransactOpts, _name, _symbol, _decimals, _feeNumerator, _addr)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) RemoveOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "removeOwner", _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.RemoveOwner(&_OwnerTokenHmy.TransactOpts, _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.RemoveOwner(&_OwnerTokenHmy.TransactOpts, _owner)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) SetFeeRatio(opts *bind.TransactOpts, _feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "setFeeRatio", _feeNumerator)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) SetFeeRatio(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.SetFeeRatio(&_OwnerTokenHmy.TransactOpts, _feeNumerator)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) SetFeeRatio(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.SetFeeRatio(&_OwnerTokenHmy.TransactOpts, _feeNumerator)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _tokenName, uint256 _value, string _to) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) Transfer(opts *bind.TransactOpts, _tokenName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "transfer", _tokenName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _tokenName, uint256 _value, string _to) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) Transfer(_tokenName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Transfer(&_OwnerTokenHmy.TransactOpts, _tokenName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _tokenName, uint256 _value, string _to) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) Transfer(_tokenName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Transfer(&_OwnerTokenHmy.TransactOpts, _tokenName, _value, _to)
}

// UpdateBSHImplementation is a paid mutator transaction binding the contract method 0x5f58b429.
//
// Solidity: function updateBSHImplementation(address _bshImpl) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) UpdateBSHImplementation(opts *bind.TransactOpts, _bshImpl common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "updateBSHImplementation", _bshImpl)
}

// UpdateBSHImplementation is a paid mutator transaction binding the contract method 0x5f58b429.
//
// Solidity: function updateBSHImplementation(address _bshImpl) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) UpdateBSHImplementation(_bshImpl common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.UpdateBSHImplementation(&_OwnerTokenHmy.TransactOpts, _bshImpl)
}

// UpdateBSHImplementation is a paid mutator transaction binding the contract method 0x5f58b429.
//
// Solidity: function updateBSHImplementation(address _bshImpl) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) UpdateBSHImplementation(_bshImpl common.Address) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.UpdateBSHImplementation(&_OwnerTokenHmy.TransactOpts, _bshImpl)
}

// Withdraw is a paid mutator transaction binding the contract method 0x30b39a62.
//
// Solidity: function withdraw(string _tokenName, uint256 _value) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactor) Withdraw(opts *bind.TransactOpts, _tokenName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.contract.Transact(opts, "withdraw", _tokenName, _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x30b39a62.
//
// Solidity: function withdraw(string _tokenName, uint256 _value) returns()
func (_OwnerTokenHmy *OwnerTokenHmySession) Withdraw(_tokenName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Withdraw(&_OwnerTokenHmy.TransactOpts, _tokenName, _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x30b39a62.
//
// Solidity: function withdraw(string _tokenName, uint256 _value) returns()
func (_OwnerTokenHmy *OwnerTokenHmyTransactorSession) Withdraw(_tokenName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerTokenHmy.Contract.Withdraw(&_OwnerTokenHmy.TransactOpts, _tokenName, _value)
}

// OwnerTokenHmyRegisterIterator is returned from FilterRegister and is used to iterate over the raw logs and unpacked data for Register events raised by the OwnerTokenHmy contract.
type OwnerTokenHmyRegisterIterator struct {
	Event *OwnerTokenHmyRegister // Event containing the contract specifics and raw log

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
func (it *OwnerTokenHmyRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerTokenHmyRegister)
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
		it.Event = new(OwnerTokenHmyRegister)
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
func (it *OwnerTokenHmyRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerTokenHmyRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerTokenHmyRegister represents a Register event raised by the OwnerTokenHmy contract.
type OwnerTokenHmyRegister struct {
	Name common.Hash
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRegister is a free log retrieval operation binding the contract event 0x3e49fb7efebefffe4fc2e2c193820fac9b11de4bf185e0d14de13e0068c2ac34.
//
// Solidity: event Register(string indexed name, address addr)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) FilterRegister(opts *bind.FilterOpts, name []string) (*OwnerTokenHmyRegisterIterator, error) {

	var nameRule []interface{}
	for _, nameItem := range name {
		nameRule = append(nameRule, nameItem)
	}

	logs, sub, err := _OwnerTokenHmy.contract.FilterLogs(opts, "Register", nameRule)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmyRegisterIterator{contract: _OwnerTokenHmy.contract, event: "Register", logs: logs, sub: sub}, nil
}

// WatchRegister is a free log subscription operation binding the contract event 0x3e49fb7efebefffe4fc2e2c193820fac9b11de4bf185e0d14de13e0068c2ac34.
//
// Solidity: event Register(string indexed name, address addr)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) WatchRegister(opts *bind.WatchOpts, sink chan<- *OwnerTokenHmyRegister, name []string) (event.Subscription, error) {

	var nameRule []interface{}
	for _, nameItem := range name {
		nameRule = append(nameRule, nameItem)
	}

	logs, sub, err := _OwnerTokenHmy.contract.WatchLogs(opts, "Register", nameRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerTokenHmyRegister)
				if err := _OwnerTokenHmy.contract.UnpackLog(event, "Register", log); err != nil {
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

// ParseRegister is a log parse operation binding the contract event 0x3e49fb7efebefffe4fc2e2c193820fac9b11de4bf185e0d14de13e0068c2ac34.
//
// Solidity: event Register(string indexed name, address addr)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) ParseRegister(log types.Log) (*OwnerTokenHmyRegister, error) {
	event := new(OwnerTokenHmyRegister)
	if err := _OwnerTokenHmy.contract.UnpackLog(event, "Register", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnerTokenHmyRemoveOwnershipIterator is returned from FilterRemoveOwnership and is used to iterate over the raw logs and unpacked data for RemoveOwnership events raised by the OwnerTokenHmy contract.
type OwnerTokenHmyRemoveOwnershipIterator struct {
	Event *OwnerTokenHmyRemoveOwnership // Event containing the contract specifics and raw log

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
func (it *OwnerTokenHmyRemoveOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerTokenHmyRemoveOwnership)
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
		it.Event = new(OwnerTokenHmyRemoveOwnership)
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
func (it *OwnerTokenHmyRemoveOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerTokenHmyRemoveOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerTokenHmyRemoveOwnership represents a RemoveOwnership event raised by the OwnerTokenHmy contract.
type OwnerTokenHmyRemoveOwnership struct {
	Remover     common.Address
	FormerOwner common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveOwnership is a free log retrieval operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) FilterRemoveOwnership(opts *bind.FilterOpts, remover []common.Address, formerOwner []common.Address) (*OwnerTokenHmyRemoveOwnershipIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _OwnerTokenHmy.contract.FilterLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmyRemoveOwnershipIterator{contract: _OwnerTokenHmy.contract, event: "RemoveOwnership", logs: logs, sub: sub}, nil
}

// WatchRemoveOwnership is a free log subscription operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) WatchRemoveOwnership(opts *bind.WatchOpts, sink chan<- *OwnerTokenHmyRemoveOwnership, remover []common.Address, formerOwner []common.Address) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _OwnerTokenHmy.contract.WatchLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerTokenHmyRemoveOwnership)
				if err := _OwnerTokenHmy.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
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

// ParseRemoveOwnership is a log parse operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) ParseRemoveOwnership(log types.Log) (*OwnerTokenHmyRemoveOwnership, error) {
	event := new(OwnerTokenHmyRemoveOwnership)
	if err := _OwnerTokenHmy.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnerTokenHmySetOwnershipIterator is returned from FilterSetOwnership and is used to iterate over the raw logs and unpacked data for SetOwnership events raised by the OwnerTokenHmy contract.
type OwnerTokenHmySetOwnershipIterator struct {
	Event *OwnerTokenHmySetOwnership // Event containing the contract specifics and raw log

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
func (it *OwnerTokenHmySetOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerTokenHmySetOwnership)
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
		it.Event = new(OwnerTokenHmySetOwnership)
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
func (it *OwnerTokenHmySetOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerTokenHmySetOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerTokenHmySetOwnership represents a SetOwnership event raised by the OwnerTokenHmy contract.
type OwnerTokenHmySetOwnership struct {
	Promoter common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetOwnership is a free log retrieval operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) FilterSetOwnership(opts *bind.FilterOpts, promoter []common.Address, newOwner []common.Address) (*OwnerTokenHmySetOwnershipIterator, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnerTokenHmy.contract.FilterLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnerTokenHmySetOwnershipIterator{contract: _OwnerTokenHmy.contract, event: "SetOwnership", logs: logs, sub: sub}, nil
}

// WatchSetOwnership is a free log subscription operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) WatchSetOwnership(opts *bind.WatchOpts, sink chan<- *OwnerTokenHmySetOwnership, promoter []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnerTokenHmy.contract.WatchLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerTokenHmySetOwnership)
				if err := _OwnerTokenHmy.contract.UnpackLog(event, "SetOwnership", log); err != nil {
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

// ParseSetOwnership is a log parse operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_OwnerTokenHmy *OwnerTokenHmyFilterer) ParseSetOwnership(log types.Log) (*OwnerTokenHmySetOwnership, error) {
	event := new(OwnerTokenHmySetOwnership)
	if err := _OwnerTokenHmy.contract.UnpackLog(event, "SetOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
