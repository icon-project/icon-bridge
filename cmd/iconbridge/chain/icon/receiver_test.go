package icon

import (
	"context"
	"encoding/json"
	ethc "github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/mocks"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/mock"
	"net/http"
	"strings"
	"testing"

	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestReceiver(t *testing.T) {
	srcAddress := "btp://0x1.icon/cx997849d3920d338ed81800833fbb270c785e743d"
	dstAddress := "btp://0x63564c40.hmny/0xa69712a3813d0505bbD55AeD3fd8471Bc2f722DD"
	srcEndpoint := []string{"https://ctz.solidwallet.io/api/v3/icon_dex"}
	var height uint64 = 0x307f54a
	var seq uint64 = 628
	opts := map[string]interface{}{
		"verifier": map[string]interface{}{
			"blockHeight":    0x307f540,
			"validatorsHash": "0xa6760c547c3f76b7071658ef383d69ec01e11ea71d695600788695b50659e409",
		},
	}
	rawOpts, err := json.Marshal(&opts)
	if err != nil {
		panic(err)
	}
	l := log.New()
	log.SetGlobalLogger(l)

	client := NewClient(srcEndpoint[0], l)

	// log.AddForwarder(&log.ForwarderConfig{Vendor: log.HookVendorSlack, Address: "https://hooks.slack.com/services/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", Level: "info"})
	recv, err := NewReceiver(chain.BTPAddress(srcAddress), chain.BTPAddress(dstAddress), client, rawOpts, l)
	if err != nil {
		panic(err)
	}
	msgCh := make(chan *chain.Message)
	if errCh, err := recv.Subscribe(
		context.Background(), msgCh, chain.SubscribeOptions{Height: height, Seq: seq}); err != nil {
		panic(err)
	} else {
		for {
			select {
			case err := <-errCh:
				panic(err)
			case msg := <-msgCh:
				if len(msg.Receipts) > 0 && msg.Receipts[0].Height == 50853195 {
					// found event
					return
				}
			}
		}
	}
}

func TestNextValidatorHashFetch(t *testing.T) {

	var conUrl string = "https://ctz.solidwallet.io/api/v3/icon_dex" //devnet
	height := 50852431
	con := jsonrpc.NewJsonRpcClient(&http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1000}}, conUrl)
	getBlockHeaderByHeight := func(height int64, con *jsonrpc.Client) (*types.BlockHeader, error) {
		var header types.BlockHeader
		var result []byte
		_, err := con.Do("icx_getBlockHeaderByHeight", &types.BlockHeightParam{
			Height: types.NewHexInt(int64(height)),
		}, &result)
		require.NoError(t, err)

		_, err = vlcodec.RLP.UnmarshalFromBytes(result, &header)
		require.NoError(t, err)
		return &header, nil
	}

	getDatabyHash := func(req interface{}, resp interface{}, con *jsonrpc.Client) (interface{}, error) {
		_, err := con.Do("icx_getDataByHash", req, resp)
		require.NoError(t, err)
		return resp, nil
	}

	header, err := getBlockHeaderByHeight(int64(height), con)
	require.NoError(t, err)

	var validatorDataBytes []byte
	_, err = getDatabyHash(&types.DataHashParam{Hash: types.NewHexBytes(header.NextValidatorsHash)}, &validatorDataBytes, con)
	require.NoError(t, err)

	var validators [][]byte
	_, err = vlcodec.BC.UnmarshalFromBytes(validatorDataBytes, &validators)
	require.NoError(t, err)

	if common.HexBytes(header.NextValidatorsHash).String() != "0xa6760c547c3f76b7071658ef383d69ec01e11ea71d695600788695b50659e409" {
		err := errors.New("Invalid Validator Hash")
		require.NoError(t, err)
	}
}

func TestReceiver_newVerifier_NoValidators(t *testing.T) {
	clientMock := new(mocks.ClientMock)

	errMessage := "error NoValidators"
	// setup expectations
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, errors.New(errMessage))

	receiverOb := Receiver{
		Client: clientMock,
	}
	opts := types.VerifierOptions{
		//Can be any value
		ValidatorsHash: common.HexHash(ethc.Hex2Bytes("34d4ab43f7351fab97f93bc72d2e02c823b08a7c469c5da6ef01ccdd91f881f4")),
	}

	_, err := receiverOb.newVerifier(&opts)
	require.Error(t, err)
	require.Equal(t, errMessage, err.Error())
	clientMock.AssertExpectations(t)
}

