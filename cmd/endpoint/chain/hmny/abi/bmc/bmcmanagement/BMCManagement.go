// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bmcmanagement

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

// TypesLink is an auto generated low-level Go binding around an user-defined struct.
type TypesLink struct {
	Relays           []common.Address
	Reachable        []string
	RxSeq            *big.Int
	TxSeq            *big.Int
	BlockIntervalSrc *big.Int
	BlockIntervalDst *big.Int
	MaxAggregation   *big.Int
	DelayLimit       *big.Int
	RelayIdx         *big.Int
	RotateHeight     *big.Int
	RxHeight         *big.Int
	RxHeightSrc      *big.Int
	IsConnected      bool
}

// TypesRelayStats is an auto generated low-level Go binding around an user-defined struct.
type TypesRelayStats struct {
	Addr       common.Address
	BlockCount *big.Int
	MsgCount   *big.Int
}

// TypesRoute is an auto generated low-level Go binding around an user-defined struct.
type TypesRoute struct {
	Dst  string
	Next string
}

// TypesService is an auto generated low-level Go binding around an user-defined struct.
type TypesService struct {
	Svc  string
	Addr common.Address
}

// BmcmanagementABI is the input ABI used to generate the binding from.
const BmcmanagementABI = "[{\"inputs\":[],\"name\":\"serialNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setBMCPeriphery\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"addService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"}],\"name\":\"removeService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getServices\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"svc\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"internalType\":\"structTypes.Service[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"}],\"name\":\"addLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"}],\"name\":\"removeLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinks\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"name\":\"setLinkRxHeight\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_blockInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxAggregation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_delayLimit\",\"type\":\"uint256\"}],\"name\":\"setLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_currentHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_relayMsgHeight\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_hasMsg\",\"type\":\"bool\"}],\"name\":\"rotateRelay\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dst\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"}],\"name\":\"addRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dst\",\"type\":\"string\"}],\"name\":\"removeRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRoutes\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"dst\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"next\",\"type\":\"string\"}],\"internalType\":\"structTypes.Route[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"_addr\",\"type\":\"address[]\"}],\"name\":\"addRelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"removeRelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"}],\"name\":\"getRelays\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_serviceName\",\"type\":\"string\"}],\"name\":\"getBshServiceByName\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"}],\"name\":\"getLink\",\"outputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"relays\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"reachable\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"rxSeq\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txSeq\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockIntervalSrc\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockIntervalDst\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxAggregation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delayLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"relayIdx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rotateHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rxHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rxHeightSrc\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isConnected\",\"type\":\"bool\"}],\"internalType\":\"structTypes.Link\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"}],\"name\":\"getLinkRxSeq\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"}],\"name\":\"getLinkTxSeq\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"}],\"name\":\"getLinkRxHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"}],\"name\":\"getLinkRelays\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"}],\"name\":\"getRelayStatusByLink\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"blockCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"msgCount\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.RelayStats[]\",\"name\":\"_relays\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_val\",\"type\":\"uint256\"}],\"name\":\"updateLinkRxSeq\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"}],\"name\":\"updateLinkTxSeq\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_val\",\"type\":\"uint256\"}],\"name\":\"updateLinkRxHeight\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_to\",\"type\":\"string[]\"}],\"name\":\"updateLinkReachable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"deleteLinkReachable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relay\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_blockCountVal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_msgCountVal\",\"type\":\"uint256\"}],\"name\":\"updateRelayStats\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstNet\",\"type\":\"string\"}],\"name\":\"resolveRoute\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Bmcmanagement is an auto generated Go binding around an Ethereum contract.
type Bmcmanagement struct {
	BmcmanagementCaller     // Read-only binding to the contract
	BmcmanagementTransactor // Write-only binding to the contract
	BmcmanagementFilterer   // Log filterer for contract events
}

