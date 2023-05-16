import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")
rlp_encode = sp.io.import_script_from_url("file:./contracts/src/RLP_encode_struct.py")
# bmc_periphery_file = sp.io.import_script_from_url(
#     "file:./contracts/src/bmc_periphery.py")
# helper_file = sp.io.import_script_from_url(
#     "file:./contracts/src/helper.py")
# parse_address_file = sp.io.import_script_from_url(
#     "file:./contracts/src/parse_address.py")

class BMCManagement(sp.Contract, rlp_encode.EncodeLibrary):
    BLOCK_INTERVAL_MSEC = sp.nat(1000)

    def __init__(self, owner_address, helper_contract):
        self.init(
            owners=sp.map(l={owner_address:True}),
            number_of_owners=sp.nat(1),
            bsh_services=sp.map(),
            relay_stats=sp.map(),
            routes=sp.map(),
            links=sp.map(),
            list_bsh_names=sp.set(),
            list_route_keys=sp.set(),
            list_link_names=sp.set(),
            bmc_periphery=sp.none,
            serial_no=sp.nat(0),
            addrs=sp.set(),
            get_route_dst_from_net=sp.map(),
            get_link_from_net=sp.map(),
            get_link_from_reachable_net=sp.map(),
            helper=helper_contract
        )

        self.init_type(sp.TRecord(
            owners=sp.TMap(sp.TAddress, sp.TBool),
            number_of_owners=sp.TNat,
            bsh_services=sp.TMap(sp.TString, sp.TAddress),
            relay_stats=sp.TMap(sp.TAddress, types.Types.RelayStats),
            routes=sp.TMap(sp.TString, sp.TString),
            links=sp.TMap(sp.TString, types.Types.Link),
            list_bsh_names=sp.TSet(sp.TString),
            list_route_keys=sp.TSet(sp.TString),
            list_link_names=sp.TSet(sp.TString),
            bmc_periphery=sp.TOption(sp.TAddress),
            serial_no=sp.TNat,
            addrs=sp.TSet(sp.TAddress),
            get_route_dst_from_net=sp.TMap(sp.TString, sp.TString),
            get_link_from_net=sp.TMap(sp.TString, sp.TString),
            get_link_from_reachable_net=sp.TMap(sp.TString, types.Types.Tuple),
            helper=sp.TAddress
        ))

    def only_owner(self):
        with sp.if_(self.data.owners.contains(sp.sender)):
            sp.verify(self.data.owners[sp.sender] == True, "Unauthorized")
        with sp.else_():
            sp.failwith("Unauthorized")

    def only_bmc_periphery(self):
        sp.verify(sp.sender == self.data.bmc_periphery.open_some("BMCAddressNotSet"), "Unauthorized")

    @sp.entry_point
    def set_bmc_periphery(self, addr):
        """

        :param addr: address of bmc_periphery
        :return:
        """
        sp.set_type(addr, sp.TAddress)
        self.only_owner()
        sp.if self.data.bmc_periphery.is_some():
            sp.verify(addr != self.data.bmc_periphery.open_some("Address not set"), "AlreadyExistsBMCPeriphery")
        self.data.bmc_periphery = sp.some(addr)

    @sp.entry_point
    def set_bmc_btp_address(self, network):
        sp.set_type(network, sp.TString)

        sp.verify(self.data.owners[sp.sender] == True, "Unauthorized")
        # call set_btp_address on BMCPeriphery
        set_btp_address_entry_point = sp.contract(sp.TString,
                                                  self.data.bmc_periphery.open_some("Address not set"),
                                                  "set_bmc_btp_address").open_some()
        sp.transfer(network, sp.tez(0), set_btp_address_entry_point)

    @sp.entry_point
    def add_owner(self, owner):
        """
        :param owner: owner address to set
        :return:
        """
        sp.set_type(owner, sp.TAddress)

        self.only_owner()
        sp.verify(self.data.owners.contains(owner) == False, "Already Exists")
        self.data.owners[owner] = True
        self.data.number_of_owners += sp.nat(1)

    @sp.entry_point
    def remove_owner(self, owner):
        """

        :param owner: owner address to remove
        :return:
        """
        sp.set_type(owner, sp.TAddress)

        self.only_owner()
        sp.verify(self.data.number_of_owners > sp.nat(1), "LastOwner")
        sp.verify(self.data.owners[owner] == True, "NotExistsPermission")
        del self.data.owners[owner]
        self.data.number_of_owners = sp.as_nat(self.data.number_of_owners - sp.nat(1))

    @sp.onchain_view()
    def is_owner(self, owner):
        """

        :param owner: address to check
        :return:
        """
        sp.result(self.data.owners.get(owner))

    @sp.entry_point
    def add_service(self, svc, addr):
        """
        Add the smart contract for the service.
        :param svc: Name of the service
        :param addr: Service's contract address
        :return:
        """
        self.only_owner()
        sp.verify(addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"), "InvalidAddress")
        sp.verify(self.data.bsh_services.contains(svc) == False, "AlreadyExistsBSH")
        self.data.bsh_services[svc] = addr
        self.data.list_bsh_names.add(svc)

    @sp.entry_point
    def remove_service(self, svc):
        """
        Unregisters the smart contract for the service.
        :param svc: Name of the service
        :return:
        """
        self.only_owner()
        sp.verify(self.data.bsh_services.contains(svc), "NotExistsBSH")
        del self.data.bsh_services[svc]
        self.data.list_bsh_names.remove(svc)

    @sp.onchain_view()
    def get_services(self):
        """
        Get registered services.
        :return: An array of Service.
        """

        services = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.Service))
        i = sp.local("i", sp.nat(0))
        sp.for item in self.data.list_bsh_names.elements():
            services[i.value] = sp.record(svc=item, addr=self.data.bsh_services.get(item))
            i.value += 1
        sp.result(services)

    @sp.entry_point
    def add_link(self, link):
        """
        Initializes status information for the link.
        :param link:
        :return: BTP Address of connected BMC
        """
        sp.set_type(link, sp.TString)

        self.only_owner()
        net, addr= sp.match_pair(strings.split_btp_address(link, "prev_idx", "result", "my_list", "last", "penultimate"))

        with sp.if_(self.data.links.contains(link)):
            sp.verify(self.data.links.get(link).is_connected == False, "AlreadyExistsLink")
        self.data.links[link] = sp.record(
            relays=sp.set([]),
            reachable=sp.set([]),
            rx_seq=sp.nat(0),
            tx_seq=sp.nat(0),
            block_interval_src=self.BLOCK_INTERVAL_MSEC,
            block_interval_dst=sp.nat(0),
            max_aggregation=sp.nat(10),
            delay_limit=sp.nat(3),
            relay_idx=sp.nat(0),
            rotate_height=sp.nat(0),
            rx_height=sp.nat(0),
            rx_height_src=sp.nat(0),
            is_connected=True
        )

        self._propagate_internal("Link", link)
        links = self.data.list_link_names

        self.data.list_link_names.add(link)
        self.data.get_link_from_net[net] = link
        self._send_internal(link, "Init", links.elements())
        sp.trace("in add_link")

    @sp.entry_point
    def remove_link(self, link):
        """
        Removes the link and status information.
        :param link:  BTP Address of connected BMC
        :return:
        """
        sp.set_type(link, sp.TString)

        self.only_owner()
        with sp.if_(self.data.links.contains(link)):
            sp.verify(self.data.links.get(link).is_connected == True, "NotExistsLink")
        with sp.else_():
            sp.failwith("NotExistsLink")
        self._propagate_internal("Unlink", link)
        del self.data.links[link]
        net, addr= sp.match_pair(strings.split_btp_address(link, "prev_idx", "result", "my_list", "last", "penultimate"))
        del self.data.get_link_from_net[net]
        self.data.list_link_names.remove(link)
        sp.trace("in remove_link")

    @sp.onchain_view()
    def get_links(self):
        """
        Get registered links.
        :return: An array of links ( BTP Addresses of the BMCs ).
        """
        sp.result(self.data.list_link_names.elements())

    @sp.entry_point
    def set_link_rx_height(self, link, height):
        """

        :param link:
        :param height:
        :return:
        """

        sp.set_type(link, sp.TString)
        sp.set_type(height, sp.TNat)

        self.only_owner()
        with sp.if_(self.data.links.contains(link)):
            sp.verify(self.data.links.get(link).is_connected == True, "NotExistsLink")
        with sp.else_():
            sp.failwith("NotExistsKey")
        sp.verify(height > sp.nat(0), "InvalidRxHeight")
        self.data.links[link].rx_height = height

        sp.trace("in set_link_rx_height")

    @sp.entry_point
    def set_link(self, _link, block_interval, _max_aggregation, delay_limit):
        """

        :param _link:
        :param block_interval:
        :param _max_aggregation:
        :param delay_limit:
        :return:
        """
        sp.set_type(_link, sp.TString)
        sp.set_type(block_interval, sp.TNat)
        sp.set_type(_max_aggregation, sp.TNat)
        sp.set_type(delay_limit, sp.TNat)

        self.only_owner()

        with sp.if_(self.data.links.contains(_link)):
            sp.verify(self.data.links.get(_link).is_connected == True, "NotExistsLink")
        with sp.else_():
            sp.failwith("NotExistsLink")
        sp.verify((_max_aggregation >= sp.nat(1)) & (delay_limit >= sp.nat(1)), "InvalidParam")

        link = sp.local("link", self.data.links.get(_link), t=types.Types.Link).value

        # not implemented
        # scale = sp.local("scale", utils.get_scale(link.block_interval_src, link.block_interval_dst), t=sp.TNat)
        # reset_rotate_height = sp.local("reset_rotate_height", True, t=sp.TBool)
        # sp.if utils.get_rotate_term(link.max_aggregation, scale.value) == sp.nat(0):
        #     reset_rotate_height.value = True

        link.block_interval_src = block_interval
        link.max_aggregation = _max_aggregation
        link.delay_limit = delay_limit

        # not implemented
        # scale.value = utils.get_scale(link.block_interval_src, block_interval)
        # rotate_term = sp.local("rotate_term", utils.get_rotate_term(_max_aggregation, scale.value), t=sp.TNat)
        # rotate_term = sp.local("rotate_term", sp.nat(6))

        # sp.if reset_rotate_height.value & (rotate_term.value > sp.nat(0)):
        link.rotate_height = sp.level
        link.rx_height = sp.nat(0)
        # net, addr = sp.match_pair(strings.split_btp_address(_link))

        self.data.links[_link] = link
        sp.trace("in set_links")


    def _propagate_internal(self, service_type, link):
        sp.set_type(service_type, sp.TString)
        sp.set_type(link, sp.TString)

        _bytes = sp.bytes("0x")  # can be any bytes
        rlp_bytes = sp.view("encode_string", self.data.helper, link, t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [_bytes, rlp_bytes] , t=sp.TBytes).open_some()

        #encode payload
        final_rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [rlp_bytes_with_prefix], t=sp.TBytes).open_some()
        sp.for item in self.data.list_link_names.elements():
            sp.if self.data.links.get(item).is_connected:
                net, addr = sp.match_pair(strings.split_btp_address(item, "prev_idx1", "result1", "my_list1", "last1", "penultimate1"))

                # call send_message on BMCPeriphery
                send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
                send_message_entry_point = sp.contract(send_message_args_type,
                                                                self.data.bmc_periphery.open_some("Address not set"),
                                                                "send_message").open_some()
                send_message_args = sp.record(to=net, svc="bmc", sn=sp.int(0), msg=self.encode_bmc_service(
                    sp.record(serviceType=service_type, payload=final_rlp_bytes_with_prefix)))
                sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)

    def _send_internal(self, target, service_type, links):
        sp.set_type(target, sp.TString)
        sp.set_type(service_type, sp.TString)
        sp.set_type(links, sp.TList(sp.TString))

        rlp_bytes = sp.local("rlp_bytes", sp.bytes("0x"))
        with sp.if_(sp.len(links) == sp.nat(0)):
            rlp_bytes.value = self.LIST_SHORT_START
        with sp.else_():
            sp.for item in links:
                _bytes = sp.bytes("0x")  # can be any bytes
                _rlp_bytes = _bytes + sp.view("encode_string", self.data.helper, item, t=sp.TBytes).open_some()
                rlp_bytes.value = sp.view("encode_list", self.data.helper, [rlp_bytes.value, _rlp_bytes], t=sp.TBytes).open_some()
        #encode payload
        # final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes.value, t=sp.TBytes).open_some()
        net, addr = sp.match_pair(
            strings.split_btp_address(target, "prev_idx2", "result2", "my_list2", "last2", "penultimate2"))

        # call send_message on BMCPeriphery
        send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
        send_message_entry_point = sp.contract(send_message_args_type,
                                               self.data.bmc_periphery.open_some("Address not set"),
                                               "send_message").open_some()
        send_message_args = sp.record(to=net, svc="bmc", sn=sp.int(0), msg=self.encode_bmc_service(
            sp.record(serviceType=service_type, payload=rlp_bytes.value)))
        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)


    @sp.entry_point
    def add_route(self, dst, link):
        """
        Add route to the BMC.
        :param dst: BTP Address of the destination BMC
        :param link: BTP Address of the next BMC for the destination
        :return:
        """
        sp.set_type(dst, sp.TString)
        sp.set_type(link, sp.TString)

        self.only_owner()
        sp.verify(self.data.routes.contains(dst) == False, "AlreadyExistRoute")
        net, addr= sp.match_pair(strings.split_btp_address(dst, "prev_idx", "result", "my_list", "last", "penultimate"))
        # TODO: need to verify link is only split never used
        # strings.split_btp_address(link)

        self.data.routes[dst] = link
        self.data.list_route_keys.add(dst)
        self.data.get_route_dst_from_net[net] = dst

    @sp.entry_point
    def remove_route(self, dst):
        """
        Remove route to the BMC.
        :param dst:  BTP Address of the destination BMC
        :return:
        """
        sp.set_type(dst, sp.TString)

        self.only_owner()
        sp.verify(self.data.routes.contains(dst) == True, "NotExistRoute")
        del self.data.routes[dst]
        net, addr= sp.match_pair(strings.split_btp_address(dst, "prev_idx", "result", "my_list", "last", "penultimate"))
        del self.data.get_route_dst_from_net[net]
        self.data.list_route_keys.remove(dst)

    @sp.onchain_view()
    def get_routes(self):
        """
        Get routing information.
        :return: An array of Route.
        """

        _routes = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.Route))
        i = sp.local("i", sp.nat(0))
        sp.for item in self.data.list_route_keys.elements():
            _routes[i.value] = sp.record(dst=item, next=self.data.routes.get(item))
            i.value += 1
        sp.result(_routes)

    @sp.entry_point
    def add_relay(self, link, addr):
        """
        Registers relay for the network.
        :param link: BTP Address of connected BMC
        :param addr: the address of Relay
        :return:
        """
        sp.set_type(link, sp.TString)
        sp.set_type(addr, sp.TSet(sp.TAddress))

        self.only_owner()
        sp.verify(self.data.links.contains(link), "NotExistsLink")
        sp.verify(self.data.links.get(link).is_connected == True, "NotExistsLink")
        self.data.links[link].relays = addr
        sp.for item in addr.elements():
            self.data.relay_stats[item] = sp.record(addr=item, block_count=sp.nat(0), msg_count=sp.nat(0))

    @sp.entry_point
    def remove_relay(self, link, addr):
        """
        Unregisters Relay for the network.
        :param link: BTP Address of connected BMC
        :param addr: the address of Relay
        :return:
        """
        sp.set_type(link, sp.TString)
        sp.set_type(addr, sp.TAddress)

        self.only_owner()
        sp.verify(self.data.links.contains(link), "NotExistsLink")
        sp.verify((self.data.links.get(link).is_connected == True) & (sp.len(self.data.links.get(link).relays.elements()) != sp.nat(0)),
                  "Unauthorized")

        sp.for item in self.data.links.get(link).relays.elements():
            sp.if item != addr:
                self.data.addrs.add(item)

        self.data.links[link].relays = self.data.addrs

        # delete all items from addrs set
        sp.for ele in self.data.addrs.elements():
            self.data.addrs.remove(ele)


    @sp.onchain_view()
    def get_relays(self, link):
        """
        Get registered relays.
        :param link: BTP Address of the connected BMC.
        :return: A list of relays
        """
        sp.set_type(link, sp.TString)

        sp.result(self.data.links.get(link).relays.elements())

    @sp.onchain_view()
    def get_bsh_service_by_name(self, service_name):
        sp.set_type(service_name, sp.TString)
        sp.result(self.data.bsh_services.get(service_name, default_value=sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")))

    @sp.onchain_view()
    def get_link(self, to):
        sp.set_type(to, sp.TString)
        sp.result(self.data.links.get(to))

    @sp.onchain_view()
    def get_link_rx_seq(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links.get(prev).rx_seq)

    @sp.onchain_view()
    def get_link_tx_seq(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links.get(prev).tx_seq)

    @sp.onchain_view()
    def get_link_rx_height(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links.get(prev).rx_height)

    @sp.onchain_view()
    def get_link_relays(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links.get(prev).relays.elements())

    @sp.onchain_view()
    def get_relay_status_by_link(self, prev):
        sp.set_type(prev, sp.TString)
        _relays = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.RelayStats))

        i = sp.local("i", sp.nat(0))
        sp.for item in self.data.links.get(prev).relays.elements():
            _relays[i.value] = self.data.relay_stats.get(item)
            i.value += 1
        sp.result(_relays)

    @sp.entry_point
    def update_link_rx_seq(self, prev, val):
        sp.set_type(prev, sp.TString)
        sp.set_type(val, sp.TNat)

        self.only_bmc_periphery()
        self.data.links[prev].rx_seq += val

    @sp.entry_point
    def update_link_tx_seq(self, prev):
        sp.set_type(prev, sp.TString)

        self.only_bmc_periphery()
        self.data.links[prev].tx_seq += sp.nat(1)

    @sp.entry_point
    def update_link_rx_height(self, prev, val):
        sp.set_type(prev, sp.TString)
        sp.set_type(val, sp.TNat)

        self.only_bmc_periphery()
        self.data.links[prev].rx_height += val

    @sp.entry_point
    def update_link_reachable(self, prev, to):
        sp.set_type(prev, sp.TString)
        sp.set_type(to, sp.TList(sp.TString))

        self.only_bmc_periphery()
        sp.for item in to:
            self.data.links[prev].reachable.add(item)
            net, addr = sp.match_pair(
                strings.split_btp_address(item, "prev_idx", "result", "my_list", "last", "penultimate"))
            self.data.get_link_from_reachable_net[net] = sp.record(prev=prev, to=item)

    @sp.entry_point
    def delete_link_reachable(self, prev, index):
        sp.set_type(prev, sp.TString)
        sp.set_type(index, sp.TNat)

        self.only_bmc_periphery()
        i = sp.local("i", sp.nat(0))
        sp.for item in self.data.links.get(prev).reachable.elements():
            sp.if i.value == index:
                net, addr = sp.match_pair(
                    strings.split_btp_address(item, "prev_idx", "result", "my_list", "last", "penultimate"))

                del self.data.get_link_from_reachable_net[net]
                self.data.links[prev].reachable.remove(item)
                # this is not needed when removing from set
                # self.data.links[prev].reachable[index] = self.data.links[prev].reachable[
                #     sp.as_nat(sp.len(self.data.links[prev].reachable) - 1)
                # ]
            i.value += 1

    @sp.entry_point
    def update_relay_stats(self, relay, block_count_val, msg_count_val):
        sp.set_type(relay, sp.TAddress)
        sp.set_type(block_count_val, sp.TNat)
        sp.set_type(msg_count_val, sp.TNat)

        self.only_bmc_periphery()
        self.data.relay_stats[relay].block_count += block_count_val
        self.data.relay_stats[relay].msg_count += msg_count_val

    @sp.onchain_view()
    def resolve_route(self, dst_net):
        sp.set_type(dst_net, sp.TString)

        dst = sp.local("dst", self.data.get_route_dst_from_net.get(dst_net), t=sp.TString)

        with sp.if_(sp.len(sp.pack(dst.value))!= sp.nat(0)):
            sp.result(sp.pair(self.data.routes.get(dst.value), dst.value))
        with sp.else_():
            dst_link = sp.local("dst_link", self.data.get_link_from_net.get(dst_net), t=sp.TString)
            with sp.if_(sp.len(sp.pack(dst_link.value)) != sp.nat(0)):
                sp.result(sp.pair(dst_link.value, dst_link.value))
            with sp.else_():
                res = sp.local("res", self.data.get_link_from_reachable_net.get(dst_net), t=types.Types.Tuple)
                sp.verify(sp.len(sp.pack(res.value.to)) > sp.nat(0), "Unreachable: " + dst_net + " is unreachable")

                sp.result(sp.pair(res.value.prev, res.value.to))


