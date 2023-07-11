import smartpy as sp

FA2 = sp.io.import_script_from_url("https://legacy.smartpy.io/templates/fa2_lib.py")


class SingleAssetToken(FA2.Admin, FA2.Fa2SingleAsset, FA2.MintSingleAsset, FA2.BurnSingleAsset,
                       FA2.OnchainviewBalanceOf):
    def __init__(self, admin, metadata, token_metadata):
        FA2.Fa2SingleAsset.__init__(self, metadata=metadata, token_metadata=token_metadata)
        FA2.Admin.__init__(self, admin)

    @sp.onchain_view()
    def is_admin(self, address):
        sp.result(address == self.data.administrator)


# @sp.add_test(name="FA2DummyContract")
# def test():
#     alice = sp.test_account("Alice")
#
#     c1 = SingleAssetToken(admin=alice.address, metadata=sp.big_map({"ss": sp.bytes("0x0dae11")}),
#                           token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}))
#
#     scenario = sp.test_scenario()
#     scenario += c1


sp.add_compilation_target("fa2_dummy",
                          SingleAssetToken(
                              admin=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"),
                              metadata=sp.utils.metadata_of_url(
                                  "ipfs://example"),
                              token_metadata=FA2.make_metadata(name="NativeWrappedCoin", decimals=6, symbol="WTEZ")))
