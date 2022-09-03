package icon

import (
	"encoding/json"
	"fmt"
	"testing"

	ethc "github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/stretchr/testify/require"
)

func NewSampleTestVerifier() *Verifier {
	validatorsHash := common.HexHash(ethc.Hex2Bytes("34d4ab43f7351fab97f93bc72d2e02c823b08a7c469c5da6ef01ccdd91f881f4"))
	return &Verifier{
		next:               50000001,
		nextValidatorsHash: validatorsHash,
		validators: map[string][]common.Address{
			validatorsHash.String(): getSampleValidators(),
		},
	}
}

func getCommitVoteItem(ts int64, sig string) commitVoteItem {
	cv := commitVoteItem{Timestamp: ts}
	_sig, _ := json.Marshal(sig)
	cv.Signature.UnmarshalJSON([]byte(_sig))
	return cv
}

func getSampleHeader() *BlockHeader {
	var header BlockHeader
	json.Unmarshal([]byte(`{"Version":2,"Height":50000001,"Timestamp":1652523322961762,"Proposer":"AIFxnc/o9YygcES3vt5Jzs1h+b0/","PrevID":"wJ4ZEWkntBrtt1OXVElzQXqoBayMNTYsL/I8W2SoddY=","VotesHash":"ZxdjQqB1gvmvbHH+G3USOa/9dpplfkym2+0zVwipmj4=","NextValidatorsHash":"NNSrQ/c1H6uX+TvHLS4CyCOwinxGnF2m7wHM3ZH4gfQ=","PatchTransactionsHash":null,"NormalTransactionsHash":"dZOh1mazDheqYZHXv2d9VW8WYGAZzDlwlSQ9HPnRR6o=","LogsBloom":"AQAgcEgsGg8DBEEAMIhsOh8QAEKiMGicUi8YjMHIEagcBA==","Result":"+M6gTDv1XbPHjPnCiT5Pz9kY9QfB7W+j2pjJxIVZ0/BMN5r4AKBQAFawZ+KCE9Yavl66ZIKKnDnfGS672n662//Rvi8gQ7iI+IagcwNDeKHFwzcfmUzXMozEqoqim/Kn/K9xfaMiWK2agi2gVy3PwnfzQ/w2CqBXdUPZq6avi9EqIiuupnDSBLCwS7b4AKDXVDsjgbWsLdJWrXY0tDrnq5xw07gCTkjVTMjV1cUz/6C9MDzp0VE5J3q8nZFNKlV5nFeNeZBl1ZEb5aT776wH5w=="}`), &header)
	return &header
}

func getSampleValidators() []common.Address {
	return []common.Address{
		*common.MustNewAddress(ethc.Hex2Bytes("009c63f73d3c564a54d0eed84f90718b1ebed16f09")),
		*common.MustNewAddress(ethc.Hex2Bytes("0081719dcfe8f58ca07044b7bede49cecd61f9bd3f")),
		*common.MustNewAddress(ethc.Hex2Bytes("00ed7175f73f63ce8dfeede1db8c4b66179eb7a857")),
	}
}

func getSampleCommitVoteList() *commitVoteList {
	cvl := commitVoteList{
		Round: 0,
		BlockPartSetID: &PartSetID{
			Count: 1,
			Hash:  ethc.Hex2Bytes("3b27a2dea9d1e8ecd1c94ff723f9efe8ed79e54f0708fa459a57148ff2aab3f1"),
		},
		Items: []commitVoteItem{
			getCommitVoteItem(1652523324922454, "5QIv0HrkyBU0wqqy/f6HFhPiCbqf9GK11z46LyrL9WAQD25TZdthyZfJXd4B3+4eIMxzW4i5oXicbD6+UtbtWQE="),
			getCommitVoteItem(1652523324864943, "weofhyea6ixet/a1sKH986dRgYRoQZ6PxA9is90eIuJ/036poH3Hj28PtCKJ2ayWikbjkIYhpkBxFegnIkLnMgA="),
			getCommitVoteItem(1652523324882445, "ocjI0SOiMpd3ZCDWAmPqAyqaRZK4zi5A3cg9y4OFC8Ft/4H5Gkpfc2fCSkvzJMva0rPvUNLjgnyWUyKiUWILhgE="),
		},
	}
	return &cvl
}

