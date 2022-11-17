use super::*;

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize, Clone)]
pub struct BlackList {
    address: String,
    network: Network,
}

#[derive(BorshDeserialize, BorshSerialize)]
pub struct BlackListedAccounts(UnorderedSet<(String, Network)>);

impl BlackListedAccounts {
    pub fn new() -> Self {
        Self(UnorderedSet::new(b"blacklist".to_vec()))
    }

    pub fn len(&self) -> usize {
        self.0.len() as usize
    }

    pub fn add(&mut self, address: &str, network: &str) {
        self.0.insert(&(address.to_owned(), network.to_owned()));
    }

    pub fn remove(&mut self, address: &str, network: &str) {
        self.0.remove(&(address.to_owned(), network.to_owned()));
    }

    pub fn contains(&self, address: &str, network: &str) -> bool {
        self.0.contains(&(address.to_owned(), network.to_owned()))
    }

    pub fn to_vec(&self) -> Vec<BlackList> {
        if !self.0.is_empty() {
            return self
                .0
                .iter()
                .map(|(address, network)| BlackList { address, network })
                .collect();
        }
        vec![]
    }

    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashSet;

    #[test]
    fn add_user() {
        let mut blacklisted_user = BlackListedAccounts::new();
        let network = "0x1.icon".to_owned();
        blacklisted_user.add(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4",
                &network
        );
        let result = blacklisted_user.contains(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4",
                &network
        );
        assert_eq!(result, true);
    }

    #[test]
    fn add_already_blacklisted_user() {
        let mut blacklisted_user = BlackListedAccounts::new();
        let network = "0x1.icon".to_owned();
        let user1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4";
        let owner_1_duplicate = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4";

        blacklisted_user.add(&user1, &network);
        blacklisted_user.add(&owner_1_duplicate, &network);

        let result = blacklisted_user.to_vec();
        let expected: Vec<BlackList> = vec![BlackList{
            address: "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_owned(),
            network
        }];

        assert_eq!(result, expected);
    }

    #[test]
    fn remove_user() {
        let mut blacklisted_user = BlackListedAccounts::new();
        let network = "0x1.icon".to_owned();
        let user = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4";

        blacklisted_user.add(&user, &network);
        blacklisted_user.remove(&user, &network);

        let result = blacklisted_user.contains(&user, &network);
        assert_eq!(result, false);
    }

    #[test]
    fn remove_user_non_existing() {
        let mut blacklisted_user = BlackListedAccounts::new();
        let network = "0x1.icon".to_owned();
        let user1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e8"
            .parse::<AccountId>()
            .unwrap();
        let user2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        blacklisted_user.add(&user1, &network);
        blacklisted_user.remove(&user2, &network);
        let result = blacklisted_user.contains(&user2, &network);
        assert_eq!(result, false);
    }

    #[test]
    fn to_vec_users() {
        let mut blacklisted_user = BlackListedAccounts::new();
        let network = "0x1.icon".to_owned();
        let user1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4";
        let user2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4";
        let user3 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4";

        blacklisted_user.add(&user1, &network);
        blacklisted_user.add(&user2, &network);
        blacklisted_user.add(&user3, &network);

        let blacklisted_user = blacklisted_user.to_vec();
        let expected_owners: Vec<BlackList> = vec![
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
