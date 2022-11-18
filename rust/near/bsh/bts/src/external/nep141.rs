use super::*;

#[ext_contract(ext_nep141)]
pub trait Nep141Service {
    fn new(owner_id: AccountId, total_supply: U128, metadata: FungibleTokenMetadata) -> Self;
    fn ft_transfer(&mut self, receiver_id: AccountId, amount: U128, memo: Option<String>);
    fn mint(&mut self, amount: U128, receiver_id: AccountId) -> U128;
    fn burn(&mut self, amount: U128);
}
