from pyteal import *

global_bsh_app_address = Bytes("bsh_app_address")
global_relayer_acc_address = Bytes("relayer_acc_address")

is_creator = Txn.sender() == Global.creator_address()
is_relayer = Txn.sender() == App.globalGet(global_relayer_acc_address)

router = Router(
    "bmc-handler",
    BareCallActions(
        no_op=OnCompleteAction.create_only(
            Seq(
                App.globalPut(global_relayer_acc_address, Global.creator_address()),
                Approve()
            )
        ),
        opt_in=OnCompleteAction.always(Approve()),
        update_application=OnCompleteAction.always(Return(is_creator)),
        delete_application=OnCompleteAction.always(Return(is_creator)),
        clear_state=OnCompleteAction.never(),
    ),
)   

@router.method
def registerBSHContract(bsh_address: abi.Address, svc: abi.String): 
    """
    This method store service name into BSH account local storage.
    
    The caller must be creator of BMC contract.
    Args:
        bsh_app_address: Address of BSH smart contract.
        svc: Service name of BSH contract.
    """

    return Seq(
        Assert(is_creator),
        App.localPut(bsh_address.get(), Bytes("svc"), svc.get()),
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
def sendMessage (to: abi.String, sn: abi.Uint64, msg: abi.DynamicBytes) -> Expr:
    """
    This method Log service name from registered BSH's
    
    The caller must be an registered BSH smart contract.
    Args:
        to: BTP Address of destination BMC.
        sn: Serial number of the message, it should be positive.
        msg: BSH Message in bytes to be picked up by relayer.
    """

    return Seq(
        (sender_svc := abi.String()).set(App.localGet(Txn.sender(), Bytes("svc"))),
        Log(sender_svc.get()),
    )

@router.method
def handleRelayMessage (bsh_app: abi.Application, svc: abi.String, msg: abi.DynamicBytes) -> Expr:
    i = ScratchVar(TealType.uint64)

    return Seq(
        Assert(is_relayer),

        bsh_addr := bsh_app.params().address(),
        Assert(bsh_addr.hasValue()),
        #TODO validate svc
        #Assert(App.localGet(bsh_addr.value(), Bytes("svc")) == svc.get()),

        InnerTxnBuilder.Begin(),
        InnerTxnBuilder.MethodCall(
            app_id=bsh_app.application_id(),
            method_signature="handleBTPMessage(byte[])void",
            args=[msg],
            extra_fields={
                TxnField.fee: Int(0),
                TxnField.assets: Txn.assets,
            }
        ),

        For(i.store(Int(1)), i.load() <= Txn.accounts.length(), i.store(i.load() + Int(1))).Do(
            InnerTxnBuilder.SetField(TxnField.accounts, [Txn.accounts[i.load()]]),
        ),

        InnerTxnBuilder.Submit(),
    )

