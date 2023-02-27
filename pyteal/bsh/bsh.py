from pyteal import *

global_initialized = Bytes("initialized")
global_bmc_id = Bytes("bmc_id")
global_receiver_address = Bytes("receiver_address")
global_last_received_message = Bytes("last_received_message")

is_creator = Txn.sender() == Global.creator_address()
is_init = App.globalGet(global_initialized) == Int(1)

router = Router(
    "bsh-handler",
    BareCallActions(
        no_op=OnCompleteAction.create_only(
            Seq(
                App.globalPut(global_initialized, Int(0)),
                App.globalPut(global_last_received_message, Bytes("BTP message")),
                Approve()
            )
        ),
        opt_in=OnCompleteAction.never(),
        update_application=OnCompleteAction.always(Return(is_creator)),
        delete_application=OnCompleteAction.always(Return(is_creator)),
        clear_state=OnCompleteAction.never(),
    ),
)

@router.method
def init(bmc_app: abi.Application, receiver_address: abi.String) -> Expr:
    """ Initialize Smart Contract """

    return Seq(
        Assert(App.globalGet(global_initialized) == Int(0)),
        Assert(is_creator),
        App.globalPut(global_bmc_id, bmc_app.application_id()),
        App.globalPut(global_receiver_address, receiver_address.get()),

        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.SetFields({
            TxnField.type_enum: TxnType.ApplicationCall,
            TxnField.on_completion: OnComplete.OptIn,
            TxnField.application_id: bmc_app.application_id(),
            TxnField.fee: Int(0)
        }),
        InnerTxnBuilder.Submit(),

        App.globalPut(global_initialized, Int(1)),
        Approve(),
    )

@router.method
def sendServiceMessage() -> Expr:
    """
    This method send BTP message to other chain using BMC smart contract.
    
    Args:
        bmc_app: ID of the BMC application that should process the message.
        to: BTP Address of destination BMC.
    """

    return Seq(
        Assert(is_init),
        (sn := abi.Uint64()).set(Int(1)),
        (msg := abi.String()).set("hello world"),
        (to := abi.String()).set(App.globalGet(global_receiver_address)),

        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.MethodCall(
            app_id=App.globalGet(global_bmc_id),
            method_signature="sendMessage(string,uint64,byte[])void",
            args=[to, sn, msg.encode()],
            extra_fields={
                TxnField.fee: Int(0)
            }
        ),
        InnerTxnBuilder.Submit(),
    )

@router.method
def handleBTPMessage(msg: abi.String) -> Expr:
    return Seq(
        Assert(is_init),
        Assert(App.globalGet(global_bmc_id) == Global.caller_app_id()),
        
        App.globalPut(global_last_received_message, msg.get()),
        Approve()
    )