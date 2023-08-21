import smartpy as sp

BMCManagement = sp.io.import_script_from_url("file:./contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:./contracts/src/parse_address.py")





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

    bmcManagement_contract = deploy_bmcManagement_contract(creator.address, helper_contract.address)
    sc += bmcManagement_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bmcPeriphery_contract = deploy_bmcPeriphery_contract(bmcManagement_contract.address, helper_contract.address, parse_address.address)
    sc += bmcPeriphery_contract



    bmc_periphery_address = bmcPeriphery_contract.address


    # Test case 1: bmc_periphery
    sc.h1("Test case 1: set bmc_periphery to a valid address")
    sc.verify(bmcManagement_contract.data.bmc_periphery.is_some() == False)
    bmcManagement_contract.set_bmc_periphery(bmc_periphery_address).run(sender=creator)
    # sender should be owner
    bmcManagement_contract.set_bmc_periphery(bob.address).run(sender=alice, valid=False, exception="Unauthorized")
    # repeated bmc_periphery should throw error
    bmcManagement_contract.set_bmc_periphery(bmc_periphery_address).run(sender=creator, valid=False, exception="AlreadyExistsBMCPeriphery")
    # Verify that bmc_periphery is set to the valid address
    sc.verify(bmcManagement_contract.data.bmc_periphery.is_some() == True)
    sc.verify(bmcManagement_contract.data.bmc_periphery.open_some() == bmc_periphery_address)

    # set_bmc_btp_address
    bmcManagement_contract.set_bmc_btp_address("tezos.77").run(sender=creator)


    # Test case 2: add_owner
    # throw error when adding owner by random address
    bmcManagement_contract.add_owner(alice.address).run(sender=bob, valid=False, exception="Unauthorized")
    # successfully added new owner
    bmcManagement_contract.add_owner(alice.address).run(sender=creator)
    sc.verify(bmcManagement_contract.data.owners[alice.address] == True)


    # Test case 3: remove owner
    # throw error when removing owner by random address
    bmcManagement_contract.remove_owner(alice.address).run(sender=bob, valid=False, exception="Unauthorized")
    # working
    bmcManagement_contract.remove_owner(alice.address).run(sender=creator)
    sc.verify(~bmcManagement_contract.data.owners.contains(jack.address))


    # Test case 4: is_owner
    # Add an owner
    bmcManagement_contract.add_owner(creator2.address).run(sender=creator)
    # Test the is_owner view function
    sc.verify(bmcManagement_contract.is_owner(creator2.address) == True)


    # Test case 5: add_service function
    svc1 = sp.string("service1")
    svc2 = sp.string("service2")
    svc3 = sp.string("service3")
    # add service by random address should fail
    bmcManagement_contract.add_service(sp.record(addr=service1_address.address, svc=svc1)).run(sender=bob, valid=False, exception="Unauthorized")
    # adding service
    bmcManagement_contract.add_service(sp.record(addr=service1_address.address, svc=svc1)).run(sender=creator)
    # shouldn't add same service twice
    bmcManagement_contract.add_service(sp.record(addr=service1_address.address, svc=svc1)).run(sender=creator, valid=False, exception="AlreadyExistsBSH")


    # Test case 6: remove_service function
    # remove service by random address should fail
    bmcManagement_contract.remove_service(svc2).run(sender=bob, valid=False, exception="Unauthorized")
    # removing unregistered should throw error
    bmcManagement_contract.remove_service(svc3).run(sender=creator, valid=False)
    # removing service
    bmcManagement_contract.add_service(sp.record(addr=service2_address.address, svc=svc2)).run(sender=creator)
    bmcManagement_contract.remove_service(svc2).run(sender=creator)


    #test case 7: get_services function
    services = bmcManagement_contract.get_services()
    sc.verify_equal(services, sp.map({0: sp.record(svc=svc1, addr=service1_address.address)}))

    # test case 13: add_route function
    dst = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest"
    next_link = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"
    # only owner can add routes
    bmcManagement_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=bob, valid=False, exception="Unauthorized")
    # should work
    bmcManagement_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=creator)
    # cannot add already existed route
    bmcManagement_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=creator, valid=False, exception="AlreadyExistRoute")


    # # test case 14: remove_route function
    # dst1 = "btp://78.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5DEST1"
    # next_link1 = "btp://78.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5LINK1"
    # # only owner can remove routes
    # bmcManagement_contract.remove_route(dst).run(sender=bob, valid=False,exception="Unauthorized")
    # # throw error when non-exist route is given & this should throw error but not thrown
    # bmcManagement_contract.remove_route(dst1).run(sender=creator, valid=False, exception="NotExistRoute")
    # # should work
    # bmcManagement_contract.add_route(sp.record(dst=dst1, link=next_link1)).run(sender=creator)
    # bmcManagement_contract.remove_route(dst1).run(sender=creator)


    # # test case 15: get_routes function
    # get_routes = bmcManagement_contract.get_routes()
    # sc.verify_equal(get_routes, sp.map({0: sp.record(dst=sp.string("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hDEST"), next=sp.string("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"))}))


    # test case 8: add_link function
    # add_link by random address should fail
    bmcManagement_contract.add_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW").run(sender=bob, valid=False, exception="Unauthorized")
    # should work
    bmcManagement_contract.add_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(sender=creator)
    # add_link by of same link should fail
    bmcManagement_contract.add_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(sender=creator, valid=False, exception="AlreadyExistsLink")

    #
    # # test case 9: remove_link function
    # # remove_link by random address should fail
    # bmcManagement_contract.remove_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnB").run(sender=bob, valid=False, exception="Unauthorized")
    # # remove_link should throw error when removing non-existing link
    # bmcManagement_contract.remove_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnZ").run(sender=creator, valid=False, exception="NotExistsLink")
    # # should work
    # bmcManagement_contract.add_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnA").run(sender=creator)
    # bmcManagement_contract.remove_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnA").run(sender=creator)
    #
    #
    # # test case 10: get_links function
    # link_to_compare = bmcManagement_contract.get_links()
    # added_link = sp.list(['btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'])
    # sc.verify_equal(link_to_compare, added_link)
    #
    #
    # # test case 11: set_link_rx_height
    # link = sp.string('btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW')
    # height = sp.nat(2)
    # # error when not exist link is given
    # bmcManagement_contract.set_link_rx_height(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnA", height=height).run(sender=creator, valid=False, exception="NotExistsKey")
    # # error when not invalid height is given
    # bmcManagement_contract.set_link_rx_height(link=link, height=sp.nat(0)).run(sender=creator,valid=False,exception="InvalidRxHeight")
    # # should work
    # bmcManagement_contract.set_link_rx_height(link=link, height=height).run(sender=creator)
    # sc.verify_equal(bmcManagement_contract.data.links[link].rx_height, 2)
    #
    #
    # # test case 12: set_link function
    # block_interval = sp.nat(2)
    # _max_aggregation = sp.nat(3)
    # delay_limit = sp.nat(2)
    # # only owner should set link
    # bmcManagement_contract.set_link(sp.record(_link=link, block_interval=block_interval,_max_aggregation=_max_aggregation, delay_limit=delay_limit)).run(sender=bob, valid=False, exception="Unauthorized")
    # # error when link doesnt exist
    # bmcManagement_contract.set_link(sp.record(_link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnZ", block_interval=block_interval,_max_aggregation=_max_aggregation, delay_limit=delay_limit)).run(sender=creator, valid=False, exception="NotExistsLink")
    # # error when invalid paramters were given
    # bmcManagement_contract.set_link(sp.record(_link=link, block_interval=block_interval,_max_aggregation=sp.nat(0), delay_limit=delay_limit)).run(sender=creator, valid=False, exception="InvalidParam")
    # bmcManagement_contract.set_link(sp.record(_link=link, block_interval=block_interval,_max_aggregation=_max_aggregation, delay_limit=sp.nat(0))).run(sender=creator, valid=False, exception="InvalidParam")
    # # should work
    # bmcManagement_contract.set_link(sp.record(_link=link, block_interval=block_interval,_max_aggregation=_max_aggregation, delay_limit=delay_limit)).run(sender=creator)
    #
    #
    # # test case 16: add_relay function
    # # only owner should add relay
    # bmcManagement_contract.add_relay(sp.record(link=link, addr=sp.set([sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD")]))).run(sender=bob, valid=False, exception="Unauthorized")
    # # fail when non-exist link is given
    # bmcManagement_contract.add_relay(sp.record(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUNONLINK", addr=sp.set([sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD")]))).run(sender=creator, valid=False, exception="NotExistsLink")
    # # should work
    # bmcManagement_contract.add_relay(sp.record(link=link, addr=sp.set([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")]))).run(sender=creator)
    #
    #
    # # test case 17: remove_relay function
    # # only owner should remove relay
    # bmcManagement_contract.remove_relay(sp.record(link=link, addr=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))).run(sender=bob, valid=False, exception="Unauthorized")
    # # fail when non-exist link is given
    # bmcManagement_contract.remove_relay(sp.record(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUNONLINK", addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD"))).run(sender=creator, valid=False, exception="NotExistsLink")
    # # should work
    # next_link1 = sp.string("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5LINK1")
    # bmcManagement_contract.add_link(next_link1).run(sender=creator)
    # bmcManagement_contract.add_relay(sp.record(link=next_link1, addr=sp.set([sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxADD1")]))).run(sender=creator)
    # bmcManagement_contract.remove_relay(sp.record(link=next_link1, addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxADD1"))).run(sender=creator)
    #
    #
    # # test case 18: get_relays function
    # compared_to = sp.list([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")])
    # get_relays = bmcManagement_contract.get_relays(link)
    # sc.verify_equal(get_relays, compared_to)
    #
    #
    # # test case 19: get_bsh_service_by_name function
    # get_bsh_service_by_name = bmcManagement_contract.get_bsh_service_by_name(svc1)
    # sc.verify_equal(get_bsh_service_by_name, service1_address.address)
    #
    #
    # # # test case 20: get_link function
    # get_link = bmcManagement_contract.get_link(link)
    # data = sp.record(
    #         relays=sp.set([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")]),
    #         reachable=sp.set([]),
    #         rx_seq=sp.nat(0),
    #         tx_seq=sp.nat(7),
    #         block_interval_src=sp.nat(2),
    #         block_interval_dst=sp.nat(0),
    #         max_aggregation=sp.nat(3),
    #         delay_limit=sp.nat(2),
    #         relay_idx=sp.nat(0),
    #         rotate_height=sp.nat(0),
    #         rx_height=sp.nat(0),
    #         rx_height_src=sp.nat(0),
    #         is_connected=True
    #     )
    # sc.verify_equal(get_link, data)
    #
    #
    # # test case 21: get_link_rx_seq function
    # get_link_rx_seq = bmcManagement_contract.get_link_rx_seq(link)
    # sc.verify_equal(get_link_rx_seq, 0)
    #
    #
    # # test case 22: get_link_tx_seq function
    # get_link_tx_seq = bmcManagement_contract.get_link_tx_seq(link)
    # sc.verify_equal(get_link_tx_seq, 7)
    #
    #
    # # test case 23: get_link_rx_height function
    # get_link_rx_height = bmcManagement_contract.get_link_rx_height(link)
    # sc.verify_equal(get_link_rx_height, 0)
    #
    #
    # # test case 24: get_link_relays function
    # get_link_relays = bmcManagement_contract.get_link_relays(link)
    # sc.verify_equal(get_link_relays, compared_to)
    #
    #
    # # test case 25: get_relay_status_by_link function
    # get_link_relays = bmcManagement_contract.get_relay_status_by_link(link)
    # sc.verify_equal(get_link_relays, sp.map({0: sp.record(addr=sp.address('tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9'), block_count=0, msg_count=0)}))
    #
    #
    # # test case 26: update_link_rx_seq function
    # #only bmc periphery address can run other users should get error
    # bmcManagement_contract.update_link_rx_seq(sp.record(prev=next_link1, val=sp.nat(3))).run(sender=creator, valid=False, exception="Unauthorized")
    # #working
    # bmcManagement_contract.update_link_rx_seq(sp.record(prev=next_link1, val=sp.nat(3))).run(sender=bmc_periphery_address)
    # # Check that the value of rx_seq is updated correctly
    # sc.verify_equal(bmcManagement_contract.data.links[next_link1].rx_seq, 3)
    #
    #
    # # test case 27: update_link_tx_seq function
    # # only bmc periphery address can run other users should get error
    # bmcManagement_contract.update_link_tx_seq(next_link1).run(sender=creator,  valid=False, exception="Unauthorized")
    # # working
    # bmcManagement_contract.update_link_tx_seq(next_link1).run(sender=bmc_periphery_address)
    # # Check that the value of tx_seq is updated correctly
    # sc.verify_equal(bmcManagement_contract.data.links[next_link1].tx_seq, 1)
    #
    #
    # # test case 28: update_link_rx_height function
    # #only bmc periphery address can run other users should get error
    # bmcManagement_contract.update_link_rx_height(sp.record(prev=next_link1, val=sp.nat(3))).run(sender=creator, valid=False, exception="Unauthorized")
    # #working
    # bmcManagement_contract.update_link_rx_height(sp.record(prev=next_link1, val=sp.nat(4))).run(sender=bmc_periphery_address)
    # # Check that the value of rx_seq is updated correctly
    # sc.verify_equal(bmcManagement_contract.data.links[next_link1].rx_height, 4)
    #
    #
    # # test case 29: update_link_reachable function
    # to = sp.list(["btp://net1/addr1", "btp://net2/addr2"])
    # # only bmc periphery address can run other users should get error
    # bmcManagement_contract.update_link_reachable(sp.record(prev=next_link1, to=to)).run(sender=creator, valid=False, exception="Unauthorized")
    # # should work
    # bmcManagement_contract.update_link_reachable(sp.record(prev=next_link1, to=to)).run(sender=bmc_periphery_address)
    # # value checking
    # sc.verify_equal(bmcManagement_contract.data.links[next_link1].reachable, sp.set(['btp://net1/addr1', 'btp://net2/addr2']))
    #
    #
    # # test case 30: delete_link_reachable function
    # #only bmc periphery address can run other users should get error
    # bmcManagement_contract.delete_link_reachable(sp.record(prev=next_link1, index=sp.nat(0))).run(sender=creator, valid=False, exception="Unauthorized")
    # #working
    # bmcManagement_contract.delete_link_reachable(sp.record(prev=next_link1, index=sp.nat(0))).run(sender=bmc_periphery_address)
    # # value checking
    # sc.verify_equal(bmcManagement_contract.data.links[next_link1].reachable, sp.set(['btp://net2/addr2']))
    #
    #
    # # test case 31: update_relay_stats function
    # #only bmc periphery address can run other users should get error
    # bmcManagement_contract.update_relay_stats(sp.record(relay=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD"), block_count_val=sp.nat(2), msg_count_val=sp.nat(2))).run(sender=creator, valid=False, exception="Unauthorized")
    # #working
    # bmcManagement_contract.update_relay_stats(sp.record(relay=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), block_count_val=sp.nat(2), msg_count_val=sp.nat(2))).run(sender=bmc_periphery_address)
    # # value checking
    # sc.verify_equal(bmcManagement_contract.data.relay_stats[sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")].block_count, 2)
    # sc.verify_equal(bmcManagement_contract.data.relay_stats[sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")].msg_count, 2)
    #
    #
    # # test case 32: resolve_route function
    # sc.verify_equal(bmcManagement_contract.resolve_route(sp.string('77.tezos')), sp.pair('btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW', 'btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hDEST'))
    #



def deploy_bmcManagement_contract(owner, helper):
    bmcManagement_contract = BMCManagement.BMCManagement(owner, helper)
    return bmcManagement_contract

def deploy_bmcPeriphery_contract(bmc_addres, helper, parse):
    owner = sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")
    bmcPeriphery_contract = BMCPeriphery.BMCPreiphery(bmc_addres, helper, parse, owner)
    return bmcPeriphery_contract

def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address