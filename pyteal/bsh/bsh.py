from pyteal import *

is_creator = Txn.sender() == Global.creator_address()

router = Router(
    "bsh-handler",
    BareCallActions(
        no_op=OnCompleteAction.create_only(Approve()),
        update_application=OnCompleteAction.always(Return(is_creator)),
        delete_application=OnCompleteAction.always(Return(is_creator)),
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
def handleBTPMessage(msg: abi.String) -> Expr:
    return Seq(
        Log(msg.get()),
    )