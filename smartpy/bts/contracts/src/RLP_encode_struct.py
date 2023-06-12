import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
helper_file = sp.io.import_script_from_url("file:./contracts/src/helper.py")


class EncodeLibrary(sp.Contract):

    def __init__(self, helper_contract):
        self.init(
            helper=helper_contract
        )


    @sp.onchain_view()
    def encode_service_message(self, params):
        sp.set_type(params, sp.TRecord(service_type_value = sp.TNat, data = sp.TBytes))

        encode_service_type = sp.view("encode_nat", self.data.helper, params.service_type_value, t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_service_type, params.data], t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix, t=sp.TBytes).open_some()

        sp.result(final_rlp_bytes_with_prefix)


    @sp.onchain_view()
    def encode_transfer_coin_msg(self, data):
        sp.set_type(data, types.Types.TransferCoin)

        rlp = sp.local("rlp", sp.bytes("0x"))
        rlp_list = sp.local("rlp_list", [], t=sp.TList(sp.TBytes))
        temp = sp.local("temp", sp.bytes("0x"))
        coin_name = sp.local("coin_name", sp.bytes("0x"))
        sp.for i in sp.range(0, sp.len(data.assets)):
            coin_name.value = sp.view("encode_string", self.data.helper, data.assets.get(i, default_value=sp.record(coin_name="",value=sp.nat(0))).coin_name, t=sp.TBytes).open_some()
            temp.value =  sp.view("encode_nat", self.data.helper, data.assets.get(i, default_value=sp.record(coin_name="",value=sp.nat(0))).value, t=sp.TBytes).open_some()
            rlp_list.value.push(sp.view("encode_list", self.data.helper, [coin_name.value, temp.value], t=sp.TBytes).open_some())
            # rlp.value = sp.view("with_length_prefix", self.data.helper, rlp.value,
            #                                       t=sp.TBytes).open_some()

        assets_list = sp.view("encode_list", self.data.helper, rlp_list.value, t=sp.TBytes).open_some()
        from_addr_encoded = sp.view("encode_string", self.data.helper, data.from_addr, t=sp.TBytes).open_some()
        to_addr_encoded = sp.view("encode_string", self.data.helper, data.to, t=sp.TBytes).open_some()
        rlp.value = sp.view("encode_list", self.data.helper, [from_addr_encoded, to_addr_encoded, assets_list], t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp.value, t=sp.TBytes).open_some()

        sp.result(final_rlp_bytes_with_prefix)


    @sp.onchain_view()
    def encode_response(self, params):
        sp.set_type(params, sp.TRecord(code=sp.TNat, message=sp.TString))

        encode_code = sp.view("encode_nat", self.data.helper, params.code, t=sp.TBytes).open_some()
        encode_message = sp.view("encode_string", self.data.helper, params.message, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_code, encode_message],
                                        t=sp.TBytes).open_some()
        final_rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp_bytes_with_prefix, t=sp.TBytes).open_some()

        sp.result(final_rlp_bytes_with_prefix)


@sp.add_test(name="Encoder")
def test():
    scenario = sp.test_scenario()

    helper = helper_file.Helper()
    scenario += helper

    c1 = EncodeLibrary(helper.address)
    scenario += c1


sp.add_compilation_target("RLP_encode_struct",
                          EncodeLibrary(helper_contract=sp.address("KT1HwFJmndBWRn3CLbvhUjdupfEomdykL5a6")))
