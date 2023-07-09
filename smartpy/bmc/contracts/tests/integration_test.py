import subprocess
import smartpy as sp

path = {"bts_core": [
    {"types": ['types = sp.io.import_script_from_url("file:../bts/contracts/src/Types.py")',
               'types = sp.io.import_script_from_url("file:./contracts/src/Types.py")']},
    {"FA2_contract": [
        'FA2_contract = sp.io.import_script_from_url("file:../bts/contracts/src/FA2_contract.py")',
        'FA2_contract = sp.io.import_script_from_url("file:./contracts/src/FA2_contract.py")']}
    ],
    "bts_periphery": [
        {"types": ['types = sp.io.import_script_from_url("file:../bts/contracts/src/Types.py")',
                   'types = sp.io.import_script_from_url("file:./contracts/src/Types.py")']},
        {"strings": ['strings = sp.io.import_script_from_url("file:../bts/contracts/src/String.py")',
                     'strings = sp.io.import_script_from_url("file:./contracts/src/String.py")']},
        {"rlp": ['rlp = sp.io.import_script_from_url("file:../bts/contracts/src/RLP_struct.py")',
                 'rlp = sp.io.import_script_from_url("file:./contracts/src/RLP_struct.py")']},

    ],
    "RLP_struct": [
        {"types": ['types = sp.io.import_script_from_url("file:../bts/contracts/src/Types.py")',
                   'types = sp.io.import_script_from_url("file:./contracts/src/Types.py")']},
    ]
}

path2 = {
    "RLP_struct": [
        {'_to_byte':
         ['        _to_byte = sp.view("encode_nat", self.data.helper, sp.as_nat(params.sn), t=sp.TBytes).open_some()',
          '        _to_byte = sp.view("to_byte", self.data.helper_parse_negative, params.sn,'
          ' t=sp.TBytes).open_some()'
          ]
         },
        {'_to_int':
         ['                    _to_int = sp.to_int(Utils2.Int.of_bytes(sn_in_bytes))',
          '                    _to_int = sp.view("to_int", self.data.helper_parse_negative, sn_in_bytes, t=sp.TInt)'
          '.open_some()'
          ]
         }
    ],
    "bmc_periphery": [
        {'_self_address':
         ['        _self_address = sp.address("KT1VTmeVTccqv3opzkbRVrYwaoSZTTEzfJ8b")',
          '        _self_address = sp.self_address'
          ]
         }
    ]
}


def patch_file_path(file_name, old_value, new_value):
    subprocess.call("sed -i -e 's#.*" + old_value + " =.*#" + new_value + "#' " + file_name, shell=True)


def bts_core_contract_deploy_setup():
    for key, value in path.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("../bts/contracts/src/" + key + ".py", lis1[0], y[0])


def bmc_periphery_contract_deploy_setup():
    for key, value in path2.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("./contracts/src/" + key + ".py", lis1[0], y[0])


def tear_down_bts_changes():
    for key, value in path.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("../bts/contracts/src/" + key + ".py", lis1[0], y[1])


def tear_down_bmc_changes():
    for key, value in path2.items():
        for i in value:
            lis1 = []
            for x, y in i.items():
                lis1.append(x)
                # lis1.append(y)
                patch_file_path("./contracts/src/" + key + ".py", lis1[0], y[1])


# import changes in bts_core for testing
bts_core_contract_deploy_setup()
bmc_periphery_contract_deploy_setup()

BTSCore = sp.io.import_script_from_url("file:../bts/contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:../bts/contracts/src/bts_owner_manager.py")
BTSPeriphery = sp.io.import_script_from_url("file:../bts/contracts/src/bts_periphery.py")

FA2 = sp.io.import_script_from_url("https://legacy.smartpy.io/templates/fa2_lib.py")
BMCManagement = sp.io.import_script_from_url("file:./contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:../bts/contracts/src/parse_address.py")
fa2_dummy_file = sp.io.import_script_from_url("file:./contracts/tests/fa2_dummy.py")

# revert the path changes made in bts_core for testing
tear_down_bts_changes()


