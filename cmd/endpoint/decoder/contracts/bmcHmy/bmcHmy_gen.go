// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bmcHmy

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

// TypesBMCMessage is an auto generated low-level Go binding around an user-defined struct.
type TypesBMCMessage struct {
	Src     string
	Dst     string
	Svc     string
	Sn      *big.Int
	Message []byte
}

// TypesBMCService is an auto generated low-level Go binding around an user-defined struct.
type TypesBMCService struct {
	ServiceType string
	Payload     []byte
}

// TypesGatherFeeMessage is an auto generated low-level Go binding around an user-defined struct.
type TypesGatherFeeMessage struct {
	Fa   string
	Svcs []string
}

// TypesLinkStats is an auto generated low-level Go binding around an user-defined struct.
type TypesLinkStats struct {
	RxSeq         *big.Int
	TxSeq         *big.Int
	RxHeight      *big.Int
	CurrentHeight *big.Int
}

// BmcHmyABI is the input ABI used to generate the binding from.
const BmcHmyABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_sn\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_code\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_errMsg\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_svcErrCode\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_svcErrMsg\",\"type\":\"string\"}],\"name\":\"ErrorOnBTPError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_next\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_seq\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"Message\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_network\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_bmcManagementAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBmcBtpAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_prev\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"handleRelayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_rlp\",\"type\":\"bytes\"}],\"name\":\"tryDecodeBTPMessage\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"src\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dst\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"svc\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"sn\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.BMCMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"tryDecodeBMCService\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"serviceType\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.BMCService\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"tryDecodeGatherFeeMessage\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"fa\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"svcs\",\"type\":\"string[]\"}],\"internalType\":\"structTypes.GatherFeeMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_to\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_svc\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_sn\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"sendMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_link\",\"type\":\"string\"}],\"name\":\"getStatus\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rxSeq\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txSeq\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rxHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentHeight\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.LinkStats\",\"name\":\"_linkStats\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// BmcHmy is an auto generated Go binding around an Ethereum contract.
type BmcHmy struct {
	BmcHmyCaller     // Read-only binding to the contract
	BmcHmyTransactor // Write-only binding to the contract
	BmcHmyFilterer   // Log filterer for contract events
}

