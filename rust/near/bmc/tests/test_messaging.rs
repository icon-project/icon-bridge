use bmc::{BtpMessageCenter, RelayMessage};
use near_sdk::{
    base64,
    json_types::Base64VecU8,
    serde::Deserialize,
    serde_json::{self, from_value, json},
    testing_env, AccountId, VMContext, env,
};
use std::{collections::HashSet, convert::TryFrom};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::BmcServiceMessage, messages::BmcServiceType, messages::BtpMessage,
    messages::ErrorMessage, messages::SerializedBtpMessages, messages::SerializedMessage,
    messages::TokenServiceMessage, messages::TokenServiceType, Account, Address, BTPAddress,
    HashedCollection, WrappedI128,
};
use std::convert::TryInto;

fn get_context(
    input: Vec<u8>,
    is_view: bool,
    signer_account_id: AccountId,
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
        storage_usage: env::storage_usage(),
        attached_deposit: 0,
        prepaid_gas: 10u64.pow(18),
        random_seed: vec![0, 1, 2],
        is_view,
        output_data_receivers: vec![],
        epoch_height: 19,
    }
}

#[test]
fn decode() {
    let message: RelayMessage = RelayMessage::try_from("-QEE-QEBuP_4_QG49fjz-PG4T2J0cDovLzB4MS5uZWFyL2QwYWU4NGNkYzhmZTdmMTE3M2ZhYzQ5NTY2MjZiNzRlZTNiNGQzZTFhZjkwODRlMjRiNmQ1ZjYzOWU0NmExZTYBuJ34m7g5YnRwOi8vMHgyLmljb24vY3gxNzY2OTY1MWU2YjhmMGI2YTdiMDJhOTEyMzQwZWRlYTE5MTk4ZmQ4uE9idHA6Ly8weDEubmVhci9kMGFlODRjZGM4ZmU3ZjExNzNmYWM0OTU2NjI2Yjc0ZWUzYjRkM2UxYWY5MDg0ZTI0YjZkNWY2MzllNDZhMWU2g2JtYwCJyIRJbml0gsHAhACpFhA=".to_string()).unwrap();

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

    contract.handle_relay_message(source, message)
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
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.add_link(link.clone());
    let bmc_service_message = BmcServiceMessage::new(BmcServiceType::Init {
        links: vec![
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        ],
    });
    let btp_message = <BtpMessage<SerializedMessage>>::new(
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        ),
        "bmc".to_string(),
        WrappedI128::new(1),
        <Vec<u8>>::from(bmc_service_message.clone()),
        None,
    );

    contract.handle_btp_messages(&link.clone(), vec![btp_message]);
    let reachables = contract.get_reachable_link(link.clone());
    let mut expected = HashedCollection::new();
    expected.add(BTPAddress::new(
        "btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
    ));
    expected.add(BTPAddress::new(
        "btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
    ));
    assert_eq!(reachables, expected);
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
    assert_eq!(serialized_btp_messages,vec![btp_message_2])
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
    let btp_message: BtpMessage<SerializedMessage> = contract.get_message().unwrap().try_into().unwrap();
    let error_message = ErrorMessage::new(21, "BMCRevertUnreachable at btp://0x1.pra/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string());
    assert_eq!(
        btp_message,
        BtpMessage::new(
            BTPAddress::new(
                "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .to_string()
            ),
            BTPAddress::new(
                "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string()
            ),
            "bmc".to_string(),
            WrappedI128::new(10),
            error_message.clone().into(),
            None
        )
    );

    let btp_message:BtpMessage<SerializedMessage> = BtpMessage::new(
        BTPAddress::new(
            "btp://0x1.near/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string()
        ),
        BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()
        ),
        "bmc".to_string(),
        WrappedI128::new(-10),
        error_message.clone().into(),
        None
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
