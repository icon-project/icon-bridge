package algorand

import "github.com/algorand/go-algorand-sdk/types"

// Check if the new block has any transaction meant to be sent across the relayer
func GetTxns(block *types.Block, appId uint64) *[]types.SignedTxnWithAD {
	if len(block.Payset) == 0 {
		return nil
	}
	txns := make([]types.SignedTxnWithAD, 0)
	for _, signedTxnInBlock := range block.Payset {
		signedTxnWithAD := signedTxnInBlock.SignedTxnWithAD
		//TODO review the way of properly identify a bmc txn
		if signedTxnWithAD.SignedTxn.Txn.Type == types.ApplicationCallTx &&
			signedTxnWithAD.SignedTxn.Txn.ApplicationID == types.AppIndex(appId) {

			txns = append(txns, signedTxnWithAD)
		}
	}
	return &txns
}
