use bmc::{BtpMessageCenter, RelayMessage};
use near_sdk::{
    base64, env,
    json_types::Base64VecU8,
    serde::Deserialize,
    serde_json::{self, from_value, json},
    testing_env, AccountId, VMContext,
    Gas, test_utils::VMContextBuilder
};
use std::{collections::HashSet, convert::TryFrom};
pub mod accounts;
use accounts::*;
use libraries::rlp::Encodable;
use libraries::types::{
    messages::BmcServiceMessage, messages::BmcServiceType, messages::BtpMessage,
    messages::ErrorMessage, messages::SerializedBtpMessages, messages::SerializedMessage,
    messages::TokenServiceMessage, messages::TokenServiceType, Account, Address, BTPAddress,
    HashedCollection, LinkStatus, WrappedI128,
};
use libraries::BytesMut;

use std::convert::TryInto;

fn get_context(input: Vec<u8>, is_view: bool, signer_account_id: AccountId) -> VMContext {
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
fn decode() {
    let message: RelayMessage = RelayMessage::try_from("-QEE-QEBuP_4_QG49fjz-PG4T2J0cDovLzB4MS5uZWFyLzQzMTBlOWI3YTQwMDMzMTkxM2EyYjE4NmNmNzMwODE3Njc4NmIzYTFhN2NkYzZhMzEzMjYxODAxY2NhMDliMGUBuJ34m7g5YnRwOi8vMHgyLmljb24vY3g0YzFhOWQ2MGRmMTE0MWEwODhhYTdiYTRhZDVhZDM3OGM3OTFjNmQxuE9idHA6Ly8weDEubmVhci80MzEwZTliN2E0MDAzMzE5MTNhMmIxODZjZjczMDgxNzY3ODZiM2ExYTdjZGM2YTMxMzI2MTgwMWNjYTA5YjBlg2JtYwCJyIRJbml0gsHAhACyR14=".to_string()).unwrap();

    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let source =
        BTPAddress::new("btp://0x2.icon/cxeecdbb073006a2b20232dd8f9e45078a281bc508".to_string());
    contract.add_link(source.clone());
    contract.add_relays(
        source.clone(),
        vec![
            "ed9436c4e4cfa045f38c0b579cc63d81d79617bc6b6c1e855ab18f9e11318f00"
                .parse::<AccountId>()
                .unwrap(),
            "verifier_2.near".parse::<AccountId>().unwrap(),
            charlie(),
        ],
    );
    testing_env!(context(charlie()));

    contract.handle_relay_message(source, message);

    let btp_message = <BtpMessage<TokenServiceMessage>>::new(
        BTPAddress::new("btp://0x38.bsc/0x034AaDE86BF402F023Aa17E5725fABC4ab9E9798".to_string()),
        BTPAddress::new("btp://0x1.icon/cx23a91ee3dd290486a9113a6a42429825d813de53".to_string()),
        "bts".to_string(),
        WrappedI128::new(21),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: "0x7A4341Af4995884546Bcf7e09eB98beD3eD26D28".to_owned(),
                receiver: "hx937517ac042d0a14f09d4677d302bb211184ac5f".to_owned(),
                assets: vec![],
            },
        )),
    );

    let serialized_message = <BtpMessage<SerializedMessage>>::try_from(&btp_message).unwrap();
}

#[test]
fn handle_serialized_btp_messages_service_message() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
}

#[test]
#[cfg(feature = "testable")]
fn handle_internal_service_message_init() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x2.icon/cx76ea658eb801a3f4aa37a19ad0a0676b5d6cecc9".to_string());
    contract.add_link(link.clone());
    let message: RelayMessage = RelayMessage::try_from("-QEE-QEBuP_4_QG49fjz-PG4T2J0cDovLzB4MS5uZWFyLzFkMGQwNjQ4NDYyMmY3MDYxYjAxNWY1ZWQwNjVkMjVjMGJmMDIzYmFjNzNiNGRkNDA3ODAxNzJmYjZlNDExOTQBuJ34m7g5YnRwOi8vMHgyLmljb24vY3g3NmVhNjU4ZWI4MDFhM2Y0YWEzN2ExOWFkMGEwNjc2YjVkNmNlY2M5uE9idHA6Ly8weDEubmVhci8xZDBkMDY0ODQ2MjJmNzA2MWIwMTVmNWVkMDY1ZDI1YzBiZjAyM2JhYzczYjRkZDQwNzgwMTcyZmI2ZTQxMTk0g2JtYwCJyIRJbml0gsHAhACuKaM=".to_string()).unwrap();

    contract.add_relays(link.clone(), vec![alice()]);

    contract.handle_relay_message(link.clone(), message);
    let reachables = contract.get_reachable_link(link.clone());
    let mut expected = HashedCollection::new();
    assert_eq!(reachables, expected);

    let link_status = contract.get_status(link);
    assert_eq!(link_status.rx_seq(), 1);
}