@sp.add_test("IntegrationTest")
def test():
    sc = sp.test_scenario()

    # test account
    alice = sp.test_account("Alice")
    owner = sp.test_account("Owner")
    jack = sp.test_account("Jack")
    bob = sp.test_account("Bob")
    creator2 = sp.test_account("creator2")
    service1_address = sp.test_account("service1_address")
    service2_address = sp.test_account("service2_address")
    relay = sp.test_account("Relay")
    helper_parse_neg_contract = sp.test_account("helper_parse_neg_contract")

    # Change icon_bmc_address for new environment
    icon_bmc_address = "cxb7de63db8c1fa2d9dfb6c531e6bc19402572cc23"
    icon_bmc_block_height = 10445602

    # deploy BMCManagement contract
    helper_contract = deploy_helper_contract()
    sc += helper_contract

    bmc_management = deploy_bmc_management(owner.address, helper_contract.address)
    sc += bmc_management

    parse_address = deploy_parse_address()
    sc += parse_address

    bmc_periphery = deploy_bmc_periphery(bmc_management.address, helper_contract.address,
                                         helper_parse_neg_contract.address, parse_address.address, owner.address)
    sc += bmc_periphery

    bts_owner_manager = deploy_bts_owner_manager_contract(owner.address)
    sc += bts_owner_manager

    bts_core = deploy_bts_core(bts_owner_manager.address)
    sc += bts_core

    bts_periphery = deploy_bts_periphery(bts_core.address, helper_contract.address,
                                         parse_address.address, bmc_periphery.address,
                                         owner.address)
    sc += bts_periphery

    fa2_dummy = fa2_dummy_file.SingleAssetToken(admin=owner.address,
                                                metadata=sp.utils.metadata_of_url(
                                                    "ipfs://example"),
                                                token_metadata=FA2.make_metadata(name="NativeWrappedCoin", decimals=6,
                                                                                 symbol="wTEZ"))
    sc += fa2_dummy

    # BMC_MANAGEMENT SETTERS
    # set bmc periphery
    bmc_management.set_bmc_periphery(bmc_periphery.address).run(
        sender=owner.address)

    # set bmc_btp_address
    bmc_management.set_bmc_btp_address("NetXnHfVqm9iesp.tezos").run(
        sender=owner.address)

    # tear down changes after bmc btp address is set
    tear_down_bmc_changes()

    # add_service
    svc1 = sp.string("bts")
    bmc_management.add_service(sp.record(addr=bts_periphery.address, svc=svc1)).run(
        sender=owner.address)

    # add_route
    dst = "btp://0x7.icon/" + icon_bmc_address
    next_link = "btp://0x7.icon/" + icon_bmc_address
    bmc_management.add_route(sp.record(dst=dst, link=next_link)).run(
        sender=owner.address)

    # add_link
    bmc_management.add_link("btp://0x7.icon/" + icon_bmc_address).run(
        sender=owner.address)

    # set_link_rx_height
    bmc_management.set_link_rx_height(sp.record(height=icon_bmc_block_height,
                                                link="btp://0x7.icon/" + icon_bmc_address)).run(
        sender=owner.address)

    # add_relay
    bmc_management.add_relay(sp.record(link="btp://0x7.icon/" + icon_bmc_address,
                                       addr=sp.set([relay.address]))).run(
        sender=owner.address)

    # BTS_CORE SETTERS
    # update_bts_periphery
    bts_core.update_bts_periphery(bts_periphery.address).run(sender=owner.address)

    # set_fee_ratio
    bts_core.set_fee_ratio(name=sp.string("btp-NetXnHfVqm9iesp.tezos-XTZ"), fee_numerator=sp.nat(100),
                           fixed_fee=sp.nat(450)).run(sender=owner.address)

    prev = "btp://0x7.icon/" + icon_bmc_address

    # 1: Init message from relay
    msg_byte = sp.bytes(
        "0xf8e5f8e3b8e1f8df01b8d7f8d5f8d3b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b54315"
        "6546d65565463637176336f707a6b625256725977616f535a5454457a664a386201b88ef88cb8396274703a2f2f3078372"
        "e69636f6e2f637862376465363364623863316661326439646662366335333165366263313934303235373263633233b840"
        "6274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543156546d65565463637176336f707a6b6252567"
        "25977616f535a5454457a664a386283626d630089c884496e697482c1c084009f639c")
    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # 2: Transfer 100 native coin
    bts_core.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(
        sender=alice.address, amount=sp.tez(100))

    # 3: relay msg for transfer end of step 2 with fee gathering
    msg_byte = sp.bytes(
        "0xf90218f90215b90212f9020f01b90206f90203f8e0b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b"
        "543156546d65565463637176336f707a6b625256725977616f535a5454457a664a386202b89bf899b8396274703a2f2f3078372e69636"
        "f6e2f637862376465363364623863316661326439646662366335333165366263313934303235373263633233b8406274703a2f2f4e657"
        "4586e486656716d39696573702e74657a6f732f4b543156546d65565463637176336f707a6b625256725977616f535a5454457a664a38"
        "62836274730196d50293d200905472616e736665722053756363657373f9011eb8406274703a2f2f4e6574586e486656716d396965737"
        "02e74657a6f732f4b543156546d65565463637176336f707a6b625256725977616f535a5454457a664a386203b8d9f8d7b8396274703a"
        "2f2f3078372e69636f6e2f637862376465363364623863316661326439646662366335333165366263313934303235373263633233b8"
        "406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543156546d65565463637176336f707a6b62525672597761"
        "6f535a5454457a664a386283626d6300b853f8518c466565476174686572696e67b842f840b8396274703a2f2f3078372e69636f6e2f"
        "687866383061643730393832643637636437363438396665613966653036343239626362306266646531c4836274738400a01441")

    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # transfer end of fee gathering
    msg_byte = sp.bytes("0xf8f2f8f0b8eef8ec01b8e4f8e2f8e0b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f"
                        "4b543156546d65565463637176336f707a6b625256725977616f535a5454457a664a386204b89bf899b8396274703"
                        "a2f2f3078372e69636f6e2f63786237646536336462386331666132643964666236633533316536626331393430323"
                        "5373263633233b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b543156546d655654636"
                        "37176336f707a6b625256725977616f535a5454457a664a3862836274730296d50293d200905472616e73666572205"
                        "37563636573738400a01489")

    bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg_byte)).run(sender=relay.address)

    # mint fa2 dummy
    fa2_dummy.mint([sp.record(to_=alice.address, amount=sp.nat(200000000))]).run(sender=owner.address)

    # add operator
    fa2_dummy.update_operators(
        [sp.variant("add_operator", sp.record(owner=alice.address, operator=bts_core.address, token_id=0))]).run(
        sender=alice.address)

    # register fa2
    bts_core.register(
        name=sp.string("test"),
        fee_numerator=sp.nat(100),
        fixed_fee=sp.nat(450),
        addr=fa2_dummy.address,
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    # 4: transfer fa2 token
    bts_core.transfer(sp.record(
        coin_name="test", value=50000000, to="btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d")).run(
        sender=alice.address)

    # 5: relay msg for transfer end of transfer fa2
    # msg_byte = sp.bytes(
    #     "")

    # # test 1: Test of add to blacklist function
    # bts_periphery.add_to_blacklist({0: "notaaddress"}).run(sender=bts_periphery.address, valid=False,
    #                                                        exception="InvalidAddress")
    # bts_periphery.add_to_blacklist({0: "tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=bts_periphery.address,
    #                                                                                 valid=True)
    # bts_periphery.add_to_blacklist({0: "tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=bts_periphery.address,
    #                                                                                 valid=True)
    # bts_periphery.add_to_blacklist({0: "tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=alice.address, valid=False,
    #                                                                                 exception="Unauthorized")
    # bts_periphery.add_to_blacklist({0: 'tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg'}).run(sender=bts_periphery.address,
    #                                                                                 valid=False,
    #                                                                                 exception='InvalidAddress')
    # sc.verify(bts_periphery.data.blacklist[
    #               sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW")] == True)
    #
    # # # transfer_native_coin
    # # bts_periphery.set_token_limit(
    # #     sp.record(
    # #         coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
    # #         token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
    # #     )
    # # ).run(sender = bts_core.address)
    # # bts_core.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30), valid=False, exception="FailCheckTransfer")
    # # bts_core.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(
    # #     sender=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30))
    #
    # # # handle_relay_message
    # # msg = sp.bytes(
    # #     "0xf8e1f8dfb8ddf8db01b8d3f8d1f8cfb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b54314657446a4338435941777268424a6d43355554596d594a544a6457776f3447676203b88af888b8396274703a2f2f3078372e69636f6e2f637833643436306163643535356336373034303566396562323934333833356366643132326662323938b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b54314657446a4338435941777268424a6d43355554596d594a544a6457776f344767628362747381ff84c328f8008400886513")
    # # prev = sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b")
    # # bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg)).run(
    # #     sender=owner.address)
    #
    # # test 2 : Test of remove from blacklist function
    # bts_periphery.remove_from_blacklist({0: 'notaaddress'}).run(sender=bts_periphery.address, valid=False,
    #                                                             exception="InvalidAddress")  # invalid address
    # bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address,
    #                                                                                      valid=False,
    #                                                                                      exception="UserNotFound")  # address not black-listed
    # bts_periphery.add_to_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
    #     sender=bts_periphery.address)  # adding to blacklist
    # bts_periphery.remove_from_blacklist({0: 'tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
    #     sender=bts_periphery.address)  # valid process
    # bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address,
    #                                                                                      valid=False,
    #                                                                                      exception='UserNotFound')  # cannot remove from blacklist twice
    # bts_periphery.add_to_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
    #     sender=bts_periphery.address)  # adding to blacklist
    # bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
    #     sender=bts_periphery.address)  # can only be called from btseperiphery contract
    # bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=alice.address,
    #                                                                                      valid=False,
    #                                                                                      exception="Unauthorized")  # can only be called from btseperiphery contract
    #
    # # # transfer_native_coin
    # # bts_periphery.set_token_limit(
    # #     sp.record(
    # #         coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
    # #         token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
    # #     )
    # # ).run(sender = bts_core.address)
    # # bts_core.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30))
    #
    # # # transfer_native_coin
    # # bts_periphery.set_token_limit(
    # #     sp.record(
    # #         coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
    # #         token_limit=sp.map({0: 5})
    # #     )
    # # ).run(sender = bts_core.address)
    # # bts_core.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30),  valid=False, exception="FailCheckTransfer")
    #
    # # # bmc_periphery get_status
    # # sc.verify_equal(bmc_periphery.get_status("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"), sp.record(current_height = 0, rx_height = 2, rx_seq = 0, tx_seq = 6))
    #
    # # test 3 : set token  limit
    # bts_periphery.set_token_limit(
    #     sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(5), 1: sp.nat(2)})).run(sender=alice.address,
    #                                                                                               valid=False,
    #                                                                                               exception='Unauthorized')  # can only be called from btsperiphery contract
    # bts_periphery.set_token_limit(
    #     sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(5), 1: sp.nat(2)})).run(
    #     sender=bts_periphery.address)  # set token limit for Tok2 coin to 5 and BB coin to 2
    # sc.verify(bts_periphery.data.token_limit["Tok2"] == sp.nat(5))  # test of token_limit for tok2 token
    # bts_periphery.set_token_limit(sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(5)})).run(
    #     valid=False, exception='InvalidParams', sender=bts_periphery.address)  # invalid parameters
    # # cannot set more than 15 token limit at once
    # bts_periphery.set_token_limit(
    #     sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(15), 1: sp.nat(22)})).run(
    #     sender=bts_periphery.address)  # can modify already set data
    # sc.verify(bts_periphery.data.token_limit["BB"] == sp.nat(22))  # test of token_limit for tok2 token
    #
    # # # handle_relay_message
    # # msg=sp.bytes("0xf90157f90154b90151f9014e01b90145f90142f9013fb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431454e5a76546f507838374c68756f774a315669786a6a715168536e597263594c6907b8faf8f8b8396274703a2f2f3078372e69636f6e2f637864633238393434343037363539393733666539393438376437356335646433326337396265303533b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431454e5a76546f507838374c68756f774a315669786a6a715168536e597263594c698362747303b874f87200b86ff86daa687839643138316431336634376335616165353535623730393831346336623232393738373937363139a4747a3165703766664b7351434e64676e6b504443566e566b67626d465a50386d464e3147dcdb906274702d3078372e69636f6e2d4943588900d71b0fe0a28e000084008502ba")
    # # prev=sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b")
    # # bmc_periphery.handle_relay_message(sp.record(prev=prev, msg=msg)).run(sender=owner.address)


def deploy_bmc_management(owner, helper):
    bmc_management = BMCManagement.BMCManagement(owner, helper)
    return bmc_management


def deploy_bmc_periphery(bmc_address, helper, helper_parse_neg_contract, parse, owner):
    bmc_periphery = BMCPeriphery.BMCPreiphery(bmc_address, helper, helper_parse_neg_contract, parse, owner)
    return bmc_periphery


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address


def deploy_bts_core(bts_owner_manager_contract):
    bts_core = BTSCore.BTSCore(
        owner_manager=bts_owner_manager_contract,
        _native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
        _fee_numerator=sp.nat(100),
        _fixed_fee=sp.nat(450)
    )
    return bts_core


def deploy_bts_owner_manager_contract(owner):
    bts_owner_manager_contract = BTSOwnerManager.BTSOwnerManager(owner)
    return bts_owner_manager_contract


def deploy_bts_periphery(core_address, helper, parse, bmc, owner):
    bts_periphery_contract = BTSPeriphery.BTSPeriphery(bmc_address=bmc, bts_core_address=core_address,
                                                       helper_contract=helper, parse_address=parse,
                                                       owner_address=owner,
                                                       native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ")
    return bts_periphery_contract
