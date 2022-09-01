use std::{collections::HashSet, iter::FromIterator, result, str::FromStr};

use btp_common::errors::BshError;
use bts::{BtpTokenService, Coin};
use libraries::types::BTPAddress;
use near_sdk::{
    env, serde_json::to_value, test_utils::test_env::alice, testing_env, AccountId, PromiseResult,
    VMContext,
};
mod token;
use token::*;
pub mod accounts;
use accounts::*;

fn get_context(
    input: Vec<u8>,
    is_view: bool,
    signer_account_id: AccountId,
    attached_deposit: u128,
    storage_usage: u64,
    account_balance: u128,
) -> VMContext {
    VMContext {
        current_account_id: alice().to_string(),
        signer_account_id: signer_account_id.to_string(),
        signer_account_pk: vec![0, 1, 2],
        predecessor_account_id: signer_account_id.to_string(),
        input,
        block_index: 0,
        block_timestamp: 0,
        account_balance,
        account_locked_balance: 0,
        storage_usage,
        attached_deposit,
        prepaid_gas: 10u64.pow(18),
        random_seed: vec![0, 1, 2],
        is_view,
        output_data_receivers: vec![],
        epoch_height: 19,
    }
}

#[test]
fn add_user_to_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let chuck_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", chuck().clone()).as_str()).unwrap();

    let charlie_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", charlie().clone()).as_str()).unwrap();

    let users = vec![chuck_btpaddress, charlie_btpaddress];

    let result = contract.add_to_blacklist(users);
    match result {
        Ok(()) => {
            let users = contract.get_blacklisted_user();
            let result: HashSet<_> = users.iter().collect();
            let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
            let expected: HashSet<_> = expected_users.iter().collect();
            assert_eq!(expected, result)
        }
        Err(_) => todo!(),
    }
}

#[test]
fn add_already_blacklisted_user_to_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let chuck_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", chuck().clone()).as_str()).unwrap();

    let charlie_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", charlie().clone()).as_str()).unwrap();

    let users = vec![chuck_btpaddress, charlie_btpaddress];

    let result = contract.add_to_blacklist(users.clone());
    match result {
        Ok(()) => {
            let users = contract.get_blacklisted_user();
            let result: HashSet<_> = users.iter().collect();
            let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
            let expected: HashSet<_> = expected_users.iter().collect();
            assert_eq!(expected, result)
        }
        Err(_) => todo!(),
    }

    let result = contract.add_to_blacklist(users.clone());
    match result {
        Ok(()) => {}
        Err(err) => {
            assert_eq!(BshError::UserAlreadyBlacklisted, err)
        }
    }
}

#[test]
fn remove_blacklisted_user_from_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let chuck_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", chuck().clone()).as_str()).unwrap();

    let charlie_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", charlie().clone()).as_str()).unwrap();

    let users = vec![chuck_btpaddress.clone(), charlie_btpaddress.clone()];

    let result = contract.add_to_blacklist(users.clone());
    match result {
        Ok(()) => {
            let users = contract.get_blacklisted_user();
            let result: HashSet<_> = users.iter().collect();
            let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
            let expected: HashSet<_> = expected_users.iter().collect();
            assert_eq!(expected, result)
        }
        Err(_) => todo!(),
    }

    let users = vec![chuck_btpaddress.clone()];
    let result = contract.remove_from_blacklist(users.clone());
    match result {
        Ok(()) => {
            let result = contract.get_blacklisted_user().contains(&chuck());

            assert_eq!(false, result)
        }
        Err(_) => todo!(),
    }
}

#[test]
fn remove_non_blacklisted_user_from_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let chuck_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", chuck().clone()).as_str()).unwrap();

    let charlie_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", charlie().clone()).as_str()).unwrap();

    let users = vec![chuck_btpaddress.clone(), charlie_btpaddress.clone()];

    let result = contract.add_to_blacklist(users.clone());
    match result {
        Ok(()) => {
            let users = contract.get_blacklisted_user();
            let result: HashSet<_> = users.iter().collect();
            let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
            let expected: HashSet<_> = expected_users.iter().collect();
            assert_eq!(expected, result)
        }
        Err(_) => todo!(),
    }
    let carol_btpaddress =
        BTPAddress::from_str(format!("btp://0x1.near/{}", carol().clone()).as_str()).unwrap();
    let users = vec![carol_btpaddress.clone()];
    let result = contract.remove_from_blacklist(users.clone());
    match result {
        Ok(()) => {}
        Err(err) => {
            assert_eq!(BshError::UserNotBlacklisted, err)
        }
    }
}
