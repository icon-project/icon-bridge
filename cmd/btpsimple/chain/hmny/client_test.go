package hmny

// func newTestClient(t *testing.T) *relayClient {
// 	url := "https://rpc.s0.b.hmny.io"
// 	cl, err := newRelayClient([]string{url}, nil, "", log.New())
// 	require.NoError(t, err)
// 	return cl
// }

// func getDefaultContext() (context.Context, context.CancelFunc) {
// 	return context.WithTimeout(context.Background(), time.Minute)
// }

// func TestGetTransactionReceipt(t *testing.T) {
// 	cl := newTestClient(t)
// 	txh := common.HexToHash("0x04c3009eb637b8871cfc3732bfe6c23bca1b6e850a6e8bb47dd32ac521d7af7b")
// 	// txh := common.HexToHash("0xa3cb6a6f7530d8541f7f6a17d5244e3ff9a6c98b61d5ec21a8b0c12eebb89809")

// 	ctx, cancel := getDefaultContext()
// 	defer cancel()
// 	tx, _, err := cl.client().eth.TransactionByHash(ctx, txh)
// 	require.NoError(t, err)

// 	ctx, cancel = getDefaultContext()
// 	defer cancel()
// 	txr, err := cl.client().eth.TransactionReceipt(ctx, txh)
// 	require.NoError(t, err)

// 	if txr.Status == 0 {
// 		callMsg := ethereum.CallMsg{
// 			From:       common.HexToAddress("0x5f7043477705a4b4a5cb612c76715aec35c26afc"),
// 			To:         tx.To(),
// 			Gas:        tx.Gas(),
// 			GasPrice:   tx.GasPrice(),
// 			Value:      tx.Value(),
// 			AccessList: tx.AccessList(),
// 			Data:       tx.Data(),
// 		}

// 		ctx, cancel = getDefaultContext()
// 		defer cancel()
// 		data, err := cl.client().eth.CallContract(ctx, callMsg, txr.BlockNumber)
// 		require.NoError(t, err)

// 		fmt.Println(revertReason(data))
// 		/*
// 			08c379a0
// 			0000000000000000000000000000000000000000000000000000000000000020
// 			000000000000000000000000000000000000000000000000000000000000002b
// 			526576657274496e76616c696452785365713a2065762e736571203e206578706563746564207278536571000000000000000000000000000000000000000000
// 		*/
// 	}
// }
