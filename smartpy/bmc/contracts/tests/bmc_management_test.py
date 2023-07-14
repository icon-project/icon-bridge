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
    creator = sp.test_account("Creator")
    jack = sp.test_account("Jack")
    bob = sp.test_account("Bob")
    bmc_periphery_address = sp.test_account("bmc_periphery_address")
    creator2 = sp.test_account("creator2")
    service1_address = sp.test_account("service1_address")
    service2_address = sp.test_account("service2_address")
    ZERO_ADDRESS = sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")

    # deploy BMCManagement contract
    helper_contract = deploy_helper_contract()
    sc += helper_contract

    bmc_management_contract = deploy_bmc_management_contract(creator.address, helper_contract.address)
    sc += bmc_management_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bmc_periphery_address = bmc_periphery_address.address

    # Scenario 1: Contract setters

    # Test cases:
    # 1: set_bmc_periphery address
    sc.verify(bmc_management_contract.data.bmc_periphery == ZERO_ADDRESS)
    bmc_management_contract.set_bmc_periphery(bmc_periphery_address).run(sender=creator)

    # 2: sender non-owner
    bmc_management_contract.set_bmc_periphery(bob.address).run(sender=alice, valid=False, exception="Unauthorized")

    # 3: bmc_periphery already set
    bmc_management_contract.set_bmc_periphery(bmc_periphery_address).run(sender=creator, valid=False,
                                                                         exception="AlreadyExistsBMCPeriphery")
    # 4: Verify valid bmc_periphery  address
    sc.verify(bmc_management_contract.data.bmc_periphery != ZERO_ADDRESS)
    sc.verify(bmc_management_contract.data.bmc_periphery == bmc_periphery_address)

    # # 5: sender non-owner for set_bmc_btp_address
    # bmc_management_contract.set_bmc_btp_address("tezos.77").run(sender=alice, valid=False, exception="Unauthorized")

    # 6: sender is owner for set_bmc_btp_address
    bmc_management_contract.set_bmc_btp_address("tezos.77").run(sender=creator)

    # 7 : setting helper address by non-owner
    bmc_management_contract.set_helper_address(helper_contract.address).run(sender=jack, valid=False,
                                                                            exception="Unauthorized")

    # 8 : setting helper address by owner
    bmc_management_contract.set_helper_address(helper_contract.address).run(sender=creator)

    # 9: verifying address
    sc.verify(bmc_management_contract.data.helper == helper_contract.address)

    # Scenario 2: add / remove owner

    # Test cases:
    # 1: set owner by non-owner
    bmc_management_contract.add_owner(alice.address).run(sender=bob, valid=False, exception="Unauthorized")

    # 2: set new owner by owner
    bmc_management_contract.add_owner(alice.address).run(sender=creator)
    sc.verify(bmc_management_contract.data.owners[alice.address] == True)

    # 3: remove owner by non-owner
    bmc_management_contract.remove_owner(alice.address).run(sender=bob, valid=False, exception="Unauthorized")

    # 4: remove owner by owner
    bmc_management_contract.remove_owner(alice.address).run(sender=creator)
    sc.verify(~bmc_management_contract.data.owners.contains(jack.address))

    # 5: add owner
    bmc_management_contract.add_owner(creator2.address).run(sender=creator)

    # 6: verify is_owner
    sc.verify(bmc_management_contract.is_owner(creator2.address) == True)

    # Scenario 3: add / remove services and get_services

    # Test cases:
    svc1 = sp.string("bmc")
    svc2 = sp.string("bts")
    # 1: add service by non-owner
    bmc_management_contract.add_service(
        sp.record(addr=service1_address.address, svc=svc1)).run(sender=bob, valid=False, exception="Unauthorized")
    # 2: adding service by owner
    bmc_management_contract.add_service(sp.record(addr=service1_address.address, svc=svc1)).run(sender=creator)

    # 3: adding same service
    bmc_management_contract.add_service(
        sp.record(addr=service1_address.address, svc=svc1)).run(sender=creator, valid=False,
                                                                exception="AlreadyExistsBSH")

    # 4: remove service by non-owner
    bmc_management_contract.remove_service(svc2).run(sender=bob, valid=False, exception="Unauthorized")

    # 5: remove unregistered service
    bmc_management_contract.remove_service(svc2).run(sender=creator, valid=False)

    # 6: removing service
    bmc_management_contract.add_service(sp.record(addr=service2_address.address, svc=svc2)).run(sender=creator)
    bmc_management_contract.remove_service(svc2).run(sender=creator)

    # 7: verify get_services
    services = bmc_management_contract.get_services()
    sc.verify_equal(services, sp.map({svc1 : service1_address.address}))

    # Scenario 4: add / remove route and get_routes

    # Test case:
    dst = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest"
    next_link = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"
    # 1: adding route by non-owner
    bmc_management_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=bob, valid=False,
                                                                              exception="Unauthorized")

    # 2: adding route by owner
    bmc_management_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=creator)

    # 3: adding same routes
    bmc_management_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=creator, valid=False,
                                                                              exception="AlreadyExistRoute")

    dst1 = "btp://78.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5DEST1"
    next_link1 = "btp://78.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5LINK1"
    # 4: removing routes by non-owner
    bmc_management_contract.remove_route(dst).run(sender=bob, valid=False, exception="Unauthorized")

    # 5: removing non-exist routes
    bmc_management_contract.remove_route(dst1).run(sender=creator, valid=False, exception="NotExistRoute")

    # 6: removing existed routes by owner
    bmc_management_contract.add_route(sp.record(dst=dst1, link=next_link1)).run(sender=creator)
    bmc_management_contract.remove_route(dst1).run(sender=creator)

    # 7: verify get_routes
    get_routes = bmc_management_contract.get_routes()
    sc.verify_equal(get_routes, sp.map({0: sp.record(
        dst=sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest"),
        next=sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"))}))

    # Scenario 5: add / remove link and get_links

    # Test case:
    # 1: adding link by non-owner
    bmc_management_contract.add_link(
        "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(sender=bob, valid=False,
                                                                         exception="Unauthorized")

    # 2: adding link by owner
    bmc_management_contract.add_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(sender=creator)

    # 3: adding existed link
    bmc_management_contract.add_link(
        "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(sender=creator, valid=False,
                                                                         exception="AlreadyExistsLink")

    # 4: removing link by non-owner
    bmc_management_contract.remove_link(
        "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dead").run(sender=bob, valid=False,
                                                                         exception="Unauthorized")

    # 5: removing non-exist link
    bmc_management_contract.remove_link(
        "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dead").run(sender=creator, valid=False,
                                                                         exception="NotExistsLink")

    # 6: removing existed link by owner
    bmc_management_contract.add_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dead").run(sender=creator)
    bmc_management_contract.remove_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dead").run(sender=creator)

    # 7: verify get_links
    link_to_compare = bmc_management_contract.get_links()
    added_link = sp.list(['btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b'])
    sc.verify_equal(link_to_compare, added_link)

    # Scenario 6: set_link_rx_height

    # Test case:
    link = sp.string('btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b')
    height = sp.nat(2)
    # 1: non-exist link is given
    bmc_management_contract.set_link_rx_height(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnA",
                                               height=height).run(sender=creator,
                                                                  valid=False, exception="NotExistsKey")
    # 2: invalid height is given
    bmc_management_contract.set_link_rx_height(link=link, height=sp.nat(0)).run(sender=creator, valid=False,
                                                                                exception="InvalidRxHeight")

    # 3: set_link_rx_height by non-owner
    bmc_management_contract.set_link_rx_height(link=link, height=height).run(sender=bob, valid=False,
                                                                             exception="Unauthorized")

    # 4: set_link_rx_height by owner
    bmc_management_contract.set_link_rx_height(link=link, height=height).run(sender=creator)

    # 5: verify rx_height value
    sc.verify_equal(bmc_management_contract.data.links[link].rx_height, 2)

    # Scenario 7: set_link

    # Test case:
    block_interval = sp.nat(2)
    _max_aggregation = sp.nat(3)
    delay_limit = sp.nat(2)
    # 1: setting link by non-owner
    bmc_management_contract.set_link(
        sp.record(_link=link, block_interval=block_interval, _max_aggregation=_max_aggregation,
                  delay_limit=delay_limit)).run(sender=bob, valid=False, exception="Unauthorized")

    # 2: setting non-exist link
    bmc_management_contract.set_link(
        sp.record(_link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnZ", block_interval=block_interval,
                  _max_aggregation=_max_aggregation, delay_limit=delay_limit)).run(sender=creator, valid=False,
                                                                                   exception="NotExistsLink")
    # 3: setting link with invalid paramter
    bmc_management_contract.set_link(
        sp.record(_link=link, block_interval=block_interval, _max_aggregation=sp.nat(0),
                  delay_limit=delay_limit)).run(
        sender=creator, valid=False, exception="InvalidParam")
    bmc_management_contract.set_link(
        sp.record(_link=link, block_interval=block_interval, _max_aggregation=_max_aggregation,
                  delay_limit=sp.nat(0))).run(sender=creator, valid=False, exception="InvalidParam")

    # 4: setting link with valid paramter by owner
    bmc_management_contract.set_link(
        sp.record(_link=link, block_interval=block_interval, _max_aggregation=_max_aggregation,
                  delay_limit=delay_limit)).run(sender=creator)

    # Scenario 8: add / remove relay and get_relays

    # Test case:
    # 1: adding relay by non-owner
    bmc_management_contract.add_relay(
        sp.record(link=link, addr=sp.set(
            [sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD")]))).run(sender=bob, valid=False,
                                                                        exception="Unauthorized")

    # 2: adding relay to non-exist link
    bmc_management_contract.add_relay(sp.record(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUNONLINK",
                                                addr=sp.set([sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD")]))).run(
        sender=creator, valid=False, exception="NotExistsLink")

    # 3: adding relay by owner
    bmc_management_contract.add_relay(
        sp.record(link=link, addr=sp.set([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")]))).run(sender=creator)

    # 4: remove relay by non-owner
    bmc_management_contract.remove_relay(
        sp.record(link=link, addr=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))).run(sender=bob, valid=False,
                                                                                           exception="Unauthorized")

    # 5: removing relay with non-exist link
    bmc_management_contract.remove_relay(sp.record(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUNONLINK",
                                                   addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD"))).run(
        sender=creator, valid=False, exception="NotExistsLink")

    # 6: removing relay by owner
    next_link1 = sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a625link1")
    bmc_management_contract.add_link(next_link1).run(sender=creator)
    bmc_management_contract.add_relay(
        sp.record(link=next_link1, addr=sp.set([sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxADD1")]))).run(
        sender=creator)
    bmc_management_contract.remove_relay(
        sp.record(link=next_link1, addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxADD1"))).run(sender=creator)

    # 7: verify get_relays
    compared_to = sp.list([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")])
    get_relays = bmc_management_contract.get_relays(link)
    sc.verify_equal(get_relays, compared_to)

    # Scenario 9: Contract getters

    # Test cases:
    # 1: verify get_bsh_service_by_name
    get_bsh_service_by_name = bmc_management_contract.get_bsh_service_by_name(svc1)
    sc.verify_equal(get_bsh_service_by_name, service1_address.address)

    # 2: verify get_link
    get_link = bmc_management_contract.get_link(link)
    data = sp.record(
        relays=sp.set([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")]),
        reachable=sp.set([]),
        rx_seq=sp.nat(0),
        tx_seq=sp.nat(0),
        block_interval_src=sp.nat(30000),
        block_interval_dst=sp.nat(2),
        max_aggregation=sp.nat(3),
        delay_limit=sp.nat(2),
        relay_idx=sp.nat(0),
        rotate_height=sp.nat(0),
        rx_height=sp.nat(0),
        rx_height_src=sp.nat(0),
        is_connected=True
    )
    sc.verify_equal(get_link, data)

    # 3: verify get_link_rx_seq
    get_link_rx_seq = bmc_management_contract.get_link_rx_seq(link)
    sc.verify_equal(get_link_rx_seq, 0)

    # 4: verify get_link_tx_seq
    get_link_tx_seq = bmc_management_contract.get_link_tx_seq(link)
    sc.verify_equal(get_link_tx_seq, 0)

    # 5: verify get_link_rx_height
    get_link_rx_height = bmc_management_contract.get_link_rx_height(link)
    sc.verify_equal(get_link_rx_height, 0)

    # 6: verify get_link_relays
    get_link_relays = bmc_management_contract.get_link_relays(link)
    sc.verify_equal(get_link_relays, compared_to)

    # 7: verify get_relay_status_by_link
    get_link_relays = bmc_management_contract.get_relay_status_by_link(link)
    sc.verify_equal(get_link_relays, sp.map(
        {0: sp.record(addr=sp.address('tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9'), block_count=0, msg_count=0)}))

    # Scenario 10: update_link_rx_seq function

    # Test cases:
    # 1: update_link_rx_seq by non-owner
    bmc_management_contract.update_link_rx_seq(
        sp.record(prev=next_link1, val=sp.nat(3))).run(sender=creator, valid=False, exception="Unauthorized")

    # 2: update_link_rx_seq by owner
    bmc_management_contract.update_link_rx_seq(sp.record(prev=next_link1, val=sp.nat(3))).run(
        sender=bmc_periphery_address)

    # 3: verifying value
    sc.verify_equal(bmc_management_contract.data.links[next_link1].rx_seq, 3)

    # Scenario 11: update_link_tx_seq function

    # Test cases:
    # 1: update_link_tx_seq by non-bmc_periphery
    bmc_management_contract.update_link_tx_seq(sp.record(prev=next_link1, serialized_msg=sp.bytes("0x64"))).run(
        sender=creator, valid=False, exception="Unauthorized")

    # 2: update_link_tx_seq by bmc_periphery
    bmc_management_contract.update_link_tx_seq(sp.record(prev=next_link1, serialized_msg=sp.bytes("0x64"))).run(
        sender=bmc_periphery_address)

    # 3: verifying value
    sc.verify_equal(bmc_management_contract.data.links[next_link1].tx_seq, 1)

    # Scenario 12: update_link_rx_height function

    # Test cases:
    # 1: update_link_rx_height by non-bmc_periphery
    bmc_management_contract.update_link_rx_height(
        sp.record(prev=next_link1, val=sp.nat(3))).run(sender=creator, valid=False, exception="Unauthorized")

    # 2: update_link_rx_height by bmc_periphery
    bmc_management_contract.update_link_rx_height(sp.record(prev=next_link1, val=sp.nat(4))).run(
        sender=bmc_periphery_address)

    # 3: verifying value
    sc.verify_equal(bmc_management_contract.data.links[next_link1].rx_height, 4)

    # Scenario 13: update_link_reachable and delete_link_reachable function

    # Test cases:
    to = sp.list(["btp://net1/addr1", "btp://net2/addr2"])
    # 1: update_link_reachable by non-bmc_periphery
    bmc_management_contract.update_link_reachable(sp.record(prev=next_link1, to=to)).run(sender=creator, valid=False,
                                                                                         exception="Unauthorized")

    # 2: update_link_reachable by bmc_periphery
    bmc_management_contract.update_link_reachable(sp.record(prev=next_link1, to=to)).run(sender=bmc_periphery_address)

    # 3: verifying value
    sc.verify_equal(bmc_management_contract.data.links[next_link1].reachable,
                    sp.set(['btp://net1/addr1', 'btp://net2/addr2']))

    # 4: delete_link_reachable by non-bmc_periphery
    bmc_management_contract.delete_link_reachable(
        sp.record(prev=next_link1, index=sp.nat(0))).run(sender=creator, valid=False, exception="Unauthorized")

    # 5: delete_link_reachable by bmc_periphery
    bmc_management_contract.delete_link_reachable(sp.record(prev=next_link1, index=sp.nat(0))).run(
        sender=bmc_periphery_address)

    # 6: verifying value
    sc.verify_equal(bmc_management_contract.data.links[next_link1].reachable, sp.set(['btp://net2/addr2']))

    # # 7: delete non-exist link
    # next_link2 = sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a625link2")
    # bmc_management_contract.delete_link_reachable(sp.record(prev=next_link2, index=sp.nat(0))).run(
    #     sender=bmc_periphery_address)

    # Scenario 13: update_relay_stats and resolve_route function

    # Test cases:
    # 1: update_relay_stats by non-bmc_periphery
    bmc_management_contract.update_relay_stats(
        sp.record(relay=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiADD"), block_count_val=sp.nat(2),
                  msg_count_val=sp.nat(2))).run(sender=creator, valid=False, exception="Unauthorized")

    # 2: update_relay_stats by bmc_periphery
    bmc_management_contract.update_relay_stats(
        sp.record(relay=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), block_count_val=sp.nat(2),
                  msg_count_val=sp.nat(2))).run(sender=bmc_periphery_address)

    # 3: verifying value
    sc.verify_equal(
        bmc_management_contract.data.relay_stats[sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")].block_count, 2)
    sc.verify_equal(
        bmc_management_contract.data.relay_stats[sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")].msg_count, 2)

    # 4: verifying value of resolve_route function
    sc.verify_equal(bmc_management_contract.resolve_route(sp.string('0x7.icon')),
                    sp.pair('btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b',
                            'btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest'))

    # 4: non-exist link resolve_route
    sc.verify_equal(bmc_management_contract.resolve_route(sp.string('0x8.icon')),
                    sp.pair('Unreachable',
                            'Unreachable: 0x8.icon is unreachable'))


def deploy_bmc_management_contract(owner, helper):
    bmc_management_contract = BMCManagement.BMCManagement(owner, helper)
    return bmc_management_contract


def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address
