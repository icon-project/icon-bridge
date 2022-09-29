use crate::types::{Math, OldAssetId};
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::serde::{Deserialize, Serialize};
use near_sdk::AccountId;
use near_sdk::Balance;
use std::collections::HashMap;

type AssetFee = u128;

#[derive(Clone, BorshDeserialize, BorshSerialize)]
pub struct OldAssetFees(HashMap<OldAssetId, AssetFee>);

impl OldAssetFees {
    pub fn new() -> Self {
        Self(HashMap::new())
    }

    pub fn add(&mut self, asset_id: &OldAssetId) {
        self.0.insert(asset_id.clone(), u128::default());
    }

    pub fn remove(&mut self, asset_id: &OldAssetId) {
        self.0.remove(asset_id);
    }

    pub fn get(&self, asset_id: &OldAssetId) -> Option<&AssetFee> {
        if let Some(asset_fee) = self.0.get(asset_id) {
            return Some(asset_fee);
        }
        None
    }

    pub fn set(&mut self, asset_id: &OldAssetId, asset_fee: AssetFee) {
        self.0.insert(asset_id.clone(), asset_fee);
    }
}
