use libraries::types::{FungibleToken, TokenLimit};

use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Coin Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    /// Register Coin, Accept coin meta(name, symbol, network, denominator) as parameters

    ///

    #[payable]
    pub fn register(&mut self, coin: Coin) {
        self.assert_have_permission();
        self.assert_coin_does_not_exists(&coin);
        let coin_id = Self::hash_coin_id(coin.name());
        if coin.network() == &self.network {
            if let Some(uri) = coin.metadata().uri_deref() {
                env::promise_create(
                    uri,
                    "storage_deposit",
                    &json!({}).to_string().as_bytes(),
                    env::attached_deposit(),
                    estimate::GAS_FOR_TOKEN_STORAGE_DEPOSIT,
                );
            };

            self.coin_ids.add(coin.name(), coin_id);
            self.register_coin(coin);
        } else {
            let coin_metadata = coin.extras().clone().expect("Coin Metadata Missing");
            let promise_idx = env::promise_batch_create(
                &coin.metadata().uri_deref().expect("Coin Account Missing"),
            );
            env::promise_batch_action_create_account(promise_idx);
            env::promise_batch_action_transfer(promise_idx, env::attached_deposit());
            env::promise_batch_action_deploy_contract(promise_idx, NEP141_CONTRACT);
            env::promise_batch_action_function_call(
                promise_idx,
                "new",
                &json!({
                    "owner_id": env::current_account_id(),
                    "total_supply": U128(0),
                    "metadata": {
                        "spec": coin_metadata.spec.clone(),
                        "name": coin.label(),
                        "symbol": coin.symbol(),
                        "icon": coin_metadata.icon.clone(),
                        "reference": coin_metadata.reference.clone(),
                        "reference_hash": coin_metadata.reference_hash.clone(),
                        "decimals": coin_metadata.decimals.clone()
                    }
                })
                .to_string()
                .into_bytes(),
                estimate::NO_DEPOSIT,
                estimate::GAS_FOR_RESOLVE_TRANSFER,
            );
            env::promise_then(
                promise_idx,
                env::current_account_id(),
                "register_coin_callback",
                &json!({ "coin": coin,"coin_id": coin_id })
                    .to_string()
                    .into_bytes(),
                0,
                estimate::GAS_FOR_RESOLVE_TRANSFER,
            );
        }
    }

    // TODO: Unregister Token

