use std::convert::{TryFrom, TryInto};

use bts::{BtpTokenService, Coin};
use near_sdk::{env, serde_json::to_value, testing_env, AccountId, PromiseResult, VMContext};
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

fn get_context(
    input: Vec<u8>,
    is_view: bool,
    signer_account_id: AccountId,
    attached_deposit: u128,
) -> VMContext {
    VMContext {
        current_account_id: alice().to_string(),
        signer_account_id: signer_account_id.to_string(),
        signer_account_pk: vec![0, 1, 2],
        predecessor_account_id: signer_account_id.to_string(),
        input,
        block_index: 0,
        block_timestamp: 0,
        account_balance: 0,
        account_locked_balance: 0,
        storage_usage: 0,
        attached_deposit,
        prepaid_gas: 10u64.pow(18),
        random_seed: vec![0, 1, 2],
        is_view,
        output_data_receivers: vec![],
        epoch_height: 19,
    }
}

#[test]
fn register_token() {
    let context = |v: AccountId, d: u128| (get_context(vec![], false, v, d));
    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
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
    let context = |v: AccountId, d: u128| (get_context(vec![], false, v, d));
    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
        Default::default(),
        Default::default(),
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
    let context = |v: AccountId, d: u128| (get_context(vec![], false, v, d));
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
    let context = |v: AccountId, d: u128| (get_context(vec![], false, v, d));
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
    let context = |v: AccountId, d: u128| (get_context(vec![], false, v, d));
    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
        Default::default(),
        Default::default(),
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
