import smartpy as sp

BMCManagement = sp.io.import_script_from_url("file:./bmc/contracts/src/bmc_management.py")
BMCPeriphery = sp.io.import_script_from_url("file:./bmc/contracts/src/bmc_periphery.py")
BMCHelper = sp.io.import_script_from_url("file:./bmc//contracts/src/helper.py")
ParseAddress = sp.io.import_script_from_url("file:./contracts/src/parse_address.py")

BTSCore = sp.io.import_script_from_url("file:./contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:./contracts/src/bts_owner_manager.py")
BTSPeriphery = sp.io.import_script_from_url("file:./contracts/src/bts_periphery.py")





@sp.add_test("BMCManagementTest")
def test():
    sc = sp.test_scenario()

    # test account
    alice = sp.test_account("Alice")
    creator = sp.test_account("Creator")
    jack = sp.test_account("Jack")
    bob = sp.test_account("Bob")
    creator2 = sp.test_account("creator2")
    service1_address = sp.test_account("service1_address")
    service2_address = sp.test_account("service2_address")


    # deploy BMCManagement contract

    helper_contract = deploy_helper_contract()
    sc += helper_contract

    bmcManagement_contract = deploy_bmcManagement_contract(creator.address, helper_contract.address)
    sc += bmcManagement_contract

    parse_address = deploy_parse_address()
    sc += parse_address

    bmcPeriphery_contract = deploy_bmcPeriphery_contract(bmcManagement_contract.address, helper_contract.address, parse_address.address,creator.address)
    sc += bmcPeriphery_contract

    bts_owner_manager = deploy_btsOwnerManager_Contract(creator.address)
    sc += bts_owner_manager

    btsCore_contract = deploy_btsCore_contract(bts_owner_manager.address)
    sc += btsCore_contract


    bts_periphery = deploy_btsPeriphery_Contract(btsCore_contract.address, helper_contract.address, parse_address.address, bmcPeriphery_contract.address,creator.address)
    sc += bts_periphery

    fa2 = deploy_fa2_Contract(bts_periphery.address)
    sc += fa2

    #set bmc periphery
    bmcManagement_contract.set_bmc_periphery(bmcPeriphery_contract.address).run(sender=creator.address)
    #set bmc_btp_address(netwrk address)
    bmcManagement_contract.set_bmc_btp_address("NetXnHfVqm9iesp.tezos").run(sender=creator.address)
    #update_bts_periphery
    btsCore_contract.update_bts_periphery(bts_periphery.address).run(sender=creator.address)

    #add_service
    svc1 = sp.string("bts")
    bmcManagement_contract.add_service(sp.record(addr=bts_periphery.address, svc=svc1)).run(sender=creator.address)

    #add_route
    dst = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest"
    next_link = "btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b"
    bmcManagement_contract.add_route(sp.record(dst=dst, link=next_link)).run(sender=creator.address)

    #add_link
    bmcManagement_contract.add_link("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b").run(sender=creator.address)


    #set_link_rx_height
    # link = sp.string('btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b')
    # height = sp.nat(2)
    # bmcManagement_contract.set_link_rx_height(link=link, height=height).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # #add_relay
    # bmcManagement_contract.add_relay(sp.record(link=link, addr=sp.set([sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9")]))).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    #test 1: Test of add to blacklist function
    # bts_periphery.add_to_blacklist({0:"notaaddress"}).run(sender=bts_periphery.address,valid=False, exception="InvalidAddress") # invalid address 
    bts_periphery.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=bts_periphery.address,valid=True) #add a address to blacklist
    # bts_periphery.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=bts_periphery.address,valid=True)  # can be called twice 
    # bts_periphery.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}) .run(sender=alice.address,valid=False,exception ="Unauthorized")# only btsperiphery contract call this function
    # bts_periphery.add_to_blacklist({0:'tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg'}).run(sender=bts_periphery.address,valid=False,exception='InvalidAddress')#invalid address
    # sc.verify(bts_periphery.data.blacklist[sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW")] == True) # checking the blacklist[] map

    # transfer_native_coin
    bts_periphery.set_token_limit(
        sp.record(
            coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
            token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
        )
    ).run(sender = btsCore_contract.address)
    btsCore_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), amount=sp.tez(30), valid=False, exception="Blacklisted")


   

    #test 2 : Test of remove from blacklist function
    # bts_periphery.remove_from_blacklist({0:'notaaddress'}).run(sender=bts_periphery.address,valid=False, exception="InvalidAddress") # invalid address 
    # bts_periphery.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address,valid=False, exception="UserNotFound") # address not black-listed
    # bts_periphery.add_to_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address) # adding to blacklist
    bts_periphery.remove_from_blacklist({0:'tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address) # valid process
    btsCore_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= bts_periphery.address, amount=sp.tez(30))

    # bts_periphery.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address,valid=False ,exception ='UserNotFound') # cannot remove from blacklist twice
    # bts_periphery.add_to_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address) # adding to blacklist
    # bts_periphery.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=bts_periphery.address) # can only be called from btseperiphery contract
    # bts_periphery.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=alice.address, valid=False, exception ="Unauthorized") # can only be called from btseperiphery contract

    # transfer_native_coin
    bts_periphery.set_token_limit(
        sp.record(
            coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
            token_limit=sp.map({0: 115792089237316195423570985008687907853269984665640564039457584007913129639935})
        )
    ).run(sender = btsCore_contract.address)
    btsCore_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= bmcManagement_contract.address, amount=sp.tez(30))

    sc.verify_equal(
        btsCore_contract.balance_of(
            sp.record(owner=bmcManagement_contract.address, coin_name='btp-NetXnHfVqm9iesp.tezos-XTZ')
        ),sp.record(usable_balance=0, locked_balance=30000000, refundable_balance=0, user_balance=0))

   

    # transfer_native_coin


    #bts core function test

    #test of transfer native coin
    sc.verify(bts_periphery.data.number_of_pending_requests == 2)#2 pending request
    sc.verify(bts_periphery.data.serial_no == 2)
    btsCore_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender= bmcManagement_contract.address, amount=sp.tez(30))
    #this calls send service message of bts core
    # in send service message function locked balance is called
    sc.show(btsCore_contract.balance_of(
            sp.record(owner=bmcManagement_contract.address, coin_name='btp-NetXnHfVqm9iesp.tezos-XTZ')
        ))
    sc.verify_equal(
        btsCore_contract.balance_of(
            sp.record(owner=bmcManagement_contract.address, coin_name='btp-NetXnHfVqm9iesp.tezos-XTZ')
        ),sp.record(usable_balance=0, locked_balance=60000000, refundable_balance=0, user_balance=0))

    #this calls btsperiphery sendservice message
    #inside bts periphery send service message serial_no is incremented by 1 ,no of pending request is increased by 1, bmc send message is called
    sc.verify(bts_periphery.data.number_of_pending_requests == 3)#2 pending request
    sc.verify(bts_periphery.data.serial_no == 3)
    #bmc management update link tx sq is called links[prev].tx_seq += sp.nat(1)


    #test of non native coin transfer
    btsCore_contract.register(
        name=sp.string("new_coin"),
        fee_numerator=sp.nat(10),
        fixed_fee=sp.nat(2),
        addr=fa2.address,
        token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}),
        metadata=sp.big_map({"ff": sp.bytes("0x0dae11")})
    ).run(sender=creator.address)

    fa2.mint([sp.record(to_=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), amount=sp.nat(100))]).run(sender=bts_periphery.address)
    fa2.set_allowance([sp.record(spender=btsCore_contract.address, amount=sp.nat(100))]).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    fa2.update_operators(
         [sp.variant("add_operator", sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), operator=btsCore_contract.address, token_id=0))]).run(
         sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
    sc.verify_equal(btsCore_contract.is_valid_coin('new_coin'), True)
    btsCore_contract.transfer(coin_name='new_coin', value=10,  to="btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender = sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))    
    sc.verify_equal(
        btsCore_contract.balance_of(
            sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_name='new_coin')
        ),
        sp.record(usable_balance=90, locked_balance=10, refundable_balance=0, user_balance=90)
    )

    #transfer batch
    # test case 16: transfer_batch function
    btsCore_contract.transfer_batch(
        coin_names={0: 'new_coin', 1: 'new_coin'},
        values={0: 10, 1: 10},
        to="btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d"
    ).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    sc.verify_equal(
        btsCore_contract.balance_of(
            sp.record(owner=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"), coin_name='new_coin')
        ),
        sp.record(usable_balance=70, locked_balance=30, refundable_balance=0, user_balance=70)
    )


    
    bts_periphery.handle_request_service(sp.record(to= "tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc", assets={0: sp.record(coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ", value=sp.nat(4))})).run(sender=bts_periphery.address)
    #core's mint is called
    # sc.verify_equal(
    #     btsCore_contract.balance_of(
    #         sp.record(owner=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"), coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ")
    #     ),
    #     sp.record(usable_balance=0, locked_balance=0, refundable_balance=0, user_balance=0)
    # )





    bts_periphery.set_token_limit(
        sp.record(
            coin_names=sp.map({0: "btp-NetXnHfVqm9iesp.tezos-XTZ"}),
            token_limit=sp.map({0: 5})
        )
    ).run(sender = btsCore_contract.address)
    btsCore_contract.transfer_native_coin("btp://0x7.icon/cx4419cb43f1c53db85c4647e4ef0707880309726d").run(sender=bmcManagement_contract.address, amount=sp.tez(30),  valid=False, exception="LimitExceed")


     #Test : handle fee gathering 
     #this calls transfer fees of btscore
    # bts_periphery.handle_fee_gathering(sp.record(fa="btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest", svc="bts")).run(sender=bmcPeriphery_contract.address) # handle_fee_gathering function call
    # bts_periphery.handle_fee_gathering(sp.record(fa="btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest", svc="btc")).run(sender=bmcPeriphery_contract.address, valid=False, exception='InvalidSvc') # svc must match hardcoded service name 'bts'
    # bts_periphery.handle_fee_gathering(sp.record(fa="btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258dest", svc="bts")).run(sender=bts_periphery.address, valid=False, exception='Unauthorized') # can only be called from bmc contract

    

    # #test 3 : set token  limit
    # bts_periphery.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(5),1:sp.nat(2)})).run(sender=alice.address,valid=False,exception='Unauthorized') #can only be called from btsperiphery contract
    # bts_periphery.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(5),1:sp.nat(2)})).run(sender=bts_periphery.address) #set token limit for Tok2 coin to 5 and BB coin to 2
    # sc.verify(bts_periphery.data.token_limit["Tok2"] == sp.nat(5))#test of token_limit for tok2 token
    # bts_periphery.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(5)} )).run(valid=False,exception='InvalidParams',sender=bts_periphery.address) #invalid parameters
    # #cannot set more than 15 token limit at once
    # bts_periphery.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(15),1:sp.nat(22)})).run(sender=bts_periphery.address) #can modify already set data
    # sc.verify(bts_periphery.data.token_limit["BB"] == sp.nat(22))#test of token_limit for tok2 token

    # # handle_relay_message
    # msg=sp.bytes("0xf90236f90233b8ddf8db01b8d3f8d1f8cfb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431417259755046465a596341766f396b71793251673167327667465965676f52714a04b88af888b8396274703a2f2f3078372e69636f6e2f637831666637646432636639373836316262653462666536386232663463313834376562666132663534b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431417259755046465a596341766f396b71793251673167327667465965676f52714a8362747381ff84c300f80084008449a0b90151f9014e01b90145f90142f9013fb8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431417259755046465a596341766f396b71793251673167327667465965676f52714a05b8faf8f8b8396274703a2f2f3078372e69636f6e2f637831666637646432636639373836316262653462666536386232663463313834376562666132663534b8406274703a2f2f4e6574586e486656716d39696573702e74657a6f732f4b5431417259755046465a596341766f396b71793251673167327667465965676f52714a8362747303b874f87200b86ff86daa687839643138316431336634376335616165353535623730393831346336623232393738373937363139a4747a3157615078716f375868556e56344c346669324457424e4e51384a6231777445716edcdb906274702d3078372e69636f6e2d4943588900d71b0fe0a28e000084008449bf")
    # prev=sp.string("btp://0x7.icon/cxff8a87fde8971a1d10d93dfed3416b0a6258d67b")
    # bmcPeriphery_contract.handle_relay_message(sp.record(prev=prev, msg=msg)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))

    # #set_fee_ratio
    # btsCore_contract.set_fee_ratio(name=sp.string("btp-NetXnHfVqm9iesp.tezos-XTZ"),fee_numerator=sp.nat(100),fixed_fee=sp.nat(450)).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))




def deploy_bmcManagement_contract(owner, helper):
    bmcManagement_contract = BMCManagement.BMCManagement(owner, helper)
    return bmcManagement_contract

def deploy_bmcPeriphery_contract(bmc_addres, helper, parse,owner):
    bmcPeriphery_contract = BMCPeriphery.BMCPreiphery(bmc_addres, helper,helper, parse, owner)
    return bmcPeriphery_contract

def deploy_helper_contract():
    helper_contract = BMCHelper.Helper()
    return helper_contract


def deploy_parse_address():
    parse_address = ParseAddress.ParseAddress()
    return parse_address


def deploy_btsCore_contract(bts_OwnerManager_Contract):
    btsCore_contract = BTSCore.BTSCore(
        owner_manager=bts_OwnerManager_Contract,
        _native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
        _fee_numerator=sp.nat(100),
        _fixed_fee=sp.nat(450)
    )
    return btsCore_contract

def deploy_btsOwnerManager_Contract(owner):
    bts_OwnerManager_Contract = BTSOwnerManager.BTSOwnerManager(owner)
    return bts_OwnerManager_Contract

def deploy_btsPeriphery_Contract(core_address, helper, parse, bmc,owner):
    btsPeriphery_Contract = BTSPeriphery.BTPPreiphery(bmc_address= bmc, bts_core_address=core_address, helper_contract=helper, parse_address=parse, owner_address=owner, native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ")
    return btsPeriphery_Contract

def deploy_fa2_Contract(creator):
        fa2_contract = BTSCore.FA2_contract.SingleAssetToken(admin=creator, metadata=sp.big_map({"ss": sp.bytes("0x0dae11")}), token_metadata=sp.map({"ff": sp.bytes("0x0dae11")}))
        return fa2_contract


#Core - function with interscore call

# Transfer_fees - called from bts_periphery - checked from bts periphery handle fee gathering
# Handle_response_service- called from bts_periphery 
# Mint - called from bts_periphery
# Refund - refunds 
# Transfer batch -
# Send service message
# Transfer
# Transfer native coin

# Periphery
# Send service message - bts core
# Handle btp message - only bmc - done
# Handle btp error- only bmc -calls handle response service of periphery - done
# Handle response service -called from handle_btp_message-calls core handle response service
# Handle fee gathering - only bmc - calls core transfer_fees
# Send response message - called from handle btp message/ calls bmc send message
# –Handle request service - called from handle btp message /calls core mint - can be called for mint