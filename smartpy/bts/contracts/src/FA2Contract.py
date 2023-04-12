import smartpy as sp

FA2 = sp.io.import_script_from_url("https://smartpy.io/templates/fa2_lib.py")


class SingleAssetToken(FA2.Admin, FA2.Fa2SingleAsset, FA2.BurnSingleAsset, FA2.OnchainviewBalanceOf):
    def __init__(self, admin, metadata, token_metadata):
        FA2.Fa2SingleAsset.__init__(self, metadata=metadata, token_metadata=token_metadata)
        FA2.Admin.__init__(self, admin)

    @sp.entry_point
    def mint(self, batch):
        """Admin can mint tokens."""
        sp.set_type(
            batch,
            sp.TList(
                sp.TRecord(to_=sp.TAddress, amount=sp.TNat).layout(("to_", "amount"))
            ),
        )
        sp.verify(self.is_administrator(sp.sender), "FA2_NOT_ADMIN")
        with sp.for_("action", batch) as action:
            sp.verify(self.is_defined(0), "FA2_TOKEN_UNDEFINED")
            self.data.supply += action.amount
            self.data.ledger[action.to_] = (
                self.data.ledger.get(action.to_, 0) + action.amount
            )
            
    @sp.onchain_view()
    def is_admin(self, address):
        sp.result(address == self.data.administrator)

    @sp.onchain_view()
    def balance_of(self, param):
        sp.set_type(param, sp.TRecord(owner=sp.TAddress, token_id=sp.TNat))
        sp.result(self.balance_(param.owner, param.token_id))


sp.add_compilation_target("fa2_single_asset",
                          SingleAssetToken(
                              admin=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
                              metadata=sp.utils.metadata_of_url(
                                  "ipfs://example"),
                              token_metadata=FA2.make_metadata(name="Token Zero", decimals=1, symbol="Tok0")))