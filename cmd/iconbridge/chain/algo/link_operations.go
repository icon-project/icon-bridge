package algo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

const filePath = "chain/algo/linkStatus.json"

func incrementSeq(fieldName string) error {
	// Read the contents of the file into a byte slice.
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map[string]interface{}.
	var data map[string]interface{}
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return err
	}

	// Increment the specified field.
	fieldValue, ok := data[fieldName].(float64)
	if !ok {
		return fmt.Errorf("%s is not a float64", fieldName)
	}
	data[fieldName] = fieldValue + 1

	// Marshal the updated data back into a JSON string.
	updatedBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the updated JSON string back to the file.
	err = ioutil.WriteFile(filePath, updatedBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func updateHeight(fieldName string, newValue uint64) error {
	// Read the contents of the file into a byte slice.
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a map[string]interface{}.
	var data map[string]interface{}
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return err
	}

	// Replace the specified field with the new value.
	data[fieldName] = newValue

	// Marshal the updated data back into a JSON string.
	updatedBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the updated JSON string back to the file.
	err = ioutil.WriteFile(filePath, updatedBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getStatus() (*chain.BMCLinkStatus, error) {
	f, err := os.Open("chain/algo/linkStatus.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	link := &bmcLink{}
	if err := json.NewDecoder(f).Decode(&link); err != nil {
		return nil, err
	}

	return &chain.BMCLinkStatus{
		TxSeq:         link.TxSeq,
		RxSeq:         link.RxSeq,
		RxHeight:      link.RxHeight,
		CurrentHeight: link.TxHeight,
	}, nil
}