// BmcHmyCaller is an auto generated read-only Go binding around an Ethereum contract.
type BmcHmyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BmcHmyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BmcHmyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BmcHmyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BmcHmyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BmcHmySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BmcHmySession struct {
	Contract     *BmcHmy           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BmcHmyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BmcHmyCallerSession struct {
	Contract *BmcHmyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BmcHmyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BmcHmyTransactorSession struct {
	Contract     *BmcHmyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BmcHmyRaw is an auto generated low-level Go binding around an Ethereum contract.
type BmcHmyRaw struct {
	Contract *BmcHmy // Generic contract binding to access the raw methods on
}

// BmcHmyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BmcHmyCallerRaw struct {
	Contract *BmcHmyCaller // Generic read-only contract binding to access the raw methods on
}

// BmcHmyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BmcHmyTransactorRaw struct {
	Contract *BmcHmyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBmcHmy creates a new instance of BmcHmy, bound to a specific deployed contract.
func NewBmcHmy(address common.Address, backend bind.ContractBackend) (*BmcHmy, error) {
	contract, err := bindBmcHmy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BmcHmy{BmcHmyCaller: BmcHmyCaller{contract: contract}, BmcHmyTransactor: BmcHmyTransactor{contract: contract}, BmcHmyFilterer: BmcHmyFilterer{contract: contract}}, nil
}

// NewBmcHmyCaller creates a new read-only instance of BmcHmy, bound to a specific deployed contract.
func NewBmcHmyCaller(address common.Address, caller bind.ContractCaller) (*BmcHmyCaller, error) {
	contract, err := bindBmcHmy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BmcHmyCaller{contract: contract}, nil
}

// NewBmcHmyTransactor creates a new write-only instance of BmcHmy, bound to a specific deployed contract.
func NewBmcHmyTransactor(address common.Address, transactor bind.ContractTransactor) (*BmcHmyTransactor, error) {
	contract, err := bindBmcHmy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BmcHmyTransactor{contract: contract}, nil
}

// NewBmcHmyFilterer creates a new log filterer instance of BmcHmy, bound to a specific deployed contract.
func NewBmcHmyFilterer(address common.Address, filterer bind.ContractFilterer) (*BmcHmyFilterer, error) {
	contract, err := bindBmcHmy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BmcHmyFilterer{contract: contract}, nil
}

// bindBmcHmy binds a generic wrapper to an already deployed contract.
func bindBmcHmy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BmcHmyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BmcHmy *BmcHmyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BmcHmy.Contract.BmcHmyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BmcHmy *BmcHmyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BmcHmy.Contract.BmcHmyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BmcHmy *BmcHmyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BmcHmy.Contract.BmcHmyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BmcHmy *BmcHmyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BmcHmy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BmcHmy *BmcHmyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BmcHmy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BmcHmy *BmcHmyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BmcHmy.Contract.contract.Transact(opts, method, params...)
}

// GetBmcBtpAddress is a free data retrieval call binding the contract method 0x2a4011e9.
//
// Solidity: function getBmcBtpAddress() view returns(string)
func (_BmcHmy *BmcHmyCaller) GetBmcBtpAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BmcHmy.contract.Call(opts, &out, "getBmcBtpAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetBmcBtpAddress is a free data retrieval call binding the contract method 0x2a4011e9.
//
// Solidity: function getBmcBtpAddress() view returns(string)
func (_BmcHmy *BmcHmySession) GetBmcBtpAddress() (string, error) {
	return _BmcHmy.Contract.GetBmcBtpAddress(&_BmcHmy.CallOpts)
}

// GetBmcBtpAddress is a free data retrieval call binding the contract method 0x2a4011e9.
//
// Solidity: function getBmcBtpAddress() view returns(string)
func (_BmcHmy *BmcHmyCallerSession) GetBmcBtpAddress() (string, error) {
	return _BmcHmy.Contract.GetBmcBtpAddress(&_BmcHmy.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x22b05ed2.
//
// Solidity: function getStatus(string _link) view returns((uint256,uint256,uint256,uint256) _linkStats)
func (_BmcHmy *BmcHmyCaller) GetStatus(opts *bind.CallOpts, _link string) (TypesLinkStats, error) {
	var out []interface{}
	err := _BmcHmy.contract.Call(opts, &out, "getStatus", _link)

	if err != nil {
		return *new(TypesLinkStats), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesLinkStats)).(*TypesLinkStats)

	return out0, err

}

// GetStatus is a free data retrieval call binding the contract method 0x22b05ed2.
//
// Solidity: function getStatus(string _link) view returns((uint256,uint256,uint256,uint256) _linkStats)
func (_BmcHmy *BmcHmySession) GetStatus(_link string) (TypesLinkStats, error) {
	return _BmcHmy.Contract.GetStatus(&_BmcHmy.CallOpts, _link)
}

// GetStatus is a free data retrieval call binding the contract method 0x22b05ed2.
//
// Solidity: function getStatus(string _link) view returns((uint256,uint256,uint256,uint256) _linkStats)
func (_BmcHmy *BmcHmyCallerSession) GetStatus(_link string) (TypesLinkStats, error) {
	return _BmcHmy.Contract.GetStatus(&_BmcHmy.CallOpts, _link)
}

// TryDecodeBMCService is a free data retrieval call binding the contract method 0x2294c488.
//
// Solidity: function tryDecodeBMCService(bytes _msg) pure returns((string,bytes))
func (_BmcHmy *BmcHmyCaller) TryDecodeBMCService(opts *bind.CallOpts, _msg []byte) (TypesBMCService, error) {
	var out []interface{}
	err := _BmcHmy.contract.Call(opts, &out, "tryDecodeBMCService", _msg)

	if err != nil {
		return *new(TypesBMCService), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesBMCService)).(*TypesBMCService)

	return out0, err

}

// TryDecodeBMCService is a free data retrieval call binding the contract method 0x2294c488.
//
// Solidity: function tryDecodeBMCService(bytes _msg) pure returns((string,bytes))
func (_BmcHmy *BmcHmySession) TryDecodeBMCService(_msg []byte) (TypesBMCService, error) {
	return _BmcHmy.Contract.TryDecodeBMCService(&_BmcHmy.CallOpts, _msg)
}

// TryDecodeBMCService is a free data retrieval call binding the contract method 0x2294c488.
//
// Solidity: function tryDecodeBMCService(bytes _msg) pure returns((string,bytes))
func (_BmcHmy *BmcHmyCallerSession) TryDecodeBMCService(_msg []byte) (TypesBMCService, error) {
	return _BmcHmy.Contract.TryDecodeBMCService(&_BmcHmy.CallOpts, _msg)
}

// TryDecodeBTPMessage is a free data retrieval call binding the contract method 0x23c31a43.
//
// Solidity: function tryDecodeBTPMessage(bytes _rlp) pure returns((string,string,string,int256,bytes))
func (_BmcHmy *BmcHmyCaller) TryDecodeBTPMessage(opts *bind.CallOpts, _rlp []byte) (TypesBMCMessage, error) {
	var out []interface{}
	err := _BmcHmy.contract.Call(opts, &out, "tryDecodeBTPMessage", _rlp)

	if err != nil {
		return *new(TypesBMCMessage), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesBMCMessage)).(*TypesBMCMessage)

	return out0, err

}

// TryDecodeBTPMessage is a free data retrieval call binding the contract method 0x23c31a43.
//
// Solidity: function tryDecodeBTPMessage(bytes _rlp) pure returns((string,string,string,int256,bytes))
func (_BmcHmy *BmcHmySession) TryDecodeBTPMessage(_rlp []byte) (TypesBMCMessage, error) {
	return _BmcHmy.Contract.TryDecodeBTPMessage(&_BmcHmy.CallOpts, _rlp)
}

// TryDecodeBTPMessage is a free data retrieval call binding the contract method 0x23c31a43.
//
// Solidity: function tryDecodeBTPMessage(bytes _rlp) pure returns((string,string,string,int256,bytes))
func (_BmcHmy *BmcHmyCallerSession) TryDecodeBTPMessage(_rlp []byte) (TypesBMCMessage, error) {
	return _BmcHmy.Contract.TryDecodeBTPMessage(&_BmcHmy.CallOpts, _rlp)
}

// TryDecodeGatherFeeMessage is a free data retrieval call binding the contract method 0x9624379f.
//
// Solidity: function tryDecodeGatherFeeMessage(bytes _msg) pure returns((string,string[]))
func (_BmcHmy *BmcHmyCaller) TryDecodeGatherFeeMessage(opts *bind.CallOpts, _msg []byte) (TypesGatherFeeMessage, error) {
	var out []interface{}
	err := _BmcHmy.contract.Call(opts, &out, "tryDecodeGatherFeeMessage", _msg)

	if err != nil {
		return *new(TypesGatherFeeMessage), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesGatherFeeMessage)).(*TypesGatherFeeMessage)

	return out0, err

}

