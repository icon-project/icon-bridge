import smartpy as sp

Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")

# contract address to deal with negative values
#TODO: change to mainnet address
HELPER_CONTRACT_ADDRESS = sp.address("KT1W6dU9xpKwMwHXopVhW5PB1NdZFmVZPKbK")

class DecodeLibrary:

    def decode_bmc_message(self, rlp):
        rlp_bm = sp.local("rlp_bm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_bm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_bm.value
        temp_map_string = sp.compute(sp.map(tkey=sp.TString, tvalue=sp.TString))
        temp_int = sp.local("int_value", 0)
        temp_byt = sp.local("byt_value", sp.bytes("0x"))
        counter = sp.local("counter", 0)
        sp.for k in rlp_.items():
            sp.if counter.value == 0:
                temp_map_string["src"] = sp.view("decode_string", self.data.helper, k.value, t=sp.TString).open_some()
            sp.if counter.value == 1:
                temp_map_string["dst"] = sp.view("decode_string", self.data.helper, k.value, t=sp.TString).open_some()
            sp.if counter.value == 2:
                temp_map_string["svc"] = sp.view("decode_string", self.data.helper, k.value, t=sp.TString).open_some()
            sp.if counter.value == 3:
                sn_in_bytes = sp.view("without_length_prefix", self.data.helper, k.value, t=sp.TBytes).open_some()
                temp_int.value = sp.view("to_int", HELPER_CONTRACT_ADDRESS, sn_in_bytes, t=sp.TInt).open_some()
            sp.if counter.value == 4:
                temp_byt.value = k.value
            counter.value = counter.value + 1
        temp_byt.value = sp.view("without_length_prefix", self.data.helper, temp_byt.value, t=sp.TBytes).open_some()
        return sp.record(src=temp_map_string.get("src"),
                         dst=temp_map_string.get("dst"),
                         svc=temp_map_string.get("svc"),
                         sn=temp_int.value,
                         message=temp_byt.value)


    def decode_response(self, rlp):
        temp_int = sp.local("int1", 0)
        temp_byt = sp.local("byt1", sp.bytes("0x"))
        rlp_dr = sp.local("rlp_dr", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_dr.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_dr.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_dr.value
        counter = sp.local("counter_response", 0)
        sp.for m in rlp_.items():
            sp.if counter.value == 0:
                temp_int.value = Utils2.Int.of_bytes(m.value)
            sp.if counter.value == 1:
                temp_byt.value = m.value
            counter.value = counter.value + 1

        return sp.record(code=temp_int.value, message=sp.view("decode_string", self.data.helper, temp_byt.value, t=sp.TString).open_some())

    def decode_propagate_message(self, rlp):
        rlp_pm = sp.local("rlp_pm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_pm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_pm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_pm.value
        counter = sp.local("counter_propagate", 0)
        temp_string = sp.local("temp_string", "")
        sp.for d in rlp_.items():
            sp.if counter.value == 0:
                temp_string.value = sp.view("decode_string", self.data.helper, d.value, t=sp.TString).open_some()
            counter.value = counter.value + 1
        return temp_string.value

    def decode_init_message(self, rlp):
        rlp_im = sp.local("rlp_im", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_im.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_im.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_im.value
        counter = sp.local("counter_init", 0)
        temp_bytes = sp.local("byt_init", sp.bytes("0x"))
        sp.for g in rlp_.items():
            sp.if counter.value == 0:
                temp_bytes.value = g.value
            counter.value = counter.value + 1

        starts_with = sp.slice(temp_bytes.value, 0, 2).open_some()
        sub_list = sp.local("sub_list_init", temp_bytes.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_bytes.value, 2, sp.as_nat(sp.len(temp_bytes.value) - 2)).open_some()
        nsl_im = sp.local("nsl_im", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl_im.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl_im.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_im.value

        _links = sp.local("links_init", [], sp.TList(sp.TString))
        counter.value = 0
        sp.for x in new_sub_list.items():
            _links.value.push(sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some())
            counter.value = counter.value + 1
        return _links.value

    def decode_bmc_service(self, rlp):
        rlp_bs = sp.local("rlp_bs", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bs.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_bs.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_bs.value
        temp_string = sp.local("str_value", "")
        temp_byt = sp.local("byt_value_bmc", sp.bytes("0x"))
        counter = sp.local("counter_service", 0)
        sp.for b in rlp_.items():
            sp.if counter.value == 0:
                temp_string.value = sp.view("decode_string", self.data.helper, b.value, t=sp.TString).open_some()
            sp.if counter.value == 1:
                temp_byt.value = b.value
            counter.value = counter.value + 1
        temp_byt.value = sp.view("without_length_prefix", self.data.helper, temp_byt.value, t=sp.TBytes).open_some()
        return sp.record(serviceType=temp_string.value,
                         payload=temp_byt.value)

    def decode_gather_fee_message(self, rlp):
        rlp_gm = sp.local("rlp_gm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_gm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_gm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_gm.value
        temp_byt = sp.local("byt4", sp.bytes("0x"))
        counter = sp.local("counter_gather", 0)
        temp_str = sp.local("str_gather", "")
        sp.for c in rlp_.items():
            sp.if counter.value == 1:
                temp_byt.value = c.value
            sp.if counter.value == 0:
                temp_str.value = sp.view("decode_string", self.data.helper, c.value, t=sp.TString).open_some()
            counter.value = counter.value + 1
        starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
        sub_list = sp.local("sub_list", temp_byt.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()
        nsl_gm = sp.local("nsl_gm", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl_gm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl_gm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_gm.value

        _svcs = sp.local("_svcs", {}, sp.TMap(sp.TNat, sp.TString))
        counter.value = 0
        sp.for x in new_sub_list.items():
            _svcs.value[counter.value] = sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some()
            counter.value = counter.value + 1
        return sp.record(fa=temp_str.value,
                         svcs=_svcs.value)

    def to_message_event(self, rlp):
        rlp_me = sp.local("rlp_me", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_me.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_me.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_me.value
        counter = sp.local("counter_event", 0)
        rv1 = sp.local("rv1_event", "")
        rv2 = sp.local("rv2_event", sp.nat(0))
        rv3 = sp.local("rv3_event", sp.bytes("0x"))
        sp.for i in rlp_.items():
            sp.if counter.value == 2:
                rv3.value = i.value
            sp.if counter.value == 0:
                rv1.value = sp.view("decode_string", self.data.helper, i.value, t=sp.TString).open_some()
            sp.if counter.value == 1:
                rv2.value = Utils2.Int.of_bytes(i.value)
            counter.value = counter.value + 1
        rv3.value = sp.view("without_length_prefix", self.data.helper, rv3.value, t=sp.TBytes).open_some()
        return sp.record(next_bmc= rv1.value, seq= rv2.value, message = rv3.value)

    def decode_receipt_proof(self, rlp):
        rlp_rp = sp.local("rlp_rp", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_rp.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_rp.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_rp.value
        temp_byt = sp.local("byt_receipt", sp.bytes("0x"))
        rv_int = sp.local("rv_int_receipt", 0)
        rv_int2 = sp.local("rv_int2_receipt", 0)
        counter = sp.local("counter", 0)
        sp.for i in rlp_.items():
            sp.if counter.value == 1:
                temp_byt.value = sp.view("without_length_prefix", self.data.helper, i.value, t=sp.TBytes).open_some()
            sp.if counter.value == 0:
                rv_int.value = Utils2.Int.of_bytes(i.value)
            sp.if counter.value == 2:
                rv_int2.value =Utils2.Int.of_bytes(i.value)
            counter.value = counter.value + 1

        starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
        sub_list = sp.local("sub_list", temp_byt.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()

        nsl_rp = sp.local("nsl_rp", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl_rp.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                    t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl_rp.value = sp.view("decode_list", self.data.helper, decode_len,
                                    t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_rp.value
        counter.value = 0
        events = sp.local("events_receipt", sp.map({}, tkey=sp.TNat,
                                                    tvalue=sp.TRecord(next_bmc= sp.TString,
                                                              seq= sp.TNat,
                                                              message = sp.TBytes)))
        sp.for z in new_sub_list.items():
            events.value[counter.value] = self.to_message_event(z.value)
            counter.value = counter.value + 1
        return sp.record(index = rv_int.value, events = events.value, height = rv_int2.value)


    def decode_receipt_proofs(self, rlp):
        rlp_rps = sp.local("rlp_rps", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_rps.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            rlp_rps.value = sp.view("decode_list", self.data.helper, decode_len, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_rps.value
        counter = sp.local("counter_receipt_proofs", 0)
        receipt_proofs = sp.local("events_receipt_proofs", sp.map({}, tkey=sp.TNat,
                tvalue=sp.TRecord(index = sp.TNat,
                          events = sp.TMap(sp.TNat, sp.TRecord(next_bmc=sp.TString, seq=sp.TNat, message=sp.TBytes)),
                          height = sp.TNat,
                        )
                                                                   )
                                  )
        temp_byt = sp.local("temp_byt_proofs", sp.bytes("0x"))
        sp.for i in rlp_.items():
            sp.if counter.value == 0:
                temp_byt.value = i.value
            counter.value = counter.value + 1
        starts_with = sp.slice(temp_byt.value, 0, 2).open_some()
        sub_list = sp.local("sub_list_proofs", temp_byt.value)
        sp.if starts_with == sp.bytes("0xb846"):
            sub_list.value = sp.slice(temp_byt.value, 2, sp.as_nat(sp.len(temp_byt.value) - 2)).open_some()

        nsl_rps = sp.local("nsl_rps", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl_rps.value = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            nsl_rps.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_rps.value
        counter.value = 0
        sp.if sp.len(new_sub_list) > 0:
            sp.for x in new_sub_list.items():
                receipt_proofs.value[counter.value] = self.decode_receipt_proof(x.value)
                counter.value = counter.value + 1
        return receipt_proofs.value

