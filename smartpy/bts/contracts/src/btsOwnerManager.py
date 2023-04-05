import smartpy as sp

class BTSOwnerManager(sp.Contract):
    def __init__(self,owner_address):
        self.update_initial_storage(
            owners=sp.map(tkey=sp.TAddress, tvalue=sp.TBool),
            listOfOwners=sp.list([owner_address]),
         
        )
     

    #Entry point for adding a new owner to the contract
    @sp.entry_point
    def add_owner(self, new_owner):

        sp.verify(self.data.owners[sp.sender], message="Unauthorized")

        # Verifies that the new owner does not already exist in the owners map
        sp.verify(self.data.owners[new_owner] == False,message="ExistedOwner")
        # Adds the new owner to the owners map and to the listOfOwners list
        self.data.owners[new_owner] == True
        # self.data.listOfOwners.push(new_owner)
        # Emits an event with the new owner addresses
        sp.emit(new_owner, tag="NewOwnerAdded")


    # Entry point for removing an owner from the contract
    @sp.entry_point
    def remove_owner(self, removed_owner):
        sp.verify(self.data.owners[sp.sender], message="Unauthorized")
        # Verifies that there are more than one owners in the contract
        sp.verify(
            sp.len(self.data.listOfOwners) > 1,
            message="CannotRemoveMinOwner"
        )
        # Verifies that the removed owner exists in the owners map
        sp.verify(
            self.data.owners[removed_owner] == True,
            message="NotanOwner"
        )
        # Deletes the removed owner from the owners map and the listOfOwners list
        del self.data.owners[removed_owner]
        # self._remove(removed_owner)
        # Emits an event with the removed owner addresses
        sp.emit(removed_owner, tag="OwnerRemoved")

    # Internal function for removing an owner from the listOfOwners list
    # @sp.private_lambda
    def _remove(self, removed_owner):
        # Loops through the listOfOwners list to find the removed owner and removes it
        sp.for i in sp.range(0,sp.len(self.data.listOfOwners)):
            sp.if self.data.listOfOwners[i] == removed_owner:
                self.data.listOfOwners[i] = self.data.listOfOwners[sp.len(self.data.listOfOwners) - 1]
                self.data.listOfOwners.pop()
        

    # External view method to check if an address is an owner of the contract
    @sp.onchain_view()
    def isOwner(self, _owner):
        sp.result( self.data.owners[_owner])

    # External view method to get a list of all the owners in the contract
    @sp.onchain_view()
    def getOwners(self):
        sp.result(self.data.listOfOwners)

@sp.add_test(name = "BTSOwnerManager")
def test():
    alice = sp.test_account("Alice")
    c1 = BTSOwnerManager(alice.address)
    scenario = sp.test_scenario()
    scenario.h1("BTSOwnerManager")
    scenario += c1
   
    # scenario.verify(val == "o")

sp.add_compilation_target("OwnerManager", BTSOwnerManager(sp.address("tz1UVtzTTE1GatMoXhs46hbdp1143a195kXh")))

    
