package mock

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

func (m *MockApi) BlockFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var block types.Block

		blockHashParam, isBlockHashParam := args.Get(0).(struct {
			BlockId string `json:"block_id"`
		})

		blockHeightParam, isBlockHeightParam := args.Get(0).(struct {
			BlockId int64 `json:"block_id"`
		})

		if isBlockHashParam && m.BlockByHashMap[blockHashParam.BlockId] != emptyResponse {
			if m.BlockByHashMap[blockHashParam.BlockId].Error != nil {
				return []interface{}{block, m.BlockByHashMap[blockHashParam.BlockId].Error}
			}

			if response, Ok := (m.BlockByHashMap[blockHashParam.BlockId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &block)
				if err != nil {
					return []interface{}{block, err}
				}

				return []interface{}{block, nil}
			}

		}

		if isBlockHeightParam && m.BlockByHeightMap[blockHeightParam.BlockId] != emptyResponse {
			if m.BlockByHeightMap[blockHeightParam.BlockId].Error != nil {
				return []interface{}{block, m.BlockByHeightMap[blockHeightParam.BlockId].Error}
			}

			if response, Ok := (m.BlockByHeightMap[blockHeightParam.BlockId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &block)
				if err != nil {
					return []interface{}{block, err}
				}

				return []interface{}{block, nil}
			}

		}

		return []interface{}{block, errors.New("invalid Param")}
	}
}

func (m *MockApi) BlockProducersFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var blockProducers types.BlockProducers
		param, Ok := args.Get(0).([]string)

		if Ok && m.BlockProducersMap[param[0]] != emptyResponse {
			if m.BlockProducersMap[param[0]].Error != nil {
				return []interface{}{blockProducers, m.BlockProducersMap[param[0]].Error}
			}

			if response, Ok := (m.BlockProducersMap[param[0]].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &blockProducers)
				if err != nil {
					return []interface{}{blockProducers, err}
				}

				return []interface{}{blockProducers, nil}
			}
		}

		return []interface{}{types.BlockProducers{}, errors.New("invalid Param")}
	}
}

func (m *MockApi) BroadcastTxAsyncFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var response types.CryptoHash
		if m.TransactionHash != emptyResponse {
			if m.TransactionHash.Error != nil {
				return []interface{}{types.CryptoHash{}, m.TransactionHash.Error}
			}

			if transactionHash, Ok := (m.TransactionHash.Reponse).(string); Ok {
				return []interface{}{types.NewCryptoHash(transactionHash), nil}
			}
		}

		return []interface{}{response, errors.New("invalid Param")}
	}
}

func (m *MockApi) CallFunctionFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var response types.CallFunctionResponse

		param, Ok := args.Get(0).(types.CallFunction)

		if Ok && param.MethodName == "get_status" {
			var linkParam struct {
				Link chain.BTPAddress `json:"link"`
			}
			data, err := base64.URLEncoding.DecodeString(param.ArgumentsB64)
			if err != nil {
				return []interface{}{response, err}
			}

			err = json.Unmarshal(data, &linkParam)
			if err != nil {
				return []interface{}{response, err}
			}

			if m.BmcLinkStatusMap[linkParam.Link.ContractAddress()] != emptyResponse {
				if m.BmcLinkStatusMap[linkParam.Link.ContractAddress()].Error != nil {
					return []interface{}{response, m.BmcLinkStatusMap[linkParam.Link.ContractAddress()].Error}
				}

				if data, Ok := (m.BmcLinkStatusMap[linkParam.Link.ContractAddress()].Reponse).([]byte); Ok {
					err := json.Unmarshal(data, &response)
					if err != nil {
						return []interface{}{response, err}
					}

					return []interface{}{response, nil}
				}
			}

		}

		return []interface{}{response, errors.New("invalid Param")}
	}
}

func (m *MockApi) ChangesFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var changes types.ContractStateChange

		param, Ok := args.Get(0).(struct {
			ChangeType string   `json:"changes_type"`
			AccountIds []string `json:"account_ids"`
			KeyPrefix  string   `json:"key_prefix_base64"`
			BlockId    int64    `json:"block_id"`
		})

		if Ok && m.ContractStateChangeMap[param.BlockId] != emptyResponse {
			if m.ContractStateChangeMap[param.BlockId].Error != nil {
				return []interface{}{changes, m.ContractStateChangeMap[param.BlockId].Error}
			}

			if response, Ok := (m.ContractStateChangeMap[param.BlockId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &changes)
				if err != nil {
					return []interface{}{changes, err}
				}

				return []interface{}{changes, nil}
			}
		}

		return []interface{}{changes, errors.New("invalid Param")}
	}
}

func (m *MockApi) StatusFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var chainStatus types.ChainStatus

		if m.LatestChainStatus != emptyResponse {
			if m.LatestChainStatus.Error != nil {
				return []interface{}{chainStatus, m.LatestChainStatus.Error}
			}

			if response, Ok := (m.LatestChainStatus.Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &chainStatus)
				if err != nil {
					return []interface{}{chainStatus, err}
				}

				return []interface{}{chainStatus, nil}
			}
		}

		return []interface{}{chainStatus, errors.New("invalid Param")}
	}
}

func (m *MockApi) TransactionFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var transactionResult types.TransactionResult

		param, Ok := args.Get(0).([]string)

		if Ok && m.TransactionResultMap[param[0]] != emptyResponse {
			if m.TransactionResultMap[param[0]].Error != nil {
				return []interface{}{transactionResult, m.TransactionResultMap[param[0]].Error}
			}

			if response, Ok := (m.TransactionResultMap[param[0]].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &transactionResult)
				if err != nil {
					return []interface{}{transactionResult, err}
				}

				return []interface{}{transactionResult, nil}
			}
		}

		return []interface{}{transactionResult, errors.New("invalid Param")}
	}
}

func (m *MockApi) ViewAccessKeyFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var accessKeyResponse types.AccessKeyResponse

		param, Ok := args.Get(0).(struct {
			AccountId    string `json:"account_id"`
			PublicKey    string `json:"public_key"`
			Finality     string `json:"finality"`
			Request_type string `json:"request_type"`
		})

		if Ok && m.AccessKeyMap[param.AccountId] != emptyResponse {
			if m.AccessKeyMap[param.AccountId].Error != nil {
				return []interface{}{accessKeyResponse, m.AccessKeyMap[param.AccountId].Error}
			}

			if response, Ok := (m.AccessKeyMap[param.AccountId].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &accessKeyResponse)
				if err != nil {
					return []interface{}{accessKeyResponse, err}
				}

				return []interface{}{accessKeyResponse, nil}
			}
		}

		return []interface{}{accessKeyResponse, errors.New("invalid Param")}
	}
}

func (m *MockApi) ViewAccountFactory() func(args mock.Arguments) mock.Arguments {
	return func(args mock.Arguments) mock.Arguments {
		var account types.Account

		param, Ok := args.Get(0).(struct {
			AccountId    types.AccountId `json:"account_id"`
			Finality     string          `json:"finality"`
			Request_type string          `json:"request_type"`
		})

		if Ok && m.AccountMap[string(param.AccountId)] != emptyResponse {
			if m.AccountMap[string(param.AccountId)].Error != nil {
				return []interface{}{nil, m.AccountMap[string(param.AccountId)].Error}
			}

			if response, Ok := (m.AccountMap[string(param.AccountId)].Reponse).([]byte); Ok {
				err := json.Unmarshal(response, &account)
				if err != nil {
					return []interface{}{account, err}
				}

				return []interface{}{account, nil}
			}

		}
		return []interface{}{account, errors.New("invalid Param")}
	}
}
