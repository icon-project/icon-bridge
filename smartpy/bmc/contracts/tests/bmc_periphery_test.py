import smartpy as sp

BMCManagement = sp.io.import_script_from_url("file:./contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:../bts/contracts/src/parse_address.py")


@sp.add_test("BMCPeripheryTest")
def test():
    sc = sp.test_scenario()

    # test account
    alice = sp.test_account("Alice")
    jack = sp.test_account("Jack")
    helper2 = sp.test_account("Helper2")
    owner = sp.test_account("owner")

    helper_contract = deploy_helper_contract()
    sc += helper_contract

    # deploy BMCManagement contract
    bmc_management_contract = deploy_bmc_management_contract(owner.address, helper_contract.address)
    sc += bmc_management_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bmc_periphery_contract = deploy_bmc_periphery_contract(
        bmc_management_contract.address, helper_contract.address, helper2.address, parse_address.address)
    sc += bmc_periphery_contract

    # Scenario 1: Contract setters

    # Test cases:
    # 1: set_bmc_btp_address by non-owner
    bmc_periphery_contract.set_bmc_btp_address("tezos.77").run(sender=alice, valid=False, exception="Unauthorized")

    # 2: set_bmc_btp_address by owner
    bmc_periphery_contract.set_bmc_btp_address("tezos.77").run(sender=bmc_management_contract.address)

    # 3: verify value for bmc_btp_address
    sc.verify(bmc_periphery_contract.data.bmc_btp_address ==
        sp.string("btp://tezos.77/KT1Tezooo3zzSmartPyzzSTATiCzzzseJjWC"))

    # 4: verify value for get_bmc_btp_address
    sc.verify_equal(bmc_periphery_contract.get_bmc_btp_address(),
                    sp.string("btp://tezos.77/KT1Tezooo3zzSmartPyzzSTATiCzzzseJjWC"))

    # 5: set_helper_address by non-owner
    bmc_periphery_contract.set_helper_address(helper_contract.address).run(sender=jack, valid=False,
                                                                           exception="Unauthorized")

    # 6: set_helper_address by non-owner
    bmc_periphery_contract.set_helper_address(helper_contract.address).run(sender=owner.address)

    # 7: verify helper address
    sc.verify(bmc_periphery_contract.data.helper == helper_contract.address)

    # 8: set_parse_address by non-owner
    bmc_periphery_contract.set_parse_address(parse_address.address).run(sender=jack, valid=False,
                                                                        exception="Unauthorized")

    # 8: set_parse_address by owner
    bmc_periphery_contract.set_parse_address(parse_address.address).run(sender=owner.address)

    # 9: verify parse_contract address
    sc.verify(bmc_periphery_contract.data.parse_contract == parse_address.address)

    # 10: set_bmc_management_addr by non-owner
    bmc_periphery_contract.set_bmc_management_addr(bmc_management_contract.address).run(sender=jack, valid=False,
                                                                                        exception="Unauthorized")

    # 11: set_bmc_management_addr by owner
    bmc_periphery_contract.set_bmc_management_addr(bmc_management_contract.address).run(
        sender=owner.address)

    # 9: verifying bmc_management address
    sc.verify(bmc_periphery_contract.data.bmc_management == bmc_management_contract.address)


def deploy_bmc_management_contract(owner, helper):
    bmc_management_contract = BMCManagement.BMCManagement(owner, helper)
    return bmc_management_contract


def deploy_bmc_periphery_contract(bmc_address, helper, helper2, parse):
    bmc_periphery_contract = BMCPeriphery.BMCPreiphery(bmc_address, helper, helper2, parse)
    return bmc_periphery_contract


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address
