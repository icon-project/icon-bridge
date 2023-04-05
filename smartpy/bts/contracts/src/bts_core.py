import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
Coin = sp.TRecord(addr=sp.TAddress, fee_numerator=sp.TNat, fixed_fee=sp.TNat,
                  coin_type=sp.TNat)


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

    def __init__(self, _native_coin_name, _fee_numerator, _fixed_fee, owner_manager, bts_periphery_addr):
        self.update_initial_storage(
            bts_owner_manager=owner_manager,
            bts_periphery_address=bts_periphery_addr,
            native_coin_name=_native_coin_name,

            list_of_owners=sp.list(t=sp.TAddress),
            charged_amounts=sp.map(tkey=sp.TNat, tvalue=sp.TNat),
            coins_name=sp.list([_native_coin_name], t=sp.TString),
            charged_coins=sp.map(tkey=sp.TNat, tvalue=sp.TString),

            owners=sp.map({}, tkey=sp.TAddress, tvalue=sp.TBool),
            aggregation_fee=sp.map({}, tkey=sp.TString, tvalue=sp.TNat),
            balances=sp.big_map(tkey=sp.TRecord(address=sp.TAddress, coin_name=sp.TString), tvalue= types.Types.Balance),
            coins=sp.map({_native_coin_name: self.NATIVE_COIN_ADDRESS}, tkey=sp.TString, tvalue=sp.TAddress),
            coin_details=sp.map({_native_coin_name: sp.record(addr=self.NATIVE_COIN_ADDRESS,
                                                              fee_numerator=_fee_numerator,
                                                              fixed_fee=_fixed_fee,
                                                              coin_type=self.NATIVE_COIN_TYPE)},
                                tkey=sp.TString, tvalue=Coin),
            coins_address=sp.map({}, tkey=sp.TAddress, tvalue=sp.TString),
        )

    #is this necessary? can we check against owners map in line 37?
    def only_owner(self):
        # call Owner Manager Contract for checking owner
        is_owner = sp.view("is_owner", self.data.bts_owner_manager, sp.sender, t=sp.TBool).open_some(
            "OwnerNotFound")
        sp.verify(is_owner == True, message="Unauthorized")

    def only_bts_periphery(self):
        sp.verify(sp.sender == self.data.bts_periphery_address, "Unauthorized")

    @sp.onchain_view()
    def get_native_coin_name(self):
        """
        Get name of nativecoin
        :return: Name of nativecoin
        """
        sp.result(self.data.native_coin_name)

    @sp.entry_point
    def update_bts_periphery(self, bts_periphery):
        """
        update BTS Periphery address.
        :param bts_periphery:  BTSPeriphery contract address.
        :return:
        """
        sp.set_type(bts_periphery, sp.TAddress)

        self.only_owner()
        # TODO: verify zero address
        sp.verify(bts_periphery != sp.address("tz1000000"), message="InvalidSetting")
        sp.if self.data.bts_periphery_address != sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"):
            has_requests = sp.view("has_pending_request", self.data.bts_periphery_address, sp.none,
                                            t=sp.TBool).open_some("OwnerNotFound")
            sp.verify(has_requests == False, message="HasPendingRequest")
        self.data.bts_periphery_address = bts_periphery

    #set fee ratio, Caller must be the owner of this contract
    #The transfer fee is calculated by feeNumerator/FEE_DEMONINATOR.
    #_feeNumerator if it is set to `10`, which means the default fee ratio is 0.1%.
    @sp.entry_point
    def set_fee_ratio(self, name, fee_numerator, fixed_fee):
        """
        set fee ratio.
        :param name:
        :param fee_numerator: the fee numerator
        :param fixed_fee:
        :return:
        """
        sp.set_type(name, sp.TString)
        sp.set_type(fee_numerator, sp.TNat)
        sp.set_type(fixed_fee, sp.TNat)

        self.only_owner()
        sp.verify(fee_numerator < self.FEE_DENOMINATOR, message="InvalidSetting")
        sp.verify((name == self.data.native_coin_name) |
                  (self.data.coins[name] != sp.address("tz1000")),
                  message = "TokenNotExists")
        sp.verify((fixed_fee > sp.nat(0)) & (fee_numerator >= sp.nat(0)), message = "LessThan0")
        self.data.coin_details[name].fee_numerator = fee_numerator
        self.data.coin_details[name].fixed_fee = fixed_fee


    @sp.entry_point
    def register(self, name, symbol, decimals, fee_numerator, fixed_fee, addr):
        """
        Registers a wrapped coin and id number of a supporting coin.
        :param name: Must be different with the native coin name.
        :param symbol: symbol name for wrapped coin.
        :param decimals: decimal number
        :param fee_numerator:
        :param fixed_fee:
        :param addr: address of the coin
        :return:
        """
        sp.set_type(name, sp.TString)
        sp.set_type(symbol, sp.TString)
        sp.set_type(decimals, sp.TNat)
        sp.set_type(fee_numerator, sp.TNat)
        sp.set_type(fixed_fee, sp.TNat)
        sp.set_type(addr, sp.TAddress)

        self.only_owner()
        sp.verify(name != self.data.native_coin_name, message="ExistNativeCoin")
        sp.verify(self.data.coins[name] == sp.address("tz10000"), message= "ExistCoin")
        sp.verify(self.data.coins_address[addr] == "", message="AddressExists")
        sp.verify(fee_numerator <= self.FEE_DENOMINATOR, message="InvalidSetting")
        sp.verify((fixed_fee >= sp.nat(0)) & (fee_numerator >= sp.nat(0)), message="LessThan0")
    #     TODO: confirm zero addr for tezos
        with sp.if_(addr == sp.address("tz10000")):
            # TODO:deploy FA2 contract and set the deployed address
            deployed_fa2 = sp.address("tz10000")
            self.data.coins[name] = deployed_fa2
            self.data.coins_name.push(name)
            self.data.coins_address[deployed_fa2] = name
            self.data.coin_details[name] = sp.record(
                addr = deployed_fa2,
                fee_numerator = fee_numerator,
                fixed_fee = fixed_fee,
                coin_type = self.NATIVE_WRAPPED_COIN_TYPE
            )
        with sp.else_():
            self.data.coins[name] = addr
            self.data.coins_name.push(name)
            self.data.coins_address[addr] = name
            self.data.coin_details[name] = sp.record(
                addr = addr,
                fee_numerator = fee_numerator,
                fixed_fee = fixed_fee,
                coin_type = self.NON_NATIVE_TOKEN_TYPE
            )
        # ToDO: initialise string and make interscore call.
        token_map = sp.map({0:name})
        val_map = sp.map({0:1})

        # call set_token_limit on bts_periphery
        set_token_limit_args_type = sp.TRecord(coin_names=sp.TMap(sp.TNat, sp.TString), token_limit=sp.TMap(sp.TNat, sp.TNat))
        set_token_limit_entry_point = sp.contract(set_token_limit_args_type, self.data.bts_periphery_address, "set_token_limit").open_some()
        set_token_limit_args = sp.record(coin_names=token_map, token_limit=val_map)
        sp.transfer(set_token_limit_args, sp.tez(0), set_token_limit_entry_point)


    @sp.onchain_view()
    def coin_names(self):
       """
       Return all supported coins names
       :return: An array of strings.
       """
       sp.result(self.data.coins_name)

    @sp.onchain_view()
    def coin_id(self, coin_name):
        """
        Return address of Coin whose name is the same with given coin_ame.
        :param coin_name:
        :return: An address of coin_name.
        """
        sp.result(self.data.coins[coin_name])

    @sp.onchain_view()
    def is_valid_coin(self, coin_name):
        """
        Check Validity of a coin_name
        :param coin_name:
        :return: true or false
        """
        sp.result((self.data.coins[coin_name] != sp.address("tz10000"))|( coin_name == self.data.native_coin_name))


    @sp.onchain_view()
    def fee_ratio(self, coin_name):
        """
        Get fee numerator and fixed fee
        :param coin_name: Coin name
        :return: a record (Fee numerator for given coin, Fixed fee for given coin)
        """

        coin = self.data.coin_details[coin_name]
        fee_numerator = coin.fee_numerator
        fixed_fee = coin.fixed_fee
        sp.result(sp.record(fee_numerator =fee_numerator, fixed_fee = fixed_fee))


    @sp.onchain_view()
    def balance_of(self, params):
        """
        Return a usable/locked/refundable balance of an account based on coinName.
        usable_balance the balance that users are holding.
        locked_balance when users transfer the coin,it will be locked until getting the Service Message Response.
        refundable_balance refundable balance is the balance that will be refunded to users.
        :param params:
        :return: a record of (usable_balance, locked_balance, refundable_balance)
        """
        sp.set_type(params, sp.TRecord(owner=sp.TAddress, coin_name=sp.TString))

        # TODO: confirm sp.balance_Of() user_balance
        with sp.if_(params.coin_name == self.data.native_coin_name):
            sp.result(sp.record(usable_balance = sp.nat(0),
                                locked_balance = self.data.balances[sp.record(address=params.owner, coin_name=params.coin_name)].locked_balance,
                                refundable_balance = self.data.balances[sp.record(address=params.owner, coin_name=params.coin_name)].refundable_balance,
                                user_balance = sp.nat(2)))
        with sp.else_():
            fa2_address = self.data.coins[params.coin_name]
            usable_balance=sp.nat(1)
            # IERC20 ierc20 = IERC20(_erc20Address);
            #TODO: userbalance = token balance of a user?
                # allowance?
            sp.if fa2_address != sp.address("tz10000"):
                pass
                # return token balance of _owner
                # _user_balance =
            user_balance = sp.nat(0)
            sp.if fa2_address == sp.address("tz10000"):
                pass
            # TODO: userbalance and allowance operations
            sp.result(sp.record(usable_balance=usable_balance,
                                locked_balance=self.data.balances[sp.record(address=params.owner, coin_name=params.coin_name)].locked_balance,
                                refundable_balance=self.data.balances[sp.record(address=params.owner, coin_name=params.coin_name)].refundable_balance,
                                user_balance=user_balance))

    @sp.onchain_view()
    def balance_of_batch(self, params):
        """
        Return a list Balance of an account.
        :param params:
        :return: a record of (usableBalances, lockedBalances, refundableBalances)
        """
        sp.set_type(params, sp.TRecord(owner=sp.TAddress, coin_names=sp.TList(sp.TString)))


        sp.verify((sp.len(params.coin_names) > sp.nat(0)) & (sp.len(params.coin_names) <= self.MAX_BATCH_SIZE), message = "BatchMaxSizeExceed")
        usable_balances =sp.local("usable_balances", {}, t=sp.TMap(sp.TNat, sp.TNat))
        locked_balances =sp.local("locked_balances", {}, t=sp.TMap(sp.TNat, sp.TNat))
        refundable_balances =sp.local("refundable_balances", {}, t=sp.TMap(sp.TNat, sp.TNat))
        user_balances =sp.local("user_balances", {}, t=sp.TMap(sp.TNat, sp.TNat))

        sp.for item in params.coin_names:
            i = sp.local("i", sp.nat(0))
            balance= sp.view("balance_of", sp.self_address,
                    sp.record(owner=params.owner, coin_name=item)).open_some()
            usable_balances.value[i.value] = balance.usable_balance
            locked_balances.value[i.value] = balance.locked_balance
            refundable_balances.value[i.value] = balance.refundable_balance
            user_balances.value[i.value] = balance.user_balance
            i.value += sp.nat(1)
        sp.result(sp.record(usable_balances=usable_balances.value, locked_balances=locked_balances.value,
                            refundable_balances=refundable_balances.value, user_balances=user_balances.value))

    @sp.onchain_view()
    def get_accumulated_fees(self):
        """
        Return a map with record of accumulated Fees.
        :return: An map of Asset
        """

        accumulated_fees = sp.local("accumulated_fees", sp.map(tkey=sp.TNat, tvalue=types.Types.Asset))
        sp.for item in self.data.coins_name:
            i = sp.local("i", sp.nat(0))
            accumulated_fees.value[i.value] = sp.record(coin_name=item, value=self.data.aggregation_fee[item])
            i.value += sp.nat(1)
        sp.result(accumulated_fees.value)


    @sp.entry_point(check_no_incoming_transfer=False)
    def transfer_native_coin (self, to):
        """
        Allow users to deposit `sp.amount` native coin into a BTSCore contract.
        :param to: An address that a user expects to receive an amount of tokens.
        :return: 
        """
        sp.set_type(to, sp.TString)

        amount_in_nat = sp.local("amount_in_nat", sp.utils.mutez_to_nat(sp.amount), t=sp.TNat)

        # call check_transfer_restrictions on bts_periphery
        check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address,
                                 sp.record(coin_name=self.data.native_coin_name, user=sp.sender, value=amount_in_nat.value),
                                 t=sp.TBool).open_some()
        sp.verify(check_transfer == True, "FailCheckTransfer")

        #TODO: confirm data type for amount

        charge_amt = amount_in_nat.value * self.data.coin_details[self.data.native_coin_name].fee_numerator / self.FEE_DENOMINATOR + self.data.coin_details[self.data.native_coin_name].fixed_fee
        #Confirm the type for this calculation

        self._send_service_message(sp.sender, to, self.data.native_coin_name, amount_in_nat.value, charge_amt)

    @sp.entry_point
    def transfer(self, coin_name, value, to):
        """
        Allow users to deposit an amount of wrapped native coin `coin_name` from the `sp.sender` address into the BTSCore contract.
        :param coin_name: A given name of a wrapped coin
        :param value: An amount request to transfer from a Requester (include fee)
        :param to: Target BTP address.
        :return:
        """
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)
        sp.set_type(to, sp.TString)


        sp.verify(coin_name != self.data.native_coin_name, message="InvalidWrappedCoin")
        fa2_address = self.data.coins[coin_name]
        sp.verify(fa2_address != sp.address("tz10000"), message= "CoinNotRegistered")

        # call check_transfer_restrictions on bts_periphery
        check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address,
                                 sp.record(coin_name=coin_name, user=sp.sender, value=value),
                                 t=sp.TBool).open_some()
        sp.verify(check_transfer == True, "FailCheckTransfer")

        charge_amt = value * self.data.coin_details[coin_name].fee_numerator / self.FEE_DENOMINATOR + self.data.coin_details[coin_name].fixed_fee
        
        #TODO: implement transferFrom function of fa2contract

        self._send_service_message(sp.sender, to, coin_name, value, charge_amt)


    def _send_service_message(self, _from, to, coin_name, value, charge_amt):
        """
        This private function handles overlapping procedure before sending a service message to BTSPeriphery
        :param _from: An address of a Requester
        :param to: BTP address of of Receiver on another chain
        :param coin_name: A given name of a requested coin
        :param value: A requested amount to transfer from a Requester (include fee)
        :param charge_amt: An amount being charged for this request
        :return:
        """
        sp.set_type(_from, sp.TAddress)
        sp.set_type(to, sp.TString)
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)
        sp.set_type(charge_amt, sp.TNat)

        sp.verify(value > charge_amt, message = "ValueGreaterThan0")
        self._lock_balance(_from, coin_name, value)

        coins = sp.local("coins", {0 : coin_name}, t=sp.TMap(sp.TNat, sp.TString))
        amounts = sp.local("amounts", {0 : sp.as_nat(value - charge_amt)}, t=sp.TMap(sp.TNat, sp.TNat))
        fees = sp.local("fees", {0: charge_amt}, t=sp.TMap(sp.TNat, sp.TNat))

        # call send_service_message on bts_periphery
        send_service_message_args_type = sp.TRecord(_from = sp.TAddress, to = sp.TString, coin_names = sp.TMap(sp.TNat, sp.TString), values = sp.TMap(sp.TNat, sp.TNat), fees = sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type, self.data.bts_periphery_address, "send_service_message").open_some()
        send_service_message_args = sp.record(_from = _from, to = to, coin_names = coins.value, values = amounts.value, fees = fees.value)
        sp.transfer(send_service_message_args, sp.tez(0), send_service_message_entry_point)

    @sp.entry_point(check_no_incoming_transfer=False)
    def transfer_batch(self, coin_names, values, to):
        """
        Allow users to transfer multiple coins/wrapped coins to another chain.
        It MUST revert if the balance of the holder for token `_coinName` is lower than the `_value` sent.
        In case of transferring a native coin, it also checks `msg.value`
        The number of requested coins MUST be as the same as the number of requested values
        The requested coins and values MUST be matched respectively
        :param coin_names: A list of requested transferring wrapped coins
        :param values: A list of requested transferring values respectively with its coin name
        :param to: Target BTP address.
        :return:
        """
        sp.set_type(coin_names, sp.TMap(sp.TNat, sp.TString))
        sp.set_type(values, sp.TMap(sp.TNat, sp.TNat))
        sp.set_type(to, sp.TString)

        sp.verify(sp.len(coin_names) == sp.len(values), message ="InvalidRequest")
        sp.verify(sp.len(coin_names) > sp.nat(0), message = "Zero length arguments")

        amount_in_nat = sp.local("amount_in_nat", sp.utils.mutez_to_nat(sp.amount), t=sp.TNat)

        size = sp.local("size", sp.nat(0))
        sp.if amount_in_nat.value != sp.nat(0):
            size.value = sp.len(coin_names) + sp.nat(1)
        sp.if amount_in_nat.value == sp.nat(0):
            size.value = sp.len(coin_names)
        sp.verify(size.value <= self.MAX_BATCH_SIZE, message ="InvalidRequest")
    
        coins = sp.local("_coins", {}, t=sp.TMap(sp.TNat, sp.TString))
        amounts = sp.local("_amounts", {}, t=sp.TMap(sp.TNat, sp.TNat))
        charge_amts = sp.local("_charge_amts", {}, t=sp.TMap(sp.TNat, sp.TNat))


        # coin = sp.TRecord(addr=sp.TAddress, fee_numerator=sp.TNat, fixed_fee=sp.TNat, coin_type=sp.TNat)
        coin_name = sp.local("coin_name", "", t= sp.TString)
        value = sp.local("value", 0, t= sp.TNat)
        
        sp.for i in sp.range(sp.nat(0), sp.len(coin_names)):
            fa2_addresses = self.data.coins[coin_names[i]]
            sp.verify(fa2_addresses != sp.address("tz10000"), message= "CoinNotRegistered")
            coin_name.value = coin_names[i]
            value.value = values[i]
            sp.verify(value.value > sp.nat(0), message ="ZeroOrLess")

            # call check_transfer_restrictions on bts_periphery
            check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address,
                                     sp.record(coin_name=coin_name.value, user=sp.sender, value=value.value),
                                     t=sp.TBool).open_some()
            sp.verify(check_transfer == True, "FailCheckTransfer")

            #TODO: implement transferFrom function of fa2contract

            coin = sp.local("coin", self.data.coin_details[coin_name.value], t=Coin)
            coins.value[i] = coin_name.value
            charge_amts.value[i] = value.value *coin.value.fee_numerator //self.FEE_DENOMINATOR + coin.value.fixed_fee
            amounts.value[i] = sp.as_nat(value.value - charge_amts.value[i])

            self._lock_balance(sp.sender, coin_name.value, value.value)

        sp.if amount_in_nat.value !=sp.nat(0):
            # call check_transfer_restrictions on bts_periphery
            check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address,
                                     sp.record(coin_name=self.data.native_coin_name, user=sp.sender, value=amount_in_nat.value),
                                     t=sp.TBool).open_some()
            sp.verify(check_transfer == True, "FailCheckTransfer")

            coins.value[sp.as_nat(size.value - 1)] = self.data.native_coin_name
            charge_amts.value[sp.as_nat(size.value - 1)] = amount_in_nat.value * self.data.coin_details[coin_name.value].fee_numerator\
                                                     / self.FEE_DENOMINATOR + self.data.coin_details[coin_name.value].fixed_fee
            amounts.value[sp.as_nat(size.value - 1)] = sp.as_nat(sp.utils.mutez_to_nat(sp.amount) - charge_amts.value[sp.as_nat(size.value - 1)])

            self._lock_balance(sp.sender, self.data.native_coin_name, sp.utils.mutez_to_nat(sp.amount))

        
        # call send_service_message on bts_periphery
        send_service_message_args_type = sp.TRecord(_from=sp.TAddress, to=sp.TString,
                                                    coin_names=sp.TMap(sp.TNat, sp.TString),
                                                    values=sp.TMap(sp.TNat, sp.TNat),
                                                    fees=sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type,
                                                       self.data.bts_periphery_address,
                                                       "send_service_message").open_some()
        send_service_message_args = sp.record(_from=sp.sender, to=to, coin_names=coins.value, values=amounts.value,
                                              fees=charge_amts.value)
        sp.transfer(send_service_message_args, sp.tez(0), send_service_message_entry_point)


    @sp.entry_point
    #TODO: implement nonReentrant
    def reclaim(self, coin_name, value):
        """
        Reclaim the token's refundable balance by an owner.
        The amount to claim must be smaller or equal than refundable balance
        :param coin_name: A given name of coin
        :param value: An amount of re-claiming tokens
        :return:
        """
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)

        sp.verify(self.data.balances[sp.record(address=sp.sender,coin_name=coin_name)].refundable_balance >= value, message="Imbalance")
        self.data.balances[sp.record(address=sp.sender,coin_name=coin_name)].refundable_balance = sp.as_nat(self.data.balances[sp.record(address=sp.sender,coin_name=coin_name)].refundable_balance - value)
        self.refund(sp.sender, coin_name, value)

    def refund(self, to, coin_name, value):
        """

        :param to:
        :param coin_name:
        :param value:
        :return:
        """
        sp.verify(sp.sender == sp.self_address, message="Unauthorized")

        with sp.if_(coin_name == self.data.native_coin_name):
            self.payment_transfer(to, value)
        with sp.else_():
            pass
            #TODO: implement transfer on fa2

    def payment_transfer(self, to, amount):
        pass
        #TODO: implement the following:

        # (bool sent,) = _to.call{value : _amount}("");
        #     require(sent, "PaymentFailed");

    @sp.entry_point
    def mint(self, to, coin_name, value):
        """
        mint the wrapped coin.

        :param to: the account receive the minted coin
        :param coin_name: coin name
        :param value: the minted amount
        :return:
        """
        sp.set_type(to, sp.TAddress)
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)

        self.only_bts_periphery()
        sp.if coin_name == self.data.native_coin_name:
            self.payment_transfer(to, value)
        sp.if self.data.coin_details[coin_name].coin_type == self.NATIVE_WRAPPED_COIN_TYPE:
            pass
            #TODO : implement mint?
            #  IERC20Tradable(coins[_coinName]).mint(_to, _value)
        sp.if self.data.coin_details[coin_name].coin_type == self.NON_NATIVE_TOKEN_TYPE:
            pass
            #TODO: implement transfer
            # IERC20(coins[_coinName]).transfer(_to, _value)


    @sp.entry_point
    def handle_response_service(self, requester, coin_name, value, fee, rsp_code):
        """
        Handle a response of a requested service.
        :param requester: An address of originator of a requested service
        :param coin_name: A name of requested coin
        :param value: An amount to receive on a destination chain
        :param fee: An amount of charged fee
        :param rsp_code:
        :return:
        """
        sp.set_type(requester, sp.TAddress)
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)
        sp.set_type(fee, sp.TNat)
        sp.set_type(rsp_code, sp.TNat)

        self.only_bts_periphery()
        sp.if requester == sp.self_address:
            sp.if rsp_code == self.RC_ERR:
                  self.data.aggregation_fee[coin_name] = self.data.aggregation_fee[coin_name] + value
            return

        amount = sp.local("amount", value + fee, t=sp.TNat)
        self.data.balances[sp.record(address=requester, coin_name=coin_name)].locked_balance = sp.as_nat(self.data.balances[sp.record(address=requester, coin_name=coin_name)].locked_balance - amount.value)

        sp.if rsp_code == self.RC_ERR:
            # TODO: implement try catch
            fa2_address = self.data.coins[coin_name]
            sp.if (coin_name != self.data.native_coin_name) & (self.data.coin_details[coin_name].coin_type == self.NATIVE_WRAPPED_COIN_TYPE):
                  pass
                  #TODO:implement burn
                  #IERC20Tradable(_erc20Address).burn(address(this), _value)

        self.data.aggregation_fee[coin_name] = self.data.aggregation_fee[coin_name] + fee
        
    
    @sp.entry_point
    def transfer_fees(self, fa):
        """
        Handle a request of Fee Gathering. Usage: Copy all charged fees to an array
        :param fa:
        :return:
        """
        sp.set_type(fa, sp.TString)

        self.only_bts_periphery()
        sp.for item in self.data.coins_name:
            i = sp.local("i", sp.nat(0))
            sp.if self.data.aggregation_fee[item] != sp.nat(0):
                self.data.charged_coins[i.value] = item
                self.data.charged_amounts[i.value] = self.data.aggregation_fee[item]
                del self.data.aggregation_fee[item]
            i.value += sp.nat(1)

        # call send_service_message on bts_periphery
        send_service_message_args_type = sp.TRecord(_from=sp.TAddress, to=sp.TString,
                                                    coin_names=sp.TMap(sp.TNat, sp.TString),
                                                    values=sp.TMap(sp.TNat, sp.TNat),
                                                    fees=sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type,
                                                       self.data.bts_periphery_address,
                                                       "send_service_message").open_some()
        send_service_message_args = sp.record(_from=sp.self_address, to=fa, coin_names=self.data.charged_coins,
                                              values=self.data.charged_amounts,
                                              fees=sp.map({}))
        sp.transfer(send_service_message_args, sp.tez(0), send_service_message_entry_point)

        sp.for i in sp.range(0, sp.len(self.data.charged_coins)):
            del self.data.charged_coins[i]

        sp.for i in sp.range(0, sp.len(self.data.charged_amounts)):
            del self.data.charged_amounts[i]

    def _lock_balance(self, to, coin_name, value):
          self.data.balances[sp.record(address=to, coin_name=coin_name)].locked_balance = self.data.balances[sp.record(address=to, coin_name=coin_name)].locked_balance + value

    @sp.entry_point    
    def update_coin_db(self):
        self.only_owner()
        self.data.coins[self.data.native_coin_name] = self.NATIVE_COIN_ADDRESS
        self.data.coins_address[self.NATIVE_COIN_ADDRESS] = self.data.native_coin_name
        self.data.coin_details[self.data.native_coin_name].addr = self.NATIVE_COIN_ADDRESS

        sp.for item in self.data.coins_name:
             sp.if item != self.data.native_coin_name:
                  self.data.coins_address[self.data.coin_details[item].addr] = item

    @sp.entry_point
    def set_bts_owner_manager(self, owner_manager):
        sp.verify(self.data.owners[sp.sender] == True , message= "Unauthorized")
        sp.verify(owner_manager != sp.address("tz10000"), message= "InvalidAddress")
        self.data.bts_owner_manager = owner_manager

@sp.add_test(name="BTSCore")
def test():
    c1 = BTSCore(
        owner_manager=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        bts_periphery_addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
        _native_coin_name="NativeCoin",
        _fee_numerator=sp.nat(1000),
        _fixed_fee=sp.nat(10)
    )
    scenario = sp.test_scenario()
    scenario.h1("BTSCore")
    scenario += c1


sp.add_compilation_target("bts_core", BTSCore(
    owner_manager=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
    bts_periphery_addr=sp.address("tz1VA29GwaSA814BVM7AzeqVzxztEjjxiMEc"),
    _native_coin_name="NativeCoin",
    _fee_numerator=sp.nat(1000),
    _fixed_fee=sp.nat(10)
))






