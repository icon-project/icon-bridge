import smartpy as sp

Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class DecodeEncodeLibrary:

    def decode_response(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        code = sp.local("code_bts_decode_response", 0)
        is_error = sp.local("error_in_bts_decode_response", sp.string("Success"))
        rlp_message_list = sp.local("rlp_message_list_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_message_list.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInBTSDecoding"
        rlp_ = rlp_message_list.value
        counter = sp.local("counter_response_bts", 0)
        msg = sp.local("message_in_bts", sp.string(""))
        with sp.if_(is_error.value == "Success"):
            sp.for i in rlp_.items():
                sp.if counter.value == 0:
                    code.value = Utils2.Int.of_bytes(i.value)
                sp.if counter.value == 1:
                    msg.value = sp.view("decode_string", self.data.helper, i.value, t=sp.TString).open_some()
                counter.value = counter.value + 1
        return sp.record(rv = sp.record(code=code.value, message = msg.value),status = is_error.value)

    def decode_service_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_service_message = sp.local("rlp_decode_service_message_bts", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bts_decode_service_message_bts", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_service_message.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInBTSDecoding"
        _service_type = sp.local("_service_type_decode_service_message_bts", sp.variant("ERROR", sp.nat(10)))
        data = sp.local("data_decode_service_message_bts", sp.bytes("0x"))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_service_message.value
            temp_int = sp.local("temp_int_decode_service_message_bts", 0)
            counter = sp.local("counter_decode_service_message_bts", 0)
            sp.for i in rlp_.items():
                sp.if counter.value == 0:
                    temp_int.value = Utils2.Int.of_bytes(i.value)
                sp.if counter.value == 1:
                    data.value = sp.view("without_length_prefix", self.data.helper, i.value,
                                             t=sp.TBytes).open_some()
                counter.value = counter.value + 1

            with sp.if_ (temp_int.value == 0):
                _service_type.value = sp.variant("REQUEST_COIN_TRANSFER", temp_int.value)
            with sp.if_ (temp_int.value == 1):
                _service_type.value = sp.variant("REQUEST_COIN_REGISTER", temp_int.value)
            with sp.if_ (temp_int.value == 2):
                _service_type.value = sp.variant("RESPONSE_HANDLE_SERVICE", temp_int.value)
            with sp.if_ (temp_int.value == 3):
                _service_type.value = sp.variant("BLACKLIST_MESSAGE", temp_int.value)
            with sp.if_ (temp_int.value == 4):
                    _service_type.value = sp.variant("CHANGE_TOKEN_LIMIT", temp_int.value)
            with sp.if_ (temp_int.value == 5):
                _service_type.value = sp.variant("UNKNOWN_TYPE", temp_int.value)

        return sp.record(rv=sp.record(serviceType=_service_type.value,
                         data=data.value), status=is_error.value)

    def decode_transfer_coin_msg(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        is_error = sp.local("error_in_bts_decode_transfer_coin_msg", sp.string("Success"))
        rlp_transfer_coin = sp.local("rlp_transfer_coin_msg_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_transfer_coin.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInBTSDecoding"
        rlp_ = rlp_transfer_coin.value

        temp_byt = sp.local("byt_transfer_coin_msg_bts", sp.bytes("0x"))
        counter = sp.local("counter_transfer_coin_msg_bts", sp.nat(0))
        from_address = sp.local("from_address_transfer_coin_msg_bts", sp.string(""))
        to_address = sp.local("to_address_transfer_coin_msg_bts", sp.string(""))
        rv_assets = sp.local("assets", {}, sp.TMap(sp.TNat, types.Types.Asset))
        with sp.if_(is_error.value == "Success"):
            sp.for i in rlp_.items():
                with sp.if_ (counter.value == 2):
                    temp_byt.value = i.value
                with sp.if_ (counter.value == 0):
                    from_address.value = sp.view("decode_string", self.data.helper, i.value, t=sp.TString).open_some()
                with sp.if_ (counter.value == 1):
                    to_address.value = sp.view("decode_string", self.data.helper, i.value, t=sp.TString).open_some()
                counter.value = counter.value + 1
            sub_list = sp.local("sub_list_transfer_coin_msg_bts", temp_byt.value)
            new_sub_list_tcm = sp.local("nsl_transfer_coin_msg_bts", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                new_sub_list_tcm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInBTSDecoding"
            with sp.if_(is_error.value == "Success"):
                new_sub_list = new_sub_list_tcm.value
                counter.value = 0
                new_temp_byt = sp.local("new_temp_byt_transfer_coin_msg_bts", sp.bytes("0x"))
                nsl3_tcm = sp.local("nsl3_transfer_coin_msg_bts", sp.map(tkey=sp.TNat))
                view_value = sp.local("view_value_transfer_coin_msg_bts", sp.map(tkey=sp.TNat))
                counter_nested = sp.local("counter_nested_transfer_coin_msg_bts", sp.nat(0), t=sp.TNat)
                temp_byt = sp.local("temp_byte2_transfer_coin_msg_bts", sp.bytes("0x"))
                temp_byt_nested = sp.local("temp_byte_nested_transfer_coin_msg_bts", sp.bytes("0x"))
                sp.for x in new_sub_list.items():
                    new_temp_byt.value = x.value
                    is_list_lambda = sp.view("is_list", self.data.helper, new_temp_byt.value,
                                             t=sp.TBool).open_some()
                    with sp.if_(is_list_lambda):
                        nsl3_tcm.value = sp.view("decode_list", self.data.helper, new_temp_byt.value,
                                                t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                    with sp.else_():
                        is_error.value = "ErrorInBTSDecoding"
                    with sp.if_(is_error.value == "Success"):
                        view_value.value = nsl3_tcm.value
                        counter_nested.value = sp.nat(0)
                        sp.for i in view_value.value.items():
                            with sp.if_ (counter_nested.value == 1):
                                temp_byt_nested.value = sp.view("without_length_prefix", self.data.helper, i.value,
                                                                t=sp.TBytes).open_some()
                            with sp.if_ (counter_nested.value == 0):
                                temp_byt.value = i.value
                            counter_nested.value += 1

                        rv_assets.value[counter.value] = sp.record(coin_name=sp.view("decode_string", self.data.helper, temp_byt.value, t=sp.TString).open_some()
                                                                   , value=Utils2.Int.of_bytes(temp_byt_nested.value))

                        counter.value = counter.value + 1

        return sp.record(value=sp.record(from_addr= from_address.value, to = to_address.value, assets = rv_assets.value), status = is_error.value)

    def decode_blacklist_msg(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_blacklist_msg = sp.local("rlp_blacklist_msg_bts", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bts_decode_blacklist_message", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        _service_type = sp.local("_service_type_blacklist_msg_bts", sp.variant("ERROR", sp.nat(10)))
        rv_blacklist_address = sp.local("blacklist_data", [], sp.TList(sp.TString))
        with sp.if_(is_list_lambda):
            rlp_blacklist_msg.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInBTSDecoding"
        rlp_ = rlp_blacklist_msg.value

        temp_byt = sp.local("byt_transfer", sp.bytes("0x"))
        rv1_byt = sp.local("rv1_byt", sp.bytes("0x"))
        rv2_byt = sp.local("rv2_byt", sp.bytes("0x"))
        counter = sp.local("counter_blacklist", sp.nat(0))
        with sp.if_(is_error.value == "Success"):
            sp.for i in rlp_.items():
                sp.if counter.value == 2:
                    rv2_byt.value = i.value
                    rv2 = sp.view("decode_string", self.data.helper, rv2_byt.value, t=sp.TString).open_some()
                sp.if counter.value == 0:
                    rv1_byt.value = sp.view("without_length_prefix", self.data.helper, i.value,
                                            t=sp.TBytes).open_some()
                    rv1 = Utils2.Int.of_bytes(rv1_byt.value)
                    with sp.if_(rv1 == 0):
                        _service_type.value = sp.variant("ADD_TO_BLACKLIST", rv1)
                    with sp.else_():
                        _service_type.value = sp.variant("REMOVE_FROM_BLACKLIST", rv1)
                sp.if counter.value == 1:
                    temp_byt.value = i.value
                counter.value = counter.value + 1
            sub_list = sp.local("sub_list", temp_byt.value)
            nsl_bm = sp.local("nsl_bts_bm", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl_bm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                         t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInBTSDecoding"
            with sp.if_(is_error.value == "Success"):
                new_sub_list = nsl_bm.value
                counter.value = 0
                addr_string = sp.local("addr_string", "")
                counter.value = 0
                sp.for x in new_sub_list.items():
                    addr_string.value = sp.view("decode_string", self.data.helper, x.value,
                                             t=sp.TString).open_some()
                    rv_blacklist_address.value.push(addr_string.value)
                    counter.value = counter.value + 1

        return sp.record(rv=sp.record(serviceType = _service_type.value , addrs = rv_blacklist_address.value ,
                                      net = rv2), status = is_error.value)

    def decode_token_limit_msg(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_tlm = sp.local("rlp_tlm_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        is_error = sp.local("error_in_bts_decode_token_limit_msg", sp.string("Success"))
        with sp.if_(is_list_lambda):
            rlp_tlm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInBTSDecoding"
        rlp_ = rlp_tlm.value

        temp_byt = sp.local("byt_transfer", sp.bytes("0x"))
        temp_byt1 = sp.local("byt_transfer_temp1", sp.bytes("0x"))
        rv1_byt = sp.local("rv1_byt", sp.bytes("0x"))
        counter = sp.local("counter_token_limit", sp.nat(0))
        net = sp.local("network", sp.string(""))
        rv_names = sp.local("names", {}, sp.TMap(sp.TNat, sp.TString))
        rv_limit = sp.local("limit", {}, sp.TMap(sp.TNat, sp.TNat))
        coin_name_limit = sp.local("coin_name_limit", {}, sp.TMap(sp.TString, sp.TNat))
        with sp.if_(is_error.value == "Success"):
            sp.for i in rlp_.items():
                sp.if counter.value == 0:
                    temp_byt.value = i.value
                sp.if counter.value == 1:
                    temp_byt1.value = i.value
                sp.if counter.value == 2:
                    rv1_byt.value = i.value
                counter.value = counter.value + 1
            sub_list = sp.local("sub_list", temp_byt.value)
            net.value = sp.view("decode_string", self.data.helper, rv1_byt.value, t=sp.TString).open_some()
            sub_list_limit = sp.local("sub_list_limit", temp_byt1.value)
            nsl1_dtlm = sp.local("nsl1_bts_dtlm", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl1_dtlm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                         t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInBTSDecoding"
            new_sub_list = nsl1_dtlm.value
            counter.value = 0
            with sp.if_(is_error.value == "Success"):
                sp.for x in new_sub_list.items():
                    rv_names.value[counter.value] = sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some()
                    counter.value += 1
                nsl_dtlm = sp.local("nsl_bts_dtlm", sp.map(tkey=sp.TNat))
                is_list_lambda = sp.view("is_list", self.data.helper, sub_list_limit.value, t=sp.TBool).open_some()
                counter.value = 0
                with sp.if_(is_list_lambda):
                    nsl_dtlm.value = sp.view("decode_list", self.data.helper, sub_list_limit.value,
                                            t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                with sp.else_():
                    is_error.value = "ErrorInBTSDecoding"
                with sp.if_(is_error.value == "Success"):
                    new_sub_list1 = nsl_dtlm.value
                    limit = sp.local("limit_val", sp.bytes("0x"), t=sp.TBytes)
                    sp.for y in new_sub_list1.items():
                        limit.value = sp.view("without_length_prefix", self.data.helper, y.value,
                                          t=sp.TBytes).open_some()
                        rv_limit.value[counter.value] = Utils2.Int.of_bytes(limit.value)
                        counter.value += 1
                sp.for elem in sp.range(sp.nat(0), sp.len(rv_names.value)):
                    coin_name_limit.value[rv_names.value.get(elem)] = rv_limit.value.get(elem)
        return sp.record(rv = sp.record(coin_name_limit = coin_name_limit.value,
                         net = net.value), status = is_error.value)

    # encoding starts here

    def encode_service_message(self, params):
        sp.set_type(params, sp.TRecord(service_type_value=sp.TNat, data=sp.TBytes))

        encode_service_type = sp.view("encode_nat", self.data.helper, params.service_type_value,
                                      t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_service_type, params.data],
                                        t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix,
                                              t=sp.TBytes).open_some()

        return final_rlp_bytes_with_prefix

    def encode_transfer_coin_msg(self, data):
        sp.set_type(data, sp.TRecord(from_addr=sp.TString, to=sp.TString,
        assets=sp.TList(sp.TRecord(coin_name=sp.TString, value=sp.TNat, fee=sp.TNat))))

        rlp = sp.local("rlp", sp.bytes("0x"))
        rlp_list = sp.local("rlp_list", [], t=sp.TList(sp.TBytes))
        temp = sp.local("temp", sp.bytes("0x"))
        coin_name = sp.local("coin_name", sp.bytes("0x"))
        sp.for item in data.assets:
            coin_name.value = sp.view("encode_string", self.data.helper, item.coin_name, t=sp.TBytes).open_some()
            temp.value = sp.view("encode_nat", self.data.helper, item.value, t=sp.TBytes).open_some()
            rlp_list.value.push(
                sp.view("encode_list", self.data.helper, [coin_name.value, temp.value], t=sp.TBytes).open_some())
            # rlp.value = sp.view("with_length_prefix", self.data.helper, rlp.value,
            #                                       t=sp.TBytes).open_some()

        assets_list = sp.view("encode_list", self.data.helper, rlp_list.value, t=sp.TBytes).open_some()
        from_addr_encoded = sp.view("encode_string", self.data.helper, data.from_addr, t=sp.TBytes).open_some()
        to_addr_encoded = sp.view("encode_string", self.data.helper, data.to, t=sp.TBytes).open_some()
        rlp.value = sp.view("encode_list", self.data.helper, [from_addr_encoded, to_addr_encoded, assets_list],
                            t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp.value,
                                              t=sp.TBytes).open_some()

        return final_rlp_bytes_with_prefix

    def encode_response(self, params):
        sp.set_type(params, sp.TRecord(code=sp.TNat, message=sp.TString))

        encode_code = sp.view("encode_nat", self.data.helper, params.code, t=sp.TBytes).open_some()
        encode_message = sp.view("encode_string", self.data.helper, params.message, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_code, encode_message],
                                        t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix,
                                              t=sp.TBytes).open_some()

        return final_rlp_bytes_with_prefix