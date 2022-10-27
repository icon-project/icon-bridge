use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::AccountId;
use std::collections::HashSet;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Owners(HashSet<AccountId>);

impl Owners {
    pub fn new() -> Self {
        Self(HashSet::new())
    }

    pub fn len(&self) -> usize {
        self.0.len()
    }

    pub fn add(&mut self, address: &AccountId) {
        self.0.insert(address.to_owned());
    }

    pub fn remove(&mut self, address: &AccountId) {
        self.0.remove(&address);
    }

    pub fn contains(&self, address: &AccountId) -> bool {
        self.0.contains(&address)
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
    use std::collections::HashSet;

    #[test]
    fn add_owner() {

        let mut owners = Owners::new();
        owners.add(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        );
        let result = owners.contains(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        );
        assert_eq!(result, true);
    }

    #[test]
    fn add_existing_owner() {

        let mut owners = Owners::new();
        let owner_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let owner_1_duplicate = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();

        owners.add(&owner_1);
        owners.add(&owner_1_duplicate);
        let result = owners.to_vec();
        let expected: Vec<AccountId> = vec![
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        ];
        assert_eq!(result, expected);
    }

    #[test]
    fn remove_owner() {

        let mut owners = Owners::new();
        let owner = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        owners.add(&owner);
        owners.remove(&owner);
        let result = owners.contains(&owner);
        assert_eq!(result, false);
    }

    #[test]
    fn remove_owner_non_existing() {

        let mut owners = Owners::new();
        let owner_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e8"
            .parse::<AccountId>()
            .unwrap();
        let owner_2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        owners.add(&owner_1);
        owners.remove(&owner_2);
        let result = owners.contains(&owner_2);
        assert_eq!(result, false);
    }

    #[test]
    fn to_vec_owners() {

        let mut owners = Owners::new();
        let owner_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let owner_2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let owner_3 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        owners.add(&owner_1);
        owners.add(&owner_2);
        owners.add(&owner_3);
        let owners = owners.to_vec();
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
        let result: HashSet<_> = owners.iter().collect();
        let expected: HashSet<_> = expected_owners.iter().collect();
        assert_eq!(result, expected);
    }
}
