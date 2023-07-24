import smartpy as sp


@sp.module
def main():
    class Convert(sp.Contract):

        @sp.onchain_view()
        def check_negative(self, x):
            sp.cast(x, sp.bytes)
            return sp.to_int(x) < 0

        @sp.onchain_view()
        def to_int(self, x):
            sp.cast(x, sp.bytes)
            return sp.to_int(x)

        @sp.onchain_view()
        def to_byte(self, x):
            return sp.to_bytes(x)


@sp.add_test(name="test")
def test():
    scenario = sp.test_scenario(main)
    c = main.Convert()
    scenario += c
