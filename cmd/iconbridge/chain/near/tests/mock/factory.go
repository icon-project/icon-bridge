package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

const mockDataPath = "./tests/mock/data"

func loadFiles(files []string, directory string) [][]byte {
	var fileBuffers = make([][]byte, 0)
	for _, f := range files {
		file, err := os.Open(directory + "/" + f + ".json")

		if err != nil {
			panic(fmt.Errorf("error [LoadFile]: %v", err))
		}

		buffer, err := ioutil.ReadAll(file)

		if err != nil {
			panic(fmt.Errorf("error [ReadFile]: %v", err))
		}

		fileBuffers = append(fileBuffers, buffer)
		defer file.Close()
	}
	return fileBuffers
}

func validateDirectory(directory string) {
	_, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(fmt.Errorf("error [ValidateDirectory]: %v", err))
	}
}

func LoadBlockFromFile(names []string) (map[int64]Response, map[string]Response) {
	sectionDir := mockDataPath + "/blocks"
	validateDirectory(sectionDir)

	blockByHashMap := map[string]Response{}
	blockByHeightMap := map[int64]Response{}

	for _, buffer := range loadFiles(names, sectionDir) {
		var block types.Block
		err := json.Unmarshal(buffer, &block)
		if err != nil {
			panic(fmt.Errorf("error [LoadBlock][ParseJson]: %v", err))
		}

		blockByHashMap[block.Header.Hash.Base58Encode()] = Response{
			Reponse: buffer,
			Error:   nil,
		}

		blockByHeightMap[block.Header.Height] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}

	return blockByHeightMap, blockByHashMap
}

func LoadAccessKeyFromFile(names []string) map[string]Response {
	sectionDir := mockDataPath + "/access_key"
	validateDirectory(sectionDir)

	var accessKeyMap = map[string]Response{}

	for index, buffer := range loadFiles(names, sectionDir) {
		accessKeyMap[names[index]] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}

	return accessKeyMap
}

func LoadAccountsFromFile(accounts []string) map[string]Response {
	sectionDir := mockDataPath + "/accounts"
	validateDirectory(sectionDir)

	var accountMap = map[string]Response{}

	for index, buffer := range loadFiles(accounts, sectionDir) {
		accountMap[accounts[index]] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}

	return accountMap
}

func LoadBmcStatusFromFile(names []string) map[string]Response {
	sectionDir := mockDataPath + "/contractsdata/bmc"
	validateDirectory(sectionDir)

	var bmcStatusMap = map[string]Response{}

	for index, buffer := range loadFiles(names, sectionDir) {
		var bmcstatus types.BmcStatus
		err := json.Unmarshal(buffer, &bmcstatus)
		if err != nil {
			panic(fmt.Errorf("error [LoadBlock][ParseJson]: %v", err))
		}

		bmcStatusMap[names[index]] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}

	return bmcStatusMap
}

func LoadEventsFromFile(names []string) map[int64]Response {
	sectionDir := mockDataPath + "/events"
	validateDirectory(sectionDir)

	var getEventsMap = map[int64]Response{}
	for index, buffer := range loadFiles(names, sectionDir) {
		var contractStateChange types.ContractStateChange

		err := json.Unmarshal(buffer, &contractStateChange)
		if err != nil {
			panic(fmt.Errorf("error [LoadEvents][ParseJson]: %v", err))
		}

		blockHeight, err := strconv.Atoi(names[index])
		if err != nil {
			panic(fmt.Errorf("error [LoadEvents][ParseBlockHeight]: %v", err))
		}
		getEventsMap[int64(blockHeight)] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}
	return getEventsMap
}

func LoadReceiptsFromFile(names []string) map[string]Response {
	sectionDir := mockDataPath + "/receipts"
	validateDirectory(sectionDir)
	var receiptProofMap = map[string]Response{}
	for index, buffer := range loadFiles(names, sectionDir) {
		var receiptProofs types.ReceiptProof

		err := json.Unmarshal(buffer, &receiptProofs)
		if err != nil {
			panic(fmt.Errorf("error [LoadReceipts][ParseJson]: %v", err))
		}

		receiptProofMap[names[index]] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}
	return receiptProofMap
}

func LoadTransactionResultFromFile(names []string) map[string]Response {
	sectionDir := mockDataPath + "/transaction"
	validateDirectory(sectionDir)

	var transactionResultMap = map[string]Response{}
	for index, buffer := range loadFiles(names, sectionDir) {
		var transactionResult types.TransactionResult

		err := json.Unmarshal(buffer, &transactionResult)
		if err != nil {
			panic(fmt.Errorf("error [LoadTransactionResult][ParseJson]: %v", err))
		}

		transactionResultMap[names[index]] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}
	return transactionResultMap
}

func LoadChainStatusFromFile(blocks []string) Response {
	sectionDir := mockDataPath + "/status"
	validateDirectory(sectionDir)

	var LatestChainStatus = Response{}
	for _, buffer := range loadFiles(blocks, sectionDir) {
		var ChainStatus types.ChainStatus

		err := json.Unmarshal(buffer, &ChainStatus)
		if err != nil {
			panic(fmt.Errorf("error [LoadChainStatus][ParseJson]: %v", err))
		}

		LatestChainStatus = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}

	return LatestChainStatus
}

func LoadBlockProducersFromFile(epochs []string) map[string]Response {
	sectionDir := mockDataPath + "/block_producers"
	validateDirectory(sectionDir)

	var blockProducersMap = map[string]Response{}
	for index, buffer := range loadFiles(epochs, sectionDir) {
		var bps types.BlockProducers

		err := json.Unmarshal(buffer, &bps)
		if err != nil {
			panic(fmt.Errorf("error [LoadBlockProducers][ParseJson]: %v", err))
		}

		blockProducersMap[epochs[index]] = Response{
			Reponse: buffer,
			Error:   nil,
		}
	}

	return blockProducersMap
}
