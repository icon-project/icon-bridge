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
        let token_account = env::predecessor_account_id();

        self.assert_have_minimum_amount(amount);
        self.assert_token_registered(&token_account);

        let initial_storage_usage = env::storage_usage();

        let token_id = *self.registered_tokens.get(&token_account).unwrap();
        let mut balance = match self.balances.get(&sender_id, &token_id) {
            Some(balance) => balance,
            None => AccountBalance::default(),
        };
        self.process_deposit(amount, &mut balance);

        self.balances.set(&sender_id, &token_id, balance);

        // calculate storage cost for the account
        let total_storage_cost = self.calculate_storage_cost(initial_storage_usage);

        self.storage_balances
            .set(&sender_id, &token_id, total_storage_cost.into());

        PromiseOrValue::Value(U128::from(0))
    }

    #[payable]
    pub fn deposit(&mut self) {
        let account = env::predecessor_account_id();
        let amount = env::attached_deposit();
        self.assert_have_minimum_amount(amount);
        let token_id = Self::hash_token_id(&self.native_coin_name);

        let initial_storage_usage = env::storage_usage();
        let mut balance = match self.balances.get(&account, &token_id) {
            Some(balance) => balance,
            None => AccountBalance::default(),
        };
        self.process_deposit(amount, &mut balance);
        self.balances.set(&account, &token_id, balance);

        let total_storage_cost = self.calculate_storage_cost(initial_storage_usage);
        self.storage_balances
            .set(&account, &token_id, total_storage_cost.into());
    }

    #[payable]
    pub fn withdraw(&mut self, token_name: String, amount: U128) {
        let amount: u128 = amount.into();
        let account = env::predecessor_account_id();
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        self.assert_have_minimum_amount(amount);
        self.assert_have_sufficient_deposit(&account, &token_id, amount, None);

        self.assert_minimum_one_yocto();
        // Check for attached storage usage cost
        self.assert_have_sufficient_storage_deposit(&account, &token_id);

        // Check if current account have sufficient balance
        self.assert_have_sufficient_balance(amount);

        let token = self.tokens.get(&token_id).unwrap();

        let transfer_promise = if token.network() != &self.network {
            ext_nep141::ext(token.metadata().uri().to_owned().unwrap())
                .with_attached_deposit(1)
                .ft_transfer(account.clone(), amount.into(), None)
        } else if let Some(uri) = token.metadata().uri_deref() {
            ext_ft::ext(uri).with_attached_deposit(1).ft_transfer(
                account.clone(),
                U128::from(amount),
                None,
            )
        } else {
            Promise::new(account.clone()).transfer(amount)
        };

        transfer_promise.then(
            Self::ext(env::current_account_id()).on_withdraw(account, amount, token_name, token_id),
        );
    }

    pub fn reclaim(&mut self, token_name: String, amount: U128) {
        let amount: u128 = amount.into();
        let account = env::predecessor_account_id();
        self.assert_have_minimum_amount(amount);
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        self.assert_have_sufficient_refundable(&account, &token_id, amount);

        let mut balance = self.balances.get(&account, &token_id).unwrap();
        balance.refundable_mut().sub(amount).unwrap();
        balance.deposit_mut().add(amount).unwrap();

        self.balances.set(&account, &token_id, balance);
    }

    pub fn locked_balance_of(&self, account_id: AccountId, token_name: String) -> U128 {
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let balance = self
            .balances
            .get(&account_id, &token_id)
            .unwrap_or_else(|| env::panic_str(format!("{}", BshError::AccountNotExist).as_str()));
        balance.locked().into()
    }

    pub fn refundable_balance_of(&self, account_id: AccountId, token_name: String) -> U128 {
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let balance = self
            .balances
            .get(&account_id, &token_id)
            .unwrap_or_else(|| env::panic_str(format!("{}", BshError::AccountNotExist).as_str()));
        balance.refundable().into()
    }

    #[cfg(feature = "testable")]
    pub fn account_balance(
        &self,
        owner_id: AccountId,
        token_name: String,
    ) -> Option<AccountBalance> {
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        self.balances.get(&owner_id, &token_id)
    }

    pub fn balance_of(&self, account_id: AccountId, token_name: String) -> U128 {
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let balance = self
            .balances
            .get(&account_id, &token_id)
            .unwrap_or_else(|| env::panic_str(format!("{}", BshError::AccountNotExist).as_str()));
        balance.deposit().into()
    }

    #[private]
    pub fn on_withdraw(
        &mut self,
        account: AccountId,
        amount: u128,
        token_name: String,
        token_id: TokenId,
    ) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                let mut balance = self.balances.get(&account, &token_id).unwrap();
                balance.deposit_mut().sub(amount).unwrap();
                self.balances.set(&account.clone(), &token_id, balance);
                let log = json!(
                {
                    "event": "Withdraw",
                    "code": "0",
                    "by": account,
                    "amount": amount.to_string(),
                    "token_name": token_name
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());

                self.storage_balances.set(&account, &token_id, 0)
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
                    "token_name": token_name
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    pub fn get_storage_balance(&self, account: AccountId, token_name: String) -> U128 {
        let token_id = self.token_id(&token_name).unwrap();
        match self.storage_balances.get(&account, &token_id) {
            Some(storage_cost) => U128::from(storage_cost),
            None => U128::from(0),
        }
    }
}
