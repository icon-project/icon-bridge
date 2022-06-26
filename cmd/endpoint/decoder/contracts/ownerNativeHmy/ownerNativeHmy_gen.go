// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ownerNativeHmy

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

// OwnerNativeHmyABI is the input ABI used to generate the binding from.
const OwnerNativeHmyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"formerOwner\",\"type\":\"address\"}],\"name\":\"RemoveOwnership\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"promoter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwnership\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"feeNumerator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fixedFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_nativeCoinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bshPeriphery\",\"type\":\"address\"}],\"name\":\"updateBSHPeriphery\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"}],\"name\":\"setFeeRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"name\":\"setFixedFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coinNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_names\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"coinId\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"isValidCoin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_valid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"getBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_usableBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lockedBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_refundableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"}],\"name\":\"getBalanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_usableBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_lockedBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_refundableBalances\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccumulatedFees\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_accumulatedFees\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transferNativeCoin\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transferBatch\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rspCode\",\"type\":\"uint256\"}],\"name\":\"handleResponseService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_fa\",\"type\":\"string\"}],\"name\":\"transferFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnerNativeHmy is an auto generated Go binding around an Ethereum contract.
type OwnerNativeHmy struct {
	OwnerNativeHmyCaller     // Read-only binding to the contract
	OwnerNativeHmyTransactor // Write-only binding to the contract
	OwnerNativeHmyFilterer   // Log filterer for contract events
}

// OwnerNativeHmyCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnerNativeHmyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerNativeHmyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnerNativeHmyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerNativeHmyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnerNativeHmyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnerNativeHmySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnerNativeHmySession struct {
	Contract     *OwnerNativeHmy   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnerNativeHmyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnerNativeHmyCallerSession struct {
	Contract *OwnerNativeHmyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OwnerNativeHmyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnerNativeHmyTransactorSession struct {
	Contract     *OwnerNativeHmyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OwnerNativeHmyRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnerNativeHmyRaw struct {
	Contract *OwnerNativeHmy // Generic contract binding to access the raw methods on
}

// OwnerNativeHmyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnerNativeHmyCallerRaw struct {
	Contract *OwnerNativeHmyCaller // Generic read-only contract binding to access the raw methods on
}

// OwnerNativeHmyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnerNativeHmyTransactorRaw struct {
	Contract *OwnerNativeHmyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnerNativeHmy creates a new instance of OwnerNativeHmy, bound to a specific deployed contract.
func NewOwnerNativeHmy(address common.Address, backend bind.ContractBackend) (*OwnerNativeHmy, error) {
	contract, err := bindOwnerNativeHmy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnerNativeHmy{OwnerNativeHmyCaller: OwnerNativeHmyCaller{contract: contract}, OwnerNativeHmyTransactor: OwnerNativeHmyTransactor{contract: contract}, OwnerNativeHmyFilterer: OwnerNativeHmyFilterer{contract: contract}}, nil
}

// NewOwnerNativeHmyCaller creates a new read-only instance of OwnerNativeHmy, bound to a specific deployed contract.
func NewOwnerNativeHmyCaller(address common.Address, caller bind.ContractCaller) (*OwnerNativeHmyCaller, error) {
	contract, err := bindOwnerNativeHmy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerNativeHmyCaller{contract: contract}, nil
}

// NewOwnerNativeHmyTransactor creates a new write-only instance of OwnerNativeHmy, bound to a specific deployed contract.
func NewOwnerNativeHmyTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnerNativeHmyTransactor, error) {
	contract, err := bindOwnerNativeHmy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnerNativeHmyTransactor{contract: contract}, nil
}

// NewOwnerNativeHmyFilterer creates a new log filterer instance of OwnerNativeHmy, bound to a specific deployed contract.
func NewOwnerNativeHmyFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnerNativeHmyFilterer, error) {
	contract, err := bindOwnerNativeHmy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnerNativeHmyFilterer{contract: contract}, nil
}

