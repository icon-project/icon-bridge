import smartpy as sp

Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class DecodeEncodeLibrary:

    def decode_bmc_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_bm = sp.local("rlp_bm", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bmc_message", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                rlp_bm.value = sp.view("decode_list", self.data.helper, decode_len,
                                   t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInBMCMessageDecoding"
        temp_map_string = sp.compute(sp.map(tkey=sp.TString, tvalue=sp.TString))
        temp_int = sp.local("int_value", 0)
        temp_byt = sp.local("byt_value", sp.bytes("0x"))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_bm.value
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
                    _to_int = sp.view("to_int", self.data.helper_parse_negative, sn_in_bytes, t=sp.TInt).open_some()
                    temp_int.value = _to_int
                sp.if counter.value == 4:
                    temp_byt.value = k.value
                counter.value = counter.value + 1
            with sp.if_(is_error.value == "Success"):
                temp_byt.value = sp.view("without_length_prefix", self.data.helper, temp_byt.value, t=sp.TBytes).open_some()

        return sp.record(bmc_dec_rv = sp.record(src=temp_map_string.get("src", default_value = "NoKey"),
                         dst=temp_map_string.get("dst", default_value = "NoKey"),
                         svc=temp_map_string.get("svc", default_value = "NoKey"),
                         sn=temp_int.value,
                         message=temp_byt.value),
                         status = is_error.value)

    def decode_response(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        temp_int = sp.local("int1", sp.nat(0))
        temp_byt = sp.local("byt1", sp.bytes("0x"))
        rlp_dr = sp.local("rlp_dr", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bmc_decode_response", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_dr.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
                decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
                is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
                with sp.if_(is_list_lambda):
                    rlp_dr.value = sp.view("decode_list", self.data.helper, decode_len,
                                           t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                with sp.else_():
                    is_error.value = "ErrorInDecodingBMCResponse"
        rlp_ = rlp_dr.value
        counter = sp.local("counter_response", 0)
        with sp.if_(is_error.value == "Success"):
            sp.for m in rlp_.items():
                sp.if counter.value == 0:
                    temp_int.value = Utils2.Int.of_bytes(m.value)
                sp.if counter.value == 1:
                    temp_byt.value = m.value
                counter.value = counter.value + 1

            # message in case of error is null which cannot be decoded into string
        return sp.record(code=temp_int.value, message="Error", status = is_error.value)

    def decode_init_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_im = sp.local("rlp_im", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bmc_decode_init_message", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_im.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
                decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
                is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
                with sp.if_(is_list_lambda):
                    rlp_im.value = sp.view("decode_list", self.data.helper, decode_len,
                                           t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                with sp.else_():
                    is_error.value = "ErrorInDecodingIniteMessage"
        _links = sp.local("links_init", [], sp.TList(sp.TString))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_im.value
            counter = sp.local("counter_init", 0)
            temp_bytes = sp.local("byt_init", sp.bytes("0x"))
            sp.for g in rlp_.items():
                sp.if counter.value == 0:
                    temp_bytes.value = g.value
                counter.value = counter.value + 1

            sub_list = sp.local("sub_list_init", temp_bytes.value)
            nsl_im = sp.local("nsl_im", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl_im.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                    decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value,
                                         t=sp.TBytes).open_some()
                    is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
                    with sp.if_(is_list_lambda):
                        nsl_im.value = sp.view("decode_list", self.data.helper, decode_len,
                                               t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                    with sp.else_():
                        is_error.value = "ErrorInDecodingInitMessage"

            with sp.if_(is_error.value == "Success"):
                new_sub_list = nsl_im.value

                counter.value = 0
                sp.for x in new_sub_list.items():
                    _links.value.push(sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some())
                    counter.value = counter.value + 1
        return sp.record(links_list = _links.value, status = is_error.value)

    def decode_bmc_service(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_bs = sp.local("rlp_bs", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bmc_decode_service", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bs.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                rlp_bs.value = sp.view("decode_list", self.data.helper, decode_len,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecodingBMCService"

        temp_string = sp.local("str_value", "")
        temp_byt = sp.local("byt_value_bmc", sp.bytes("0x"))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_bs.value
            counter = sp.local("counter_service", 0)
            sp.for b in rlp_.items():
                sp.if counter.value == 0:
                    temp_string.value = sp.view("decode_string", self.data.helper, b.value, t=sp.TString).open_some()
                sp.if counter.value == 1:
                    temp_byt.value = b.value
                counter.value = counter.value + 1
            temp_byt.value = sp.view("without_length_prefix", self.data.helper, temp_byt.value, t=sp.TBytes).open_some()
        return sp.record(bmc_service_rv=sp.record(serviceType=temp_string.value,
                         payload=temp_byt.value), status = is_error.value)

    def decode_gather_fee_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_gm = sp.local("rlp_gm", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bmc_fee_message", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_gm.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                rlp_gm.value = sp.view("decode_list", self.data.helper, decode_len,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecodingFeeMessage"
        temp_str = sp.local("str_gather", "")
        _svcs = sp.local("_svcs", {}, sp.TMap(sp.TNat, sp.TString))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_gm.value
            temp_byt = sp.local("byt4", sp.bytes("0x"))
            counter = sp.local("counter_gather", 0)
            sp.for c in rlp_.items():
                sp.if counter.value == 1:
                    temp_byt.value = c.value
                sp.if counter.value == 0:
                    temp_str.value = sp.view("decode_string", self.data.helper, c.value, t=sp.TString).open_some()
                counter.value = counter.value + 1
            sub_list = sp.local("sub_list", temp_byt.value)
            nsl_gm = sp.local("nsl_gm", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl_gm.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
                is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
                with sp.if_(is_list_lambda):
                    nsl_gm.value = sp.view("decode_list", self.data.helper, decode_len,
                                           t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                with sp.else_():
                    is_error.value = "ErrorInDecodingFeeMessage"
            with sp.if_(is_error.value == "Success"):
                new_sub_list = nsl_gm.value

                counter.value = 0
                sp.for x in new_sub_list.items():
                    _svcs.value[counter.value] = sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some()
                    counter.value = counter.value + 1
        return sp.record(fee_decode_rv = sp.record(fa=temp_str.value,
                         svcs=_svcs.value), status = is_error.value)

    def to_message_event(self, rlp):
        rlp_me = sp.local("rlp_me", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        is_error = sp.local("error_in_Event_message", sp.string("Success"))
        with sp.if_(is_list_lambda):
            rlp_me.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                rlp_me.value = sp.view("decode_list", self.data.helper, decode_len,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecoding"
        rv1 = sp.local("rv1_event", "")
        rv2 = sp.local("rv2_event", sp.nat(0))
        rv3 = sp.local("rv3_event", sp.bytes("0x"))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_me.value
            counter = sp.local("counter_event", 0)
            sp.for i in rlp_.items():
                sp.if counter.value == 2:
                    rv3.value = i.value
                sp.if counter.value == 0:
                    rv1.value = sp.view("decode_string", self.data.helper, i.value, t=sp.TString).open_some()
                sp.if counter.value == 1:
                    rv2.value = Utils2.Int.of_bytes(i.value)
                counter.value = counter.value + 1
            rv3.value = sp.view("without_length_prefix", self.data.helper, rv3.value, t=sp.TBytes).open_some()
        return sp.record(event_rv=sp.record(next_bmc= rv1.value, seq= rv2.value, message = rv3.value), status = is_error.value)

    def decode_receipt_proof(self, rlp):
        rlp_rp = sp.local("rlp_rp", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        is_error = sp.local("error_in_receipt_proof", sp.string("Success"))
        with sp.if_(is_list_lambda):
            rlp_rp.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                rlp_rp.value = sp.view("decode_list", self.data.helper, decode_len,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecoding"
        rv_int = sp.local("rv_int_receipt", 0)
        rv_int2 = sp.local("rv_int2_receipt", 0)
        events = sp.local("events_receipt", sp.map({}, tkey=sp.TNat,
                                                   tvalue=sp.TRecord(next_bmc=sp.TString,
                                                                     seq=sp.TNat,message=sp.TBytes)))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_rp.value
            temp_byt = sp.local("byt_receipt", sp.bytes("0x"))
            counter = sp.local("counter", 0)
            sp.for i in rlp_.items():
                sp.if counter.value == 1:
                    temp_byt.value = sp.view("without_length_prefix", self.data.helper, i.value, t=sp.TBytes).open_some()
                sp.if counter.value == 0:
                    rv_int.value = Utils2.Int.of_bytes(i.value)
                sp.if counter.value == 2:
                    wl_prefix = sp.view("without_length_prefix", self.data.helper, i.value, t=sp.TBytes).open_some()
                    rv_int2.value =Utils2.Int.of_bytes(wl_prefix)
                counter.value = counter.value + 1

            sub_list = sp.local("sub_list", temp_byt.value)

            nsl_rp = sp.local("nsl_rp", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl_rp.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                        t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
                is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
                with sp.if_(is_list_lambda):
                    nsl_rp.value = sp.view("decode_list", self.data.helper, decode_len,
                                            t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
                with sp.else_():
                    is_error.value = "ErrorInDecoding"
            with sp.if_(is_error.value == "Success"):
                new_sub_list = nsl_rp.value
                counter.value = 0

                sp.for z in new_sub_list.items():
                    from_event = self.to_message_event(z.value)
                    with sp.if_(from_event.status == "Success"):
                        events.value[counter.value] = from_event.event_rv
                        counter.value = counter.value + 1
                    with sp.else_():
                        is_error.value = "ErrorInDecoding"
        return sp.record(rv = sp.record(index = rv_int.value, events = events.value, height = rv_int2.value), status = is_error.value)

    def decode_receipt_proofs(self, rlp):
        sp.set_type(rlp, sp.TBytes)
        is_error = sp.local("error_in_receipt_proofs", sp.string("Success"))
        rlp_rps = sp.local("rlp_rps", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_rps.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                rlp_rps.value = sp.view("decode_list", self.data.helper, decode_len, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecoding"
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
        sub_list = sp.local("sub_list_proofs", temp_byt.value)

        nsl_rps = sp.local("nsl_rps", sp.map(tkey=sp.TNat))
        is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            nsl_rps.value = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            decode_len = sp.view("without_length_prefix", self.data.helper, sub_list.value, t=sp.TBytes).open_some()
            is_list_lambda = sp.view("is_list", self.data.helper, decode_len, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                nsl_rps.value = sp.view("decode_list", self.data.helper, decode_len,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecoding"
        new_sub_list = nsl_rps.value
        counter.value = 0
        sp.if sp.len(new_sub_list) > 0:
            sp.for x in new_sub_list.items():
                from_receipt_proofs = self.decode_receipt_proof(x.value)
                with sp.if_(from_receipt_proofs.status == "Success"):
                    receipt_proofs.value[counter.value] = from_receipt_proofs.rv
                    counter.value = counter.value + 1
                with sp.else_():
                    is_error.value = "ErrorInDecoding"
        return sp.record(receipt_proof = receipt_proofs.value, status = is_error.value)

    # rlp encoding starts here

    def encode_bmc_service(self, params):
        sp.set_type(params, sp.TRecord(serviceType=sp.TString, payload=sp.TBytes))

        encode_service_type = sp.view("encode_string", self.data.helper, params.serviceType, t=sp.TBytes).open_some()

        payload_rlp = sp.view("encode_list", self.data.helper, [params.payload], t=sp.TBytes).open_some()
        payload_rlp = sp.view("with_length_prefix", self.data.helper, payload_rlp, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_service_type, payload_rlp],
                                        t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix,
                                        t=sp.TBytes).open_some()
        return rlp_bytes_with_prefix

    def encode_bmc_message(self, params):
        sp.set_type(params, sp.TRecord(src=sp.TString, dst=sp.TString, svc=sp.TString, sn=sp.TInt, message=sp.TBytes))

        rlp = sp.local("rlp_sn", sp.bytes("0x"))
        encode_src = sp.view("encode_string", self.data.helper, params.src, t=sp.TBytes).open_some()
        encode_dst = sp.view("encode_string", self.data.helper, params.dst, t=sp.TBytes).open_some()
        encode_svc = sp.view("encode_string", self.data.helper, params.svc, t=sp.TBytes).open_some()
        _to_byte = sp.view("to_byte", self.data.helper_parse_negative, params.sn, t=sp.TBytes).open_some()
        rlp.value = _to_byte

        sp.if params.sn < sp.int(0):
            encode_sn = sp.view("with_length_prefix", self.data.helper, rlp.value, t=sp.TBytes).open_some()
            rlp.value = encode_sn

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper,
                                        [encode_src, encode_dst, encode_svc, rlp.value, params.message],
                                        t=sp.TBytes).open_some()
        return rlp_bytes_with_prefix

    def encode_response(self, params):
        sp.set_type(params, sp.TRecord(code=sp.TNat, message=sp.TString))

        encode_code = sp.view("encode_nat", self.data.helper, params.code, t=sp.TBytes).open_some()
        encode_message = sp.view("encode_string", self.data.helper, params.message, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_code, encode_message],
                                        t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix,
                                              t=sp.TBytes).open_some()
        return final_rlp_bytes_with_prefix

