package algo

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/icon-project/icon-bridge/common/errors"
)

func EncodeBlockHash(block *types.Block) [32]byte {
	blockHeaderBytes := msgpack.Encode(block.BlockHeader)
	toBeHashed := []byte("BH")
	toBeHashed = append(toBeHashed, blockHeaderBytes...)
	hash := sha512.Sum512_256(toBeHashed)
	return hash
}

func RlpDecodeHex(str string, out interface{}) error {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	input, err := hex.DecodeString(str)
	if err != nil {
		return errors.Wrap(err, "hex.DecodeString ")
	}
	err = rlp.Decode(bytes.NewReader(input), out)
	if err != nil {
		return errors.Wrap(err, "rlp.Decode ")
	}
	return nil
}

func RlpEncodeHex(in interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := rlp.Encode(&buf, in)
	if err != nil {
		return nil, errors.Wrap(err, "rlp.Encode ")
	}
	return buf.Bytes(), nil
}

// Decode receiving rlp encoded msg to identify which service is it requesting
func DecodeRelayMessage(rlpMsg string) (string, interface{}, error) {
	bmcMsg := BMCMessage{}
	if err := RlpDecodeHex(rlpMsg, &bmcMsg); err != nil {
		err = errors.Wrapf(err, "Failed to decode bmc msg: %v", err)
		return "", nil, err
	}
	return bmcMsg.Svc, bmcMsg.Message, nil
}
