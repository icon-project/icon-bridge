import smartpy as sp

BTSPeriphery = sp.io.import_script_from_url("file:./contracts/src/bts_periphery.py")
BTSCore = sp.io.import_script_from_url("file:./contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:./contracts/src/bts_owner_manager.py")
ParseAddress = sp.io.import_script_from_url("file:./contracts/src/parse_address.py")
Helper = sp.io.import_script_from_url("file:./contracts/src/helper.py")


@sp.add_test("BTSPeripheryTest")
def test():
    sc = sp.test_scenario()

    # test account
    admin = sp.test_account('admin')

    def deploy_bts_periphery_contract(core_address, helper, parse):
        bmc = sp.test_account("bmc")
        _bts_periphery_contract = BTSPeriphery.BTSPeriphery(
            bmc_address=bmc.address, bts_core_address=core_address, helper_contract=helper, parse_address=parse,
            native_coin_name='NativeCoin', owner_address=admin.address)
        return _bts_periphery_contract

    def deploy_parse_contract():
        _bts_parse_contract = ParseAddress.ParseAddress()
        return _bts_parse_contract

    def deploy_helper():
        _bts_helper = Helper.Helper()
        return _bts_helper

    def deploy_bts_core_contract(_bts_owner_manager_contract):
        _bts_core_contract = BTSCore.BTSCore(
            owner_manager=_bts_owner_manager_contract,
            _native_coin_name="Tok1",
            _fee_numerator=sp.nat(1000),
            _fixed_fee=sp.nat(10))
        return _bts_core_contract

    def deploy_bts_owner_manager_contract():
        _bts_owner_manager_contract = BTSOwnerManager.BTSOwnerManager(admin.address)
        return _bts_owner_manager_contract

    bts_owner_manager_contract = deploy_bts_owner_manager_contract()
    sc += bts_owner_manager_contract

    bts_core_contract = deploy_bts_core_contract(bts_owner_manager_contract.address)
    sc += bts_core_contract

    bts_helper_contract = deploy_helper()
    sc += bts_helper_contract

    bts_parse_contract = deploy_parse_contract()
    sc += bts_parse_contract

    # deploy bts_periphery contract
    bts_periphery_contract = deploy_bts_periphery_contract(bts_core_contract.address, bts_helper_contract.address,
                                                           bts_parse_contract.address)
    sc += bts_periphery_contract

    bts_core_contract.update_bts_periphery(bts_periphery_contract.address).run(sender=admin.address)

    # Scenario 1: set token  limit

    # Test cases:
    # 1: set_token_limit from non-bts_periphery contract
    bts_periphery_contract.set_token_limit(sp.map({"Tok2": sp.nat(5), "BB": sp.nat(2)})).run(sender=admin.address,
                                                                                             valid=False,
                                                                                             exception='Unauthorized')

    # 2: set token limit for Tok2 coin to 5 and BB coin to 2 from bts_periphery_contract
    bts_periphery_contract.set_token_limit(sp.map({"Tok2": sp.nat(5), "BB": sp.nat(2)})).run(
        sender=bts_core_contract.address)

    # 3: verifying the value of token limit
    sc.verify(bts_periphery_contract.data.token_limit["Tok2"] == sp.nat(5))

    # 4: modify already set data
    bts_periphery_contract.set_token_limit(sp.map({"Tok2": sp.nat(15), "BB": sp.nat(22)})).run(
        sender=bts_core_contract.address)

    # 5: verifying the value of token limit after change
    sc.verify(bts_periphery_contract.data.token_limit["BB"] == sp.nat(22))

    # Scenario 2: send service message

    # Test cases:
    # 1 : send_service_message called by bts-core
    bts_periphery_contract.send_service_message(
        sp.record(
            _from=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"),
            to="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW",
            coin_details=[
                sp.record(coin_name="Tok1", value=sp.nat(10), fee=sp.nat(2))
            ]
        )).run(sender=bts_core_contract.address).run(sender=bts_core_contract.address)

    # 2 : send_service_message called by non-bts-core
    bts_periphery_contract.send_service_message(
        sp.record(
            _from=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"),
            to="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW",
            coin_details=[
                sp.record(coin_name="Tok1", value=sp.nat(10), fee=sp.nat(2))
            ]
        )).run(sender=admin, valid=False, exception='Unauthorized')

    # 3: verify if request message is correct
    sc.show(bts_periphery_contract.data.requests[1])
    sc.verify_equal(
        bts_periphery_contract.data.requests[1],
        sp.record(
            coin_details=[
                sp.record(
                    coin_name='Tok1',
                    value=10,
                    fee=2
                )
            ],
            from_='tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW',
            to='btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'))

    # 4: verify request data
    sc.verify(bts_periphery_contract.data.number_of_pending_requests == 1)
    sc.verify(bts_periphery_contract.data.serial_no == 1)
