import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class EncodeLibrary:
    def encode_service_message(self, params):
        sp.set_type(params, sp.TRecord(service_type_value = sp.TNat, data = sp.TBytes))

        encode_service_type = sp.view("encode_nat", self.data.helper, params.service_type_value, t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_service_type, params.data], t=sp.TBytes).open_some()

        return rlp_bytes_with_prefix


    def encode_transfer_coin_msg(self, data):
        sp.set_type(data, types.Types.TransferCoin)

        rlp = sp.local("rlp", sp.bytes("0x"))
        temp = sp.local("temp", sp.bytes("0x"))
        coin_name = sp.local("coin_name", sp.bytes("0x"))

        sp.for i in sp.range(0, sp.len(data.assets)):
            coin_name.value = sp.view("encode_string", self.data.helper, data.assets.get(i, default_value=sp.record(coin_name="",value=sp.nat(0))).coin_name, t=sp.TBytes).open_some()
            temp.value =  sp.view("encode_nat", self.data.helper, data.assets.get(i, default_value=sp.record(coin_name="",value=sp.nat(0))).value, t=sp.TBytes).open_some()
            rlp.value = sp.view("encode_list", self.data.helper, [rlp.value, coin_name.value, temp.value], t=sp.TBytes).open_some()

        from_addr_encoded = sp.view("encode_string", self.data.helper, data.from_addr, t=sp.TBytes).open_some()
        to_addr_encoded = sp.view("encode_string", self.data.helper, data.to, t=sp.TBytes).open_some()
        rlp.value = sp.view("encode_list", self.data.helper, [from_addr_encoded, to_addr_encoded, rlp.value], t=sp.TBytes).open_some()

        return rlp.value
