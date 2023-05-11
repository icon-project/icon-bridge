import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")
rlp_encode = sp.io.import_script_from_url("file:./contracts/src/RLP_encode_struct.py")
rlp_decode = sp.io.import_script_from_url("file:./contracts/src/RLP_decode_struct.py")


class BTPPreiphery(sp.Contract, rlp_decode.DecodeLibrary, rlp_encode.EncodeLibrary):
    service_name = sp.string("bts")

    RC_OK = sp.nat(0)
    RC_ERR = sp.nat(1)

    MAX_BATCH_SIZE = sp.nat(15)

    def __init__(self, bmc_address, bts_core_address, helper_contract, parse_address):
        self.update_initial_storage(
            bmc=bmc_address,
            bts_core=bts_core_address,
            blacklist=sp.map(tkey=sp.TAddress, tvalue=sp.TBool),
            token_limit=sp.map(tkey=sp.TString, tvalue=sp.TNat),
            requests=sp.big_map(tkey=sp.TNat, tvalue=types.Types.PendingTransferCoin),
            serial_no = sp.nat(0),
            number_of_pending_requests = sp.nat(0),
            helper=helper_contract,
            parse_contract=parse_address
        )

    def only_bmc(self):
        sp.verify(sp.sender == self.data.bmc, "Unauthorized")

    def only_bts_core(self):
        sp.verify(sp.sender == self.data.bts_core, "Unauthorized")

    @sp.onchain_view()
    def has_pending_request(self):
        """

        :return: boolean
        """
        sp.result(self.data.number_of_pending_requests != sp.nat(0))

    @sp.entry_point
    def add_to_blacklist(self, params):
        """

        :param params: List of addresses to be blacklisted
        :return:
        """
        sp.set_type(params, sp.TMap(sp.TNat, sp.TString))
        sp.verify(sp.sender == sp.self_address, "Unauthorized")
        sp.verify(sp.len(params) <= self.MAX_BATCH_SIZE, "BatchMaxSizeExceed")

        sp.for i in sp.range(sp.nat(0), sp.len(params)):
            parsed_addr = sp.view("str_to_addr", self.data.parse_contract, params.get(i), t=sp.TAddress).open_some()
            with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                self.data.blacklist[parsed_addr] = True
            with sp.else_():
                sp.failwith("InvalidAddress")


    @sp.entry_point
    def remove_from_blacklist(self, params):
        """
        :param params: list of address strings
        :return:
        """
        sp.set_type(params, sp.TMap(sp.TNat, sp.TString))

        sp.verify(sp.sender == sp.self_address, "Unauthorized")
        sp.verify(sp.len(params) <= self.MAX_BATCH_SIZE, "BatchMaxSizeExceed")

        sp.for i in sp.range(sp.nat(0), sp.len(params)):
            parsed_addr = sp.view("str_to_addr", self.data.parse_contract, params.get(i), t=sp.TAddress).open_some()
            with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                sp.verify(self.data.blacklist.contains(parsed_addr), "UserNotFound")
                sp.verify(self.data.blacklist.get(parsed_addr) == True, "UserNotBlacklisted")
                del self.data.blacklist[parsed_addr]
            with sp.else_():
                sp.failwith("InvalidAddress")

    @sp.entry_point
    def set_token_limit(self, coin_names, token_limit):
        """
        :param coin_names: list of coin names
        :param token_limit: list of token limits
        :return:
        """
        sp.set_type(coin_names, sp.TMap(sp.TNat, sp.TString))
        sp.set_type(token_limit, sp.TMap(sp.TNat, sp.TNat))

        sp.verify((sp.sender == sp.self_address )| (sp.sender == self.data.bts_core), "Unauthorized")
        sp.verify(sp.len(coin_names) == sp.len(token_limit), "InvalidParams")
        sp.verify(sp.len(coin_names) <= self.MAX_BATCH_SIZE, "BatchMaxSizeExceed")

        sp.for i in sp.range(0, sp.len(coin_names)):
            self.data.token_limit[coin_names[i]] = token_limit.get(i)


    @sp.entry_point
    def send_service_message(self, _from, to, coin_names, values, fees):
        """
        Send service message to BMC
        :param _from: from address
        :param to: to address
        :param coin_names:
        :param values:
        :param fees:
        :return:
        """

        sp.set_type(_from, sp.TAddress)
        sp.set_type(to, sp.TString)
        sp.set_type(coin_names, sp.TMap(sp.TNat, sp.TString))
        sp.set_type(values, sp.TMap(sp.TNat, sp.TNat))
        sp.set_type(fees, sp.TMap(sp.TNat, sp.TNat))

        self.only_bts_core()

        to_network, to_address = sp.match_pair(strings.split_btp_address(to))

        assets = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.Asset))
        assets_details = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.AssetTransferDetail))
        sp.for i in sp.range(sp.nat(0), sp.len(coin_names)):
            assets[i]=sp.record(
                coin_name=coin_names[i],
                value=values[i]
            )
            assets_details[i] = sp.record(
                coin_name=coin_names[i],
                value=values[i],
                fee=fees[i]
            )

        self.data.serial_no += 1

        start_from = sp.view("add_to_str", self.data.parse_contract, _from, t=sp.TString).open_some()

        send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TNat, msg=sp.TBytes)
        send_message_entry_point = sp.contract(send_message_args_type, self.data.bmc, "send_message").open_some()
        send_message_args = sp.record(
            to=to_network, svc=self.service_name, sn=self.data.serial_no,
            msg=self.encode_service_message(sp.compute(sp.record(serviceType=(sp.variant("REQUEST_COIN_TRANSFER", 0)),
                                  data=self.encode_transfer_coin_msg(sp.compute(sp.record(from_addr=start_from, to=to_address, assets=assets)))
                                            )
                                    )
                        )
        )
        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)

        # push pending tx into record list
        self.data.requests[self.data.serial_no] = sp.record(
            from_=start_from, to=to, coin_names=coin_names, amounts=values, fees=fees
        )
        self.data.number_of_pending_requests +=sp.nat(1)
        sp.emit(sp.record(from_address=_from, to=to, serial_no=self.data.serial_no, assets_details=assets_details), tag="TransferStart")


    @sp.entry_point
    def handle_btp_message(self, _from, svc, sn, msg):
        """
        BSH handle BTP message from BMC contract
        :param _from: An originated network address of a request
        :param svc: A service name of BSH contract
        :param sn: A serial number of a service request
        :param msg: An RLP message of a service request/service response
        :return:
        """

        sp.set_type(_from, sp.TString)
        sp.set_type(svc, sp.TString)
        sp.set_type(sn, sp.TNat)
        sp.set_type(msg, sp.TBytes)

        self.only_bmc()

        sp.verify(svc == self.service_name, "InvalidSvc")
        err_msg = sp.local("error", "")
        sm = self.decode_service_message(msg)

        # TODO: implement try in below cases
        with sm.serviceType.match_cases() as arg:
            with arg.match("REQUEST_COIN_TRANSFER") as a1:
                tc = self.decode_transfer_coin_msg(sm.data)
                parsed_addr = sp.view("str_to_addr", self.data.parse_contract, tc.to, t=sp.TAddress).open_some()

                with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                    contract = sp.self_entry_point(entry_point="handle_request_service")
                    _param = sp.record(to=tc.to, assets=tc.assets)
                    sp.transfer(_param, sp.tez(0), contract)
                    self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), _from, sn, "", self.RC_OK)
                    sp.emit(sp.record(from_address=_from, to=parsed_addr, serial_no=self.data.serial_no, assets_details=tc.assets), tag="TransferReceived")
                    # return
                with sp.else_():
                    err_msg.value = "InvalidAddress"

                self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), _from, sn, "", self.RC_ERR)

            with arg.match("BLACKLIST_MESSAGE") as a2:
                bm = self.decode_blacklist_msg(sm.data)
                addresses = bm.addrs

                with bm.serviceType.match_cases() as b_agr:
                    with b_agr.match("ADD_TO_BLACKLIST") as b_val_1:
                        contract = sp.self_entry_point(entry_point="add_to_blacklist")
                        sp.transfer(addresses, sp.tez(0), contract)
                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), _from, sn, "AddedToBlacklist", self.RC_OK)

                    with b_agr.match("REMOVE_FROM_BLACKLIST") as b_val_2:
                        contract = sp.self_entry_point(entry_point="remove_from_blacklist")
                        sp.transfer(addresses, sp.tez(0), contract)
                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), _from, sn, "RemovedFromBlacklist", self.RC_OK)

            with arg.match("CHANGE_TOKEN_LIMIT") as a3:
                tl = self.decode_token_limit_msg(sm.data)
                coin_names = tl.coin_name
                token_limits = tl.token_limit

                contract = sp.self_entry_point(entry_point="set_token_limit")
                _param = sp.record(coin_names=coin_names, token_limit=token_limits)
                sp.transfer(_param, sp.tez(0), contract)
                self.send_response_message(sp.variant("CHANGE_TOKEN_LIMIT", 4), _from, sn, "ChangeTokenLimit",self.RC_OK)

            with arg.match("RESPONSE_HANDLE_SERVICE") as a4:
                sp.verify(sp.len(sp.pack(self.data.requests.get(sn).from_)) != 0, "InvalidSN")
                response = self.decode_response(sm.data)
                self.handle_response_service(sn, response.code, response.message)

            with arg.match("UNKNOWN_TYPE") as a5:
                sp.emit(sp.record(_from=_from, sn=sn), tag= "UnknownResponse")

        # using if else
        # with sp.if_(sm.serviceType == types.Types.ServiceType.open_variant("REQUEST_COIN_TRANSFER")):
            # tc = sp.unpack(sm.data, t=types.Types.TransferCoin)
            # TDo: check address and handle request service
            # # with sp.if_(self.check_parse_address(tc.to)):
            # #     with sp.if_(self.handle_request_service(tc.to, tc.assets)):
            # self.send_response_message(types.Types.ServiceType.open_variant("REQUEST_COIN_TRANSFER"), _from, sn, "", self.RC_OK)
            # sp.emit(sp.record(from_address=_from, to=tc.to, serial_no=self.data.serial_no, assets_details=tc.assets))
            #     # with sp.else_():
            #     #     err_msg = "ErrorWhileMinting"
            # # with sp.else_():
            # err_msg = "InvalidAddress"
            # sp.trace(err_msg)
            # self.send_response_message(types.Types.ServiceType.open_variant("REQUEST_COIN_TRANSFER"), _from, sn, err_msg, self.RC_ERR)


    @sp.entry_point
    def handle_btp_error(self, svc, sn, code, msg):
        """
        BSH handle BTP Error from BMC contract
        :param svc: A service name of BSH contract
        :param sn: A serial number of a service request
        :param code: A response code of a message (RC_OK / RC_ERR)
        :param msg: A response message
        :return:
        """

        sp.set_type(svc, sp.TString)
        sp.set_type(sn, sp.TNat)
        sp.set_type(code, sp.TNat)
        sp.set_type(msg, sp.TString)

        self.only_bmc()

        sp.verify(svc == self.service_name, "InvalidSvc")
        sp.verify(sp.len(sp.pack(self.data.requests.get(sn).from_)) != 0, "InvalidSN")

        emit_msg= sp.concat(["errCode: ", sp.view("string_of_int", self.data.parse_contract, sp.to_int(code), t=sp.TString).open_some(),", errMsg: ", msg])
        self.handle_response_service(sn, self.RC_ERR, emit_msg)

    def handle_response_service(self, sn, code, msg):
        """

        :param sn:
        :param code:
        :param msg:
        :return:
        """
        sp.set_type(sn, sp.TNat)
        sp.set_type(code, sp.TNat)
        sp.set_type(msg, sp.TString)

        caller = sp.local("caller", sp.view("str_to_addr", self.data.parse_contract, self.data.requests.get(sn).from_, t=sp.TAddress).open_some()
                          , sp.TAddress).value
        loop = sp.local("loop", sp.len(self.data.requests.get(sn).coin_names), sp.TNat).value
        sp.verify(loop <= self.MAX_BATCH_SIZE, "BatchMaxSizeExceed")

        sp.for i in sp.range(0, loop):
            # inter score call
            handle_response_service_args_type = sp.TRecord(
                requester=sp.TAddress, coin_name=sp.TString, value=sp.TNat, fee=sp.TNat, rsp_code=sp.TNat
            )
            handle_response_service_entry_point = sp.contract(handle_response_service_args_type, self.data.bts_core, "handle_response_service").open_some("invalid call")
            handle_response_service_args = sp.record(
                requester=caller, coin_name=self.data.requests.get(sn).coin_names.get(i), value=self.data.requests.get(sn).amounts.get(i),
                fee=self.data.requests.get(sn).fees.get(i), rsp_code=code
            )
            sp.transfer(handle_response_service_args, sp.tez(0), handle_response_service_entry_point)

        del self.data.requests[sn]
        self.data.number_of_pending_requests = sp.as_nat(self.data.number_of_pending_requests-1)

        sp.emit(sp.record(caller=caller, sn=sn, code=code, msg=msg), tag="TransferEnd")

    @sp.entry_point
    def handle_request_service(self, to, assets):
        """
        Handle a list of minting/transferring coins/tokens
        :param to: An address to receive coins/tokens
        :param assets:  A list of requested coin respectively with an amount
        :return:
        """
        sp.set_type(to, sp.TString)
        sp.set_type(assets, sp.TMap(sp.TNat, types.Types.Asset))

        sp.verify(sp.sender == sp.self_address, "Unauthorized")
        sp.verify(sp.len(assets) <= self.MAX_BATCH_SIZE, "BatchMaxSizeExceed")

        parsed_to = sp.view("str_to_addr", self.data.parse_contract, to, t=sp.TAddress).open_some()
        sp.for i in sp.range(0, sp.len(assets)):
            valid_coin = sp.view("is_valid_coin", self.data.bts_core, assets[i].coin_name, t=sp.TBool).open_some()
            sp.verify(valid_coin == True, "UnregisteredCoin")

            check_transfer = sp.view("check_transfer_restrictions", sp.self_address, sp.record(
                coin_name=assets[i].coin_name,user=parsed_to, value=assets[i].value), t=sp.TBool).open_some()
            sp.verify(check_transfer == True, "FailCheckTransfer")

            # TODO: implement try for mint
            # inter score call
            mint_args_type = sp.TRecord(
                to=sp.TAddress, coin_name=sp.TString, value=sp.TNat
            )
            mint_args_type_entry_point = sp.contract(mint_args_type, self.data.bts_core, "mint").open_some()
            mint_args = sp.record(
                to=parsed_to, coin_name=assets[i].coin_name, value=assets[i].value
            )
            sp.transfer(mint_args, sp.tez(0), mint_args_type_entry_point)


    def send_response_message(self, service_type, to, sn, msg, code):
        """

        :param service_type:
        :param to:
        :param sn:
        :param msg:
        :param code:
        :return:
        """
        sp.set_type(service_type, types.Types.ServiceType)
        sp.set_type(to, sp.TString)
        sp.set_type(sn, sp.TNat)
        sp.set_type(msg, sp.TString)
        sp.set_type(code, sp.TNat)

        sp.trace("in send_response_message")

        send_message_args_type = sp.TRecord(
            to=sp.TString, svc=sp.TString, sn=sp.TNat, msg=sp.TBytes
        )
        send_message_entry_point = sp.contract(send_message_args_type, self.data.bmc, "send_message").open_some()
        send_message_args = sp.record(to=to, svc=self.service_name, sn=self.data.serial_no,
                                      msg=sp.pack(sp.record(serviceType=service_type, data=sp.pack(sp.record(code=code, message=msg))))
                                      )
        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)


    @sp.entry_point
    def handle_fee_gathering(self, fa, svc):
        """
        BSH handle Gather Fee Message request from BMC contract
        :param fa: A BTP address of fee aggregator
        :param svc: A name of the service
        :return:
        """
        sp.set_type(fa, sp.TString)
        sp.set_type(svc, sp.TString)

        self.only_bmc()
        sp.verify(svc == self.service_name, "InvalidSvc")

        strings.split_btp_address(fa)

        # call transfer_fees of BTS_Core
        transfer_fees_args_type = sp.TString
        transfer_fees_entry_point = sp.contract(transfer_fees_args_type, self.data.bts_core, "transfer_fees").open_some()
        sp.transfer(fa, sp.tez(0), transfer_fees_entry_point)


    @sp.onchain_view()
    def check_transfer_restrictions(self, params):
        """

        :param params: Record of coin transfer details
        :return:
        """
        sp.set_type(params, sp.TRecord(coin_name=sp.TString, user=sp.TAddress, value=sp.TNat))

        sp.verify(self.data.blacklist.contains(params.user) == False, "Blacklisted")
        sp.verify(self.data.token_limit.get(params.coin_name) >= params.value, "LimitExceed")
        sp.result(True)



