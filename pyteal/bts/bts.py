from pyteal import *

# Create a simple Expression to use later
is_creator = Txn.sender() == Global.creator_address()

# Main router class
router = Router(
    # Name of the contract
    "bts-handler",
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
def sendServiceMessage(bmc_app: abi.Application, to: abi.String, svc: abi.String, sn: abi.Uint64) -> Expr:
    return Seq(
        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.MethodCall(
            app_id=bmc_app.application_id(),
            method_signature="sendMessage(string,string,uint64)string",
            args=[to, svc, sn],
            extra_fields={
                TxnField.fee: Int(0)
            }
        ),
        InnerTxnBuilder.Submit(),
    )

@router.method
def handleBTPMessage(msg: abi.String, *, output: abi.String) -> Expr:
    return Seq(
        output.set("start:handle BTP Message")
    )

