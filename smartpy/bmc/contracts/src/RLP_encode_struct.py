import smartpy as sp
Utils = sp.io.import_script_from_url("https://raw.githubusercontent.com/Acurast/acurast-hyperdrive/main/contracts/tezos/libs/utils.py")

LIST_SHORT_START = sp.bytes("0xc0")


def encode_bmc_service(params):
    sp.set_type(params, sp.TRecord(serviceType=sp.TString, payload=sp.TBytes))

    encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
    _rlpBytes = params.payload + encode_string_packed(params.serviceType)
    with_length_prefix = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)
    rlp_bytes_with_prefix = with_length_prefix(_rlpBytes)
    return rlp_bytes_with_prefix

