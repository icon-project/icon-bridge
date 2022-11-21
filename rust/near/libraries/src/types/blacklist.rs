use super::*;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct BlacklistedAccounts(UnorderedSet<AccountId>);

impl BlacklistedAccounts {
    pub fn new() -> Self {
        Self(UnorderedSet::new(StorageKey::BlacklistedAccounts))
    }

    pub fn len(&self) -> usize {
        self.0.len() as usize
    }

    pub fn add(&mut self, account_id: &AccountId) {
        self.0.insert(account_id);
    }

    pub fn remove(&mut self, account_id: &AccountId) {
        self.0.remove(account_id);
    }

    pub fn contains(&self, account_id: &AccountId) -> bool {
        self.0.contains(account_id)
    }

    pub fn to_vec(&self) -> Vec<AccountId> {
        if !self.0.is_empty() {
            return self.0.to_vec();
        }
        vec![]
    }

    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }
}

impl Default for BlacklistedAccounts {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashSet;

    #[test]
    fn add_user() {
        let mut blacklisted_user = BlacklistedAccounts::new();
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
        let mut blacklisted_user = BlacklistedAccounts::new();
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
        let mut blacklisted_user = BlacklistedAccounts::new();
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
        let mut blacklisted_user = BlacklistedAccounts::new();
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
        let mut blacklisted_user = BlacklistedAccounts::new();
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
