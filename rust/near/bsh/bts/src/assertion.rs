use super::*;

impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Internal Validations  * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn assert_predecessor_is_bmc(&self) {
        require!(
            env::predecessor_account_id() == *self.bmc(),
            format!("{}", BshError::NotBmc)
        )
    }

    pub fn assert_coin_id_len_match_amount_len(&self, coin_ids: &Vec<CoinId>, amounts: &Vec<U128>) {
        require!(
            coin_ids.len() == amounts.len(),
            format!(
                "{}",
                BshError::InvalidCount {
                    message: "Coin Ids and amounts".to_string()
                }
            ),
        );
    }

    pub fn assert_transfer_amounts_len_match_returned_amount_len(
        &self,
        amounts: &Vec<U128>,
        returned_amount: &Vec<U128>,
    ) {
        require!(
            returned_amount.len() == amounts.len(),
            format!(
                "{}",
                BshError::InvalidCount {
                    message: "Transfer amounts and returned amounts".to_string()
                }
            ),
        );
    }

    pub fn assert_valid_fee_ratio(&self, fee_numerator: u128, fixed_fee: u128) {
        require!(
            fee_numerator <= FEE_DENOMINATOR,
            format!("{}", BshError::InvalidSetting),
        );

        require!(
            fee_numerator >= 0 && fixed_fee >= 0,
            format!("{}", BshError::LessThanZero),
        );
    }

    pub fn assert_valid_service(&self, service: &String) {
        require!(
            self.name() == service,
            format!("{}", BshError::InvalidService)
        )
    }

    /// Check whether signer account id is an owner
    pub fn assert_have_permission(&self) {
        require!(
            self.owners.contains(&env::predecessor_account_id()),
            format!("{}", BshError::PermissionNotExist)
        );
    }

    pub fn assert_have_minimum_amount(&self, amount: u128) {
        require!(amount > 0, format!("{}", BshError::NotMinimumAmount));
    }

    pub fn assert_have_sufficient_balance(&self, amount: u128) {
        require!(
            env::account_balance() > amount,
            format!(
                "{}",
                BshError::NotMinimumBalance {
                    account: env::current_account_id().to_string()
                }
            )
        );
    }

    pub fn assert_have_sufficient_deposit(
        &self,
        account: &AccountId,
        coin_id: &CoinId,
        amount: u128,
        fees: Option<u128>,
    ) {
        let amount = std::cmp::max(amount, fees.unwrap_or_default());
        if let Some(balance) = self.balances.get(&account, &coin_id) {
            require!(
                balance.deposit() >= amount,
                format!("{}", BshError::NotMinimumDeposit)
            );
        } else {
            env::panic_str(format!("{}", BshError::NotMinimumDeposit).as_str());
        }
    }

    pub fn assert_have_sufficient_refundable(
        &self,
        account: &AccountId,
        coin_id: &CoinId,
        amount: u128,
    ) {
        if let Some(balance) = self.balances.get(&account, &coin_id) {
            require!(
                balance.refundable() >= amount,
                format!("{}", BshError::NotMinimumRefundable)
            );
        } else {
            env::panic_str(format!("{}", BshError::NotMinimumRefundable).as_str());
        }
    }

    pub fn assert_sender_is_not_receiver(&self, sender_id: &AccountId, receiver_id: &AccountId) {
        require!(
            sender_id != receiver_id,
            format!("{}", BshError::SameSenderReceiver)
        );
    }

    pub fn assert_owner_exists(&self, account: &AccountId) {
        require!(
            self.owners.contains(&account),
            format!("{}", BshError::OwnerNotExist)
        );
    }

    pub fn assert_owner_does_not_exists(&self, account: &AccountId) {
        require!(
            !self.owners.contains(account),
            format!("{}", BshError::OwnerExist)
        );
    }

    pub fn assert_owner_is_not_last_owner(&self) {
        require!(self.owners.len() > 1, format!("{}", BshError::LastOwner));
    }

    pub fn assert_coin_does_not_exists(&self, coin: &Coin) {
        let coin = self.coins.get(&Self::hash_coin_id(coin.name()));
        require!(coin.is_none(), format!("{}", BshError::TokenExist))
    }

    pub fn assert_coins_exists(&self, coin_ids: &Vec<CoinId>) {
        let mut unregistered_coins: Vec<CoinId> = vec![];
        coin_ids.iter().for_each(|coin_id| {
            if !self.coins.contains(&coin_id) {
                unregistered_coins.push(coin_id.to_owned())
            }
        });

        require!(
            unregistered_coins.len() == 0,
            format!(
                "{}",
                BshError::TokenNotExist {
                    message: unregistered_coins
                        .iter()
                        .map(|coin_id| format!("{:x?}", coin_id))
                        .collect::<Vec<String>>()
                        .join(", "),
                }
            ),
        );
    }
    pub fn assert_coin_registered(&self, coin_account: &AccountId) {
        require!(
            self.registered_coins.contains(coin_account),
            format!("{}", BshError::TokenNotRegistered)
        )
    }

    pub fn ensure_user_blacklisted(&self, user: &AccountId) -> Result<(), BshError> {
        if !self.blacklisted_accounts.contains(user) {
            return Err(BshError::UserNotBlacklisted);
        }
        Ok(())
    }

    pub fn ensure_length_matches(
        &self,
        coin_names: &Vec<String>,
        token_limits: &Vec<u128>,
    ) -> Result<(), BshError> {
        if coin_names.len() != token_limits.len() {
            return Err(BshError::InvalidParams);
        }
        Ok(())
    }

    pub fn ensure_coin_exists(&self, coin_name: &String) -> bool {
        match self.coin_ids.get(coin_name) {
            Some(coin) => true,
            None => false,
        }
    }

    pub fn get_coin_ids(&self, coin_name: &Vec<String>) -> Result<Vec<CoinId>, BshError> {
        let mut invlaid_coins: Vec<String> = vec![];
        let mut coin_ids: Vec<Vec<u8>> = vec![];
        coin_name
            .into_iter()
            .for_each(|name| match self.coin_ids.get(&name) {
                Some(coin_id) => coin_ids.push(coin_id.to_vec()),
                None => invlaid_coins.push(name.to_string()),
            });

        require!(
            invlaid_coins.len() == 0,
            format!(
                "{}",
                BshError::TokenNotExist {
                    message: invlaid_coins
                        .iter()
                        .map(|coin_id| format!("{:x?}", coin_id))
                        .collect::<Vec<String>>()
                        .join(", "),
                }
            ),
        );
        Ok(coin_ids)
    }

    pub fn get_coin_id(&self, coin_name: &String) -> Result<CoinId, BshError> {
        let mut id: Vec<u8> = vec![];
        let is_valid = match self.coin_ids.get(&coin_name) {
            Some(coin_id) => {
                id = coin_id.to_vec();
                1
            }
            None => 0,
        };

        require!(
            is_valid != 0,
            format!(
                "{}",
                BshError::TokenNotExist {
                    message: coin_name.to_string()
                }
            ),
        );

        Ok(id)
    }
}
