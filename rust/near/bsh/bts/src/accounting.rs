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

        let initial_storage_usage = env::storage_usage();

        let coin_id = *self.registered_coins.get(&coin_account).unwrap();
        let mut balance = match self.balances.get(&sender_id, &coin_id) {
            Some(balance) => balance,
            None => AccountBalance::default(),
        };
        self.process_deposit(amount, &mut balance);

        self.balances.set(&sender_id, &coin_id, balance);

        // calculate storage cost for the account
        let total_storage_cost = self.calculate_storage_cost(initial_storage_usage);

        self.storage_balances
            .set(&sender_id, &coin_id, total_storage_cost.into());

        PromiseOrValue::Value(U128::from(0))
    }

    #[payable]
    pub fn deposit(&mut self) {
        let account = env::predecessor_account_id();
        let amount = env::attached_deposit();
        self.assert_have_minimum_amount(amount);
        let coin_id = Self::hash_coin_id(&self.native_coin_name);

        let initial_storage_usage = env::storage_usage();
        let mut balance = match self.balances.get(&account, &coin_id) {
            Some(balance) => balance,
            None => AccountBalance::default(),
        };
        self.process_deposit(amount, &mut balance);
        self.balances.set(&account, &coin_id, balance);

        let total_storage_cost = self.calculate_storage_cost(initial_storage_usage);
        self.storage_balances
            .set(&account, &coin_id, total_storage_cost.into());
    }

    #[payable]
    pub fn withdraw(&mut self, coin_name: String, amount: U128) {
        let amount: u128 = amount.into();
        let account = env::predecessor_account_id();
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        self.assert_have_minimum_amount(amount);
        self.assert_have_sufficient_deposit(&account, &coin_id, amount, None);

        self.assert_minimum_one_yocto();
        // Check for attached storage usage cost
        self.assert_have_sufficient_storage_deposit(&account, &coin_id);

        // Check if current account have sufficient balance
        self.assert_have_sufficient_balance(amount);

        let coin = self.coins.get(&coin_id).unwrap();

        let transfer_promise = if coin.network() != &self.network {
            ext_nep141::ext(coin.metadata().uri().to_owned().unwrap())
                .with_attached_deposit(1)
                .ft_transfer(account.clone(), amount.into(), None)
        } else if let Some(uri) = coin.metadata().uri_deref() {
            ext_ft::ext(uri).with_attached_deposit(1).ft_transfer(
                account.clone(),
                U128::from(amount),
                None,
            )
        } else {
            Promise::new(account.clone()).transfer(amount)
        };

        transfer_promise.then(
            Self::ext(env::current_account_id()).on_withdraw(account, amount, coin_name, coin_id),
        );
    }

    pub fn reclaim(&mut self, coin_name: String, amount: U128) {
        let amount: u128 = amount.into();
        let account = env::predecessor_account_id();
        self.assert_have_minimum_amount(amount);
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        self.assert_have_sufficient_refundable(&account, &coin_id, amount);

        let mut balance = self.balances.get(&account, &coin_id).unwrap();
        balance.refundable_mut().sub(amount).unwrap();
        balance.deposit_mut().add(amount).unwrap();

        self.balances.set(&account, &coin_id, balance);
    }

    pub fn locked_balance_of(&self, account_id: AccountId, coin_name: String) -> U128 {
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let balance = self
            .balances
            .get(&account_id, &coin_id)
            .unwrap_or_else(|| env::panic_str(format!("{}", BshError::AccountNotExist).as_str()));
        balance.locked().into()
    }

    pub fn refundable_balance_of(&self, account_id: AccountId, coin_name: String) -> U128 {
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let balance = self
            .balances
            .get(&account_id, &coin_id)
            .unwrap_or_else(|| env::panic_str(format!("{}", BshError::AccountNotExist).as_str()));
        balance.refundable().into()
    }

    #[cfg(feature = "testable")]
    pub fn account_balance(
        &self,
        owner_id: AccountId,
        coin_name: String,
    ) -> Option<AccountBalance> {
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        self.balances.get(&owner_id, &coin_id)
    }

    pub fn balance_of(&self, account_id: AccountId, coin_name: String) -> U128 {
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let balance = self
            .balances
            .get(&account_id, &coin_id)
            .unwrap_or_else(|| env::panic_str(format!("{}", BshError::AccountNotExist).as_str()));
        balance.deposit().into()
    }

    #[private]
    pub fn on_withdraw(
        &mut self,
        account: AccountId,
        amount: u128,
        coin_name: String,
        coin_id: CoinId,
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
                    "token_name": coin_name
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());

                self.storage_balances.set(&account, &coin_id, 0)
            }
            PromiseResult::NotReady => {
                log!("Not Ready")
            }
            PromiseResult::Failed => {
                let log = json!(
                {
                    "event": "Withdraw",
                    "code": "1",
                    "by": account,
                    "amount": amount.to_string(),
                    "token_name": coin_name
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    pub fn get_storage_balance(&self, account: AccountId, coin_name: String) -> U128 {
        let coin_id = self.coin_id(&coin_name).unwrap();
        match self.storage_balances.get(&account, &coin_id) {
            Some(storage_cost) => U128::from(storage_cost),
            None => U128::from(0),
        }
    }
}
