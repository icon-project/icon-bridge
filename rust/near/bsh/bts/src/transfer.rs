use super::*;

#[near_bindgen]
impl BtpTokenService {
    #[payable]
    pub fn transfer(&mut self, token_name: String, destination: BTPAddress, amount: U128) {
        let sender_id = env::predecessor_account_id();
        self.assert_have_minimum_amount(amount.into());

        let token_id = self
            .token_id(&token_name)
            .map_err(|error| error.to_string())
            .unwrap();

        //check for enough attached deposit to bear storage cost
        self.assert_have_sufficient_storage_deposit(&sender_id, &token_id);

        let asset = self
            .process_external_transfer(&token_id, &sender_id, amount.into())
            .unwrap();
        self.send_request(sender_id.clone(), destination, vec![asset]);

        self.storage_balances.set(&sender_id, &token_id, 0)
    }

    #[payable]
    pub fn transfer_batch(
        &mut self,
        token_names: Vec<String>,
        destination: BTPAddress,
        amounts: Vec<U128>,
    ) {
        let sender_id = env::predecessor_account_id();

        let token_ids = token_names
            .iter()
            .map(|token_name| self.token_id(token_name))
            .collect::<Result<Vec<TokenId>, BshError>>()
            .map_err(|error| error.to_string())
            .unwrap();

        let mut storage_cost = u128::default();

        token_ids.clone().into_iter().for_each(|token_id| {
            let storage_balance = match self.storage_balances.get(&sender_id.clone(), &token_id) {
                Some(balance) => balance,
                None => u128::default(),
            };

            storage_cost.add(storage_balance).unwrap();
        });

        //check for enough attached deposit to bear storage cost
        self.assert_have_sufficient_storage_deposit_for_batch(storage_cost);

        let assets = token_ids
            .iter()
            .enumerate()
            .map(|(index, token_id)| {
                self.assert_have_minimum_amount(amounts[index].into());
                self.process_external_transfer(token_id, &sender_id, amounts[index].into())
                    .unwrap()
            })
            .collect::<Vec<TransferableAsset>>();

        self.send_request(sender_id.clone(), destination, assets);

        token_ids
            .into_iter()
            .for_each(|token_id| self.storage_balances.set(&sender_id, &token_id, 0));
    }
}

impl BtpTokenService {
    pub fn process_external_transfer(
        &mut self,
        token_id: &TokenId,
        sender_id: &AccountId,
        mut amount: u128,
    ) -> Result<TransferableAsset, String> {
        let token = self.tokens.get(token_id).unwrap();
        let fees = self.calculate_token_transfer_fee(amount, &token)?;

        self.assert_have_sufficient_deposit(sender_id, token_id, amount, Some(fees));

        amount.sub(fees)?;
        let mut balance = self.balances.get(sender_id, token_id).unwrap();

        // Handle Fees
        balance.locked_mut().add(fees)?;
        balance.deposit_mut().sub(fees)?;

        // Handle Deposit
        balance.deposit_mut().sub(amount)?;
        balance.locked_mut().add(amount)?;

        self.balances.set(sender_id, token_id, balance);

        Ok(TransferableAsset::new(token.name().clone(), amount, fees))
    }

    pub fn internal_transfer(
        &mut self,
        sender_id: &AccountId,
        receiver_id: &AccountId,
        token_id: &TokenId,
        amount: u128,
    ) {
        self.assert_sender_is_not_receiver(sender_id, receiver_id);
        self.assert_have_sufficient_deposit(sender_id, token_id, amount, None); //TODO: Convert to ensure

        let mut sender_balance = self.balances.get(sender_id, token_id).unwrap();
        sender_balance.deposit_mut().sub(amount).unwrap();

        let receiver_balance = match self.balances.get(receiver_id, token_id) {
            Some(mut balance) => {
                balance.deposit_mut().add(amount).unwrap();
                balance
            }
            None => {
                let mut balance = AccountBalance::default();
                let amount = amount;
                balance.deposit_mut().add(amount).unwrap();
                balance
            }
        };

        self.balances.set(sender_id, token_id, sender_balance);

        self.balances.set(receiver_id, token_id, receiver_balance);
    }

    pub fn verify_internal_transfer(
        &self,
        sender_id: &AccountId,
        receiver_id: &AccountId,
        token_id: &TokenId,
        amount: u128,
        sender_balance: &mut AccountBalance,
    ) -> Result<(), String> {
        self.assert_sender_is_not_receiver(sender_id, receiver_id);
        sender_balance.deposit_mut().sub(amount)?;

        match self.balances.get(receiver_id, token_id) {
            Some(mut balance) => {
                balance.deposit_mut().add(amount)?;
                balance
            }
            None => {
                let mut balance = AccountBalance::default();
                balance.deposit_mut().add(amount)?;
                balance
            }
        };

        Ok(())
    }

    pub fn internal_transfer_batch(
        &mut self,
        sender_id: &AccountId,
        receiver_id: &AccountId,
        token_ids: &[TokenId],
        amounts: &[U128],
    ) {
        token_ids.iter().enumerate().for_each(|(index, token_id)| {
            self.internal_transfer(sender_id, receiver_id, token_id, amounts[index].into());
        });
    }