func TestReceiver_newVerifier_NoBlockHeader(t *testing.T) {
	clientMock := new(mocks.ClientMock)

	errMessage := "error NoBlockHeader"
	// setup expectations
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, nil)
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, errors.New(errMessage))

	receiverOb := Receiver{
		Client: clientMock,
	}
	opts := types.VerifierOptions{
		//Can be any value
		ValidatorsHash: common.HexHash(ethc.Hex2Bytes("34d4ab43f7351fab97f93bc72d2e02c823b08a7c469c5da6ef01ccdd91f881f4")),
	}

	_, err := receiverOb.newVerifier(&opts)
	require.Error(t, err)
	require.Equal(t, errMessage, err.Error())
	clientMock.AssertExpectations(t)
}

func TestReceiver_newVerifier_NoVotes(t *testing.T) {
	clientMock := new(mocks.ClientMock)

	errMessage := "error NoVotes"
	// setup expectations
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, nil)
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, nil)
	clientMock.On("GetVotesByHeight", mock.Anything).Return(nil, errors.New(errMessage))

	receiverOb := Receiver{
		Client: clientMock,
	}
	opts := types.VerifierOptions{
		//Can be any value
		ValidatorsHash: common.HexHash(ethc.Hex2Bytes("34d4ab43f7351fab97f93bc72d2e02c823b08a7c469c5da6ef01ccdd91f881f4")),
	}

	_, err := receiverOb.newVerifier(&opts)
	require.Error(t, err)
	require.Equal(t, errMessage, err.Error())
	clientMock.AssertExpectations(t)
}

func TestReceiver_newVerifier_VerificationFailed(t *testing.T) {
	clientMock := new(mocks.ClientMock)

	errMessage := "verification failed"
	// setup expectations
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, nil)
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, nil)
	clientMock.On("GetVotesByHeight", mock.Anything).Return(nil, nil)

	receiverOb := Receiver{
		Client: clientMock,
	}
	opts := types.VerifierOptions{
		//Can be any value
		ValidatorsHash: common.HexHash(ethc.Hex2Bytes("34d4ab43f7351fab97f93bc72d2e02c823b08a7c469c5da6ef01ccdd91f881f4")),
	}

	_, err := receiverOb.newVerifier(&opts)
	require.Error(t, err)
	require.Equal(t, errMessage, err.Error())
	clientMock.AssertExpectations(t)
}

func TestReceiver_ReceiverOptions_Unmarshal(t *testing.T) {
	var opts ReceiverOptions

	jsonReceiverOptions := `{"syncConcurrency":100,"verifier":{"blockHeight":50853184,"validatorsHash":"0xa6760c547c3f76b7071658ef383d69ec01e11ea71d695600788695b50659e409"}}`

	json.Unmarshal([]byte(jsonReceiverOptions), &opts)

	require.NotNil(t, opts)
	require.NotNil(t, opts.SyncConcurrency)
	require.EqualValues(t, 100, opts.SyncConcurrency)
	require.NotNil(t, opts.Verifier)
	require.NotNil(t, opts.Verifier.BlockHeight)
	require.EqualValues(t, 50853184, opts.Verifier.BlockHeight)
	require.NotNil(t, opts.Verifier.ValidatorsHash)
}

func TestReceiver_SyncVerifier_InvalidHeight(t *testing.T) {
	receiverOb := &Receiver{}
	verifier := &Verifier{
		next: 101,
	}

	err := receiverOb.syncVerifier(verifier, 100)

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "invalid target height"))
}

func TestReceiver_SyncVerifier_HeightNextAreEql(t *testing.T) {
	receiverOb := &Receiver{}
	verifier := &Verifier{
		next: 100,
	}

	err := receiverOb.syncVerifier(verifier, 100)

	require.NoError(t, err)
}


