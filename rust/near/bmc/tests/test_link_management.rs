use bmc::BtpMessageCenter;
use near_sdk::{env, serde_json::json, testing_env, AccountId, VMContext, Gas, test_utils::VMContextBuilder};
use std::collections::HashSet;
pub mod accounts;
use accounts::*;
use libraries::types::{Address, BTPAddress, VerifierStatus};

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
fn add_link_new_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));

    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link);

    let result = contract.get_links();
    assert_eq!(
        result,
        json!(["btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b"])
    );
}

#[test]
#[should_panic(expected = "BMCRevertAlreadyExistsLink")]
fn add_link_existing_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    contract.add_link(link);
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn add_link_permission() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    testing_env!(context(chuck()));
    contract.add_link(link);
}

#[test]
fn remove_link_existing_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link_1 =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let link_2 =
        BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link_1.clone());

    contract.add_link(link_2);
    contract.remove_link(link_1);

    let result = contract.get_links();
    assert_eq!(
        result,
        json!(["btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b"])
    );
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsLink")]
fn remove_link_non_exisitng_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link_1 =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let link_2 =
        BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link_1.clone());

    contract.remove_link(link_2);
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn remove_link_permission() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());
    testing_env!(context(chuck()));
    contract.remove_link(link);
}

#[test]
fn set_link_existing_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());

    contract.set_link(link.clone(), 1500, 100, 5);
    let link_status = contract.get_status(link);

    assert_eq!(link_status.block_interval_dst(), 1500);
    assert_eq!(link_status.delay_limit(), 5);
    assert_eq!(link_status.max_aggregation(), 100);
}

#[test]
fn set_link_rx_height() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link.clone());

    contract.set_link_rx_height(link.clone(), 11001);
    let link_status = contract.get_status(link);

    assert_eq!(link_status.rx_height(), 11001);
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsLink")]
fn set_link_non_exisitng_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.set_link(link, 1, 1, 1);
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsPermission")]
fn set_link_permission() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());
    testing_env!(context(chuck()));
    contract.set_link(link.clone(), 1, 1, 1);
}

#[test]
#[should_panic(expected = "BMCRevertInvalidParam")]
fn set_link_invalid_param() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());
    contract.set_link(link, 0, 0, 0);
}

#[test]
fn get_links() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link_1 =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let link_2 =
        BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.add_link(link_1.clone());

    contract.add_link(link_2);

    let result = contract.get_links();
    assert_eq!(
        result,
        json!([
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b",
            "btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b"
        ])
    );
}

#[test]
fn get_status_exisitng_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());

    contract.set_link(link.clone(), 1, 10, 1);

    testing_env!(context(charlie()));
    let link_status = contract.get_status(link);

    assert_eq!(link_status.block_interval_dst(), 1);
    assert_eq!(link_status.delay_limit(), 1);
    assert_eq!(link_status.max_aggregation(), 10);
}

#[test]
#[should_panic(expected = "BMCRevertNotExistsLink")]
fn get_status_non_exisitng_link() {
    let context = |v: AccountId| (get_context(false, v));
    testing_env!(context(alice()));
    let contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    contract.get_status(link);
}
