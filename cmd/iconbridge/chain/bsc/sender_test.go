package bsc

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	dst_options_json =`{"step_limit":24, "gas_limit":24000000,"tx_data_size_limit":8192,"balance_threshold":"100000000000000000000","boost_gas_price":1.0}`
)

func Test_unmarshalOpt(t *testing.T) {
	var senderOpt senderOptions

	err := unmarshalOpt([]byte(dst_options_json), &senderOpt)

	if err != nil {
		t.Errorf("unmarshalOpt() error = %v", err)
	}

	assert.EqualValues(t, 24000000, senderOpt.GasLimit)
	assert.EqualValues(t, 1.0, senderOpt.BoostGasPrice)
	assert.EqualValues(t, 8192, senderOpt.TxDataSizeLimit)
	threshold := new(big.Int)
	threshold.SetString("100000000000000000000",10)
	assert.EqualValues(t, 0, threshold.Cmp(&senderOpt.BalanceThreshold))
}