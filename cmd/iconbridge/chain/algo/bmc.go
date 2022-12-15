/* package algo

func (s *sender) HandleRelayMessage(ctx context.Context, _prev []byte, _msg []byte) (
	*types.Transaction, error) {

	s.callAbi("HandleRelayMessage", []interface{[]struct{_prev, _msg ,}})

	var dummyPk ed25519.PrivateKey

	txID, signedTxn, _ := crypto.SignTransaction(dummyPk, txn)

	txID, err = cl.algod.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return nil, "", err
	}

	return signedTxn, txID, nil
}

*/