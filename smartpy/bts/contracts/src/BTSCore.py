import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
Coin = sp.TRecord(addr=sp.TAddress, feeNumerator=sp.TNat, fixedFee=sp.TNat,
                  coinType=sp.TNat)


class BTSCore(sp.Contract):
    FEE_DENOMINATOR = sp.nat(10000)
    RC_OK = sp.nat(0)
    RC_ERR = sp.nat(1)
    NATIVE_COIN_TYPE = sp.nat(0)
    NATIVE_WRAPPED_COIN_TYPE = sp.nat(1)
    NON_NATIVE_TOKEN_TYPE = sp.nat(2)

    MAX_BATCH_SIZE = sp.nat(15)
    NATIVE_COIN_ADDRESS = sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc")

    # TODO: change the native coin addr

    def __init__(self, ownerManager_address, btsPeriphery_address, _nativeCoinName, _feeNumerator, _fixedFee,
                 ):
        # Sets the initial storage of the contract with an empty map and a list containing the owner address
        self.update_initial_storage(

            ownerManager_contract_address=ownerManager_address,
            btsPeriphery_contract_address=btsPeriphery_address,
            nativeCoinName=_nativeCoinName,
            listOfOwners=sp.list(t=sp.TAddress),
            # a list of amounts have been charged so far (use this when Fee Gathering occurs)
            chargedAmounts=sp.list(t=sp.TNat),
            # a string array stores names of supported coins
            coinsName=sp.list([_nativeCoinName], t=sp.TString),
            #  a list of coins' names have been charged so far (use this when Fee Gathering occurs)
            chargedCoins=sp.list(t=sp.TString),

            owners=sp.map({}, tkey=sp.TAddress, tvalue=sp.TBool),
            aggregationFee=sp.map({}, tkey=sp.TString, tvalue=sp.TInt),
            balances=sp.big_map(tkey=sp.TRecord(sp.TAddress,sp.TString), tvalue= types.Type.Balance),
            coins=sp.map({_nativeCoinName: self.NATIVE_COIN_ADDRESS}, tkey=sp.TString, tvalue=sp.TAddress),

            coinDetails=sp.map({_nativeCoinName: sp.record(addr=self.NATIVE_COIN_ADDRESS,
                                                           feeNumerator=_feeNumerator,
                                                           fixedFee=_fixedFee,
                                                           coinType=sp.nat(0))},
                               tkey=sp.TString, tvalue=Coin),

            coinsAddress=sp.map({}, tkey=sp.TAddress, tvalue=sp.TString),
        )
    #is this necessary? can we check against owners map in line 37?
    def OnlyOwner(self, _owner):
        # call Owner Manager Contract for checking owner
        isOwner = sp.view("isOwner", self.data.ownerManager_contract_address, sp.sender, t=sp.TBool).open_some(
            "OwnerNotFound")

        # ToDO: find a way to transfer function parameter to another contract

        sp.verify(isOwner == sp.sender, message="Unauthorized")


    # Get the name of Native Coin , Caller can be any
    @sp.onchain_view()
    def getNativeCoinName(self):
        sp.result(self.data.nativeCoinName)

    # update BTS Periphery address
    # TODO: verify zero address
    @sp.entry_point
    def updateBTSPeriphery(self, _btsPeripheryAddr):
        sp.set_type(_btsPeripheryAddr, sp.TAddress)
        sp.verify(self.data.owners[sp.sender], message="Unauthorized")
        sp.verify(_btsPeripheryAddr != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"), message="InvalidSetting")
        sp.if(self.data.btsPeriphery_contract_address != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc")):
        btsPeriphery_contract = sp.view("hasPendingRequest", self.data.btsPeriphery_contract_address, sp.none,
                                        t=sp.TBool).open_some("OwnerNotFound")

        sp.verify(btsPeriphery_contract == False, message="HasPendingRequest")
        self.data.btsPeriphery_contract_address = _btsPeripheryAddr

    #set fee ratio, Caller must be the owner of this contract
    #The transfer fee is calculated by feeNumerator/FEE_DEMONINATOR.
    #_feeNumerator if it is set to `10`, which means the default fee ratio is 0.1%.
    @sp.entry_point
    def set_fee_ratio(self,_name,_fee_numerator,_fixed_fee):
        sp.set_type(_name, sp.TString)
        sp.set_type(_fee_numerator, sp.TNat)
        sp.set_type(_fixed_fee, sp.TNat)
        sp.verify(self.data.owners[sp.sender], message="Unauthorized")
        sp.verify(_fee_numerator < self.FEE_DENOMINATOR, message="InvalidSetting")
        sp.verify((_name == self.data.nativeCoinName) |
                  (self.data.coins[_name] != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc")),
                  message = "TokenNotExists")
        sp.verify((_fixed_fee > 0) & (_fee_numerator >= 0), message = "LessThan0")
        self.data.coinDetails[_name].feeNumerator = _fee_numerator
        self.data.coinDetails[_name].fixedFee = _fixed_fee


    @sp.entry_point
    def register(self, _name, _symbol, _decimals, _fee_numerator, _fixed_fee, _addr ):
        sp.set_type(_name, sp.TString)
        sp.set_type(_symbol, sp.TString)
        sp.set_type(_decimals, sp.TNat)
        sp.set_type(_fee_numerator, sp.TNat)
        sp.set_type(_fixed_fee, sp.TNat)
        sp.set_type(_addr, sp.TAddress)
        sp.verify(self.data.owners[sp.sender], message="Unauthorized")
        sp.verify(_name == self.data.nativeCoinName, message= "ExistNativeCoin")
        sp.verify(self.data.coins[_name] == sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"), message= "ExistCoin")
        sp.verify(_fee_numerator <= self.FEE_DENOMINATOR, message="InvalidSetting")
        sp.verify((_fixed_fee > 0) & (_fee_numerator >= 0), message="LessThan0")
    #     TODO: confirm zero addr for tezos
        sp.if _addr == sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"):
            # deploy FA2 contract and set the deployed address
            # deployedFA2 = sp.address()
            self.data.coins[_name] = deployedFA2
            self.data.coinsName.push(_name)
            self.data.coinsAddress[deployedFA2] = _name
            self.data.coinDetails[_name] = Coin (
                addr = deployedFA2,
                feeNumerator = _fee_numerator,
                fixedFee = _fixed_fee,
                coinType = self.NATIVE_WRAPPED_COIN_TYPE
            )
        sp.else:
            self.data.coins[_name] = _addr
            self.data.coinsName.push(_name)
            self.data.coinsAddress[_addr] = _name
            self.data.coinDetails[_name] = Coin (
                addr = _addr,
                feeNumerator = _fee_numerator,
                fixedFee = _fixed_fee,
                coinType = self.NON_NATIVE_TOKEN_TYPE
            )
    # ToDO: initialise string and make interscore call.
        token_arr = sp.list[_name]
        val_arr = sp.list[]
        #TODO: confirm the following interscore call is correct or not
        set_token_limit_args_type = sp.TRecord(coin_names = sp.TMap(sp.TNat, sp.TString), token_limit =sp.TMap(sp.TNat, sp.TNat))
        set_token_limit_entry_point = sp.contract(set_token_limit_args_type, self.data.btsPeriphery_contract_address,"setTokenLimit").open_some()
        set_token_limit_args = sp.record(coin_names = token_arr, token_limit = val_arr)
        sp.transfer(set_token_limit_args,sp.tez(0),set_token_limit_entry_point)


# Return all supported coins names , Caller can be any
    @sp.onchain_view()
    def coinNames(self):
        sp.result(self.data.coinsName)

    # Return the _id number of Coin whose name is the same with given _coinName
    @sp.onchain_view()
    def coinId(self, _coinName):
        sp.result(self.data.coins[_coinName])


    @sp.onchain_view()
    def isValidCoin(self, _coin_name):
        sp.if (self.data.coins[_coin_name] != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"))|( _coin_name == self.data.nativeCoinName):
            sp.result(True)
        sp.else:
            sp.result(False)

    @sp.onchain_view()
    def feeRatio(self, _coin_name):
        coin = self.data.coinDetails[_coin_name]
        _fee_numerator = coin.feeNumerator
        _fixed_fee = coin.fixedFee
        sp.result(sp.record(fee_numerator =_fee_numerator,fixed_fee = _fixed_fee))


    @sp.onchain_view()
    # TODO: confirm sp.balance_Of()
    def balanceOf(self, _owner, _coin_name):
        sp.if _coin_name == self.data.nativeCoinName:
            sp.result(sp.record(_usable_balance = sp.nat(0),
                                _locked_balane = self.data.balances[_owner][_coin_name].lockedBalance,
                                _refundable_balance = self.data.balances[_owner][_coin_name].refundableBalance,
                                _user_balance = sp.balance_of(sp.address(_owner))))
        _fa2_address = self.data.coins[_coin_name]
            # IERC20 ierc20 = IERC20(_erc20Address);
    #TODO: userbalance = token balance of a user?
            # allowance?
        sp.if _fa2_address != self.data.zero_addr:
                # return token balance of _owner
                # _user_balance =
        sp.if _fa2_address == self.data.zero_addr:
                _user_balance = sp.nat(0)
                sp.result(_user_balance)
        # TODO: userbalance and allowance operations
        sp.result(sp.record(_usable_balance,
                            _locked_balane=self.data.balances[_owner][_coin_name].lockedBalance,
                            _refundable_balance=self.data.balances[_owner][_coin_name].refundableBalance,
                            _user_balance))


    @sp.onchain_view()
    def balanceOfBatch(self, _owner, _coin_names):
        sp.verify((sp.len(_coin_names) >0 ) & (sp.len(_coin_names) <= self.MAX_BATCH_SIZE), message = "BatchMaxSizeExceed")
        _usable_balances =sp.local("_usable_balances", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))
        _locked_balances =sp.local("_locked_balances", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))
        _refundable_balances =sp.local("_refundable_balances", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))
        _user_balances =sp.local("_user_balances", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))
        sp.for i in sp.range(0, sp.len(_coin_names)):
            (_usable_balances[i],
             _locked_balances[i],
             _refundable_balances[i],
             _user_balances[i]) = balanceOf(_owner, _coin_names[i]) #Confirm this part
        sp.result(sp.record(_usable_balances, _locked_balances, _refundable_balances, _user_balances))

    @sp.onchain_view()
    #TODO: dynamic array?
    def getAccumulatedFees():
        sp.for i in sp.range(0, sp.len( self.data.coinsName)):
            _accumulatedFees[i] = Types.Asset(self.data.coinsName[i], aggregationFee[self.data.coinsName[i]])
            sp.result(_accumulatedFees)

    @sp.entry_point
    def transferNativeCoin (self, _to):
        #TODO: confirm data type for amount
        check_transfer_restrictions = sp.contract(sp.string,sp.address,sp.nat, self.data.btsPeriphery_contract_address,"checkTransferRestrictions").open_some()
        sp.transfer(self.data.nativeCoinName,sp.sender, sp.amount, check_transfer_restrictions)
        #  Aggregation Fee will be charged on BSH Contract
        #  `fixedFee` is introduced as charging fee
        # charge_amt = fixedFee + msg.value * feeNumerator / FEE_DENOMINATOR
        charge_amt = sp.amount
                     * self.data.coinDetails[self.data.nativeCoinName].feeNumerator
                    //self.FEE_DENOMINATOR
                    + self.data.coinDetails[self.data.nativeCoinName].fixedFee #Confirm the type for this calculation

        # @dev sp.sender is an amount request to transfer (include fee)
        # Later on, it will be calculated a true amount that should be received at a destination
        _sendServiceMessage(sp.sender, _to, self.data.coinsName[0],sp.amount, charge_amt)

    @sp.entry_point
    def transfer(self, _coin_name, value, _to):
        sp.verify(_coin_name == self.data.nativeCoinName, message="InvalidWrappedCoin")
        _fa2_address = self.data.coins[_coinName]
        sp.verify(_fa2_address != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"), message= "CoinNotRegistered")
        check_transfer_restrictions_args_type = sp.TRecord(_coinName =sp.Tstring, _user=sp.TAddress , _value=sp.TNat)
        check_transfer_restrictions_entry_point = sp.contract(check_transfer_restrictions_args_type, self.data.btsPeriphery_contract_address,"checkTransferRestrictions").open_some()
        check_transfer_restrictions_args = sp.record(_coinName = _coin_name,_user= sp.sender,_value = sp.amount )
        sp.transfer(check_transfer_restrictions_args,sp.tez(0), check_transfer_restrictions_entry_point)
        charge_amt = sp.amount
                        * self.data.coinDetails[_coin_name].feeNumerator
                        //self.FEE_DENOMINATOR
                        + self.data.coinDetails[_coin_name].fixedFee
        
        #TODO: implement transferFrom function of fa2contract

        _sendServiceMessage(sp.sender, _to, self.data.coinsName[0],sp.amount, charge_amt)


    def _sendServiceMessage(self, _from, _to, _coin_name, _value, _charge_amt)
        sp.set_type(_from, sp.TAddress)
        sp.set_type(_to, sp.TAddress)
        sp.set_type(_coin_name, sp.TString)
        sp.set_type(_value, sp.TNat)
        sp.set_type(_charge_amt, sp.TNat)
        sp.verify(_value > _charge_amt, message = "ValueGreaterThan0")
        lockBalance(_from, _coin_name, _value)

        _coins= sp.local("_coins", {0 : "_coinName" }, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TString))
        _amounts= sp.local("_amounts", {0 : sp.as_nat(_value - _charge_amt)}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))
        _fees= sp.local("_fees", {0: _charge_amt}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))

