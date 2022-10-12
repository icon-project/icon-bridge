use std::collections::HashMap;

use crate::types::AssetId;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::{self, LookupMap};
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
pub struct Balances(HashMap<(AccountId, AssetId), AccountBalance>);

impl Balances {
    pub fn new() -> Self {
        Self(HashMap::new())
    }

    pub fn add(&mut self, account: &AccountId, asset_id: &AssetId) {
        if !self.contains(account, asset_id) {
            self.0.insert(
                (account.to_owned(), asset_id.to_owned()),
                AccountBalance::default(),
            );
        }
    }

    pub fn remove(&mut self, account: &AccountId, asset_id: &AssetId) {
        self.0.remove(&(account.to_owned(), asset_id.to_owned()));
    }

    pub fn get(&self, account: &AccountId, asset_id: &AssetId) -> Option<AccountBalance> {
        if let Some(balance) = self.0.get(&(account.to_owned(), asset_id.to_owned())) {
            return Some(balance.to_owned());
        }
        None
    }

    pub fn contains(&self, account: &AccountId, asset_id: &AssetId) -> bool {
        return self
            .0
            .contains_key(&(account.to_owned(), asset_id.to_owned()));
    }

    pub fn set(
        &mut self,
        account: &AccountId,
        asset_id: &AssetId,
        account_balance: AccountBalance,
    ) {
        self.0
            .insert((account.to_owned(), asset_id.to_owned()), account_balance);
    }

    pub fn to_vec(&self) -> Vec<((AccountId, AssetId), AccountBalance)> {
        if !self.0.is_empty() {
            return self
                .0
                .clone()
                .into_iter()
                .map(|((accound_id, asset_id), account_balance)| {
                    (
                        (accound_id, asset_id),
                        AccountBalance {
                            deposit: account_balance.deposit(),
                            refundable: account_balance.refundable(),
                            locked: account_balance.locked(),
                        },
                    )
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
    use near_sdk::{testing_env, VMContext};
    use std::convert::TryInto;
    use std::vec;

    fn get_context(input: Vec<u8>, is_view: bool) -> VMContext {
        VMContext {
            current_account_id: "alice.testnet".to_string(),
            signer_account_id: "robert.testnet".to_string(),
            signer_account_pk: vec![0, 1, 2],
            predecessor_account_id: "jane.testnet".to_string(),
            input,
            block_index: 0,
            block_timestamp: 0,
            account_balance: 0,
            account_locked_balance: 0,
            storage_usage: 0,
            attached_deposit: 0,
            prepaid_gas: 10u64.pow(18),
            random_seed: vec![0, 1, 2],
            is_view,
            output_data_receivers: vec![],
            epoch_height: 19,
        }
    }

    #[test]
    fn add_balance() {
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
        let context = get_context(vec![], false);
        testing_env!(context);
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
