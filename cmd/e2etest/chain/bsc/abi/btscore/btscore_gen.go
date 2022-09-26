// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package btscore

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

// BtscoreABI is the input ABI used to generate the binding from.
const BtscoreABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remover\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"formerOwner\",\"type\":\"address\"}],\"name\":\"RemoveOwnership\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"promoter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwnership\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_nativeCoinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNativeCoinName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_btsPeriphery\",\"type\":\"address\"}],\"name\":\"updateBTSPeriphery\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"name\":\"setFeeRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"coinNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_names\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"coinId\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"isValidCoin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_valid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"feeRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeNumerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fixedFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_usableBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lockedBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_refundableBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_userBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_usableBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_lockedBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_refundableBalances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_userBalances\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAccumulatedFees\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.Asset[]\",\"name\":\"_accumulatedFees\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transferNativeCoin\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_coinNames\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_values\",\"type\":\"uint256[]\"},{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"transferBatch\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_coinName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rspCode\",\"type\":\"uint256\"}],\"name\":\"handleResponseService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_fa\",\"type\":\"string\"}],\"name\":\"transferFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Btscore is an auto generated Go binding around an Ethereum contract.
type Btscore struct {
	BtscoreCaller     // Read-only binding to the contract
	BtscoreTransactor // Write-only binding to the contract
	BtscoreFilterer   // Log filterer for contract events
}

// BtscoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type BtscoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtscoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BtscoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtscoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BtscoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BtscoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BtscoreSession struct {
	Contract     *Btscore          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BtscoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BtscoreCallerSession struct {
	Contract *BtscoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BtscoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BtscoreTransactorSession struct {
	Contract     *BtscoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BtscoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type BtscoreRaw struct {
	Contract *Btscore // Generic contract binding to access the raw methods on
}

// BtscoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BtscoreCallerRaw struct {
	Contract *BtscoreCaller // Generic read-only contract binding to access the raw methods on
}

// BtscoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BtscoreTransactorRaw struct {
	Contract *BtscoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBtscore creates a new instance of Btscore, bound to a specific deployed contract.
func NewBtscore(address common.Address, backend bind.ContractBackend) (*Btscore, error) {
	contract, err := bindBtscore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Btscore{BtscoreCaller: BtscoreCaller{contract: contract}, BtscoreTransactor: BtscoreTransactor{contract: contract}, BtscoreFilterer: BtscoreFilterer{contract: contract}}, nil
}

// NewBtscoreCaller creates a new read-only instance of Btscore, bound to a specific deployed contract.
func NewBtscoreCaller(address common.Address, caller bind.ContractCaller) (*BtscoreCaller, error) {
	contract, err := bindBtscore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BtscoreCaller{contract: contract}, nil
}

// NewBtscoreTransactor creates a new write-only instance of Btscore, bound to a specific deployed contract.
func NewBtscoreTransactor(address common.Address, transactor bind.ContractTransactor) (*BtscoreTransactor, error) {
	contract, err := bindBtscore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BtscoreTransactor{contract: contract}, nil
}

// NewBtscoreFilterer creates a new log filterer instance of Btscore, bound to a specific deployed contract.
func NewBtscoreFilterer(address common.Address, filterer bind.ContractFilterer) (*BtscoreFilterer, error) {
	contract, err := bindBtscore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BtscoreFilterer{contract: contract}, nil
}

// bindBtscore binds a generic wrapper to an already deployed contract.
func bindBtscore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BtscoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Btscore *BtscoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Btscore.Contract.BtscoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Btscore *BtscoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Btscore.Contract.BtscoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Btscore *BtscoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Btscore.Contract.BtscoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Btscore *BtscoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Btscore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Btscore *BtscoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Btscore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Btscore *BtscoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Btscore.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0xb9b092c8.
//
// Solidity: function balanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance, uint256 _userBalance)
func (_Btscore *BtscoreCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
	UserBalance       *big.Int
}, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "balanceOf", _owner, _coinName)

	outstruct := new(struct {
		UsableBalance     *big.Int
		LockedBalance     *big.Int
		RefundableBalance *big.Int
		UserBalance       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UsableBalance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LockedBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RefundableBalance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UserBalance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// BalanceOf is a free data retrieval call binding the contract method 0xb9b092c8.
//
// Solidity: function balanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance, uint256 _userBalance)
func (_Btscore *BtscoreSession) BalanceOf(_owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
	UserBalance       *big.Int
}, error) {
	return _Btscore.Contract.BalanceOf(&_Btscore.CallOpts, _owner, _coinName)
}

// BalanceOf is a free data retrieval call binding the contract method 0xb9b092c8.
//
// Solidity: function balanceOf(address _owner, string _coinName) view returns(uint256 _usableBalance, uint256 _lockedBalance, uint256 _refundableBalance, uint256 _userBalance)
func (_Btscore *BtscoreCallerSession) BalanceOf(_owner common.Address, _coinName string) (struct {
	UsableBalance     *big.Int
	LockedBalance     *big.Int
	RefundableBalance *big.Int
	UserBalance       *big.Int
}, error) {
	return _Btscore.Contract.BalanceOf(&_Btscore.CallOpts, _owner, _coinName)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x17d55ad6.
//
// Solidity: function balanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances, uint256[] _userBalances)
func (_Btscore *BtscoreCaller) BalanceOfBatch(opts *bind.CallOpts, _owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
	UserBalances       []*big.Int
}, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "balanceOfBatch", _owner, _coinNames)

	outstruct := new(struct {
		UsableBalances     []*big.Int
		LockedBalances     []*big.Int
		RefundableBalances []*big.Int
		UserBalances       []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UsableBalances = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.LockedBalances = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.RefundableBalances = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)
	outstruct.UserBalances = *abi.ConvertType(out[3], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x17d55ad6.
//
// Solidity: function balanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances, uint256[] _userBalances)
func (_Btscore *BtscoreSession) BalanceOfBatch(_owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
	UserBalances       []*big.Int
}, error) {
	return _Btscore.Contract.BalanceOfBatch(&_Btscore.CallOpts, _owner, _coinNames)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x17d55ad6.
//
// Solidity: function balanceOfBatch(address _owner, string[] _coinNames) view returns(uint256[] _usableBalances, uint256[] _lockedBalances, uint256[] _refundableBalances, uint256[] _userBalances)
func (_Btscore *BtscoreCallerSession) BalanceOfBatch(_owner common.Address, _coinNames []string) (struct {
	UsableBalances     []*big.Int
	LockedBalances     []*big.Int
	RefundableBalances []*big.Int
	UserBalances       []*big.Int
}, error) {
	return _Btscore.Contract.BalanceOfBatch(&_Btscore.CallOpts, _owner, _coinNames)
}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_Btscore *BtscoreCaller) CoinId(opts *bind.CallOpts, _coinName string) (common.Address, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "coinId", _coinName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_Btscore *BtscoreSession) CoinId(_coinName string) (common.Address, error) {
	return _Btscore.Contract.CoinId(&_Btscore.CallOpts, _coinName)
}

// CoinId is a free data retrieval call binding the contract method 0x8506a74e.
//
// Solidity: function coinId(string _coinName) view returns(address)
func (_Btscore *BtscoreCallerSession) CoinId(_coinName string) (common.Address, error) {
	return _Btscore.Contract.CoinId(&_Btscore.CallOpts, _coinName)
}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_Btscore *BtscoreCaller) CoinNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "coinNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_Btscore *BtscoreSession) CoinNames() ([]string, error) {
	return _Btscore.Contract.CoinNames(&_Btscore.CallOpts)
}

// CoinNames is a free data retrieval call binding the contract method 0x9bda00cd.
//
// Solidity: function coinNames() view returns(string[] _names)
func (_Btscore *BtscoreCallerSession) CoinNames() ([]string, error) {
	return _Btscore.Contract.CoinNames(&_Btscore.CallOpts)
}

