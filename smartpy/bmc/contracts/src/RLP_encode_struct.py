import smartpy as sp

helper_file = sp.io.import_script_from_url("file:./contracts/src/helper.py")


class EncodeLibrary(sp.Contract):

    def __init__(self, helper_contract, helper_negative_address):
        self.init(
            helper=helper_contract,
            helper_parse_negative=helper_negative_address
        )

    @sp.onchain_view()
    def encode_bmc_service(self, params):
        sp.set_type(params, sp.TRecord(serviceType=sp.TString, payload=sp.TBytes))

        encode_service_type = sp.view("encode_string", self.data.helper, params.serviceType, t=sp.TBytes).open_some()

        payload_rlp = sp.view("encode_list", self.data.helper, [params.payload], t=sp.TBytes).open_some()
        payload_rlp = sp.view("with_length_prefix", self.data.helper, payload_rlp, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_service_type, payload_rlp],
                                        t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix,
                                        t=sp.TBytes).open_some()
        sp.result(rlp_bytes_with_prefix)

    @sp.onchain_view()
    def encode_bmc_message(self, params):
        sp.set_type(params, sp.TRecord(src=sp.TString, dst=sp.TString, svc=sp.TString, sn=sp.TInt, message=sp.TBytes))

        encode_src = sp.view("encode_string", self.data.helper, params.src, t=sp.TBytes).open_some()
        encode_dst = sp.view("encode_string", self.data.helper, params.dst, t=sp.TBytes).open_some()
        encode_svc = sp.view("encode_string", self.data.helper, params.svc, t=sp.TBytes).open_some()
        encode_sn = sp.view("to_byte", self.data.helper_parse_negative, params.sn, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_src, encode_dst, encode_svc, encode_sn, params.message], t=sp.TBytes).open_some()
        sp.result(rlp_bytes_with_prefix)

    @sp.onchain_view()
    def encode_response(self, params):
        sp.set_type(params, sp.TRecord(code=sp.TNat, message=sp.TString))

        encode_code = sp.view("encode_nat", self.data.helper, params.code, t=sp.TBytes).open_some()
        encode_message = sp.view("encode_string", self.data.helper, params.message, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_code, encode_message], t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix, t=sp.TBytes).open_some()
        sp.result(final_rlp_bytes_with_prefix)


@sp.add_test(name="Encoder")
def test():
    helper_nev = sp.test_account("Helper Negative")
    scenario = sp.test_scenario()

    helper = helper_file.Helper()
    scenario += helper

    c1 = EncodeLibrary(helper.address, helper_nev.address)
    scenario += c1


sp.add_compilation_target("RLP_encode_struct", EncodeLibrary(helper_contract=sp.address("KT1HwFJmndBWRn3CLbvhUjdupfEomdykL5a6"),
                                                             helper_negative_address=sp.address("KT1DHptHqSovffZ7qqvSM9dy6uZZ8juV88gP")))