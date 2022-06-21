// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BEP

/*
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

// TokenBSHABI is the input ABI used to generate the binding from.
const TokenBSHABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenBSH is an auto generated Go binding around an Ethereum contract.
type TokenBSH struct {
	TokenBSHCaller     // Read-only binding to the contract
	TokenBSHTransactor // Write-only binding to the contract
	TokenBSHFilterer   // Log filterer for contract events
}

// TokenBSHCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenBSHCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenBSHTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenBSHTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenBSHFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenBSHFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenBSHSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenBSHSession struct {
	Contract     *TokenBSH         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenBSHCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenBSHCallerSession struct {
	Contract *TokenBSHCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TokenBSHTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenBSHTransactorSession struct {
	Contract     *TokenBSHTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TokenBSHRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenBSHRaw struct {
	Contract *TokenBSH // Generic contract binding to access the raw methods on
}

// TokenBSHCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenBSHCallerRaw struct {
	Contract *TokenBSHCaller // Generic read-only contract binding to access the raw methods on
}

// TokenBSHTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenBSHTransactorRaw struct {
	Contract *TokenBSHTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenBSH creates a new instance of TokenBSH, bound to a specific deployed contract.
func NewTokenBSH(address common.Address, backend bind.ContractBackend) (*TokenBSH, error) {
	contract, err := bindTokenBSH(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenBSH{TokenBSHCaller: TokenBSHCaller{contract: contract}, TokenBSHTransactor: TokenBSHTransactor{contract: contract}, TokenBSHFilterer: TokenBSHFilterer{contract: contract}}, nil
}

// NewTokenBSHCaller creates a new read-only instance of TokenBSH, bound to a specific deployed contract.
func NewTokenBSHCaller(address common.Address, caller bind.ContractCaller) (*TokenBSHCaller, error) {
	contract, err := bindTokenBSH(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenBSHCaller{contract: contract}, nil
}

// NewTokenBSHTransactor creates a new write-only instance of TokenBSH, bound to a specific deployed contract.
func NewTokenBSHTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenBSHTransactor, error) {
	contract, err := bindTokenBSH(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenBSHTransactor{contract: contract}, nil
}

// NewTokenBSHFilterer creates a new log filterer instance of TokenBSH, bound to a specific deployed contract.
func NewTokenBSHFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenBSHFilterer, error) {
	contract, err := bindTokenBSH(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenBSHFilterer{contract: contract}, nil
}

// bindTokenBSH binds a generic wrapper to an already deployed contract.
func bindTokenBSH(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenBSHABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenBSH *TokenBSHRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenBSH.Contract.TokenBSHCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenBSH *TokenBSHRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenBSH.Contract.TokenBSHTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenBSH *TokenBSHRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenBSH.Contract.TokenBSHTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenBSH *TokenBSHCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenBSH.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenBSH *TokenBSHTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenBSH.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenBSH *TokenBSHTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenBSH.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenBSH *TokenBSHCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenBSH *TokenBSHSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TokenBSH.Contract.Allowance(&_TokenBSH.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenBSH *TokenBSHCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TokenBSH.Contract.Allowance(&_TokenBSH.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenBSH *TokenBSHCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenBSH *TokenBSHSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TokenBSH.Contract.BalanceOf(&_TokenBSH.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenBSH *TokenBSHCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TokenBSH.Contract.BalanceOf(&_TokenBSH.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenBSH *TokenBSHCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenBSH *TokenBSHSession) Decimals() (uint8, error) {
	return _TokenBSH.Contract.Decimals(&_TokenBSH.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenBSH *TokenBSHCallerSession) Decimals() (uint8, error) {
	return _TokenBSH.Contract.Decimals(&_TokenBSH.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_TokenBSH *TokenBSHCaller) IsOwner(opts *bind.CallOpts, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "isOwner", _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_TokenBSH *TokenBSHSession) IsOwner(_owner common.Address) (bool, error) {
	return _TokenBSH.Contract.IsOwner(&_TokenBSH.CallOpts, _owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_TokenBSH *TokenBSHCallerSession) IsOwner(_owner common.Address) (bool, error) {
	return _TokenBSH.Contract.IsOwner(&_TokenBSH.CallOpts, _owner)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenBSH *TokenBSHCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenBSH *TokenBSHSession) Name() (string, error) {
	return _TokenBSH.Contract.Name(&_TokenBSH.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenBSH *TokenBSHCallerSession) Name() (string, error) {
	return _TokenBSH.Contract.Name(&_TokenBSH.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenBSH *TokenBSHCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenBSH *TokenBSHSession) Symbol() (string, error) {
	return _TokenBSH.Contract.Symbol(&_TokenBSH.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenBSH *TokenBSHCallerSession) Symbol() (string, error) {
	return _TokenBSH.Contract.Symbol(&_TokenBSH.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenBSH *TokenBSHCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenBSH.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenBSH *TokenBSHSession) TotalSupply() (*big.Int, error) {
	return _TokenBSH.Contract.TotalSupply(&_TokenBSH.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenBSH *TokenBSHCallerSession) TotalSupply() (*big.Int, error) {
	return _TokenBSH.Contract.TotalSupply(&_TokenBSH.CallOpts)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_TokenBSH *TokenBSHTransactor) AddOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "addOwner", _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_TokenBSH *TokenBSHSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _TokenBSH.Contract.AddOwner(&_TokenBSH.TransactOpts, _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_TokenBSH *TokenBSHTransactorSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _TokenBSH.Contract.AddOwner(&_TokenBSH.TransactOpts, _owner)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.Approve(&_TokenBSH.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.Approve(&_TokenBSH.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenBSH *TokenBSHTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenBSH *TokenBSHSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.DecreaseAllowance(&_TokenBSH.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenBSH *TokenBSHTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.DecreaseAllowance(&_TokenBSH.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenBSH *TokenBSHTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenBSH *TokenBSHSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.IncreaseAllowance(&_TokenBSH.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenBSH *TokenBSHTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.IncreaseAllowance(&_TokenBSH.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactor) Mint(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "mint", amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.Mint(&_TokenBSH.TransactOpts, amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactorSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.Mint(&_TokenBSH.TransactOpts, amount)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_TokenBSH *TokenBSHTransactor) RemoveOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "removeOwner", _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_TokenBSH *TokenBSHSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _TokenBSH.Contract.RemoveOwner(&_TokenBSH.TransactOpts, _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_TokenBSH *TokenBSHTransactorSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _TokenBSH.Contract.RemoveOwner(&_TokenBSH.TransactOpts, _owner)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.Transfer(&_TokenBSH.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.Transfer(&_TokenBSH.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.TransferFrom(&_TokenBSH.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TokenBSH *TokenBSHTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBSH.Contract.TransferFrom(&_TokenBSH.TransactOpts, sender, recipient, amount)
}

// TokenBSHApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TokenBSH contract.
type TokenBSHApprovalIterator struct {
	Event *TokenBSHApproval // Event containing the contract specifics and raw log

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
func (it *TokenBSHApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBSHApproval)
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
		it.Event = new(TokenBSHApproval)
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
func (it *TokenBSHApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBSHApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBSHApproval represents a Approval event raised by the TokenBSH contract.
type TokenBSHApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenBSH *TokenBSHFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TokenBSHApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TokenBSH.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TokenBSHApprovalIterator{contract: _TokenBSH.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenBSH *TokenBSHFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TokenBSHApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TokenBSH.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBSHApproval)
				if err := _TokenBSH.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenBSH *TokenBSHFilterer) ParseApproval(log types.Log) (*TokenBSHApproval, error) {
	event := new(TokenBSHApproval)
	if err := _TokenBSH.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenBSHTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TokenBSH contract.
type TokenBSHTransferIterator struct {
	Event *TokenBSHTransfer // Event containing the contract specifics and raw log

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
func (it *TokenBSHTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBSHTransfer)
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
		it.Event = new(TokenBSHTransfer)
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
func (it *TokenBSHTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBSHTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBSHTransfer represents a Transfer event raised by the TokenBSH contract.
type TokenBSHTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenBSH *TokenBSHFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenBSHTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenBSH.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenBSHTransferIterator{contract: _TokenBSH.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenBSH *TokenBSHFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TokenBSHTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenBSH.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBSHTransfer)
				if err := _TokenBSH.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenBSH *TokenBSHFilterer) ParseTransfer(log types.Log) (*TokenBSHTransfer, error) {
	event := new(TokenBSHTransfer)
	if err := _TokenBSH.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
*/