func TestReceiver_validateRequests_RetryTimes_NotFoundBlockHeader(t *testing.T) {
	retryTimes  := 3

	verifier := &Verifier{
		next: 100,
	}
	requestCh := make(chan *request, 1)
	requestCh <- &request{height: 100, retry: retryTimes}

	clientMock := new(mocks.ClientMock)
	errMessage := "BlockHeader not found"

	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, errors.New(errMessage)).Times(retryTimes)

	responseCh := validateRequests(requestCh, clientMock, verifier, log.New())

	require.NotNil(t, responseCh)
	response := responseCh[0]

	require.NotNil(t, response)
	require.Error(t, response.err)
	require.True(t, strings.Contains(response.err.Error(), errMessage))
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t,"GetBlockHeaderByHeight", retryTimes)
}

func TestReceiver_validateRequests_RetryTimes_NotFoundVotes(t *testing.T) {
	retryTimes  := 1

	verifier := &Verifier{
		next: 100,
	}
	requestCh := make(chan *request, 1)
	requestCh <- &request{height: 100, retry: retryTimes}

	clientMock := new(mocks.ClientMock)
	errMessage := "Votes not found"

	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, nil).Times(retryTimes)
	clientMock.On("GetVotesByHeight", mock.Anything).Return(nil, errors.New(errMessage)).Times(retryTimes)

	responseCh := validateRequests(requestCh, clientMock, verifier, log.New())

	require.NotNil(t, responseCh)
	response := responseCh[0]
	require.NotNil(t, response)
	require.Error(t, response.err)
	require.True(t, strings.Contains(response.err.Error(), errMessage))
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t,"GetVotesByHeight", retryTimes)
}

func TestReceiver_validateRequests_RetryTimes_NotFoundValidators(t *testing.T) {
	retryTimes  := 2
	errMessage := "Validators not found"
	blockHeader := &types.BlockHeader{
		NextValidatorsHash: common.HexHash(ethc.Hex2Bytes("34d4ab43f7351fab97f93bc72d2e02c823b08a7c469c5da6ef01ccdd91f881f4")),
	}

	requestCh := make(chan *request, 1)
	requestCh <- &request{height: 100, retry: retryTimes}

	// setup expectations
	verifier := new(mocks.VerifierMock)
	verifier.On("Validators", mock.Anything).Return(make([]common.Address, 0)).Times(retryTimes)

	clientMock := new(mocks.ClientMock)
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(blockHeader, nil).Times(retryTimes)
	clientMock.On("GetVotesByHeight", mock.Anything).Return(nil, nil).Times(retryTimes)
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, errors.New(errMessage)).Times(retryTimes)


	responseCh := validateRequests(requestCh, clientMock, verifier, log.New())

	require.NotNil(t, responseCh)
	response := responseCh[0]
	require.NotNil(t, response)
	require.Error(t, response.err)
	require.True(t, strings.Contains(response.err.Error(), errMessage))
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t,"GetValidatorsByHash", retryTimes)
}

func TestReceiver_receiveLoop_CantCreateNewVerifier(t *testing.T) {
	clientMock := new(mocks.ClientMock)

	errMessage := "error NoValidators"
	// setup expectations
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, errors.New(errMessage))

	receiverOb := Receiver{
		Client: clientMock,
		// Verifier will be updated if  not nil.
		opts: ReceiverOptions{
			Verifier: &types.VerifierOptions{
				//Mock
			},
		},
	}

	err := receiverOb.receiveLoop(nil, 100, 50, nil)

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), errMessage))
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t,"GetValidatorsByHash", 1)
}


func TestReceiver_receiveLoop_processBlockResult_VerifyFalse(t *testing.T) {
	vrMock := new(mocks.VerifierMock)
	vrMock.On("Verify", mock.Anything, mock.Anything).Return(false, nil)

	blockResponse := receiverResponse{
		Height: 100,
 	}
	var next int64 = 110
	 var reconnectCallCounter = 0

	var reconnect = func() {
		reconnectCallCounter++
	}

	err := processBlockResult(&blockResponse, vrMock, &next, reconnect, nil,nil, log.New())

	require.NoError(t, err)
	require.EqualValues(t, 1, reconnectCallCounter)
	vrMock.AssertExpectations(t)
	vrMock.AssertNumberOfCalls(t,"Verify", 1)
}

