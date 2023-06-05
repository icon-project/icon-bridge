import smartpy as sp

Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class DecodeLibrary:
    def decode_response(self, rlp):
        temp_int = sp.local("int1", 0)
        temp_byt = sp.local("byt1", sp.bytes("0x"))
        rlp_dr = sp.local("rlp_dr_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_dr.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_dr.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_dr.value
        counter = sp.local("counter_response", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 0:
                # check_length = sp.view("prefix_length", self.data.helper, i.value, t=sp.TNat).open_some()
                # with sp.if_(check_length > 0):
                #     i.value = sp.view("without_length_prefix", self.data.helper, i.value,
                #                       t=sp.TBytes).open_some()
                temp_int.value = Utils2.Int.of_bytes(i.value)
            sp.if counter.value == 1:
                temp_byt.value = i.value
            counter.value = counter.value + 1

        return sp.record(code=temp_int.value, message=sp.view("decode_string", self.data.helper, temp_byt.value, t=sp.TString).open_some())

    def decode_service_message(self, rlp):
        rlp_sm = sp.local("rlp_sm_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_sm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_sm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_sm.value
        temp_int = sp.local("int2", 0)
        temp_byt = sp.local("byte1", sp.bytes("0x"))
        counter = sp.local("counter", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 0:
                # check_length = sp.view("prefix_length", self.data.helper, i.value, t=sp.TNat).open_some()
                # with sp.if_(check_length > 0):
                #     i.value = sp.view("without_length_prefix", self.data.helper, i.value,
                #                       t=sp.TBytes).open_some()
                temp_int.value = Utils2.Int.of_bytes(i.value)
            sp.if counter.value == 1:
                temp_byt.value = i.value
            counter.value = counter.value + 1

        _service_type = sp.local("_service_type", sp.variant("", 10))
        sp.if temp_int.value == 0:
            _service_type.value = sp.variant("REQUEST_COIN_TRANSFER", temp_int.value)
        sp.if temp_int.value == 1:
            _service_type.value = sp.variant("REQUEST_COIN_REGISTER", temp_int.value)
        sp.if temp_int.value == 2:
            _service_type.value = sp.variant("RESPONSE_HANDLE_SERVICE", temp_int.value)
        sp.if temp_int.value == 3:
            _service_type.value = sp.variant("BLACKLIST_MESSAGE", temp_int.value)
        sp.if temp_int.value == 4:
                _service_type.value = sp.variant("CHANGE_TOKEN_LIMIT", temp_int.value)
        sp.if temp_int.value == 5:
            _service_type.value = sp.variant("UNKNOWN_TYPE", temp_int.value)
        temp_byt.value = sp.view("without_length_prefix", self.data.helper, temp_byt.value, t=sp.TBytes).open_some()

        return sp.record(serviceType=_service_type.value,
                         data=temp_byt.value)

    def decode_transfer_coin_msg(self, rlp):
        rlp_tcm = sp.local("rlp_tcm_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_tcm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_tcm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_tcm.value

        temp_byt = sp.local("byt_transfer", sp.bytes("0x"))
        rv1_byt = sp.local("rv1_byt", sp.bytes("0x"))
        rv2_byt = sp.local("rv2_byt", sp.bytes("0x"))
        counter = sp.local("counter_coin", sp.nat(0))
        sp.for i in rlp_.items():
            sp.if counter.value == 2:
                temp_byt.value = i.value
            sp.if counter.value == 0:
                rv1_byt.value = i.value
            sp.if counter.value == 1:
                rv2_byt.value = i.value
            counter.value = counter.value + 1
        sub_list = sp.local("sub_list", temp_byt.value)
        nsl_tcm = sp.local("nsl_bts_tcm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl_tcm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl_tcm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_tcm.value
        counter.value = 0
        new_temp_byt = sp.local("new_temp_byt", sp.bytes("0x"))
        rv_assets = sp.local("assets", {}, sp.TMap(sp.TNat, types.Types.Asset))
        nsl3_tcm = sp.local("nsl3_bts_tcm", sp.map(tkey=sp.TNat))
        view_value = sp.local("view_value", sp.map(tkey=sp.TNat))
        counter_nested = sp.local("counter_nested", sp.nat(0), t=sp.TNat)
        temp_byt = sp.local("tempByt2", sp.bytes("0x"))
        temp_byt_nested = sp.local("tempByt2nested", sp.bytes("0x"))
        sp.for x in new_sub_list.items():
            new_temp_byt.value = x.value
            # sp.if sp.slice(new_temp_byt.value, 0, 2).open_some() == sp.bytes("0xb846"):
            #     new_temp_byt.value = sp.slice(new_temp_byt.value, 2, sp.as_nat(sp.len(new_temp_byt.value) - 2)).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, new_temp_byt.value,
                                     t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl3_tcm.value = sp.view("decode_list", self.data.helper, new_temp_byt.value,
                                        t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                decode_len = sp.view("without_length_prefix", self.data.helper, new_temp_byt.value,
                                     t=sp.TBytes).open_some()
                nsl3_tcm.value = sp.view("decode_list", self.data.helper, decode_len,
                                         t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            view_value.value = nsl3_tcm.value
            counter_nested.value = sp.nat(0)
            sp.for i in view_value.value.items():
                sp.if counter_nested.value == 1:
                    # check_length = sp.view("prefix_length", self.data.helper, i.value, t=sp.TNat).open_some()
                    # with sp.if_ (check_length > 0):
                    temp_byt_nested.value = sp.view("without_length_prefix", self.data.helper, i.value,
                                                    t=sp.TBytes).open_some()
                sp.if counter_nested.value == 0:
                    temp_byt.value = i.value
                counter_nested.value += 1

            rv_assets.value[counter.value] = sp.record(coin_name=sp.view("decode_string", self.data.helper, temp_byt.value, t=sp.TString).open_some()
                                                       , value=Utils2.Int.of_bytes(temp_byt_nested.value))

            counter.value = counter.value + 1

        rv1 = sp.view("decode_string", self.data.helper, rv1_byt.value, t=sp.TString).open_some()
        rv2 = sp.view("decode_string", self.data.helper, rv2_byt.value, t=sp.TString).open_some()
        return sp.record(from_= rv1, to = rv2 , assets = rv_assets.value)

    def decode_blacklist_msg(self, rlp):
        rlp_bm = sp.local("rlp_bm_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_bm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_bm.value

        temp_byt = sp.local("byt_transfer", sp.bytes("0x"))
        rv1_byt = sp.local("rv1_byt", sp.bytes("0x"))
        rv2_byt = sp.local("rv2_byt", sp.bytes("0x"))
        counter = sp.local("counter_blacklist", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 2:
                rv2_byt.value = i.value
            sp.if counter.value == 0:
                rv1_byt.value = i.value
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
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl_bm.value = sp.view("decode_list", self.data.helper, decode_len,
                                     t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_bm.value
        counter.value = 0
        new_temp_byt = sp.local("new_temp_byt", sp.bytes("0x"))
        rv_blacklist_address = sp.local("blacklist_data", {}, sp.TMap(sp.TNat, sp.TString))
        addr_string = sp.local("addr_string", "")
        nsl2_bm = sp.local("nsl2_bts_bm", sp.map(tkey=sp.TNat))
        counter.value = 0
        sp.for x in new_sub_list.items():
            new_temp_byt.value = x.value
            # sp.if sp.slice(new_temp_byt.value, 0, 2).open_some() == sp.bytes("0xb846"):
            #     new_temp_byt.value = sp.slice(new_temp_byt.value, 2, sp.as_nat(sp.len(new_temp_byt.value) - 2)).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, new_temp_byt.value,
                                     t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl2_bm.value = sp.view("decode_list", self.data.helper, new_temp_byt.value,
                                         t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                decode_len = sp.view("without_length_prefix", self.data.helper, new_temp_byt.value, t=sp.TBytes).open_some()
                addr_string.value = sp.view("decode_string", self.data.helper, decode_len,
                                         t=sp.TString).open_some()
            # _decode_list = nsl2_bm.value
            # sp.for j in _decode_list.items():
            rv_blacklist_address.value[counter.value] = addr_string.value
            counter.value = counter.value + 1
        check_length = sp.view("prefix_length", self.data.helper, rv1_byt.value, t=sp.TNat).open_some()
        with sp.if_(check_length > 0):
            rv1_byt.value = sp.view("without_length_prefix", self.data.helper, rv1_byt.value,
                              t=sp.TBytes).open_some()
        rv1 = Utils2.Int.of_bytes(rv1_byt.value)
        rv2 = sp.view("decode_string", self.data.helper, rv2_byt.value, t=sp.TString).open_some()
        _service_type = sp.local("_service_type_blacklist", sp.variant("", 10))
        with sp.if_(rv1 == 0):
            _service_type.value = sp.variant("ADD_TO_BLACKLIST", rv1)
        with sp.else_():
            _service_type.value = sp.variant("REMOVE_FROM_BLACKLIST", rv1)
        return sp.record(serviceType = _service_type.value , addrs = rv_blacklist_address.value , net = rv2)

    def decode_token_limit_msg(self, rlp):
        rlp_tlm = sp.local("rlp_tlm_bts", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_tlm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_tlm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_tlm.value

        temp_byt = sp.local("byt_transfer", sp.bytes("0x"))
        temp_byt1 = sp.local("byt_transfer_temp1", sp.bytes("0x"))
        rv1_byt = sp.local("rv1_byt", sp.bytes("0x"))
        counter = sp.local("counter_token_limit", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 0:
                temp_byt.value = i.value
            sp.if counter.value == 1:
                temp_byt1.value = i.value
            sp.if counter.value == 2:
                rv1_byt.value = i.value
            counter.value = counter.value + 1
        sub_list = sp.local("sub_list", temp_byt.value)
        sub_list_limit = sp.local("sub_list_limit", temp_byt1.value)
        nsl1_dtlm = sp.local("nsl1_bts_dtlm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl1_dtlm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                     t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl1_dtlm.value = sp.view("decode_list", self.data.helper, decode_len,
                                     t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl1_dtlm.value
        counter.value = 0
        rv_names = sp.local("names", {}, sp.TMap(sp.TNat, sp.TString))
        rv_limit = sp.local("limit", {}, sp.TMap(sp.TNat, sp.TNat))
        sp.for x in new_sub_list.items():
            rv_names.value[counter.value] = sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some()
            counter.value += 1
        nsl_dtlm = sp.local("nsl_bts_dtlm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        counter.value = 0
        with sp.if_(is_list_lambda):
            nsl_dtlm.value = sp.view("decode_list", self.data.helper, sub_list_limit.value,
                                    t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list_limit.value, t=sp.TBytes).open_some()
            nsl_dtlm.value = sp.view("decode_list", self.data.helper, decode_len,
                                    t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list1 = nsl_dtlm.value
        limit = sp.local("limit_val", sp.bytes("0x"), t=sp.TBytes)
        sp.for y in new_sub_list1.items():
            check_length = sp.view("prefix_length", self.data.helper, y.value, t=sp.TNat).open_some()
            with sp.if_(check_length > 0):
                limit.value = sp.view("without_length_prefix", self.data.helper, y.value,
                                  t=sp.TBytes).open_some()
            rv_limit.value[counter.value] = Utils2.Int.of_bytes(limit.value)
            counter.value += 1
        return sp.record(coin_name = rv_names.value, token_limit = rv_limit.value ,
                         net = sp.view("decode_string", self.data.helper, rv1_byt.value, t=sp.TString).open_some())