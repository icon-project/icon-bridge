use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::UnorderedSet;
use near_sdk::AccountId;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct BlackListedAccounts(UnorderedSet<AccountId>);

impl BlackListedAccounts {
    pub fn new() -> Self {
        Self(UnorderedSet::new(b"blacklist".to_vec()))
    }

    pub fn len(&self) -> usize {
        self.0.len() as usize
    }

    pub fn add(&mut self, user: &AccountId) {
        self.0.insert(&user.to_owned());
    }

    pub fn remove(&mut self, user: &AccountId) {
        self.0.remove(&user);
    }

    pub fn contains(&self, user: &AccountId) -> bool {
        self.0.contains(&user)
    }

    pub fn to_vec(&self) -> Vec<AccountId> {
        if !self.0.is_empty() {
            return self.0.to_vec();
        }
        vec![]
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashSet;

    #[test]
    fn add_user() {
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