// BmcmanagementCaller is an auto generated read-only Go binding around an Ethereum contract.
type BmcmanagementCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BmcmanagementTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BmcmanagementTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BmcmanagementFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BmcmanagementFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BmcmanagementSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BmcmanagementSession struct {
	Contract     *Bmcmanagement    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BmcmanagementCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BmcmanagementCallerSession struct {
	Contract *BmcmanagementCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BmcmanagementTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BmcmanagementTransactorSession struct {
	Contract     *BmcmanagementTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BmcmanagementRaw is an auto generated low-level Go binding around an Ethereum contract.
type BmcmanagementRaw struct {
	Contract *Bmcmanagement // Generic contract binding to access the raw methods on
}

// BmcmanagementCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BmcmanagementCallerRaw struct {
	Contract *BmcmanagementCaller // Generic read-only contract binding to access the raw methods on
}

// BmcmanagementTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BmcmanagementTransactorRaw struct {
	Contract *BmcmanagementTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBmcmanagement creates a new instance of Bmcmanagement, bound to a specific deployed contract.
func NewBmcmanagement(address common.Address, backend bind.ContractBackend) (*Bmcmanagement, error) {
	contract, err := bindBmcmanagement(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bmcmanagement{BmcmanagementCaller: BmcmanagementCaller{contract: contract}, BmcmanagementTransactor: BmcmanagementTransactor{contract: contract}, BmcmanagementFilterer: BmcmanagementFilterer{contract: contract}}, nil
}

// NewBmcmanagementCaller creates a new read-only instance of Bmcmanagement, bound to a specific deployed contract.
func NewBmcmanagementCaller(address common.Address, caller bind.ContractCaller) (*BmcmanagementCaller, error) {
	contract, err := bindBmcmanagement(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BmcmanagementCaller{contract: contract}, nil
}

// NewBmcmanagementTransactor creates a new write-only instance of Bmcmanagement, bound to a specific deployed contract.
func NewBmcmanagementTransactor(address common.Address, transactor bind.ContractTransactor) (*BmcmanagementTransactor, error) {
	contract, err := bindBmcmanagement(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BmcmanagementTransactor{contract: contract}, nil
}

// NewBmcmanagementFilterer creates a new log filterer instance of Bmcmanagement, bound to a specific deployed contract.
func NewBmcmanagementFilterer(address common.Address, filterer bind.ContractFilterer) (*BmcmanagementFilterer, error) {
	contract, err := bindBmcmanagement(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BmcmanagementFilterer{contract: contract}, nil
}

// bindBmcmanagement binds a generic wrapper to an already deployed contract.
func bindBmcmanagement(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BmcmanagementABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bmcmanagement *BmcmanagementRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bmcmanagement.Contract.BmcmanagementCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bmcmanagement *BmcmanagementRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.BmcmanagementTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bmcmanagement *BmcmanagementRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.BmcmanagementTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bmcmanagement *BmcmanagementCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bmcmanagement.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bmcmanagement *BmcmanagementTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bmcmanagement *BmcmanagementTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.contract.Transact(opts, method, params...)
}

// GetBshServiceByName is a free data retrieval call binding the contract method 0xcd881dbd.
//
// Solidity: function getBshServiceByName(string _serviceName) view returns(address)
func (_Bmcmanagement *BmcmanagementCaller) GetBshServiceByName(opts *bind.CallOpts, _serviceName string) (common.Address, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getBshServiceByName", _serviceName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBshServiceByName is a free data retrieval call binding the contract method 0xcd881dbd.
//
// Solidity: function getBshServiceByName(string _serviceName) view returns(address)
func (_Bmcmanagement *BmcmanagementSession) GetBshServiceByName(_serviceName string) (common.Address, error) {
	return _Bmcmanagement.Contract.GetBshServiceByName(&_Bmcmanagement.CallOpts, _serviceName)
}

// GetBshServiceByName is a free data retrieval call binding the contract method 0xcd881dbd.
//
// Solidity: function getBshServiceByName(string _serviceName) view returns(address)
func (_Bmcmanagement *BmcmanagementCallerSession) GetBshServiceByName(_serviceName string) (common.Address, error) {
	return _Bmcmanagement.Contract.GetBshServiceByName(&_Bmcmanagement.CallOpts, _serviceName)
}

// GetLink is a free data retrieval call binding the contract method 0x7476452a.
//
// Solidity: function getLink(string _to) view returns((address[],string[],uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool))
func (_Bmcmanagement *BmcmanagementCaller) GetLink(opts *bind.CallOpts, _to string) (TypesLink, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getLink", _to)

	if err != nil {
		return *new(TypesLink), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesLink)).(*TypesLink)

	return out0, err

}

// GetLink is a free data retrieval call binding the contract method 0x7476452a.
//
// Solidity: function getLink(string _to) view returns((address[],string[],uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool))
func (_Bmcmanagement *BmcmanagementSession) GetLink(_to string) (TypesLink, error) {
	return _Bmcmanagement.Contract.GetLink(&_Bmcmanagement.CallOpts, _to)
}

// GetLink is a free data retrieval call binding the contract method 0x7476452a.
//
// Solidity: function getLink(string _to) view returns((address[],string[],uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool))
func (_Bmcmanagement *BmcmanagementCallerSession) GetLink(_to string) (TypesLink, error) {
	return _Bmcmanagement.Contract.GetLink(&_Bmcmanagement.CallOpts, _to)
}

// GetLinkRelays is a free data retrieval call binding the contract method 0x401a2c57.
//
// Solidity: function getLinkRelays(string _prev) view returns(address[])
func (_Bmcmanagement *BmcmanagementCaller) GetLinkRelays(opts *bind.CallOpts, _prev string) ([]common.Address, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getLinkRelays", _prev)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetLinkRelays is a free data retrieval call binding the contract method 0x401a2c57.
//
// Solidity: function getLinkRelays(string _prev) view returns(address[])
func (_Bmcmanagement *BmcmanagementSession) GetLinkRelays(_prev string) ([]common.Address, error) {
	return _Bmcmanagement.Contract.GetLinkRelays(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkRelays is a free data retrieval call binding the contract method 0x401a2c57.
//
// Solidity: function getLinkRelays(string _prev) view returns(address[])
func (_Bmcmanagement *BmcmanagementCallerSession) GetLinkRelays(_prev string) ([]common.Address, error) {
	return _Bmcmanagement.Contract.GetLinkRelays(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkRxHeight is a free data retrieval call binding the contract method 0xb589082d.
//
// Solidity: function getLinkRxHeight(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementCaller) GetLinkRxHeight(opts *bind.CallOpts, _prev string) (*big.Int, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getLinkRxHeight", _prev)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLinkRxHeight is a free data retrieval call binding the contract method 0xb589082d.
//
// Solidity: function getLinkRxHeight(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementSession) GetLinkRxHeight(_prev string) (*big.Int, error) {
	return _Bmcmanagement.Contract.GetLinkRxHeight(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkRxHeight is a free data retrieval call binding the contract method 0xb589082d.
//
// Solidity: function getLinkRxHeight(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementCallerSession) GetLinkRxHeight(_prev string) (*big.Int, error) {
	return _Bmcmanagement.Contract.GetLinkRxHeight(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkRxSeq is a free data retrieval call binding the contract method 0x535d362b.
//
// Solidity: function getLinkRxSeq(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementCaller) GetLinkRxSeq(opts *bind.CallOpts, _prev string) (*big.Int, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getLinkRxSeq", _prev)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLinkRxSeq is a free data retrieval call binding the contract method 0x535d362b.
//
// Solidity: function getLinkRxSeq(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementSession) GetLinkRxSeq(_prev string) (*big.Int, error) {
	return _Bmcmanagement.Contract.GetLinkRxSeq(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkRxSeq is a free data retrieval call binding the contract method 0x535d362b.
//
// Solidity: function getLinkRxSeq(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementCallerSession) GetLinkRxSeq(_prev string) (*big.Int, error) {
	return _Bmcmanagement.Contract.GetLinkRxSeq(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkTxSeq is a free data retrieval call binding the contract method 0x0a05d5ae.
//
// Solidity: function getLinkTxSeq(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementCaller) GetLinkTxSeq(opts *bind.CallOpts, _prev string) (*big.Int, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getLinkTxSeq", _prev)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLinkTxSeq is a free data retrieval call binding the contract method 0x0a05d5ae.
//
// Solidity: function getLinkTxSeq(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementSession) GetLinkTxSeq(_prev string) (*big.Int, error) {
	return _Bmcmanagement.Contract.GetLinkTxSeq(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinkTxSeq is a free data retrieval call binding the contract method 0x0a05d5ae.
//
// Solidity: function getLinkTxSeq(string _prev) view returns(uint256)
func (_Bmcmanagement *BmcmanagementCallerSession) GetLinkTxSeq(_prev string) (*big.Int, error) {
	return _Bmcmanagement.Contract.GetLinkTxSeq(&_Bmcmanagement.CallOpts, _prev)
}

// GetLinks is a free data retrieval call binding the contract method 0xf66ddcbb.
//
// Solidity: function getLinks() view returns(string[])
func (_Bmcmanagement *BmcmanagementCaller) GetLinks(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getLinks")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetLinks is a free data retrieval call binding the contract method 0xf66ddcbb.
//
// Solidity: function getLinks() view returns(string[])
func (_Bmcmanagement *BmcmanagementSession) GetLinks() ([]string, error) {
	return _Bmcmanagement.Contract.GetLinks(&_Bmcmanagement.CallOpts)
}

// GetLinks is a free data retrieval call binding the contract method 0xf66ddcbb.
//
// Solidity: function getLinks() view returns(string[])
func (_Bmcmanagement *BmcmanagementCallerSession) GetLinks() ([]string, error) {
	return _Bmcmanagement.Contract.GetLinks(&_Bmcmanagement.CallOpts)
}

// GetRelayStatusByLink is a free data retrieval call binding the contract method 0x5fb09c13.
//
// Solidity: function getRelayStatusByLink(string _prev) view returns((address,uint256,uint256)[] _relays)
func (_Bmcmanagement *BmcmanagementCaller) GetRelayStatusByLink(opts *bind.CallOpts, _prev string) ([]TypesRelayStats, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getRelayStatusByLink", _prev)

	if err != nil {
		return *new([]TypesRelayStats), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesRelayStats)).(*[]TypesRelayStats)

	return out0, err

}

// GetRelayStatusByLink is a free data retrieval call binding the contract method 0x5fb09c13.
//
// Solidity: function getRelayStatusByLink(string _prev) view returns((address,uint256,uint256)[] _relays)
func (_Bmcmanagement *BmcmanagementSession) GetRelayStatusByLink(_prev string) ([]TypesRelayStats, error) {
	return _Bmcmanagement.Contract.GetRelayStatusByLink(&_Bmcmanagement.CallOpts, _prev)
}

// GetRelayStatusByLink is a free data retrieval call binding the contract method 0x5fb09c13.
//
// Solidity: function getRelayStatusByLink(string _prev) view returns((address,uint256,uint256)[] _relays)
func (_Bmcmanagement *BmcmanagementCallerSession) GetRelayStatusByLink(_prev string) ([]TypesRelayStats, error) {
	return _Bmcmanagement.Contract.GetRelayStatusByLink(&_Bmcmanagement.CallOpts, _prev)
}

// GetRelays is a free data retrieval call binding the contract method 0x40926734.
//
// Solidity: function getRelays(string _link) view returns(address[])
func (_Bmcmanagement *BmcmanagementCaller) GetRelays(opts *bind.CallOpts, _link string) ([]common.Address, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getRelays", _link)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRelays is a free data retrieval call binding the contract method 0x40926734.
//
// Solidity: function getRelays(string _link) view returns(address[])
func (_Bmcmanagement *BmcmanagementSession) GetRelays(_link string) ([]common.Address, error) {
	return _Bmcmanagement.Contract.GetRelays(&_Bmcmanagement.CallOpts, _link)
}

// GetRelays is a free data retrieval call binding the contract method 0x40926734.
//
// Solidity: function getRelays(string _link) view returns(address[])
func (_Bmcmanagement *BmcmanagementCallerSession) GetRelays(_link string) ([]common.Address, error) {
	return _Bmcmanagement.Contract.GetRelays(&_Bmcmanagement.CallOpts, _link)
}

// GetRoutes is a free data retrieval call binding the contract method 0x7e928072.
//
// Solidity: function getRoutes() view returns((string,string)[])
func (_Bmcmanagement *BmcmanagementCaller) GetRoutes(opts *bind.CallOpts) ([]TypesRoute, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getRoutes")

	if err != nil {
		return *new([]TypesRoute), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesRoute)).(*[]TypesRoute)

	return out0, err

}

// GetRoutes is a free data retrieval call binding the contract method 0x7e928072.
//
// Solidity: function getRoutes() view returns((string,string)[])
func (_Bmcmanagement *BmcmanagementSession) GetRoutes() ([]TypesRoute, error) {
	return _Bmcmanagement.Contract.GetRoutes(&_Bmcmanagement.CallOpts)
}

// GetRoutes is a free data retrieval call binding the contract method 0x7e928072.
//
// Solidity: function getRoutes() view returns((string,string)[])
func (_Bmcmanagement *BmcmanagementCallerSession) GetRoutes() ([]TypesRoute, error) {
	return _Bmcmanagement.Contract.GetRoutes(&_Bmcmanagement.CallOpts)
}

// GetServices is a free data retrieval call binding the contract method 0x75417851.
//
// Solidity: function getServices() view returns((string,address)[])
func (_Bmcmanagement *BmcmanagementCaller) GetServices(opts *bind.CallOpts) ([]TypesService, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "getServices")

	if err != nil {
		return *new([]TypesService), err
	}

	out0 := *abi.ConvertType(out[0], new([]TypesService)).(*[]TypesService)

	return out0, err

}

// GetServices is a free data retrieval call binding the contract method 0x75417851.
//
// Solidity: function getServices() view returns((string,address)[])
func (_Bmcmanagement *BmcmanagementSession) GetServices() ([]TypesService, error) {
	return _Bmcmanagement.Contract.GetServices(&_Bmcmanagement.CallOpts)
}

// GetServices is a free data retrieval call binding the contract method 0x75417851.
//
// Solidity: function getServices() view returns((string,address)[])
func (_Bmcmanagement *BmcmanagementCallerSession) GetServices() ([]TypesService, error) {
	return _Bmcmanagement.Contract.GetServices(&_Bmcmanagement.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Bmcmanagement *BmcmanagementCaller) IsOwner(opts *bind.CallOpts, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "isOwner", _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Bmcmanagement *BmcmanagementSession) IsOwner(_owner common.Address) (bool, error) {
	return _Bmcmanagement.Contract.IsOwner(&_Bmcmanagement.CallOpts, _owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address _owner) view returns(bool)
func (_Bmcmanagement *BmcmanagementCallerSession) IsOwner(_owner common.Address) (bool, error) {
	return _Bmcmanagement.Contract.IsOwner(&_Bmcmanagement.CallOpts, _owner)
}

// ResolveRoute is a free data retrieval call binding the contract method 0xbe7f8676.
//
// Solidity: function resolveRoute(string _dstNet) view returns(string, string)
func (_Bmcmanagement *BmcmanagementCaller) ResolveRoute(opts *bind.CallOpts, _dstNet string) (string, string, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "resolveRoute", _dstNet)

	if err != nil {
		return *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// ResolveRoute is a free data retrieval call binding the contract method 0xbe7f8676.
//
// Solidity: function resolveRoute(string _dstNet) view returns(string, string)
func (_Bmcmanagement *BmcmanagementSession) ResolveRoute(_dstNet string) (string, string, error) {
	return _Bmcmanagement.Contract.ResolveRoute(&_Bmcmanagement.CallOpts, _dstNet)
}

// ResolveRoute is a free data retrieval call binding the contract method 0xbe7f8676.
//
// Solidity: function resolveRoute(string _dstNet) view returns(string, string)
func (_Bmcmanagement *BmcmanagementCallerSession) ResolveRoute(_dstNet string) (string, string, error) {
	return _Bmcmanagement.Contract.ResolveRoute(&_Bmcmanagement.CallOpts, _dstNet)
}

// SerialNo is a free data retrieval call binding the contract method 0x660f17fe.
//
// Solidity: function serialNo() view returns(uint256)
func (_Bmcmanagement *BmcmanagementCaller) SerialNo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bmcmanagement.contract.Call(opts, &out, "serialNo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SerialNo is a free data retrieval call binding the contract method 0x660f17fe.
//
// Solidity: function serialNo() view returns(uint256)
func (_Bmcmanagement *BmcmanagementSession) SerialNo() (*big.Int, error) {
	return _Bmcmanagement.Contract.SerialNo(&_Bmcmanagement.CallOpts)
}

// SerialNo is a free data retrieval call binding the contract method 0x660f17fe.
//
// Solidity: function serialNo() view returns(uint256)
func (_Bmcmanagement *BmcmanagementCallerSession) SerialNo() (*big.Int, error) {
	return _Bmcmanagement.Contract.SerialNo(&_Bmcmanagement.CallOpts)
}

// AddLink is a paid mutator transaction binding the contract method 0x22a618fa.
//
// Solidity: function addLink(string _link) returns()
func (_Bmcmanagement *BmcmanagementTransactor) AddLink(opts *bind.TransactOpts, _link string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "addLink", _link)
}

// AddLink is a paid mutator transaction binding the contract method 0x22a618fa.
//
// Solidity: function addLink(string _link) returns()
func (_Bmcmanagement *BmcmanagementSession) AddLink(_link string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddLink(&_Bmcmanagement.TransactOpts, _link)
}

// AddLink is a paid mutator transaction binding the contract method 0x22a618fa.
//
// Solidity: function addLink(string _link) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) AddLink(_link string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddLink(&_Bmcmanagement.TransactOpts, _link)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Bmcmanagement *BmcmanagementTransactor) AddOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "addOwner", _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Bmcmanagement *BmcmanagementSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddOwner(&_Bmcmanagement.TransactOpts, _owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address _owner) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) AddOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddOwner(&_Bmcmanagement.TransactOpts, _owner)
}

// AddRelay is a paid mutator transaction binding the contract method 0x0748ea7a.
//
// Solidity: function addRelay(string _link, address[] _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactor) AddRelay(opts *bind.TransactOpts, _link string, _addr []common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "addRelay", _link, _addr)
}

// AddRelay is a paid mutator transaction binding the contract method 0x0748ea7a.
//
// Solidity: function addRelay(string _link, address[] _addr) returns()
func (_Bmcmanagement *BmcmanagementSession) AddRelay(_link string, _addr []common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddRelay(&_Bmcmanagement.TransactOpts, _link, _addr)
}

// AddRelay is a paid mutator transaction binding the contract method 0x0748ea7a.
//
// Solidity: function addRelay(string _link, address[] _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) AddRelay(_link string, _addr []common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddRelay(&_Bmcmanagement.TransactOpts, _link, _addr)
}

// AddRoute is a paid mutator transaction binding the contract method 0x065a9e9b.
//
// Solidity: function addRoute(string _dst, string _link) returns()
func (_Bmcmanagement *BmcmanagementTransactor) AddRoute(opts *bind.TransactOpts, _dst string, _link string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "addRoute", _dst, _link)
}

// AddRoute is a paid mutator transaction binding the contract method 0x065a9e9b.
//
// Solidity: function addRoute(string _dst, string _link) returns()
func (_Bmcmanagement *BmcmanagementSession) AddRoute(_dst string, _link string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddRoute(&_Bmcmanagement.TransactOpts, _dst, _link)
}

// AddRoute is a paid mutator transaction binding the contract method 0x065a9e9b.
//
// Solidity: function addRoute(string _dst, string _link) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) AddRoute(_dst string, _link string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddRoute(&_Bmcmanagement.TransactOpts, _dst, _link)
}

// AddService is a paid mutator transaction binding the contract method 0x6d8c73d7.
//
// Solidity: function addService(string _svc, address _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactor) AddService(opts *bind.TransactOpts, _svc string, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "addService", _svc, _addr)
}

// AddService is a paid mutator transaction binding the contract method 0x6d8c73d7.
//
// Solidity: function addService(string _svc, address _addr) returns()
func (_Bmcmanagement *BmcmanagementSession) AddService(_svc string, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddService(&_Bmcmanagement.TransactOpts, _svc, _addr)
}

// AddService is a paid mutator transaction binding the contract method 0x6d8c73d7.
//
// Solidity: function addService(string _svc, address _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) AddService(_svc string, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.AddService(&_Bmcmanagement.TransactOpts, _svc, _addr)
}

// DeleteLinkReachable is a paid mutator transaction binding the contract method 0xb6b09484.
//
// Solidity: function deleteLinkReachable(string _prev, uint256 _index) returns()
func (_Bmcmanagement *BmcmanagementTransactor) DeleteLinkReachable(opts *bind.TransactOpts, _prev string, _index *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "deleteLinkReachable", _prev, _index)
}

// DeleteLinkReachable is a paid mutator transaction binding the contract method 0xb6b09484.
//
// Solidity: function deleteLinkReachable(string _prev, uint256 _index) returns()
func (_Bmcmanagement *BmcmanagementSession) DeleteLinkReachable(_prev string, _index *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.DeleteLinkReachable(&_Bmcmanagement.TransactOpts, _prev, _index)
}

// DeleteLinkReachable is a paid mutator transaction binding the contract method 0xb6b09484.
//
// Solidity: function deleteLinkReachable(string _prev, uint256 _index) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) DeleteLinkReachable(_prev string, _index *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.DeleteLinkReachable(&_Bmcmanagement.TransactOpts, _prev, _index)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Bmcmanagement *BmcmanagementTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Bmcmanagement *BmcmanagementSession) Initialize() (*types.Transaction, error) {
	return _Bmcmanagement.Contract.Initialize(&_Bmcmanagement.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) Initialize() (*types.Transaction, error) {
	return _Bmcmanagement.Contract.Initialize(&_Bmcmanagement.TransactOpts)
}

// RemoveLink is a paid mutator transaction binding the contract method 0x6e4060d7.
//
// Solidity: function removeLink(string _link) returns()
func (_Bmcmanagement *BmcmanagementTransactor) RemoveLink(opts *bind.TransactOpts, _link string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "removeLink", _link)
}

// RemoveLink is a paid mutator transaction binding the contract method 0x6e4060d7.
//
// Solidity: function removeLink(string _link) returns()
func (_Bmcmanagement *BmcmanagementSession) RemoveLink(_link string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveLink(&_Bmcmanagement.TransactOpts, _link)
}

// RemoveLink is a paid mutator transaction binding the contract method 0x6e4060d7.
//
// Solidity: function removeLink(string _link) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) RemoveLink(_link string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveLink(&_Bmcmanagement.TransactOpts, _link)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Bmcmanagement *BmcmanagementTransactor) RemoveOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "removeOwner", _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Bmcmanagement *BmcmanagementSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveOwner(&_Bmcmanagement.TransactOpts, _owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address _owner) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) RemoveOwner(_owner common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveOwner(&_Bmcmanagement.TransactOpts, _owner)
}

// RemoveRelay is a paid mutator transaction binding the contract method 0xdef59f5e.
//
// Solidity: function removeRelay(string _link, address _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactor) RemoveRelay(opts *bind.TransactOpts, _link string, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "removeRelay", _link, _addr)
}

// RemoveRelay is a paid mutator transaction binding the contract method 0xdef59f5e.
//
// Solidity: function removeRelay(string _link, address _addr) returns()
func (_Bmcmanagement *BmcmanagementSession) RemoveRelay(_link string, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveRelay(&_Bmcmanagement.TransactOpts, _link, _addr)
}

// RemoveRelay is a paid mutator transaction binding the contract method 0xdef59f5e.
//
// Solidity: function removeRelay(string _link, address _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) RemoveRelay(_link string, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveRelay(&_Bmcmanagement.TransactOpts, _link, _addr)
}

// RemoveRoute is a paid mutator transaction binding the contract method 0xbd0a0bb3.
//
// Solidity: function removeRoute(string _dst) returns()
func (_Bmcmanagement *BmcmanagementTransactor) RemoveRoute(opts *bind.TransactOpts, _dst string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "removeRoute", _dst)
}

// RemoveRoute is a paid mutator transaction binding the contract method 0xbd0a0bb3.
//
// Solidity: function removeRoute(string _dst) returns()
func (_Bmcmanagement *BmcmanagementSession) RemoveRoute(_dst string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveRoute(&_Bmcmanagement.TransactOpts, _dst)
}

// RemoveRoute is a paid mutator transaction binding the contract method 0xbd0a0bb3.
//
// Solidity: function removeRoute(string _dst) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) RemoveRoute(_dst string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveRoute(&_Bmcmanagement.TransactOpts, _dst)
}

// RemoveService is a paid mutator transaction binding the contract method 0xf51acaea.
//
// Solidity: function removeService(string _svc) returns()
func (_Bmcmanagement *BmcmanagementTransactor) RemoveService(opts *bind.TransactOpts, _svc string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "removeService", _svc)
}

// RemoveService is a paid mutator transaction binding the contract method 0xf51acaea.
//
// Solidity: function removeService(string _svc) returns()
func (_Bmcmanagement *BmcmanagementSession) RemoveService(_svc string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveService(&_Bmcmanagement.TransactOpts, _svc)
}

// RemoveService is a paid mutator transaction binding the contract method 0xf51acaea.
//
// Solidity: function removeService(string _svc) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) RemoveService(_svc string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RemoveService(&_Bmcmanagement.TransactOpts, _svc)
}

// RotateRelay is a paid mutator transaction binding the contract method 0xe88f3cd1.
//
// Solidity: function rotateRelay(string _link, uint256 _currentHeight, uint256 _relayMsgHeight, bool _hasMsg) returns(address)
func (_Bmcmanagement *BmcmanagementTransactor) RotateRelay(opts *bind.TransactOpts, _link string, _currentHeight *big.Int, _relayMsgHeight *big.Int, _hasMsg bool) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "rotateRelay", _link, _currentHeight, _relayMsgHeight, _hasMsg)
}

