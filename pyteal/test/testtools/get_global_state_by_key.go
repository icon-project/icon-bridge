package testtools

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)

func GetGlobalStateByKey(t *testing.T,client *algod.Client, appId uint64, key string) []byte {
	app, err := client.GetApplicationByID(appId).Do(context.Background())

	if err != nil {
		t.Fatalf("Failed to get application with id: %v. details: %+v", appId, err)
	}

	encodedKey := base64.StdEncoding.EncodeToString([]byte(key))

	globalState := app.Params.GlobalState

	for _, modelTealKeyValue := range globalState {
		if modelTealKeyValue.Key == encodedKey {
			decodedValue, err := base64.StdEncoding.DecodeString(modelTealKeyValue.Value.Bytes)

			if err != nil {
				t.Fatalf("Failed to decode string to base64: %+v", err)
			}

			return decodedValue
		}
	}

	t.Fatalf("Failed to find %s in global state of application by id: %v", key, appId)

	return nil
}