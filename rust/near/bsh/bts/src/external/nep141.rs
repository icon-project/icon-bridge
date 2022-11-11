use near_contract_standards::fungible_token::metadata::FungibleTokenMetadata;
use near_sdk::json_types::U128;
use near_sdk::{ext_contract, AccountId, Balance};

#[ext_contract(ext_nep141)]
pub trait Nep141Service {
    fn new(owner_id: AccountId, total_supply: U128, metadata: FungibleTokenMetadata) -> Self;
    fn ft_transfer(&mut self, receiver_id: AccountId, amount: U128, memo: Option<String>);
    fn mint(&mut self, amount: U128, receiver_id: AccountId) -> U128;
    fn burn(&mut self, amount: U128);
}
