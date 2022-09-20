package bsc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/bsc/mocks"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	ICON_BMC          = "btp://0x7.icon/cx8a6606d526b96a16e6764aee5d9abecf926689df"
	BSC_BMC_PERIPHERY = "btp://0x61.bsc/0xB4fC4b3b4e3157448B7D279f06BC8e340d63e2a9"
	BlockHeight       = 21447824
)

func newTestReceiver(t *testing.T, src, dst chain.BTPAddress) chain.Receiver {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	mp := map[string]interface{}{"syncConcurrency": 2}
	res, err := json.Marshal(mp)
	require.NoError(t, err)
	receiver, err := NewReceiver(src, dst, []string{url}, res, log.New())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return receiver
}

func newTestClient(t *testing.T, bmcAddr string) IClient {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	cls, _, err := newClients([]string{url}, bmcAddr, log.New())
	require.NoError(t, err)
	return cls[0]
}

func TestMedianGasPrice(t *testing.T) {
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	cls, _, err := newClients([]string{url}, BSC_BMC_PERIPHERY, log.New())
	require.NoError(t, err)

	_, _, err = cls[0].GetMedianGasPriceForBlock(context.Background())
	require.NoError(t, err)
}

func TestFilterLogs(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set(BSC_BMC_PERIPHERY)
	require.NoError(t, err)
	err = dst.Set(ICON_BMC)
	require.NoError(t, err)

	recv := newTestReceiver(t, src, dst).(*receiver)
	if recv == nil {
		t.Fatal(errors.New("Receiver is nil"))
	}
	exists, err := recv.hasBTPMessage(context.Background(), big.NewInt(BlockHeight))
	require.NoError(t, err)
	if !exists {
		require.NoError(t, errors.New("Expected true"))
	}
}

func TestSubscribeMessage(t *testing.T) {
	var src, dst chain.BTPAddress
	err := src.Set(BSC_BMC_PERIPHERY)
	err = dst.Set(ICON_BMC)
	if err != nil {
		fmt.Println(err)
	}

	recv := newTestReceiver(t, src, dst).(*receiver)

	ctx, cancel := context.Background(), func() {}
	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()
	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := recv.Subscribe(ctx,
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    75,
			Height: uint64(BlockHeight),
		})
	require.NoError(t, err, "failed to subscribe")

	for {
		defer cancel()
		select {
		case err := <-srcErrCh:
			t.Logf("subscription closed: %v", err)
			t.FailNow()
		case msg := <-srcMsgCh:
			if len(msg.Receipts) > 0 && msg.Receipts[0].Height == 21447824 {
				// received event exit
				return
			}
		}
	}
}

func TestReceiver_GetReceiptProofs(t *testing.T) {
	cl := newTestClient(t, BSC_BMC_PERIPHERY)
	header, err := cl.GetHeaderByHeight(big.NewInt(BlockHeight))
	require.NoError(t, err)
	hash := header.Hash()
	receipts, err := cl.GetBlockReceipts(hash)
	require.NoError(t, err)
	receiptsRoot := ethTypes.DeriveSha(receipts, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptsRoot.Bytes(), header.ReceiptHash.Bytes()) {
		err = fmt.Errorf(
			"invalid receipts: remote=%v, local=%v",
			header.ReceiptHash, receiptsRoot)
		require.NoError(t, err)
	}
}

func TestVerify(t *testing.T) {
	height := uint64(22169979)
	blockHash, err := hexutil.Decode("0x489b5865c1b015fa03177c30a4286533f02d2086c3db5f751180519f872fc37f")
	require.NoError(t, err)
	validatorData, err := hexutil.Decode("0xd98301010b846765746889676f312e31362e3130856c696e75780000de3b3a04049153b8dae0a232ac90d20c78f1a5d1de7b7dc51284214b9b9c85549ab3d2b972df0deef66ac2c935552c16704d214347f29fa77f77da6d75d7c7526d6247501b822fd4eaa76fcb64baea360279497f96c5d20b2a975c050e4220be276ace4892f4b41a980a75ecd1309ea12fa2ed87a8744fbfc9b863d5a2959d3f95eae5dc7d70144ce1b73b403b7eb6e0b71b214cb885500844365e95cd9942c7276e7fd833329df8450664d5960414752117d15811254efed1fb30e82660f82ce03df6536cc69315173fea12f202c1c1d0d165d5efb87dc2882d1602fdd3c1a11a03c86e01")
	require.NoError(t, err)
	opts := VerifierOptions{
		BlockHeight:   height,
		BlockHash:     blockHash,
		ValidatorData: validatorData,
	}
	vr := &Verifier{
		mu:         sync.RWMutex{},
		next:       big.NewInt(int64(opts.BlockHeight)),
		parentHash: common.BytesToHash(opts.BlockHash),
		validators: map[ethCommon.Address]bool{},
		chainID:    big.NewInt(97),
	}
	vr.validators, err = getValidatorMapFromHex(opts.ValidatorData)
	require.NoError(t, err)
	cl := newTestClient(t, BSC_BMC_PERIPHERY)
	header, err := cl.GetHeaderByHeight(big.NewInt(int64(opts.BlockHeight)))
	require.NoError(t, err)
	newHeader, err := cl.GetHeaderByHeight(big.NewInt(int64(opts.BlockHeight + 1)))
	require.NoError(t, err)
	err = vr.Verify(header, newHeader, nil)
	require.NoError(t, err)
}

func TestReceiver_newVerifier_NoValidators(t *testing.T) {
	cl := new(mocks.IClient)
	cl.On("GetBalance", mock.Anything, mock.Anything).Return(big.NewInt(1), errors.New("hex_addrs "))
	res, err := cl.GetBalance(context.TODO(), "hex_addr")
	fmt.Println(err)
	fmt.Println(res)
}
