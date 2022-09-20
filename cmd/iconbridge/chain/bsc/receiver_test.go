package bsc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/bsc/mocks"
	"github.com/icon-project/icon-bridge/common/intconv"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
	"github.com/pkg/errors"
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

func TestReceiver_MockReceiverOptions_UnmarshalWithVerifier(t *testing.T) {
	var opts ReceiverOptions
	jsonReceiverOptions := `{"syncConcurrency":100,"verifier":{"blockHeight":22169979,"parentHash":"0x489b5865c1b015fa03177c30a4286533f02d2086c3db5f751180519f872fc37f", "validatorData":"0xd98301010b846765746889676f312e31362e3130856c696e75780000de3b3a04049153b8dae0a232ac90d20c78f1a5d1de7b7dc51284214b9b9c85549ab3d2b972df0deef66ac2c935552c16704d214347f29fa77f77da6d75d7c7526d6247501b822fd4eaa76fcb64baea360279497f96c5d20b2a975c050e4220be276ace4892f4b41a980a75ecd1309ea12fa2ed87a8744fbfc9b863d5a2959d3f95eae5dc7d70144ce1b73b403b7eb6e0b71b214cb885500844365e95cd9942c7276e7fd833329df8450664d5960414752117d15811254efed1fb30e82660f82ce03df6536cc69315173fea12f202c1c1d0d165d5efb87dc2882d1602fdd3c1a11a03c86e01"}}`

	json.Unmarshal([]byte(jsonReceiverOptions), &opts)
	require.NotNil(t, opts)
	require.NotNil(t, opts.Verifier)
	require.NotNil(t, opts.SyncConcurrency)
	require.EqualValues(t, 100, opts.SyncConcurrency)
	require.NotNil(t, opts.Verifier.BlockHeight)
	require.EqualValues(t, 22169979, opts.Verifier.BlockHeight)
	require.NotNil(t, opts.Verifier.BlockHash)
	require.EqualValues(t, "0x489b5865c1b015fa03177c30a4286533f02d2086c3db5f751180519f872fc37f", opts.Verifier.BlockHash.String())
	require.NotNil(t, opts.Verifier.ValidatorData)
	require.EqualValues(t, "0xd98301010b846765746889676f312e31362e3130856c696e75780000de3b3a04049153b8dae0a232ac90d20c78f1a5d1de7b7dc51284214b9b9c85549ab3d2b972df0deef66ac2c935552c16704d214347f29fa77f77da6d75d7c7526d6247501b822fd4eaa76fcb64baea360279497f96c5d20b2a975c050e4220be276ace4892f4b41a980a75ecd1309ea12fa2ed87a8744fbfc9b863d5a2959d3f95eae5dc7d70144ce1b73b403b7eb6e0b71b214cb885500844365e95cd9942c7276e7fd833329df8450664d5960414752117d15811254efed1fb30e82660f82ce03df6536cc69315173fea12f202c1c1d0d165d5efb87dc2882d1602fdd3c1a11a03c86e01", opts.Verifier.ValidatorData.String())
}

func TestReceiver_MockReceiverOptions_UnmarshalWithoutVerifier(t *testing.T) {
	// Verifier should be nil if not passed
	var empty_opts ReceiverOptions
	jsonReceiverOptions := `{"syncConcurrency":100}`
	json.Unmarshal([]byte(jsonReceiverOptions), &empty_opts)
	require.NotNil(t, empty_opts)
	require.Nil(t, empty_opts.Verifier)
	require.NotNil(t, empty_opts.SyncConcurrency)
	require.EqualValues(t, 100, empty_opts.SyncConcurrency)
}

