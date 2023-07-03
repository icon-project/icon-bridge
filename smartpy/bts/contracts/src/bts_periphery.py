import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")
rlp = sp.io.import_script_from_url("file:./contracts/src/RLP_struct.py")
t_balance_of_request = sp.TRecord(owner=sp.TAddress, token_id=sp.TNat).layout(
    ("owner", "token_id")
)

t_balance_of_response = sp.TRecord(
    request=t_balance_of_request, balance=sp.TNat
).layout(("request", "balance"))

class BTSPeriphery(sp.Contract, rlp.DecodeEncodeLibrary):
    service_name = sp.string("bts")

    RC_OK = sp.nat(0)
    RC_ERR = sp.nat(1)
    UINT_CAP = sp.nat(115792089237316195423570985008687907853269984665640564039457584007913129639935)

    MAX_BATCH_SIZE = sp.nat(15)

    def __init__(self, bmc_address, bts_core_address, helper_contract, parse_address, native_coin_name, owner_address):
        self.update_initial_storage(
            bmc=bmc_address,
            owner=owner_address,
            bts_core=bts_core_address,
            blacklist=sp.map(tkey=sp.TAddress, tvalue=sp.TBool),
            token_limit=sp.map({native_coin_name : self.UINT_CAP}, tkey=sp.TString, tvalue=sp.TNat),
            requests=sp.big_map(tkey=sp.TInt, tvalue=types.Types.PendingTransferCoin),
            serial_no = sp.int(0),
            number_of_pending_requests = sp.nat(0),
            helper=helper_contract,
            parse_contract=parse_address
        )

    def only_bmc(self):
        check_access = sp.local("check_access", "Unauthorized")
        with sp.if_(sp.sender == self.data.bmc):
            check_access.value = "Authorized"
        return check_access.value

    def only_owner(self):
        sp.verify(sp.sender == self.data.owner, "Unauthorized")

    def only_bts_core(self):
        sp.verify(sp.sender == self.data.bts_core, "Unauthorized")

    @sp.entry_point
    def set_helper_address(self, address):
        sp.set_type(address, sp.TAddress)
        self.only_owner()
        self.data.helper = address

    @sp.entry_point
    def set_parse_address(self, address):
        sp.set_type(address, sp.TAddress)
        self.only_owner()
        self.data.parse_contract = address

    @sp.entry_point
    def set_bmc_address(self, params):
        sp.set_type(params, sp.TAddress)
        self.only_owner()
        self.data.bmc = params

    @sp.entry_point
    def set_bts_core_address(self, params):
        sp.set_type(params, sp.TAddress)
        self.only_owner()
        self.data.bts_core = params

    @sp.onchain_view()
    def has_pending_request(self):
        """
        :return: boolean
        """
        sp.result(self.data.number_of_pending_requests != sp.nat(0))

    # private function for adding blacklist address
    def _add_to_blacklist(self, params):
        """
        :param params: List of addresses to be blacklisted
        :return:
        """
        sp.set_type(params, sp.TList(sp.TString))

        add_blacklist_status = sp.local("add_blacklist_status", "success")
        with sp.if_(sp.len(params) <= self.MAX_BATCH_SIZE):
            sp.for item in params:
                parsed_addr = sp.view("str_to_addr", self.data.parse_contract, item, t=sp.TAddress).open_some()
                with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                    self.data.blacklist[parsed_addr] = True
                with sp.else_():
                    add_blacklist_status.value = "InvalidAddress"
        with sp.else_():
            add_blacklist_status.value = "error"
        return add_blacklist_status.value

    # private function for removing blacklist address
    def _remove_from_blacklist(self, params):
        """
        :param params: list of address strings
        :return:
        """
        sp.set_type(params, sp.TList(sp.TString))

        remove_blacklist_status = sp.local("remove_blacklist_status", "success")
        with sp.if_(sp.len(params) <= self.MAX_BATCH_SIZE):
            sp.for item in params:
                parsed_addr = sp.view("str_to_addr", self.data.parse_contract, item, t=sp.TAddress).open_some()
                with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                    with sp.if_(self.data.blacklist.contains(parsed_addr)):
                        del self.data.blacklist[parsed_addr]
                    with sp.else_():
                        remove_blacklist_status.value = "UserNotBlacklisted"
                with sp.else_():
                    remove_blacklist_status.value = "InvalidAddress"
        with sp.else_():
            remove_blacklist_status.value = "error"
        return remove_blacklist_status.value

    @sp.entry_point
    def set_token_limit(self, coin_names_limit):
        """
        :param coin_names_limit: map of coin names and its limits
        :return:
        """
        sp.set_type(coin_names_limit, sp.TMap(sp.TString, sp.TNat))
        # sp.set_type(token_limit, sp.TMap(sp.TNat, sp.TNat))

        sp.verify((sp.sender == sp.self_address )| (sp.sender == self.data.bts_core), "Unauthorized")
        # sp.verify(sp.len(coin_names_limits) == sp.len(token_limit), "InvalidParams")
        sp.verify(sp.len(coin_names_limit) <= self.MAX_BATCH_SIZE, "BatchMaxSizeExceed")

        coin_names_limit_items = coin_names_limit.items()
        sp.for item in coin_names_limit_items:
            self.data.token_limit[item.key] = item.value

    # private function to return the tx status made from handle_btp_message
    def _set_token_limit(self, coin_name_limit):
        """
        :param coin_names: list of coin names
        :param token_limit: list of token limits
        :return:
        """
        sp.set_type(coin_name_limit, sp.TMap(sp.TString, sp.TNat))

        set_limit_status = sp.local("set_limit_status", "success")
        with sp.if_(sp.len(coin_name_limit) <= self.MAX_BATCH_SIZE):
            sp.for item in coin_name_limit.items():
                self.data.token_limit[item.key] = item.value
        with sp.else_():
            set_limit_status.value = "error"
        return set_limit_status.value

    @sp.entry_point(lazify=False)
    def update_send_service_message(self, ep):
        self.only_owner()
        sp.set_entry_point("send_service_message", ep)

    @sp.entry_point(lazify=True)
    def send_service_message(self, _from, to, coin_details):
        """
        Send service message to BMC
        :param _from: from address
        :param to: to address
        :param coin_details:
        :return:
        """

        sp.set_type(_from, sp.TAddress)
        sp.set_type(to, sp.TString)
        sp.set_type(coin_details, sp.TList(sp.TRecord(coin_name=sp.TString, value=sp.TNat, fee=sp.TNat)))

        self.only_bts_core()

        to_network, to_address = sp.match_pair(strings.split_btp_address(to))

        assets_details = sp.compute(sp.map(tkey=sp.TNat, tvalue=types.Types.AssetTransferDetail))

        i=sp.local("i_", sp.nat(0))
        sp.for item in coin_details:
            assets_details[i.value]= item
            i.value += sp.nat(1)

        self.data.serial_no += 1

        start_from = sp.view("add_to_str", self.data.parse_contract, _from, t=sp.TString).open_some()

        send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
        send_message_entry_point = sp.contract(send_message_args_type, self.data.bmc, "send_message").open_some()
        send_message_args = sp.record(
            to=to_network, svc=self.service_name, sn=self.data.serial_no,
            msg = self.encode_service_message(sp.record(service_type_value=sp.nat(0),
            data=self.encode_transfer_coin_msg(sp.record(from_addr=start_from, to=to_address,
                                                         assets=coin_details))
                                                        )
                                              )
        )

        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)

        # push pending tx into record list
        self.data.requests[self.data.serial_no] = sp.record(
            from_=start_from, to=to, coin_details=coin_details)
        self.data.number_of_pending_requests += sp.nat(1)
        sp.emit(sp.record(from_address=_from, to=to, serial_no=self.data.serial_no,
                          assets_details=assets_details), tag="TransferStart")

    @sp.entry_point(lazify=False)
    def update_handle_btp_message(self, ep):
        self.only_owner()
        sp.set_entry_point("handle_btp_message", ep)

    @sp.entry_point(lazify=True)
    def handle_btp_message(self, _from, svc, sn, msg, callback, bsh_addr, prev, callback_msg):
        """
        BSH handle BTP message from BMC contract
        :param _from: An originated network address of a request
        :param svc: A service name of BSH contract
        :param sn: A serial number of a service request
        :param msg: An RLP message of a service request/service response
        :param callback: callback function type in bmc_periphery
        :param bsh_addr: param for callback function in bmc_periphery
        :param prev: param for callback function in bmc_periphery
        :param callback_msg: param for callback function in bmc_periphery
        :return:
        """

        sp.set_type(_from, sp.TString)
        sp.set_type(svc, sp.TString)
        sp.set_type(sn, sp.TInt)
        sp.set_type(msg, sp.TBytes)
        sp.set_type(callback, sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress,
                                                    prev=sp.TString, callback_msg=sp.TRecord(src=sp.TString,
                                                    dst=sp.TString, svc=sp.TString, sn=sp.TInt, message=sp.TBytes)
                                                      )))
        sp.set_type(bsh_addr, sp.TAddress)

        check_caller = self.only_bmc()

        callback_string = sp.local("callback_string", "success")
        with sp.if_((svc == self.service_name) & (check_caller == "Authorized")):
            decode_call = self.decode_service_message(msg)
            with sp.if_(decode_call.status == "Success"):
                sm = decode_call.rv

                with sm.serviceType.match_cases() as arg:
                    with arg.match("REQUEST_COIN_TRANSFER") as a1:
                        tc_call = self.decode_transfer_coin_msg(sm.data)
                        with sp.if_(tc_call.status == "Success"):
                            tc = tc_call.value
                            parsed_addr = sp.view("str_to_addr", self.data.parse_contract,
                                                  tc.to, t=sp.TAddress).open_some()

                            with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                                handle_request_call= self._handle_request_service(parsed_addr, tc.assets)
                                with sp.if_(handle_request_call == "success"):
                                    self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), sp.nat(2),
                                                               _from, sn, "", self.RC_OK)
                                    sp.emit(sp.record(from_address=_from, to=parsed_addr, serial_no=self.data.serial_no,
                                                      assets_details=tc.assets), tag="TransferReceived")
                                with sp.else_():
                                    self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), sp.nat(2),
                                                               _from, sn, handle_request_call, self.RC_ERR)
                            with sp.else_():
                                self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), sp.nat(2), _from,
                                                           sn, "InvalidAddress", self.RC_ERR)
                        with sp.else_():
                            callback_string.value = "ErrorInDecodingTransferCoin"

                    with arg.match("BLACKLIST_MESSAGE") as a2:
                        bm_call = self.decode_blacklist_msg(sm.data)
                        with sp.if_(bm_call.status == "Success"):
                            bm = bm_call.rv
                            addresses = bm.addrs

                            with bm.serviceType.match_cases() as b_agr:
                                with b_agr.match("ADD_TO_BLACKLIST") as b_val_1:
                                    add_blacklist_call = self._add_to_blacklist(addresses)
                                    with sp.if_(add_blacklist_call == "success"):
                                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from,
                                                                   sn, "AddedToBlacklist", self.RC_OK)
                                    with sp.else_():
                                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from,
                                                                   sn, "ErrorAddToBlackList", self.RC_ERR)

                                with b_agr.match("REMOVE_FROM_BLACKLIST") as b_val_2:
                                    remove_blacklist_call = self._remove_from_blacklist(addresses)
                                    with sp.if_(remove_blacklist_call == "success"):
                                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from,
                                                                   sn, "RemovedFromBlacklist", self.RC_OK)
                                    with sp.else_():
                                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from,
                                                                   sn, "ErrorRemoveFromBlackList", self.RC_ERR)

                                with b_agr.match("ERROR") as b_val_2:
                                    self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from, sn,
                                                               "BlacklistServiceTypeErr", self.RC_ERR)
                        with sp.else_():
                            callback_string.value = "ErrorInDecodingBlacklist"

                    with arg.match("CHANGE_TOKEN_LIMIT") as a3:
                        tl_call = self.decode_token_limit_msg(sm.data)
                        with sp.if_(tl_call.status == "Success"):
                            tl = tl_call.rv
                            coin_name_limit = tl.coin_name_limit

                            set_limit_call = self._set_token_limit(coin_name_limit)
                            with sp.if_(set_limit_call == "success"):
                                self.send_response_message(sp.variant("CHANGE_TOKEN_LIMIT", 4), sp.nat(4), _from, sn,
                                                           "ChangeTokenLimit", self.RC_OK)
                            with sp.else_():
                                self.send_response_message(sp.variant("CHANGE_TOKEN_LIMIT", 4), sp.nat(4), _from, sn,
                                                           "ErrorChangeTokenLimit", self.RC_ERR)
                        with sp.else_():
                            callback_string.value = "ErrorInDecodingTokenLimit"

                    with arg.match("RESPONSE_HANDLE_SERVICE") as a4:
                        with sp.if_(sp.len(sp.pack(self.data.requests.get(sn).from_)) != 0):
                            fn_call = self.decode_response(sm.data)
                            response = fn_call.rv
                            with sp.if_(fn_call.status == "Success"):
                                handle_response = self.handle_response_service(sn, response.code, response.message)
                                with sp.if_(handle_response != "success"):
                                    callback_string.value = "fail"
                            with sp.else_():
                                callback_string.value = "ErrorInDecoding"
                        with sp.else_():
                            callback_string.value = "InvalidSN"

                    with arg.match("UNKNOWN_TYPE") as a5:
                        sp.emit(sp.record(_from=_from, sn=sn), tag= "UnknownResponse")

                    with arg.match("ERROR") as a5:
                        self.send_response_message(sp.variant("UNKNOWN_TYPE", 5), sp.nat(5), _from, sn,
                                                   "Unknown",self.RC_ERR)
            with sp.else_():
                callback_string.value = "ErrorInDecoding"
        with sp.else_():
            callback_string.value = "UnAuthorized"

        return_value = sp.record(string=sp.some(callback_string.value), bsh_addr=bsh_addr, prev=prev,
                                 callback_msg=callback_msg)
        sp.transfer(return_value, sp.tez(0), callback)

    @sp.entry_point(lazify=False)
    def update_handle_btp_error(self, ep):
        self.only_owner()
        sp.set_entry_point("handle_btp_error", ep)

    @sp.entry_point(lazify=True)
    def handle_btp_error(self, svc, sn, code, msg, callback, bsh_addr):
        """
        BSH handle BTP Error from BMC contract
        :param svc: A service name of BSH contract
        :param sn: A serial number of a service request
        :param code: A response code of a message (RC_OK / RC_ERR)
        :param msg: A response message
        :param callback: callback function type in bmc_periphery
        :param bsh_addr: param for callback function in bmc_periphery
        :return:
        """

        sp.set_type(svc, sp.TString)
        sp.set_type(sn, sp.TInt)
        sp.set_type(code, sp.TNat)
        sp.set_type(msg, sp.TString)
        sp.set_type(callback, sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress,
                                                      svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)))
        sp.set_type(bsh_addr, sp.TAddress)

        check_caller = self.only_bmc()
        handle_btp_error_status = sp.local("handle_btp_error_status", "success")
        with sp.if_((svc == self.service_name) & (check_caller == "Authorized") & (sp.len(sp.pack(
                self.data.requests.get(sn).from_)) != 0)):
            emit_msg= sp.concat(["errCode: ", sp.view("string_of_int", self.data.parse_contract, sp.to_int(code),
                                                      t=sp.TString).open_some(),", errMsg: ", msg])
            handle_response_serv = self.handle_response_service(sn, self.RC_ERR, emit_msg)
            with sp.if_(handle_response_serv != "success"):
                handle_btp_error_status.value = "fail"
        with sp.else_():
            handle_btp_error_status.value = "UnAuthorized"

        return_value = sp.record(string=sp.some(handle_btp_error_status.value), bsh_addr=bsh_addr, svc=svc, sn=sn,
                                 code=code, msg=msg)
        sp.transfer(return_value, sp.tez(0), callback)

    def handle_response_service(self, sn, code, msg):
        """
        :param sn:
        :param code:
        :param msg:
        :return:
        """
        sp.set_type(sn, sp.TInt)
        sp.set_type(code, sp.TNat)
        sp.set_type(msg, sp.TString)

        caller = sp.local("caller", sp.view("str_to_addr", self.data.parse_contract,
                                            self.data.requests.get(sn).from_,
                                            t=sp.TAddress).open_some(), sp.TAddress)
        loop = sp.local("loop", sp.len(self.data.requests.get(sn).coin_details), sp.TNat)
        response_call_status = sp.local("response_call_status", "success")
        check_valid = sp.local("check_valid", True)
        bts_core_fa2_balance = sp.local("fa2_token_balance_response_service", sp.nat(0))

        with sp.if_(loop.value <= self.MAX_BATCH_SIZE):
            sp.for item in self.data.requests.get(sn).coin_details:
                bts_core_address = self.data.bts_core
                coin_name = item.coin_name
                value = item.value
                fee = item.fee
                amount = value + fee
                bts_core_balance = sp.view("balance_of", bts_core_address,
                                                  sp.record(owner=caller.value, coin_name=coin_name), t=
                                                  sp.TRecord(usable_balance=sp.TNat, locked_balance=sp.TNat,
                                                            refundable_balance=sp.TNat, user_balance=sp.TNat)
                                                  ).open_some()
                # check if caller has enough locked in bts_core
                with sp.if_(bts_core_balance.locked_balance < amount):
                    check_valid.value = False
                coin_type = sp.view("coin_type", bts_core_address, coin_name, t=sp.TNat).open_some()
                with sp.if_(coin_type == sp.nat(1)):
                    coin_address = sp.view("coin_id", bts_core_address, coin_name, t=sp.TAddress).open_some()
                    bts_core_fa2 = sp.view("get_balance_of", coin_address,
                                                   [sp.record(owner=bts_core_address, token_id=sp.nat(0))],
                                                   t=sp.TList(t_balance_of_response)).open_some("Invalid view")
                    sp.for elem in bts_core_fa2:
                        bts_core_fa2_balance.value = elem.balance
                    # check if bts_core has enough NATIVE_WRAPPED_COIN_TYPE to burn
                    with sp.if_(bts_core_fa2_balance.value < value):
                        check_valid.value = False

            with sp.if_(check_valid.value == True):
                sp.for _item in self.data.requests.get(sn).coin_details:
                    # inter score call
                    handle_response_service_args_type = sp.TRecord(
                        requester=sp.TAddress, coin_name=sp.TString, value=sp.TNat, fee=sp.TNat, rsp_code=sp.TNat)
                    handle_response_service_entry_point = sp.contract(handle_response_service_args_type,
                                            self.data.bts_core, "handle_response_service").open_some("invalid call")
                    handle_response_service_args = sp.record(
                        requester=caller.value, coin_name=_item.coin_name,
                        value=_item.value,
                        fee=_item.fee, rsp_code=code
                    )
                    sp.transfer(handle_response_service_args, sp.tez(0), handle_response_service_entry_point)

                del self.data.requests[sn]
                self.data.number_of_pending_requests = sp.as_nat(self.data.number_of_pending_requests-1)

                sp.emit(sp.record(caller=caller.value, sn=sn, code=code, msg=msg), tag="TransferEnd")
            with sp.else_():
                response_call_status.value = "Error in bts handle_response_service"
        with sp.else_():
            response_call_status.value = "BatchMaxSizeExceed"

        return response_call_status.value

    def _handle_request_service(self, to, assets):
        """
        Handle a list of minting/transferring coins/tokens
        :param to: An address to receive coins/tokens
        :param assets:  A list of requested coin respectively with an amount
        :return:
        """
        sp.set_type(to, sp.TAddress)
        sp.set_type(assets, sp.TMap(sp.TNat, types.Types.Asset))

        status = sp.local("status", "success")
        check_validity = sp.local("check_validity", True)
        bts_core_fa2_balance = sp.local("fa2_token_balance", sp.nat(0))
        bts_core_address = self.data.bts_core
        with sp.if_(sp.len(assets) <= self.MAX_BATCH_SIZE):
            parsed_to = to
            sp.for i in sp.range(0, sp.len(assets)):
                coin_name = assets[i].coin_name
                transferred_amount = assets[i].value
                valid_coin = sp.view("is_valid_coin", bts_core_address, coin_name, t=sp.TBool).open_some()
                check_transfer = sp.view("check_transfer_restrictions", sp.self_address, sp.record(
                    coin_name=coin_name, user=parsed_to, value=transferred_amount), t=sp.TBool).open_some()

                native_coin_name, bts_core_balance = sp.match_pair(sp.view("native_coin_balance_of", bts_core_address,
                                        sp.unit, t=sp.TPair(sp.TString, sp.TMutez)).open_some("Invalid view"))
                with sp.if_(native_coin_name == coin_name):
                    with sp.if_(sp.utils.mutez_to_nat(bts_core_balance) < transferred_amount):
                        check_validity.value = False
                with sp.else_():
                    coin_type = sp.view("coin_type", bts_core_address, coin_name, t=sp.TNat).open_some()
                    with sp.if_((valid_coin == True) & (coin_type == sp.nat(2))):
                        coin_address = sp.view("coin_id", bts_core_address, coin_name, t=sp.TAddress).open_some()
                        bts_core_fa2 = sp.view("get_balance_of", coin_address,
                                               [sp.record(owner=bts_core_address, token_id=sp.nat(0))],
                                               t=sp.TList(t_balance_of_response)).open_some("Invalid view")
                        sp.for elem in bts_core_fa2:
                            bts_core_fa2_balance.value = elem.balance
                        with sp.if_(bts_core_fa2_balance.value < transferred_amount):
                            check_validity.value = False

                with sp.if_((check_transfer == False) | (valid_coin == False)) :
                    check_validity.value = False
            with sp.if_(check_validity.value == True):
                sp.for i in sp.range(0, sp.len(assets)):
                    # inter score call
                    mint_args_type = sp.TRecord(to=sp.TAddress, coin_name=sp.TString, value=sp.TNat)
                    mint_args_type_entry_point = sp.contract(mint_args_type, bts_core_address, "mint").open_some()
                    mint_args = sp.record(to=parsed_to, coin_name=assets[i].coin_name, value=assets[i].value)
                    sp.transfer(mint_args, sp.tez(0), mint_args_type_entry_point)
            with sp.else_():
                status.value = "UnregisteredCoin"
        with sp.else_():
            status.value = "BatchMaxSizeExceed"

        return status.value

    def send_response_message(self, service_type, service_type_val, to, sn, msg, code):
        """

        :param service_type:
        :param service_type_val: value of service_type variant
        :param to:
        :param sn:
        :param msg:
        :param code:
        :return:
        """
        sp.set_type(service_type, types.Types.ServiceType)
        sp.set_type(service_type_val, sp.TNat)
        sp.set_type(to, sp.TString)
        sp.set_type(sn, sp.TInt)
        sp.set_type(msg, sp.TString)
        sp.set_type(code, sp.TNat)

        send_message_args_type = sp.TRecord(
            to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes
        )
        send_message_entry_point = sp.contract(send_message_args_type, self.data.bmc, "send_message").open_some()
        send_message_args = sp.record(to=to, svc=self.service_name, sn=sn,
                                      msg=self.encode_service_message(sp.record(service_type_value=service_type_val,
                                                                                data=self.encode_response(
                                        sp.record(code=code, message=msg)))))
        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)


    @sp.entry_point
    def handle_fee_gathering(self, fa, svc, bsh_addr):
        """
        BSH handle Gather Fee Message request from BMC contract
        :param fa: A BTP address of fee aggregator
        :param svc: A name of the service
        :param callback: callback function type in bmc_periphery
        :param bsh_addr: address of bts_periphery
        :return:
        """
        sp.set_type(fa, sp.TString)
        sp.set_type(svc, sp.TString)
        sp.set_type(bsh_addr, sp.TAddress)

        check_caller = self.only_bmc()
        with sp.if_((svc == self.service_name) & (check_caller == "Authorized")):
            strings.split_btp_address(fa)

            # call transfer_fees of BTS_Core
            transfer_fees_args_type = sp.TString
            transfer_fees_entry_point = sp.contract(transfer_fees_args_type, self.data.bts_core,
                                                    "transfer_fees").open_some()
            sp.transfer(fa, sp.tez(0), transfer_fees_entry_point)

    @sp.onchain_view()
    def check_transfer_restrictions(self, params):
        """

        :param params: Record of coin transfer details
        :return:
        """
        sp.set_type(params, sp.TRecord(coin_name=sp.TString, user=sp.TAddress, value=sp.TNat))

        with sp.if_((self.data.blacklist.contains(params.user) == False) &
                    (self.data.token_limit.get(params.coin_name, default_value=sp.nat(0)) >= params.value)):
            sp.result(True)
        with sp.else_():
            sp.result(False)


sp.add_compilation_target("bts_periphery", BTSPeriphery(bmc_address=sp.address("KT1VFtWq2dZDH1rTfLtgMaASMt4UX78omMs2"),
                                                        bts_core_address=sp.address("KT1JAippuMfS6Bso8DGmigmTdkgEZUxQxYyX"),
                                                        helper_contract=sp.address("KT1HwFJmndBWRn3CLbvhUjdupfEomdykL5a6"),
                                                        parse_address=sp.address("KT1Ha8LzZa7ku1F8eytY7hgNKFJ2BKFRqSDh"),
                                                        native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
                                                        owner_address = sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")
                                                        ))