#[test]
#[cfg(feature = "testable")]
fn handle_internal_service_message_link() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());
    let bmc_service_message_1 = BmcServiceMessage::new(BmcServiceType::Init {
        links: vec![
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        ],
    });
    let btp_message_1 = <BtpMessage<SerializedMessage>>::new(
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        "bmc".to_string(),
        WrappedI128::new(1),
        <Vec<u8>>::from(bmc_service_message_1.clone()),
        None,
    );

    let bmc_service_message_2 = BmcServiceMessage::new(BmcServiceType::Link {
        link: BTPAddress::new(
            "btp://0x5.bsc/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        ),
    });
    let btp_message_2 = <BtpMessage<SerializedMessage>>::new(
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        "bmc".to_string(),
        WrappedI128::new(1),
        <Vec<u8>>::from(bmc_service_message_2.clone()),
        None,
    );

    contract.handle_btp_messages(&link.clone(), vec![btp_message_1, btp_message_2]);
    let reachables = contract.get_reachable_link(link.clone());
    let mut expected = HashedCollection::new();
    expected.add(BTPAddress::new(
        "btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
    ));
    expected.add(BTPAddress::new(
        "btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
    ));
    expected.add(BTPAddress::new(
        "btp://0x5.bsc/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
    ));
    assert_eq!(reachables, expected);
}

#[test]
#[cfg(feature = "testable")]
fn handle_internal_service_message_unlink() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());
    let bmc_service_message_1 = BmcServiceMessage::new(BmcServiceType::Init {
        links: vec![
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        ],
    });
    let btp_message_1 = <BtpMessage<SerializedMessage>>::new(
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        "bmc".to_string(),
        WrappedI128::new(1),
        <Vec<u8>>::from(bmc_service_message_1.clone()),
        None,
    );

    let bmc_service_message_2 = BmcServiceMessage::new(BmcServiceType::Unlink {
        link: BTPAddress::new(
            "btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        ),
    });
    let btp_message_2 = <BtpMessage<SerializedMessage>>::new(
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        "bmc".to_string(),
        WrappedI128::new(1),
        <Vec<u8>>::from(bmc_service_message_2.clone()),
        None,
    );

    contract.handle_btp_messages(&link.clone(), vec![btp_message_1, btp_message_2]);
    let reachables = contract.get_reachable_link(link.clone());
    let mut expected = HashedCollection::new();
    expected.add(BTPAddress::new(
        "btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
    ));
    assert_eq!(reachables, expected);
}

