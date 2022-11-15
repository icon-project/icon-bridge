use super::AssetId;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::env;
use near_sdk::serde::{Deserialize, Serialize};
use std::convert::TryInto;
use std::{collections::HashMap, hash::Hash};

#[derive(BorshDeserialize, BorshSerialize, Debug)]
pub struct TokenIds(HashMap<String, AssetId>);

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize)]
pub struct TokenProperty {
    token_name: String,
    token_id: AssetId,
}

impl TokenIds {
    pub fn new() -> Self {
        let token_ids = HashMap::new();
        Self(token_ids)
    }

    pub fn add(&mut self, token_name: &str, token_id: AssetId) {
        self.0.insert(token_name.to_string(), token_id);
    }

    pub fn remove(&mut self, token_name: &str) {
        if self.0.contains_key(token_name) {
            self.0.remove(token_name);
        }
    }

    pub fn contains(&self, token_name: &str) -> bool {
        self.0.contains_key(token_name)
    }

    pub fn get(&self, token_name: &str) -> Option<&AssetId> {
        if let Some(token_id) = self.0.get(token_name) {
            return Some(token_id);
        }
        None
    }

    pub fn to_vec(&self) -> Vec<TokenProperty> {
        if !self.0.is_empty() {
            return self
                .0
                .clone()
                .into_iter()
                .map(|v| TokenProperty {
                    token_name: v.0,
                    token_id: v.1,
                })
                .collect();
        }
        vec![]
    }
}
#[cfg(test)]
mod tests {
    use super::*;
    use std::{collections::HashSet, convert::TryInto};

    #[test]
    fn add_token_property() {
        let tokens = vec!["ICX", "NEAR", "sIcx"];

        let mut token_store = TokenIds::new();
        tokens.into_iter().for_each(|token_name| {
            let token_id = env::sha256(token_name.as_bytes());
            token_store.add(token_name, token_id.try_into().unwrap())
        });

        let result: HashSet<_> = token_store.to_vec().into_iter().collect();
        let actual = vec![
            TokenProperty {
                token_name: "ICX".to_string(),
                token_id: env::sha256("ICX".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            },
            TokenProperty {
                token_name: "NEAR".to_string(),
                token_id: env::sha256("NEAR".as_bytes()).try_into().unwrap(),
            },
            TokenProperty {
                token_name: "sIcx".to_string(),
                token_id: env::sha256("sIcx".as_bytes()).try_into().unwrap(),
            },
        ];
        let actual: HashSet<_> = actual.into_iter().collect();
        assert_eq!(result, actual);
    }
    #[test]
    fn get_token_id() {
        let tokens = vec!["ICX", "NEAR", "sIcx"];

        let mut token_store = TokenIds::new();
        tokens.into_iter().for_each(|token_name| {
            let coin_id = env::sha256(token_name.as_bytes());
            token_store.add(token_name, coin_id.try_into().unwrap())
        });

        let coin_id = token_store.get("sIcx").unwrap();
        assert_eq!(coin_id, &env::sha256("sIcx".as_bytes()).as_slice());
    }
}