#TODO: confirm the following interscore call is correct or not
        send_service_message_args_type = sp.TRecord(_from = sp.TAddress, to = sp.TAddress, coin_names = sp.TMap(sp.TNat, sp.Tstring), values = sp.TMap(sp.TNat,sp.TNat), fees = sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type, self.data.btsPeriphery_contract_address,"sendServiceMessage").open_some()
        send_service_message_args = sp.record(_from = _from, to = _to, coin_names = _coins, values = _amounts, fees = _fees)
        sp.transfer(send_service_message_args,sp.tez(0),send_service_message_entry_point)        

    @sp.entry_point
    def transferBatch(self, _coin_names, _values, _to)
        sp.set_type(_coin_names, sp.TMap(tkey = TNat, tvalue = TString))
        sp.set_type(_values, sp.TMap(tkey = TNat, tvalue = TNat))
        sp.set_type(_to, sp.Taddress)
        sp.verify(sp.len(_coin_names) == sp.len(_values), message ="InvalidRequest")
        sp.verify(sp.len(_coin_names) >0, message = "Zero length arguments")
        sp.if sp.amount != sp.nat(0):
            size = sp.len(_coin_names)+sp.nat(1)
        sp.if sp.amount == sp.nat(0):
            size = sp.len(_coin_names)
        sp.verify(size <= self.MAX_BATCH_SIZE, message ="InvalidRequest")
    
        _coins= sp.local("_coins", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TString))
        _amounts= sp.local("_amounts", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))
        _charge_amts= sp.local("_charge_amts", {}, t=sp.TMap(tkey=sp.TNat, tvalue=sp.TNat))

        Coin _coin
        coin_name = sp.local("coin_name", "", t= TString)
        value = sp.local("value", 0, t= tNat)
        
        sp.for i in sp.range(0, sp.len(_coin_names)):
            _fa2_addresses = self.data.coins[_coin_names[i]]
            sp.verify(_fa2_addresses != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"), message= "CoinNotRegistered")
            coin_name.value = _coin_names[i]
            value = _values[i]
            sp.verify(value > 0, message ="ZeroOrLess")
            
            check_transfer_restrictions_args_type = sp.TRecord(_coinName =sp.Tstring, _user=sp.TAddress , _value=sp.TNat)
            check_transfer_restrictions_entry_point = sp.contract(check_transfer_restrictions_args_type, self.data.btsPeriphery_contract_address,"checkTransferRestrictions").open_some()
            check_transfer_restrictions_args = sp.record(_coinName = coin_name,_user= sp.sender,_value = sp.amount )
            sp.transfer(check_transfer_restrictions_args,sp.tez(0), check_transfer_restrictions_entry_point)

            #TODO: implement transferFrom function of fa2contract

            _coin = self.data.coinDetails[coin_name]
            _coins[i] = coin_name
            _charge_amts[i] = value
                            *_coin.feeNumerator
                            //self.FEE_DENOMINATOR
                            + _coin.fixedFee
            _amounts[i] = value - _charged_amts[i] #set proper data type for this calculation

            lockBalance(sp.sender, coin_name, value)

        sp.if sp.amount !=sp.nat(0):
            check_transfer_restrictions_args_type = sp.TRecord(_coinName =sp.Tstring, _user=sp.TAddress , _value=sp.TNat)
            check_transfer_restrictions_entry_point = sp.contract(check_transfer_restrictions_args_type, self.data.btsPeriphery_contract_address,"checkTransferRestrictions").open_some()
            check_transfer_restrictions_args = sp.record(_coinName = self.data.nativeCoinName,_user= sp.sender,_value = sp.amount )
            sp.transfer(check_transfer_restrictions_args,sp.tez(0), check_transfer_restrictions_entry_point)

            _coins[size - 1] = self.data.nativeCoinName
            _charge_amts[size -1] = sp.amount
                        * self.data.coinDetails[_coin_name].feeNumerator
                        //self.FEE_DENOMINATOR
                        + self.data.coinDetails[_coin_name].fixedFee
            
            _amounts[size -1] = sp.amount - _charged_amts[size -1]

            lockBalance(sp.sender, self.data.nativeCoinName, value)

        
        #TODO: confirm the following interscore call is correct or not
        send_service_message_args_type = sp.TRecord(_from = sp.TAddress, to = sp.TAddress, coin_names = sp.TMap(sp.TNat, sp.Tstring), values = sp.TMap(sp.TNat,sp.TNat), fees = sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type, self.data.btsPeriphery_contract_address,"sendServiceMessage").open_some()
        send_service_message_args = sp.record(_from = sp.sender, to = _to, coin_names = _coins, values = _amounts, fees = _charge_amts)
        sp.transfer(send_service_message_args,sp.tez(0),send_service_message_entry_point)
    
    @sp.entry_point
    #TODO: implement nonReentrant
    def reclaim(self, _coin_name, _value)
         sp.verify(self.data.balances[sp.sender][_coin_name].refundableBalance >= _value, message="Imbalance")
         self.data.balances[sp.sender][_coin_name].refundableBalance = self.data.balances[sp.sender][_coin_name].refundableBalance -_value
         refund(sp.sender, _coin_name, _value)

    @sp.entry_point
    #TODO: ref doc of solidity for this method to set as public or private
    def refund(self, _to, _coin_name, _value)
         #here address(this) refers to the address of this contract
         sp.verify(sp.sender == sp.self_address, message="Unauthorized")
         sp.if _coin_name == self.data.nativeCoinName:
              paymentTransfer(_to, _value)
        sp.else: 
        #TODO: implement transfer on fa2

    def paymentTransfer(self, _to, _amount)
    #TODO: implement the following:
    # (bool sent,) = _to.call{value : _amount}("");
    #     require(sent, "PaymentFailed");

    @sp.entry_point
    def mint(_to, _coin_name, _value)
        sp.verify(sp.sender == self.data.btsPeriphery_contract_address, message="Unauthorized")
        sp.if _coin_name == self.data.nativeCoinName:
              paymentTransfer(_to, _value)
        sp.if self.data.coinDetails[_coin_name].coinType == self.NATIVE_WRAPPED_COIN_TYPE:
            #TODO : implement mint?
            #  IERC20Tradable(coins[_coinName]).mint(_to, _value)
        sp.else self.data.coinDetails[_coinName].coinType == self.NON_NATIVE_TOKEN_TYPE:
            #TODO: implement transfer
            # IERC20(coins[_coinName]).transfer(_to, _value)


    @sp.entry_point
    def handleResponseService( self, _requester, _coin_name, _value, _fee, _rsp_code)    
        sp.verify(sp.sender == self.data.btsPeriphery_contract_address, message="Unauthorized")
        sp.if _requester == sp.self_address:
             sp.if _rsp_code == self.RC_ERR:
                  self.data.aggregationFee[_coin_name] = self.data.aggregationFee[_coin_name] + _value

        _amount = _value + _fee
        self.data.balances[_requester][_coin_name].lockedBalance = self.data.balances[_requester][_coin_name].lockedBalance -_amount
        sp.if _rsp_code == self.RC_ERR:
             _fa2_address = self.data.coins[_coin_name]
             sp.if (_coin_name != self.data.nativeCoinName) & (self.data.coinDetails[_coin_name].coinType == self.NATIVE_WRAPPED_COIN_TYPE):
                  #TODO:implement burn
                  #IERC20Tradable(_erc20Address).burn(address(this), _value)

        self.data.aggregationFee[_coin_name] = self.data.aggregationFee[_coin_name] + _fee
        
    
    @sp.entry_point
    #here _fa is an address, not sure why its in string
    def transferFees(self, _fa)
        sp.verify(sp.sender == self.data.btsPeriphery_contract_address, message="Unauthorized")
        sp.for i in sp.range(0, sp.len(self.data.coinsName)):
              sp.if self.data.aggregationFee[self.data.coinsName[i]] != sp.nat(0)
              self.data.chargedCoins.push(self.data.coinsName[i])
              self.data.chargedAmounts.push(self.data.aggregationFee[self.data.coinsName[i]])
              del self.data.aggregationFee[self.data.coinsName[i]]
        #TODO: confirm data type for amount, address(this),refer solidity doc
        send_service_message= sp.contract(sp.address,sp.address, sp.list, sp.list, sp.list,self.data.btsPeriphery_contract_address,"sendServiceMessage").open_some()
        sp.transfer(sp.self_address, _fa, chargedCoins, chargedAmounts, [], send_service_message )

        #TODO: confirm the following interscore call is correct or not + check for fees.
        send_service_message_args_type = sp.TRecord(_from = sp.TAddress, to = sp.TAddress, coin_names = sp.TMap(sp.TNat, sp.Tstring), values = sp.TMap(sp.TNat,sp.TNat), fees = sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type, self.data.btsPeriphery_contract_address,"sendServiceMessage").open_some()
        send_service_message_args = sp.record(_from = sp.self_address, to = _fa, coin_names = chargedCoins, values = chargedAmounts, fees = [])
        sp.transfer(send_service_message_args,sp.tez(0),send_service_message_entry_point)

        del self.data.chargedCoins
        del self.data.chargedAmounts

    def lockBalance(self, _to, _coin_name, _value)
          self.data.balances[_to][_coin_name].lockedBalance = self.data.balances[_to][_coin_name].lockedBalance + _value
                             
    @sp.entry_point    
    def updateCoinDb()
        sp.verify(onlyOwner(sp.sender) == True, message= "Unauthorized")
        self.data.coins[self.data.nativeCoinName] = sp.address(self.NATIVE_COIN_ADDRESS)
        self.data.coinsAddress[sp.address(self.NATIVE_COIN_ADDRESS)] = self.data.nativeCoinName
        coins_length = sp.len(self.data.coinsLength)
        sp.for i in sp.range(0, coins_length):
             sp.if self.data.coinsName[i] != self.data.nativeCoinName:
                  self.data.coinsAddress[self.data.coinDetails[self.data.coinsName[i]].addr] = self.data.coinsName[i]

    @sp.entry_point
    def setBTSOwnerManager(_ownerManager)
        sp.verify(self.data.owners[sp.sender] == True , message= "Unauthorized")
        sp.verify(_ownerManager != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"), message= "InvalidAddress")
        #TODO :set addr self.data.btsOwnerManager = _ownerManager

@sp.add_test(name="BTSCore")
def test():
    c1 = BTSCore(
        ownerManager_address=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        btsPeriphery_address=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        _nativeCoinName="NativeCoin",
        _feeNumerator=sp.nat(1000),
        _fixedFee=sp.nat(10)
    )
    scenario = sp.test_scenario()
    scenario.h1("BTSCore")
    scenario += c1


sp.add_compilation_target("", BTSCore(
    ownerManager_address=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
    btsPeriphery_address=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
    _nativeCoinName="NativeCoin",
    _feeNumerator=sp.nat(1000),
    _fixedFee=sp.nat(10)
))