// TryDecodeGatherFeeMessage is a free data retrieval call binding the contract method 0x9624379f.
//
// Solidity: function tryDecodeGatherFeeMessage(bytes _msg) pure returns((string,string[]))
func (_BmcHmy *BmcHmySession) TryDecodeGatherFeeMessage(_msg []byte) (TypesGatherFeeMessage, error) {
	return _BmcHmy.Contract.TryDecodeGatherFeeMessage(&_BmcHmy.CallOpts, _msg)
}

// TryDecodeGatherFeeMessage is a free data retrieval call binding the contract method 0x9624379f.
//
// Solidity: function tryDecodeGatherFeeMessage(bytes _msg) pure returns((string,string[]))
func (_BmcHmy *BmcHmyCallerSession) TryDecodeGatherFeeMessage(_msg []byte) (TypesGatherFeeMessage, error) {
	return _BmcHmy.Contract.TryDecodeGatherFeeMessage(&_BmcHmy.CallOpts, _msg)
}

// HandleRelayMessage is a paid mutator transaction binding the contract method 0x21b1e9bb.
//
// Solidity: function handleRelayMessage(string _prev, bytes _msg) returns()
func (_BmcHmy *BmcHmyTransactor) HandleRelayMessage(opts *bind.TransactOpts, _prev string, _msg []byte) (*types.Transaction, error) {
	return _BmcHmy.contract.Transact(opts, "handleRelayMessage", _prev, _msg)
}

