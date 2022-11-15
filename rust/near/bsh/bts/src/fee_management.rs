use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Fee Management  * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    /// Returns the total accumulated fees
    
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
                    accumulated_fees: coin_fee,
                }
            })
            .collect()
    }

    /// handling the fee gather method
    /// 
    /// # Arguments
    /// * `fee_aggregator` - should be in the form of btp://0x1.near/account.testnet
    /// * `service` - Name of the service should be given in string format
    /// 
    pub fn handle_fee_gathering(&mut self, fee_aggregator: BTPAddress, service: String) {
        self.assert_predecessor_is_bmc();
        self.assert_valid_service(&service);
        self.transfer_fees(&fee_aggregator);
    }


    /// setting the fee ratio in btp
    /// caller should be owner
    /// 
    /// # Arguments
    /// * `coin_name` - name of the coin shpuld be given in the string format
    /// * `fee_numerator` - unsigned number should be given
    /// * `fixed_fee` - unsigned number should be given
    /// 
    pub fn set_fee_ratio(&mut self, coin_name: String, fee_numerator: U128, fixed_fee: U128) {
        self.assert_have_permission();
        self.assert_valid_fee_ratio(fee_numerator.into(), fixed_fee.into());

        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        // Getting the coin by its id
        let mut coin = self.coins.get(&coin_id).unwrap();
        coin.metadata_mut()
            .fee_numerator_mut()
            .clone_from(&fee_numerator.into());
        coin.metadata_mut()
            .fixed_fee_mut()
            .clone_from(&fixed_fee.into());

        self.coins.set(&coin_id, &coin)
    }


    /// The method get fee is created
    /// # Arguments 
    /// * `coin_name` - name of the coin should be given in string format
    /// * `amount` - should be an unsigned number
    /// 
    pub fn get_fee(&self, coin_name: String, amount: U128) -> U128 {
        let coin_id = self
            .coin_id(&coin_name)
            .map_err(|err| format!("{}", err))
            .unwrap();
        let coin = self.coins.get(&coin_id).unwrap();

        self.calculate_coin_transfer_fee(u128::from(amount), &coin)
            .map(|e| U128(e))
            .unwrap()
    }
}

impl BtpTokenService {
    /// method transfer fees got created in btp
    /// caller should be a owner
    /// # Arguments
    /// * `fee_aggregator` - btp address should be given that should be in the form of btp://0x1.near/account.testnet
    /// 
    pub fn transfer_fees(&mut self, fee_aggregator: &BTPAddress) {
        let sender_id = env::current_account_id();
        let assets = self
            .coins
            .to_vec()
            .iter()
            .filter_map(|coin| {
                let coin_id = self.coin_ids.get(&coin.name).unwrap().clone();
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

    /// coin transfer fee is calculated
    /// # Arguments
    /// * `amount` - should be a unsigned number
    /// * `coin` - wrapped native coin should be given
    /// 
    pub fn calculate_coin_transfer_fee(
        &self,
        amount: u128,
        coin: &Asset<WrappedNativeCoin>,
    ) -> Result<u128, String> {
        let mut fee = (amount * coin.metadata().fee_numerator()) / FEE_DENOMINATOR;
        fee.add(coin.metadata().fixed_fee()).map(|fee| *fee)
    }
}
