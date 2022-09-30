use super::*;

#[near_bindgen]
impl BtpTokenService {
    #[init(ignore_state)]
    #[private]
    pub fn migrate() -> Self {
        let old: BtpTokenService = env::state_read().expect("failed to read the state");

        Self {
            owners: old.owners,
            native_coin_name: old.native_coin_name,
            network: old.network,
            coins: old.coins.into(),
            balances: old.balances.into(),
            storage_balances: old.storage_balances,
            coin_fees: old.coin_fees.into(),
            serial_no: old.serial_no,
            requests: old.requests,
            bmc: old.bmc,
            name: old.name,
            registered_coins: old.registered_coins,
            coin_ids: old.coin_ids.into(),
            blacklisted_accounts: BlackListedAccounts::new(),
            token_limits: TokenLimits::new(),

            #[cfg(feature = "testable")]
            message: old.message,
        }
    }
}
