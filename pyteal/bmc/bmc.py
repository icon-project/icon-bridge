from pyteal import *

# Create a simple Expression to use later
# Creator assume to be the relayer
is_creator = Txn.sender() == Global.creator_address()

global_bsh_app_address = Bytes("bsh_app_address")

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
        Assert(is_creator),
        App.globalPut(global_bsh_app_address, bsh_app_address.get()),
        Approve()
    )
    
@router.method
def sendMessage (to: abi.String, svc: abi.String, sn: abi.Uint64,  *, output: abi.String) -> Expr:
    return Seq(
        Assert(Txn.sender() == App.globalGet(global_bsh_app_address)),
        output.set("event:btp message")
    )

@router.method
def handleRelayMessage (bsh_app: abi.Application, msg: abi.String,  *, output: abi.String) -> Expr:
    return Seq(
        Assert(is_creator),
        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.MethodCall(
            app_id=bsh_app.application_id(),
            method_signature="handleBTPMessage(string)void",
            args=[msg],
            extra_fields={
                TxnField.fee: Int(0)
            }
        ),
        InnerTxnBuilder.Submit(),
        output.set("event:start handleBTPMessage")
    )