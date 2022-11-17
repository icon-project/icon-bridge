use super::*;

#[near_bindgen]
impl BtpTokenService {
    pub fn accumulated_fees(&self) -> Vec<AccumulatedAssetFees> {
        self.tokens
            .to_vec()
            .iter()
            .map(|token| {
                let token_id = *self.token_ids.get(&token.name).unwrap();
                let token_fee = self.token_fees.get(&token_id).unwrap();
                let token = self.tokens.get(&token_id).unwrap();

                AccumulatedAssetFees {
                    name: token.name().clone(),
                    network: token.network().clone(),
                    accumulated_fees: token_fee,
                }
            })
            .collect()
    }

    pub fn handle_fee_gathering(&mut self, fee_aggregator: BTPAddress, service: String) {
        self.assert_predecessor_is_bmc();
        self.assert_valid_service(&service);

        self.transfer_fees(&fee_aggregator);
    }

    pub fn set_fee_ratio(&mut self, token_name: String, fee_numerator: U128, fixed_fee: U128) {
        self.assert_have_permission();

        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let mut token = self.tokens.get(&token_id).unwrap();

        token
            .metadata_mut()
            .fee_numerator_mut()
            .clone_from(&fee_numerator.into());

        token
            .metadata_mut()
            .fixed_fee_mut()
            .clone_from(&fixed_fee.into());

        self.tokens.set(&token_id, &token)
    }

    pub fn get_fee(&self, token_name: String, amount: U128) -> U128 {
        let token_id = self
            .token_id(&token_name)
            .map_err(|err| format!("{}", err))
            .unwrap();

        let token = self.tokens.get(&token_id).unwrap();

        self.calculate_token_transfer_fee(u128::from(amount), &token)
            .map(U128)
            .unwrap()
    }
}

impl BtpTokenService {
    pub fn transfer_fees(&mut self, fee_aggregator: &BTPAddress) {
        let sender_id = env::current_account_id();

        let assets = self
            .tokens
            .to_vec()
            .iter()
            .filter_map(|token| {
                let token_id = *self.token_ids.get(&token.name).unwrap();
                let token_fee = self.token_fees.get(&token_id).unwrap();

                if token_fee > 0 {
                    self.token_fees.set(&token_id, 0);

                    Some(
                        self.process_external_transfer(&token_id, &sender_id, token_fee)
                            .unwrap(),
                    )
                } else {
                    None
                }
            })
            .collect::<Vec<TransferableAsset>>();

        self.send_request(sender_id, fee_aggregator.clone(), assets);
    }

    pub fn calculate_token_transfer_fee(
        &self,
        amount: u128,
        token: &Asset<FungibleToken>,
    ) -> Result<u128, String> {
        let mut fee = (amount * token.metadata().fee_numerator()) / FEE_DENOMINATOR;

        fee.add(token.metadata().fixed_fee()).map(|fee| *fee)
    }
}
