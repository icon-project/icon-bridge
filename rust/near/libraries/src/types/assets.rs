use super::*;

#[derive(Debug, Clone, PartialEq, PartialOrd, Ord, Eq, Serialize, Hash)]
#[serde(crate = "near_sdk::serde")]
pub struct AssetItem {
    pub name: String,
    pub network: String,
    pub symbol: String,
}

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Assets<T: AssetMetadata> {
    keys: UnorderedSet<AssetId>,
    values: Metadata<T>,
}
#[derive(BorshDeserialize, BorshSerialize)]
pub struct Metadata<T: AssetMetadata>(LookupMap<AssetId, Asset<T>>);

impl<T: BorshDeserialize + BorshSerialize + AssetMetadata> Metadata<T> {
    fn new() -> Self {
        Self(LookupMap::new(StorageKey::Assets(KeyType::Value)))
    }

    fn add(&mut self, asset_id: &AssetId, asset: &Asset<T>) {
        self.0.insert(asset_id, asset);
    }

    fn remove(&mut self, asset_id: &AssetId) {
        self.0.remove(asset_id);
    }

    fn get(&self, asset_id: &AssetId) -> Option<Asset<T>> {
        if let Some(asset) = self.0.get(asset_id) {
            return Some(asset);
        }
        None
    }
}

impl<T: BorshDeserialize + BorshSerialize + AssetMetadata> Assets<T> {
    pub fn new() -> Self {
        Self {
            keys: UnorderedSet::new(StorageKey::Assets(KeyType::Key)),
            values: Metadata::new(),
        }
    }

    pub fn add(&mut self, asset_id: &AssetId, asset: &Asset<T>) {
        self.keys.insert(asset_id);
        self.values.add(asset_id, asset);
    }

    pub fn remove(&mut self, asset_id: &AssetId) {
        self.keys.remove(asset_id);
        self.values.remove(asset_id);
    }

    pub fn contains(&self, asset_id: &AssetId) -> bool {
        self.keys.contains(asset_id)
    }

    pub fn get(&self, asset_id: &AssetId) -> Option<Asset<T>> {
        self.values.get(asset_id)
    }

    pub fn set(&mut self, asset_id: &AssetId, asset: &Asset<T>) {
        self.values.add(asset_id, asset)
    }

    pub fn to_vec(&self) -> Vec<AssetItem> {
        self.keys
            .to_vec()
            .iter()
            .map(|asset_id| {
                let metdata = self.values.get(asset_id).unwrap();
                AssetItem {
                    name: metdata.name().clone(),
                    network: metdata.network().clone(),
                    symbol: metdata.symbol().clone(),
                }
            })
            .collect::<Vec<AssetItem>>()
    }
}

impl<T: BorshDeserialize + BorshSerialize + AssetMetadata> Default for Assets<T>  {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use crate::types::WrappedNativeCoin;
    use crate::types::{asset::*, assets::*};
    use near_sdk::serde_json;
    use std::collections::HashSet;

    #[test]
    fn add_token() {
        let mut tokens = Assets::new();
        let native_coin = WrappedNativeCoin::new(
            "ABC Asset".to_string(),
            "ABC Asset".to_string(),
            "ABC".to_string(),
            None,
            "0x1.near".to_string(),
            10000,
            10000,
            None,
            Some(10000),
        );

        tokens.add(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin.clone()),
        );

        let result = tokens.contains(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, true);

        let result = tokens.get(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, Some(<Asset<WrappedNativeCoin>>::new(native_coin)));
    }

    #[test]
    fn add_existing_token() {
        let mut tokens = Assets::new();
        let native_coin = WrappedNativeCoin::new(
            "ABC Asset".to_string(),
            "ABC Asset".to_string(),
            "ABC".to_string(),
            None,
            "0x1.near".to_string(),
            10000,
            10000,
            None,
            None,
        );

        tokens.add(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin.clone()),
        );
        tokens.add(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin.clone()),
        );
        let result = tokens.to_vec();

        let expected: Vec<AssetItem> = vec![AssetItem {
            name: "ABC Asset".to_string(),
            network: "0x1.near".to_string(),
            symbol: "ABC".to_string(),
        }];
        assert_eq!(result, expected);
    }

    #[test]
    fn remove_token() {
        let mut tokens = Assets::new();
        let native_coin = WrappedNativeCoin::new(
            "ABC Asset".to_string(),
            "ABC Asset".to_string(),
            "ABC".to_string(),
            None,
            "0x1.near".to_string(),
            10000,
            10000,
            None,
            None,
        );

        tokens.add(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin.clone()),
        );

        tokens.remove(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        let result = tokens.contains(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, false);

        let result = tokens.get(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, None);
    }

    #[test]
    fn remove_token_non_existing() {
        let mut tokens = <Assets<WrappedNativeCoin>>::new();
        tokens.remove(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        let result = tokens.contains(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
        );
        assert_eq!(result, false);
    }

    #[test]
    fn to_vec_tokens() {
        let mut tokens = <Assets<WrappedNativeCoin>>::new();
        let native_coin_1 = WrappedNativeCoin::new(
            "ABC Asset".to_string(),
            "ABC Asset".to_string(),
            "ABC".to_string(),
            None,
            "0x1.near".to_string(),
            10000,
            10000,
            None,
            None,
        );
        let native_coin_2 = WrappedNativeCoin::new(
            "DEF Asset".to_string(),
            "DEF Asset".to_string(),
            "DEF".to_string(),
            None,
            "0x1.bsc".to_string(),
            10000,
            10000,
            None,
            None,
        );

        tokens.add(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin_1),
        );
        tokens.add(
            &env::sha256("DEF Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin_2),
        );
        let tokens = tokens.to_vec();
        let expected_tokens: Vec<AssetItem> = vec![
            AssetItem {
                name: "ABC Asset".to_string(),
                network: "0x1.near".to_string(),
                symbol: "ABC".to_string(),
            },
            AssetItem {
                name: "DEF Asset".to_string(),
                network: "0x1.bsc".to_string(),
                symbol: "DEF".to_string(),
            },
        ];
        let result: HashSet<_> = tokens.iter().collect();
        let expected: HashSet<_> = expected_tokens.iter().collect();
        assert_eq!(result, expected);
    }

    #[test]
    fn to_vec_tokens_value() {
        let mut tokens = <Assets<WrappedNativeCoin>>::new();
        let native_coin = WrappedNativeCoin::new(
            "ABC Asset".to_string(),
            "ABC Asset".to_string(),
            "ABC".to_string(),
            None,
            "0x1.near".to_string(),
            10000,
            10000,
            None,
            None,
        );
        tokens.add(
            &env::sha256("ABC Asset".to_string().as_bytes())
                .try_into()
                .unwrap(),
            &<Asset<WrappedNativeCoin>>::new(native_coin),
        );
        let tokens = serde_json::to_value(tokens.to_vec()).unwrap();
        assert_eq!(
            tokens,
            serde_json::json!(
                [
                    {
                        "name": "ABC Asset",
                        "network": "0x1.near",
                        "symbol": "ABC"
                    }
                ]
            )
        );
    }
}
