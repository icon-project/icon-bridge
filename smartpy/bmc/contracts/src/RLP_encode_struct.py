import smartpy as sp


class EncodeLibrary:
    LIST_SHORT_START = sp.bytes("0xc0")

    def encode_bmc_service(self, params):
        sp.set_type(params, sp.TRecord(serviceType=sp.TString, payload=sp.TBytes))

        encode_service_type = sp.view("encode_string", self.data.helper, params.serviceType, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_service_type, params.payload], t=sp.TBytes).open_some()
        return rlp_bytes_with_prefix

    def encode_bmc_message(self, params):
        sp.set_type(params, sp.TRecord(src=sp.TString, dst=sp.TString, svc=sp.TString, sn=sp.TInt, message=sp.TBytes))

        encode_src = sp.view("encode_string", self.data.helper, params.src, t=sp.TBytes).open_some()
        encode_dst = sp.view("encode_string", self.data.helper, params.dst, t=sp.TBytes).open_some()
        encode_svc = sp.view("encode_string", self.data.helper, params.svc, t=sp.TBytes).open_some()
        nat_sn = sp.as_nat(params.sn)
        encode_sn = sp.view("encode_nat", self.data.helper, nat_sn, t=sp.TBytes).open_some()
        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_src, encode_dst, encode_svc, encode_sn, params.message], t=sp.TBytes).open_some()
        return rlp_bytes_with_prefix

    def encode_response(self, params):
        sp.set_type(params, sp.TRecord(code=sp.TNat, message=sp.TString))

        encode_code = sp.view("encode_nat", self.data.helper, params.code, t=sp.TBytes).open_some()
        encode_message = sp.view("encode_string", self.data.helper, params.message, t=sp.TBytes).open_some()

        rlp_bytes_with_prefix = sp.view("encode_list", self.data.helper, [encode_code, encode_message], t=sp.TBytes).open_some()
        return rlp_bytes_with_prefix
