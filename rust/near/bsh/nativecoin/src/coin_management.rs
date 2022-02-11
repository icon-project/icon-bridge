use super::*;

#[near_bindgen]
impl NativeCoinService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Coin Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    /// Register Coin, Accept token meta(name, symbol, network, denominator) as parameters
    // TODO: Complete Documentation
    pub fn register(&mut self, token: Token<WrappedNativeCoin>) {
        self.assert_have_permission();
        self.assert_token_does_not_exists(&token);

        let token_id = Self::hash_token_id(token.name());
        // token id is hash(token_name)
        self.tokens.add(&token_id, &token);
        self.token_fees.add(&token_id);

        // Sets initial balance for self
        self.balances.add(&env::current_account_id(), &token_id);
    }

    // TODO: Unregister Token

    pub fn coins(&self) -> Value {
        to_value(self.tokens.to_vec()).unwrap()
    }

    // Hashing to be done out of chain
    pub fn coin_id(&self, coin_name: String) -> TokenId {
        let coin_id = Self::hash_token_id(&coin_name);
        self.assert_tokens_exists(&vec![coin_id.clone()]);
        coin_id
    }

    #[private]
    pub fn on_mint(
        &mut self,
        amount: u128,
        token_id: TokenId,
        token_symbol: String,
        receiver_id: AccountId,
    ) {
        let mut balance = self.balances.get(&receiver_id, &token_id).unwrap();
        balance.deposit_mut().add(amount).unwrap();
        self.balances.set(&receiver_id, &token_id, balance);

        log!("[Mint] {} {}", amount, token_symbol);
    }
}

impl NativeCoinService {
    pub fn mint(
        &mut self,
        token_id: &TokenId,
        amount: u128,
        token: &Token<WrappedNativeCoin>,
        receiver_id: AccountId,
    ) {
        ext_nep141::mint(
            amount.into(),
            token.metadata().uri().to_owned().unwrap(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_MT_TRANSFER_CALL,
        )
        .then(ext_self::on_mint(
            amount,
            token_id.to_vec(),
            token.symbol().to_string(),
            receiver_id,
            env::current_account_id(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_MT_TRANSFER_CALL,
        ));
    }

    pub fn burn(&mut self, token_id: &TokenId, amount: u128, token: &Token<WrappedNativeCoin>) {
        // TODO: Remove from supply
        let mut balance = self
            .balances
            .get(&env::current_account_id(), token_id)
            .unwrap();
        balance.deposit_mut().sub(amount).unwrap();
        self.balances
            .set(&env::current_account_id(), token_id, balance);

        log!("[Burn] {} {}", amount, token.symbol());
    }

    pub fn verify_mint(&self, token_id: &TokenId, amount: u128) -> Result<(), String> {
        let mut balance = self
            .balances
            .get(&env::current_account_id(), token_id)
            .unwrap();
        balance.deposit_mut().add(amount)?;
        Ok(())
    }
}