// bindOwnerNativeHmy binds a generic wrapper to an already deployed contract.
func bindOwnerNativeHmy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnerNativeHmyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerNativeHmy *OwnerNativeHmyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerNativeHmy.Contract.OwnerNativeHmyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerNativeHmy *OwnerNativeHmyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.OwnerNativeHmyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerNativeHmy *OwnerNativeHmyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.OwnerNativeHmyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnerNativeHmy *OwnerNativeHmyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnerNativeHmy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnerNativeHmy *OwnerNativeHmyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnerNativeHmy *OwnerNativeHmyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.contract.Transact(opts, method, params...)
}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) CoinId(opts *bind.CallOpts, _coinName string) (common.Address, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "coinId", _coinName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_OwnerNativeHmy *OwnerNativeHmySession) CoinId(_coinName string) (common.Address, error) {
	return _OwnerNativeHmy.Contract.CoinId(&_OwnerNativeHmy.CallOpts, _coinName)
}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) CoinId(_coinName string) (common.Address, error) {
	return _OwnerNativeHmy.Contract.CoinId(&_OwnerNativeHmy.CallOpts, _coinName)
}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) CoinNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "coinNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_OwnerNativeHmy *OwnerNativeHmySession) CoinNames() ([]string, error) {
	return _OwnerNativeHmy.Contract.CoinNames(&_OwnerNativeHmy.CallOpts)
}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) CoinNames() ([]string, error) {
	return _OwnerNativeHmy.Contract.CoinNames(&_OwnerNativeHmy.CallOpts)
}

// FeeNumerator is a free data retrieval call binding the contract method 0xe86dea4a.
//
// Solidity: function feeNumerator() view returns(uint256)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) FeeNumerator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "feeNumerator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeNumerator is a free data retrieval call binding the contract method 0xe86dea4a.
//
// Solidity: function feeNumerator() view returns(uint256)
func (_OwnerNativeHmy *OwnerNativeHmySession) FeeNumerator() (*big.Int, error) {
	return _OwnerNativeHmy.Contract.FeeNumerator(&_OwnerNativeHmy.CallOpts)
}

// FeeNumerator is a free data retrieval call binding the contract method 0xe86dea4a.
//
// Solidity: function feeNumerator() view returns(uint256)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) FeeNumerator() (*big.Int, error) {
	return _OwnerNativeHmy.Contract.FeeNumerator(&_OwnerNativeHmy.CallOpts)
}

// FixedFee is a free data retrieval call binding the contract method 0x91792d5b.
//
// Solidity: function fixedFee() view returns(uint256)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) FixedFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "fixedFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FixedFee is a free data retrieval call binding the contract method 0x91792d5b.
//
// Solidity: function fixedFee() view returns(uint256)
func (_OwnerNativeHmy *OwnerNativeHmySession) FixedFee() (*big.Int, error) {
	return _OwnerNativeHmy.Contract.FixedFee(&_OwnerNativeHmy.CallOpts)
}

// FixedFee is a free data retrieval call binding the contract method 0x91792d5b.
//
// Solidity: function fixedFee() view returns(uint256)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) FixedFee() (*big.Int, error) {
	return _OwnerNativeHmy.Contract.FixedFee(&_OwnerNativeHmy.CallOpts)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) GetAccumulatedFees(opts *bind.CallOpts) ([]TypesAsset, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "getAccumulatedFees")

	if err != nil {
		return *new([]TypesAsset), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesAsset)).(*[]TypesAsset)

	return out0, err

}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_OwnerNativeHmy *OwnerNativeHmySession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _OwnerNativeHmy.Contract.GetAccumulatedFees(&_OwnerNativeHmy.CallOpts)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _OwnerNativeHmy.Contract.GetAccumulatedFees(&_OwnerNativeHmy.CallOpts)
}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) GetBalanceOf(opts *bind.CallOpts, _owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "getBalanceOf", _owner, _coinName)

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
// Solidity: function getBalanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_OwnerNativeHmy *OwnerNativeHmySession) GetBalanceOf(_owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	return _OwnerNativeHmy.Contract.GetBalanceOf(&_OwnerNativeHmy.CallOpts, _owner, _coinName)
}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) GetBalanceOf(_owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	return _OwnerNativeHmy.Contract.GetBalanceOf(&_OwnerNativeHmy.CallOpts, _owner, _coinName)
}

