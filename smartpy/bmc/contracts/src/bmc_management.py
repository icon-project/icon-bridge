import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")


class BMCManagement(sp.Contract):
    BLOCK_INTERVAL_MSEC = sp.nat(1000)

    def __init__(self):
        self.update_initial_storage(
            owners=sp.map(l={sp.address("tz1000"): True}, tkey=sp.TAddress, tvalue=sp.TBool),
            number_of_owners=sp.nat(1),
            bsh_service=sp.map(tkey=sp.TString, tvalue=sp.TAddress),
            relay_stats=sp.map(tkey=sp.TAddress, tvalue=types.Types.RelayStats),
            routes=sp.map(tkey=sp.TString, tvalue=sp.TString),
            links=sp.map(tkey=sp.TString, tvalue=types.Types.Link),
            list_bsh_names=sp.list(),
            list_route_keys=sp.list(),
            list_link_names=sp.list(),
            bmc_periphery=sp.none,
            serial_no=sp.nat(0),
            addrs=sp.map(tkey=sp.TNat, tvalue=sp.TAddress),
            get_route_dst_from_net=sp.map(tkey=sp.TString, tvalue=sp.TString),
            get_link_from_net=sp.map(tkey=sp.TString, tvalue=sp.TString),
            get_link_from_reachable_net=sp.map(tkey=sp.TString, tvalue=types.Types.Tuple)
        )

        # self.init_type(sp.TRecord(
        #     list_bsh_names=sp.TList(sp.TAddress),
        #     list_route_keys=sp.TList(sp.TAddress),
        #     list_link_names=sp.TList(sp.TAddress),
        #     addrs=sp.TList(sp.TAddress)
        # ))

    def only_owner(self):
        sp.verify(self.data.owners[sp.sender] == True, "Unauthorized")

    def only_bmc_periphery(self):
        sp.verify(sp.sender == self.data.bmc_periphery.open_some("BMCAddressNotSet"), "Unauthorized")

    @sp.entry_point
    def set_bmc_periphery(self, addr):
        """

        :param addr: address of bmc_periphery
        :return:
        """
        sp.set_type(addr, sp.TAddress)
        # self.only_owner()
        sp.verify(addr != sp.address("tz100000000"), "InvalidAddress")
        sp.trace(self.data.bmc_periphery.is_some())
        sp.if self.data.bmc_periphery.is_some():
            sp.verify(addr != self.data.bmc_periphery.open_some("Address not set"), "AlreadyExistsBMCPeriphery")
        self.data.bmc_periphery = sp.some(addr)

    @sp.entry_point
    def add_owner(self, owner):
        """
        :param owner: owner address to set
        :return:
        """
        sp.set_type(owner, sp.TAddress)

        # self.only_owner()
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
        sp.result(self.data.owners[owner])

    @sp.entry_point
    def add_service(self, svc, addr):
        """
        Add the smart contract for the service.
        :param svc: Name of the service
        :param addr: Service's contract address
        :return:
        """
        self.only_owner()
        sp.verify(addr != sp.address("tz100000000"), "InvalidAddress")
        sp.verify(self.data.bsh_service[svc] == sp.address("tz100000000"), "AlreadyExistsBSH")
        self.data.bsh_service[svc] = addr
        self.data.list_bsh_names.push(svc)

    @sp.entry_point
    def remove_service(self, svc):
        """
        Unregisters the smart contract for the service.
        :param svc: Name of the service
        :return:
        """
        self.only_owner()
        sp.verify(self.data.bsh_service[svc] == sp.address("tz100000000"), "NotExistsBSH")
        del self.data.bsh_service[svc]

    #     need to delete item from list_bsh_names (not possible)

    # @sp.onchain_view()
    # def get_services(self):
    #     """
    #     Get registered services.
    #     :return: An array of Service.
    #     """
    #
    #     services = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.Service))
    #
    #     sp.for i in sp.range(sp.nat(0), sp.len(self.data.bsh_service)):
    #         services[i] = sp.record(
    #
    #         )

    @sp.entry_point
    def add_link(self, link):
        """
        Initializes status information for the link.
        :param link:
        :return: BTP Address of connected BMC
        """
        sp.set_type(link, sp.TString)

        self.only_owner()
        net, addr= sp.match_pair(strings.split_btp_address(link))

        with sp.if_(self.data.links.contains(link)):
            sp.verify(self.data.links[link].is_connected == False, "AlreadyExistsLink")
        #TODO:review how to add key in relays map
        self.data.links[link] = sp.record(
            relays=sp.map({0: sp.address("tz100000000")}),
            reachable=sp.map({0: "0"}),
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

        # self._propagate_internal("Link", link)
        links = sp.local("links", self.data.list_link_names, t=sp.TList(sp.TString))
        # TODO:push link to listLinkNames

        self.data.get_link_from_net[net] = link

        # self._send_internal(link, "Init", links.value)
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
        # TODO : review this if else
        with sp.if_(self.data.links.contains(link)):
            sp.verify(self.data.links[link].is_connected == True, "NotExistsLink")
        with sp.else_():
            sp.failwith("NotExistsLink")
        del self.data.links[link]
        net, addr= sp.match_pair(strings.split_btp_address(link))
        del self.data.get_link_from_net[link]
        # self._propagate_internal("Unlink", link)
        # TODO:remove link to listLinkNames
        sp.trace("in remove_link")

    @sp.onchain_view()
    def get_links(self):
        """
        Get registered links.
        :return: An array of links ( BTP Addresses of the BMCs ).
        """
        sp.result(self.data.list_link_names)


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
            sp.verify(self.data.links[link].is_connected == True, "NotExistsLink")
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

        sp.verify(self.data.links[_link].is_connected == True, "NotExistsLink")
        sp.verify((_max_aggregation >= sp.nat(1)) & (delay_limit >= sp.nat(1)), "InvalidParam")

        link = sp.local("link", self.data.links[_link], t=types.Types.Link).value
        # scale = sp.local("scale", utils.get_scale(link.block_interval_src, link.block_interval_dst), t=sp.TNat)

        reset_rotate_height = sp.local("reset_rotate_height", False, t=sp.TBool)

        # sp.if utils.get_rotate_term(link.max_aggregation, scale.value) == sp.nat(0):
        #     reset_rotate_height.value = True

        link.block_interval_src = block_interval
        link.max_aggregation = _max_aggregation
        link.delay_limit = delay_limit

        # TODO: uncomment after writing utils lib
        # scale.value = utils.get_scale(link.block_interval_src, block_interval)
        # rotate_term = sp.local("rotate_term", utils.get_rotate_term(_max_aggregation, scale.value), t=sp.TNat)
        rotate_term = sp.local("rotate_term", sp.nat(6))

        sp.if reset_rotate_height.value & (rotate_term.value > sp.nat(0)):
            link.rotate_height = sp.level + rotate_term.value
            link.rx_height = sp.nat(0)
            net, addr = sp.match_pair(strings.split_btp_address(_link))

        self.data.links[_link] = link
        sp.trace("in set_links")

    def _rotate(self, _link, rotate_term, rotate_count, base_height):
        sp.set_type(_link, sp.TString)
        sp.set_type(rotate_term, sp.TNat)
        sp.set_type(rotate_count, sp.TNat)
        sp.set_type(base_height, sp.TNat)

        link = sp.local("link", self.data.links[_link], t=types.Types.Link).value

        sp.if (rotate_term >sp.nat(0)) & (rotate_count > sp.nat(0)):
            link.rotate_height = base_height + rotate_term
            link.relay_idx = link.relay_idx + rotate_count
            sp.if link.relay_idx >= sp.len(link.relays):
                link.relay_idx = link.relay_idx % sp.len(link.relays)
            links[_link] = link
        return link.relays[link.relay_idx]

    def _propagate_internal(self, service_type, link):
        sp.set_type(service_type, sp.TString)
        sp.set_type(link, sp.TString)

        #TODO implement abi encodePacked

        # TODO: encode actual payload
        rlp_bytes =sp.bytes("rlp_bytes_here")

        sp.for i in sp.range(sp.nat(0), sp.len(self.data.list_link_names)):
            sp.for item in self.data.list_link_names:
                sp.if self.data.links[item].is_connected:
                    net, addr = sp.match_pair(strings.split_btp_address(item))

                    # call send_message on BMCPeriphery
                    send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TNat, msg=sp.TBytes)
                    send_message_entry_point = sp.contract(send_message_args_type,
                                                                    self.data.bmc_periphery,
                                                                    "send_message").open_some()
                    send_message_args = sp.record(to=net, svc="bmc", sn=sp.nat(0), msg=encodeBMCService(
                        sp.record(serviceType=service_type, payload=rlp_bytes)))
                    sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)



    def _send_internal(self, target, service_type, links):
        sp.set_type(target, sp.TString)
        sp.set_type(service_type, sp.TString)
        sp.set_type(links, sp.TMap(TNat, TString))

        with sp.if_(sp.len(links) == sp.nat(0)):
            # TODO: abi encode
            rlp_bytes = sp.bytes("rpl byte")
        with sp.else_():
            sp.for i in sp.range(0, sp.len(links)):
                #TODO
                # TODO: abi encode
                rlp_bytes = sp.bytes("rpl byte")

        #TODO: encode payload
        rlp_bytes = sp.bytes("abi encode")

        net, addr = sp.match_pair(strings.split_btp_address(target))

        # call send_message on BMCPeriphery
        send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TNat, msg=sp.TBytes)
        send_message_entry_point = sp.contract(send_message_args_type,
                                               self.data.bmc_periphery,
                                               "send_message").open_some()
        send_message_args = sp.record(to=net, svc="bmc", sn=sp.nat(0), msg=encodeBMCService(
            sp.record(serviceType=service_type, payload=rlp_bytes)))
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
        sp.verify(sp.len(sp.pack(self.data.routes[dst])) == sp.nat(0), "AlreadyExistRoute")
        net, addr= sp.match_pair(strings.split_btp_address(dst))
        # TODO: need to verify link is only split never used
        # strings.split_btp_address(link)

        self.data.routes[dst] = link
        self.data.list_route_keys.push(dst)
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
        sp.verify(sp.len(sp.pack(self.data.routes[dst])) != sp.nat(0), "NotExistRoute")
        del self.data.routes[dst]
        net, addr= sp.match_pair(strings.split_btp_address(dst))
        del self.data.get_route_dst_from_net[net]
        #TODO: remove dst from list_route_keys

    @sp.onchain_view()
    def get_routes(self):
        """
        Get routing information.
        :return: An array of Route.
        """

        routes = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.Route))
        #TODO: whether to make listRouteKeys list or map

    @sp.entry_point
    def add_relay(self, link, addr):
        """
        Registers relay for the network.
        :param link: BTP Address of connected BMC
        :param addr: the address of Relay
        :return:
        """
        sp.set_type(link, sp.TString)
        sp.set_type(addr, sp.TMap(sp.TNat, sp.TAddress))

        self.only_owner()
        sp.verify(self.data.links[link].is_connected == True, "NotExistsLink")
        self.data.links[link].relays = addr
        sp.for i in sp.range(sp.nat(0), sp.len(addr)):
            self.data.relay_stats[addr[i]] = sp.record(addr=addr[i], block_count=sp.nat(0), msg_count=sp.nat(0))

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
        sp.verify((self.data.links[link].is_connected == True) & (sp.len(self.data.links[link].relays) != sp.nat(0)),
                  "Unauthorized")

        sp.for i in sp.range(sp.nat(0), sp.len(self.data.links[link].relays)):
            sp.if self.data.links[link].relays[i] != addr:
                self.data.addrs[i] = self.data.links[link].relays[i]

        self.data.links[link].relays = self.data.addrs
        # TODO: delete addrs map
        # del self.data.addrs

    @sp.onchain_view()
    def get_relays(self, link):
        """
        Get registered relays.
        :param link: BTP Address of the connected BMC.
        :return: A map of relays
        """
        sp.set_type(link, sp.TString)

        sp.result(self.data.links[link].relays)

    @sp.onchain_view()
    def get_bsh_service_by_name(self, service_name):
        sp.set_type(service_name, sp.TString)
        sp.result(self.data.bsh_service[service_name])

    @sp.onchain_view()
    def get_link(self, to):
        sp.set_type(to, sp.TString)
        sp.result(self.data.links[to])

    @sp.onchain_view()
    def get_link_rx_seq(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links[prev].rx_seq)

    @sp.onchain_view()
    def get_link_tx_seq(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links[prev].tx_seq)

    @sp.onchain_view()
    def get_link_rx_height(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links[prev].rx_height)

    @sp.onchain_view()
    def get_link_relays(self, prev):
        sp.set_type(prev, sp.TString)
        sp.result(self.data.links[prev].relays)

    @sp.onchain_view()
    def get_relay_status_by_link(self, prev):
        sp.set_type(prev, sp.TString)
        relays = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.RelayStats))
        sp.for i in sp.range(sp.nat(0), sp.len(self.data.links[prev].relays)):
            relays[i] = self.data.relay_stats[self.data.links[prev].relays[i]]
        sp.result(relays)

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
        sp.set_type(to, sp.TMap(sp.TNat, sp.TString))

        self.only_bmc_periphery()
        sp.for i in sp.range(sp.nat(0), sp.len(to)):
            self.data.links[prev].reachable[i] = to[i]
            net, addr = sp.match_pair(strings.split_btp_address(to[i]))
            self.data.get_link_from_reachable_net[net] = sp.record(prev=prev, to=to[i])

    @sp.entry_point
    def delete_link_reachable(self, prev, index):
        sp.set_type(prev, sp.TString)
        sp.set_type(index, sp.TNat)

        self.only_bmc_periphery()
        net, addr = sp.match_pair(strings.split_btp_address(self.data.links[prev].reachable[index]))

        del self.data.get_link_from_reachable_net[net]
        del self.data.links[prev].reachable[index]
        self.data.links[prev].reachable[index] = self.data.links[prev].reachable[
            sp.as_nat(sp.len(self.data.links[prev].reachable) - 1)
        ]
        # TODO:pop from reachable

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

        self.only_bmc_periphery()
        dst = sp.local("dst", self.data.get_route_dst_from_net[dst_net], t=sp.TString)

        # TODO: calculate length of byte
        # sp.if sp.len(sp.pack(dst.value))!= sp.nat(0):
        #     sp.result(sp.pair(self.data.routes[dst.value], dst.value))

        dst_link = sp.local("dst_link", self.data.get_link_from_net[dst_net], t=sp.TString)
        # TODO: calculate length of byte
        # sp.if sp.len(sp.pack(dst_link.value)) != sp.nat(0):
        #     sp.result(sp.pair(dst_link.value, dst_link.value))

        res = sp.local("res", self.data.get_link_from_reachable_net[dst_net], t=types.Types.Tuple)
        sp.verify(sp.len(sp.pack(res.value.to)) > sp.nat(0), "Unreachable: " + dst_net + " is unreachable")

        sp.result(sp.pair(res.value.prev, res.value.to))


@sp.add_test(name="BMCM")
def test():
    alice = sp.test_account("Alice")
    bmc_periphery = sp.test_account("BMC Periphery")
    # bmc= sp.test_account("BMC")

    scenario = sp.test_scenario()
    bmc_man = BMCManagement()
    scenario += bmc_man

    bmc_man.set_bmc_periphery(bmc_periphery.address).run(sender=alice)
    bmc_man.add_owner(alice.address).run(sender=alice)

    # bmc_man.remove_owner(alice.address).run(sender=alice)

    bmc_man.add_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW").run(sender=alice)
    # bmc_man.remove_link("btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW").run(sender=alice)

    bmc_man.set_link_rx_height(sp.record(link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", height=sp.nat(2))).run(sender=alice)
    bmc_man.set_link(sp.record(_link="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", block_interval=sp.nat(2),
                                _max_aggregation=sp.nat(3), delay_limit=sp.nat(2))).run(sender=alice)


sp.add_compilation_target("bmc_management", BMCManagement())