// FeeRatio is a free data retrieval call binding the contract method 0xc40238c4.
//
// Solidity: function feeRatio(string _coinName) view returns(uint256 _feeNumerator, uint256 _fixedFee)
func (_Btscore *BtscoreCaller) FeeRatio(opts *bind.CallOpts, _coinName string) (struct {
	FeeNumerator *big.Int
	FixedFee     *big.Int
}, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "feeRatio", _coinName)

	outstruct := new(struct {
		FeeNumerator *big.Int
		FixedFee     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FeeNumerator = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FixedFee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// FeeRatio is a free data retrieval call binding the contract method 0xc40238c4.
//
// Solidity: function feeRatio(string _coinName) view returns(uint256 _feeNumerator, uint256 _fixedFee)
func (_Btscore *BtscoreSession) FeeRatio(_coinName string) (struct {
	FeeNumerator *big.Int
	FixedFee     *big.Int
}, error) {
	return _Btscore.Contract.FeeRatio(&_Btscore.CallOpts, _coinName)
}

// FeeRatio is a free data retrieval call binding the contract method 0xc40238c4.
//
// Solidity: function feeRatio(string _coinName) view returns(uint256 _feeNumerator, uint256 _fixedFee)
func (_Btscore *BtscoreCallerSession) FeeRatio(_coinName string) (struct {
	FeeNumerator *big.Int
	FixedFee     *big.Int
}, error) {
	return _Btscore.Contract.FeeRatio(&_Btscore.CallOpts, _coinName)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_Btscore *BtscoreCaller) GetAccumulatedFees(opts *bind.CallOpts) ([]TypesAsset, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "getAccumulatedFees")

	if err != nil {
		return *new([]TypesAsset), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesAsset)).(*[]TypesAsset)

	return out0, err

}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_Btscore *BtscoreSession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _Btscore.Contract.GetAccumulatedFees(&_Btscore.CallOpts)
}

// GetAccumulatedFees is a free data retrieval call binding the contract method 0x5df45a37.
//
// Solidity: function getAccumulatedFees() view returns((string,uint256)[] _accumulatedFees)
func (_Btscore *BtscoreCallerSession) GetAccumulatedFees() ([]TypesAsset, error) {
	return _Btscore.Contract.GetAccumulatedFees(&_Btscore.CallOpts)
}

// GetNativeCoinName is a free data retrieval call binding the contract method 0x71433cfb.
//
// Solidity: function getNativeCoinName() view returns(string)
func (_Btscore *BtscoreCaller) GetNativeCoinName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "getNativeCoinName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNativeCoinName is a free data retrieval call binding the contract method 0x71433cfb.
//
// Solidity: function getNativeCoinName() view returns(string)
func (_Btscore *BtscoreSession) GetNativeCoinName() (string, error) {
	return _Btscore.Contract.GetNativeCoinName(&_Btscore.CallOpts)
}

