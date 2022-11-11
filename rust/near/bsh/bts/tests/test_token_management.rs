#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_mut)]

use std::convert::{TryFrom, TryInto};

use bts::{BtpTokenService, Coin};
use near_sdk::{
    env, serde_json::to_value, test_utils::VMContextBuilder, testing_env, AccountId, Gas,
    PromiseResult, RuntimeFeesConfig, VMConfig, VMContext,
};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::{BtpMessage, SerializedMessage},
    Asset, AssetItem, WrappedNativeCoin,
};
mod token;
use token::*;

pub type Token = Asset<WrappedNativeCoin>;
pub type TokenItem = AssetItem;

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
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
    let coin_id: [u8; 32] = env::sha256(baln.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(baln.clone(), coin_id);

    let result = contract.coins();
    let expected = to_value(vec![
        AssetItem {
            name: nativecoin.name().to_owned(),
            network: nativecoin.network().to_owned(),
            symbol: nativecoin.symbol().to_owned(),
        },
        TokenItem {
            name: baln.name().to_owned(),
            network: baln.network().to_owned(),
            symbol: baln.symbol().to_owned(),
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
        context(alice(), 1_000_000_000_000_000_000_000_000),
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
        nativecoin,
    );
    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
    let coin_id: [u8; 32] = env::sha256(baln.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(baln.clone(), coin_id);

    contract.register(baln.clone());
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsPermission")]
fn register_token_permission() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin,
    );
    testing_env!(context(chuck(), 0));
    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken: ICON")]
fn get_non_exist_token_id() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin,
    );
    let coin_id = contract
        .coin_id("ICON")
        .map_err(|err| format!("{}", err))
        .unwrap();
}

#[test]
fn get_registered_token_id() {
    let context = |v: AccountId, d: u128| (get_context(false, v, d));
    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
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
        nativecoin,
    );
    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
    let coin_id: [u8; 32] = env::sha256(baln.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_coin_callback(baln.clone(), coin_id);

    let token_id = contract.coin_id("BALN").unwrap();
    let expected: [u8; 32] = env::sha256(baln.name().as_bytes()).try_into().unwrap();
    assert_eq!(token_id, expected)
}
