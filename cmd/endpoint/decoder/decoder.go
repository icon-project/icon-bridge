package decoder

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ctr "github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/bshImpl"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/bshPeriphery"
)

// Update this function for more contracts
func getNewContract(cName ctr.ContractName, url string, cAddr common.Address) (ctr.Contract, error) {
	if cName == ctr.BSHImpl {
		return bshImpl.NewContract(url, cAddr)
	} else if cName == ctr.BSHPeriphery {
		return bshPeriphery.NewContract(url, cAddr)
	}
	return nil, errors.New("Contract not registered")
}

type Decoder interface {
	Add(contractNameToAddressMap map[ctr.ContractName]common.Address) (err error)
	Remove(addr common.Address)
	DecodeEventLogData(log types.Log) (map[string]interface{}, error)
}

type decoder struct {
	url            string
	mtx            sync.RWMutex
	addrToContract map[common.Address]ctr.Contract
}

func New(url string, contractNameToAddressMap map[ctr.ContractName]common.Address) (Decoder, error) {
	var err error
	dec := &decoder{mtx: sync.RWMutex{}, url: url, addrToContract: make(map[common.Address]ctr.Contract)}
	for cName, cAddr := range contractNameToAddressMap {
		dec.addrToContract[cAddr], err = getNewContract(cName, url, cAddr)
		if err != nil {
			return nil, err
		}
	}
	return dec, nil
}

func (d *decoder) Add(contractNameToAddressMap map[ctr.ContractName]common.Address) (err error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	for cName, cAddr := range contractNameToAddressMap {
		if cn, ok := d.addrToContract[cAddr]; ok || cn.GetName() == cName {
			continue // Name or address already exists
		}
		d.addrToContract[cAddr], err = getNewContract(cName, d.url, cAddr)
		if err != nil {
			return
		}
	}
	return nil
}

func (d *decoder) Remove(addr common.Address) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	delete(d.addrToContract, addr)
}

func (d *decoder) DecodeEventLogData(log types.Log) (map[string]interface{}, error) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	cAddr := d.addrToContract[log.Address]
	return cAddr.Decode(log)
}
