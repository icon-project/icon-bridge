use super::*;

#[near_bindgen]
impl BtpTokenService {

    #[init(ignore_state)]
    #[private]
    pub fn migrate() -> Self {
        let old: BtpTokenServiceV0_9 = env::state_read().expect("failed to read the state");
        let mut coin_ids = CoinIds::new();
        old.coins.to_vec().iter().for_each(|coin| {
            coin_ids.add(&coin.name, &Self::hash_coin_id(&coin.name))
        });
        
        Self {
            owners: old.owners,
            native_coin_name: old.native_coin_name,
            network: old.network,
            coins: old.coins,
            balances: old.balances,
            storage_balances: old.storage_balances,
            coin_fees: old.coin_fees,
            serial_no: old.serial_no,
            requests: old.requests,
            bmc: old.bmc,
            name: old.name,
            registered_coins: old.registered_coins,
            coin_ids: coin_ids,
            blacklisted_accounts: BlackListedAccounts::new(),
            token_limits: TokenLimits::new(),

            #[cfg(feature = "testable")]
            message: old.message
        }
    } 
}