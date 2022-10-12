use bts::BtpTokenService;
use near_sdk::{
    env, json_types::U128, serde_json::to_value, testing_env, AccountId, PromiseResult, VMContext, test_utils::VMContextBuilder, Gas, VMConfig, RuntimeFeesConfig
};
use std::{
    collections::{HashMap, HashSet},
    convert::TryInto,
};
pub mod accounts;
use accounts::*;
use libraries::types::{AccountBalance, Asset, AssetItem, Math, TokenLimits, WrappedNativeCoin};
mod token;
use token::*;
pub type Coin = Asset<WrappedNativeCoin>;

fn get_context(is_view: bool, signer_account_id: AccountId, attached_deposit: u128) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .storage_usage(env::storage_usage())
        .prepaid_gas(Gas(10u64.pow(18)))
        .attached_deposit(attached_deposit)
        .build()
}

#[test]
fn register_token() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Coin>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(icx_coin.clone(), coin_id);

    let result = contract.coins();
    let expected = to_value(vec![
        AssetItem {
            name: nativecoin.name().to_owned(),
            network: nativecoin.network().to_owned(),
            symbol: nativecoin.symbol().to_owned(),
        },
        AssetItem {
            name: icx_coin.name().to_owned(),
            network: icx_coin.network().to_owned(),
            symbol: icx_coin.symbol().to_owned(),
        },
    ])
    .unwrap();
    assert_eq!(result, expected);
}

#[test]
#[should_panic(expected = "BSHRevertAlreadyExistsToken")]
fn register_existing_token() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Coin>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(icx_coin.clone(), coin_id);
    contract.register(icx_coin.clone());
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsPermission")]
fn register_token_permission() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 0));
    let icx_coin = <Coin>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(icx_coin.clone(), coin_id);
}

#[test]
fn get_registered_coin_id() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let coin_id = contract.coin_id("NEAR").unwrap();
    let expected: [u8; 32] = env::sha256(nativecoin.name().as_bytes())
        .try_into()
        .unwrap();
    assert_eq!(coin_id, expected)
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken: ICON")]
fn get_non_exist_coin_id() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let coin_id = contract
        .coin_id("ICON")
        .map_err(|err| format!("{}", err))
        .unwrap();
}

#[test]
#[cfg(feature = "testable")]
fn set_token_limit() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let coins = vec!["NEAR".to_string()];
    let limits = vec![10000000000000000000000_u128];
    contract.set_token_limit(coins, limits).unwrap();
    let tokenlimits = contract.get_token_limit();

    assert_eq!(
        tokenlimits.get("NEAR").unwrap(),
        &10000000000000000000000_u128
    )
}

#[test]
#[cfg(feature = "testable")]
fn update_token_limit() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let coins = vec!["NEAR".to_string()];
    let limits = vec![10000000000000000000000_u128];
    contract.set_token_limit(coins, limits).unwrap();

    let coins = vec!["NEAR".to_string()];
    let limits = vec![10000000000000000000003_u128];
    contract.set_token_limit(coins, limits).unwrap();

    let tokenlimits = contract.get_token_limit().get("NEAR").unwrap();
    assert_eq!(tokenlimits, &10000000000000000000003_u128)
}

#[test]
fn query_token_metadata() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Coin>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(icx_coin.clone(), coin_id);

    let result = contract.coin(icx_coin.name().to_string());

    assert_eq!(icx_coin, result);
}
