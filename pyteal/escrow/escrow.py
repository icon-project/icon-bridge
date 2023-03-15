from pyteal import *
from typing import Literal

global_initialized = Bytes("initialized")
global_bmc_id = Bytes("bmc_id")
global_receiver_address = Bytes("receiver_address")
global_asset_id = Bytes("asset_id")

is_creator = Txn.sender() == Global.creator_address()
is_initialized = App.globalGet(global_initialized) == Int(1)

router = Router(
    "bsh-handler",
    BareCallActions(
        no_op=OnCompleteAction.create_only(
            Seq(
                App.globalPut(global_initialized, Int(0)),
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
def init(bmc_app: abi.Application, receiver_address: abi.String, asaId: abi.Uint64) -> Expr:
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

        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.SetFields({
            TxnField.type_enum: TxnType.AssetTransfer,
            TxnField.xfer_asset: asaId.get(),
            TxnField.asset_amount: Int(0),
            TxnField.asset_receiver: Global.current_application_address()
        }),
        InnerTxnBuilder.Submit(),
        App.globalPut(global_asset_id, asaId.get()),

        App.globalPut(global_initialized, Int(1)),
        Approve(),
    )

@router.method
def deposit(amount: abi.Uint64, isContract: abi.Bool, dst: abi.StaticBytes[Literal[20]]) -> Expr:
    """
    This method send BTP message to other chain using BMC smart contract.
    
    Args:
        bmc_app: ID of the BMC application that should process the message.
        to: BTP Address of destination BMC.
    """

    return Seq(
        Assert(is_initialized),
        
        Assert(
          And(
              Txn.group_index() == Int(0),
              Gtxn[0].sender() == Gtxn[1].sender(),
              Gtxn[1].type_enum() == TxnType.AssetTransfer,
              Gtxn[1].xfer_asset() == App.globalGet(global_asset_id),
              Gtxn[1].asset_receiver() == Global.current_application_address(),
              Gtxn[1].asset_amount() == amount.get(),
          )
        ),

        (sn := abi.Uint64()).set(Int(1)),
        (to := abi.String()).set(App.globalGet(global_receiver_address)),

        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.MethodCall(
            app_id=App.globalGet(global_bmc_id),
            method_signature="sendMessage(string,uint64,byte[])void",
            args=[to, sn, Concat(amount.encode(), isContract.encode(), dst.get())],
            extra_fields={
                TxnField.fee: Int(0)
            }
        ),
        InnerTxnBuilder.Submit(),
    )

@router.method
def handleBTPMessage(msg: abi.DynamicBytes) -> Expr:
    return Seq(
        Assert(is_initialized),
        Assert(App.globalGet(global_bmc_id) == Global.caller_app_id()),

        (amount := abi.Uint64()).set(Btoi(Substring(msg.get(), Int(0), Int(8)))),
        (receiver_address := abi.DynamicBytes()).set(Substring(msg.get(), Int(8), Int(40))),

        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.SetFields({
            TxnField.type_enum: TxnType.AssetTransfer,
            TxnField.xfer_asset: App.globalGet(global_asset_id),
            TxnField.asset_receiver: receiver_address.get(),
            TxnField.asset_amount: amount.get(),
            TxnField.fee: Int(0)
        }),
        InnerTxnBuilder.Submit(),

        Approve()
    )