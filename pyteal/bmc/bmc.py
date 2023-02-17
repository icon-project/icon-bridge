from pyteal import *

global_bsh_app_address = Bytes("bsh_app_address")
global_relayer_acc_address = Bytes("relayer_acc_address")

is_creator = Txn.sender() == Global.creator_address()
is_relayer = Txn.sender() == App.globalGet(global_relayer_acc_address)
is_bsh = Txn.sender() == App.globalGet(global_bsh_app_address)

router = Router(
    "bmc-handler",
    BareCallActions(
        no_op=OnCompleteAction.create_only(
            Seq(
                App.globalPut(global_relayer_acc_address, Global.creator_address()),
                Approve()
            )
        ),
        update_application=OnCompleteAction.always(Return(is_creator)),
        delete_application=OnCompleteAction.always(Return(is_creator)),
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
def setRelayer(relayer_account: abi.Address): 
    return Seq(
        Assert(is_relayer),
        App.globalPut(global_relayer_acc_address, relayer_account.get()),
        Approve()
    )
    
@router.method
def sendMessage (to: abi.String, svc: abi.String, sn: abi.Uint64, ) -> Expr:
    return Seq(
        Log(Bytes("hello world"))
    )

@router.method
def handleRelayMessage (bsh_app: abi.Application, msg: abi.String,  *, output: abi.String) -> Expr:
    return Seq(
        Assert(is_relayer),
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