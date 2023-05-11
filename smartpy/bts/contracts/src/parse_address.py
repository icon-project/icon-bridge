import smartpy as sp
Utils = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")


class ParseAddress(sp.Contract):
    tz_prefixes = sp.map({
        sp.bytes('0x0000'): sp.string('tz1'),
        sp.bytes('0x0001'): sp.string('tz2'),
        sp.bytes('0x0002'): sp.string('tz3'),
        sp.bytes('0x0003'): sp.string('tz4')
    })
    base58_encodings = sp.list([
        sp.map({"prefix": "tz1", "elem1": "6", "elem2": "161", "elem3": "159", "len": "20"}),
        sp.map({"prefix": "tz2", "elem1": "6", "elem2": "161", "elem3": "161", "len": "20"}),
        sp.map({"prefix": "tz3", "elem1": "6", "elem2": "161", "elem3": "164", "len": "20"}),
        sp.map({"prefix": "tz4", "elem1": "6", "elem2": "161", "elem3": "16", "len": "20"}),
        sp.map({"prefix": "KT1", "elem1": "2", "elem2": "90", "elem3": "121", "len": "20"}),
    ])

    def __init__(self):
        self.init()

    def unforge_address(self, data):
        """Decode address or key_hash from bytes.

        :param data: encoded address or key_hash
        :returns: base58 encoded address
        """
        sp.set_type(data, sp.TBytes)
        byt = sp.slice(data, 6, 22).open_some()
        prefix = sp.slice(byt, 0, 2).open_some()
        starts_with = sp.slice(byt, 0, 1).open_some()
        ends_with = sp.slice(byt, 21, 1).open_some()
        sliced_byte = sp.slice(byt, 1, 20).open_some()
        local_byte = sp.local("local_byte", sp.bytes("0x"))
        return_value = sp.local("return_value", "tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg", sp.TString)

        sp.for item in self.tz_prefixes.items():
            sp.if item.key == prefix:
                return_value.value = self.base58_encode(sp.slice(byt, 2, 20).open_some(), Utils.Bytes.of_string(item.value), local_byte)

        sp.if (starts_with == sp.bytes("0x01")) & (ends_with == sp.bytes("0x00")):
            return_value.value = self.base58_encode(sliced_byte, Utils.Bytes.of_string("KT1"), local_byte)
        sp.if (starts_with == sp.bytes("0x02")) & (ends_with == sp.bytes("0x00")):
            return_value.value = self.base58_encode(sliced_byte, Utils.Bytes.of_string("txr1"), local_byte)
        sp.if (starts_with == sp.bytes("0x03")) & (ends_with == sp.bytes("0x00")):
            return_value.value = self.base58_encode(sliced_byte, Utils.Bytes.of_string("sr1"), local_byte)

        return return_value.value

    def tb(self, _list):
        byte_str = sp.local("byte_str", sp.bytes("0x"))
        sp.for num in _list:
            byte_str.value += Utils.Bytes.of_nat(num)
        return byte_str.value

    def base58_encode(self, v, prefix, _byte):
        """
        Encode data using Base58 with checksum and add an according binary prefix in the end.
        :param v: Array of bytes
        :param prefix: Human-readable prefix (use b'') e.g. b'tz', b'KT', etc
        :param local_byte: local variable

        :returns: bytes (use string.decode())
        """
        length_v = sp.to_int(sp.len(v))
        encoding = sp.local("encode", sp.map({}))
        byte_from_tbl = sp.local("byte_from_tbl", sp.bytes("0x"))
        byte_value = _byte

        sp.for enc in self.base58_encodings:
            sp.if (length_v == Utils.Int.of_string(enc["len"])) & (prefix == Utils.Bytes.of_string(enc["prefix"])):
                encoding.value = enc
                byte_from_tbl.value = self.tb([sp.as_nat(Utils.Int.of_string(enc["elem1"])),
                                          sp.as_nat(Utils.Int.of_string(enc["elem2"])),
                                          sp.as_nat(Utils.Int.of_string(enc["elem3"]))])
        sha256_encoding = sp.sha256(sp.sha256(byte_from_tbl.value + v))
        sha256_encoding = byte_from_tbl.value + v + sp.slice(sha256_encoding, 0, 4).open_some()
        acc = sp.local("for_while_loop", Utils.Int.of_bytes(sha256_encoding))
        alphabet = Utils.Bytes.of_string("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
        base = 58
        sp.while acc.value > 0:
            (acc.value, idx) = sp.match_pair(sp.ediv(acc.value, base).open_some())
            byte_value.value = sp.slice(alphabet, idx, 1).open_some() + byte_value.value

        return sp.unpack(sp.bytes("0x050100000024") + byte_value.value, sp.TString).open_some()

    @sp.onchain_view()
    def add_to_str(self, params):
        sp.set_type(params, sp.TAddress)
        sp.result(self.unforge_address(sp.pack(params)))

    def _to_addr(self, params):
        string_in_bytes = sp.pack(params)
        actual_prefix = sp.local("actual_prefix", "")
        addr = sp.local("addr", sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"))
        alphabet = Utils.Bytes.of_string("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
        sp.trace(":kljlk")
        sp.if sp.len(string_in_bytes) == sp.nat(42):
            sp.trace("lsdnglknslnd")
            string_in_bytes = sp.slice(string_in_bytes, 6, 36).open_some()

            element_list = sp.range(0, sp.len(alphabet), 1)
            temp_map = sp.local("temp_map", {})
            temp_var = sp.local("y", 0)
            sp.trace("ele")
            sp.for elem in element_list:
                temp_var.value = elem
                temp_map.value[sp.slice(alphabet, temp_var.value, 1).open_some()] = temp_var.value
            decimal = sp.local("decimal", 0)
            base = sp.len(alphabet)
            element_list_2 = sp.range(0, sp.len(string_in_bytes), 1)
            sp.trace("ele1")
            sp.for elem in element_list_2:
                decimal.value = decimal.value * base + temp_map.value[sp.slice(string_in_bytes, elem, 1).open_some()]
            byt_value = Utils.Bytes.of_nat(sp.as_nat(sp.to_int(decimal.value)))
            new_byt_value = sp.slice(byt_value, 0, sp.as_nat(sp.len(byt_value) - 4)).open_some()
            prefix = sp.slice(new_byt_value, 0, 3).open_some()
            prefix_len = sp.range(0, sp.len(prefix), 1)
            temp_var3 = sp.local("z", 0)
            list_string = sp.local("list_string", [])
            sp.for x in prefix_len:
                temp_var3.value = x
                list_string.value.push(Utils.Int.of_bytes(sp.slice(prefix, temp_var3.value, 1).open_some()))
            v = sp.slice(new_byt_value, 3, sp.as_nat(sp.len(new_byt_value) - 3))
            byte_local = sp.local("byt_old", sp.bytes("0x"))
            sp.trace("ele2")
            sp.for enc in self.base58_encodings:
                byte_local.value = self.tb([sp.as_nat(Utils.Int.of_string(enc["elem1"])),
                                       sp.as_nat(Utils.Int.of_string(enc["elem2"])),
                                       sp.as_nat(Utils.Int.of_string(enc["elem3"]))])
                sp.if byte_local.value == prefix:
                    actual_prefix.value = enc["prefix"]
            sp.trace("ele3")
            sp.for item in self.tz_prefixes.items():
                sp.trace(item.value)
                sp.trace("k")
                sp.trace(actual_prefix.value)
                sp.if item.value == actual_prefix.value:
                    sp.trace("in if")
                    decoded_address = sp.unpack(sp.bytes("0x050a00000016") + item.key + v.open_some(),
                                                sp.TAddress)
                    addr.value = decoded_address.open_some()
                    # return addr.value
            sp.trace("actual")
            sp.trace(actual_prefix.value)
            sp.if actual_prefix.value == "KT1":
                sp.trace("KTSSSS")
                decoded_address = sp.unpack(
                    sp.bytes("0x050a00000016") + sp.bytes("0x01") + v.open_some() + sp.bytes("0x00"),
                    sp.TAddress)
                addr.value = decoded_address.open_some()
            sp.if actual_prefix.value == "txr1":
                decoded_address = sp.unpack(
                    sp.bytes("0x050a00000016") + sp.bytes("0x02") + v.open_some() + sp.bytes("0x00"),
                    sp.TAddress)
                addr.value = decoded_address.open_some()
            sp.if actual_prefix.value == "sr1":
                decoded_address = sp.unpack(
                    sp.bytes("0x050a00000016") + sp.bytes("0x03") + v.open_some() + sp.bytes("0x00"),
                    sp.TAddress)
                addr.value = decoded_address.open_some()

            sp.if actual_prefix.value == "":
                sp.trace("in else")
                addr.value = sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")

        with sp.else_():
            addr.value =  sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")

        return addr.value

    @sp.onchain_view()
    def str_to_addr(self, params):
        sp.set_type(params, sp.TString)
        sp.result(self._to_addr(params))

    @sp.onchain_view()
    def string_of_int(self, params):
        sp.set_type(params, sp.TInt)
        sp.result(Utils.String.of_int(params))


@sp.add_test(name="Conversion")
def test():
    alice=sp.test_account("Alice")
    c1 = ParseAddress()
    scenario = sp.test_scenario()
    scenario.h1("Conversion")
    scenario += c1
    c1.add_to_str(sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"))
    scenario.verify(c1.add_to_str(sp.address("KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr")) == "KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr")
    scenario.verify(c1.add_to_str(sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")) == "tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")

    c1.str_to_addr("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")
    c1.str_to_addr("KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr")
    scenario.verify(c1.str_to_addr("KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr") == sp.address("KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr"))
    scenario.verify(c1.str_to_addr("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP") == sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"))
    # invalid address
    scenario.verify(c1.str_to_addr("tz1g3pJZPifxhN") == sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg"))



sp.add_compilation_target("parse_address", ParseAddress())


