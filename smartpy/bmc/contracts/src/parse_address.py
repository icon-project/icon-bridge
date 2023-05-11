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


@sp.add_test(name="Conversion")
def test():
    c1 = ParseAddress()
    scenario = sp.test_scenario()
    scenario.h1("Conversion")
    scenario += c1
    c1.add_to_str(sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"))
    scenario.verify(c1.add_to_str(sp.address("KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr")) == "KT1FfkTSts5DnvyJp2qZbPMeqm2XpMYES7Vr")
    scenario.verify(c1.add_to_str(sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")) == "tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")

sp.add_compilation_target("parse_address", ParseAddress())