// GetNativeCoinName is a free data retrieval call binding the contract method 0x71433cfb.
//
// Solidity: function getNativeCoinName() view returns(string)
func (_Btscore *BtscoreCallerSession) GetNativeCoinName() (string, error) {
	return _Btscore.Contract.GetNativeCoinName(&_Btscore.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Btscore *BtscoreCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Btscore *BtscoreSession) GetOwners() ([]common.Address, error) {
	return _Btscore.Contract.GetOwners(&_Btscore.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Btscore *BtscoreCallerSession) GetOwners() ([]common.Address, error) {
	return _Btscore.Contract.GetOwners(&_Btscore.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Btscore *BtscoreCaller) IsOwner(opts *bind.CallOpts, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "isOwner", _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Btscore *BtscoreSession) IsOwner(_owner common.Address) (bool, error) {
	return _Btscore.Contract.IsOwner(&_Btscore.CallOpts, _owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Btscore *BtscoreCallerSession) IsOwner(_owner common.Address) (bool, error) {
	return _Btscore.Contract.IsOwner(&_Btscore.CallOpts, _owner)
}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_Btscore *BtscoreCaller) IsValidCoin(opts *bind.CallOpts, _coinName string) (bool, error) {
	var out []interface{}
	err := _Btscore.contract.Call(opts, &out, "isValidCoin", _coinName)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_Btscore *BtscoreSession) IsValidCoin(_coinName string) (bool, error) {
	return _Btscore.Contract.IsValidCoin(&_Btscore.CallOpts, _coinName)
}

// IsValidCoin is a free data retrieval call binding the contract method 0xb30a072b.
//
// Solidity: function isValidCoin(string _coinName) view returns(bool _valid)
func (_Btscore *BtscoreCallerSession) IsValidCoin(_coinName string) (bool, error) {
	return _Btscore.Contract.IsValidCoin(&_Btscore.CallOpts, _coinName)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Btscore *BtscoreTransactor) AddOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "addOwner", _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Btscore *BtscoreSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.AddOwner(&_Btscore.TransactOpts, _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Btscore *BtscoreTransactorSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.AddOwner(&_Btscore.TransactOpts, _owner)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_Btscore *BtscoreTransactor) HandleResponseService(opts *bind.TransactOpts, _requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "handleResponseService", _requester, _coinName, _value, _fee, _rspCode)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_Btscore *BtscoreSession) HandleResponseService(_requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.HandleResponseService(&_Btscore.TransactOpts, _requester, _coinName, _value, _fee, _rspCode)
}

// HandleResponseService is a paid mutator transaction binding the contract method 0x69c939b6.
//
// Solidity: function handleResponseService(address _requester, string _coinName, uint256 _value, uint256 _fee, uint256 _rspCode) returns()
func (_Btscore *BtscoreTransactorSession) HandleResponseService(_requester common.Address, _coinName string, _value *big.Int, _fee *big.Int, _rspCode *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.HandleResponseService(&_Btscore.TransactOpts, _requester, _coinName, _value, _fee, _rspCode)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Btscore *BtscoreTransactor) Initialize(opts *bind.TransactOpts, _nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "initialize", _nativeCoinName, _feeNumerator, _fixedFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Btscore *BtscoreSession) Initialize(_nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Initialize(&_Btscore.TransactOpts, _nativeCoinName, _feeNumerator, _fixedFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x977d08c0.
//
// Solidity: function initialize(string _nativeCoinName, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Btscore *BtscoreTransactorSession) Initialize(_nativeCoinName string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Initialize(&_Btscore.TransactOpts, _nativeCoinName, _feeNumerator, _fixedFee)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreTransactor) Mint(opts *bind.TransactOpts, _to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "mint", _to, _coinName, _value)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreSession) Mint(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Mint(&_Btscore.TransactOpts, _to, _coinName, _value)
}

// Mint is a paid mutator transaction binding the contract method 0xba7aef43.
//
// Solidity: function mint(address _to, string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreTransactorSession) Mint(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Mint(&_Btscore.TransactOpts, _to, _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreTransactor) Reclaim(opts *bind.TransactOpts, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "reclaim", _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreSession) Reclaim(_coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Reclaim(&_Btscore.TransactOpts, _coinName, _value)
}

// Reclaim is a paid mutator transaction binding the contract method 0x4f6e30b8.
//
// Solidity: function reclaim(string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreTransactorSession) Reclaim(_coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Reclaim(&_Btscore.TransactOpts, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreTransactor) Refund(opts *bind.TransactOpts, _to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "refund", _to, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreSession) Refund(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Refund(&_Btscore.TransactOpts, _to, _coinName, _value)
}

// Refund is a paid mutator transaction binding the contract method 0x48c692d7.
//
// Solidity: function refund(address _to, string _coinName, uint256 _value) returns()
func (_Btscore *BtscoreTransactorSession) Refund(_to common.Address, _coinName string, _value *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.Refund(&_Btscore.TransactOpts, _to, _coinName, _value)
}

// Register is a paid mutator transaction binding the contract method 0xd1155d6a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals, uint256 _feeNumerator, uint256 _fixedFee, address _addr) returns()
func (_Btscore *BtscoreTransactor) Register(opts *bind.TransactOpts, _name string, _symbol string, _decimals uint8, _feeNumerator *big.Int, _fixedFee *big.Int, _addr common.Address) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "register", _name, _symbol, _decimals, _feeNumerator, _fixedFee, _addr)
}

