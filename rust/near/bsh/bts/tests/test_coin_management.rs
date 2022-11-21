#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_mut)]
use bts::BtpTokenService;
use near_sdk::{
    env, json_types::U128, serde_json::to_value, test_utils::VMContextBuilder, testing_env,
    AccountId, Gas, PromiseResult, RuntimeFeesConfig, VMConfig, VMContext,
};
use std::{
    collections::{HashMap, HashSet},
    convert::TryInto,
};
pub mod accounts;
use accounts::*;
use libraries::types::{
    AccountBalance, Asset, AssetItem, AssetMetadata, Math, TokenLimits, WrappedNativeCoin,
};
mod token;
use token::*;
pub type Token = Asset<WrappedNativeCoin>;

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
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Token>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);

    let result = contract.tokens();
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
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Token>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);
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
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 0));
    let icx_coin = <Token>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);
}

#[test]
fn get_registered_token_id() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let token_id = contract.token_id("NEAR").unwrap();
    let expected: [u8; 32] = env::sha256(nativecoin.name().as_bytes())
        .try_into()
        .unwrap();
    assert_eq!(token_id, expected)
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken: ICON")]
fn get_non_exist_token_id() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let token_id = contract
        .token_id("ICON")
        .map_err(|error| error.to_string())
        .unwrap();
}

#[test]
fn set_token_limit() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let tokens = vec!["NEAR".to_string()];
    let limits = vec![10000000000000000000000_u128];
    contract.set_token_limit(&tokens, &limits).unwrap();
    let token_limits = contract.get_token_limit("NEAR".to_string());

    assert_eq!(token_limits, U128(10000000000000000000000))
}

#[test]
fn update_token_limit() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let tokens = vec!["NEAR".to_string()];
    let limits = vec![10000000000000000000000_u128];
    contract.set_token_limit(&tokens, &limits).unwrap();

    let tokens = vec!["NEAR".to_string()];
    let limits = vec![10000000000000000000003_u128];
    contract.set_token_limit(&tokens, &limits).unwrap();

    let token_limits = contract.get_token_limit("NEAR".to_string());
    assert_eq!(token_limits, U128(10000000000000000000003))
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
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Token>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);

    let result = contract.token(icx_coin.name().to_string());

    assert_eq!(icx_coin, result);
}

#[test]
fn query_token_fee_ratio() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let icx_coin = <Token>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);

    let result = contract.get_fee_ratio(icx_coin.name().to_string());

    assert_eq!(
        (
            icx_coin.metadata().fee_numerator().into(),
            icx_coin.metadata().fixed_fee().into()
        ),
        result
    );
}
