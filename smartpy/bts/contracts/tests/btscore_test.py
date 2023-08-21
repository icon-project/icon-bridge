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
    bob = sp.test_account("Bob")
    owner = sp.test_account("owner")

    # deploy BTSCore contract
    bts_owner_manager = deploy_bts_owner_manager_contract(owner.address)
    sc += bts_owner_manager
    bts_core_contract = deploy_bts_core_contract(bts_owner_manager.address)
    sc += bts_core_contract

    helper_contract = deploy_helper_contract()
    sc += helper_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bts_periphery = deploy_bts_periphery_contract(bts_core_contract.address, helper_contract.address,
                                                  parse_address.address, owner.address)
    sc += bts_periphery
    fa2 = deploy_fa2_contract(bts_periphery.address)
    sc += fa2

    # Scenario 1: BTSCore Scenario

    # Test cases:
    # 1: verify owner manager 
    sc.verify_equal(bts_core_contract.data.bts_owner_manager, bts_owner_manager.address)

    # 2: verify update_bts_periphery function  
    bts_core_contract.update_bts_periphery(bts_periphery.address).run(sender=owner.address)
    sc.verify_equal(bts_core_contract.data.bts_periphery_address, sp.some(bts_periphery.address))

    # Scenario 2: set_fee_ratio

    # Test cases:
    # 1: token doesn't exist
    bts_core_contract.set_fee_ratio(name=sp.string("ABC"), fee_numerator=sp.nat(100), fixed_fee=sp.nat(10)).run(
        sender=owner.address, valid=False, exception='TokenNotExists')

    # 2: fee numerator is greater than denominator
    bts_core_contract.set_fee_ratio(name=sp.string("BTSCOIN"), fee_numerator=sp.nat(10000), fixed_fee=sp.nat(10)).run(
        sender=owner.address, valid=False, exception='InvalidSetting')

    # 3: fixed fee is 0
    bts_core_contract.set_fee_ratio(name=sp.string("BTSCOIN"), fee_numerator=sp.nat(100), fixed_fee=sp.nat(0)).run(
        sender=owner.address, valid=False, exception='LessThan0')

    # 4: valid condition
    bts_core_contract.set_fee_ratio(name=sp.string("BTSCOIN"), fee_numerator=sp.nat(100), fixed_fee=sp.nat(10)).run(
        sender=owner.address)

    # 5: set_fee_ratio called by non-owner
    bts_core_contract.set_fee_ratio(name=sp.string("BTSCOIN"), fee_numerator=sp.nat(100), fixed_fee=sp.nat(10)).run(
        sender=bob.address, valid=False, exception='Unauthorized')

    # 5: verify fees
    sc.verify_equal(bts_core_contract.data.coin_details["BTSCOIN"].fee_numerator, 100)
    sc.verify_equal(bts_core_contract.data.coin_details["BTSCOIN"].fixed_fee, 10)

    # Scenario 2: register

    # Test cases:
    # 1: native coin (native_coin = BTSCOIN)
    bts_core_contract.register(
        name=sp.string("BTSCOIN"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address, valid=False, exception='ExistNativeCoin')

    # 2:  fee numerator is greater than denominator
    bts_core_contract.register(
        name=sp.string("new_coin1"),
        fee_numerator=sp.nat(100000),
        fixed_fee=sp.nat(2),
        addr=fa2.address,
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address, valid=False, exception='InvalidSetting')

    # 3: valid case
    bts_core_contract.register(
        name=sp.string("new_coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=fa2.address,
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    # 4: verify the registered value
    sc.verify_equal(bts_core_contract.data.coins_name, ['new_coin', 'BTSCOIN'])
    sc.verify_equal(bts_core_contract.data.coin_details['new_coin'],
                    sp.record(addr=fa2.address, coin_type=2, fee_numerator=10, fixed_fee=2))

    # 5: existing coin name
    bts_core_contract.register(
        name=sp.string("new_coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=owner.address, valid=False, exception="ExistCoin")

    # 6: registered NON-NATIVE COIN_TYPE by non-owner
    bts_core_contract.register(
        name=sp.string("new_coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=bob.address, valid=False, exception="Unauthorized")

    # 7: register NON-NATIVE COIN_TYPE 
    bts_core_contract.register(
        name=sp.string("new_coin2"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=fa2.address,
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=owner.address, valid=False, exception="AddressExists")

    # 8: register NATIVE WRAPPED COIN_TYPE 
    bts_core_contract.register(
        name=sp.string("wrapped-coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"),
        token_metadata=sp.map({"token_metadata": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"metadata": sp.bytes("0x0dae11")})
    ).run(sender=owner.address)

    # Scenario 3: contract getters

    # Test cases:
    # 1: verify coin_id functio
    sc.verify_equal(bts_core_contract.coin_id('new_coin'), fa2.address)

    # 2: verify is_valid_coin function
    sc.verify_equal(bts_core_contract.is_valid_coin('new_coin'), True)
    sc.verify_equal(bts_core_contract.is_valid_coin('not_valid'), False)

    # 3: verify fee_ratio function
    sc.verify_equal(bts_core_contract.fee_ratio('new_coin'), sp.record(fee_numerator=10, fixed_fee=2))

    # 4: verify balance_of function
    sc.verify_equal(
        bts_core_contract.balance_of(
            sp.record(owner=owner.address, coin_name='new_coin')
        ),
        sp.record(usable_balance=0, locked_balance=0, refundable_balance=0, user_balance=0)
    )

    # 5: verify balance_of_batch function
    bts_core_contract.balance_of_batch(
        sp.record(owner=owner.address, coin_names=['new_coin', 'BTSCOIN'])
    )
    sc.verify_equal(
        bts_core_contract.balance_of_batch(
            sp.record(owner=owner.address, coin_names=['new_coin', 'BTSCOIN'])
        ),
        [
            sp.record(locked_balance=0, refundable_balance=0, usable_balance=0, user_balance=0),
            sp.record(locked_balance=0, refundable_balance=0, usable_balance=0, user_balance=0)
        ]
    )


def deploy_bts_core_contract(bts_owner_manager_contract):
    bts_core_contract = BTSCore.BTSCore(
        owner_manager=bts_owner_manager_contract,
        _native_coin_name="BTSCOIN",
        _fee_numerator=sp.nat(1000),
        _fixed_fee=sp.nat(10)
    )
    return bts_core_contract


def deploy_bts_owner_manager_contract(owner):
    bts_owner_manager_contract = BTSOwnerManager.BTSOwnerManager(owner)
    return bts_owner_manager_contract


def deploy_bts_periphery_contract(core_address, helper, parse, owner):
    bmc = sp.test_account("bmc")
    bts_periphery_contract = BTSPeriphery.BTSPeriphery(
        bmc_address=bmc.address, bts_core_address=core_address, helper_contract=helper, parse_address=parse,
        native_coin_name='BTSCoin', owner_address=owner)
    return bts_periphery_contract


def deploy_fa2_contract(admin_address):
    fa2_contract = BTSCore.FA2_contract.SingleAssetToken(admin=admin_address,
                                                         metadata=sp.big_map({"ss": sp.bytes("0x0dae11")}),
                                                         token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}))
    return fa2_contract


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address