// GetBalanceOfBatch is a free data retrieval call binding the contract method 0x4c3370f3.
//
// Solidity: function getBalanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) GetBalanceOfBatch(opts *bind.CallOpts, _owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
}, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "getBalanceOfBatch", _owner, _coinNames)

	outstruct := new(struct {
		UsableBalances     []*big.Int
		LockedBalances     []*big.Int
		RefundableBalances []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UsableBalances = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.LockedBalances = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.RefundableBalances = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetBalanceOfBatch is a free data retrieval call binding the contract method 0x4c3370f3.
//
// Solidity: function getBalanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances)
func (_OwnerNativeHmy *OwnerNativeHmySession) GetBalanceOfBatch(_owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
}, error) {
	return _OwnerNativeHmy.Contract.GetBalanceOfBatch(&_OwnerNativeHmy.CallOpts, _owner, _coinNames)
}

// GetBalanceOfBatch is a free data retrieval call binding the contract method 0x4c3370f3.
//
// Solidity: function getBalanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) GetBalanceOfBatch(_owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
}, error) {
	return _OwnerNativeHmy.Contract.GetBalanceOfBatch(&_OwnerNativeHmy.CallOpts, _owner, _coinNames)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_OwnerNativeHmy *OwnerNativeHmyCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_OwnerNativeHmy *OwnerNativeHmySession) GetOwners() ([]common.Address, error) {
	return _OwnerNativeHmy.Contract.GetOwners(&_OwnerNativeHmy.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) GetOwners() ([]common.Address, error) {
	return _OwnerNativeHmy.Contract.GetOwners(&_OwnerNativeHmy.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) IsOwner(opts *bind.CallOpts, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "isOwner", _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_OwnerNativeHmy *OwnerNativeHmySession) IsOwner(_owner common.Address) (bool, error) {
	return _OwnerNativeHmy.Contract.IsOwner(&_OwnerNativeHmy.CallOpts, _owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) IsOwner(_owner common.Address) (bool, error) {
	return _OwnerNativeHmy.Contract.IsOwner(&_OwnerNativeHmy.CallOpts, _owner)
}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_OwnerNativeHmy *OwnerNativeHmyCaller) IsValidCoin(opts *bind.CallOpts, _coinName string) (bool, error) {
	var out []interface{}
	err := _OwnerNativeHmy.contract.Call(opts, &out, "isValidCoin", _coinName)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_OwnerNativeHmy *OwnerNativeHmySession) IsValidCoin(_coinName string) (bool, error) {
	return _OwnerNativeHmy.Contract.IsValidCoin(&_OwnerNativeHmy.CallOpts, _coinName)
}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_OwnerNativeHmy *OwnerNativeHmyCallerSession) IsValidCoin(_coinName string) (bool, error) {
	return _OwnerNativeHmy.Contract.IsValidCoin(&_OwnerNativeHmy.CallOpts, _coinName)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) AddOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "addOwner", _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.AddOwner(&_OwnerNativeHmy.TransactOpts, _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.AddOwner(&_OwnerNativeHmy.TransactOpts, _owner)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) HandleResponseService(opts *bind.TransactOpts, _requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "handleResponseService", _requester, _coinName, _value, _fee, _rspCode)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) HandleResponseService(_requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.HandleResponseService(&_OwnerNativeHmy.TransactOpts, _requester, _coinName, _value, _fee, _rspCode)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) HandleResponseService(_requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.HandleResponseService(&_OwnerNativeHmy.TransactOpts, _requester, _coinName, _value, _fee, _rspCode)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) Initialize(opts *bind.TransactOpts, _nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "initialize", _nativeCoinName, _feeNumerator, _fixedFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) Initialize(_nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Initialize(&_OwnerNativeHmy.TransactOpts, _nativeCoinName, _feeNumerator, _fixedFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) Initialize(_nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Initialize(&_OwnerNativeHmy.TransactOpts, _nativeCoinName, _feeNumerator, _fixedFee)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) Mint(opts *bind.TransactOpts, _to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "mint", _to, _coinName, _value)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) Mint(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Mint(&_OwnerNativeHmy.TransactOpts, _to, _coinName, _value)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) Mint(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Mint(&_OwnerNativeHmy.TransactOpts, _to, _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) Reclaim(opts *bind.TransactOpts, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "reclaim", _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) Reclaim(_coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Reclaim(&_OwnerNativeHmy.TransactOpts, _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) Reclaim(_coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Reclaim(&_OwnerNativeHmy.TransactOpts, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) Refund(opts *bind.TransactOpts, _to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "refund", _to, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) Refund(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Refund(&_OwnerNativeHmy.TransactOpts, _to, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) Refund(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Refund(&_OwnerNativeHmy.TransactOpts, _to, _coinName, _value)
}

// Register is a paid mutator transaction binding the contract method 0xb3f90e0a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) Register(opts *bind.TransactOpts, _name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "register", _name, _symbol, _decimals)
}

// Register is a paid mutator transaction binding the contract method 0xb3f90e0a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) Register(_name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Register(&_OwnerNativeHmy.TransactOpts, _name, _symbol, _decimals)
}

// Register is a paid mutator transaction binding the contract method 0xb3f90e0a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) Register(_name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Register(&_OwnerNativeHmy.TransactOpts, _name, _symbol, _decimals)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) RemoveOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "removeOwner", _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.RemoveOwner(&_OwnerNativeHmy.TransactOpts, _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.RemoveOwner(&_OwnerNativeHmy.TransactOpts, _owner)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) SetFeeRatio(opts *bind.TransactOpts, _feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "setFeeRatio", _feeNumerator)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) SetFeeRatio(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.SetFeeRatio(&_OwnerNativeHmy.TransactOpts, _feeNumerator)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) SetFeeRatio(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.SetFeeRatio(&_OwnerNativeHmy.TransactOpts, _feeNumerator)
}

// SetFixedFee is a paid mutator transaction binding the contract method 0x37de8106.
//
// Solidity: function setFixedFee(uint256 _fixedFee) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) SetFixedFee(opts *bind.TransactOpts, _fixedFee *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "setFixedFee", _fixedFee)
}