// RotateRelay is a paid mutator transaction binding the contract method 0xe88f3cd1.
//
// Solidity: function rotateRelay(string _link, uint256 _currentHeight, uint256 _relayMsgHeight, bool _hasMsg) returns(address)
func (_Bmcmanagement *BmcmanagementSession) RotateRelay(_link string, _currentHeight *big.Int, _relayMsgHeight *big.Int, _hasMsg bool) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RotateRelay(&_Bmcmanagement.TransactOpts, _link, _currentHeight, _relayMsgHeight, _hasMsg)
}

// RotateRelay is a paid mutator transaction binding the contract method 0xe88f3cd1.
//
// Solidity: function rotateRelay(string _link, uint256 _currentHeight, uint256 _relayMsgHeight, bool _hasMsg) returns(address)
func (_Bmcmanagement *BmcmanagementTransactorSession) RotateRelay(_link string, _currentHeight *big.Int, _relayMsgHeight *big.Int, _hasMsg bool) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.RotateRelay(&_Bmcmanagement.TransactOpts, _link, _currentHeight, _relayMsgHeight, _hasMsg)
}

// SetBMCPeriphery is a paid mutator transaction binding the contract method 0xc620234c.
//
// Solidity: function setBMCPeriphery(address _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactor) SetBMCPeriphery(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "setBMCPeriphery", _addr)
}

