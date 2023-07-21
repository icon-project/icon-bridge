import smartpy as sp


class BTSOwnerManager(sp.Contract):
    def __init__(self, owner):
        self.init(
            owners=sp.map({owner: True}),
        )
        self.init_type(sp.TRecord(
            owners=sp.TMap(sp.TAddress, sp.TBool),
        ))

    def only_owner(self):
        sp.verify(self.data.owners[sp.sender] == True, message="Unauthorized")

    @sp.entry_point
    def add_owner(self, owner):
        """
        :param owner: address to set as owner
        :return:
        """
        sp.set_type(owner, sp.TAddress)

        self.only_owner()
        sp.verify(self.data.owners[owner] == False, message="ExistedOwner")

        self.data.owners[owner] = True
        sp.emit(sp.record(sender=sp.sender, owner=owner), tag="NewOwnerAdded")

    @sp.entry_point
    def remove_owner(self, owner):
        """
        :param owner: address to remove as owner
        :return:
        """
        sp.set_type(owner, sp.TAddress)

        self.only_owner()
        sp.verify(sp.len(self.data.owners) > 1, message="CannotRemoveMinOwner")
        sp.verify(self.data.owners[owner] == True, message="NotOwner")

        del self.data.owners[owner]
        sp.emit(sp.record(sender=sp.sender, owner=owner), tag="OwnerRemoved")

    @sp.onchain_view()
    def is_owner(self, owner):
        sp.set_type(owner, sp.TAddress)
        sp.result(self.data.owners.get(owner, default_value=False))

    @sp.onchain_view()
    def get_owners(self):
        sp.result(self.data.owners.keys())


@sp.add_test(name="BTSOwnerManager")
def test():
    alice = sp.test_account("Alice")
    c1 = BTSOwnerManager(alice.address)
    scenario = sp.test_scenario()
    scenario.h1("BTSOwnerManager")
    scenario += c1

    scenario.verify(c1.is_owner(alice.address) == True)


sp.add_compilation_target("bts_owner_manager", BTSOwnerManager(owner=sp.address("tz1g3pJZPifxhN49ukCZjdEQtyWgX2ERdfqP")))

    