// SetFixedFee is a paid mutator transaction binding the contract method 0x37de8106.
//
// Solidity: function setFixedFee(uint256 _fixedFee) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) SetFixedFee(_fixedFee *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.SetFixedFee(&_OwnerNativeHmy.TransactOpts, _fixedFee)
}

// SetFixedFee is a paid mutator transaction binding the contract method 0x37de8106.
//
// Solidity: function setFixedFee(uint256 _fixedFee) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) SetFixedFee(_fixedFee *big.Int) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.SetFixedFee(&_OwnerNativeHmy.TransactOpts, _fixedFee)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) Transfer(opts *bind.TransactOpts, _coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "transfer", _coinName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) Transfer(_coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Transfer(&_OwnerNativeHmy.TransactOpts, _coinName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) Transfer(_coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.Transfer(&_OwnerNativeHmy.TransactOpts, _coinName, _value, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) TransferBatch(opts *bind.TransactOpts, _coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "transferBatch", _coinNames, _values, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) TransferBatch(_coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.TransferBatch(&_OwnerNativeHmy.TransactOpts, _coinNames, _values, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) TransferBatch(_coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.TransferBatch(&_OwnerNativeHmy.TransactOpts, _coinNames, _values, _to)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) TransferFees(opts *bind.TransactOpts, _fa string) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "transferFees", _fa)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) TransferFees(_fa string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.TransferFees(&_OwnerNativeHmy.TransactOpts, _fa)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) TransferFees(_fa string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.TransferFees(&_OwnerNativeHmy.TransactOpts, _fa)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) TransferNativeCoin(opts *bind.TransactOpts, _to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "transferNativeCoin", _to)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) TransferNativeCoin(_to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.TransferNativeCoin(&_OwnerNativeHmy.TransactOpts, _to)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) TransferNativeCoin(_to string) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.TransferNativeCoin(&_OwnerNativeHmy.TransactOpts, _to)
}

// UpdateBSHPeriphery is a paid mutator transaction binding the contract method 0x2fbe21ba.
//
// Solidity: function updateBSHPeriphery(address _bshPeriphery) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactor) UpdateBSHPeriphery(opts *bind.TransactOpts, _bshPeriphery common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.contract.Transact(opts, "updateBSHPeriphery", _bshPeriphery)
}

// UpdateBSHPeriphery is a paid mutator transaction binding the contract method 0x2fbe21ba.
//
// Solidity: function updateBSHPeriphery(address _bshPeriphery) returns()
func (_OwnerNativeHmy *OwnerNativeHmySession) UpdateBSHPeriphery(_bshPeriphery common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.UpdateBSHPeriphery(&_OwnerNativeHmy.TransactOpts, _bshPeriphery)
}

