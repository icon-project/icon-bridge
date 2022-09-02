package near

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	btsp "github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc/abi/btsperiphery"
	"github.com/icon-project/icon-bridge/common/errors"
)

type parser struct {
	backend               bind.ContractBackend
	genBtsObj             *btsp.Btsperiphery
	eventIDToName         map[common.Hash]string
	addressToContractName map[string]chain.ContractName
}

func NewParser(nameToAddr map[chain.ContractName]string) (*parser, error) {
	addrToName := map[string]chain.ContractName{}
	for name, addr := range nameToAddr {
		addrToName[addr] = name
	}
	return &parser{addressToContractName: addrToName}, nil
}

func eventIDToName(abiStr string) (map[common.Hash]string, error) {
	resMap := map[common.Hash]string{}
	abi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, errors.Wrap(err, "abi.JSON ")
	}
	for _, a := range abi.Events {
		resMap[a.ID] = a.Name
	}
	return resMap, nil
}
