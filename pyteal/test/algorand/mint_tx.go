package algorand

import (
	"encoding/base64"

	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

func MintTx(
	txParams types.SuggestedParams,
	creator types.Address,
	totalIssuance uint64,
	decimals uint32,
	unitName string,
	assetName string,
	assetURL string,
	assetMetadataHash string,
) (tx types.Transaction, err error) {
	
	defaultFrozen := false
	manager := ""
	freeze := ""
	clawback := ""
	note := []byte(nil)

	return transaction.MakeAssetCreateTxn(
		creator.String(),
		uint64(txParams.Fee),
		uint64(txParams.FirstRoundValid),
		uint64(txParams.LastRoundValid),
		note,
		txParams.GenesisID,
		base64.StdEncoding.EncodeToString(txParams.GenesisHash),
		totalIssuance,
		decimals,
		defaultFrozen,
		manager,
		creator.String(), // reserve
		freeze,
		clawback,
		unitName,
		assetName,
		assetURL,
		assetMetadataHash,
	)
}