// HandleRelayMessage is a paid mutator transaction binding the contract method 0x21b1e9bb.
//
// Solidity: function handleRelayMessage(string _prev, bytes _msg) returns()
func (_BmcHmy *BmcHmySession) HandleRelayMessage(_prev string, _msg []byte) (*types.Transaction, error) {
	return _BmcHmy.Contract.HandleRelayMessage(&_BmcHmy.TransactOpts, _prev, _msg)
}

// HandleRelayMessage is a paid mutator transaction binding the contract method 0x21b1e9bb.
//
// Solidity: function handleRelayMessage(string _prev, bytes _msg) returns()
func (_BmcHmy *BmcHmyTransactorSession) HandleRelayMessage(_prev string, _msg []byte) (*types.Transaction, error) {
	return _BmcHmy.Contract.HandleRelayMessage(&_BmcHmy.TransactOpts, _prev, _msg)
}

// Initialize is a paid mutator transaction binding the contract method 0x7ab4339d.
//
// Solidity: function initialize(string _network, address _bmcManagementAddr) returns()
func (_BmcHmy *BmcHmyTransactor) Initialize(opts *bind.TransactOpts, _network string, _bmcManagementAddr common.Address) (*types.Transaction, error) {
	return _BmcHmy.contract.Transact(opts, "initialize", _network, _bmcManagementAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0x7ab4339d.
//
// Solidity: function initialize(string _network, address _bmcManagementAddr) returns()
func (_BmcHmy *BmcHmySession) Initialize(_network string, _bmcManagementAddr common.Address) (*types.Transaction, error) {
	return _BmcHmy.Contract.Initialize(&_BmcHmy.TransactOpts, _network, _bmcManagementAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0x7ab4339d.
//
// Solidity: function initialize(string _network, address _bmcManagementAddr) returns()
func (_BmcHmy *BmcHmyTransactorSession) Initialize(_network string, _bmcManagementAddr common.Address) (*types.Transaction, error) {
	return _BmcHmy.Contract.Initialize(&_BmcHmy.TransactOpts, _network, _bmcManagementAddr)
}

// SendMessage is a paid mutator transaction binding the contract method 0xbf6c1d9a.
//
// Solidity: function sendMessage(string _to, string _svc, uint256 _sn, bytes _msg) returns()
func (_BmcHmy *BmcHmyTransactor) SendMessage(opts *bind.TransactOpts, _to string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BmcHmy.contract.Transact(opts, "sendMessage", _to, _svc, _sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0xbf6c1d9a.
//
// Solidity: function sendMessage(string _to, string _svc, uint256 _sn, bytes _msg) returns()
func (_BmcHmy *BmcHmySession) SendMessage(_to string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BmcHmy.Contract.SendMessage(&_BmcHmy.TransactOpts, _to, _svc, _sn, _msg)
}

// SendMessage is a paid mutator transaction binding the contract method 0xbf6c1d9a.
//
// Solidity: function sendMessage(string _to, string _svc, uint256 _sn, bytes _msg) returns()
func (_BmcHmy *BmcHmyTransactorSession) SendMessage(_to string, _svc string, _sn *big.Int, _msg []byte) (*types.Transaction, error) {
	return _BmcHmy.Contract.SendMessage(&_BmcHmy.TransactOpts, _to, _svc, _sn, _msg)
}

// BmcHmyErrorOnBTPErrorIterator is returned from FilterErrorOnBTPError and is used to iterate over the raw logs and unpacked data for ErrorOnBTPError events raised by the BmcHmy contract.
type BmcHmyErrorOnBTPErrorIterator struct {
	Event *BmcHmyErrorOnBTPError // Event containing the contract specifics and raw log

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
func (it *BmcHmyErrorOnBTPErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BmcHmyErrorOnBTPError)
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
		it.Event = new(BmcHmyErrorOnBTPError)
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
func (it *BmcHmyErrorOnBTPErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BmcHmyErrorOnBTPErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BmcHmyErrorOnBTPError represents a ErrorOnBTPError event raised by the BmcHmy contract.
type BmcHmyErrorOnBTPError struct {
	Svc        string
	Sn         *big.Int
	Code       *big.Int
	ErrMsg     string
	SvcErrCode *big.Int
	SvcErrMsg  string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterErrorOnBTPError is a free log retrieval operation binding the contract event 0x45eab163faa71c8b113fcbc0dcc77bd39e7e3365be446895b5169bd97fc5522a.
//
// Solidity: event ErrorOnBTPError(string _svc, int256 _sn, uint256 _code, string _errMsg, uint256 _svcErrCode, string _svcErrMsg)
func (_BmcHmy *BmcHmyFilterer) FilterErrorOnBTPError(opts *bind.FilterOpts) (*BmcHmyErrorOnBTPErrorIterator, error) {

	logs, sub, err := _BmcHmy.contract.FilterLogs(opts, "ErrorOnBTPError")
	if err != nil {
		return nil, err
	}
	return &BmcHmyErrorOnBTPErrorIterator{contract: _BmcHmy.contract, event: "ErrorOnBTPError", logs: logs, sub: sub}, nil
}

// WatchErrorOnBTPError is a free log subscription operation binding the contract event 0x45eab163faa71c8b113fcbc0dcc77bd39e7e3365be446895b5169bd97fc5522a.
//
// Solidity: event ErrorOnBTPError(string _svc, int256 _sn, uint256 _code, string _errMsg, uint256 _svcErrCode, string _svcErrMsg)
func (_BmcHmy *BmcHmyFilterer) WatchErrorOnBTPError(opts *bind.WatchOpts, sink chan<- *BmcHmyErrorOnBTPError) (event.Subscription, error) {

	logs, sub, err := _BmcHmy.contract.WatchLogs(opts, "ErrorOnBTPError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BmcHmyErrorOnBTPError)
				if err := _BmcHmy.contract.UnpackLog(event, "ErrorOnBTPError", log); err != nil {
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

// ParseErrorOnBTPError is a log parse operation binding the contract event 0x45eab163faa71c8b113fcbc0dcc77bd39e7e3365be446895b5169bd97fc5522a.
//
// Solidity: event ErrorOnBTPError(string _svc, int256 _sn, uint256 _code, string _errMsg, uint256 _svcErrCode, string _svcErrMsg)
func (_BmcHmy *BmcHmyFilterer) ParseErrorOnBTPError(log types.Log) (*BmcHmyErrorOnBTPError, error) {
	event := new(BmcHmyErrorOnBTPError)
	if err := _BmcHmy.contract.UnpackLog(event, "ErrorOnBTPError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BmcHmyMessageIterator is returned from FilterMessage and is used to iterate over the raw logs and unpacked data for Message events raised by the BmcHmy contract.
type BmcHmyMessageIterator struct {
	Event *BmcHmyMessage // Event containing the contract specifics and raw log

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
func (it *BmcHmyMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BmcHmyMessage)
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
		it.Event = new(BmcHmyMessage)
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
func (it *BmcHmyMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BmcHmyMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BmcHmyMessage represents a Message event raised by the BmcHmy contract.
type BmcHmyMessage struct {
	Next string
	Seq  *big.Int
	Msg  []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMessage is a free log retrieval operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string _next, uint256 _seq, bytes _msg)
func (_BmcHmy *BmcHmyFilterer) FilterMessage(opts *bind.FilterOpts) (*BmcHmyMessageIterator, error) {

	logs, sub, err := _BmcHmy.contract.FilterLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return &BmcHmyMessageIterator{contract: _BmcHmy.contract, event: "Message", logs: logs, sub: sub}, nil
}

// WatchMessage is a free log subscription operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string _next, uint256 _seq, bytes _msg)
func (_BmcHmy *BmcHmyFilterer) WatchMessage(opts *bind.WatchOpts, sink chan<- *BmcHmyMessage) (event.Subscription, error) {

	logs, sub, err := _BmcHmy.contract.WatchLogs(opts, "Message")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BmcHmyMessage)
				if err := _BmcHmy.contract.UnpackLog(event, "Message", log); err != nil {
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

// ParseMessage is a log parse operation binding the contract event 0x37be353f216cf7e33639101fd610c542e6a0c0109173fa1c1d8b04d34edb7c1b.
//
// Solidity: event Message(string _next, uint256 _seq, bytes _msg)
func (_BmcHmy *BmcHmyFilterer) ParseMessage(log types.Log) (*BmcHmyMessage, error) {
	event := new(BmcHmyMessage)
	if err := _BmcHmy.contract.UnpackLog(event, "Message", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
