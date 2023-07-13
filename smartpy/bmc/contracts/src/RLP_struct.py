import smartpy as sp

Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")
types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class DecodeEncodeLibrary:

    def decode_bmc_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_bmc_message = sp.local("rlp_decode_bmc_message", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_decode_bmc_message", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bmc_message.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInBMCMessageDecoding"
        return_value_map = sp.compute(sp.map(tkey=sp.TString, tvalue=sp.TString))
        sn_no = sp.local("sn_no_decode_bmc_message", 0)
        decoded_message = sp.local("message_decode_bmc_message", sp.bytes("0x"))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_bmc_message.value
            counter = sp.local("counter", 0)
            sp.for k in rlp_.items():
                with sp.if_ (counter.value == 0):
                    return_value_map["src"] = sp.view("decode_string", self.data.helper, k.value, t=sp.TString).open_some()
                with sp.if_ (counter.value == 1):
                    return_value_map["dst"] = sp.view("decode_string", self.data.helper, k.value, t=sp.TString).open_some()
                with sp.if_ (counter.value == 2):
                    return_value_map["svc"] = sp.view("decode_string", self.data.helper, k.value, t=sp.TString).open_some()
                with sp.if_ (counter.value == 3):
                    sn_in_bytes = sp.view("without_length_prefix", self.data.helper, k.value, t=sp.TBytes).open_some()
                    _to_int = sp.view("to_int", self.data.helper_parse_negative, sn_in_bytes, t=sp.TInt).open_some()
                    sn_no.value = _to_int
                with sp.if_ (counter.value == 4):
                    decoded_message.value = sp.view("without_length_prefix", self.data.helper, k.value,
                                             t=sp.TBytes).open_some()
                counter.value = counter.value + 1
        return sp.record(bmc_dec_rv = sp.record(src=return_value_map.get("src", default_value = "NoKey"),
                         dst=return_value_map.get("dst", default_value = "NoKey"),
                         svc=return_value_map.get("svc", default_value = "NoKey"),
                         sn=sn_no.value,
                         message=decoded_message.value),
                         status = is_error.value)

    def decode_response(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        code_num = sp.local("code_decode_response", sp.nat(0))
        rlp_ = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        counter = sp.local("counter_decode_response", 0)
        sp.for m in rlp_.items():
            with sp.if_ (counter.value == 0):
                code_num.value = Utils2.Int.of_bytes(m.value)
            counter.value = counter.value + 1

            # message in case of error is null which cannot be decoded into string
        return sp.record(code=code_num.value, message="Error")

    def decode_init_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_init = sp.local("rlp_init_message", sp.map(tkey=sp.TNat))
        rlp_init.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        _links = sp.local("links_init_message", [], sp.TList(sp.TString))
        rlp_ = rlp_init.value
        counter = sp.local("counter_init_message", 0)
        bytes_message = sp.local("byte_init_message", sp.bytes("0x"))
        sp.for g in rlp_.items():
            with sp.if_ (counter.value == 0):
                bytes_message.value = g.value
            counter.value = counter.value + 1

        sub_list = sp.local("sub_list_init_message", bytes_message.value)
        nsl_init = sp.local("nsl_init_message", sp.map(tkey=sp.TNat))
        nsl_init.value = sp.view("decode_list", self.data.helper, sub_list.value,
                               t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = nsl_init.value

        counter.value = 0
        sp.for x in new_sub_list.items():
            _links.value.push(sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some())
            counter.value = counter.value + 1
        return sp.record(links_list = _links.value)

    def decode_bmc_service(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_bmc_message = sp.local("rlp_decode_bmc_service", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_decode_bmc_service", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_bmc_message.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInDecodingBMCService"

        service_type = sp.local("str_value_decode_bmc_service", "")
        payload = sp.local("byt_value_decode_bmc_service", sp.bytes("0x"))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_bmc_message.value
            counter = sp.local("counter_service", 0)
            sp.for b in rlp_.items():
                with sp.if_ (counter.value == 0):
                    service_type.value = sp.view("decode_string", self.data.helper, b.value, t=sp.TString).open_some()
                with sp.if_ (counter.value == 1):
                    payload.value = sp.view("without_length_prefix", self.data.helper, b.value,
                                             t=sp.TBytes).open_some()
                counter.value = counter.value + 1
        return sp.record(bmc_service_rv=sp.record(serviceType=service_type.value,
                         payload=payload.value), status = is_error.value)

    def decode_gather_fee_message(self, rlp):
        sp.set_type(rlp, sp.TBytes)

        rlp_gather_message = sp.local("rlp_gather_message", sp.map(tkey=sp.TNat))
        is_error = sp.local("error_in_bmc_fee_message", sp.string("Success"))
        is_list_lambda = sp.view("is_list", self.data.helper, rlp, t=sp.TBool).open_some()
        with sp.if_(is_list_lambda):
            rlp_gather_message.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        with sp.else_():
            is_error.value = "ErrorInDecodingFeeMessage"
        fa = sp.local("fee_aggregator_gather_message", "")
        _svcs = sp.local("_svcs", {}, sp.TMap(sp.TNat, sp.TString))
        with sp.if_(is_error.value == "Success"):
            rlp_ = rlp_gather_message.value
            byte_message = sp.local("byte_gather_fee", sp.bytes("0x"))
            counter = sp.local("counter_gather_message", 0)
            sp.for c in rlp_.items():
                with sp.if_ (counter.value == 1):
                    byte_message.value = c.value
                with sp.if_ (counter.value == 0):
                    fa.value = sp.view("decode_string", self.data.helper, c.value, t=sp.TString).open_some()
                counter.value = counter.value + 1
            sub_list = sp.local("sub_list_gather_message", byte_message.value)
            new_sub_list_gather = sp.local("new_sub_list_gather_message", sp.map(tkey=sp.TNat))
            is_list_lambda = sp.view("is_list", self.data.helper, sub_list.value, t=sp.TBool).open_some()
            with sp.if_(is_list_lambda):
                new_sub_list_gather.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                       t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
            with sp.else_():
                is_error.value = "ErrorInDecodingFeeMessage"
            with sp.if_(is_error.value == "Success"):
                new_sub_list = new_sub_list_gather.value

                counter.value = 0
                sp.for x in new_sub_list.items():
                    _svcs.value[counter.value] = sp.view("decode_string", self.data.helper, x.value, t=sp.TString).open_some()
                    counter.value = counter.value + 1
        return sp.record(fee_decode_rv = sp.record(fa=fa.value,
                         svcs=_svcs.value), status = is_error.value)

    def to_message_event(self, rlp):
        rlp_message_event = sp.local("rlp_message_event", sp.map(tkey=sp.TNat))
        rlp_message_event.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        next_bmc = sp.local("next_bmc_message_event", "")
        seq = sp.local("seq_bmc_message_event", sp.nat(0))
        message = sp.local("message_bmc_message_event", sp.bytes("0x"))
        rlp_ = rlp_message_event.value
        counter = sp.local("counter_message_event", 0)
        sp.for i in rlp_.items():
            with sp.if_ (counter.value == 2):
                message.value = sp.view("without_length_prefix", self.data.helper, i.value, t=sp.TBytes).open_some()
            with sp.if_ (counter.value == 0):
                next_bmc.value = sp.view("decode_string", self.data.helper, i.value, t=sp.TString).open_some()
            with sp.if_ (counter.value == 1):
                seq.value = Utils2.Int.of_bytes(i.value)
            counter.value = counter.value + 1
        return sp.record(event_rv=sp.record(next_bmc= next_bmc.value, seq= seq.value, message = message.value))

    def decode_receipt_proof(self, rlp):
        rlp_receipt_proof = sp.local("rlp_receipt_proof", sp.map(tkey=sp.TNat))
        without_prefix = sp.view("without_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
        rlp_receipt_proof.value = sp.view("decode_list", self.data.helper, without_prefix, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        index = sp.local("index_receipt_proof", 0)
        height = sp.local("height_receipt_proof", 0)
        events = sp.local("events_receipt_proof", sp.map({}, tkey=sp.TNat,
                                                   tvalue=sp.TRecord(next_bmc=sp.TString,
                                                                     seq=sp.TNat,message=sp.TBytes)))
        rlp_ = rlp_receipt_proof.value
        byte_message_receipt_proof = sp.local("byte_message_receipt_proof", sp.bytes("0x"))
        counter = sp.local("counter", 0)
        sp.for i in rlp_.items():
            with sp.if_ (counter.value == 1):
                byte_message_receipt_proof.value = sp.view("without_length_prefix", self.data.helper, i.value, t=sp.TBytes).open_some()
            with sp.if_ (counter.value == 0):
                index.value = Utils2.Int.of_bytes(i.value)
            with sp.if_ (counter.value == 2):
                wl_prefix = sp.view("without_length_prefix", self.data.helper, i.value, t=sp.TBytes).open_some()
                height.value =Utils2.Int.of_bytes(wl_prefix)
            counter.value = counter.value + 1

        sub_list = sp.local("sub_list_receipt_proof", byte_message_receipt_proof.value)

        new_sub_list_receipt_proof = sp.local("new_sub_list_receipt_proof", sp.map(tkey=sp.TNat))
        new_sub_list_receipt_proof.value = sp.view("decode_list", self.data.helper, sub_list.value,
                                t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = new_sub_list_receipt_proof.value
        counter.value = 0

        sp.for z in new_sub_list.items():
            from_event = self.to_message_event(z.value)
            events.value[counter.value] = from_event.event_rv
            counter.value = counter.value + 1
        return sp.record(rv = sp.record(index = index.value, events = events.value, height = height.value))

    def decode_receipt_proofs(self, rlp):
        sp.set_type(rlp, sp.TBytes)
        rlp_receipt_proofs = sp.local("rlp_receipt_proofs", sp.map(tkey=sp.TNat))
        rlp_receipt_proofs.value = sp.view("decode_list", self.data.helper, rlp, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        rlp_ = rlp_receipt_proofs.value
        counter = sp.local("counter_receipt_proofs", 0)
        receipt_proofs = sp.local("events_receipt_proofs", sp.map({}, tkey=sp.TNat,
                tvalue=sp.TRecord(index = sp.TNat,
                          events = sp.TMap(sp.TNat, sp.TRecord(next_bmc=sp.TString, seq=sp.TNat, message=sp.TBytes)),
                          height = sp.TNat,
                        )
                                                                   )
                                  )
        message_byte_receipt_proofs = sp.local("message_byte_receipt_proofs", sp.bytes("0x"))
        sp.for i in rlp_.items():
            with sp.if_ (counter.value == 0):
                message_byte_receipt_proofs.value = i.value
            counter.value = counter.value + 1
        sub_list = sp.local("sub_list_receipt_proofs", message_byte_receipt_proofs.value)

        new_sub_list_receipt_proofs = sp.local("new_sub_list_receipt_proofs", sp.map(tkey=sp.TNat))
        new_sub_list_receipt_proofs.value = sp.view("decode_list", self.data.helper, sub_list.value, t=sp.TMap(sp.TNat, sp.TBytes)).open_some()
        new_sub_list = new_sub_list_receipt_proofs.value
        counter.value = 0
        with sp.if_ (sp.len(new_sub_list) > 0):
            sp.for x in new_sub_list.items():
                from_receipt_proofs = self.decode_receipt_proof(x.value)
                receipt_proofs.value[counter.value] = from_receipt_proofs.rv
                counter.value = counter.value + 1
        return sp.record(receipt_proof = receipt_proofs.value)

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

        with sp.if_ (params.sn < sp.int(0)):
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