/// Querying the coins in btp

    pub fn coins(&self) -> Value {
        to_value(self.coins.to_vec()).unwrap()
    }

    #[private]
    /// Method on mint is created
    /// # Arguments
    /// * `amount` - should be in unsigned integer
    /// * `coin_id` - coin id should be provided
    /// * `coin_symbol` - Should be in a string format
    /// * `receiver_id` - The Account Id of the user should be given
    /// 
    pub fn on_mint(
        &mut self,
        amount: u128,
        coin_id: CoinId,
        coin_symbol: String,
        receiver_id: AccountId,
        #[callback_result] storage_cost: Result<U128, near_sdk::PromiseError>,
    ) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                if storage_cost.is_ok() {
                    let mut balance = self
                        .balances
                        .get(&env::current_account_id(), &coin_id)
                        .unwrap();
                    balance.deposit_mut().add(amount).unwrap();
                    self.balances
                        .set(&env::current_account_id(), &coin_id, balance);
                    //initial storage useage can be get
                    let inital_storage_used = env::storage_usage();

                    self.internal_transfer(
                        &env::current_account_id(),
                        &receiver_id,
                        &coin_id,
                        amount,
                    );
                    // calculate storage cost for the account
                    let total_storage_cost = self.calculate_storage_cost(inital_storage_used);
                    let mut storage_balance =
                        match self.storage_balances.get(&receiver_id.clone(), &coin_id) {
                            Some(balance) => balance,
                            None => u128::default(),
                        };

                    storage_balance
                        .add(storage_cost.unwrap().0)
                        .unwrap()
                        .add(total_storage_cost.0)
                        .unwrap();

                    self.storage_balances
                        .set(&receiver_id, &coin_id, storage_balance);

                    let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();
                    let log = json!(
                    {
                        "event": "Mint",
                        "code": "0",
                        "amount": amount.to_string(),
                        "token_name": coin_name,
                        "token_account": env::signer_account_id().to_string()

                    });
                    log!(near_sdk::serde_json::to_string(&log).unwrap());
                } else {
                    let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();

                    let log = json!(
                    {
                        "event": "Mint",
                        "code": "1",
                        "amount": amount.to_string(),
                        "token_name": coin_name,
                        "token_account": env::signer_account_id().to_string()

                    });
                    log!(near_sdk::serde_json::to_string(&log).unwrap());
                }
            }
            PromiseResult::NotReady => {
                log!("Not Ready")
            }
            PromiseResult::Failed => {
                let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();

                let log = json!(
                {
                    "event": "Mint",
                    "code": "1",
                    "amount": amount.to_string(),
                    "token_name": coin_name,
                    "token_account": env::signer_account_id().to_string()

                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    #[private]
    /// on burn method got created in btp
    /// # Arguments
    /// * `amount` - should be in unsigned integer
    /// * `coin_id` - Coin id should be given
    /// * `coin_symbol` - It should be in a string format
    /// 
    pub fn on_burn(&mut self, amount: u128, coin_id: CoinId, coin_symbol: String) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                let mut balance = self
                    .balances
                    .get(&env::current_account_id(), &coin_id)
                    .unwrap();
                balance.deposit_mut().sub(amount).unwrap();
                self.balances
                    .set(&env::current_account_id(), &coin_id, balance);
                    // getting the coin name
                let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();
                let log = json!(
                {
                    "event": "Burn",
                    "code": "0",
                    "amount": amount.to_string(),
                    "token_name": coin_name,
                    "token_account": env::signer_account_id().to_string()
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
            PromiseResult::NotReady => log!("Not Ready"),
            PromiseResult::Failed => {
                let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();
                let log = json!(
                {
                    "event": "Burn",
                    "code": "1",
                    "amount": amount.to_string(),
                    "token_name": coin_name,
                    "token_account": env::signer_account_id().to_string()
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    #[private]
    /// Call back the registered coin
    /// # Arguments
    /// * `coin` : Should give the parameters like name, symbol, fee_numerator, denominator, network, fixed fee, uri, extras
    /// 
    pub fn register_coin_callback(&mut self, coin: Coin, coin_id: CoinId) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                self.coin_ids.add(coin.name(), coin_id);
                self.register_coin(coin)
            }
            PromiseResult::NotReady => log!("Not Ready"),
            PromiseResult::Failed => {
                log!("Failed to register the coin")
            }
        }
    }
/// Method coin is created 
/// # Arguments
/// * `coin_name` - name of the coin should be given in string format
/// 
    pub fn coin(&self, coin_name: String) -> Asset<FungibleToken> {
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        self.coins.get(&coin_id).unwrap()
    }
/// returns the token limits 
    pub fn get_token_limits(&self) -> Vec<TokenLimit> {
        self.token_limits.to_vec()
    }
/// getting the token limit by giving the particular coin name
/// # Arguments
/// * `coin_name` - name of the coin should be in the string format
/// 
    pub fn get_token_limit(&self, coin_name: String) -> U128 {
        self.token_limits
            .get(&coin_name)
            .map(|token_limit| U128(token_limit))
            .expect(&format!("{}", BshError::LimitNotSet))
    }
}

impl BtpTokenService {
    /// mint method is created
    /// # Arguments
    /// * `coin_id` - coin id should be given
    /// * `amount` - should be a unsigned number
    /// * `coin` - Should give the parameters like name, symbol, fee_numerator, denominator, network, fixed fee, uri, extras
    /// * `receiver_id` - Should give the existence account id.
    /// 
    pub fn mint(&mut self, coin_id: &CoinId, amount: u128, coin: &Coin, receiver_id: AccountId) {
        ext_nep141::ext(coin.metadata().uri().to_owned().unwrap())
            .mint(amount.into(), receiver_id.clone())
            .then(Self::ext(env::current_account_id()).on_mint(
                amount,
                *coin_id,
                coin.symbol().to_string(),
                receiver_id,
            ));
    }

    /// method burn got created
    /// # Arguments
    /// * `coin_id` - should give the coin_id of the particular coin
    /// * `amount` - should be an unsigned number
    /// * `coin` - Should give the parameters like name, symbol, fee_numerator, denominator, network, fixed fee, uri, extras
    /// 
    pub fn burn(&mut self, coin_id: &CoinId, amount: u128, coin: &Coin) {
        ext_nep141::ext(coin.metadata().uri().to_owned().unwrap())
            .burn(amount.into())
            .then(Self::ext(env::current_account_id()).on_burn(
                amount,
                coin_id.to_owned(),
                coin.symbol().to_string(),
            ));
    }
    /// verify mint method got created
    /// # Arguments
    /// * `coin_id` - Should give the coin id of the particular coin
    /// * `amount` - should be an unsigned number
    /// 
    pub fn verify_mint(&self, coin_id: &CoinId, amount: u128) -> Result<(), String> {
        let mut balance = self
            .balances
            .get(&env::current_account_id(), coin_id)
            .unwrap();
        balance.deposit_mut().add(amount)?;
        Ok(())
    }

    /// Coin got registered
    /// # Arguments
    /// * `coin` - Should give the parameters like name, symbol, fee_numerator, denominator, network, fixed fee, uri, extras
    /// 
    pub fn register_coin(&mut self, coin: Coin) {
        let coin_id = Self::hash_coin_id(coin.name());

        self.coins.add(&coin_id, &coin);
        self.coin_fees.add(&coin_id);

        self.registered_coins.add(
            &coin.metadata().uri_deref().expect("Coin Account Missing"),
            &coin_id,
        );

        self.balances.add(&env::current_account_id(), &coin_id);
        let log = json!(
        {
            "event": "Register",
            "code": "0",
            "token_name": coin.name(),
            "token_account": coin.metadata().uri()

        });
        log!(near_sdk::serde_json::to_string(&log).unwrap());
    }
    /// setting the token limit of the coin
    /// caller should be a owner
    /// # Arguments
    /// * `coin_names` - should be in the vector format
    /// * `token-limits` - Given in the vector format
    /// 
    pub fn set_token_limit(
        &mut self,
        coin_names: Vec<String>,
        token_limits: Vec<u128>,
    ) -> Result<(), BshError> {
        match self.ensure_length_matches(&coin_names, &token_limits) {
            Ok(()) => {
                let mut invalid_coins: Vec<String> = Vec::new();
                let mut valid_coins: Vec<String> = Vec::new();
                coin_names.into_iter().for_each(|coin_name| {
                    match self.ensure_coin_exists(&coin_name) {
                        true => valid_coins.push(coin_name),
                        false => invalid_coins.push(coin_name),
                    }
                });

                if !invalid_coins.is_empty() {
                    return Err(BshError::TokenNotExist {
                        message: invalid_coins.join(", "),
                    });
                }

                for (index, coin_name) in valid_coins.iter().enumerate() {
                    self.token_limits.add(coin_name, &token_limits[index])
                }
                return Ok(());
            }
            Err(err) => return Err(err),
        }
    }
    /// method coin_id is created by using coin name
    /// # Arguments
    /// `coin_name` - name of the coin should be given in the string format
    /// 
    pub fn coin_id(&self, coin_name: &str) -> Result<CoinId, BshError> {
        self.coin_ids
            .get(coin_name)
            .map(|coin_id| coin_id.to_owned())
            .ok_or(BshError::TokenNotExist {
                message: coin_name.to_string(),
            })
    }
}
