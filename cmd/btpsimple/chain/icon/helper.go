package icon

import (
	"bytes"
	"io"

	"github.com/gorilla/websocket"
	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/common/trie/ompt"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/pkg/errors"
)

func getBlockHeader(cl *client, height HexInt) (*BlockHeader, error) {
	p := &BlockHeightParam{Height: height}
	b, err := cl.GetBlockHeaderByHeight(p)
	if err != nil {
		return nil, mapError(err)
	}
	var bh BlockHeader
	_, err = codec.RLP.UnmarshalFromBytes(b, &bh)
	if err != nil {
		return nil, err
	}
	bh.serialized = b
	return &bh, nil
}

func getValidatorsFromHash(cl *client, hash []byte) ([][]byte, error) {
	vBytes, err := cl.GetDataByHash(&DataHashParam{Hash: NewHexBytes(hash)})
	if err != nil {
		return nil, errors.Wrap(err, "verifyHeader; GetDataByHash Validators; Err: ")
	}
	var vs [][]byte
	_, err = vlcodec.BC.UnmarshalFromBytes(vBytes, &vs)
	if err != nil {
		return nil, errors.Wrap(err, "verifyHeader; Unmarshal Validators; Err: ")
	}
	return vs, nil
}

func mptProve(key HexInt, proofs [][]byte, hash []byte) ([]byte, error) {
	db := db.NewMapDB()
	defer db.Close()
	index, err := key.Value()
	if err != nil {
		return nil, err
	}
	indexKey, err := vlcodec.RLP.MarshalToBytes(index)
	if err != nil {
		return nil, err
	}
	mpt := ompt.NewMPTForBytes(db, hash)
	trie, err1 := mpt.Prove(indexKey, proofs)
	if err1 != nil {
		return nil, err1

	}
	return trie, nil
}

func addressesContains(data []byte, list [][]byte) bool {
	for _, current := range list {
		if bytes.Equal(data, current) {
			return true
		}
	}
	return false
}

// Websocket connection is closed by peer abruptly with EOF message. The function checks and verifies if the error thrown is unexpected EOF
func isUnexpectedEOFError(err error) bool {
	//websocket/conn.go 	errUnexpectedEOF       = &CloseError{Code: CloseAbnormalClosure, Text: io.ErrUnexpectedEOF.Error()}
	if cErr, ok := err.(*websocket.CloseError); ok && cErr.Code == websocket.CloseAbnormalClosure && cErr.Text == io.ErrUnexpectedEOF.Error() {
		return true
	} else if err.Error() == io.ErrUnexpectedEOF.Error() {
		return true
	}
	return false
}
