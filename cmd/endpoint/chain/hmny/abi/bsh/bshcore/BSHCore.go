// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bshcore

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

// BshcoreABI is the input ABI used to generate the binding from.
const BshcoreABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"formerOwner\",\"type\":\"address\"}],\"name\":\"RemoveOwnership\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"promoter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwnership\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"feeNumerator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fixedFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_nativeCoinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bshPeriphery\",\"type\":\"address\"}],\"name\":\"updateBSHPeriphery\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"}],\"name\":\"setFeeRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"name\":\"setFixedFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coinNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_names\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"coinId\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"isValidCoin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_valid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"getBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_usableBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lockedBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_refundableBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"}],\"name\":\"getBalanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_usableBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_lockedBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_refundableBalances\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccumulatedFees\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_accumulatedFees\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transferNativeCoin\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transferBatch\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rspCode\",\"type\":\"uint256\"}],\"name\":\"handleResponseService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_fa\",\"type\":\"string\"}],\"name\":\"transferFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Bshcore is an auto generated Go binding around an Ethereum contract.
type Bshcore struct {
	BshcoreCaller     // Read-only binding to the contract
	BshcoreTransactor // Write-only binding to the contract
	BshcoreFilterer   // Log filterer for contract events
}

// BshcoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type BshcoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshcoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BshcoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshcoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BshcoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BshcoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BshcoreSession struct {
	Contract     *Bshcore          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BshcoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BshcoreCallerSession struct {
	Contract *BshcoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BshcoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BshcoreTransactorSession struct {
	Contract     *BshcoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BshcoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type BshcoreRaw struct {
	Contract *Bshcore // Generic contract binding to access the raw methods on
}

// BshcoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BshcoreCallerRaw struct {
	Contract *BshcoreCaller // Generic read-only contract binding to access the raw methods on
}

// BshcoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BshcoreTransactorRaw struct {
	Contract *BshcoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBshcore creates a new instance of Bshcore, bound to a specific deployed contract.
func NewBshcore(address common.Address, backend bind.ContractBackend) (*Bshcore, error) {
	contract, err := bindBshcore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bshcore{BshcoreCaller: BshcoreCaller{contract: contract}, BshcoreTransactor: BshcoreTransactor{contract: contract}, BshcoreFilterer: BshcoreFilterer{contract: contract}}, nil
}

// NewBshcoreCaller creates a new read-only instance of Bshcore, bound to a specific deployed contract.
func NewBshcoreCaller(address common.Address, caller bind.ContractCaller) (*BshcoreCaller, error) {
	contract, err := bindBshcore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BshcoreCaller{contract: contract}, nil
}

// NewBshcoreTransactor creates a new write-only instance of Bshcore, bound to a specific deployed contract.
func NewBshcoreTransactor(address common.Address, transactor bind.ContractTransactor) (*BshcoreTransactor, error) {
	contract, err := bindBshcore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BshcoreTransactor{contract: contract}, nil
}

// NewBshcoreFilterer creates a new log filterer instance of Bshcore, bound to a specific deployed contract.
func NewBshcoreFilterer(address common.Address, filterer bind.ContractFilterer) (*BshcoreFilterer, error) {
	contract, err := bindBshcore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BshcoreFilterer{contract: contract}, nil
}

// bindBshcore binds a generic wrapper to an already deployed contract.
func bindBshcore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BshcoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bshcore *BshcoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bshcore.Contract.BshcoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bshcore *BshcoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bshcore.Contract.BshcoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bshcore *BshcoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bshcore.Contract.BshcoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bshcore *BshcoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bshcore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bshcore *BshcoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bshcore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bshcore *BshcoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bshcore.Contract.contract.Transact(opts, method, params...)
}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_Bshcore *BshcoreCaller) CoinId(opts *bind.CallOpts, _coinName string) (common.Address, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "coinId", _coinName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_Bshcore *BshcoreSession) CoinId(_coinName string) (common.Address, error) {
	return _Bshcore.Contract.CoinId(&_Bshcore.CallOpts, _coinName)
}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_Bshcore *BshcoreCallerSession) CoinId(_coinName string) (common.Address, error) {
	return _Bshcore.Contract.CoinId(&_Bshcore.CallOpts, _coinName)
}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_Bshcore *BshcoreCaller) CoinNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "coinNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_Bshcore *BshcoreSession) CoinNames() ([]string, error) {
	return _Bshcore.Contract.CoinNames(&_Bshcore.CallOpts)
}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_Bshcore *BshcoreCallerSession) CoinNames() ([]string, error) {
	return _Bshcore.Contract.CoinNames(&_Bshcore.CallOpts)
}