// Register is a paid mutator transaction binding the contract method 0xd1155d6a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals, uint256 _feeNumerator, uint256 _fixedFee, address _addr) returns()
func (_Btscore *BtscoreSession) Register(_name string, _symbol string, _decimals uint8, _feeNumerator *big.Int, _fixedFee *big.Int, _addr common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.Register(&_Btscore.TransactOpts, _name, _symbol, _decimals, _feeNumerator, _fixedFee, _addr)
}

// Register is a paid mutator transaction binding the contract method 0xd1155d6a.
//
// Solidity: function register(string _name, string _symbol, uint8 _decimals, uint256 _feeNumerator, uint256 _fixedFee, address _addr) returns()
func (_Btscore *BtscoreTransactorSession) Register(_name string, _symbol string, _decimals uint8, _feeNumerator *big.Int, _fixedFee *big.Int, _addr common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.Register(&_Btscore.TransactOpts, _name, _symbol, _decimals, _feeNumerator, _fixedFee, _addr)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Btscore *BtscoreTransactor) RemoveOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "removeOwner", _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Btscore *BtscoreSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.RemoveOwner(&_Btscore.TransactOpts, _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Btscore *BtscoreTransactorSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.RemoveOwner(&_Btscore.TransactOpts, _owner)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x030de8f1.
//
// Solidity: function setFeeRatio(string _name, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Btscore *BtscoreTransactor) SetFeeRatio(opts *bind.TransactOpts, _name string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "setFeeRatio", _name, _feeNumerator, _fixedFee)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x030de8f1.
//
// Solidity: function setFeeRatio(string _name, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Btscore *BtscoreSession) SetFeeRatio(_name string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.SetFeeRatio(&_Btscore.TransactOpts, _name, _feeNumerator, _fixedFee)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x030de8f1.
//
// Solidity: function setFeeRatio(string _name, uint256 _feeNumerator, uint256 _fixedFee) returns()
func (_Btscore *BtscoreTransactorSession) SetFeeRatio(_name string, _feeNumerator *big.Int, _fixedFee *big.Int) (*types.Transaction, error) {
	return _Btscore.Contract.SetFeeRatio(&_Btscore.TransactOpts, _name, _feeNumerator, _fixedFee)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_Btscore *BtscoreTransactor) Transfer(opts *bind.TransactOpts, _coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "transfer", _coinName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_Btscore *BtscoreSession) Transfer(_coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _Btscore.Contract.Transfer(&_Btscore.TransactOpts, _coinName, _value, _to)
}

// Transfer is a paid mutator transaction binding the contract method 0xd5823df0.
//
// Solidity: function transfer(string _coinName, uint256 _value, string _to) returns()
func (_Btscore *BtscoreTransactorSession) Transfer(_coinName string, _value *big.Int, _to string) (*types.Transaction, error) {
	return _Btscore.Contract.Transfer(&_Btscore.TransactOpts, _coinName, _value, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_Btscore *BtscoreTransactor) TransferBatch(opts *bind.TransactOpts, _coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "transferBatch", _coinNames, _values, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_Btscore *BtscoreSession) TransferBatch(_coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _Btscore.Contract.TransferBatch(&_Btscore.TransactOpts, _coinNames, _values, _to)
}

// TransferBatch is a paid mutator transaction binding the contract method 0x48c6c8e6.
//
// Solidity: function transferBatch(string[] _coinNames, uint256[] _values, string _to) payable returns()
func (_Btscore *BtscoreTransactorSession) TransferBatch(_coinNames []string, _values []*big.Int, _to string) (*types.Transaction, error) {
	return _Btscore.Contract.TransferBatch(&_Btscore.TransactOpts, _coinNames, _values, _to)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_Btscore *BtscoreTransactor) TransferFees(opts *bind.TransactOpts, _fa string) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "transferFees", _fa)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_Btscore *BtscoreSession) TransferFees(_fa string) (*types.Transaction, error) {
	return _Btscore.Contract.TransferFees(&_Btscore.TransactOpts, _fa)
}

