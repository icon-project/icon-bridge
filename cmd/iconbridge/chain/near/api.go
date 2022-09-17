package near

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
)

type api struct {
	host string
	*jsonrpc.Client
}

func (api *api) Block(param interface{}) (response types.Block, err error) {
	if _, err := api.Do("block", param, &response); err != nil {
		return types.Block{}, err
	}

	return response, nil
}

func (api *api) BroadcastTxCommit(param interface{}) (response types.TransactionResult, err error) {
	if _, err := api.Do("broadcast_tx_commit", param, &response); err != nil {
		return types.TransactionResult{}, err
	}

	return response, nil
}

func (api *api) BroadcastTxAsync(param interface{}) (response types.CryptoHash, err error) {
	if _, err := api.Do("broadcast_tx_async", param, &response); err != nil {
		return types.CryptoHash{}, err
	}

	return response, nil
}

func (api *api) CallFunction(param interface{}) (response types.CallFunctionResponse, err error) {
	if _, err := api.Do("query", param, &response); err != nil {
		return types.CallFunctionResponse{}, err
	}

	return response, nil
}

func (api *api) Changes(param interface{}) (response types.ContractStateChange, err error) {
	if _, err := api.Do("EXPERIMENTAL_changes", param, &response); err != nil {
		return types.ContractStateChange{}, err
	}

	return response, nil
}

func (api *api) LightClientProof(param interface{}) (response types.ReceiptProof, err error) {
	if _, err := api.Do("EXPERIMENTAL_light_client_proof", param, &response); err != nil {
		return types.ReceiptProof{}, err
	}

	return response, nil
}

func (api *api) Status(param interface{}) (response types.ChainStatus, err error) {
	if _, err := api.Do("status", param, &response); err != nil {
		return types.ChainStatus{}, err
	}

	return response, nil
}

func (api *api) Transaction(param interface{}) (response types.TransactionResult, err error) {
	if _, err := api.Do("tx", param, &response); err != nil {
		return types.TransactionResult{}, err
	}

	return response, nil
}

func (api *api) ViewAccessKey(param interface{}) (response types.AccessKeyResponse, err error) {
	if _, err := api.Do("query", param, &response); err != nil {
		return types.AccessKeyResponse{}, err
	}

	return response, nil
}

func (api *api) ViewAccount(param interface{}) (response types.Account, err error) {
	if _, err := api.Do("query", param, &response); err != nil {
		return types.Account{}, err
	}

	return response, nil
}