    pub fn finalize_external_transfer(
        &mut self,
        sender_id: &AccountId,
        assets: &[TransferableAsset],
    ) {
        assets.iter().for_each(|asset| {
            let token_id = self.token_id(asset.name()).unwrap();
            let token = self.tokens.get(&token_id).unwrap();

            let mut token_fee = self.token_fees.get(&token_id).unwrap().to_owned();
            let mut sender_balance = self.balances.get(sender_id, &token_id).unwrap();

            sender_balance
                .locked_mut()
                .sub(asset.amount() + asset.fees())
                .unwrap();

            self.balances.set(sender_id, &token_id, sender_balance);

            let mut current_account_balance = self
                .balances
                .get(&env::current_account_id(), &token_id)
                .unwrap();

            current_account_balance
                .deposit_mut()
                .add(asset.amount() + asset.fees())
                .unwrap();

            self.balances.set(
                &env::current_account_id(),
                &token_id,
                current_account_balance,
            );

            token_fee.add(asset.fees()).unwrap();
            self.token_fees.set(&token_id, token_fee);

            if token.network() != &self.network {
                self.burn(&token_id, asset.amount(), &token);
            }
        });
    }

    pub fn rollback_external_transfer(
        &mut self,
        sender_id: &AccountId,
        assets: &[TransferableAsset],
    ) {
        assets.iter().for_each(|asset| {
            let token_id = self.token_id(asset.name()).unwrap();
            let mut token_fee = self.token_fees.get(&token_id).unwrap().to_owned();
            let mut sender_balance = self.balances.get(sender_id, &token_id).unwrap();

            sender_balance
                .locked_mut()
                .sub(asset.amount() + asset.fees())
                .unwrap();

            sender_balance.refundable_mut().add(asset.amount()).unwrap();
            self.balances.set(sender_id, &token_id, sender_balance);

            let mut current_account_balance = self
                .balances
                .get(&env::current_account_id(), &token_id)
                .unwrap();

            current_account_balance
                .deposit_mut()
                .add(asset.fees())
                .unwrap();

            self.balances.set(
                &env::current_account_id(),
                &token_id,
                current_account_balance,
            );

            token_fee.add(asset.fees()).unwrap();

            self.token_fees.set(&token_id, token_fee);
        });
    }

    pub fn handle_token_transfer(
        &mut self,
        receiver_id: &String,
        assets: &[TransferableAsset],
    ) -> Result<Option<TokenServiceMessage>, BshError> {
        let receiver_id = AccountId::try_from(receiver_id.to_owned()).map_err(|error| {
            BshError::InvalidAddress {
                message: error.to_string(),
            }
        })?;

        let mut unregistered_tokens: Vec<String> = Vec::new();

        let token_ids: Vec<(usize, TokenId)> = assets
            .iter()
            .map(|asset| {
                self.token_ids
                    .get(asset.name())
                    .copied()
                    .unwrap_or_default()
            })
            .enumerate()
            .filter(|(index, _)| {
                return if self
                    .ensure_token_exists(assets[index.to_owned()].name())
                    .is_err()
                {
                    unregistered_tokens.push(assets[index.to_owned()].name().to_owned());
                    false
                } else {
                    true
                };
            })
            .collect();

        if !unregistered_tokens.is_empty() {
            return Err(BshError::TokenNotExist {
                message: unregistered_tokens.join(", "),
            });
        }

        let tokens = token_ids
            .into_iter()
            .map(|(asset_index, token_id)| {
                (asset_index, token_id, self.tokens.get(&token_id).unwrap())
            })
            .collect::<Vec<(usize, TokenId, Token)>>();

        let transferable =
            self.is_tokens_transferable(&env::current_account_id(), &receiver_id, &tokens, assets);

        if transferable.is_err() {
            return Err(BshError::Reverted {
                message: format!("Coins not transferable: {}", transferable.unwrap_err()),
            });
        }

        tokens.iter().for_each(|(index, token_id, token)| {
            if token.network() != &self.network {
                self.mint(
                    token_id,
                    assets[index.to_owned()].amount(),
                    token,
                    receiver_id.clone(),
                );
            } else {
                self.internal_transfer(
                    &env::current_account_id(),
                    &receiver_id,
                    token_id,
                    assets[index.to_owned()].amount(),
                );
            }
        });

        Ok(Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 0,
                message: "Transfer Success".to_string(),
            },
        )))
    }

    fn is_tokens_transferable(
        &self,
        sender_id: &AccountId,
        receiver_id: &AccountId,
        tokens: &[(usize, TokenId, Asset<FungibleToken>)],
        assets: &[TransferableAsset],
    ) -> Result<(), String> {
        tokens
            .iter()
            .map(|(index, token_id, token)| -> Result<(), String> {
                self.ensure_user_not_blacklisted(receiver_id)
                    .map_err(|error| error.to_string())?;

                let mut sender_balance = self.balances.get(sender_id, token_id).unwrap();

                if let &Some(token_limit) = token.metadata().token_limit() {
                    if assets[index.to_owned()].amount() > token_limit {
                        return Err("limit exceeded".to_string());
                    }
                };

                if token.network() != &self.network {
                    self.verify_mint(token_id, assets[index.to_owned()].amount())?;

                    sender_balance
                        .deposit_mut()
                        .add(assets[index.to_owned()].amount())?;
                };

                self.verify_internal_transfer(
                    &env::current_account_id(),
                    receiver_id,
                    token_id,
                    assets[index.to_owned()].amount(),
                    &mut sender_balance,
                )?;

                Ok(())
            })
            .collect()
    }
}
