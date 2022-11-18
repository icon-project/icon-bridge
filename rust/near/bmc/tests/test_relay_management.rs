#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_mut)]

use bmc::BtpMessageCenter;
use near_sdk::{
    env, serde_json::json, test_utils::VMContextBuilder, testing_env, AccountId, Gas, VMContext,
};
pub mod accounts;
use accounts::*;
use libraries::types::BTPAddress;

fn get_context(
    input: Vec<u8>,
    is_view: bool,
    signer_account_id: AccountId,
    storage_usage: u64,
    block_index: u64,
) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .storage_usage(storage_usage)
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .prepaid_gas(Gas(10u64.pow(18)))
        .block_index(block_index)
        .build()
}

#[test]
fn add_relay_new_relay() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );

    let relays = contract.get_relays(link);
    assert_eq!(relays, json!(["verifier_1.near", "verifier_2.near"]));
}

#[test]
fn add_relays_new_relay() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
    contract.add_relay(
        link.clone(),
        "verifier_3.near".parse::<AccountId>().unwrap(),
    );
    let relays = contract.get_relays(link);
    assert_eq!(
        relays,
        json!(["verifier_1.near", "verifier_2.near", "verifier_3.near"])
    );
}

#[test]
#[should_panic(expected = "BMCRevertRelayExist")]
fn add_relay_existing_relay() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
    contract.add_relay(
        link.clone(),
        "verifier_2.near".parse::<AccountId>().unwrap(),
    );
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsLink")]
fn add_relay_non_existing_link() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
}

#[test]
fn get_relays() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
    contract.add_relay(
        link.clone(),
        "verifier_3.near".parse::<AccountId>().unwrap(),
    );
    testing_env!(context(bob()));
    let relays = contract.get_relays(link);
    assert_eq!(
        relays,
        json!(["verifier_1.near", "verifier_2.near", "verifier_3.near"])
    );
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn add_relays_permission() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    testing_env!(context(chuck()));
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn add_relay_permission() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
    testing_env!(context(chuck()));
    contract.add_relay(
        link.clone(),
        "verifier_3.near".parse::<AccountId>().unwrap(),
    );
}

#[test]
fn remove_relay_existing_relay() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
    contract.remove_relay(
        link.clone(),
        "verifier_1.near".parse::<AccountId>().unwrap(),
    );
    let relays = contract.get_relays(link);
    assert_eq!(relays, json!(["verifier_2.near"]));
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsLink")]
fn remove_relay_non_existing_link() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.remove_relay(
        link.clone(),
        "verifier_3.near".parse::<AccountId>().unwrap(),
    );
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn remove_relay_permission() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );

    testing_env!(context(chuck()));
    contract.remove_relay(
        link.clone(),
        "verifier_1.near".parse::<AccountId>().unwrap(),
    );
}

#[test]
#[should_panic(expected = "BMCRevertNotExistRelay")]
fn remove_relay_non_existing_relay() {
    let context = |v: AccountId| (get_context(vec![], false, v, 0, 0));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
        ],
    );
    contract.remove_relay(
        link.clone(),
        "verifier_3.near".parse::<AccountId>().unwrap(),
    );
}

#[test]
fn rotate_relay() {
    let context = |v: AccountId, block_index: u64| {
        get_context(vec![], false, v, env::storage_usage(), block_index)
    };
    testing_env!(context(alice(), 0));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.set_link(link.clone(), 2000, 50, 1);
    contract.add_relays(
        link.clone(),
        vec![
            "verifier_1.near".parse::<AccountId>().unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
            "verifier_3.near".parse::<AccountId>().unwrap(),
            "verifier_4.near".parse::<AccountId>().unwrap(),
            "verifier_5.near".parse::<AccountId>().unwrap(),
            "verifier_6.near".parse::<AccountId>().unwrap(),
            "verifier_7.near".parse::<AccountId>().unwrap(),
        ],
    );

    let mut link_property = contract.get_link(link.clone());

    testing_env!(context(alice(), 51));

    let mut link = contract.get_link(link);
    let account_id = link.rotate_relay(1200, true);
    let account_id = link.rotate_relay(3200, true);
}
