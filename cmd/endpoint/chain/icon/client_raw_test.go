package icon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/icon-project/icon-bridge/common/crypto"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	btp_icon_demo_wallet_address = "hx51ecae93216cb6d58bbdc51e2b5d790da94f738a"
	btp_icon_god_wallet_address  = "hxff0ea998b84ab9955157ab27915a9dc1805edd35"
	btp_icon_step_limit          = 3500000000
	btp_icon_nid                 = 0x5b9a77
	Version                      = 3
)

func getDemoIconWallet(walFile string, password string) module.Wallet {
	btp_icon_sender_wallet := walFile
	btp_icon_sender_wallet_password := password
	keyReader, err := os.Open(btp_icon_sender_wallet)
	if err != nil {
		log.Fatal(err)
	}
	defer keyReader.Close()

	keyStore, err := ioutil.ReadAll(keyReader)
	if err != nil {
		log.Fatal(err)
	}
	w, err := wallet.NewFromKeyStore(keyStore, []byte(btp_icon_sender_wallet_password))
	if err != nil {
		log.Fatal(err)
	}
	return w
}

func TestGetIconBalance(t *testing.T) {
	const URL = "http://127.0.0.1:9080/api/v3/default"
	l := log.New()
	log.SetGlobalLogger(l)
	walletAddress := "hx51ecae93216cb6d58bbdc51e2b5d790da94f738a"
	cMap := map[string]string{
		"btp_icon_irc2":           "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f",
		"btp_icon_irc2_tradeable": "cx7831ba8969943d96c375261f8245e6e9964389c9",
		"btp_icon_nativecoin_bsh": "cx6fa46cc92fcf8e3135bc6645ef3a258b0ce73602",
		"btp_icon_token_bsh":      "cx2194782b7951d8abf26ea88204e45862e9821bbc",
	}
	cAddr := &contractAddress{}
	cAddr.FromMap(cMap)
	client, _ := newClient(URL, l, cAddr, "0x5b9a77")
	if val, err := client.GetBalance(&AddressParam{Address: Address(walletAddress)}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(val)
	}
	//curl -X POST 'http://127.0.0.1:9080/api/v3/default' -H 'Content-Type:application/json' -d '{"id":"1001", "jsonrpc":"2.0", "method": "icx_getBalance", "params":{"address":"hxff0ea998b84ab9955157ab27915a9dc1805edd35"} }'
}

func TestIconTransfer(t *testing.T) {
	/*
			goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" sendtx transfer --to "hx267ed8d02bae84ada9f6ab486d4557aa4763b33a" --value "20" --key_store devnet/docker/icon-hmny/src/icon.god.wallet.json --key_password "gochain" --nid "6003319" --step_limit "3500000000"
		"0x0ee02a03a85faead44393212d1b77df05c467f855ffab7b723bb97efa1560338"
	*/
	cl := newLocalClient()
	rpcWallet := getDemoIconWallet("/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.god.wallet.json", "gochain")

	param := TransactionParam{
		Version:     NewHexInt(Version),
		ToAddress:   Address(btp_icon_demo_wallet_address),
		Value:       NewHexInt(100000000000),
		FromAddress: Address(rpcWallet.Address().String()),
		StepLimit:   NewHexInt(btp_icon_step_limit),
		Timestamp:   NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   NewHexInt(btp_icon_nid),
	}
	js, err := json.Marshal(param)
	if err != nil {
		log.Fatal(err)
	}
	var txSerializeExcludes = map[string]bool{"signature": true}
	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		log.Fatal(err)
	}
	bs = append([]byte("icx_sendTransaction."), bs...)
	sig, err := rpcWallet.Sign(crypto.SHA3Sum256(bs))
	if err != nil {
		log.Fatal(err)
	}

	param.Signature = base64.StdEncoding.EncodeToString(sig)
	txHash, txr, err := cl.SendTransactionAndGetResult(&param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TxHash ", txHash, " Status ", txr.Status, " Failure ", txr.Failure)
}