@sp.add_test(name="Counter")
def test():
    bmc = sp.test_account("Alice")
    admin = sp.test_account("Admin")
    bts_core = sp.test_account("BTS")
    helper = sp.test_account("Helper")

    scenario = sp.test_scenario()
    counter = BTPPreiphery(bmc.address, bts_core.address, helper.address, admin.address)
    scenario += counter

    # counter.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=counter.address)
    # counter.send_service_message(sp.record(_from=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), to="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW",
    #                                        coin_names={0:"Tok1"}, values={0:sp.nat(10)}, fees={0:sp.nat(2)})).run(
    #     sender=bts_core
    # )
    # counter.handle_btp_error(sp.record(svc= "bts", code=sp.nat(2), sn=sp.nat(1), msg="test 1")).run(
    #     sender=bmc
    # )

    # counter.set_token_limit(sp.record(coin_names={0:"Tok2"}, token_limit={0:sp.nat(5)})).run(sender=counter.address)

    # counter.handle_request_service(sp.record(to= "tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc", assets={0:
    #                                          sp.record(coin_name="Tok2", value=sp.nat(4))})).run(
    #     sender=counter.address
    # )

    # counter.handle_fee_gathering(sp.record(fa="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="bts")).run(sender=bmc)

    # counter.handle_btp_message(sp.record(_from="tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="bts", sn=sp.nat(4),
    #                                      msg=sp.bytes("0x0507070a000000030dae110000") )).run(sender=admin)


sp.add_compilation_target("bts_periphery", BTPPreiphery(bmc_address=sp.address("KT1BhqGDND5JC8wSrPSR7hA8LtvhaesqCvAq"),
                                                        bts_core_address=sp.address("KT1Sf2Hrs8hTuhwwywExNTrBt2dND3YGDprR"),
                                                        helper_contract=sp.address("KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr"),
                                                        parse_address=sp.address("KT1EKPrSLWjWViZQogFgbc1QmztkR5UGXEWa")))