use super::*;

type AssetFee = u128;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct AssetFees(LookupMap<AssetId, AssetFee>);

impl AssetFees {
    pub fn new() -> Self {
        Self(LookupMap::new(StorageKey::AssetFees))
    }

    pub fn add(&mut self, asset_id: &AssetId) {
        self.0.insert(asset_id, &u128::default());
    }

    pub fn remove(&mut self, asset_id: &AssetId) {
        self.0.remove(asset_id);
    }

    pub fn get(&self, asset_id: &AssetId) -> Option<AssetFee> {
        if let Some(asset_fee) = self.0.get(asset_id) {
            return Some(asset_fee);
        }
        None
    }

    pub fn set(&mut self, asset_id: &AssetId, asset_fee: AssetFee) {
        self.0.insert(asset_id, &asset_fee);
    }
}

impl Default for AssetFees {
    fn default() -> Self {
        Self::new()
    }
}