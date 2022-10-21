package substrate_eth

/*
const (
	ICON_BMC          = "btp://0x7.icon/cx8a6606d526b96a16e6764aee5d9abecf926689df"
	BSC_BMC_PERIPHERY = "btp://0x61.bsc/0xB4fC4b3b4e3157448B7D279f06BC8e340d63e2a9"
	BlockHeight       = 21447824
)

func newTestReceiver(t *testing.T, src, dst chain.BTPAddress) chain.Receiver {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	mp := map[string]interface{}{"syncConcurrency": 2}
	res, err := json.Marshal(mp)
	require.NoError(t, err)
	receiver, err := NewReceiver(src, dst, []string{url}, res, log.New())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return receiver
}

func newTestClient(t *testing.T, bmcAddr string) IClient {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	cls, _, err := newClients([]string{url}, bmcAddr, log.New())
	require.NoError(t, err)
	return cls[0]
}

func TestFilterLogs(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set(BSC_BMC_PERIPHERY)
	require.NoError(t, err)
	err = dst.Set(ICON_BMC)
	require.NoError(t, err)

	recv := newTestReceiver(t, src, dst).(*receiver)
	if recv == nil {
		t.Fatal(errors.New("Receiver is nil"))
	}
	exists, err := recv.hasBTPMessage(context.Background(), big.NewInt(BlockHeight))
	require.NoError(t, err)
	if !exists {
		require.NoError(t, errors.New("Expected true"))
	}
}

func TestSubscribeMessage(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set(BSC_BMC_PERIPHERY)
	err = dst.Set(ICON_BMC)
	if err != nil {
		fmt.Println(err)
	}

	recv := newTestReceiver(t, src, dst).(*receiver)

	ctx, cancel := context.Background(), func() {}
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()
	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := recv.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    75,
			Height: uint64(BlockHeight),
		})
	require.NoError(t, err, "failed to subscribe")

	for {
		defer cancel()
		select {
		case err := <-srcErrCh:
			t.Logf("subscription closed: %v", err)
			t.FailNow()
		case msg := <-srcMsgCh:
			if len(msg.Receipts) > 0 && msg.Receipts[0].Height == 21447824 {
				// received event exit
				return
			}
		}
	}
}

func TestReceiver_GetReceiptProofs(t *testing.T) {
	cl := newTestClient(t, BSC_BMC_PERIPHERY)
	header, err := cl.GetHeaderByHeight(big.NewInt(BlockHeight))
	require.NoError(t, err)
	//hash := header.Hash()
	receipts, err := cl.GetBlockReceiptsFromHeight(big.NewInt(BlockHeight))
	require.NoError(t, err)
	receiptsRoot := ethTypes.DeriveSha(receipts, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptsRoot.Bytes(), header.ReceiptHash.Bytes()) {
		err = fmt.Errorf(
			"invalid receipts: remote=%v, local=%v",
			header.ReceiptHash, receiptsRoot)
		require.NoError(t, err)
	}
}

func TestReceiver_MockReceiverOptions_UnmarshalWithoutVerifier(t *testing.T) {
	// Verifier should be nil if not passed
	var empty_opts ReceiverOptions
	jsonReceiverOptions := `{"syncConcurrency":100}`
	json.Unmarshal([]byte(jsonReceiverOptions), &empty_opts)
	require.NotNil(t, empty_opts)
	require.Nil(t, empty_opts.Verifier)
	require.NotNil(t, empty_opts.SyncConcurrency)
	require.EqualValues(t, 100, empty_opts.SyncConcurrency)
}

func TestSender_NewObj(t *testing.T) {
	//senderOpts := `{"gas_limit": 24000000,"tx_data_size_limit": 8192,"balance_threshold": "100000000000000000000","boost_gas_price": 1}`
	thres := intconv.BigInt{}
	thres.SetString("100000000000000000000", 10)
	sopts := senderOptions{
		GasLimit:         24000000,
		TxDataSizeLimit:  8192,
		BalanceThreshold: thres,
		BoostGasPrice:    1,
	}
	raw, err := json.Marshal(sopts)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	s, err := NewSender(
		chain.BTPAddress(BSC_BMC_PERIPHERY),
		chain.BTPAddress(ICON_BMC),
		[]string{url}, &wallet.EvmWallet{Skey: privKey, Pkey: &privKey.PublicKey},
		raw,
		log.New(),
	)
	balance, threshold, err := s.Balance(context.TODO())
	require.NoError(t, err)
	require.Equal(t, balance.Cmp(big.NewInt(0)), 0)
	require.Equal(t, threshold.String(), thres.String())

	msg := &chain.Message{
		From: "",
		Receipts: []*chain.Receipt{{
			Index:  0,
			Height: 1,
			Events: []*chain.Event{},
		}},
	}
	tx, _, err := s.Segment(context.TODO(), msg)
	require.NoError(t, err)
	err = tx.Send(context.TODO())
	require.Equal(t, err.Error(), "InsufficientBalance")
}
*/
