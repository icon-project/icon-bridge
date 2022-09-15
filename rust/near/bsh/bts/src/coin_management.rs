use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Coin Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    /// Register Coin, Accept coin meta(name, symbol, network, denominator) as parameters
    // TODO: Complete Documentation
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

            self.coin_ids.add(coin.name(), &coin_id);
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
                        "name": coin.name(),
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

    pub fn coins(&self) -> Value {
        to_value(self.coins.to_vec()).unwrap()
    }

    #[private]
    pub fn on_mint(
        &mut self,
        amount: u128,
        coin_id: CoinId,
        coin_symbol: String,
        receiver_id: AccountId,
    ) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                let mut balance = self
                    .balances
                    .get(&env::current_account_id(), &coin_id)
                    .unwrap();
                balance.deposit_mut().add(amount).unwrap();
                self.balances
                    .set(&env::current_account_id(), &coin_id, balance);

                self.internal_transfer(&env::current_account_id(), &receiver_id, &coin_id, amount);
                let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();
                let log = json!(
                {

                    "event": "Mint",
                    "code": "0",
                    "amount": amount.to_string(),
                    "token_name": coin_name

                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
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
                    "token_name": coin_name

                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    #[private]
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
                let coin_name = self.coins.get(&coin_id).unwrap().name().to_string();
                let log = json!(
                {
                    "event": "Burn",
                    "code": "0",
                    "amount": amount.to_string(),
                    "token_name": coin_name
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
                    "token_name": coin_name
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    #[private]
    pub fn register_coin_callback(&mut self, coin: Coin, coin_id: CoinId) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                self.coin_ids.add(coin.name(), &coin_id);
                self.register_coin(coin)
            }
            PromiseResult::NotReady => log!("Not Ready"),
            PromiseResult::Failed => {
                log!("Failed to register the coin")
            }
        }
    }
}

impl BtpTokenService {
    pub fn mint(&mut self, coin_id: &CoinId, amount: u128, coin: &Coin, receiver_id: AccountId) {
        ext_nep141::mint(
            amount.into(),
            coin.metadata().uri().to_owned().unwrap(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_MINT,
        )
        .then(ext_self::on_mint(
            amount,
            coin_id.to_vec(),
            coin.symbol().to_string(),
            receiver_id,
            env::current_account_id(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_ON_MINT,
        ));
    }

    pub fn burn(&mut self, coin_id: &CoinId, amount: u128, coin: &Coin) {
        ext_nep141::burn(
            amount.into(),
            coin.metadata().uri().to_owned().unwrap(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_BURN,
        )
        .then(ext_self::on_burn(
            amount,
            coin_id.to_owned(),
            coin.symbol().to_string(),
            env::current_account_id(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_FT_TRANSFER_CALL,
        ));
    }

    pub fn verify_mint(&self, coin_id: &CoinId, amount: u128) -> Result<(), String> {
        let mut balance = self
            .balances
            .get(&env::current_account_id(), coin_id)
            .unwrap();
        balance.deposit_mut().add(amount)?;
        Ok(())
    }

    pub fn register_coin(&mut self, coin: Coin) {
        let coin_id = Self::hash_coin_id(coin.name());

        self.coins.add(&coin_id, &coin);
        self.coin_fees.add(&coin_id);

        self.registered_coins.add(
            &coin.metadata().uri_deref().expect("Coin Account Missing"),
            &coin_id,
        );

        self.balances.add(&env::current_account_id(), &coin_id);
    }

    pub fn set_token_limit(
        &mut self,
        coin_names: Vec<String>,
        token_limits: Vec<u128>,
    ) -> Result<(), BshError> {
        self.assert_have_permission();
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
                    self.tokenlimits.add(coin_name, &token_limits[index])
                }
                return Ok(());
            }
            Err(err) => return Err(err),
        }
    }

    #[cfg(feature = "testable")]
    pub fn get_token_limit(&self) -> &TokenLimits {
        &self.tokenlimits
    }

    pub fn coin_id(&self, coin_name: &str) -> Result<CoinId, BshError> {
        self.coin_ids.get(coin_name).map(|coin_id| coin_id.to_owned()).ok_or(BshError::TokenNotExist {
            message: coin_name.to_string(),
        })
    }
}
