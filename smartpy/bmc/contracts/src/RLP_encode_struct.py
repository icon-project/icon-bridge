import smartpy as sp

Utils = sp.io.import_script_from_url(
    "https://raw.githubusercontent.com/Acurast/acurast-hyperdrive/main/contracts/tezos/libs/utils.py")

LIST_SHORT_START = sp.bytes("0xc0")


def encode_bmc_service(params):
    sp.set_type(params, sp.TRecord(serviceType=sp.TString, payload=sp.TBytes))

    encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
    _rlpBytes = params.payload + encode_string_packed(params.serviceType)
    with_length_prefix = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)
    rlp_bytes_with_prefix = with_length_prefix(_rlpBytes)
    return rlp_bytes_with_prefix


def encode_bmc_message(params):
    sp.set_type(params, sp.TRecord(src=sp.TString, dst=sp.TString, svc=sp.TString, sn=sp.TNat, message=sp.TBytes))

    encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
    encode_nat_packed = sp.build_lambda(Utils.Bytes.of_nat)
    rlp = encode_string_packed(params.src) + encode_string_packed(params.dst)\
          + encode_string_packed(params.svc) + encode_nat_packed(params.sn) + params.message
    with_length_prefix = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)
    rlp_bytes_with_prefix = with_length_prefix(rlp)
    return rlp_bytes_with_prefix


def encode_response(params):
    sp.set_type(params, sp.TRecord(code=sp.TNat, message=sp.TString))

    encode_string_packed = sp.build_lambda(Utils.Bytes.of_string)
    encode_nat_packed = sp.build_lambda(Utils.Bytes.of_nat)
    rlp = encode_nat_packed(params.code) + encode_string_packed(params.message)
    with_length_prefix = sp.build_lambda(Utils.RLP.Encoder.with_length_prefix)
    rlp_bytes_with_prefix = with_length_prefix(rlp)
    return rlp_bytes_with_prefix
