use near_contract_standards::fungible_token::metadata::{
    FungibleTokenMetadata, FungibleTokenMetadataProvider,
};
use near_contract_standards::fungible_token::FungibleToken;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::LazyOption;
use near_sdk::json_types::U128;
use near_sdk::{
    env, log, near_bindgen, require, AccountId, Balance, BorshStorageKey, PanicOnDefault,
    PromiseOrValue,
};

#[near_bindgen]
#[derive(BorshDeserialize, BorshSerialize, PanicOnDefault)]
pub struct Contract {
    token: FungibleToken,
    metadata: LazyOption<FungibleTokenMetadata>,
    owner: AccountId,
}

#[derive(BorshSerialize, BorshStorageKey)]
enum StorageKey {
    FungibleToken,
    Metadata,
}
#[near_bindgen]
impl Contract {
    /// Initializes the contract with the given total supply owned by the given `owner_id` with
    /// the given fungible token metadata.
    #[init]
    pub fn new(owner_id: AccountId, total_supply: U128, metadata: FungibleTokenMetadata) -> Self {
        require!(!env::state_exists(), "Already initialized");
        metadata.assert_valid();
        let mut this = Self {
            token: FungibleToken::new(StorageKey::FungibleToken),
            metadata: LazyOption::new(StorageKey::Metadata, Some(&metadata)),
            owner: owner_id.clone(),
        };
        this.token.internal_register_account(&owner_id);
        this.token.internal_deposit(&owner_id, total_supply.into());

        this
    }

    pub fn mint(&mut self, amount: U128, receiver_id: AccountId) -> U128 {
        require!(env::predecessor_account_id() == self.owner);
        self.token
            .internal_deposit(&env::predecessor_account_id(), amount.into());

        match self.storage_balance_of(receiver_id.clone()) {
            Some(_) => U128::from(0),
            None => {
                let inital_storage_used = env::storage_usage();
                self.token.internal_register_account(&receiver_id);
                self.get_storage_cost(inital_storage_used)
            }
        }
    }

    pub fn burn(&mut self, amount: U128) {
        require!(env::predecessor_account_id() == self.owner);
        self.token
            .internal_withdraw(&env::predecessor_account_id(), amount.into())
    }

    fn on_account_closed(&mut self, account_id: AccountId, balance: Balance) {
        log!("Closed @{} with {}", account_id, balance);
    }

    fn on_tokens_burned(&mut self, account_id: AccountId, amount: Balance) {
        log!("Account @{} burned {}", account_id, amount);
    }

    fn get_storage_cost(&self, initial_storage_usage: u64) -> U128 {
        let total_storage_usage = env::storage_usage() - initial_storage_usage;
        let storage_cost =
            total_storage_usage as u128 * env::storage_byte_cost() + 669547687500000000;
        U128(storage_cost)
    }
}

near_contract_standards::impl_fungible_token_core!(Contract, token, on_tokens_burned);
near_contract_standards::impl_fungible_token_storage!(Contract, token, on_account_closed);

#[near_bindgen]
impl FungibleTokenMetadataProvider for Contract {
    fn ft_metadata(&self) -> FungibleTokenMetadata {
        self.metadata.get().unwrap()
    }
}
