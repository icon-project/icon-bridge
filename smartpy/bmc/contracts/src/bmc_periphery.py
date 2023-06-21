import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")


class BMCPreiphery(sp.Contract):
    BMC_ERR = sp.nat(10)
    BSH_ERR = sp.nat(40)
    UNKNOWN_ERR = sp.nat(0)

    BMCRevertUnauthorized = sp.string("Unauthorized")
    BMCRevertParseFailure = sp.string("ParseFailure")
    BMCRevertNotExistsBSH = sp.string("NotExistsBSH")
    BMCRevertNotExistsLink = sp.string("NotExistsLink")
    BMCRevertInvalidSn = sp.string("InvalidSn")
    BMCRevertInvalidSeqNumber = sp.string("InvalidSeqNumber")
    BMCRevertNotExistsInternalHandler = sp.string("NotExistsInternalHandler")
    BMCRevertUnknownHandleBTPError = sp.string("UnknownHandleBTPError")
    BMCRevertUnknownHandleBTPMessage = sp.string("UnknownHandleBTPMessage")

    def __init__(self, bmc_management_addr, helper_contract, helper_parse_neg_contract, parse_address, owner_address, rlp_contract):
        self.init(
            helper=helper_contract,
            helper_parse_negative=helper_parse_neg_contract,
            bmc_btp_address=sp.none,
            bmc_management=bmc_management_addr,
            parse_contract=parse_address,
            handle_btp_message_status=sp.none,
            handle_btp_error_status=sp.none,
            handle_fee_gathering_status=sp.none,
            owner_address = owner_address,
            rlp_contract=rlp_contract,
        )

    def only_owner(self):
        sp.verify(sp.sender == self.data.owner_address, "Unauthorized")

    @sp.entry_point
    def set_helper_address(self, address):
        sp.set_type(address, sp.TAddress)
        self.only_owner()
        self.data.helper = address

    @sp.entry_point
    def set_helper_parse_negative_address(self, address):
        sp.set_type(address, sp.TAddress)
        self.only_owner()
        self.data.helper_parse_negative = address

    @sp.entry_point
    def set_parse_address(self, address):
        sp.set_type(address, sp.TAddress)
        self.only_owner()
        self.data.parse_contract = address

    @sp.entry_point
    def set_bmc_management_addr(self, params):
        sp.set_type(params, sp.TAddress)
        self.only_owner()
        self.data.bmc_management = params

    @sp.entry_point(lazify=False)
    def update_set_bmc_btp_address(self, ep):
        self.only_owner()
        sp.set_entry_point("set_bmc_btp_address", ep)

    @sp.entry_point(lazify=True)
    def set_bmc_btp_address(self, network):
        sp.set_type(network, sp.TString)

        sp.verify(sp.sender == self.data.bmc_management, "Unauthorized")
        with sp.if_(self.data.bmc_btp_address == sp.none):
            self.data.bmc_btp_address = sp.some(sp.string("btp://") + network + "/" + sp.view("add_to_str", self.data.parse_contract, sp.self_address, t=sp.TString).open_some())
        with sp.else_():
            sp.failwith("Address already set")

    @sp.entry_point
    def set_rlp_contract_address(self, param):
        sp.set_type(param, sp.TAddress)
        self.only_owner()
        self.data.rlp_contract = param

    @sp.onchain_view()
    def get_bmc_btp_address(self):
        sp.result(self.data.bmc_btp_address.open_some("Address not set"))

    @sp.onchain_view()
    def get_bmc_management_address(self):
        sp.result(self.data.bmc_management)

    def _require_registered_relay(self, prev):
        sp.set_type(prev, sp.TString)

        relays = sp.view("get_link_relays", self.data.bmc_management, prev, t=sp.TList(sp.TAddress)).open_some()
        valid = sp.local("valid", False)
        sp.for x in relays:
            sp.if sp.sender == x:
                valid.value = True
        sp.verify(valid.value, self.BMCRevertUnauthorized)

    @sp.entry_point(lazify=False)
    def update_callback_btp_message(self, ep):
        self.only_owner()
        sp.set_entry_point("callback_btp_message", ep)

    @sp.entry_point(lazify=True)
    def callback_btp_message(self, string, bsh_addr, prev, callback_msg):
        sp.set_type(string, sp.TOption(sp.TString))
        sp.set_type(bsh_addr, sp.TAddress)
        sp.set_type(prev, sp.TString)
        sp.set_type(callback_msg, types.Types.BMCMessage)

        sp.verify(sp.sender == bsh_addr, "Unauthorized")
        self.data.handle_btp_message_status = string

        with sp.if_(self.data.handle_btp_message_status.open_some() == "success"):
            pass
        with sp.else_():
            self._send_error(prev, callback_msg, self.BSH_ERR, self.BMCRevertUnknownHandleBTPMessage)
        self.data.handle_btp_message_status = sp.none

    @sp.entry_point(lazify=False)
    def update_callback_btp_error(self, ep):
        self.only_owner()
        sp.set_entry_point("callback_btp_error", ep)

    @sp.entry_point(lazify=True)
    def callback_btp_error(self, string, bsh_addr, svc, sn, code, msg):
        sp.set_type(string, sp.TOption(sp.TString))
        sp.set_type(bsh_addr, sp.TAddress)
        sp.set_type(svc, sp.TString)
        sp.set_type(sn, sp.TInt)
        sp.set_type(code, sp.TNat)
        sp.set_type(msg, sp.TString)

        sp.verify(sp.sender == bsh_addr, "Unauthorized")
        self.data.handle_btp_error_status = string

        with sp.if_(self.data.handle_btp_error_status.open_some() == "success"):
            pass
        with sp.else_():
            error_code = self.UNKNOWN_ERR
            err_msg = self.BMCRevertUnknownHandleBTPError
            sp.emit(sp.record(svc=svc, sn=sn, code=code, msg=msg, err_code=error_code, err_msg=err_msg), tag="ErrorOnBTPError")
        self.data.handle_btp_error_status = sp.none

    @sp.entry_point(lazify=False)
    def update_callback_handle_fee_gathering(self, ep):
        self.only_owner()
        sp.set_entry_point("callback_handle_fee_gathering", ep)

    @sp.entry_point(lazify=True)
    def callback_handle_fee_gathering(self, string, bsh_addr):
        sp.set_type(string, sp.TOption(sp.TString))
        sp.set_type(bsh_addr, sp.TAddress)

        sp.verify(sp.sender == bsh_addr, "Unauthorized")
        self.data.handle_fee_gathering_status = string
        with sp.if_(self.data.handle_fee_gathering_status.open_some() == "success"):
            pass
        self.data.handle_fee_gathering_status = sp.none

    @sp.entry_point(lazify=False)
    def update_handle_relay_message(self, ep):
        self.only_owner()
        sp.set_entry_point("handle_relay_message", ep)

    @sp.entry_point(lazify=True)
    def handle_relay_message(self, prev, msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(msg, sp.TBytes)

        self._require_registered_relay(prev)

        link_rx_seq = sp.view("get_link_rx_seq", self.data.bmc_management, prev, t=sp.TNat).open_some()
        link_rx_height = sp.view("get_link_rx_height", self.data.bmc_management, prev, t=sp.TNat).open_some()

        rx_seq = sp.local("rx_seq", link_rx_seq, t=sp.TNat)
        rx_height = sp.local("rx_height", link_rx_height, t=sp.TNat)
        rps = sp.view("decode_receipt_proofs", self.data.rlp_contract, msg, t=sp.TMap(sp.TNat, types.Types.ReceiptProof)).open_some()
        bmc_msg = sp.local("bmc_msg", sp.record(src="", dst="", svc="", sn=sp.int(0), message=sp.bytes("0x")), t=types.Types.BMCMessage)
        ev = sp.local("ev", sp.record(next_bmc="", seq=sp.nat(0), message=sp.bytes("0x")), t=types.Types.MessageEvent)
        sp.for i in sp.range(sp.nat(0), sp.len(rps)):
            with sp.if_(rps[i].height < rx_height.value):
               pass
            with sp.else_():
                rx_height.value = rps[i].height
                sp.for j in sp.range(sp.nat(0), sp.len(rps[i].events)):
                    ev.value = rps[i].events[j]
                    sp.verify(ev.value.next_bmc == self.data.bmc_btp_address.open_some("Address not set"), "Invalid Next BMC")
                    rx_seq.value +=sp.nat(1)
                    with sp.if_(ev.value.seq < rx_seq.value):
                        rx_seq.value = sp.as_nat(rx_seq.value-sp.nat(1))
                    with sp.else_():
                        sp.if ev.value.seq > rx_seq.value:
                            sp.failwith(self.BMCRevertInvalidSeqNumber)

                        _decoded = sp.view("decode_bmc_message", self.data.rlp_contract, ev.value.message, t=types.Types.BMCMessage).open_some()
                        bmc_msg.value = _decoded

                        sp.if bmc_msg.value.src != "":
                            with sp.if_(bmc_msg.value.dst == self.data.bmc_btp_address.open_some("Address not set")):
                                self._handle_message(prev, bmc_msg.value)
                            with sp.else_():
                                net, addr = sp.match_pair(strings.split_btp_address(bmc_msg.value.dst, "prev_idx", "result", "my_list", "last", "penultimate"))
                                next_link, prev_link = sp.match_pair(sp.view("resolve_route", self.data.bmc_management, net, t=sp.TPair(sp.TString, sp.TString)).open_some("Invalid Call"))
                                with sp.if_(next_link != "Unreachable"):
                                    self._send_message(next_link, ev.value.message)
                                with sp.else_():
                                    self._send_error(prev, bmc_msg.value, self.BMC_ERR, "Unreachable_"+ net)
        # call update_link_rx_seq on BMCManagement
        update_link_rx_seq_args_type = sp.TRecord(prev=sp.TString, val=sp.TNat)
        update_link_rx_seq_entry_point = sp.contract(update_link_rx_seq_args_type,
                                                          self.data.bmc_management,
                                                          "update_link_rx_seq").open_some()
        update_link_rx_seq_args = sp.record(prev=prev, val=sp.as_nat(rx_seq.value - link_rx_seq))
        sp.transfer(update_link_rx_seq_args, sp.tez(0), update_link_rx_seq_entry_point)

        # call update_relay_stats on BMCManagement
        update_relay_stats_args_type = sp.TRecord(relay=sp.TAddress, block_count_val=sp.TNat, msg_count_val=sp.TNat)
        update_relay_stats_entry_point = sp.contract(update_relay_stats_args_type,
                                                     self.data.bmc_management,
                                                     "update_relay_stats").open_some()
        update_relay_stats_args = sp.record(relay=sp.sender, block_count_val=sp.nat(0), msg_count_val=sp.as_nat(rx_seq.value - link_rx_seq))
        sp.transfer(update_relay_stats_args, sp.tez(0), update_relay_stats_entry_point)

        # call update_link_rx_height on BMCManagement
        update_link_rx_height_args_type = sp.TRecord(prev=sp.TString, val=sp.TNat)
        update_link_rx_height_entry_point = sp.contract(update_link_rx_height_args_type,
                                                     self.data.bmc_management,
                                                     "update_link_rx_height").open_some()
        update_link_rx_height_args = sp.record(prev=prev, val=sp.as_nat(rx_height.value - link_rx_height))
        sp.transfer(update_link_rx_height_args, sp.tez(0), update_link_rx_height_entry_point)



    def _handle_message(self, prev, msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(msg, types.Types.BMCMessage)

        # bsh_addr = sp.local("bsh_addr",sp.TAddress)
        with sp.if_(msg.svc == "bmc"):
            sm = sp.local("sm", sp.record(serviceType="", payload=sp.bytes("0x")))
            sm.value = sp.view("decode_bmc_service", self.data.rlp_contract, msg.message, t=types.Types.BMCService).open_some()
            with sp.if_(sm.value.serviceType == ""):
                self._send_error(prev, msg, self.BMC_ERR, self.BMCRevertParseFailure)
            with sp.else_():
                sp.if sm.value.serviceType == "FeeGathering":
                    gather_fee =sp.local("gather_fee", sp.record(fa="error", svcs=sp.map({0:""})))
                    gather_fee.value = sp.view("decode_gather_fee_message", self.data.rlp_contract, sm.value.payload, t=types.Types.GatherFeeMessage).open_some()

                    with sp.if_(gather_fee.value.fa == "error"):
                        self._send_error(prev, msg, self.BMC_ERR, self.BMCRevertParseFailure)

                    with sp.else_():
                        sp.for c in sp.range(sp.nat(0), sp.len(gather_fee.value.svcs)):
                            bsh_addr = sp.view("get_bsh_service_by_name", self.data.bmc_management, gather_fee.value.svcs[c], t=sp.TAddress).open_some("Invalid Call")

                            sp.if bsh_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"):
                                # call handle_fee_gathering of bts periphery
                                handle_fee_gathering_args_type = sp.TRecord(callback=sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress)), bsh_addr=sp.TAddress, fa=sp.TString, svc=sp.TString)
                                handle_fee_gathering_entry_point = sp.contract(handle_fee_gathering_args_type,
                                                                                bsh_addr,
                                                                                "handle_fee_gathering").open_some()
                                t = sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress)
                                callback = sp.contract(t, sp.self_address, "callback_handle_fee_gathering")
                                handle_fee_gathering_args = sp.record(callback=callback.open_some(), bsh_addr=bsh_addr, fa=gather_fee.value.fa, svc=gather_fee.value.svcs[c])
                                sp.transfer(handle_fee_gathering_args, sp.tez(0), handle_fee_gathering_entry_point)

                sp.if sm.value.serviceType == "Init":
                    links = sp.view("decode_init_message", self.data.rlp_contract, sm.value.payload, t=sp.TList(sp.TString)).open_some()
                    # call update_link_reachable on BMCManagement
                    update_link_reachable_args_type = sp.TRecord(prev=sp.TString, to=sp.TList(sp.TString))
                    update_link_reachable_entry_point = sp.contract(update_link_reachable_args_type,
                                                                    self.data.bmc_management,
                                                                    "update_link_reachable").open_some()
                    update_link_reachable_args = sp.record(prev=prev, to=links)
                    sp.transfer(update_link_reachable_args, sp.tez(0), update_link_reachable_entry_point)
        with sp.else_():
            bsh_addr = sp.view("get_bsh_service_by_name", self.data.bmc_management, msg.svc, t=sp.TAddress).open_some("Invalid view")

            with sp.if_(bsh_addr == sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                self._send_error(prev, msg, self.BMC_ERR, self.BMCRevertNotExistsBSH)

            with sp.else_():
                with sp.if_(msg.sn >= sp.int(0)):
                    net, addr = sp.match_pair(strings.split_btp_address(msg.src, "prev_idx", "result", "my_list", "last", "penultimate"))
                    # implemented callback
                    # call handle_btp_message on bts periphery
                    handle_btp_message_args_type = sp.TRecord(callback=sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress, prev=sp.TString, callback_msg=types.Types.BMCMessage)),
                                                              bsh_addr=sp.TAddress, prev=sp.TString, callback_msg=types.Types.BMCMessage, _from=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
                    handle_btp_message_entry_point = sp.contract(handle_btp_message_args_type,
                                                                    bsh_addr,
                                                                    "handle_btp_message").open_some()
                    t = sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress, prev=sp.TString, callback_msg=types.Types.BMCMessage )
                    callback = sp.contract(t, sp.self_address, "callback_btp_message")

                    handle_btp_message_args = sp.record(callback=callback.open_some(), bsh_addr=bsh_addr, prev=prev, callback_msg=msg, _from=net, svc=msg.svc, sn=msg.sn, msg=msg.message)
                    sp.transfer(handle_btp_message_args, sp.tez(0), handle_btp_message_entry_point)

                with sp.else_():
                    res = sp.view("decode_response", self.data.rlp_contract, msg.message, t=types.Types.Response).open_some()
                    # implemented callback
                    # call handle_btp_error on bts periphery
                    handle_btp_error_args_type = sp.TRecord(callback=sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress, svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)),
                        bsh_addr=sp.TAddress, svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)
                    handle_btp_error_entry_point = sp.contract(handle_btp_error_args_type,
                                                                 bsh_addr,
                                                                 "handle_btp_error").open_some()

                    t = sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress, svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)
                    callback = sp.contract(t, sp.self_address, "callback_btp_error")
                    handle_btp_error_args = sp.record(callback=callback.open_some(), bsh_addr=bsh_addr,
                                                      svc=msg.svc, sn=msg.sn * -1, code=res.code, msg=res.message)
                    sp.transfer(handle_btp_error_args, sp.tez(0), handle_btp_error_entry_point)

    def _send_message(self, to ,serialized_msg):
        sp.set_type(to, sp.TString)
        sp.set_type(serialized_msg, sp.TBytes)

        # call update_link_tx_seq on BMCManagement
        update_link_tx_seq_entry_point = sp.contract(sp.TRecord(prev=sp.TString, serialized_msg=sp.TBytes),
                                                     self.data.bmc_management,
                                                     "update_link_tx_seq").open_some()
        sp.transfer(sp.record(prev=to, serialized_msg=serialized_msg), sp.tez(0), update_link_tx_seq_entry_point)

    def _send_error(self, prev, message, err_code, err_msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(message, types.Types.BMCMessage)
        sp.set_type(err_code, sp.TNat)
        sp.set_type(err_msg, sp.TString)

        sp.if message.sn > sp.int(0):
            serialized_msg = sp.view("encode_bmc_message", self.data.rlp_contract, sp.record(
                src=self.data.bmc_btp_address.open_some("Address not set"),
                dst=message.src,
                svc=message.svc,
                sn=message.sn * -1,
                message=sp.view("encode_response", self.data.rlp_contract, sp.record(code=err_code, message=err_msg), t=sp.TBytes).open_some()), t=sp.TBytes).open_some()
            self._send_message(prev, serialized_msg)

    @sp.entry_point(lazify=False)
    def update_send_message(self, ep):
        self.only_owner()
        sp.set_entry_point("send_message", ep)

    @sp.entry_point(lazify=True)
    def send_message(self, to, svc, sn, msg):
        """
        Send the message to a specific network
        :param to: Network Address of destination network
        :param svc: Name of the service
        :param sn: Serial number of the message, it should be positive
        :param msg: Serialized bytes of Service Message
        :return:
        """
        sp.set_type(to, sp.TString)
        sp.set_type(svc, sp.TString)
        sp.set_type(sn, sp.TInt)
        sp.set_type(msg, sp.TBytes)

        sp.verify((sp.sender == self.data.bmc_management) |
                  (sp.view("get_bsh_service_by_name", self.data.bmc_management, svc, t=sp.TAddress).open_some() == sp.sender),
                  self.BMCRevertUnauthorized)
        sp.verify(sn >= sp.int(0), self.BMCRevertInvalidSn)

        next_link, dst = sp.match_pair(sp.view("resolve_route", self.data.bmc_management, to, t=sp.TPair(sp.TString, sp.TString)).open_some())

        rlp = sp.view("encode_bmc_message", self.data.rlp_contract, sp.record(
                src=self.data.bmc_btp_address.open_some("Address not set"),
                dst=dst,
                svc=svc,
                sn=sn,
                message=msg), t=sp.TBytes).open_some()
        self._send_message(next_link, rlp)

    @sp.onchain_view()
    def get_status(self, _link):
        """
        Get status of BMC
        :param _link: BTP Address of the connected BMC
        :return:
        """
        sp.set_type(_link, sp.TString)

        link = sp.view("get_link", self.data.bmc_management, _link, t=types.Types.Link).open_some()
        sp.verify(link.is_connected == True, self.BMCRevertNotExistsLink)

        sp.result(sp.record(
            rx_seq=link.rx_seq,
            tx_seq=link.tx_seq,
            rx_height=link.rx_height,
            current_height=sp.level #block height
        ))


sp.add_compilation_target("bmc_periphery", BMCPreiphery(bmc_management_addr=sp.address("KT1G3R9VqESejtsFnvjHSjzXYfuKuHMeaiE3"),
                                                        helper_contract=sp.address("KT1HwFJmndBWRn3CLbvhUjdupfEomdykL5a6"),
                                                        helper_parse_neg_contract=sp.address("KT1DHptHqSovffZ7qqvSM9dy6uZZ8juV88gP"),
                                                        parse_address=sp.address("KT1VJn3WNXDsyFxeSExjSWKBs9JYqRCJ1LFN"),
                                                        owner_address=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"),
                                                        rlp_contract=sp.address("KT1VyqgqhrcJZwML1WbgHYeUK6WSGt4qyj8Z")))