// FeeNumerator is a free data retrieval call binding the contract method 0xe86dea4a.
//
// Solidity: function feeNumerator() view returns(uint256)
func (_Bshcore *BshcoreCaller) FeeNumerator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "feeNumerator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeNumerator is a free data retrieval call binding the contract method 0xe86dea4a.
//
// Solidity: function feeNumerator() view returns(uint256)
func (_Bshcore *BshcoreSession) FeeNumerator() (*big.Int, error) {
	return _Bshcore.Contract.FeeNumerator(&_Bshcore.CallOpts)
}

// FeeNumerator is a free data retrieval call binding the contract method 0xe86dea4a.
//
// Solidity: function feeNumerator() view returns(uint256)
func (_Bshcore *BshcoreCallerSession) FeeNumerator() (*big.Int, error) {
	return _Bshcore.Contract.FeeNumerator(&_Bshcore.CallOpts)
}

// FixedFee is a free data retrieval call binding the contract method 0x91792d5b.
//
// Solidity: function fixedFee() view returns(uint256)
func (_Bshcore *BshcoreCaller) FixedFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "fixedFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FixedFee is a free data retrieval call binding the contract method 0x91792d5b.
//
// Solidity: function fixedFee() view returns(uint256)
func (_Bshcore *BshcoreSession) FixedFee() (*big.Int, error) {
	return _Bshcore.Contract.FixedFee(&_Bshcore.CallOpts)
}

// FixedFee is a free data retrieval call binding the contract method 0x91792d5b.
//
// Solidity: function fixedFee() view returns(uint256)
func (_Bshcore *BshcoreCallerSession) FixedFee() (*big.Int, error) {
	return _Bshcore.Contract.FixedFee(&_Bshcore.CallOpts)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_Bshcore *BshcoreCaller) GetAccumulatedFees(opts *bind.CallOpts) ([]TypesAsset, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "getAccumulatedFees")

	if err != nil {
		return *new([]TypesAsset), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesAsset)).(*[]TypesAsset)

	return out0, err

}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_Bshcore *BshcoreSession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _Bshcore.Contract.GetAccumulatedFees(&_Bshcore.CallOpts)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_Bshcore *BshcoreCallerSession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _Bshcore.Contract.GetAccumulatedFees(&_Bshcore.CallOpts)
}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_Bshcore *BshcoreCaller) GetBalanceOf(opts *bind.CallOpts, _owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "getBalanceOf", _owner, _coinName)

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
func (_Bshcore *BshcoreSession) GetBalanceOf(_owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	return _Bshcore.Contract.GetBalanceOf(&_Bshcore.CallOpts, _owner, _coinName)
}

// GetBalanceOf is a free data retrieval call binding the contract method 0xc5975f1d.
//
// Solidity: function getBalanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance)
func (_Bshcore *BshcoreCallerSession) GetBalanceOf(_owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
}, error) {
	return _Bshcore.Contract.GetBalanceOf(&_Bshcore.CallOpts, _owner, _coinName)
}

// GetBalanceOfBatch is a free data retrieval call binding the contract method 0x4c3370f3.
//
// Solidity: function getBalanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances)
func (_Bshcore *BshcoreCaller) GetBalanceOfBatch(opts *bind.CallOpts, _owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
}, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "getBalanceOfBatch", _owner, _coinNames)

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
func (_Bshcore *BshcoreSession) GetBalanceOfBatch(_owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
}, error) {
	return _Bshcore.Contract.GetBalanceOfBatch(&_Bshcore.CallOpts, _owner, _coinNames)
}

