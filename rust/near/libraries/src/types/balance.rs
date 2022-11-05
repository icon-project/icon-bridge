use std::collections::HashMap;

use crate::types::AssetId;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::{self, LookupMap, UnorderedSet};
use near_sdk::serde::{Deserialize, Serialize};
use near_sdk::AccountId;
use near_sdk::Balance;

#[derive(
    Debug, Default, BorshDeserialize, BorshSerialize, PartialEq, Eq, Clone, Serialize, Deserialize,
)]
#[serde(crate = "near_sdk::serde")]
pub struct AccountBalance {
    deposit: Balance,
    refundable: Balance,
    locked: Balance,
}

impl AccountBalance {
    pub fn deposit(&self) -> Balance {
        self.deposit
    }

    pub fn locked(&self) -> Balance {
        self.locked
    }

    pub fn refundable(&self) -> Balance {
        self.refundable
    }

    pub fn deposit_mut(&mut self) -> &mut Balance {
        &mut self.deposit
    }

    pub fn locked_mut(&mut self) -> &mut Balance {
        &mut self.locked
    }

    pub fn refundable_mut(&mut self) -> &mut Balance {
        &mut self.refundable
    }
}

#[derive(BorshDeserialize, BorshSerialize)]

pub struct Balances {
    keys: UnorderedSet<(AccountId, AssetId)>,
    values: LookupMap<(AccountId, AssetId), AccountBalance>,
}

impl Balances {
    pub fn new() -> Self {
        Self {
            keys: UnorderedSet::new(b"balances_keys".to_vec()),
            values: LookupMap::new(b"balance_values".to_vec()),
        }
    }

    pub fn add(&mut self, account: &AccountId, asset_id: &AssetId) {
        if !self.contains(account, asset_id) {
            self.keys.insert(&(account.to_owned(), asset_id.to_owned()));
            self.values.insert(
                &(account.to_owned(), asset_id.to_owned()),
                &AccountBalance::default(),
            );
        }
    }

    pub fn remove(&mut self, account: &AccountId, asset_id: &AssetId) {
        self.keys.remove(&(account.to_owned(), asset_id.to_owned()));
        self.values
            .remove(&(account.to_owned(), asset_id.to_owned()));
    }

    pub fn get(&self, account: &AccountId, asset_id: &AssetId) -> Option<AccountBalance> {
        if let Some(balance) = self.values.get(&(account.to_owned(), asset_id.to_owned())) {
            return Some(balance.to_owned());
        }
        None
    }

    pub fn contains(&self, account: &AccountId, asset_id: &AssetId) -> bool {
        return self
            .keys
            .contains(&(account.to_owned(), asset_id.to_owned()));
    }

    pub fn set(
        &mut self,
        account: &AccountId,
        asset_id: &AssetId,
        account_balance: AccountBalance,
    ) {
        self.values
            .insert(&(account.to_owned(), asset_id.to_owned()), &account_balance);
    }

    pub fn to_vec(&self) -> Vec<((AccountId, AssetId), AccountBalance)> {
        if !self.keys.is_empty() {
            return self
                .keys
                .to_vec()
                .into_iter()
                .map(|(account_id, asset_id)| {
                    let account_balance = self.values.get(&(account_id.clone(), asset_id)).unwrap();

                    ((account_id, asset_id), account_balance)
                })
                .collect();
        }

        vec![]
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::types::Math;
    use near_sdk::{env, AccountId};
    use std::convert::TryInto;
    use std::vec;

    #[test]
    fn add_balance() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        let result = balances.contains(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, true);
    }

    #[test]
    fn add_balance_exisitng() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        account_balance.deposit_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        assert_eq!(result.deposit(), 1000);
    }

    #[test]
    fn remove_balance() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        balances.remove(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let result = balances.contains(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, false);
    }

    #[test]
    fn remove_balance_non_exisitng() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        balances.remove(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let result = balances.contains(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, false);
    }

    #[test]
    fn set_balance() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        account_balance.deposit_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        assert_eq!(result.deposit(), 1000);
    }

    #[test]
    fn deposit_add() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        account_balance.deposit_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        assert_eq!(result.deposit(), 1000);
    }

    #[test]
    fn deposit_add_overflow() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();
        account_balance.deposit_mut().add(u128::MAX).unwrap();

        assert_eq!(
            account_balance.deposit_mut().add(1),
            Err("overflow occured".to_string())
        )
    }

    #[test]
    fn deposit_sub_underflow() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(
            account_balance.deposit_mut().sub(1),
            Err("underflow occured".to_string())
        )
    }

    #[test]
    fn locked_balance_add() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        account_balance.locked_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(result.locked(), 1000);
    }

    #[test]
    fn locked_balance_sub() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        account_balance.locked_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(result.locked(), 1000);

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        account_balance.locked_mut().sub(1).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(result.locked(), 999);
    }

    #[test]
    fn refundable_balance_add() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        account_balance.refundable_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(result.refundable(), 1000);
    }

    #[test]
    fn refundable_balance_sub() {
        let mut balances = Balances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        balances.add(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        account_balance.refundable_mut().add(1000).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(result.refundable(), 1000);

        let mut account_balance = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        account_balance.refundable_mut().sub(1).unwrap();

        balances.set(
            &account,
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            account_balance,
        );

        let result = balances
            .get(
                &account,
                &env::sha256("ABC Asset".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            )
            .unwrap();

        assert_eq!(result.refundable(), 999);
    }
}
