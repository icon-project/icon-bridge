use super::*;

type StorageBalance = u128;
#[derive(BorshDeserialize, BorshSerialize)]
pub struct StorageBalances {
    keys: UnorderedSet<(AccountId, AssetId)>,
    values: LookupMap<(AccountId, AssetId), StorageBalance>,
}

impl StorageBalances {
    pub fn new() -> Self {
        Self {
            keys: UnorderedSet::new(b"storage_keys".to_vec()),
            values: LookupMap::new(b"storage_values".to_vec()),
        }
    }

    pub fn add(&mut self, account: AccountId, asset_id: AssetId) {
        if !self.contains(&account, &asset_id) {
            self.keys.insert(&(account.to_owned(), asset_id));
            self.values.insert(&(account, asset_id), &u128::default());
        }
    }

    pub fn remove(&mut self, account: AccountId, asset_id: AssetId) {
        self.keys.remove(&(account.to_owned(), asset_id));
        self.values.remove(&(account, asset_id));
    }

    pub fn get(&self, account: &AccountId, asset_id: &AssetId) -> Option<StorageBalance> {
        self.values.get(&(account.to_owned(), asset_id.to_owned()))
    }

    pub fn contains(&self, account: &AccountId, asset_id: &AssetId) -> bool {
        self.keys
            .contains(&(account.to_owned(), asset_id.to_owned()))
    }

    pub fn set(
        &mut self,
        account: &AccountId,
        asset_id: &AssetId,
        storage_balance: StorageBalance,
    ) {
        self.values
            .insert(&(account.to_owned(), asset_id.to_owned()), &storage_balance);
    }

    pub fn to_vec(&self) -> Vec<((AccountId, AssetId), StorageBalance)> {
        if !self.keys.is_empty() {
            return self
                .keys
                .iter()
                .map(|keys| {
                    let storage_balance = self.values.get(&keys).unwrap();
                    (keys, storage_balance)
                })
                .collect::<Vec<((AccountId, AssetId), StorageBalance)>>();
        };

        vec![]
    }
}

#[cfg(test)]
mod tests {
    use std::convert::TryInto;

    use near_sdk::env;

    use super::*;
    use crate::types::Math;

    #[test]
    fn add_storage_balance() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let asset_id: [u8; 32] = env::sha256("asset1".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id.clone());
        let result = storage_balances.get(&account, &asset_id).unwrap();
        assert_eq!(result, 0);
    }

    #[test]
    fn add_storage_balance_existing() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let asset_id: [u8; 32] = env::sha256("asset1".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id.clone());
        let mut storage_balance = storage_balances.get(&account, &asset_id).unwrap();
        storage_balance.add(100).unwrap();

        let result = storage_balance.clone();
        storage_balances.set(&account.clone(), &asset_id, storage_balance);
        assert_eq!(result, storage_balances.get(&account, &asset_id).unwrap());
    }

    #[test]
    fn udpate_storage_balance() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let asset_id: [u8; 32] = env::sha256("asset1".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id.clone());
        storage_balances.set(&account.clone(), &asset_id, 1000);

        let mut storage_balance = storage_balances
            .get(&account.clone(), &asset_id.clone())
            .unwrap();

        storage_balance.add(3000000).unwrap();
        storage_balances.set(&account.clone(), &asset_id, storage_balance);

        assert_eq!(
            storage_balances
                .get(&account.clone(), &asset_id.clone())
                .unwrap(),
            3001000
        );

        let asset_id: [u8; 32] = env::sha256("asset2".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id.clone());
    }

    #[test]
    fn remove_storage_balance() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let asset_id: [u8; 32] = env::sha256("asset1".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id.clone());
        storage_balances.remove(account.clone(), asset_id);
        let result = storage_balances.get(&account, &asset_id);
        assert_eq!(result, None);
    }

    #[test]
    fn remove_storage_balance_non_existing() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let asset_id: [u8; 32] = env::sha256("asset1".as_bytes()).try_into().unwrap();
        storage_balances.remove(account.clone(), asset_id);
        let result = storage_balances.get(&account, &asset_id);
        assert_eq!(result, None);
    }

    #[test]
    fn to_vec_test() {
        let mut storage_balances = StorageBalances::new();
        let account = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let asset_id1: [u8; 32] = env::sha256("asset1".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id1.clone());
        storage_balances.set(&account.clone(), &asset_id1, 1000);

        let asset_id2: [u8; 32] = env::sha256("asset2".as_bytes()).try_into().unwrap();
        storage_balances.add(account.clone(), asset_id2.clone());

        let expected = vec![
            ((account.clone(), asset_id1), 1000),
            ((account.clone(), asset_id2), 0),
        ];

        assert_eq!(expected, storage_balances.to_vec())
    }
}
