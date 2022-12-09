from pyteal import *

# Create a simple Expression to use later
is_creator = Txn.sender() == Global.creator_address()

global_bsh_app_address = Bytes("bsh_app_address")
global_relayer_acc_address = Bytes("relayer_acc_address")

# Main router class
router = Router(
    # Name of the contract
    "bmc-handler",
    # What to do for each on-complete type when no arguments are passed (bare call)
    BareCallActions(
        # On create only, just approve
        no_op=OnCompleteAction.create_only(Approve()),
        # Always let creator update/delete but only by the creator of this contract
        update_application=OnCompleteAction.always(Return(is_creator)),
        delete_application=OnCompleteAction.always(Return(is_creator)),
        # No local state, dont bother handling it
        # close_out=OnCompleteAction.never(),
        # opt_in=OnCompleteAction.never(),
        # Just be nice, we _must_ provide _something_ for clear state becuase it is its own
        # program and the router needs _something_ to build
        # clear_state=OnCompleteAction.call_only(Approve()),
        clear_state=OnCompleteAction.never(),
    ),
)   

@router.method
def registerBSHContract(bsh_app_address: abi.Address): 
    return Seq(
        App.globalPut(global_bsh_app_address, bsh_app_address.get()),
        Approve()
    )
    
@router.method
def registerRelayer(relayer_account: abi.Address): 
    return Seq(
        App.globalPut(global_relayer_acc_address, relayer_account.get()),
        Approve()
    )

@router.method
def sendMessage (to: abi.String, svc: abi.String, sn: abi.Uint64,  *, output: abi.String) -> Expr:
    return Seq(
        Assert(Txn.sender() == App.globalGet(global_bsh_app_address)),
        output.set("btp:message")
    )