#[ignore]
#[test]
fn deserialize_serialized_btp_messages_from_json() {
    let btp_message = json!(["uNz42rg5YnRwOi8vMHgxLmljb24vY3g4N2VkOTA0OGI1OTRiOTUxOTlmMzI2ZmM3NmU3NmE5ZDMzZGQ2NjViuE9idHA6Ly8weDEubmVhci84OGJkMDU0NDI2ODZiZTBhNWRmN2RhMzNiNmYxMDg5ZWJmZWEzNzY5YjE5ZGJiMjQ3N2ZlMGNkNmUwZjEyNmU0g2JtYwG4R_hFhlVubGlua7g8-Dq4OGJ0cDovLzB4NS5wcmEvY3g4N2VkOTA0OGI1OTRiOTUxOTlmMzI2ZmM3NmU3NmE5ZDMzZGQ2NjVi"]);
    let serialized_btp_messages: SerializedBtpMessages = from_value(btp_message).unwrap();
    //let error_message = ErrorMessage::new(21, "BMCRevertUnreachable at btp://0x1.pra/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string());
    let bmc_service_message_2 = BmcServiceMessage::new(BmcServiceType::Unlink {
        link: BTPAddress::new(
            "btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        ),
    });
    let btp_message_2 = <BtpMessage<SerializedMessage>>::new(
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        "bmc".to_string(),
        WrappedI128::new(1),
        <Vec<u8>>::from(bmc_service_message_2.clone()),
        None,
    );
    assert_eq!(serialized_btp_messages, vec![btp_message_2])
    // TODO: Add;
}

#[ignore]
#[test]
#[cfg(feature = "testable")]
fn handle_route_message() {
    let btp_message = json!(["uQEB-P-4T2J0cDovLzB4MS5uZWFyLzg4YmQwNTQ0MjY4NmJlMGE1ZGY3ZGEzM2I2ZjEwODllYmZlYTM3NjliMTlkYmIyNDc3ZmUwY2Q2ZTBmMTI2ZTS4OWJ0cDovLzB4MS5pY29uL2N4ODdlZDkwNDhiNTk0Yjk1MTk5ZjMyNmZjNzZlNzZhOWQzM2RkNjY1YoNibWOB9rhr-GkVuGZCTUNSZXZlcnRVbnJlYWNoYWJsZSBhdCBidHA6Ly8weDEucHJhLzg4YmQwNTQ0MjY4NmJlMGE1ZGY3ZGEzM2I2ZjEwODllYmZlYTM3NjliMTlkYmIyNDc3ZmUwY2Q2ZTBmMTI2ZTQ"]);
    let serialized_btp_messages: SerializedBtpMessages = from_value(btp_message).unwrap();
    let context = |v: AccountId| (get_context(vec![], false, v));
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    let link =
        BTPAddress::new("btp://0x2.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());
    contract.handle_btp_messages(&link.clone(), serialized_btp_messages);
    let btp_message: BtpMessage<SerializedMessage> =
        contract.get_message().unwrap().try_into().unwrap();
    let error_message = ErrorMessage::new(21, "BMCRevertUnreachable at btp://0x1.pra/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string());
    assert_eq!(
        btp_message,
        BtpMessage::new(
            BTPAddress::new(
                "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .to_string()
            ),
            BTPAddress::new(
                "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .to_string()
            ),
            "bmc".to_string(),
            WrappedI128::new(10),
            error_message.clone().into(),
            None
        )
    );

    let btp_message: BtpMessage<SerializedMessage> = BtpMessage::new(
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        "bmc".to_string(),
        WrappedI128::new(-10),
        error_message.clone().into(),
        None,
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_external_service_message_existing_service() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd675b".to_string());
    let destination =
        BTPAddress::new("btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let btp_message = BtpMessage::new(
        link.clone(),
        destination.clone(),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: chuck().to_string(),
                receiver: destination.account_id().to_string(),
                assets: vec![],
            },
        )),
    );
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    contract.add_link(link.clone());

    contract.add_service(
        "nativecoin".to_string(),
        "nativecoin.near".parse::<AccountId>().unwrap(),
    );

    contract.handle_service_message_testable(
        link.clone(),
        <BtpMessage<SerializedMessage>>::try_from(&btp_message).unwrap(),
    );
}

// #[ignore]
#[test]
#[cfg(feature = "testable")]
fn handle_external_service_message_non_existing_service() {
    let context = |v: AccountId| (get_context(vec![], false, v));
    let link =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd675b".to_string());
    let destination =
        BTPAddress::new("btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let btp_message = BtpMessage::new(
        link.clone(),
        destination.clone(),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: chuck().to_string(),
                receiver: destination.account_id().to_string(),
                assets: vec![],
            },
        )),
    );
    testing_env!(context(alice()));
    let mut contract = BtpMessageCenter::new("0x1.near".into(), 1500);
    contract.add_link(link.clone());

    contract.handle_service_message_testable(
        link.clone(),
        <BtpMessage<SerializedMessage>>::try_from(&btp_message).unwrap(),
    );

    let btp_message: BtpMessage<ErrorMessage> = contract.get_message().unwrap().try_into().unwrap();
    let error_message = ErrorMessage::new(16, "BMCRevertNotExistBSH".to_string());
    assert_eq!(
        btp_message,
        BtpMessage::new(
            BTPAddress::new(
                "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .to_string()
            ),
            BTPAddress::new(
                "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd675b".to_string()
            ),
            "nativecoin".to_string(),
            WrappedI128::new(-1),
            error_message.clone().into(),
            Some(error_message)
        )
    );
}
