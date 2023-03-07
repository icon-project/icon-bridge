package near

import (
	"context"
	"crypto/ed25519"
	"math/big"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearSender(t *testing.T) {
	if test, err := tests.GetTest("GetBmcLinkStatus", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api: testData.MockApi,
					}

					input, Ok := (testData.Input).(struct {
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					assert.True(f, Ok)

					sender, err := NewSender(SenderConfig{source: input.Source, destination: input.Destination}, log.New(), client)
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
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						PrivateKey  string
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					assert.True(f, Ok)

					privateKeyBytes := base58.Decode(input.PrivateKey)
					privateKey := ed25519.PrivateKey(privateKeyBytes)
					wallet, err := wallet.NewNearwalletFromPrivateKey(&privateKey)

					assert.NoError(f, err)
					sender, err := NewSender(SenderConfig{source: input.Source, destination: input.Destination, wallet: wallet}, log.New(), client)
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
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						PrivateKey  string
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					assert.True(f, Ok)

					privateKeyBytes := base58.Decode(input.PrivateKey)
					privateKey := ed25519.PrivateKey(privateKeyBytes)
					wallet, err := wallet.NewNearwalletFromPrivateKey(&privateKey)

					assert.NoError(f, err)
					sender, err := NewSender(SenderConfig{source: input.Source, destination: input.Destination, wallet: wallet}, log.New(), client)
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

	if test, err := tests.GetTest("GetBalance", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						PrivateKey  string
						Source      chain.BTPAddress
						Destination chain.BTPAddress
					})
					assert.True(f, Ok)

					privateKeyBytes := base58.Decode(input.PrivateKey)
					privateKey := ed25519.PrivateKey(privateKeyBytes)
					wallet, err := wallet.NewNearwalletFromPrivateKey(&privateKey)

					assert.NoError(f, err)
					sender, err := NewSender(SenderConfig{source: input.Source, destination: input.Destination, wallet: wallet}, log.New(), client)
					assert.Nil(f, err)

					balance, _, err := sender.Balance(context.Background())
					assert.Nil(f, err)

					if testData.Expected.Success != nil {
						expected, Ok := (testData.Expected.Success).(*big.Int)
						assert.True(f, Ok)
						assert.Equal(f, expected, balance)
					} else {
						assert.Error(f, err)
					}
				})
			}
		})
	}

	if test, err := tests.GetTest("SenderSegment", t); err == nil {
		t.Run(test.Description(), func(f *testing.T) {
			for _, testData := range test.TestDatas() {
				f.Run(testData.Description, func(f *testing.T) {
					client := &Client{
						api:    testData.MockApi,
						logger: log.New(),
					}

					input, Ok := (testData.Input).(struct {
						PrivateKey  string
						Source      chain.BTPAddress
						Destination chain.BTPAddress
						Message     *chain.Message
						Options     types.SenderOptions
					})
					require.True(f, Ok)

					privateKeyBytes := base58.Decode(input.PrivateKey)
					privateKey := ed25519.PrivateKey(privateKeyBytes)
					nearWallet, err := wallet.NewNearwalletFromPrivateKey(&privateKey)
					require.NoError(f, err)

					sender, err := NewSender(SenderConfig{source: input.Source, destination: input.Destination, wallet: nearWallet, options: input.Options}, log.New(), client)
					require.NoError(f, err)

					relayTx, newMsg, err := sender.Segment(context.Background(), input.Message)
					require.NoError(f, err)

					if testData.Expected.Success != nil {
						message := make([]byte, 0)

						if relayTx != nil {
							message = (relayTx).(*RelayTransaction).message
						}

						testData.Expected.Success.(func(*testing.T, []byte, *chain.Message))(f, message, newMsg)
						assert.Nil(f, err)
					} else {
						testData.Expected.Fail.(func(*testing.T, error))(f, err)
					}
				})
			}
		})
	}
}
