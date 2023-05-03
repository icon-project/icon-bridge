import smartpy as sp

Utils = sp.io.import_script_from_url("https://raw.githubusercontent.com/Acurast/acurast-hyperdrive/main/contracts/tezos/libs/utils.py")
Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")


def decode_bmc_message(self, rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
    rlp_ = decode_list(rlp)
    temp_map_string = sp.compute(sp.map(tkey=sp.TString, tvalue=sp.TString))
    temp_int = sp.local("int_value", 0)
    temp_byt = sp.local("byt_value", sp.bytes("0x"))
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        # sp.trace(sp.bytes("0xf8c8") + i.value)
        sp.if counter.value == 0:
            temp_map_string["src"] = decode_string(i.value)
            # sp.trace(decode_string(i.value))
        sp.if counter.value == 1:
            temp_map_string["dst"] = decode_string(i.value)
            # sp.trace(decode_string(i.value))
        sp.if counter.value == 2:
            temp_map_string["svc"] = decode_string(i.value)
            # sp.trace(decode_string(i.value))
        sp.if counter.value == 3:
            temp_int.value = Utils2.Int.of_bytes(i.value)
            # sp.trace(Utils2.Int.of_bytes(i.value))
        sp.if counter.value == 4:
            temp_byt.value = i.value
            # sp.trace(i.value)
        counter.value = counter.value + 1

    return sp.record(src=temp_map_string["src"],
                     dst=temp_map_string["dst"],
                     svc=temp_map_string["svc"],
                     sn=temp_int.value,
                     message=temp_byt.value)


def decode_response(self, rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
    temp_int = sp.local("int1", 0)
    temp_byt = sp.local("byt1", sp.bytes("0x"))
    rlp_ = decode_list(rlp)
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        sp.if counter.value == 0:
            temp_int.value = Utils2.Int.of_bytes(i.value)
        sp.if counter.value == 1:
            temp_byt.value = i.value
        counter.value = counter.value + 1

    return sp.record(code=temp_int.value, message=decode_string(temp_byt.value))


def decode_bmc_service(self, rlp):
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


def decode_service_message(self, rlp):
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

    return sp.record(serviceType=temp_int.value,
                     data=temp_byt.value)


def decode_gather_fee_message(self, rlp):
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
    sp.
    for x in new_sub_list.items():
        _svcs.value[counter.value] = decode_string(x.value)
        counter.value = counter.value + 1
    return sp.record(fa=temp_str,
                     svcs=_svcs)

def decode_event_log(self, rlp):
    decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
    eventMptNode = sp.local("eventMptNode", sp.TMap(sp.TNat, sp.TString))
    rlp_ = decode_list(rlp)
    counter = sp.local("counter", 0)
    sp.for i in rlp_.items():
        eventMptNode.value[counter.value] = i.value
        counter.value = counter.value + 1
    return eventMptNode


def decode_receipt_proof(self, rlp):
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
                temp_map = self.decode_event_log(i.value)
            counter.value = counter.value + 1
            ep.value[counter.value] = sp.record(index=temp_int, eventMptNode=temp_map)
    return sp.record(index=rv_int.value, txReceipts=txReceipts.value, eventMptNode=ep.value)