// UpdateBSHPeriphery is a paid mutator transaction binding the contract method 0x2fbe21ba.
//
// Solidity: function updateBSHPeriphery(address _bshPeriphery) returns()
func (_OwnerNativeHmy *OwnerNativeHmyTransactorSession) UpdateBSHPeriphery(_bshPeriphery common.Address) (*types.Transaction, error) {
	return _OwnerNativeHmy.Contract.UpdateBSHPeriphery(&_OwnerNativeHmy.TransactOpts, _bshPeriphery)
}

// OwnerNativeHmyRemoveOwnershipIterator is returned from FilterRemoveOwnership and is used to iterate over the raw logs and unpacked data for RemoveOwnership events raised by the OwnerNativeHmy contract.
type OwnerNativeHmyRemoveOwnershipIterator struct {
	Event *OwnerNativeHmyRemoveOwnership // Event containing the contract specifics and raw log

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
func (it *OwnerNativeHmyRemoveOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerNativeHmyRemoveOwnership)
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
		it.Event = new(OwnerNativeHmyRemoveOwnership)
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
func (it *OwnerNativeHmyRemoveOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerNativeHmyRemoveOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerNativeHmyRemoveOwnership represents a RemoveOwnership event raised by the OwnerNativeHmy contract.
type OwnerNativeHmyRemoveOwnership struct {
	Remover     common.Address
	FormerOwner common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveOwnership is a free log retrieval operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_OwnerNativeHmy *OwnerNativeHmyFilterer) FilterRemoveOwnership(opts *bind.FilterOpts, remover []common.Address, formerOwner []common.Address) (*OwnerNativeHmyRemoveOwnershipIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _OwnerNativeHmy.contract.FilterLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnerNativeHmyRemoveOwnershipIterator{contract: _OwnerNativeHmy.contract, event: "RemoveOwnership", logs: logs, sub: sub}, nil
}

// WatchRemoveOwnership is a free log subscription operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_OwnerNativeHmy *OwnerNativeHmyFilterer) WatchRemoveOwnership(opts *bind.WatchOpts, sink chan<- *OwnerNativeHmyRemoveOwnership, remover []common.Address, formerOwner []common.Address) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _OwnerNativeHmy.contract.WatchLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerNativeHmyRemoveOwnership)
				if err := _OwnerNativeHmy.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
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
func (_OwnerNativeHmy *OwnerNativeHmyFilterer) ParseRemoveOwnership(log types.Log) (*OwnerNativeHmyRemoveOwnership, error) {
	event := new(OwnerNativeHmyRemoveOwnership)
	if err := _OwnerNativeHmy.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnerNativeHmySetOwnershipIterator is returned from FilterSetOwnership and is used to iterate over the raw logs and unpacked data for SetOwnership events raised by the OwnerNativeHmy contract.
type OwnerNativeHmySetOwnershipIterator struct {
	Event *OwnerNativeHmySetOwnership // Event containing the contract specifics and raw log

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
func (it *OwnerNativeHmySetOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnerNativeHmySetOwnership)
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
		it.Event = new(OwnerNativeHmySetOwnership)
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
func (it *OwnerNativeHmySetOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnerNativeHmySetOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnerNativeHmySetOwnership represents a SetOwnership event raised by the OwnerNativeHmy contract.
type OwnerNativeHmySetOwnership struct {
	Promoter common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetOwnership is a free log retrieval operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_OwnerNativeHmy *OwnerNativeHmyFilterer) FilterSetOwnership(opts *bind.FilterOpts, promoter []common.Address, newOwner []common.Address) (*OwnerNativeHmySetOwnershipIterator, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnerNativeHmy.contract.FilterLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnerNativeHmySetOwnershipIterator{contract: _OwnerNativeHmy.contract, event: "SetOwnership", logs: logs, sub: sub}, nil
}

// WatchSetOwnership is a free log subscription operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_OwnerNativeHmy *OwnerNativeHmyFilterer) WatchSetOwnership(opts *bind.WatchOpts, sink chan<- *OwnerNativeHmySetOwnership, promoter []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnerNativeHmy.contract.WatchLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnerNativeHmySetOwnership)
				if err := _OwnerNativeHmy.contract.UnpackLog(event, "SetOwnership", log); err != nil {
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
func (_OwnerNativeHmy *OwnerNativeHmyFilterer) ParseSetOwnership(log types.Log) (*OwnerNativeHmySetOwnership, error) {
	event := new(OwnerNativeHmySetOwnership)
	if err := _OwnerNativeHmy.contract.UnpackLog(event, "SetOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