// SetBMCPeriphery is a paid mutator transaction binding the contract method 0xc620234c.
//
// Solidity: function setBMCPeriphery(address _addr) returns()
func (_Bmcmanagement *BmcmanagementSession) SetBMCPeriphery(_addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.SetBMCPeriphery(&_Bmcmanagement.TransactOpts, _addr)
}

// SetBMCPeriphery is a paid mutator transaction binding the contract method 0xc620234c.
//
// Solidity: function setBMCPeriphery(address _addr) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) SetBMCPeriphery(_addr common.Address) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.SetBMCPeriphery(&_Bmcmanagement.TransactOpts, _addr)
}

// SetLink is a paid mutator transaction binding the contract method 0xf216b155.
//
// Solidity: function setLink(string _link, uint256 _blockInterval, uint256 _maxAggregation, uint256 _delayLimit) returns()
func (_Bmcmanagement *BmcmanagementTransactor) SetLink(opts *bind.TransactOpts, _link string, _blockInterval *big.Int, _maxAggregation *big.Int, _delayLimit *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "setLink", _link, _blockInterval, _maxAggregation, _delayLimit)
}

// SetLink is a paid mutator transaction binding the contract method 0xf216b155.
//
// Solidity: function setLink(string _link, uint256 _blockInterval, uint256 _maxAggregation, uint256 _delayLimit) returns()
func (_Bmcmanagement *BmcmanagementSession) SetLink(_link string, _blockInterval *big.Int, _maxAggregation *big.Int, _delayLimit *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.SetLink(&_Bmcmanagement.TransactOpts, _link, _blockInterval, _maxAggregation, _delayLimit)
}

