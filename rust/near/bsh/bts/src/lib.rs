mod accounting;
mod assertion;
mod blacklist_management;
mod estimate;
mod external;
mod fee_management;
mod messaging;
mod owner_management;
mod token_management;
mod transfer;
mod types;
mod util;

use btp_common::{btp_address::Address, errors::BshError};
use libraries::types::{
    messages::{
        BlackListType, BtpMessage, ErrorMessage, SerializedMessage, TokenServiceMessage,
        TokenServiceType,
    },
    AccountBalance, AccumulatedAssetFees, Asset, AssetFees, AssetId, AssetMetadata, Assets,
    BTPAddress, Balances, BlackListedAccounts, FungibleToken, Math, Network, Owners, Request,
    Requests, StorageBalances, TokenIds, TokenLimit, TokenLimits, TransferableAsset, WrappedI128,
    WrappedNativeCoin,
};

use std::str::FromStr;

use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};

#[cfg(feature = "testable")]
use near_sdk::{collections::LazyOption, json_types::Base64VecU8};

use near_sdk::{
    env,
    json_types::U128,
    log, near_bindgen, require,
    serde_json::{json, to_value, Value},
    AccountId, Balance, Gas, PanicOnDefault, Promise, PromiseOrValue, PromiseResult,
};

use std::convert::{TryFrom, TryInto};

use external::*;

pub use types::RegisteredTokens;
pub type TokenFees = AssetFees;
pub type TokenId = AssetId;
pub type Token = Asset<WrappedNativeCoin>;
pub type Tokens = Assets<WrappedNativeCoin>;

pub static NEP141_CONTRACT: &[u8] = include_bytes!("../res/NEP141_CONTRACT.wasm");
pub static FEE_DENOMINATOR: u128 = 10_u128.pow(4);

pub static RC_ERROR: u8 = 1;
pub static RC_OK: u8 = 0;

#[near_bindgen]
#[derive(BorshDeserialize, BorshSerialize, PanicOnDefault)]
pub struct BtpTokenService {
    native_coin_name: String,
    network: Network,
    owners: Owners,
    tokens: Tokens,
    balances: Balances,
    storage_balances: StorageBalances,
    token_fees: TokenFees,
    requests: Requests,
    serial_no: i128,
    bmc: AccountId,
    name: String,
    blacklisted_accounts: BlackListedAccounts,
    token_limits: TokenLimits,
    token_ids: TokenIds,
    registered_tokens: RegisteredTokens,

    #[cfg(feature = "testable")]
    pub message: LazyOption<Base64VecU8>,
}

#[near_bindgen]
impl BtpTokenService {
    #[init]
    pub fn new(service_name: String, bmc: AccountId, network: String, native_coin: Token) -> Self {
        require!(!env::state_exists(), "Already initialized");
        let mut owners = Owners::new();
        owners.add(&env::current_account_id());

        let mut tokens = <Tokens>::new();
        let mut balances = Balances::new();
        let native_coin_id = Self::hash_token_id(native_coin.name());

        balances.add(&env::current_account_id(), &native_coin_id);
        tokens.add(&native_coin_id, &native_coin);
        let blacklisted_accounts = BlackListedAccounts::new();
        let mut token_fees = TokenFees::new();
        token_fees.add(&native_coin_id);
        let mut token_ids = TokenIds::new();
        token_ids.add(native_coin.name(), native_coin_id);
        Self {
            native_coin_name: native_coin.name().to_owned(),
            network,
            owners,
            tokens,
            balances,
            storage_balances: StorageBalances::new(),
            token_fees,
            serial_no: Default::default(),
            requests: Requests::new(),
            bmc,
            name: service_name,
            blacklisted_accounts,
            registered_tokens: RegisteredTokens::new(),
            token_limits: TokenLimits::new(),
            token_ids,

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
