use crate::types::btp_address::Network;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::serde::{Deserialize, Serialize};

use super::fungible_token::AssetMetadataExtras;
pub type AssetId = [u8; 32];

pub trait AssetMetadata {
    fn name(&self) -> &String;
    fn label(&self) -> &String;
    fn network(&self) -> &Network;
    fn symbol(&self) -> &String;
    fn fee_numerator(&self) -> u128;
    fn fee_numerator_mut(&mut self) -> &mut u128;
    fn fixed_fee(&self) -> u128;
    fn fixed_fee_mut(&mut self) -> &mut u128;
    fn metadata(&self) -> &Self;
    fn extras(&self) -> &Option<AssetMetadataExtras>;
}

#[derive(BorshDeserialize, BorshSerialize, Clone, Deserialize, Serialize, Debug, PartialEq)]
#[serde(crate = "near_sdk::serde")]
pub struct Asset<T: AssetMetadata> {
    pub metadata: T,
}

impl<T: AssetMetadata> Asset<T> {
    pub fn new(asset: T) -> Self {
        Self { metadata: asset }
    }

    pub fn name(&self) -> &String {
        self.metadata.name()
    }

    pub fn label(&self) -> &String {
        self.metadata.label()
    }

    pub fn network(&self) -> &String {
        self.metadata.network()
    }

    pub fn symbol(&self) -> &String {
        self.metadata.symbol()
    }

    pub fn metadata(&self) -> &T {
        &self.metadata
    }

    pub fn metadata_mut(&mut self) -> &mut T {
        &mut self.metadata
    }

    pub fn extras(&self) -> &Option<AssetMetadataExtras> {
        &self.metadata.extras()
    }
}