// SetLink is a paid mutator transaction binding the contract method 0xf216b155.
//
// Solidity: function setLink(string _link, uint256 _blockInterval, uint256 _maxAggregation, uint256 _delayLimit) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) SetLink(_link string, _blockInterval *big.Int, _maxAggregation *big.Int, _delayLimit *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.SetLink(&_Bmcmanagement.TransactOpts, _link, _blockInterval, _maxAggregation, _delayLimit)
}

// SetLinkRxHeight is a paid mutator transaction binding the contract method 0x0f95230f.
//
// Solidity: function setLinkRxHeight(string _link, uint256 _height) returns()
func (_Bmcmanagement *BmcmanagementTransactor) SetLinkRxHeight(opts *bind.TransactOpts, _link string, _height *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "setLinkRxHeight", _link, _height)
}

// SetLinkRxHeight is a paid mutator transaction binding the contract method 0x0f95230f.
//
// Solidity: function setLinkRxHeight(string _link, uint256 _height) returns()
func (_Bmcmanagement *BmcmanagementSession) SetLinkRxHeight(_link string, _height *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.SetLinkRxHeight(&_Bmcmanagement.TransactOpts, _link, _height)
}

// SetLinkRxHeight is a paid mutator transaction binding the contract method 0x0f95230f.
//
// Solidity: function setLinkRxHeight(string _link, uint256 _height) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) SetLinkRxHeight(_link string, _height *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.SetLinkRxHeight(&_Bmcmanagement.TransactOpts, _link, _height)
}

