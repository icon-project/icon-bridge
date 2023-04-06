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

    def __init__(self, network, bmc_management_addr):
        self.init(
            bmc_btp_address=sp.string("btp://") + network + "/" + "jj",
            bmc_management=bmc_management_addr
        )

    @sp.onchain_view()
    def get_bmc_btp_address(self):
        sp.result(self.data.bmc_btp_address)

    def _require_registered_relay(self, prev):
        sp.set_type(prev, sp.TString)

        sp.trace(prev)
        # relay = sp.view("get_link_relays", self.data.bmc_management, prev, t=sp.TList(sp.TAddress)).open_some()
        relay = []
        sp.for x in relay:
            sp.if sp.sender == x:
                return
        sp.fail_with(self.BMCRevertUnauthorized)

    @sp.entry_point
    def handle_relay_message(self, prev, msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(msg, sp.TBytes)

        self._require_registered_relay(prev)

        # link_rx_seq = sp.view("get_link_rx_seq", self.data.bmc_management, prev, t=sp.TNat).open_some()
        link_rx_seq= sp.nat(2)
        # link_rx_height = sp.view("get_link_rx_height", self.data.bmc_management, prev, t=sp.TNat).open_some()
        link_rx_height= sp.nat(3)

        rx_seq = sp.local("rx_seq", link_rx_seq, t=sp.TNat)
        rx_height = sp.local("rx_height", link_rx_height, t=sp.TNat)

        rps = sp.map(tkey=sp.TNat, tvalue=types.Types.ReceiptProof)
        # rsp = decodeReceiptProofs(msg)
        # rps = sp.map({sp.nat(0):sp.record(index=1, events=[], height=sp.nat(3))})
        # bmc_msg = sp.local("bmc_msg")

        # ev = sp.local("ev")

        sp.for i in sp.range(sp.nat(0), sp.len(rps)):
            sp.trace("ll")
            with sp.if_(rps[i].height < rx_height.value):
                sp.trace("ggg")
                # sp.continue

            rx_height.value = rps[i].height
            sp.for j in sp.range(sp.nat(0), sp.len(rps[i].events)):
                ev = rps[i].events[j]
                sp.verify(ev.next_bmc == self.data.bmc_btp_address, "Invalid Next BMC")
                rx_seq.value +=sp.nat(1)
                sp.if ev.seq < rx_seq.value:
                    rx_seq.value = sp.as_nat(rx_seq.value-sp.nat(1))
                    # sp.continue

                sp.if ev.seq > rx_seq.value:
                    sp.failwith(self.BMCRevertInvalidSeqNumber)
                # TODO: implement code inside of try catch

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
            sm = self.try_decode_bmc_service(msg.message)
            #TODO: implement try catch

            sp.if sm.serviceType == "FeeGathering":
                gather_fee = self.try_decode_gather_fee_message(sm.payload)

                sp.for i in sp.range(sp.nat(0), len(gather_fee.svcs)):
                    bsh_addr = sp.view("get_bsh_service_by_name", self.data.bmc_management, gather_fee.svcs[i], t=sp.TAddress)

                    sp.if bsh_addr != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"):
                        pass
                        #TODO: call BSH handleFeeGathering

            sp.if sm.serviceType == "Link":
                # to = decodePropagateMessage(sm.payload)
                to = "to"
                link = sp.view("get_link", self.data.bmc_management, prev, t=types.Types.Link).open_some()

                check = sp.local("check", False).value
                sp.if link.is_connected:
                    sp.for i in sp.range(sp.nat(0), len(link.reachable)):
                        sp.if to == link.reachable[i]:
                            check = True
                            # sp.break

                    sp.if not check:
                        links = sp.list([to], t=sp.TString)

                        # call update_link_reachable on BMCManagement
                        update_link_reachable_args_type = sp.TRecord(prev=sp.TString, to=sp.TList(sp.TString))
                        update_link_reachable_entry_point = sp.contract(update_link_reachable_args_type,
                                                                        self.data.bmc_management,
                                                                        "update_link_reachable").open_some()
                        update_link_reachable_args = sp.record(prev=prev, to=links)
                        sp.transfer(update_link_reachable_args, sp.tez(0), update_link_reachable_entry_point)

            sp.if sm.serviceType == "Unlink":
                # to = decodePropagateMessage(sm.payload)
                to = "to"
                link = sp.view("get_link", self.data.bmc_management, prev, t=types.Types.Link).open_some()

                sp.if link.is_connected:
                    sp.for i in sp.range(sp.nat(0), len(link.reachable)):
                        sp.if to == link.reachable[i]:

                            # call delete_link_reachable on BMCManagement
                            delete_link_reachable_args_type = sp.TRecord(prev=sp.TString, index=sp.TNat)
                            delete_link_reachable_entry_point = sp.contract(delete_link_reachable_args_type,
                                                                            self.data.bmc_management,
                                                                            "delete_link_reachable").open_some()
                            delete_link_reachable_args = sp.record(prev=prev, to=i)
                            sp.transfer(delete_link_reachable_args, sp.tez(0), delete_link_reachable_entry_point)

            sp.if sm.serviceType == "Init":
                # links = decodeInitMessage(sm.payload)
                links = ["link"]
                # call update_link_reachable on BMCManagement
                update_link_reachable_args_type = sp.TRecord(prev=sp.TString, to=sp.TList(sp.TString))
                update_link_reachable_entry_point = sp.contract(update_link_reachable_args_type,
                                                                self.data.bmc_management,
                                                                "update_link_reachable").open_some()
                update_link_reachable_args = sp.record(prev=prev, to=links)
                sp.transfer(update_link_reachable_args, sp.tez(0), update_link_reachable_entry_point)
        with sp.else_():
            bsh_addr = sp.view("get_bsh_service_by_name", self.data.bmc_management, msg.svc, t=sp.TAddress)

            sp.if bsh_addr == sp.addrress("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"):
                self._send_error(prev, msg, self.BMC_ERR, self.BMCRevertNotExistsBSH)
                return

            with sp.if_(msg.sn >= sp.nat(0)):
                net, = strings.split_btp_address(msg.src)
                #TODO: implement try catch, call handleBTPMessage from BSH


            with sp.else_():
                res = decodeResponse(msg.message)
                # TODO: implement try catch, call handleBTPError from BSH


    def try_decode_btp_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        return decodeBMCMessage(rpl)

    def try_decode_bmc_service(self, msg):
        sp.set_type(msg, sp.TBytes)

        return decodeBMCService(msg)

    def try_decode_gather_fee_message(self, msg):
        sp.set_type(msg, sp.TBytes)

        return decodeGatherFeeMessage(msg)

    def _send_message(self, to ,serialized_msg):
        sp.set_type(to, sp.TString)
        sp.set_type(serialized_msg, sp.TBytes)

        # call update_link_tx_seq on BMCManagement
        update_link_tx_seq_args_type = sp.TRecord(prev=sp.TString)
        update_link_tx_seq_entry_point = sp.contract(update_link_tx_seq_args_type,
                                                     self.data.bmc_management,
                                                     "update_link_tx_seq").open_some()
        update_link_tx_seq_args = sp.record(prev=to)
        sp.transfer(update_link_tx_seq_args, sp.tez(0), update_link_tx_seq_entry_point)

        sp.emit(sp.record(next=to, seq=sp.view("get_link_tx_seq", self.data.bmc_management, to, t=sp.TNat).open_some(), msg=serialized_msg))

    def _send_error(self, prev, message, err_code, err_msg):
        sp.set_type(prev, sp.TString)
        sp.set_type(message, types.Types.BMCMessage)
        sp.set_type(err_code, sp.TNat)
        sp.set_type(err_msg, sp.TString)

        if message.sn > sp.nat(0):
            serialized_msg = encode_bmc_message(sp.record(
                src=self.data.bmc_btp_address,
                dst=message.src,
                svc=message.svc,
                sn=message.sn * - sp.nat(1),
                message=encode_response(sp.record(code=err_code, message=err_msg))
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
        sp.set_type(sn, sp.TNat)
        sp.set_type(msg, sp.TBytes)

        sp.verify((sp.sender == self.data.bmc_management) |
                  (sp.view("get_bsh_service_by_name", self.data.bmc_management, svc, t=sp.TAddress).open_some() == sp.sender),
                  self.BMCRevertUnauthorized)
        sp.verify(sn >= sp.nat(0), self.BMCRevertInvalidSn)

        next_link, dst = sp.match_pair(sp.view("resolve_route", self.data.bmc_management, to, t=sp.TPair(sp.TString, sp.TString)).open_some())

        # need to import encode_bmc_message from library
        # rlp = encode_bmc_message(sp.record(
        #         src=self.data.bmc_btp_address,
        #         dst=dst,
        #         svc=svc,
        #         sn=sn,
        #         message=msg
        # ))
        next_link = "next_link"
        rlp = sp.bytes("0x0dae11")
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
    bmc_management = sp.test_account("BMC Management")
    # bmc= sp.test_account("BMC")

    scenario = sp.test_scenario()
    bmc = BMCPreiphery("tezos", bmc_management.address)
    scenario += bmc

    bmc.handle_relay_message(sp.record(prev="demo string", msg=sp.bytes("0x0dae11"))).run(sender=alice)

sp.add_compilation_target("bmc_periphery", BMCPreiphery(network="tezos",
                                                        bmc_management_addr=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW")))
