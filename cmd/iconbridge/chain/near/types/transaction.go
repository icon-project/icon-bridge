package types

import (
	"crypto/sha256"
	"fmt"

	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/near/borsh-go"
)

type TransactionResult struct {
	Status             ExecutionStatus              `json:"status"`
	Transaction        Transaction                  `json:"transaction"`
	TransactionOutcome ExecutionOutcomeWithIdView   `json:"transaction_outcome"`
	ReceiptsOutcome    []ExecutionOutcomeWithIdView `json:"receipts_outcome"`
}

type Transaction struct {
	SignerId   AccountId `json:"signer_id"`
	PublicKey  PublicKey `json:"public_key"`
	Nonce      int       `json:"nonce"`
	ReceiverId AccountId `json:"receiver_id"`
	BlockHash  CryptoHash
	Actions    []Action   `json:"actions"` // TODO: ActionView
	Signature  Signature  `json:"signature"`
	Txid       CryptoHash `json:"hash"`
}

func (t *Transaction) serialize() ([]byte, error) {
	return borsh.Serialize(struct {
		SignerId   AccountId `json:"signer_id"`
		PublicKey  PublicKey `json:"public_key"`
		Nonce      int       `json:"nonce"`
		ReceiverId AccountId `json:"receiver_id"`
		BlockHash  CryptoHash
		Actions    []Action `json:"actions"` // TODO: ActionView
	}{
		SignerId:   t.SignerId,
		PublicKey:  t.PublicKey,
		Nonce:      t.Nonce,
		ReceiverId: t.ReceiverId,
		BlockHash:  t.BlockHash,
		Actions:    t.Actions,
	})
}

func (t *Transaction) Sign(wallet *wallet.NearWallet) error {
	serializedTransaction, err := t.serialize()
	if err != nil {
		return err
	}
	preSigndata := sha256.Sum256(serializedTransaction)

	signature, err := wallet.Sign(preSigndata[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction")
	}

	if len(signature) != 64 {
		return fmt.Errorf("signature error,length is not equal 64, length=%d", len(signature))
	}

	t.Signature = Signature{
		KeyType: ED25519,
		Data:    signature,
	}
	return nil
}

type RelayMessageParam struct {
	Previous string `json:"_prev"`
	Message  string `json:"_msg"`
}

type TransactionParam struct {
	From              string
	To                string
	RelayMessage      RelayMessageParam
	Base64encodedData string
}