// GetBalanceOfBatch is a free data retrieval call binding the contract method 0x4c3370f3.
//
// Solidity: function getBalanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances)
func (_Bshcore *BshcoreCallerSession) GetBalanceOfBatch(_owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
}, error) {
	return _Bshcore.Contract.GetBalanceOfBatch(&_Bshcore.CallOpts, _owner, _coinNames)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Bshcore *BshcoreCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Bshcore *BshcoreSession) GetOwners() ([]common.Address, error) {
	return _Bshcore.Contract.GetOwners(&_Bshcore.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Bshcore *BshcoreCallerSession) GetOwners() ([]common.Address, error) {
	return _Bshcore.Contract.GetOwners(&_Bshcore.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Bshcore *BshcoreCaller) IsOwner(opts *bind.CallOpts, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "isOwner", _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Bshcore *BshcoreSession) IsOwner(_owner common.Address) (bool, error) {
	return _Bshcore.Contract.IsOwner(&_Bshcore.CallOpts, _owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Bshcore *BshcoreCallerSession) IsOwner(_owner common.Address) (bool, error) {
	return _Bshcore.Contract.IsOwner(&_Bshcore.CallOpts, _owner)
}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_Bshcore *BshcoreCaller) IsValidCoin(opts *bind.CallOpts, _coinName string) (bool, error) {
	var out []interface{}
	err := _Bshcore.contract.Call(opts, &out, "isValidCoin", _coinName)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_Bshcore *BshcoreSession) IsValidCoin(_coinName string) (bool, error) {
	return _Bshcore.Contract.IsValidCoin(&_Bshcore.CallOpts, _coinName)
}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_Bshcore *BshcoreCallerSession) IsValidCoin(_coinName string) (bool, error) {
	return _Bshcore.Contract.IsValidCoin(&_Bshcore.CallOpts, _coinName)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Bshcore *BshcoreTransactor) AddOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "addOwner", _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Bshcore *BshcoreSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bshcore.Contract.AddOwner(&_Bshcore.TransactOpts, _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Bshcore *BshcoreTransactorSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bshcore.Contract.AddOwner(&_Bshcore.TransactOpts, _owner)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_Bshcore *BshcoreTransactor) HandleResponseService(opts *bind.TransactOpts, _requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "handleResponseService", _requester, _coinName, _value, _fee, _rspCode)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_Bshcore *BshcoreSession) HandleResponseService(_requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.HandleResponseService(&_Bshcore.TransactOpts, _requester, _coinName, _value, _fee, _rspCode)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_Bshcore *BshcoreTransactorSession) HandleResponseService(_requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.HandleResponseService(&_Bshcore.TransactOpts, _requester, _coinName, _value, _fee, _rspCode)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Bshcore *BshcoreTransactor) Initialize(opts *bind.TransactOpts, _nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "initialize", _nativeCoinName, _feeNumerator, _fixedFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Bshcore *BshcoreSession) Initialize(_nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Initialize(&_Bshcore.TransactOpts, _nativeCoinName, _feeNumerator, _fixedFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Bshcore *BshcoreTransactorSession) Initialize(_nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Initialize(&_Bshcore.TransactOpts, _nativeCoinName, _feeNumerator, _fixedFee)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreTransactor) Mint(opts *bind.TransactOpts, _to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "mint", _to, _coinName, _value)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreSession) Mint(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Mint(&_Bshcore.TransactOpts, _to, _coinName, _value)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreTransactorSession) Mint(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Mint(&_Bshcore.TransactOpts, _to, _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreTransactor) Reclaim(opts *bind.TransactOpts, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "reclaim", _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreSession) Reclaim(_coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Reclaim(&_Bshcore.TransactOpts, _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreTransactorSession) Reclaim(_coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Reclaim(&_Bshcore.TransactOpts, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreTransactor) Refund(opts *bind.TransactOpts, _to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "refund", _to, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreSession) Refund(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Refund(&_Bshcore.TransactOpts, _to, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_Bshcore *BshcoreTransactorSession) Refund(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.Refund(&_Bshcore.TransactOpts, _to, _coinName, _value)
}

// Register is a paid mutator transaction binding the contract method 0xb3f90e0a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals) returns()
func (_Bshcore *BshcoreTransactor) Register(opts *bind.TransactOpts, _name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "register", _name, _symbol, _decimals)
}

// Register is a paid mutator transaction binding the contract method 0xb3f90e0a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals) returns()
func (_Bshcore *BshcoreSession) Register(_name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _Bshcore.Contract.Register(&_Bshcore.TransactOpts, _name, _symbol, _decimals)
}

// Register is a paid mutator transaction binding the contract method 0xb3f90e0a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals) returns()
func (_Bshcore *BshcoreTransactorSession) Register(_name string, _symbol string, _decimals uint8) (*types.Transaction, error) {
	return _Bshcore.Contract.Register(&_Bshcore.TransactOpts, _name, _symbol, _decimals)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Bshcore *BshcoreTransactor) RemoveOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "removeOwner", _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Bshcore *BshcoreSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bshcore.Contract.RemoveOwner(&_Bshcore.TransactOpts, _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Bshcore *BshcoreTransactorSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bshcore.Contract.RemoveOwner(&_Bshcore.TransactOpts, _owner)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_Bshcore *BshcoreTransactor) SetFeeRatio(opts *bind.TransactOpts, _feeNumerator *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "setFeeRatio", _feeNumerator)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_Bshcore *BshcoreSession) SetFeeRatio(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.SetFeeRatio(&_Bshcore.TransactOpts, _feeNumerator)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x19f4ff2f.
//
// Solidity: function setFeeRatio(uint256 _feeNumerator) returns()
func (_Bshcore *BshcoreTransactorSession) SetFeeRatio(_feeNumerator *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.SetFeeRatio(&_Bshcore.TransactOpts, _feeNumerator)
}

// SetFixedFee is a paid mutator transaction binding the contract method 0x37de8106.
//
// Solidity: function setFixedFee(uint256 _fixedFee) returns()
func (_Bshcore *BshcoreTransactor) SetFixedFee(opts *bind.TransactOpts, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "setFixedFee", _fixedFee)
}

// SetFixedFee is a paid mutator transaction binding the contract method 0x37de8106.
//
// Solidity: function setFixedFee(uint256 _fixedFee) returns()
func (_Bshcore *BshcoreSession) SetFixedFee(_fixedFee *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.SetFixedFee(&_Bshcore.TransactOpts, _fixedFee)
}

// SetFixedFee is a paid mutator transaction binding the contract method 0x37de8106.
//
// Solidity: function setFixedFee(uint256 _fixedFee) returns()
func (_Bshcore *BshcoreTransactorSession) SetFixedFee(_fixedFee *big.Int) (*types.Transaction, error) {
	return _Bshcore.Contract.SetFixedFee(&_Bshcore.TransactOpts, _fixedFee)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_Bshcore *BshcoreTransactor) Transfer(opts *bind.TransactOpts, _coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "transfer", _coinName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_Bshcore *BshcoreSession) Transfer(_coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _Bshcore.Contract.Transfer(&_Bshcore.TransactOpts, _coinName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_Bshcore *BshcoreTransactorSession) Transfer(_coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _Bshcore.Contract.Transfer(&_Bshcore.TransactOpts, _coinName, _value, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_Bshcore *BshcoreTransactor) TransferBatch(opts *bind.TransactOpts, _coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "transferBatch", _coinNames, _values, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_Bshcore *BshcoreSession) TransferBatch(_coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _Bshcore.Contract.TransferBatch(&_Bshcore.TransactOpts, _coinNames, _values, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_Bshcore *BshcoreTransactorSession) TransferBatch(_coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _Bshcore.Contract.TransferBatch(&_Bshcore.TransactOpts, _coinNames, _values, _to)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_Bshcore *BshcoreTransactor) TransferFees(opts *bind.TransactOpts, _fa string) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "transferFees", _fa)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_Bshcore *BshcoreSession) TransferFees(_fa string) (*types.Transaction, error) {
	return _Bshcore.Contract.TransferFees(&_Bshcore.TransactOpts, _fa)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_Bshcore *BshcoreTransactorSession) TransferFees(_fa string) (*types.Transaction, error) {
	return _Bshcore.Contract.TransferFees(&_Bshcore.TransactOpts, _fa)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_Bshcore *BshcoreTransactor) TransferNativeCoin(opts *bind.TransactOpts, _to string) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "transferNativeCoin", _to)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_Bshcore *BshcoreSession) TransferNativeCoin(_to string) (*types.Transaction, error) {
	return _Bshcore.Contract.TransferNativeCoin(&_Bshcore.TransactOpts, _to)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_Bshcore *BshcoreTransactorSession) TransferNativeCoin(_to string) (*types.Transaction, error) {
	return _Bshcore.Contract.TransferNativeCoin(&_Bshcore.TransactOpts, _to)
}

// UpdateBSHPeriphery is a paid mutator transaction binding the contract method 0x2fbe21ba.
//
// Solidity: function updateBSHPeriphery(address _bshPeriphery) returns()
func (_Bshcore *BshcoreTransactor) UpdateBSHPeriphery(opts *bind.TransactOpts, _bshPeriphery common.Address) (*types.Transaction, error) {
	return _Bshcore.contract.Transact(opts, "updateBSHPeriphery", _bshPeriphery)
}

// UpdateBSHPeriphery is a paid mutator transaction binding the contract method 0x2fbe21ba.
//
// Solidity: function updateBSHPeriphery(address _bshPeriphery) returns()
func (_Bshcore *BshcoreSession) UpdateBSHPeriphery(_bshPeriphery common.Address) (*types.Transaction, error) {
	return _Bshcore.Contract.UpdateBSHPeriphery(&_Bshcore.TransactOpts, _bshPeriphery)
}

// UpdateBSHPeriphery is a paid mutator transaction binding the contract method 0x2fbe21ba.
//
// Solidity: function updateBSHPeriphery(address _bshPeriphery) returns()
func (_Bshcore *BshcoreTransactorSession) UpdateBSHPeriphery(_bshPeriphery common.Address) (*types.Transaction, error) {
	return _Bshcore.Contract.UpdateBSHPeriphery(&_Bshcore.TransactOpts, _bshPeriphery)
}

// BshcoreRemoveOwnershipIterator is returned from FilterRemoveOwnership and is used to iterate over the raw logs and unpacked data for RemoveOwnership events raised by the Bshcore contract.
type BshcoreRemoveOwnershipIterator struct {
	Event *BshcoreRemoveOwnership // Event containing the contract specifics and raw log

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
func (it *BshcoreRemoveOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshcoreRemoveOwnership)
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
		it.Event = new(BshcoreRemoveOwnership)
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
func (it *BshcoreRemoveOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshcoreRemoveOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshcoreRemoveOwnership represents a RemoveOwnership event raised by the Bshcore contract.
type BshcoreRemoveOwnership struct {
	Remover     common.Address
	FormerOwner common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveOwnership is a free log retrieval operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_Bshcore *BshcoreFilterer) FilterRemoveOwnership(opts *bind.FilterOpts, remover []common.Address, formerOwner []common.Address) (*BshcoreRemoveOwnershipIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _Bshcore.contract.FilterLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BshcoreRemoveOwnershipIterator{contract: _Bshcore.contract, event: "RemoveOwnership", logs: logs, sub: sub}, nil
}

// WatchRemoveOwnership is a free log subscription operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_Bshcore *BshcoreFilterer) WatchRemoveOwnership(opts *bind.WatchOpts, sink chan<- *BshcoreRemoveOwnership, remover []common.Address, formerOwner []common.Address) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _Bshcore.contract.WatchLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshcoreRemoveOwnership)
				if err := _Bshcore.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
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
func (_Bshcore *BshcoreFilterer) ParseRemoveOwnership(log types.Log) (*BshcoreRemoveOwnership, error) {
	event := new(BshcoreRemoveOwnership)
	if err := _Bshcore.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BshcoreSetOwnershipIterator is returned from FilterSetOwnership and is used to iterate over the raw logs and unpacked data for SetOwnership events raised by the Bshcore contract.
type BshcoreSetOwnershipIterator struct {
	Event *BshcoreSetOwnership // Event containing the contract specifics and raw log

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
func (it *BshcoreSetOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BshcoreSetOwnership)
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
		it.Event = new(BshcoreSetOwnership)
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
func (it *BshcoreSetOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BshcoreSetOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BshcoreSetOwnership represents a SetOwnership event raised by the Bshcore contract.
type BshcoreSetOwnership struct {
	Promoter common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetOwnership is a free log retrieval operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_Bshcore *BshcoreFilterer) FilterSetOwnership(opts *bind.FilterOpts, promoter []common.Address, newOwner []common.Address) (*BshcoreSetOwnershipIterator, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bshcore.contract.FilterLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BshcoreSetOwnershipIterator{contract: _Bshcore.contract, event: "SetOwnership", logs: logs, sub: sub}, nil
}

// WatchSetOwnership is a free log subscription operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_Bshcore *BshcoreFilterer) WatchSetOwnership(opts *bind.WatchOpts, sink chan<- *BshcoreSetOwnership, promoter []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bshcore.contract.WatchLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BshcoreSetOwnership)
				if err := _Bshcore.contract.UnpackLog(event, "SetOwnership", log); err != nil {
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
func (_Bshcore *BshcoreFilterer) ParseSetOwnership(log types.Log) (*BshcoreSetOwnership, error) {
	event := new(BshcoreSetOwnership)
	if err := _Bshcore.contract.UnpackLog(event, "SetOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