// UpdateLinkReachable is a paid mutator transaction binding the contract method 0xfc21d1c3.
//
// Solidity: function updateLinkReachable(string _prev, string[] _to) returns()
func (_Bmcmanagement *BmcmanagementTransactor) UpdateLinkReachable(opts *bind.TransactOpts, _prev string, _to []string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "updateLinkReachable", _prev, _to)
}

// UpdateLinkReachable is a paid mutator transaction binding the contract method 0xfc21d1c3.
//
// Solidity: function updateLinkReachable(string _prev, string[] _to) returns()
func (_Bmcmanagement *BmcmanagementSession) UpdateLinkReachable(_prev string, _to []string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkReachable(&_Bmcmanagement.TransactOpts, _prev, _to)
}

// UpdateLinkReachable is a paid mutator transaction binding the contract method 0xfc21d1c3.
//
// Solidity: function updateLinkReachable(string _prev, string[] _to) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) UpdateLinkReachable(_prev string, _to []string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkReachable(&_Bmcmanagement.TransactOpts, _prev, _to)
}

// UpdateLinkRxHeight is a paid mutator transaction binding the contract method 0x2f623dd1.
//
// Solidity: function updateLinkRxHeight(string _prev, uint256 _val) returns()
func (_Bmcmanagement *BmcmanagementTransactor) UpdateLinkRxHeight(opts *bind.TransactOpts, _prev string, _val *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "updateLinkRxHeight", _prev, _val)
}