@sp.add_test(name="BMCM")
def test():
    alice = sp.test_account("Alice")
    owner = sp.test_account("Owner")
    helper = sp.test_account("Helper")
    # bmc_periphery = sp.test_account("BMC Periphery")
    # bmc= sp.test_account("BMC")


    scenario = sp.test_scenario()
    # deploy helper
    # helper_con = helper_file.Helper()
    # scenario += helper_con
    #
    # parse_addr_con = parse_address_file.ParseAddress()
    # scenario += parse_addr_con
    #
    bmc_man = BMCManagement(owner.address, helper.address)
    scenario += bmc_man
    # # deploy bmc periphery
    # bmc_per = bmc_periphery_file.BMCPreiphery(
    #     bmc_man.address, helper_con.address, parse_addr_con.address
    # )
    #
    # scenario += bmc_per
    #
    # bmc_man.set_bmc_periphery(bmc_per.address).run(sender=owner)
    # bmc_man.set_bmc_btp_address("tezos.77").run(sender=owner)
    # # bmc_man.add_owner(alice.address).run(sender=owner)
    #
    # # bmc_man.remove_owner(alice.address).run(sender=alice)
    #
    # bmc_man.add_route(sp.record(dst = "btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hDEST",
    # link = "btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW")).run(sender= owner)
    # bmc_man.add_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW").run(sender=owner)
    # # decode_string = sp.build_lambda(Utils.RLP.Decoder.without_length_prefix)
    #
    # # bmc_man.remove_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW").run(sender=alice)
    # # bmc_man.add_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW").run(sender=alice)
    #
    # # bmc_man.set_link_rx_height(sp.record(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", height=sp.nat(2))).run(sender=alice)
    # # bmc_man.set_link(sp.record(_link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", block_interval=sp.nat(2),
    # #                             _max_aggregation=sp.nat(3), delay_limit=sp.nat(2))).run(sender=alice)
    #
    #
sp.add_compilation_target("bmc_management", BMCManagement(owner_address=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"),
                                                          helper_contract=sp.address("KT1Q5erZm7Pp8UJywK1nkiP8QPCRmyUotUMq")))
