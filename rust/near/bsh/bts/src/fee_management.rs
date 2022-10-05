use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Fee Management  * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn accumulated_fees(&self) -> Vec<AccumulatedAssetFees> {
        self.coins
            .to_vec()
            .iter()
            .map(|coin| {
                let coin_id = Self::hash_coin_id(&coin.name);
                let coin_fee = self.coin_fees.get(&coin_id).unwrap();
                let coin = self.coins.get(&coin_id).unwrap();
                AccumulatedAssetFees {
                    name: coin.name().clone(),
                    network: coin.network().clone(),
                    accumulated_fees: *coin_fee,
                }
            })
            .collect()
    }

    pub fn handle_fee_gathering(&mut self, fee_aggregator: BTPAddress, service: String) {
        self.assert_predecessor_is_bmc();
        self.assert_valid_service(&service);
        self.transfer_fees(&fee_aggregator);
    }

    pub fn set_fee_ratio(&mut self, coin_id: &CoinId, fee_numerator: U128, fixed_fee: U128) {
        self.assert_have_permission();
        self.assert_valid_fee_ratio(fee_numerator.into(), fixed_fee.into());

        let mut coin = self.coins.get(&coin_id).unwrap();
        coin.metadata_mut()
            .fee_numerator_mut()
            .add(fee_numerator.into())
            .unwrap();
        coin.metadata_mut()
            .fixed_fee_mut()
            .add(fixed_fee.into())
            .unwrap();

        self.coins.set(coin_id, &coin)
    }

    pub fn get_fee(&self, coin_name: String, amount: U128) -> Result<U128, String> {
        let coin_id = self
        .coin_id(&coin_name)
        .map_err(|err| format!("{}", err))
        .unwrap();
        let coin = self.coins.get(&coin_id).unwrap();

        self.calculate_coin_transfer_fee(u128::from(amount), &coin).map(|e| U128(e))
    }
}

impl BtpTokenService {
    pub fn transfer_fees(&mut self, fee_aggregator: &BTPAddress) {
        let sender_id = env::current_account_id();
        let assets = self
            .coins
            .to_vec()
            .iter()
            .filter_map(|coin| {
                let coin_id = Self::hash_coin_id(&coin.name);
                let coin_fee = self.coin_fees.get(&coin_id).unwrap().clone();

                if coin_fee > 0 {
                    self.coin_fees.set(&coin_id, 0);

                    Some(
                        self.process_external_transfer(&coin_id, &sender_id, coin_fee)
                            .unwrap(),
                    )
                } else {
                    None
                }
            })
            .collect::<Vec<TransferableAsset>>();

        self.send_request(sender_id, fee_aggregator.clone(), assets);
    }

    pub fn calculate_coin_transfer_fee(
        &self,
        amount: u128,
        coin: &Asset<WrappedNativeCoin>,
    ) -> Result<u128, String> {
        let mut fee = (amount * coin.metadata().fee_numerator()) / FEE_DENOMINATOR;
        fee.add(coin.metadata().fixed_fee()).map(|fee| *fee)
    }
}
