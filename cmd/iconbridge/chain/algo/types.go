package algo

import "github.com/algorand/go-algorand-sdk/types"

type Bmc struct {
	appID        uint64
	approvalProg []byte
	clearProg    []byte
	globalSchema types.StateSchema
	localSchema  types.StateSchema
}
