package tests

import (
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test interface {
	Description() string
	TestDatas() []TestData
}

type TestData struct {
	Description string
	Input       interface{}
	Expected    struct {
		Success interface{}
		Fail    interface{}
	}
	MockApi     *mock.MockApi
	MockStorage mock.Storage
}

var Tests = map[string]Test{}

func RegisterTest(module string, test Test) {
	Tests[module] = test
}

func GetTest(module string, t *testing.T) (Test, error) {
	err := fmt.Errorf("not supported test:%s", module)
	if test := Tests[module]; test != nil {
		return test, nil
	}

	assert.NoError(t, err)
	return nil, err
}
