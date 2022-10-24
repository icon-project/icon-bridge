use bts::BtpTokenService;
use near_sdk::{testing_env, AccountId, VMContext, test_utils::VMContextBuilder, Gas, env};
use std::collections::HashSet;
pub mod accounts;
use accounts::*;
use libraries::types::{Asset, WrappedNativeCoin};
mod token;
use token::*;
pub type Coin = Asset<WrappedNativeCoin>;

fn get_context(is_view: bool, signer_account_id: AccountId) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .storage_usage(env::storage_usage())
        .prepaid_gas(Gas(10u64.pow(18)))
        .build()
}

#[test]
fn add_owner_new_owner() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );

    contract.add_owner(carol());

    let owners = contract.get_owners();
    let result: HashSet<_> = owners.iter().collect();
    let expected_owners: Vec<AccountId> = vec![alice(), carol()];
    let expected: HashSet<_> = expected_owners.iter().collect();
    assert_eq!(result, expected);
}

#[test]
#[should_panic(expected = "BSHRevertAlreadyExistsOwner")]
fn add_owner_existing_owner() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );

    contract.add_owner(alice());
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsPermission")]
fn add_owner_permission() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );
    testing_env!(context(chuck()));
    contract.add_owner(carol());
}

#[test]
fn remove_owner_existing_owner() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );

    contract.add_owner(carol());
    contract.add_owner(charlie());

    contract.remove_owner(alice());
    let owners = contract.get_owners();
    let result: HashSet<_> = owners.iter().collect();
    let expected_owners: Vec<AccountId> = vec![carol(), charlie()];
    let expected: HashSet<_> = expected_owners.iter().collect();
    assert_eq!(result, expected);
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsPermission")]
fn remove_owner_permission() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );

    contract.add_owner(carol());

    testing_env!(context(chuck()));
    contract.add_owner(charlie());
}

#[test]
#[should_panic(expected = "BSHRevertLastOwner")]
fn remove_owner_last_owner() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );

    contract.remove_owner(alice());
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsOwner")]
fn remove_owner_non_existing_owner() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        <Coin>::new(NATIVE_COIN.to_owned()),
    );

    contract.remove_owner(carol());
}
