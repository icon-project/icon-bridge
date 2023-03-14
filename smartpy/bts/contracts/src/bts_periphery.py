import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
strings = sp.io.import_script_from_url("file:./contracts/src/String.py")
rlp_encode = sp.io.import_script_from_url("file:./contracts/src/RLP_encode_struct.py")
rlp_decode = sp.io.import_script_from_url("file:./contracts/src/RLP_decode_struct.py")


class BTPPreiphery(sp.Contract, rlp_decode.DecodeLibrary, rlp_encode.EncodeLibrary):
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
            parse_contract=parse_address,
            mint_status=sp.none
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

    # private function to return the tx status made from handle_btp_message
    def _add_to_blacklist(self, params):
        """
        :param params: List of addresses to be blacklisted
        :return:
        """
        sp.set_type(params, sp.TMap(sp.TNat, sp.TString))

        add_blacklist_status = sp.local("add_blacklist_status", "")
        with sp.if_(sp.len(params) <= self.MAX_BATCH_SIZE):
            sp.for i in sp.range(sp.nat(0), sp.len(params)):
                parsed_addr = sp.view("str_to_addr", self.data.parse_contract, params.get(i), t=sp.TAddress).open_some()
                with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                    self.data.blacklist[parsed_addr] = True
                    add_blacklist_status.value = "success"
                with sp.else_():
                    add_blacklist_status.value = "InvalidAddress"
        with sp.else_():
            add_blacklist_status.value = "error"
        return add_blacklist_status.value

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

    # private function to return the tx status made from handle_btp_message
    def _remove_from_blacklist(self, params):
        """
        :param params: list of address strings
        :return:
        """
        sp.set_type(params, sp.TMap(sp.TNat, sp.TString))

        remove_blacklist_status = sp.local("remove_blacklist_status", "")
        with sp.if_(sp.len(params) <= self.MAX_BATCH_SIZE):
            sp.for i in sp.range(sp.nat(0), sp.len(params)):
                parsed_addr = sp.view("str_to_addr", self.data.parse_contract, params.get(i), t=sp.TAddress).open_some()
                with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                    with sp.if_(self.data.blacklist.contains(parsed_addr)):
                        del self.data.blacklist[parsed_addr]
                        remove_blacklist_status.value = "success"
                    with sp.else_():
                        remove_blacklist_status.value = "UserNotBlacklisted"
                with sp.else_():
                    remove_blacklist_status.value = "InvalidAddress"
        with sp.else_():
            remove_blacklist_status.value = "error"
        return remove_blacklist_status.value

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

    # private function to return the tx status made from handle_btp_message
    def _set_token_limit(self, coin_names, token_limit):
        """
        :param coin_names: list of coin names
        :param token_limit: list of token limits
        :return:
        """
        sp.set_type(coin_names, sp.TMap(sp.TNat, sp.TString))
        sp.set_type(token_limit, sp.TMap(sp.TNat, sp.TNat))

        set_limit_status = sp.local("set_limit_status", "")
        with sp.if_((sp.len(coin_names) == sp.len(token_limit)) & (sp.len(coin_names) <= self.MAX_BATCH_SIZE)):
            sp.for i in sp.range(0, sp.len(coin_names)):
                self.data.token_limit[coin_names[i]] = token_limit.get(i)
            set_limit_status.value = "success"
        with sp.else_():
            set_limit_status.value = "error"
        return set_limit_status.value

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

        send_message_args_type = sp.TRecord(to=sp.TString, svc=sp.TString, sn=sp.TInt, msg=sp.TBytes)
        send_message_entry_point = sp.contract(send_message_args_type, self.data.bmc, "send_message").open_some()
        send_message_args = sp.record(
            to=to_network, svc=self.service_name, sn=self.data.serial_no,
            msg=self.encode_service_message(sp.compute(sp.record(service_type_value=sp.nat(0),
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
        sp.set_type(callback, sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress, prev=sp.TString, callback_msg=sp.TRecord(
        src=sp.TString, dst=sp.TString, svc=sp.TString, sn=sp.TInt, message=sp.TBytes)
                                                      )))
        sp.set_type(bsh_addr, sp.TAddress)

        check_caller = self.only_bmc()

        callback_string = sp.local("callback_string", "")
        with sp.if_((svc == self.service_name) & (check_caller == "Authorized")):
            err_msg = sp.local("error", "")
            sm = self.decode_service_message(msg)

            service_type_variant_match = sp.local("serviceType_variant", False, t=sp.TBool)
            with sm.serviceType.match_cases() as arg:
                with arg.match("REQUEST_COIN_TRANSFER") as a1:
                    service_type_variant_match.value = True
                    callback_string.value = "success"
                    tc = self.decode_transfer_coin_msg(sm.data)
                    parsed_addr = sp.view("str_to_addr", self.data.parse_contract, tc.to, t=sp.TAddress).open_some()

                    with sp.if_(parsed_addr != sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")):
                        handle_request_call = self._handle_request_service(tc.to, tc.assets)
                        with sp.if_(handle_request_call == "success"):
                            self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), sp.nat(2), _from, sn, "", self.RC_OK)
                            sp.emit(sp.record(from_address=_from, to=parsed_addr, serial_no=self.data.serial_no, assets_details=tc.assets), tag="TransferReceived")
                        with sp.else_():
                            err_msg.value = handle_request_call
                            self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), sp.nat(2), _from, sn, err_msg.value,
                                                       self.RC_ERR)
                    with sp.else_():
                        err_msg.value = "InvalidAddress"
                        self.send_response_message(sp.variant("RESPONSE_HANDLE_SERVICE", 2), sp.nat(2), _from, sn, err_msg.value, self.RC_ERR)

                with arg.match("BLACKLIST_MESSAGE") as a2:
                    service_type_variant_match.value = True
                    callback_string.value = "success"
                    bm = self.decode_blacklist_msg(sm.data)
                    addresses = bm.addrs

                    blacklist_service_called = sp.local("blacklist_service", False, t=sp.TBool)
                    with bm.serviceType.match_cases() as b_agr:
                        with b_agr.match("ADD_TO_BLACKLIST") as b_val_1:
                            blacklist_service_called.value = True

                            add_blacklist_call = self._add_to_blacklist(addresses)
                            with sp.if_(add_blacklist_call == "success"):
                                self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from, sn, "AddedToBlacklist", self.RC_OK)
                            with sp.else_():
                                self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from, sn, "ErrorAddToBlackList", self.RC_ERR)

                        with b_agr.match("REMOVE_FROM_BLACKLIST") as b_val_2:
                            blacklist_service_called.value = True

                            remove_blacklist_call = self._remove_from_blacklist(addresses)
                            with sp.if_(remove_blacklist_call == "success"):
                                self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from, sn, "RemovedFromBlacklist", self.RC_OK)
                            with sp.else_():
                                self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from, sn, "ErrorRemoveFromBlackList", self.RC_ERR)

                    sp.if blacklist_service_called.value == False:
                        self.send_response_message(sp.variant("BLACKLIST_MESSAGE", 3), sp.nat(3), _from, sn, "BlacklistServiceTypeErr", self.RC_ERR)

                with arg.match("CHANGE_TOKEN_LIMIT") as a3:
                    service_type_variant_match.value = True
                    callback_string.value = "success"
                    tl = self.decode_token_limit_msg(sm.data)
                    coin_names = tl.coin_name
                    token_limits = tl.token_limit

                    set_limit_call = self._set_token_limit(coin_names, token_limits)
                    with sp.if_(set_limit_call == "success"):
                        self.send_response_message(sp.variant("CHANGE_TOKEN_LIMIT", 4), sp.nat(4), _from, sn, "ChangeTokenLimit", self.RC_OK)
                    with sp.else_():
                        self.send_response_message(sp.variant("CHANGE_TOKEN_LIMIT", 4), sp.nat(4), _from, sn, "ErrorChangeTokenLimit", self.RC_ERR)

                with arg.match("RESPONSE_HANDLE_SERVICE") as a4:
                    service_type_variant_match.value = True
                    with sp.if_(sp.len(sp.pack(self.data.requests.get(sn).from_)) != 0):
                        response = self.decode_response(sm.data)
                        handle_response = self.handle_response_service(sn, response.code, response.message)
                        with sp.if_(handle_response == "success"):
                            callback_string.value = "success"
                        with sp.else_():
                            callback_string.value = "fail"
                    with sp.else_():
                        callback_string.value = "InvalidSN"
                with arg.match("UNKNOWN_TYPE") as a5:
                    service_type_variant_match.value = True
                    callback_string.value = "success"
                    sp.emit(sp.record(_from=_from, sn=sn), tag= "UnknownResponse")

            sp.if service_type_variant_match.value == False:
                callback_string.value = "success"
                self.send_response_message(sp.variant("UNKNOWN_TYPE", 5), sp.nat(5), _from, sn, "Unknown",self.RC_ERR)
        with sp.else_():
            callback_string.value = "fail"

        return_value = sp.record(string=sp.some(callback_string.value), bsh_addr=bsh_addr, prev=prev,
                                 callback_msg=callback_msg)
        sp.transfer(return_value, sp.tez(0), callback)

    @sp.entry_point
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
        sp.set_type(callback, sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress,svc=sp.TString, sn=sp.TInt, code=sp.TNat, msg=sp.TString)))
        sp.set_type(bsh_addr, sp.TAddress)

        check_caller = self.only_bmc()
        handle_btp_error_status = sp.local("handle_btp_error_statue", "")
        with sp.if_((svc == self.service_name) & (check_caller == "Authorized") & (sp.len(sp.pack(self.data.requests.get(sn).from_)) != 0)):
            emit_msg= sp.concat(["errCode: ", sp.view("string_of_int", self.data.parse_contract, sp.to_int(code), t=sp.TString).open_some(),", errMsg: ", msg])
            handle_response_serv = self.handle_response_service(sn, self.RC_ERR, emit_msg)
            with sp.if_(handle_response_serv == "success"):
                handle_btp_error_status.value = "success"
        with sp.else_():
            handle_btp_error_status.value = "fail"

        return_value = sp.record(string=sp.some(handle_btp_error_status.value), bsh_addr=bsh_addr, svc=svc, sn=sn, code=code, msg=msg)
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

        caller = sp.local("caller", sp.view("str_to_addr", self.data.parse_contract, self.data.requests.get(sn).from_, t=sp.TAddress).open_some()
                          , sp.TAddress).value
        loop = sp.local("loop", sp.len(self.data.requests.get(sn).coin_names), sp.TNat).value
        response_call_status = sp.local("response_call_status", "")
        with sp.if_(loop <= self.MAX_BATCH_SIZE):
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
            response_call_status.value = "success"
        with sp.else_():
            response_call_status.value = "BatchMaxSizeExceed"
        return response_call_status.value

    @sp.entry_point
    def callback_mint(self, string):
        sp.set_type(string, sp.TOption(sp.TString))

        sp.verify(sp.sender == self.data.bts_core, "Unauthorized")
        self.data.mint_status = string

        sp.verify(self.data.mint_status.open_some() == "success", "TransferFailed")
        self.data.mint_status = sp.none

    def _handle_request_service(self, to, assets):
        """
        Handle a list of minting/transferring coins/tokens
        :param to: An address to receive coins/tokens
        :param assets:  A list of requested coin respectively with an amount
        :return:
        """
        sp.set_type(to, sp.TString)
        sp.set_type(assets, sp.TMap(sp.TNat, types.Types.Asset))

        status = sp.local("status", "error")
        with sp.if_(sp.len(assets) <= self.MAX_BATCH_SIZE):
            parsed_to = sp.view("str_to_addr", self.data.parse_contract, to, t=sp.TAddress).open_some()
            sp.for i in sp.range(0, sp.len(assets)):
                valid_coin = sp.view("is_valid_coin", self.data.bts_core, assets[i].coin_name, t=sp.TBool).open_some()

                with sp.if_(valid_coin == True):
                    check_transfer = sp.view("check_transfer_restrictions", sp.self_address, sp.record(
                        coin_name=assets[i].coin_name, user=parsed_to, value=assets[i].value), t=sp.TBool).open_some()
                    with sp.if_(check_transfer == True):
                        # inter score call
                        mint_args_type = sp.TRecord(callback=sp.TContract(sp.TOption(sp.TString)),
                                                    to=sp.TAddress, coin_name=sp.TString, value=sp.TNat
                                                    )
                        mint_args_type_entry_point = sp.contract(mint_args_type, self.data.bts_core, "mint").open_some()
                        mint_args = sp.record(callback=sp.self_entry_point("callback_mint"),
                                              to=parsed_to, coin_name=assets[i].coin_name, value=assets[i].value
                                              )
                        sp.transfer(mint_args, sp.tez(0), mint_args_type_entry_point)
                        status.value = "success"
                    with sp.else_():
                        status.value = "FailCheckTransfer"
                with sp.else_():
                    status.value = "UnregisteredCoin"
        with sp.else_():
            status.value = "BatchMaxSizeExceed"

        return status.value

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

            # inter score call
            mint_args_type = sp.TRecord(to=sp.TAddress, coin_name=sp.TString, value=sp.TNat
            )
            mint_args_type_entry_point = sp.contract(mint_args_type, self.data.bts_core, "mint").open_some()
            mint_args = sp.record(
                to=parsed_to, coin_name=assets[i].coin_name, value=assets[i].value
            )
            sp.transfer(mint_args, sp.tez(0), mint_args_type_entry_point)


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
                                      msg=self.encode_service_message(sp.record(service_type_value=service_type_val, data=self.encode_response(sp.record(code=code, message=msg))))
                                      )
        sp.transfer(send_message_args, sp.tez(0), send_message_entry_point)


    @sp.entry_point
    def handle_fee_gathering(self, fa, svc, callback, bsh_addr):
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
        sp.set_type(callback, sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), bsh_addr=sp.TAddress)))
        sp.set_type(bsh_addr, sp.TAddress)

        check_caller = self.only_bmc()
        with sp.if_((svc == self.service_name) & (check_caller == "Authorized")):
            strings.split_btp_address(fa)

            # call transfer_fees of BTS_Core
            transfer_fees_args_type = sp.TString
            transfer_fees_entry_point = sp.contract(transfer_fees_args_type, self.data.bts_core, "transfer_fees").open_some()
            sp.transfer(fa, sp.tez(0), transfer_fees_entry_point)

            sp.transfer(sp.record(string=sp.some("success"), bsh_addr=bsh_addr), sp.tez(0), callback)
        with sp.else_():
            sp.transfer(sp.record(string=sp.some("fail"), bsh_addr=bsh_addr), sp.tez(0), callback)

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



sp.add_compilation_target("bts_periphery", BTPPreiphery(bmc_address=sp.address("KT1UrLqhQHDC3mJw9BUrqsiix7JRbxTsvWJu"),
                                                        bts_core_address=sp.address("KT1JAippuMfS6Bso8DGmigmTdkgEZUxQxYyX"),
                                                        helper_contract=sp.address("KT1HwFJmndBWRn3CLbvhUjdupfEomdykL5a6"),
                                                        parse_address=sp.address("KT1EKPrSLWjWViZQogFgbc1QmztkR5UGXEWa"),
                                                        native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
                                                        owner_address = sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"))    )