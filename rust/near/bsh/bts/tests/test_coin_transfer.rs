#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_mut)]

use bts::BtpTokenService;
use near_sdk::{
    env, json_types::U128, test_utils::VMContextBuilder, testing_env, AccountId, Gas,
    PromiseResult, RuntimeFeesConfig, VMConfig, VMContext,
};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::{BtpMessage, TokenServiceMessage, TokenServiceType},
    Account, AccountBalance, AccumulatedAssetFees, Asset, BTPAddress, Math, WrappedI128,
    WrappedNativeCoin,
};
mod token;
use libraries::types::{Request, TransferableAsset};
use std::convert::TryInto;
use token::*;
pub type Token = Asset<WrappedNativeCoin>;

fn get_context(
    is_view: bool,
    signer_account_id: AccountId,
    attached_deposit: u128,
    account_balance: u128,
) -> VMContext {
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
fn deposit_native_coin() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 100));

    contract.deposit();

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    assert_eq!(result, U128::from(expected.deposit()))
}

#[test]
fn withdraw_native_coin() {
    let context =
        |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 1000);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let token_id = contract.token_id(nativecoin.name()).unwrap();
    testing_env!(context(chuck(), 1000));

    contract.deposit();

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1000).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));

    let storage_cost = contract
        .get_storage_balance(chuck(), nativecoin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.withdraw(nativecoin.name().to_string(), U128::from(999));

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_withdraw(
        chuck(),
        999,
        nativecoin.name().to_string(),
        token_id.clone(),
    );

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn withdraw_native_coin_higher_amount() {
    let context =
        |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 1000);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let token_id = contract.token_id(nativecoin.name()).unwrap();
    testing_env!(context(chuck(), 100));

    contract.deposit();

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));

    testing_env!(context(chuck(), 1));
    contract.withdraw(nativecoin.name().to_string(), U128::from(1000));

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let expected = AccountBalance::default();
    assert_eq!(result, U128::from(expected.deposit()));
}

