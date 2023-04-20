import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
Utils = sp.io.import_script_from_url(
    "https://raw.githubusercontent.com/Acurast/acurast-hyperdrive/main/contracts/tezos/libs/utils.py")


def encode_service_message(params):
    sp.set_type(params, types.Types.ServiceMessage)

    encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
    rlp = params.data + encode_string_packed(params.serviceType)
    with_length_prefix = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)
    rlp_bytes_with_prefix = with_length_prefix(rlp)
    return rlp_bytes_with_prefix


def encode_transfer_coin_msg(data):
    sp.set_type(data, types.Types.TransferCoin)

    rlp = sp.local("rlp", sp.bytes("0x80"))
    temp = sp.local("temp", sp.bytes("0x80"))
    encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
    encode_nat_packed = sp.build_lambda(Utils.Bytes.of_nat)
    with_length_prefix = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)

    sp.for i in sp.range(0, sp.len(data.assets)):
        temp.value = encode_string_packed(data.assets[i].coin_name) + encode_nat_packed(data.assets[i].value)
        rlp.value = rlp.value + with_length_prefix(temp.value)

    rlp.value = encode_string_packed(data.from_addr) + encode_string_packed(data.to) + with_length_prefix(rlp.value)

    return with_length_prefix(rlp.value)
