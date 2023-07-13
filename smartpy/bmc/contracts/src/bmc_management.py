import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")
rlp = sp.io.import_script_from_url("file:./contracts/src/RLP_struct.py")


class BMCManagement(sp.Contract, rlp.DecodeEncodeLibrary):
    BLOCK_INTERVAL_MSEC = sp.nat(30000)
    LIST_SHORT_START = sp.bytes("0xc0")
    ZERO_ADDRESS = sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")

    def __init__(self, owner_address, helper_contract):
        self.init(
            owners=sp.map(l={owner_address:Tue}),
            number_of_owners=sp.nat(1),
            bsh_services=sp.map(),
            relay_stats=sp.map(),
            routes=sp.map(),
            links=sp.map(),
            list_route_keys=sp.set(),
            list_link_names=sp.set(),
            bmc_periphery=self.ZERO_ADDRESS,
            serial_no=sp.nat(0),
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
            list_route_keys=sp.TSet(sp.TString),
            list_link_names=sp.TSet(sp.TString),
            bmc_periphery=sp.TAddress,
            serial_no=sp.TNat,
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
        sp.verify(sp.sender == self.data.bmc_periphery, "Unauthorized")

    @sp.entry_point
    def set_helper_address(self, address):
        sp.set_type(address, sp.TAddress)
        self.only_owner()
        self.data.helper = address

    @sp.entry_point
    def set_bmc_periphery(self, addr):
        """
        :param addr: address of bmc_periphery
        :return:
        """
        sp.set_type(addr, sp.TAddress)

        self.only_owner()
        sp.verify(addr != self.ZERO_ADDRESS, "Invalid Address")
        sp.verify(addr != self.data.bmc_periphery, "AlreadyExistsBMCPeriphery")
        self.data.bmc_periphery = addr

    @sp.entry_point
    def set_bmc_btp_address(self, network):
        sp.set_type(network, sp.TString)

        self.only_owner()
        # call set_btp_address on BMCPeriphery
        set_btp_address_entry_point = sp.contract(sp.TString,
                                                  self.data.bmc_periphery,
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

    @sp.entry_point
    def remove_owner(self, owner):
        """

        :param owner: owner address to remove
        :return:
        """
        sp.set_type(owner, sp.TAddress)

        self.only_owner()
        sp.verify(sp.len(self.data.owners) > sp.nat(1), "LastOwner")
        sp.verify(self.data.owners[owner] == True, "NotExistsPermission")
        del self.data.owners[owner]

    @sp.onchain_view()
    def is_owner(self, owner):
        """
        :param owner: address to check
        :return:
        """
        sp.set_type(owner, sp.TAddress)
        sp.result(self.data.owners.get(owner))

    @sp.entry_point
    def add_service(self, svc, addr):
        """
        Add the smart contract for the service.
        :param svc: Name of the service
        :param addr: Service's contract address
        :return:
        """
        sp.set_type(svc, sp.TString)
        sp.set_type(addr, sp.TAddress)

        self.only_owner()
        sp.verify(addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"), "InvalidAddress")
        sp.verify(self.data.bsh_services.contains(svc) == False, "AlreadyExistsBSH")
        self.data.bsh_services[svc] = addr

    @sp.entry_point
    def remove_service(self, svc):
        """
        Unregisters the smart contract for the service.
        :param svc: Name of the service
        :return:
        """
        sp.set_type(svc, sp.TString)

        self.only_owner()
        sp.verify(self.data.bsh_services.contains(svc), "NotExistsBSH")
        del self.data.bsh_services[svc]

    @sp.onchain_view()
    def get_services(self):
        """
        Get registered services.
        :return: An array of Service.
        """
        sp.result(self.data.bsh_services)

    @sp.entry_point(lazify=False)
    def update_add_link(self, ep):
        self.only_owner()
        sp.set_entry_point("add_link", ep)

    @sp.entry_point(lazify=True)
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
        links = sp.compute(self.data.list_link_names.elements())

        self.data.list_link_names.add(link)
        self.data.get_link_from_net[net] = link
        self._send_internal(link, "Init", links)

    @sp.entry_point(lazify=False)
    def update_remove_link(self, ep):
        self.only_owner()
        sp.set_entry_point("remove_link", ep)

    @sp.entry_point(lazify=True)
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

        del self.data.links[link]
        net, addr= sp.match_pair(strings.split_btp_address(link, "prev_idx", "result", "my_list", "last", "penultimate"))
        del self.data.get_link_from_net[net]
        self._propagate_internal("Unlink", link)
        self.data.list_link_names.remove(link)

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

        link.block_interval_dst = block_interval
        link.max_aggregation = _max_aggregation
        link.delay_limit = delay_limit

        link.rotate_height = sp.level
        link.rx_height = sp.nat(0)

        self.data.links[_link] = link


    def _propagate_internal(self, service_type, link):
        sp.set_type(service_type, sp.TString)
        sp.set_type(link, sp.TString)

        _bytes = sp.bytes("0x")  # can be any bytes
        rlp_bytes = sp.view("encode_string", self.data.helper, link, t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [_bytes, rlp_bytes] , t=sp.TBytes).open_some()

        #encode payload
        final_rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [rlp_bytes_with_prefix],
                                              t=sp.TBytes).open_some()
        sp.for item in self.data.list_link_names.elements():
            with sp.if_(self.data.links.contains(item)):
                with sp.if_(self.data.links.get(item).is_connected):
                    net, addr = sp.match_pair(strings.split_btp_address(item, "prev_idx1", "result1",
                                                                        "my_list1", "last1", "penultimate1"))

                    # call send_message on BMCPeriphery
                    send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
                    send_message_entry_point = sp.contract(send_message_args_type,
                                                                    self.data.bmc_periphery,
                                                                    "send_message").open_some()
                    send_message_args = sp.record(to=net, svc="bmc", sn=sp.int(0), msg=self.encode_bmc_service(
                                                sp.record(serviceType=service_type,payload=final_rlp_bytes_with_prefix)))
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
                rlp_bytes.value = sp.view("encode_list", self.data.helper, [rlp_bytes.value, _rlp_bytes],
                                          t=sp.TBytes).open_some()
        #encode payload
        net, addr = sp.match_pair(
            strings.split_btp_address(target, "prev_idx2", "result2", "my_list2", "last2", "penultimate2"))

        # call send_message on BMCPeriphery
        send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
        send_message_entry_point = sp.contract(send_message_args_type,
                                               self.data.bmc_periphery,
                                               "send_message").open_some()
        send_message_args = sp.record(to=net, svc="bmc", sn=sp.int(0),
                                        msg=self.encode_bmc_service(sp.record(serviceType=service_type,
                                                                              payload=rlp_bytes.value)))
        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)

    @sp.entry_point(lazify=False)
    def update_add_route(self, ep):
        self.only_owner()
        sp.set_entry_point("add_route", ep)

    @sp.entry_point(lazify=True)
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
        strings.split_btp_address(link, "prev_idx1", "result1", "my_list1", "last1", "penultimate1")

        self.data.routes[dst] = link
        self.data.list_route_keys.add(dst)
        self.data.get_route_dst_from_net[net] = dst

    @sp.entry_point(lazify=False)
    def update_remove_route(self, ep):
        self.only_owner()
        sp.set_entry_point("remove_route", ep)

    @sp.entry_point(lazify=True)
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

    @sp.entry_point(lazify=False)
    def update_add_relay(self, ep):
        self.only_owner()
        sp.set_entry_point("add_relay", ep)

    @sp.entry_point(lazify=True)
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

    @sp.entry_point(lazify=False)
    def update_remove_relay(self, ep):
        self.only_owner()
        sp.set_entry_point("remove_relay", ep)

    @sp.entry_point(lazify=True)
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
        sp.verify((self.data.links.get(link).is_connected == True) &
                  (sp.len(self.data.links.get(link).relays.elements()) != sp.nat(0)), "Unauthorized")
        addr_set = sp.local("addr_set", sp.set(), t=sp.TSet(sp.TAddress))
        sp.for item in self.data.links.get(link).relays.elements():
            with sp.if_(item != addr):
                addr_set.value.add(item)

        self.data.links[link].relays = addr_set.value

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
        sp.result(self.data.bsh_services.get(service_name,
                                             default_value=sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")))

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

    @sp.entry_point(lazify=False)
    def update_update_link_tx_seq(self, ep):
        self.only_owner()
        sp.set_entry_point("update_link_tx_seq", ep)

    @sp.entry_point(lazify=True)
    def update_link_tx_seq(self, prev, serialized_msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(serialized_msg, sp.TBytes)

        self.only_bmc_periphery()
        self.data.links[prev].tx_seq += sp.nat(1)

        sp.emit(sp.record(next=prev, seq=self.data.links.get(prev).tx_seq, msg=serialized_msg), tag="Message")

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
            with sp.if_(i.value == index):
                net, addr = sp.match_pair(
                    strings.split_btp_address(item, "prev_idx", "result", "my_list", "last", "penultimate"))

                del self.data.get_link_from_reachable_net[net]
                self.data.links[prev].reachable.remove(item)
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

        dst = sp.local("dst", self.data.get_route_dst_from_net.get(dst_net, default_value=sp.string("")), t=sp.TString)
        with sp.if_(sp.len(dst.value)!= sp.nat(0)):
            sp.result(sp.pair(self.data.routes.get(dst.value), dst.value))
        with sp.else_():
            dst_link = sp.local("dst_link", self.data.get_link_from_net.get(dst_net,
                                                                            default_value=sp.string("")), t=sp.TString)
            with sp.if_(sp.len(dst_link.value) != sp.nat(0)):
                sp.result(sp.pair(dst_link.value, dst_link.value))
            with sp.else_():
                res = sp.local("res", self.data.get_link_from_reachable_net.get(dst_net, default_value=
                sp.record(prev="", to="")), t=types.Types.Tuple)
                with sp.if_(sp.len(res.value.to) > sp.nat(0)):
                    sp.result(sp.pair(res.value.prev, res.value.to))
                with sp.else_():
                    sp.result(sp.pair("Unreachable", "Unreachable: " + dst_net + " is unreachable"))


sp.add_compilation_target("bmc_management", BMCManagement(owner_address=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"),
                                                          helper_contract=sp.address("KT1HwFJmndBWRn3CLbvhUjdupfEomdykL5a6")
                                                          ))