#[test]
#[cfg(feature = "testable")]
fn external_transfer() {
    use btp_common::btp_address::Address;

    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    testing_env!(
        context(chuck(), 10000000000000000000000000),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );

    contract.deposit();

    let token_id = contract.token_id(nativecoin.name()).unwrap();

    let storage_cost = contract
        .get_storage_balance(chuck(), nativecoin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));

    contract.transfer(
        nativecoin.name().to_string(),
        destination.clone(),
        U128::from(9000000000000000000000000),
    );

    let message = TokenServiceMessage::new(TokenServiceType::RequestTokenTransfer {
        sender: chuck().to_string(),
        receiver: destination.account_id().to_string(),
        assets: vec![TransferableAsset::new(
            nativecoin.name().to_owned(),
            8099999999999999999999999,
            900000000000000000000001,
        )],
    });

    testing_env!(
        context(chuck(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );

    contract.send_service_message_callback(destination.network_address().unwrap(), message, 1);

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();

    expected
        .deposit_mut()
        .add(1000000000000000000000000)
        .unwrap();
    expected
        .locked_mut()
        .add(9000000000000000000000000)
        .unwrap();

    assert_eq!(result, Some(expected));

    let request = contract.last_request().unwrap();
    assert_eq!(
        request,
        Request::new(
            chuck().to_string(),
            destination.account_id().to_string(),
            vec![TransferableAsset::new(
                nativecoin.name().to_owned(),
                8099999999999999999999999,
                900000000000000000000001
            )]
        )
    )
}

#[test]
#[cfg(feature = "testable")]
fn handle_success_response_native_coin_external_transfer() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));

    contract.deposit();

    let token_id = contract.token_id(nativecoin.name()).unwrap();

    let storage_cost = contract
        .get_storage_balance(chuck(), nativecoin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));

    contract.transfer(
        nativecoin.name().to_string(),
        destination.clone(),
        U128::from(999),
    );

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(900).unwrap();
    expected.locked_mut().add(99).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 0,
                message: "Transfer Success".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(999));

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![AccumulatedAssetFees {
            name: nativecoin.name().to_string(),
            network: nativecoin.network().to_string(),
            accumulated_fees: 100
        }]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_success_response_icx_coin_external_transfer() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let icx_coin = Token::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(chuck(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let token_id = contract.token_id(icx_coin.name()).unwrap();

    contract.on_mint(
        900,
        token_id.clone(),
        chuck().clone(),
        Ok(U128::from(700000)),
    );

    let storage_cost = contract
        .get_storage_balance(chuck(), icx_coin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        icx_coin.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.locked_mut().add(800).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(contract.serial_no()),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 0,
                message: "Transfer Success".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(chuck(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_burn(720, token_id.clone());

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(80));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![
            AccumulatedAssetFees {
                name: nativecoin.name().to_string(),
                network: nativecoin.network().to_string(),
                accumulated_fees: 0
            },
            AccumulatedAssetFees {
                name: icx_coin.name().to_string(),
                network: icx_coin.network().to_string(),
                accumulated_fees: 81
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_failure_response_native_coin_external_transfer() {
    use libraries::types::AccumulatedAssetFees;

    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));
    let token_id = contract.token_id(nativecoin.name()).unwrap();

    contract.deposit();

    let storage_cost = contract
        .get_storage_balance(chuck(), nativecoin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        nativecoin.name().to_string(),
        destination.clone(),
        U128::from(999),
    );

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(900).unwrap();
    expected.locked_mut().add(99).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 1,
                message: "Transfer Failed".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(100));

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();
    expected.refundable_mut().add(899).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![AccumulatedAssetFees {
            name: nativecoin.name().to_string(),
            network: nativecoin.network().to_string(),
            accumulated_fees: 100
        }]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_failure_response_icx_coin_external_transfer() {
    use libraries::types::AccumulatedAssetFees;

    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let icx_coin = Token::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );
    testing_env!(
        context(bmc(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let token_id = contract.token_id(icx_coin.name()).unwrap();

    contract.on_mint(
        900,
        token_id.clone(),
        chuck().clone(),
        Ok(U128::from(700000)),
    );

    let storage_cost = contract
        .get_storage_balance(chuck(), icx_coin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));

    contract.transfer(
        icx_coin.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    testing_env!(context(chuck(), 0));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.locked_mut().add(800).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(contract.serial_no()),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 1,
                message: "Transfer Failed".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(81));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.refundable_mut().add(719).unwrap();

    assert_eq!(result, Some(expected));
    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![
            AccumulatedAssetFees {
                name: nativecoin.name().to_string(),
                network: nativecoin.network().to_string(),
                accumulated_fees: 0
            },
            AccumulatedAssetFees {
                name: icx_coin.name().to_string(),
                network: icx_coin.network().to_string(),
                accumulated_fees: 81
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn reclaim_icx_coin() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let icx_coin = Token::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(chuck(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let token_id = contract.token_id(icx_coin.name()).unwrap();
    contract.on_mint(900, token_id.clone(), chuck(), Ok(U128::from(700000)));
    let storage_cost = contract
        .get_storage_balance(chuck(), icx_coin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        icx_coin.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(contract.serial_no()),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 1,
                message: "Transfer Failed".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(context(chuck(), 0));
    contract.reclaim(icx_coin.name().to_string(), U128::from(700));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(800).unwrap();
    expected.refundable_mut().add(19).unwrap();

    assert_eq!(result, Some(expected));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer(nativecoin.name().to_string(), destination, U128::from(1001));
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken: ICON")]
fn external_transfer_unregistered_coin() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let icx_coin = Token::new(ICON_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer(icx_coin.name().to_string(), destination, U128::from(1001));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_nil_balance() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let icx_coin = Token::new(ICON_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer(icx_coin.name().to_string(), destination, U128::from(1001));
}

#[test]
#[cfg(feature = "testable")]
fn external_transfer_batch() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));
    let token_id = contract.token_id(nativecoin.name()).unwrap();

    contract.deposit();

    let token_id = contract.token_id(nativecoin.name()).unwrap();

    let storage_cost = contract
        .get_storage_balance(chuck(), nativecoin.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer_batch(
        vec![nativecoin.name().to_string()],
        destination,
        vec![U128::from(999)],
    );
    // TODO: Add other tokens
    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(900).unwrap();
    expected.locked_mut().add(99).unwrap();

    assert_eq!(result, Some(expected));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_batch_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer_batch(
        vec![nativecoin.name().to_string()],
        destination,
        vec![U128::from(1001)],
    );
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken")]
fn external_transfer_batch_unregistered_coin() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(context(alice(), 0));
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let icx_coin = Token::new(ICON_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer_batch(
        vec![nativecoin.name().to_string(), icx_coin.name().to_string()],
        destination,
        vec![U128::from(900), U128::from(1)],
    );
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_batch_nil_balance() {
    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit, 0);
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let icx_coin = Token::new(ICON_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    contract.register(icx_coin.clone());
    let token_id: [u8; 32] = env::sha256(icx_coin.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(icx_coin.clone(), token_id);
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer_batch(
        vec![nativecoin.name().to_string(), icx_coin.name().to_string()],
        destination,
        vec![U128::from(900), U128::from(1)],
    );
}
