import smartpy as sp

FA2 = sp.io.import_script_from_url("https://legacy.smartpy.io/templates/fa2_lib.py")

t_transfer_batch = sp.TRecord(
    callback=sp.TContract(
        sp.TRecord(string=sp.TOption(sp.TString), requester=sp.TAddress, coin_name=sp.TString, value=sp.TNat)),
    from_=sp.TAddress,
    coin_name=sp.TString,
    txs=sp.TList(
        sp.TRecord(
            to_=sp.TAddress,
            token_id=sp.TNat,
            amount=sp.TNat,
        ).layout(("to_", ("token_id", "amount")))
    ),
).layout((("from_", "coin_name"), ("callback", "txs")))

t_transfer_params = sp.TList(t_transfer_batch)


class SingleAssetToken(FA2.Admin, FA2.Fa2SingleAsset, FA2.MintSingleAsset, FA2.BurnSingleAsset,
                       FA2.OnchainviewBalanceOf):
    def __init__(self, admin, metadata, token_metadata):
        FA2.Fa2SingleAsset.__init__(self, metadata=metadata, token_metadata=token_metadata)
        FA2.Admin.__init__(self, admin)
        self.update_initial_storage(
            allowances=sp.big_map(tkey=sp.TRecord(spender=sp.TAddress, owner=sp.TAddress), tvalue=sp.TNat)
        )

    @sp.onchain_view()
    def is_admin(self, address):
        sp.result(address == self.data.administrator)

    @sp.entry_point
    def set_allowance(self, batch):
        sp.set_type(batch, sp.TList(sp.TRecord(spender=sp.TAddress, amount=sp.TNat)))

        with sp.for_("params", batch) as params:
            allowance = sp.compute(sp.record(spender=params.spender, owner=sp.sender))
            current_allowance = self.data.allowances.get(allowance, default_value=0)
            sp.verify((params.amount == 0) | (current_allowance == 0), "FA2_UnsafeAllowanceChange")
            self.data.allowances[allowance] = params.amount

    @sp.onchain_view()
    def get_allowance(self, allowance):
        sp.set_type(allowance, sp.TRecord(spender=sp.TAddress, owner=sp.TAddress))
        sp.result(self.data.allowances.get(allowance, default_value=0))

    @sp.onchain_view()
    def transfer_permissions(self, params):
        sp.set_type(params, sp.TRecord(from_=sp.TAddress, token_id=sp.TNat))

        with sp.if_((self.policy.supports_transfer) & (self.is_defined(params.token_id))):
            with sp.if_((sp.sender == params.from_) | (self.data.operators.contains(
                    sp.record(owner=params.from_, operator=sp.sender, token_id=params.token_id)))):
                sp.result(True)
            with sp.else_():
                sp.result(False)
        with sp.else_():
            sp.result(False)

    @sp.entry_point
    def transfer(self, batch):
        """Accept a list of transfer operations between a source and multiple
        destinations.
        Custom version with allowance system.
        `transfer_tx_` must be defined in the child class.
        """
        sp.set_type(batch, FA2.t_transfer_params)
        if self.policy.supports_transfer:
            with sp.for_("transfer", batch) as transfer:
                with sp.for_("tx", transfer.txs) as tx:
                    # The ordering of sp.verify is important: 1) token_undefined, 2) transfer permission 3) balance
                    sp.verify(self.is_defined(tx.token_id), "FA2_TOKEN_UNDEFINED")
                    self.policy.check_tx_transfer_permissions(
                        self, transfer.from_, tx.to_, tx.token_id
                    )
                    with sp.if_(sp.sender != transfer.from_):
                        self.update_allowance_(sp.sender, transfer.from_, tx.token_id, tx.amount)
                    with sp.if_(tx.amount > 0):
                        self.transfer_tx_(transfer.from_, tx)
        else:
            sp.failwith("FA2_TX_DENIED")

    def update_allowance_(self, spender, owner, token_id, amount):
        allowance = sp.record(spender=spender, owner=owner)
        self.data.allowances[allowance] = sp.as_nat(self.data.allowances.get(allowance, default_value=0) - amount,
                                                    message=sp.pair("FA2_NOT_OPERATOR", "NoAllowance"))


@sp.add_test(name="FA2Contract")
def test():
    alice = sp.test_account("Alice")
    bob = sp.test_account("Bob")
    spender = sp.test_account("spender")
    receiver = sp.test_account("receiver")

    c1 = SingleAssetToken(admin=alice.address, metadata=sp.big_map({"ss": sp.bytes("0x0dae11")}),
                          token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}))

    scenario = sp.test_scenario()
    scenario.h1("FA2Contract")
    scenario += c1

    c1.mint([sp.record(to_=bob.address, amount=sp.nat(200))]).run(sender=alice)
    scenario.verify(c1.data.ledger.get(bob.address) == sp.nat(200))

    scenario.verify(c1.get_allowance(sp.record(spender=spender.address, owner=bob.address)) == 0)
    # set allowance
    c1.set_allowance([sp.record(spender=spender.address, amount=sp.nat(100))]).run(sender=bob)
    c1.set_allowance([sp.record(spender=spender.address, amount=sp.nat(100))]).run(sender=bob, valid=False,
                                                                                   exception="FA2_UnsafeAllowanceChange")

    scenario.verify(c1.get_allowance(sp.record(spender=spender.address, owner=bob.address)) == 100)
    c1.update_operators(
        [sp.variant("add_operator", sp.record(owner=bob.address, operator=spender.address, token_id=0))]).run(
        sender=bob)
    # transfer more than allowance
    # c1.transfer([sp.record(callback=sp.self_entry_point("callback"), from_=bob.address, txs=[sp.record(to_=receiver.address, token_id=0, amount=101)])]).run(
    #     sender=spender, valid=False, exception=('FA2_NOT_OPERATOR', 'NoAllowance'))
    # # transfer all allowance
    # c1.transfer([sp.record(from_=bob.address, txs=[sp.record(to_=receiver.address, token_id=0, amount=100)])]).run(
    #     sender=spender)

    # # verify remaining allowance
    # scenario.verify(c1.get_allowance(sp.record(spender=spender.address, owner=bob.address)) == 0)
    # # again set allowance after prev-allowance is 0
    # c1.set_allowance([sp.record(spender=spender.address, amount=sp.nat(100))]).run(sender=bob)


sp.add_compilation_target("FA2_contract",
                          SingleAssetToken(
                              admin=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP"),
                              metadata=sp.utils.metadata_of_url(
                                  "ipfs://example"),
                              token_metadata=FA2.make_metadata(name="NativeWrappedCoin", decimals=6, symbol="WTEZ")))
