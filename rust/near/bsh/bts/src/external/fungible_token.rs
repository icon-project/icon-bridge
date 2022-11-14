use near_sdk::json_types::U128;
use near_sdk::{ext_contract, AccountId};

#[ext_contract(ext_ft)]
pub trait FungibleToken {
    fn ft_transfer(receiver_id: AccountId, amount: U128, memo: Option<String>);
}