// UpdateLinkRxHeight is a paid mutator transaction binding the contract method 0x2f623dd1.
//
// Solidity: function updateLinkRxHeight(string _prev, uint256 _val) returns()
func (_Bmcmanagement *BmcmanagementSession) UpdateLinkRxHeight(_prev string, _val *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkRxHeight(&_Bmcmanagement.TransactOpts, _prev, _val)
}

// UpdateLinkRxHeight is a paid mutator transaction binding the contract method 0x2f623dd1.
//
// Solidity: function updateLinkRxHeight(string _prev, uint256 _val) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) UpdateLinkRxHeight(_prev string, _val *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkRxHeight(&_Bmcmanagement.TransactOpts, _prev, _val)
}

// UpdateLinkRxSeq is a paid mutator transaction binding the contract method 0x7c38d594.
//
// Solidity: function updateLinkRxSeq(string _prev, uint256 _val) returns()
func (_Bmcmanagement *BmcmanagementTransactor) UpdateLinkRxSeq(opts *bind.TransactOpts, _prev string, _val *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "updateLinkRxSeq", _prev, _val)
}

// UpdateLinkRxSeq is a paid mutator transaction binding the contract method 0x7c38d594.
//
// Solidity: function updateLinkRxSeq(string _prev, uint256 _val) returns()
func (_Bmcmanagement *BmcmanagementSession) UpdateLinkRxSeq(_prev string, _val *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkRxSeq(&_Bmcmanagement.TransactOpts, _prev, _val)
}

