use std::convert::TryFrom;

use lazy_static::lazy_static;
use libraries::types::{AssetMetadataExtras, WrappedNativeCoin};
use near_sdk::AccountId;

lazy_static! {
    pub static ref NATIVE_COIN: WrappedNativeCoin = WrappedNativeCoin::new(
        "NEAR".into(),
        "Wrapped TOKEN From SOURCE_CHAIN".into(),
        "NEAR".into(),
        None,
        "0x1.near".into(),
        1000,
        1,
        None,
        None
    );
    pub static ref ICON_COIN: WrappedNativeCoin = WrappedNativeCoin::new(
        "ICON".into(),
        "Wrapped TOKEN From SOURCE_CHAIN".into(),
        "ICX".into(),
        Some(AccountId::try_from("icx.near".to_string()).unwrap()),
        "0x1.icon".into(),
        1000,
        1,
        Some(AssetMetadataExtras {
            spec: "ft-1.0.0".to_string(),
            icon: None,
            reference: None,
            reference_hash: None,
            decimals: 24
        }),
        None
    );
    pub static ref WNEAR: WrappedNativeCoin = WrappedNativeCoin::new(
        "WNEAR".into(),
        "Wrapped TOKEN From SOURCE_CHAIN".into(),
        "wNEAR".into(),
        Some(AccountId::try_from("wnear.near".to_string()).unwrap()),
        "0x1.near".into(),
        1000,
        1,
        Some(AssetMetadataExtras {
            icon: None,
            decimals: 24,
            reference: None,
            reference_hash: None,
            spec: "ft-1.0.0".to_string()
        }),
        None
    );
    pub static ref BALN: WrappedNativeCoin = WrappedNativeCoin::new(
        "BALN".into(),
        "Wrapped TOKEN From SOURCE_CHAIN".into(),
        "BALN".into(),
        Some(AccountId::try_from("baln.icon".to_string()).unwrap()),
        "0x1.icon".into(),
        1000,
        1,
        Some(AssetMetadataExtras {
            icon: None,
            decimals: 24,
            reference: None,
            reference_hash: None,
            spec: "ft-1.0.0".to_string()
        }),
        None
    );
    pub static ref NEAR_NATIVE_COIN: WrappedNativeCoin = WrappedNativeCoin::new(
        "btp-0x1.near-NEAR".into(),
        "Wrapped TOKEN From SOURCE_CHAIN".into(),
        "NEAR".into(),
        None,
        "0x1.near".into(),
        1000,
        1,
        None,
        None
    );
}
