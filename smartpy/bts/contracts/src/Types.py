import smartpy as sp


class Types:

    Asset = sp.TRecord(
        coin_name=sp.TString,
        value=sp.TNat
    )

    AssetTransferDetail = sp.TRecord(
        coin_name=sp.TString,
        value=sp.TNat,
        fee=sp.TNat
    )

    Response = sp.TRecord(
        code=sp.TNat,
        message=sp.TString
    )

    ServiceType = sp.TVariant(
        REQUEST_COIN_TRANSFER=sp.TNat,
        REQUEST_COIN_REGISTER=sp.TNat,
        RESPONSE_HANDLE_SERVICE=sp.TNat,
        BLACKLIST_MESSAGE=sp.TNat,
        CHANGE_TOKEN_LIMIT=sp.TNat,
        UNKNOWN_TYPE=sp.TNat
    )

    BlacklistService = sp.TVariant(
        ADD_TO_BLACKLIST=sp.TNat,
        REMOVE_FROM_BLACKLIST=sp.TNat
    )

    ServiceMessage = sp.TRecord(
        serviceType=ServiceType,
        data=sp.TBytes
    )

    TransferCoin = sp.TRecord(
        from_addr=sp.TString,
        to=sp.TString,
        assets=sp.TList(Asset)
    )

    PendingTransferCoin = sp.TRecord(
        # from_addr=sp.TString, (change to original later)
        from_addr=sp.TAddress,
        to=sp.TString,
        coin_names=sp.TMap(sp.TNat, sp.TString),
        amounts=sp.TMap(sp.TNat, sp.TNat),
        fees=sp.TMap(sp.TNat, sp.TNat)
    )

    BlacklistMessage = sp.TRecord(
        serviceType=BlacklistService,
        addrs=sp.TMap(sp.TNat, sp.TString),
        net=sp.TString
    )

    TokenLimitMessage = sp.TRecord(
        coin_name=sp.TMap(sp.TNat, sp.TString),
        token_limit=sp.TMap(sp.TNat, sp.TNat),
        net=sp.TString
    )

