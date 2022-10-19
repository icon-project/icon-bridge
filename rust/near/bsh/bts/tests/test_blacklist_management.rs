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
    messages::{BtpMessage, SerializedMessage},
    BTPAddress,
};
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
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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

    let users = vec![chuck(), charlie()];

    contract.add_to_blacklist(users);

    let is_user_blacklisted = contract.is_user_black_listed(charlie());

    assert_eq!(true, is_user_blacklisted)
}
