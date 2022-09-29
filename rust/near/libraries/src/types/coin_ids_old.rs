use super::OldAssetId;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::env;
use near_sdk::serde::{Deserialize, Serialize};
use std::convert::TryInto;
use std::{collections::HashMap, hash::Hash};

#[derive(BorshDeserialize, BorshSerialize, Debug)]
pub struct OldCoinIds(HashMap<String, OldAssetId>);

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize)]
pub struct OldCoinProperty {
    coin_name: String,
    coin_id: OldAssetId,
}

impl OldCoinIds {
    pub fn new() -> Self {
        let coin_ids = HashMap::new();
        Self(coin_ids)
    }

    pub fn add(&mut self, coin_name: &str, coin_id: &Vec<u8>) {
        self.0
            .insert(coin_name.to_string(), coin_id.to_vec().try_into().unwrap());
    }

    pub fn remove(&mut self, coin_name: &str) {
        if self.0.contains_key(coin_name) {
            self.0.remove(coin_name);
        }
    }

    pub fn contains(&self, coin_name: &str) -> bool {
        return self.0.contains_key(coin_name);
    }

    pub fn get(&self, coin_name: &str) -> Option<&OldAssetId> {
        if let Some(coin_id) = self.0.get(coin_name) {
            return Some(coin_id);
        }
        None
    }

    pub fn to_vec(&self) -> Vec<OldCoinProperty> {
        if !self.0.is_empty() {
            return self
                .0
                .clone()
                .into_iter()
                .map(|v| OldCoinProperty {
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
    use near_sdk::{env, testing_env, VMContext};
    use std::{collections::HashSet, convert::TryInto};

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
    fn add_coin_property() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let coins = vec!["ICX", "NEAR", "sIcx"];

        let mut coin_store = OldCoinIds::new();
        coins.into_iter().for_each(|coin_name| {
            let coin_id = env::sha256(coin_name.as_bytes());
            coin_store.add(coin_name, &coin_id)
        });

        let result: HashSet<_> = coin_store.to_vec().into_iter().collect();
        let actual = vec![
            OldCoinProperty {
                coin_name: "ICX".to_string(),
                coin_id: env::sha256("ICX".to_string().as_bytes())
                    .try_into()
                    .unwrap(),
            },
            OldCoinProperty {
                coin_name: "NEAR".to_string(),
                coin_id: env::sha256("NEAR".as_bytes()).try_into().unwrap(),
            },
            OldCoinProperty {
                coin_name: "sIcx".to_string(),
                coin_id: env::sha256("sIcx".as_bytes()).try_into().unwrap(),
            },
        ];
        let actual: HashSet<_> = actual.into_iter().collect();
        assert_eq!(result, actual);
    }
    #[test]
    fn get_coin_id() {
        let context = get_context(vec![], false);
        testing_env!(context);
        let coins = vec!["ICX", "NEAR", "sIcx"];

        let mut coin_store = OldCoinIds::new();
        coins.into_iter().for_each(|coin_name| {
            let coin_id = env::sha256(coin_name.as_bytes());
            coin_store.add(coin_name, &coin_id)
        });

        let coin_id = coin_store.get("sIcx").unwrap();
        assert_eq!(coin_id, &env::sha256("sIcx".as_bytes()).as_slice());
    }
}
