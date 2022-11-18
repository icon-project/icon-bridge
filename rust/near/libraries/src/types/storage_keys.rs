use super::*;

#[derive(BorshSerialize, BorshStorageKey)]
pub enum KeyType {
    Key,
    Value,
}

#[derive(BorshSerialize, BorshStorageKey)]
pub enum BmcEventType {
    Message,
    Error,
}

#[derive(BorshSerialize, BorshStorageKey)]
pub enum StorageKey {
    Owners,
    Services,
    Links(KeyType),
    Routes(KeyType),
    Connections,
    BmcEvent(BmcEventType),
    Assets(KeyType),
    Balances(KeyType),
    StorageBalances(KeyType),
    AssetFees,
    Requests,
    BlacklistedAccounts,
}
