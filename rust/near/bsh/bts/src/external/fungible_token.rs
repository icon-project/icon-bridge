use super::*;
#[ext_contract(ext_ft)]
pub trait FungibleToken {
    fn ft_transfer(receiver_id: AccountId, amount: U128, memo: Option<String>);
}