func TestNextValidatorHash(t *testing.T) {
	raw := HexBytes("0xf86e950038f35eff5e5516b48a713fe3c8031c94124191f09500f526cc053c33a7c3a48b70111834cf3a71609f0c950014d4c29c4bd2bb2cc79f1284d7b6a403ad6a677a950024791b621e1f25bbac71e2bab8294ff38294a2c69500ed5f818ba1486f996b92cf02db32e4920bfc095f")
	data, err := raw.Value()
	require.NoError(t, err, "failed to decode raw")

	rawh := HexBytes("0xb10fc0dce4c066322dbca49cf76f162026ee5b632da2cb1e060503c398729a4b")
	hash, err := rawh.Value()
	require.NoError(t, err, "failed to decode rawh")

	h := crypto.SHA3Sum256(data)
	require.Equal(t, hash, h, "hash should match")

	fmt.Println(NewHexBytes(h))
}

func TestVerifierSufficientVotes(t *testing.T) {
	h := getSampleHeader()
	vr := NewSampleTestVerifier()
	cvl := getSampleCommitVoteList()
	cvl.Items = cvl.Items[:2]

	rawVotes, err := codec.BC.MarshalToBytes(cvl)
	require.NoError(t, err)

	ok, err := vr.Verify(h, rawVotes)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestVerifierInsufficientVotes(t *testing.T) {
	h := getSampleHeader()
	vr := NewSampleTestVerifier()
	cvl := getSampleCommitVoteList()
	cvl.Items = cvl.Items[:1]

	rawVotes, err := codec.BC.MarshalToBytes(cvl)
	require.NoError(t, err)

	ok, err := vr.Verify(h, rawVotes)
	require.EqualError(t, err, "insufficient votes")
	require.False(t, ok)
}

func TestVerifierInvalidValidatorVotes(t *testing.T) {
	h := getSampleHeader()
	vr := NewSampleTestVerifier()
	cvl := getSampleCommitVoteList()
	cvl.Items = append(
		cvl.Items[:1],
		getCommitVoteItem(1652523324898246, "L4NkrE96T9Bf8wsb5xvqpOVLkCFbgFIjKGl3W66AUJQyKra6QDhRLH37XB2ckLrVJ75LbIv1e+eGRLxFqyG0VAE="),
		getCommitVoteItem(1652523324866924, "3lvwtyCNNX+w1X+3E/N8Hu1rqEyoHgJFH1uB5XOgln9iodU7OjXw4pnyxmln4rdje/0icYgmyTPwgdwKmbk1iAE="),
	)

	rawVotes, err := codec.BC.MarshalToBytes(cvl)
	require.NoError(t, err)

	ok, err := vr.Verify(h, rawVotes)
	require.EqualError(t, err, "insufficient votes")
	require.False(t, ok)
}

func TestVerifierDuplicateVotes(t *testing.T) {
	h := getSampleHeader()
	vr := NewSampleTestVerifier()
	cvl := getSampleCommitVoteList()
	cvl.Items = append(cvl.Items[:1], cvl.Items[:1]...)

	rawVotes, err := codec.BC.MarshalToBytes(cvl)
	require.NoError(t, err)

	ok, err := vr.Verify(h, rawVotes)
	require.EqualError(t, err, "insufficient votes")
	require.False(t, ok)
}

func TestVerifierMinimumRequiredValidators(t *testing.T) {
	h := getSampleHeader()
	vr := NewSampleTestVerifier()
	vr.validators[vr.nextValidatorsHash.String()] = vr.validators[vr.nextValidatorsHash.String()][:1]
	cvl := getSampleCommitVoteList()
	cvl.Items = cvl.Items[:0]

	rawVotes, err := codec.BC.MarshalToBytes(cvl)
	require.NoError(t, err)

	ok, err := vr.Verify(h, rawVotes)
	require.EqualError(t, err, "insufficient votes")
	require.False(t, ok)
}
