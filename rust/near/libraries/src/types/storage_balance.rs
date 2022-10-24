use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::AccountId;
use std::collections::HashMap;

type StorageBalance = u128;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct StorageBalances(HashMap<AccountId, StorageBalance>);

impl StorageBalances {
    pub fn new() -> Self {
        Self(HashMap::new())
    }

    pub fn add(&mut self, account: AccountId) {
        if !self.contains(&account) {
            self.0.insert(account, u128::default());
        }
    }

    pub fn remove(&mut self, account: &AccountId) {
        self.0.remove(account);
    }

    pub fn get(&self, account: &AccountId) -> Option<&StorageBalance> {
        if let Some(storage_balance) = self.0.get(account) {
            return Some(storage_balance);
        }
        None
    }

    pub fn get_mut(&mut self, account: &AccountId) -> Option<&mut StorageBalance> {
        if let Some(storage_balance) = self.0.get_mut(account) {
            return Some(storage_balance);
        }
        None
    }

    pub fn contains(&self, account: &AccountId) -> bool {
        return self.0.contains_key(account)
    }

    pub fn set(&mut self, account: AccountId, storage_balance: StorageBalance) {
        self.0.insert(account, storage_balance);
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::types::Math;


    #[test]
    fn add_storage_balance() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        storage_balances.add(account.clone());
        let result = storage_balances.get(&account).unwrap();
        assert_eq!(*result, 0);
    }

    #[test]
    fn add_storage_balance_existing() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        storage_balances.add(account.clone());
        let storage_balance = storage_balances.get_mut(&account).unwrap();
        storage_balance.add(100).unwrap();
        let result = storage_balance.clone();
        storage_balances.add(account.clone());
        assert_eq!(result, 100);
    }

    #[test]
    fn remove_storage_balance() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        storage_balances.add(account.clone());
        storage_balances.remove(&account);
        let result = storage_balances.get(&account);
        assert_eq!(result, None);
    }

    #[test]
    fn remove_storage_balance_non_existing() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        storage_balances.remove(&account);
        let result = storage_balances.get(&account);
        assert_eq!(result, None);
    }
}
