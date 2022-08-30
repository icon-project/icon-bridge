package cmd

import (
	"encoding/json"
	"os"
)

type Config struct {
	Network    string `json:"name"`
	KeyStore   string `json:"key_path"`
	KeySecret  string `json:"key_secret"`
	Sender     string `json:"sender"`
	Reciever   string `json:"reciever"`
	BtsAddress string `json:"bts"`
	StepLimit  int64  `json:"step_limit,omitempty"`
	NetworkId  string `json:"network_id"`
	Uri        string `json:"uri"`
	Value      string `json:"value"`
}

func NewConfig(filePath string) ([]Config, error) {

	var configs []Config
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contents, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil

}