func TestReceiver_MockNewVerifier(t *testing.T) {
	// verifier options
	height := int64(22169979)
	blockHash, err := hexutil.Decode("0x489b5865c1b015fa03177c30a4286533f02d2086c3db5f751180519f872fc37f")
	require.NoError(t, err)
	validatorData, err := hexutil.Decode("0xd98301010b846765746889676f312e31362e3130856c696e75780000de3b3a04049153b8dae0a232ac90d20c78f1a5d1de7b7dc51284214b9b9c85549ab3d2b972df0deef66ac2c935552c16704d214347f29fa77f77da6d75d7c7526d6247501b822fd4eaa76fcb64baea360279497f96c5d20b2a975c050e4220be276ace4892f4b41a980a75ecd1309ea12fa2ed87a8744fbfc9b863d5a2959d3f95eae5dc7d70144ce1b73b403b7eb6e0b71b214cb885500844365e95cd9942c7276e7fd833329df8450664d5960414752117d15811254efed1fb30e82660f82ce03df6536cc69315173fea12f202c1c1d0d165d5efb87dc2882d1602fdd3c1a11a03c86e01")
	require.NoError(t, err)
	opts := &VerifierOptions{
		BlockHeight:   uint64(height),
		BlockHash:     blockHash,
		ValidatorData: validatorData,
	}
	validatorMap := map[ethCommon.Address]bool{
		ethCommon.HexToAddress("0x049153b8DAe0a232Ac90D20C78f1a5D1dE7B7dc5"): true,
		ethCommon.HexToAddress("0x1284214b9b9c85549aB3D2b972df0dEEf66aC2c9"): true,
		ethCommon.HexToAddress("0x35552c16704d214347f29Fa77f77DA6d75d7C752"): true,
		ethCommon.HexToAddress("0x6d6247501b822FD4Eaa76FCB64bAEa360279497f"): true,
		ethCommon.HexToAddress("0x96C5D20b2a975c050e4220BE276ACe4892f4b41A"): true,
		ethCommon.HexToAddress("0x980A75eCd1309eA12fa2ED87A8744fBfc9b863D5"): true,
		ethCommon.HexToAddress("0xA2959D3F95eAe5dC7D70144Ce1b73b403b7EB6E0"): true,
		ethCommon.HexToAddress("0xB71b214Cb885500844365E95CD9942C7276E7fD8"): true,
	}

	// mock client
	cl := new(mocks.IClient)
	cl.On("GetChainID").Return(big.NewInt(97))
	cl.On("GetHeaderByHeight", big.NewInt(height)).Return(&ethTypes.Header{ParentHash: ethCommon.BytesToHash(blockHash)}, nil)
	cl.On("GetHeaderByHeight", big.NewInt(height-height%int64(defaultEpochLength))).Return(&ethTypes.Header{Extra: validatorData}, nil)

	rx := &receiver{
		cls: []IClient{cl},
	}
	vr, err := rx.newVerifier(opts)
	require.NoError(t, err)

	require.NotNil(t, vr)
	require.Nil(t, err)
	require.Equal(t, vr.Next().Cmp(big.NewInt(int64(opts.BlockHeight))), 0)
	require.Equal(t, vr.ParentHash().String(), opts.BlockHash.String())
	for k := range validatorMap {
		require.Equal(t, vr.IsValidator(k), true)
	}
	require.Equal(t, vr.IsValidator(ethCommon.HexToAddress("abc")), false)
}

