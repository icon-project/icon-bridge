use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::UnorderedMap;
use serde::{Deserialize, Serialize};

#[derive(BorshDeserialize, BorshSerialize, Debug)]
pub struct TokenLimits(UnorderedMap<String, u128>);

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize, Clone)]
pub struct TokenLimit {
    coin_name: String,
    token_limit: u128,
}

impl TokenLimit {
    pub fn new(coin_name: String, token_limit: u128) -> Self {
        TokenLimit {
            coin_name,
            token_limit,
        }
    }
}

impl TokenLimits {
    pub fn new() -> Self {
        Self(UnorderedMap::new(b"token_limits".to_vec()))
    }

    pub fn add(&mut self, coin_name: &str, token_limit: &u128) {
        self.0
            .insert(&coin_name.to_string(), &token_limit.to_owned());
    }

    pub fn remove(&mut self, coin_name: &str) {
        self.0.remove(&coin_name.to_string());
    }

    pub fn get(&self, coin_name: &str) -> Option<u128> {
        if let Some(token_limit) = self.0.get(&coin_name.to_string()) {
            return Some(token_limit);
        }
        None
    }

    pub fn contains(&self, coin_name: &str) -> bool {
        if let Some(_) = self.0.keys().into_iter().find(|s| s == coin_name) {
            return true;
        }
        false
    }

    pub fn to_vec(&self) -> Vec<TokenLimit> {
        if !self.0.is_empty() {
            return self
                .0
                .to_vec()
                .into_iter()
                .map(|values| TokenLimit::new(values.0, values.1))
                .collect::<Vec<TokenLimit>>();
        }
        vec![]
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashSet;

    #[test]
    fn set_tokenlimit() {
        let coins = vec!["ICX", "NEAR", "sIcx"];
        let limits = vec![1000000000_u128, 100000002, 1000000003];

        let mut token_limits = TokenLimits::new();
        for (index, coin) in coins.into_iter().enumerate() {
            token_limits.add(coin, &limits[index]);
        }

        let result: HashSet<_> = token_limits.to_vec().into_iter().collect();
        let actual = vec![
            TokenLimit::new("ICX".to_string(), 1000000000_u128),
            TokenLimit::new("NEAR".to_string(), 100000002),
            TokenLimit::new("sIcx".to_string(), 1000000003),
        ];
        let actual: HashSet<_> = actual.into_iter().collect();
        assert_eq!(result, actual);
    }
    #[test]
    fn update_tokenlimit() {
        let coins = vec!["ICX", "NEAR", "sIcx"];
        let limits = vec![1000000000_u128, 100000002, 1000000003];

        let mut token_limits = TokenLimits::new();
        for (index, coin) in coins.into_iter().enumerate() {
            token_limits.add(coin, &limits[index]);
        }

        let token_limit = token_limits.get("ICX").unwrap();
        assert_eq!(1000000000_u128, token_limit.clone());

        token_limits.add("ICX", &10000000333_u128);

        assert_eq!(token_limits.get("ICX").unwrap(), 10000000333_u128);
    }

    #[test]
    fn check_for_non_existing_token_limit() {
        let coins = vec!["ICX", "NEAR", "sIcx"];
        let limits = vec![1000000000_u128, 100000002, 1000000003];

        let mut token_limits = TokenLimits::new();
        for (index, coin) in coins.into_iter().enumerate() {
            token_limits.add(coin, &limits[index]);
        }

        let expected = token_limits.contains("ICXV");

        assert_eq!(false, expected);
    }

    #[test]
    fn remove_token_from_token_limits() {
        let coins = vec!["ICX", "NEAR", "sIcx"];
        let limits = vec![1000000000_u128, 100000002, 1000000003];

        let mut token_limits = TokenLimits::new();
        for (index, coin) in coins.into_iter().enumerate() {
            token_limits.add(coin, &limits[index]);
        }

        token_limits.remove("ICX");
        let actual: HashSet<_> = vec![
            TokenLimit::new("NEAR".to_string(), 100000002_u128),
            TokenLimit::new("sIcx".to_string(), 1000000003_u128),
        ]
        .into_iter()
        .collect();

        assert_eq!(actual, token_limits.to_vec().into_iter().collect());
    }
}
