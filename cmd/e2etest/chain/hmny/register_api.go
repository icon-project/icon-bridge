//go:build hmny
// +build hmny

package hmny

import (
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/executor"
)

func init() {
	executor.APICallerFunc[chain.HMNY] = NewApi
}