func TestReceiver_MockVerifyAndUpdate_CorrectHeader(t *testing.T) {
	height := int64(22169979)
	blockHash, err := hexutil.Decode("0x489b5865c1b015fa03177c30a4286533f02d2086c3db5f751180519f872fc37f")
	require.NoError(t, err)
	validatorData, err := hexutil.Decode("0xd98301010b846765746889676f312e31362e3130856c696e75780000de3b3a04049153b8dae0a232ac90d20c78f1a5d1de7b7dc51284214b9b9c85549ab3d2b972df0deef66ac2c935552c16704d214347f29fa77f77da6d75d7c7526d6247501b822fd4eaa76fcb64baea360279497f96c5d20b2a975c050e4220be276ace4892f4b41a980a75ecd1309ea12fa2ed87a8744fbfc9b863d5a2959d3f95eae5dc7d70144ce1b73b403b7eb6e0b71b214cb885500844365e95cd9942c7276e7fd833329df8450664d5960414752117d15811254efed1fb30e82660f82ce03df6536cc69315173fea12f202c1c1d0d165d5efb87dc2882d1602fdd3c1a11a03c86e01")
	require.NoError(t, err)
	opts := &VerifierOptions{
		BlockHeight:   uint64(height),
		BlockHash:     blockHash,
		ValidatorData: validatorData,
	}

	// Header
	headerStr := "7b22706172656e7448617368223a22307834383962353836356331623031356661303331373763333061343238363533336630326432303836633364623566373531313830353139663837326663333766222c2273686133556e636c6573223a22307831646363346465386465633735643761616238356235363762366363643431616433313234353162393438613734313366306131343266643430643439333437222c226d696e6572223a22307836643632343735303162383232666434656161373666636236346261656133363032373934393766222c227374617465526f6f74223a22307863336334343462666261656634333061666633376463613830363765656333653831656230663766633561653737656132663438383539363762613862666533222c227472616e73616374696f6e73526f6f74223a22307834626466323861626539373931373561393961613035323438363139643163393662653332623136343533306266333339363662326264373164376466643535222c227265636569707473526f6f74223a22307866396231623165616534383737353236373031316163613030616533353264656233663431353237346237333462353765643935633736643863656232643762222c226c6f6773426c6f6f6d223a2230783030323030383030323030323030613031303030303032303231303030323030303030383030383030303138303030303430313030303130313330303034313030343038303030303030343030303030303038383438303430343130303030303034303030303030313030363030303434303034303430323030323032393038313030313830313030303030303238303830303430343239303032303030303132383130383631313030323031303332303134323138303231383830303134313030383331303230306230323032343030303030383330383030323034383030306130303030303030313031303030303030343030303130303038306535343431303030383030343030303030303430323030303430343230303030303030323430303430343830313030303031313230303030303031303030323030303230303330303030353161303030343130303061303030323030303032303030303030303030613030303030303030303030303030303030313831383030303033383030303032353032303030303030303030383832303031303030303030313434313034303230303030303030303230383038303063383032323430383730313230303130303230303030303032303532303930303230383863303030303030316330303030303830343631303030303230303030323830303632303030343034222c22646966666963756c7479223a22307832222c226e756d626572223a22307831353234393762222c226761734c696d6974223a22307831636133343235222c2267617355736564223a223078316361376336222c2274696d657374616d70223a2230783633303335313931222c22657874726144617461223a2230786439383330313031306138343637363537343638383936373666333132653331333632653331333538353663363936653735373830303030646533623361303466613632643864323262353265636333366436646638336432376333663932643462373931626566336466366263643838326332616366343262663830336130333636653639356433646365333366623630393337656137323730353966383361643337383032336139323435316639366338363034626339316466633935343030222c226d697848617368223a22307830303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030222c226e6f6e6365223a22307830303030303030303030303030303030222c2262617365466565506572476173223a6e756c6c2c2268617368223a22307830356333646335303335633335643431396336333064366138386537383431616335633937653232373164633966303936623139643561316266316536353036227d"
	nextHeaderStr := "7b22706172656e7448617368223a22307830356333646335303335633335643431396336333064366138386537383431616335633937653232373164633966303936623139643561316266316536353036222c2273686133556e636c6573223a22307831646363346465386465633735643761616238356235363762366363643431616433313234353162393438613734313366306131343266643430643439333437222c226d696e6572223a22307839366335643230623261393735633035306534323230626532373661636534383932663462343161222c227374617465526f6f74223a22307866653232626564613965356162386464363563316538323363376337643438326662306633653766653439363433323035663134613865356463636636306632222c227472616e73616374696f6e73526f6f74223a22307866323931613135306461386137366332383032363434313061326635306339366462353536633130656336313462653136313732633632313165656335656466222c227265636569707473526f6f74223a22307861366337653333633162376631646561663766393162343464393834363164626339326163353136333837373237326632353239373561363363336138653263222c226c6f6773426c6f6f6d223a2230783032323030343230303035303032383031303030303030303331303032363030303832303030343038313161303031303630313030306230313830303034303430303032303030303030363034303030303038303430303030303030323038313034383230313032303038323030303430303036306338303830323231393463313031396334333031303430303330316330343630633064303032303030303033623130303431383030303830323132613034323038303231633031303030313030613230303230306132323032343030323030303330393061303234383331386130323432303030303030303032303438303030313930303039303463303031303032383030343030303030303531303030303430363238303430303232303031383030343830313430303030306130303030303030303030333030303230303231313038343030303030303830313061303030303030323030303030393230613038323030303031303030303430303031303030303030303030303133383238303032373265303434303030303030383030383230303030303030313634313034303030303230303434303630383038383065303032303032313230363238303138323230303030303030303032303930313231313830303430313130393038323030303030303034303030303030383230323030303638303830343034222c22646966666963756c7479223a22307832222c226e756d626572223a22307831353234393763222c226761734c696d6974223a22307831633963333830222c2267617355736564223a223078393162363162222c2274696d657374616d70223a2230783633303335313934222c22657874726144617461223a2230786438383330313031306238343637363537343638383836373666333132653331333832653332383536633639366537353738303030303030646533623361303430646537653531623362656236613461343862663564366462613139326436633366383632623638373965316439363532623562313134653164393137373733333065386161346130373963306135316263396336303233333433323731303761626339323130396231626537643030333133373239376336366637376633343031222c226d697848617368223a22307830303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030222c226e6f6e6365223a22307830303030303030303030303030303030222c2262617365466565506572476173223a6e756c6c2c2268617368223a22307831396434316131373561343734656236356435386335316136323239336338343563333364666361633637343761303665613338306134373561356534383861227d"
	headerBytes, err := hex.DecodeString(headerStr)
	require.NoError(t, err)
	nextHeaderBytes, err := hex.DecodeString(nextHeaderStr)
	require.NoError(t, err)
	header := new(ethTypes.Header)
	nextHeader := new(ethTypes.Header)
	err = json.Unmarshal(headerBytes, header)
	require.NoError(t, err)
	err = json.Unmarshal(nextHeaderBytes, nextHeader)
	require.NoError(t, err)
	validatorMap := map[ethCommon.Address]bool{
		ethCommon.HexToAddress("0x049153b8DAe0a232Ac90D20C78f1a5D1dE7B7dc5"): true,
		ethCommon.HexToAddress("0x1284214b9b9c85549aB3D2b972df0dEEf66aC2c9"): true,
		ethCommon.HexToAddress("0x35552c16704d214347f29Fa77f77DA6d75d7C752"): true,
		ethCommon.HexToAddress("0x6d6247501b822FD4Eaa76FCB64bAEa360279497f"): true,
		ethCommon.HexToAddress("0x96C5D20b2a975c050e4220BE276ACe4892f4b41A"): true,
		ethCommon.HexToAddress("0x980A75eCd1309eA12fa2ED87A8744fBfc9b863D5"): true,
		ethCommon.HexToAddress("0xA2959D3F95eAe5dC7D70144Ce1b73b403b7EB6E0"): true,
		ethCommon.HexToAddress("0xB71b214Cb885500844365E95CD9942C7276E7fD8"): true,
	}

	// Client
	cl := new(mocks.IClient)
	cl.On("GetChainID").Return(big.NewInt(97))
	cl.On("GetHeaderByHeight", big.NewInt(height)).Return(&ethTypes.Header{ParentHash: ethCommon.BytesToHash(blockHash)}, nil)
	cl.On("GetHeaderByHeight", big.NewInt(height-height%int64(defaultEpochLength))).Return(&ethTypes.Header{Extra: validatorData}, nil)

	rx := &receiver{
		cls: []IClient{cl},
	}
	vr, err := rx.newVerifier(opts)

	err = vr.Verify(header, nextHeader, nil)
	require.NoError(t, err)
	err = vr.Update(header) // should not update because header.Number % defaultEpochLength != 0
	require.NoError(t, err)
	require.Equal(t, vr.ParentHash().String(), header.Hash().String())
	require.Equal(t, vr.Next().Cmp(header.Number.Add(header.Number, big.NewInt(1))), 0)
	for k := range validatorMap {
		require.Equal(t, vr.IsValidator(k), true)
	}
	require.Equal(t, vr.IsValidator(ethCommon.HexToAddress("abc")), false)
}

