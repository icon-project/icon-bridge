import smartpy as sp

Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class DecodeLibrary:
    def decode_response(self, rlp):
        temp_int = sp.local("int1", 0)
        temp_byt = sp.local("byt1", sp.bytes("0x"))
        rlp_ = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        counter = sp.local("counter_response", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 0:
                temp_int.value = Utils2.Int.of_bytes(i.value)
            sp.if counter.value == 1:
                temp_byt.value = i.value
            counter.value = counter.value + 1

        return sp.record(code=temp_int.value, message=sp.view("decode_string", self.data.helper, temp_byt.value, t=sp.TString).open_some())

    def decode_service_message(self, rlp):
        rlp_ = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        temp_int = sp.local("int2", 0)
        temp_byt = sp.local("byte1", sp.bytes("0x"))
        counter = sp.local("counter", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 0:
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

        return sp.record(serviceType=_service_type.value,
                         data=temp_byt.value)

    def decode_transfer_coin_msg(self, rlp):
        rlp_ = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()

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
        starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
        sub_list = sp.local("sub_list", temp_byt.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
        new_sub_list = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        counter.value = 0
        new_temp_byt = sp.local("new_temp_byt", sp.bytes("0x"))
        rv_assets = sp.local("assets", {}, sp.TMap(sp.TNat, types.Types.Asset))
        sp.for x in new_sub_list.items():
            new_temp_byt.value = x.value
            sp.if sp.slice(new_temp_byt.value, 0, 2).open_some() == sp.bytes("0xb846"):
                new_temp_byt.value = sp.slice(new_temp_byt.value, 2, sp.as_nat(sp.len(new_temp_byt.value) - 2)).open_some()
            temp_byt = sp.local("tempByt2", sp.bytes("0x"))
            temp_int = sp.local("tempInt", sp.nat(0))
            counter.value = 0
            view_value = sp.view("decode_list", self.data.helper, new_temp_byt.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            sp.for i in view_value.items():
                sp.if counter.value == 1:
                    temp_int.value = Utils2.Int.of_bytes(i.value)
                sp.if counter.value == 0:
                    temp_byt.value = i.value
                rv_assets.value[counter.value] = sp.record(coin_name=sp.view("decode_string", self.data.helper, temp_byt.value, t=sp.TString).open_some()
                                                           , value=temp_int.value)
                counter.value = counter.value + 1
        rv1 = sp.view("decode_string", self.data.helper, rv1_byt.value, t=sp.TString).open_some()
        rv2 = sp.view("decode_string", self.data.helper, rv2_byt.value, t=sp.TString).open_some()
        return sp.record(from_= rv1, to = rv2 , assets = rv_assets.value)

    def decode_blacklist_msg(self, rlp):
        rlp_ = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()

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
        starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
        sub_list = sp.local("sub_list", temp_byt.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
        new_sub_list = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        counter.value = 0
        new_temp_byt = sp.local("new_temp_byt", sp.bytes("0x"))
        rv_blacklist_address = sp.local("blacklist_data", {}, sp.TMap(sp.TNat, sp.TString))
        sp.for x in new_sub_list.items():
            new_temp_byt.value = x.value
            sp.if sp.slice(new_temp_byt.value, 0, 2).open_some() == sp.bytes("0xb846"):
                new_temp_byt.value = sp.slice(new_temp_byt.value, 2, sp.as_nat(sp.len(new_temp_byt.value) - 2)).open_some()
            counter.value = 0
            _decode_list = sp.view("decode_list", self.data.helper, new_temp_byt.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            sp.for j in _decode_list.items():
                rv_blacklist_address.value[counter.value] = sp.view("decode_string", self.data.helper, j.value, t=sp.TString).open_some()
                counter.value = counter.value + 1
        rv1 = Utils2.Int.of_bytes(rv1_byt.value)
        rv2 = sp.view("decode_string", self.data.helper, rv2_byt.value, t=sp.TString).open_some()
        _service_type = sp.local("_service_type_blacklist", sp.variant("", 10))
        with sp.if_(rv1 == 0):
            _service_type.value = sp.variant("ADD_TO_BLACKLIST", rv1)
        with sp.else_():
            _service_type.value = sp.variant("REMOVE_FROM_BLACKLIST", rv1)
        return sp.record(serviceType = _service_type.value , addrs = rv_blacklist_address.value ,
                         net = sp.view("decode_string", self.data.helper, rv2, t=sp.TString).open_some())

    def decode_token_limit_msg(self, rlp):
        rlp_ = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()

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
        starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
        sub_list = sp.local("sub_list", temp_byt.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
        new_sub_list = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        counter.value = 0
        rv_names = sp.local("names", {}, sp.TMap(sp.TNat, sp.TString))
        rv_limit = sp.local("limit", {}, sp.TMap(sp.TNat, sp.TNat))
        sp.for x in new_sub_list.items():
            rv_names.value[counter.value] = sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some()

        new_sub_list1 = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        sp.for y in new_sub_list1.items():
            rv_limit.value[counter.value] = Utils2.Int.of_bytes(y.value)
        return sp.record(coin_name = rv_names.value, token_limit = rv_limit.value ,
                         net = sp.view("decode_string", self.data.helper, rv1_byt.value, t=sp.TString).open_some())