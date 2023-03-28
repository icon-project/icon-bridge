import smartpy as sp


class Types:

    MessageEvent = sp.TRecord(
        next_bmc=sp.TString,
        seq=sp.TNat,
        message=sp.TBytes
    )

    ReceiptProof = sp.TRecord(
        index=sp.TNat,
        events=sp.TList(MessageEvent),
        height=sp.TNat
    )

    BMCMessage = sp.TRecord(
        src=sp.TString,
        dst=sp.TString,
        svc=sp.TString,
        sn=sp.TNat,
        message=sp.TBytes
    )

    MessageEvent = sp.TRecord(
        next_bmc=sp.TString,
        seq=sp.TNat,
        message=sp.TBytes
    )

    Response = sp.TRecord(
        code=sp.TNat,
        message=sp.TString
    )

    Route = sp.TRecord(
        dst=sp.TString,
        next=sp.TString
    )

    Link = sp.TRecord(
        relays=sp.TMap(sp.TNat, sp.TAddress),
        reachable=sp.TMap(sp.TNat, sp.TString),
        rx_seq=sp.TNat,
        tx_seq=sp.TNat,
        block_interval_src=sp.TNat,
        block_interval_dst=sp.TNat,
        max_aggregation=sp.TNat,
        delay_limit=sp.TNat,
        relay_idx=sp.TNat,
        rotate_height=sp.TNat,
        rx_height=sp.TNat,
        rx_height_src=sp.TNat,
        is_connected=sp.TBool
    )

    LinkStats = sp.TRecord(
        rx_seq=sp.TNat,
        tx_seq=sp.TNat,
        rx_height=sp.TNat,
        current_height=sp.TNat
    )

    BMCService = sp.TRecord(
        serviceType=sp.TString,
        payload=sp.TBytes
    )

    GatherFeeMessage = sp.TRecord(
        fa=sp.TString,
        svcs=sp.TMap(sp.TNat, sp.TString)
    )

    RelayStats = sp.TRecord(
        addr=sp.TAddress,
        block_count=sp.TNat,
        msg_count=sp.TNat
    )

    Tuple = sp.TRecord(
        prev=sp.TString,
        to=sp.TString
    )

    Service = sp.TRecord(
        svc=sp.TString,
        addr=sp.TAddress
    )

