import smartpy as sp
BTSPeriphery = sp.io.import_script_from_url("file:./contracts/src/bts_periphery.py")
BTSCore = sp.io.import_script_from_url("file:./contracts/src/bts_core.py")
BTSOwnerManager = sp.io.import_script_from_url("file:./contracts/src/bts_owner_manager.py")
ParseAddress = sp.io.import_script_from_url("file:./contracts/src/parse_address.py")
Helper= sp.io.import_script_from_url("file:./contracts/src/helper.py")


@sp.add_test("BTSPeripheryTest")
def test():
    sc = sp.test_scenario()

    # test account
    alice=sp.test_account("Alice")
    bmc_address = sp.test_account('bmc')
    admin=sp.test_account('admin')
    helper = sp.test_account("Helper")
    

    def deploy_btsperiphery_contract():
     btsperiphery_contract = BTSPeriphery.BTPPreiphery(bmc_address.address,btscore_contract.address,btshelpercontract.address,btsparsecontract.address,'NativeCoin',admin.address)
     return btsperiphery_contract
    
    def deploy_parsecontract():
       btsparsecontract = ParseAddress.ParseAddress()
       return btsparsecontract
    
    def deploy_helper():
       btshelper = Helper.Helper()
       return btshelper
    
       
    
    def deploy_btscore_contract():
        btscore_contract= BTSCore.BTSCore(
         owner_manager= bts_OwnerManager_contract.address,
         _native_coin_name="Tok1",
         _fee_numerator=sp.nat(1000),
         _fixed_fee=sp.nat(10))
        return btscore_contract
    
    def deploy_btsOwnerManager_Contract():
     bts_OwnerManager_Contract = BTSOwnerManager.BTSOwnerManager(sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))
     return bts_OwnerManager_Contract
    
    bts_OwnerManager_contract = deploy_btsOwnerManager_Contract()
    sc+= bts_OwnerManager_contract

    btscore_contract = deploy_btscore_contract()
    sc+= btscore_contract

    btshelpercontract = deploy_helper()
    sc+= btshelpercontract

    btsparsecontract = deploy_parsecontract()
    sc+= btsparsecontract
    
    # deploy btsperiphery contract
    btsperiphery_contract = deploy_btsperiphery_contract()
    sc += btsperiphery_contract
    
    btscore_contract.update_bts_periphery(btsperiphery_contract.address).run(sender=sp.address("tz1XGbmLYhqcigxFuBCJrgyJejnwkySE4Sk9"))


    

    #test 1: Test of add to blacklist function
    btsperiphery_contract.add_to_blacklist({0:"notaaddress"}).run(sender=btsperiphery_contract.address,valid=False, exception="InvalidAddress") # invalid address 
    btsperiphery_contract.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=btsperiphery_contract.address,valid=True) #add a address to blacklist
    btsperiphery_contract.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}).run(sender=btsperiphery_contract.address,valid=True)  # can be called twice 
    btsperiphery_contract.add_to_blacklist({0:"tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"}) .run(sender=admin.address,valid=False,exception ="Unauthorized")# only btsperiphery contract call this function
    btsperiphery_contract.add_to_blacklist({0:'tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg'}).run(sender=btsperiphery_contract.address,valid=False,exception='InvalidAddress')#invalid address
    sc.verify(btsperiphery_contract.data.blacklist[sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW")] == True) # checking the blacklist[] map
   

    #test 2 : Test of remove from blacklist function
    btsperiphery_contract.remove_from_blacklist({0:'notaaddress'}).run(sender=btsperiphery_contract.address,valid=False, exception="InvalidAddress") # invalid address 
    btsperiphery_contract.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=btsperiphery_contract.address,valid=False, exception="UserNotFound") # address not black-listed
    btsperiphery_contract.add_to_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=btsperiphery_contract.address) # adding to blacklist
    btsperiphery_contract.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=btsperiphery_contract.address) # valid process
    btsperiphery_contract.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=btsperiphery_contract.address,valid=False ,exception ='UserNotFound') # cannot remove from blacklist twice
    btsperiphery_contract.add_to_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=btsperiphery_contract.address) # adding to blacklist
    btsperiphery_contract.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=btsperiphery_contract.address) # can only be called from btseperiphery contract
    btsperiphery_contract.remove_from_blacklist({0:'tz1d2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW'}).run(sender=admin.address, valid=False, exception ="Unauthorized") # can only be called from btseperiphery contract


    #test 3 : set token  limit
    btsperiphery_contract.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(5),1:sp.nat(2)})).run(sender=admin.address,valid=False,exception='Unauthorized') #can only be called from btsperiphery contract
    btsperiphery_contract.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(5),1:sp.nat(2)})).run(sender=btsperiphery_contract.address) #set token limit for Tok2 coin to 5 and BB coin to 2
    sc.verify(btsperiphery_contract.data.token_limit["Tok2"] == sp.nat(5))#test of token_limit for tok2 token
    btsperiphery_contract.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(5)} )).run(valid=False,exception='InvalidParams',sender=btsperiphery_contract.address) #invalid parameters
    #cannot set more than 15 token limit at once
    btsperiphery_contract.set_token_limit(sp.record(coin_names = {0:"Tok2" , 1:'BB'} ,token_limit ={0:sp.nat(15),1:sp.nat(22)})).run(sender=btsperiphery_contract.address) #can modify already set data
    sc.verify(btsperiphery_contract.data.token_limit["BB"] == sp.nat(22))#test of token_limit for tok2 token


    #test 4 :send service message

    btsperiphery_contract.send_service_message(sp.record(_from=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), to="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW",coin_names={0:"Tok1"}, values={0:sp.nat(10)}, fees={0:sp.nat(2)})).run(sender=btscore_contract.address )
    btsperiphery_contract.send_service_message(sp.record(_from=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), to="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW",coin_names={0:"Tok1"}, values={0:sp.nat(10)}, fees={0:sp.nat(2)})).run( sender=btscore_contract.address ) # test of function
    btsperiphery_contract.send_service_message(sp.record(_from=sp.address("tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW"), to="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW",coin_names={0:"Tok1"}, values={0:sp.nat(10)}, fees={0:sp.nat(2)})).run( sender= admin,valid=False, exception='Unauthorized' ) # only message from bts-core is authorized
    sc.show(btsperiphery_contract.data.requests[1])#test to verify if request message is correct
    sc.verify_equal(btsperiphery_contract.data.requests[1] ,sp.record(amounts = {0 : 10}, coin_names = {0 : 'Tok1'}, fees = {0 : 2}, from_ = 'tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW', to = 'btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW')) # request data verified
    sc.verify(btsperiphery_contract.data.number_of_pending_requests == 2)#2 pending request
    sc.verify(btsperiphery_contract.data.serial_no == 2) #serial no of request increased to 2


    #Test 5: handle btp message
   #  btsperiphery_contract.handle_btp_message(sp.record(_from="tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="bts", sn=sp.nat(4),msg=sp.bytes("0xf8cfb83d6274703a2f2f30783232382e6172637469632f307831313131313131313131313131313131313131313131313131313131313131313131313131313131b8396274703a2f2f3078312e69636f6e2f637830303030303030303030303030303030303030303030303030303030303030303030303030303036836274730ab84ef84cb83d6274703a2f2f30783232382e6172637469632f307861313434326339303132304138393163336465393739336361433730393638436162313133323335cc8b627470206d657373616765") )).run(sender=bmc_address) #test of handle_btp_message function
   #  btsperiphery_contract.handle_btp_message(sp.record(_from="tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="btc", sn=sp.nat(4),msg=sp.bytes("0xf8cfb83d6274703a2f2f30783232382e6172637469632f307831313131313131313131313131313131313131313131313131313131313131313131313131313131b8396274703a2f2f3078312e69636f6e2f637830303030303030303030303030303030303030303030303030303030303030303030303030303036836274730ab84ef84cb83d6274703a2f2f30783232382e6172637469632f307861313434326339303132304138393163336465393739336361433730393638436162313133323335cc8b627470206d657373616765") )).run(sender=bmc_address,valid=False, exception='InvalidSvc') #svc name must match hardcoded service name "btc"
   #  btsperiphery_contract.handle_btp_message(sp.record(_from="tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="btc", sn=sp.nat(4),msg=sp.bytes("0xf8cfb83d6274703a2f2f30783232382e6172637469632f307831313131313131313131313131313131313131313131313131313131313131313131313131313131b8396274703a2f2f3078312e69636f6e2f637830303030303030303030303030303030303030303030303030303030303030303030303030303036836274730ab84ef84cb83d6274703a2f2f30783232382e6172637469632f307861313434326339303132304138393163336465393739336361433730393638436162313133323335cc8b627470206d657373616765") )).run(sender=btsperiphery_contract.address,valid=False, exception='Unauthorized') #can only be called from bmc contract




    



   


    #Test : handle btp error
   #  sc.verify(btsperiphery_contract.data.number_of_pending_requests == 2)#pending request is 2 here
   #  btsperiphery_contract.handle_btp_error(sp.record(svc= "bts", code=sp.nat(2), sn=sp.nat(1), msg="test 1")).run(sender=bmc_address)
   #  btsperiphery_contract.handle_btp_error(sp.record(svc= "btc", code=sp.nat(2), sn=sp.nat(1), msg="test 1")).run(sender=bmc_address,valid=False,exception='InvalidSvc') #Invalid Svc
   #  btsperiphery_contract.handle_btp_error(sp.record(svc= "bts", code=sp.nat(2), sn=sp.nat(111), msg="test 1")).run(sender=bmc_address,valid=False,exception='Missing item in map') # Invalid sn , sn must be serial number of service request
   #  btsperiphery_contract.handle_btp_error(sp.record(svc= "bts", code=sp.nat(2), sn=sp.nat(1), msg="test 1")).run(sender=btsperiphery_contract.address,valid=False,exception='Unauthorized')#Only bmc contract can call this fucntion
   #  sc.verify(btsperiphery_contract.data.number_of_pending_requests == 0) #pending request decreased


    #Test : handle request service
    btsperiphery_contract.handle_request_service(sp.record(to= "tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc", assets={0: sp.record(coin_name="BB", value=sp.nat(4))})).run(sender=btsperiphery_contract.address,valid=False,exception='UnregisteredCoin')


    #Test : handle fee gathering 
   #  btsperiphery_contract.handle_fee_gathering(sp.record(fa="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="bts")).run(sender=bmc_address) # handle_fee_gathering function call
   #  btsperiphery_contract.handle_fee_gathering(sp.record(fa="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="btc")).run(sender=bmc_address, valid=False, exception='InvalidSvc') # svc must match hardcoded service name 'bts'
   #  btsperiphery_contract.handle_fee_gathering(sp.record(fa="btp://77.tezos/tz1e2HPzZWBsuExFSM4XDBtQiFnaUB5hiPnW", svc="bts")).run(sender=btsperiphery_contract.address, valid=False, exception='Unauthorized') # can only be called from bmc contract
    



 

















   







    


