package main

import (
	"context"
	"fmt"
	"os"

	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
	"github.com/echa/log"
	"github.com/joho/godotenv"
)

const (
	tzZeroAddress = "tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"
)

func main() {
	rpc.UseLogger(log.Log)

	err := godotenv.Load(".env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := rpc.NewClient("https://ghostnet.tezos.marigold.dev", nil)
	fmt.Println("new client")

	fmt.Println(c.ChainId)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Init(ctx)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.Listen()

	// deployment options
	opts := rpc.DefaultOptions
	opts.Signer = signer.NewFromKey(tezos.MustParsePrivateKey(os.Getenv("secret_deployer")))
	fmt.Println(os.Getenv("secret_deployer"))

	bmc_periphery := os.Getenv("BMC_PERIPHERY")
	bmc_management := os.Getenv("BMC_MANAGEMENT")
	btsCore := os.Getenv("BTS_CORE")
	btsPeriphery := os.Getenv("BTS_PERIPHERY")
	prim := micheline.Prim{}

	// bmc_periphery

	contractAddress := tezos.MustParseAddress(bmc_periphery)
	bmcPeripheryClient := contract.NewContract(contractAddress, c)

	in := "{\"string\": \"" + bmc_management + "\" }"
	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args := contract.NewTxArgs()

	entrypoint := "set_bmc_management_addr"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	from := tezos.MustParseAddress("tz1ZPVxKiybvbV1GvELRJJpyE1xj1UpNpXMv")

	argument := args.WithSource(from).WithDestination(contractAddress)

	fmt.Println("setting bmc management address in periphery....")

	res, err := bmcPeripheryClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}
	// *************************************************************************************************************************************
	// bts periphery

	contractAddress = tezos.MustParseAddress(btsPeriphery)

	btsPeripheryClient := contract.NewContract(contractAddress, c)

	in = "{\"string\": \"" + bmc_periphery + "\" }"
	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "set_bmc_address"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	from = tezos.MustParseAddress("tz1ZPVxKiybvbV1GvELRJJpyE1xj1UpNpXMv")

	argument = args.WithSource(from).WithDestination(contractAddress)

	fmt.Println("setting bmc periphery in bts core....")

	res, err = btsPeripheryClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	// *************************************************************************************************************************************
	// bts periphery

	in = "{\"string\": \"" + btsCore + "\" }"
	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "set_bts_core_address"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	from = tezos.MustParseAddress("tz1ZPVxKiybvbV1GvELRJJpyE1xj1UpNpXMv")

	argument = args.WithSource(from).WithDestination(contractAddress)

	fmt.Println("setting setting bts core in bts periphery....")

	res, err = btsPeripheryClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	// *************************************************************************************************************************************
	// bmc management

	contractAddress = tezos.MustParseAddress(bmc_management)

	bmcManagementClient := contract.NewContract(contractAddress, c)

	in = "{\"string\": \"" + bmc_periphery + "\" }"
	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "set_bmc_periphery"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	from = tezos.MustParseAddress("tz1ZPVxKiybvbV1GvELRJJpyE1xj1UpNpXMv")

	argument = args.WithSource(from).WithDestination(contractAddress)

	fmt.Println("setting bmc periphery in bmc management....")

	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	// *************************************************************************************************************************************
	// set btp address

	prim = micheline.Prim{}

	in = "{ \"string\": \"" + os.Getenv("TZ_NETWORK") + "\" }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	fmt.Println("setting bmcBTP address in bmcManagement...")

	args = contract.NewTxArgs()

	entrypoint = "set_bmc_btp_address"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(contractAddress)

	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	// update bts periphery

	contractAddress = tezos.MustParseAddress(btsCore)
	btsCoreClient := contract.NewContract(contractAddress, c)

	prim = micheline.Prim{}

	in = "{ \"string\": \"" + btsPeriphery + "\" }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	fmt.Println("setting bts periphery in btsCoreClient...")

	args = contract.NewTxArgs()

	entrypoint = "update_bts_periphery"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(contractAddress)

	res, err = btsCoreClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	// add service

	prim = micheline.Prim{}

	in = "{ \"prim\": \"Pair\", \"args\": [ { \"string\": \"" + btsPeriphery + "\" }, { \"string\": \"bts\" } ] }"

	contractAddress = tezos.MustParseAddress(bmc_management)
	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "add_service"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(contractAddress)
	fmt.Println("adding service...")

	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	// set fee ratio

	prim = micheline.Prim{}
	contractAddress = tezos.MustParseAddress(btsCore)

	in = "{ \"prim\": \"Pair\", \"args\": [ { \"int\": \"100\" }, { \"prim\": \"Pair\", \"args\": [ { \"int\": \"450\" }, { \"string\": \"" + os.Getenv("TZ_NATIVE_COIN_NAME") + "\" } ] } ] }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "set_fee_ratio"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(contractAddress)

	fmt.Println("setting fee ratio...")
	res, err = btsCoreClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	// add route

	prim = micheline.Prim{}

	link := "btp://" + os.Getenv("ICON_NETWORK") + "/" + os.Getenv("ICON_BMC")
	fmt.Println(link)

	in = "{ \"prim\": \"Pair\", \"args\": [ { \"string\": \"" + link + "\" }, { \"string\": \"" + link + "\" } ] }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "add_route"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(bmcManagementClient.Address())

	fmt.Println("adding route...")
	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	//  add link

	prim = micheline.Prim{}

	fmt.Println(link)

	in = "{ \"string\": \"" + link + "\" }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "add_link"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(bmcManagementClient.Address())

	fmt.Println("adding link....")

	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	// set link rx height

	prim = micheline.Prim{}

	fmt.Println(os.Getenv("ICON_RX_HEIGHT"))

	in = "{ \"prim\": \"Pair\", \"args\": [ { \"int\": \"" + os.Getenv("ICON_RX_HEIGHT") + "\" }, { \"string\": \"" + link + "\" } ] }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "set_link_rx_height"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(bmcManagementClient.Address())

	fmt.Println("setting link_rx_height...")
	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)

	//***********************************************************************************************************************************
	// add relay

	prim = micheline.Prim{}

	fmt.Println(os.Getenv("RELAYER_ADDRESS"))

	in = "{ \"prim\": \"Pair\", \"args\": [ [ { \"string\": \"" + os.Getenv("RELAYER_ADDRESS") + "\" } ], { \"string\": \"" + link + "\" } ] }"

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return
	}

	args = contract.NewTxArgs()

	entrypoint = "add_relay"

	args.WithParameters(micheline.Parameters{Entrypoint: entrypoint, Value: prim})

	argument = args.WithSource(from).WithDestination(bmcManagementClient.Address())

	fmt.Println("adding relay...")
	res, err = bmcManagementClient.Call(ctx, argument, &opts)

	if err != nil {
		fmt.Println("error while calling")
		fmt.Println(err)
	}

	fmt.Println(res)
}
