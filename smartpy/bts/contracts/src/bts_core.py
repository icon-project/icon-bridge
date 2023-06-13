import smartpy as sp

types = sp.io.import_script_from_url("file:./contracts/src/Types.py")
FA2_contract = sp.io.import_script_from_url("file:./contracts/src/FA2_contract.py")

Coin = sp.TRecord(addr=sp.TAddress,
                  fee_numerator=sp.TNat,
                  fixed_fee=sp.TNat,
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
    ZERO_ADDRESS = sp.address("tz1ZZZZZZZZZZZZZZZZZZZZZZZZZZZZNkiRg")
    # Nat:(TWO.pow256 - 1)
    UINT_CAP = sp.nat(115792089237316195423570985008687907853269984665640564039457584007913129639935)

    def __init__(self, _native_coin_name, _fee_numerator, _fixed_fee, owner_manager):
        self.update_initial_storage(
            bts_owner_manager=owner_manager,
            bts_periphery_address=sp.none,
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
            transfer_status=sp.none
        )

    def only_owner(self):
        is_owner = sp.view("is_owner", self.data.bts_owner_manager, sp.sender, t=sp.TBool).open_some(
            "OwnerNotFound")
        sp.verify(is_owner == True, message="Unauthorized")

    def only_bts_periphery(self):
        sp.verify(sp.sender == self.data.bts_periphery_address.open_some("Address not set"), "Unauthorized")

    @sp.entry_point
    def update_bts_periphery(self, bts_periphery):
        """
        update BTS Periphery address.
        :param bts_periphery:  BTSPeriphery contract address.
        :return:
        """
        sp.set_type(bts_periphery, sp.TAddress)

        self.only_owner()
        sp.if self.data.bts_periphery_address.is_some():
            has_requests = sp.view("has_pending_request", self.data.bts_periphery_address.open_some("Address not set"), sp.unit, t=sp.TBool).open_some("OwnerNotFound")
            sp.verify(has_requests == False, "HasPendingRequest")
        self.data.bts_periphery_address = sp.some(bts_periphery)

    @sp.entry_point
    def set_fee_ratio(self, name, fee_numerator, fixed_fee):
        """
        set fee ratio. The transfer fee is calculated by fee_numerator/FEE_DEMONINATOR.
        fee_numerator if it is set to `10`, which means the default fee ratio is 0.1%.
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
        sp.verify((name == self.data.native_coin_name) | self.data.coins.contains(name), message = "TokenNotExists")
        sp.verify((fixed_fee > sp.nat(0)) & (fee_numerator >= sp.nat(0)), message = "LessThan0")
        self.data.coin_details[name].fee_numerator = fee_numerator
        self.data.coin_details[name].fixed_fee = fixed_fee

    @sp.entry_point(lazify=False)
    def update_register(self, ep):
        self.only_owner()
        sp.set_entry_point("register", ep)

    @sp.entry_point(lazify=True)
    def register(self, name, fee_numerator, fixed_fee, addr, token_metadata, metadata):
        """
        Registers a wrapped coin and id number of a supporting coin.
        :param name: Must be different with the native coin name.
        :param fee_numerator:
        :param fixed_fee:
        :param addr: address of the coin
        :param token_metadata: Token metadata name, symbol and decimals of wrapped token
        :param metadata: metadata of the token contract
        :return:
        """
        sp.set_type(name, sp.TString)
        sp.set_type(fee_numerator, sp.TNat)
        sp.set_type(fixed_fee, sp.TNat)
        sp.set_type(addr, sp.TAddress)
        sp.set_type(token_metadata, sp.TMap(sp.TString, sp.TBytes))
        sp.set_type(metadata, sp.TBigMap(sp.TString, sp.TBytes))

        self.only_owner()
        sp.verify(name != self.data.native_coin_name, message="ExistNativeCoin")
        sp.verify(self.data.coins.contains(name) == False, message= "ExistCoin")
        sp.verify(self.data.coins_address.contains(addr) == False, message="AddressExists")
        sp.verify(fee_numerator <= self.FEE_DENOMINATOR, message="InvalidSetting")
        sp.verify((fixed_fee >= sp.nat(0)) & (fee_numerator >= sp.nat(0)), message="LessThan0")
        with sp.if_(addr == self.ZERO_ADDRESS):
            deployed_fa2 = sp.create_contract_operation(contract=FA2_contract.SingleAssetToken(admin=sp.self_address, metadata=metadata,
                                                              token_metadata=token_metadata
                                                              ))
            sp.operations().push(deployed_fa2.operation)
            self.data.coins[name] = deployed_fa2.address
            self.data.coins_name.push(name)
            self.data.coins_address[deployed_fa2.address] = name
            self.data.coin_details[name] = sp.record(
                addr = deployed_fa2.address,
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

        token_map = sp.map({0:name}, tkey=sp.TNat, tvalue=sp.TString)
        val_map = sp.map({0:self.UINT_CAP}, tkey=sp.TNat, tvalue=sp.TNat)

        # call set_token_limit on bts_periphery
        set_token_limit_args_type = sp.TRecord(coin_names=sp.TMap(sp.TNat, sp.TString), token_limit=sp.TMap(sp.TNat, sp.TNat))
        set_token_limit_entry_point = sp.contract(set_token_limit_args_type, self.data.bts_periphery_address.open_some("Address not set"),
                                                  "set_token_limit").open_some("ErrorINCALL")
        set_token_limit_args = sp.record(coin_names=token_map, token_limit=val_map)
        sp.transfer(set_token_limit_args, sp.tez(0), set_token_limit_entry_point)

    @sp.onchain_view()
    def coin_id(self, coin_name):
        """
        Return address of Coin whose name is the same with given coin_ame.
        :param coin_name:
        :return: An address of coin_name.
        """
        sp.set_type(coin_name, sp.TString)

        sp.result(self.data.coins.get(coin_name))

    @sp.onchain_view()
    def is_valid_coin(self, coin_name):
        """
        Check Validity of a coin_name
        :param coin_name:
        :return: true or false
        """
        sp.set_type(coin_name, sp.TString)

        sp.result((self.data.coins.contains(coin_name))|( coin_name == self.data.native_coin_name))

    @sp.onchain_view()
    def fee_ratio(self, coin_name):
        """
        Get fee numerator and fixed fee
        :param coin_name: Coin name
        :return: a record (Fee numerator for given coin, Fixed fee for given coin)
        """
        sp.set_type(coin_name, sp.TString)

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

        # sending user_balance as 0 because in smartpy we cannot get Tez balance of address
        locked_balance = sp.local("locked_balance", sp.nat(0))
        refundable_balance = sp.local("refundable_balance", sp.nat(0))

        sp.if self.data.balances.contains(sp.record(address=params.owner, coin_name=params.coin_name)):
            locked_balance.value = self.data.balances[
                sp.record(address=params.owner, coin_name=params.coin_name)].locked_balance
            refundable_balance.value = self.data.balances[
                sp.record(address=params.owner, coin_name=params.coin_name)].refundable_balance

        with sp.if_(params.coin_name == self.data.native_coin_name):
            sp.result(sp.record(usable_balance=sp.nat(0),
                                locked_balance=locked_balance.value,
                                refundable_balance=refundable_balance.value,
                                user_balance=sp.nat(0)))
        with sp.else_():
            fa2_address = self.data.coins.get(params.coin_name)
            user_balance= sp.view("balance_of", fa2_address, sp.record(owner=params.owner, token_id=sp.nat(0)), t=sp.TNat).open_some("Invalid view")

            allowance = sp.view("get_allowance", fa2_address, sp.record(spender=sp.self_address, owner=params.owner), t=sp.TNat).open_some("Invalid view")
            usable_balance = sp.local("usable_balance", allowance, t=sp.TNat)
            sp.if allowance > user_balance:
                usable_balance.value = user_balance

            sp.result(sp.record(usable_balance=usable_balance.value,
                                locked_balance=locked_balance.value,
                                refundable_balance=refundable_balance.value,
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

        i = sp.local("i", sp.nat(0))
        sp.for item in params.coin_names:
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
        i = sp.local("i", sp.nat(0))
        sp.for item in self.data.coins_name:
            accumulated_fees.value[i.value] = sp.record(coin_name=item, value=self.data.aggregation_fee.get(item, default_value=sp.nat(0)))
            i.value += sp.nat(1)
        sp.result(accumulated_fees.value)

    @sp.entry_point(lazify=False)
    def update_transfer_native_coin(self, ep):
        self.only_owner()
        sp.set_entry_point("transfer_native_coin", ep)

    @sp.entry_point(check_no_incoming_transfer=False, lazify=True)
    def transfer_native_coin(self, to):
        """
        Allow users to deposit `sp.amount` native coin into a BTSCore contract.
        :param to: An address that a user expects to receive an amount of tokens.
        :return: 
        """
        sp.set_type(to, sp.TString)

        amount_in_nat = sp.local("amount_in_nat", sp.utils.mutez_to_nat(sp.amount), t=sp.TNat)
        # call check_transfer_restrictions on bts_periphery
        check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address.open_some("Address not set"),
                                 sp.record(coin_name=self.data.native_coin_name, user=sp.sender, value=amount_in_nat.value),
                                 t=sp.TBool).open_some()
        sp.verify(check_transfer == True, "FailCheckTransfer")

        charge_amt = amount_in_nat.value * self.data.coin_details[self.data.native_coin_name].fee_numerator / self.FEE_DENOMINATOR + self.data.coin_details[self.data.native_coin_name].fixed_fee

        self._send_service_message(sp.sender, to, self.data.native_coin_name, amount_in_nat.value, charge_amt)

    @sp.entry_point(lazify=False)
    def update_transfer(self, ep):
        self.only_owner()
        sp.set_entry_point("transfer", ep)

    @sp.entry_point(lazify=True)
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
        sp.verify(self.data.coins.contains(coin_name), message= "CoinNotRegistered")
        fa2_address = self.data.coins[coin_name]

        # call check_transfer_restrictions on bts_periphery
        check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address.open_some("Address not set"),
                                 sp.record(coin_name=coin_name, user=sp.sender, value=value),
                                 t=sp.TBool).open_some()
        sp.verify(check_transfer == True, "FailCheckTransfer")

        charge_amt = value * self.data.coin_details[coin_name].fee_numerator / self.FEE_DENOMINATOR + self.data.coin_details[coin_name].fixed_fee
        
        # call transfer from in FA2
        transfer_args_type = sp.TList(sp.TRecord(from_=sp.TAddress, txs=sp.TList(sp.TRecord(
            to_=sp.TAddress, token_id=sp.TNat, amount=sp.TNat).layout(("to_", ("token_id", "amount"))))
                                                 ).layout(("from_", "txs")))
        transfer_entry_point = sp.contract(transfer_args_type, fa2_address, "transfer").open_some()
        transfer_args = [sp.record(from_=sp.sender, txs=[sp.record(to_=sp.self_address, token_id=sp.nat(0), amount=value)])]
        sp.transfer(transfer_args, sp.tez(0), transfer_entry_point)

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

        coins = sp.local("coins", {sp.nat(0) : coin_name}, t=sp.TMap(sp.TNat, sp.TString))
        amounts = sp.local("amounts", {sp.nat(0) : sp.as_nat(value - charge_amt)}, t=sp.TMap(sp.TNat, sp.TNat))
        fees = sp.local("fees", {sp.nat(0): charge_amt}, t=sp.TMap(sp.TNat, sp.TNat))

        # call send_service_message on bts_periphery
        send_service_message_args_type = sp.TRecord(_from = sp.TAddress, to = sp.TString, coin_names = sp.TMap(sp.TNat, sp.TString), values = sp.TMap(sp.TNat, sp.TNat), fees = sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type, self.data.bts_periphery_address.open_some("Address not set"), "send_service_message").open_some()
        send_service_message_args = sp.record(_from = _from, to = to, coin_names = coins.value, values = amounts.value, fees = fees.value)
        sp.transfer(send_service_message_args, sp.tez(0), send_service_message_entry_point)

    @sp.entry_point(lazify=False)
    def update_transfer_batch(self, ep):
        self.only_owner()
        sp.set_entry_point("transfer_batch", ep)

    @sp.entry_point(check_no_incoming_transfer=False, lazify=True)
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
        size = sp.local("size", sp.nat(0), t=sp.TNat)
        with sp.if_(amount_in_nat.value != sp.nat(0)):
            size.value = sp.len(coin_names) + sp.nat(1)
        with sp.else_():
            size.value = sp.len(coin_names)
        sp.verify(size.value <= self.MAX_BATCH_SIZE, message ="Batch maxSize Exceeds")
    
        coins = sp.local("_coins", {}, t=sp.TMap(sp.TNat, sp.TString))
        amounts = sp.local("_amounts", {}, t=sp.TMap(sp.TNat, sp.TNat))
        charge_amts = sp.local("_charge_amts", {}, t=sp.TMap(sp.TNat, sp.TNat))

        coin_name = sp.local("coin_name", "", t= sp.TString)
        value = sp.local("value", sp.nat(0), t= sp.TNat)
        
        sp.for i in sp.range(sp.nat(0), sp.len(coin_names)):
            sp.verify(coin_names[i] != self.data.native_coin_name, message="InvalidCoin")
            sp.verify(self.data.coins.contains(coin_names[i]), message= "CoinNotRegistered")
            fa2_address = self.data.coins[coin_names[i]]

            coin_name.value = coin_names[i]
            value.value = values[i]
            sp.verify(value.value > sp.nat(0), message ="ZeroOrLess")

            # call check_transfer_restrictions on bts_periphery
            check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address.open_some("Address not set"),
                                     sp.record(coin_name=coin_name.value, user=sp.sender, value=value.value),
                                     t=sp.TBool).open_some()
            sp.verify(check_transfer == True, "FailCheckTransfer")

            # call transfer from in FA2
            transfer_args_type = sp.TList(sp.TRecord(from_=sp.TAddress, txs=sp.TList(sp.TRecord(
                to_=sp.TAddress, token_id=sp.TNat, amount=sp.TNat).layout(("to_", ("token_id", "amount"))))
                                                     ).layout(("from_", "txs")))
            transfer_entry_point = sp.contract(transfer_args_type, fa2_address, "transfer").open_some()
            transfer_args = [
                sp.record(from_=sp.sender, txs=[sp.record(to_=sp.self_address, token_id=sp.nat(0), amount=value.value)])]
            sp.transfer(transfer_args, sp.tez(0), transfer_entry_point)

            coin = sp.local("coin", self.data.coin_details[coin_name.value], t=Coin)
            coins.value[i] = coin_name.value
            charge_amts.value[i] = value.value * coin.value.fee_numerator / self.FEE_DENOMINATOR + coin.value.fixed_fee
            amounts.value[i] = sp.as_nat(value.value - charge_amts.value[i])

            self._lock_balance(sp.sender, coin_name.value, value.value)

        sp.if amount_in_nat.value !=sp.nat(0):
            # call check_transfer_restrictions on bts_periphery
            check_transfer = sp.view("check_transfer_restrictions", self.data.bts_periphery_address.open_some("Address not set"),
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
                                                       self.data.bts_periphery_address.open_some("Address not set"),
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

        with sp.if_(self.data.balances.contains(sp.record(address=sp.sender,coin_name=coin_name))):
            sp.verify(self.data.balances[sp.record(address=sp.sender,coin_name=coin_name)].refundable_balance >= value, message="Imbalance")
            self.data.balances[sp.record(address=sp.sender,coin_name=coin_name)].refundable_balance = sp.as_nat(self.data.balances[sp.record(address=sp.sender,coin_name=coin_name)].refundable_balance - value)
            self.refund(sp.sender, coin_name, value)
        with sp.else_():
            sp.failwith("NoRefundableBalance")

    def refund(self, to, coin_name, value):
        """
        :param to:
        :param coin_name:
        :param value:
        :return:
        """
        sp.set_type(to, sp.TAddress)
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)

        sp.verify(sp.sender == sp.self_address, message="Unauthorized")

        with sp.if_(coin_name == self.data.native_coin_name):
            self.payment_transfer(to, value)
        with sp.else_():
            # call transfer in FA2
            transfer_args_type = sp.TList(sp.TRecord(from_=sp.TAddress, txs=sp.TList(sp.TRecord(
                to_=sp.TAddress, token_id=sp.TNat, amount=sp.TNat).layout(("to_", ("token_id", "amount"))))
                                                     ).layout(("from_", "txs")))
            transfer_entry_point = sp.contract(transfer_args_type, self.data.coins[coin_name], "transfer").open_some()
            transfer_args = [sp.record(from_=sp.sender, txs=[sp.record(to_=to, token_id=sp.nat(0), amount=value)])]
            sp.transfer(transfer_args, sp.tez(0), transfer_entry_point)

    def payment_transfer(self, to, amount):
        sp.set_type(to, sp.TAddress)
        sp.set_type(amount, sp.TNat)

        sp.send(to, sp.utils.nat_to_mutez(amount), message="PaymentFailed")

    @sp.entry_point(lazify=False)
    def update_mint(self, ep):
        self.only_owner()
        sp.set_entry_point("mint", ep)

    @sp.entry_point(lazify=True)
    def mint(self, to, coin_name, value, callback):
        """
        mint the wrapped coin.
        :param to: the account receive the minted coin
        :param coin_name: coin name
        :param value: the minted amount
        :param callback: callback function type in bts_periphery
        :return:
        """
        sp.set_type(to, sp.TAddress)
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)
        sp.set_type(callback, sp.TContract(sp.TOption(sp.TString)))

        self.only_bts_periphery()
        with sp.if_(coin_name == self.data.native_coin_name):
            self.payment_transfer(to, value)
        with sp.else_():
            with sp.if_(self.data.coin_details[coin_name].coin_type == self.NATIVE_WRAPPED_COIN_TYPE):
                # call mint in FA2
                mint_args_type = sp.TList(sp.TRecord(to_=sp.TAddress, amount=sp.TNat).layout(("to_", "amount")))
                mint_entry_point = sp.contract(mint_args_type, self.data.coins[coin_name], "mint").open_some()
                mint_args = [sp.record(to_=to, amount=value)]
                sp.transfer(mint_args, sp.tez(0), mint_entry_point)
            with sp.else_():
                sp.if self.data.coin_details[coin_name].coin_type == self.NON_NATIVE_TOKEN_TYPE:
                    # call transfer in FA2
                    transfer_args_type = sp.TList(sp.TRecord(from_=sp.TAddress, txs=sp.TList(sp.TRecord(
                        to_=sp.TAddress, token_id=sp.TNat, amount=sp.TNat).layout(("to_", ("token_id", "amount"))))
                                                             ).layout(("from_", "txs")))
                    transfer_entry_point = sp.contract(transfer_args_type, self.data.coins[coin_name], "transfer").open_some()
                    transfer_args = [sp.record(from_=sp.self_address, txs=[sp.record(to_=to, token_id=sp.nat(0), amount=value)])]
                    sp.transfer(transfer_args, sp.tez(0), transfer_entry_point)
        sp.transfer(sp.some("success"), sp.tez(0), callback)

    @sp.entry_point
    def transfer_callback(self, string, requester, coin_name, value):
        sp.set_type(string, sp.TOption(sp.TString))
        sp.set_type(requester, sp.TAddress)
        sp.set_type(coin_name, sp.TString)
        sp.set_type(value, sp.TNat)

        sp.verify(sp.sender == self.data.coins[coin_name], "Unauthorized")
        self.data.transfer_status = string

        with sp.if_(self.data.transfer_status.open_some() == "success"):
            pass
        with sp.else_():
            self.data.balances[sp.record(address=requester, coin_name=coin_name)].refundable_balance = \
                self.data.balances.get(sp.record(address=requester, coin_name=coin_name),
                                       default_value=sp.record(locked_balance=sp.nat(0),refundable_balance=sp.nat(0))
                                       ).refundable_balance + value
        self.data.transfer_status = sp.none

    @sp.entry_point(lazify=False)
    def update_handle_response_service(self, ep):
        self.only_owner()
        sp.set_entry_point("handle_response_service", ep)

    @sp.entry_point(lazify=True)
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
        return_flag = sp.local("return_flag", False, t=sp.TBool)
        sp.if requester == sp.self_address:
            sp.if rsp_code == self.RC_ERR:
                self.data.aggregation_fee[coin_name] = self.data.aggregation_fee.get(coin_name,
                                                                                     default_value=sp.nat(0)) + value
            return_flag.value = True

        sp.if return_flag.value == False:
            amount = sp.local("amount", value + fee, t=sp.TNat)
            sp.if self.data.balances.contains(sp.record(address=requester, coin_name=coin_name)):
                self.data.balances[sp.record(address=requester, coin_name=coin_name)].locked_balance = \
                    sp.as_nat(self.data.balances.get(sp.record(address=requester, coin_name=coin_name),
                                                     default_value=sp.record(locked_balance=sp.nat(0), refundable_balance=sp.nat(0))).locked_balance - amount.value)

            sp.if rsp_code == self.RC_ERR:
                with sp.if_(coin_name == self.data.native_coin_name):
                    with sp.if_(sp.utils.mutez_to_nat(sp.balance) >= value):
                        self.payment_transfer(requester, value)
                    with sp.else_():
                        self.data.balances[sp.record(address=requester, coin_name=coin_name)].refundable_balance = self.data.balances.get(
                            sp.record(address=requester, coin_name=coin_name),
                            default_value=sp.record(locked_balance=sp.nat(0), refundable_balance=sp.nat(0))).refundable_balance + value
                with sp.else_():
                    # call transfer in FA2
                    transfer_args_type = sp.TList(sp.TRecord(
                        callback=sp.TContract(sp.TRecord(string=sp.TOption(sp.TString), requester=sp.TAddress, coin_name=sp.TString, value=sp.TNat)),
                        from_=sp.TAddress,
                        coin_name=sp.TString,
                        txs=sp.TList(sp.TRecord(to_=sp.TAddress, token_id=sp.TNat, amount=sp.TNat).layout(("to_", ("token_id", "amount"))))
                                                             ).layout((("from_", "coin_name"), ("callback", "txs"))))
                    transfer_entry_point = sp.contract(transfer_args_type, self.data.coins[coin_name],
                                                       "transfer_bts").open_some()
                    t = sp.TRecord(string=sp.TOption(sp.TString), requester=sp.TAddress, coin_name=sp.TString, value=sp.TNat)
                    callback = sp.contract(t, sp.self_address, "transfer_callback")
                    transfer_args = [
                        sp.record(callback=callback.open_some(), from_=sp.self_address, coin_name=coin_name, txs=[sp.record(to_=requester, token_id=sp.nat(0), amount=value)])]
                    sp.transfer(transfer_args, sp.tez(0), transfer_entry_point)

            sp.if rsp_code == self.RC_OK:
                fa2_address = self.data.coins[coin_name]
                sp.if (coin_name != self.data.native_coin_name) & (self.data.coin_details[coin_name].coin_type == self.NATIVE_WRAPPED_COIN_TYPE):
                    # call burn in FA2
                    burn_args_type = sp.TList(sp.TRecord(from_=sp.TAddress, token_id=sp.TNat, amount=sp.TNat).layout(("from_", ("token_id", "amount"))))
                    burn_entry_point = sp.contract(burn_args_type, fa2_address, "burn").open_some()
                    burn_args = [sp.record(from_=sp.self_address, token_id=sp.nat(0), amount=value)]
                    sp.transfer(burn_args, sp.tez(0), burn_entry_point)

            self.data.aggregation_fee[coin_name] = self.data.aggregation_fee.get(coin_name, default_value=sp.nat(0)) + fee

    @sp.entry_point(lazify=False)
    def update_transfer_fees(self, ep):
        self.only_owner()
        sp.set_entry_point("transfer_fees", ep)

    @sp.entry_point(lazify=True)
    def transfer_fees(self, fa):
        """
        Handle a request of Fee Gathering. Usage: Copy all charged fees to an array
        :param fa:
        :return:
        """
        sp.set_type(fa, sp.TString)

        self.only_bts_periphery()
        l = sp.local("l", sp.nat(0))
        sp.for item in self.data.coins_name:
            sp.if self.data.aggregation_fee.get(item, default_value=sp.nat(0)) != sp.nat(0):
                self.data.charged_coins[l.value] = item
                self.data.charged_amounts[l.value] = self.data.aggregation_fee.get(item, default_value=sp.nat(0))
                del self.data.aggregation_fee[item]
                l.value += sp.nat(1)

        _charged_coins = sp.local("_charged_coins", self.data.charged_coins)
        _charged_amounts = sp.local("_charged_amounts", self.data.charged_amounts)

        # call send_service_message on bts_periphery
        send_service_message_args_type = sp.TRecord(_from=sp.TAddress, to=sp.TString,
                                                    coin_names=sp.TMap(sp.TNat, sp.TString),
                                                    values=sp.TMap(sp.TNat, sp.TNat),
                                                    fees=sp.TMap(sp.TNat, sp.TNat))
        send_service_message_entry_point = sp.contract(send_service_message_args_type,
                                                       self.data.bts_periphery_address.open_some("Address not set"),
                                                       "send_service_message").open_some()
        send_service_message_args = sp.record(_from=sp.self_address, to=fa, coin_names=_charged_coins.value,
                                              values=_charged_amounts.value,
                                              fees=sp.map({}))
        sp.transfer(send_service_message_args, sp.tez(0), send_service_message_entry_point)

        sp.for i in sp.range(0, sp.len(self.data.charged_coins)):
            del self.data.charged_coins[i]

        sp.for i in sp.range(0, sp.len(self.data.charged_amounts)):
            del self.data.charged_amounts[i]

    def _lock_balance(self, to, coin_name, value):
        new_balance = self.data.balances.get(sp.record(address=to, coin_name=coin_name),
                                             default_value=sp.record(locked_balance=sp.nat(0),
                                                                     refundable_balance=sp.nat(0)))
        self.data.balances[sp.record(address=to, coin_name=coin_name)] = sp.record(
            locked_balance=new_balance.locked_balance + value, refundable_balance=new_balance.refundable_balance)

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
        sp.set_type(owner_manager, sp.TAddress)

        sp.verify(self.data.owners.get(sp.sender) == True , message= "Unauthorized")
        self.data.bts_owner_manager = owner_manager


sp.add_compilation_target("bts_core", BTSCore(
    owner_manager=sp.address("KT1MxuVecS7HRRRZrJM7juddJg1HZZ4SGA5B"),
    _native_coin_name="btp-NetXnHfVqm9iesp.tezos-XTZ",
    _fee_numerator=sp.nat(100),
    _fixed_fee=sp.nat(450)
))






