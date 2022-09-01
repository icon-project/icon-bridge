use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::AccountId;
use std::collections::HashSet;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct BlackListedAccounts(HashSet<AccountId>);

impl BlackListedAccounts {
    pub fn new() -> Self {
        Self(HashSet::new())
    }

    pub fn len(&self) -> usize {
        self.0.len()
    }

    pub fn add(&mut self, user: &AccountId) {
        self.0.insert(user.to_owned());
    }

    pub fn remove(&mut self, user: &AccountId) {
        self.0.remove(&user);
    }

    pub fn contains(&self, user: &AccountId) -> bool {
        self.0.contains(&user)
    }

    pub fn to_vec(&self) -> Vec<AccountId> {
        if !self.0.is_empty() {
            return self.0.clone().into_iter().collect::<Vec<AccountId>>();
        }
        vec![]
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use near_sdk::{testing_env, VMContext};
    use std::collections::HashSet;

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
    fn add_user() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let mut blacklisted_user = BlackListedAccounts::new();
        blacklisted_user.add(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        );
        let result = blacklisted_user.contains(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        );
        assert_eq!(result, true);
    }

    #[test]
    fn add_already_blacklisted_user() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let mut blacklisted_user = BlackListedAccounts::new();
        let user1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let owner_1_duplicate = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        blacklisted_user.add(&user1);
        blacklisted_user.add(&owner_1_duplicate);
        let result = blacklisted_user.to_vec();
        let expected: Vec<AccountId> = vec![
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        ];
        assert_eq!(result, expected);
    }

    #[test]
    fn remove_user() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let mut blacklisted_user = BlackListedAccounts::new();
        let user = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        blacklisted_user.add(&user);
        blacklisted_user.remove(&user);
        let result = blacklisted_user.contains(&user);
        assert_eq!(result, false);
    }

    #[test]
    fn remove_user_non_existing() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let mut blacklisted_user = BlackListedAccounts::new();
        let user1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e8"
            .parse::<AccountId>()
            .unwrap();
        let user2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        blacklisted_user.add(&user1);
        blacklisted_user.remove(&user2);
        let result = blacklisted_user.contains(&user2);
        assert_eq!(result, false);
    }

    #[test]
    fn to_vec_users() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let mut blacklisted_user = BlackListedAccounts::new();
        let user1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let user2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let user3 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        blacklisted_user.add(&user1);
        blacklisted_user.add(&user2);
        blacklisted_user.add(&user3);
        let blacklisted_user = blacklisted_user.to_vec();
        let expected_owners: Vec<AccountId> = vec![
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
            "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
            "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        ];
        let result: HashSet<_> = blacklisted_user.iter().collect();
        let expected: HashSet<_> = expected_owners.iter().collect();
        assert_eq!(result, expected);
    }
}
