use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * Transactions  * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn ft_on_transfer(
        &mut self,
        sender_id: AccountId,
        amount: U128,
        #[allow(unused_variables)] msg: String,
    ) -> PromiseOrValue<U128> {
        let amount = amount.into();
        let coin_account = env::predecessor_account_id();

        self.assert_have_minimum_amount(amount);
        self.assert_coin_registered(&coin_account);

        let coin_id = self.registered_coins.get(&coin_account).unwrap().clone();
        let mut balance = match self.balances.get(&sender_id, &coin_id) {
            Some(balance) => balance,
            None => AccountBalance::default(),
        };

        self.process_deposit(amount, &mut balance);
        self.balances.set(&sender_id, &coin_id, balance);

        PromiseOrValue::Value(U128::from(0))
    }

    #[payable]
    pub fn deposit(&mut self) {
        let account = env::predecessor_account_id();
        let amount = env::attached_deposit();
        self.assert_have_minimum_amount(amount);
        let coin_id = Self::hash_coin_id(&self.native_coin_name);

        let mut balance = match self.balances.get(&account, &coin_id) {
            Some(balance) => balance,
            None => AccountBalance::default(),
        };

        self.process_deposit(amount, &mut balance);
        self.balances.set(&account, &coin_id, balance);
    }

    #[payable]
    pub fn withdraw(&mut self, coin_id: CoinId, amount: U128) {
        // To Prevent Spam
        assert_one_yocto();

        let amount: u128 = amount.into();
        let account = env::predecessor_account_id();

        self.assert_have_minimum_amount(amount);
        self.assert_have_sufficient_deposit(&account, &coin_id, amount, None);

        // Check if current account have sufficient balance
        self.assert_have_sufficient_balance(1 + amount);

        let coin = self.coins.get(&coin_id).unwrap();

        let transfer_promise = if coin.network() != &self.network {
            ext_nep141::ft_transfer_with_storage_check(
                account.clone(),
                amount,
                None,
                coin.metadata().uri().to_owned().unwrap(),
                estimate::NO_DEPOSIT,
                estimate::GAS_FOR_FT_TRANSFER_CALL,
            )
        } else {
            if let Some(uri) = coin.metadata().uri_deref() {
                ext_ft::ft_transfer(
                    account.clone(),
                    U128::from(amount),
                    None,
                    uri,
                    estimate::ONE_YOCTO,
                    estimate::GAS_FOR_FT_TRANSFER_CALL,
                )
            } else {
                Promise::new(account.clone()).transfer(amount)
            }
        };

        transfer_promise
            .then(ext_self::on_withdraw(
                account.clone(),
                amount,
                coin_id,
                coin.symbol().to_owned(),
                env::current_account_id(),
                estimate::NO_DEPOSIT,
                estimate::GAS_FOR_FT_TRANSFER_CALL,
            ))
            .then(Promise::new(account.clone()).transfer(1));
    }

    pub fn reclaim(&mut self, coin_id: CoinId, amount: U128) {
        let amount: u128 = amount.into();
        let account = env::predecessor_account_id();
        self.assert_have_minimum_amount(amount.into());
        self.assert_coins_exists(&vec![coin_id.clone()]);
        self.assert_have_sufficient_refundable(&account, &coin_id, amount);

        let mut balance = self.balances.get(&account, &coin_id).unwrap();
        balance.refundable_mut().sub(amount).unwrap();
        balance.deposit_mut().add(amount).unwrap();

        self.balances.set(&account, &coin_id, balance);
    }

    pub fn locked_balance_of(&self, owner_id: AccountId, coin_id: CoinId) -> U128 {
        self.assert_coins_exists(&vec![coin_id.clone()]);
        let balance = self
            .balances
            .get(&owner_id, &coin_id)
            .expect(format!("{}", BshError::AccountNotExist).as_str());
        balance.locked().into()
    }

    pub fn refundable_balance_of(&self, owner_id: AccountId, coin_id: CoinId) -> U128 {
        self.assert_coins_exists(&vec![coin_id.clone()]);
        let balance = self
            .balances
            .get(&owner_id, &coin_id)
            .expect(format!("{}", BshError::AccountNotExist).as_str());
        balance.refundable().into()
    }

    #[cfg(feature = "testable")]
    pub fn account_balance(&self, owner_id: AccountId, coin_id: CoinId) -> Option<AccountBalance> {
        self.balances.get(&owner_id, &coin_id)
    }

    pub fn balance_of(&self, owner_id: AccountId, coin_id: CoinId) -> U128 {
        self.assert_coins_exists(&vec![coin_id.clone()]);
        let balance = self
            .balances
            .get(&owner_id, &coin_id)
            .expect(format!("{}", BshError::AccountNotExist).as_str());
        balance.deposit().into()
    }

    #[private]
    pub fn on_withdraw(
        &mut self,
        account: AccountId,
        amount: u128,
        coin_id: CoinId,
        coin_symbol: String,
    ) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                let mut balance = self.balances.get(&account, &coin_id).unwrap();
                balance.deposit_mut().sub(amount).unwrap();
                self.balances.set(&account.clone(), &coin_id, balance);
                let log = json!(
                {

                    "event": "Withdraw",
                    "code": "0",
                    "by": account,
                    "amount": amount.to_string(),
                    "token_name": coin_symbol

                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
            PromiseResult::NotReady => {
                log!("Not Ready")
            }
            PromiseResult::Failed => {
                println!("{}", amount.to_string());
                let log = json!(
                {

                    "event": "Withdraw",
                    "code": "1",
                    "by": account,
                    "amount": amount.to_string(),
                    "token_name": coin_symbol

                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }
}
