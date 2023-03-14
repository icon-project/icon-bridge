import smartpy as sp

BMCManagement = sp.io.import_script_from_url("file:./contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:../bts/contracts/src/parse_address.py")


@sp.add_test("BMCManagementTest")
def test():
    sc = sp.test_scenario()

    # test account
    alice = sp.test_account("Alice")
    jack = sp.test_account("Jack")
    bob = sp.test_account("Bob")
    helper2 = sp.test_account("Helper2")
    owner = sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")

    helper_contract = deploy_helper_contract()
    sc += helper_contract

    # deploy BMCManagement contract
    bmc_management_contract = deploy_bmc_management_contract(sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), helper_contract.address)
    sc += bmc_management_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bmc_periphery_contract = deploy_bmc_periphery_contract(bmc_management_contract.address, helper_contract.address, helper2.address, parse_address.address, owner)
    sc += bmc_periphery_contract

    # set_bmc_btp_address
    bmc_periphery_contract.set_bmc_btp_address("tezos.77").run(sender=alice, valid=False, exception="Unauthorized")
    bmc_periphery_contract.set_bmc_btp_address("tezos.77").run(sender=bmc_management_contract.address)
    sc.verify(bmc_periphery_contract.data.bmc_btp_address == sp.some(sp.string("btp://tezos.77/KT1Tezooo3zzSmartPyzzSTATiCzzzseJjWC")))

    # get_bmc_btp_address
    sc.verify_equal(bmc_periphery_contract.get_bmc_btp_address(), sp.string("btp://tezos.77/KT1Tezooo3zzSmartPyzzSTATiCzzzseJjWC"))

    # set_helper_address
    bmc_periphery_contract.set_helper_address(sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXCXPzhelper")).run(sender=jack, valid=False, exception="Unauthorized")
    bmc_periphery_contract.set_helper_address(sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXCXPzhelper")).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    sc.verify(bmc_periphery_contract.data.helper == sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXCXPzhelper"))

    # set_parse_address
    bmc_periphery_contract.set_parse_address(sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXCXPzhparse")).run(sender=jack, valid=False, exception="Unauthorized")
    bmc_periphery_contract.set_parse_address(sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXCXPzhparse")).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    sc.verify(bmc_periphery_contract.data.parse_contract == sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXCXPzhparse"))

    # set_bmc_management_addr
    bmc_periphery_contract.set_bmc_management_addr(sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXmanagement")).run(sender=jack, valid=False, exception="Unauthorized")
    bmc_periphery_contract.set_bmc_management_addr(sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXmanagement")).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    sc.verify(bmc_periphery_contract.data.bmc_management == sp.address("KT1EXYXNGdbh4uvdKc8hh7ETQXmanagement"))


def deploy_bmc_management_contract(owner, helper):
    bmc_management_contract = BMCManagement.BMCManagement(owner, helper)
    return bmc_management_contract


def deploy_bmc_periphery_contract(bmc_address, helper,helper2, parse, owner):
    bmc_periphery_contract = BMCPeriphery.BMCPreiphery(bmc_address, helper, helper2, parse, owner)
    return bmc_periphery_contract


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address
