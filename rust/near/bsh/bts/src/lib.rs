use btp_common::btp_address::Address;
use btp_common::errors::BshError;
use libraries::types::{
    Account, AccountBalance, AccumulatedAssetFees, AssetId, BTPAddress, CoinIds, TransferableAsset,
    WrappedNativeCoin,
};
use libraries::{
    types::messages::BlackListType,
    types::messages::SerializedMessage,
    types::messages::TokenServiceMessage,
    types::messages::TokenServiceType,
    types::messages::{BtpMessage, ErrorMessage},
    types::Asset,
    types::AssetFees,
    types::AssetMetadata,
    types::Assets,
    types::Balances,
    types::BlackListedAccounts,
    types::Math,
    types::Network,
    types::Owners,
    types::Request,
    types::Requests,
    types::StorageBalances,
    types::TokenLimits,
    types::WrappedI128,
};

use std::str::FromStr;

use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::LazyOption;
use near_sdk::serde_json::json;
use near_sdk::serde_json::{to_value, Value};
use near_sdk::PromiseOrValue;
use near_sdk::{assert_one_yocto, AccountId};
use near_sdk::{
    env,
    json_types::{Base64VecU8, U128},
    log, near_bindgen, require, Gas, PanicOnDefault, Promise, PromiseResult,
};
use std::collections::HashSet;
use std::convert::TryFrom;
use std::convert::TryInto;
mod external;
use external::*;
mod accounting;
mod assertion;
mod blacklist_management;
mod coin_management;
mod estimate;
mod fee_management;
mod messaging;
mod owner_management;
mod transfer;
mod types;
mod util;
pub use types::RegisteredCoins;
pub type CoinFees = AssetFees;
pub type CoinId = AssetId;
pub type Coin = Asset<WrappedNativeCoin>;
pub type Coins = Assets<WrappedNativeCoin>;

pub static NEP141_CONTRACT: &'static [u8] = include_bytes!("../res/NEP141_CONTRACT.wasm");
pub static FEE_DENOMINATOR: u128 = 10_u128.pow(4);

pub static RC_ERROR: u8 = 1;
pub static RC_OK: u8 = 0;

#[near_bindgen]
#[derive(BorshDeserialize, BorshSerialize, PanicOnDefault)]
pub struct BtpTokenService {
    native_coin_name: String,
    network: Network,
    owners: Owners,
    coins: Coins,
    balances: Balances,
    storage_balances: StorageBalances,
    coin_fees: CoinFees,
    requests: Requests,
    serial_no: i128,
    bmc: AccountId,
    name: String,
    blacklisted_accounts: BlackListedAccounts,
    token_limits: TokenLimits,
    coin_ids: CoinIds,
    registered_coins: RegisteredCoins,

    #[cfg(feature = "testable")]
    pub message: LazyOption<Base64VecU8>,
}

#[near_bindgen]
impl BtpTokenService {
    #[init]
    pub fn new(service_name: String, bmc: AccountId, network: String, native_coin: Coin) -> Self {
        require!(!env::state_exists(), "Already initialized");
        let mut owners = Owners::new();
        owners.add(&env::current_account_id());

        let mut coins = <Coins>::new();
        let mut balances = Balances::new();
        let native_coin_id = Self::hash_coin_id(native_coin.name());

        balances.add(&env::current_account_id(), &native_coin_id);
        coins.add(&native_coin_id, &native_coin);
        let blacklisted_accounts = BlackListedAccounts::new();
        let mut coin_fees = CoinFees::new();
        coin_fees.add(&native_coin_id);
        let mut coin_ids = CoinIds::new();
        coin_ids.add(native_coin.name(), native_coin_id);
        Self {
            native_coin_name: native_coin.name().to_owned(),
            network,
            owners,
            coins,
            balances,
            storage_balances: StorageBalances::new(),
            coin_fees,
            serial_no: Default::default(),
            requests: Requests::new(),
            bmc,
            name: service_name,
            blacklisted_accounts,
            registered_coins: RegisteredCoins::new(),
            token_limits: TokenLimits::new(),
            coin_ids,

            #[cfg(feature = "testable")]
            message: LazyOption::new(b"message".to_vec(), None),
        }
    }

    fn bmc(&self) -> &AccountId {
        &self.bmc
    }

    fn name(&self) -> &String {
        &self.name
    }

    fn requests(&self) -> &Requests {
        &self.requests
    }

    fn requests_mut(&mut self) -> &mut Requests {
        &mut self.requests
    }

    fn process_deposit(&mut self, amount: u128, balance: &mut AccountBalance) {
        balance.deposit_mut().add(amount).unwrap();
    }

    fn calculate_storage_cost(&self, initial_storage_usage: u64) -> U128 {
        let total_storage_usage = env::storage_usage() - initial_storage_usage;
        let storage_cost =
            total_storage_usage as u128 * env::storage_byte_cost() + 669547687500000000;
        U128(storage_cost)
    }
}

impl BtpTokenService {
    pub fn serial_no(&self) -> i128 {
        self.serial_no
    }
}
