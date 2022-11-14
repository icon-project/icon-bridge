package mock

var LatestChainStatus = []string{"377840"}
var Blocks = []string{"377825", "377943"}
var Nonce = []string{"69c003c3b80ed12ea02f5c67c9e8167f0ce3b2e8020a0f43b1029c4d787b0d21", "94a5a3fc9bc948a7f4b1c6210518b4afe1744ebe33188eb91d17c863dfe200a8"}
var BmcLinkStatus = []string{"c294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5"}
var BmvLinkStatus = []string{"dev-20211205172325-28827597417784"}
var GetEvents = []string{"377825", "377826", "377827", "377828", "377829"}
var GetReceiptsProof = []string{"2VWWEfyg5BzyDBVRsHRXApQJdZ37Bdtj8GtkH7UvNm7G"}
var TransactionResult = []string{"6zgh2u9DqHHiXzdy9ouTP7oGky2T4nugqzqt9wJZwNFm"}
var accounts = []string{"69c003c3b80ed12ea02f5c67c9e8167f0ce3b2e8020a0f43b1029c4d787b0d21", "94a5a3fc9bc948a7f4b1c6210518b4afe1744ebe33188eb91d17c863dfe200a8"}
var bps = []string{"84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"}

func Default() Storage {
	latestChainStatus := LoadChainStatusFromFile(LatestChainStatus)
	blockByHeightMap, blockByHashMap := LoadBlockFromFile(Blocks)
	accessKeyMap := LoadAccessKeyFromFile(Nonce)
	bmcLinkStatusMap := LoadBmcStatusFromFile(BmcLinkStatus)
	contractStateChangeMap := LoadEventsFromFile(GetEvents)
	receiptProofMap := LoadReceiptsFromFile(GetReceiptsProof)
	transactionHashMap := Response{
		Reponse: "6zgh2u9DqHHiXzdy9ouTP7oGky2T4nugqzqt9wJZwNFm",
		Error:   nil,
	}
	transactionResultMap := LoadTransactionResultFromFile(TransactionResult)
	accountMap := LoadAccountsFromFile(accounts)
	blockProducersMap := LoadBlockProducersFromFile(bps)

	return Storage{
		LatestChainStatus:      latestChainStatus,
		BlockByHeightMap:       blockByHeightMap,
		BlockByHashMap:         blockByHashMap,
		AccessKeyMap:           accessKeyMap,
		BmcLinkStatusMap:       bmcLinkStatusMap,
		ContractStateChangeMap: contractStateChangeMap,
		ReceiptProofMap:        receiptProofMap,
		TransactionHash:        transactionHashMap,
		TransactionResultMap:   transactionResultMap,
		AccountMap:             accountMap,
		BlockProducersMap:      blockProducersMap,
	}
}
