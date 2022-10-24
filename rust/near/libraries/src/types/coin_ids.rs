use super::AssetId;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::env;
use near_sdk::serde::{Deserialize, Serialize};
use std::convert::TryInto;
use std::{collections::HashMap, hash::Hash};

#[derive(BorshDeserialize, BorshSerialize, Debug)]
pub struct CoinIds(HashMap<String, AssetId>);

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize)]
pub struct CoinProperty {
    coin_name: String,
    coin_id: AssetId,
}

impl CoinIds {
    pub fn new() -> Self {
        let coin_ids = HashMap::new();
        Self(coin_ids)
    }

    pub fn add(&mut self, coin_name: &str, coin_id: AssetId) {
        self.0.insert(coin_name.to_string(), coin_id);
    }

    pub fn remove(&mut self, coin_name: &str) {
        if self.0.contains_key(coin_name) {
            self.0.remove(coin_name);
        }
    }

    pub fn contains(&self, coin_name: &str) -> bool {
        return self.0.contains_key(coin_name);
    }

    pub fn get(&self, coin_name: &str) -> Option<&AssetId> {
        if let Some(coin_id) = self.0.get(coin_name) {
            return Some(coin_id);
        }
        None
    }

    pub fn to_vec(&self) -> Vec<CoinProperty> {
        if !self.0.is_empty() {
            return self
                .0
                .clone()
                .into_iter()
                .map(|v| CoinProperty {
                    coin_name: v.0,
                    coin_id: v.1,
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
    fn add_coin_property() {

        let coins = vec!["ICX", "NEAR", "sIcx"];

        let mut coin_store = CoinIds::new();
        coins.into_iter().for_each(|coin_name| {
            let coin_id = env::sha256(coin_name.as_bytes());
            coin_store.add(coin_name, coin_id.try_into().unwrap())
        });

        let result: HashSet<_> = coin_store.to_vec().into_iter().collect();
        let actual = vec![
            CoinProperty {
                coin_name: "ICX".to_string(),
                coin_id: env::sha256("ICX".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            },
            CoinProperty {
                coin_name: "NEAR".to_string(),
                coin_id: env::sha256("NEAR".as_bytes()).try_into().unwrap(),
            },
            CoinProperty {
                coin_name: "sIcx".to_string(),
                coin_id: env::sha256("sIcx".as_bytes()).try_into().unwrap(),
            },
        ];
        let actual: HashSet<_> = actual.into_iter().collect();
        assert_eq!(result, actual);
    }
    #[test]
    fn get_coin_id() {

        let coins = vec!["ICX", "NEAR", "sIcx"];

        let mut coin_store = CoinIds::new();
        coins.into_iter().for_each(|coin_name| {
            let coin_id = env::sha256(coin_name.as_bytes());
            coin_store.add(coin_name, coin_id.try_into().unwrap())
        });

        let coin_id = coin_store.get("sIcx").unwrap();
        assert_eq!(coin_id, &env::sha256("sIcx".as_bytes()).as_slice());
    }
}
