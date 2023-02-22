package algo

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"
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
// Return the name of the abi call implementing this service and its input parameters
func DecodeRelayMessage(rlpMsg string) (string, interface{}, error) {
	bmcMsg := BMCMessage{}
	if err := RlpDecodeHex(rlpMsg, &bmcMsg); err != nil {
		err = errors.Wrapf(err, "Failed to decode bmc msg: %v", err)
		return "", nil, err
	}
	//msgSn := (&big.Int{}).SetBytes(bmcMsg.Sn)
	svcMessage := ServiceMessage{}
	if err := RlpDecodeHex(hex.EncodeToString(bmcMsg.Message), &svcMessage); err != nil {
		err = errors.Wrapf(err, "Failed to decode service message: %v", err)
		return "", nil, err
	}

	if bmcMsg.Svc == "bmc" {
		svcTypeStr := bmcService(string(svcMessage.ServiceType))
		switch svcTypeStr {
		case FEE_GATHERING:
			svcData := FeeGatheringSvc{}
			if err := RlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcData); err != nil {
				err = errors.Wrapf(err, "Failed to decode fee gathering: %v", err)
				return "", nil, err
			}
			return "FeeGathering", svcData, nil
		case LINK, UNLINK:
			svcArgs := LinkSvc{}
			if err := RlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcArgs); err != nil {
				err = errors.Wrapf(err, "Failed to decode link: %v", err)
				return "", nil, err
			}
			return "HandleLink", svcArgs, nil
		case INIT:
			svcArgs := InitSvc{}
			if err := RlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcArgs); err != nil {
				err = errors.Wrapf(err, "Failed to decode init: %v", err)
				return "", nil, err
			}
			return "Init", svcArgs, nil
		default:
			return "", nil, fmt.Errorf("Unexpected bmc service: %v", svcTypeStr)
		}
	}

	if bmcMsg.Svc == "bts" {
		svcTypeNum := btsService((&big.Int{}).SetBytes(svcMessage.ServiceType).Int64())
		switch svcTypeNum {
		case REQUEST_COIN_TRANSFER:
			svcArgs := CoinTransferSvc{}
			if err := RlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcArgs); err != nil {
				err = errors.Wrapf(err, "Failed to decode coin transfer: %v", err)
				return "", nil, err
			}
			return "CoinTransfer", svcArgs, nil
		case BLACKLIST_MESSAGE:
			svcArgs := BlacklistSvc{}
			if err := RlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcArgs); err != nil {
				err = errors.Wrapf(err, "Failed to decode blacklist: %v", err)
				return "", nil, err
			}
			requestType := blacklistSvc((&big.Int{}).SetBytes(svcArgs.RequestType).Int64())
			if requestType == ADD_TO_BLACKLIST {
				return "BlacklistAdd", svcArgs, nil
			} else if requestType == REMOVE_FROM_BLACKLIST {
				return "BlacklistRemove", svcArgs, nil
			} else {
				return "", nil, fmt.Errorf("Unknown blacklist request type: %v", requestType)
			}
		case CHANGE_TOKEN_LIMIT:
			svcArgs := TokenLimitSvc{}
			if err := RlpDecodeHex(hex.EncodeToString(svcMessage.Payload), &svcArgs); err != nil {
				err = errors.Wrapf(err, "Failed to decode token limit: %v", err)
				return "", nil, err
			}
			return "ChangeTokenLimit", svcArgs, nil
		default:
			return "", nil, fmt.Errorf("Unexpected bts service: %v", svcTypeNum)
		}
	}
	return "", nil, fmt.Errorf("Service not address to either BMC or BTS, but %v", bmcMsg.Svc)
}
