use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Token Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    /// Register Token, Accept Token meta(name, symbol, network, denominator) as parameters
    // TODO: Complete Documentation
    #[payable]
    pub fn register(&mut self, token: Token) {
        self.assert_have_permission();
        self.assert_token_does_not_exists(&token);
        let token_id = Self::hash_token_id(token.name());
        if token.network() == &self.network {
            if let Some(uri) = token.metadata().uri_deref() {
                env::promise_create(
                    uri,
                    "storage_deposit",
                    json!({}).to_string().as_bytes(),
                    env::attached_deposit(),
                    estimate::GAS_FOR_TOKEN_STORAGE_DEPOSIT,
                );
            };

            self.token_ids.add(token.name(), token_id);
            self.register_token(token);
        } else {
            let token_metadata = token.extras().clone().expect("Token Metadata Missing");
            let promise_idx = env::promise_batch_create(
                &token.metadata().uri_deref().expect("Token Account Missing"),
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
                        "spec": token_metadata.spec,
                        "name": token.label(),
                        "symbol": token.symbol(),
                        "icon": token_metadata.icon,
                        "reference": token_metadata.reference,
                        "reference_hash": token_metadata.reference_hash,
                        "decimals": token_metadata.decimals
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
                "register_token_callback",
                &json!({ "token": token,"token_id": token_id })
                    .to_string()
                    .into_bytes(),
                0,
                estimate::GAS_FOR_RESOLVE_TRANSFER,
            );
        }
    }

    // TODO: Unregister Token

    pub fn tokens(&self) -> Value {
        to_value(self.tokens.to_vec()).unwrap()
    }

    #[private]
    pub fn on_mint(
        &mut self,
        amount: u128,
        token_id: TokenId,
        receiver_id: AccountId,
        #[callback_result] storage_cost: Result<U128, near_sdk::PromiseError>,
    ) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                if let Ok(storage_cost) = storage_cost {
                    let mut balance = self
                        .balances
                        .get(&env::current_account_id(), &token_id)
                        .unwrap();
                    balance.deposit_mut().add(amount).unwrap();
                    self.balances
                        .set(&env::current_account_id(), &token_id, balance);

                    let inital_storage_used = env::storage_usage();

                    self.internal_transfer(
                        &env::current_account_id(),
                        &receiver_id,
                        &token_id,
                        amount,
                    );
                    // calculate storage cost for the account
                    let total_storage_cost = self.calculate_storage_cost(inital_storage_used);
                    let mut storage_balance =
                        match self.storage_balances.get(&receiver_id.clone(), &token_id) {
                            Some(balance) => balance,
                            None => u128::default(),
                        };

                    storage_balance
                        .add(storage_cost.0)
                        .unwrap()
                        .add(total_storage_cost.0)
                        .unwrap();

                    self.storage_balances
                        .set(&receiver_id, &token_id, storage_balance);

                    let token_name = self.tokens.get(&token_id).unwrap().name().to_string();
                    let log = json!(
                    {
                        "event": "Mint",
                        "code": "0",
                        "amount": amount.to_string(),
                        "token_name": token_name,
                        "token_account": env::signer_account_id().to_string()

                    });
                    log!(near_sdk::serde_json::to_string(&log).unwrap());
                } else {
                    let token_name = self.tokens.get(&token_id).unwrap().name().to_string();

                    let log = json!(
                    {
                        "event": "Mint",
                        "code": "1",
                        "amount": amount.to_string(),
                        "token_name": token_name,
                        "token_account": env::signer_account_id().to_string()

                    });
                    log!(near_sdk::serde_json::to_string(&log).unwrap());
                }
            }
            PromiseResult::NotReady => {
                log!("Not Ready")
            }
            PromiseResult::Failed => {
                let token_name = self.tokens.get(&token_id).unwrap().name().to_string();

                let log = json!(
                {
                    "event": "Mint",
                    "code": "1",
                    "amount": amount.to_string(),
                    "token_name": token_name,
                    "token_account": env::signer_account_id().to_string()

                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    #[private]
    pub fn on_burn(&mut self, amount: u128, token_id: TokenId) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                let mut balance = self
                    .balances
                    .get(&env::current_account_id(), &token_id)
                    .unwrap();
                balance.deposit_mut().sub(amount).unwrap();
                self.balances
                    .set(&env::current_account_id(), &token_id, balance);
                let token_name = self.tokens.get(&token_id).unwrap().name().to_string();
                let log = json!(
                {
                    "event": "Burn",
                    "code": "0",
                    "amount": amount.to_string(),
                    "token_name": token_name,
                    "token_account": env::signer_account_id().to_string()
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
            PromiseResult::NotReady => log!("Not Ready"),
            PromiseResult::Failed => {
                let token_name = self.tokens.get(&token_id).unwrap().name().to_string();
                let log = json!(
                {
                    "event": "Burn",
                    "code": "1",
                    "amount": amount.to_string(),
                    "token_name": token_name,
                    "token_account": env::signer_account_id().to_string()
                });
                log!(near_sdk::serde_json::to_string(&log).unwrap());
            }
        }
    }

    #[private]
    pub fn register_token_callback(&mut self, token: Token, token_id: TokenId) {
        match env::promise_result(0) {
            PromiseResult::Successful(_) => {
                self.token_ids.add(token.name(), token_id);
                self.register_token(token)
            }
            PromiseResult::NotReady => log!("Not Ready"),
            PromiseResult::Failed => {
                log!("Failed to register the token")
            }
        }
    }

    pub fn token(&self, token_name: String) -> Asset<FungibleToken> {
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        self.tokens.get(&token_id).unwrap()
    }

    pub fn get_token_limits(&self) -> Vec<TokenLimit> {
        self.token_limits.to_vec()
    }

    pub fn get_token_limit(&self, token_name: String) -> U128 {
        self.token_limits
            .get(&token_name)
            .map(U128)
            .unwrap_or_else(|| env::panic_str(&format!("{}", BshError::LimitNotSet)))
    }
}

impl BtpTokenService {
    pub fn mint(
        &mut self,
        token_id: &TokenId,
        amount: u128,
        token: &Token,
        receiver_id: AccountId,
    ) {
        ext_nep141::ext(token.metadata().uri().to_owned().unwrap())
            .mint(amount.into(), receiver_id.clone())
            .then(Self::ext(env::current_account_id()).on_mint(amount, *token_id, receiver_id));
    }

    pub fn burn(&mut self, token_id: &TokenId, amount: u128, token: &Token) {
        ext_nep141::ext(token.metadata().uri().to_owned().unwrap())
            .burn(amount.into())
            .then(Self::ext(env::current_account_id()).on_burn(amount, token_id.to_owned()));
    }

    pub fn verify_mint(&self, token_id: &TokenId, amount: u128) -> Result<(), String> {
        let mut balance = self
            .balances
            .get(&env::current_account_id(), token_id)
            .unwrap();
        balance.deposit_mut().add(amount)?;
        Ok(())
    }

    pub fn register_token(&mut self, token: Token) {
        let token_id = Self::hash_token_id(token.name());

        self.tokens.add(&token_id, &token);
        self.token_fees.add(&token_id);

        self.registered_tokens.add(
            &token.metadata().uri_deref().expect("Token Account Missing"),
            &token_id,
        );

        self.balances.add(&env::current_account_id(), &token_id);
        let log = json!(
        {
            "event": "Register",
            "code": "0",
            "token_name": token.name(),
            "token_account": token.metadata().uri()

        });
        log!(near_sdk::serde_json::to_string(&log).unwrap());
    }

    pub fn set_token_limit(
        &mut self,
        token_names: Vec<String>,
        token_limits: Vec<u128>,
    ) -> Result<(), BshError> {
        match self.ensure_length_matches(&token_names, &token_limits) {
            Ok(()) => {
                let mut invalid_tokens: Vec<String> = Vec::new();
                let mut valid_tokens: Vec<String> = Vec::new();
                token_names.into_iter().for_each(|token_name| {
                    match self.ensure_token_exists(&token_name) {
                        true => valid_tokens.push(token_name),
                        false => invalid_tokens.push(token_name),
                    }
                });

                if !invalid_tokens.is_empty() {
                    return Err(BshError::TokenNotExist {
                        message: invalid_tokens.join(", "),
                    });
                }

                for (index, token_name) in valid_tokens.iter().enumerate() {
                    self.token_limits.add(token_name, &token_limits[index])
                }
                Ok(())
            }
            Err(err) => Err(err),
        }
    }

    pub fn token_id(&self, token_name: &str) -> Result<TokenId, BshError> {
        self.token_ids
            .get(token_name)
            .map(|token_id| token_id.to_owned())
            .ok_or(BshError::TokenNotExist {
                message: token_name.to_string(),
            })
    }
}
