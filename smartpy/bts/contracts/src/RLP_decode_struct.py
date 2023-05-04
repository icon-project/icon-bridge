import smartpy as sp

Utils = sp.io.import_script_from_url("https://raw.githubusercontent.com/Acurast/acurast-hyperdrive/main/contracts/tezos/libs/utils.py")
Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")

def decode_bmc_message(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
    rlp_ = decode_list(rlp)
    temp_map_string = sp.compute(sp.map(tkey=sp.TString, tvalue=sp.TString))
    temp_int = sp.local("int_value", 0)
    temp_byt = sp.local("byt_value", sp.bytes("0x"))
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        sp.if counter.value == 0:
            temp_map_string["src"] = decode_string(i.value)
        sp.if counter.value == 1:
            temp_map_string["dst"] = decode_string(i.value)
        sp.if counter.value == 2:
            temp_map_string["svc"] = decode_string(i.value)
        sp.if counter.value == 3:
            temp_int.value = Utils2.Int.of_bytes(i.value)
        sp.if counter.value == 4:
            temp_byt.value = i.value
        counter.value = counter.value + 1

    return sp.record(src=temp_map_string.get("src"),
                     dst=temp_map_string.get("dst"),
                     svc=temp_map_string.get("svc"),
                     sn=temp_int.value,
                     message=temp_byt.value)


def decode_response(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
    temp_int = sp.local("int1", 0)
    temp_byt = sp.local("byt1", sp.bytes("0x"))
    rlp_ = decode_list(rlp)
    counter = sp.local("counter_response", 0)
    sp.for i in rlp_.items():
        sp.if counter.value == 0:
            temp_int.value = Utils2.Int.of_bytes(i.value)
        sp.if counter.value == 1:
            temp_byt.value = i.value
        counter.value = counter.value + 1

    return sp.record(code=temp_int.value, message=decode_string(temp_byt.value))


def decode_bmc_service(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
    rlp_ = decode_list(rlp)
    temp_string = sp.local("str_value", 0)
    temp_byt = sp.local("byt_value", sp.bytes("0x"))
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        sp.if counter.value == 0:
            temp_string.value = decode_string(i.value)
        sp.if counter.value == 1:
            temp_byt.value = i.value
        counter.value = counter.value + 1

    return sp.record(serviceType=temp_string.value,
                     payload=temp_byt.value)


def decode_service_message(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    rlp_ = decode_list(rlp)
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


def decode_gather_fee_message(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
    rlp_ = decode_list(rlp)
    temp_byt = sp.local("byt4", sp.bytes("0x"))
    counter = sp.local("counter", 0)
    temp_str = sp.local("str_gather", "")
    sp.for i in rlp_.items():
        sp.if counter.value == 1:
            temp_byt.value = i.value
        sp.if counter.value == 0:
            temp_str.value = decode_string(i.value)
        counter.value = counter.value + 1
    starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
    sub_list = sp.local("sub_list", temp_byt.value)
    sp.if starts_with == sp.bytes("0xb846"):
        sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
    new_sub_list = decode_list(sub_list.value)
    _svcs = sp.local("_svcs", sp.TMap(sp.TNat, sp.TString))
    counter.value = 0
    sp.for x in new_sub_list.items():
        _svcs.value[counter.value] = decode_string(x.value)
        counter.value = counter.value + 1
    return sp.record(fa=temp_str,
                     svcs=_svcs)

def decode_event_log(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    event_mpt_node = sp.local("eventMptNode", sp.TMap(sp.TNat, sp.TString))
    rlp_ = decode_list(rlp)
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        event_mpt_node.value[counter.value] = i.value
        counter.value = counter.value + 1
    return event_mpt_node


def decode_receipt_proof(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    rlp_ = decode_list(rlp)
    temp_byt = sp.local("byt_receipt", sp.bytes("0x"))
    temp_byt2 = sp.local("byt_receipt2", sp.bytes("0x"))
    rv_int = sp.local("rv_int", 0)
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        sp.if counter.value == 1:
            temp_byt.value = i.value
        sp.if counter.value == 0:
            rv_int.value = Utils2.Int.of_bytes(i.value)
        sp.if counter.value == 2:
            temp_byt2 = i.value
        counter.value = counter.value + 1

    starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
    sub_list = sp.local("sub_list", temp_byt.value)
    sp.if starts_with == sp.bytes("0xb846"):
        sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
    new_sub_list = decode_list(sub_list.value)
    counter.value = 0
    txReceipts = sp.local("txReceipts", sp.TMap(sp.TNat, sp.TBytes))
    sp.for x in new_sub_list.items():
        txReceipts.value[counter.value] = x.value
        counter.value = counter.value + 1

    starts_with_second_elem = sp.slice(temp_byt2.value, 0, 2).open_some()
    sub_list_second_elem = sp.local("sub_list_second_elem", temp_byt.value)
    sp.if starts_with_second_elem == sp.bytes("0xb846"):
        sub_list_second_elem.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
    new_sub_list = decode_list(sub_list_second_elem.value)
    ep = sp.local("ep", sp.TMap(sp.TNat, sp.TBytes))
    new_temp_byt = sp.local("new_temp_byt", sp.bytes("0x"))
    sp.for x in new_sub_list.items():
        new_temp_byt.value = x.value
        sp.if sp.slice(new_temp_byt.value, 0, 2).open_some() == sp.bytes("0xb846"):
            new_temp_byt.value = sp.slice(new_temp_byt.value, 2, sp.as_nat(sp.len(new_temp_byt.value) - 2)).open_some()
        temp_map = sp.local("tempMap", sp.TMap(sp.TNat, sp.TString))
        temp_int = sp.local("tempInt", 0)
        counter.value = 0
        sp.for i in decode_list(new_temp_byt.value):
            sp.if counter.value == 0:
                temp_int = Utils2.Int.of_bytes(i.value)
            sp.if counter.value == 1:
                temp_map = decode_event_log(i.value)
            counter.value = counter.value + 1
            ep.value[counter.value] = sp.record(index=temp_int, eventMptNode=temp_map)
    return sp.record(index=rv_int.value, txReceipts=txReceipts.value, eventMptNode=ep.value)

def decode_transfer_coin_msg(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    rlp_ = decode_list(rlp)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)

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
    new_sub_list = decode_list(sub_list.value)
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
        sp.for i in decode_list(new_temp_byt.value).items():
            sp.if counter.value == 1:
                temp_int.value = Utils2.Int.of_bytes(i.value)
            sp.if counter.value == 0:
                temp_byt.value = i.value
            rv_assets.value[counter.value] = sp.record(coin_name=decode_string(temp_byt.value), value=temp_int.value)
            counter.value = counter.value + 1
    rv1 = decode_string(rv1_byt.value)
    rv2 = decode_string(rv2_byt.value)
    return sp.record(from_= rv1, to = rv2 , assets = rv_assets.value)

def decode_blacklist_msg(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    rlp_ = decode_list(rlp)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)

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
    new_sub_list = decode_list(sub_list.value)
    counter.value = 0
    new_temp_byt = sp.local("new_temp_byt", sp.bytes("0x"))
    rv_blacklist_address = sp.local("blacklist_data", {}, sp.TMap(sp.TNat, sp.TString))
    sp.for x in new_sub_list.items():
        new_temp_byt.value = x.value
        sp.if sp.slice(new_temp_byt.value, 0, 2).open_some() == sp.bytes("0xb846"):
            new_temp_byt.value = sp.slice(new_temp_byt.value, 2, sp.as_nat(sp.len(new_temp_byt.value) - 2)).open_some()
        counter.value = 0
        sp.for j in decode_list(new_temp_byt.value).items():
            rv_blacklist_address.value[counter.value] = decode_string(j.value)
            counter.value = counter.value + 1
    rv1 = Utils2.Int.of_bytes(rv1_byt.value)
    rv2 = decode_string(rv2_byt.value)
    _service_type = sp.local("_service_type_blacklist", sp.variant("", 10))
    with sp.if_(rv1 == 0):
        _service_type.value = sp.variant("ADD_TO_BLACKLIST", rv1)
    with sp.else_():
        _service_type.value = sp.variant("REMOVE_FROM_BLACKLIST", rv1)
    return sp.record(serviceType = _service_type.value , addrs = rv_blacklist_address.value , net = decode_string(rv2))

def decode_token_limit_msg(rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    rlp_ = decode_list(rlp)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)

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
    new_sub_list = decode_list(sub_list.value)
    counter.value = 0
    rv_names = sp.local("names", {}, sp.TMap(sp.TNat, sp.TString))
    rv_limit = sp.local("limit", {}, sp.TMap(sp.TNat, sp.TNat))
    sp.for x in new_sub_list.items():
        rv_names.value[counter.value] = decode_string(x.value)

    new_sub_list1 = decode_list(sub_list.value)
    sp.for y in new_sub_list1.items():
        rv_limit.value[counter.value] = Utils2.Int.of_bytes(y.value)
    return sp.record(coin_name = rv_names.value, token_limit = rv_limit.value , net = decode_string(rv1_byt.value))