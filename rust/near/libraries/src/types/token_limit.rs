use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::serde::{Deserialize, Serialize};
use near_sdk::AccountId;
use std::{collections::HashMap, hash::Hash};

#[derive(BorshDeserialize, BorshSerialize, Debug)]
pub struct TokenLimits(HashMap<String, u128>);

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize)]
pub struct TokenLimit {
    coin_name: String,
    token_limit: u128,
}

impl TokenLimits {
    pub fn new() -> Self {
        let token_limit = HashMap::new();
        Self(token_limit)
    }

    pub fn add(&mut self, coin_name: &str, token_limit: &u128) {
        self.0.insert(coin_name.to_string(), token_limit.to_owned());
    }

    pub fn remove(&mut self, coin_name: &str) {
        if self.0.contains_key(coin_name) {
            self.0.remove(coin_name);
        }
    }

    pub fn contains(&self, coin_name: &str) -> bool {
        return self.0.contains_key(coin_name);
    }

    pub fn get(&self, coin_name: &str) -> Option<&u128> {
        if let Some(token_limit) = self.0.get(coin_name) {
            return Some(token_limit);
        }
        None
    }

    pub fn to_vec(&self) -> Vec<TokenLimit> {
        if !self.0.is_empty() {
            return self
                .0
                .clone()
                .into_iter()
                .map(|v| TokenLimit {
                    coin_name: v.0,
                    token_limit: v.1,
                })
                .collect();
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
    fn set_tokenlimit() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let coins = vec!["ICX", "NEAR", "sIcx"];
        let limits = vec![1000000000_u128, 100000002, 1000000003];

        let mut token_limits = TokenLimits::new();
        for (index, coin) in coins.into_iter().enumerate() {
            token_limits.add(coin, &limits[index]);
        }

        let result: HashSet<_> = token_limits.to_vec().into_iter().collect();
        let actual = vec![
            TokenLimit {
                coin_name: "ICX".to_string(),
                token_limit: 1000000000_u128,
            },
            TokenLimit {
                coin_name: "NEAR".to_string(),
                token_limit: 100000002,
            },
            TokenLimit {
                coin_name: "sIcx".to_string(),
                token_limit: 1000000003,
            },
        ];
        let actual: HashSet<_> = actual.into_iter().collect();
        assert_eq!(result, actual);
    }
    #[test]
    fn update_tokenlimit() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let coins = vec!["ICX", "NEAR", "sIcx"];
        let limits = vec![1000000000_u128, 100000002, 1000000003];

        let mut token_limits = TokenLimits::new();
        for (index, coin) in coins.into_iter().enumerate() {
            token_limits.add(coin, &limits[index]);
        }

        let token_limit = token_limits.get("ICX").unwrap();
        assert_eq!(1000000000_u128, token_limit.clone());

        token_limits.add("ICX", &10000000333_u128);

        assert_eq!(token_limits.get("ICX").unwrap(), &10000000333_u128);
    }
}