func TestGetIconIrc2Balance(t *testing.T) {
	/*
		balance=$(icon_callsc "$btp_icon_irc2" balanceOf "_owner=$btp_icon_demo_wallet_address" | jq -r .)
		goloop rpc --uri "http://127.0.0.1:9080/api/v3/default" call --to "cxf559e2ab2d3a69d8b1c0f1c44f1a2c45bdc4424f" --method "balanceOf" --param _owner=hx51ecae93216cb6d58bbdc51e2b5d790da94f738a
			"0x130e5dd8e6a6dad60000
	*/
	cl := newLocalClient()
	param := &CallParam{
		FromAddress: Address(btp_icon_god_wallet_address),
		ToAddress:   Address(cl.contractAddress.btp_icon_irc2),
		DataType:    "call",
	}
	argMap := map[string]interface{}{}
	argMap["method"] = "balanceOf"
	argMap["params"] = map[string]string{"_owner": btp_icon_demo_wallet_address}
	param.Data = argMap

	var res interface{}
	err := cl.Call(param, &res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TxHash %v", res.(string))
}

func TestTransIrc(t *testing.T) {
	/*
			bal=$(get_icon_irc2_balance)
		bal=$(echo "scale=18;$irc2_target-$bal" | bc)
		if (($(echo "$bal > 0" | bc -l))); then
			WALLET=$btp_icon_wallet \
				PASSWORD=$btp_icon_wallet_password \
				icon_sendtx_call >/dev/null \
				"$btp_icon_irc2" transfer 0 \
				"_to=$btp_icon_demo_wallet_address" \
				"_value=$bal"
	*/

	cl := newLocalClient()
	rpcWallet := getDemoIconWallet("/home/manish/go/src/work/icon-bridge/devnet/docker/icon-hmny/src/icon.god.wallet.json", "gochain")

	param := TransactionParam{
		Version:     NewHexInt(Version),
		ToAddress:   Address(cl.contractAddress.btp_icon_irc2),
		Value:       NewHexInt(0),
		FromAddress: Address(rpcWallet.Address().String()),
		StepLimit:   NewHexInt(btp_icon_step_limit),
		Timestamp:   NewHexInt(time.Now().UnixNano() / int64(time.Microsecond)),
		NetworkID:   NewHexInt(btp_icon_nid),
		DataType:    "call",
	}
	argMap := map[string]interface{}{}
	argMap["method"] = "transfer"
	argMap["params"] = map[string]string{"_to": btp_icon_demo_wallet_address, "_value": "1000000000"}
	param.Data = argMap

	js, err := json.Marshal(param)
	if err != nil {
		log.Fatal(err)
	}
	var txSerializeExcludes = map[string]bool{"signature": true}
	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		log.Fatal(err)
	}
	bs = append([]byte("icx_sendTransaction."), bs...)
	sig, err := rpcWallet.Sign(crypto.SHA3Sum256(bs))
	if err != nil {
		log.Fatal(err)
	}
	param.Signature = base64.StdEncoding.EncodeToString(sig)
	txHash, txr, err := cl.SendTransactionAndGetResult(&param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TxHash ", txHash, " Status ", txr.Status, " Failure ", txr.Failure)
}

/*
   i2h_nativecoin_transfer_amount=2000000000000000000 # 2 ICX
   echo "Transfer Native ICX (ICON -> HMNY):"
   echo "    amount=$(format_token $i2h_nativecoin_transfer_amount)"
   echo -n "    "
   WALLET=$btp_icon_demo_wallet \
       PASSWORD=$btp_icon_demo_wallet_password \
       run_exec iconTransferNativeCoin \
       $i2h_nativecoin_transfer_amount \
       "btp://$btp_hmny_net/$btp_hmny_demo_wallet_address" >/dev/null

*/

// }

/*btp_icon_nativecoin_bsh

icx_getBalance
btp_icon_demo_wallet_address (from btp_icon_demo_wallet.json

rootPFlags.String("key_store", "", "KeyStore file for wallet")
	rootPFlags.String("key_secret", "", "Secret(password) file for KeyStore")
rootPFlags.String("key_password", "", "Password for the KeyStore file")
rootPFlags.String("nid", "", "Network ID")
rootPFlags.Int64("step_limit", 0, "StepLimit")
	rootPFlags.Bool("wait", false, "Wait transaction result")
	rootPFlags.Int("wait_interval", 1000, "Polling interval(msec) for wait transaction result")
	rootPFlags.Int("wait_timeout", 10, "Timeout(sec) for wait transaction result")
	rootPFlags.Bool("estimate", false, "Just estimate steps for the tx")
	rootPFlags.String("save", "", "Store transaction to the file")

Fund demo wallets
icon_transfer -> icon_sendtx_transfer
	args: wallet, walletPassword, walletAddress, balance
	func: (get_icon_balance), icon_transfer
	        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx transfer \
            --to "$address" \
            --value "$value" \
            --key_store "$WALLET" \
            --key_password "$PASSWORD" \
            --nid "$btp_icon_nid" \
            --step_limit "$btp_icon_step_limit" | jq -r .

		    goloop rpc \
            --uri "$btp_icon_uri" \
            txresult "$tx_hash" &>/dev/null && break || sleep 1

	icx_getBalance
    curl "$btp_icon_uri" -s -X POST \
        -H 'Content-Type:application/json' \
        -d "$(jq <<<{} -c \
            '.id=1|.jsonrpc="2.0"|.method=$method|.params=$params' \
            --arg method "$1" --argjson params "$2")"

get_icon_irc2_balance -> icon_callsc
	args: wallet, password, , wallet_address, balance
	func: get_icon_icr2_balance, icon_sendtx_call, transfer

    goloop rpc \
        --uri "$btp_icon_uri" \
        call \
        --to "$address" \
        --method "$method" ${params[@]}

iconTransferNativeCoin -> icon_sendtx_call
        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx call \
            --to "$address" \
            --key_store "$WALLET" \
            --key_password "$PASSWORD" \
            --nid "$btp_icon_nid" \
            --step_limit "$btp_icon_step_limit" \
            --value "$value" \
            --method "$method" \
            ${params[@]} | jq -r .

iconBSHApprove
iconTransferWrappedCoin
*/
