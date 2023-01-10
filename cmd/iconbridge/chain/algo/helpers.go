package algo

import (
	"crypto/sha512"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

func BlockHash(block *types.Block) [32]byte {
	blockHeaderBytes := msgpack.Encode(block.BlockHeader)
	toBeHashed := []byte("BH")
	toBeHashed = append(toBeHashed, blockHeaderBytes...)
	hash := sha512.Sum512_256(toBeHashed)
	return hash
}
