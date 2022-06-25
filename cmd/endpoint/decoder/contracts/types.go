package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/harmony/accounts/abi"
)

type Contract interface {
	Decode(log interface{}) (map[string]interface{}, error)
	GetName() ContractName
}

type ContractName string

const (
	TokenHmy   ContractName = "TokenHmy"
	NativeHmy  ContractName = "NativeHmy"
	TokenIcon  ContractName = "TokenIcon"
	NativeIcon ContractName = "NativeIcon"
)

// For Hmy only
func EventIDToName(abiStr string) (map[common.Hash]string, error) {
	resMap := map[common.Hash]string{}
	abi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}
	for _, a := range abi.Events {
		resMap[a.ID] = a.Name
	}
	return resMap, nil
}
