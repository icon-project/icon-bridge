package hmny

/*
import (
	"context"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestHmnyReceiver(t *testing.T) {
	const (
		src = "btp://0x6357d2e0.hmny/0x0169AE3f21b67e798fd4AdF50d0FA9FB83d72651"
		dst = "btp://0x5b9a77.icon/cx8015df5623344958af75ba1598b4ee14b8574bee"
		url = "http://localhost:9500"
	)

	var opts map[string]interface{} = map[string]interface{}{
		// "options": map[string]interface{}{
		"syncConcurrency": 100,
		"verifier": map[string]interface{}{
			"blockHeight":     29075, //1171,
			"commitBitmap":    "0xff",
			"commitSignature": "0x44e933d799e2e44f5f7d84fc2f8400b429337a4779271798bb45b6e07e94d311e6de89272f3a41ba766316b3a850b20701c207926c263a8d009264629a22412d1b8e6ad7905444b4da81c81362abc5dca42b970d0233e0ed980dd79677a96204",
		},
		// },
	}

	l := log.New()
	log.SetGlobalLogger(l)

	rx, err := NewReceiver(chain.BTPAddress(src), chain.BTPAddress(dst), []string{url}, opts, l)
	if err != nil {
		log.Fatal((err))
	}

	srcMsgCh := make(chan *uint64)
	srcErrCh, err := rx.Subscribe(context.TODO(),
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    138,
			Height: 29076,
		})
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-srcErrCh:
			log.Fatal(err)

		case <-srcMsgCh:
		}
	}
}



	// const (
	// 	src = "btp://0x6357d2e0.hmny/0x0169AE3f21b67e798fd4AdF50d0FA9FB83d72651"
	// 	dst = "btp://0x5b9a77.icon/cx8015df5623344958af75ba1598b4ee14b8574bee"
	// 	url = "http://localhost:9500"
	// )

	// var opts map[string]interface{} = map[string]interface{}{
	// 	// "options": map[string]interface{}{
	// 	"syncConcurrency": 100,
	// 	"verifier": map[string]interface{}{
	// 		"blockHeight":     29075, //1171,
	// 		"commitBitmap":    "0xff",
	// 		"commitSignature": "0x44e933d799e2e44f5f7d84fc2f8400b429337a4779271798bb45b6e07e94d311e6de89272f3a41ba766316b3a850b20701c207926c263a8d009264629a22412d1b8e6ad7905444b4da81c81362abc5dca42b970d0233e0ed980dd79677a96204",
	// 	},
	// 	// },
	// }

	// l := log.New()
	// log.SetGlobalLogger(l)


	func (r *receiver) getFilteredReceipts(v *BlockNotification) []*LogResult {
	const (
		TransferStartSignature         = "0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a" //"TransferStart(address,string,uint256,(string,uint256,uint256)[])" //
		TransferEndSignature           = "0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2" //"TransferEnd(address,uint256,uint256,string)"                      //
		TransferReceivedSignature      = "0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680" //"TransferReceived(string,address,uint256,(string,uint256)[])"      //
		TransferReceivedSignatureToken = "0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c" //"TransferReceived(string,address,uint256,(string,uint256,uint256)[])" //
	)
	signatureMap := map[string]string{
		"0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a": "TransferStart",
		"0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2": "TransferEnd",
		"0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680": "TransferReceived",
		"0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c": "TransferReceived",
	}

	newResults := []*LogResult{}
	for _, receipt := range v.Receipts {
		for _, log := range receipt.Logs {
			var newTopic *common.Hash
			for _, topic := range log.Topics {
				if topic == common.HexToHash(TransferStartSignature) ||
					topic == common.HexToHash(TransferReceivedSignature) ||
					topic == common.HexToHash(TransferReceivedSignatureToken) ||
					topic == common.HexToHash(TransferEndSignature) {
					newTopic = &topic
					break
				}
			}
			if newTopic != nil {

				if res, err := decodeLogData(log.Data, signatureMap[newTopic.String()]); err == nil && res != nil {
					newResults = append(newResults, &LogResult{
						TxHash:   log.TxHash,
						LogIndex: log.Index,
						Address:  log.Address,
						Topic:    signatureMap[newTopic.String()],
						Logs:     res,
					})
				} else if err != nil {
					r.log.Error(err)
				} else if res == nil {
					r.log.Error("Returned nil interface")
				}
			}
		}
	}
	for i, r := range newResults {
		fmt.Println("New ", i, "  ", *r)
	}
	return newResults
}

func decodeLogData(data []byte, topicType string) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty Log Data input to decode")
	}

	abi, err := abi.JSON(strings.NewReader(bshPeripherABI))
	if err != nil {
		return nil, err
	}

	if topicType == "TransferStart" {
		var ev TransferStart
		err = abi.UnpackIntoInterface(&ev, topicType, data)
		if err != nil {
			return nil, err
		}
		return ev, nil
	} else if topicType == "TransferEnd" {
		var ev TransferEnd
		err = abi.UnpackIntoInterface(&ev, topicType, data)
		if err != nil {
			return nil, err
		}
		return ev, nil
	} else if topicType == "TransferReceived" {
		var ev TransferReceived
		err = abi.UnpackIntoInterface(&ev, topicType, data)
		if err != nil {
			return nil, err
		}
		return ev, nil
	}
	return nil, errors.New("Doesn't match any signature")
}

*/
