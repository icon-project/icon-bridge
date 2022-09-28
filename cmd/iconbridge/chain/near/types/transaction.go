package types

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/near/borsh-go"
)

type TransactionResult struct {
	Status             ExecutionStatus              `json:"status"`
	Transaction        TransactionView              `json:"transaction"`
	TransactionOutcome ExecutionOutcomeWithIdView   `json:"transaction_outcome"`
	ReceiptsOutcome    []ExecutionOutcomeWithIdView `json:"receipts_outcome"`
}

type TransactionView struct {
	SignerId   AccountId    `json:"signer_id"`
	PublicKey  PublicKey    `json:"public_key"`
	Nonce      int          `json:"nonce"`
	ReceiverId AccountId    `json:"receiver_id"`
	Actions    []ActionView `json:"actions"` // TODO: ActionView
	Signature  Signature    `json:"signature"`
	Txid       CryptoHash   `json:"hash"`
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

func (t *Transaction) Payload(wallet *wallet.NearWallet) (string, error) {
	if err := t.sign(wallet); err != nil {
		return "", err
	}

	serializedSignedTransaction, err := borsh.Serialize(struct {
		Transaction struct {
			SignerId   AccountId
			PublicKey  PublicKey
			Nonce      uint64
			ReceiverId AccountId
			BlockHash  [32]byte
			Actions    []Action
		}
		Signature Signature
	}{
		Transaction: struct {
			SignerId   AccountId
			PublicKey  PublicKey
			Nonce      uint64
			ReceiverId AccountId
			BlockHash  [32]byte
			Actions    []Action
		}{
			SignerId:   t.SignerId,
			PublicKey:  t.PublicKey,
			Nonce:      uint64(t.Nonce),
			ReceiverId: t.ReceiverId,
			BlockHash:  t.BlockHash,
			Actions:    t.Actions,
		},
		Signature: t.Signature,
	})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(serializedSignedTransaction[:]), nil
}

func (t *Transaction) sign(wallet *wallet.NearWallet) error {
	serializedTransaction, err := borsh.Serialize(struct {
		SignerId   AccountId
		PublicKey  PublicKey
		Nonce      uint64
		ReceiverId AccountId
		BlockHash  CryptoHash
		Actions    []Action
	}{
		SignerId:   t.SignerId,
		PublicKey:  t.PublicKey,
		Nonce:      uint64(t.Nonce),
		ReceiverId: t.ReceiverId,
		BlockHash:  t.BlockHash,
		Actions:    t.Actions,
	})
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
	}

	copy(t.Signature.Data[:], signature)

	return nil
}
