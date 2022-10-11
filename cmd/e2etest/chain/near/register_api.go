package near

import (
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
)

func init() {
	executor.APICallerFunc[chain.NEAR] = NewApi
}
