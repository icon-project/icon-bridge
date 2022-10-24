use std::{
    collections::HashSet,
    convert::{TryFrom, TryInto},
    iter::FromIterator,
    result,
    str::FromStr,
};

use btp_common::errors::BshError;
use bts::{BtpTokenService, Coin};
use libraries::types::{
    messages::{BtpMessage, ErrorMessage, SerializedMessage},
    BTPAddress, WrappedI128,
};
use near_sdk::{
    env, serde_json::to_value, test_utils::test_env::alice, testing_env, AccountId, PromiseResult,
    VMContext, test_utils::VMContextBuilder, Gas
};
mod token;
use token::*;
pub mod accounts;
use accounts::*;

fn get_context(is_view: bool, signer_account_id: AccountId, attached_deposit: u128, account_balance: u128) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .storage_usage(env::storage_usage())
        .prepaid_gas(Gas(10u64.pow(18)))
        .attached_deposit(attached_deposit)
        .account_balance(account_balance)
        .build()
}

#[test]
fn add_user_to_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let users = vec![chuck(), charlie()];

    contract.add_to_blacklist(users);
    let users = contract.get_blacklisted_users();
    let result: HashSet<_> = users.iter().collect();
    let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
    let expected: HashSet<_> = expected_users.iter().collect();
    assert_eq!(expected, result)
}

#[test]
fn remove_blacklisted_user_from_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let users = vec![chuck().clone(), charlie().clone()];

    contract.add_to_blacklist(users.clone());
    let users = contract.get_blacklisted_users();
    let result: HashSet<_> = users.iter().collect();
    let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
    let expected: HashSet<_> = expected_users.iter().collect();
    assert_eq!(expected, result);

    let users = vec![chuck().clone()];
    let result = contract.remove_from_blacklist(users.clone());
    match result {
        Ok(()) => {
            let result = contract.get_blacklisted_users().contains(&chuck());

            assert_eq!(false, result)
        }
        Err(_) => todo!(),
    }
}

#[test]
fn remove_non_blacklisted_user_from_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let users = vec![chuck().clone(), charlie().clone()];

    contract.add_to_blacklist(users.clone());
    let users = contract.get_blacklisted_users();
    let result: HashSet<_> = users.iter().collect();
    let expected_users: Vec<AccountId> = vec![charlie(), chuck()];
    let expected: HashSet<_> = expected_users.iter().collect();
    assert_eq!(expected, result);

    let users = vec![carol().clone()];
    let result = contract.remove_from_blacklist(users.clone());
    match result {
        Ok(()) => {}
        Err(err) => {
            assert_eq!(
                BshError::NonBlacklistedUsers {
                    message: carol().to_string()
                },
                err
            )
        }
    }
}

#[test]

fn handle_btp_message_to_add_user_to_blacklist() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "bts".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let message: BtpMessage<SerializedMessage> = BtpMessage::try_from("-K-4OWJ0cDovLzB4Mi5pY29uL2N4NjE5M2U2OTI3NzIzZWNiMzJkYWNiMGExMjVhOTg2NjMzNzY4N2IwM7hPYnRwOi8vMHgxLm5lYXIvN2ZlN2VkMGY4YjI2MTdmYjRlMTA4NWY3YzQzYTM0OWFjZDNmMzMwMGVlYTZiODgxODc2NDZhNDU4ZWNhYzIwY4NidHMIndwDmtkAzo1hbGljZS50ZXN0bmV0iDB4MS5uZWFy".to_string()).unwrap();

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(message);

    let blacklisted_user = contract.get_blacklisted_users();

    assert_eq!(
        blacklisted_user,
        vec![AccountId::from_str("alice.testnet").unwrap()]
    )
}

#[test]
#[cfg(feature = "testable")]
fn handle_btp_message_to_change_token_limit() {
    use libraries::types::TokenLimit;

    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NEAR_NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "bts".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let message: BtpMessage<SerializedMessage> = BtpMessage::try_from("-L64OWJ0cDovLzB4Mi5pY29uL2N4NjE5M2U2OTI3NzIzZWNiMzJkYWNiMGExMjVhOTg2NjMzNzY4N2IwM7hPYnRwOi8vMHgxLm5lYXIvN2ZlN2VkMGY4YjI2MTdmYjRlMTA4NWY3YzQzYTM0OWFjZDNmMzMwMGVlYTZiODgxODc2NDZhNDU4ZWNhYzIwY4NidHMKrOsEqejSkWJ0cC0weDEubmVhci1ORUFSy4oCHhngybqyQAAAiDB4MS5uZWFy".to_string()).unwrap();

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(message);

    let token_limits = contract.get_token_limits().to_vec();

    assert_eq!(
        token_limits,
        vec![TokenLimit::new(
            "btp-0x1.near-NEAR".to_string(),
            10000000000000000000000
        )]
    )
}

#[test]
fn is_user_blacklisted() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let users = vec![chuck(), charlie()];

    contract.add_to_blacklist(users);

    let is_user_blacklisted = contract.is_user_black_listed(charlie());

    assert_eq!(true, is_user_blacklisted)
}

#[test]

fn handle_external_service_error_message() {
    use near_sdk::json_types::Base64VecU8;

    let message = "-P_4_bj7-PkBuPH47_jtuE9idHA6Ly8weDIubmVhci83MjcwYTc5YmU3ODlkNzcwZjJkZTAxNTA0NzY4NGUyODA2NTk3ZWVlZTk2ZWUzY2E4N2IxNzljNjM5OWRlYWFmNriZ-Je4OWJ0cDovLzB4Ny5pY29uL2N4MWFkNmZjYzQ2NWQxYjg2NDRjYTM3NWY5ZTEwYmFiZWVhNGMzODMxNbhPYnRwOi8vMHgyLm5lYXIvNzI3MGE3OWJlNzg5ZDc3MGYyZGUwMTUwNDc2ODRlMjgwNjU5N2VlZWU5NmVlM2NhODdiMTc5YzYzOTlkZWFhZoNidHOB3ITDKPgAhADNaJY=";
    let btp_message: BtpMessage<SerializedMessage> = BtpMessage::new(
        BTPAddress::new("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315".to_string()),
        BTPAddress::new(
            "btp://0x2.near/7270a79be789d770f2de015047684e2806597eeee96ee3ca87b179c6399deaaf"
                .to_string(),
        ),
        "bts".to_string(),
        WrappedI128::new(-36),
        vec![195, 40, 248, 0],
        None,
    );

    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, 0, deposit)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = <Coin>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "bts".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let link =
        BTPAddress::new("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315".to_string());
    let destination = BTPAddress::new(
        "btp://0x2.near/7270a79be789d770f2de015047684e2806597eeee96ee3ca87b179c6399deaaf"
            .to_string(),
    );
    testing_env!(context(bmc(), 0));
    contract.handle_btp_error(link.clone(), "bts".to_string(), -36, btp_message)
}
