use bmc::BtpMessageCenter;
use near_sdk::{env, testing_env, AccountId, VMContext, Gas, test_utils::VMContextBuilder};
use std::collections::HashSet;
pub mod accounts;
use accounts::*;

fn get_context(input: Vec<u8>, is_view: bool, signer_account_id: AccountId) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .prepaid_gas(Gas(10u64.pow(18)))
        .build()
}

#[test]
fn add_owner_new_owner() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);

    contract.add_owner(carol());

    let owners = contract.get_owners();
    let result: HashSet<_> = owners.iter().collect();
    let expected_owners: Vec<AccountId> = vec![alice(), carol()];
    let expected: HashSet<_> = expected_owners.iter().collect();
    assert_eq!(result, expected);
}

#[test]
#[should_panic(expected = "BMCRevertAlreadyExistsOwner")]
fn add_owner_exisinting_owner() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);

    contract.add_owner(alice());
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn add_owner_permission() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    testing_env!(context(chuck()));
    contract.add_owner(carol());
}

#[test]
fn remove_owner_existing_owner() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);

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
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn remove_owner_permission() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);

    contract.add_owner(carol());

    testing_env!(context(chuck()));
    contract.add_owner(charlie());
}

#[test]
#[should_panic(expected = "BMCRevertLastOwner")]
fn remove_owner_last_owner() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);

    contract.remove_owner(alice());
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsOwner")]
fn remove_owner_non_exisitng_owner() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);

    contract.remove_owner(carol());
}