func TestReceiver_receiveLoop_processBlockResult_VerifyFalse_2(t *testing.T) {
	vrMock := new(mocks.VerifierMock)
	errMessage := "error NoValidators"
	vrMock.On("Verify", mock.Anything, mock.Anything).Return(false, errors.New(errMessage))

	blockResponse := receiverResponse{
		Height: 100,
	}
	var reconnectCallCounter = 0
	var next int64 = 110
	var reconnect = func() {
		reconnectCallCounter++
	}

	err := processBlockResult(&blockResponse, vrMock, &next, reconnect, nil,nil, log.New())

	require.NoError(t, err)
	require.EqualValues(t, 1, reconnectCallCounter)
	vrMock.AssertExpectations(t)
	vrMock.AssertNumberOfCalls(t,"Verify", 1)
}

func TestReceiver_receiveLoop_processBlockResult_UpdateFalse(t *testing.T) {
	vrMock := new(mocks.VerifierMock)
	errMessage := "Update error"
	vrMock.On("Verify", mock.Anything, mock.Anything).Return(true, nil)
	vrMock.On("Update", mock.Anything, mock.Anything).Return(errors.New(errMessage))

	blockResponse := receiverResponse{
		Height: 100,
	}
	var next int64 = 110
	var reconnectCallCounter = 0

	var reconnect = func() {
		reconnectCallCounter++
	}

	err := processBlockResult(&blockResponse, vrMock, &next, reconnect, nil,nil, log.New())

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), errMessage))
	vrMock.AssertExpectations(t)
	vrMock.AssertNumberOfCalls(t,"Verify", 1)
	vrMock.AssertNumberOfCalls(t,"Update", 1)
}

func TestReceiver_receiveLoop_processBlockResult_countNext(t *testing.T) {
	blockResponse := receiverResponse{
		Height: 100,
	}
	var next int64 = 110
	var callBackMock = func(rs []*chain.Receipt) error {
		return nil
	}
	blockResultCh := make(chan *receiverResponse, 3)

	blockResponse2 := receiverResponse{
		Height: 101,
	}
	blockResultCh <- &blockResponse2

	err := processBlockResult(&blockResponse, nil, &next, nil, callBackMock, blockResultCh, log.New())

	require.NoError(t, err)
	require.EqualValues(t, 112, next)
}

func TestReceiver_receiveLoop_processBlockResult_callbackError(t *testing.T) {
	blockResponse := receiverResponse{
		Height: 100,
	}
	errMessage := "Callback error"
	var callBackCallCounter = 0
	var next int64 = 110

	var callBack = func(rs []*chain.Receipt) error {
		callBackCallCounter++
		return errors.New(errMessage)
	}

	err := processBlockResult(&blockResponse, nil, &next, nil, callBack,nil, log.New())

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), errMessage))
	require.EqualValues(t, 1, callBackCallCounter)
}

func TestReceiver_requestProcessor_invalidHash(t *testing.T) {
	blockResponse := receiverRequest{
		height: 1000,
		hash:types.HexBytes("Invalid_Hash"),
	}

	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}

	requestProcessor(&blockResponse, requestCh, nil, nil, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	require.True(t, strings.Contains(blockResponse.err.Error(), "invalid hash"))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), "invalid hash"))
}

func TestReceiver_requestProcessor_blockHeaderNotFound(t *testing.T) {
	blockResponse := receiverRequest{
		height: 1000,
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}

	clientMock := new(mocks.ClientMock)
	errMessage := "BlockHeader not found"
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, errors.New(errMessage)).Times(1)


	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}
	requestProcessor(&blockResponse, requestCh, nil, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	require.True(t, strings.Contains(blockResponse.err.Error(), errMessage))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), errMessage))
}

func TestReceiver_requestProcessor_votesNotFound(t *testing.T) {
	blockResponse := receiverRequest{
		height: 1000,
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}

	vrMock := new(mocks.VerifierMock)
	clientMock := new(mocks.ClientMock)
	errMessage := "Votes not found"
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(nil, nil).Times(1)
	clientMock.On("GetVotesByHeight", mock.Anything).Return(nil, errors.New(errMessage)).Times(1)


	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}
	requestProcessor(&blockResponse, requestCh, vrMock, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	require.True(t, strings.Contains(blockResponse.err.Error(), errMessage))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), errMessage))
}