// TransferFees is a paid mutator transaction binding the contract method 0x173e4045.
//
// Solidity: function transferFees(string _fa) returns()
func (_Btscore *BtscoreTransactorSession) TransferFees(_fa string) (*types.Transaction, error) {
	return _Btscore.Contract.TransferFees(&_Btscore.TransactOpts, _fa)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_Btscore *BtscoreTransactor) TransferNativeCoin(opts *bind.TransactOpts, _to string) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "transferNativeCoin", _to)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_Btscore *BtscoreSession) TransferNativeCoin(_to string) (*types.Transaction, error) {
	return _Btscore.Contract.TransferNativeCoin(&_Btscore.TransactOpts, _to)
}

// TransferNativeCoin is a paid mutator transaction binding the contract method 0x74e518c5.
//
// Solidity: function transferNativeCoin(string _to) payable returns()
func (_Btscore *BtscoreTransactorSession) TransferNativeCoin(_to string) (*types.Transaction, error) {
	return _Btscore.Contract.TransferNativeCoin(&_Btscore.TransactOpts, _to)
}

// UpdateBTSPeriphery is a paid mutator transaction binding the contract method 0xc9332478.
//
// Solidity: function updateBTSPeriphery(address _btsPeriphery) returns()
func (_Btscore *BtscoreTransactor) UpdateBTSPeriphery(opts *bind.TransactOpts, _btsPeriphery common.Address) (*types.Transaction, error) {
	return _Btscore.contract.Transact(opts, "updateBTSPeriphery", _btsPeriphery)
}

// UpdateBTSPeriphery is a paid mutator transaction binding the contract method 0xc9332478.
//
// Solidity: function updateBTSPeriphery(address _btsPeriphery) returns()
func (_Btscore *BtscoreSession) UpdateBTSPeriphery(_btsPeriphery common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.UpdateBTSPeriphery(&_Btscore.TransactOpts, _btsPeriphery)
}

// UpdateBTSPeriphery is a paid mutator transaction binding the contract method 0xc9332478.
//
// Solidity: function updateBTSPeriphery(address _btsPeriphery) returns()
func (_Btscore *BtscoreTransactorSession) UpdateBTSPeriphery(_btsPeriphery common.Address) (*types.Transaction, error) {
	return _Btscore.Contract.UpdateBTSPeriphery(&_Btscore.TransactOpts, _btsPeriphery)
}

// BtscoreInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Btscore contract.
type BtscoreInitializedIterator struct {
	Event *BtscoreInitialized // Event containing the contract specifics and raw log

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
func (it *BtscoreInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtscoreInitialized)
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
		it.Event = new(BtscoreInitialized)
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
func (it *BtscoreInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtscoreInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtscoreInitialized represents a Initialized event raised by the Btscore contract.
type BtscoreInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Btscore *BtscoreFilterer) FilterInitialized(opts *bind.FilterOpts) (*BtscoreInitializedIterator, error) {

	logs, sub, err := _Btscore.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BtscoreInitializedIterator{contract: _Btscore.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Btscore *BtscoreFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BtscoreInitialized) (event.Subscription, error) {

	logs, sub, err := _Btscore.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtscoreInitialized)
				if err := _Btscore.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Btscore *BtscoreFilterer) ParseInitialized(log types.Log) (*BtscoreInitialized, error) {
	event := new(BtscoreInitialized)
	if err := _Btscore.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtscoreRemoveOwnershipIterator is returned from FilterRemoveOwnership and is used to iterate over the raw logs and unpacked data for RemoveOwnership events raised by the Btscore contract.
type BtscoreRemoveOwnershipIterator struct {
	Event *BtscoreRemoveOwnership // Event containing the contract specifics and raw log

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
func (it *BtscoreRemoveOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtscoreRemoveOwnership)
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
		it.Event = new(BtscoreRemoveOwnership)
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
func (it *BtscoreRemoveOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtscoreRemoveOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtscoreRemoveOwnership represents a RemoveOwnership event raised by the Btscore contract.
type BtscoreRemoveOwnership struct {
	Remover     common.Address
	FormerOwner common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveOwnership is a free log retrieval operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_Btscore *BtscoreFilterer) FilterRemoveOwnership(opts *bind.FilterOpts, remover []common.Address, formerOwner []common.Address) (*BtscoreRemoveOwnershipIterator, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _Btscore.contract.FilterLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BtscoreRemoveOwnershipIterator{contract: _Btscore.contract, event: "RemoveOwnership", logs: logs, sub: sub}, nil
}

// WatchRemoveOwnership is a free log subscription operation binding the contract event 0xda94804c6fea691edd453996746b93f789375a915c17acf1d1460944dffb9b37.
//
// Solidity: event RemoveOwnership(address indexed remover, address indexed formerOwner)
func (_Btscore *BtscoreFilterer) WatchRemoveOwnership(opts *bind.WatchOpts, sink chan<- *BtscoreRemoveOwnership, remover []common.Address, formerOwner []common.Address) (event.Subscription, error) {

	var removerRule []interface{}
	for _, removerItem := range remover {
		removerRule = append(removerRule, removerItem)
	}
	var formerOwnerRule []interface{}
	for _, formerOwnerItem := range formerOwner {
		formerOwnerRule = append(formerOwnerRule, formerOwnerItem)
	}

	logs, sub, err := _Btscore.contract.WatchLogs(opts, "RemoveOwnership", removerRule, formerOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtscoreRemoveOwnership)
				if err := _Btscore.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
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
func (_Btscore *BtscoreFilterer) ParseRemoveOwnership(log types.Log) (*BtscoreRemoveOwnership, error) {
	event := new(BtscoreRemoveOwnership)
	if err := _Btscore.contract.UnpackLog(event, "RemoveOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BtscoreSetOwnershipIterator is returned from FilterSetOwnership and is used to iterate over the raw logs and unpacked data for SetOwnership events raised by the Btscore contract.
type BtscoreSetOwnershipIterator struct {
	Event *BtscoreSetOwnership // Event containing the contract specifics and raw log

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
func (it *BtscoreSetOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BtscoreSetOwnership)
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
		it.Event = new(BtscoreSetOwnership)
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
func (it *BtscoreSetOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BtscoreSetOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BtscoreSetOwnership represents a SetOwnership event raised by the Btscore contract.
type BtscoreSetOwnership struct {
	Promoter common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetOwnership is a free log retrieval operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_Btscore *BtscoreFilterer) FilterSetOwnership(opts *bind.FilterOpts, promoter []common.Address, newOwner []common.Address) (*BtscoreSetOwnershipIterator, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Btscore.contract.FilterLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BtscoreSetOwnershipIterator{contract: _Btscore.contract, event: "SetOwnership", logs: logs, sub: sub}, nil
}

// WatchSetOwnership is a free log subscription operation binding the contract event 0x8a566e8b76ab6f8a031711472b4fdc77432d6f59c804e4e0811a1c3bbfa74771.
//
// Solidity: event SetOwnership(address indexed promoter, address indexed newOwner)
func (_Btscore *BtscoreFilterer) WatchSetOwnership(opts *bind.WatchOpts, sink chan<- *BtscoreSetOwnership, promoter []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var promoterRule []interface{}
	for _, promoterItem := range promoter {
		promoterRule = append(promoterRule, promoterItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Btscore.contract.WatchLogs(opts, "SetOwnership", promoterRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BtscoreSetOwnership)
				if err := _Btscore.contract.UnpackLog(event, "SetOwnership", log); err != nil {
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
func (_Btscore *BtscoreFilterer) ParseSetOwnership(log types.Log) (*BtscoreSetOwnership, error) {
	event := new(BtscoreSetOwnership)
	if err := _Btscore.contract.UnpackLog(event, "SetOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
