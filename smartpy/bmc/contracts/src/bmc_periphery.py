import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")
rlp_encode = sp.io.import_script_from_url("file:./contracts/src/RLP_encode_struct.py")
rlp_decode = sp.io.import_script_from_url("file:./contracts/src/RLP_decode_struct.py")


class BMCPreiphery(sp.Contract, rlp_decode.DecodeLibrary, rlp_encode.EncodeLibrary):
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

    def __init__(self, bmc_management_addr, helper_contract, parse_address, owner_address):
        self.init(
            helper=helper_contract,
            bmc_btp_address=sp.none,
            bmc_management=bmc_management_addr,
            parse_contract=parse_address,
            handle_btp_message_status=sp.none,
            handle_btp_error_status=sp.none,
            owner_address = owner_address
        )

    def only_owner(self):
        sp.verify(sp.sender == self.data.owner_address, "Unauthorized")

    @sp.entry_point
    def set_bmc_management_addr(self, params):
        sp.set_type(params, sp.TAddress)
        self.only_owner()
        self.data.bmc_management = params

    @sp.entry_point
    def set_bmc_btp_address(self, network):
        sp.set_type(network, sp.TString)

        sp.verify(sp.sender == self.data.bmc_management, "Unauthorized")
        with sp.if_(self.data.bmc_btp_address == sp.none):
            self.data.bmc_btp_address = sp.some(sp.string("btp://") + network + "/" + sp.view("add_to_str", self.data.parse_contract, sp.self_address, t=sp.TString).open_some())
        with sp.else_():
            sp.failwith("Address already set")

    @sp.onchain_view()
    def get_bmc_btp_address(self):
        sp.result(self.data.bmc_btp_address.open_some("Address not set"))

    def _require_registered_relay(self, prev):
        sp.set_type(prev, sp.TString)

        relays = sp.view("get_link_relays", self.data.bmc_management, prev, t=sp.TList(sp.TAddress)).open_some()
        valid = sp.local("valid", False)
        sp.for x in relays:
            sp.if sp.sender == x:
                valid.value = True
        sp.verify(valid.value, self.BMCRevertUnauthorized)

    @sp.entry_point
    def callback(self, string, bsh_addr, prev, callback_msg):
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

    @sp.entry_point
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
            sp.emit(sp.record(svc=svc, sn=sn * -1, code=code, msg=msg, err_code=error_code, err_msg=err_msg), tag="ErrorOnBTPError")

    @sp.entry_point
    def handle_relay_message(self, prev, msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(msg, sp.TBytes)

        self._require_registered_relay(prev)

        link_rx_seq = sp.view("get_link_rx_seq", self.data.bmc_management, prev, t=sp.TNat).open_some()
        link_rx_height = sp.view("get_link_rx_height", self.data.bmc_management, prev, t=sp.TNat).open_some()

        rx_seq = sp.local("rx_seq", link_rx_seq, t=sp.TNat)
        rx_height = sp.local("rx_height", link_rx_height, t=sp.TNat)

        rps = self.decode_receipt_proofs(msg)
        bmc_msg = sp.local("bmc_msg", sp.record(src="", dst="", svc="", sn=sp.int(0), message=sp.bytes("0x")), t=types.Types.BMCMessage)
        ev = sp.local("ev", sp.record(next_bmc="", seq=sp.nat(0), message=sp.bytes("0x")), t=types.Types.MessageEvent)

        sp.for i in sp.range(sp.nat(0), sp.len(rps)):
            sp.trace("ll")
            with sp.if_(rps[i].height < rx_height.value):
                sp.trace("ggg")
                # sp.continue

            rx_height.value = rps[i].height
            sp.for j in sp.range(sp.nat(0), sp.len(rps[i].events)):
                ev.value = rps[i].events[j]
                sp.verify(ev.value.next_bmc == self.data.bmc_btp_address.open_some("Address not set"), "Invalid Next BMC")
                rx_seq.value +=sp.nat(1)
                sp.if ev.value.seq < rx_seq.value:
                    rx_seq.value = sp.as_nat(rx_seq.value-sp.nat(1))
                    # sp.continue

                sp.if ev.value.seq > rx_seq.value:
                    sp.failwith(self.BMCRevertInvalidSeqNumber)

                _decoded = self.decode_bmc_message(ev.value.message)
                bmc_msg.value = _decoded

                sp.if bmc_msg.value.src != "":
                    with sp.if_(bmc_msg.value.dst == self.data.bmc_btp_address.open_some("Address not set")):
                        self._handle_message(prev, bmc_msg.value)
                    with sp.else_():
                        net, addr = sp.match_pair(strings.split_btp_address(bmc_msg.value.dst, "prev_idx", "result", "my_list", "last", "penultimate"))
                        # resolve route inside try catch
                        next_link, prev_link = sp.match_pair(sp.view("resolve_route", self.data.bmc_management, net, t=sp.TPair(sp.TString, sp.TString)).open_some("Invalid Call"))
                        self._send_message(next_link, ev.value.message)

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
            sm.value = self.decode_bmc_service(msg.message)
            with sp.if_(sm.value.serviceType == ""):
                self._send_error(prev, msg, self.BMC_ERR, self.BMCRevertParseFailure)
            with sp.else_():
                sp.if sm.value.serviceType == "FeeGathering":
                    gather_fee =sp.local("gather_fee", sp.record(fa="", svcs=sp.map({0:""})))
                    gather_fee.value = self.decode_gather_fee_message(sm.value.payload)

                    with sp.if_(gather_fee.value.fa == ""):
                        self._send_error(prev, msg, self.BMC_ERR, self.BMCRevertParseFailure)

                    with sp.else_():
                        sp.for c in sp.range(sp.nat(0), sp.len(gather_fee.value.svcs)):
                            bsh_addr = sp.view("get_bsh_service_by_name", self.data.bmc_management, gather_fee.value.svcs[c], t=sp.TAddress).open_some("Invalid Call")

                            sp.if bsh_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"):
                                # call handle_fee_gathering of bts periphery
                                handle_fee_gathering_args_type = sp.TRecord(fa=sp.TString, svc=sp.TString)
                                handle_fee_gathering_entry_point = sp.contract(handle_fee_gathering_args_type,
                                                                                bsh_addr,
                                                                                "handle_fee_gathering").open_some()
                                handle_fee_gathering_args = sp.record(fa=gather_fee.value.fa, svc=gather_fee.value.svcs[c])
                                sp.transfer(handle_fee_gathering_args, sp.tez(0), handle_fee_gathering_entry_point)

                sp.if sm.value.serviceType == "Link":
                    to = self.decode_propagate_message(sm.value.payload)
                    link = sp.view("get_link", self.data.bmc_management, prev, t=types.Types.Link).open_some()

                    check = sp.local("check", False)
                    sp.if link.is_connected:
                        sp.for e in link.reachable.elements():
                            sp.if to == e:
                                check.value = True
                                # sp.break

                        sp.if check.value == False:
                            links = sp.list([to], t=sp.TString)

                            # call update_link_reachable on BMCManagement
                            update_link_reachable_args_type = sp.TRecord(prev=sp.TString, to=sp.TList(sp.TString))
                            update_link_reachable_entry_point = sp.contract(update_link_reachable_args_type,
                                                                            self.data.bmc_management,
                                                                            "update_link_reachable").open_some()
                            update_link_reachable_args = sp.record(prev=prev, to=links)
                            sp.transfer(update_link_reachable_args, sp.tez(0), update_link_reachable_entry_point)

                sp.if sm.value.serviceType == "Unlink":
                    to = self.decode_propagate_message(sm.value.payload)
                    link = sp.view("get_link", self.data.bmc_management, prev, t=types.Types.Link).open_some()

                    sp.if link.is_connected:
                        f = sp.local("f", sp.nat(0))
                        sp.for itm in link.reachable.elements():
                            sp.if to == itm:

                                # call delete_link_reachable on BMCManagement
                                delete_link_reachable_args_type = sp.TRecord(prev=sp.TString, index=sp.TNat)
                                delete_link_reachable_entry_point = sp.contract(delete_link_reachable_args_type,
                                                                                self.data.bmc_management,
                                                                                "delete_link_reachable").open_some()
                                delete_link_reachable_args = sp.record(prev=prev, index=f.value)
                                sp.transfer(delete_link_reachable_args, sp.tez(0), delete_link_reachable_entry_point)
                                f.value += sp.nat(1)

                sp.if sm.value.serviceType == "Init":
                    links = self.decode_init_message(sm.value.payload)
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
                    handle_btp_message_args = sp.record(callback=sp.self_entry_point("callback"), bsh_addr=bsh_addr, prev=prev, callback_msg=msg, _from=net, svc=msg.svc, sn=msg.sn, msg=msg.message)
                    sp.transfer(handle_btp_message_args, sp.tez(0), handle_btp_message_entry_point)

                with sp.else_():
                    res = self.decode_response(msg.message)
                    # implemented callback
                    # call handle_btp_error on bts periphery
                    handle_btp_error_args_type = sp.TRecord(callback=sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress, svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)),
                        bsh_addr=sp.TAddress, svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)
                    handle_btp_error_entry_point = sp.contract(handle_btp_error_args_type,
                                                                 bsh_addr,
                                                                 "handle_btp_error").open_some()
                    handle_btp_error_args = sp.record(callback=sp.self_entry_point("callback_btp_error"), bsh_addr=bsh_addr,
                                                      svc=msg.svc, sn=msg.sn * -1, code=res.code, msg=res.message)
                    sp.transfer(handle_btp_error_args, sp.tez(0), handle_btp_error_entry_point)

    def _send_message(self, to ,serialized_msg):
        sp.set_type(to, sp.TString)
        sp.set_type(serialized_msg, sp.TBytes)

        # call update_link_tx_seq on BMCManagement
        update_link_tx_seq_entry_point = sp.contract(sp.TString,
                                                     self.data.bmc_management,
                                                     "update_link_tx_seq").open_some()
        sp.transfer(to, sp.tez(0), update_link_tx_seq_entry_point)

        sp.emit(sp.record(next=to, seq=sp.view("get_link_tx_seq", self.data.bmc_management, to, t=sp.TNat).open_some(), msg=serialized_msg),
                tag="Message")

    def _send_error(self, prev, message, err_code, err_msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(message, types.Types.BMCMessage)
        sp.set_type(err_code, sp.TNat)
        sp.set_type(err_msg, sp.TString)

        sp.if message.sn > sp.int(0):
            serialized_msg = self.encode_bmc_message(sp.record(
                src=self.data.bmc_btp_address.open_some("Address not set"),
                dst=message.src,
                svc=message.svc,
                sn=message.sn * -1,
                message=self.encode_response(sp.record(code=err_code, message=err_msg))
            ))
            self._send_message(prev, serialized_msg)

    @sp.entry_point
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

        rlp = self.encode_bmc_message(sp.record(
                src=self.data.bmc_btp_address.open_some("Address not set"),
                dst=dst,
                svc=svc,
                sn=sn,
                message=msg
        ))
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

@sp.add_test(name="BMC")
def test():
    alice = sp.test_account("Alice")
    helper = sp.test_account("Helper")
    parse_contract = sp.test_account("Parser")
    owner = sp.test_account("Owner")
    bmc_management = sp.test_account("BMC Management")
    # bmc= sp.test_account("BMC")

    scenario = sp.test_scenario()
    bmc = BMCPreiphery(bmc_management.address, helper.address, parse_contract.address, owner.address)
    scenario += bmc

    # bmc.handle_relay_message(sp.record(prev="demo string", msg=sp.bytes("0x0dae11"))).run(sender=alice)

sp.add_compilation_target("bmc_periphery", BMCPreiphery(bmc_management_addr=sp.address("KT1Uiycjx4iXdjKFfR2kAo2NUdEtQ6PmDX4Y"),
                                                        helper_contract=sp.address("KT1Q5erZm7Pp8UJywK1nkiP8QPCRmyUotUMq"),
                                                        parse_address=sp.address("KT1XgRyjQPfpfwNrvYYpgERpYpCrGh24aoPX"),
                                                        owner_address=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")))
