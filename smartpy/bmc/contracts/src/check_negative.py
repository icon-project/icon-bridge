import smartpy as sp
Utils2 = sp.io.import_script_from_url("https://raw.githubusercontent.com/RomarQ/tezos-sc-utils/main/smartpy/utils.py")

#
# @sp.module
# def main():
#     class C(sp.Contract):
#
#         @sp.onchain_view()
#         def check_negative(self, x):
#             sp.cast(x, sp.bytes)
#             return (sp.to_int(x) < 0)
#
#         @sp.onchain_view()
#         def to_int(self, x):
#             sp.cast(x, sp.bytes)
#             return (sp.to_int(x))
#
#
# @sp.add_test(name="test")
# def test():
#     scenario = sp.test_scenario(main)
#     c = main.C()
#     scenario += c

class Sample(sp.Contract):
    def __init__(self):
        self.init(
            tf = False
        )

    @sp.entry_point
    def test(self, addr):
        sp.set_type(addr, sp.TAddress)
        x = sp.view("check_negative", addr, sp.bytes("0xf6"), t=sp.TBool).open_some()
        self.data.tf = x

@sp.add_test("Tests")
def test():
    sc = sp.test_scenario()
    c = Sample()
    sc += c

sp.add_compilation_target("check_negative", Sample())