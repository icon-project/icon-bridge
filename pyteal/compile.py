import importlib
import sys
import os
import json
from pyteal import *

def application(pyteal: Expr) -> str:
    return compileTeal(pyteal, mode=Mode.Application, version=MAX_TEAL_VERSION)

if __name__ == "__main__":
    mod = sys.argv[1]

    try:
        approval_out = sys.argv[2]
    except IndexError:
        approval_out = None

    try:
        clear_out = sys.argv[3]
    except IndexError:
        clear_out = None
    try:
        abi_out = sys.argv[4]
    except IndexError:
        abi_out = None

    contract = importlib.import_module(mod)

    approval, clear, contract = contract.router.compile_program(
        version=6, optimize=OptimizeOptions(scratch_slots=True)
    )

    # Dump out the contract as json that can be read in by any of the SDKs
    if abi_out is None:
        print('Please provide a out dir for ABI contract.json file as last argument')
    else:
        with open(os.path.join(abi_out, "contract.json"), "w") as f:
            f.write(json.dumps(contract.dictify(), indent=2))

    if approval_out is None:
        print(approval)
    else:
        with open(approval_out, "w") as h:
            h.write(approval)

    if clear_out is not None:
        with open(clear_out, "w") as h:
            h.write(clear)