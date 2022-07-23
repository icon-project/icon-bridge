package tests

import "github.com/icon-project/icon-bridge/cmd/iconbridge/chain"

type RelayTxSendTest struct {
	description string
	testData    []TestData
}

func (t RelayTxSendTest) Description() string {
	return t.description
}

func (t RelayTxSendTest) TestDatas() []TestData {
	return t.testData
}

func init() {
	source := chain.BTPAddress("btp://0x1.icon/0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5")
	destination := chain.BTPAddress("btp://0x1.near/dev-20211206025826-24100687319598")
	var testData = []TestData{
		{
			Description: "RelayTx Send Success",
			Input:       []chain.BTPAddress{destination, source},
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: "6zgh2u9DqHHiXzdy9ouTP7oGky2T4nugqzqt9wJZwNFm",
			},
		},
	}

	RegisterTest("RelayTxSend", RelayTxSendTest{
		description: "Test RelayTx Send",
		testData:    testData,
	})
}
