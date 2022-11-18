package tests

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/tests/mock"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
)

type VerifierValidateHeaderTest struct {
	description string
	testData    []TestData
}

func (t VerifierValidateHeaderTest) Description() string {
	return t.description
}

func (t VerifierValidateHeaderTest) TestDatas() []TestData {
	return t.testData
}

func init() {
	var testData = []TestData{
		{
			Description: "Valid BlockHeader",
			Input: struct {
				Offset  uint64
				Options types.ReceiverOptions
			}{
				Offset: 377826,
				Options: types.ReceiverOptions{
					Verifier: &types.VerifierConfig{
						BlockHeight:       377825,
						PreviousBlockHash: types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						CurrentBpsHash:    types.NewCryptoHash("C4zVnMf27hRJYoWEC816Pttyz122TWZN7zjUMoZCNkuw"),
						CurrentEpochId:    types.NewCryptoHash("FtrJuAXqH5oXDVADh6QkUyacf2MGmLHYbHCHKSZ8C7KS"),
						NextEpochId:       types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						NextBpsHash:       types.NewCryptoHash("5QouG4ceHjyARjVTaySWXcXdsQduDqExVRKdwLjeANi"),
					},
				},
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
				})

				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				mockApi.On("BlockProducers", mock.MockParam).Return(mockApi.BlockProducersFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: true,
			},
		},
		{
			Description: "Invalid previous hash",
			Input: struct {
				Offset  uint64
				Options types.ReceiverOptions
			}{
				Offset: 377826,
				Options: types.ReceiverOptions{
					Verifier: &types.VerifierConfig{
						BlockHeight:       377825,
						PreviousBlockHash: types.NewCryptoHash("74toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						CurrentBpsHash:    types.NewCryptoHash("G2TyLP33XfqndppUzipoTWTs6XnKjmUhCQg1tH44isAG"),
						CurrentEpochId:    types.NewCryptoHash("FtrJuAXqH5oXDVADh6QkUyacf2MGmLHYbHCHKSZ8C7KS"),
						NextEpochId:       types.NewCryptoHash("84toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95"),
						NextBpsHash:       types.NewCryptoHash("5QouG4ceHjyARjVTaySWXcXdsQduDqExVRKdwLjeANi"),
					},
				},
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"377825", "377826", "377827", "377828", "377829", "377830", "377831"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
					BlockProducersMap: map[string]mock.Response{
						"74toXNMo2p5ttdjkV6RHdJFrgxrnTLRkCTjb7aA8Dh95": {
							Reponse: []byte(`[]`),
						},
					},
				})

				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				mockApi.On("BlockProducers", mock.MockParam).Return(mockApi.BlockProducersFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())
				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Fail: func(err error) bool {
					return err.Error() == string("expected hash: 5YqjrSoiQjqrmrMHQJVZB25at7yQ2BZEC2exweLFmc6w, got hash: E78zFTZ21jN4nFaiygYqZNLw8iKMSYxrTNFG8X7eAJhE for block: 377826")
				},
			},
		},
		{
			Description: "Block from new epoch",
			Input: struct {
				Offset  uint64
				Options types.ReceiverOptions
			}{
				Offset: 102814341,
				Options: types.ReceiverOptions{
					Verifier: &types.VerifierConfig{
						BlockHeight:       102814335,
						PreviousBlockHash: types.NewCryptoHash("7XVimvTroNs9eTncfWruSURER1Uy3C1CaphZVD2eZPLX"),
						CurrentBpsHash:    types.NewCryptoHash("4wsFrLRD9qCBT6jng6k8b3R6JTNBY7qRrFv9piZMPmmw"),
						CurrentEpochId:    types.NewCryptoHash("FtrJuAXqH5oXDVADh6QkUyacf2MGmLHYbHCHKSZ8C7KS"),
						NextEpochId:       types.NewCryptoHash("E5PsDpHsAG5b6pYvbs9HxzzCmaXyG3MaWknnzZadpcrZ"),
						NextBpsHash:       types.NewCryptoHash("F5s6UcitJ6uq3PXsbJPLpFTgMLnoBQdmMqoR91GyMMSt"),
					},
				},
			},
			MockApi: func() *mock.MockApi {
				blockByHeightMap, blockByHashMap := mock.LoadBlockFromFile([]string{"102814335", "102814336", "102814337", "102814338", "102814339", "102814340", "102814341"})

				mockApi := mock.NewMockApi(mock.Storage{
					BlockByHeightMap: blockByHeightMap,
					BlockByHashMap:   blockByHashMap,
					BlockProducersMap: map[string]mock.Response{
						"7XVimvTroNs9eTncfWruSURER1Uy3C1CaphZVD2eZPLX": {
							Reponse: []byte(`[{
								"account_id": "node0",
								"public_key": "ed25519:7PGseFbWxvYVgZ89K1uTJKYoKetWs7BJtbyXDzfbAcqX",
								"stake": "32049497097804540271330918142596",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "node2",
								"public_key": "ed25519:GkDv7nSMS3xcqA45cpMvFmfV1o4fRF6zYo1JRR6mNqg5",
								"stake": "32046892864435253988970887269589",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "node1",
								"public_key": "ed25519:6DSjZ8mvsRZDvFqFxo8tCKePG96omXW7eVYVSySmDk8e",
								"stake": "32039371460921428273767465143706",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "node3",
								"public_key": "ed25519:ydgzeXHJ5Xyt7M1gXLxqLBW1Ejx6scNV5Nx2pxFM8su",
								"stake": "32034292776141323339281285108484",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "aurora.pool.f863973.m0",
								"public_key": "ed25519:9c7mczZpNzJz98V1sDeGybfD4gMybP4JKHotH8RrrHTm",
								"stake": "13334014284409172181833745110510",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "01node.pool.f863973.m0",
								"public_key": "ed25519:3iNqnvBgxJPXCxu6hNdvJso1PEAc1miAD35KQMBCA3aL",
								"stake": "9479648576824608411094559666801",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "legends.pool.f863973.m0",
								"public_key": "ed25519:AhQ6sUifJYgjqarXSAzdDZU9ZixpUesP9JEH1Vr7NbaF",
								"stake": "7828123351461489284989450457543",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "everstake.pool.f863973.m0",
								"public_key": "ed25519:4LDN8tZUTRRc4siGmYCPA67tRyxStACDchdGDZYKdFsw",
								"stake": "7400314092316390016929343987164",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "spectrum.pool.f863973.m0",
								"public_key": "ed25519:ASecMN9e28vtCJn7rD2noNwL5c3odzQgAfbfHrUnbSVe",
								"stake": "7360690520495033742328707927957",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "nodeasy.pool.f863973.m0",
								"public_key": "ed25519:25Dhg8NBvQhsVTuugav3t1To1X1zKiomDmnh8yN9hHMb",
								"stake": "7197442933924577829850251760816",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "chorusone.pool.f863973.m0",
								"public_key": "ed25519:3TkUuDpzrq75KtJhkuLfNNJBPHR5QEWpDxrter3znwto",
								"stake": "6920450831929719424138114239051",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "ni.pool.f863973.m0",
								"public_key": "ed25519:GfCfFkLk2twbAWdsS3tr7C2eaiHN3znSfbshS5e8NqBS",
								"stake": "6644296385466072550043169727944",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "stakely_v2.pool.f863973.m0",
								"public_key": "ed25519:7BanKZKGvFjK5Yy83gfJ71vPhqRwsDDyVHrV2FMJCUWr",
								"stake": "6262871152843303923708119657452",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "masternode24.pool.f863973.m0",
								"public_key": "ed25519:9E3JvrQN6VGDGg1WJ3TjBsNyfmrU6kncBcDvvJLj6qHr",
								"stake": "4693521833962415518524139697922",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "staked.pool.f863973.m0",
								"public_key": "ed25519:D2afKYVaKQ1LGiWbMAZRfkKLgqimTR74wvtESvjx5Ft2",
								"stake": "4678102913440886171585992595514",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "foundryusa.pool.f863973.m0",
								"public_key": "ed25519:ABGnMW8c87ZKWxvZLLWgvrNe72HN7UoSf4cTBxCHbEE5",
								"stake": "1706282977165769959670548491994",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "tribe-pool.pool.f863973.m0",
								"public_key": "ed25519:CRS4HTSAeiP8FKD3c3ZrCL5pC92Mu1LQaWj22keThwFY",
								"stake": "1616532585015444534688348119597",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "sweden.pool.f863973.m0",
								"public_key": "ed25519:2RVUnsMEZhGCj1A3vLZBGjj3i9SQ2L46Z1Z41aEgBzXg",
								"stake": "1546703609118387752645151521610",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "chorus-one.pool.f863973.m0",
								"public_key": "ed25519:6LFwyEEsqhuDxorWfsKcPPs324zLWTaoqk4o6RDXN7Qc",
								"stake": "1541166959717019871695690779347",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "lunanova2.pool.f863973.m0",
								"public_key": "ed25519:9Jv6e9Kye4wM9EL1XJvXY8CYsLi1HLdRKnTzXBQY44w9",
								"stake": "1534453244469297292666471744528",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "hotones.pool.f863973.m0",
								"public_key": "ed25519:2fc5xtbafKiLtxHskoPL2x7BpijxSZcwcAjzXceaxxWt",
								"stake": "1315107169518961469863675410025",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "pathrocknetwork.pool.f863973.m0",
								"public_key": "ed25519:CGzLGZEMb84nRSRZ7Au1ETAoQyN7SQXQi55fYafXq736",
								"stake": "1055197978776919097890139330540",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "stakesstone.pool.f863973.m0",
								"public_key": "ed25519:3aAdsKUuzZbjW9hHnmLWFRKwXjmcxsnLNLfNL4gP1wJ8",
								"stake": "908655797584936045251281227138",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "leadnode.pool.f863973.m0",
								"public_key": "ed25519:CdP6CBFETfWYzrEedmpeqkR6rsJNeT22oUFn2mEDGk5i",
								"stake": "904995694434597259777870102019",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "blockscope.pool.f863973.m0",
								"public_key": "ed25519:6K6xRp88BCQX5pcyrfkXDU371awMAmdXQY4gsxgjKmZz",
								"stake": "895431226115882998649728009238",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "grassets.pool.f863973.m0",
								"public_key": "ed25519:3S4967Dt1VeeKrwBdTTR5tFEUFSwh17hEFLATRmtUNYV",
								"stake": "894540957158715177643365597822",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "al3c5.pool.f863973.m0",
								"public_key": "ed25519:BoYixTjyBePQ1VYP3s29rZfjtz1FLQ9og4FWZB5UgWCZ",
								"stake": "892989468567239925546401924678",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "shurik.pool.f863973.m0",
								"public_key": "ed25519:9zEn7DVpvQDxWdj5jSgrqJzqsLo8T9Wv37t83NXBiWi6",
								"stake": "860405997792048643155757037073",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "baziliknear.pool.f863973.m0",
								"public_key": "ed25519:9Rbzfkhkk6RSa1HoPnJXS4q2nn1DwYeB4HMfJBB4WQpU",
								"stake": "857133423027981457653220143522",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "optimusvalidatornetwork.pool.f863973.m0",
								"public_key": "ed25519:BGoxGmpvN7HdUSREQXfjH6kw5G6ph7NBXVfBVfUSH85V",
								"stake": "812700279269706342022593906327",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "chelovek_iz_naroda.pool.f863973.m0",
								"public_key": "ed25519:89aWsXXytjAZxyefXuGN73efnM9ugKTjPEGV4hDco8AZ",
								"stake": "805688434970125310262089629985",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "basilisk-stake.pool.f863973.m0",
								"public_key": "ed25519:CFo8vxoEUZoxbs87mGtG8qWUvSBHB91Vc6qWsaEXQ5cY",
								"stake": "782911221569156167328795553831",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "ou812.pool.f863973.m0",
								"public_key": "ed25519:2APjYBPnQ7CGDxFpsUHeAcyhpjZRWWpjciq2Pdzk31uQ",
								"stake": "631717642275210246129061598644",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "blazenet.pool.f863973.m0",
								"public_key": "ed25519:DiogP36wBXKFpFeqirrxN8G2Mq9vnakgBvgnHdL9CcN3",
								"stake": "473288657824744176090140658707",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "infstones.pool.f863973.m0",
								"public_key": "ed25519:BLP6HB8tcwYRTxswQ2YRaJ5sGj1dgGpUUfcNwbnWFGCU",
								"stake": "471885824702620439244621471143",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "prophet.pool.f863973.m0",
								"public_key": "ed25519:HYJ9mUhxLhzSVtbjj89smAaZkMqXca68iCumZy3gySoB",
								"stake": "355997620603810940411666064214",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "kiln.pool.f863973.m0",
								"public_key": "ed25519:Bq8fe1eUgDRexX2CYDMhMMQBiN13j8vTAVFyTNhEfh1W",
								"stake": "334456617923856893652278116518",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "idtcn3.pool.f863973.m0",
								"public_key": "ed25519:DtkY9WtkWweSrF13BJi5k4c6xyk3tBAC9y92AEY4Ayfb",
								"stake": "155652812877690310946020030531",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "pennyvalidators.pool.f863973.m0",
								"public_key": "ed25519:HiHdwq9rxi9hyxaGkazDHbYu4XL1j3J4TjgHQioyhEva",
								"stake": "133571031935075707692398860500",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "shardlabs.pool.f863973.m0",
								"public_key": "ed25519:DxmhGQZ6oqdxw7qGBvzLuBzE6XQjEh67hk5tt66vhLqL",
								"stake": "88870331168887059380490819947",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "gettingnear.pool.f863973.m0",
								"public_key": "ed25519:5QzHuNZ4stznMwf3xbDfYGUbjVt8w48q8hinDRmVx41z",
								"stake": "77707438897240102168396287062",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "pandateam.pool.f863973.m0",
								"public_key": "ed25519:7n426KJocZpJ5496UHp6onwYqWyt5xuiAyzvTGwCQLTN",
								"stake": "69577639281507536583579919468",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "sevennines-t0.pool.f863973.m0",
								"public_key": "ed25519:BHKMMc1t7F6B26BbaBVMhex3riPtWoGL5CricgXiSFC4",
								"stake": "63751488962259034471551570689",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "guardia.pool.f863973.m0",
								"public_key": "ed25519:2b5AQqcf8PHAUzqxWYoaFiTEa7QuEkpihzeoDittQaPL",
								"stake": "63133600280024637714801256889",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "redhead.pool.f863973.m0",
								"public_key": "ed25519:5qPvLhc86TDdof4YEjBMKrENzT3UA9mEonKWQXuFvaHX",
								"stake": "56730087909083846639756701318",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "tandemk.pool.f863973.m0",
								"public_key": "ed25519:8zqx8dzqsxXivPMcSk2e7gezMvooZb1RWMGHA1FzFrak",
								"stake": "53340550512009977033302559677",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "azetsi.pool.f863973.m0",
								"public_key": "ed25519:2MFKLj9E2kRdJoQqUgaY9KtheebzLv9ntdgTGsZxLaE1",
								"stake": "50641772538973291262999333590",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "wackazong.pool.f863973.m0",
								"public_key": "ed25519:EK1bdY5F6prLush2aKnJBe5neHdK52wwTjg3qY4v3cjX",
								"stake": "47873777943373528239318206994",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "pero_val.pool.f863973.m0",
								"public_key": "ed25519:J43JCHe2XKU7wiDi7PSS8exsdVgnYHcbnC8R14SpF6HV",
								"stake": "45787360236563807456853169642",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "nodebull.pool.f863973.m0",
								"public_key": "ed25519:8tgk15x5XU15Ka9UCVCmAxeJweRrTbA96rnigs2ARQP4",
								"stake": "41723879196565485499078685560",
								"validator_stake_struct_version": "V1"
							},
							{
								"account_id": "lastnode.pool.f863973.m0",
								"public_key": "ed25519:811gesxXYdYeThry96ZiWn8chgWYNyreiScMkmxg4U9u",
								"stake": "38187423354513923800722592298",
								"validator_stake_struct_version": "V1"
							}]`),
						},
						"CPtZxDKtuNJGf77z7ror7FcnAJZrLosW4d1jWDosmVXd": {
							Reponse: []byte(`[
								{
									"account_id": "node0",
									"public_key": "ed25519:7PGseFbWxvYVgZ89K1uTJKYoKetWs7BJtbyXDzfbAcqX",
									"stake": "32071566353900138729788131511111",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "node2",
									"public_key": "ed25519:GkDv7nSMS3xcqA45cpMvFmfV1o4fRF6zYo1JRR6mNqg5",
									"stake": "32068960327257997522685217081709",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "node1",
									"public_key": "ed25519:6DSjZ8mvsRZDvFqFxo8tCKePG96omXW7eVYVSySmDk8e",
									"stake": "32061433744512190775097726128321",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "node3",
									"public_key": "ed25519:ydgzeXHJ5Xyt7M1gXLxqLBW1Ejx6scNV5Nx2pxFM8su",
									"stake": "32056351562554120803778313322153",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "cryptogarik.pool.f863973.m0",
									"public_key": "ed25519:FyFYc2MVwgitVf4NDLawxVoiwUZ1gYsxGesGPvaZcv6j",
									"stake": "14325379360348667914817212139793",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "aurora.pool.f863973.m0",
									"public_key": "ed25519:9c7mczZpNzJz98V1sDeGybfD4gMybP4JKHotH8RrrHTm",
									"stake": "13343197253368794794911171517610",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "01node.pool.f863973.m0",
									"public_key": "ed25519:3iNqnvBgxJPXCxu6hNdvJso1PEAc1miAD35KQMBCA3aL",
									"stake": "9486176254669280395190766006794",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "legends.pool.f863973.m0",
									"public_key": "ed25519:AhQ6sUifJYgjqarXSAzdDZU9ZixpUesP9JEH1Vr7NbaF",
									"stake": "7793409763418864380895788334147",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "everstake.pool.f863973.m0",
									"public_key": "ed25519:4LDN8tZUTRRc4siGmYCPA67tRyxStACDchdGDZYKdFsw",
									"stake": "7405411942356139348260594492054",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "spectrum.pool.f863973.m0",
									"public_key": "ed25519:ASecMN9e28vtCJn7rD2noNwL5c3odzQgAfbfHrUnbSVe",
									"stake": "7360690520701357304707707927957",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "nodeasy.pool.f863973.m0",
									"public_key": "ed25519:25Dhg8NBvQhsVTuugav3t1To1X1zKiomDmnh8yN9hHMb",
									"stake": "7197442934096787083677751760816",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "chorusone.pool.f863973.m0",
									"public_key": "ed25519:3TkUuDpzrq75KtJhkuLfNNJBPHR5QEWpDxrter3znwto",
									"stake": "6925216248411630981598526783849",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "ni.pool.f863973.m0",
									"public_key": "ed25519:GfCfFkLk2twbAWdsS3tr7C2eaiHN3znSfbshS5e8NqBS",
									"stake": "6644296386969014753665669727944",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "stakely_v2.pool.f863973.m0",
									"public_key": "ed25519:7BanKZKGvFjK5Yy83gfJ71vPhqRwsDDyVHrV2FMJCUWr",
									"stake": "6268202151882553952989099765599",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "masternode24.pool.f863973.m0",
									"public_key": "ed25519:9E3JvrQN6VGDGg1WJ3TjBsNyfmrU6kncBcDvvJLj6qHr",
									"stake": "4695709841767262026290452539198",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "staked.pool.f863973.m0",
									"public_key": "ed25519:D2afKYVaKQ1LGiWbMAZRfkKLgqimTR74wvtESvjx5Ft2",
									"stake": "4678102913671524075387092595514",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "foundryusa.pool.f863973.m0",
									"public_key": "ed25519:ABGnMW8c87ZKWxvZLLWgvrNe72HN7UoSf4cTBxCHbEE5",
									"stake": "1707457901496345240513560219745",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "tribe-pool.pool.f863973.m0",
									"public_key": "ed25519:CRS4HTSAeiP8FKD3c3ZrCL5pC92Mu1LQaWj22keThwFY",
									"stake": "1616532585178731827215848119597",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "chorus-one.pool.f863973.m0",
									"public_key": "ed25519:6LFwyEEsqhuDxorWfsKcPPs324zLWTaoqk4o6RDXN7Qc",
									"stake": "1542228206135376416333977438902",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "lunanova2.pool.f863973.m0",
									"public_key": "ed25519:9Jv6e9Kye4wM9EL1XJvXY8CYsLi1HLdRKnTzXBQY44w9",
									"stake": "1535510591656788451274253930651",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "hotones.pool.f863973.m0",
									"public_key": "ed25519:2fc5xtbafKiLtxHskoPL2x7BpijxSZcwcAjzXceaxxWt",
									"stake": "1315107169683243495728775410025",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "pathrocknetwork.pool.f863973.m0",
									"public_key": "ed25519:CGzLGZEMb84nRSRZ7Au1ETAoQyN7SQXQi55fYafXq736",
									"stake": "1055924587491205836633946214781",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "grassets.pool.f863973.m0",
									"public_key": "ed25519:3S4967Dt1VeeKrwBdTTR5tFEUFSwh17hEFLATRmtUNYV",
									"stake": "935242474781481476322268340415",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "stakesstone.pool.f863973.m0",
									"public_key": "ed25519:3aAdsKUuzZbjW9hHnmLWFRKwXjmcxsnLNLfNL4gP1wJ8",
									"stake": "909281497384996368820733078344",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "leadnode.pool.f863973.m0",
									"public_key": "ed25519:CdP6CBFETfWYzrEedmpeqkR6rsJNeT22oUFn2mEDGk5i",
									"stake": "905618873690263465459292715321",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "blockscope.pool.f863973.m0",
									"public_key": "ed25519:6K6xRp88BCQX5pcyrfkXDU371awMAmdXQY4gsxgjKmZz",
									"stake": "896047896979597050666985615574",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "al3c5.pool.f863973.m0",
									"public_key": "ed25519:BoYixTjyBePQ1VYP3s29rZfjtz1FLQ9og4FWZB5UgWCZ",
									"stake": "892989468728326357074301924678",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "shurik.pool.f863973.m0",
									"public_key": "ed25519:9zEn7DVpvQDxWdj5jSgrqJzqsLo8T9Wv37t83NXBiWi6",
									"stake": "860998472902654259868574822778",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "baziliknear.pool.f863973.m0",
									"public_key": "ed25519:9Rbzfkhkk6RSa1HoPnJXS4q2nn1DwYeB4HMfJBB4WQpU",
									"stake": "857723644364752944835556159412",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "optimusvalidatornetwork.pool.f863973.m0",
									"public_key": "ed25519:BGoxGmpvN7HdUSREQXfjH6kw5G6ph7NBXVfBVfUSH85V",
									"stake": "812700279565474879609093906327",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "chelovek_iz_naroda.pool.f863973.m0",
									"public_key": "ed25519:89aWsXXytjAZxyefXuGN73efnM9ugKTjPEGV4hDco8AZ",
									"stake": "805688435324658036814989629985",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "basilisk-stake.pool.f863973.m0",
									"public_key": "ed25519:CFo8vxoEUZoxbs87mGtG8qWUvSBHB91Vc6qWsaEXQ5cY",
									"stake": "782911221569156167328795553831",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "ou812.pool.f863973.m0",
									"public_key": "ed25519:2APjYBPnQ7CGDxFpsUHeAcyhpjZRWWpjciq2Pdzk31uQ",
									"stake": "631717642543179431613061598644",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "infstones.pool.f863973.m0",
									"public_key": "ed25519:BLP6HB8tcwYRTxswQ2YRaJ5sGj1dgGpUUfcNwbnWFGCU",
									"stake": "471885824911510679525121471143",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "kiln.pool.f863973.m0",
									"public_key": "ed25519:Bq8fe1eUgDRexX2CYDMhMMQBiN13j8vTAVFyTNhEfh1W",
									"stake": "334786926780084247228644505360",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "idtcn3.pool.f863973.m0",
									"public_key": "ed25519:DtkY9WtkWweSrF13BJi5k4c6xyk3tBAC9y92AEY4Ayfb",
									"stake": "155759995431387821259760930195",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "shardlabs.pool.f863973.m0",
									"public_key": "ed25519:DxmhGQZ6oqdxw7qGBvzLuBzE6XQjEh67hk5tt66vhLqL",
									"stake": "88870331313370762221390819947",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "gettingnear.pool.f863973.m0",
									"public_key": "ed25519:5QzHuNZ4stznMwf3xbDfYGUbjVt8w48q8hinDRmVx41z",
									"stake": "78148062589783451476140853866",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "sergo.pool.f863973.m0",
									"public_key": "ed25519:3uV2DGyNVfSgEPi9UdwFWBHtLZyPHmsN4iNtF7nEvZnT",
									"stake": "75508412425203533516705477393",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "leadnode-shard.pool.f863973.m0",
									"public_key": "ed25519:CzWox9TE1AR4xfDracD5eN4xy5f91ZNu5CTRsPcdH45C",
									"stake": "72965503363078502911434712939",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "pandateam.pool.f863973.m0",
									"public_key": "ed25519:7n426KJocZpJ5496UHp6onwYqWyt5xuiAyzvTGwCQLTN",
									"stake": "69577639450742001777679919468",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "sevennines-t0.pool.f863973.m0",
									"public_key": "ed25519:BHKMMc1t7F6B26BbaBVMhex3riPtWoGL5CricgXiSFC4",
									"stake": "63751489124387441582951570689",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "guardia.pool.f863973.m0",
									"public_key": "ed25519:2b5AQqcf8PHAUzqxWYoaFiTEa7QuEkpihzeoDittQaPL",
									"stake": "63133600633477268084701256889",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "redhead.pool.f863973.m0",
									"public_key": "ed25519:5qPvLhc86TDdof4YEjBMKrENzT3UA9mEonKWQXuFvaHX",
									"stake": "56730088229248773751756701318",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "tandemk.pool.f863973.m0",
									"public_key": "ed25519:8zqx8dzqsxXivPMcSk2e7gezMvooZb1RWMGHA1FzFrak",
									"stake": "53240451020555633927106476314",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "azetsi.pool.f863973.m0",
									"public_key": "ed25519:2MFKLj9E2kRdJoQqUgaY9KtheebzLv9ntdgTGsZxLaE1",
									"stake": "50676644892701662115298282618",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "wackazong.pool.f863973.m0",
									"public_key": "ed25519:EK1bdY5F6prLush2aKnJBe5neHdK52wwTjg3qY4v3cjX",
									"stake": "47906743971967209276054625040",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "pero_val.pool.f863973.m0",
									"public_key": "ed25519:J43JCHe2XKU7wiDi7PSS8exsdVgnYHcbnC8R14SpF6HV",
									"stake": "45787360402032110036453169642",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "do0k13.pool.f863973.m0",
									"public_key": "ed25519:A7wSnvPTLQcGVBsoLgyLqSw98mLRug4z6DKNETeFGzKx",
									"stake": "45290214292892250146292674687",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "qibaocenter.pool.f863973.m0",
									"public_key": "ed25519:ETL13ggWC6zGt5pPaQ2KyRaVmNYnhRBVGMPLTLAR5Fpa",
									"stake": "44531677218892266351417624560",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "twintest1.pool.f863973.m0",
									"public_key": "ed25519:7DKaSRvjyniVLkuyaKAiYA5Y2zZGLz69siHKYrqJVgqs",
									"stake": "43295592964525333762316994581",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "mlnear.pool.f863973.m0",
									"public_key": "ed25519:CknBNuCaqYmgxXB9grSS2aCxLRqaGEcsF6mTAKfVv9Sp",
									"stake": "41030066979177211258677827658",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "kt2.staking-farm-factory.testnet",
									"public_key": "ed25519:HzJtjiAzXrhG6mGGgNCLx4JTivaP6hwssd7gD2ve2Ej6",
									"stake": "40030001380790398998731784100",
									"validator_stake_struct_version": "V1"
								},
								{
									"account_id": "bgpntx.pool.f863973.m0",
									"public_key": "ed25519:DcCH7h9B4YhLc6BvdLTG5Mm62ZeLGT1yRrC7szbxRDaK",
									"stake": "39416268156247694908560833474",
									"validator_stake_struct_version": "V1"
								}
							]`),
						},
					},
					LatestChainStatus: mock.Response{
						Reponse: []byte(`{
							"chain_id": "testnet",
							"latest_protocol_version": 56,
							"node_key": null,
							"protocol_version": 56,
							"rpc_addr": "0.0.0.0:4040",
							"sync_info": {
								"earliest_block_hash": "FWJ9kR6KFWoyMoNjpLXXGHeuiy7tEY6GmoFeCA5yuc6b",
								"earliest_block_height": 42376888,
								"earliest_block_time": "2020-07-31T03:39:42.911378Z",
								"epoch_id": "85pkJkFvHtqjkwTRAeu9kz4RXg2kJSZj1oc7UAtS6dr5",
								"epoch_start_height": 99963122,
								"latest_block_hash": "CE3NVJbb5tStSPXNhuNEJc1UkHbuGdeL9UBFFcBrxHx3",
								"latest_block_height": 102814344,
								"latest_block_time": "2022-09-12T06:39:36.637866480Z",
								"latest_state_root": "3Xee355aQfunXZqQ2VKmhLAEjBkeHxzxYd2EycLxMYuT",
								"syncing": false
							},
							"uptime_sec": 928122,
							"validator_account_id": null,
							"validators": [
								{
									"account_id": "node0",
									"is_slashed": false
								},
								{
									"account_id": "node3",
									"is_slashed": false
								},
								{
									"account_id": "node2",
									"is_slashed": false
								},
								{
									"account_id": "node1",
									"is_slashed": false
								},
								{
									"account_id": "aurora.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "legends.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "masternode24.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "01node.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "p2p.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "nodeasy.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "chorusone.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "tribe-pool.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "foundryusa.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "sweden.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "chorus-one.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "ni.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "cryptogarik.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "pathrocknetwork.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "stakely_v2.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "everstake.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "freshtest.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "stakesstone.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "leadnode.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "dsrvlabs.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "al3c5.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "grassets.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "baziliknear.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "shurik.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "chelovek_iz_naroda.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "basilisk-stake.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "tayang.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "zetsi.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "infiniteloop.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "ou812.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "bflame.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "blazenet.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "dimasik.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "prophet.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "stingray.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "g2.pool.devnet",
									"is_slashed": false
								},
								{
									"account_id": "kiln.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "spectrum.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "dysprosium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "idtcn3.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "tin-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "neodymium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "beryllium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "barium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "vanadium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "gadolinium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "rhodium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "ruthenium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "bromine-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "gettingnear.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "chlorine-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "xenon-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "thulium-chunk.testnet",
									"is_slashed": false
								},
								{
									"account_id": "azetsi.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "pinpin.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "wackazong.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "guardia.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "alxvoy.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "redhead.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "domanodes.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "sergo.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "twintest1.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "tandemk.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "pero_val.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "leadnode-shard.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "cryptobtcbuyer.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "pandateam.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "meduza.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "adel0515.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "qibaocenter.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "gruberx.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "nodebull.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "makil.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "mlnear.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "p2pstaking.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "cunum.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "lastnode.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "omnistake_v5.factory01.littlefarm.testnet",
									"is_slashed": false
								},
								{
									"account_id": "jstaking.pool.f863973.m0",
									"is_slashed": false
								},
								{
									"account_id": "cryptolions.pool.f863973.m0",
									"is_slashed": false
								}
							],
							"version": {
								"build": "1.29.0-rc.2",
								"rustc_version": "1.62.1",
								"version": "1.29.0-rc.2"
							}
						}`),
					},
				})

				mockApi.On("Block", mock.MockParam).Return(mockApi.BlockFactory())
				mockApi.On("BlockProducers", mock.MockParam).Return(mockApi.BlockProducersFactory())
				mockApi.On("Status", mock.MockParam).Return(mockApi.StatusFactory())

				return mockApi
			}(),
			Expected: struct {
				Success interface{}
				Fail    interface{}
			}{
				Success: true,
			},
		},
	}

	RegisterTest("ValidateHeader", VerifierValidateHeaderTest{
		description: "Test ValidateHeader",
		testData:    testData,
	})
}
