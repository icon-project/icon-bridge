import smartpy as sp

BMCManagement = sp.io.import_script_from_url("file:./contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:../bts/contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:../bts/contracts/src/parse_address.py")

BTSCore = sp.io.import_script_from_url("file:../bts/contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:../bts/contracts/src/bts_owner_manager.py")
BTSPeriphery = sp.io.import_script_from_url("file:../bts/contracts/src/bts_periphery.py")


@sp.add_test("BMCManagementTest")
def test():
    sc = sp.test_scenario()

    # test account
    alice = sp.test_account("Alice")
    creator = sp.test_account("Creator")
    jack = sp.test_account("Jack")
    bob = sp.test_account("Bob")
    creator2 = sp.test_account("creator2")
    service1_address = sp.test_account("service1_address")
    service2_address = sp.test_account("service2_address")

    # deploy BMCManagement contract
    helper_contract = deploy_helper_contract()
    sc += helper_contract

    bmc_management_contract = deploy_bmc_management_contract(sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"),
                                                             helper_contract.address)
    sc += bmc_management_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bmc_periphery_contract = deploy_bmc_periphery_contract(bmc_management_contract.address, helper_contract.address,
                                                           parse_address.address)
    sc += bmc_periphery_contract

    bts_owner_manager = deploy_bts_owner_manager_contract()
    sc += bts_owner_manager
    bts_core_contract = deploy_bts_core_contract(bts_owner_manager.address)
    sc += bts_core_contract

    bts_periphery = deploy_bts_periphery_contract(bts_core_contract.address, helper_contract.address,
                                                  parse_address.address, bmc_periphery_contract.address)
    sc += bts_periphery

    # set bmc periphery
    bmc_management_contract.set_bmc_periphery(bmc_periphery_contract.address).run(
        sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    # set bmc_btp_address(netwrk address)
    bmc_management_contract.set_bmc_btp_address("NetXnHfVqm9iesp.tezos").run(
        sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    # update_bts_periphery
    bts_core_contract.update_bts_periphery(bts_periphery.address).run(
        sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # add_service
    svc1 = sp.string("bts")
    bmc_management_contract.add_service(sp.record(addr=bts_periphery.address, svc=svc1)).run(
        sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # add_route
    dst = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest"
    next_link = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"
    bmc_management_contract.add_route(sp.record(dst=dst, link=next_link)).run(
        sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # # add_link
    # bmc_management_contract.add_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(
    #     sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # # add_relay
    # bmc_management_contract.add_relay(sp.record(link="btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b",
    #                                             addr=sp.set([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")]))).run(
    #     sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # test 1: Test of add to blacklist function
    bts_periphery.add_to_blacklist({0: "notaaddress"}).run(sender=bts_periphery.address, valid=False,
                                                           exception="InvalidAddress")
    bts_periphery.add_to_blacklist({0: "tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=bts_periphery.address,
                                                                                    valid=True)
    bts_periphery.add_to_blacklist({0: "tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=bts_periphery.address,
                                                                                    valid=True)
    bts_periphery.add_to_blacklist({0: "tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=alice.address, valid=False,
                                                                                    exception="Unauthorized")
    bts_periphery.add_to_blacklist({0: 'tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg'}).run(sender=bts_periphery.address,
                                                                                    valid=False,
                                                                                    exception='InvalidAddress')
    sc.verify(bts_periphery.data.blacklist[
                  sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW")] == True)

    # # transfer_native_coin
    # bts_periphery.set_token_limit(
    #     sp.record(
    #         coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
    #         token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
    #     )
    # ).run(sender = bts_core_contract.address)
    # bts_core_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30), valid=False, exception="FailCheckTransfer")
    # bts_core_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(
    #     sender=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30))

    # # handle_relay_message
    # msg = sp.bytes(
    #     "0xf8e1f8dfb8ddf8db01b8d3f8d1f8cfb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b54314657446a4338435941777268424a6d43355554596d594a544a6457776f3447676203b88af888b8396274703a2f2f3078372e69636f6e2f637833643436306163643535356336373034303566396562323934333833356366643132326662323938b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b54314657446a4338435941777268424a6d43355554596d594a544a6457776f344767628362747381ff84c328f8008400886513")
    # prev = sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b")
    # bmc_periphery_contract.handle_relay_message(sp.record(prev=prev, msg=msg)).run(
    #     sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # test 2 : Test of remove from blacklist function
    bts_periphery.remove_from_blacklist({0: 'notaaddress'}).run(sender=bts_periphery.address, valid=False,
                                                                exception="InvalidAddress")  # invalid address
    bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address,
                                                                                         valid=False,
                                                                                         exception="UserNotFound")  # address not black-listed
    bts_periphery.add_to_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
        sender=bts_periphery.address)  # adding to blacklist
    bts_periphery.remove_from_blacklist({0: 'tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
        sender=bts_periphery.address)  # valid process
    bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address,
                                                                                         valid=False,
                                                                                         exception='UserNotFound')  # cannot remove from blacklist twice
    bts_periphery.add_to_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
        sender=bts_periphery.address)  # adding to blacklist
    bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(
        sender=bts_periphery.address)  # can only be called from btseperiphery contract
    bts_periphery.remove_from_blacklist({0: 'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=alice.address,
                                                                                         valid=False,
                                                                                         exception="Unauthorized")  # can only be called from btseperiphery contract

    # # transfer_native_coin
    # bts_periphery.set_token_limit(
    #     sp.record(
    #         coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
    #         token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
    #     )
    # ).run(sender = bts_core_contract.address)
    # bts_core_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30))

    # # transfer_native_coin
    # bts_periphery.set_token_limit(
    #     sp.record(
    #         coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
    #         token_limit=sp.map({0: 5})
    #     )
    # ).run(sender = bts_core_contract.address)
    # bts_core_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30),  valid=False, exception="FailCheckTransfer")

    # # bmc_periphery get_status
    # sc.verify_equal(bmc_periphery_contract.get_status("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"), sp.record(current_height = 0, rx_height = 2, rx_seq = 0, tx_seq = 6))

    # test 3 : set token  limit
    bts_periphery.set_token_limit(
        sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(5), 1: sp.nat(2)})).run(sender=alice.address,
                                                                                                  valid=False,
                                                                                                  exception='Unauthorized')  # can only be called from btsperiphery contract
    bts_periphery.set_token_limit(
        sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(5), 1: sp.nat(2)})).run(
        sender=bts_periphery.address)  # set token limit for Tok2 coin to 5 and BB coin to 2
    sc.verify(bts_periphery.data.token_limit["Tok2"] == sp.nat(5))  # test of token_limit for tok2 token
    bts_periphery.set_token_limit(sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(5)})).run(
        valid=False, exception='InvalidParams', sender=bts_periphery.address)  # invalid parameters
    # cannot set more than 15 token limit at once
    bts_periphery.set_token_limit(
        sp.record(coin_names={0: "Tok2", 1: 'BB'}, token_limit={0: sp.nat(15), 1: sp.nat(22)})).run(
        sender=bts_periphery.address)  # can modify already set data
    sc.verify(bts_periphery.data.token_limit["BB"] == sp.nat(22))  # test of token_limit for tok2 token

    # # handle_relay_message
    # msg=sp.bytes("0xf90157f90154b90151f9014e01b90145f90142f9013fb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431454e5a76546f507838374c68756f774a315669786a6a715168536e597263594c6907b8faf8f8b8396274703a2f2f3078372e69636f6e2f637864633238393434343037363539393733666539393438376437356335646433326337396265303533b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431454e5a76546f507838374c68756f774a315669786a6a715168536e597263594c698362747303b874f87200b86ff86daa687839643138316431336634376335616165353535623730393831346336623232393738373937363139a4747a3165703766664b7351434e64676e6b504443566e566b67626d465a50386d464e3147dcdb906274702d3078372e69636f6e2d4943588900d71b0fe0a28e000084008502ba")
    # prev=sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b")
    # bmc_periphery_contract.handle_relay_message(sp.record(prev=prev, msg=msg)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # set_fee_ratio
    bts_core_contract.set_fee_ratio(name=sp.string("btp-NetXnHfVqm9iesp.tezos-XTZ"), fee_numerator=sp.nat(100),
                                    fixed_fee=sp.nat(450)).run(
        sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))


def deploy_bmc_management_contract(owner, helper):
    bmc_management_contract = BMCManagement.BMCManagement(owner, helper)
    return bmc_management_contract


def deploy_bmc_periphery_contract(bmc_addres, helper, parse):
    owner = sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")
    bmc_periphery_contract = BMCPeriphery.BMCPreiphery(bmc_addres, helper, parse, owner)
    return bmc_periphery_contract


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address


def deploy_bts_core_contract(bts_owner_manager_contract):
    bts_core_contract = BTSCore.BTSCore(
        owner_manager=bts_owner_manager_contract,
        _native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
        _fee_numerator=sp.nat(100),
        _fixed_fee=sp.nat(450)
    )
    return bts_core_contract


def deploy_bts_owner_manager_contract():
    bts_owner_manager_contract = BTSOwnerManager.BTSOwnerManager(sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    return bts_owner_manager_contract


def deploy_bts_periphery_contract(core_address, helper, parse, bmc):
    bts_periphery_contract = BTSPeriphery.BTPPreiphery(bmc_address=bmc, bts_core_address=core_address,
                                                       helper_contract=helper, parse_address=parse,
                                                       owner_address=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"),
                                                       native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ")
    return bts_periphery_contract