func TestReceiver_requestProcessor_validatorsNotFound(t *testing.T) {
	blockResponse := &receiverRequest{
		height: 1000,
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}
	blockHeader := &types.BlockHeader{
		Height: 1000,
	}

	clientMock := new(mocks.ClientMock)
	errMessage := "Validators not found"
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(blockHeader, nil).Times(1)
	clientMock.On("GetVotesByHeight", mock.Anything).Return(nil, nil).Times(1)

	vrMock := new(mocks.VerifierMock)
	clientMock.On("GetValidatorsByHash", mock.Anything).Return(nil, errors.New(errMessage)).Times(1)
	vrMock.On("Validators", mock.Anything).Return([]common.Address{})


	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}
	requestProcessor(blockResponse, requestCh, vrMock, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	require.True(t, strings.Contains(blockResponse.err.Error(), errMessage))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), errMessage))
}

func TestReceiver_requestProcessor_HeaderUnmarshalError(t *testing.T) {
	blockResponse := receiverRequest{
		height: 1000,
		indexes: [][]types.HexInt{{types.NewHexInt(125)}},
		events: [][][]types.HexInt{{{types.NewHexInt(125)}}},
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}
	blockHeader := &types.BlockHeader{
		Height: 1000,
	}

	clientMock := new(mocks.ClientMock)
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(blockHeader, nil).Times(1)

	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}

	requestProcessor(&blockResponse, requestCh, nil, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	require.True(t, strings.Contains(blockResponse.err.Error(), "BlockHeaderResult.UnmarshalFromBytes"))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), "BlockHeaderResult.UnmarshalFromBytes"))
}

func TestReceiver_requestProcessor_ProofForEventsNotFound(t *testing.T) {
	blockResponse := receiverRequest{
		height: 1000,
		indexes: [][]types.HexInt{{types.NewHexInt(125)}},
		events: [][][]types.HexInt{{{types.NewHexInt(125)}}},
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}

	errMessage := "ProofForEvents not found"
	clientMock := new(mocks.ClientMock)
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(getSampleHeader(), nil).Times(1)
	clientMock.On("GetProofForEvents", mock.Anything).Return(nil, errors.New(errMessage)).Times(1)

	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}


	requestProcessor(&blockResponse, requestCh, nil, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	require.True(t, strings.Contains(blockResponse.err.Error(), errMessage))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), errMessage))
}

func TestReceiver_requestProcessor_ProofsError(t *testing.T) {
	blockResponse := receiverRequest{
		height: 1000,
		indexes: [][]types.HexInt{{types.NewHexInt(125)}},
		events: [][][]types.HexInt{{{types.NewHexInt(125)}}},
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}

	clientMock := new(mocks.ClientMock)
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(getSampleHeader(), nil).Times(1)
	clientMock.On("GetProofForEvents", mock.Anything).Return(nil, nil).Times(1)

	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}


	requestProcessor(&blockResponse, requestCh, nil, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	errMessage := "Proof does not include all events"
	require.True(t, strings.Contains(blockResponse.err.Error(), errMessage))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), errMessage))
}

func TestReceiver_requestProcessor_MPTProveError(t *testing.T) {
    events := [][][]types.HexInt{{{types.NewHexInt(125)}}, {{types.NewHexInt(125)}}}
	blockResponse := receiverRequest{
		height: 1000,
		indexes: [][]types.HexInt{{types.NewHexInt(125)}},
		events: events,
		hash:types.HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b"),
	}

	clientMock := new(mocks.ClientMock)
	// setup expectations
	clientMock.On("GetBlockHeaderByHeight", mock.Anything).Return(getSampleHeader(), nil).Times(1)
	proofs := [][][]byte{{{}}, {{}}}
	clientMock.On("GetProofForEvents", mock.Anything).Return(proofs, nil).Times(1)

	requestCh := make(chan *receiverRequest, 1)
	eventLogRawFilter := eventLogRawFilter{}


	requestProcessor(&blockResponse, requestCh, nil, clientMock, eventLogRawFilter, log.New())

	requestResult := <-requestCh
	require.Error(t, blockResponse.err)
	errMessage := "MPTProve Receipt"
	require.True(t, strings.Contains(blockResponse.err.Error(), errMessage))
	require.NotNil(t, requestResult)
	require.True(t, strings.Contains(requestResult.err.Error(), errMessage))
}









