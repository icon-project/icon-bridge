import smartpy as sp

Utils = sp.io.import_script_from_url("https://raw.githubusercontent.com/Acurast/acurast-hyperdrive/main/contracts/tezos/libs/utils.py")


class Helper(sp.Contract):
    def __init__(self):
        self.init()

    @sp.onchain_view()
    def decode_string(self, params):
        sp.set_type(params, sp.TBytes)
        decode_string = sp.build_lambda(Utils.RLP.Decoder.decode_string)
        sp.result(decode_string(params))

    @sp.onchain_view()
    def decode_list(self, params):
        sp.set_type(params, sp.TBytes)
        decode_list = sp.build_lambda(Utils.RLP.Decoder.decode_list)
        sp.result(decode_list(params))

    @sp.onchain_view()
    def of_string(self, params):
        sp.set_type(params, sp.TString)
        encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
        sp.result(encode_string_packed(params))

    @sp.onchain_view()
    def encode_string(self, params):
        sp.set_type(params, sp.TString)
        encode_string_packed = sp.build_lambda(Utils.RLP.Encoder.encode_string)
        sp.result(encode_string_packed(params))

    @sp.onchain_view()
    def encode_nat(self, params):
        sp.set_type(params, sp.TNat)
        encode_nat_packed = sp.build_lambda(Utils.RLP.Encoder.encode_nat)
        sp.result(encode_nat_packed(params))

    @sp.onchain_view()
    def of_nat(self, params):
        sp.set_type(params, sp.TNat)
        encode_nat_packed = sp.build_lambda(Utils.Bytes.of_nat)
        sp.result(encode_nat_packed(params))

    @sp.onchain_view()
    def with_length_prefix(self, params):
        sp.set_type(params, sp.TBytes)
        encode_length_packed = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)
        sp.result(encode_length_packed(params))

    @sp.onchain_view()
    def without_length_prefix(self, params):
        sp.set_type(params, sp.TBytes)
        decode = sp.build_lambda(Utils.RLP.Decoder.without_length_prefix)
        sp.result(decode(params))

    @sp.onchain_view()
    def encode_list(self, params):
        sp.set_type(params, sp.TList(sp.TBytes))
        encode_list_packed = sp.build_lambda(Utils.RLP.Encoder.encode_list)
        sp.result(encode_list_packed(params))


@sp.add_test(name="Helper")
def test():
    scenario = sp.test_scenario()
    helper = Helper()
    scenario += helper


sp.add_compilation_target("helper", Helper())