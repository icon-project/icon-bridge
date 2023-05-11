import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")


class EncodeLibrary:
    def encode_service_message(self, params):
        sp.set_type(params, types.Types.ServiceMessage)

        encode_service_type = sp.view("of_string", self.data.helper, params.serviceType, t=sp.TBytes).open_some()
        rlp = params.data + encode_service_type
        rlp_bytes_with_prefix = sp.view("with_length_prefix", self.data.helper, rlp, t=sp.TBytes).open_some()
        return rlp_bytes_with_prefix


    def encode_transfer_coin_msg(self, data):
        sp.set_type(data, types.Types.TransferCoin)

        rlp = sp.local("rlp", sp.bytes("0x80"))
        temp = sp.local("temp", sp.bytes("0x80"))

        sp.for i in sp.range(0, sp.len(data.assets)):
            temp.value = sp.view("of_string", self.data.helper, data.assets.get(i, default_value=sp.record(coin_name="",value=sp.nat(0))).coin_name, t=sp.TBytes).open_some()\
                         + sp.view("of_nat", self.data.helper, data.assets.get(i, default_value=sp.record(coin_name="",value=sp.nat(0))).value, t=sp.TBytes).open_some()
            rlp.value = rlp.value + sp.view("with_length_prefix", self.data.helper, temp.value, t=sp.TBytes).open_some()

        rlp.value = sp.view("of_string", self.data.helper, data.from_addr, t=sp.TBytes).open_some() + \
                    sp.view("of_string", self.data.helper, data.to, t=sp.TBytes).open_some() + \
                    sp.view("with_length_prefix", self.data.helper, rlp.value, t=sp.TBytes).open_some()

        return sp.view("with_length_prefix", self.data.helper, rlp.value, t=sp.TBytes).open_some()
