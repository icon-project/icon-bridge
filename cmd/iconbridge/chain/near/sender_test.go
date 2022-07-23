package near

import (
	"context"
	"crypto/ed25519"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/assert"
)

func TestNearSender(t *testing.T) {
	if test, err := tests.GetTest("GetBmcLinkStatus", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					mockApi := NewMockApi(testData.MockStorage)
					client := &Client{
						api: &mockApi,
					}

					links, Ok := (testData.Input).([]chain.BTPAddress)
					assert.True(f, Ok)

					sender, err := newMockSender(links[1], links[0], client, nil, nil, nil)
					assert.Nil(f, err)
					
					status, err := sender.Status(context.Background())
					assert.Nil(f, err)
					assert.NotNil(f, status)
					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(int)
						assert.True(f, Ok)
						assert.Equal(f, uint64(expected), status.RotateHeight)
					} else {
						assert.Error(f, err)
					}

				})
			}
		})
	}

	if test, err := tests.GetTest("RelayTxSend", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					mockApi := NewMockApi(testData.MockStorage)
					client := &Client{
						api: &mockApi,
						logger: log.New(),
					}

					links, Ok := (testData.Input).([]chain.BTPAddress)
					assert.True(f, Ok)

					privateKeyBytes := base58.Decode("22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J")
					privateKey := ed25519.PrivateKey(privateKeyBytes)
					nearWallet, err := wallet.NewNearwalletFromPrivateKey(&privateKey)
			
					assert.NoError(f, err)
					sender, err := newMockSender(links[1], links[0], client, nearWallet, nil, log.New())
					assert.Nil(f, err)

					relayTx, err := sender.newRelayTransaction(context.Background(), "", []byte{})
					assert.Nil(f, err)

					err = relayTx.Send(context.Background())
					assert.Nil(f, err)

					assert.NotNil(f, relayTx.Transaction.Txid)

					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(string)
						assert.True(f, Ok)
						assert.Equal(f, expected, relayTx.Transaction.Txid.Base58Encode())
					} else {
						assert.Error(f, err)
					}
				})
			}
		})
	}

	if test, err := tests.GetTest("RelayTxReceipt", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					mockApi := NewMockApi(testData.MockStorage)
					client := &Client{
						api: &mockApi,
						logger: log.New(),
					}

					links, Ok := (testData.Input).([]chain.BTPAddress)
					assert.True(f, Ok)

					privateKeyBytes := base58.Decode("22yx6AjQgG1jGuAmPuEwLnVKFnuq5LU23dbU3JBZodKxrJ8dmmqpDZKtRSfiU4F8UQmv1RiZSrjWhQMQC3ye7M1J")
					privateKey := ed25519.PrivateKey(privateKeyBytes)
					nearWallet, err := wallet.NewNearwalletFromPrivateKey(&privateKey)
			
					assert.NoError(f, err)
					sender, err := newMockSender(links[1], links[0], client, nearWallet, nil, log.New())
					assert.Nil(f, err)

					relayTx, err := sender.newRelayTransaction(context.Background(), "", []byte{})
					assert.Nil(f, err)

					err = relayTx.Send(context.Background())
					assert.Nil(f, err)

					assert.NotNil(f, relayTx.Transaction.Txid)

					blockHeight, err := relayTx.Receipt(context.Background())
					assert.Nil(f, err)

					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(int)
						assert.True(f, Ok)
						assert.Equal(f, uint64(expected), blockHeight)
					} else {
						assert.Error(f, err)
					}
				})
			}
		})
	}
}