// UpdateLinkRxSeq is a paid mutator transaction binding the contract method 0x7c38d594.
//
// Solidity: function updateLinkRxSeq(string _prev, uint256 _val) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) UpdateLinkRxSeq(_prev string, _val *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkRxSeq(&_Bmcmanagement.TransactOpts, _prev, _val)
}

// UpdateLinkTxSeq is a paid mutator transaction binding the contract method 0xa98978cb.
//
// Solidity: function updateLinkTxSeq(string _prev) returns()
func (_Bmcmanagement *BmcmanagementTransactor) UpdateLinkTxSeq(opts *bind.TransactOpts, _prev string) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "updateLinkTxSeq", _prev)
}

// UpdateLinkTxSeq is a paid mutator transaction binding the contract method 0xa98978cb.
//
// Solidity: function updateLinkTxSeq(string _prev) returns()
func (_Bmcmanagement *BmcmanagementSession) UpdateLinkTxSeq(_prev string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkTxSeq(&_Bmcmanagement.TransactOpts, _prev)
}

// UpdateLinkTxSeq is a paid mutator transaction binding the contract method 0xa98978cb.
//
// Solidity: function updateLinkTxSeq(string _prev) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) UpdateLinkTxSeq(_prev string) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateLinkTxSeq(&_Bmcmanagement.TransactOpts, _prev)
}

// UpdateRelayStats is a paid mutator transaction binding the contract method 0x754cd745.
//
// Solidity: function updateRelayStats(address relay, uint256 _blockCountVal, uint256 _msgCountVal) returns()
func (_Bmcmanagement *BmcmanagementTransactor) UpdateRelayStats(opts *bind.TransactOpts, relay common.Address, _blockCountVal *big.Int, _msgCountVal *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.contract.Transact(opts, "updateRelayStats", relay, _blockCountVal, _msgCountVal)
}

// UpdateRelayStats is a paid mutator transaction binding the contract method 0x754cd745.
//
// Solidity: function updateRelayStats(address relay, uint256 _blockCountVal, uint256 _msgCountVal) returns()
func (_Bmcmanagement *BmcmanagementSession) UpdateRelayStats(relay common.Address, _blockCountVal *big.Int, _msgCountVal *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateRelayStats(&_Bmcmanagement.TransactOpts, relay, _blockCountVal, _msgCountVal)
}

// UpdateRelayStats is a paid mutator transaction binding the contract method 0x754cd745.
//
// Solidity: function updateRelayStats(address relay, uint256 _blockCountVal, uint256 _msgCountVal) returns()
func (_Bmcmanagement *BmcmanagementTransactorSession) UpdateRelayStats(relay common.Address, _blockCountVal *big.Int, _msgCountVal *big.Int) (*types.Transaction, error) {
	return _Bmcmanagement.Contract.UpdateRelayStats(&_Bmcmanagement.TransactOpts, relay, _blockCountVal, _msgCountVal)
}
