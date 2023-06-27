import smartpy as sp


BTSCore = sp.io.import_script_from_url("file:./contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:./contracts/src/bts_owner_manager.py")
BTSPeriphery = sp.io.import_script_from_url("file:./contracts/src/bts_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:./contracts/src/parse_address.py")



@sp.add_test("BTSCoreTest")
def test():
    sc = sp.test_scenario()
    sc.h1("BTSCore")


    # test account
    alice = sp.test_account("Alice")
    jack = sp.test_account("Jack")
    bob = sp.test_account("Bob")
    creator = sp.test_account("Creator")

    # deploy BTSCore contract
    bts_owner_manager = deploy_btsOwnerManager_Contract()
    sc += bts_owner_manager
    btsCore_contract = deploy_btsCore_contract(bts_owner_manager.address)
    sc += btsCore_contract

    helper_contract = deploy_helper_contract()
    sc += helper_contract

    parse_address = deploy_parse_address()
    sc += parse_address


    bts_periphery = deploy_btsPeriphery_Contract(btsCore_contract.address, helper_contract.address, parse_address.address)
    sc += bts_periphery
    fa2 = deploy_fa2_Contract(bts_periphery.address)
    sc += fa2




    bts_owner_manager.is_owner(sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))


    # test case 1: check the name of coin
    sc.verify_equal(btsCore_contract.get_native_coin_name(), "BTSCOIN")


    # test case 2: check the owner manager of the contract
    sc.verify_equal(btsCore_contract.data.bts_owner_manager, bts_owner_manager.address)


    # test case 3: update_bts_periphery function
    btsCore_contract.update_bts_periphery(bts_periphery.address).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    sc.verify_equal(btsCore_contract.data.bts_periphery_address,sp.some(bts_periphery.address))

    # test case 4: set_fee_ratio function
    #throws error if coin name is different
    btsCore_contract.set_fee_ratio(name=sp.string("coindiff"),fee_numerator=sp.nat(100),fixed_fee=sp.nat(10)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception='TokenNotExists')
    #throws error when fee numerator is greater than denominator
    btsCore_contract.set_fee_ratio(name=sp.string("BTSCOIN"),fee_numerator=sp.nat(10000),fixed_fee=sp.nat(10)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception='InvalidSetting')
    # throws error when fixed fee is less than 0
    btsCore_contract.set_fee_ratio(name=sp.string("BTSCOIN"),fee_numerator=sp.nat(100),fixed_fee=sp.nat(0)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception='LessThan0')
    #works
    btsCore_contract.set_fee_ratio(name=sp.string("BTSCOIN"),fee_numerator=sp.nat(100),fixed_fee=sp.nat(10)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    #checking value
    sc.verify_equal(btsCore_contract.data.coin_details["BTSCOIN"].fee_numerator, 100)
    sc.verify_equal(btsCore_contract.data.coin_details["BTSCOIN"].fixed_fee, 10)


    # test case 5: register function
    # name shouldn't be native coin name
    btsCore_contract.register(
        name=sp.string("BTSCOIN"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception='ExistNativeCoin')
    #throws error when fee numerator is greater than denominator
    btsCore_contract.register(
        name=sp.string("new_coin1"),
        fee_numerator=sp.nat(100000),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception='InvalidSetting')
    # works
    btsCore_contract.register(
        name=sp.string("new_coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=fa2.address,
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # verifying the value
    sc.verify_equal(btsCore_contract.data.coins_name, ['new_coin', 'BTSCOIN'])
    sc.verify_equal(btsCore_contract.data.coin_details['new_coin'], sp.record(addr = fa2.address, coin_type = 2, fee_numerator = 10, fixed_fee = 2))
    # throws error when existed coin name is given
    btsCore_contract.register(
        name=sp.string("new_coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception="ExistCoin")
    # throws error when existed address  is given
    btsCore_contract.register(
        name=sp.string("new_coin2"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=fa2.address,
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), valid=False, exception="AddressExists")


    # test case 6: coin_names function
    sc.verify_equal(btsCore_contract.coin_names(), ['new_coin', 'BTSCOIN'])

    # test case 7: coin_id function
    sc.verify_equal(btsCore_contract.coin_id('new_coin'), fa2.address)

    # test case 8: is_valid_coin function
    sc.verify_equal(btsCore_contract.is_valid_coin('new_coin'), True)
    sc.verify_equal(btsCore_contract.is_valid_coin('not_valid'), False)


    # test case 9: fee_ratio function
    sc.verify_equal(btsCore_contract.fee_ratio('new_coin'), sp.record(fee_numerator=10, fixed_fee=2))


    #test case 10: balance_of function
    sc.verify_equal(
        btsCore_contract.balance_of(
            sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_name='new_coin')
        ),
        sp.record(usable_balance=0, locked_balance=0, refundable_balance=0, user_balance=0)
    )

    #test case 11: balance_of_batch function
    btsCore_contract.balance_of_batch(
        sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_names=['new_coin', 'BTSCOIN'])
    )
    sc.verify_equal(
        btsCore_contract.balance_of_batch(
            sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_names=['new_coin', 'BTSCOIN'])
        ),
        sp.record(locked_balances = {0 : 0, 1 : 0}, refundable_balances = {0 : 0, 1 : 0}, usable_balances = {0 : 0, 1 : 0}, user_balances = {0 : 0, 1 : 0}))


    # #test case 13: get_accumulated_fees function
    # # sc.verify_equal(btsCore_contract.get_accumulated_fees(), {})


    # #test case 14: transfer_native_coin function
    bts_periphery.set_token_limit(
        sp.record(
            coin_names=sp.map({0: "BTSCOIN"}),
            token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
        )
    ).run(sender = btsCore_contract.address)
    btsCore_contract.transfer_native_coin("tz1eZMrKqCNPrHzykdTuqKRyySoDv4QRSo7d").run(sender= sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), amount=sp.tez(30))



    # # test case 15: transfer function
    # fa2.mint([sp.record(to_=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), amount=sp.nat(100))]).run(sender=bts_periphery.address)
    # fa2.set_allowance([sp.record(spender=btsCore_contract.address, amount=sp.nat(100))]).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    # fa2.update_operators(
    #      [sp.variant("add_operator", sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), operator=btsCore_contract.address, token_id=0))]).run(
    #      sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))


    # btsCore_contract.transfer(coin_name='new_coin', value=10,  to="tz1eZMrKqCNPrHzykdTuqKRyySoDv4QRSo7d").run(sender = sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # sc.verify_equal(
    #     btsCore_contract.balance_of(
    #         sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_name='new_coin')
    #     ),
    #     sp.record(usable_balance=90, locked_balance=10, refundable_balance=0, user_balance=90)
    # )

    #test case 12: minting and checking balance
    # mint for native coin
    # btsCore_contract.mint(to=sp.address("tz1UHiurbSsDFXnTfMQkeni4GWCuMpq6JFRB"), coin_name="BTSCOIN", value=1000).run(sender = bts_periphery.address)
    # sc.verify_equal(
    #     btsCore_contract.balance_of(
    #         sp.record(owner=sp.address("tz1UHiurbSsDFXnTfMQkeni4GWCuMpq6JFRB"), coin_name='BTSCOIN')
    #     ),
    #     sp.record(usable_balance=0, locked_balance=0, refundable_balance=0, user_balance=1000)
    # )
    # error balance_of native_coin balance ma halne condition

    # for NON_NATIVE_TOKEN_TYPE
    # fa2.mint([sp.record(to_=bts_periphery.address, amount=sp.nat(2000))]).run(
    #     sender=bts_periphery.address)
    # fa2.update_operators(
    #      [sp.variant("add_operator", sp.record(owner=bts_periphery.address, operator=btsCore_contract.address, token_id=0))]).run(
    #      sender=bts_periphery.address)
    # fa2.set_allowance([sp.record(spender=btsCore_contract.address, amount=sp.nat(1000))]).run(
    #     sender=bts_periphery.address)

    # btsCore_contract.mint(to=sp.address("tz1UHiurbSsDFXnTfMQkeni4GWCuMpq6JFRB"), coin_name="new_coin", value=1000).run(sender = bts_periphery.address)
    # sc.verify_equal(
    #     btsCore_contract.balance_of(
    #         sp.record(owner=sp.address("tz1UHiurbSsDFXnTfMQkeni4GWCuMpq6JFRB"), coin_name='new_coin')
    #     ),
    #     sp.record(usable_balance=0, locked_balance=0, refundable_balance=0, user_balance=1000)
    # )


    # # test case 16: transfer_batch function
    # btsCore_contract.transfer_batch(
    #     coin_names={0: 'new_coin', 1: 'new_coin'},
    #     values={0: 10, 1: 10},
    #     to="tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"
    # ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # sc.verify_equal(
    #     btsCore_contract.balance_of(
    #         sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_name='new_coin')
    #     ),
    #     sp.record(usable_balance=70, locked_balance=30, refundable_balance=0, user_balance=70)
    # )

    # # test case 17: handle_response_service function
    # btsCore_contract.handle_response_service(sp.record(requester=sp.address("KT1VCbyNieUsQsCShkxtTz9ZbLmE9oowmJPm"), coin_name="BTSCOIN",value=sp.nat(44), fee=sp.nat(3), rsp_code=sp.nat(1))).run(sender=bts_periphery.address)





#
#     #
#     # # test case 16: reclaim function
#     #
#     # # btsCore_contract.reclaim(coin_name="new_coin", value=25).run(sender = sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
#     #
#     #
#     # # btsCore_contract.refund(to=sp.address("tz1UHiurbSsDFXnTfMQkeni4GWCuMpq6JFRB"), coin_name="new_coin", value=30).run(sender = sp.address("tz1UHiurbSsDFXnTfMQkeni4GWCuMpq6JFRB"))
#     #
#     #
#     #
#     #
#     #
#

#
#
#

def deploy_btsCore_contract(bts_OwnerManager_Contract):
    btsCore_contract = BTSCore.BTSCore(
        owner_manager=bts_OwnerManager_Contract,
        _native_coin_name="BTSCOIN",
        _fee_numerator=sp.nat(1000),
        _fixed_fee=sp.nat(10)
    )
    return btsCore_contract

def deploy_btsOwnerManager_Contract():
    bts_OwnerManager_Contract = BTSOwnerManager.BTSOwnerManager(sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    return bts_OwnerManager_Contract

def deploy_btsPeriphery_Contract(core_address, helper, parse):
    bmc = sp.test_account("bmc")
    btsPeriphery_Contract = BTSPeriphery.BTPPreiphery(bmc_address= bmc.address, bts_core_address=core_address, helper_contract=helper, parse_address=parse,native_coin_name= 'BTSCoin',owner_address=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    return btsPeriphery_Contract

def deploy_fa2_Contract(admin_address):
    fa2_contract = BTSCore.FA2_contract.SingleAssetToken(admin=admin_address, metadata=sp.big_map({"ss": sp.bytes("0x0dae11")}), token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}))
    return fa2_contract

def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract

def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address


