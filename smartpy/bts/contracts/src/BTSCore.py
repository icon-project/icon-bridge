import smartpy as sp

# types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
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

    def OnlyOwner(self, _owner):
        # call Owner Manager Contract for checking owner
        isOwner = sp.view("isOwner", self.data.ownerManager_contract_address, sp.sender, t=sp.TBool).open_some(
            "OwnerNotFound")

        # ToDO: find a way to transfer function parameter to another contract

        sp.verify(isOwner == sp.sender, message="Unauthorized")

    def OnlyBtsPeriphery(self, _owner):
        sp.verify(sp.sender == self.data.btsPeriphery_contract_address, message="Unauthorized")

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
        # sp.if(self.data.btsPeriphery_contract_address != sp.none):
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
    #     TODO:
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
        set_token_limit = sp.contract(sp.list,sp.list, self.data.btsPeriphery_contract_address,"setTokenLimit").open_some()
        sp.transfer(token_arr, val_arr, sp.tez(0), setTokenLimit)


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
                                _user_balance = sp.balance_of(sp.address(_owner)))
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
                            _user_balance)


    @sp.onchain_view()
    def balanceOfBatch(self, _owner, _coin_names):
        sp.verify((sp.len(_coin_names) >0 ) & (sp.len(_coin_names) <= self.MAX_BATCH_SIZE), message = "BatchMaxSizeExceed")
        #ToDO: make a dynamic array
        sp.for i in sp.range(0, sp.len(_coin_names)):
            #TODO

    @sp.onchain_view()
    #TODO: dynamic array?
    def getAccumulatedFees():
        sp.for i in sp.range(0, sp.len( self.data.coinsName)):
            _accumulatedFees[i] = Types.Asset(self.data.coinsName[i], aggregationFee[self.data.coinsName[i])
            sp.result(_accumulatedFees)

@sp.add_test(name="Calculator")
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