func TestSender_NewObj(t *testing.T) {
	//senderOpts := `{"gas_limit": 24000000,"tx_data_size_limit": 8192,"balance_threshold": "100000000000000000000","boost_gas_price": 1}`
	thres := intconv.BigInt{}
	thres.SetString("100000000000000000000", 10)
	sopts := senderOptions{
		GasLimit:         24000000,
		TxDataSizeLimit:  8192,
		BalanceThreshold: thres,
		BoostGasPrice:    1,
	}
	raw, err := json.Marshal(sopts)
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	url := "https://data-seed-prebsc-1-s1.binance.org:8545"
	s, err := NewSender(
		chain.BTPAddress(BSC_BMC_PERIPHERY),
		chain.BTPAddress(ICON_BMC),
		[]string{url}, &wallet.EvmWallet{Skey: privKey, Pkey: &privKey.PublicKey},
		raw,
		log.New(),
	)
	balance, threshold, err := s.Balance(context.TODO())
	require.NoError(t, err)
	require.Equal(t, balance.Cmp(big.NewInt(0)), 0)
	require.Equal(t, threshold.String(), thres.String())

	msg := &chain.Message{
		From: "SourceName",
		Receipts: []*chain.Receipt{{
			Index:  0,
			Height: 1,
			Events: []*chain.Event{},
		}},
	}
	_, _, err = s.Segment(context.TODO(), msg)
	require.NoError(t, err)
}
