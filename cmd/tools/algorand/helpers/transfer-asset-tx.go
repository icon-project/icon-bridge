package helpers

import (
	"encoding/base64"

	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

func TransferAssetTx(
	txParams types.SuggestedParams,
	from types.Address,
	to types.Address,
	assetID uint64,
	amount uint64,
) (types.Transaction, error) {
	
	return transaction.MakeAssetTransferTxn(
		from.String(),
		to.String(),
		"",
		amount,
		uint64(txParams.Fee),
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		[]byte(nil),
		txParams.GenesisID,
		base64.StdEncoding.EncodeToString(txParams.GenesisHash),
		assetID,
	)
}